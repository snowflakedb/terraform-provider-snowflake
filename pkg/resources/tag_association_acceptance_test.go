//go:build !account_level_tests

package resources_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TagAssociation(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tagId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	tag2Id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	tagValue := "foo"
	tagValue2 := "bar"
	databaseId := acc.TestClient().Ids.DatabaseId()

	tag1Model := model.TagBase("test", tagId).WithAllowedValues("bar", "foo", "external")
	tag2Model := model.TagBase("test", tag2Id).WithAllowedValues("bar", "foo", "external")
	tagAssociationModel := model.TagAssociation("test", []sdk.ObjectIdentifier{databaseId}, string(sdk.ObjectTypeDatabase), tagId.FullyQualifiedName(), tagValue).
		WithDependsOn(tag1Model.ResourceReference())
	tagAssociationModelDifferentValue := model.TagAssociation("test", []sdk.ObjectIdentifier{databaseId}, string(sdk.ObjectTypeDatabase), tagId.FullyQualifiedName(), tagValue2).
		WithDependsOn(tag1Model.ResourceReference())
	tagAssociationModelTag2 := model.TagAssociation("test", []sdk.ObjectIdentifier{databaseId}, string(sdk.ObjectTypeDatabase), tag2Id.FullyQualifiedName(), tagValue2).
		WithDependsOn(tag2Model.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckResourceTagUnset(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tag1Model, tagAssociationModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), tagValue, string(sdk.ObjectTypeDatabase))),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_type", string(sdk.ObjectTypeDatabase)),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(tagAssociationModel.ResourceReference(), "object_identifiers.*", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_value", tagValue),
				),
			},
			// external change - unset tag
			{
				PreConfig: func() {
					acc.TestClient().Tag.Unset(t, sdk.ObjectTypeDatabase, databaseId, []sdk.ObjectIdentifier{tagId})
				},
				Config: accconfig.FromModels(t, tag1Model, tagAssociationModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), tagValue, string(sdk.ObjectTypeDatabase))),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_type", string(sdk.ObjectTypeDatabase)),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(tagAssociationModel.ResourceReference(), "object_identifiers.*", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_value", tagValue),
				),
			},
			// external change - set a different value
			{
				PreConfig: func() {
					acc.TestClient().Tag.Set(t, sdk.ObjectTypeDatabase, databaseId, []sdk.TagAssociation{
						{
							Name:  tagId,
							Value: "external",
						},
					})
				},
				Config: accconfig.FromModels(t, tag1Model, tagAssociationModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), tagValue, string(sdk.ObjectTypeDatabase))),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_type", string(sdk.ObjectTypeDatabase)),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(tagAssociationModel.ResourceReference(), "object_identifiers.*", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_value", tagValue),
				),
			},
			// change tag value
			{
				Config: accconfig.FromModels(t, tag1Model, tagAssociationModelDifferentValue),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(tagAssociationModelDifferentValue.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModelDifferentValue.ResourceReference(), "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), tagValue2, string(sdk.ObjectTypeDatabase))),
					resource.TestCheckResourceAttr(tagAssociationModelDifferentValue.ResourceReference(), "object_type", string(sdk.ObjectTypeDatabase)),
					resource.TestCheckResourceAttr(tagAssociationModelDifferentValue.ResourceReference(), "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(tagAssociationModelDifferentValue.ResourceReference(), "object_identifiers.*", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModelDifferentValue.ResourceReference(), "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModelDifferentValue.ResourceReference(), "tag_value", tagValue2),
				),
			},
			// change tag id
			{
				Config: accconfig.FromModels(t, tag2Model, tagAssociationModelTag2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(tagAssociationModelTag2.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModelTag2.ResourceReference(), "id", helpers.EncodeSnowflakeID(tag2Id.FullyQualifiedName(), tagValue2, string(sdk.ObjectTypeDatabase))),
					resource.TestCheckResourceAttr(tagAssociationModelTag2.ResourceReference(), "object_type", string(sdk.ObjectTypeDatabase)),
					resource.TestCheckResourceAttr(tagAssociationModelTag2.ResourceReference(), "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(tagAssociationModelTag2.ResourceReference(), "object_identifiers.*", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModelTag2.ResourceReference(), "tag_id", tag2Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModelTag2.ResourceReference(), "tag_value", tagValue2),
					acc.CheckTagUnset(t, tagId, databaseId, sdk.ObjectTypeDatabase),
				),
			},
			{
				Config:            accconfig.FromModels(t, tag2Model, tagAssociationModelTag2),
				ResourceName:      tagAssociationModelTag2.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				// object_identifiers does not get set because during the import, the configuration is considered as empty
				ImportStateVerifyIgnore: []string{"skip_validation", "object_identifiers.#", "object_identifiers.0"},
			},
			// after refreshing the state, object_identifiers is correct
			{
				RefreshState: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModelTag2.ResourceReference(), "id", helpers.EncodeSnowflakeID(tag2Id.FullyQualifiedName(), tagValue2, string(sdk.ObjectTypeDatabase))),
					resource.TestCheckResourceAttr(tagAssociationModelTag2.ResourceReference(), "object_type", string(sdk.ObjectTypeDatabase)),
					resource.TestCheckResourceAttr(tagAssociationModelTag2.ResourceReference(), "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(tagAssociationModelTag2.ResourceReference(), "object_identifiers.*", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModelTag2.ResourceReference(), "tag_id", tag2Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModelTag2.ResourceReference(), "tag_value", tagValue2),
				),
			},
		},
	})
}

