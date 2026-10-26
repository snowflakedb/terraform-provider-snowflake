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
	IdType                              string
	LegacyIdEncoding                    bool // Keeps the pipe-delimited resource id (db|schema|name) that older resources store.
	IdentityNameForceNew                bool
	IdentityHierarchyForceNew           bool // database and schema identifiers are ForceNew; the identity schema template does not emit this yet
	AdditionalIdentityNameDescription   string
	SdkCreateRequestObjectName          string  // SDK Create<Name>Request stem; defaults to Name
	VariantOf                           *Object // Shares SDK parameter details, set/unset request, and ParameterType names with this object; catalog membership still uses Name
	// VariantEntityCreateName overrides the SQL object words a variant CREATE emits.
	// Empty uses VariantOf.ObjectType, so WarehouseAdaptive emits WAREHOUSE. Interactive CREATE sets INTERACTIVE WAREHOUSE.
	VariantEntityCreateName string
	// VariantEntityAlterName overrides the SQL object words a variant ALTER emits.
	// Empty uses VariantOf.ObjectType, the same default as CREATE.
	VariantEntityAlterName   string
	Attributes               []Attribute
	AttributesModification   bool
	ResourceHooks            []ResourceHookOption
	Generators               []GeneratorOption
	SchemaStateUpgraders     []StateUpgraderFunc
	AdditionalCustomizeDiffs []CustomizeDiffFunc
	Show                     OutputSpec
	// ShowGeneric is the temporary union of every variant SHOW column.
	// It is the one paired SHOW. Each variant, including the main one, declares its own show in Show.
	// TODO [next PRs]: Each variant'sShow                  .Columns becomes a paired SHOW with its dedicated type, and ShowGeneric goes away.
	ShowGeneric OutputSpec
	Describe    OutputSpec // ShowPlainName is the plain SHOW struct name when it differs from Name.
	// Empty uses Name. Set it when ShowGeneric uses Name for the paired SHOW struct.
	ShowPlainName string
	// ShowVariants are plain SHOW structs emitted with this object's DESCRIBE.
	// Each variant's Show.Columns is that struct. The variant does not set SdkShow.
	// TODO [next PRs]: Emit these from ShowOperationFromObject. See DescribeOperationFromObject.
	ShowVariants []*Object
	// Create, Set, Unset, and Alter are the SQL field order for those statements.
	// Presence in a slice is membership. The slice order is the order generation emits.
	// Set and Unset are the fields inside those blocks. Alter is the order of the blocks and action keywords after the name.
	// A field in those slices may be absent from Attributes when it is not a resource attribute.
	Create CreateSpec
	Set    []*Attribute
	Unset  []*Attribute
	Alter  AlterSpec
	// CustomOperations are operations other than CREATE and ALTER.
	// Clause order is the emission order. SdkCustomOperations turns them on.
	CustomOperations []CustomOperationSpec
	// CustomInterfaceMethods are interface declarations whose bodies stay in the SDK ext file.
	// SdkCustomInterfaceMethods turns them on.
	CustomInterfaceMethods []CustomInterfaceMethodSpec
	// ShowParametersDetailsAdditionalLevels adds parameters from these levels onto ShowParametersDetails.
	// A name already in ParameterLevel() is kept once.
	ShowParametersDetailsAdditionalLevels []ParameterLevel
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
	AlterClauseKeyword   AlterClauseKind = "keyword"
	AlterClauseRenameTo  AlterClauseKind = "rename_to"
	AlterClauseSwapWith  AlterClauseKind = "swap_with"
	AlterClauseSet       AlterClauseKind = "set"
	AlterClauseUnset     AlterClauseKind = "unset"
	AlterClauseSetTags   AlterClauseKind = "set_tags"
	AlterClauseUnsetTags AlterClauseKind = "unset_tags"
	AlterClauseList      AlterClauseKind = "list"
)

// AlterClause is one ALTER entry. Keyword and list use SQL. The other kinds are the SET, UNSET, SET TAG, and UNSET TAG blocks.
type AlterClause struct {
	Kind AlterClauseKind
	SQL  string
	// NotExactlyOne leaves a keyword out of ExactlyOneValueSet. It may accompany another action.
	NotExactlyOne bool
	// ListItem is the Go type of a list element, without a pointer.
	ListItem string
}

// AlterKeyword is an optional action keyword, such as RESUME or ENABLE MANAGED ACCESS.
// It is one branch of ExactlyOneValueSet.
func AlterKeyword(sql string) AlterClause {
	return AlterClause{Kind: AlterClauseKeyword, SQL: sql}
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

// AlterRenameTo is RENAME TO. The identifier type is the object's IdType.
func AlterRenameTo() AlterClause { return AlterClause{Kind: AlterClauseRenameTo} }

// AlterSwapWith is SWAP WITH. The identifier type is the object's IdType.
func AlterSwapWith() AlterClause { return AlterClause{Kind: AlterClauseSwapWith} }

// AlterSet is the SET block. Its fields are Object.Set.
func AlterSet() AlterClause { return AlterClause{Kind: AlterClauseSet} }

// AlterUnset is the UNSET block. Its fields are Object.Unset.
func AlterUnset() AlterClause { return AlterClause{Kind: AlterClauseUnset} }

// AlterSetTags is the SET TAG block.
func AlterSetTags() AlterClause { return AlterClause{Kind: AlterClauseSetTags} }

// AlterUnsetTags is the UNSET TAG block.
func AlterUnsetTags() AlterClause { return AlterClause{Kind: AlterClauseUnsetTags} }

// AlterSpec is the ALTER statement order after the name.
// Doc overrides the generated doc URL. AdditionalValidations marks the query struct.
// A variant's SQL words are VariantEntityAlterName, or VariantOf.ObjectType when that override is empty.
type AlterSpec struct {
	Doc                      string
	AdditionalValidations    bool
	SetAdditionalValidations bool
	Clauses                  []AlterClause
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
// A variant does not set Sql. Variant CREATE uses VariantEntityCreateName, or VariantOf.ObjectType when that override is empty.
type CreateSpec struct {
	Doc                   string
	Sql                   string
	AdditionalValidations bool
	Clauses               []*Attribute
	KeywordsBefore        []*Attribute
}

// HasCreateSpec reports whether this object defines a CREATE statement.
// Presence is Clauses. Doc, Sql, KeywordsBefore, and AdditionalValidations do not define it.
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
// VariantOf does not change this level. A variant that needs a different catalog level
// will require an explicit SSOT parameter-level override.
func (obj Object) ParameterLevel() ParameterLevel {
	return ParameterLevel(strings.ToUpper(genhelpers.ToSnakeCase(obj.Name)))
}

// ParameterAssertionNames is the SDK parameter constant prefix and object type name for assertions.
// Both are empty when the object is not a variant. A variant uses VariantOf.Name for both.
// The catalog level stays ParameterLevel.
type ParameterAssertionNames struct {
	ConstantPrefix string
	ObjectTypeName string
}

// ParameterAssertionNames returns the assertion names for this object.
func (obj Object) ParameterAssertionNames() ParameterAssertionNames {
	if obj.VariantOf == nil {
		return ParameterAssertionNames{}
	}
	if obj.VariantOf.Name == "" {
		panic(fmt.Sprintf("object %q: VariantOf has an empty Name", obj.Name))
	}
	return ParameterAssertionNames{
		ConstantPrefix: obj.VariantOf.Name,
		ObjectTypeName: obj.VariantOf.Name,
	}
}
