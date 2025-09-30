## generator commons

Because we generate a bunch of code in the project, and we tend to copy-paste similar setup for the generators, templates, etc., we decided to introduce common generator creation utils that can be reused by variety of generators.

### Description

The main building blocks of this package are:
- `generator.go` defining `Generator[T ObjectNameProvider, M GenerationModel] struct` allowing to create new generators
- `generation_part_filters.go` and `object_filers.go` containing the common generation filters
- `flags.go` defining custom flag types that are used to alter generators behavior
- `mappers.go` defining mappers that can be reused in the generated objects
- `models.go` defining common models that can be reused throughout generators
- `struct_details_extractor.go` allowing to parse any struct to retrieve its information (for the later generation purposes)
- `template_commons.go` containing template helper functions and the easy way to use them without providing their name everytime
- `util.go` with a variety of util functions

### How it works

Each generator consists of:
- Name and version - version should be bumped whenever the templates or generation logic changes.
- List of input objects - they are transformed to models used in templates.
- List of generation parts - they contain the needed templates and file naming logic.
- (optional) Additional object/generation part filters - generator allows defining custom filters altering the given generators logic further.
- (currently optional) Description - used for printing the usage for the given generator (`-h` or `--help` flag, check the TODO section below).
- (currently optional) Makefile command part - we usually invoke the generator as make commands; it's used for printing the usage for the given generator (`-h` or `--help` flag, check the TODO section below).

Each generator workflow looks as follows:
- The list of input objects is filtered using all defined filters (default and additional).
- The list of generation parts is filtered using all defined filters (default and additional).
- For each remaining object, each remaining generation part is run.
- The generator returns to OS.

#### Defining and running a new generator

Before proceeding with the following steps check [objectassert/gen](../../acceptance/bettertestspoc/assert/objectassert/gen) package for reference.

To create a new generator:
1. Create `gen` package in the destination package with:
    - `main/main.go` file
    - `templates` directory
    - `model.go` containing the model definition and conversion
    - `templates.go` containing the templates definitions and helper functions 
2. Create `generate.go` file on the same level as the `gen` package above with the following content only (in addition to the package name) `//go:generate go run ./gen/main/main.go $SF_TF_GENERATOR_ARGS`.
3. In the `gen/main/main.go` create and run a new generator. This means invoking the `genhelpers.NewGenerator` and:
   - providing the name and version for the generator
   - (currently optional) providing a short description and Makefile command part
   - providing an input definition for the source objects
   - defining all needed templates
   - providing all generation parts using the defined templates and containing the file naming logic
   - providing method to enrich source object definitions with the necessary content
   - providing method to translate enriched objects to the models used inside the templates
   - (optional) providing additional debug output you want to run for each of the objects
   - (optional) providing additional filters to limit the generation to only specific objects or generation parts
   - ending with `RunAndHandleOsReturn()`
4. Add two entries to our Makefile:
   - first for a cleanup, in form of `clean-<makefile-command-part>`, e.g.
   ```makefile
   clean-snowflake-object-parameters-assertions:
       rm -f ./pkg/acceptance/bettertestspoc/assert/objectparametersassert/*_gen.go
   ```
   - second for a generation itself, in form of `generate-<makefile-command-part>`, e.g.
   ```makefile
   generate-snowflake-object-parameters-assertions:
       go generate ./pkg/acceptance/bettertestspoc/assert/objectparametersassert/generate.go
   ```
5. By default, generator support the following command line flags (invokable with e.g. `make generate-<makefile-command-part> SF_TF_GENERATOR_ARGS='<space-separated flags>'`)
   - `--dry-run` allowing to print the generated content to the command line instead of saving it to files
   - `-h` or `--help` printing the usage for the given generator
   - `--filter-generation-part-names <name1>,<name2>,...` allowing to generate only for the given generation part names; defaults to empty list meaning all generation parts are used
   - `--filter-object-names <name1>,<name2>,...` allowing to generate only for the given object names; defaults to empty list meaning all objects are used
   - `--verbose` allowing to see the all the additional debug logs

### Next steps

#### Improvements

Functional improvements:
- add a generic terraform schema reader, to allow later generation from schemas
- handle the missing types (TODOs in [struct_details_extractor_test.go](./struct_details_extractor_test.go))
- add support for custom command line flags

Implementation improvements:
- add acceptance test for a `testStruct` (the one from [struct_details_extractor_test.go](./struct_details_extractor_test.go)) for the whole generation flow
- add description to all publicly available structs and functions (multiple TODOs left)
- describe and test all the template helpers (TODOs left in `templates_commons.go`)
- test writing to file (TODO left in `util.go`)
