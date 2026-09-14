# Deploy and use AppDeployer on OpenShift

This procedure uses Red Hat OpenShift Pipelines for Tekton, OpenShift GitOps for Argo CD and OpenShift Serverless for Knative Serving.

## 1. Install the OpenShift Operators

As a cluster administrator, install these Operators from **Administrator → Operators → OperatorHub**:

1. Red Hat OpenShift Pipelines.
2. Red Hat OpenShift GitOps.
3. Red Hat OpenShift Serverless, when `deployment.type: serverless` will be used.

Verify the APIs instead of assuming that a successful subscription means the operands are ready:

```bash
oc get csv --namespace openshift-operators
oc get crd pipelines.tekton.dev pipelineruns.tekton.dev applications.argoproj.io
oc get pods --namespace openshift-pipelines
oc get pods --namespace openshift-gitops
```

For serverless workloads, create Knative Serving after the Serverless Operator is ready:

```bash
oc create namespace knative-serving
oc apply --filename - <<'EOF'
apiVersion: operator.knative.dev/v1beta1
kind: KnativeServing
metadata:
  name: knative-serving
  namespace: knative-serving
spec: {}
EOF

oc get knativeserving knative-serving --namespace knative-serving
oc get pods --namespace knative-serving
oc get crd services.serving.knative.dev
```

## 2. Build and deploy the Operator

Log in to the image registry and run these commands from `operators/appdeployer`:

```bash
export APPDEPLOYER_IMAGE=quay.io/your-org/appdeployer-operator:0.2.0

make test
make docker-build IMG="$APPDEPLOYER_IMAGE" CONTAINER_TOOL=podman
make docker-push IMG="$APPDEPLOYER_IMAGE" CONTAINER_TOOL=podman
make install KUBECTL=oc
make deploy IMG="$APPDEPLOYER_IMAGE" KUBECTL=oc

oc rollout status deployment/appdeployer-controller-manager \
  --namespace appdeployer-system --timeout=3m
oc logs deployment/appdeployer-controller-manager --namespace appdeployer-system
```

The generated manifests use restricted container settings for the controller. S2I/Buildah build pods run through the `appdeployer-pipeline` ServiceAccount created in each application namespace. Confirm that the OpenShift Pipelines SCC policy admits the requested build capabilities; do not bind `privileged` globally.

## 3. Prepare credentials

```bash
oc new-project orders-dev

oc create secret docker-registry registry-push-credentials \
  --docker-server=quay.io \
  --docker-username="$REGISTRY_USERNAME" \
  --docker-password="$REGISTRY_PASSWORD"

oc create secret generic gitops-credentials \
  --from-literal=username="$GIT_USERNAME" \
  --from-literal=password="$GIT_TOKEN"

oc create secret generic sonarqube-token \
  --from-literal=token="$SONAR_TOKEN"

oc create secret generic maven-settings \
  --from-file=settings.xml=/secure/path/settings.xml
```

Create secrets in every namespace containing an AppDeployer. Secrets are referenced, never copied into status or generated manifests.

For a private GitOps repository, configure repository access in the `openshift-gitops` Argo CD instance as a separate administrative action.

## 4. Standard OpenShift deployment

Review the development sample, keep `deployment.type: standard`, and ensure the GitOps path contains a `Deployment`/`Service` Kustomize base:

```bash
oc apply --filename config/samples/platform_v1alpha1_appdeployer_development.yaml
oc get appdep,pipeline,pipelinerun --namespace orders-dev
oc get application --namespace openshift-gitops
```

The development PipelineRun clones the source, injects artifact settings, runs enabled tests and SonarQube, generates an S2I Dockerfile, builds and pushes the image, updates GitOps and lets OpenShift GitOps synchronize it.

## 5. OpenShift Serverless deployment

Prepare the GitOps path using `examples/gitops/serverless`, then create the namespace and its registry/Git credentials:

```bash
oc new-project orders-serverless
# Create registry-push-credentials and gitops-credentials here.
oc apply --filename config/samples/platform_v1alpha1_appdeployer_serverless.yaml

oc get appdeployer --namespace orders-serverless --watch
oc get pipelinerun --namespace orders-serverless
oc get ksvc,revision,route --namespace orders-serverless
```

Example serverless configuration:

```yaml
deployment:
  type: serverless
  serverless:
    platform: knative
    serviceName: orders-functions
    manifestPath: service.yaml
    containerName: user-container
    minScale: 0
    maxScale: 20
    scaleTarget: 50
    metric: concurrency
    containerConcurrency: 100
    timeoutSeconds: 300
    visibility: external
```

Use `visibility: cluster-local` for internal cross-microservice endpoints. For externally visible workloads, obtain the route from Knative:

```bash
oc get ksvc orders-functions --namespace orders-serverless \
  --output jsonpath='{.status.url}{"\n"}'
```

## 6. Production promotion

The production CR needs only the image reference, GitOps configuration and deployment mode. It does not need `source`, `builderImage`, tests, Nexus or SonarQube:

```bash
oc new-project orders-prod
# Create registry-pull-credentials and gitops-credentials here.
oc apply --filename config/samples/platform_v1alpha1_appdeployer_production.yaml

oc get job --namespace orders-prod
oc get appdeployer orders-api-production --namespace orders-prod --output yaml
```

For serverless production, copy the `deployment.serverless` block into the production CR and point `gitOps.path` to the production Knative Kustomize base.

## 7. Day-2 operations

Trigger a new immutable delivery execution by changing `spec.registry.tag`. Follow status conditions rather than relying only on pod state:

```bash
oc patch appdeployer orders-api --namespace orders-dev --type merge \
  --patch '{"spec":{"registry":{"tag":"dev-1.0.1"}}}'
oc get appdeployer orders-api --namespace orders-dev \
  --output jsonpath='{range .status.conditions[*]}{.type}{"="}{.status}{" "}{.reason}{"\n"}{end}'
```

Useful diagnostics:

```bash
oc logs deployment/appdeployer-controller-manager --namespace appdeployer-system
oc describe pipelinerun PIPELINE_RUN --namespace APP_NAMESPACE
oc logs --namespace APP_NAMESPACE -l tekton.dev/pipelineRun=PIPELINE_RUN --all-containers
oc describe application APP_NAMESPACE-APP_NAME --namespace openshift-gitops
oc get events --namespace APP_NAMESPACE --sort-by=.lastTimestamp
```
