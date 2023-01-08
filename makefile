PROJECT_NAME=restful-api-demo
MAIN_FILE=main.go
PKG := "github.com/jacknotes/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
PKG_CONFIG := "etc/restful-api.toml"

.PHONY: all dep lint vet test test-coverage build clean

all: build

dep: ## Get the dependencies
	@go mod tidy

lint: ## Lint Golang files
	@golint -set_exit_status ${PKG_LIST}

vet: ## Run go vet
	@go vet ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}

test-coverage: ## Run tests with coverage
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST} 
	@cat cover.out >> coverage.txt

build: dep ## Build the binary file
	@go build -ldflags "-s -w" -o dist/${PROJECT_NAME} $(MAIN_FILE)

linux: dep ## Build the binary file
	@GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${PROJECT_NAME} $(MAIN_FILE)

windows: dep ## Build the binary file
	@GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ${PROJECT_NAME} $(MAIN_FILE)

mac: dep ## Build the binary file
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ${PROJECT_NAME} $(MAIN_FILE)

freeBSD: dep ## Build the binary file
	@GOOS=freebsd GOARCH=amd64 go build -ldflags "-s -w" -o ${PROJECT_NAME} $(MAIN_FILE)

install: dep  ## install grpc gen tools
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/favadi/protoc-go-inject-tag@latest

gen: ## generate code # 对生成的protobuf go文件注入tag
	@protoc -I=. -I='/usr/local/include' --go_out=. --go_opt=module=${PKG} --go-grpc_out=. --go-grpc_opt=module=${PKG} apps/*/pb/*.proto
	@protoc-go-inject-tag -input="*apps/*/*.pb.go"	

run: ## Run Develop server
	@go run $(MAIN_FILE) start -f ${PKG_CONFIG}

clean: ## Remove previous build
	@rm -f dist/*
	@ rm -f cover.out coverage.txt 

push: # push git to multi repo
	@git push

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'