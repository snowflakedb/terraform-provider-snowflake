package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

func (p *PostgresMirrorModel) WithPostgresTables(tables ...string) *PostgresMirrorModel {
	p.PostgresTables = postgresMirrorNameSet(tables)
	return p
}

func (p *PostgresMirrorModel) WithPostgresSchemas(schemas ...string) *PostgresMirrorModel {
	p.PostgresSchemas = postgresMirrorNameSet(schemas)
	return p
}

func (p *PostgresMirrorModel) WithRefreshIntervalSeconds(seconds int) *PostgresMirrorModel {
	return p.WithRefreshIntervalValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"seconds": tfconfig.IntegerVariable(seconds),
		}),
	)
}

func (p *PostgresMirrorModel) WithRefreshIntervalMinutes(minutes int) *PostgresMirrorModel {
	return p.WithRefreshIntervalValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"minutes": tfconfig.IntegerVariable(minutes),
		}),
	)
}

func postgresMirrorNameSet(names []string) tfconfig.Variable {
	return tfconfig.SetVariable(collections.Map(names, func(name string) tfconfig.Variable { return tfconfig.StringVariable(name) })...)
}