func TestAcc_TagAssociation_objectIdentifiers(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tag, tagCleanup := acc.TestClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)
	dbRole1, dbRole1Cleanup := acc.TestClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(dbRole1Cleanup)
	dbRole2, dbRole2Cleanup := acc.TestClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(dbRole2Cleanup)
	dbRole3, dbRole3Cleanup := acc.TestClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(dbRole3Cleanup)

	model12 := model.TagAssociation("test", []sdk.ObjectIdentifier{dbRole1.ID(), dbRole2.ID()}, string(sdk.ObjectTypeDatabaseRole), tag.ID().FullyQualifiedName(), "foo")
	model123 := model.TagAssociation("test", []sdk.ObjectIdentifier{dbRole1.ID(), dbRole2.ID(), dbRole3.ID()}, string(sdk.ObjectTypeDatabaseRole), tag.ID().FullyQualifiedName(), "foo")
	model13 := model.TagAssociation("test", []sdk.ObjectIdentifier{dbRole1.ID(), dbRole3.ID()}, string(sdk.ObjectTypeDatabaseRole), tag.ID().FullyQualifiedName(), "foo")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			acc.CheckResourceTagUnset(t),
		),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model12),
				Check: assertThat(t, resourceassert.TagAssociationResource(t, model12.ResourceReference()).
					HasObjectTypeString(string(sdk.ObjectTypeDatabaseRole)).
					HasTagIdString(tag.ID().FullyQualifiedName()).
					HasObjectIdentifiersLength(2).
					HasTagValueString("foo"),
					assert.Check(resource.TestCheckTypeSetElemAttr(model12.ResourceReference(), "object_identifiers.*", dbRole1.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckTypeSetElemAttr(model12.ResourceReference(), "object_identifiers.*", dbRole2.ID().FullyQualifiedName())),
				),
			},
			{
				Config: accconfig.FromModels(t, model123),
				Check: assertThat(t, resourceassert.TagAssociationResource(t, model12.ResourceReference()).
					HasObjectTypeString(string(sdk.ObjectTypeDatabaseRole)).
					HasTagIdString(tag.ID().FullyQualifiedName()).
					HasObjectIdentifiersLength(3).
					HasTagValueString("foo"),
					assert.Check(resource.TestCheckTypeSetElemAttr(model12.ResourceReference(), "object_identifiers.*", dbRole1.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckTypeSetElemAttr(model12.ResourceReference(), "object_identifiers.*", dbRole2.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckTypeSetElemAttr(model12.ResourceReference(), "object_identifiers.*", dbRole3.ID().FullyQualifiedName())),
				),
			},
			{
				Config: accconfig.FromModels(t, model13),
				Check: assertThat(t, resourceassert.TagAssociationResource(t, model13.ResourceReference()).
					HasObjectTypeString(string(sdk.ObjectTypeDatabaseRole)).
					HasTagIdString(tag.ID().FullyQualifiedName()).
					HasObjectIdentifiersLength(2).
					HasTagValueString("foo"),
					assert.Check(resource.TestCheckTypeSetElemAttr(model13.ResourceReference(), "object_identifiers.*", dbRole1.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckTypeSetElemAttr(model13.ResourceReference(), "object_identifiers.*", dbRole3.ID().FullyQualifiedName())),
					assert.Check(acc.CheckTagUnset(t, tag.ID(), dbRole2.ID(), sdk.ObjectTypeDatabaseRole)),
				),
			},
		},
	})
}

