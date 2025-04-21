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
	$(NULL)
