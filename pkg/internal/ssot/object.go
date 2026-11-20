package ssot

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

// StateUpgraderFunc names a function in package resources with the schema.StateUpgraderFunc signature.
// Object.SchemaStateUpgraders lists these functions in version order.
// The version each function upgrades from is its index.
// The resource schema version is the length of the slice.
type StateUpgraderFunc string

// CustomizeDiffFunc names a schema.CustomizeDiffFunc in package resources.
// Object.AdditionalCustomizeDiffs lists these functions in call order.
// They are appended to the generated customdiff.All.
type CustomizeDiffFunc string

type Object struct {
	Name                                string
	NamePlural                          string
	DescribeFailRead                    bool
	DescribeUsesShowId                  bool
	HasPartialOnUpdate                  bool
	AdditionalDescribeOutputDescription string
	// DescribeDetailsName is the SDK type in schemas.Describe<Name>Schema.
	// Empty uses Name plus Details.
	DescribeDetailsName string
	IdType              string
	// LegacyIdEncoding keeps the pipe-delimited resource id (db|schema|name) that older resources store.
	LegacyIdEncoding bool
	// PreviewFeature marks a preview resource. Generated Delete wraps DeleteContext with previewfeatures.<Name>Resource.
	PreviewFeature       bool
	IdentityNameForceNew bool
	// IdentityHierarchyForceNew marks database and schema ForceNew in the identity schema.
	IdentityHierarchyForceNew bool
	// AdditionalIdentityNameDescription is appended to the generated name sentence.
	AdditionalIdentityNameDescription string
	// IdentityNamesCreatedBySnowflake are names of objects that Snowflake creates itself, such as the PUBLIC schema of a
	// database. With WithBeforeCreate(Generated()), creating one of them adopts the existing object instead: the generated
	// code sets the id and runs Update. Names are compared ignoring case and outer spaces.
	IdentityNamesCreatedBySnowflake []string
	// SdkCreateRequestObjectName is the CREATE request stem when it differs from Name.
	// Empty uses Name. A variant sets Variant.RequestObjectName instead. PostgresMirrorArguments is one.
	SdkCreateRequestObjectName string
	// Variant groups the settings of one role of a base object. Nil means this object is not a variant.
	Variant *ObjectVariant
	// BaseVariant groups the settings of an object that is the base of variants.
	// Nil means this object is not such a base.
	BaseVariant            *BaseVariant
	Attributes             []Attribute
	ComputedAttributes     []Attribute // Read-only resource state; never used for CREATE/ALTER.
	AttributesModification bool
	ResourceHooks          []ResourceHookOption
	// ResourceParameters is the subset of the object's parameter catalog that is a resource attribute.
	// Empty uses the full ParameterLevel() catalog. A non-empty list is still emitted in catalog order.
	ResourceParameters       []ParameterDef
	DataSourceHooks          []DataSourceHookOption // Exceptions to the defaults of EffectiveDataSourceHooks.
	Generators               []GeneratorOption
	SchemaStateUpgraders     []StateUpgraderFunc
	AdditionalCustomizeDiffs []CustomizeDiffFunc
	// SuspendForAlter says the object must be suspended while it is altered. SdkSuspendIfStarted generates the SDK method
	// that does it, and WithSuspendForAlter(Generated()) calls it once per update.
	SuspendForAlter *SuspendForAlter
	Show            OutputSpec
	Describe        OutputSpec
	// DescribeRows is the DESCRIBE row when the SDK Describe method returns raw rows that differ from Describe.Columns.
	// A property-list DESCRIBE is one: the rows are (property, property_type, ...) and the SDK parses them into the details
	// that Describe.Columns models. Only the SDK paired struct reads DescribeRows. Schema and assertions read Describe.Columns.
	DescribeRows OutputSpec
	// Create and Alter are the SQL field order for those statements.
	// Presence in a slice is membership. The slice order is the order generation emits.
	// Alter is the order of the blocks and action keywords after the name. The SET and UNSET blocks carry their fields
	// as groups, see SetAttributes and UnsetAttributes.
	// A field in those slices may be absent from Attributes when it is not a resource attribute.
	Create CreateSpec
	Alter  AlterSpec
	// AdditionalEnums are SDK enums that the generator does not find from the object's attributes and clauses, such as enums
	// of nested ALTER blocks and constraints, or the accepted values of a field that the resource validates but the SQL does not
	// send. They are emitted with the enums the object refers to.
	AdditionalEnums []*Enum
	// CustomOperations are operations other than the object's CREATE and ALTER, including extra CREATE shapes.
	// Clause order is the emission order. SdkCustomOperations turns them on.
	CustomOperations []CustomOperationSpec
	// CustomInterfaceMethods are interface declarations whose bodies stay in the SDK ext file.
	// SdkCustomInterfaceMethods turns them on.
	CustomInterfaceMethods []CustomInterfaceMethodSpec
	// ExactlyOneOf are groups of optional top-level attributes of which exactly one is set. Each group is declared once,
	// with pointers to the attributes. SDK CREATE adds ExactlyOneValueSet to a statement that has all of the group,
	// and the generated schema gets ExactlyOneOf on each attribute of the group.
	ExactlyOneOf [][]*Attribute
	// AtLeastOneOf are groups of optional top-level attributes of which at least one is set. Each group is declared once,
	// with pointers to the attributes. SDK CREATE adds AtLeastOneValueSet to a statement that has all of the group.
	// The schema-side AtLeastOneOf is not generated from it.
	AtLeastOneOf [][]*Attribute
	// ShowParametersDetailsAdditionalLevels adds parameters from these levels onto ShowParametersDetails.
	// A name already in ParameterLevel() is kept once.
	ShowParametersDetailsAdditionalLevels []ParameterLevel
}

