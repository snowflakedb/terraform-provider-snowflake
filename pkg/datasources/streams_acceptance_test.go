//go:build !account_level_tests

package datasources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Streams(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	table, cleanupTable := acc.TestClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	streamModel := model.StreamOnTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), table.ID().FullyQualifiedName()).
		WithAppendOnly(datasources.BooleanTrue).
		WithComment(comment)
	streamsModel := datasourcemodel.Streams("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(streamModel.ResourceReference())
	streamsModelWithoutDescribe := datasourcemodel.Streams("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithWithDescribe(false).
		WithDependsOn(streamModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, streamModel, streamsModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.#", "1")),

					resourceshowoutputassert.StreamsDatasourceShowOutput(t, "snowflake_streams.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(table.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeTable).
						HasBaseTables(table.ID()).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeAppendOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(streamsModel.DatasourceReference(), "streams.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.table_name", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.source_type", string(sdk.StreamSourceTypeTable))),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.base_tables.0", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.mode", string(sdk.StreamModeAppendOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(streamsModel.DatasourceReference(), "streams.0.describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.owner_role_type", "ROLE")),
				),
			},
			{
				Config: accconfig.FromModels(t, streamModel, streamsModelWithoutDescribe),

				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.#", "1")),

					resourceshowoutputassert.StreamsDatasourceShowOutput(t, "snowflake_streams.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(table.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeTable).
						HasBaseTables(table.ID()).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeAppendOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_StreamOnTable(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	table, cleanupTable := acc.TestClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	streamModel := model.StreamOnTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), table.ID().FullyQualifiedName()).
		WithAppendOnly(datasources.BooleanTrue).
		WithComment(comment)
	streamsModel := datasourcemodel.Streams("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(streamModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, streamModel, streamsModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.#", "1")),
					resourceshowoutputassert.StreamsDatasourceShowOutput(t, "snowflake_streams.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(table.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeTable).
						HasBaseTables(table.ID()).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeAppendOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(streamsModel.DatasourceReference(), "streams.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.table_name", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.source_type", string(sdk.StreamSourceTypeTable))),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.base_tables.0", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.mode", string(sdk.StreamModeAppendOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(streamsModel.DatasourceReference(), "streams.0.describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.owner_role_type", "ROLE")),
				),
			},
		},
	})
}

func TestAcc_StreamOnExternalTable(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stage, stageCleanup := acc.TestClient().Stage.CreateStageWithURL(t)
	t.Cleanup(stageCleanup)

	externalTable, externalTableCleanup := acc.TestClient().ExternalTable.CreateWithLocation(t, stage.Location())
	t.Cleanup(externalTableCleanup)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	streamModel := model.StreamOnExternalTableBase("test", id, externalTable.ID()).
		WithCopyGrants(true).
		WithComment(comment)
	streamsModel := datasourcemodel.Streams("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(streamModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, streamModel, streamsModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.#", "1")),
					resourceshowoutputassert.StreamsDatasourceShowOutput(t, "snowflake_streams.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(externalTable.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeExternalTable).
						HasBaseTables(externalTable.ID()).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeInsertOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(streamsModel.DatasourceReference(), "streams.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.table_name", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.source_type", string(sdk.StreamSourceTypeExternalTable))),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.base_tables.0", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.mode", string(sdk.StreamModeInsertOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(streamsModel.DatasourceReference(), "streams.0.describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.owner_role_type", "ROLE")),
				),
			},
		},
	})
}

func TestAcc_StreamOnDirectoryTable(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stage, cleanupStage := acc.TestClient().Stage.CreateStageWithDirectory(t)
	t.Cleanup(cleanupStage)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	streamModel := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.SchemaName(), id.Name(), stage.ID().FullyQualifiedName()).
		WithComment(comment)
	streamsModel := datasourcemodel.Streams("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(streamModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, streamModel, streamsModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.#", "1")),
					resourceshowoutputassert.StreamsDatasourceShowOutput(t, "snowflake_streams.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(stage.ID().Name()).
						HasSourceType(sdk.StreamSourceTypeStage).
						HasBaseTablesPartiallyQualified(stage.ID().Name()).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeDefault).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(streamsModel.DatasourceReference(), "streams.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.table_name", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.source_type", string(sdk.StreamSourceTypeStage))),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.base_tables.0", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.mode", string(sdk.StreamModeDefault))),
					assert.Check(resource.TestCheckResourceAttrSet(streamsModel.DatasourceReference(), "streams.0.describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.owner_role_type", "ROLE")),
				),
			},
		},
	})
}

func TestAcc_StreamOnView(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	table, cleanupTable := acc.TestClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)

	statement := fmt.Sprintf("SELECT * FROM %s", table.ID().FullyQualifiedName())
	view, cleanupView := acc.TestClient().View.CreateView(t, statement)
	t.Cleanup(cleanupView)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	streamModel := model.StreamOnView("test", id.DatabaseName(), id.SchemaName(), id.Name(), view.ID().FullyQualifiedName()).
		WithComment(comment).
		WithAppendOnly(datasources.BooleanTrue)
	streamsModel := datasourcemodel.Streams("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(streamModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, streamModel, streamsModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.#", "1")),
					resourceshowoutputassert.StreamsDatasourceShowOutput(t, "snowflake_streams.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(view.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeView).
						HasBaseTables(table.ID()).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeAppendOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(streamsModel.DatasourceReference(), "streams.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.table_name", view.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.source_type", string(sdk.StreamSourceTypeView))),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.base_tables.0", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.mode", string(sdk.StreamModeAppendOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(streamsModel.DatasourceReference(), "streams.0.describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(streamsModel.DatasourceReference(), "streams.0.describe_output.0.owner_role_type", "ROLE")),
				),
			},
		},
	})
}

func TestAcc_Streams_Filtering(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	table, cleanupTable := acc.TestClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)

	prefix := random.AlphaN(4)
	id1 := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	model1 := model.StreamOnTable("test_1", id1.DatabaseName(), id1.SchemaName(), id1.Name(), table.ID().FullyQualifiedName())
	model2 := model.StreamOnTable("test_2", id2.DatabaseName(), id2.SchemaName(), id2.Name(), table.ID().FullyQualifiedName())
	model3 := model.StreamOnTable("test_3", id3.DatabaseName(), id3.SchemaName(), id3.Name(), table.ID().FullyQualifiedName())
	streamsModelLikeFirstOne := datasourcemodel.Streams("test").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	streamsModelLikePrefix := datasourcemodel.Streams("test").
		WithLike(prefix+"%").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() { acc.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, streamsModelLikeFirstOne),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(streamsModelLikeFirstOne.DatasourceReference(), "streams.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, streamsModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(streamsModelLikePrefix.DatasourceReference(), "streams.#", "2"),
				),
			},
		},
	})
}

func TestAcc_Streams_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      streamsDatasourceEmptyIn(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func streamsDatasourceEmptyIn() string {
	return `
data "snowflake_streams" "test" {
  in {
  }
}
`
}

func TestAcc_Streams_NotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streams/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one stream"),
			},
		},
	})
}
