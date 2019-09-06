PROJECT_NAME := "GraphQL_Cello"
PKG_LIST := $(shell go list ${PROJECT_NAME}/...)

.PHONY: all dep build docker clean gofmt goreport goreport_deb test coverage coverhtml lint

all: dep build

copy_dir: ## Copy project folder to GOPATH
	@rm -rf $(GOPATH)/src/${PROJECT_NAME}
	@cp -Rp `pwd` $(GOPATH)/src/${PROJECT_NAME}

lint_dep: ## Get the dependencies for golint
	@$(GOROOT)/bin/go get -u golang.org/x/lint/golint
	@$(GOROOT)/bin/go install golang.org/x/lint/golint

lint: ## Lint the files
	@$(GOPATH)/bin/golint -set_exit_status ${PKG_LIST}

test: ## Run unittests
	@sudo -E $(GOROOT)/bin/go test -v ${PKG_LIST}

race: ## Run data race detector
	@sudo -E $(GOROOT)/bin/go test -race -v ${PKG_LIST}

coverage: ## Generate global code coverage report
	@sudo -E $(GOROOT)/bin/go test -v -coverprofile=coverage.out ${PKG_LIST}
	@$(GOROOT)/bin/go tool cover -func=coverage.out

coverhtml: coverage ## Generate global code coverage report in HTML
	@$(GOROOT)/bin/go tool cover -html=coverage.out

dep: ## Get the dependencies for build
	@$(GOROOT)/bin/go get -u github.com/nu7hatch/gouuid
	@$(GOROOT)/bin/go get -u github.com/go-sql-driver/mysql
	@$(GOROOT)/bin/go get -u github.com/graphql-go/graphql
	@$(GOROOT)/bin/go get -u github.com/graphql-go/handler

gofmt: ## Run gofmt for go files
	@find -name '*.go' -exec $(GOROOT)/bin/gofmt -s -w {} \;

goreport_dep: ## Get the dependencies for goreport
	@$(GOROOT)/bin/go get -u github.com/gojp/goreportcard/cmd/goreportcard-cli
	@$(GOROOT)/bin/go install github.com/gojp/goreportcard/cmd/goreportcard-cli
	@$(GOROOT)/bin/go get -u github.com/alecthomas/gometalinter
	@$(GOROOT)/bin/go install github.com/alecthomas/gometalinter
	@$(GOROOT)/bin/go get -u github.com/fzipp/gocyclo
	@$(GOROOT)/bin/go install github.com/fzipp/gocyclo
	@$(GOROOT)/bin/go get -u github.com/gordonklaus/ineffassign
	@$(GOROOT)/bin/go install github.com/gordonklaus/ineffassign
	@$(GOROOT)/bin/go get -u github.com/client9/misspell/cmd/misspell
	@$(GOROOT)/bin/go install github.com/client9/misspell/cmd/misspell

goreport: goreport_dep ## Make goreport
	@git submodule sync --recursive
	@git submodule update --init --recursive
	@./graphql_cello_badge/update_goreport_grade.sh

build: ## Build the binary file
	@$(GOROOT)/bin/go build -o $(PROJECT_NAME) main.go

docker: ## Build docker image and push it to private docker registry
	@sudo docker build -t graphql_cello .
	@sudo docker tag graphql_cello:latest 192.168.110.230:5000/graphql_cello:latest
	@sudo docker push 192.168.110.230:5000/graphql_cello:latest

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

