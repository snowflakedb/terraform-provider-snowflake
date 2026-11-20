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
	ParametersReducedOutput              ExperimentalFeature = "PARAMETERS_REDUCED_OUTPUT"
	TagsAllowEmptyAllowedValues          ExperimentalFeature = "TAGS_ALLOW_EMPTY_ALLOWED_VALUES"
	ImportBooleanDefault                 ExperimentalFeature = "IMPORT_BOOLEAN_DEFAULT"
	ImportPermissiveMode                 ExperimentalFeature = "IMPORT_PERMISSIVE_MODE"
	GrantsSafeDestroy                    ExperimentalFeature = "GRANTS_SAFE_DESTROY"
	TagAssociationSafeDestroy            ExperimentalFeature = "TAG_ASSOCIATION_SAFE_DESTROY"
	GrantAccountRoleSafePublicRole       ExperimentalFeature = "GRANT_ACCOUNT_ROLE_SAFE_PUBLIC_ROLE"
	GrantAccountRoleShowCaching          ExperimentalFeature = "GRANT_ACCOUNT_ROLE_SHOW_CACHING"
	AccountRoleShowCaching               ExperimentalFeature = "ACCOUNT_ROLE_SHOW_CACHING"
	GrantsShowCaching                    ExperimentalFeature = "GRANTS_SHOW_CACHING"
	HierarchyRenames                     ExperimentalFeature = "HIERARCHY_RENAMES"
	InheritedGrants                      ExperimentalFeature = "INHERITED_GRANTS"
	ObjectParameterUnsetOnDelete         ExperimentalFeature = "OBJECT_PARAMETER_UNSET_ON_DELETE"
	AuthenticatorExplicitOnly            ExperimentalFeature = "AUTHENTICATOR_EXPLICIT_ONLY"
	ProviderConfigurationAccountFallback ExperimentalFeature = "PROVIDER_CONFIGURATION_ACCOUNT_FALLBACK"
	AccountParameterCatalogWritePath     ExperimentalFeature = "ACCOUNT_PARAMETER_CATALOG_WRITE_PATH"
	ObjectParameterParseOnSet            ExperimentalFeature = "OBJECT_PARAMETER_PARSE_ON_SET"
	WriteOnlyAttributes                  ExperimentalFeature = "WRITE_ONLY_ATTRIBUTES"
	FailoverGroupReadOnlyAutoSwitch      ExperimentalFeature = "FAILOVER_GROUP_READ_ONLY_AUTO_SWITCH"
)

type experimentalFeatureState string

const (
	ExperimentalFeatureStateOptIn            experimentalFeatureState = "OPT_IN"
	ExperimentalFeatureStateEnabledByDefault experimentalFeatureState = "ENABLED_BY_DEFAULT"
	ExperimentalFeatureStatePromoted         experimentalFeatureState = "PROMOTED"
	ExperimentalFeatureStateDiscontinued     experimentalFeatureState = "DISCONTINUED"
)

type Experiment struct {
	name             ExperimentalFeature
	state            experimentalFeatureState
	description      string
	removeInVersions []string
}

func (e *Experiment) Name() ExperimentalFeature {
	return e.name
}

func (e *Experiment) Description() string {
	if len(e.removeInVersions) == 0 {
		return e.description
	}
	sentence := fmt.Sprintf("It will be promoted in %s.", e.RemoveInVersionsPhrase())
	if e.description == "" {
		return sentence
	}
	return joinWithDoubleNewline(e.description, sentence)
}

func (e *Experiment) RemoveInVersionsPhrase() string {
	return joinWithOr(e.removeInVersions)
}

func NewOptInExperiment(name ExperimentalFeature, description ...string) Experiment {
	return Experiment{
		name:        name,
		state:       ExperimentalFeatureStateOptIn,
		description: joinWithDoubleNewline(description...),
	}
}

const enabledByDefaultOptOutDescription = "This experiment is enabled by default. Opt out by listing it in `experimental_features_disabled`."

func NewEnabledByDefaultExperiment(name ExperimentalFeature, removeInVersions []string, description ...string) Experiment {
	return Experiment{
		name:             name,
		state:            ExperimentalFeatureStateEnabledByDefault,
		description:      joinWithDoubleNewline(append(slices.Clone(description), enabledByDefaultOptOutDescription)...),
		removeInVersions: removeInVersions,
	}
}

