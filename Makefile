TAG?=latest

.PHONY: all
all: clean build

.PHONY: build
build:
	./build.sh

.PHONY: build_redist
build_redist:
	./build_redist.sh

.PHONY: ci-armhf-push
ci-armhf-push:
	(docker push alexellis2/inlets:$(TAG)-armhf)

.PHONY: ci-armhf-build
ci-armhf-build:
	(./build.sh $(TAG)-armhf)

.PHONY: ci-arm64-push
ci-arm64-push:
	(docker push alexellis2/inlets:$(TAG)-arm64)

.PHONY: ci-arm64-build
ci-arm64-build:
	(./build.sh $(TAG)-arm64)

.PHONY: clean
clean:
	(rm -f inlets inlets-arm64 inlets-armhf inlets-darwin)
