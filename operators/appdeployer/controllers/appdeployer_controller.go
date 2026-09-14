package controllers

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	platformv1alpha1 "github.com/Open-Industries/open-ecosystem-services/operators/appdeployer/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Embedded templates keep generated Tekton and Argo CD resources versioned with the controller.
//
//go:embed templates/tekton-pipeline.yaml.gotmpl
var pipelineTemplate string

//go:embed templates/tekton-pipelinerun.yaml.gotmpl
var pipelineRunTemplate string

//go:embed templates/argocd-application.yaml.gotmpl
var applicationTemplate string

//go:embed templates/image-validation-job.yaml.gotmpl
var imageValidationJobTemplate string

//go:embed templates/gitops-update-job.yaml.gotmpl
var gitOpsUpdateJobTemplate string

const requeueInterval = 20 * time.Second

// AppDeployerReconciler reconciles AppDeployer resources.
type AppDeployerReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

type templateData struct {
	Name, Namespace, ApplicationName, RunName, Generation                                                           string
	ImageRef, SourceRevision, SourceContextDir                                                                      string
	TestImage, UnitTestCommand, IntegrationTestCommand                                                              string
	RegistrySecret, GitOpsSecret, GitSourceSecret, ArtifactCredentialsSecret, MavenSettingsSecret, SonarTokenSecret string
	ArgoNamespace, ArgoProject, DestinationServer, GitOpsRevision                                                   string
	ImageCheckJobName, GitOpsJobName                                                                                string
	AfterClone, AfterUnitTests, AfterQuality, AfterPush                                                             string
	RegistryTLSVerify                                                                                               string
	SonarEnabled                                                                                                    bool
	ServerlessEnabled                                                                                               bool
	ServerlessServiceName, ServerlessManifestPath, ServerlessContainerName, ServerlessMinScale                      string
	ServerlessMaxScale, ServerlessScaleTarget, ServerlessMetric, ServerlessContainerConcurrency                     string
	ServerlessTimeoutSeconds, ServerlessVisibility                                                                  string
	SonarTimeout                                                                                                    int32
	Spec                                                                                                            platformv1alpha1.AppDeployerSpec
}

// +kubebuilder:rbac:groups=platform.openindustries.org,resources=appdeployers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=platform.openindustries.org,resources=appdeployers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=platform.openindustries.org,resources=appdeployers/finalizers,verbs=update
// +kubebuilder:rbac:groups=tekton.dev,resources=pipelines;pipelineruns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=argoproj.io,resources=applications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;delete
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *AppDeployerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var app platformv1alpha1.AppDeployer
	if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := validateSpec(&app); err != nil {
		r.setCondition(&app, "Ready", metav1.ConditionFalse, "InvalidSpec", err.Error())
		app.Status.Phase = "Invalid"
		return ctrl.Result{}, r.updateStatus(ctx, &app)
	}
	if err := r.verifyReferencedSecrets(ctx, &app); err != nil {
		return r.fail(ctx, &app, "SecretValidationFailed", err)
	}

	data := buildTemplateData(&app)
	if app.Status.ObservedGeneration == app.Generation && conditionTrue(app.Status.Conditions, "Ready", app.Generation) {
		if err := r.ensureArgoApplication(ctx, &app, data); err != nil {
			return r.fail(ctx, &app, "ArgoApplicationError", err)
		}
		app.Status.ArgoApplicationName = data.ApplicationName
		return ctrl.Result{}, r.updateStatus(ctx, &app)
	}

	var result ctrl.Result
	var err error
	switch app.Spec.Env {
	case platformv1alpha1.EnvironmentDevelopment:
		if err = r.ensureArgoApplication(ctx, &app, data); err != nil {
			return r.fail(ctx, &app, "ArgoApplicationError", err)
		}
		app.Status.ArgoApplicationName = data.ApplicationName
		result, err = r.reconcileDevelopment(ctx, &app, data)
	case platformv1alpha1.EnvironmentProduction:
		result, err = r.reconcileProduction(ctx, &app, data)
	default:
		err = fmt.Errorf("unsupported environment %q", app.Spec.Env)
	}
	if err != nil {
		logger.Error(err, "reconciliation failed", "environment", app.Spec.Env)
		return r.fail(ctx, &app, "ReconcileError", err)
	}

	app.Status.ObservedGeneration = app.Generation
	app.Status.ImageRef = data.ImageRef
	if statusErr := r.updateStatus(ctx, &app); statusErr != nil {
		return ctrl.Result{}, statusErr
	}
	return result, nil
}

