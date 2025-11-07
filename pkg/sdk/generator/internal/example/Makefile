SF_TF_GENERATOR_ARGS ?= --filter-object-names=Sequences,DatabaseRoles

generate-sdk: ## Generate all SDK objects
	go generate ./generate.go

clean-sdk: ## Clean all generated SDK objects
	rm -f ./*_gen.go
	rm -f ./*_gen_test.go

generate-examples: ## Generate only the predefined subset of examples
	$(MAKE) clean-sdk generate-sdk SF_TF_GENERATOR_ARGS='$(SF_TF_GENERATOR_ARGS)'
