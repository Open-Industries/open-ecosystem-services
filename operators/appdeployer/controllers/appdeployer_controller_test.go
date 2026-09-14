package controllers

import (
	"testing"

	platformv1alpha1 "github.com/Open-Industries/open-ecosystem-services/operators/appdeployer/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestDevelopmentPipelineContainsConfiguredStages(t *testing.T) {
	app := developmentApp()
	data := buildTemplateData(app)
	pipeline, err := renderUnstructured("pipeline", pipelineTemplate, data)
	if err != nil {
		t.Fatalf("render pipeline: %v", err)
	}
	tasks, found, err := unstructured.NestedSlice(pipeline.Object, "spec", "tasks")
	if err != nil || !found {
		t.Fatalf("pipeline tasks missing: found=%v err=%v", found, err)
	}
	want := map[string]bool{"clone-source": false, "configure-dependencies": false, "unit-tests": false, "sonarqube-quality-gate": false, "s2i-build": false, "push-image": false, "integration-tests": false, "update-gitops": false}
	for _, task := range tasks {
		if item, ok := task.(map[string]any); ok {
			if name, ok := item["name"].(string); ok {
				if _, tracked := want[name]; tracked {
					want[name] = true
				}
			}
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("expected task %q", name)
		}
	}
}

func TestProductionRequiresNoSourceOrBuilder(t *testing.T) {
	app := productionApp()
	if err := validateSpec(app); err != nil {
		t.Fatalf("production spec rejected: %v", err)
	}
	data := buildTemplateData(app)
	if _, err := renderJob("image-validation", imageValidationJobTemplate, data); err != nil {
		t.Fatalf("render image validation Job: %v", err)
	}
	if _, err := renderJob("gitops-update", gitOpsUpdateJobTemplate, data); err != nil {
		t.Fatalf("render GitOps Job: %v", err)
	}
	if _, err := renderUnstructured("application", applicationTemplate, data); err != nil {
		t.Fatalf("render Application: %v", err)
	}
}

func developmentApp() *platformv1alpha1.AppDeployer {
	return &platformv1alpha1.AppDeployer{ObjectMeta: metav1.ObjectMeta{Name: "orders", Namespace: "orders-dev", Generation: 7}, Spec: platformv1alpha1.AppDeployerSpec{
		Env: platformv1alpha1.EnvironmentDevelopment, BuilderImage: "example.com/s2i-java:latest",
		Source:             &platformv1alpha1.GitSource{URL: "https://github.com/example/orders.git", Revision: "main", ContextDir: ".", CredentialsSecretRef: &platformv1alpha1.SecretReference{Name: "source-git"}},
		ArtifactRepository: &platformv1alpha1.ArtifactRepository{Type: "nexus", URL: "https://nexus.example.com", CredentialsSecretRef: &platformv1alpha1.SecretReference{Name: "nexus"}, MavenSettingsSecretRef: &platformv1alpha1.SecretReference{Name: "maven-settings"}},
		SonarQube:          &platformv1alpha1.SonarQubeConfig{Enabled: true, URL: "https://sonar.example.com", ProjectKey: "orders", TokenSecretRef: &platformv1alpha1.SecretReference{Name: "sonar-token"}, QualityGateWait: true},
		Registry:           platformv1alpha1.RegistryConfig{Image: "quay.io/example/orders", Tag: "dev-7", CredentialsSecretRef: &platformv1alpha1.SecretReference{Name: "registry-auth"}},
		Tests:              platformv1alpha1.TestConfig{UnitTests: true, IntegrationTests: true},
		GitOps:             platformv1alpha1.GitOpsConfig{RepositoryURL: "https://github.com/example/gitops.git", Path: "environments/dev", KustomizationImage: "quay.io/example/orders", CredentialsSecretRef: &platformv1alpha1.SecretReference{Name: "gitops-auth"}, DestinationNamespace: "orders-dev", AutoSync: true},
	}}
}

func productionApp() *platformv1alpha1.AppDeployer {
	app := developmentApp()
	app.Name = "orders-production"
	app.Namespace = "orders-prod"
	app.Spec.Env = platformv1alpha1.EnvironmentProduction
	app.Spec.Source = nil
	app.Spec.BuilderImage = ""
	app.Spec.ArtifactRepository = nil
	app.Spec.SonarQube = nil
	app.Spec.Tests = platformv1alpha1.TestConfig{}
	app.Spec.GitOps.Path = "environments/production"
	app.Spec.GitOps.DestinationNamespace = "orders-prod"
	return app
}
