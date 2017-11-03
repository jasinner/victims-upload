# Passed into cmd/main.go at build time
VERSION := $(shell cat ./VERSION)
COMMIT_HASH := $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date +%s)

# Used in tagging images
IMAGE_VERSION_TAG := victims-upload:$(VERSION)
IMAGE_DATE_TAG := victims-upload:$(BUILD_TIME)

# Used during all builds
LDFLAGS := -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH} -X main.buildTime=${BUILD_TIME}

.PHONY: help clean victims-upload image

default: help

help:
	@echo "Targets:"
	@echo " deps: Install dependencies with govendor"
	@echo "	victims-upload: Builds a victims-upload binary"
	@echo "	clean: cleans up and removes built files"
	@echo "	image: builds a container image"

deps:
	go get github.com/kardianos/govendor
	govendor sync

victims-upload:
	govendor build -ldflags '${LDFLAGS}' -o victims-upload cmd/main.go

static-victims-upload:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 govendor build --ldflags '-extldflags "-static" ${LDFLAGS}' -a -o victims-upload cmd/main.go

clean:
	go clean
	rm -f victims-upload

image: clean deps static-victims-api
	sudo docker build -t $(IMAGE_VERSION_TAG) -t $(IMAGE_DATE_TAG) .

gofmt:
	gofmt -l api/ cmd/

golint:
	go get github.com/golang/lint/golint
	golint api/ cmd/

lint: gofmt golint
