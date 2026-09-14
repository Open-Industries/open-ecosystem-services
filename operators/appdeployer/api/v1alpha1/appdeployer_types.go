package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Environment determines whether AppDeployer builds an image or only promotes one.
// +kubebuilder:validation:Enum=development;production
type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentProduction  Environment = "production"
)

// DeploymentType selects a conventional Kubernetes workload or Knative Serving.
// +kubebuilder:validation:Enum=standard;serverless
type DeploymentType string

const (
	DeploymentTypeStandard   DeploymentType = "standard"
	DeploymentTypeServerless DeploymentType = "serverless"
)

// SecretReference points to a Secret in the AppDeployer namespace.
type SecretReference struct {
	// Name is the Secret name.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}

// GitSource defines the application source repository.
type GitSource struct {
	// URL is the HTTPS or SSH clone URL.
	// +kubebuilder:validation:Pattern=`^(https://|ssh://|git@).+`
	URL string `json:"url"`
	// Revision is a branch, tag, or commit.
	// +kubebuilder:default=main
	Revision string `json:"revision,omitempty"`
	// ContextDir is the application directory inside the repository.
	// +kubebuilder:default=.
	ContextDir string `json:"contextDir,omitempty"`
	// CredentialsSecretRef optionally references a Tekton-compatible Git Secret.
	CredentialsSecretRef *SecretReference `json:"credentialsSecretRef,omitempty"`
}

// ArtifactRepository configures Nexus or JFrog Artifactory.
type ArtifactRepository struct {
	// Type selects the repository implementation.
	// +kubebuilder:validation:Enum=nexus;artifactory
	Type string `json:"type"`
	// URL is the repository base URL.
	// +kubebuilder:validation:Pattern=`^https?://.+`
	URL string `json:"url"`
	// CredentialsSecretRef references username/password or token credentials.
	CredentialsSecretRef *SecretReference `json:"credentialsSecretRef,omitempty"`
	// MavenSettingsSecretRef references a Secret containing settings.xml.
	MavenSettingsSecretRef *SecretReference `json:"mavenSettingsSecretRef,omitempty"`
}

// SonarQubeConfig configures static analysis and Quality Gate enforcement.
type SonarQubeConfig struct {
	// Enabled controls execution of SonarQube analysis.
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty"`
	// URL is the SonarQube server URL.
	// +kubebuilder:validation:Pattern=`^https?://.+`
	URL string `json:"url,omitempty"`
	// ProjectKey identifies the project in SonarQube.
	ProjectKey string `json:"projectKey,omitempty"`
	// TokenSecretRef references a Secret whose token key contains the API token.
	TokenSecretRef *SecretReference `json:"tokenSecretRef,omitempty"`
	// QualityGateWait makes the scanner wait and fail the task when the gate fails.
	// +kubebuilder:default=true
	QualityGateWait bool `json:"qualityGateWait,omitempty"`
	// QualityGateTimeoutSeconds is the maximum wait time for the gate result.
	// +kubebuilder:default=300
	// +kubebuilder:validation:Minimum=30
	QualityGateTimeoutSeconds int32 `json:"qualityGateTimeoutSeconds,omitempty"`
}

// RegistryConfig defines the OCI image destination and authentication.
type RegistryConfig struct {
	// Image is the repository without a tag, for example quay.io/acme/orders.
	Image string `json:"image"`
	// Tag is the immutable or release tag to build/promote.
	// +kubebuilder:validation:MinLength=1
	Tag string `json:"tag"`
	// CredentialsSecretRef references a dockerconfigjson Secret.
	CredentialsSecretRef *SecretReference `json:"credentialsSecretRef,omitempty"`
	// InsecureSkipTLSVerify is intended only for development registries.
	// +kubebuilder:default=false
	InsecureSkipTLSVerify bool `json:"insecureSkipTLSVerify,omitempty"`
}

// TestConfig controls optional pipeline test stages.
type TestConfig struct {
	// UnitTests enables the unit-test task.
	// +kubebuilder:default=true
	UnitTests bool `json:"unitTests,omitempty"`
	// UnitTestCommand is executed from the source context directory.
	// +kubebuilder:default="./mvnw test"
	UnitTestCommand string `json:"unitTestCommand,omitempty"`
	// IntegrationTests enables post-build integration tests.
	// +kubebuilder:default=false
	IntegrationTests bool `json:"integrationTests,omitempty"`
	// IntegrationTestCommand receives IMAGE_URL in its environment.
	IntegrationTestCommand string `json:"integrationTestCommand,omitempty"`
	// TestImage provides the tools needed by the configured commands.
	// +kubebuilder:default="registry.access.redhat.com/ubi9/openjdk-21:latest"
	TestImage string `json:"testImage,omitempty"`
}

// ServerlessConfig controls the Knative Service stored in the GitOps repository.
type ServerlessConfig struct {
	// Platform identifies the serverless API implementation.
	// +kubebuilder:default=knative
	// +kubebuilder:validation:Enum=knative
	Platform string `json:"platform,omitempty"`
	// ServiceName is the metadata.name of the Knative Service.
	ServiceName string `json:"serviceName"`
	// ManifestPath is relative to gitOps.path and identifies the Knative Service YAML.
	// +kubebuilder:default=service.yaml
	ManifestPath string `json:"manifestPath,omitempty"`
	// ContainerName identifies the container receiving the built image.
	// +kubebuilder:default=user-container
	ContainerName string `json:"containerName,omitempty"`
	// MinScale allows zero so the service can scale to zero.
	// +kubebuilder:validation:Minimum=0
	MinScale *int32 `json:"minScale,omitempty"`
	// MaxScale limits the number of replicas.
	// +kubebuilder:validation:Minimum=1
	MaxScale *int32 `json:"maxScale,omitempty"`
	// ScaleTarget is the per-replica target for the selected metric.
	// +kubebuilder:validation:Minimum=1
	ScaleTarget *int32 `json:"scaleTarget,omitempty"`
	// Metric selects the Knative Pod Autoscaler metric.
	// +kubebuilder:default=concurrency
	// +kubebuilder:validation:Enum=concurrency;rps
	Metric string `json:"metric,omitempty"`
	// ContainerConcurrency is the hard concurrent-request limit per replica; zero means unlimited.
	// +kubebuilder:validation:Minimum=0
	ContainerConcurrency *int64 `json:"containerConcurrency,omitempty"`
	// TimeoutSeconds is the maximum request execution time.
	// +kubebuilder:validation:Minimum=1
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty"`
	// Visibility controls whether the Knative Route is externally reachable.
	// +kubebuilder:default=external
	// +kubebuilder:validation:Enum=external;cluster-local
	Visibility string `json:"visibility,omitempty"`
}

// DeploymentConfig defines how Argo CD materializes the application runtime.
// +kubebuilder:validation:XValidation:rule="self.type != 'serverless' || has(self.serverless)",message="serverless configuration is required when deployment.type is serverless"
type DeploymentConfig struct {
	// Type selects Deployment/Service manifests or a Knative Service.
	// +kubebuilder:default=standard
	Type DeploymentType `json:"type,omitempty"`
	// Serverless contains Knative-specific settings.
	Serverless *ServerlessConfig `json:"serverless,omitempty"`
}

// GitOpsConfig defines the deployment source of truth and Argo CD destination.
type GitOpsConfig struct {
	// RepositoryURL is the GitOps repository clone URL.
	RepositoryURL string `json:"repositoryURL"`
	// Revision is the GitOps branch updated and watched by Argo CD.
	// +kubebuilder:default=main
	Revision string `json:"revision,omitempty"`
	// Path is the environment-specific manifests or Kustomize directory.
	Path string `json:"path"`
	// KustomizationImage is the image name replaced by `kustomize edit set image`.
	KustomizationImage string `json:"kustomizationImage"`
	// CredentialsSecretRef references username/password or SSH Git credentials.
	CredentialsSecretRef *SecretReference `json:"credentialsSecretRef,omitempty"`
	// ArgoCDNamespace contains the Application resource.
	// +kubebuilder:default=openshift-gitops
	ArgoCDNamespace string `json:"argoCDNamespace,omitempty"`
	// Project is the Argo CD project.
	// +kubebuilder:default=default
	Project string `json:"project,omitempty"`
	// DestinationServer is the Kubernetes API endpoint known to Argo CD.
	// +kubebuilder:default="https://kubernetes.default.svc"
	DestinationServer string `json:"destinationServer,omitempty"`
	// DestinationNamespace receives the workload.
	DestinationNamespace string `json:"destinationNamespace"`
	// AutoSync enables prune and self-heal in the generated Application.
	// +kubebuilder:default=true
	AutoSync bool `json:"autoSync,omitempty"`
}

// AppDeployerSpec defines the desired state.
// +kubebuilder:validation:XValidation:rule="self.env != 'development' || (has(self.source) && has(self.builderImage) && size(self.builderImage) > 0)",message="source and builderImage are required when env is development"
type AppDeployerSpec struct {
	// Env selects CI+CD (development) or promotion-only CD (production).
	Env Environment `json:"env"`
	// BuilderImage is any S2I-compatible builder image.
	BuilderImage string `json:"builderImage,omitempty"`
	// Source is required for development builds and ignored in production.
	Source             *GitSource          `json:"source,omitempty"`
	ArtifactRepository *ArtifactRepository `json:"artifactRepository,omitempty"`
	SonarQube          *SonarQubeConfig    `json:"sonarQube,omitempty"`
	Registry           RegistryConfig      `json:"registry"`
	Tests              TestConfig          `json:"tests,omitempty"`
	// Deployment selects standard Kubernetes or Knative serverless delivery.
	Deployment DeploymentConfig `json:"deployment,omitempty"`
	GitOps     GitOpsConfig     `json:"gitOps"`
}

// AppDeployerStatus reports reconciliation and delivery progress.
type AppDeployerStatus struct {
	ObservedGeneration  int64              `json:"observedGeneration,omitempty"`
	Phase               string             `json:"phase,omitempty"`
	ImageRef            string             `json:"imageRef,omitempty"`
	PipelineRunName     string             `json:"pipelineRunName,omitempty"`
	GitOpsJobName       string             `json:"gitOpsJobName,omitempty"`
	ArgoApplicationName string             `json:"argoApplicationName,omitempty"`
	Conditions          []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=appdep,categories=open-industries
// +kubebuilder:printcolumn:name="Environment",type=string,JSONPath=`.spec.env`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Image",type=string,JSONPath=`.status.imageRef`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
type AppDeployer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AppDeployerSpec   `json:"spec,omitempty"`
	Status            AppDeployerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type AppDeployerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppDeployer `json:"items"`
}
