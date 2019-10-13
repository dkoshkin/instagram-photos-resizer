.DEFAULT_GOAL := build

PKG = github.com/dkoshkin/instagram-photos-resizer
# Set the build version
ifeq ($(origin VERSION), undefined)
	VERSION := $(shell git describe --tags --always)
endif
ifeq ($(origin GOOS), undefined)
	GOOS := $(shell go env GOOS)
endif

GO_MOD_SHASUM := $$(shasum go.sum | awk '{ print $$1 }' | cut -c1-3)
DOCKERFILE_SHASUM := $$(shasum ./build/docker/Dockerfile.gomod | awk '{ print $$1 }' | cut -c1-3)
# DOCKER_IMG_TAG is the tag of the builder image based on the go.sum and Dockerfile
DOCKER_IMG_TAG := $(shell echo $(DOCKERFILE_SHASUM)$(GO_MOD_SHASUM))

DOCKER_GOMOD_IMG_NAME := dkoshkin/instagram-photos-resizer-gomod
DOCKER_GOMOD_IMG := $(DOCKER_GOMOD_IMG_NAME):$(DOCKER_IMG_TAG)

DOCKER_IMG_NAME := dkoshkin/instagram-photos-resizer
DOCKER_IMG := $(DOCKER_IMG_NAME):$(VERSION)

.PHONY: build
build: builder.check
	docker run \
	--rm \
	-v "$(shell pwd)":"/src/$(PKG)" \
	-v "$(shell go env GOCACHE)":"/root/.cache/go-build" \
	-e CGO_ENABLED=0 \
	-e GOOS=$(GOOS) \
	$(DOCKER_GOMOD_IMG) \
	go build -ldflags '-w -extldflags "-static"' -v -o bin/$(GOOS)/instagram-photos-resizer ./cmd/cli/

.PHONY: cross
cross: builder.check
	@$(MAKE) GOOS=darwin build
	@$(MAKE) GOOS=windows build
	@$(MAKE) GOOS=linux build

.PHONY: image
image: builder.check
	docker build -f ./build/docker/Dockerfile -t $(DOCKER_IMG) .
	docker tag $(DOCKER_IMG) $(DOCKER_IMG_NAME):latest

.PHONY: image.push
image.push:
	docker push $(DOCKER_IMG)
	docker push $(DOCKER_IMG_NAME):latest

.PHONY: test
test: builder.check
	docker run \
	--rm \
	-v "$(shell pwd)":"/src/$(PKG)" \
	-v "$(shell go env GOCACHE)":"/root/.cache/go-build" \
	$(DOCKER_GOMOD_IMG) \
	go test -v ./cmd/... ./pkg/...

.PHONY: goimports
goimports:
	goimports -local $(shell go list -e) -w $(shell find . -name '*.go')

.PHONY: builder.check
builder.check:
	@docker image inspect $(DOCKER_GOMOD_IMG) > /dev/null || $(MAKE) gomod

.PHONY: gomod
gomod:
	docker build --target gomod -f build/docker/Dockerfile.gomod -t $(DOCKER_GOMOD_IMG) .
	docker tag $(DOCKER_GOMOD_IMG) $(DOCKER_GOMOD_IMG_NAME):latest

