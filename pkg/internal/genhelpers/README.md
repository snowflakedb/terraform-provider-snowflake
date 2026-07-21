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
   - `--exclude-object-names <name1>,<name2>,...` allowing to exclude the given object names from generation; composes with inclusion filters (AND logic)
   - `--exclude-generation-part-names <name1>,<name2>,...` allowing to exclude the given generation part names from generation; composes with inclusion filters (AND logic)
   - `--enable-generation-part-names <name1>,<name2>,...` allowing to force-enable optional (disabled-by-default) generation parts for all objects during a run
   - `--verbose` allowing to see the all the additional debug logs

### Next steps

#### Improvements

Functional improvements:
- Currently, generation part filters are applied to all filtered objects. There are situations (like adding a generation part) in which we would like to use this part only for a subset of objects without disrupting the generation of other objects within the given generator. For that reason, generation part filtering could be handled on the object-by-object basis. It means:
  - Extracting a common generation part setting on the object level. Consider:
    - Setting it as additional field to the `PreambleModel` (in which case, maybe it's good to finally rename this struct to a more appropriate name).
    - Setting it as a separate struct as `PreambleModel` is currently reused throughout the whole generator.
    - Extracting a dedicated common struct that would be nested in each object definition; it could contain the default generation parts for the given object and some other common object-level generation settings in the future. Also, the object name currently acquired from the `ObjectNameProvider` could be a part of this dedicated struct and handled like the `HasPreambleModel` interface.
  - Providing a generator-wide default for cases when generation part filtering is not set on the given object (e.g. when the new generation part is being developed it can be set on the chosen object without affecting the others and without the need to alter their definitions).
    - Similarly to the defaults for the generation parts, we could provide the default object filtering for the given generator (e.g. for situations where given object requires some additional attention but all others can be automatically regenerated without problems).
  - Rethinking how the filtering should be handled on the generator level and how it should behave with combination of the command-line argument (i.e. should it override filters for all objects? should it be treated as an additional filter in addition to the setting on the object level? if the latter, should we have another option to generate the given part even if it's not listed for the given object?).
  - Ensuring that the listed parts are available for the given generator (compile-time or runtime validation?).
- add a generic terraform schema reader, to allow later generation from schemas
- handle the missing types (TODOs in [struct_details_extractor_test.go](./struct_details_extractor_test.go))
- add support for custom command line flags
- Currently, there is no friendly error handling in the created generators. Some of them, use panic to prevent obvious configuration mistakes. It would be better to handle errors programmatically (e.g. do not fail on the first one but collect errors from all object declarations and present nicely to the user invoking).

Implementation improvements:
- add acceptance test for a `testStruct` (the one from [struct_details_extractor_test.go](./struct_details_extractor_test.go)) for the whole generation flow
- add description to all publicly available structs and functions (multiple TODOs left)
- describe and test all the template helpers (TODOs left in `templates_commons.go`)
- test writing to file (TODO left in `util.go`)

## Running Generators

This section is for developers who need to re-run generators after changing inputs (SDK definitions, resource schemas, show result structs, etc.). For building a new generator from scratch, see [Defining and running a new generator](#defining-and-running-a-new-generator) above.

### File Conventions

| Pattern | Editable? | Purpose |
|---|---|---|
| `*_gen.go` | **No** | Generated; contains `// Code generated by … DO NOT EDIT.` |
| `*_ext.go` / `*.ext.go` | **Yes** | Manual extensions alongside generated code |
| `*_gen_test.go` | **Yes** | SDK unit test skeletons — currently hand-editable |

Never edit `*_gen.go` files. Changes belong in definition files or `*_ext.go` extensions.

### Generator Reference

| Make target | Output location | State |
|---|---|---|
| `generate-sdk-no-tests` | `pkg/sdk/*_gen.go` | **Enforced** |
| `generate-sdk` | `pkg/sdk/*_gen.go` + `*_gen_test.go` | Converging |
| `generate-show-output-schemas` | `pkg/schemas/*_gen.go` | Converging |
| `generate-snowflake-object-assertions` | `pkg/acceptance/bettertestspoc/assert/objectassert/*_gen.go` | **Enforced** |
| `generate-snowflake-object-parameters-assertions` | `pkg/acceptance/bettertestspoc/assert/objectparametersassert/*_gen.go` | **Enforced** |
| `generate-resource-assertions` | `pkg/acceptance/bettertestspoc/assert/resourceassert/*_gen.go` | **Enforced** |
| `generate-resource-parameters-assertions` | `pkg/acceptance/bettertestspoc/assert/resourceparametersassert/*_gen.go` | **Enforced** |
| `generate-resource-show-output-assertions` | `pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert/*_gen.go` | **Enforced** |
| `generate-resource-model-builders` | `pkg/acceptance/bettertestspoc/config/model/*_gen.go` | **Enforced** |
| `generate-datasource-model-builders` | `pkg/acceptance/bettertestspoc/config/datasourcemodel/*_gen.go` | **Enforced** |
| `generate-provider-model-builders` | `pkg/acceptance/bettertestspoc/config/providermodel/*_gen.go` | **Enforced** |

**Enforced** — regeneration is idempotent; verified in CI by `make pre-push-check`.
**Converging** — still being cleaned up; `*_gen.go` files may contain manual edits; treat with care.

Aggregate targets: `generate-all-config-model-builders` (all model builders), `generate-all-assertions-and-config-models` (all assertions + model builders), `make pre-push` (everything + formatting + linting + docs).

### When to Regenerate

Always scope to the object you changed — avoid regenerating everything unless preparing a PR.

| After this change… | Scoped command |
|---|---|
| Modify an SDK definition in `pkg/sdk/generator/defs/` | `make generate-sdk-no-tests SF_TF_GENERATOR_ARGS='--filter-object-names=<Object>'` |
| Modify a show result struct in `pkg/sdk/` | `make generate-show-output-schemas SF_TF_GENERATOR_ARGS='--filter-object-names=<Object>'` |
| Add or modify a resource schema | `make generate-resource-assertions SF_TF_GENERATOR_ARGS='--filter-object-names=<Object>'` and `make generate-resource-model-builders SF_TF_GENERATOR_ARGS='--filter-object-names=<Object>'` |
| Add or modify a data source schema | `make generate-datasource-model-builders SF_TF_GENERATOR_ARGS='--filter-object-names=<Object>'` |
| Before pushing any PR | `make pre-push` |

> **Multi-PR note:** Each PR only needs its own generators. Do not run resource assertion or model builder generators in an SDK-only PR — those belong in the resource PR.

### CLI Flags

All genhelpers-based generators accept flags via `SF_TF_GENERATOR_ARGS`:

```bash
make generate-<target> SF_TF_GENERATOR_ARGS='<flags>'
```

| Flag | Purpose |
|---|---|
| `--dry-run` | Print output to stdout, write nothing |
| `-h` / `--help` | List valid object names, generation parts, and flags |
| `--filter-object-names=A,B` | Only generate for listed objects (case-sensitive) |
| `--exclude-object-names=A,B` | Skip listed objects |
| `--filter-generation-part-names=A,B` | Only run listed generation parts |
| `--exclude-generation-part-names=A,B` | Exclude listed generation parts |
| `--verbose` | Print additional debug output |

Run `make generate-<target> SF_TF_GENERATOR_ARGS='-h'` to discover valid object names for any generator.

> **Note:** `generate-issue-labels` is a standalone tool and does not support `SF_TF_GENERATOR_ARGS`.

### Troubleshooting

| Symptom | Cause | Fix |
|---|---|---|
| Panic with stack trace | Invalid definition or generator config | Use `--dry-run --filter-object-names=<One>` to isolate the failing object |
| No output files created | Object name mismatch | Run with `-h`; names are case-sensitive |
| Partial/corrupt files after interruption | Generator killed mid-run | `git checkout -- <output-dir>/` then regenerate the affected object |
| `SF_TF_GENERATOR_ARGS` ignored | Target is not genhelpers-based | Only applies to genhelpers generators |