// ObjectVariant is one role of a base object.
// Name starts with the base Name, and NamePlural matches the base.
// The variant shares the base SDK interface, parameter details, SET and UNSET requests, and ParameterType names.
// Catalog membership still uses the variant Name.
type ObjectVariant struct {
	// Of is the base.
	Of *Object
	// CreateName overrides the SQL object words CREATE emits.
	// Empty uses Of.ObjectType, so WarehouseAdaptive emits WAREHOUSE. Interactive CREATE sets INTERACTIVE WAREHOUSE.
	CreateName string
	// AlterName overrides the SQL object words ALTER emits.
	// Empty uses Of.ObjectType, the same default as CREATE.
	AlterName string
	// RequestObjectName is the CREATE and ALTER query-struct stem when it differs from Name.
	// Empty uses Name. AdaptiveWarehouse yields CreateAdaptiveWarehouse.
	RequestObjectName string
	// SchemaShowFromBase reuses the base SHOW schema. SchemaShow stays on.
	// No Show<Variant>Schema is generated. The variant resource uses Show<Base>Schema.
	SchemaShowFromBase bool
	// ObjectAssertShowRowFromBase reuses the base SHOW object assertion. ObjectAssertShow stays on.
	// No <Variant> assertion is generated. Tests use the base assertion.
	ObjectAssertShowRowFromBase bool
	// ShowOutputShowRowFromBase reuses the base SHOW show-output assertion. ShowOutputShow stays on.
	// Only the constructors <Variant>ShowOutput and Imported<Variant>ShowOutput are generated, and they return the base type.
	ShowOutputShowRowFromBase bool
	// ObjectAssertShowFromBase is the SHOW fetch method on the base client, such as ShowByIDAdaptive.
	// Empty keeps testClient.<Object>.Show. The row type stays sdk.<Variant>.
	// The object type constant is sdk.ObjectType<Base> when ObjectAssertShowObjectTypeName is empty.
	ObjectAssertShowFromBase string
}

// BaseVariant groups the settings of an object that is the base of variants.
type BaseVariant struct {
	// NotResource marks a base that is not a Terraform resource itself: the variants are separate resources
	// and the base owns what they share, such as the SDK interface. The resource generators skip it.
	NotResource bool
	// ShowGeneric is the temporary union of every variant SHOW column.
	// It is the one paired SHOW. Each variant, including the main one, declares its own show in Show.
	// TODO [next PRs]: Each variant's Show.Columns becomes a paired SHOW with its dedicated type, and ShowGeneric goes away.
	ShowGeneric OutputSpec
	// ShowPlainName is the plain SHOW struct name when it differs from Name.
	// Empty uses Name. Set it when ShowGeneric uses Name for the paired SHOW struct.
	// The projection on the paired SHOW struct is As plus this name with Name removed, so WarehouseRegular is AsRegular.
	// When ShowByID is generated, that suffix also emits ShowByIDRegular and ShowByIDSafelyRegular.
	ShowPlainName string
	// ShowVariants are plain SHOW structs emitted with this object's DESCRIBE.
	// Each variant's Show.Columns is that struct. The variant does not set SdkShow.
	// The projection is As plus the variant name with Name removed, so WarehouseAdaptive is AsAdaptive.
	// When ShowByID is generated, that suffix also emits ShowByIDAdaptive and ShowByIDSafelyAdaptive.
	// TODO [next PRs]: Emit these structs from ShowOperationFromObject. Show() still returns the union row.
	ShowVariants []*Object
}

