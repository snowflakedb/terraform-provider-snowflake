## ssot

Shared object definitions for generators. Types live in this package; instances in [defs](defs). This package must not import `pkg/sdk` or any meaningful logic package.

Generators consume ssot; they do not all consume `defs.All()` the same way. Resource generation uses every object in `All()`. SHOW/DESCRIBE schema generation, resource config model builders, resource assertions, resource show output assertions, and SDK paired structs are opt-in per object ([Generators](#generators)). More generators will follow. Ultimately, all the generators should default to using ssot and `defs.All()`.

### Object

`Object` is the written intent for one Snowflake object:

- `Attributes` - CREATE/ALTER sequence (resource schema)
- `Show` / `Describe` - output columns ([SHOW and DESCRIBE](#show-and-describe))
- `Generators` - which generators may consume those lists ([Generators](#generators))
- `ResourceHooks` - resource lifecycle/schema overrides (resource generator)

Identity (`Name`, `IdType`, ...) is separate from `Attributes`.

### Attribute

`Attribute` is one field. Current types are in [attribute.go](attribute.go).

- `Readonly` - non-writable output column. Must not sit on `Attributes`. CREATE-only fields are writable and absent from `Show` / `Describe` `Columns`.
- `ShowName` - SHOW/DESCRIBE key when it differs from the resource name (`for_application` → `application`). `OutputKey()` is `ShowName` or `Name`.
- `StoredAsText` - the db column is text, and the plain kind stays the semantic type. For a bool, the paired-struct generator calls `BoolFromText`. A bool without the mark is still db `bool`. Other generators ignore the mark. See [PairedStructs](../../sdk/generator/README.md#pairedstructs).

Def files keep attributes as vars (`ComputePoolAttrs`) so `Columns` can point at them instead of restating snake_case strings.

### SHOW and DESCRIBE

`Attributes` cannot drive SHOW/DESCRIBE: more columns than settable fields, CREATE-only fields, renames, mixed Snowflake order.

`OutputSpec.Columns` is membership and Snowflake order. Presence in `Columns` is the only SHOW/DESCRIBE membership.

### OptionalInShow vs OptionalInDb

- `OptionalInShow` - the SHOW/DESCRIBE SDK Go field is a pointer (`*string`, `*AccountObjectIdentifier`). SHOW schema generation prefixes `*` so `ToSchema` emits a nil-guard.
- `OptionalInDb` - the DB row is nullable (`sql.NullString`, …). Unused by SHOW schema generation. Live for the SDK paired-struct db and plain kinds. See [PairedStructs](../../sdk/generator/README.md#pairedstructs).

They are independent: a column can be optional in SHOW, in the DB row, both, or neither.

### Generators

`Generators` is per-object enablement, not an implicit `All()` feed:

- `WithSchemaShow()` / `WithSchemaDescribe()` - SHOW/DESCRIBE schema generator (`make generate-show-output-schemas`). See [that README](../../schemas/gen/README.md).
- `WithModel()` - resource config model builders (`make generate-resource-model-builders`). See [that section](../../acceptance/bettertestspoc/README.md#adding-new-resource-config-model-builders). Fields come from `Attributes` and identity, not from [SHOW and DESCRIBE](#show-and-describe) `Columns`. Types are in [attribute.go](attribute.go). An optional bool is stored as a boolean-string; the extra setter takes the ssot type. `fully_qualified_name` is not emitted.
- `WithResourceAssert()` - resource assertions (`make generate-resource-assertions`). See [that section](../../acceptance/bettertestspoc/README.md#adding-new-resource-assertions). Fields come from `Attributes` and identity, not from [SHOW and DESCRIBE](#show-and-describe) `Columns`. Types are in [attribute.go](attribute.go). An optional bool is stored as a boolean-string; the extra method takes the ssot type. `fully_qualified_name` is emitted.
- `WithShowOutputAssert()` / `WithDescribeOutputAssert()` - resource show output assertions (`make generate-resource-show-output-assertions`). See [that section](../../acceptance/bettertestspoc/README.md#adding-new-resource-show-output-assertions). Fields come from [SHOW and DESCRIBE](#show-and-describe) `Columns`, not from `Attributes`.
- `WithSdkShow()` / `WithSdkDescribe()` - SDK paired structs (`make generate-sdk`). See [PairedStructs](../../sdk/generator/README.md#pairedstructs). Fields come from [SHOW and DESCRIBE](#show-and-describe) `Columns`, not from `Attributes`. `OptionalInDb` is live for the mapped db and plain kinds. `StoredAsText` selects `BoolFromText` for a bool. A bool without the mark is still db `bool`. Other generators ignore the mark. See [Attribute](#attribute).

Fill `Show` / `Describe` first; enable the generator parts when the lists are ready. Enabling SHOW or DESCRIBE schema generation, show output assertions, or SDK paired structs without `Columns` panics. Enabling `WithModel()` or `WithResourceAssert()` without `Attributes` is allowed when the object only needs identity.

### Later

The SDK field chain can come from `Columns` when `WithSdkShow()` or `WithSdkDescribe()` is on. The def keeps the query struct, ShowByID filtering, and `DescriptionMappingKind`. A handwritten field chain is still the source when the matching option is off; keep that chain in sync with `Columns`. See [PairedStructs](../../sdk/generator/README.md#pairedstructs).