func (r *AppDeployerReconciler) reconcileDevelopment(ctx context.Context, app *platformv1alpha1.AppDeployer, data templateData) (ctrl.Result, error) {
	if err := r.ensurePipelineServiceAccount(ctx, app); err != nil {
		return ctrl.Result{}, err
	}
	pipeline, err := renderUnstructured("pipeline", pipelineTemplate, data)
	if err != nil {
		return ctrl.Result{}, err
	}
	if err = r.upsertUnstructured(ctx, app, pipeline, true); err != nil {
		return ctrl.Result{}, err
	}

	run, err := renderUnstructured("pipelinerun", pipelineRunTemplate, data)
	if err != nil {
		return ctrl.Result{}, err
	}
	created, err := r.createUnstructuredIfAbsent(ctx, app, run, true)
	if err != nil {
		return ctrl.Result{}, err
	}
	if created && r.Recorder != nil {
		r.Recorder.Eventf(app, "Normal", "PipelineStarted", "Created PipelineRun %s", data.RunName)
	}
	app.Status.PipelineRunName = data.RunName

	current := &unstructured.Unstructured{}
	current.SetAPIVersion("tekton.dev/v1")
	current.SetKind("PipelineRun")
	if err = r.Get(ctx, types.NamespacedName{Name: data.RunName, Namespace: app.Namespace}, current); err != nil {
		return ctrl.Result{}, err
	}
	done, succeeded, reason, message := unstructuredCompletion(current)
	if !done {
		app.Status.Phase = "Building"
		r.setCondition(app, "Ready", metav1.ConditionFalse, "PipelineRunning", "Tekton PipelineRun is still running")
		return ctrl.Result{RequeueAfter: requeueInterval}, nil
	}
	if !succeeded {
		app.Status.Phase = "Failed"
		r.setCondition(app, "Ready", metav1.ConditionFalse, reason, message)
		return ctrl.Result{}, nil
	}
	app.Status.Phase = "Ready"
	r.setCondition(app, "Ready", metav1.ConditionTrue, "DeliverySucceeded", "CI pipeline and GitOps update completed")
	return ctrl.Result{}, nil
}

func (r *AppDeployerReconciler) ensurePipelineServiceAccount(ctx context.Context, owner *platformv1alpha1.AppDeployer) error {
	sa := &corev1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: "appdeployer-pipeline", Namespace: owner.Namespace}}
	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, sa, func() error {
		if sa.Labels == nil {
			sa.Labels = map[string]string{}
		}
		sa.Labels["app.kubernetes.io/managed-by"] = "appdeployer-operator"
		return controllerutil.SetControllerReference(owner, sa, r.Scheme)
	})
	return err
}

func (r *AppDeployerReconciler) verifyReferencedSecrets(ctx context.Context, app *platformv1alpha1.AppDeployer) error {
	type secretCheck struct {
		ref     *platformv1alpha1.SecretReference
		keys    []string
		purpose string
	}
	checks := []secretCheck{
		{app.Spec.Registry.CredentialsSecretRef, []string{corev1.DockerConfigJsonKey}, "registry"},
		{app.Spec.GitOps.CredentialsSecretRef, []string{"username", "password"}, "GitOps"},
	}
	if app.Spec.Source != nil {
		checks = append(checks, secretCheck{app.Spec.Source.CredentialsSecretRef, []string{"username", "password"}, "source Git"})
	}
	if app.Spec.ArtifactRepository != nil {
		checks = append(checks,
			secretCheck{app.Spec.ArtifactRepository.CredentialsSecretRef, nil, "artifact repository"},
			secretCheck{app.Spec.ArtifactRepository.MavenSettingsSecretRef, []string{"settings.xml"}, "Maven settings"},
		)
	}
	if app.Spec.SonarQube != nil && app.Spec.SonarQube.Enabled {
		checks = append(checks, secretCheck{app.Spec.SonarQube.TokenSecretRef, []string{"token"}, "SonarQube"})
	}
	for _, check := range checks {
		if check.ref == nil {
			continue
		}
		var secret corev1.Secret
		if err := r.Get(ctx, types.NamespacedName{Name: check.ref.Name, Namespace: app.Namespace}, &secret); err != nil {
			return fmt.Errorf("%s Secret %q: %w", check.purpose, check.ref.Name, err)
		}
		for _, key := range check.keys {
			if _, found := secret.Data[key]; !found {
				return fmt.Errorf("%s Secret %q must contain key %q", check.purpose, check.ref.Name, key)
			}
		}
	}
	return nil
}

