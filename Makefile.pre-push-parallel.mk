# Parallel pre-push graph. The edges in this file are not on the targets in Makefile.
#
#   make pre-push-parallel
#   make -f Makefile.pre-push-parallel.mk -j pre-push-parallel-run
#
# Group 1, no edges: generate-sdk, generate-sdk-examples, generate-resource.
# Group 2 waits for generate-sdk: show-output schemas, object assertions,
# object parameter assertions, resource parameter assertions, show-output assertions.
# Group 3 waits for generate-sdk and, separately, for generate-resource.
# generate-resource does not wait for generate-sdk. The nodes are resource
# assertions, resource/datasource/provider model builders, docs additional
# files, and issue labels. Show-output schemas is the sdk side of the wait.
# Group 4: docs waits for docs additional files. mod and gofumpt wait for every
# generator above. terraform-fmt waits for docs.
# Group 5: check-compilation waits for mod and gofumpt.
# Group 6: lint-fix waits for docs, terraform-fmt, and check-compilation.
# Group 7: test-architecture waits for lint-fix.
include Makefile

.PHONY: pre-push-parallel-run pre-push-gofumpt
pre-push-parallel-run: test-architecture

PP_AFTER_SDK_SCHEMAS_AND_RESOURCES = generate-sdk generate-show-output-schemas generate-resource

generate-show-output-schemas: generate-sdk
generate-snowflake-object-assertions: generate-sdk
generate-snowflake-object-parameters-assertions: generate-sdk
generate-resource-parameters-assertions: generate-sdk
generate-resource-show-output-assertions: generate-sdk

generate-resource-assertions: $(PP_AFTER_SDK_SCHEMAS_AND_RESOURCES)
generate-resource-model-builders: $(PP_AFTER_SDK_SCHEMAS_AND_RESOURCES)
generate-datasource-model-builders: $(PP_AFTER_SDK_SCHEMAS_AND_RESOURCES)
generate-provider-model-builders: $(PP_AFTER_SDK_SCHEMAS_AND_RESOURCES)
generate-docs-additional-files: $(PP_AFTER_SDK_SCHEMAS_AND_RESOURCES)
generate-issue-labels: $(PP_AFTER_SDK_SCHEMAS_AND_RESOURCES)

# Leaves whose own recipes are not already behind another leaf in this list.
# generate-sdk, generate-resource, and show-output schemas are covered through
# the edges above.
PP_GENERATORS = generate-sdk-examples \
	generate-snowflake-object-assertions \
	generate-snowflake-object-parameters-assertions \
	generate-resource-show-output-assertions \
	generate-resource-parameters-assertions \
	generate-resource-assertions \
	generate-resource-model-builders \
	generate-datasource-model-builders \
	generate-provider-model-builders \
	generate-docs-additional-files \
	generate-issue-labels

mod: $(PP_GENERATORS)
pre-push-gofumpt: $(PP_GENERATORS)
	tools/bin/gofumpt -l -w .
terraform-fmt: docs

check-compilation: mod pre-push-gofumpt
lint-fix: docs terraform-fmt check-compilation
test-architecture: lint-fix
