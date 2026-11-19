//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Tables_BasicUseCase(t *testing.T) {
	table, tableCleanup := testClient().Table.Create(t)
	t.Cleanup(tableCleanup)

	id := table.ID()

	tablesModel := datasourcemodel.Tables("test").
		WithLike(id.Name()).
		WithInSchema(id.SchemaId())
	tablesModelWithoutAdditional := datasourcemodel.Tables("test").
		WithWithDescribe(false).
		WithWithParameters(false).
		WithLike(id.Name()).
		WithInSchema(id.SchemaId())
	tablesModelWithoutDescribe := datasourcemodel.Tables("test").
		WithWithDescribe(false).
		WithLike(id.Name()).
		WithInSchema(id.SchemaId())
	tablesModelWithoutParameters := datasourcemodel.Tables("test").
		WithWithParameters(false).
		WithLike(id.Name()).
		WithInSchema(id.SchemaId())

	showOutputAssertions := resourceshowoutputassert.TablesDatasourceShowOutput(t, tablesModel.DatasourceReference()).
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasKind("TABLE").
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasOwnerRoleType("ROLE")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tablesModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(tablesModel.DatasourceReference(), "tables.#", "1")),
					showOutputAssertions,
					resourceshowoutputassert.TablesDatasourceDescribeOutput(t, tablesModel.DatasourceReference()).
						HasName("ID").
						HasType(testdatatypes.DataTypeNumber_38_0),
					resourceparametersassert.TablesDatasourceParameters(t, tablesModel.DatasourceReference()).
						HasDataRetentionTimeInDays(table.RetentionTime),
				),
			},
			{
				Config: accconfig.FromModels(t, tablesModelWithoutAdditional),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(tablesModelWithoutAdditional.DatasourceReference(), "tables.#", "1")),
					showOutputAssertions,
					assert.Check(resource.TestCheckResourceAttr(tablesModelWithoutAdditional.DatasourceReference(), "tables.0.describe_output.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(tablesModelWithoutAdditional.DatasourceReference(), "tables.0.parameters.#", "0")),
				),
			},
			{
				Config: accconfig.FromModels(t, tablesModelWithoutDescribe),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(tablesModelWithoutDescribe.DatasourceReference(), "tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(tablesModelWithoutDescribe.DatasourceReference(), "tables.0.describe_output.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(tablesModelWithoutDescribe.DatasourceReference(), "tables.0.parameters.#", "1")),
				),
			},
			{
				Config: accconfig.FromModels(t, tablesModelWithoutParameters),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(tablesModelWithoutParameters.DatasourceReference(), "tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(tablesModelWithoutParameters.DatasourceReference(), "tables.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(tablesModelWithoutParameters.DatasourceReference(), "tables.0.parameters.#", "0")),
				),
			},
		},
	})
}

func TestAcc_Tables_Filtering(t *testing.T) {
	secondSchema, secondSchemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, sdk.NewAccountObjectIdentifier(TestDatabaseName))
	t.Cleanup(secondSchemaCleanup)

	prefix := strings.ToUpper(random.AlphaN(4))
	table1, cleanup1 := testClient().Table.CreateWithName(t, testClient().Ids.AlphaWithPrefix(prefix))
	t.Cleanup(cleanup1)
	_, cleanup2 := testClient().Table.CreateWithName(t, testClient().Ids.AlphaWithPrefix(prefix))
	t.Cleanup(cleanup2)
	_, cleanup3 := testClient().Table.CreateWithRequest(t, sdk.NewCreateTableRequest(
		testClient().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, secondSchema.ID()),
		*sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns([]sdk.TableColumnRequest{*sdk.NewTableColumnRequest("ID", testdatatypes.DataTypeNumber)}),
	))
	t.Cleanup(cleanup3)

	likeFirst := datasourcemodel.Tables("test").
		WithLike(table1.ID().Name()).
		WithInDatabase(table1.ID().DatabaseId())
	likePrefixInSchema := datasourcemodel.Tables("test").
		WithLike(prefix + "%").
		WithInSchema(table1.ID().SchemaId())
	likePrefixInDatabase := datasourcemodel.Tables("test").
		WithLike(prefix + "%").
		WithInDatabase(table1.ID().DatabaseId())
	likePrefixInAccount := datasourcemodel.Tables("test").
		WithLike(prefix + "%").
		WithInAccount()
	startsWith := datasourcemodel.Tables("test").
		WithStartsWith(prefix).
		WithInDatabase(table1.ID().DatabaseId())
	withLimit := datasourcemodel.Tables("test").
		WithLike(prefix + "%").
		WithInSchema(table1.ID().SchemaId()).
		WithLimit(1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, likeFirst),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(likeFirst.DatasourceReference(), "tables.#", "1")),
			},
			{
				Config: accconfig.FromModels(t, likePrefixInSchema),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(likePrefixInSchema.DatasourceReference(), "tables.#", "2")),
			},
			{
				Config: accconfig.FromModels(t, likePrefixInDatabase),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(likePrefixInDatabase.DatasourceReference(), "tables.#", "3")),
			},
			{
				Config: accconfig.FromModels(t, likePrefixInAccount),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(likePrefixInAccount.DatasourceReference(), "tables.#", "3")),
			},
			{
				Config: accconfig.FromModels(t, startsWith),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(startsWith.DatasourceReference(), "tables.#", "3")),
			},
			{
				Config: accconfig.FromModels(t, withLimit),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(withLimit.DatasourceReference(), "tables.#", "1")),
			},
		},
	})
}