func (r *AppDeployerReconciler) reconcileProduction(ctx context.Context, app *platformv1alpha1.AppDeployer, data templateData) (ctrl.Result, error) {
	check, err := renderJob("image-validation", imageValidationJobTemplate, data)
	if err != nil {
		return ctrl.Result{}, err
	}
	if _, err = r.createJobIfAbsent(ctx, app, check); err != nil {
		return ctrl.Result{}, err
	}
	done, succeeded, message, err := r.jobCompletion(ctx, check.Namespace, check.Name)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !done {
		app.Status.Phase = "ValidatingImage"
		r.setCondition(app, "ImageAvailable", metav1.ConditionUnknown, "ValidationRunning", "OCI registry validation is running")
		return ctrl.Result{RequeueAfter: requeueInterval}, nil
	}
	if !succeeded {
		app.Status.Phase = "Failed"
		r.setCondition(app, "ImageAvailable", metav1.ConditionFalse, "ImageValidationFailed", message)
		r.setCondition(app, "Ready", metav1.ConditionFalse, "ImageValidationFailed", message)
		return ctrl.Result{}, nil
	}
	r.setCondition(app, "ImageAvailable", metav1.ConditionTrue, "ImageValidated", "Production image exists in the registry")

	job, err := renderJob("gitops-update", gitOpsUpdateJobTemplate, data)
	if err != nil {
		return ctrl.Result{}, err
	}
	if _, err = r.createJobIfAbsent(ctx, app, job); err != nil {
		return ctrl.Result{}, err
	}
	app.Status.GitOpsJobName = job.Name
	done, succeeded, message, err = r.jobCompletion(ctx, job.Namespace, job.Name)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !done {
		app.Status.Phase = "Promoting"
		r.setCondition(app, "Ready", metav1.ConditionFalse, "GitOpsUpdateRunning", "Production GitOps update is running")
		return ctrl.Result{RequeueAfter: requeueInterval}, nil
	}
	if !succeeded {
		app.Status.Phase = "Failed"
		r.setCondition(app, "Ready", metav1.ConditionFalse, "GitOpsUpdateFailed", message)
		return ctrl.Result{}, nil
	}
	if err = r.ensureArgoApplication(ctx, app, data); err != nil {
		return ctrl.Result{}, fmt.Errorf("ensure Argo CD Application after promotion: %w", err)
	}
	app.Status.ArgoApplicationName = data.ApplicationName
	app.Status.Phase = "Ready"
	r.setCondition(app, "Ready", metav1.ConditionTrue, "PromotionSucceeded", "Image validated, GitOps updated, and Argo CD Application ensured")
	return ctrl.Result{}, nil
}

func (r *AppDeployerReconciler) ensureArgoApplication(ctx context.Context, owner *platformv1alpha1.AppDeployer, data templateData) error {
	application, err := renderUnstructured("application", applicationTemplate, data)
	if err != nil {
		return err
	}
	// Argo CD commonly runs in a different namespace, where Kubernetes forbids ownerReferences.
	return r.upsertUnstructured(ctx, owner, application, application.GetNamespace() == owner.Namespace)
}