func NewPromotedExperiment(name ExperimentalFeature) Experiment {
	return Experiment{
		name:  name,
		state: ExperimentalFeatureStatePromoted,
	}
}

func NewDiscontinuedExperiment(name ExperimentalFeature) Experiment {
	return Experiment{
		name:  name,
		state: ExperimentalFeatureStateDiscontinued,
	}
}

var allExperiments = []Experiment{
	NewPromotedExperiment(WarehouseShowImprovedPerformance),
	NewPromotedExperiment(GrantsStrictPrivilegeManagement),
	NewPromotedExperiment(ParametersIgnoreValueChangesIfNotOnObjectLevel),
	NewOptInExperiment(
		ParametersReducedOutput,
		"Currently, the `parameters` field in various resources contains a verbatim output for the `SHOW PARAMETERS IN <object>` command. One of the fields contained in the output is the `description`. It does not change and is repeated for all objects containing the given parameter. It leads to an excessive output (check e.g., [#3118](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3118)).",
		"To mitigate the problem, we are adding this option to reduce the output to only `value` and `level` fields, which should significantly reduce the state size. **Note**: it's also affecting the `parameters` output for data sources.",
		"We considered the option to remove the `parameters` output completely, however, we plan to change the external change logic detection to use it (to make it consistent with other attributes using `show_output` and because we won't be able to implement the current logic when switching to the Terraform Plugin Framework) and it still allows referencing the parameter value/level from other parts of the configuration.",
	),
	NewEnabledByDefaultExperiment(
		UserEnableDefaultWorkloadIdentity,
		[]string{"v2.24.0", "v2.25.0"},
		"The `default_workload_identity` field on the `snowflake_legacy_service_user` and `snowflake_service_user` resources manages workload identity federation (WIF).",
		"When the field is omitted, a workload identity configured outside Terraform is unset on the next apply.",
		"Read more in our [migration guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/dev/MIGRATION_GUIDE.md#new-feature-user_enable_default_workload_identity-enabled-by-default).",
	),
	NewPromotedExperiment(GrantsImportValidation),
	NewOptInExperiment(
		TagsAllowEmptyAllowedValues,
		"Enables behavior changes for the `allowed_values` field in the `snowflake_tag` resource.",
		"When enabled, the three possible states in Snowflake for allowed values will be supported: `nil` (any value is allowed; whenever `allowed_values` are empty), `empty` (no value is allowed; handled by the `no_allowed_values` field), and `set` (all values defined in `allowed_values` are allowed).",
		"Otherwise, the `no_allowed_values` field will be ignored (explicit changes will cause updates, but without any effect) and the `allowed_values` field will follow the old behavior: `nil` (any value is allowed; only available whenever tag resource is created without `allowed_values`), `empty` (no value is allowed; always set when updating from filled `allowed_values` set to empty one or completely removed from config), `set` (all values defined in `allowed_values` are allowed).",
	),
	NewOptInExperiment(
		ImportBooleanDefault,
		fmt.Sprintf("This experiment is deprecated in favor of `%s`. Resources are migrated to `%s` gradually; when all of them are migrated, this experiment will be discontinued.", ImportPermissiveMode, ImportPermissiveMode),
		"Changes import behavior for boolean fields using the special `\"default\"` value.",
		"When enabled, boolean fields using the special `\"default\"` value are set to `\"default\"` during import instead of the actual Snowflake value (e.g., `\"false\"`). This prevents unavoidable diffs on every plan after import.",
		"Note: this is supported on all stage resources (`snowflake_stage_external_s3`, `snowflake_stage_external_azure`, `snowflake_stage_external_gcs`, `snowflake_stage_external_s3_compatible`, and `snowflake_stage_internal`), stream resources (`snowflake_stream_on_table` and `snowflake_stream_on_view`), and storage integration resources (`snowflake_storage_integration_aws` and `snowflake_storage_integration_azure`).",
	),
	NewOptInExperiment(
		ImportPermissiveMode,
		"Switches imports to the permissive mode.",
		"By default, imports use the enforcing mode: the plan after import shows a change for every attribute that the configuration omits but that has a value in Snowflake. Applying that plan resets these attributes to their Snowflake defaults.",
		"When enabled, the provider accepts the values the object already has for the attributes that the configuration omits, so the plan after import shows no change for them. This applies only to the import itself: external changes made after the import are detected as usual.",
		"Note: this is currently supported on the `snowflake_authentication_policy`, `snowflake_backup_policy`, `snowflake_backup_set`, `snowflake_catalog_integration_aws_glue`, `snowflake_catalog_integration_iceberg_rest`, `snowflake_catalog_integration_object_storage`, `snowflake_catalog_integration_open_catalog`, `snowflake_compute_pool`, `snowflake_external_volume`, `snowflake_openflow_deployment_snowflake_managed`, `snowflake_pipe`, `snowflake_postgres_mirror`, `snowflake_share`, `snowflake_warehouse`, `snowflake_warehouse_adaptive`, and `snowflake_warehouse_interactive` resources.",
	),
	NewPromotedExperiment(GrantsSafeDestroy),
	NewPromotedExperiment(TagAssociationSafeDestroy),
	NewOptInExperiment(
		GrantAccountRoleShowCaching,
		"Enables per-plan in-memory caching of `SHOW GRANTS OF ROLE` results for the `snowflake_grant_account_role` resource.",
		"Without caching, every resource instance issues an independent `SHOW GRANTS OF ROLE <name>` call during Read. In configurations with many grants sharing the same role, this results in N identical round-trips returning the same full result set — only 1 is needed.",
		"When enabled, the first Read for a given role fetches and caches the result; subsequent Reads in the same plan reuse it. The cache is invalidated on Create and Delete so mutations within a single apply remain visible to subsequent Reads.",
		"Additionally, the trailing Read at the end of Create is skipped (this resource has no computed or server-default fields to populate), removing a redundant `SHOW GRANTS OF ROLE` call per grant during apply.",
		"Intended for large configurations (thousands of `snowflake_grant_account_role` resources) where plan and apply time is dominated by redundant `SHOW GRANTS OF ROLE` calls.",
	),
	NewOptInExperiment(
		AccountRoleShowCaching,
		"When enabled, the result of looking up an account role by identifier (`SHOW ROLES LIKE '<name>'`, via the underlying `ShowByID`/`ShowByIDSafely` calls) is cached in memory for the duration of a single plan or apply cycle.",
		"Currently supported by: `snowflake_account_role`, `snowflake_grant_application_role`, `snowflake_grant_privileges_to_account_role`.",
		"Without caching, every lookup of a given role — whether it's the role's own `snowflake_account_role` Read, or an existence check performed by a grant resource before granting to/from it — issues an independent round trip, even when many resource instances reference the same role. When enabled, the first lookup for a given role fetches and caches the result; subsequent lookups in the same plan reuse it. The cache is invalidated on `snowflake_account_role` Update (rename or comment change) and Delete, since only that resource can change what a cached lookup would return.",
		fmt.Sprintf("This is a separate flag from `%s`: enabling this does not enable caching for `snowflake_grant_account_role`'s `SHOW GRANTS OF ROLE` calls, and vice versa. Both can be enabled together.", GrantAccountRoleShowCaching),
		"Intended for large configurations (thousands of role or grant resources) where plan and apply time is dominated by redundant role lookups.",
	),
	NewOptInExperiment(
		GrantsShowCaching,
		"When enabled, `SHOW GRANTS` results are cached in memory for the duration of a single plan or apply cycle, so multiple resource instances resolving to the same underlying SHOW statement share one round-trip instead of each issuing their own.",
		"Currently supported by: `snowflake_grant_privileges_to_account_role`, `snowflake_grant_ownership`.",
		"Without caching, every resource instance issues an independent `SHOW GRANTS ON <object>` / `SHOW FUTURE GRANTS IN <container>` call during Read. In configurations with many grants resolving to the same underlying SHOW statement (e.g. many privilege grants on the same schema, or many future-grant roles on the same database), this results in N identical round-trips returning the same full result set — only 1 is needed.",
		"The first Read for a given SHOW statement fetches and caches the result; subsequent Reads in the same plan reuse it. The cache is invalidated on Create, Update, and Delete of the resources listed above so mutations within a single apply remain visible to subsequent Reads.",
		fmt.Sprintf("This is a separate flag from `%s`: enabling this does not enable caching for `snowflake_grant_account_role`, and vice versa. Both can be enabled together.", GrantAccountRoleShowCaching),
		"Intended for large configurations (thousands of grant resources) where plan and apply time is dominated by redundant `SHOW GRANTS` calls.",
	),
	NewPromotedExperiment(GrantAccountRoleSafePublicRole),
	NewOptInExperiment(
		HierarchyRenames,
		"When enabled, allows in-place handling of hierarchy renames and moves for supported resources.",
		"Currently supported by: `snowflake_authentication_policy`, `snowflake_backup_policy`, `snowflake_backup_set`, `snowflake_file_format_avro`, `snowflake_file_format_csv`, `snowflake_file_format_json`, `snowflake_file_format_orc`, `snowflake_file_format_parquet`, `snowflake_file_format_xml`, `snowflake_hybrid_table`, `snowflake_openflow_runtime`, `snowflake_password_policy`, `snowflake_projection_policy`, `snowflake_row_access_policy`, `snowflake_schema`, `snowflake_semantic_view`, `snowflake_session_policy`, `snowflake_stage_external_azure`, `snowflake_stage_external_gcs`, `snowflake_stage_external_s3`, `snowflake_stage_external_s3_compatible`, `snowflake_stage_internal`, `snowflake_storage_lifecycle_policy`, `snowflake_streamlit`, `snowflake_table`, `snowflake_tag`, `snowflake_view`.",
		"Without this experiment, changing the `database` field on `snowflake_schema` or the `database`/`schema` fields on the other supported resources may force resource recreation. With this experiment, the provider detects whether the parent was renamed or the object should be moved, and handles it without recreation.",
		"For more information, see the [object renaming guide](./guides/object_renaming_guide).",
	),
	NewPromotedExperiment(InheritedGrants),
	NewOptInExperiment(
		ObjectParameterUnsetOnDelete,
		"Changes the delete behavior of the `snowflake_object_parameter` resource to use `ALTER <OBJECT_TYPE> <identifier> UNSET <PARAMETER>` instead of resetting the parameter to its default value.",
		"Without this experiment, deleting the resource fetches the parameter's default value and explicitly sets it back, which is fragile and doesn't truly remove the object-level override.",
		"When enabled, the parameter is properly unset, allowing the inherited value from the higher hierarchy level to take effect.",
	),
	NewEnabledByDefaultExperiment(
		ObjectParameterParseOnSet,
		[]string{"v2.25.0", "v2.26.0"},
		"Changes how `snowflake_object_parameter` handles the `value` field for covered parameters.",
		"With this experiment, covered parameters are parsed with the same rules as the typed SDK fields. A value that does not match the parameter fails in the provider. String values are quoted by the provider, so a configuration such as `DEFAULT_DDL_COLLATION = \"en-ci-upper\"` is stored as written. Identifier parameters stay a bare name in configuration. Numeric and boolean values need no quoting either way.",
		"Covered today: account-level (`on_account = true`) parameters and Database, Schema, Warehouse, User, Task, Service, Hybrid Table, and Iceberg Table. A parameter or object type that is not covered yet keeps the previous behavior, including Function and Procedure, because this resource does not capture argument types. An enum value the catalog does not recognize is logged and also keeps the previous behavior. This experiment also decides if `snowflake_object_parameter` is forward-compatible: with the experiment enabled, the `value_type` field sends a parameter or an object type that is not covered yet, rendering `value` by that type. It is behind this experiment, because the rules for the value (the provider quotes it by the parameter type) follow the behavior of this experiment. Without the experiment, `value_type` is not used.",
		"Without this experiment, `value` is sent as a raw SQL token for parameters set on objects, and account-level parameters use the v2.22.x legacy write path. You supply any quoting, and a value that does not match the parameter is not checked by the provider.",
		"See the migration guide for a detailed view: https://github.com/snowflakedb/terraform-provider-snowflake/blob/dev/MIGRATION_GUIDE.md#new-feature-account_parameter_catalog_write_path-and-object_parameter_parse_on_set-enabled-by-default.",
	),
	NewOptInExperiment(
		AuthenticatorExplicitOnly,
		"Removes implicit authenticator derivation from other provider configuration fields.",
		"Without this experiment, the provider automatically sets the authenticator to `OAUTH` when the `token` or `token_accessor` field is configured, even if `authenticator` is not explicitly set. This implicit behavior can be confusing and will be removed in v3.",
		"When enabled, the `authenticator` field must be set explicitly in the provider configuration or TOML profile. The `SNOWFLAKE` default (when no authenticator is configured anywhere) is preserved.",
	),
	NewOptInExperiment(
		ProviderConfigurationAccountFallback,
		"Re-introduces the `account` field as a fallback for `organization_name` and `account_name` in both the provider configuration and TOML profiles.",
		"When enabled, you can set `account` instead of setting `organization_name` and `account_name` separately. The field accepts both the `org-name` format (e.g. `\"myorg-myaccount\"`) and an account locator (e.g. `\"xy12345\"`). If both `organization_name` and `account_name` are set, they take precedence over `account`. The `SNOWFLAKE_ACCOUNT` environment variable is used as the `account` value only when this experiment is enabled.",
		"Without this experiment, setting the `account` field in the provider configuration or in a TOML profile results in an error directing you to enable this experiment. A value coming from the `SNOWFLAKE_ACCOUNT` environment variable is ignored with a warning instead, because this experiment will be enabled by default in v3 and the variable will be taken into account from that version on.",
	),
	NewEnabledByDefaultExperiment(
		AccountParameterCatalogWritePath,
		[]string{"v2.25.0", "v2.26.0"},
		"Changes how `snowflake_account_parameter` builds the SQL it sends for the `value` field of string-typed and identifier-typed parameters.",
		"With this experiment, string-typed parameter values must not be manually wrapped in single quotes within `value`, because quoting is handled by the provider. Identifier-typed parameters (e.g. referencing a network policy by name) are sent as a quoted identifier, so `value` must match the name stored in Snowflake.",
		"Without this experiment, the v2.22.x legacy write path is used: some string-typed values must be wrapped in single quotes manually, identifier-typed values are sent as a single-quoted string literal or as written, and the parameters that the legacy write path does not support (e.g. the parameters added in v2.22.0, like `DISABLE_UI_DOWNLOAD_BUTTON`) fail with an error.",
		"This experiment also decides if `snowflake_account_parameter` is forward-compatible (new Snowflake parameters usable without a provider upgrade). With the experiment enabled, the `value_type` field sends a parameter that the provider does not support yet, rendering `value` by that type. It is behind this experiment, because the rules for the value (the provider quotes it by the parameter type) follow the behavior of this experiment. Without the experiment, `value_type` is not used. See the migration guide for a detailed, before/after view of the configuration changes: https://github.com/snowflakedb/terraform-provider-snowflake/blob/dev/MIGRATION_GUIDE.md#new-feature-account_parameter_catalog_write_path-and-object_parameter_parse_on_set-enabled-by-default.",
	),
	NewOptInExperiment(
		WriteOnlyAttributes,
		"Enables write-only attributes. Their values are never saved in the Terraform plan or state, so they can be set from an ephemeral resource. They require Terraform 1.11 or later.",
		"Currently supported by: `rsa_public_key_wo` and `rsa_public_key_2_wo`, with their `rsa_public_key_wo_version` and `rsa_public_key_2_wo_version` version counters, on `snowflake_user`, `snowflake_service_user`, and `snowflake_legacy_service_user`.",
		"Without this experiment, using a write-only attribute that is behind it results in an error.",
	),
	NewOptInExperiment(
		FailoverGroupReadOnlyAutoSwitch,
		"Lets the `snowflake_failover_group` resource change the membership of a failover group that stays in the read-only mode.",
		"Snowflake rejects changes to `object_types`, `allowed_databases`, `allowed_external_volumes`, `allowed_shares` and `allowed_integration_types` while a failover group is read-only. Without this experiment, the plan fails when `read_only` stays `true` and one of these attributes changes.",
		"When enabled, the plan passes. The provider leaves the read-only mode before the change and enters it again afterwards, in the same apply.",
	),
}

