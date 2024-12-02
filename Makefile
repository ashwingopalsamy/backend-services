APP_NAME = backend-services
DOCKER_IMAGE = $(APP_NAME):latest
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

all: docker-build

dep:
	go mod tidy
	go mod vendor

docker-build: dep
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .
	@echo "Running Docker container..."
	docker run -p 8080:8080 --rm $(DOCKER_IMAGE)

clean:
	rm -rf bin vendor

# Run tests
test:
	go test ./... -v
