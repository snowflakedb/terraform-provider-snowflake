package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (t *UserProgrammaticAccessTokenModel) WithKeepers(value map[string]string) *UserProgrammaticAccessTokenModel {
	keepersVariables := make(map[string]tfconfig.Variable)
	for k, v := range value {
		keepersVariables[k] = tfconfig.StringVariable(v)
	}
	t.Keepers = tfconfig.MapVariable(keepersVariables)
	return t
}
