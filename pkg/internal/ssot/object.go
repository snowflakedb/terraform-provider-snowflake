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
	IdentityNameForceNew                bool
	AdditionalIdentityNameDescription   string
	Attributes                          []Attribute
	AttributesModification              bool
	ResourceHooks                       []ResourceHookOption
	Generators                          []GeneratorOption
	SchemaStateUpgraders                []StateUpgraderFunc
	Show                                OutputSpec
	Describe                            OutputSpec
	// Create, Set, and Unset are the SQL field order for those statements.
	// Presence in the slice is membership. The slice order is the order generation emits.
	Create []*Attribute
	Set    []*Attribute
	Unset  []*Attribute
}

func (obj Object) ObjectName() string { return obj.Name }

// ParameterLevel is the catalog level for this object, derived from Name (Schema is ParameterLevelSchema).
func (obj Object) ParameterLevel() ParameterLevel {
	return ParameterLevel(strings.ToUpper(genhelpers.ToSnakeCase(obj.Name)))
}
