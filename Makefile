PROGRAM_NAME = ftpdts

COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
TAG=$(shell git describe --tags |cut -d- -f1)

LDFLAGS = -ldflags "-X main.gitTag=${TAG} -X main.gitCommit=${COMMIT} -X main.gitBranch=${BRANCH}"

.PHONY: help clean dep build install uninstall

.DEFAULT_GOAL := help

help: ## Display this help screen.
	@echo "Makefile available targets:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  * \033[36m%-15s\033[0m %s\n", $$1, $$2}'

dep: ## Download the dependencies.
	go mod download

build: dep ## Build executable.
	mkdir -p ./bin
	CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${PROGRAM_NAME} ./src

clean: ## Clean build directory.
	rm -f ./bin/${PROGRAM_NAME}
	rmdir ./bin

#golangci-lint and gosec must be installed, see details:
#https://golangci-lint.run/usage/install/#local-installation
#https://github.com/securego/gosec
lint: dep ## Lint the source files
	golangci-lint run --timeout 5m -E golint
	gosec -quiet ./...

test: dep ## Run tests
	go test -race -p 1 -timeout 300s -coverprofile=.test_coverage.txt ./... && \
	go tool cover -func=.test_coverage.txt | tail -n1 | awk '{print "Total test coverage: " $$3}'
	@rm .test_coverage.txt

docker-build: build ## Build docker image
	docker build -t starshiptroopers/${PROGRAM_NAME}:latest -t starshiptroopers/${PROGRAM_NAME}:${TAG} .
	docker image prune --force --filter label=stage=intermediate

docker-image: docker-build ## Save docker image to file
	rm -f ./build/*.docker.tar
	docker save -o build/${PROGRAM_NAME}:${TAG}.docker.tar starshiptroopers/${PROGRAM_NAME}:latest

docker-app: docker-image ## Build application tar with docker image
	mkdir -p ./build/config
	mkdir -p ./build/tmpl
	mkdir -p ./build/logs
	cp *.ini ./build/config/
	cp tmpl/default.tmpl ./build/tmpl
	cp docker-compose.yml ./build/
	echo "Load the docker image:\n> docker load -i ${PROGRAM_NAME}:${TAG}.docker.tar\nRun the container with:\n> docker-compose up" > ./build/readme.txt
	mv build ${PROGRAM_NAME}-${TAG}
	tar -c ${PROGRAM_NAME}-${TAG} | gzip > bin/${PROGRAM_NAME}-${TAG}.tar.gz
	mv ${PROGRAM_NAME}-${TAG} build
	echo "Application archive has been saved to bin/${PROGRAM_NAME}-${TAG}.tar.gz"