func TestAcc_Tables_TerseAndHistory(t *testing.T) {
	schema, schemaCleanup := testClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	table, tableCleanup := testClient().Table.CreateInSchema(t, schema.ID())
	t.Cleanup(tableCleanup)

	droppedTable, droppedTableCleanup := testClient().Table.CreateInSchema(t, schema.ID())
	t.Cleanup(droppedTableCleanup)
	testClient().Table.DropFunc(t, droppedTable.ID())()

	terse := datasourcemodel.Tables("test").
		WithInSchema(schema.ID()).
		WithLike(table.ID().Name()).
		WithTerse(true)
	withoutHistory := datasourcemodel.Tables("test").
		WithInSchema(schema.ID()).
		WithLike(droppedTable.ID().Name())
	withHistory := datasourcemodel.Tables("test").
		WithInSchema(schema.ID()).
		WithLike(droppedTable.ID().Name()).
		WithHistory(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, terse),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(terse.DatasourceReference(), "tables.#", "1")),
					resourceshowoutputassert.TablesDatasourceShowOutput(t, terse.DatasourceReference()).
						HasName(table.ID().Name()).
						HasKind("TABLE").
						HasOwner(""),
					resourceshowoutputassert.TablesDatasourceDescribeOutput(t, terse.DatasourceReference()).
						HasName("ID"),
					assert.Check(resource.TestCheckResourceAttr(terse.DatasourceReference(), "tables.0.parameters.#", "1")),
				),
			},
			{
				Config: accconfig.FromModels(t, withoutHistory),
				Check:  resource.TestCheckResourceAttr(withoutHistory.DatasourceReference(), "tables.#", "0"),
			},
			{
				Config: accconfig.FromModels(t, withHistory),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(withHistory.DatasourceReference(), "tables.#", "1")),
					resourceshowoutputassert.TablesDatasourceShowOutput(t, withHistory.DatasourceReference()).
						HasName(droppedTable.ID().Name()).
						HasDroppedOnNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(withHistory.DatasourceReference(), "tables.0.describe_output.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(withHistory.DatasourceReference(), "tables.0.parameters.#", "0")),
				),
			},
		},
	})
}

func TestAcc_Tables_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.Tables("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Tables_InApplication(t *testing.T) {
	applicationPackage, cleanup := testClient().ApplicationPackage.CreateApplicationPackage(t)
	t.Cleanup(cleanup)

	inApplicationPackage := datasourcemodel.Tables("test").
		WithInApplicationPackage(applicationPackage.ID())
	inMissingApplication := datasourcemodel.Tables("test").
		WithInApplication(testClient().Ids.RandomAccountObjectIdentifier())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, inApplicationPackage),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(inApplicationPackage.DatasourceReference(), "tables.#", "0")),
			},
			{
				Config:      accconfig.FromModels(t, inMissingApplication),
				ExpectError: regexp.MustCompile("does not exist"),
			},
		},
	})
}

func TestAcc_Tables_EmptyResult(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	tablesModel := datasourcemodel.Tables("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tablesModel),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(tablesModel.DatasourceReference(), "tables.#", "0")),
			},
		},
	})
}
