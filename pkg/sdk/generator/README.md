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

Additional example definitions covering specific generator features:
- [sequences_def.go](example/defs/sequences_def.go) — full CRUD with `ShowOperationWithPairedStructs` and `DescribeOperationWithPairedStructs`
- [paired_struct_def.go](example/defs/paired_struct_def.go) — all `PairedStructs` field methods and options (see [PairedStructs](#pairedstructs) below)
- [to_opts_optional_example_def.go](example/defs/to_opts_optional_example_def.go) — `ListQueryStructField` with slice-of-structs toOpts, optional nested fields
- [drop_safely_example_def.go](example/defs/drop_safely_example_def.go) — `DropOperation` with `WithDropSafelyHook()` and `WithDropSafelyForce()` options
- [instance_method_example_def.go](example/defs/instance_method_example_def.go) — `InstanceMethodOperation` and `InstanceMethodOperationScalar`
- [enum_example_def.go](example/defs/enum_examples_def.go) — enum definitions with `Enum` and `OptionalEnum`

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

##### Known issues/limitations
- The generator was added after parts of the SDK were implemented manually. Some objects don't have the generator definitions which make it harder to keep the up-to-date. All of them should be gradually migrated to the definition-based generation implementation.
- The implementation of nested fields causes problems when reusing nested definitions (the same `[]Fields` slice is reused causing parent redefinition and incorrect mapping; the root cause being the lack of separation between the definition and model structs). It's currently validated programmatically and the panic is raised (`Field <field> already has a parent`). When it happens, create a function wrapper instead of directly creating a `var` with a definition.
- Currently, multiple `convert` method signatures are generated when multiple interface methods use the same object pairs. They are manually removed but such situations should be recognized automatically.

##### High-priority Improvements

This section aims mostly at reducing the manual labor when using the generator in its current state.

- Generate `ID()` methods for `Request` structs as already done for the `Opts` structs.
- Allow overriding the `ObjectType()` method returned `ObjectType`.

[//]: # (TODO [next PRs]: update next sections)

> ⚠️ **Disclaimer**: The following sections may contain the deprecated information. They will be cleaned up shortly.

##### Old High-priority improvements/changes

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
- Generate nested Request structs for fields that use slices of Opt objects (see the following fields Create operation in semantic_view.def: LogicalTables, semanticViewRelationships, etc.)

---

### DropSafely options

`DropOperation` accepts optional functional options that configure how the generated `DropSafely` method behaves:

- **`WithDropSafelyHook()`** — the generated `DropSafely` calls `v.dropSafelyHook(ctx, id)` before issuing the drop. Useful when a pre-drop side-effect is needed (e.g. revoking grants). The hook function must be implemented manually in a `_ext.go` file.

- **`WithDropSafelyForce()`** — the generated `DropSafely` appends `.WithForce(true)` to the Drop request, so dependent objects are also removed.

Example usage:
```go
// Calls v.dropSafelyHook(ctx, id) before dropping.
).DropOperation("https://...", dropStruct, g.WithDropSafelyHook())

// Appends .WithForce(true) to the Drop request.
).DropOperation("https://...", dropStruct, g.WithDropSafelyForce())
```

See [drop_safely_example_def.go](example/defs/drop_safely_example_def.go) for complete examples.

---

### PairedStructs

`PairedStructs` is a single-definition approach for declaring both the database row struct (`dbStruct`) and the plain SDK struct (`plainStruct`) in one field-by-field chain. It replaces the older pattern of calling `ShowOperation`/`DescribeOperation` with separate `DbStruct` and `PlainStruct` builders.

See [paired_struct_def.go](example/defs/paired_struct_def.go) for a complete usage example covering all supported field methods and options.

#### Constructor

**`StructPair(dbName, plainName string)`** — creates a new `PairedStructs` builder. `dbName` is the Go name for the database row struct (used in SQL scanning); `plainName` is the plain SDK struct name (returned by `Show`/`Describe`).

#### Field methods

Each method accepts zero or more `PairedFieldOption`s (see below).

| Method | db struct type | plain struct type |
|---|---|---|
| `Text(col)` | `string` | `string` |
| `OptionalText(col)` | `sql.NullString` | `*string` |
| `Bool(col)` | `bool` | `bool` |
| `OptionalBool(col)` | `sql.NullBool` | `*bool` |
| `BoolFromText(col)` | `string` | `bool` (compared to `"Y"`) |
| `OptionalBoolFromText(col)` | `sql.NullString` | `*bool` |
| `Number(col)` | `int` | `int` |
| `OptionalNumber(col)` | `sql.NullInt64` | `*int` |
| `Time(col)` | `time.Time` | `time.Time` |
| `OptionalTime(col)` | `sql.NullTime` | `*time.Time` |
| `PlainField(col, plainKind)` | `string` | `<plainKind>` |
| `OptionalPlainField(col, plainKind)` | `sql.NullString` | `<plainKind>` |
| `DataType(col)` | `string` | `datatypes.DataType` (via `ParseDataType`) |
| `StringList(col)` | `string` | `[]string` |
| `AccountObjectIdentifier(col)` | `string` | `AccountObjectIdentifier` (plain defaults to `"Id"`) |
| `OptionalAccountObjectIdentifier(col)` | `sql.NullString` | `*AccountObjectIdentifier` |
| `DatabaseObjectIdentifier(col)` | `string` | `DatabaseObjectIdentifier` (plain defaults to `"Id"`) |
| `SchemaObjectIdentifier(col)` | `string` | `SchemaObjectIdentifier` (plain defaults to `"Id"`) |
| `OptionalSchemaObjectIdentifier(col)` | `sql.NullString` | `*SchemaObjectIdentifier` |
| `NullableSchemaObjectIdentifierArray(col)` | `sql.NullString` | `[]SchemaObjectIdentifier` |
| `AccountIdentifierArray(col)` | `string` | `[]AccountIdentifier` |
| `SchemaObjectIdentifierWithArguments(col)` | `string` | `SchemaObjectIdentifierWithArguments` |
| `OptionalSchemaObjectIdentifierWithArguments(col)` | `sql.NullString` | `*SchemaObjectIdentifierWithArguments` |
| `Enum(col, enumDef)` | `string` | `<enumType>` |
| `OptionalEnum(col, enumDef)` | `sql.NullString` | `*<enumType>` |
| `JsonField(col, kind)` | `string` | `<kind>` (via `json.Unmarshal`) |
| `Field(col, dbKind, plainKind)` | explicit | explicit |

#### PairedFieldOption options

- **`WithDbFieldName(name)`** — override the Go field name in the db row struct (default: derived from `col` via `ToSnakeCase` → `ToCamelCase`).
- **`WithPlainFieldName(name)`** — override the plain struct field name.
- **`WithRequiredInPlain()`** — strip the pointer from the plain kind (e.g. `sql.NullString` db → `string` plain instead of `*string`).
- **`WithCustomParser(funcName)`** — use a custom parse function `func(string) (T, error)` to convert the db value.
- **`WithValueAdjuster(funcName)`** — apply an adjustment function `func(T) T` to the converted value after assignment.
- **`WithBoolTrueValue(v)`** — override the truthy string for `BoolFromText`/`OptionalBoolFromText` (default `"Y"`).
- **`WithBoolParsed()`** — use `strconv.ParseBool` instead of a fixed string comparison for bool fields.
- **`WithManualConvert()`** — skip this field in the generated `convert()` body; handle it manually in `additionalConvert()` in a `_ext.go` file.

#### PairedStructs modifiers

- **`WithoutConvertGeneration()`** — disable `convert()` body generation for this pair entirely.
- **`WithShowResultFilterHook()`** — enable row filtering; the generated code calls `excludeFromShow()` which must be implemented in a `_ext.go` file.
