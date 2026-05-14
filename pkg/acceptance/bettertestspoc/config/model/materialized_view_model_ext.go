package model

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func MaterializedViewWithId(
	resourceName string,
	id sdk.SchemaObjectIdentifier,
	statement string,
	warehouse string,
) *MaterializedViewModel {
	return MaterializedView(resourceName, id.DatabaseName(), id.SchemaName(), id.Name(), statement, warehouse)
}

func (m *MaterializedViewModel) WithTagReference(tagReference string, value string) *MaterializedViewModel {
	m.Tag = tfconfig.ListVariable(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"name":     config.UnquotedWrapperVariable(fmt.Sprintf("%s.name", tagReference)),
			"schema":   config.UnquotedWrapperVariable(fmt.Sprintf("%s.schema", tagReference)),
			"database": config.UnquotedWrapperVariable(fmt.Sprintf("%s.database", tagReference)),
			"value":    tfconfig.StringVariable(value),
		}),
	)
	return m
}
