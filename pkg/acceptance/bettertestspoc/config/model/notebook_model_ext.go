package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (n *NotebookModel) WithFrom(path string, stageId sdk.SchemaObjectIdentifier) *NotebookModel {
	n.From = tfconfig.ListVariable(
		tfconfig.MapVariable(map[string]tfconfig.Variable{
			"stage": tfconfig.StringVariable(stageId.FullyQualifiedName()),
			"path":  tfconfig.StringVariable(path),
		}))
	return n
}
