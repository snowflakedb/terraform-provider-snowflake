//go:build account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// Use only parameters that can be set only on the account level for the time-being.
// TODO [SNOW-1866453]: add more acc tests for the remaining parameters

func TestAcc_AccountParameter(t *testing.T) {
	testCases := []struct {
		param        sdk.AccountParameter
		value        string
		defaultLevel sdk.ParameterType
	}{
		{sdk.AccountParameterAllowIDToken, "true", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterPreventLoadFromInlineURL, "true", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterRequireStorageIntegrationForStageCreation, "true", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterInitialReplicationSizeLimitInTB, "3.0", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterCsvTimestampFormat, "YYYY-MM-DD", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterDisableUserPrivilegeGrants, resources.BooleanTrue, sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterAllowBindValuesAccess, "true", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterAllowedSpcsWorkloadTypes, "ALL", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterDataMetricSchedule, "60 MINUTES", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterDefaultDbtVersion, "1.9.4", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterDisallowedSpcsWorkloadTypes, "", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterEnableBudgetEventLogging, "true", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterEnableDataCompaction, "true", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterEnableGetDdlUseDataTypeAlias, "false", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterEnableIcebergMergeOnRead, "true", sdk.ParameterTypeSnowflakeDefault},
		{sdk.AccountParameterEnableNotebookCreationInPersonalDb, "false", sdk.ParameterTypeSystem},
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
						Check: assertThat(t, resourceassert.AccountParameterResource(t, accountParameterModel.ResourceReference()).
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
	accountParameterModel := model.AccountParameter("test", string(sdk.AccountParameterInitialReplicationSizeLimitInTB), "3.0")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountParameterUnset(t, sdk.AccountParameterTimezone),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, accountParameterModel),
				Check: assertThat(t, resourceassert.AccountParameterResource(t, accountParameterModel.ResourceReference()).
					HasKeyString(string(sdk.AccountParameterInitialReplicationSizeLimitInTB)).
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
	accountParameterModel := model.AccountParameter("test", string(sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList), "true")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountParameterUnset(t, sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, accountParameterModel),
				Check: assertThat(t, resourceassert.AccountParameterResource(t, accountParameterModel.ResourceReference()).
					HasKeyString(string(sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList)).
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
