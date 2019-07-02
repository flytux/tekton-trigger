
# Image URL to use all building/pushing image targets
IMG=registry.gitlab.com/pongsatt/githook/controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

WEBHOOK_IMG=registry.gitlab.com/pongsatt/githook/webhook:latest

export WEBHOOK_IMG

all: manager

# Run tests on manager
test: generate fmt vet manifests
	go test ./api/... ./controllers/... -coverprofile cover.out

# Run tests on webhook
wh-test:
	go test ./pkg/... ./cmd/hook/... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager cmd/manager/main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./cmd/manager/main.go

# Run webhook server
wh-run:
	go run ./cmd/hook/main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crd/bases

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy:
	kubectl apply -f config/crd/bases
	kustomize build config/default | kubectl apply -f -

# release all required kubernetes yamls in one file
release:
	kustomize build config/default > release.yaml

# Undeploy controller in the configured Kubernetes cluster in ~/.kube/config
undeploy:
	kustomize build config/default | kubectl delete -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths=./api/...

# Build the docker image on manager
docker-build: test
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e "s@image: .*@image: ${IMG}@" ./config/default/manager_image_patch.yaml

# Build the docker image on webhook
wh-docker-build: wh-test
	docker build . -f Dockerfile.wh -t ${WEBHOOK_IMG}
	@echo "updating kustomize image patch file for webhook image resource"
	sed -i'' -e "s@value: .*@value: ${WEBHOOK_IMG}@" ./config/default/manager_image_patch.yaml

# Push the docker image on manager
docker-push: docker-build
	docker push ${IMG}

# Push the docker image on webhook
wh-docker-push: wh-docker-build
	docker push ${WEBHOOK_IMG}

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.0-beta.2
CONTROLLER_GEN=$(shell go env GOPATH)/bin/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
