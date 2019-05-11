Version := $(shell git describe --tags --dirty)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)"
PLATFORM := $(shell ./hack/platform-tag.sh)
DOCKER_HUB_REPO ?= alexellis2/inlets
BUILDKIT_PUSH_TO_REPO ?= true

BUILDKIT_HOST ?= tcp://0.0.0.0:1234
BUILDKITD_CONTAINER_NAME ?= buildkitd
BUILDKITD_CONTAINER_STOPPED := $(shell docker ps --filter name=$(BUILDKITD_CONTAINER_NAME) --filter status=exited --format='{{.Names}}' 2>/dev/null)
BUILDKITD_CONTAINER_RUNNING := $(shell docker ps --filter name=$(BUILDKITD_CONTAINER_NAME) --filter status=running --format='{{.Names}}' 2>/dev/null)

BUILDKIT_COMMON_ARGS =  --progress=plain
BUILDKIT_COMMON_ARGS += --frontend dockerfile.v0
BUILDKIT_COMMON_ARGS += --local dockerfile=.
BUILDKIT_COMMON_ARGS += --local context=.
BUILDKIT_COMMON_ARGS += --opt filename=./Dockerfile.multi-arch
BUILDKIT_COMMON_ARGS += --opt build-arg:GIT_COMMIT=$(GitCommit)
BUILDKIT_COMMON_ARGS += --opt build-arg:VERSION=$(Version)

# This option is for running docker manifest command
export DOCKER_CLI_EXPERIMENTAL := enabled
export BUILDKIT_HOST

.PHONY: all
all: docker

.PHONY: dist
dist:
	CGO_ENABLED=0 GOOS=linux go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/inlets
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/inlets-darwin
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/inlets-armhf
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/inlets-arm64

.PHONY: docker
docker:
	docker build --build-arg Version=$(Version) --build-arg GIT_COMMIT=$(GitCommit) -t $(DOCKER_HUB_REPO):$(Version)$(PLATFORM) .


.PHONY: buildkitd
buildkitd:
	docker run --privileged linuxkit/binfmt:v0.7
ifeq (tcp://0.0.0.0:1234,$(findstring tcp://0.0.0.0:1234,$(BUILDKIT_HOST)))
ifeq ($(BUILDKITD_CONTAINER_STOPPED),$(BUILDKITD_CONTAINER_NAME))
	@echo "Removing exited buildkitd container"
	docker rm $(BUILDKITD_CONTAINER_NAME)
endif
ifneq ($(BUILDKITD_CONTAINER_RUNNING),$(BUILDKITD_CONTAINER_NAME))
	@echo "Starting buildkitd container"
	docker run -d --privileged -p 1234:1234 --name $(BUILDKITD_CONTAINER_NAME) moby/buildkit:latest \
          --addr $(BUILDKIT_HOST) \
          --oci-worker-platform linux/amd64 \
          --oci-worker-platform linux/arm64 \
          --oci-worker-platform linux/armhf
	sudo docker cp $(BUILDKITD_CONTAINER_NAME):/usr/bin/buildctl /usr/bin/
endif
endif

# push=true assumes valid authenticated login into registry
.PHONY: docker-buildkit-images
docker-buildkit-images:
	for arch in amd64 arm64 armhf; do \
		buildctl build $(BUILDKIT_COMMON_ARGS) \
          --opt platform=linux/$${arch} \
          --output type=image,name=docker.io/$(DOCKER_HUB_REPO):$(Version)-$${arch},push=$(BUILDKIT_PUSH_TO_REPO) ; \
	done

# push tags for latest and create multi-arch image manifest
.PHONY: docker-buildkit-push-manifests
docker-buildkit-push-manifests:
ifeq (true,$(BUILDKIT_PUSH_TO_REPO))
	for arch in amd64 arm64 armhf; do \
		docker pull $(DOCKER_HUB_REPO):$(Version)-$${arch} && \
		docker tag $(DOCKER_HUB_REPO):$(Version)-$${arch} $(DOCKER_HUB_REPO):latest-$${arch} && \
		docker push $(DOCKER_HUB_REPO):latest-$${arch} ; \
	done
	for version in $(Version) latest; do \
		docker manifest create --amend $(DOCKER_HUB_REPO):$${version} \
			$(DOCKER_HUB_REPO):$${version}-amd64 \
			$(DOCKER_HUB_REPO):$${version}-arm64 \
			$(DOCKER_HUB_REPO):$${version}-armhf && \
		docker manifest annotate $(DOCKER_HUB_REPO):$${version} \
			$(DOCKER_HUB_REPO):$${version}-armhf --os linux --arch arm   --variant v6 && \
		docker manifest annotate $(DOCKER_HUB_REPO):$${version} \
			$(DOCKER_HUB_REPO):$${version}-arm64 --os linux --arch arm64 --variant v8 && \
		docker manifest push --purge $(DOCKER_HUB_REPO):$${version} ;\
	done

endif




.PHONY: docker-login
docker-login:
	echo -n "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin

.PHONY: push
push:
	docker push $(DOCKER_HUB_REPO):$(Version)$(PLATFORM)