func TestAcc_TagAssociation_objectType(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tag, tagCleanup := acc.TestClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)
	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	dbRole, dbRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(dbRoleCleanup)

	baseModel := model.TagAssociation("test", []sdk.ObjectIdentifier{role.ID()}, string(sdk.ObjectTypeRole), tag.ID().FullyQualifiedName(), "foo")
	modelWithDifferentObjectType := model.TagAssociation("test", []sdk.ObjectIdentifier{dbRole.ID()}, string(sdk.ObjectTypeDatabaseRole), tag.ID().FullyQualifiedName(), "foo")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			acc.CheckResourceTagUnset(t),
		),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, baseModel),
				Check: assertThat(t, resourceassert.TagAssociationResource(t, baseModel.ResourceReference()).
					HasObjectTypeString(string(sdk.ObjectTypeRole)).
					HasTagIdString(tag.ID().FullyQualifiedName()).
					HasObjectIdentifiersLength(1).
					HasTagValueString("foo"),
				),
			},
			{
				Config: accconfig.FromModels(t, modelWithDifferentObjectType),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithDifferentObjectType.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t, resourceassert.TagAssociationResource(t, baseModel.ResourceReference()).
					HasObjectTypeString(string(sdk.ObjectTypeDatabaseRole)).
					HasTagIdString(tag.ID().FullyQualifiedName()).
					HasObjectIdentifiersLength(1).
					HasTagValueString("foo"),
					assert.Check(acc.CheckTagUnset(t, tag.ID(), role.ID(), sdk.ObjectTypeRole)),
				),
			},
		},
	})
}

func TestAcc_TagAssociationSchema(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tagId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	schemaId := acc.TestClient().Ids.SchemaId()
	tagValue := "TAG_VALUE"

	tagModel := model.TagBase("test", tagId)
	tagAssociationModel := model.TagAssociation("test", []sdk.ObjectIdentifier{schemaId}, string(sdk.ObjectTypeSchema), tagId.FullyQualifiedName(), tagValue).
		WithDependsOn(tagModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tagModel, tagAssociationModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), tagValue, string(sdk.ObjectTypeSchema))),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_type", string(sdk.ObjectTypeSchema)),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(tagAssociationModel.ResourceReference(), "object_identifiers.*", schemaId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_value", tagValue),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3235 is fixed
func TestAcc_TagAssociation_lowercaseObjectType(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tag, tagCleanup := acc.TestClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	objectType := strings.ToLower(string(sdk.ObjectTypeSchema))
	objectId := acc.TestClient().Ids.SchemaId()
	tagValue := "TAG_VALUE"

	tagAssociationModel := model.TagAssociation("test", []sdk.ObjectIdentifier{objectId}, objectType, tag.ID().FullyQualifiedName(), tagValue)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tagAssociationModel),
				Check: assertThat(t, resourceassert.TagAssociationResource(t, tagAssociationModel.ResourceReference()).
					HasIdString(helpers.EncodeSnowflakeID(tag.ID().FullyQualifiedName(), tagValue, string(sdk.ObjectTypeSchema))).
					HasObjectTypeString(string(sdk.ObjectTypeSchema)).
					HasTagIdString(tag.ID().FullyQualifiedName()).
					HasObjectIdentifiersLength(1).
					HasTagValueString(tagValue),
				),
			},
		},
	})
}

