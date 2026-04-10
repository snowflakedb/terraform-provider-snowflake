//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAcc_Experimental_TagAssociation_SafeDestroy_MissingColumn verifies that destroying a tag_association
// for a column fails when the column is deleted externally (default behavior), and succeeds
// when the TAG_ASSOCIATION_SAFE_DESTROY experiment is enabled.
//
// This is a regression test for https://github.com/snowflakedb/terraform-provider-snowflake/issues/3869.
func TestAcc_Experimental_TagAssociation_SafeDestroy_MissingColumn(t *testing.T) {
	tag, tagCleanup := testClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	table, tableCleanup := testClient().Table.CreateWithPredefinedColumns(t)
	t.Cleanup(tableCleanup)

	columnId := sdk.NewTableColumnIdentifier(table.ID().DatabaseName(), table.ID().SchemaName(), table.ID().Name(), "ID")
	columnId2 := sdk.NewTableColumnIdentifier(table.ID().DatabaseName(), table.ID().SchemaName(), table.ID().Name(), "SOME_TEXT_COLUMN")

	tagAssociationModel := model.TagAssociation("test", []sdk.ObjectIdentifier{columnId, columnId2}, string(sdk.ObjectTypeColumn), tag.ID().FullyQualifiedName(), "v1").
		WithSkipValidation(true)

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.TagAssociationSafeDestroy)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create the tag association with default provider.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, tagAssociationModel),
			},
			// Drop the column externally (SYSTEM$GET_TAG on missing object fails).
			// Without experiment, Read propagates the error and the destroy fails.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig: func() {
					testClient().Table.AlterWithRequest(t, sdk.NewAlterTableRequest(table.ID()).
						WithColumnAction(sdk.NewTableColumnActionRequest().
							WithDropColumns([]string{columnId.Name()})))
				},
				Config:      config.FromModels(t, tagAssociationModel),
				Destroy:     true,
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
			// Destroy with TAG_ASSOCIATION_SAFE_DESTROY experiment.
			{
				ProtoV6ProviderFactories: tagAssociationSafeDestroyProviderFactory,
				Config:                   config.FromModels(t, experimentProviderModel, tagAssociationModel),
				Destroy:                  true,
			},
		},
	})
}

// TestAcc_Experimental_TagAssociation_SafeDestroy_MissingTable verifies that destroying a tag_association
// for a column fails when the column's table is deleted externally (default behavior), and succeeds when
// the TAG_ASSOCIATION_SAFE_DESTROY experiment is enabled.
//
// This is a regression test for https://github.com/snowflakedb/terraform-provider-snowflake/issues/3869.
func TestAcc_Experimental_TagAssociation_SafeDestroy_MissingTable(t *testing.T) {
	tag, tagCleanup := testClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	table, tableCleanup := testClient().Table.CreateWithPredefinedColumns(t)
	t.Cleanup(tableCleanup)

	columnId := sdk.NewTableColumnIdentifier(table.ID().DatabaseName(), table.ID().SchemaName(), table.ID().Name(), "ID")
	columnId2 := sdk.NewTableColumnIdentifier(table.ID().DatabaseName(), table.ID().SchemaName(), table.ID().Name(), "SOME_TEXT_COLUMN")

	tagAssociationModel := model.TagAssociation("test", []sdk.ObjectIdentifier{columnId, columnId2}, string(sdk.ObjectTypeColumn), tag.ID().FullyQualifiedName(), "v1").
		WithSkipValidation(true)

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.TagAssociationSafeDestroy)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create the tag association with default provider.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, tagAssociationModel),
			},
			// Drop the table externally (SYSTEM$GET_TAG on missing object fails).
			// Without experiment, Read propagates the error and the destroy fails.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig:                testClient().Table.DropFunc(t, table.ID()),
				Config:                   config.FromModels(t, tagAssociationModel),
				Destroy:                  true,
				ExpectError:              regexp.MustCompile("does not exist or not authorized"),
			},
			// Destroy with TAG_ASSOCIATION_SAFE_DESTROY experiment.
			{
				ProtoV6ProviderFactories: tagAssociationSafeDestroyProviderFactory,
				Config:                   config.FromModels(t, experimentProviderModel, tagAssociationModel),
				Destroy:                  true,
			},
		},
	})
}

// TestAcc_Experimental_TagAssociation_SafeDestroy_MissingSchema verifies that destroying a tag_association
// for a column fails when the column's parent schema is deleted externally (default behavior), and succeeds
// when the TAG_ASSOCIATION_SAFE_DESTROY experiment is enabled.
//
// This is a regression test for https://github.com/snowflakedb/terraform-provider-snowflake/issues/3869,
// where skipColumnIfDoesNotExist propagated ErrDoesNotExistOrOperationCannotBePerformed (from
// SHOW TABLES IN SCHEMA <dropped_schema>) instead of silently succeeding.
func TestAcc_Experimental_TagAssociation_SafeDestroy_MissingSchema(t *testing.T) {
	tag, tagCleanup := testClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	table, _ := testClient().Table.CreateWithRequest(t,
		sdk.NewCreateTableRequest(
			testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID()),
			[]sdk.TableColumnRequest{*sdk.NewTableColumnRequest("ID", sdk.DataTypeInt), *sdk.NewTableColumnRequest("ID2", sdk.DataTypeInt)},
		),
	)

	columnId := sdk.NewTableColumnIdentifier(table.ID().DatabaseName(), table.ID().SchemaName(), table.ID().Name(), "ID")
	columnId2 := sdk.NewTableColumnIdentifier(table.ID().DatabaseName(), table.ID().SchemaName(), table.ID().Name(), "ID2")

	tagAssociationModel := model.TagAssociation("test", []sdk.ObjectIdentifier{columnId, columnId2}, string(sdk.ObjectTypeColumn), tag.ID().FullyQualifiedName(), "v1").
		WithSkipValidation(true)

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.TagAssociationSafeDestroy)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create the tag association with default provider.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, tagAssociationModel),
			},
			// Drop the schema externally (SHOW TABLES IN SCHEMA <dropped_schema> fails).
			// Without experiment:
			//   - Read: SYSTEM$GET_TAG fails
			// OR (if Read somehow keeps state):
			//   - Delete: skipColumnIfDoesNotExist calls ShowByID on dropped schema
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig:                testClient().Schema.DropSchemaFunc(t, table.ID().SchemaId()),
				Config:                   config.FromModels(t, tagAssociationModel),
				Destroy:                  true,
				ExpectError:              regexp.MustCompile("does not exist"),
			},
			// Destroy with TAG_ASSOCIATION_SAFE_DESTROY experiment.
			{
				ProtoV6ProviderFactories: tagAssociationSafeDestroyProviderFactory,
				Config:                   config.FromModels(t, experimentProviderModel, tagAssociationModel),
				Destroy:                  true,
			},
		},
	})
}