func joinWithDoubleNewline(parts ...string) string {
	return strings.Join(parts, "\n\n")
}

func joinWithOr(parts []string) string {
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0]
	case 2:
		return parts[0] + " or " + parts[1]
	default:
		return strings.Join(parts[:len(parts)-1], ", ") + ", or " + parts[len(parts)-1]
	}
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
	OptInExperiments              = collections.Filter(allExperiments, filterByStateFuncProvider(ExperimentalFeatureStateOptIn))
	EnabledByDefaultExperiments   = collections.Filter(allExperiments, filterByStateFuncProvider(ExperimentalFeatureStateEnabledByDefault))
	PromotedExperiments           = collections.Filter(allExperiments, filterByStateFuncProvider(ExperimentalFeatureStatePromoted))
	DiscontinuedExperiments       = collections.Filter(allExperiments, filterByStateFuncProvider(ExperimentalFeatureStateDiscontinued))
	acceptedInDisabledExperiments = allExperiments
)

var (
	allExperimentalFeatureNames                = collections.Map(allExperiments, mapToName)
	optInExperimentalFeatureNames              = collections.Map(OptInExperiments, mapToName)
	enabledByDefaultExperimentalFeatureNames   = collections.Map(EnabledByDefaultExperiments, mapToName)
	acceptedInDisabledExperimentalFeatureNames = collections.Map(acceptedInDisabledExperiments, mapToName)
)