// VariantBase is the base object, or nil when this object is not a variant.
func (obj Object) VariantBase() *Object {
	if obj.Variant == nil {
		return nil
	}
	return obj.Variant.Of
}

// IsVariantBaseOnly says the object only is the base of its variants and is not a Terraform resource itself.
func (obj Object) IsVariantBaseOnly() bool {
	return obj.BaseVariant != nil && obj.BaseVariant.NotResource
}

// ShowGenericColumns is the union SHOW of a variant base, or nil when this object is not one.
func (obj Object) ShowGenericColumns() []*Attribute {
	if obj.BaseVariant == nil {
		return nil
	}
	return obj.BaseVariant.ShowGeneric.Columns
}

// ShowPlainName is the plain SHOW struct name of a variant base. Empty uses Name.
func (obj Object) ShowPlainName() string {
	if obj.BaseVariant == nil {
		return ""
	}
	return obj.BaseVariant.ShowPlainName
}

// ShowVariants are the plain SHOW structs a variant base emits with DESCRIBE.
func (obj Object) ShowVariants() []*Object {
	if obj.BaseVariant == nil {
		return nil
	}
	return obj.BaseVariant.ShowVariants
}

// RequestObjectName is the CREATE request stem, and a variant's CREATE and ALTER query-struct stem.
// A variant uses Variant.RequestObjectName. Any other object uses SdkCreateRequestObjectName. Empty uses Name.
func (obj Object) RequestObjectName() string {
	if obj.Variant != nil && obj.Variant.RequestObjectName != "" {
		return obj.Variant.RequestObjectName
	}
	return obj.SdkCreateRequestObjectName
}

// ParameterLevelsFor is the catalog levels a Parameters entry spreads.
// CREATE, SET, and UNSET are ParameterLevel().
// ShowParametersDetails appends ShowParametersDetailsAdditionalLevels itself.
func (obj Object) ParameterLevelsFor(string) []ParameterLevel {
	return []ParameterLevel{obj.ParameterLevel()}
}

// AlterClauseKind is one entry in an ALTER statement after the name.
type AlterClauseKind string

const (
	AlterClauseKeyword     AlterClauseKind = "keyword"
	AlterClauseRenameTo    AlterClauseKind = "rename_to"
	AlterClauseSwapWith    AlterClauseKind = "swap_with"
	AlterClauseSet         AlterClauseKind = "set"
	AlterClauseUnset       AlterClauseKind = "unset"
	AlterClauseSetTags     AlterClauseKind = "set_tags"
	AlterClauseUnsetTags   AlterClauseKind = "unset_tags"
	AlterClauseList        AlterClauseKind = "list"
	AlterClauseKeywordList AlterClauseKind = "keyword_list"
	AlterClauseText        AlterClauseKind = "text"
	AlterClausePlainText   AlterClauseKind = "plain_text"
	AlterClauseBlock       AlterClauseKind = "block"
)

// AlterBlockJoin is how a nested block is tagged on the ALTER options: the payload is Block, this is the join.
type AlterBlockJoin string

const (
	// AlterBlockJoinParameter is a parameter clause, such as ADD STORAGE_LOCATION.
	AlterBlockJoinParameter AlterBlockJoin = ""
	// AlterBlockJoinKeyword space-joins the nested struct as a keyword field. SQL may be empty: the nested struct then
	// carries the words, and the field tag is a bare keyword. Clustering and add/drop row access policy use this.
	AlterBlockJoinKeyword AlterBlockJoin = "keyword"
	// AlterBlockJoinComma comma-joins the nested struct's children with no parentheses (ddl list,no_parentheses),
	// such as DROP ROW ACCESS POLICY …, ADD ROW ACCESS POLICY …. UNSET uses the same join.
	AlterBlockJoinComma AlterBlockJoin = "comma"
)

// AlterClause is one ALTER entry. Keyword and list use SQL. The other kinds are the SET, UNSET, SET TAG, and UNSET TAG blocks.
type AlterClause struct {
	Kind AlterClauseKind
	SQL  string
	// NotExactlyOne leaves a keyword out of ExactlyOneValueSet. It may accompany another action.
	NotExactlyOne bool
	// ListItem is the Go type of a list element, without a pointer.
	ListItem string
	// Field is the Go field name of a block.
	Field string
	// Block is the nested struct: a list of objects, whose element is the block, or a SqlOnly object attribute whose
	// Children are the block fields. It is not the ALTER field join; that is BlockJoin.
	Block *Attribute
	// BlockJoin is how Block is tagged on the ALTER options. See AlterBlockJoin.
	BlockJoin AlterBlockJoin
	// Attribute and AttributeValue link a keyword to the bool attribute value that selects it, such as RESUME for an
	// enabled attribute that is true. The SDK ignores them. See AlterKeywordFor.
	Attribute      *Attribute
	AttributeValue bool
	// Groups are the statements of a SET or UNSET block. See AlterGroup.
	Groups []AlterGroupSpec
	// MixedForms records an AlterSet or AlterUnset call that mixed attributes and groups. Validation rejects it.
	MixedForms bool
	// SharedToOpts generates a toOpts method for the block struct, because a struct of another SDK interface calls it.
	// WithSharedToOpts sets it.
	SharedToOpts bool
}

// AlterKeyword is an optional action keyword, such as RESUME or ENABLE MANAGED ACCESS.
// It is one branch of ExactlyOneValueSet.
func AlterKeyword(sql string) AlterClause {
	return AlterClause{Kind: AlterClauseKeyword, SQL: sql}
}

// AlterKeywordFor is AlterKeyword selected by a value of a bool attribute, such as RESUME when enabled is true and
// SUSPEND when it is false. The resource generator reads the link; the SDK sees a plain keyword.
func AlterKeywordFor(attribute *Attribute, value bool, sql string) AlterClause {
	return AlterClause{Kind: AlterClauseKeyword, SQL: sql, Attribute: attribute, AttributeValue: value}
}

// AlterKeywordNotExactlyOne is an optional action keyword that may accompany another branch.
// IF SUSPENDED is one: it can be set together with RESUME.
func AlterKeywordNotExactlyOne(sql string) AlterClause {
	return AlterClause{Kind: AlterClauseKeyword, SQL: sql, NotExactlyOne: true}
}

// AlterList is a list assignment after the name, such as ADD TABLES.
// itemType is the element Go type, without a pointer. The clause uses parentheses and no equals.
// It is one branch of ExactlyOneValueSet, and each element is a valid identifier.
func AlterList(sql, itemType string) AlterClause {
	return AlterClause{Kind: AlterClauseList, SQL: sql, ListItem: itemType}
}

// AlterText is an optional single-quoted text value after the name, with no equals, such as REMOVE STORAGE_LOCATION 'name'.
// It is one branch of ExactlyOneValueSet.
func AlterText(sql string) AlterClause {
	return AlterClause{Kind: AlterClauseText, SQL: sql}
}

// AlterKeywordList is an optional keyword list in parentheses with no comma, such as MODIFY CONDITION EXISTS.
// field is the Go field name. itemType is the element Go type, without a pointer.
// The elements are not checked as identifiers. It is one branch of ExactlyOneValueSet.
func AlterKeywordList(field, sql, itemType string) AlterClause {
	return AlterClause{Kind: AlterClauseKeywordList, Field: field, SQL: sql, ListItem: itemType}
}

// AlterPlainText is an optional unquoted value with no equals, such as MODIFY ACTION.
// field is the Go field name. It is one branch of ExactlyOneValueSet.
func AlterPlainText(field, sql string) AlterClause {
	return AlterClause{Kind: AlterClausePlainText, Field: field, SQL: sql}
}

