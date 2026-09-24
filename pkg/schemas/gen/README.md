## show output schemas generation

These schemas are necessary to include SHOW output in every resource and datasource. The work is repetitive, so it's
easier to just generate all the needed schemas and mappers.

### Description

File [generate.go](../generate.go) invokes the generation logic from [main.go](./main/main.go). By default, all
definitions in `SdkShowResultStructs` ([sdk_show_result_structs.go](./sdk_show_result_structs.go)) are used — SHOW
structs plus struct DESCRIBE (`IsDescribe`) and property-row list entries (`UsedAsListEntry`). After successful
generation all SDK objects will have:

- show output schema that can be used in the resource/datasource (e.g. [user_gen](../user_gen.go))
- describe output schema when `IsDescribe` is set (`DescribeXSchema` in `{snake}_desc_gen.go`, `_details` suffix trimmed)
- mapper from the SDK object to the generated schema (e.g. [user_gen](../user_gen.go))

Unscoped generate and `generate-show-output-schemas-check` skip objects in `SHOW_OUTPUT_SCHEMAS_EXCLUDE` (Makefile). This generator is **Converging**. Customizations belong in `*_ext.go`: `SkipFields` omits a generated key; `AdditionalMapping` plus `additionalSchema` / `additionalToSchema` on the generated mapper type re-adds public keys. Callers always use generated `XToSchema`.

### How it works

##### Invoking the generation

To regenerate show outputs (skips `SHOW_OUTPUT_SCHEMAS_EXCLUDE`):

```shell
make generate-show-output-schemas
```

`make generate-show-output-schemas` / `generate-show-output-schemas-check` run in `pre-push` / `pre-push-check` with that exclude list. `make clean-show-output-schemas` deletes every `pkg/schemas/*_gen.go` file.

To generate only a chosen subset:
```shell
make generate-show-output-schemas SF_TF_GENERATOR_ARGS="--filter-object-names=sdk.User"
```

To generate an object that is in `SHOW_OUTPUT_SCHEMAS_EXCLUDE`:
```shell
make generate-show-output-schemas SHOW_OUTPUT_SCHEMAS_EXCLUDE= SF_TF_GENERATOR_ARGS='--filter-object-names=sdk.Warehouse'
```

```shell
# show usage
make generate-show-output-schemas SF_TF_GENERATOR_ARGS='-h'
```

##### Supported types

The following types are supported currently in the generator (schema and mappings):

- basic types (`string`, `int`, `float64`, `bool`)
- pointers to basic types (the same as above)
- `time.Time` (pointer too)
- enums based on `string` and `int` like `sdk.WarehouseType` or `sdk.ResourceMonitorLevel` (pointers too)
- identifiers (pointers too):
    - `sdk.AccountIdentifier`
    - `sdk.ExternalObjectIdentifier`
    - `sdk.AccountObjectIdentifier`
    - `sdk.DatabaseObjectIdentifier`
    - `sdk.SchemaObjectIdentifier`
    - `sdk.TableColumnIdentifier`
- `sdk.ObjectIdentifier` interface (nil-guarded; no dereference)

##### To schema mappings

Given SDK struct field can be mapped to the generated schema depending on its type:
- no mapping (`Identity`) - used for `string` and other basic types
- string value mapping (`ToString`) - used e.g. for `time.Time`
- fully qualified name mapping (`FullyQualifiedName`) - used for all identifiers and `sdk.ObjectIdentifier` interface (nil-guarded when the field is an interface)
- casting (`CastToString` and `CastToInt`) - used for enums with underlying type `string` or `int`

##### Changing the SDK object's show output

If you change the show output struct in the SDK:

1. Check if you don't introduce a type that is unsupported (check [supported types](#supported-types)
   and [known limitations](#known-limitations)).
2. Run generation according to [instructions](#invoking-the-generation).

##### Adding a new object to the SDK

1. Add a `ShowResultSchemaDef` to `SdkShowResultStructs` in [sdk_show_result_structs.go](./sdk_show_result_structs.go):
   - SHOW: `{ObjectStruct: sdk.<Singular>{}}` → `Show<Singular>Schema` in `<singular>_gen.go`
   - struct DESCRIBE: `{ObjectStruct: sdk.<Singular>Details{}, IsDescribe: true}` → `Describe<Singular>DetailsSchema` in `<singular>_desc_gen.go`
   - property-row list entry: `{ObjectStruct: sdk.<Type>{}, UsedAsListEntry: true}` → `<Type>Schema` in `<type>_gen.go`
   - `SkipFields: []string{"snake_case_key"}` omits that key from the schema map and `ToSchema`. Use when `*_ext.go` owns the field now, when the field is SDK-internal (keep omitted), or when a real SHOW/DESCRIBE property is missing from today’s public schema (skip now, add in a follow-up PR).
   - `AdditionalMapping: true` generates an empty mapper type that must implement `additionalSchemaMapper[T]` in `*_ext.go` (schema keys + mapper). Callers always use generated `XToSchema`. Not inferred from `SkipFields`.
2. Check if you don't introduce a type that is unsupported (check [supported types](#supported-types)
   and [known limitations](#known-limitations)).
3. Run generation according to [instructions](#invoking-the-generation).

### Next steps

##### Known limitations

- The following types (already existing in the SDK show output structs) are not yet supported (for all of them the
  schema will be generated with `schema.TypeInvalid`:
    - other basic types (e.g. `int8`, etc.)
    - slices of basic types (`[]int`, `[]string`)
    - slices of identifiers (`[]sdk.AccountIdentifier`, `[]sdk.SchemaObjectIdentifier`)
    - slices of enums (`[]sdk.IntegrationType`, `[]sdk.PluralObjectType`)
    - structs (`sdk.FileFormatTypeOptions`)

##### Improvements

Functional improvements:
- handle the missing types (TODOs in [schema_field_mapper.go](./schema_field_mapper.go))
  - handle nested structs with identifiers / slices of identifiers
- parametrize the generation, e.g.:
  - (optional) parametrize the output directory - currently, it's always written to `schemas` package

Implementation improvements:
- (optional) consider different implementations of `Mapper` (e.g. TODO in [schema_field_mapper_test.go](./schema_field_mapper_test.go): `ugly comparison of functions with the current implementation of mapper` and not ideal implementation in the [to_schema_mapper.tmpl](./templates/to_schema_mapper.tmpl): `runMapper .Mapper $nameLowerCase "." .OriginalName`)