var (
	AllExperimentalFeatureNames                = sdk.AsStringList(allExperimentalFeatureNames)
	OptInExperimentalFeatureNames              = sdk.AsStringList(optInExperimentalFeatureNames)
	EnabledByDefaultExperimentalFeatureNames   = sdk.AsStringList(enabledByDefaultExperimentalFeatureNames)
	AcceptedInDisabledExperimentalFeatureNames = sdk.AsStringList(acceptedInDisabledExperimentalFeatureNames)
)

// Experiments contains the configured experiment set.
// Construct with New; do not assemble it by hand.
type Experiments struct {
	userEnabled  []string
	userDisabled []string
	effective    []string
}

// New builds the experiment set from the user-configured enabled and disabled lists.
func New(userEnabled, userDisabled []string) Experiments {
	return Experiments{
		userEnabled:  userEnabled,
		userDisabled: userDisabled,
		effective:    ComputeEffectiveExperimentsEnabled(userEnabled, userDisabled),
	}
}

// IsEnabled reports whether the given experiment is on for this provider instance.
func (e Experiments) IsEnabled(feature ExperimentalFeature) bool {
	return isListed(feature, e.effective)
}

// UserEnabled returns a copy of the user-configured experimental_features_enabled list.
func (e Experiments) UserEnabled() []string {
	return slices.Clone(e.userEnabled)
}

