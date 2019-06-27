GIT_VERSION?=$(shell git describe --tags --always --abbrev=42 --dirty)

build: bin vendor
	go build -ldflags "-X gitlab.bearstech.com/factory/gitlab-cli/version.version=$(GIT_VERSION)" \
	-o bin/factory \
	main.go

bin:
	mkdir -p bin

vendor: dep

dep:
	dep ensure -v

test:
	go test -v \
		github.com/factorysh/factory-cli/client \
		github.com/factorysh/factory-cli/gitlab

docker-build:
	docker run --rm \
		-u `id -u` \
		-v ~/.cache:/.cache \
		-v `pwd`:/go/src/github.com/factorysh/factory-cli \
		-w /go/src/github.com/factorysh/factory-cli \
		bearstech/golang-dep \
		make build

docker-test:
	docker run --rm \
		-u `id -u` \
		-v ~/.cache:/.cache \
		-v `pwd`:/go/src/github.com/factorysh/factory-cli \
		-w /go/src/github.com/factorysh/factory-cli \
		bearstech/golang-dep \
		make test

test-sftp: docker-build
	echo get ./data/volume/test /tmp/test | \
		PRIVATE_TOKEN=$(PRIVATE_TOKEN) ./bin/factory \
		volume sftp \
		-p factory/factory-canary -e staging

test-exec: docker-build
	PRIVATE_TOKEN=$(PRIVATE_TOKEN) ./bin/factory \
		container exec \
		-p factory/factory-canary -e staging \
		web -- ls -l

clean:
	rm -rf bin vendor