func (r *AppDeployerReconciler) upsertUnstructured(ctx context.Context, owner *platformv1alpha1.AppDeployer, desired *unstructured.Unstructured, setOwner bool) error {
	if setOwner {
		if err := controllerutil.SetControllerReference(owner, desired, r.Scheme); err != nil {
			return err
		}
	}
	current := &unstructured.Unstructured{}
	current.SetAPIVersion(desired.GetAPIVersion())
	current.SetKind(desired.GetKind())
	key := types.NamespacedName{Name: desired.GetName(), Namespace: desired.GetNamespace()}
	err := r.Get(ctx, key, current)
	if apierrors.IsNotFound(err) {
		return r.Create(ctx, desired)
	}
	if err != nil {
		return err
	}
	desired.SetResourceVersion(current.GetResourceVersion())
	desired.SetUID(current.GetUID())
	return r.Update(ctx, desired)
}

func (r *AppDeployerReconciler) createUnstructuredIfAbsent(ctx context.Context, owner *platformv1alpha1.AppDeployer, desired *unstructured.Unstructured, setOwner bool) (bool, error) {
	if setOwner {
		if err := controllerutil.SetControllerReference(owner, desired, r.Scheme); err != nil {
			return false, err
		}
	}
	current := &unstructured.Unstructured{}
	current.SetAPIVersion(desired.GetAPIVersion())
	current.SetKind(desired.GetKind())
	err := r.Get(ctx, client.ObjectKeyFromObject(desired), current)
	if err == nil {
		return false, nil
	}
	if !apierrors.IsNotFound(err) {
		return false, err
	}
	return true, r.Create(ctx, desired)
}

func (r *AppDeployerReconciler) createJobIfAbsent(ctx context.Context, owner *platformv1alpha1.AppDeployer, desired *batchv1.Job) (bool, error) {
	if err := controllerutil.SetControllerReference(owner, desired, r.Scheme); err != nil {
		return false, err
	}
	current := &batchv1.Job{}
	err := r.Get(ctx, client.ObjectKeyFromObject(desired), current)
	if err == nil {
		return false, nil
	}
	if !apierrors.IsNotFound(err) {
		return false, err
	}
	return true, r.Create(ctx, desired)
}

func (r *AppDeployerReconciler) jobCompletion(ctx context.Context, namespace, name string) (bool, bool, string, error) {
	var job batchv1.Job
	if err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, &job); err != nil {
		return false, false, "", err
	}
	for _, condition := range job.Status.Conditions {
		if condition.Status != "True" {
			continue
		}
		switch condition.Type {
		case batchv1.JobComplete:
			return true, true, condition.Message, nil
		case batchv1.JobFailed:
			return true, false, condition.Message, nil
		}
	}
	return false, false, "", nil
}

func unstructuredCompletion(object *unstructured.Unstructured) (bool, bool, string, string) {
	conditions, found, _ := unstructured.NestedSlice(object.Object, "status", "conditions")
	if !found {
		return false, false, "PipelineRunning", "PipelineRun has not reported a terminal condition"
	}
	for _, item := range conditions {
		condition, ok := item.(map[string]any)
		if !ok || condition["type"] != "Succeeded" {
			continue
		}
		status, _ := condition["status"].(string)
		reason, _ := condition["reason"].(string)
		message, _ := condition["message"].(string)
		if status == "True" {
			return true, true, reason, message
		}
		if status == "False" {
			return true, false, reason, message
		}
	}
	return false, false, "PipelineRunning", "PipelineRun is still running"
}

