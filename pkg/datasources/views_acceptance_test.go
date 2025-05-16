//go:build !account_level_tests

package datasources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO(SNOW-1423486): Fix using warehouse in all tests and remove unsetting testenvs.ConfigureClientOnce.
func TestAcc_Views(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	schema, schemaCleanup := acc.TestClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	prefix := random.AlphaN(6)
	viewId := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix + "1")
	viewId2 := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix+"2", schema.ID())
	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"

	viewModel1 := model.View("v1", viewId.DatabaseName(), viewId.SchemaName(), viewId.Name(), statement)
	viewModel2 := model.View("v2", viewId2.DatabaseName(), viewId2.SchemaName(), viewId2.Name(), statement)
	viewsModelInSchema := datasourcemodel.Views("in_schema").
		WithInSchema(schema.ID()).
		WithDependsOn(viewModel1.ResourceReference(), viewModel2.ResourceReference())
	viewsModelFiltering := datasourcemodel.Views("filtering").
		WithLike(prefix+"%").
		WithStartsWith(prefix).
		WithInDatabase(viewId.DatabaseId()).
		WithLimitRowsAndFrom(1, viewId.Name()).
		WithDependsOn(viewModel1.ResourceReference(), viewModel2.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, viewModel1, viewModel2, viewsModelInSchema),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.#", "1"),

					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.show_output.0.name", viewId2.Name()),
					resource.TestCheckResourceAttrSet(viewsModelInSchema.DatasourceReference(), "views.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.show_output.0.kind", ""),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.show_output.0.reserved", ""),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.show_output.0.database_name", schema.ID().DatabaseName()),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.show_output.0.schema_name", schema.ID().Name()),
					resource.TestCheckResourceAttrSet(viewsModelInSchema.DatasourceReference(), "views.0.show_output.0.owner"),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet(viewsModelInSchema.DatasourceReference(), "views.0.show_output.0.text"),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.show_output.0.is_secure", "false"),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.show_output.0.is_materialized", "false"),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.show_output.0.owner_role_type", "ROLE"),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.show_output.0.change_tracking", "OFF"),

					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.describe_output.#", "2"),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.describe_output.0.name", "ROLE_NAME"),
					resource.TestCheckResourceAttrSet(viewsModelInSchema.DatasourceReference(), "views.0.describe_output.0.type"),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.describe_output.0.kind", "COLUMN"),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.describe_output.0.is_nullable", "true"),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.describe_output.0.default", ""),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.describe_output.0.is_primary", "false"),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.describe_output.0.is_unique", "false"),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.describe_output.0.check", ""),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.describe_output.0.expression", ""),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.describe_output.0.comment", ""),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.describe_output.0.policy_name", ""),
					resource.TestCheckNoResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.describe_output.0.policy_domain"),
					resource.TestCheckResourceAttr(viewsModelInSchema.DatasourceReference(), "views.0.describe_output.1.name", "ROLE_OWNER"),
				),
			},
			{
				Config: accconfig.FromModels(t, viewModel1, viewModel2, viewsModelFiltering),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(viewsModelFiltering.DatasourceReference(), "views.#", "1"),
					resource.TestCheckResourceAttr(viewsModelFiltering.DatasourceReference(), "views.0.show_output.0.name", viewId2.Name()),
				),
			},
		},
	})
}
