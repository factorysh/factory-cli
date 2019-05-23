GIT_VERSION?=$(shell git describe --tags --always --abbrev=42 --dirty)

build: bin vendor
	go build -ldflags "-X gitlab.bearstech.com/factory/gitlab-cli/version.version=$(GIT_VERSION)" \
	-o bin/factory-cli

bin:
	mkdir -p bin

vendor: dep

dep:
	dep ensure -v

test:
	go test -v \
		gitlab.bearstech.com/factory/factory-cli/client \
		gitlab.bearstech.com/factory/factory-cli/gitlab

docker-build:
	docker run --rm \
		-u `id -u` \
		-v ~/.cache:/.cache \
		-v `pwd`:/go/src/gitlab.bearstech.com/factory/factory-cli \
		-w /go/src/gitlab.bearstech.com/factory/factory-cli \
		bearstech/golang-dep \
		make build

docker-test:
	docker run --rm \
		-u `id -u` \
		-v ~/.cache:/.cache \
		-v `pwd`:/go/src/gitlab.bearstech.com/factory/factory-cli \
		-w /go/src/gitlab.bearstech.com/factory/factory-cli \
		bearstech/golang-dep \
		make test

clean:
	rm -rf bin vendor

