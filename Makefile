GIT_VERSION?=$(shell git describe --tags --always --abbrev=42 --dirty)

GOOS?=linux
GOARCH?=amd64

build: bin vendor
	go build -ldflags "-X github.com/factorysh/factory-cli/version.version=$(GIT_VERSION)" \
	-o bin/factory \
	main.go

binaries: vendor
	make binary GOOS=windows GOARCH=amd64
	make binary GOOS=linux GOARCH=amd64
	make binary GOOS=darwin GOARCH=amd64

build/factory-$(GOOS)-$(GOARCH)-$(GIT_VERSION):
	env GOOS=$(GOOS) GOARCH=$(GOARCH) \
	go build -ldflags "-X github.com/factorysh/factory-cli/version.version=$(GIT_VERSION)" \
	-o build/factory-$(GOOS)-$(GOARCH)-$(GIT_VERSION) \
	main.go

dist/factory-$(GOOS)-$(GOARCH)-$(GIT_VERSION).gz: build/factory-$(GOOS)-$(GOARCH)-$(GIT_VERSION)
	gzip -c build/factory-$(GOOS)-$(GOARCH)-$(GIT_VERSION) > \
			dist/factory-$(GOOS)-$(GOARCH)-$(GIT_VERSION).gz

binary: build dist dist/factory-$(GOOS)-$(GOARCH)-$(GIT_VERSION).gz

bin:
	mkdir -p bin

dist:
	mkdir -p dist

build:
	mkdir -p build

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

docker-binaries:
	docker run --rm \
		-u `id -u` \
		-v ~/.cache:/.cache \
		-v `pwd`:/go/src/github.com/factorysh/factory-cli \
		-w /go/src/github.com/factorysh/factory-cli \
		bearstech/golang-dep \
		make binaries

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
	rm -rf bin vendor build dist
