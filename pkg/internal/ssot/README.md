## ssot

Shared object definitions for generators. Types live in this package; instances in [defs](defs). This package must not
import `pkg/sdk` or any meaningful logic package.

Generators consume ssot; they do not all consume `defs.All()` the same way. Resource generation uses every object in
`All()`. SHOW/DESCRIBE schema generation, resource config model builders, resource assertions, resource show output
assertions, Snowflake object assertions, object and resource parameters assertions, and SDK paired structs are opt-in per object ([Generators](#generators)). More
generators will follow. Ultimately, all the generators should default to using ssot and `defs.All()`.

### Object

`Object` is the written intent for one Snowflake object:

- `Attributes` - CREATE/ALTER sequence (resource schema)
- `Show` / `Describe` - output columns ([SHOW and DESCRIBE](#show-and-describe))
- `Generators` - which generators may consume those lists ([Generators](#generators))
- `ResourceHooks` - resource lifecycle/schema overrides (resource generator)
- `SdkCreateRequestObjectName` - SDK CREATE request type-name stem when it differs from `Name`
- `VariantOf` - base object whose SDK parameter details, SET/UNSET requests, and `ParameterType` names a variant shares; parameter catalog membership still comes from the variant's own `Name`
- `VariantEntityCreateName` - SQL object words a variant CREATE emits. `ObjectType` builds those words from `Name`, which does not match a variant (`WarehouseAdaptive` would be `WAREHOUSE ADAPTIVE`). Other statements keep `ObjectType`
- `SchemaStateUpgraders` - ordered upgrade function names ([Resource generator](#resource-generator))

Identity (`Name`, `IdType`, ...) is separate from `Attributes`. `LegacyIdEncoding` keeps the pipe-delimited resource id
(`db|schema|name`, `helpers.EncodeSnowflakeID`) that older resources store. The generated Create calls it instead of
`helpers.EncodeResourceIdentifier`. `WithParseId(Generated())` emits a named parser used by Import, Read, and Update.
It calls `sdk.Parse<IdType>` for current ids and `helpers.DecodeSnowflakeIDErrLegacy[<IdType>]` for legacy ids.
The default `ExtCall` keeps the handwritten parser.

`VariantOf` keeps SDK parameter naming simple because variants are expected to share their base object's SDK
parameter types and usually the same parameter set. Different sets can already use a variant-specific catalog level
derived from `Name`, as interactive warehouses do. If a variant needs a catalog level other than the one derived from
its `Name`, extend SSOT with an explicit parameter-level override.

### Attribute

`Attribute` is one field. Current types are in [attribute.go](attribute.go).

- `Readonly` - non-writable output column. Must not sit on `Attributes`. CREATE-only fields are writable and absent from
  `Show` / `Describe` `Columns`.
- `SqlOnly` - the attribute stays on a SQL sequence and is left out of schema, resource, model, and assertion generation.
  [`ResourceAttributes`](object.go) is that filtered list. A fixed CREATE value uses [`StaticSQL`](sql_clause.go): the
  field is the unexported Go name and the SQL is the ddl text, such as `WAREHOUSE_TYPE = 'ADAPTIVE'`.
- `ShowName` - SHOW/DESCRIBE key when it differs from the resource name (`for_application` → `application`).
  `OutputKey()` is `ShowName` or `Name`.
- `CreateMapping` - the create helper from `pkg/resources/resource_helpers_create.go` used by a generated
  `ApplyCreateOptionals`. Use a `HelperMapping` constant (`StringAttributeCreateBuilder`, ...) or
  `ListValueCreate("<mapperFunc>")`. A required attribute with no mapping is a `NewCreate<Object>Request` constructor
  argument (int or enum only). The resource identifier is not an attribute. A generated `ParseIdFromConfig` builds it
  from `IdType`: `name` for an account object, `database` and `name` for a database object, and `database`, `schema`,
  and `name` for a schema object.
- `UpdateMapping` - the update helper from `pkg/resources/resource_helpers_update.go` used by a generated
  `ApplySetUnset`. Use a `SetOnlyMapping` constant (`IntAttributeUpdateSetOnlyBuilder`), a `SetUnsetMapping` constant
  (`StringAttributeUpdateBuilder`, ...), or `ListValueUpdate("<mapperFunc>")`. An attribute with no mapping is left out
  of set/unset.
- `NativeBool` - an optional bool stays Terraform `TypeBool`. Without the mark, an optional bool is a boolean string.
- `IncludeAliasedInTheDescription` - the resource description lists enum aliases. Without the mark, the description
  lists canonical values. Warehouse size sets it.
- `StoredAsText` - the db column is text, and the plain kind stays the semantic type. For a bool, the paired-struct
  generator calls `BoolFromText`. A bool without the mark is still db `bool`. Other generators ignore the mark.
  See [PairedStructs](../../sdk/generator/README.md#pairedstructs).
- `DiffSuppress` - a suppressor composed with inferred identifier and enum normalization. Set it directly to
  `ssot.SuppressShow`, `ssot.SuppressDescribe`, or `ssot.SuppressStatement`; combine them with
  `ssot.SuppressIfAny(ssot.SuppressStatement, ssot.SuppressDescribe)`. Nested combinations are flattened, and nil
  means no configured suppression. Output suppressors use `OutputKey()` unless a leaf supplies an override, e.g.
  `ssot.SuppressDescribe.WithKey("definition")`. Output strategies apply to fields, not collection elems.
  Statement suppression uses the resource's `DiffSuppressStatement` helper for strings, including string elems.

Def files keep attributes as vars (`ComputePoolAttrs`) so `Columns` can point at them instead of restating snake_case
strings.

### SHOW and DESCRIBE

`Attributes` cannot drive SHOW/DESCRIBE: more columns than settable fields, CREATE-only fields, renames, mixed Snowflake
order.

`OutputSpec.Columns` is membership and Snowflake order. Presence in `Columns` is the only SHOW/DESCRIBE membership.

Each warehouse variant declares its show in `Show.Columns`. `Warehouse.ShowGeneric` is the temporary union and the one paired SHOW. Next, each variant's `Show.Columns` becomes a paired SHOW with its dedicated type, and `ShowGeneric` goes away.

### OptionalInShow vs OptionalInDb

- `OptionalInShow` - the SHOW/DESCRIBE SDK Go field is a pointer (`*string`, `*AccountObjectIdentifier`). SHOW schema
  generation prefixes `*` so `ToSchema` emits a nil-guard.
- `OptionalInDb` - the DB row is nullable (`sql.NullString`, …). Unused by SHOW schema generation. Live for the SDK
  paired-struct db and plain kinds. See [PairedStructs](../../sdk/generator/README.md#pairedstructs).

They are independent: a column can be optional in SHOW, in the DB row, both, or neither.

### Generators

`Generators` is per-object enablement, not an implicit `All()` feed. Each generator is one entry. SHOW and DESCRIBE are
settings of that entry. `WithSchema()`, `WithSdk()`, `WithShowOutputAssert()`, and `WithObjectAssert()` with no settings
enable nothing.

- `WithSchema(SchemaShow(), SchemaDescribe())` - SHOW/DESCRIBE schema generator (`make generate-show-output-schemas`).
  See [that README](../../schemas/gen/README.md).
- `WithModel()` - resource config model builders (`make generate-resource-model-builders`).
  See [that section](../../acceptance/bettertestspoc/README.md#adding-new-resource-config-model-builders). Fields come
  from `Attributes` and identity, not from [SHOW and DESCRIBE](#show-and-describe) `Columns`. Types are
  in [attribute.go](attribute.go). An optional bool is stored as a boolean-string; the extra setter takes the ssot type.
  `fully_qualified_name` is not emitted.
- `WithResourceAssert()` - resource assertions (`make generate-resource-assertions`).
  See [that section](../../acceptance/bettertestspoc/README.md#adding-new-resource-assertions). Fields come from
  `Attributes` and identity, not from [SHOW and DESCRIBE](#show-and-describe) `Columns`. Types are
  in [attribute.go](attribute.go). An optional bool is stored as a boolean-string; the extra method takes the ssot type.
  `fully_qualified_name` is emitted.
- `WithShowOutputAssert(ShowOutputShow(), ShowOutputDescribe())` - resource show output assertions (
  `make generate-resource-show-output-assertions`).
  See [that section](../../acceptance/bettertestspoc/README.md#adding-new-resource-show-output-assertions). Fields come
  from [SHOW and DESCRIBE](#show-and-describe) `Columns`, not from `Attributes`.
  `ShowOutputShow(ShowOutputShowSkipFields(&attrs.Tables))` copies `SkipFields` from the old show def. The argument is
  the same `*Attribute` that sits in `Show.Columns`. The column stays in `Columns`, and `Has*` is not generated.
  `ShowOutputDescribe(ShowOutputDescribeSkipFields(&attrs.Password))` is the describe-side copy.
- `WithObjectAssert(ObjectAssertShow(), ObjectAssertDescribe())` - Snowflake object assertions (
  `make generate-snowflake-object-assertions`).
  See [that section](../../acceptance/bettertestspoc/README.md#adding-new-snowflake-object-assertions). Fields come
  from [SHOW and DESCRIBE](#show-and-describe) `Columns`, not from `Attributes`. `StoredAsText` is ignored. A bool stays
  `bool`. `ObjectAssertShow(ObjectAssertShowSkipFields(&attrs.Tables))` copies `SkipFields` from the old show def. The
  argument is the same `*Attribute` that sits in `Show.Columns`. The column stays in `Columns`, and `Has*` is not
  generated. `ObjectAssertDescribe(ObjectAssertDescribeSkipFields(&attrs.Password))` is the describe-side copy.
  `ObjectAssertDescribe(WithNoIdentifiableObject())` suppresses the fetch constructor when `Describe` returns a
  slice and the row has no `ID()`. `ObjectAssertShow(ObjectAssertShowFromBase("AsAdaptive"))` fetches through
  `testClient.<VariantOf.Name>.Show` and projects the row with `AsAdaptive()`. `Object.ShowFromBase` resolves that
  client from `VariantOf`. The object type constant is `sdk.ObjectType<VariantOf.Name>`. These are settings of
  object-assertion show or describe, not their own generators.
- `WithObjectParametersAssert()` - Snowflake object parameters assertions (
  `make generate-snowflake-object-parameters-assertions`).
  See [that section](../../acceptance/bettertestspoc/README.md#adding-new-snowflake-object-parameters-assertions).
  Parameters come from the [catalog](defs/parameters.go) entries for the object's `ParameterLevel()`, not from
  `Attributes`. Constants are `sdk.<Name>Parameter`. A variant uses `VariantOf.Name` for the constant and the object
  type. `Object.ParameterAssertionNames` returns those names. The catalog level stays `ParameterLevel()`.
- `WithResourceParametersAssert()` - resource parameters assertions (`make generate-resource-parameters-assertions`).
  See [that section](../../acceptance/bettertestspoc/README.md#adding-new-resource-parameters-assertions).
  Parameters come from the same catalog entries, and constants are `sdk.<Name>Parameter`. The two switches are
  independent. Leave the handwritten entry while only one of them is opted in.
- `WithSdk(SdkShow(), SdkDescribe(), SdkDrop(), SdkCreate(), SdkAlter())` - SDK paired structs, DROP, CREATE, and ALTER (`make generate-sdk`).
  See [PairedStructs](../../sdk/generator/README.md#pairedstructs). Show and Describe fields come
  from [SHOW and DESCRIBE](#show-and-describe) `Columns`, not from `Attributes`. `OptionalInDb` is live for the mapped
  db and plain kinds. `StoredAsText` selects `BoolFromText` for a bool. A bool without the mark is still db `bool`.
  Other generators ignore the mark. See [Attribute](#attribute). `SdkDrop()` builds DROP from `Name`: `DROP`, the
  uppercase words, `IF EXISTS`, and the name. `SdkDropSafelyHook()`, `SdkDropSafelyForce()`, and
  `SdkDropCascadeRestrict()` are settings of that call. `SdkDropCascadeRestrict()` adds optional `CASCADE` and
  `RESTRICT` and marks those fields as conflicting. The hook body stays in the SDK ext file. `SdkCreate()` builds
  CREATE from `Name`: `CREATE`, optional `OR REPLACE`, `Create.KeywordsBefore`, the uppercase words, `IF NOT EXISTS`,
  and the name, then `Create.Clauses`. `SdkCreateWithoutOrReplace()` and `SdkCreateWithoutIfNotExists()` leave those
  clauses out. When both clauses are
  present they are marked conflicting. `SdkCreateAdditionalValidations()` adds that validation call. `SdkAlter()` builds
  ALTER from `Name`: `ALTER`, the uppercase words, `IF EXISTS`, and the name, then [`Alter`](object.go) in list order.
  `SdkAlterWithoutIfExists()` and `SdkAlterWithoutExactlyOneValueSet()` leave those pieces out.
  `SdkAlterSetAdditionalValidations()` marks the `SET` struct. `SdkAlterSetList()` makes that `SET` a list with no
  parentheses. Without it, `SET` is a keyword. `Create` is a [`CreateSpec`](object.go). `Alter` is an
  [`AlterSpec`](object.go): action keywords, `RenameTo`, `SwapWith`, `SET`, `UNSET`, list assignments, `SET TAG`, and
  `UNSET TAG` are entries, so a keyword can follow the tag clauses. `AlterKeywordNotExactlyOne` leaves a keyword out of
  `ExactlyOneValueSet`. `AlterList` is a parentheses list with no equals, one exactly-one branch, and `ValidIdentifier`
  on each element. `SdkAlterAdditionalValidations()` marks the ALTER query struct. `RenameTo` and `SwapWith` use the
  object's `IdType` and get `ValidIdentifierIfSet`. `Set` and `Unset` are the fields inside those blocks. [`Tags`](attribute.go) is the
  optional `TAG` clause and belongs on `Create.Clauses`, not on `KeywordsBefore`. [`Parameters`](attribute.go) spreads
  the parameter catalog for `ParameterLevel()`. `AlterAdditionalParametersForLevels` adds parameters from more levels
  onto SET and UNSET. The union is emitted in catalog order, and a name already present is kept once. CREATE does not
  take those levels.
  CREATE and SET use `WithParameters`. UNSET uses `WithParametersUnset`. `AtLeastOneValueSet` lists those parameter
  fields first, then the other fields in the struct.
  An identifier parameter gets `ValidIdentifierIfSet` on CREATE and SET. Comment stays an attribute clause beside that
  entry. Clause shape is `CreateSql`, `SetSql`,
  and `UnsetSql` on [Attribute](attribute.go). A zero create or set clause uses the attribute type. `Clause` applies
  functional options on top of that default. A zero unset clause is a bool keyword. Leave the attribute out of `Unset`
  to omit it, or set `UnsetSql` to a parameter clause to unset by setting a value. A list whose element is an enum uses
  an item query struct named after that enum, with `ListItem` appended. Its field is `Value`, a required single-quoted
  keyword. When `SdkCreate()` or `SdkAlter()` is off, that call stays in the def. `SdkCreateForVariants()` emits
  each variant CREATE as a custom method. The variant fills `Create` (`Doc`, clauses, and
  `AdditionalValidations`). The SQL object words come from `VariantEntityCreateName`, not from `Create.Sql` and not from
  `ObjectType`. A fixed value is an `SqlOnly` attribute with `StaticSQL` in that clause list. The method name is `Create` plus the variant name with the base name removed. The query
  struct name is `Create` plus `SdkCreateRequestObjectName`.
  TODO [next PRs]: the same switch for variant ALTER, and for SHOW/DESCRIBE. `SdkCustomOperations()` emits
  [`CustomOperations`](custom_operation.go) from one call, in slice order. Place that call after the dedicated operation
  methods and before custom interface methods. The generated interface lists every operation, then the custom interface
  methods, so a custom operation stays with the operations. Reordering in the generated SDK is expected while objects
  are aligned. When the setting is off, the custom operation stays in the SDK def. `SdkCustomInterfaceMethods()` emits
  [`CustomInterfaceMethods`](custom_interface_method.go) from one call, in slice order. Each entry is an interface
  declaration. The body stays in the SDK ext file. `ctx` is always first and is not listed. A nil parameter slice is
  `id` of `IdType`. A non-nil slice is the parameter list. An empty parameter `Kind` uses `IdType`. An empty return
  slice is `error`. An empty doc emits no comment.
  When the setting is off, the method stays in the SDK def. `SdkShowParameters()` emits `ShowParameters` as that kind of
  declaration: an empty doc, `id` of `IdType`, and returns `[]*Parameter` and `error`. The body stays in the SDK ext
  file. `SdkShowParametersDetails()` emits `ShowParametersDetails` and the `{Name}ParametersDetails` struct. The fields
  come from `ParameterDefsForLevel` for the object's `ParameterLevel()`. The generated body calls `ShowParameters`.
  When either setting is off, that call stays in the SDK def. `SdkDefineFromSsot()` defines this object's SDK
  interface from this ssot object alone. The interface contains the operation parts whose flags are on. An object
  without it stays on the handwritten SDK def.

Fill `Show` / `Describe` first; enable the matching SHOW or DESCRIBE setting when the lists are ready. Enabling SHOW or
DESCRIBE schema generation, show output assertions, Snowflake object assertions, or SDK paired structs without `Columns`
panics. Enabling `WithModel()` or `WithResourceAssert()` without `Attributes` is allowed when the object only needs
identity.

### Assumptions

Rules a generator may rely on. The calls are in [Generators](#generators). Types are in [attribute.go](attribute.go).

#### Imports

`pkg/internal/ssot` may import `pkg/internal/genhelpers` and other `pkg/internal` helpers. [defs](defs) may import
`pkg/internal/ssot` only.

ssot must not import `pkg/sdk`, `pkg/resources`, `pkg/schemas`, `pkg/acceptance`, or any generator package. Generation
behavior stays in the generator. ssot holds the written object. `IdType` has no `sdk.` prefix.

Generator defs may import `ssot/defs` as `ssotdefs`. A generator may import `ssot` and `ssot/defs`. That import points
one way.

#### Resource generator

[`defs.All()`](defs/all.go) is the set of objects in ssot. The resource generator (`make generate-resource`) walks that
slice with no `Generators` flag. Every object in `All()` is a resource-generation input. That is the exception to
opt-in.

A subtemplate calls the ext function by default (`ExtCall`). Set that hook to `Generated` to take the step from the
template. Until the subtemplate body is written, generation panics through `Unimplemented`. Each hook is switched on its
own, with its `With` option. See [resource_hooks.go](resource_hooks.go). `DropSafelyInSdk` defaults to omitted, which
keeps `DeleteExt`. Set it to `Generated` to emit `Delete` and point `DeleteContext` at that function. `Delete` calls
`ResourceDeleteContextFunc` with `sdk.Parse` plus the object `IdType`.

`AfterCreate` and `AfterAlter` are `Omitted` by default. They run after `SetId` and after the alters, and their
diagnostics (usually warnings) are added to the final Read. Use them for a step that must not drop the object from
state when it fails, such as resuming a created object. `Import<Name>` is generated only when `ShowByIdInSdk` is
present. A resource with a passthrough importer omits it.

`WithCreateInSdk(Generated())` calls `Client.<NamePlural>.Create(ctx, req)` after applying the create fields and
parameters. A failure returns `diag.FromErr` before setting the resource id or reading the object.
Use `ExtCall` when an object needs a different SDK method or custom error diagnostics.

`SchemaStateUpgraders` on the object is an ordered list of upgrade function names. `WithStateUpgraders` is the
`ResourceHook` for that list and for the schema version. `Generated` emits `SchemaVersion` as the length and one
`StateUpgrader` per name. `Version` is the index. `Type` is `cty.EmptyObject`. The functions stay in package
resources. `ExtCall` keeps the ext `SchemaVersion` and `StateUpgraders` pair. `ValidateObject` panics when
`Generated` has no function names, so the template does not check the list.

`CustomizeDiff` defaults to `ExtCall`, which keeps `<prefix>CustomizeDiffExt` as the whole `schema.CustomizeDiffFunc`.
`WithCustomizeDiff(Generated())` emits `<prefix>CustomizeDiff` and points `CustomizeDiff` at it. The value is
`TrackingCustomDiffWrapper` around `customdiff.All`. Call order is `show_output`, `describe_output`, the parameters
computed-if, `fully_qualified_name`, then `<prefix>ParametersCustomDiff` when that hook is present.
`AdditionalCustomizeDiffs` on the object is an ordered list of `schema.CustomizeDiffFunc` names in package resources.
Those names are appended after the generated calls. The function bodies stay in the ext file. `ValidateObject` panics
when the list is non-empty without `WithCustomizeDiff(Generated())`, when an entry is empty, or when an entry is not
a function name. An output call is emitted only when that output's schema
hook is present and the call has changed keys. Its schema argument is `<prefix>Schema`. The parameters computed-if is
emitted when both `ParametersAttributesSchema` and `ParametersFieldNames` are present. Its schema argument is the
parameters attributes schema, and its changed keys are `<prefix>ParametersFieldNames...`. `show_output` and
`describe_output` take the top-level columns of that output, in column order. A column is a changed key when its
`Name` is a top-level resource or identifier attribute. `Children` marks a nested column or a nested resource
attribute, and both are left out. For now, the set/list fields have to be excluded
(`TODO [SNOW-1648997]`): a column or resource attribute whose type is a set or a list is left out, including a set or
list of scalars. A readonly column that is not a resource attribute is left out: a diff cannot change it. The key is
`Name`, the resource attribute `ComputedIfAnyAttributeChanged` watches, including when `ShowName` differs.
`fully_qualified_name` is emitted only when at least one identifier attribute is not ForceNew. Its keys are those
attributes, in `IdentityFields` order: `database`, then `schema`, then `name`. `name` is ForceNew when
`IdentityNameForceNew` is set. `database` and `schema` are ForceNew together when `IdentityHierarchyForceNew` is set.
The identity schema template does not emit that yet. An unknown `IdType` panics.

Optional bools are boolean-strings on purpose (`applyResourceAttributeRules`
in [resource_attribute_rules.go](../../resources/gen/resource_attribute_rules.go)): an optional bool is Terraform
`TypeString`, a required bool is `TypeBool`. That classification lives in [classified_field.go](classified_field.go)
and applies to the resource schema only.

#### Other generators

Every plugged generator other than the resource generator is opt-in. It walks `defs.All()`, calls `ObjectGenerators`,
and reads that generator's bools. An object with the flag off stays on the generator's old input. Setting the flag on
nothing must not change generated files. Do not feed `All()` in as "generate every object" for those generators.
Generators default to ssot and `All()` only once the definitions cover the objects. The calls are
in [Generators](#generators).

`SchemaShow`, `SchemaDescribe`, `SdkShow`, `SdkDescribe`, `SdkDrop`, `SdkCreate`, `SdkAlter`, `ShowOutputShow`,
`ShowOutputDescribe`, `ObjectAssertShow`, and `ObjectAssertDescribe` are not `GeneratorOption`. They cannot sit alone in
`Generators`. A bare `WithSchema()`, `WithSdk()`, `WithShowOutputAssert()`, or `WithObjectAssert()` enables nothing.
`SdkDrop()` with no options still sets `Drop`. `SdkCreate()` with no options still sets `Create`. `SdkAlter()` with no
options still sets `Alter`. `ShowOutputShow()` with no options sets `Show` and leaves `ShowSkipFields` nil.
`ShowOutputDescribe()` with no options sets `Describe` and leaves `DescribeSkipFields` nil. `ObjectAssertShow()` with no
options sets `Show` and leaves `ShowSkipFields` nil and `ShowProjection` empty. `ObjectAssertShowFromBase` sets `ShowProjection` on that same call. The client name is `VariantOf.Name`. `ObjectAssertDescribe()` with no options sets `Describe` and leaves
`NoIdentifiableObject` false and `DescribeSkipFields` nil. `WithModel()` and `WithResourceAssert()` take a variadic of
`ModelOption` and `ResourceAssertOption`. No constructors exist for those option types. Zero options still sets
`Enabled`. A later setting is another option on that call, not a new `Generators` entry. Consumers read the fields
`Show`, `Describe`, `Drop`, `Create`, `Alter`, `Enabled`, `NoIdentifiableObject`, `ShowSkipFields`, and
`DescribeSkipFields`.

Do not apply resource-attribute classification. A SHOW bool is `bool`. A paired bool without `StoredAsText` is db
`bool`.

Marks belong to the generator that defined them. `StoredAsText` is a paired-struct db mark. `OptionalInShow` is the
SHOW/DESCRIBE pointer. `OptionalInDb` is the db row. `OptionalInDb` is unused by SHOW schema generation. They are
independent. See [OptionalInShow vs OptionalInDb](#optionalinshow-vs-optionalindb).

An unmapped kind panics. Do not guess a mapping. Tests build a fake `ssot.Object`. They do not walk `defs.All()` to
assert which live objects are opted in. Tests of clause shape, emission, and validation use that fake object. They do
not assert those paths against a live definition such as Compute Pool or Schema.

The old path for an object that is not opted in stays byte-identical. Do not edit `*_gen.go`.

#### Statement model

Rules for CREATE and ALTER. The calls are in [Generators](#generators).

- Order is membership in a slice, and the slice order is the emission order. SHOW and DESCRIBE use `Columns`. CREATE uses [`CreateSpec`](object.go). ALTER uses [`AlterSpec`](object.go). `Set` and `Unset` are the fields inside those blocks.
- A new position or switch on a statement is a field of that statement's spec. `KeywordsBefore` is the CREATE field for clauses between `OR REPLACE` and the object words.
- `KeywordsBefore` holds attributes. `TRANSIENT` is the `is_transient` attribute. An alter keyword that is not an attribute, such as `RESUME` or `STOP ALL`, is an `AlterKeyword` entry on `Alter.Clauses`.
- A SQL constructor does not take an `Attribute`. Set the clause on the attribute. `ValidateObject` checks it while walking the object.
- The SDK generator calls the existing query-struct builders (`NumberAssignment`, `OptionalSQL`, `ListAssignment`, and the rest). It does not record a Go type on the clause, and it does not build `ddl` tags.
- Clauses Snowflake statements usually accept are on unless the object opts out: `OR REPLACE`, `IF NOT EXISTS`, `IF EXISTS`, `SET`, `UNSET`, `SET TAG`, `UNSET TAG`, and the exactly-one alter check. `SET` is a keyword. `SdkAlterSetList` opts into `list,no_parentheses`.
- A custom operation is an entry on [`CustomOperations`](custom_operation.go). `SdkCustomOperations()` emits that slice from one call, in slice order. Clone and undrop are custom operations. Generating CREATE or ALTER does not take them over. Place that call after the dedicated operation methods and before custom interface methods. The generated interface lists every operation, then the custom interface methods, so a custom operation stays with the operations. Reordering in the generated SDK is expected while objects are aligned.
- A custom interface method is an entry on [`CustomInterfaceMethods`](custom_interface_method.go). `SdkCustomInterfaceMethods()` emits that slice from one call, in slice order. The declaration is generated. The body stays in the SDK ext file. `ctx` is always first and is not listed. A nil parameter slice is `id` of `IdType`. A non-nil slice is the parameter list, and an empty parameter `Kind` uses `IdType`. An empty return slice is `error`. An empty doc emits no comment. A method with another parameter list, return list, or comment sets those fields. When the setting is off, the method stays in the SDK def.
- `ShowParameters` is that same declaration with a fixed name, an empty doc, `id` of `IdType`, and returns `[]*Parameter` and `error`. The body stays in the SDK ext file. `SdkShowParameters()` enables it. The object does not pass the identifier kind.
- `ShowParametersDetails` fills the details struct from `ParameterDefsForLevel` for `ParameterLevel()`, which is derived from `Name`. The struct name is `Name` plus `ParametersDetails`. The method has an empty doc, takes `id` of the interface identifier kind, and returns a pointer to that struct and `error`. The body is generated and calls `ShowParameters`. `SdkShowParametersDetails()` enables it. When the setting is off, the call stays in the SDK def.
- One generator fits every object. When an object's statement does not match the current emission, extend the spec.
- A repeated generated name has one spelling. An enum list item is `{Enum}ListItem` with the field `Value`. Rename the SDK type to that spelling. The definition does not pass the struct name or the field name.

#### Incremental coverage

Object and attribute configuration belongs in ssot, including which generators run for an object. The definitions do
not cover that yet. One change that moves every object is the wrong shape: there are many objects, and some mappings
are still unsolved. `All()`, per-object `Generators`, and the old input lists are there so one object or one generator
can move while another change does the same for a different object.

Until a piece has moved, do not assume it is already in ssot:

- An object absent from `All()` still enters through the generator's old list.
- Parameters are not `Attributes` yet. A model builder or resource assertion that needs them still joins
  `ParameterDefsForLevel`.
- When `SdkShow()` or `SdkDescribe()` is on, that Show or Describe method comes from the option. The query struct,
  ShowByID filtering, and `DescriptionMappingKind` are part of it. An object with the setting off keeps the handwritten
  call. When `SdkDrop()` is on, that Drop method comes from the option. An object with the setting off keeps the
  handwritten call. When `SdkCreate()` or `SdkAlter()` is on, that Create or Alter method comes from the option. An
  object with the setting off keeps the handwritten call. A zero create or set clause uses the attribute type. A zero unset clause is a bool keyword. Clause shape is `CreateSql`, `SetSql`, and `UnsetSql`
  ([Attribute](attribute.go)). `Create` is a [`CreateSpec`](object.go). `Set` and `Unset` are the statement order. A
  `Parameters` entry spreads the catalog on that statement. `Alter` is the order of the ALTER blocks after the name.
  When `SdkCustomOperations()` is on, those operations come from `CustomOperations` and are emitted together. When it is
  off, they stay in the SDK def. When `SdkCustomInterfaceMethods()` is on, those declarations come from
  `CustomInterfaceMethods`. The body stays in the SDK ext file. When it is off, the method stays in the SDK def. When
  `SdkShowParameters()` or `SdkShowParametersDetails()` is on, that method comes from the option. The details fields
  come from `ParameterDefsForLevel` for `ParameterLevel()`. When the setting is off, the call stays in the SDK def.
  Every other operation stays in the SDK def. They are not
  generated from `Columns`.
- Schema `SkipFields`, `ManualFields`, `TypeOverrides`, and `UsedAsListEntry` have no ssot equivalent yet. Object-assertion
  `SkipFields` is `ObjectAssertShowSkipFields` and `ObjectAssertDescribeSkipFields`. Show-output `SkipFields` is
  `ShowOutputShowSkipFields` and `ShowOutputDescribeSkipFields`.
- `ResourceHooks` do not yet generate the whole lifecycle, the SDK calls, or the mapping. A missing step stays an ext
  call, or panics through `Unimplemented`, until that hook is implemented. See [Resource generator](#resource-generator).
- Datasource model builders, provider model builders, resource parameter assertions, and Snowflake object parameter
  assertions do not read `Generators` yet.
- `g.Enum` in the SDK generator is a type alias of the ssot enum type ([enum.go](../../sdk/generator/gen/enum.go)). The
  alias does not mean every enum is defined in ssot yet.

These are the model, not missing coverage. Identity (`Name`, `IdType`, `IdentityNameForceNew`) is not an `Attribute`.
A `Readonly` column must not sit on `Attributes`. It may still appear in `Set` or `Unset` when that SHOW column is also
an ALTER field. A field in `Create`, `Set`, or `Unset` may be absent from `Attributes` when it is not a resource
attribute. A CREATE-only field is an `Attribute` and is absent from `Columns`.
A paired-field option that ssot does not map panics. There is no silent fallback. `StoredAsText` is mapped for a bool
with neither optional flag. Anything else panics. A list of enums is plain `[]string` in the SDK struct, not a slice of
the enum type.

#### Limitations

- A list that is not a list of enums takes the item query struct name as a string. That struct stays in the SDK def.
- A list of enums is built in the SDK generator: struct name `{Enum}ListItem`, field `Value`, required single-quoted keyword. ssot should hold that description, and the generator should only translate it into the query struct. Do that together with more complex list items, not for this enum case alone.
- Computed fields, `Deprecated`, and an explicit schema `Default` are not attribute fields. `NativeBool`
  keeps an optional bool as `TypeBool` and does not emit `Default`. `Min` and `Max` emit `ValidateDiagFunc`.
  Chrome the generator does not emit stays in the attributes modification function.

#### Why Show and Describe columns are explicit

`Attributes` is the resource schema sequence. `Create.KeywordsBefore` is the CREATE order between `OR REPLACE` and the
object words. `Create.Clauses`, `Set`, and `Unset` are the SQL field order inside those statements. `Alter.Clauses` is the ALTER order
after the name.
`OutputSpec.Columns` is membership and Snowflake order for one operation.
Presence in that slice is the only membership. The same `Attribute` var is pointed at from `Columns` so the snake_case
name is not restated. `ShowName` is the output key when it differs from `Name`. `OutputKey()` is `ShowName` or `Name`.
See [SHOW and DESCRIBE](#show-and-describe).

`Attributes` cannot drive those lists:

- The output has columns that are not settable.
- Some settable fields never appear in the output.
- Some settable fields appear under another key.
- DESCRIBE is not SHOW plus a flag. The two sets differ.
- Order in `Columns` is the order Snowflake returns and the order generation emits. It is not the `Attributes` order.

A generator that needs the paired SHOW reads `PairedShowColumns`: `ShowGeneric` when that list is set, otherwise `Show.Columns`. A generator that needs DESCRIBE reads `Describe.Columns`. It does not
derive one from the other, and it does not derive either from `Attributes`.

#### Generators do not depend on each other

A generator reads ssot. It does not call another generator's adapter, field mapping, or templates. Naming a generated
schema symbol is allowed: the resource generator uses `schemas.ShowXSchema` and `schemas.DescribeXDetailsSchema`. That
is the published output of the schema generator, not an import of `pkg/schemas/gen`.

These calls are still there. Each one is removed on its own. Do not add another.

- Resource config model builders call `GetResourceSchemaDetails` from the resource assertion generator
  ([ssot_model_plugs.go](../../acceptance/bettertestspoc/config/model/gen/ssot_model_plugs.go)). That function reflects
  `allResourceSchemaDefs`.
- Resource show output assertions call `GetSdkObjectDetails` from the Snowflake object assertion generator.
  `GetFilteredSdkObjectDetails`
  in [sdk_object_def.go](../../acceptance/bettertestspoc/assert/resourceshowoutputassert/gen/sdk_object_def.go) does that
  call, then drops structs that are not resources. `GetSdkObjectDetails` reflects `allStructs`.
- Resource parameter assertions call `GetAllSnowflakeObjectParameters` from the Snowflake object parameter assertion
  generator
  ([main.go](../../acceptance/bettertestspoc/assert/resourceparametersassert/gen/main/main.go)).
- Datasource model builders and provider model builders call `ModelFromResourceSchemaDetails` from the resource model
  generator
  ([datasource main.go](../../acceptance/bettertestspoc/config/datasourcemodel/gen/main/main.go),
  [provider main.go](../../acceptance/bettertestspoc/config/providermodel/gen/main/main.go)).

If two generators must see the same objects that are still on an old list, that is the shared-list case below, not a
reason to call the other generator.

#### One old list per object

An object is entered through one list. Opting it into ssot removes it from that list.

[`SdkShowResultStructs`](../../schemas/gen/sdk_show_result_structs.go) has no skip. `GetShowResultSchemaDetails` appends
the reflected structs and the ssot details, then panics on a repeated name. The struct leaves the slice in the same
change that enables `SchemaShow` or `SchemaDescribe`.

A shared list keeps the entry until every generator that reads it emits that object from ssot. The merge skips a name
ssot already emits, so generation stays correct while another reader still reflects it. After the last reader has moved,
the old list must not contain the object.

| Shared list             | Readers                                                                                                              | Skip                                                                      |
|-------------------------|----------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------|
| `allStructs`            | Snowflake object assertions (`GetSdkObjectDetails`) and show output assertions (`GetFilteredSdkObjectDetails`)       | `mergeSnowflakeObjectAssertionDetails`, `mergeShowOutputAssertionDetails` |
| `allResourceSchemaDefs` | Resource assertions (`GetResourceSchemaDetails`) and model builders (the same function, called from the model plugs) | the name skip in `ResourceAssertionModels` and `ModelBuilderModels`       |

Show output does not emit `sdk.<Name>Details` unless `ShowOutputDescribe` is on. Do not leave that name in `allStructs`
when show output describe is off.

The skip is not permission to leave the object on the list forever. A repeated ssot `ObjectName()` panics in
`mergeSnowflakeObjectAssertionDetails`, `mergeShowOutputAssertionDetails`, and `GetShowResultSchemaDetails`.
`ResourceAssertionModels` and `ModelBuilderModels` skip a reflected name and do not panic on a repeated ssot name. A
repeat for an object that ssot does not emit stays. Do not delete it.

SDK interfaces come from two lists. The old list is the handwritten defs. ssot adds an interface when `SdkDefineFromSsot()`
is on and `NamePlural` is unique in [`defs.All()`](defs/all.go). That option defines the interface from this ssot
object alone. It contains the operation parts whose flags are on. A part whose flag is off is absent. The SDK generator may then change that
interface through an override keyed by the ssot object `Name`. A shared `NamePlural`, such as the catalog
integration variants, does not get its own interface. The merge skips an old name that ssot emits.
`ShowOperationFromObject` and the other `FromObject` calls stay the translators for one operation. They panic when that
operation's flag is off.

### Parameters

Parameter types live in [parameter.go](parameter.go). Catalog entries live in [defs/parameters.go](defs/parameters.go).

`AccountLevels` and `SessionLevels` are the fixed hierarchies from the [parameter hierarchy](https://docs.snowflake.com/en/sql-reference/parameters#parameter-hierarchy-and-types): account, and account plus user and session. `ObjectLevels` adds one or more object levels on top of the account levels. Provider-only consumer sets (`ACCOUNT_EXT`, `HYBRID_TABLE`, warehouse variants) stay `ParameterLevel` values inside those sets.

### Later

The SDK field chain comes from `Columns` when `SdkShow()` or `SdkDescribe()` is passed to `WithSdk()`. That option also
supplies the Show or Describe method: the query struct, ShowByID filtering, and `DescriptionMappingKind`. A handwritten
field chain, and a handwritten Show or Describe call, are still the source when the matching setting is off. Keep that
chain in sync with `Columns`.
See [PairedStructs](../../sdk/generator/README.md#pairedstructs).