func TestAcc_TagAssociationColumn(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tag, tagCleanup := acc.TestClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	table, tableCleanup := acc.TestClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("id", sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest("TEST_COLUMN", sdk.DataTypeNumber),
	})
	t.Cleanup(tableCleanup)

	tagId := tag.ID()
	tableId := table.ID()
	columnId := sdk.NewTableColumnIdentifier(tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), "TEST_COLUMN")
	tagValue := "TAG_VALUE"

	tagAssociationModel := model.TagAssociation("test", []sdk.ObjectIdentifier{columnId}, string(sdk.ObjectTypeColumn), tag.ID().FullyQualifiedName(), tagValue)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tagAssociationModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), tagValue, string(sdk.ObjectTypeColumn))),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_type", string(sdk.ObjectTypeColumn)),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(tagAssociationModel.ResourceReference(), "object_identifiers.*", columnId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_value", tagValue),
				),
			},
		},
	})
}

func TestAcc_TagAssociationIssue1202(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tag, tagCleanup := acc.TestClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	table, tableCleanup := acc.TestClient().Table.CreateWithPredefinedColumns(t)
	t.Cleanup(tableCleanup)

	tagId := tag.ID()
	tableId := table.ID()
	tagValue := "v1"

	tagAssociationModel := model.TagAssociation("test", []sdk.ObjectIdentifier{tableId}, string(sdk.ObjectTypeTable), tag.ID().FullyQualifiedName(), tagValue)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tagAssociationModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_type", "TABLE"),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_value", "v1"),
				),
			},
		},
	})
}

func TestAcc_TagAssociationIssue1909(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tag, tagCleanup := acc.TestClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	table, tableCleanup := acc.TestClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("id", sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest(`"test.column"`, sdk.DataTypeNumber),
	})
	t.Cleanup(tableCleanup)

	table2, table2Cleanup := acc.TestClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("id", sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest(`"test.column"`, sdk.DataTypeNumber),
	})
	t.Cleanup(table2Cleanup)

	tagId := tag.ID()
	tableId := table.ID()
	table2Id := table2.ID()
	columnId := sdk.NewTableColumnIdentifier(tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), "test.column")
	column2Id := sdk.NewTableColumnIdentifier(table2Id.DatabaseName(), table2Id.SchemaName(), table2Id.Name(), "test.column")
	tagValue := "v1"

	tagAssociationModel := model.TagAssociation("test", []sdk.ObjectIdentifier{columnId, column2Id}, string(sdk.ObjectTypeColumn), tag.ID().FullyQualifiedName(), tagValue)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tagAssociationModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_type", string(sdk.ObjectTypeColumn)),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_value", tagValue),
					testAccCheckTableColumnTagAssociation(tagId, columnId, tagValue),
					testAccCheckTableColumnTagAssociation(tagId, column2Id, tagValue),
				),
			},
		},
	})
}

func testAccCheckTableColumnTagAssociation(tagID sdk.SchemaObjectIdentifier, objectID sdk.ObjectIdentifier, tagValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()
		tv, err := client.SystemFunctions.GetTag(ctx, tagID, objectID, sdk.ObjectTypeColumn)
		if err != nil {
			return err
		}
		if tv == nil {
			return fmt.Errorf("expected tag value %s, got nil", tagValue)
		}
		if tagValue != *tv {
			return fmt.Errorf("expected tag value %s, got %s", tagValue, *tv)
		}
		return nil
	}
}

// TODO(SNOW-1165821): use a separate account with ORGADMIN in CI

func TestAcc_TagAssociationAccountIssues1910(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tagId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	accountId := acc.TestClient().Context.CurrentAccountIdentifier(t)
	tagValue := "v1"

	tagModel := model.TagBase("test", tagId)
	tagAssociationModel := model.TagAssociation("test", []sdk.ObjectIdentifier{accountId}, string(sdk.ObjectTypeAccount), tagId.FullyQualifiedName(), tagValue).
		WithDependsOn(tagModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tagModel, tagAssociationModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_type", string(sdk.ObjectTypeAccount)),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(tagAssociationModel.ResourceReference(), "object_identifiers.*", accountId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_value", tagValue),
				),
			},
		},
	})
}

