## ssot

Shared object definitions for generators. Types live in this package; instances in [defs](defs). This package must not
import `pkg/sdk` or any meaningful logic package.

Generators consume ssot; they do not all consume `defs.All()` the same way. Resource generation, SHOW/DESCRIBE schema generation, resource config model builders, resource assertions, resource show output
assertions, Snowflake object assertions, object and resource parameters assertions, and SDK paired structs are opt-in per object ([Generators](#generators)). More
generators will follow. Ultimately, all the generators should default to using ssot and `defs.All()`.

### Object

`Object` is the written intent for one Snowflake object:

- `Attributes` - CREATE/ALTER sequence (resource schema)
- `ComputedAttributes` - read-only resource state used by resource assertions; excluded from CREATE/ALTER and config models
- `Show` / `Describe` - output columns ([SHOW and DESCRIBE](#show-and-describe))
- `Generators` - which generators may consume those lists ([Generators](#generators))
- `ResourceHooks` - resource lifecycle/schema overrides (resource generator)
- `SdkCreateRequestObjectName` - SDK CREATE request type-name stem when it differs from `Name`
- `VariantOf` - base object whose SDK parameter details, SET/UNSET requests, and `ParameterType` names a variant shares; parameter catalog membership still comes from the variant's own `Name`
- `VariantEntityCreateName` - optional SQL object words for a variant CREATE. Empty uses `VariantOf.ObjectType` (`WarehouseAdaptive` emits `WAREHOUSE`). Set it when the words differ, as interactive CREATE does with `INTERACTIVE WAREHOUSE`
- `VariantEntityAlterName` - optional SQL object words for a variant ALTER. Empty uses the same base default as CREATE
- `VariantBaseOnly` - the object only is the base of its `VariantOf` variants and is not a Terraform resource itself (catalog integration: each source type is its own resource). The base owns what the variants share, such as the SDK interface. The resource generator skips it
- `ExactlyOneOf` - groups of optional top-level attributes of which exactly one is set, declared once with pointers to the
  attributes: `ExactlyOneOf: [][]*ssot.Attribute{{&oauth, &bearer, &sigv4}}`. Validation needs at least two attributes per group,
  all of them optional attributes of the object, and none in two groups. SDK CREATE adds `ExactlyOneValueSet` over the Go fields
  of the group to a statement that has all of its attributes in `Create.Clauses`, and panics when it has only some of them.
  The Go field name is the one the clause gets, so `FieldName` on the clause renames it. The generated schema gets Terraform
  `ExactlyOneOf` on each attribute of the group, with the names of the whole group.
- `AtLeastOneOf` - groups of optional top-level attributes of which at least one is set, declared the same way as
  `ExactlyOneOf`. SDK CREATE adds `AtLeastOneValueSet` over the Go fields of a group that is entirely in `Create.Clauses`.
  A custom operation that includes every attribute of a group gets the same check. The schema-side `AtLeastOneOf` is not
  generated from it.
- `Alter.SetConflictingFields` - pairs of `Set` attributes that cannot be set together. The SET struct gets `ConflictingFields`
  for each pair, with the Go field names of the clauses.
- `UnsetConflicts` - groups of `Unset` attributes that cannot be set together. Each group has at least two attributes from
  `Unset`. SDK UNSET adds `ConflictingFields` for the group.
- `AdditionalEnums` - SDK enums that no attribute, column, or clause refers to, such as the accepted values of a field the
  resource validates but the SQL does not send. They are emitted with the enums the object refers to, which
  `SdkDefineFromSsot()` needs so the hand-written SDK definition has nothing left to list
- `SchemaStateUpgraders` - ordered upgrade function names ([Resource generator](#resource-generator))

Identity (`Name`, `IdType`, ...) is separate from `Attributes`. `LegacyIdEncoding` keeps the pipe-delimited resource id
(`db|schema|name`, `helpers.EncodeSnowflakeID`) that older resources store. The generated Create calls it instead of
`helpers.EncodeResourceIdentifier`. The identity schema for those resources stays a plain string and omits
`suppressIdentifierQuoting`. `WithParseId(Generated())` emits a named parser used by Import, Read, and Update.
It calls `sdk.Parse<IdType>` for current ids and `helpers.DecodeSnowflakeIDErrLegacy[<IdType>]` for legacy ids.
The default `ExtCall` keeps the handwritten parser. `PreviewFeature` marks a preview resource. Generated `Delete`
wraps `DeleteContext` with `PreviewFeatureDeleteContextWrapper` and `previewfeatures.<Name>Resource`.

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
  `Name` stays the resource attribute and the CREATE field. `OutputKey()` is `ShowName` or `Name`.
  `FieldName` on a clause renames the SDK request field. It is not how the output column is named.
- `CreateMapping` - overrides the create helper a generated `ApplyCreateOptionals` deduces for the attribute. Leave it
  nil, including on nested fields, unless the deduced call does not match the SDK request. The helpers are in
  `pkg/resources/resource_helpers_create.go`. Without a mapping, an optional field is deduced from its type and the
  schema default: a string, int, bool, enum (`attributeMappedValueCreateBuilder` with `sdk.To<Enum>`), or
  `AccountObjectIdentifier` (by its own `IdentifierType`) uses its builder. A list with `Elem` calls `listValueCreate`,
  and a set, of scalars or of objects, calls `setMappedValueCreateBuilder`. Both take a hand-written per-element mapper
  in package resources named `to<Stem><Field>`, where the stem is the object name plus any enclosing block names:
  `toComputePoolBackupInstanceFamilies`, `toAuthenticationPolicyMfaPolicyAllowedMethods`. A required field with no
  mapping is a constructor argument (string, int, bool, or enum). To override, use a `HelperMapping` constant
  (`StringAttributeCreateBuilder`, ...), `ListValueCreate("<mapperFunc>")`, or
  `AttributeMappedValueCreate("<mapperFunc>")` when the SDK setter takes the whole value. A generated
  `ApplyCreateOptionals` covers only `Attributes`, so a resource field that is not an attribute keeps the
  ext call. The resource identifier is not an attribute. A generated `ParseIdFromConfig` builds it from `IdType`:
  `name` for an account object, `database` and `name` for a database object, and `database`, `schema`, and `name` for
  a schema object.
- `SkipCreateBlock` - leaves a list block out of generated create even when its fields could be mapped. Ext create code
  builds the value, and the generator does not emit a parse function.
- `UpdateMapping` - the update helper from `pkg/resources/resource_helpers_update.go` used by a generated
  `ApplySetUnset`. Use a `SetOnlyMapping` constant (`IntAttributeUpdateSetOnlyBuilder`), a `SetUnsetMapping` constant
  (`StringAttributeUpdateBuilder`, ...), or `ListValueUpdate("<mapperFunc>")`. An attribute with no mapping is left out
  of set/unset.
- `NativeBool` - an optional bool stays Terraform `TypeBool`. Without the mark, an optional bool is a boolean string.
- `IncludeAliasedInTheDescription` - the resource description lists enum aliases. Without the mark, the description
  lists canonical values. Warehouse size sets it.
- `StoredAsText` - the db column is text, and the plain kind stays the semantic type. For a bool, the paired-struct
  generator calls `BoolFromText`, which is true when the column is `Y`. `BoolParsed` with `StoredAsText` parses
  `true`/`false` through `strconv.ParseBool`. A page that says `Y` or `N` can still return `true` or `false`; use
  `BoolParsed` for that column and record the gap. A bool without `StoredAsText` is still db `bool`. Other generators ignore both marks.
  See [PairedStructs](../../sdk/generator/README.md#pairedstructs).
- `DiffSuppress` - a suppressor composed with inferred identifier and enum normalization. Set it directly to
  `ssot.SuppressShow`, `ssot.SuppressDescribe`, or `ssot.SuppressStatement`; combine them with
  `ssot.SuppressIfAny(ssot.SuppressStatement, ssot.SuppressDescribe)`. Nested combinations are flattened, and nil
  means no configured suppression. Output suppressors use `OutputKey()` unless a leaf supplies an override, e.g.
  `ssot.SuppressDescribe.WithKey("definition")`. Output strategies apply to fields, not collection elems.
  Statement suppression uses the resource's `DiffSuppressStatement` helper for strings, including string elems.
- `ConflictsWith` names other attributes on the same level. The generator emits `schema.ConflictsWith` as a path from
  the resource root: the attribute name at the top, and `parent.0.name` for a child of a single block (an `AttrObject`,
  or a list of objects with `MaxItems: 1`).
  A set, or a list without that bound, cannot name these relations. `AtLeastOneOf` uses the same paths and may include
  this attribute.
- `Sensitive` marks the schema field sensitive.
- `SkipInSchema` - the SHOW/DESCRIBE schema generator (`WithSchema`) omits the column from the generated schema map
  and `ToSchema`; the key is dropped entirely (comment left in place). The column still participates in every other
  generator (assertions, SDK). This is the ssot equivalent of the old-list `SkipFields`.

Def files keep attributes as vars (`ComputePoolAttrs`) so `Columns` can point at them instead of restating snake_case
strings.

### Lists and nested objects

How an attribute models a collection or a nested object. The Terraform side (schema, model, resource) and the SDK side
(CREATE and ALTER query structs) read the same attribute, and neither needs a second definition.

- **List or set of scalars.** `Elem` is the scalar element attribute: a string, bool, int, enum, or identifier. A list of
  enums takes a `<Enum>ListItem` query struct in SQL. Create and update use `ListValueCreate` / `ListValueUpdate` with a
  hand-written per-element mapper.
- **Single object.** An attribute of type `AttrObject` with `Children` is one nested block, and one struct in SDK SQL.
  Terraform sees it as a list of at most one object. The resource attribute rule `single_object_is_list_of_one` makes it a
  `TypeList` with `MaxItems: 1`, and `MinItems: 1` when it is `Required`. The definition sets no `MinItems`, `MaxItems`, or
  `Elem`. Its SQL is its own `CreateSql`, `ObjectParameter(Sql("REST_CONFIG ="), Parentheses())`, and each child has the
  `CreateSql` it needs. Generated create parses it like a list block with `MaxItems: 1` (`Attribute.IsSingleBlock()`). A single
  block that is a `SqlOnly` attribute is the same struct and is not a Terraform field. Child `ConflictsWith` and
  `AtLeastOneOf` use the same `name.0.child` paths as a list of objects with `MaxItems: 1`. A list of exactly one object
  (`MaxItems: 1`) is still a list: it is a slice in SQL.
  `ChildrenRule` on the object (`ChildrenAtLeastOne` or `ChildrenExactlyOne`) adds `AtLeastOneValueSet` or `ExactlyOneValueSet` over
  the `Children` that are fields of its SQL struct, so not the fixed texts, in `Children` order. It does not generate the
  schema-side `AtLeastOneOf`.
  A fixed text of its struct, such as `TYPE = OAUTH`, is a `SqlOnly` child with `StaticSQL`, in SQL position. It is not a
  Terraform field: `NestedFields()` skips it and `NestedSqlFields()` has it. Only SDK SQL reads `NestedSqlFields()`.
  The same holds for the `Children` of a list element.
  A `SqlOnly` child of an object element can also be a nested struct (`ObjectParameter(Sql("="), Parentheses())` with its own
  `Children`), for SQL that is nested deeper than the Terraform block. A child of the element that such a struct holds is
  written only inside it, so `client_policy` stays `{client_type, minimum_version}` in Terraform and is
  `CLIENT_TYPE = (MINIMUM_VERSION = '...')` in SQL. `EnumKeyword(NoQuotes(), Required())` writes an enum as a bare keyword
  (the driver name before its settings), and `IdentifierListKeyword()` writes a list of identifiers with no name, in
  parentheses. A `SqlOnly` child can also be a keyword (`BoolKeyword(Sql("('ALL')"), FieldName("All"))` is an optional flag
  with that text). An attribute that a `SqlOnly` nested struct holds counts as listed in the CREATE sequence when the struct
  is: this is how an attribute that is a Terraform field is also the field of a struct that is not (`security_integrations`
  is the list of `SECURITY_INTEGRATIONS = ('ALL')` or `= (a, b)`, and `ChildrenExactlyOne` requires one of the two).
- **List or set of objects.** The list attribute has `Elem` of type `AttrObject`, and the `Children` of that `Elem` are the
  object fields. Terraform sees a nested block (`MaxItems: 1` makes it a single block). `Attribute.NestedFields()` returns the
  fields and `IsObjectList()` tells a list of objects from a list of scalars, so generators ask the attribute for its fields
  and do not read `Children` or `Elem` directly.

  ```go
  {
      Name:     "mfa_policy",
      Type:     ssot.AttrList,
      MaxItems: new(1),
      Elem: ssot.ObjectElem([]ssot.Attribute{
          {Name: "enforce_mfa_on_external_authentication", Type: ssot.AttrEnum, Enum: ...},
          {Name: "allowed_methods", Type: ssot.AttrSet, Elem: &ssot.Attribute{Type: ssot.AttrEnum, Enum: ...}},
      }),
  }
  ```

  `Required`, `MinItems`, `MaxItems`, `Description`, `ForceNew`, and the other settings of the block stay on the list
  attribute. The `Elem` carries the fields and, for SDK SQL, the element clause. `Children` is valid only on an object:
  validation panics when a list or set has `Children`, and an `AttrObject` with no `Children` or with an `Elem`.

  - SQL of the list is the list attribute's `CreateSql`: `ListOfObjects(Sql(text), Parentheses())`.
  - SQL of one element is the `Elem`'s `CreateSql`: `ObjectParameter(...)`. `Parentheses()` puts each element in
    parentheses, through a generated `<Struct>Item` that holds the struct.
  - The struct is named after the object and the list attribute, so `ExternalVolume` and `storage_location` give
    `ExternalVolumeStorageLocation`, and its Request type follows.
  - Each field's SQL is the `CreateSql` the child attribute already has, in `Children` order.
  - Provider-like alternatives are variants of the element: `WithVariants(Variant("S3", &a, &b, ...), ...)` lists, for each
    one, the attributes it consists of, as pointers, in SQL order. Exactly one variant is set. A child that is in no variant
    is a field of the element itself.
- **Object that is not a Terraform attribute.** A `SqlOnly` attribute of type `AttrObject` with `Children` is a nested SQL
  struct (`ENCRYPTION = (...)`, `UPDATE ... CREDENTIALS = (...)`). It is not in the schema, the resource, or the assertions.
  A variant or an ALTER block refers to it by pointer. A fixed text is a `SqlOnly` attribute with `StaticSQL`, and
  `WithCreateSql(&attr, clause)` is a copy of an attribute with another clause for the one place that needs it.
- **Generated create of a block.** A generated `ApplyCreateOptionals` parses an optional list block with a generated
  `parse<Stem>` function. `MaxItems: 1` calls `singleBlockAttributeCreate`, and any other list calls
  `listBlockAttributeCreate`. The block request is built like the object request, with the same deduction and templates:
  required fields are `sdk.New<Stem>Request` arguments in field order, and other fields use their deduced create helper on
  that request. Fields are read from `prefix+"<name>"`. The stem of a block is the parent stem plus the block name:
  `Widget` + `policy` + `limit` is `sdk.WidgetPolicyLimitRequest`. A set block is not parsed, because a set element has no
  index for the prefix. It uses the set rule of `CreateMapping`. A list block stays out of generated create when it is
  required or any field has no create mapping, including a required field the constructor cannot take. Keep such a block
  in an ext `ApplyCreateOptionals`, or use `SkipCreateBlock`. Resource update generation (`ApplySetUnset`) skips nested
  attributes.
- **Where each part is used.** CREATE reads the list attribute through `Create.Clauses`. ALTER reuses the same element with
  `AlterParameterBlock(field, sql, &listAttr)` (`ADD STORAGE_LOCATION`), or a `SqlOnly` object with
  `AlterKeywordBlock(field, sql, &object)`. Enums on the children, on their clauses (`EnumOf`), and in variants are emitted
  with the object's enums. See [Generators](#generators) for the clause options.

[external_volume_def.go](defs/external_volume_def.go) is the reference: a list of objects with four provider variants, nested
`ENCRYPTION` and `CREDENTIALS` structs, static provider texts, and an ALTER `UPDATE` block.

### SHOW and DESCRIBE

`Attributes` cannot drive SHOW/DESCRIBE: more columns than settable fields, CREATE-only fields, renames, mixed Snowflake
order.

`OutputSpec.Columns` is membership and Snowflake order. Presence in `Columns` is the only SHOW/DESCRIBE membership.

Each warehouse variant declares its show in `Show.Columns`. `Warehouse.ShowGeneric` is the temporary union and the one paired SHOW. `ShowPlainName` and `ShowVariants` are plain structs emitted with DESCRIBE. Each one gets a projection on the paired SHOW struct: `As` plus the plain name with the object name removed, so `WarehouseRegular` is `AsRegular` and `WarehouseAdaptive` is `AsAdaptive`. When `ShowByID` is generated, the same suffix also emits `ShowByIDRegular` and `ShowByIDSafelyRegular`. `Show()` still returns the union row. Next, each variant's `Show.Columns` becomes a paired SHOW with its dedicated type, and `ShowGeneric` goes away.

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

- `WithResource()` - resource generator (`make generate-resource`).
  See [Resource generator](#resource-generator). An object in `All()` is a resource-generation input when this call is present.
  Zero options still sets `Enabled`. Generation parts stay on the resource generator's per-object settings.
- `WithSchema(SchemaShow(), SchemaDescribe())` - SHOW/DESCRIBE schema generator (`make generate-show-output-schemas`).
  See [that README](../../schemas/gen/README.md). `SchemaShow(SchemaShowManualFields(&attrs.Tables))` copies
  `ManualFields` from the old show schema def. The argument is the same `*Attribute` that sits in `Show.Columns`. The
  key stays in the generated schema map, and its mapping line is not generated — `*_ext.go` must implement the
  generated `additionalSchemaMapper[T]` stub instead. `SchemaDescribe(SchemaDescribeManualFields(&attrs.Password))`
  is the describe-side copy.
- `WithModel()` - resource config model builders (`make generate-resource-model-builders`).
  See [that section](../../acceptance/bettertestspoc/README.md#adding-new-resource-config-model-builders). Fields come
  from `Attributes` and identity, not from [SHOW and DESCRIBE](#show-and-describe) `Columns`. Types are
  in [attribute.go](attribute.go). An optional bool is stored as a boolean-string; the extra setter takes the ssot type.
  `fully_qualified_name` is not emitted. `WithModel(ModelRequestType(&attrs.StorageLocation, "[]sdk.StorageLocationRequest"))`
  supplies the constructor parameter type for a Required nested-block attribute (a list or set of objects): there is no
  mechanical mapping from a nested block to a concrete Go request type, so the generator cannot build one on its own.
  This is the ssot equivalent of the old-list `complexListAttributesOverrides`, scoped to this object instead of by
  resource name. Without it, a Required nested-block attribute stays a stub (`WithModel()` omitted entirely, as
  `CatalogIntegrationIcebergRest` does) since the constructor parameter would otherwise have no type.
  `WithModel(ModelMultiline(&attrs.CopyStatement))` marks string attributes whose config model builder emits
  `MultilineWrapperVariable`. This is the ssot equivalent of the old-list `multilineAttributesOverrides`, scoped to
  this object. The argument is the same `*Attribute` that sits in `Attributes`, and it must be a string.
- `ComputedAttributes` are included in generated resource assertions, but not in the resource schema or config model.
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
  `ObjectAssertShow(ObjectAssertShowNestedField(&attrs.Locations))` copies `NestedAssertFields` from the old show def.
  The attribute is the one that sits in `Show.Columns`. The column gets `Has<Field>With(...)` taking the nested
  assertion type instead of a plain `Has*`. A nested column has only a placeholder scalar `Elem`, so the assertion type
  cannot be derived from it, and the SDK field's Go type defaults for a list or set to `[]sdk.` + the object `Name` + the column name in camel case with
  one trailing "s" dropped + `Details` (`storage_locations` on `ExternalVolume` is
  `[]sdk.ExternalVolumeStorageLocationDetails`). Pass `ObjectAssertNestedGoType("[]sdk.Other")` when the default does
  not match. `ObjectAssertDescribeNestedField` is the describe-side copy.
  A DESCRIBE assertion uses `sdk.<Object>Details` and fetches through `testClient.<Object>.Describe`.
  Rename the SDK describe type and the test client method to that. There is no type or method override.
  Identified rows provide `ID()` in SDK extensions when it is not generated.
  `ObjectAssertDescribe(WithNoIdentifiableObject())` suppresses the fetch constructor when `Describe` returns a
  slice and the row has no `ID()`.   `ObjectAssertShow(ObjectAssertShowFromBase("ShowByIDAdaptive"))` fetches through
  `testClient.<VariantOf.Name>.ShowByIDAdaptive`. `Object.ShowFromBase` resolves that
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
  from [SHOW and DESCRIBE](#show-and-describe) `Columns`, not from `Attributes`. `SdkDescribePlainOnlyField(name, kind)` adds a field that is only on the DESCRIBE plain struct. `name` is the Go field name and `kind` is its type, as in `PlainOnlyField`. It is not a column. `convert` skips it, and `additionalConvert` in the SDK ext file fills it.
  `SdkShowOptionalOn()` adds the optional `ON` clause, after `IN` and before `STARTS WITH`.
  `SdkShowResultFilterHook()` makes the generated `Show` drop each row for which `excludeFromShow()` returns true, a method
  of the SHOW row that the SDK ext file implements. `SdkShowPlainOnly(&attr, kind)` marks a SHOW column as only on the plain
  struct, with Go type `kind`. `attr` is the same `*Attribute` that sits in `Show.Columns`: the column stays there for the
  schema and the assertions, the SHOW row does not read it, and `additionalConvert` fills the field. A string column with
  `ManualConvert` is skipped by `convert` the same way. `OptionalInDb` is live for the mapped
  db and plain kinds. `StoredAsText` selects `BoolFromText` for a bool. A bool without the mark is still db `bool`.
  Other generators ignore the mark. See [Attribute](#attribute). `SdkDrop()` builds DROP from `Name`: `DROP`, the
  uppercase words, `IF EXISTS`, and the name. `SdkDropDoc()` overrides the default DROP doc URL. `SdkDropSafelyHook()`, `SdkDropSafelyForce()`, and
  `SdkDropCascadeRestrict()` are settings of that call. `SdkDropCascadeRestrict()` adds optional `CASCADE` and
  `RESTRICT` and marks those fields as conflicting. The hook body stays in the SDK ext file. `SdkCreate()` builds
  CREATE from `Name`: `CREATE`, optional `OR REPLACE`, `Create.KeywordsBefore`, the uppercase words, `IF NOT EXISTS`,
  and the name, then `Create.Clauses`. `SdkCreateWithoutOrReplace()` and `SdkCreateWithoutIfNotExists()` leave those
  clauses out. When both clauses are
  present they are marked conflicting. `AtLeastOneOf` adds `AtLeastOneValueSet` for each group that is entirely in
  `Create.Clauses`. `SdkCreateAdditionalValidations()` adds that validation call. `SdkAlter()` builds
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
  the parameter catalog for `ParameterLevel()`. `ShowParametersDetailsAdditionalLevels` adds parameters from more levels
  onto `ShowParametersDetails`. ALTER SET and UNSET stay `ParameterLevel()`. The union is emitted in catalog order, and a name already present is kept once. CREATE does not
  take those levels.
  CREATE and SET use `WithParameters`. UNSET uses `WithParametersUnset`. `AtLeastOneValueSet` lists those parameter
  fields first, then the other fields in the struct.
  An identifier parameter gets `ValidIdentifierIfSet` on CREATE and SET. Comment stays an attribute clause beside that
  entry. Clause shape is `CreateSql`, `SetSql`,
  and `UnsetSql` on [Attribute](attribute.go). A zero create or set clause uses the attribute type. `AsParameter()` on an object clause writes the nested struct as a parameter field (`NAME = (...)`) instead of a list. `PlainParameter()` is a parameter with no quote tag, for an attribute whose type default adds one (an enum gets `no_quotes`). `Clause` applies
  functional options on top of that default. A zero unset clause is a bool keyword. Leave the attribute out of `Unset`
  to omit it, or set `UnsetSql` to a parameter clause to unset by setting a value. A list or set whose element is an enum uses
  an item query struct named after that enum, with `ListItem` appended. Its field is `Value`, a required single-quoted
  keyword.
  `AlterText(sql)` is an optional single-quoted value with no equals, such as `REMOVE STORAGE_LOCATION 'name'`.
  `AlterParameterBlock(field, sql, &attr)` and `AlterKeywordBlock(field, sql, &attr)` are optional nested blocks, as a
  parameter or a keyword clause. `attr` is a list of objects, and the block is one of its elements, with the same struct as in
  CREATE (`ADD STORAGE_LOCATION`), or a SqlOnly object attribute whose `Children` are the fields and whose `CreateSql` are
  the field clauses (`UPDATE`). All three are branches of `ExactlyOneValueSet`.
  A list of objects, its variants, and its nested structs are described on the attributes: see
  [Lists and nested objects](#lists-and-nested-objects). A `Clause(...)` on a child takes `Sql`, `Required`, `FieldName(name)`
  for the Go field name, and `EnumOf(enum)`. An identifier child with a parameter clause is quoted text. `AllowEmpty()` on a
  string parameter makes the SDK field `StringAllowEmpty`, which takes an empty string and quotes itself, so SET can clear it.
  A SET block that changes part of a nested struct is a `SqlOnly` object in `Set` with a `SetSql` object clause and the changed
  fields as `Children`.
  When `SdkCreate()` or `SdkAlter()` is off, that call stays in the def. `SdkCreateForVariants()` emits
  each variant CREATE as a custom method. The variant fills `Create` (`Doc`, clauses, and
  `AdditionalValidations`). The SQL object words are `VariantEntityCreateName`, or `VariantOf.ObjectType` when that override is empty. A fixed value is an `SqlOnly` attribute with `StaticSQL` in that clause list. The method name is `Create` plus the variant name with the base name removed. The query
  struct name is `Create` plus `SdkCreateRequestObjectName`. `SdkAlterForVariants()` emits each variant ALTER the same way.
  The variant fills `Alter` (`Doc`, clauses, and `AdditionalValidations`). SQL words are `VariantEntityAlterName`, or the same base default.
  The method name is `Alter` plus the variant name with the base name removed. The query struct name is `Alter` plus
  `SdkCreateRequestObjectName`. SET and UNSET structs are `Name` plus `Set` or `Unset`.
  TODO [next PRs]: the same switch for variant SHOW/DESCRIBE. `SdkCustomOperations()` emits
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
  come from `ParameterLevel()` plus `ShowParametersDetailsAdditionalLevels`. The generated body calls `ShowParameters`.
  When either setting is off, that call stays in the SDK def. `SdkDefineFromSsot()` defines this object's SDK
  interface from this ssot object alone. The interface contains the operation parts whose flags are on. An object
  without it stays on the handwritten SDK def. An object that already has one does not set this flag. It replaces one
  operation at a time with the matching `FromObject` call. See [Adding a new entity](#adding-a-new-entity).

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
slice and keeps objects with `WithResource()` on. An object without that call is left out. `WithResource()` with no
options still sets `Enabled`.

A subtemplate calls the ext function by default (`ExtCall`). Set that hook to `Generated` to take the step from the
template. Until the subtemplate body is written, generation panics through `Unimplemented`. Each hook is switched on its
own, with its `With` option. See [resource_hooks.go](resource_hooks.go). `DropSafelyInSdk` defaults to omitted, which
keeps `DeleteExt`. Set it to `Generated` to emit `Delete` and point `DeleteContext` at that function. `Delete` calls
`ResourceDeleteContextFunc` with `sdk.Parse` plus the object `IdType`. An object with `LegacyIdEncoding` uses
`helpers.DecodeSnowflakeIDErrLegacy` instead. An object with `PreviewFeature` wraps that `DeleteContext` in
`PreviewFeatureDeleteContextWrapper` using `previewfeatures.<Name>Resource`.

`ShowOutputSet`, `DescribeOutputSet`, and `FqnSet` default to `Generated`. They emit the `d.Set` call for `show_output`,
`describe_output`, and `fully_qualified_name` from the generated `schemas` mappers and the id. `ExtCall` calls
`<name>ShowOutputSetExt(d, <name>)`, `<name>DescribeOutputSetExt(d, <name>Details)`, or `<name>FqnSetExt(d, id)`. Each
returns an `error`. `Omitted` drops the call. A resource does not list these hooks to opt in. `ShowOutputSet` is omitted by default when
`ShowOutputSchema` is omitted. `DescribeOutputSet` is omitted by default when `DescribeInSdk` or `DescribeOutputSchema`
is omitted. `SetDescribeOutput` does not omit it, so that ext only sets what the generated call does not. Warehouse omits `ShowOutputSet` itself because its
`SetConfigFieldsExt` sets it.

`AfterCreate` and `AfterAlter` are `Omitted` by default. They run after `SetId` and after the alters, and their
diagnostics (usually warnings) are added to the final Read. Use them for a step that must not drop the object from
state when it fails, such as resuming a created object. `Import<Name>` is generated only when `ShowByIdInSdk` is
present. A resource with a passthrough importer omits it.

The constructor emits `UpdateContext` only when a field can change in place. An identity field counts when its ForceNew
flag is off. An attribute counts when `ForceNew` is off. An ext identity schema, an attributes-schema extension, an
attributes modification, or parameters attributes count, because those fields are not visible on the model. When nothing
counts, `Update` and `UpdateContext` are left out. The provider rejects an Update function when every field is ForceNew
or computed without Optional. `UpdateContext` is the line directly above `DeleteContext`.

`WithCreateInSdk(Generated())` calls `Client.<NamePlural>.Create(ctx, req)` after applying the create fields and
parameters. A failure returns `diag.FromErr` before setting the resource id or reading the object.
Use `ExtCall` when an object needs a different SDK method or custom error diagnostics.

`WithDescribeInSdk(Generated())` calls `Client.<NamePlural>.Describe(ctx, id)` during Read and assigns the result to
`<prefix>Details`. A failure returns `diag.FromErr` and stops Read. Use `ExtCall` when an object needs a different SDK
method, a transformed result, or custom error handling.

`WithShowByIdInSdk(Generated())` calls `Client.<NamePlural>.ShowByID(ctx, id)` during Import and assigns the result to
`<prefix>`. A failure returns that error. `WithShowByIdSafelyInSdk(Generated())` calls
`Client.<NamePlural>.ShowByIDSafely(ctx, id)` during Read. `ErrObjectNotFound` clears the id and returns a warning that
the resource was removed. Any other error returns `diag.FromErr`. Use `ExtCall` when an object needs a different SDK
method or a different not-found diagnostic.

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
`Name` is a top-level resource or identifier attribute. An object `Elem` marks a nested column or a nested resource
attribute, and both are left out. For now, the set/list fields have to be excluded
(`TODO [SNOW-1648997]`): a column or resource attribute whose type is a set or a list is left out, including a set or
list of scalars. A readonly column that is not a resource attribute is left out: a diff cannot change it, unless its `ShowName` is a
resource attribute (a DESCRIBE-only column that reads that attribute's property in another type, such as a plain string): the key
is then that attribute. Otherwise the key is `Name`, the resource attribute `ComputedIfAnyAttributeChanged` watches.
`fully_qualified_name` is emitted only when at least one identifier attribute is not ForceNew. Its keys are those
attributes, in `IdentityFields` order: `database`, then `schema`, then `name`. `name` is ForceNew when
`IdentityNameForceNew` is set. `database` and `schema` are ForceNew together when `IdentityHierarchyForceNew` is set.
The identity schema template emits that `ForceNew`. An unknown `IdType` panics.
`AdditionalIdentityNameDescription` is appended to the generated name sentence.
`DescribeDetailsName` is the SDK type in `schemas.Describe<Name>Schema`. Empty uses `Name` plus `Details`.
`ObjectType()` is `Sdk.DescribeSql` when that is set, and otherwise the uppercase words of `Name`. The words after
"for the given" are the lowercase of `ObjectType()`. A shared command such as `DESCRIBE CATALOG INTEGRATION` sets
`SdkDescribeSql` so every variant emits the same sentence.

Optional bools are boolean-strings on purpose (`applyResourceAttributeRules`
in [resource_attribute_rules.go](../../resources/gen/resource_attribute_rules.go)): an optional bool is Terraform
`TypeString`, a required bool is `TypeBool`. That classification lives in [classified_field.go](classified_field.go)
and applies to the resource schema only. The schema default is `default`. Create omits the clause for that value. Import writes `true` or `false` from SHOW.

An optional int whose `Min` is 0 or less emits `Default = IntDefault`, because 0 is a real value and cannot mean unset.
Warehouse `auto_suspend` is that case. A minimum above 0 does not emit that default.

A top-level set of identifiers emits `NormalizeAndCompareIdentifiersInSet`. A top-level set of enums emits
`NormalizeAndCompareEnumsInSet`. A list of enums, and a nested set of enums, emit `NormalizeAndCompare` on the
collection. A resource that must keep a nested set case-sensitive clears that suppressor in the attributes
modification function. A list or set element is validated: an enum element uses `sdkValidation`, and an identifier
element uses `IsValidIdentifier` for the element's `IdentifierType`. `ConflictsWith`, `AtLeastOneOf`, and `Sensitive`
are emitted on the schema field. See [Attribute](#attribute).

`WithAlterSetInSdk(Generated())` and `WithAlterUnsetInSdk(Generated())` call `Client.<NamePlural>.Alter` with
`NewAlter<Name>Request(id).WithSet` or `WithUnset`. They compare the corresponding request against
`New<Name>SetRequest()` or `New<Name>UnsetRequest()` with `reflect.DeepEqual`, so empty requests are skipped
and requests containing slices remain supported. SET runs before UNSET; either failure returns `diag.FromErr`
and marks partial state when `HasPartialOnUpdate` is set. Keep `ExtCall` for suspend-aware alters, additional
operations, custom diagnostics, or request names that differ from the object name. Omit UNSET for SET-only objects.

`SetStateToValuesFromConfig` defaults to `Generated`. Read inlines the field names in `setStateToValuesFromConfig`.
The names are optional attributes in `Attributes` order that the helper can copy: string (an optional bool, an enum,
an identifier, or a float), int, and a set of string or int. An optional bool element is stored as a string, so a set
of those bools is included. `comment` is left out, because `SetConfigFields` writes it from SHOW. A nested block
(`Children`) is left out. A list is left out, because the helper does not handle `TypeList`. `CreateOnly` is left out.
Skipping a nested block does not panic. An empty list omits the call. A resource does not list this hook to opt in.
`Omitted` drops the call. `ExtCall` keeps a hand-written `<prefix>SetStateToValueFromConfigFieldsExt` list, for an
object that needs a nested name or a list name.

`WithSetConfigFields(Generated())` emits `<prefix>SetConfigFields(d *tfschema.ResourceData, obj *sdk.<Name>) error`.
Read passes the SHOW row. A variant's row uses the parent type. A field is set when it is a resource attribute in
`Show.Columns`, not in the set-state list, not `CreateOnly`, and not a nested block. Membership matches `Name`.
The SDK field is `obj.` plus the camel-case of the show column's output key (`ShowName`, otherwise `Name`).
A string, int, float, enum, native bool, list, or set is passed through, including a SHOW pointer such as `comment`.
An optional boolean-string uses `booleanStringFromBool`. An identifier uses `FullyQualifiedName()`.
A SHOW pointer is passed through. No matching fields returns nil.
Skipping a nested block does not panic. An object that also writes `show_output` in this hook, or that copies
values from DESCRIBE, keeps `ExtCall`.

#### Other generators

Every plugged generator is opt-in, including the resource generator. It walks `defs.All()`, calls `ObjectGenerators`,
and reads that generator's bools. An object with the flag off stays on the generator's old input. The resource generator
has no old list, so an object with `WithResource()` off is left out. Setting the flag on
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
`NoIdentifiableObject` false and `DescribeSkipFields` nil. `WithModel()`, `WithResourceAssert()`, and `WithResource()` take a variadic of
`ModelOption`, `ResourceAssertOption`, and `ResourceOption`. No constructors exist for those option types. Zero options still sets
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
- Clauses Snowflake statements usually accept are on unless the object opts out: `OR REPLACE`, `IF NOT EXISTS`, ALTER `IF EXISTS`, `SET`, `UNSET`, `SET TAG`, `UNSET TAG`, and the exactly-one alter check. DROP `IF EXISTS` is always emitted. `SET` is a keyword. `SdkAlterSetList` opts into `list,no_parentheses`.
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
- Schema `SkipFields` has an ssot equivalent: the `SkipInSchema` attribute mark ([Attribute](#attribute)). Schema
  `ManualFields` has an ssot equivalent too, but as a generator option rather than an attribute mark:
  `SchemaShowManualFields` / `SchemaDescribeManualFields` ([Generators](#generators)), the same shape as the
  show-output and object-assertion `SkipFields` options below. Schema `TypeOverrides` and `UsedAsListEntry` have no
  ssot equivalent yet. Object-assertion
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
with neither optional flag. Anything else panics. A list of enums is plain `[]string` in a SHOW row, and a slice of
the enum type in a details struct (also in its nested structs).

#### Limitations

- Paired struct field names have no override, for SHOW or for DESCRIBE. The db field and the plain field are both the
  camel-cased `OutputKey()`, so a column cannot map to a differently named plain field (the handwritten defs use
  `WithPlainFieldName` for that, such as `property_value` to `Value`). An object moved to ssot may end up with slightly
  different plain struct fields than its handwritten def had. That is acceptable for now: the SDK is internal, so callers
  are renamed with the move. Add an override to both paths together if a real need appears.
- A list that is not a list of enums or of objects takes the item query struct name as a string. That struct stays in the SDK
  def. A list of objects is generated from its attributes: see [Lists and nested objects](#lists-and-nested-objects).
- A list of enums is built in the SDK generator: struct name `{Enum}ListItem`, field `Value`, required single-quoted keyword. ssot should hold that description, and the generator should only translate it into the query struct. Do that together with more complex list items, not for this enum case alone.
- The generated DESCRIBE parse reads a column in three ways: `DescribeProperty` (a string as it is), `DescribePropertyPrefix` (a nested
  list) and `DescribePropertyParsed` (a hand-written `Parse<Name><Column>Described`). Some columns that need a hand-written parser
  could be generated from the column itself: a list of enums (the items validated with the `To<Enum>` function of the element
  enum), a single enum (one `To<Enum>` call), and a string whose `null` value means empty. They are hand-written for now.
- Computed fields and `Deprecated` are not attribute fields. An explicit schema `Default` is not an attribute field.
  An optional int whose `Min` is 0 or less still emits `Default = IntDefault`. `NativeBool`
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

A generator that needs the paired DESCRIBE row reads `PairedDescribeColumns`: `DescribeRows` when that list is set,
otherwise `Describe.Columns`. `DescribeRows` is for a DESCRIBE that returns raw rows the SDK then parses into the
details `Describe.Columns` models, such as a property list (`parent_property`, `property`, `property_type`, ...). Its
columns are Readonly, and the plain struct is named `<Name>Property` instead of `<Name>Details`, which stays free for
the parsed details. It needs `SdkDescribe()`, and only the SDK paired struct reads it.
The DESCRIBE mapping defaults to a slice when `DescribeRows` is set, since raw property rows are a list. Plain field names
are the camel-cased column names (`property_value` is `PropertyValue`).
`SdkDescribeDetails()` emits the parsed details struct, `<Name>Details`, from `Describe.Columns` as a helper of the
DESCRIBE operation. It needs `DescribeRows`, because that frees the `<Name>Details` name. `Id` is the first field unless
a column is named `id`. `SdkDescribeDetailsNestedField(&attr)` marks a list or set column whose elements are a
hand-written struct, with the element type `[]<Name><Column key in camel case without the plural s>Details`
(`storage_locations` of `ExternalVolume` is `[]ExternalVolumeStorageLocationDetails`). A column that is an object attribute (an attribute that has `Children`, such as `rest_config`) is a nested struct
`<Name><Column key in camel case>Details`, generated from the `Children` of the attribute, as plain values. The field is a
pointer when the attribute is not `Required`. A variant lists its own attributes as DESCRIBE columns, so the nested details are
described by the same attribute as the Terraform block and the SQL. Generators that cannot render a nested column skip it
(`ShowOutputDescribeSkipFields`). Parsing the rows into the struct
stays in the SDK ext file. `SdkDescribeDetailsForVariants()` on the base object (under `SdkDescribe`) emits the details of each
variant that turns on `SdkDescribeDetails` (the variant needs no `DescribeRows`, it uses the ones of the base): the struct
`<Variant>Details` from the Describe columns of the variant and the method `Describe<Variant>Details`, with the generated body of
`SdkDescribeDetailsMethod()` and the hand-written `Parse<Variant>Described`. `DescribePropertyParsed(&attr, "NAME")` reads a column of any type from the property row with exactly that name, through the hand-written `Parse<Name><Column key in camel case>Described(value string) (<column type>, error)`. `SdkDescribeDetailsParseIgnoreUnknown()` makes the generated parse skip a property that has no mapping, with a warning, instead of returning an error: use it when Snowflake can add DESCRIBE properties without a provider release. `SdkDescribeDetailsRowsField(name, kind)` adds a field that keeps the original rows to the end of the details struct, and the generated parse fills it (for a legacy string `describe_output` that typed details cannot rebuild). `SdkDescribeDetailsMethod()` also emits the `DescribeDetails` method: the interface declaration
and a generated body that calls `Describe`, parses the rows with `Parse<Name>Described(rows []<Name>Property)
(<Name>Details, error)`, and sets `Id`. Without `SdkDescribeDetailsParse` that function is hand-written.
`SdkDescribeDetailsParse(...)` generates `Parse<Name>Described`. Every `Describe.Columns` entry except `id` has one
mapping: `DescribeProperty(&attr, "NAME")` reads a required string column from the row with exactly that property
name, and `DescribePropertyPrefix(&attr, "PREFIX_")` reads a nested list column from every row whose property name
contains the prefix. Exact properties are matched first, then prefixes, and any other property is an error
(`Unrecognized <object words> property: <name>`). Two functions stay in the SDK ext file: the element parser of
each prefix mapping, `Parse<Name><Column key in camel case without the plural s>Described(value string) (<element>,
error)`, and `validate<Name>Described(details <Name>Details) error`, which runs last. Only string columns and nested
lists are mapped by the first two, and any other type is mapped by `DescribePropertyParsed`.

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
is on and `NamePlural` is unique among objects that are not variants in [`defs.All()`](defs/all.go). A variant shares the base interface: the variants of
`Warehouse`, and the four catalog integration sources of `CatalogIntegration`, which is a `VariantBaseOnly` base. The base supplies the
interface and each variant supplies its own part of it: `SdkCreateForVariants()` the `Create<Variant>` methods,
`SdkDescribeDetailsForVariants()` the `Describe<Variant>Details` methods and their structs. That option defines the interface from this
ssot object and its variants. It contains the operation parts whose flags are on. A part whose flag is off is absent. An object that already has a handwritten SDK definition does not set this flag, so the merge does not drop that interface. The SDK
generator may then change that interface through an override keyed by the ssot object `Name`. Two objects with the same
`NamePlural` that are not variants of one base do not get an interface. The merge skips an old name that ssot emits.
`ShowOperationFromObject` and the other `FromObject` calls stay the translators for one operation. They panic when that
operation's flag is off.

### Parameters

Parameter types live in [parameter.go](parameter.go). Catalog entries live in [defs/parameters.go](defs/parameters.go).

`AccountLevels` and `SessionLevels` are the fixed hierarchies from the [parameter hierarchy](https://docs.snowflake.com/en/sql-reference/parameters#parameter-hierarchy-and-types): account, and account plus user and session. `ObjectLevels` adds one or more object levels on top of the account levels. Provider-only consumer sets (`ACCOUNT_EXT`, `HYBRID_TABLE`, warehouse variants) stay `ParameterLevel` values inside those sets.

### Adding a new entity

A new Snowflake object is one definition in [defs](defs), registered in [`All()`](defs/all.go). Generator flags stay off until a step below turns one on. `WithResource()` stays off until the resource step, so membership in `All()` does not emit a resource. On the definition step, read `pkg/sdk/generator/defs/<plural_snake_case>_def.go` and `pkg/sdk/<plural>_*.go` when they exist. The SSOT definition matches that object. Where the handwritten definition and a public page disagree, follow the page and record the gap. A new object has no handwritten SDK definition. `SdkDefineFromSsot()` then builds the interface from the flags that are on, and a flag that is off is absent. Do not add a handwritten file for that object. An object that already has a handwritten SDK definition does not set `SdkDefineFromSsot()`. That flag replaces the whole interface. Migrate one operation at a time by replacing that operation on the handwritten definition with the matching `FromObject` call: `CreateOperationFromObject`, `ShowOperationFromObject`, `DropOperationFromObject`, `AlterOperationFromObject`, or `DescribeOperationFromObject`. The other operations stay as they are. Do not add `CREATE OR ALTER` when the existing interface does not already have it. A create page that documents the variant is not enough. Record that gap. After generation, diff `pkg/sdk/<plural>_*.go` with the previous files. The operation just migrated may differ by an expected rename or a structural change. Every other operation is unchanged. Set `SdkDefineFromSsot()` and delete the handwritten definition only after every operation has moved.

Public documentation only. Implement a command, clause, column, or parameter when a public Snowflake page for this object shows it: the SQL reference, or the user guide when no SQL reference page exists. A missing page is not filled from another object, from the parser, or from a live query. Record what was left out, which page was checked, and where the capability was seen. Parameters and tags stay off unless that page shows them. When the default SQL reference URL does not exist, override it with `SdkCreateDoc`, `SdkShowDoc`, `SdkDropDoc`, `SdkAlterDoc`, or `SdkDescribeDoc`.

When a public page shows a command and no column list, the columns come from an uncommitted integration test. Give it a build tag so `go test ./...` does not run it. The test uses the default Snowflake profile, creates only what that command needs, runs the command, logs column names, database types, and every row, and drops what it created. It is not kept in the repository. After that test has been run, put those columns on the matching `Columns` list in the logged order and enable the generator setting for that command. Record that the list came from the probe. A string column that came back null is optional in SHOW and in the DB row. A time column that is nullable in the driver is `OptionalInDb` only: an optional plain time is not mapped, and the plain field stays `time.Time`. Do not mark `name`, `database_name`, or `schema_name` optional when the SHOW row's `ID()` is built from them.

When a public page shows ALTER and does not list SET, UNSET, or RENAME, those clauses come from the same kind of uncommitted test. It runs each clause the way other objects spell it, logs the statement and the error or the SHOW row after it, and drops what it created. After that test has been run, add a clause that succeeded and leave out a clause Snowflake rejected. A field that SET and UNSET can change is not ForceNew.

Each step enables one generator flag, or one operation setting inside `WithSdk()`. SHOW and DESCRIBE are separate settings of schema generation, object assertions, and show-output assertions. Regenerate with `--filter-object-names` after reading that generator's `--help`. `make generate-sdk` filters on `NamePlural`. The resource generator, model builders, resource assertions, and parameter assertions filter on `Name`. Snowflake object assertions, show-output assertions, and show-output schemas filter on `sdk.<Name>` for SHOW and `sdk.<Name>Details` for DESCRIBE. Do not edit `*_gen.go`. Extension files hold hook bodies and other custom code. A `-check` target restores uncommitted generated diffs, so run it only after those files are committed.

1. **Definition.** Add `defs/<plural_snake_case>_def.go` and register the object in `All()`. Read an existing handwritten SDK definition first, as the paragraph above describes. Registration runs `ValidateObject`, so the definition is valid with every generator still off. Attributes, [`Create`](object.go) / [`Alter`](object.go) clauses, and `Show` / `Describe` `Columns` may be filled before the generator that reads them is enabled. A `Readonly` SHOW or DESCRIBE column stays off `Attributes`. A writable output column is that same attribute. An explicit SQL clause is listed on that statement. No `Generators` entry yet.
2. **SDK.** One operation at a time: `SdkCreate()`, `SdkShow()`, `SdkDrop()`, `SdkAlter()` (`make generate-sdk`). A new object also sets `SdkDefineFromSsot()`. An existing handwritten SDK definition does not. It replaces that one operation with the matching `FromObject` call, as the paragraph above describes, and diffs the generated files. Show needs `Columns` first. Enabling show or describe without `Columns` panics. The generator does not add the interface to `Client`; that field and its constructor stay in `pkg/sdk/client.go`. Show emits `ObjectType()` as `ObjectType<Name>`. That constant, its plural, and the singular-to-plural map stay in `pkg/sdk/object_types.go`. `SdkShowOptionalIn` emits `IN` for account, database, and schema, and `SdkShowByIDIn` is the matching ShowByID filter. `SdkShowOptionalExtendedIn` also emits application and application package. When the public page omits `IF EXISTS`, probe the clause before leaving it out. `SdkDrop()` still emits `IF EXISTS`. Add an opt-out only after Snowflake rejects the clause. `SdkDropSafelyHook()` is for a pre-drop step the public page shows. Leave it off when the page shows none. `DropSafely` then drops directly. Write the SDK ext for hook bodies and the `*_ext_test.go` companion for the operation just enabled.
3. **DESCRIBE.** When the published page lists columns, enable `SdkDescribe()` from that list. When the page has the command and no columns, use the uncommitted probe above, then enable `SdkDescribe()` from that output. `SdkDescribeDoc()` points at the page that shows the command. The default result is one value. A property-list result uses `DescribeRows` and `SdkDescribeMapping(SdkDescriptionMappingSlice)`. When DESCRIBE omits a column that `ID()` needs, add that field with `SdkDescribePlainOnlyField(name, kind)`. It is not a `Columns` entry. `convert` skips it. `additionalConvert` leaves it unset, because the row has no value for it. A method such as `DescribeDetails` copies it from the identifier passed to `Describe`, and `ID()` is built from those fields. When the row already has every column `ID()` needs, write `ID()` in the SDK ext from those fields. Do not add plain-only fields for them. The test client `Describe` calls the SDK `Describe`.
4. **Object assertions.** `WithObjectAssert(ObjectAssertShow())` (`make generate-snowflake-object-assertions`). `ObjectAssertDescribe()` is its own step, and only when describe columns exist. The generated show fetch calls `testClient.<Name>.Show`, which is `ShowByID`. The describe fetch calls `testClient.<Name>.Describe`. Each method has to exist on the test client or the assertion package does not compile. The details value needs `ID()`, from the DESCRIBE step: the row columns when they are present, or the plain-only fields when a column is missing. `WithNoIdentifiableObject()` is for a describe value that still has no `ID()`, such as a property-list row. Creating and dropping the object stays with the integration tests. `WithObjectParametersAssert()` (`make generate-snowflake-object-parameters-assertions`) only when the [catalog](defs/parameters.go) has entries for `ParameterLevel()`.
5. **SDK integration tests.** In `pkg/sdk/testint/`, with build tag `non_account_level_tests` unless the test changes account settings. Use the object assertions from the previous step. Cover each operation the SDK has: create, alter, drop of an existing object, drop of a missing object, `IF EXISTS` on a missing object when the statement has it, `DropSafely` when the SDK has it, each SHOW clause the statement has, `ShowByID` for the same name in two schemas, `ShowByIDSafely`, and describe of an existing object and of a missing one. Leave out a clause the statement does not have. An optional text column that comes back as `""` is asserted as that empty string. `HasNoX` means the pointer is nil. A stage path in a string clause is `stage.Location()`, which quotes the fully qualified name. An unquoted `@db.schema.stage` is folded to uppercase and misses a mixed-case test database or schema.
6. **SHOW and DESCRIBE schemas.** `WithSchema(SchemaShow())` (`make generate-show-output-schemas`). `SchemaDescribe()` is its own step, and only when describe columns exist. The schema and mapper are built from `Columns`. A `SdkDescribePlainOnlyField` stays on the SDK details struct and is left out of the DESCRIBE schema and mapper. SHOW writes `schemas.Show<Name>Schema` and `<Name>ToSchema` in `<snake>_gen.go`. DESCRIBE writes `schemas.Describe<Name>DetailsSchema` and `<Name>DetailsToSchema` in `<snake>_desc_gen.go`. The resource step needs the show schema: the generated show-output hook references `schemas.Show<Name>Schema`.
7. **Resource.** `WithResource()` (`make generate-resource`) and the hooks that step generates. Generate the describe SDK call, describe output schema, and describe output set when describe columns exist. Omit those hooks when there is no describe. Omit parameter hooks and the parameters generation part when the object has no parameters. Omit set and unset when the statement has no SET or UNSET. When every field is ForceNew, the constructor omits `Update`, so there is no update id to keep used. `WithAfterAlter(ExtCall())` with a body that returns nil is for an `Update` that is still generated and would otherwise leave the parsed id unused. A config field that SHOW does not return stays in the set-state-from-config list so Read keeps the configured value. `IdentityNameForceNew` and `IdentityHierarchyForceNew` are both emitted on the identity schema. `WithIdentitySchema(ExtCall())` is for an identity schema that differs from that. See [Resource generator](#resource-generator). Do not register the resource in this step. An entry in `acceptanceTestsProvider()` is visible to every acceptance test. Do not add a preview feature constant. Add the constant, list it in `allPreviewFeatures`, and switch create, read, and update to the preview wrappers when the resource moves to `pkg/provider/provider.go`.
8. **Resource helpers.** One flag per step: `WithModel()` (`make generate-resource-model-builders`), `WithResourceAssert()` (`make generate-resource-assertions`), `WithShowOutputAssert(ShowOutputShow())` (`make generate-resource-show-output-assertions`). `ShowOutputDescribe()` is its own step when describe columns exist. `WithResourceParametersAssert()` (`make generate-resource-parameters-assertions`) only when the catalog has entries for the object.
9. **Acceptance test.** One `TestAcc_<Object>_BasicUseCase` in `pkg/testacc/`, with build tag `non_account_level_tests` unless the test changes account settings, plus a helper client. Use the generated model and assertions, including an object assertion. Register the resource in `acceptanceTestsProvider()` in `pkg/testacc/1_setup_provider_test.go` in this step, not earlier. Do not add a preview feature constant or preview wrappers on that provider. Those wrappers call `StringToFeature`, which accepts only `ValidPreviewFeatures`, the same allow-list as `preview_features_enabled`. `AllPreviewFeatures` is what `docs/index.md` writes into that field. Create, read, and update stay tracking wrappers. Generated delete calls `DropSafely` through the tracking wrapper. Register destroy checking in `pkg/testacc/2_check_destroy_test.go` through `ShowByID`. A config field that SHOW does not return is listed in `ImportStateVerifyIgnore`. So is an optional SHOW string that comes back nil: Read leaves it absent when the config omits it, and import writes that nil as `""`. A nested show or describe string that the mapper leaves unset is stored as `""`, so that output assertion expects the empty string. A resource string that `SetConfigFields` copies from a non-nil SHOW pointer is present, including a pointer to `""`. The SDK object assertion still treats a nil value as unset. The test does not enable a preview feature.

A data source, moving registration from `acceptanceTestsProvider()` to `pkg/provider/provider.go`, docs templates, and examples come after the acceptance test.

### Later

The SDK field chain comes from `Columns` when `SdkShow()` or `SdkDescribe()` is passed to `WithSdk()`. That option also
supplies the Show or Describe method: the query struct, ShowByID filtering, and `DescriptionMappingKind`. A handwritten
field chain, and a handwritten Show or Describe call, are still the source when the matching setting is off. Keep that
chain in sync with `Columns`.
See [PairedStructs](../../sdk/generator/README.md#pairedstructs).
