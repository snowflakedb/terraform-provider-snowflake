package ssot

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

// StateUpgraderFunc names a function in package resources with the schema.StateUpgraderFunc signature.
// Object.SchemaStateUpgraders lists these functions in version order.
// The version each function upgrades from is its index.
// The resource schema version is the length of the slice.
type StateUpgraderFunc string

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
	AdditionalIdentityNameDescription   string
	SdkCreateRequestObjectName          string  // SDK Create<Name>Request stem; defaults to Name
	VariantOf                           *Object // Shares SDK parameter details, set/unset request, and ParameterType names with this object; catalog membership still uses Name
	Attributes                          []Attribute
	AttributesModification              bool
	ResourceHooks                       []ResourceHookOption
	Generators                          []GeneratorOption
	SchemaStateUpgraders                []StateUpgraderFunc
	Show                                OutputSpec
	Describe                            OutputSpec
	// Create, Set, Unset, and Alter are the SQL field order for those statements.
	// Presence in a slice is membership. The slice order is the order generation emits.
	// Set and Unset are the fields inside those blocks. Alter is the order of the blocks and action keywords after the name.
	Create CreateSpec
	Set    []*Attribute
	Unset  []*Attribute
	Alter  AlterSpec
	// CustomOperations are operations other than CREATE and ALTER.
	// Clause order is the emission order. SdkCustomOperations turns them on.
	CustomOperations []CustomOperationSpec
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
)

// AlterClause is one ALTER entry. Keyword uses SQL. The other kinds are the SET, UNSET, SET TAG, and UNSET TAG blocks.
type AlterClause struct {
	Kind AlterClauseKind
	SQL  string
}

// AlterKeyword is an optional action keyword, such as RESUME or ENABLE MANAGED ACCESS.
func AlterKeyword(sql string) AlterClause {
	return AlterClause{Kind: AlterClauseKeyword, SQL: sql}
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
type AlterSpec struct {
	Clauses []AlterClause
}

// CreateSpec is the CREATE statement.
// Clauses is the order after the name.
// KeywordsBefore is the order after OR REPLACE and before the object words.
type CreateSpec struct {
	Clauses        []*Attribute
	KeywordsBefore []*Attribute
}

func (obj Object) ObjectName() string { return obj.Name }

// ParameterLevel is the catalog level for this object, derived from Name (Schema is ParameterLevelSchema).
// VariantOf does not change this level. A variant that needs a different catalog level
// will require an explicit SSOT parameter-level override.
func (obj Object) ParameterLevel() ParameterLevel {
	return ParameterLevel(strings.ToUpper(genhelpers.ToSnakeCase(obj.Name)))
}
