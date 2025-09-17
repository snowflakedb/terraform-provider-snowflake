//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Function_gh3823_bcr2025_03_proof(t *testing.T) {
	// TODO(SNOW-2251351): Resolve these tests after the change rollout is clarified.
	t.Skip("Skipping because the changes have been reverted from the BCR")

	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	schema, schemaCleanup := secondaryTestClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	funcName := "some_function"
	argName := "x"
	dataType := testdatatypes.DataTypeNumber_36_2

	id := secondaryTestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsInSchemaNewDataTypes(schema.ID(), dataType)
	definition := secondaryTestClient().Function.SamplePythonDefinition(t, funcName, argName)

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary).
		WithPreviewFeaturesEnabled(string(previewfeatures.FunctionPythonResource))
	functionModel := model.FunctionPythonBasicInline("test", id, testvars.PythonRuntime, dataType, funcName, definition).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.2.0"),
				PreConfig: func() {
					secondaryTestClient().BcrBundles.DisableBcrBundle(t, "2025_03")
				},
				Config: config.FromModels(t, providerModel, functionModel),
				Check: assertThat(t,
					resourceassert.FunctionPythonResource(t, functionModel.ResourceReference()).HasNameString(id.Name()),
				),
			},
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.2.0"),
				PreConfig: func() {
					secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2025_03")
				},
				Config: config.FromModels(t, providerModel, functionModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(functionModel.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				// This error is pretty unclear, so these are the reasons why it happens:
				// - function is created without the bundle, bundle is enabled after
				// - when terraform tries to read the state of this function:
				//   - function output from SHOW cannot be parsed
				//   - `convert` function does not return error
				//   - `NewSchemaObjectIdentifierWithArguments` function does not return error
				//   - comparison in ShowByID method does not find a match (as the arguments list is empty)
				//   - ErrObjectNotFound is returned
				//   - terraform marks the object as removed from state
				//   - terraform tries to create it again (but it exists in Snowflake already)
				//   - the compilation error is thrown by Snowflake (as object already exists)
				ExpectError: regexp.MustCompile("SQL compilation error"),
			},
		},
	})
}

func TestAcc_Function_gh3823_bcr2025_03_fix(t *testing.T) {
	// TODO(SNOW-2251351): Resolve these tests after the change rollout is clarified.
	t.Skip("Skipping because the changes have been reverted from the BCR")

	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	schema, schemaCleanup := secondaryTestClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	funcName := "some_function"
	argName := "x"
	dataType := testdatatypes.DataTypeNumber_36_2

	id := secondaryTestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsInSchemaNewDataTypes(schema.ID(), dataType)
	definition := secondaryTestClient().Function.SamplePythonDefinition(t, funcName, argName)

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary).
		WithPreviewFeaturesEnabled(string(previewfeatures.FunctionPythonResource))
	functionModel := model.FunctionPythonBasicInline("test", id, testvars.PythonRuntime, dataType, funcName, definition).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.2.0"),
				PreConfig: func() {
					secondaryTestClient().BcrBundles.DisableBcrBundle(t, "2025_03")
				},
				Config: config.FromModels(t, providerModel, functionModel),
				Check: assertThat(t,
					resourceassert.FunctionPythonResource(t, functionModel.ResourceReference()).HasNameString(id.Name()),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig: func() {
					secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2025_03")
				},
				Config: config.FromModels(t, providerModel, functionModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: assertThat(t,
					resourceassert.FunctionPythonResource(t, functionModel.ResourceReference()).HasNameString(id.Name()),
				),
			},
		},
	})
}
