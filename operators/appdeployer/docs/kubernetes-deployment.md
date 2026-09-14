# Deploy and use AppDeployer on Kubernetes

This procedure installs AppDeployer on a conformant Kubernetes cluster and uses upstream Tekton, Argo CD and, optionally, Knative Serving. Pin and approve component versions through the platform lifecycle process before using them in production.

## 1. Prerequisites

The administration workstation requires `kubectl`, Git, Go 1.24+, Kustomize and a container client. The target cluster must provide a default `StorageClass`, outbound access to the source/artifact/registry endpoints, and enough resources for Tekton build pods.

Confirm the active cluster before making changes:

```bash
kubectl config current-context
kubectl get nodes
kubectl get storageclass
```

## 2. Install the delivery dependencies

Install Tekton Pipelines using the release approved by the cluster administrator. The upstream quick-start uses the latest release manifest:

```bash
kubectl apply --filename \
  https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
kubectl wait --for=condition=Available deployment/tekton-pipelines-controller \
  --namespace tekton-pipelines --timeout=5m
```

Install Argo CD. Server-side apply is used because its CRDs can exceed the client-side annotation limit:

```bash
kubectl create namespace argocd
kubectl apply --namespace argocd --server-side --force-conflicts \
  --filename https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
kubectl rollout status statefulset/argocd-application-controller \
  --namespace argocd --timeout=5m
```

For a serverless target, install Knative Serving and one supported networking layer. The following upstream example uses Knative 1.23 and Kourier:

```bash
kubectl apply --filename \
  https://github.com/knative/serving/releases/download/knative-v1.23.0/serving-crds.yaml
kubectl apply --filename \
  https://github.com/knative/serving/releases/download/knative-v1.23.0/serving-core.yaml
kubectl apply --filename \
  https://github.com/knative-extensions/net-kourier/releases/download/knative-v1.23.0/kourier.yaml
kubectl patch configmap/config-network --namespace knative-serving --type merge \
  --patch '{"data":{"ingress-class":"kourier.ingress.networking.knative.dev"}}'
kubectl apply --filename \
  https://github.com/knative/serving/releases/download/knative-v1.23.0/serving-default-domain.yaml
kubectl get pods --namespace knative-serving
```

Use Knative's compatibility matrix rather than copying these version numbers into a managed production cluster without review.

## 3. Build and install AppDeployer

From `operators/appdeployer`:

```bash
export APPDEPLOYER_IMAGE=quay.io/your-org/appdeployer-operator:0.2.0

make test
make docker-build IMG="$APPDEPLOYER_IMAGE"
make docker-push IMG="$APPDEPLOYER_IMAGE"
make install
make deploy IMG="$APPDEPLOYER_IMAGE"

kubectl rollout status deployment/appdeployer-controller-manager \
  --namespace appdeployer-system --timeout=3m
kubectl get crd appdeployers.platform.openindustries.org
```

## 4. Prepare the application namespace and credentials

```bash
kubectl create namespace orders-dev

kubectl create secret docker-registry registry-push-credentials \
  --namespace orders-dev \
  --docker-server=quay.io \
  --docker-username="$REGISTRY_USERNAME" \
  --docker-password="$REGISTRY_PASSWORD"

kubectl create secret generic gitops-credentials \
  --namespace orders-dev \
  --from-literal=username="$GIT_USERNAME" \
  --from-literal=password="$GIT_TOKEN"
```

For private source code, create a second `username`/`password` Secret. For SonarQube, use a Secret with a `token` key. Maven settings must be stored under a `settings.xml` key.

Register a private GitOps repository with Argo CD separately. AppDeployer intentionally does not copy namespaced credentials into the Argo CD control-plane namespace.

## 5. Prepare the GitOps repository

Every selected `gitOps.path` must contain `kustomization.yaml`. A standard application normally contains a Kubernetes `Deployment` and `Service`. A serverless application must contain a `serving.knative.dev/v1` `Service`; a ready-to-copy example is available under `examples/gitops/serverless`.

For serverless:

```bash
cp -R examples/gitops/serverless /path/to/gitops-repository/environments/serverless-development
git -C /path/to/gitops-repository add environments/serverless-development
git -C /path/to/gitops-repository commit -m "feat: add serverless application base"
git -C /path/to/gitops-repository push
```

## 6. Deploy a standard application

Copy the development sample and change the repository, registry and destination values. On upstream Kubernetes, `argoCDNamespace` must be `argocd`:

```yaml
spec:
  deployment:
    type: standard
  gitOps:
    argoCDNamespace: argocd
```

Apply and observe it:

```bash
kubectl apply --filename config/samples/platform_v1alpha1_appdeployer_development.yaml
kubectl get appdeployer --namespace orders-dev --watch
kubectl get pipelines,pipelineruns --namespace orders-dev
kubectl get applications.argoproj.io --namespace argocd
```

## 7. Deploy a serverless application

Edit `config/samples/platform_v1alpha1_appdeployer_serverless.yaml`, set `argoCDNamespace: argocd`, and apply it:

```bash
kubectl create namespace orders-serverless
# Create registry-push-credentials and gitops-credentials in orders-serverless.
kubectl apply --filename config/samples/platform_v1alpha1_appdeployer_serverless.yaml

kubectl get appdeployer --namespace orders-serverless --watch
kubectl get ksvc --namespace orders-serverless
kubectl get revision,route --namespace orders-serverless
```

The operator updates both the Kustomize image and these Knative settings in `service.yaml`: minimum and maximum scale, scaling target/metric, container concurrency, request timeout and route visibility. `minScale: 0` permits scale-to-zero. Set `visibility: cluster-local` when the service must not expose an external route.

## 8. Promote an existing image to production

Production never creates Tekton resources. It runs an OCI inspection Job, commits the requested tag to GitOps and then ensures the Argo CD Application:

```bash
kubectl create namespace orders-prod
# Create registry-pull-credentials and gitops-credentials in orders-prod.
kubectl apply --filename config/samples/platform_v1alpha1_appdeployer_production.yaml

kubectl get jobs --namespace orders-prod
kubectl get appdeployer orders-api-production --namespace orders-prod --output yaml
```

## 9. Troubleshooting

```bash
kubectl logs deployment/appdeployer-controller-manager --namespace appdeployer-system
kubectl describe appdeployer NAME --namespace APP_NAMESPACE
kubectl logs --namespace APP_NAMESPACE -l tekton.dev/pipelineRun=PIPELINE_RUN_NAME --all-containers
kubectl describe application APP_NAMESPACE-APP_NAME --namespace argocd
```

If Buildah is rejected, inspect the cluster's Pod Security admission policy. AppDeployer does not weaken cluster-wide policy automatically; authorize the narrow capabilities needed by the pipeline ServiceAccount through the platform's security process.
