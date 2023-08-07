export GO111MODULE=on
export TF_ACC_TERRAFORM_VERSION=1.4.1
export SKIP_EXTERNAL_TABLE_TESTS=true
export SKIP_SCIM_INTEGRATION_TESTS=true

BASE_BINARY_NAME=terraform-provider-snowflake
TERRAFORM_PLUGINS_DIR=$(HOME)/.terraform.d/plugins
TERRAFORM_PLUGIN_LOCAL_INSTALL=$(TERRAFORM_PLUGINS_DIR)/$(BASE_BINARY_NAME)
COVERAGE_REPORT_FILE=coverage.txt
COVERAGE_FLAGS=-coverprofile=$(COVERAGE_REPORT_FILE) -covermode=atomic

help: ## display help for this makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

setup: ## setup development dependencies
	@which ./bin/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.53.3
	@which ./bin/reviewdog || curl -sSfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh | sh -s -- -b ./bin v0.14.2
.PHONY: setup

cleanup: ## cleanup development dependencies
	rm -rf bin/*
.PHONY: cleanup

sweep: ## destroy architecture; USE ONLY FOR DEVELOPMENT ACCOUNTS
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	@read -p "Are you sure? [y/n]" -n 1 REPLY; echo; \
		if [[ $$REPLY =~ ^[yY]$$ ]]; then \
		  SNOWFLAKE_ENABLE_SWEEP=1 go test -timeout 300s -run ^TestSweepAll ./pkg/sdk -v; \
		  else echo "Aborting..."; \
		fi;
.PHONY: sweep

lint-ci: ## run the fast go linters
	./bin/reviewdog -conf .reviewdog.yml -reporter=github-pr-review -tee -fail-on-error=true
.PHONY: lint-ci

test: ## run the tests (except sdk tests)
	CGO_ENABLED=1 go test -race $(COVERAGE_FLAGS) ./pkg/resources/...
	CGO_ENABLED=1 go test -race $(COVERAGE_FLAGS) ./pkg/provider/...
	CGO_ENABLED=1 go test -race $(COVERAGE_FLAGS) ./pkg/snowflake/...
.PHONY: test

test-acceptance: ## runs all tests, including the acceptance tests which create and destroys real resources
	SKIP_MANAGED_ACCOUNT_TEST=1 SKIP_EMAIL_INTEGRATION_TESTS=1 TF_ACC=1 go test -timeout 1200s -v $(COVERAGE_FLAGS) ./...
.PHONY: test-acceptance

build-local: ## build the binary locally
	go build -o $(BASE_BINARY_NAME) .
.PHONY: build-local

install-tf: build-local ## installs plugin where terraform can find it
	mkdir -p $(TERRAFORM_PLUGINS_DIR)
	cp ./$(BASE_BINARY_NAME) $(TERRAFORM_PLUGIN_LOCAL_INSTALL)
.PHONY: install-tf

uninstall-tf: ## uninstalls plugin from where terraform can find it
	rm -f $(TERRAFORM_PLUGIN_LOCAL_INSTALL)
.PHONY: uninstall-tf

clean: ## clean local binaries
	rm -f $(BASE_BINARY_NAME)
	go clean
	@#dist dir is removed because goreleaser puts artifacts there: https://goreleaser.com/customization/dist/
	rm -rf dist
.PHONY: clean

docs: ## generate docs for terraform plugin
	SNOWFLAKE_USER= SNOWFLAKE_ACCOUNT= go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
.PHONY: docs

docs-check: docs ## check that docs have been generated
	git diff --exit-code -- docs
.PHONY: docs-check

mod: ## add missing and remove unused modules
	go mod tidy -compat=1.20
.PHONY: mod

mod-check: mod ## check if there are any missing/unused modules
	git diff --exit-code -- go.mod go.sum
.PHONY: mod-check

fmt-check: ## Check formatting
	./bin/golangci-lint run ./... -v
.PHONY: fmt-check

fmt-fix: ## Check and fix formatting
	./bin/golangci-lint run ./... -v --fix
.PHONY: fmt-fix

pre-push: mod-check docs-check fmt-check; ## Run a few checks before pushing a change (docs, fmt, mod, etc.)
.PHONY: pre-push
