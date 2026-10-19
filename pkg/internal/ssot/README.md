## ssot

Shared object definitions for generators. Types live in this package; instances in [defs](defs). This package must not
import `pkg/sdk` or any meaningful logic package.

Generators consume ssot; they do not all consume `defs.All()` the same way. Resource generation uses every object in
`All()`. SHOW/DESCRIBE schema generation, resource config model builders, resource assertions, resource show output
assertions, Snowflake object assertions, and SDK paired structs are opt-in per object ([Generators](#generators)). More
generators will follow. Ultimately, all the generators should default to using ssot and `defs.All()`.

### Object

`Object` is the written intent for one Snowflake object:

- `Attributes` - CREATE/ALTER sequence (resource schema)
- `Show` / `Describe` - output columns ([SHOW and DESCRIBE](#show-and-describe))
- `Generators` - which generators may consume those lists ([Generators](#generators))
- `ResourceHooks` - resource lifecycle/schema overrides (resource generator)

Identity (`Name`, `IdType`, ...) is separate from `Attributes`.

### Attribute

`Attribute` is one field. Current types are in [attribute.go](attribute.go).

- `Readonly` - non-writable output column. Must not sit on `Attributes`. CREATE-only fields are writable and absent from
  `Show` / `Describe` `Columns`.
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
- `StoredAsText` - the db column is text, and the plain kind stays the semantic type. For a bool, the paired-struct
  generator calls `BoolFromText`. A bool without the mark is still db `bool`. Other generators ignore the mark.
  See [PairedStructs](../../sdk/generator/README.md#pairedstructs).

Def files keep attributes as vars (`ComputePoolAttrs`) so `Columns` can point at them instead of restating snake_case
strings.

### SHOW and DESCRIBE

`Attributes` cannot drive SHOW/DESCRIBE: more columns than settable fields, CREATE-only fields, renames, mixed Snowflake
order.

`OutputSpec.Columns` is membership and Snowflake order. Presence in `Columns` is the only SHOW/DESCRIBE membership.

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
- `WithObjectAssert(ObjectAssertShow(), ObjectAssertDescribe())` - Snowflake object assertions (
  `make generate-snowflake-object-assertions`).
  See [that section](../../acceptance/bettertestspoc/README.md#adding-new-snowflake-object-assertions). Fields come
  from [SHOW and DESCRIBE](#show-and-describe) `Columns`, not from `Attributes`. `StoredAsText` is ignored. A bool stays
  `bool`. `ObjectAssertDescribe(WithNoIdentifiableObject())` suppresses the fetch constructor when `Describe` returns a
  slice and the row has no `ID()`. `WithNoIdentifiableObject()` is a setting of object-assertion describe, not its own
  generator.
- `WithSdk(SdkShow(), SdkDescribe())` - SDK paired structs (`make generate-sdk`).
  See [PairedStructs](../../sdk/generator/README.md#pairedstructs). Fields come
  from [SHOW and DESCRIBE](#show-and-describe) `Columns`, not from `Attributes`. `OptionalInDb` is live for the mapped
  db and plain kinds. `StoredAsText` selects `BoolFromText` for a bool. A bool without the mark is still db `bool`.
  Other generators ignore the mark. See [Attribute](#attribute).

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
own, with its `With` option. See [resource_hooks.go](resource_hooks.go).

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

`SchemaShow`, `SchemaDescribe`, `SdkShow`, `SdkDescribe`, `ShowOutputShow`, `ShowOutputDescribe`, `ObjectAssertShow`,
and `ObjectAssertDescribe` are not `GeneratorOption`. They cannot sit alone in `Generators`. A bare `WithSchema()`,
`WithSdk()`, `WithShowOutputAssert()`, or `WithObjectAssert()` is not both. `ObjectAssertDescribe()` with no options
sets `Describe` and leaves `NoIdentifiableObject` false. `WithModel()` and `WithResourceAssert()` take a variadic of
`ModelOption` and `ResourceAssertOption`. No constructors exist for those option types. Zero options still sets
`Enabled`. A later setting is another option on that call, not a new `Generators` entry. Consumers read the fields
`Show`, `Describe`, `Enabled`, and `NoIdentifiableObject`.

Do not apply resource-attribute classification. A SHOW bool is `bool`. A paired bool without `StoredAsText` is db
`bool`.

Marks belong to the generator that defined them. `StoredAsText` is a paired-struct db mark. `OptionalInShow` is the
SHOW/DESCRIBE pointer. `OptionalInDb` is the db row. `OptionalInDb` is unused by SHOW schema generation. They are
independent. See [OptionalInShow vs OptionalInDb](#optionalinshow-vs-optionalindb).

An unmapped kind panics. Do not guess a mapping. Tests build a fake `ssot.Object`. They do not walk `defs.All()` to
assert which live objects are opted in.

The old path for an object that is not opted in stays byte-identical. Do not edit `*_gen.go`.

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
  call. CREATE, ALTER, DROP, parameters, and every other operation stay in the SDK def. They are not generated from
  `Columns`.
- `SkipFields`, `ManualFields`, `TypeOverrides`, and `UsedAsListEntry` have no ssot equivalent yet.
- `ResourceHooks` do not yet generate the whole lifecycle, the SDK calls, or the mapping. A missing step stays an ext
  call, or panics through `Unimplemented`, until that hook is implemented. See [Resource generator](#resource-generator).
- Datasource model builders, provider model builders, resource parameter assertions, and Snowflake object parameter
  assertions do not read `Generators` yet.
- `g.Enum` in the SDK generator is a type alias of the ssot enum type ([enum.go](../../sdk/generator/gen/enum.go)). The
  alias does not mean every enum is defined in ssot yet.

These are the model, not missing coverage. Identity (`Name`, `IdType`, `IdentityNameForceNew`) is not an `Attribute`.
A `Readonly` column must not sit on `Attributes`. A CREATE-only field is an `Attribute` and is absent from `Columns`.
A paired-field option that ssot does not map panics. There is no silent fallback. `StoredAsText` is mapped for a bool
with neither optional flag. Anything else panics. A list of enums is plain `[]string` in the SDK struct, not a slice of
the enum type.

#### Why Show and Describe columns are explicit

`Attributes` is the CREATE/ALTER sequence. `OutputSpec.Columns` is membership and Snowflake order for one operation.
Presence in that slice is the only membership. The same `Attribute` var is pointed at from `Columns` so the snake_case
name is not restated. `ShowName` is the output key when it differs from `Name`. `OutputKey()` is `ShowName` or `Name`.
See [SHOW and DESCRIBE](#show-and-describe).

`Attributes` cannot drive those lists:

- The output has columns that are not settable.
- Some settable fields never appear in the output.
- Some settable fields appear under another key.
- DESCRIBE is not SHOW plus a flag. The two sets differ.
- Order in `Columns` is the order Snowflake returns and the order generation emits. It is not the `Attributes` order.

A generator that needs SHOW reads `Show.Columns`. A generator that needs DESCRIBE reads `Describe.Columns`. It does not
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

SDK paired structs have no second object list. The SDK def calls `ShowOperationFromObject` / `DescribeOperationFromObject`
when the setting is on, and keeps the handwritten Show or Describe call when it is off.

### Parameters

Parameter types live in [parameter.go](parameter.go). Catalog entries live in [defs/parameters.go](defs/parameters.go).

`AccountLevels` and `SessionLevels` are the fixed hierarchies from the [parameter hierarchy](https://docs.snowflake.com/en/sql-reference/parameters#parameter-hierarchy-and-types): account, and account plus user and session. `ObjectLevels` adds one or more object levels on top of the account levels. Provider-only consumer sets (`ACCOUNT_EXT`, `HYBRID_TABLE`, warehouse variants) stay `ParameterLevel` values inside those sets.

### Later

The SDK field chain comes from `Columns` when `SdkShow()` or `SdkDescribe()` is passed to `WithSdk()`. That option also
supplies the Show or Describe method: the query struct, ShowByID filtering, and `DescriptionMappingKind`. A handwritten
field chain, and a handwritten Show or Describe call, are still the source when the matching setting is off. Keep that
chain in sync with `Columns`.
See [PairedStructs](../../sdk/generator/README.md#pairedstructs).
