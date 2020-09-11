# Current Controller version
VERSION ?= 0.0.1

# Image URL to use all building/pushing image targets
IMG ?= aws-shield-controller:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif


# Run tests
test: fmt vet
	go test ./... -coverprofile cover.out

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	go run ./main.go

# Install CRDs into a cluster
install:  kustomize
	kustomize build kustomize/ | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall:  kustomize
	kustomize build kustomize/ | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: kustomize
	cd kustomize && kustomize edit set image lehtux/aws-shield-controller=${IMG}
	kustomize build kustomize | kubectl apply -f -

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

build:
	go build .

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}
