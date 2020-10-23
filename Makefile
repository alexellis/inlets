Version := $(shell git describe --tags --dirty)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)"

# docker manifest command will work with Docker CLI 18.03 or newer
# but for now it's still experimental feature so we need to enable that
export DOCKER_CLI_EXPERIMENTAL=enabled

PLATFORMS?=linux/amd64,linux/arm/v6,linux/arm64
OUTPUT=
PROGRESS=plain

REGISTRY=docker.io

.PHONY: all
all: docker

.PHONY: dist
dist:
	CGO_ENABLED=0 GOOS=linux go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/inlets
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/inlets-darwin
	CGO_ENABLED=0 GOOS=freebsd go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/inlets-freebsd
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/inlets-armhf
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/inlets-arm64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/inlets.exe

.PHONY: docker
docker:
	docker buildx build \
		--platform=${PLATFORMS} $(OUTPUT) \
		--progress=$(PROGRESS) \
		--pull \
		--build-arg VERSION=$(Version) --build-arg GIT_COMMIT=$(GitCommit) \
		-t $(REGISTRY)/inlets/inlets:$(Version) .

.PHONY: docker-login
docker-login:
	echo -n "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin

.PHONY: docker-login-ghcr
docker-login-ghcr:
	echo -n "${GHCR_PASSWORD}" | docker login -u "${GHCR_USERNAME}" --password-stdin ghcr.io

.PHONY: push
push: OUTPUT=--push
push: docker

.PHONY: push-ghcr
push-ghcr: REGISTRY=ghcr.io
push-ghcr: OUTPUT=--push
push-ghcr: docker

# build local docker image for amd64
.PHONY: local-image
local-image: PLATFORMS=amd64
local-image: OUTPUT=--load
local-image: docker
