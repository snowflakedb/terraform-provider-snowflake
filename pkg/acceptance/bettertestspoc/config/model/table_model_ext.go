package model

import (
	"fmt"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

func TableWithId(
	resourceName string,
	tableId sdk.SchemaObjectIdentifier,
	column []sdk.TableColumnSignature,
) *TableModel {
	return Table(resourceName, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), column)
}

func TableWithImplicitDependencies(
	resourceName string,
	tableName string,
	column []sdk.TableColumnSignature,
	schemaModel *SchemaModel,
	databaseModel *DatabaseModel,
) *TableModel {
	return Table(resourceName, "", "", tableName, column).
		WithDatabaseValue(config.UnquotedWrapperVariable(fmt.Sprintf("%s.name", databaseModel.ResourceReference()))).
		WithSchemaValue(config.UnquotedWrapperVariable(fmt.Sprintf("%s.name", schemaModel.ResourceReference())))
}

func (t *TableModel) WithColumn(column []sdk.TableColumnSignature) *TableModel {
	columns := make([]tfconfig.Variable, len(column))
	for idx, v := range column {
		columns[idx] = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"name": tfconfig.StringVariable(v.Name),
			"type": tfconfig.StringVariable(v.Type.ToSql()),
		})
	}
	t.Column = tfconfig.ListVariable(columns...)
	return t
}

// TableColumnConfig is a richer column definition used in tests that need
// column fields beyond name and type (comment, as, not_null, default expression, collate, masking policy). A nil Collate leaves the field out of the config.
type TableColumnConfig struct {
	Name    string
	Type    datatypes.DataType
	Comment string
	// As is the expression of a virtual column. An empty value leaves the field out of the config.
	As string
	// DefaultExpression is a pointer to tell an unset default from an empty one ("" is an empty string default).
	DefaultExpression  *string
	NotNull            *bool
	Collate            *string
	MaskingPolicy      *sdk.SchemaObjectIdentifier
	MaskingPolicyUsing []string
	ProjectionPolicy   *sdk.SchemaObjectIdentifier
	// MaskingPolicyUnquoted renders the masking policy name without double quotes (database.schema.name).
	MaskingPolicyUnquoted bool
}

func (t *TableModel) WithColumnConfigs(columns ...TableColumnConfig) *TableModel {
	vars := make([]tfconfig.Variable, len(columns))
	for idx, c := range columns {
		m := map[string]tfconfig.Variable{
			"name": tfconfig.StringVariable(c.Name),
			"type": tfconfig.StringVariable(c.Type.ToSql()),
		}
		if c.Comment != "" {
			m["comment"] = tfconfig.StringVariable(c.Comment)
		}
		if c.As != "" {
			m["as"] = tfconfig.StringVariable(c.As)
		}
		if c.DefaultExpression != nil {
			m["default"] = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"expression": tfconfig.StringVariable(*c.DefaultExpression),
			})
		}
		if c.NotNull != nil {
			m["not_null"] = tfconfig.BoolVariable(*c.NotNull)
		}
		if c.Collate != nil {
			m["collate"] = tfconfig.StringVariable(*c.Collate)
		}
		if c.MaskingPolicy != nil {
			policyName := c.MaskingPolicy.FullyQualifiedName()
			if c.MaskingPolicyUnquoted {
				policyName = fmt.Sprintf("%s.%s.%s", c.MaskingPolicy.DatabaseName(), c.MaskingPolicy.SchemaName(), c.MaskingPolicy.Name())
			}
			maskingPolicy := map[string]tfconfig.Variable{
				"policy_name": tfconfig.StringVariable(policyName),
			}
			if len(c.MaskingPolicyUsing) > 0 {
				maskingPolicy["using"] = tfconfig.ListVariable(collections.Map(c.MaskingPolicyUsing, func(s string) tfconfig.Variable {
					return tfconfig.StringVariable(s)
				})...)
			}
			m["masking_policy"] = tfconfig.ObjectVariable(maskingPolicy)
		}
		if c.ProjectionPolicy != nil {
			m["projection_policy"] = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"policy_name": tfconfig.StringVariable(c.ProjectionPolicy.FullyQualifiedName()),
			})
		}
		vars[idx] = tfconfig.ObjectVariable(m)
	}
	t.Column = tfconfig.ListVariable(vars...)
	return t
}

func (t *TableModel) WithClusterBy(clusterBy ...string) *TableModel {
	return t.WithClusterByValue(tfconfig.ListVariable(collections.Map(clusterBy, func(s string) tfconfig.Variable {
		return tfconfig.StringVariable(s)
	})...))
}

func (t *TableModel) WithRowAccessPolicy(rap sdk.SchemaObjectIdentifier, on string) *TableModel {
	return t.WithRowAccessPolicyValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"policy_name": tfconfig.StringVariable(rap.FullyQualifiedName()),
				"on":          tfconfig.ListVariable(tfconfig.StringVariable(on)),
			},
		),
	)
}

func (t *TableModel) WithAggregationPolicy(ap sdk.SchemaObjectIdentifier, entityKey ...string) *TableModel {
	m := map[string]tfconfig.Variable{
		"policy_name": tfconfig.StringVariable(ap.FullyQualifiedName()),
	}
	if len(entityKey) > 0 {
		m["entity_key"] = tfconfig.ListVariable(collections.Map(entityKey, func(s string) tfconfig.Variable {
			return tfconfig.StringVariable(s)
		})...)
	}
	return t.WithAggregationPolicyValue(
		tfconfig.ObjectVariable(
			m,
		),
	)
}

func (t *TableModel) WithJoinPolicy(jp sdk.SchemaObjectIdentifier, allowedJoinKeys ...string) *TableModel {
	m := map[string]tfconfig.Variable{
		"policy_name": tfconfig.StringVariable(jp.FullyQualifiedName()),
	}
	if len(allowedJoinKeys) > 0 {
		m["allowed_join_keys"] = tfconfig.ListVariable(collections.Map(allowedJoinKeys, func(s string) tfconfig.Variable {
			return tfconfig.StringVariable(s)
		})...)
	}
	return t.WithJoinPolicyValue(
		tfconfig.ObjectVariable(
			m,
		),
	)
}