func TestAcc_TagAssociationIssue1926(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tagId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	tableId1 := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	tableId2 := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix("table.test")
	columnId1 := sdk.NewTableColumnIdentifier(tableId1.DatabaseName(), tableId1.SchemaName(), tableId1.Name(), "init")
	columnId2 := sdk.NewTableColumnIdentifier(tableId2.DatabaseName(), tableId2.SchemaName(), tableId2.Name(), "column")
	columnId3 := sdk.NewTableColumnIdentifier(tableId2.DatabaseName(), tableId2.SchemaName(), tableId2.Name(), "column.test")

	columns := func(columnIdentifier sdk.TableColumnIdentifier) []sdk.TableColumnSignature {
		return []sdk.TableColumnSignature{
			// TODO(SNOW-1348114): use only one column, if possible.
			// We need a dummy column here because a table must have at least one column, and when we rename the second one in the config, it gets dropped for a moment.
			{Name: "DUMMY", Type: testdatatypes.DataTypeVariant},
			{Name: columnIdentifier.Name(), Type: testdatatypes.DataTypeVariant},
		}
	}

	tagModel := model.TagBase("test", tagId)
	tableModel1 := model.TableWithId("test", tableId1, columns(columnId1))
	tableModel2 := model.TableWithId("test", tableId2, columns(columnId2))
	tableModel3 := model.TableWithId("test", tableId2, columns(columnId3))
	tagAssociationModel1 := model.TagAssociation("test", []sdk.ObjectIdentifier{columnId1}, string(sdk.ObjectTypeColumn), tagId.FullyQualifiedName(), "TAG_VALUE").
		WithDependsOn(tableModel1.ResourceReference())
	tagAssociationModel2 := model.TagAssociation("test", []sdk.ObjectIdentifier{columnId2}, string(sdk.ObjectTypeColumn), tagId.FullyQualifiedName(), "TAG_VALUE").
		WithDependsOn(tableModel2.ResourceReference())
	tagAssociationModel3 := model.TagAssociation("test", []sdk.ObjectIdentifier{columnId3}, string(sdk.ObjectTypeColumn), tagId.FullyQualifiedName(), "TAG_VALUE").
		WithDependsOn(tableModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tagModel, tableModel1, tagAssociationModel1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModel1.ResourceReference(), "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), "TAG_VALUE", string(sdk.ObjectTypeColumn))),
					resource.TestCheckResourceAttr(tagAssociationModel1.ResourceReference(), "object_type", string(sdk.ObjectTypeColumn)),
					resource.TestCheckResourceAttr(tagAssociationModel1.ResourceReference(), "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(tagAssociationModel1.ResourceReference(), "object_identifiers.*", columnId1.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel1.ResourceReference(), "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel1.ResourceReference(), "tag_value", "TAG_VALUE"),
				),
			},
			{
				Config: accconfig.FromModels(t, tagModel, tableModel2, tagAssociationModel2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModel2.ResourceReference(), "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), "TAG_VALUE", string(sdk.ObjectTypeColumn))),
					resource.TestCheckResourceAttr(tagAssociationModel2.ResourceReference(), "object_type", string(sdk.ObjectTypeColumn)),
					resource.TestCheckResourceAttr(tagAssociationModel2.ResourceReference(), "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(tagAssociationModel2.ResourceReference(), "object_identifiers.*", columnId2.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel2.ResourceReference(), "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel2.ResourceReference(), "tag_value", "TAG_VALUE"),
				),
			},
			{
				Config: accconfig.FromModels(t, tagModel, tableModel3, tagAssociationModel3),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModel3.ResourceReference(), "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), "TAG_VALUE", string(sdk.ObjectTypeColumn))),
					resource.TestCheckResourceAttr(tagAssociationModel3.ResourceReference(), "object_type", string(sdk.ObjectTypeColumn)),
					resource.TestCheckResourceAttr(tagAssociationModel3.ResourceReference(), "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(tagAssociationModel3.ResourceReference(), "object_identifiers.*", columnId3.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel3.ResourceReference(), "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel3.ResourceReference(), "tag_value", "TAG_VALUE"),
				),
			},
		},
	})
}

