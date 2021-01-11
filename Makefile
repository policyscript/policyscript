.PHONY: help
help: ## Prints out all make commands and their actions
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: clean ## Build golang binaries
	@echo \\nGenerating function binaries...
	@mkdir bin
	@for function in `ls functions`; do \
		GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$$function functions/$$function/main.go ; \
		echo "  -> Generated: bin/$$function" ; \
	done

.PHONY: clean
clean: ## Cleans build binaries
	@echo \\nRemoving build files...
	@rm -rf ./bin
	@echo "  -> Removed: ./bin"
	@rm -rf ./zip
	@echo "  -> Removed: ./zip"

.PHONY: deepclean
deepclean: clean ## Calls clean plus removes all git ignored Terraform files
	@echo \\nRemoving local Terraform files and folders...
	@rm -rf ./.terraform
	@echo "  -> Removed: ./.terraform"
	@rm terraform.tfstate.backup
	@echo "  -> Removed: terraform.tfstate.backup"

.PHONY: test
test: ## Run tests in the package
	@echo \\nTesting with Ginko...
	@go run github.com/onsi/ginkgo/ginkgo ./...

.PHONY: init
init: ## Generates a test at the directory you choose
	cd $(filter-out $@,$(MAKECMDGOALS)) ; \
	go run github.com/onsi/ginkgo/ginkgo bootstrap ; \
	go run github.com/onsi/ginkgo/ginkgo generate ;

# Does nothing when job doesn't match, rather than throwing error. This supports
# `init` allowing a second argument for path.
# https://stackoverflow.com/a/47008498/10577245
%:
    @: