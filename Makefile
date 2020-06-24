.PHONY: upload_dists

GIT_VERSION?=$(shell git describe --tags --always --abbrev=42 --dirty)

GOOS?=linux
GOARCH?=amd64
export COMPOSE=docker-compose -f docker-compose.yml

binary: bin
	go build \
	-ldflags " \
	-X github.com/factorysh/factory-cli/version.version=$(GIT_VERSION) \
	-X github.com/factorysh/factory-cli/version.os=$(GOOS) \
	-X github.com/factorysh/factory-cli/version.arch=$(GOARCH) \
	" \
	-o bin/factory \
	main.go

build/factory-$(GOOS)-$(GOARCH)-$(GIT_VERSION):
	env GOOS=$(GOOS) GOARCH=$(GOARCH) \
	go build \
	-ldflags " \
	-X github.com/factorysh/factory-cli/version.version=$(GIT_VERSION) \
	-X github.com/factorysh/factory-cli/version.os=$(GOOS) \
	-X github.com/factorysh/factory-cli/version.arch=$(GOARCH) \
	" \
	-o build/factory-$(GOOS)-$(GOARCH)-$(GIT_VERSION) \
	main.go

dist/factory-$(GOOS)-$(GOARCH)-$(GIT_VERSION).gz: build/factory-$(GOOS)-$(GOARCH)-$(GIT_VERSION)
	gzip -c build/factory-$(GOOS)-$(GOARCH)-$(GIT_VERSION) > \
			dist/factory-$(GOOS)-$(GOARCH)-$(GIT_VERSION).gz

release: build dist dist/factory-$(GOOS)-$(GOARCH)-$(GIT_VERSION).gz

releases:
	make release GOOS=windows GOARCH=amd64
	make release GOOS=linux GOARCH=amd64
	make release GOOS=darwin GOARCH=amd64

venv:
	python3 -m venv venv
	./venv/bin/pip install requests

upload_dists: venv
	./venv/bin/python upload_dists

new_tag:
	# check invalid tag name
	echo $(GIT_VERSION) | grep -E "^v[0-9].[0-9].[0-9]" || false
	# fail if tag exists
	git tag | grep $(GIT_VERSION) || true
	git tag $(GIT_VERSION)
	git checkout $(GIT_VERSION)

new_release: new_tag docker-binaries upload_dists
	git checkout master

bin:
	mkdir -p bin

dist:
	mkdir -p dist

build:
	mkdir -p build

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
		-e GIT_VERSION=$(GIT_VERSION) \
		bearstech/golang-dep \
		make binary

docker-binaries:
	docker run --rm \
		-u `id -u` \
		-v ~/.cache:/.cache \
		-v `pwd`:/go/src/github.com/factorysh/factory-cli \
		-w /go/src/github.com/factorysh/factory-cli \
		-e GIT_VERSION=$(GIT_VERSION) \
		bearstech/golang-dep \
		make releases

docker-test:
	docker run --rm \
		-u `id -u` \
		-v ~/.cache:/.cache \
		-v `pwd`:/go/src/github.com/factorysh/factory-cli \
		-w /go/src/github.com/factorysh/factory-cli \
		bearstech/golang-dep \
		make test

docker-test-runjob: docker-build
	bin/factory runjob -h
	bin/factory runjob -D job1
	bin/factory runjob -v job1
	bin/factory runjob -D job2
	bin/factory runjob job2
	bin/factory runjob -D job3
	bin/factory runjob job3
	bin/factory runjob notfound || true

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
	rm -rf bin vendor build dist mysql-*.gz
