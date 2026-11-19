//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// Use only parameters that can be set only on the account level for the time-being.
// TODO [SNOW-1866453]: add more acc tests for the remaining parameters

func TestAcc_AccountParameter_BasicUseCase(t *testing.T) {
	testCases := []struct {
		param        sdk.AccountParameter
		value        string
		defaultLevel sdk.ParameterType
	}{
		{sdk.AccountParameterAllowIdToken, "true", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterPreventLoadFromInlineUrl, "true", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterRequireStorageIntegrationForStageCreation, "true", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterInitialReplicationSizeLimitInTb, "3.0", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterCsvTimestampFormat, "YYYY-MM-DD", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterDisableUserPrivilegeGrants, resources.BooleanTrue, sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterAllowBindValuesAccess, "true", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterAllowedSpcsWorkloadTypes, "ALL", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterDataMetricSchedule, "60 MINUTES", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterDefaultDbtVersion, "1.9.4", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterDefaultStreamlitComputePool, "SYSTEM_COMPUTE_POOL_GPU", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterDisallowedSpcsWorkloadTypes, "", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterEnableBudgetEventLogging, "true", sdk.ParameterTypeSnowflakeDefault},
		// A StringAllowEmpty parameter must accept an empty string as a real value.
		{sdk.AccountParameterDefaultDdlCollation, "", sdk.ParameterTypeSnowflakeDefault},
		// TOOD(SNOW-3953840): Unskip or remove {sdk.AccountParameterCortexModelsAllowlist, "All", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterCortexCodeCliDailyEstCreditLimitPerUser, "10", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterCortexCodeDesktopDailyEstCreditLimitPerUser, "20", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterCortexCodeSnowsightDailyEstCreditLimitPerUser, "30", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterEnableCortexAnalyst, "true", testClient().SnowflakeDefaults.DefaultEnableCortexAnalystLevel(t)},
		{sdk.AccountParameterEnableDataCompaction, "true", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterEnablePerAccountAppServicePrivatelinkUrl, "true", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterEnableGetDdlUseDataTypeAlias, "false", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterEnableIcebergMergeOnRead, "true", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterEnableNotebookCreationInPersonalDb, "false", testClient().SnowflakeDefaults.DefaultEnableNotebookCreationInPersonalDbLevel(t)},
		{sdk.AccountParameterEnableSpcsBlockStorageSnowflakeFullEncryptionEnforcement, "false", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterEnableTagPropagationEventLogging, "false", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterIcebergVersionDefault, "2", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterReadConsistencyMode, "SESSION", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterRowTimestampDefault, "false", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterSqlTraceQueryText, "OFF", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterUseWorkspacesForSql, "false", sdk.ParameterTypeSnowflakeDefault},
	}

	for _, tt := range testCases {
		t.Run(string(tt.param), func(t *testing.T) {
			accountParameterModel := model.AccountParameter("test", string(tt.param), tt.value)

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: CheckAccountParameterUnsetToDefaultLevel(t, tt.param, tt.defaultLevel),
				Steps: []resource.TestStep{
					{
						Config: config.FromModels(t, accountParameterModel),
						Check: assertThat(
							t, resourceassert.AccountParameterResource(t, accountParameterModel.ResourceReference()).
								HasKeyString(string(tt.param)).
								HasValueString(tt.value),
						),
					},
				},
			})
		})
	}
}

// Proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2573 is solved.
// Instead of TIMEZONE, we used INITIAL_REPLICATION_SIZE_LIMIT_IN_TB which is only settable on account so it does not mess with other tests.
func TestAcc_AccountParameter_Issue2573(t *testing.T) {
	accountParameterModel := model.AccountParameter("test", string(sdk.AccountParameterInitialReplicationSizeLimitInTb), "3.0")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountParameterUnset(t, sdk.AccountParameterTimezone),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, accountParameterModel),
				Check: assertThat(
					t, resourceassert.AccountParameterResource(t, accountParameterModel.ResourceReference()).
						HasKeyString(string(sdk.AccountParameterInitialReplicationSizeLimitInTb)).
						HasValueString("3.0"),
				),
			},
			{
				ResourceName:            "snowflake_account_parameter.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAcc_AccountParameter_Issue3025(t *testing.T) {
	accountParameterModel := model.AccountParameter("test", string(sdk.AccountParameterOauthAddPrivilegedRolesToBlockedList), "true")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountParameterUnset(t, sdk.AccountParameterOauthAddPrivilegedRolesToBlockedList),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, accountParameterModel),
				Check: assertThat(
					t, resourceassert.AccountParameterResource(t, accountParameterModel.ResourceReference()).
						HasKeyString(string(sdk.AccountParameterOauthAddPrivilegedRolesToBlockedList)).
						HasValueString("true"),
				),
			},
			{
				ResourceName:            "snowflake_account_parameter.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// The parameter catalog supports SERVICE_CALLER_TOKEN_VALIDITY_SECS only on services, so on the account level it can be set only through the generic write path.
func TestAcc_AccountParameter_GenericWritePath_UnknownParameter(t *testing.T) {
	param := sdk.AccountParameter("SERVICE_CALLER_TOKEN_VALIDITY_SECS")
	t.Cleanup(func() { testClient().Parameter.UnsetAccountParameter(t, param) })

	withoutType := model.AccountParameter("unknown", string(param), "120")
	withType := model.AccountParameter("unknown", string(param), "120").WithValueType("integer")
	withTypeInOtherCase := model.AccountParameter("unknown", string(param), "120").WithValueType("INTEGER")
	changedValue := model.AccountParameter("unknown", string(param), "240").WithValueType("INTEGER")
	invalidValue := model.AccountParameter("unknown", string(param), "abc").WithValueType("INTEGER")
	experimentsDisabledProviderModel := providermodel.SnowflakeProvider().
		WithAllEnabledByDefaultExperimentsDisabled()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountParameterUnset(t, param),
		Steps: []resource.TestStep{
			// The generic path needs the experiment.
			{
				ProtoV6ProviderFactories: enabledByDefaultExperimentsDisabledProviderFactory,
				Config:                   config.FromModels(t, experimentsDisabledProviderModel, withType),
				ExpectError:              regexp.MustCompile("remove ACCOUNT_PARAMETER_CATALOG_WRITE_PATH from `experimental_features_disabled`"),
			},
			// Without value_type, the provider cannot type the value.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, withoutType),
				ExpectError:              regexp.MustCompile("set the `value_type` field"),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, invalidValue),
				ExpectError:              regexp.MustCompile("integer value expected"),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, withType),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: assertThat(t,
					resourceassert.AccountParameterResource(t, withType.ResourceReference()).
						HasKeyString(string(param)).
						HasValueString("120").
						HasValueTypeString("integer"),
				),
			},
			// value_type is case-insensitive, so changing only its case shows no difference.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, withTypeInOtherCase),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, changedValue),
				Check: assertThat(t,
					resourceassert.AccountParameterResource(t, changedValue.ResourceReference()).
						HasValueString("240"),
				),
			},
		},
	})
}

// When a provider version starts to support a parameter, value_type stops being used, and the configuration does not have
// to change: with or without value_type, nothing is planned. INITIAL_REPLICATION_SIZE_LIMIT_IN_TB is in the catalog.
func TestAcc_AccountParameter_GenericWritePath_ValueTypeIsIgnoredForSupportedParameter(t *testing.T) {
	param := sdk.AccountParameterInitialReplicationSizeLimitInTb
	t.Cleanup(func() { testClient().Parameter.UnsetAccountParameter(t, param) })

	withType := model.AccountParameter("supported", string(param), "3.0").WithValueType("STRING")
	withOtherType := model.AccountParameter("supported", string(param), "3.0").WithValueType("INTEGER")
	withoutType := model.AccountParameter("supported", string(param), "3.0")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountParameterUnset(t, param),
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, withType),
				Check: assertThat(t,
					resourceassert.AccountParameterResource(t, withType.ResourceReference()).
						HasValueString("3.0"),
				),
			},
			// value_type is ignored for a supported parameter, so changing it to another type shows no difference.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, withOtherType),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply:             []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, withoutType),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply:             []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}