func buildTemplateData(app *platformv1alpha1.AppDeployer) templateData {
	spec := app.Spec
	data := templateData{
		Name: app.Name, Namespace: app.Namespace, ApplicationName: dnsName(app.Namespace + "-" + app.Name), Generation: strconv.FormatInt(app.Generation, 10), Spec: spec,
		ImageRef:          spec.Registry.Image + ":" + spec.Registry.Tag,
		RunName:           generationName(app.Name+"-run", app.Generation),
		ImageCheckJobName: generationName(app.Name+"-image-check", app.Generation),
		GitOpsJobName:     generationName(app.Name+"-gitops", app.Generation),
		SourceRevision:    "main", SourceContextDir: ".", TestImage: "registry.access.redhat.com/ubi9/openjdk-21:latest",
		UnitTestCommand: "./mvnw test", IntegrationTestCommand: "./mvnw verify -Pintegration",
		ArgoNamespace: valueOr(spec.GitOps.ArgoCDNamespace, "openshift-gitops"),
		ArgoProject:   valueOr(spec.GitOps.Project, "default"), DestinationServer: valueOr(spec.GitOps.DestinationServer, "https://kubernetes.default.svc"),
		GitOpsRevision: valueOr(spec.GitOps.Revision, "main"), RegistryTLSVerify: strconv.FormatBool(!spec.Registry.InsecureSkipTLSVerify),
		AfterClone: "clone-source", AfterUnitTests: "clone-source", AfterQuality: "clone-source", AfterPush: "push-image", SonarTimeout: 300,
	}
	if spec.Source != nil {
		data.SourceRevision = valueOr(spec.Source.Revision, "main")
		data.SourceContextDir = valueOr(spec.Source.ContextDir, ".")
	}
	if spec.Source != nil && spec.Source.CredentialsSecretRef != nil {
		data.GitSourceSecret = spec.Source.CredentialsSecretRef.Name
	}
	if spec.Tests.TestImage != "" {
		data.TestImage = spec.Tests.TestImage
	}
	if spec.Tests.UnitTestCommand != "" {
		data.UnitTestCommand = spec.Tests.UnitTestCommand
	}
	if spec.Tests.IntegrationTestCommand != "" {
		data.IntegrationTestCommand = spec.Tests.IntegrationTestCommand
	}
	if spec.ArtifactRepository != nil {
		data.AfterClone = "configure-dependencies"
		if spec.ArtifactRepository.CredentialsSecretRef != nil {
			data.ArtifactCredentialsSecret = spec.ArtifactRepository.CredentialsSecretRef.Name
		}
		if spec.ArtifactRepository.MavenSettingsSecretRef != nil {
			data.MavenSettingsSecret = spec.ArtifactRepository.MavenSettingsSecretRef.Name
		}
	}
	data.AfterUnitTests = data.AfterClone
	if spec.Tests.UnitTests {
		data.AfterUnitTests = "unit-tests"
	}
	data.AfterQuality = data.AfterUnitTests
	if spec.SonarQube != nil && spec.SonarQube.Enabled {
		data.SonarEnabled = true
		data.AfterQuality = "sonarqube-quality-gate"
		data.SonarTokenSecret = spec.SonarQube.TokenSecretRef.Name
		if spec.SonarQube.QualityGateTimeoutSeconds > 0 {
			data.SonarTimeout = spec.SonarQube.QualityGateTimeoutSeconds
		}
	}
	if spec.Tests.IntegrationTests {
		data.AfterPush = "integration-tests"
	}
	if spec.Registry.CredentialsSecretRef != nil {
		data.RegistrySecret = spec.Registry.CredentialsSecretRef.Name
	}
	if spec.GitOps.CredentialsSecretRef != nil {
		data.GitOpsSecret = spec.GitOps.CredentialsSecretRef.Name
	}
	if spec.Deployment.Type == platformv1alpha1.DeploymentTypeServerless && spec.Deployment.Serverless != nil {
		serverless := spec.Deployment.Serverless
		data.ServerlessEnabled = true
		data.ServerlessServiceName = serverless.ServiceName
		data.ServerlessManifestPath = valueOr(serverless.ManifestPath, "service.yaml")
		data.ServerlessContainerName = valueOr(serverless.ContainerName, "user-container")
		data.ServerlessMetric = valueOr(serverless.Metric, "concurrency")
		data.ServerlessVisibility = valueOr(serverless.Visibility, "external")
		data.ServerlessMinScale = int32PointerString(serverless.MinScale)
		data.ServerlessMaxScale = int32PointerString(serverless.MaxScale)
		data.ServerlessScaleTarget = int32PointerString(serverless.ScaleTarget)
		data.ServerlessContainerConcurrency = int64PointerString(serverless.ContainerConcurrency)
		data.ServerlessTimeoutSeconds = int64PointerString(serverless.TimeoutSeconds)
	}
	return data
}

