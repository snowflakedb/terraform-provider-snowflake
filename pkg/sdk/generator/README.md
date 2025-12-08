> ⚠️ **Disclaimer**: The SDK generator started as PoC but was widely used to speed up the development of the SQL abstraction over Snowflake. It requires a lot of changes as improvements as working with it is not always the easiest. Additionally, we are currently considering the move to REST API (check [this roadmap entry](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#snowflake-rest-apis)), which may ultimately lead to deprecation of this generator as SQL abstraction may not be needed anymore.

## SDK generator

Generating full SDK object implementation based on object definition.

### How it works
##### Adding new object to the SDK

To add definition for the new SDK object:

1. Create file `<object_name_plural>_def.go` in the [defs directory](defs) (e.g., [sequences_def.go](defs/sequences_def.go)).
2. Create object definition in the created file. Base it on the existing definitions and the [example directory](example).
3. Add the created definition to the [0_init.go](defs/0_init.go) file.
4. You are all set to run generation.

##### Invoking generation

> **Important**: All the commands should be run from the main project directory.

The generator offers filtering by object name and by generation part. To list all available objects and generation parts run:
```shell
make generate-sdk SF_TF_GENERATOR_ARGS='--help'
```

Available generation parts:
- default
- dto
- dto_builders
- impl
- unit_tests
- validations

Generator is built on top of our common generator (read more in its [README](../../internal/genhelpers/README.md)). Experiment with the following commands:

```shell
# generate all objects and all files
make generate-sdk
```
```shell
# remove all generated files first; generate all objects and all files
make clean-generated-sdk generate-sdk
```
```shell
# generate all objects and chosen files only
make generate-sdk SF_TF_GENERATOR_ARGS='--filter-generation-part-names=default,dto,validations'
```
```shell
# generate chosen objects only and all files
make generate-sdk SF_TF_GENERATOR_ARGS='--filter-object-names=Sequences'
```
```shell
# generate chosen objects and chosen files only
make generate-sdk SF_TF_GENERATOR_ARGS='--filter-generation-part-names=default,impl --filter-object-names=Sequences'
```

##### Examples

There are example files ready for generation, e.g. [database_role_def.go](example/defs/database_role_def.go), which creates files:
- [database_role_gen.go](example/database_roles_gen.go) - SDK interface, options structs
- [database_role_dto_gen.go](example/database_roles_dto_gen.go) - SDK Request DTOs
- [database_role_dto_builders_gen.go](example/database_roles_dto_builders_gen.go) - SDK Request DTOs constructors and builder methods
- [database_role_validations_gen.go](example/database_roles_validations_gen.go) - options structs validations
- [database_role_impl_gen.go](example/database_roles_impl_gen.go) - SDK interface implementation
- [database_role_gen_test.go](example/database_roles_gen_test.go) - unit tests placeholders with guidance comments (at least for now)

The commands follow the same format as the official SDK ones:

```shell
# generate all example objects and all files
make generate-sdk-examples
```
```shell
# remove all example generated files first; generate all example objects and all files
make clean-generated-sdk-examples generate-sdk-examples
```
```shell
# generate all example objects and chosen files only
make generate-sdk-examples SF_TF_GENERATOR_ARGS='--filter-generation-part-names=default,dto,validations'
```
```shell
# generate chosen example objects only and all files
make generate-sdk-examples SF_TF_GENERATOR_ARGS='--filter-object-names=Sequences'
```
```shell
# generate chosen example objects and chosen files only
make generate-sdk-examples SF_TF_GENERATOR_ARGS='--filter-generation-part-names=default,impl --filter-object-names=Sequences'
```
```shell
# show usage
make generate-sdk-examples SF_TF_GENERATOR_ARGS='--help'
```

##### Known issues
- The implementation of nested fields causes problems when reusing nested definitions (the same `[]Fields` slice is reused causing parent redefinition and incorrect mapping; the root cause being the lack of separation between the definition and model structs). It's currently validated programmatically and the panic is raised (`Field <field> already has a parent`). When it happens, create a function wrapper instead of directly creating a `var` with a definition.

[//]: # (TODO [next PRs]: update this section)
### Next steps

> ⚠️ **Disclaimer**: This section may contain the deprecated essentials/improvements/known issues. It will be cleaned up shortly.

##### High-priority improvements/changes

This section was introduced after more than 1.5 year of using this generator for development.
It aims to be the up-to-date section covering the next changes that should be applied in this generator.
Most of the topics listed here aim to reduce the manual work needed after the generation.
- Generate ID() methods for `Request` structs as already done for the `Opts` structs.
- Fully generate `convert` function body (from objectDBRow to object) as it needs to be added manually.
  - Start with 1-to-1 conversion
  - Support custom mappers (for lists or identifiers)
- Improve lists handling
  - wrong generated validations for `validIdentifierIfSet` for cases like
    ```go
    A := QueryStruct("A").
        Identifier("Name").
        Validation(ValidIdentifierIfSet, "Name")
    B := QueryStruct("B")
        .ListQueryStructField(A) // A []A - validations will be wrong because this is array
    ```
  - generation for lists is not recursive, so we're only supporting one-level mapping
    ```go
    []Request1{ foo Request2, bar int } // won't be converted
    []Request1{ foo string, bar int } // will
    ```
  - optional structs inside slices are not being checked for nil in the toOpts() methods in the implementation file
  - pointer fields are not being handled correctly in lists
- add support for enums generation as for now, we add them and conversion methods manually (e.g. [`UserType`](https://github.com/snowflakedb/terraform-provider-snowflake/blob/5bdcd127d9288212b10ea7b138bebc0cb770c5b9/pkg/sdk/users.go#L746))
  - Add new generator struct `EnumMapping` to generate mappings
  - Add new template injected to interface.tmpl creating mappings functions
  - Generate unit tests for the enum type converters
  - Handle synonyms (e.g. [`ToWarehouseSize`](https://github.com/snowflakedb/terraform-provider-snowflake/blob/5bdcd127d9288212b10ea7b138bebc0cb770c5b9/pkg/sdk/warehouses.go#L77))
- Move manually added changes to the separate `_ext.go` files (e.g. [`procedures_ext.go`](https://github.com/snowflakedb/terraform-provider-snowflake/blob/5bdcd127d9288212b10ea7b138bebc0cb770c5b9/pkg/sdk/procedures_ext.go#L19))
- Add definitions for the objects that were written by hand before the generator was created
- In dto builder generator, use first letter lowercase for method arguments (e.g. [`OrReplace`](https://github.com/snowflakedb/terraform-provider-snowflake/blob/5bdcd127d9288212b10ea7b138bebc0cb770c5b9/pkg/sdk/procedures_dto_builders_gen.go#L26) -> `orReplace`)
- Reuse existing `genhelpers.NewGenerator()` or other useful tools from the `genhelpers` package (to e.g. allow running it with command line arguments to regenerate only the chosen files)
- Refresh the example

##### Essentials
- generate each branch of alter in tests (instead of basic and all options)
- clean up predefined operations in generator (now casting to string)
- handle more validation types
  - validating numbers in a given range constrained by another variable (e.g. `x <= y`, `x > y`, etc.)
  - validating number relations in a sequence (e.g. `x <= y <= z`, `x < y < z`)
- write new `valueSet` function (see validations.go) that will have better defaults or more parameters that will determine
checking behaviour which should get rid of edge cases that may cause bugs in the future
   - right now, we have `valueSet` function that doesn't take into consideration edge cases, e.g. with slice where sometimes
   we would like to do something like `alter x set y = ()` (set empty array to unset `y`). Those edge cases have cause on our
   validation, and it determines sometimes if we'll return an error or not, which can lead to bugs!
- refactor generation of `Describe`, so it will tak context and request as arguments
  - all the interface functions should have context and request as arguments for the sake of API consistency and generation simplicity
- check if SelfIdentifier implementation is correct (mostly type, because it's derived from interface obj) by checking
if there's a resource with different types of identifiers across queries (e.g. Create <AccountObjectIdentifier>, Alter <SchemaObjectIdentifier>)
- we should specify prefix / postfix standard for top-level items in _def.go files to avoid any conflicts in the package
- remove name argument from QueryStruct in the Operation, because Opt structs in the Operation will have name from op name + interface field and not query struct itself
- Derive field name from QueryStruct, e.g. see network_policies_def where we can remove "Set" field, but we have to make a convention of creating nested struct with
name pattern like <interface name><name> e.g. NetworkPoliciesSet or NetworkPolicySet, then we could automatically remove prefix and we'll name field with postfix, so "Set" in this case
- Add more operations (every operation ?) in the database_role_def.go example
- Divide into packages or add common prefix for similar files (e.g. struct_plain.go, struct_db.go or builders_keyword.go, builders_parameter.go)
- Make a clear division between DSL files and model files (etc. QueryStruct(DSL) and Field(Model)) and divide them into separate packages (?)
- Add parameter to DtoTemplate (templates.go) to generate the right path to the dto generator's main.go file
- Right now to avoid generated structs duplication, arrays containing struct names have been introduced (template_executors.go),
find a better solution to solve the issue (add more logic to the templates ?)

##### Improvements
- automatic names of nested `struct`s (e.g. `DatabaseRoleRename`)
- check if generating with package name + invoking format removes unnecessary qualifier
- consider merging templates `StructTemplate` and `OptionsTemplate` (requires moving Doc to Field)
- expand unit tests generation
- experiment with Snowflake table (any table) representation in Go in order to implement DbStruct -> PlainStruct convert function
  - see if *string can have similar effect as sql.NullString (check go-snowflake connector ?)
     - if yes, then we should be using pointers instead of abstractions like sql.NullString and we can
     modify ShowMapping and DescribeMapping to generate convert function with automatic conversion (as we have in DTOs).
     warehouses.go is a good place to start with when planning mapping strategy, because there's a lot of different mapping cases.
- when calling .SelfIdentifier we can implicitly also add validateObjectIdentifier validation rule
- enforce user to use KindOf... functions with interface
  - example implementation - StringTyper implements Typer and all the KindOf... functions use StringTyper to return Typer easily - https://go.dev/play/p/TZZgSkkHw_M
- `queryStruct` should be spilled into `Operation` interface file, because the idea was to have model which is unaware of DSL used to create it.
- generate full tests for common types (e.g. setting/unsetting tags)
- generate common resources for integration tests
- cleanup the design of builders in DSL (e.g. why transformer has to be always added?)
- generate getters for requests, at least for identifier/name
- generate integration tests in child package (because now we keep them in `testint` package)
- struct_to_builder is not supporting templated-like values. See stages_def.go where in SQL there could be value, where 'n' can be replaced with any number
  - `SKIP_FILE_n` - this looks more like keyword without a space between SQL prefix and int
  - `SKIP_FILE_n%` (e.g. `SKIP_FILE_123%`) - this is more template-like behaviour, notice that 'n' is inside the value (we cannot reproduce that right now with struct_to_builder capabilities)
- fix builder generation
  - we can add `flatten` option in cases where some sql structs had to be nested to create correct sql representation
    - for example encryption options in `stages_def.go` (instead of calling `.WithEncryption(NewEncryptionRequest(encryption))` we could call `.WithEncryption(encryption)`)
  - operation names (or their sql struct names) should dictate more how constructors are made
- better handling of list of strings/identifiers
  - there should be no need to define custom types every time
  - more clear definition of lists that can be empty vs cannot be empty
- add empty ids in generated tests (TODO in random_test.go)
- add optional imports (currently they have to be added manually, e.g. `datatypes.DataType`)
- handle objects that do not have ids
  - ShowById should take more customizable attributes, instead of only object ID
  - Add a possibility to generate a non-sql method with a custom implementation. Currently, it is done only in `ShowById...` functions with `newNoSqlOperation`.
- improve handling operations that return one row
- add more context to validated identifiers, so that error contains the affected field
- add custom identifier wrapping, like it's used in security integrations' network policies
- Generate nested Request structs for fields that use slices of Opt objects (see the following fields Create operation in semantic_view.def: LogicalTables, semanticViewRelationships, etc.)

##### Known issues
- generating two converts when Show and Desc use the same data structure
