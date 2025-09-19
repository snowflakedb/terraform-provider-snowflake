//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Functions_gh3822_bcr2025_03(t *testing.T) {
	// TODO(SNOW-2196333): Resolve these tests after the change rollout is clarified.
	t.Skip("Skipping because the changes have been reverted from the BCR")

	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	schema, schemaCleanup := secondaryTestClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	_, functionCleanup := secondaryTestClient().Function.CreatePythonInSchema(t, schema.ID())
	t.Cleanup(functionCleanup)

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary).
		WithPreviewFeaturesEnabled(string(previewfeatures.FunctionsDatasource))
	functionsModel := datasourcemodel.Functions("test", schema.ID().DatabaseName(), schema.ID().Name())

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
				Config: config.FromModels(t, providerModel, functionsModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(functionsModel.DatasourceReference(), "database", schema.ID().DatabaseName()),
					resource.TestCheckResourceAttr(functionsModel.DatasourceReference(), "schema", schema.ID().Name()),
					resource.TestCheckResourceAttr(functionsModel.DatasourceReference(), "functions.#", "1"),
				),
			},
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.2.0"),
				PreConfig: func() {
					secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2025_03")
				},
				Config:      config.FromModels(t, providerModel, functionsModel),
				ExpectError: regexp.MustCompile("could not parse arguments"),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, providerModel, functionsModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(functionsModel.DatasourceReference(), "database", schema.ID().DatabaseName()),
					resource.TestCheckResourceAttr(functionsModel.DatasourceReference(), "schema", schema.ID().Name()),
					resource.TestCheckResourceAttr(functionsModel.DatasourceReference(), "functions.#", "1"),
				),
			},
		},
	})
}
