# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOVET=$(GOCMD) vet
GOFMT=gofmt
GOLINT=golint
BINARY_NAME=user-microservice

all: clean build lint tool test

test:
	$(GOTEST) -short -mod=vendor ./... 

test-integration: 
	export GIN_MODE=release \
	&& $(GOTEST) -short -tags=integration  -mod=vendor ./...

build: ensure
	$(GOBUILD) -mod=vendor . 

build-linux: ensure
	GOOS=linux $(GOBUILD) -mod=vendor . 

tool:
	$(GOVET) ./...; true
	$(GOFMT) -w .

coverage:
	scripts/coverage.sh

clean:
	go clean -i .
	rm -rf docs
	rm -rf vendor
	rm -f $(BINARY_NAME)

lint:
	GO111MODULE=off $(GOGET) golang.org/x/lint/golint
	$(GOLINT) -set_exit_status $($(GOCMD) list ./... | grep -v /vendor/)

generate-swagger:
	GO111MODULE=off $(GOGET) github.com/swaggo/swag/cmd/swag
	swag init

ensure: generate-swagger
	go mod vendor

docker-build: build-linux
	docker build -t $(BINARY_NAME) .

run-as-lambda: build-linux
	sam local start-api

continuous-integration: ensure build lint tool test-integration coverage


help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make lint: golint ./..."
	@echo "make clean: remove object files and cached files"
	@echo "make ensure: get the deployment tools"
	@echo "make coverage: get the coverage of my files"
	@echo "make docker-build: build a docker image and run the container"
