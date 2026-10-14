## ssot

Shared object definitions for generators. Types live in this package; instances in [defs](defs). This package must not import `pkg/sdk` or any meaningful logic package.

Generators consume ssot; they do not all consume `defs.All()` the same way. Resource generation uses every object in `All()`. SHOW/DESCRIBE schema generation and resource config model builders are opt-in per object ([Generators](#generators)). More generators will follow. Ultimately, all the generators should default to using ssot and `defs.All()`.

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

Def files keep attributes as vars (`ComputePoolAttrs`) so `Columns` can point at them instead of restating snake_case strings.

### SHOW and DESCRIBE

`Attributes` cannot drive SHOW/DESCRIBE: more columns than settable fields, CREATE-only fields, renames, mixed Snowflake order.

`OutputSpec.Columns` is membership and Snowflake order. Presence in `Columns` is the only SHOW/DESCRIBE membership.

### OptionalInShow vs OptionalInDb

- `OptionalInShow` - the SHOW/DESCRIBE SDK Go field is a pointer (`*string`, `*AccountObjectIdentifier`). SHOW schema generation prefixes `*` so `ToSchema` emits a nil-guard.
- `OptionalInDb` - the DB row is nullable (`sql.NullString`, …). Unused by SHOW schema generation. Reserved for later SDK paired-struct generation from ssot.

They are independent: a column can be optional in SHOW, in the DB row, both, or neither.

### Generators

`Generators` is per-object enablement, not an implicit `All()` feed:

- `WithSchemaShow()` / `WithSchemaDescribe()` - SHOW/DESCRIBE schema generator (`make generate-show-output-schemas`). See [that README](../../schemas/gen/README.md).
- `WithModel()` - resource config model builders (`make generate-resource-model-builders`). See [that section](../../acceptance/bettertestspoc/README.md#adding-new-resource-config-model-builders). Fields come from `Attributes` and identity, not from [SHOW and DESCRIBE](#show-and-describe) `Columns`. Types are in [attribute.go](attribute.go). An optional bool is stored as a boolean-string; the extra setter takes the ssot type. `fully_qualified_name` is not emitted.
- SDK generator parts are empty for now.

Fill `Show` / `Describe` first; enable the generator parts when the lists are ready. Enabling SHOW or DESCRIBE without `Columns` panics. Enabling `WithModel()` without `Attributes` is allowed when the object only needs identity.

### Later

SDK `ShowOperationWithPairedStructs` / `DescribeOperationWithPairedStructs` still own `sdk.<Name>` / `sdk.<Name>Details`. Until paired structs are generated from ssot, keep those field lists in sync with `Columns`. `OptionalInDb` is the ssot flag that work will consume.
