package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

// TODO(SNOW-2357735): Remove after complex non-list type overrides are handled
func DynamicTableWithoutTargetLag(
	resourceName string,
	database string,
	schema string,
	name string,
	query string,
	warehouse string,
) *DynamicTableModel {
	d := &DynamicTableModel{ResourceModelMeta: config.Meta(resourceName, resources.DynamicTable)}
	d.WithDatabase(database)
	d.WithSchema(schema)
	d.WithName(name)
	d.WithQuery(query)
	d.WithWarehouse(warehouse)
	return d
}

// TODO(SNOW-2357735): Remove after complex non-list type overrides are handled
func (d *DynamicTableModel) WithTargetLag(tl []string) *DynamicTableModel {
	_ = tl
	return d
}

func (d *DynamicTableModel) WithMaximumDurationTargetLag(value string) *DynamicTableModel {
	return d.WithTargetLagValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"maximum_duration": tfconfig.StringVariable(value),
	}))
}
