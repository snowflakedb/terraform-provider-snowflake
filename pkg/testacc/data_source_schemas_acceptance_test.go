//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Schemas_BasicUseCase_DifferentFiltering(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	prefix := strings.ToUpper(random.AlphaN(4))

	idOne := testClient().Ids.RandomDatabaseObjectIdentifierWithPrefix(prefix + "1")
	idTwo := testClient().Ids.RandomDatabaseObjectIdentifierWithPrefix(prefix + "2")
	idThree := testClient().Ids.RandomDatabaseObjectIdentifier()
	idFour := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(database.ID())

	schemaModel1 := model.Schema("test_1", idOne.DatabaseName(), idOne.Name())
	schemaModel2 := model.Schema("test_2", idTwo.DatabaseName(), idTwo.Name())
	schemaModel3 := model.Schema("test_3", idThree.DatabaseName(), idThree.Name())
	schemaModel4 := model.Schema("test_4", idFour.DatabaseName(), idFour.Name())

	schemasModelLike := datasourcemodel.Schemas("test1").
		WithLike(prefix+"%").
		WithInDatabase(idOne.DatabaseId()).
		WithWithDescribe(false).
		WithWithParameters(false).
		WithDependsOn(schemaModel1.ResourceReference(), schemaModel2.ResourceReference(), schemaModel3.ResourceReference(), schemaModel4.ResourceReference())

	schemasModelStartsWith := datasourcemodel.Schemas("test2").
		WithStartsWith(prefix+"1").
		WithInDatabase(idOne.DatabaseId()).
		WithWithDescribe(false).
		WithWithParameters(false).
		WithDependsOn(schemaModel1.ResourceReference(), schemaModel2.ResourceReference(), schemaModel3.ResourceReference(), schemaModel4.ResourceReference())

	schemasModelLimit := datasourcemodel.Schemas("test3").
		WithRowsAndFrom(1, prefix+"1").
		WithInDatabase(idOne.DatabaseId()).
		WithWithDescribe(false).
		WithWithParameters(false).
		WithDependsOn(schemaModel1.ResourceReference(), schemaModel2.ResourceReference(), schemaModel3.ResourceReference(), schemaModel4.ResourceReference())

	schemasModelIn := datasourcemodel.Schemas("test4").
		WithInDatabase(idFour.DatabaseId()).
		WithStartsWith(idFour.Name()).
		WithWithDescribe(false).
		WithWithParameters(false).
		WithDependsOn(schemaModel1.ResourceReference(), schemaModel2.ResourceReference(), schemaModel3.ResourceReference(), schemaModel4.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, schemaModel1, schemaModel2, schemaModel3, schemaModel4, schemasModelLike),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(schemasModelLike.DatasourceReference(), "schemas.#", "2"),
				),
			},
			{
				Config: config.FromModels(t, schemaModel1, schemaModel2, schemaModel3, schemaModel4, schemasModelStartsWith),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(schemasModelStartsWith.DatasourceReference(), "schemas.#", "1"),
					resource.TestCheckResourceAttr(schemasModelStartsWith.DatasourceReference(), "schemas.0.show_output.0.name", idOne.Name()),
				),
			},
			{
				Config: config.FromModels(t, schemaModel1, schemaModel2, schemaModel3, schemaModel4, schemasModelLimit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(schemasModelLimit.DatasourceReference(), "schemas.#", "1"),
					resource.TestCheckResourceAttr(schemasModelLimit.DatasourceReference(), "schemas.0.show_output.0.name", idOne.Name()),
				),
			},
			{
				Config: config.FromModels(t, schemaModel1, schemaModel2, schemaModel3, schemaModel4, schemasModelIn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(schemasModelIn.DatasourceReference(), "schemas.#", "1"),
					resource.TestCheckResourceAttr(schemasModelIn.DatasourceReference(), "schemas.0.show_output.0.name", idFour.Name()),
				),
			},
		},
	})
}

