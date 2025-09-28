## generator commons

Because we generate a bunch of code in the project, and we tend to copy-paste similar setup for the generators, templates, etc., we decided to introduce common generator creation utils that can be reused by variety of generators.

### Description

The main building blocks of this package are:
- `generator.go` defining `Generator[T ObjectNameProvider, M GenerationModel] struct` allowing to create new generators
- `mappers.go` defining mappers that can be reused in the generated objects
- `struct_details_extractor.go` allowing to parse any struct to retrieve its information (for the later generation purposes)
- `template_commons.go` containing template helper functions and the easy way to use them without providing their name everytime
- `util.go` with a variety of util functions

### How it works

##### Defining and running a new generator

Before proceeding with the following steps check [objectassert/gen](../../acceptance/bettertestspoc/assert/objectassert/gen) package for reference.

To create a new generator:
1. Create `gen` package in the destination package with:
    - `main/main.go` file
    - `templates` directory
    - `model.go` containing the model definition and conversion
    - `templates.go` containing the templates definitions and helper functions 
2. Create `generate.go` file on the same level as the `gen` package above with the following content only (in addition to the package name) `//go:generate go run ./gen/main/main.go $SF_TF_GENERATOR_ARGS`.
3. In the `gen/main/main.go` create and run a new generator. This means:
   - providing an input definition for the source objects
   - method to enrich source object definitions with the necessary content
   - method to translate enriched objects to the models used inside the templates
   - method with the generated files naming strategy
   - list of all the needed templates
   - (optionally) additional debug output you want to run for each of the objects
   - (optionally) a filter to limit the generation to only specific objects
4. Add two entries to our Makefile:
   - first for a cleanup, e.g. `rm -f ./pkg/acceptance/bettertestspoc/assert/objectparametersassert/*_gen.go`
   - second for a generation itself, e.g. `go generate ./pkg/acceptance/bettertestspoc/assert/objectparametersassert/generate.go`
5. By default, generator support the following command line flags (invokable with e.g. `make generate-show-output-schemas SF_TF_GENERATOR_ARGS='--dry-run --verbose'`)
   - `--dry-run` allowing to print the generated content to the command line instead of saving it to files
   - `--verbose` allowing to see the all the additional debug logs

[//]: # (TODO [this PR]: Update this description)

### Next steps

##### Improvements

Functional improvements:
- add a generic terraform schema reader, to allow later generation from schemas
- handle the missing types (TODOs in [struct_details_extractor_test.go](./struct_details_extractor_test.go))

Implementation improvements:
- add acceptance test for a `testStruct` (the one from [struct_details_extractor_test.go](./struct_details_extractor_test.go)) for the whole generation flow
- add description to all publicly available structs and functions (multiple TODOs left)
- describe and test all the template helpers (TODOs left in `templates_commons.go`)
- test writing to file (TODO left in `util.go`)
