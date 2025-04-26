SHELL := bash

REPOSITORY ?= localhost
CONTAINER_NAME ?= go-wol
TAG ?= latest

# Build the binary
build:
	hack/build.sh

# Run the server on port 8080 to quickly test changes
run: build
	bin/go-wol server --log debug

# Build the container image
image:
	podman build -t $(REPOSITORY)/$(CONTAINER_NAME):$(TAG) .

# Run unit-tests with race detection and coverage
test:
	go test -v -race -coverprofile=coverprofile.out -coverpkg "./..." ./...

# Update project dependencies
update-deps:
	hack/update-deps.sh

# Generate coverage profile
coverprofile:
	hack/coverprofile.sh

# Run linter
lint:
	golangci-lint run -v

# Format the code
fmt:
	gofmt -s -w ./cmd ./pkg

# Validate that all generated files are up to date.
validate:
	hack/validate.sh

# Generate the bootstrap.css file
generate-bootstrap:
	hack/generate-bootstrap.sh

# Scan code for vulnerabilities using gosec
gosec:
	gosec ./...

# Clean up generated files
clean:
	rm -rf bin coverprofiles coverprofile.out

# Show this help message
help:
	@echo "Available targets:"
	@echo ""
	@awk '/^#/{c=substr($$0,3);next}c&&/^[[:alpha:]][[:alnum:]_-]+:/{print substr($$1,1,index($$1,":")),c}1{c=0}' $(MAKEFILE_LIST) | column -s: -t
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