func TestAcc_TagAssociation_migrateFromVersion_0_98_0(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tagId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	schemaId := acc.TestClient().Ids.SchemaId()
	tagValue := "TAG_VALUE"

	tagModel := model.TagBase("test", tagId)
	tagAssociationModel := model.TagAssociation("test", []sdk.ObjectIdentifier{schemaId}, string(sdk.ObjectTypeSchema), tagId.FullyQualifiedName(), tagValue).
		WithDependsOn(tagModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetLegacyConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.98.0"),
				Config:            tagAssociationV098(tagId, tagValue, schemaId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "id", helpers.EncodeSnowflakeID(tagId.DatabaseName(), tagId.SchemaName(), tagId.Name())),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_type", string(sdk.ObjectTypeSchema)),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_identifier.#", "1"),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_identifier.0.name", schemaId.Name()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_identifier.0.database", schemaId.DatabaseName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_identifier.0.schema", ""),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_value", tagValue),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, tagModel, tagAssociationModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(tagAssociationModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(tagAssociationModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), tagValue, string(sdk.ObjectTypeSchema))),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_type", string(sdk.ObjectTypeSchema)),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(tagAssociationModel.ResourceReference(), "object_identifiers.*", schemaId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_value", tagValue),
				),
			},
		},
	})
}

func tagAssociationV098(tagId sdk.SchemaObjectIdentifier, tagValue string, schemaId sdk.DatabaseObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_tag" "test" {
  database       = "%[1]s"
  schema         = "%[2]s"
  name           = "%[3]s"
}

resource "snowflake_tag_association" "test" {
	tag_id					= snowflake_tag.test.fully_qualified_name
	tag_value				= "%[4]s"
	object_type				= "%[5]s"
	object_identifier {
		database = "%[6]s"
		name     = "%[7]s"
	}
}
`, tagId.DatabaseName(), tagId.SchemaName(), tagId.Name(), tagValue, sdk.ObjectTypeSchema, schemaId.DatabaseName(), schemaId.Name())
}

// proves https://github.com/snowflakedb/terraform-provider-snowflake/issues/3622 is fixed
func TestAcc_TagAssociation_issue_3622(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tagId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	table, cleanupTable := acc.TestClient().Table.Create(t)
	t.Cleanup(cleanupTable)
	tagValue := "TAG_VALUE"
	newTagValue := "NEW_TAG_VALUE"

	tagModel := model.TagBase("test", tagId)
	tagAssociationModel := model.TagAssociation("test", []sdk.ObjectIdentifier{table.ID()}, string(sdk.ObjectTypeTable), tagId.FullyQualifiedName(), tagValue).
		WithDependsOn(tagModel.ResourceReference())
	newTagAssociationModel := model.TagAssociation("test", []sdk.ObjectIdentifier{table.ID()}, string(sdk.ObjectTypeTable), tagId.FullyQualifiedName(), newTagValue).
		WithDependsOn(tagModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tagModel, tagAssociationModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), tagValue, string(sdk.ObjectTypeTable))),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_value", tagValue),
				),
			},
			{
				// Couldn't assert it, but with `ExternalProviders: acc.ExternalProviderWithExactVersion("2.0.0")`,
				// the plugin crashes with the same panic as in https://github.com/snowflakedb/terraform-provider-snowflake/issues/3622.
				// The new version handles the SchemaObjectIdentifier correctly, so the panic is not triggered.
				Config: accconfig.FromModels(t, tagModel, newTagAssociationModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), newTagValue, string(sdk.ObjectTypeTable))),
					resource.TestCheckResourceAttr(tagAssociationModel.ResourceReference(), "tag_value", newTagValue),
				),
			},
		},
	})
}
