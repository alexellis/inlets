TAG?=latest
Version := $(shell git describe --tags --dirty)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-X main.Version=$(Version) -X main.GitCommit=$(GitCommit)"

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
	docker build --build-arg Version=$(Version) --build-arg GIT_COMMIT=$(GitCommit) -t alexellis2/inlets:$(TAG) .
