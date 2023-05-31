GOOS ?= linux
GOARCH ?= amd64

IMAGE_REPO_SERVER ?= prasadg193/sample-csi-cbt-service
IMAGE_TAG_SERVER ?= latest


.PHONY: build
build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build  -o grpc-server ./cmd/server/main.go

image: build
	docker buildx build --platform=linux/amd64 -t $(IMAGE_REPO_SERVER):$(IMAGE_TAG_SERVER) -f Dockerfile-grpc .

push:
	docker push $(IMAGE_REPO_SERVER):$(IMAGE_TAG_SERVER)
