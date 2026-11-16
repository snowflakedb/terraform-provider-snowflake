//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_BackupSets_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	table, cleanupTable := testClient().Table.CreateWithPredefinedColumns(t)
	t.Cleanup(cleanupTable)
	tableID := table.ID()

	backupSetModel := model.BackupSet("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithForTable(tableID.FullyQualifiedName())

	backupSetsModel := datasourcemodel.BackupSets("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(backupSetModel.ResourceReference())

	backupSetsModelWithoutDescribe := datasourcemodel.BackupSets("test").
		WithWithDescribe(false).
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(backupSetModel.ResourceReference())

	showOutputAssertions := resourceshowoutputassert.BackupSetsDatasourceShowOutput(t, backupSetsModel.DatasourceReference()).
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasObjectKind("TABLE").
		HasObjectName(tableID.Name()).
		HasObjectDatabaseName(tableID.DatabaseName()).
		HasObjectSchemaName(tableID.SchemaName()).
		HasOwnerRole(snowflakeroles.Accountadmin.Name()).
		HasOwnerRoleType("ROLE")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, backupSetModel, backupSetsModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(backupSetsModel.DatasourceReference(), "backup_sets.#", "1")),
					showOutputAssertions,
					resourceshowoutputassert.BackupSetsDatasourceDescribeOutput(t, backupSetsModel.DatasourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasObjectKind("TABLE").
						HasObjectName(tableID.Name()).
						HasObjectDatabaseName(tableID.DatabaseName()).
						HasObjectSchemaName(tableID.SchemaName()).
						HasOwnerRole(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE"),
				),
			},
			{
				Config: accconfig.FromModels(t, backupSetModel, backupSetsModelWithoutDescribe),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(backupSetsModelWithoutDescribe.DatasourceReference(), "backup_sets.#", "1")),
					resourceshowoutputassert.BackupSetsDatasourceShowOutput(t, backupSetsModelWithoutDescribe.DatasourceReference()).
						HasName(id.Name()).
						HasObjectKind("TABLE"),
					assert.Check(resource.TestCheckResourceAttr(backupSetsModelWithoutDescribe.DatasourceReference(), "backup_sets.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_BackupSets_Filtering(t *testing.T) {
	secondSchema, secondSchemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, sdk.NewAccountObjectIdentifier(TestDatabaseName))
	t.Cleanup(secondSchemaCleanup)

	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, secondSchema.ID())

	// One backup set per table.
	table1, cleanupTable1 := testClient().Table.CreateWithPredefinedColumns(t)
	t.Cleanup(cleanupTable1)
	table2, cleanupTable2 := testClient().Table.CreateWithPredefinedColumns(t)
	t.Cleanup(cleanupTable2)
	table3, cleanupTable3 := testClient().Table.CreateWithPredefinedColumns(t)
	t.Cleanup(cleanupTable3)

	model1 := model.BackupSet("test1", id1.DatabaseName(), id1.SchemaName(), id1.Name()).WithForTable(table1.ID().FullyQualifiedName())
	model2 := model.BackupSet("test2", id2.DatabaseName(), id2.SchemaName(), id2.Name()).WithForTable(table2.ID().FullyQualifiedName())
	model3 := model.BackupSet("test3", id3.DatabaseName(), id3.SchemaName(), id3.Name()).WithForTable(table3.ID().FullyQualifiedName())

	dependsOn := []string{model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference()}
	likeFirst := datasourcemodel.BackupSets("test").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(dependsOn...)
	likePrefix := datasourcemodel.BackupSets("test").
		WithLike(prefix + "%").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(dependsOn...)
	inSchema := datasourcemodel.BackupSets("test").
		WithInSchema(id1.SchemaId()).
		WithDependsOn(dependsOn...)
	inDatabase := datasourcemodel.BackupSets("test").
		WithStartsWith(prefix).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(dependsOn...)
	withLimit := datasourcemodel.BackupSets("test").
		WithStartsWith(prefix).
		WithInDatabase(id1.DatabaseId()).
		WithLimit(1).
		WithDependsOn(dependsOn...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, likeFirst),
				Check:  resource.TestCheckResourceAttr(likeFirst.DatasourceReference(), "backup_sets.#", "1"),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, likePrefix),
				Check:  resource.TestCheckResourceAttr(likePrefix.DatasourceReference(), "backup_sets.#", "3"),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, inSchema),
				Check:  resource.TestCheckResourceAttr(inSchema.DatasourceReference(), "backup_sets.#", "2"),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, inDatabase),
				Check:  resource.TestCheckResourceAttr(inDatabase.DatasourceReference(), "backup_sets.#", "3"),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, withLimit),
				Check:  resource.TestCheckResourceAttr(withLimit.DatasourceReference(), "backup_sets.#", "1"),
			},
		},
	})
}

func TestAcc_BackupSets_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.BackupSets("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}
