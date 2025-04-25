SHELL := bash

REPOSITORY ?= localhost
CONTAINER_NAME ?= go-wol
TAG ?= latest

build: ## Build the binary
	hack/build.sh

run: build ## Run the server on port 8080 to quickly test changes
	bin/go-wol server --log debug

image: ## Build the container image
	podman build -t $(REPOSITORY)/$(CONTAINER_NAME):$(TAG) .

test: ## Run unit-tests with race detection and coverage
	go test -v -race -coverprofile=coverprofile.out -coverpkg "./..." ./...

update-deps: ## Update project dependencies
	hack/update-deps.sh

coverprofile: ## Generate coverage profile
	hack/coverprofile.sh

lint: ## Run linter
	golangci-lint run -v

fmt: ## Format the code
	gofmt -s -w ./cmd ./pkg

validate: ## Validate that all generated files are up to date.
	hack/validate.sh

generate-bootstrap: ## Generate the bootstrap.css file
	hack/generate-bootstrap.sh

gosec: ## Scan code for vulnerabilities using gosec
	gosec ./...

clean: ## Clean up generated files
	rm -rf bin coverprofiles coverprofile.out

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "%-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "Run 'make <target>' to execute a specific target."

.PHONY: \
	build \
	run \
	image \
	test \
	update-deps \
	coverprofile \
	lint \
	fmt \
	validate \
	generate-bootstrap \
	gosec \
	clean \
	help \
	$(NULL)