// AlterParameterBlock is an optional nested block after the name, as a parameter clause, such as ADD STORAGE_LOCATION.
// field is the Go field name. block is the attribute the block is built from: a list of objects (the block is one element,
// as in the CREATE list), or a SqlOnly object attribute whose Children are the fields.
// It is one branch of ExactlyOneValueSet.
func AlterParameterBlock(field, sql string, block *Attribute) AlterClause {
	return AlterClause{Kind: AlterClauseBlock, Field: field, SQL: sql, Block: block}
}

// AlterKeywordBlock is AlterParameterBlock as a keyword clause (space-join). UPDATE followed by the struct.
// An empty sql leaves the SQL prefix off. The nested struct carries the words.
func AlterKeywordBlock(field, sql string, block *Attribute) AlterClause {
	return AlterClause{Kind: AlterClauseBlock, Field: field, SQL: sql, Block: block, BlockJoin: AlterBlockJoinKeyword}
}

// AlterCommaSeparatedBlock is a nested block whose children are comma-joined with no parentheses and no SQL prefix,
// such as DROP ROW ACCESS POLICY …, ADD ROW ACCESS POLICY ….
func AlterCommaSeparatedBlock(field string, block *Attribute) AlterClause {
	return AlterClause{Kind: AlterClauseBlock, Field: field, Block: block, BlockJoin: AlterBlockJoinComma}
}

// BlockProjection is the statement of Block this clause embeds. A field named like the FieldName of the block's
// SetSql or UnsetSql is that projection, such as SetAggregationPolicy. Any other field is the CREATE struct.
func (c AlterClause) BlockProjection() ClauseUse {
	switch {
	case c.Block == nil:
		return ClauseCreate
	case c.Block.SetSql.Present && c.Block.SetSql.Field == c.Field:
		return ClauseSet
	case c.Block.UnsetSql.Present && c.Block.UnsetSql.Field == c.Field:
		return ClauseUnset
	default:
		return ClauseCreate
	}
}

// WithSharedToOpts makes the block struct generate a toOpts method, for a block that another SDK interface shares.
func (c AlterClause) WithSharedToOpts() AlterClause {
	c.SharedToOpts = true
	return c
}

// AlterRenameTo is RENAME TO. The identifier type is the object's IdType.
func AlterRenameTo() AlterClause { return AlterClause{Kind: AlterClauseRenameTo} }

// AlterSwapWith is SWAP WITH. The identifier type is the object's IdType.
func AlterSwapWith() AlterClause { return AlterClause{Kind: AlterClauseSwapWith} }

// MainAlterGroup is the name of the main group of a SET or UNSET block. Its request is set (unset) in the resource ext
// struct, and the generated update also sends the parameters and every attribute that no other group holds with it.
const MainAlterGroup = ""

// AlterItem is what AlterSet and AlterUnset accept: an attribute, or a group made by AlterGroup.
type AlterItem interface{ alterItem() }

func (*Attribute) alterItem()     {}
func (AlterGroupSpec) alterItem() {}

// AlterSet is the SET block. Pass its attributes when one ALTER statement sends all of them. When Snowflake needs
// several statements, enclose every attribute in an AlterGroup instead. The two forms cannot be mixed in one call.
func AlterSet(items ...AlterItem) AlterClause {
	return alterBlock(AlterClauseSet, items)
}

// AlterUnset is the UNSET block. It takes attributes or groups the same way as AlterSet.
func AlterUnset(items ...AlterItem) AlterClause {
	return alterBlock(AlterClauseUnset, items)
}

// alterBlock builds a SET or UNSET block. Attributes passed directly become one main group.
// A call that mixes attributes and groups is recorded, and ValidateObject rejects it with the object name.
func alterBlock(kind AlterClauseKind, items []AlterItem) AlterClause {
	clause := AlterClause{Kind: kind}
	var direct []*Attribute
	hasDirect := false
	for _, item := range items {
		switch typed := item.(type) {
		case *Attribute:
			direct = append(direct, typed)
			hasDirect = true
		case AlterGroupSpec:
			clause.Groups = append(clause.Groups, typed)
		case nil:
			direct = append(direct, nil)
			hasDirect = true
		}
	}
	switch {
	case hasDirect && len(clause.Groups) > 0:
		clause.MixedForms = true
	case hasDirect:
		clause.Groups = []AlterGroupSpec{{Name: MainAlterGroup, Attributes: direct}}
	}
	return clause
}

