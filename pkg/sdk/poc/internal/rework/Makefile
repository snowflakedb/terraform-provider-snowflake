generate-sdk: ## Generate all SDK objects
	go generate ./generate.go
	-go generate ./*_dto_gen.go

clean-sdk: ## Clean all generated SDK objects
	rm -f ./*_gen.go
	rm -f ./*_gen_test.go
