//go:build account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
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

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.AccountParameterCatalogWritePath)

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
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, legacyStringModel, identifierModel, floatModel),
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
				ProtoV6ProviderFactories: accountParameterCatalogWritePathProviderFactory,
				Config:                   config.FromModels(t, experimentProviderModel, catalogStringModel, identifierModel, floatModel),
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
