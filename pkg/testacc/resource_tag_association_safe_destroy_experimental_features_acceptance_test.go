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

// TestAcc_Experimental_TagAssociation_SafeDestroy verifies that destroying a tag_association for columns
// fails when the backing object is removed externally (default behavior), and succeeds when the
// TAG_ASSOCIATION_SAFE_DESTROY experiment is enabled.
//
// Regression test for https://github.com/snowflakedb/terraform-provider-snowflake/issues/3869.
func TestAcc_Experimental_TagAssociation_SafeDestroy(t *testing.T) {
	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.TagAssociationSafeDestroy)

	cases := []struct {
		name  string
		setup func(t *testing.T) (tagAssociationModel *model.TagAssociationModel, preDestroy func())
	}{
		{
			name: "missing_column",
			setup: func(t *testing.T) (*model.TagAssociationModel, func()) {
				tag, tagCleanup := testClient().Tag.CreateTag(t)
				t.Cleanup(tagCleanup)

				table, tableCleanup := testClient().Table.CreateWithPredefinedColumns(t)
				t.Cleanup(tableCleanup)

				columnId := sdk.NewTableColumnIdentifier(table.ID().DatabaseName(), table.ID().SchemaName(), table.ID().Name(), "ID")
				columnId2 := sdk.NewTableColumnIdentifier(table.ID().DatabaseName(), table.ID().SchemaName(), table.ID().Name(), "SOME_TEXT_COLUMN")

				tagAssociationModel := model.TagAssociation("test", []sdk.ObjectIdentifier{columnId, columnId2}, string(sdk.ObjectTypeColumn), tag.ID().FullyQualifiedName(), "v1").
					WithSkipValidation(true)

				return tagAssociationModel, func() {
					testClient().Table.AlterWithRequest(t, sdk.NewAlterTableRequest(table.ID()).
						WithColumnAction(sdk.NewTableColumnActionRequest().
							WithDropColumns([]string{columnId.Name()})))
				}
			},
		},
		{
			name: "missing_table",
			setup: func(t *testing.T) (*model.TagAssociationModel, func()) {
				tag, tagCleanup := testClient().Tag.CreateTag(t)
				t.Cleanup(tagCleanup)

				table, tableCleanup := testClient().Table.CreateWithPredefinedColumns(t)
				t.Cleanup(tableCleanup)

				columnId := sdk.NewTableColumnIdentifier(table.ID().DatabaseName(), table.ID().SchemaName(), table.ID().Name(), "ID")
				columnId2 := sdk.NewTableColumnIdentifier(table.ID().DatabaseName(), table.ID().SchemaName(), table.ID().Name(), "SOME_TEXT_COLUMN")

				tagAssociationModel := model.TagAssociation("test", []sdk.ObjectIdentifier{columnId, columnId2}, string(sdk.ObjectTypeColumn), tag.ID().FullyQualifiedName(), "v1").
					WithSkipValidation(true)

				return tagAssociationModel, testClient().Table.DropFunc(t, table.ID())
			},
		},
		{
			// Without experiment, Delete can surface ErrDoesNotExistOrOperationCannotBePerformed from
			// SHOW TABLES IN SCHEMA on a dropped schema (see skipColumnIfDoesNotExist).
			name: "missing_schema",
			setup: func(t *testing.T) (*model.TagAssociationModel, func()) {
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

				return tagAssociationModel, testClient().Schema.DropSchemaFunc(t, table.ID().SchemaId())
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tagAssociationModel, preDestroy := tc.setup(t)

			resource.Test(t, resource.TestCase{
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				Steps: []resource.TestStep{
					{
						ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
						Config:                   config.FromModels(t, tagAssociationModel),
					},
					{
						ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
						PreConfig:                preDestroy,
						Config:                   config.FromModels(t, tagAssociationModel),
						Destroy:                  true,
						ExpectError:              regexp.MustCompile("does not exist or not authorized"),
					},
					{
						ProtoV6ProviderFactories: tagAssociationSafeDestroyProviderFactory,
						Config:                   config.FromModels(t, experimentProviderModel, tagAssociationModel),
						Destroy:                  true,
					},
				},
			})
		})
	}
}
