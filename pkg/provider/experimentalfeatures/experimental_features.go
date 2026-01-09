package experimentalfeatures

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ExperimentalFeature string

const (
	ParametersIgnoreValueChangesIfNotOnObjectLevel ExperimentalFeature = "PARAMETERS_IGNORE_VALUE_CHANGES_IF_NOT_ON_OBJECT_LEVEL"
	WarehouseShowImprovedPerformance               ExperimentalFeature = "WAREHOUSE_SHOW_IMPROVED_PERFORMANCE"
	GrantsStrictPrivilegeManagement                ExperimentalFeature = "GRANTS_STRICT_PRIVILEGE_MANAGEMENT"
)

type experimentalFeatureState string

const (
	ExperimentalFeatureStateActive       experimentalFeatureState = "ACTIVE"
	ExperimentalFeatureStateDiscontinued experimentalFeatureState = "DISCONTINUED"
)

type Experiment struct {
	name        ExperimentalFeature
	state       experimentalFeatureState
	description string
}

func (e *Experiment) Name() ExperimentalFeature {
	return e.name
}

func (e *Experiment) Description() string {
	return e.description
}

var allExperiments = []Experiment{
	{
		WarehouseShowImprovedPerformance,
		ExperimentalFeatureStateActive,
		joinWithDoubleNewline(
			"It's meant to improve the performance for accounts with many warehouses.",
			"When enabled, it uses a slightly different SHOW query to read warehouse details (`SHOW WAREHOUSES LIKE '<identifier>' STARTS WITH '<identifier>' LIMIT 1`).",
			"**Important**: to benefit from this improvement, you need to have it enabled also on your Snowflake account. To do this, please reach out to us through your Snowflake Account Manager.",
		),
	},
	{
		GrantsStrictPrivilegeManagement,
		ExperimentalFeatureStateActive,
		joinWithDoubleNewline(
			"The new `strict_privilege_management` flag was added to the `snowflake_grant_privileges_to_account_role` resource.",
			"It has similar behavior to the `enable_multiple_grants` flag present in the old grant resources, and it makes the resource able to detect external changes for privileges other than those present in the configuration which can make the `snowflake_grant_privileges_to_account_role` resource a central point of knowledge privilege management for a given object and role.",
			"Read more in our [strict privilege management](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/strict_privilege_management) guide.",
		),
	},
	{
		ParametersIgnoreValueChangesIfNotOnObjectLevel,
		ExperimentalFeatureStateActive,
		joinWithDoubleNewline(
			"Currently, not setting the parameter value on the object level can unnecessarily react to external changes to this parameter's value on the higher levels (e.g. not setting `data_retention_time_in_days` on `snowflake_schema` can result in non-empty plan when the parameter value changes on the database/account level).",
			"When enabled, the provider ignores changes to the parameter value happening on the higher hierarchy levels.",
		),
	},
}

func joinWithDoubleNewline(parts ...string) string {
	return strings.Join(parts, "\n\n")
}

var mapToName = func(e Experiment) ExperimentalFeature {
	return e.name
}

var filterByStateFuncProvider = func(state experimentalFeatureState) func(Experiment) bool {
	return func(e Experiment) bool {
		return e.state == state
	}
}

var ActiveExperiments = collections.Filter(allExperiments, filterByStateFuncProvider(ExperimentalFeatureStateActive))
var DiscontinuedExperiments = collections.Filter(allExperiments, filterByStateFuncProvider(ExperimentalFeatureStateDiscontinued))

var allExperimentalFeatureNames = collections.Map(allExperiments, mapToName)
var activeExperimentalFeatureNames = collections.Map(ActiveExperiments, mapToName)

var AllExperimentalFeatureNames = sdk.AsStringList(allExperimentalFeatureNames)
var ActiveExperimentalFeatureNames = sdk.AsStringList(activeExperimentalFeatureNames)

func IsExperimentEnabled(experiment ExperimentalFeature, enabledExperiments []string) bool {
	return slices.ContainsFunc(enabledExperiments, func(s string) bool {
		return strings.EqualFold(string(experiment), s)
	})
}
