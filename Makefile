Version := $(shell git describe --tags --dirty)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)"

BINARIES := inlets inlets-darwin inlets-armhf inlets-arm64
ALL_BINARIES := $(addprefix bin/,$(BINARIES))

PREFIX := /usr/local
BINDIR := $(PREFIX)/bin

.PHONY: all
all: docker

.PHONY: clean
clean:
	$(RM) -r bin/

.PHONY: build
build: $(ALL_BINARIES)

# TODO: kept for backwards compat, but is now redundant with the build target
.PHONY: dist
dist: $(ALL_BINARIES)

.PHONY: install
install: bin/inlets
	install -d $(BINDIR)
	install -p -m 0755 $^ $(BINDIR)

.PHONY: uninstall
uninstall:
	$(RM) $(PREFIX)$(BINDIR)/inlets

bin/inlets:
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o $@

bin/inlets-darwin:
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o $@

bin/inlets-armhf:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o $@

bin/inlets-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/inlets-arm64

.PHONY: docker
docker:
	docker build --build-arg Version=$(Version) --build-arg GIT_COMMIT=$(GitCommit) -t alexellis2/inlets:$(Version) .

.PHONY: docker-login
docker-login:
	echo -n "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin

.PHONY: push
push:
	docker push alexellis2/inlets:$(Version)
