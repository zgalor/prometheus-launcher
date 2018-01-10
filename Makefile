SOURCE_FILES?=$$(go list ./... | grep -v /vendor/)
DOCKER_REGISTRY := zgalor
DOCKER_IMAGE_NAME := prometheus-launcher

setup: ## Install all the build and lint dependencies
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/golang/dep/...
	dep ensure
	gometalinter --install --update

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

lint: ## Run all the linters
	gometalinter --vendor --disable-all \
		--enable=deadcode \
		--enable=ineffassign \
		--enable=gosimple \
		--enable=staticcheck \
		--enable=gofmt \
		--enable=goimports \
		--enable=dupl \
		--enable=misspell \
		--enable=errcheck \
		--enable=vet \
		--enable=vetshadow \
		--deadline=10m \
		./...

build: ## Build a beta version
	mkdir -p ./bin
	go build -o ./bin/lnw ./cmd/launch-n-watch/main.go

container: build ## Build container image
	docker build --rm -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME) .

push: container ## Push container image to registry
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME)

clean: # Clean files
	rm -rf ./bin

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build