func validateSpec(app *platformv1alpha1.AppDeployer) error {
	s := app.Spec
	if s.Env != platformv1alpha1.EnvironmentDevelopment && s.Env != platformv1alpha1.EnvironmentProduction {
		return errors.New("spec.env must be development or production")
	}
	if s.Registry.Image == "" || s.Registry.Tag == "" {
		return errors.New("spec.registry.image and spec.registry.tag are required")
	}
	if !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/@-]*$`).MatchString(s.Registry.Image) || !regexp.MustCompile(`^[A-Za-z0-9_][A-Za-z0-9_.-]{0,127}$`).MatchString(s.Registry.Tag) {
		return errors.New("spec.registry contains an invalid OCI image name or tag")
	}
	if s.Registry.CredentialsSecretRef == nil {
		return errors.New("spec.registry.credentialsSecretRef is required")
	}
	if s.GitOps.RepositoryURL == "" || s.GitOps.Path == "" || s.GitOps.KustomizationImage == "" || s.GitOps.DestinationNamespace == "" {
		return errors.New("repositoryURL, path, kustomizationImage, and destinationNamespace are required in spec.gitOps")
	}
	if s.GitOps.CredentialsSecretRef == nil {
		return errors.New("spec.gitOps.credentialsSecretRef is required for automatic commits")
	}
	if !strings.HasPrefix(s.GitOps.RepositoryURL, "https://") {
		return errors.New("spec.gitOps.repositoryURL must use HTTPS")
	}
	if !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/@-]*$`).MatchString(s.GitOps.KustomizationImage) {
		return errors.New("spec.gitOps.kustomizationImage is invalid")
	}
	if strings.HasPrefix(s.GitOps.Path, "/") || strings.Contains(s.GitOps.Path, "..") || strings.ContainsAny(s.GitOps.Path, "\"'`$;|&\n\r") {
		return errors.New("spec.gitOps.path must be a safe relative path")
	}
	argoNamespace := valueOr(s.GitOps.ArgoCDNamespace, "openshift-gitops")
	if errs := validation.IsDNS1123Subdomain(argoNamespace); len(errs) > 0 {
		return fmt.Errorf("invalid spec.gitOps.argoCDNamespace: %s", strings.Join(errs, ", "))
	}
	if s.Env == platformv1alpha1.EnvironmentDevelopment {
		if s.Source == nil || s.Source.URL == "" || s.BuilderImage == "" {
			return errors.New("spec.source.url and spec.builderImage are required for development")
		}
		if !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/@-]*$`).MatchString(s.BuilderImage) {
			return errors.New("spec.builderImage is invalid")
		}
		if strings.HasPrefix(s.Source.ContextDir, "/") || strings.Contains(s.Source.ContextDir, "..") || strings.ContainsAny(s.Source.ContextDir, "\"'`$;|&\n\r") {
			return errors.New("spec.source.contextDir must be a safe relative path")
		}
		if s.Source.CredentialsSecretRef != nil && !strings.HasPrefix(s.Source.URL, "https://") {
			return errors.New("authenticated source repositories must use HTTPS")
		}
		if s.SonarQube != nil && s.SonarQube.Enabled && (s.SonarQube.URL == "" || s.SonarQube.ProjectKey == "" || s.SonarQube.TokenSecretRef == nil) {
			return errors.New("enabled SonarQube requires url, projectKey, and tokenSecretRef")
		}
	}
	if s.Deployment.Type != "" && s.Deployment.Type != platformv1alpha1.DeploymentTypeStandard && s.Deployment.Type != platformv1alpha1.DeploymentTypeServerless {
		return errors.New("spec.deployment.type must be standard or serverless")
	}
	if s.Deployment.Type == platformv1alpha1.DeploymentTypeServerless {
		serverless := s.Deployment.Serverless
		if serverless == nil || serverless.ServiceName == "" {
			return errors.New("spec.deployment.serverless.serviceName is required for serverless deployment")
		}
		if errs := validation.IsDNS1123Subdomain(serverless.ServiceName); len(errs) > 0 {
			return fmt.Errorf("invalid serverless serviceName: %s", strings.Join(errs, ", "))
		}
		manifestPath := valueOr(serverless.ManifestPath, "service.yaml")
		if strings.HasPrefix(manifestPath, "/") || strings.Contains(manifestPath, "..") || strings.ContainsAny(manifestPath, "\"'`$;|&\n\r") {
			return errors.New("spec.deployment.serverless.manifestPath must be a safe relative path")
		}
		containerName := valueOr(serverless.ContainerName, "user-container")
		if errs := validation.IsDNS1123Label(containerName); len(errs) > 0 {
			return fmt.Errorf("invalid serverless containerName: %s", strings.Join(errs, ", "))
		}
		if serverless.MinScale != nil && serverless.MaxScale != nil && *serverless.MinScale > *serverless.MaxScale {
			return errors.New("spec.deployment.serverless.minScale cannot exceed maxScale")
		}
	}
	return nil
}

