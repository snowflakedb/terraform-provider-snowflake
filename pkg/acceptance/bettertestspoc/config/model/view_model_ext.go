package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
)

func (v *ViewModel) WithColumnNames(columnNames ...string) *ViewModel {
	return v.WithColumnValue(config.TupleVariable(
		collections.Map(columnNames, func(columnName string) config.Variable {
			return config.ObjectVariable(
				map[string]config.Variable{
					"column_name": config.StringVariable(columnName),
				},
			)
		})...,
	))
}

func (v *ViewModel) WithRowAccessPolicy(rap sdk.SchemaObjectIdentifier, on string) *ViewModel {
	return v.WithRowAccessPolicyValue(
		config.ObjectVariable(
			map[string]config.Variable{
				"policy_name": config.StringVariable(rap.FullyQualifiedName()),
				"on":          config.ListVariable(config.StringVariable(on)),
			},
		),
	)
}

func (v *ViewModel) WithAggregationPolicy(ap sdk.SchemaObjectIdentifier, key string) *ViewModel {
	return v.WithAggregationPolicyValue(
		config.ObjectVariable(
			map[string]config.Variable{
				"policy_name": config.StringVariable(ap.FullyQualifiedName()),
				"entity_key":  config.ListVariable(config.StringVariable(key)),
			},
		),
	)
}

func (v *ViewModel) WithDataMetricFunction(functionId sdk.SchemaObjectIdentifier, on string, scheduleStatus sdk.DataMetricScheduleStatusOption) *ViewModel {
	return v.WithDataMetricFunctionValue(
		config.ObjectVariable(
			map[string]config.Variable{
				"function_name":   config.StringVariable(functionId.FullyQualifiedName()),
				"on":              config.ListVariable(config.StringVariable(on)),
				"schedule_status": config.StringVariable(string(scheduleStatus)),
			},
		),
	)
}

func (v *ViewModel) WithDataMetricSchedule(cron string) *ViewModel {
	return v.WithDataMetricScheduleValue(
		config.ObjectVariable(
			map[string]config.Variable{
				"using_cron": config.StringVariable(cron),
			},
		),
	)
}
