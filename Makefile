SHELL := bash

REPOSITORY ?= localhost
CONTAINER_NAME ?= go-wol
TAG ?= latest

build:
	hack/build.sh

run: build
	bin/go-wol server --log debug

image:
	podman build -t $(REPOSITORY)/$(CONTAINER_NAME):$(TAG) .

test:
	go test -v -race -coverprofile=coverprofile.out -coverpkg "./..." ./...

update-deps:
	hack/update-deps.sh

coverprofile:
	hack/coverprofile.sh

lint:
	golangci-lint run -v

fmt:
	gofmt -s -w ./cmd ./pkg

validate:
	hack/validate.sh

generate-bootstrap:
	hack/generate-bootstrap.sh

clean:
	rm -rf bin coverprofiles coverprofile.out

help:
	@echo "Available targets:"
	@echo "  build               Build the project"
	@echo "  run                 Build and run the server"
	@echo "  image               Build the container image"
	@echo "  test                Run tests with coverage"
	@echo "  update-deps         Update project dependencies"
	@echo "  coverprofile        Generate coverage profile"
	@echo "  lint                Run linter"
	@echo "  fmt                 Format the code"
	@echo "  validate            Validate the project"
	@echo "  generate-bootstrap  Generate bootstrap files"
	@echo "  clean               Clean up generated files"

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
	clean \
	help \
	$(NULL)
