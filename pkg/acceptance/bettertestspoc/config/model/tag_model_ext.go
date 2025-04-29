package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func TagBase(resourceName string, tagId sdk.SchemaObjectIdentifier) *TagModel {
	return Tag(resourceName, tagId.DatabaseName(), tagId.SchemaName(), tagId.Name())
}

func (t *TagModel) WithAllowedValues(allowedValues ...string) *TagModel {
	allowedValuesStringVariables := make([]tfconfig.Variable, len(allowedValues))
	for i, v := range allowedValues {
		allowedValuesStringVariables[i] = tfconfig.StringVariable(v)
	}

	t.AllowedValues = tfconfig.SetVariable(allowedValuesStringVariables...)
	return t
}

func (t *TagModel) WithMaskingPolicies(maskingPolicies ...sdk.SchemaObjectIdentifier) *TagModel {
	maskingPoliciesStringVariables := make([]tfconfig.Variable, len(maskingPolicies))
	for i, v := range maskingPolicies {
		maskingPoliciesStringVariables[i] = tfconfig.StringVariable(v.FullyQualifiedName())
	}

	t.MaskingPolicies = tfconfig.SetVariable(maskingPoliciesStringVariables...)
	return t
}
