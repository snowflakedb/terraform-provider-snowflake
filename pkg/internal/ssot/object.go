package ssot

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

type Object struct {
	Name                                string
	NamePlural                          string
	DescribeFailRead                    bool
	DescribeUsesShowId                  bool
	HasPartialOnUpdate                  bool
	AdditionalDescribeOutputDescription string
	IdType                              string
	// ParameterSnowflakeType corresponds to the sdk.ParameterType suffix Snowflake reports for this
	// object's parameters ("Schema" emits sdk.ParameterTypeSchema).
	ParameterSnowflakeType            string
	IdentityNameForceNew              bool
	AdditionalIdentityNameDescription string
	Attributes                        []Attribute
	AttributesModification            bool
	ResourceHooks                     []ResourceHookOption
	Generators                        []GeneratorOption
	Show                              OutputSpec
	Describe                          OutputSpec
}

func (obj Object) ObjectName() string { return obj.Name }

// ParameterLevel is the catalog consumer level for this object (Schema is ParameterLevelSchema).
func (obj Object) ParameterLevel() ParameterLevel {
	return ParameterLevel(strings.ToUpper(genhelpers.ToSnakeCase(obj.Name)))
}
