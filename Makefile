export TEST_SF_TF_SKIP_SAML_INTEGRATION_TEST=true
export TEST_SF_TF_SKIP_MANAGED_ACCOUNT_TEST=true
export BASE_BINARY_NAME=terraform-provider-snowflake
export TERRAFORM_PLUGINS_DIR=$(HOME)/.terraform.d/plugins
export TERRAFORM_PLUGIN_LOCAL_INSTALL=$(TERRAFORM_PLUGINS_DIR)/$(BASE_BINARY_NAME)
export LATEST_GIT_TAG=$(shell git tag --sort=-version:refname | head -n 1)
export CURRENT_OS := $(shell uname -s)
export CURRENT_ARCH := $(shell arch)

default: help

dev-setup: ## setup development dependencies
# TODO(SNOW-2002208): Upgrade to the latest version of golangci-lint.
	@which ./bin/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.64.8
	cd tools && mkdir -p bin/
	cd tools && env GOBIN=$$PWD/bin go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	cd tools && env GOBIN=$$PWD/bin go install mvdan.cc/gofumpt

dev-cleanup: ## cleanup development dependencies
	rm -rf bin/*
	rm -rf tools/bin/*

docs: generate-docs-additional-files ## generate docs
	tools/bin/tfplugindocs generate

docs-check: docs ## check that docs have been generated
	git diff --exit-code -- docs

fmt: terraform-fmt ## Run terraform fmt and gofumpt
	tools/bin/gofumpt -l -w .

terraform-fmt: ## Run terraform fmt
	terraform fmt -recursive ./examples/
	terraform fmt -recursive ./pkg/resources/testdata/
	terraform fmt -recursive ./pkg/datasources/testdata/

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-23s\033[0m %s\n", $$1, $$2}'

install: ## install the binary
	go install -v ./...

lint: # Run static code analysis, check formatting. See https://golangci-lint.run/
	./bin/golangci-lint run ./... -v

lint-fix: ## Run static code analysis, check formatting and try to fix findings
	./bin/golangci-lint run ./... -v --fix

mod: ## add missing and remove unused modules
	go mod tidy -compat=1.23.6

mod-check: mod ## check if there are any missing/unused modules
	git diff --exit-code -- go.mod go.sum

pre-push: generate-all-config-model-builders-check mod fmt generate-docs-additional-files docs lint test-architecture ## Run a few checks before pushing a change (docs, fmt, mod, etc.)

pre-push-check: pre-push mod-check generate-docs-additional-files-check docs-check ## Run checks before pushing a change (docs, fmt, mod, etc.)

sweep: ## destroy the whole architecture; USE ONLY FOR DEVELOPMENT ACCOUNTS
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	@echo "Are you sure? [y/n]" >&2
	@read -r REPLY; \
		if echo "$$REPLY" | grep -qG "^[yY]$$"; then \
			TEST_SF_TF_ENABLE_SWEEP=1 go test -timeout=10m -run "^(TestSweepAll|Test_Sweeper_NukeStaleObjects)" ./pkg/sdk -v; \
			else echo "Aborting..."; \
		fi;

test: ## run unit and integration tests
	go test -v -cover -timeout=60m ./...

test-acceptance: ## run acceptance tests
	TF_ACC=1 SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE=true TEST_SF_TF_REQUIRE_TEST_OBJECT_SUFFIX=1 TEST_SF_TF_REQUIRE_GENERATED_RANDOM_VALUE=1 SF_TF_ACC_TEST_ENABLE_ALL_PREVIEW_FEATURES=true go test -run "^TestAcc_" -v -cover -timeout=120m ./...

test-account-level-features: ## run integration and acceptance test modifying account
	TF_ACC=1 SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE=true TEST_SF_TF_REQUIRE_TEST_OBJECT_SUFFIX=1 TEST_SF_TF_REQUIRE_GENERATED_RANDOM_VALUE=1 SF_TF_ACC_TEST_ENABLE_ALL_PREVIEW_FEATURES=true go test --tags=account_level_tests -run "^(TestAcc_|TestInt_)" -v -cover -timeout=30m ./...

test-integration: ## run SDK integration tests
	TEST_SF_TF_REQUIRE_TEST_OBJECT_SUFFIX=1 TEST_SF_TF_REQUIRE_GENERATED_RANDOM_VALUE=1 go test -run "^TestInt_" -v -cover -timeout=60m ./...

test-architecture: ## check architecture constraints between packages
	go test ./pkg/architests/... -v

test-object-renaming: ## runs tests in object_renaming_acceptance_test.go
	TEST_SF_TF_ENABLE_OBJECT_RENAMING=1 go test ./pkg/resources/object_renaming_acceptace_test.go -v

test-acceptance-%: ## run acceptance tests for the given resource only, e.g. test-acceptance-Warehouse
	TF_ACC=1 TF_LOG=DEBUG SNOWFLAKE_DRIVER_TRACING=debug SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE=true SF_TF_ACC_TEST_ENABLE_ALL_PREVIEW_FEATURES=true go test -run ^TestAcc_$*_ -v -timeout=20m ./pkg/resources

build-local: ## build the binary locally
	go build -o $(BASE_BINARY_NAME) .

install-tf: build-local ## installs plugin where terraform can find it
	mkdir -p $(TERRAFORM_PLUGINS_DIR)
	cp ./$(BASE_BINARY_NAME) $(TERRAFORM_PLUGIN_LOCAL_INSTALL)

release-local: ## use GoReleaser to build the binary locally for the current OS and ARCH
	goreleaser build --clean --skip=validate --single-target

release-local-all: ## use GoReleaser to build the binary locally
	goreleaser build --clean --skip=validate

install-locally-released-tf: release-local ## installs plugin (built by the GoReleaser) where terraform can find it
	mkdir -p $(TERRAFORM_PLUGINS_DIR)
	cp ./dist/terraform-provider-snowflake_$(CURRENT_OS)_$(CURRENT_ARCH)/terraform-provider-snowflake_$(LATEST_GIT_TAG) $(TERRAFORM_PLUGIN_LOCAL_INSTALL)

uninstall-tf: ## uninstalls plugin from where terraform can find it
	rm -f $(TERRAFORM_PLUGIN_LOCAL_INSTALL)

generate-all-dto: ## Generate all DTOs for SDK interfaces
	go generate ./pkg/sdk/*_dto.go

generate-dto-%: ./pkg/sdk/%_dto.go ## Generate DTO for given SDK interface
	go generate $<

clean-generator-poc:
	rm -f ./pkg/sdk/poc/example/*_gen.go
	rm -f ./pkg/sdk/poc/example/*_gen_test.go

clean-generator-%: ## Clean generated files for specified resource
	rm -f ./pkg/sdk/$**_gen.go
	rm -f ./pkg/sdk/$**_gen_*test.go

run-generator-poc:
	go generate ./pkg/sdk/poc/example/*_def.go
	go generate ./pkg/sdk/poc/example/*_dto_gen.go

run-generator-%: ./pkg/sdk/%_def.go ## Run go generate on given object definition
	go generate $<
	go generate ./pkg/sdk/$*_dto_gen.go

generate-docs-additional-files: ## generate docs additional files
	go run ./pkg/internal/tools/doc-gen-helper/ $$PWD

generate-docs-additional-files-check: generate-docs-additional-files ## check that docs additional files have been generated
	git diff --exit-code -- examples/additional

generate-show-output-schemas: ## Generate show output schemas with mappers
	go generate ./pkg/schemas/generate.go

clean-show-output-schemas: ## Clean generated show output schemas
	rm -f ./pkg/schemas/*_gen.go

generate-snowflake-object-assertions: ## Generate snowflake object assertions
	go generate ./pkg/acceptance/bettertestspoc/assert/objectassert/generate.go

clean-snowflake-object-assertions: ## Clean snowflake object assertions
	rm -f ./pkg/acceptance/bettertestspoc/assert/objectassert/*_gen.go

generate-snowflake-object-parameters-assertions: ## Generate snowflake object parameters assertions
	go generate ./pkg/acceptance/bettertestspoc/assert/objectparametersassert/generate.go

clean-snowflake-object-parameters-assertions: ## Clean snowflake object parameters assertions
	rm -f ./pkg/acceptance/bettertestspoc/assert/objectparametersassert/*_gen.go

generate-resource-assertions: ## Generate resource assertions
	go generate ./pkg/acceptance/bettertestspoc/assert/resourceassert/generate.go

clean-resource-assertions: ## Clean resource assertions
	rm -f ./pkg/acceptance/bettertestspoc/assert/resourceassert/*_gen.go

generate-resource-parameters-assertions: ## Generate resource parameters assertions
	go generate ./pkg/acceptance/bettertestspoc/assert/resourceparametersassert/generate.go

clean-resource-parameters-assertions: ## Clean resource parameters assertions
	rm -f ./pkg/acceptance/bettertestspoc/assert/resourceparametersassert/*_gen.go

generate-resource-show-output-assertions: ## Generate resource parameters assertions
	go generate ./pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert/generate.go

clean-resource-show-output-assertions: ## Clean resource parameters assertions
	rm -f ./pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert/*_gen.go

generate-resource-model-builders: ## Generate resource model builders
	go generate ./pkg/acceptance/bettertestspoc/config/model/generate.go

clean-resource-model-builders: ## Clean resource model builders
	rm -f ./pkg/acceptance/bettertestspoc/config/model/*_gen.go

generate-provider-model-builders: ## Generate provider model builders
	go generate ./pkg/acceptance/bettertestspoc/config/providermodel/generate.go

clean-provider-model-builders: ## Clean provider model builders
	rm -f ./pkg/acceptance/bettertestspoc/config/providermodel/*_gen.go

generate-datasource-model-builders: ## Generate datasource model builders
	go generate ./pkg/acceptance/bettertestspoc/config/datasourcemodel/generate.go

clean-datasource-model-builders: ## Clean datasource model builders
	rm -f ./pkg/acceptance/bettertestspoc/config/datasourcemodel/*_gen.go

clean-all-config-model-builders: clean-resource-model-builders clean-datasource-model-builders clean-provider-model-builders ## clean all generated config model builders

generate-all-config-model-builders: generate-resource-model-builders generate-datasource-model-builders generate-provider-model-builders ## generate all config model builders

generate-all-config-model-builders-check: clean-all-config-model-builders generate-all-config-model-builders ## check that generated config model builders are up-to-date
	git diff --exit-code -- pkg/acceptance/bettertestspoc/config/model
	git diff --exit-code -- pkg/acceptance/bettertestspoc/config/datasourcemodel
	git diff --exit-code -- pkg/acceptance/bettertestspoc/config/providelmodel

clean-all-assertions-and-config-models: clean-snowflake-object-assertions clean-snowflake-object-parameters-assertions clean-resource-assertions clean-resource-parameters-assertions clean-resource-show-output-assertions clean-resource-model-builders clean-provider-model-builders clean-datasource-model-builders ## clean all generated assertions and config models

generate-all-assertions-and-config-models: generate-snowflake-object-assertions generate-snowflake-object-parameters-assertions generate-resource-assertions generate-resource-parameters-assertions generate-resource-show-output-assertions generate-resource-model-builders generate-provider-model-builders generate-datasource-model-builders ## generate all assertions and config models

.PHONY: build-local clean-generator-poc dev-setup dev-cleanup docs docs-check fmt fmt-check fumpt help install lint lint-fix mod mod-check pre-push pre-push-check sweep test test-acceptance uninstall-tf
