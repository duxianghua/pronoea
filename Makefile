.PHONY: version git-version init run

_GIT_LAST_COMMIT_TIME=$(shell TZ=UTC git log --pretty=format:'%cd' -1 --date=format-local:'%Y%m%d-%H%M%S')
_GIT_LAST_COMMIT_HASH=$(shell git rev-parse --short HEAD)
_GIT_VERSION=$(_GIT_LAST_COMMIT_TIME).$(_GIT_LAST_COMMIT_HASH)

_VERSION=$(shell cat Version)
_VERSION2=""
DOCKER_IMAGE_PREFIX=xingba/
DOCKER_IMAGE_NAME=$(DOCKER_IMAGE_PREFIX)pronoea:$(_VERSION)

GOCMD=go
GOOS=linux
GOARCH=amd64
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

GOLDFLAGS += -X main.Version=$(_VERSION)
GOFLAGS = -ldflags "$(GOLDFLAGS)"

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)


## Tool Binaries
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest

## Tool Versions
KUSTOMIZE_VERSION ?= v3.8.7
CONTROLLER_TOOLS_VERSION ?= v0.10.0


version:
	@echo ${_VERSION}

git-version:
	@echo ${_GIT_VERSION}

init:
	@mkdir -p bin

build: image-build

run: generate manifests apply-manifests
	go run cmd/main.go --debug

stop: container-stop

# Local
local-run:
	@$(GORUN) $(GOFLAGS) -v ./cmd

local-build: init
	@$(GOBUILD) $(GOFLAGS) -o bin/pronoea -v ./cmd

# Docker
image-build:
	docker build -t $(DOCKER_IMAGE_NAME) .

image-build-push: image-build
	docker push $(DOCKER_IMAGE_NAME)

container-run: build
	DOCKER_IMAGE_NAME=$(DOCKER_IMAGE_NAME) docker-compose up -d

container-stop:
	DOCKER_IMAGE_NAME=$(DOCKER_IMAGE_NAME) docker-compose down -v

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./internal/api/..."

.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./internal/api/..." output:crd:artifacts:config=manifests/crd output:rbac:artifacts:config=manifests/rbac

apply-manifests:
	kubectl apply -f manifests/crd



load-image: image-build
	kind kind load docker-image $(DOCKER_IMAGE_NAME)