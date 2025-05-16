//go:build !account_level_tests

package datasources_test

import (
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Databases_Complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseName := databaseId.Name()
	comment := random.Comment()
	secondaryAccountId := acc.SecondaryTestClient().Account.GetAccountIdentifier(t)

	databaseModel := model.DatabaseWithParametersSet("test", databaseName).
		WithComment(comment).
		WithReplication(secondaryAccountId, true, true)
	databasesModel := datasourcemodel.Databases("test").
		WithLike(databaseName).
		WithStartsWith(databaseName).
		WithLimit(1).
		WithDependsOn(databaseModel.ResourceReference())
	databasesWithoutOptionalsModel := datasourcemodel.Databases("test").
		WithLike(databaseName).
		WithStartsWith(databaseName).
		WithLimit(1).
		WithDependsOn(databaseModel.ResourceReference()).
		WithWithDescribe(false).
		WithWithParameters(false)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, databaseModel, databasesModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.#", "1"),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.name", databaseName),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.kind", "STANDARD"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.transient", "false"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.is_current", "true"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.origin", ""),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.show_output.0.owner"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.options", ""),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.show_output.0.retention_time"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.resource_group", ""),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.show_output.0.owner_role_type"),

					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.describe_output.#", "2"),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.describe_output.0.created_on"),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.describe_output.0.name"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.describe_output.0.kind", "SCHEMA"),

					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.parameters.#", "1"),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.parameters.0.data_retention_time_in_days.0.value"),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.parameters.0.max_data_extension_time_in_days.0.value"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.parameters.0.external_volume.0.value", ""),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.parameters.0.catalog.0.value"),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.parameters.0.replace_invalid_characters.0.value"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.parameters.0.default_ddl_collation.0.value", ""),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.parameters.0.storage_serialization_policy.0.value"),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.parameters.0.log_level.0.value"),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.parameters.0.trace_level.0.value"),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.parameters.0.suspend_task_after_num_failures.0.value"),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.parameters.0.task_auto_retry_attempts.0.value"),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.parameters.0.user_task_managed_initial_warehouse_size.0.value"),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.parameters.0.user_task_minimum_trigger_interval_in_seconds.0.value"),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.parameters.0.quoted_identifiers_ignore_case.0.value"),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.parameters.0.enable_console_output.0.value"),
				),
			},
			{
				Config: accconfig.FromModels(t, databaseModel, databasesWithoutOptionalsModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.#", "1"),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.name", databaseName),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.kind", "STANDARD"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.transient", "false"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.is_current", "true"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.origin", ""),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.show_output.0.owner"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.options", ""),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.show_output.0.retention_time"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.show_output.0.resource_group", ""),
					resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.show_output.0.owner_role_type"),

					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.describe_output.#", "0"),
					resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.parameters.#", "0"),
				),
			},
		},
	})
}

func TestAcc_Databases_DifferentFiltering(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	prefix := random.AlphaN(4)
	idOne := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	databaseModel1 := model.DatabaseWithParametersSet("test", idOne.Name())
	databaseModel2 := model.DatabaseWithParametersSet("test1", idTwo.Name())
	databaseModel3 := model.DatabaseWithParametersSet("test2", idThree.Name())
	databasesWithLikeModel := datasourcemodel.Databases("test").
		WithLike(idOne.Name()).
		WithDependsOn(databaseModel1.ResourceReference(), databaseModel2.ResourceReference(), databaseModel3.ResourceReference())
	databasesWithStartsWithModel := datasourcemodel.Databases("test").
		WithStartsWith(prefix).
		WithDependsOn(databaseModel1.ResourceReference(), databaseModel2.ResourceReference(), databaseModel3.ResourceReference())
	databasesWithLimitModel := datasourcemodel.Databases("test").
		WithRowsAndFrom(1, prefix).
		WithDependsOn(databaseModel1.ResourceReference(), databaseModel2.ResourceReference(), databaseModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, databaseModel1, databaseModel2, databaseModel3, databasesWithLikeModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(databasesWithLikeModel.DatasourceReference(), "databases.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, databaseModel1, databaseModel2, databaseModel3, databasesWithStartsWithModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(databasesWithLikeModel.DatasourceReference(), "databases.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, databaseModel1, databaseModel2, databaseModel3, databasesWithLimitModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(databasesWithLikeModel.DatasourceReference(), "databases.#", "1"),
				),
			},
		},
	})
}

func TestAcc_Databases_DatabaseNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      databasesWithPostcondition(),
				ExpectError: regexp.MustCompile("there should be at least one database"),
			},
		},
	})
}

func databasesWithPostcondition() string {
	return `
data "snowflake_databases" "test" {
  like = "non-existing-database"

  lifecycle {
    postcondition {
      condition     = length(self.databases) > 0
      error_message = "there should be at least one database"
    }
  }
}
`
}