func TestAcc_Schemas_CompleteUseCase(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	id := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(database.ID())
	comment := random.Comment()

	viewId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(id)
	statement := "SELECT ROLE_NAME FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	columnNames := []string{"ROLE_NAME"}

	schemaModel := model.Schema("test", id.DatabaseName(), id.Name()).
		WithComment(comment).
		WithIsTransient(datasources.BooleanTrue).
		WithWithManagedAccess(datasources.BooleanTrue)

	viewModel := model.View("test", viewId.DatabaseName(), viewId.SchemaName(), viewId.Name(), statement).
		WithColumnNames(columnNames...).
		WithDependsOn(schemaModel.ResourceReference())

	schemasModel := datasourcemodel.Schemas("test").
		WithLike(id.Name()).
		WithStartsWith(id.Name()).
		WithInDatabase(database.ID()).
		WithLimit(1).
		WithDependsOn(schemaModel.ResourceReference(), viewModel.ResourceReference())

	schemasModelWithoutAdditional := datasourcemodel.Schemas("test").
		WithLike(id.Name()).
		WithStartsWith(id.Name()).
		WithInDatabase(database.ID()).
		WithLimit(1).
		WithWithDescribe(false).
		WithWithParameters(false).
		WithDependsOn(schemaModel.ResourceReference(), viewModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: activeWarehouseSetOnUserProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, schemaModel, viewModel, schemasModelWithoutAdditional),
				Check: assertThat(
					t,
					resourceshowoutputassert.SchemasDatasourceShowOutput(t, schemasModelWithoutAdditional.DatasourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIsDefault(false).
						HasIsCurrent(true).
						HasDatabaseName(id.DatabaseName()).
						HasOwnerNotEmpty().
						HasComment(comment).
						HasOptions("TRANSIENT, MANAGED ACCESS").
						HasRetentionTimeNotEmpty().
						HasOwnerRoleTypeNotEmpty(),

					assert.Check(resource.TestCheckResourceAttr(schemasModelWithoutAdditional.DatasourceReference(), "schemas.0.describe_output.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(schemasModelWithoutAdditional.DatasourceReference(), "schemas.0.parameters.#", "0")),
				),
			},
			{
				Config: config.FromModels(t, schemaModel, viewModel, schemasModel),
				Check: assertThat(
					t,
					resourceshowoutputassert.SchemasDatasourceShowOutput(t, schemasModel.DatasourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIsDefault(false).
						HasIsCurrent(true).
						HasDatabaseName(id.DatabaseName()).
						HasOwnerNotEmpty().
						HasComment(comment).
						HasOptions("TRANSIENT, MANAGED ACCESS").
						HasRetentionTimeNotEmpty().
						HasOwnerRoleTypeNotEmpty(),

					resourceparametersassert.SchemasDatasourceParameters(t, schemasModel.DatasourceReference()).
						HasAllDefaultParameters(),

					assert.Check(resource.TestCheckResourceAttr(schemasModel.DatasourceReference(), "schemas.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.describe_output.0.name")),
					assert.Check(resource.TestCheckResourceAttr(schemasModel.DatasourceReference(), "schemas.0.describe_output.0.kind", "VIEW")),
				),
			},
		},
	})
}

func TestAcc_Schemas_SchemaNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      schemasWithPostcondition(),
				ExpectError: regexp.MustCompile("there should be at least one schema"),
			},
		},
	})
}

func schemasWithPostcondition() string {
	return `
data "snowflake_schemas" "test" {
  like = "non-existing-schema"

  lifecycle {
    postcondition {
      condition     = length(self.schemas) > 0
      error_message = "there should be at least one schema"
    }
  }
}
`
}

func TestAcc_Schemas_BadCombination(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      schemasDatasourceConfigDbAndSchema(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func schemasDatasourceConfigDbAndSchema() string {
	return fmt.Sprintf(`
data "snowflake_schemas" "test" {
  in {
    database = "%s"
    application = "foo"
    application_package = "bar"
  }
}
`, TestDatabaseName)
}
