.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: cover
cover: ## Generate coverage report
	go test -trimpath -coverprofile=.coverprofile -covermode=count ./...

.PHONY: cover/func
cover/func: ## Generate coverage report by function
	go tool cover -func .coverprofile

.PHONY: cover/html
cover/html: ## Generate HTML coverage report
	go tool cover -html .coverprofile

.PHONY: lint
lint: ## Lint code
	golangci-lint run -c ./../.golangci.yaml



# Makefile config
#===============================================================
help: ## Display this help screen
	echo "Usage: make [task]\n\nTasks:"
	perl -nle 'printf("    \033[33m%-30s\033[0m %s\n",$$1,$$2) if /^([a-zA-Z0-9_\/-]*?):(?:.+?## )?(.*?)$$/' $(MAKEFILE_LIST)

.SILENT: help

.PHONY: $(shell egrep -o '^(\._)?[a-z_-]+:' $(MAKEFILE_LIST) | sed 's/://')
