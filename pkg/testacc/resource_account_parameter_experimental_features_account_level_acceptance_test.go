//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Experimental_AccountParameter_CatalogWritePath(t *testing.T) {
	stringParam := sdk.AccountParameterPythonProfilerModules
	identifierParam := sdk.AccountParameterPythonProfilerTargetStage
	floatParam := sdk.AccountParameterInitialReplicationSizeLimitInTb

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	t.Cleanup(func() {
		testClient().Parameter.UnsetAccountParameter(t, stringParam)
		testClient().Parameter.UnsetAccountParameter(t, identifierParam)
		testClient().Parameter.UnsetAccountParameter(t, floatParam)
	})

	stageName := stage.ID().FullyQualifiedName()
	legacyStringModel := model.AccountParameter("string_type", string(stringParam), "'module_a,module_b'")
	identifierModel := model.AccountParameter("identifier_type", string(identifierParam), stageName)
	floatModel := model.AccountParameter("float_type", string(floatParam), "3.0")

	catalogStringModel := model.AccountParameter("string_type", string(stringParam), "module_a,module_b")

	experimentsDisabledProviderModel := providermodel.SnowflakeProvider().
		WithAllEnabledByDefaultExperimentsDisabled()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			CheckAccountParameterUnset(t, stringParam),
			CheckAccountParameterUnset(t, identifierParam),
			CheckAccountParameterUnset(t, floatParam),
		),
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: enabledByDefaultExperimentsDisabledProviderFactory,
				Config:                   config.FromModels(t, experimentsDisabledProviderModel, legacyStringModel, identifierModel, floatModel),
				ExpectNonEmptyPlan:       true,
				// The legacy write path requires the quotes as part of the HCL value, but Snowflake reports back the
				// unquoted value, so refresh always diffs config vs. state for the string-typed parameter and plans an
				// update (a no-op in Snowflake) — exactly the quoting quirk the experiment fixes.
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(legacyStringModel.ResourceReference(), plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(legacyStringModel.ResourceReference(), tfjsonpath.New("value"), knownvalue.StringExact("'module_a,module_b'")),
						plancheck.ExpectResourceAction(identifierModel.ResourceReference(), plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction(floatModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(
					t,
					resourceassert.AccountParameterResource(t, legacyStringModel.ResourceReference()).
						HasKeyString(string(stringParam)).
						HasValueString("module_a,module_b"),
					resourceassert.AccountParameterResource(t, identifierModel.ResourceReference()).
						HasKeyString(string(identifierParam)).
						HasValueString(stageName),
					resourceassert.AccountParameterResource(t, floatModel.ResourceReference()).
						HasKeyString(string(floatParam)).
						HasValueString("3.0"),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, catalogStringModel, identifierModel, floatModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(
					t,
					resourceassert.AccountParameterResource(t, catalogStringModel.ResourceReference()).
						HasKeyString(string(stringParam)).
						HasValueString("module_a,module_b"),
					resourceassert.AccountParameterResource(t, identifierModel.ResourceReference()).
						HasKeyString(string(identifierParam)).
						HasValueString(stageName),
					resourceassert.AccountParameterResource(t, floatModel.ResourceReference()).
						HasKeyString(string(floatParam)).
						HasValueString("3.0"),
				),
			},
		},
	})
}

func TestAcc_Experimental_AccountParameter_CatalogWritePath_NewParameters(t *testing.T) {
	param := sdk.AccountParameterDisableUiDownloadButton
	t.Cleanup(func() { testClient().Parameter.UnsetAccountParameter(t, param) })

	parameterModel := model.AccountParameter("new_parameter", string(param), "true")
	experimentsDisabledProviderModel := providermodel.SnowflakeProvider().
		WithAllEnabledByDefaultExperimentsDisabled()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountParameterUnset(t, param),
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: enabledByDefaultExperimentsDisabledProviderFactory,
				Config:                   config.FromModels(t, experimentsDisabledProviderModel, parameterModel),
				ExpectError:              regexp.MustCompile(`parameter not supported by the legacy account parameter write path: DISABLE_UI_DOWNLOAD_BUTTON`),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, parameterModel),
				Check: assertThat(
					t,
					resourceassert.AccountParameterResource(t, parameterModel.ResourceReference()).
						HasKeyString(string(param)).
						HasValueString("true"),
				),
			},
		},
	})
}