func renderUnstructured(name, body string, data templateData) (*unstructured.Unstructured, error) {
	raw, err := executeTemplate(name, body, data)
	if err != nil {
		return nil, err
	}
	jsonData, err := yaml.ToJSON(raw)
	if err != nil {
		return nil, fmt.Errorf("decode %s template: %w", name, err)
	}
	var object map[string]any
	if err = json.Unmarshal(jsonData, &object); err != nil {
		return nil, fmt.Errorf("unmarshal %s: %w", name, err)
	}
	return &unstructured.Unstructured{Object: object}, nil
}

func renderJob(name, body string, data templateData) (*batchv1.Job, error) {
	raw, err := executeTemplate(name, body, data)
	if err != nil {
		return nil, err
	}
	jsonData, err := yaml.ToJSON(raw)
	if err != nil {
		return nil, err
	}
	var job batchv1.Job
	if err = json.Unmarshal(jsonData, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

func executeTemplate(name, body string, data templateData) ([]byte, error) {
	functions := template.FuncMap{
		"q":      strconv.Quote,
		"shellq": func(value string) string { return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'" },
	}
	tmpl, err := template.New(name).Funcs(functions).Option("missingkey=error").Parse(body)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	if err = tmpl.Execute(&buffer, data); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func generationName(prefix string, generation int64) string {
	name := fmt.Sprintf("%s-g%d", prefix, generation)
	if len(name) <= 63 {
		return name
	}
	return fmt.Sprintf("%s-g%d", strings.TrimRight(prefix[:max(1, 61-len(strconv.FormatInt(generation, 10)))], "-"), generation)
}

func dnsName(value string) string {
	value = strings.ToLower(value)
	value = regexp.MustCompile(`[^a-z0-9.-]+`).ReplaceAllString(value, "-")
	value = strings.Trim(value, "-.")
	if len(value) > 63 {
		value = strings.TrimRight(value[:63], "-.")
	}
	return value
}

func conditionTrue(conditions []metav1.Condition, conditionType string, generation int64) bool {
	for _, condition := range conditions {
		if condition.Type == conditionType {
			return condition.Status == metav1.ConditionTrue && condition.ObservedGeneration == generation
		}
	}
	return false
}

func valueOr(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func int32PointerString(value *int32) string {
	if value == nil {
		return ""
	}
	return strconv.FormatInt(int64(*value), 10)
}

func int64PointerString(value *int64) string {
	if value == nil {
		return ""
	}
	return strconv.FormatInt(*value, 10)
}

func (r *AppDeployerReconciler) setCondition(app *platformv1alpha1.AppDeployer, conditionType string, status metav1.ConditionStatus, reason, message string) {
	meta.SetStatusCondition(&app.Status.Conditions, metav1.Condition{Type: conditionType, Status: status, Reason: reason, Message: message, ObservedGeneration: app.Generation})
}

func (r *AppDeployerReconciler) updateStatus(ctx context.Context, app *platformv1alpha1.AppDeployer) error {
	return r.Status().Update(ctx, app)
}

func (r *AppDeployerReconciler) fail(ctx context.Context, app *platformv1alpha1.AppDeployer, reason string, cause error) (ctrl.Result, error) {
	app.Status.Phase = "Error"
	r.setCondition(app, "Ready", metav1.ConditionFalse, reason, cause.Error())
	statusErr := r.updateStatus(ctx, app)
	if statusErr != nil {
		return ctrl.Result{}, errors.Join(cause, statusErr)
	}
	return ctrl.Result{}, cause
}

func (r *AppDeployerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&platformv1alpha1.AppDeployer{}).
		Owns(&batchv1.Job{}).
		Named("appdeployer").
		Complete(r)
}