// AlterGroupSpec is a set of SET (or UNSET) attributes that one ALTER statement sends. Make it with AlterGroup.
type AlterGroupSpec struct {
	Name       string
	Attributes []*Attribute
}

// AlterGroup is one ALTER statement of a SET or UNSET block. A block sends its groups in declaration order, and
// skips a group with no values. Name is MainAlterGroup for the main group, whose request is set (unset). Any other
// name is a Go identifier that starts with a capital letter, and its request field is set<Name> (unset<Name>),
// such as setBody. A block has at most one main group. When it has none, an empty one is sent last.
// Attributes of different groups cannot share a statement: the SDK SET (UNSET) struct holds every group, and its
// validation rejects a value that sets fields of two groups.
func AlterGroup(name string, attributes ...*Attribute) AlterGroupSpec {
	return AlterGroupSpec{Name: name, Attributes: attributes}
}

// AlterSetTags is the SET TAG block.
func AlterSetTags() AlterClause { return AlterClause{Kind: AlterClauseSetTags} }

// AlterUnsetTags is the UNSET TAG block.
func AlterUnsetTags() AlterClause { return AlterClause{Kind: AlterClauseUnsetTags} }

// AlterSpec is the ALTER statement order after the name.
// Doc overrides the generated doc URL. AdditionalValidations marks the query struct.
// A variant's SQL words are Variant.AlterName, or the base object type when that override is empty.
type AlterSpec struct {
	Doc                      string
	AdditionalValidations    bool
	SetAdditionalValidations bool
	Clauses                  []AlterClause
}

// LinkedAlterKeyword returns the ALTER keyword clause linked to the value of the bool attribute with this name.
func (obj Object) LinkedAlterKeyword(attributeName string, value bool) (AlterClause, bool) {
	for _, clause := range obj.Alter.Clauses {
		if clause.Kind == AlterClauseKeyword && clause.Attribute != nil && clause.Attribute.Name == attributeName && clause.AttributeValue == value {
			return clause, true
		}
	}
	return AlterClause{}, false
}

// SetAttributes are the SET attributes, in the order of the groups of the SET block.
func (obj Object) SetAttributes() []*Attribute {
	return alterGroupAttributes(obj.alterGroups(AlterClauseSet))
}

// UnsetAttributes are the UNSET attributes, in the order of the groups of the UNSET block.
func (obj Object) UnsetAttributes() []*Attribute {
	return alterGroupAttributes(obj.alterGroups(AlterClauseUnset))
}

// AlterSetGroups are the groups of the SET block in declaration order, the main group included. A block that lists its
// attributes directly has one main group.
func (obj Object) AlterSetGroups() []AlterGroupSpec {
	return obj.alterGroups(AlterClauseSet)
}

// AlterUnsetGroups are the groups of the UNSET block, the same way as AlterSetGroups.
func (obj Object) AlterUnsetGroups() []AlterGroupSpec {
	return obj.alterGroups(AlterClauseUnset)
}

// HasAlterSetBlock reports whether the ALTER spec has a SET block.
func (obj Object) HasAlterSetBlock() bool {
	_, ok := obj.alterClause(AlterClauseSet)
	return ok
}

// HasAlterUnsetBlock reports whether the ALTER spec has an UNSET block.
func (obj Object) HasAlterUnsetBlock() bool {
	_, ok := obj.alterClause(AlterClauseUnset)
	return ok
}

// alterClause is the SET or UNSET block, if the object has one.
func (obj Object) alterClause(kind AlterClauseKind) (AlterClause, bool) {
	for _, clause := range obj.Alter.Clauses {
		if clause.Kind == kind {
			return clause, true
		}
	}
	return AlterClause{}, false
}

// alterGroups are the groups of the SET or UNSET block, or nil when the object has no such block.
func (obj Object) alterGroups(kind AlterClauseKind) []AlterGroupSpec {
	clause, _ := obj.alterClause(kind)
	return clause.Groups
}

func alterGroupAttributes(groups []AlterGroupSpec) []*Attribute {
	var out []*Attribute
	for _, group := range groups {
		out = append(out, group.Attributes...)
	}
	return out
}

// HasAlterSpec reports whether this object defines an ALTER statement.
// Presence is Clauses. Doc and the validation flags do not define it.
func (obj Object) HasAlterSpec() bool {
	return len(obj.Alter.Clauses) > 0
}

