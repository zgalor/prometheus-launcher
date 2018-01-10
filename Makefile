DOCKER_REGISTRY := zgalor
DOCKER_IMAGE_NAME := prometheus-launcher

setup: ## Install dependencies
	dep ensure

build: ## Build executable
	mkdir -p ./bin
	go build -o ./bin/lnw ./cmd/launch-n-watch/main.go

container: build ## Build container image
	docker build --rm -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME) .

push: container ## Push container image to registry
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME)

clean: ## Clean intermediate files
	rm -rf ./bin

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build