//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Procedures_gh3822_bcr2025_03(t *testing.T) {
	// TODO(SNOW-2196333): Resolve these tests after the change rollout is clarified.
	t.Skip("Skipping because the changes have been reverted from the BCR")

	schema, schemaCleanup := secondaryTestClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	_, procedureCleanup := secondaryTestClient().Procedure.CreatePythonInSchema(t, schema.ID())
	t.Cleanup(procedureCleanup)

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary).
		WithPreviewFeaturesEnabled(string(previewfeatures.ProceduresDatasource))
	proceduresModel := datasourcemodel.Procedures("test", schema.ID().DatabaseName(), schema.ID().Name())

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
				Config: config.FromModels(t, providerModel, proceduresModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(proceduresModel.DatasourceReference(), "database", schema.ID().DatabaseName()),
					resource.TestCheckResourceAttr(proceduresModel.DatasourceReference(), "schema", schema.ID().Name()),
					resource.TestCheckResourceAttrSet(proceduresModel.DatasourceReference(), "procedures.#"),
				),
			},
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.2.0"),
				PreConfig: func() {
					secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2025_03")
				},
				Config:      config.FromModels(t, providerModel, proceduresModel),
				ExpectError: regexp.MustCompile("could not parse arguments"),
			},
			{
				ProtoV6ProviderFactories: functionsAndProceduresProviderFactory,
				Config:                   config.FromModels(t, providerModel, proceduresModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(proceduresModel.DatasourceReference(), "database", schema.ID().DatabaseName()),
					resource.TestCheckResourceAttr(proceduresModel.DatasourceReference(), "schema", schema.ID().Name()),
					resource.TestCheckResourceAttrSet(proceduresModel.DatasourceReference(), "procedures.#"),
				),
			},
		},
	})
}