// CreateSpec is the CREATE statement.
// Clauses is the order after the name.
// KeywordsBefore is the order after OR REPLACE and before the object words.
// Doc and Sql override the generated doc URL and object words. Empty Sql is the uppercase words of Name.
// A variant does not set Sql. Variant CREATE uses Variant.CreateName, or Variant.Of.ObjectType when that override is empty.
// WithoutOrReplace leaves OR REPLACE out of a variant CREATE. The base uses SdkCreateWithoutOrReplace.
type CreateSpec struct {
	Doc                   string
	Sql                   string
	AdditionalValidations bool
	WithoutOrReplace      bool
	Clauses               []*Attribute
	KeywordsBefore        []*Attribute
}

// HasCreateSpec reports whether this object defines a CREATE statement.
// Presence is Clauses. Doc, Sql, KeywordsBefore, WithoutOrReplace, and AdditionalValidations do not define it.
func (obj Object) HasCreateSpec() bool {
	return len(obj.Create.Clauses) > 0
}

// ResourceAttributes is Attributes without SQL-only fields.
// Schema, resource, model, and assertion generators use this list.
func (obj Object) ResourceAttributes() []Attribute {
	return collections.Filter(obj.Attributes, func(attr Attribute) bool {
		return !attr.SqlOnly
	})
}

func (obj Object) ObjectName() string { return obj.Name }

// ParameterLevel is the catalog level for this object, derived from Name (Schema is ParameterLevelSchema).
// Variant.Of does not change this level. A variant that needs a different catalog level
// will require an explicit SSOT parameter-level override.
func (obj Object) ParameterLevel() ParameterLevel {
	return ParameterLevel(strings.ToUpper(genhelpers.ToSnakeCase(obj.Name)))
}

// ParameterAssertionNames is the SDK parameter constant prefix and object type name for assertions.
// Both are empty when the object is not a variant. A variant uses Variant.Of.Name for both.
// The catalog level stays ParameterLevel.
type ParameterAssertionNames struct {
	ConstantPrefix string
	ObjectTypeName string
}

// ParameterAssertionNames returns the assertion names for this object.
func (obj Object) ParameterAssertionNames() ParameterAssertionNames {
	base := obj.VariantBase()
	if base == nil {
		return ParameterAssertionNames{}
	}
	if base.Name == "" {
		panic(fmt.Sprintf("object %q: Variant.Of has an empty Name", obj.Name))
	}
	return ParameterAssertionNames{
		ConstantPrefix: base.Name,
		ObjectTypeName: base.Name,
	}
}

// SuspendForAlter describes an object that has to be suspended while it is altered, such as an alert or a warehouse that
// changes its type. It reuses the SHOW state column, its enum, and the object's SUSPEND and RESUME keyword clauses.
type SuspendForAlter struct {
	// When lists the attributes whose change needs the suspension. Empty means every attribute that changes in place.
	When []*Attribute
	// State is the SHOW column with the state of the object. Started and Suspended are values of its enum.
	State              *Attribute
	Started, Suspended string
	// Enabled is the bool attribute that says whether the object runs, such as an alert's enabled. After a successful
	// update the object is resumed when it is true and stays suspended otherwise. Nil restores the state it had before.
	Enabled *Attribute
	// NeededExt makes the guard also ask <prefix>SuspendForAlterNeededExt(d) of package resources whether the change needs
	// the suspension, after one of the When attributes changed. Use it when a change in Terraform can leave Snowflake as it
	// is, such as a removed warehouse_type that sets the type the warehouse already has.
	NeededExt bool
}

// SuspendForAlterWhen returns the names of the attributes whose change suspends the object: SuspendForAlter.When, or every
// attribute that changes in place when When is empty. It is nil without SuspendForAlter.
func (obj Object) SuspendForAlterWhen() []string {
	if obj.SuspendForAlter == nil {
		return nil
	}
	if len(obj.SuspendForAlter.When) > 0 {
		return collections.Map(obj.SuspendForAlter.When, func(attr *Attribute) string { return attr.Name })
	}
	inPlace := collections.Filter(obj.ResourceAttributes(), Attribute.ChangesInPlace)
	return collections.Map(inPlace, func(attr Attribute) string { return attr.Name })
}
