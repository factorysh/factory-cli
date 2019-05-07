GIT_VERSION?=$(shell git describe --tags --always --abbrev=42 --dirty)

build: bin vendor
	go build -ldflags "-X gitlab.bearstech.com/factory/gitlab-cli/version.version=$(GIT_VERSION)" \
	-o bin/factory

bin:
	mkdir -p bin

vendor: dep

dep:
	dep ensure

clean:
	rm -rf bin vendor