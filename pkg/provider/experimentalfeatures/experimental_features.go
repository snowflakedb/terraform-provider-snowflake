package experimentalfeatures

import (
	"fmt"
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
	UserEnableDefaultWorkloadIdentity              ExperimentalFeature = "USER_ENABLE_DEFAULT_WORKLOAD_IDENTITY"
	GrantsImportValidation                         ExperimentalFeature = "GRANTS_IMPORT_VALIDATION"
	// TODO [SNOW-2739299]: Discuss having an additional ParametersNoOutput experiment
	ParametersReducedOutput             ExperimentalFeature = "PARAMETERS_REDUCED_OUTPUT"
	TagNewTriValueAllowedValuesBehavior ExperimentalFeature = "TAG_NEW_TRI_VALUE_ALLOWED_VALUES_BEHAVIOR"
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
			fmt.Sprintf("This feature works independently of the `%s` flag.", GrantsImportValidation),
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
	{
		ParametersReducedOutput,
		ExperimentalFeatureStateActive,
		joinWithDoubleNewline(
			"Currently, the `parameters` field in various resources contains a verbatim output for the `SHOW PARAMETERS IN <object>` command. One of the fields contained in the output is the `description`. It does not change and is repeated for all objects containing the given parameter. It leads to an excessive output (check e.g., [#3118](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3118)).",
			"To mitigate the problem, we are adding this option to reduce the output to only `value` and `level` fields, which should significantly reduce the state size. **Note**: it's also affecting the `parameters` output for data sources.",
			"We considered the option to remove the `parameters` output completely, however, we plan to change the external change logic detection to use it (to make it consistent with other attributes using `show_output` and because we won't be able to implement the current logic when switching to the Terraform Plugin Framework) and it still allows referencing the parameter value/level from other parts of the configuration.",
		),
	},
	{
		UserEnableDefaultWorkloadIdentity,
		ExperimentalFeatureStateActive,
		joinWithDoubleNewline(
			"The new `default_workload_identity_federation` field was added to the `snowflake_legacy_service_user` and `snowflake_service_user` resources. This field allows for managing WIFs. Due to feature complexity, it requires enabling this experiment.",
			"Read more in our [migration guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/dev/MIGRATION_GUIDE.md#new-feature-workload-identity-federation-support-for-service-users).",
		),
	},
	{
		GrantsImportValidation,
		ExperimentalFeatureStateActive,
		joinWithDoubleNewline(
			"Enables import validation for the `snowflake_grant_privileges_to_account_role` resource.",
			"When enabled, importing a grant resource with a fixed set of privileges (`privileges` field) will validate that the specified privileges actually exist in Snowflake with the correct `with_grant_option` setting, and error immediately if they don't match.",
			fmt.Sprintf("This feature works independently of the `%s` flag.", GrantsStrictPrivilegeManagement),
		),
	},
	{
		TagNewTriValueAllowedValuesBehavior,
		ExperimentalFeatureStateActive,
		joinWithDoubleNewline(
			"Enables behavior changes for the `allowed_values` field in the `snowflake_tag` resource.",
			"When enabled, the three possible states in Snowflake for allowed values will be supported: `nil` (any value is allowed; whenever `allowed_values` are empty), `empty` (no value is allowed; handled by the `no_allowed_values` field), and `set` (all values defined in `allowed_values` are allowed).",
			"Otherwise, the `no_allowed_values` field will be ignored (explicit changes will cause updates, but without any effect) and the `allowed_values` field will follow the old behavior: `nil` (any value is allowed; only available whenever tag resource is created without `allowed_values`), `empty` (no value is allowed; always set when updating from filled `allowed_values` set to empty one or completely removed from config), `set` (all values defined in `allowed_values` are allowed).",
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

var (
	ActiveExperiments       = collections.Filter(allExperiments, filterByStateFuncProvider(ExperimentalFeatureStateActive))
	DiscontinuedExperiments = collections.Filter(allExperiments, filterByStateFuncProvider(ExperimentalFeatureStateDiscontinued))
)

var (
	allExperimentalFeatureNames    = collections.Map(allExperiments, mapToName)
	activeExperimentalFeatureNames = collections.Map(ActiveExperiments, mapToName)
)

var (
	AllExperimentalFeatureNames    = sdk.AsStringList(allExperimentalFeatureNames)
	ActiveExperimentalFeatureNames = sdk.AsStringList(activeExperimentalFeatureNames)
)

func IsExperimentEnabled(experiment ExperimentalFeature, enabledExperiments []string) bool {
	return slices.ContainsFunc(enabledExperiments, func(s string) bool {
		return strings.EqualFold(string(experiment), s)
	})
}