// UserDisabled returns a copy of the user-configured experimental_features_disabled list.
func (e Experiments) UserDisabled() []string {
	return slices.Clone(e.userDisabled)
}

// EffectiveEnabled returns a copy of the computed effective enabled set.
func (e Experiments) EffectiveEnabled() []string {
	return slices.Clone(e.effective)
}

func isListed(feature ExperimentalFeature, names []string) bool {
	return slices.ContainsFunc(names, func(s string) bool {
		return strings.EqualFold(string(feature), s)
	})
}

func GetExperiment(name ExperimentalFeature) (Experiment, bool) {
	return experimentByName(string(name))
}

func experimentByName(name string) (Experiment, bool) {
	for _, experiment := range allExperiments {
		if strings.EqualFold(string(experiment.name), name) {
			return experiment, true
		}
	}
	return Experiment{}, false
}

func ComputeEffectiveExperimentsEnabled(userEnabled, userDisabled []string) []string {
	return computeEffectiveExperimentsEnabled(allExperiments, userEnabled, userDisabled)
}

func computeEffectiveExperimentsEnabled(experiments []Experiment, userEnabled, userDisabled []string) []string {
	effective := make([]string, 0)
	for _, experiment := range experiments {
		switch experiment.state {
		case ExperimentalFeatureStateOptIn:
			if isListed(experiment.name, userEnabled) && !isListed(experiment.name, userDisabled) {
				effective = append(effective, string(experiment.name))
			}
		case ExperimentalFeatureStateEnabledByDefault:
			if !isListed(experiment.name, userDisabled) {
				effective = append(effective, string(experiment.name))
			}
		}
	}
	return effective
}
