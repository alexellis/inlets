Version := $(shell git describe --tags --dirty)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)"

# docker manifest command will work with Docker CLI 18.03 or newer
# but for now it's still experimental feature so we need to enable that
export DOCKER_CLI_EXPERIMENTAL=enabled

.PHONY: all
all: docker

.PHONY: dist
dist:
	CGO_ENABLED=0 GOOS=linux go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/inlets
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/inlets-darwin
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/inlets-armhf
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/inlets-arm64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/inlets.exe

.PHONY: docker
docker:
	docker build --build-arg VERSION=$(Version) --build-arg GIT_COMMIT=$(GitCommit) -t alexellis2/inlets:$(Version)-amd64 .
	docker build --build-arg VERSION=$(Version) --build-arg GIT_COMMIT=$(GitCommit) --build-arg OPTS="GOARCH=arm64" -t alexellis2/inlets:$(Version)-arm64 .
	docker build --build-arg VERSION=$(Version) --build-arg GIT_COMMIT=$(GitCommit) --build-arg OPTS="GOARCH=arm GOARM=6" -t alexellis2/inlets:$(Version)-armhf .

.PHONY: docker-login
docker-login:
	echo -n "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin

.PHONY: push
push:
	docker push alexellis2/inlets:$(Version)-amd64
	docker push alexellis2/inlets:$(Version)-arm64
	docker push alexellis2/inlets:$(Version)-armhf

.PHONY: manifest
manifest:
	docker manifest create --amend alexellis2/inlets:$(Version) alexellis2/inlets:$(Version)-amd64 alexellis2/inlets:$(Version)-arm64 alexellis2/inlets:$(Version)-armhf
	docker manifest annotate alexellis2/inlets:$(Version) alexellis2/inlets:$(Version)-arm64 --os linux --arch arm64
	docker manifest annotate alexellis2/inlets:$(Version) alexellis2/inlets:$(Version)-armhf --os linux --arch arm --variant v6
	docker manifest push alexellis2/inlets:$(Version)
