## show output schemas generation

These schemas are necessary to include SHOW output in every resource and datasource. The work is repetitive, so it's
easier to just generate all the needed schemas and mappers.

### Description

File [generate.go](../generate.go) invokes the generation logic from [main.go](./main/main.go). By default, all SDK show
output struct are used (listed in [sdk_show_result_structs.go](./sdk_show_result_structs.go)). After successful
generation all SDK objects will have:

- show output schema that can be used in the resource/datasource (e.g. [warehouse_gen](../warehouse_gen.go#L11))
- mapper from the SDK object to the generated schema (e.g. [warehouse_gen](../warehouse_gen.go#L124))

### How it works

##### Invoking the generation

To generate all show outputs (with a cleanup first) run:

```shell
make clean-show-output-schemas generate-show-output-schemas
```

To generate only chosen subset of all objects run:
```shell
make clean-show-output-schemas generate-show-output-schemas SF_TF_GENERATOR_ARGS="--filter-object-names=sdk.Warehouse,sdk.User"
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
- `sdk.ObjectIdentifier` interface

##### To schema mappings

Given SDK struct field can be mapped to the generated schema depending on its type:
- no mapping (`Identity`) - used for `string` and other basic types
- string value mapping (`ToString`) - used e.g. for `time.Time`
- fully qualified name mapping (`FullyQualifiedName`) - used for all identifiers and `sdk.ObjectIdentifier` interface
- casting (`CastToString` and `CastToInt`) - used for enums with underlying type `string` or `int`

##### Changing the SDK object's show output

If you change the show output struct in the SDK:

1. Check if you don't introduce a type that is unsupported (check [supported types](#supported-types)
   and [known limitations](#known-limitations)).
2. Run generation according to [instructions](#invoking-the-generation).

##### Adding a new object to the SDK

1. Add the new show output struct to [sdk_show_result_structs.go](./sdk_show_result_structs.go).
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
- discover a change and generate as part of a `make pre-push`

Implementation improvements:
- (optional) consider different implementations of `Mapper` (e.g. TODO in [schema_field_mapper_test.go](./schema_field_mapper_test.go): `ugly comparison of functions with the current implementation of mapper` and not ideal implementation in the [to_schema_mapper.tmpl](./templates/to_schema_mapper.tmpl): `runMapper .Mapper $nameLowerCase "." .OriginalName`)
