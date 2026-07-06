> ⚠️ **Disclaimer**: The SDK generator started as PoC but was widely used to speed up the development of the SQL abstraction over Snowflake. SDK in its current state fully depends on this generation and no manual changes are needed (besides unit tests). When adding the new SDK object, make sure the regeneration goes smoothly. Additionally, we are currently considering the move to REST API (check [this roadmap entry](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#snowflake-rest-apis)), which may ultimately lead to deprecation of this generator as SQL abstraction may not be needed anymore.

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

Generation parts can be registered as **optional** (disabled by default). Optional parts are only generated for objects that explicitly enable them via `WithEnabledGenerationParts(...)`. Use `WithOptionalGenerationPart(...)` on the generator to register such a part.

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
```shell
# generate all objects except the given ones
make generate-sdk SF_TF_GENERATOR_ARGS='--exclude-object-names=Sequences'
```
```shell
# generate all files except the given generation parts
make generate-sdk SF_TF_GENERATOR_ARGS='--exclude-generation-part-names=unit_tests'
```
```shell
# combine inclusion and exclusion filters
make generate-sdk SF_TF_GENERATOR_ARGS='--filter-object-names=Sequences,DatabaseRoles --exclude-object-names=Sequences'
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
# generate all example objects except the given ones
make generate-sdk-examples SF_TF_GENERATOR_ARGS='--exclude-object-names=Sequences'
```
```shell
# generate all example files except the given generation parts
make generate-sdk-examples SF_TF_GENERATOR_ARGS='--exclude-generation-part-names=unit_tests'
```
```shell
# show usage
make generate-sdk-examples SF_TF_GENERATOR_ARGS='--help'
```

##### Known issues/limitations
- The generator was added after parts of the SDK were implemented manually. Some objects don't have the generator definitions which make it harder to keep the up-to-date. All of them should be gradually migrated to the definition-based generation implementation.
- The implementation of nested fields causes problems when reusing nested definitions (the same `[]Fields` slice is reused causing parent redefinition and incorrect mapping; the root cause being the lack of separation between the definition and model structs). It's currently validated programmatically and the panic is raised (`Field <field> already has a parent`). When it happens, create a function wrapper instead of directly creating a `var` with a definition.

##### Remaining TODOs

- Generate `ID()` methods for `Request` structs as already done for the `Opts` structs.
- PoC of unit tests generation without manual changes in the generated files (details soon)
  - generate each branch of alter in tests (instead of basic and all options)
- Improve validation handling for nested slices (the path is built incorrectly now)
- `PlainStruct`-only fields do not currently trigger the additionalConvert creation, as they are filtered out in the iteration

##### CI guard targets

- `make generate-sdk-no-tests-check` — regenerates all non-test parts (`default`, `dto`, `dto_builders`, `impl`, `validations`) for all objects and verifies no diff. Wired into `pre-push-check`.
- `make generate-sdk-examples-check` — regenerates examples and verifies no diff. Wired into `pre-push-check`.

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
| `PlainOnlyField(fieldName, plainKind)` | _(none)_ | `<plainKind>` (must be populated in `additionalConvert()`) |

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

### ShowByID suppression

By default, `ShowOperationWithPairedStructs` auto-generates `ShowByID` and `ShowByIDSafely` methods.
When an object requires a custom `ShowByID` signature (e.g. additional parameters), pass
`g.ShowByIDSuppressed` as the filtering argument to suppress auto-generation:

```go
.ShowOperationWithPairedStructs("https://...", pairs, queryStruct, g.ShowByIDSuppressed).
    WithCustomInterfaceMethod("ShowByID", "", []*g.MethodParameter{...}, "*Object", "error").
    WithCustomInterfaceMethod("ShowByIDSafely", "", []*g.MethodParameter{...}, "*Object", "error")
```

The custom methods must be implemented manually in a `_ext.go` file.

See [user_programmatic_access_tokens_def.go](defs/user_programmatic_access_tokens_def.go) for an example.

---

### Shared struct reuse across objects

When a `QueryStruct` is shared across multiple object definitions (e.g. `SecretsList` used by
Functions, Notebooks, and Procedures), use `WithSharedToOpts()` and `OptionalSharedQueryStructField`
to avoid duplicate struct declarations.

**In the originating object** (the one that "owns" the struct):
```go
var sharedStruct = g.NewQueryStruct("SecretsList").
    List("SecretsList", "SecretReference", g.ListOptions().Required().MustParentheses()).
    WithSharedToOpts()  // generates a standalone func (r *SecretsListRequest) toOpts() *SecretsList
```

**In reusing objects:**
```go
OptionalSharedQueryStructField("Secrets", sharedStruct, g.ParameterOptions().SQL("SECRETS").Parentheses())
```

**What happens:**
- The struct type (`SecretsList`), Request type (`SecretsListRequest`), and constructor (`NewSecretsListRequest`) are generated only once — in the originating object.
- The standalone `toOpts()` method is generated in the originating object's `_impl_gen.go`.
- In all objects (including the originator), the `toOpts` mapping calls `.toOpts()` instead of inlining the field mapping.
- Reusing objects skip struct/DTO/constructor generation for the shared field.

See [functions_def.go](defs/functions_def.go) (originator) and [notebooks_def.go](defs/notebooks_def.go) (reuser).

### Potential Improvements

- handle more validation types
  - validating numbers in a given range constrained by another variable (e.g. `x <= y`, `x > y`, etc.)
  - validating number relations in a sequence (e.g. `x <= y <= z`, `x < y < z`),
  - validating inputs containing blocklisted characters, e.g. `$$`.
- remove name argument from QueryStruct in the Operation, because Opt structs in the Operation will have name from op name + interface field and not query struct itself
- Derive field name from QueryStruct, e.g. see network_policies_def where we can remove "Set" field, but we have to make a convention of creating nested struct with
  name pattern like <interface name><name> e.g. NetworkPoliciesSet or NetworkPolicySet, then we could automatically remove prefix and we'll name field with postfix, so "Set" in this case
- automatic names of nested `struct`s (e.g. `DatabaseRoleRename`)
- enforce user to use KindOf... functions with interface
  - example implementation - StringTyper implements Typer and all the KindOf... functions use StringTyper to return Typer easily - https://go.dev/play/p/TZZgSkkHw_M
- cleanup the design of builders in DSL (e.g. why transformer has to be always added?)
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
- add more context to validated identifiers, so that error contains the affected field
