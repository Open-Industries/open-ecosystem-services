IMG ?= quay.io/open-industries/appdeployer-operator:0.2.0
CONTAINER_TOOL ?= podman
KUBECTL ?= kubectl
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
LOCALBIN ?= $(CURDIR)/bin
CONTROLLER_TOOLS_VERSION ?= v0.18.0

.PHONY: manifests generate fmt vet test build docker-build docker-push install uninstall deploy undeploy

manifests: controller-gen
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="" paths="./..."

fmt:
	go fmt ./...

vet:
	go vet ./...

test: manifests generate fmt vet
	go test ./... -coverprofile cover.out

build: manifests generate fmt vet
	go build -o bin/manager ./cmd/main.go

docker-build:
	$(CONTAINER_TOOL) build -t ${IMG} .

docker-push:
	$(CONTAINER_TOOL) push ${IMG}

install: manifests
	$(KUBECTL) apply -f config/crd/bases

uninstall:
	$(KUBECTL) delete -f config/crd/bases --ignore-not-found

deploy: manifests
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default | $(KUBECTL) apply -f -

undeploy:
	kustomize build config/default | $(KUBECTL) delete --ignore-not-found -f -

controller-gen:
	@mkdir -p $(LOCALBIN)
	@test -x $(CONTROLLER_GEN) || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)
