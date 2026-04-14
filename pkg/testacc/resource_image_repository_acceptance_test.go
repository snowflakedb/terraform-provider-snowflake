//go:build non_account_level_tests

package testacc

import (
	"strings"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ImageRepository_basic(t *testing.T) {
	// TODO(SNOW-2070746): We set up a separate database and schema with capitalized ids. Remove this after fix on snowflake side.
	db, dbCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, db.ID())
	t.Cleanup(schemaCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	comment := random.Comment()
	changedComment := random.Comment()

	imageRepositoryModelBasic := model.ImageRepository("test", id.DatabaseName(), id.SchemaName(), id.Name())
	imageRepositoryModelWithComment := model.ImageRepository("test", id.DatabaseName(), id.SchemaName(), id.Name()).WithComment(comment)
	imageRepositoryModelWithChangedComment := model.ImageRepository("test", id.DatabaseName(), id.SchemaName(), id.Name()).WithComment(changedComment)
	imageRepositoryModelWithEncryption := model.ImageRepository("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithEncryptionEnum(sdk.ImageRepositoryEncryptionTypeSnowflakeFull)
	imageRepositoryModelWithDifferentEncryption := model.ImageRepository("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithEncryptionEnum(sdk.ImageRepositoryEncryptionTypeSnowflakeSse)
	imageRepositoryModelWithDifferentEncryptionLowercase := model.ImageRepository("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithEncryption(strings.ToLower(string(sdk.ImageRepositoryEncryptionTypeSnowflakeSse)))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ImageRepository),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				Config: accconfig.FromModels(t, imageRepositoryModelBasic),
				Check: assertThat(t,
					resourceassert.ImageRepositoryResource(t, imageRepositoryModelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString("").
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasNoEncryption(),
					resourceshowoutputassert.ImageRepositoryShowOutput(t, imageRepositoryModelBasic.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasRepositoryUrlNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment("").
						HasEncryption(sdk.ImageRepositoryEncryptionTypeSnowflakeFull).
						HasPrivatelinkRepositoryUrl(""),
				),
			},
			// import - without optionals
			{
				Config:            accconfig.FromModels(t, imageRepositoryModelBasic),
				ResourceName:      imageRepositoryModelBasic.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// set optionals
			{
				Config: accconfig.FromModels(t, imageRepositoryModelWithComment),
				Check: assertThat(t,
					resourceassert.ImageRepositoryResource(t, imageRepositoryModelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString(comment).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasNoEncryption(),
					resourceshowoutputassert.ImageRepositoryShowOutput(t, imageRepositoryModelBasic.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasRepositoryUrlNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(comment).
						HasEncryption(sdk.ImageRepositoryEncryptionTypeSnowflakeFull).
						HasPrivatelinkRepositoryUrl(""),
				),
			},
			// import - complete
			{
				Config:            accconfig.FromModels(t, imageRepositoryModelWithComment),
				ResourceName:      imageRepositoryModelWithComment.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// alter
			{
				Config: accconfig.FromModels(t, imageRepositoryModelWithChangedComment),
				Check: assertThat(t,
					resourceassert.ImageRepositoryResource(t, imageRepositoryModelWithChangedComment.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString(changedComment).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasNoEncryption(),
					resourceshowoutputassert.ImageRepositoryShowOutput(t, imageRepositoryModelWithChangedComment.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasRepositoryUrlNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(changedComment).
						HasEncryption(sdk.ImageRepositoryEncryptionTypeSnowflakeFull).
						HasPrivatelinkRepositoryUrl(""),
				),
			},
			// change externally
			{
				PreConfig: func() {
					testClient().ImageRepository.Alter(t, sdk.NewAlterImageRepositoryRequest(id).WithSet(
						*sdk.NewImageRepositorySetRequest().
							WithComment(sdk.StringAllowEmpty{Value: comment}),
					))
				},
				Config: accconfig.FromModels(t, imageRepositoryModelWithChangedComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(imageRepositoryModelWithChangedComment.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(imageRepositoryModelWithChangedComment.ResourceReference(), "comment", sdk.Pointer(changedComment), sdk.Pointer(comment)),
						planchecks.ExpectChange(imageRepositoryModelWithChangedComment.ResourceReference(), "comment", tfjson.ActionUpdate, sdk.Pointer(comment), sdk.Pointer(changedComment)),
					},
				},
				Check: assertThat(t,
					resourceassert.ImageRepositoryResource(t, imageRepositoryModelWithChangedComment.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString(changedComment).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasNoEncryption(),
					resourceshowoutputassert.ImageRepositoryShowOutput(t, imageRepositoryModelWithChangedComment.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasRepositoryUrlNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(changedComment).
						HasEncryption(sdk.ImageRepositoryEncryptionTypeSnowflakeFull).
						HasPrivatelinkRepositoryUrl(""),
				),
			},
			// unset
			{
				Config: accconfig.FromModels(t, imageRepositoryModelBasic),
				Check: assertThat(t,
					resourceassert.ImageRepositoryResource(t, imageRepositoryModelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString("").
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasNoEncryption(),
					resourceshowoutputassert.ImageRepositoryShowOutput(t, imageRepositoryModelBasic.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasRepositoryUrlNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment("").
						HasEncryption(sdk.ImageRepositoryEncryptionTypeSnowflakeFull).
						HasPrivatelinkRepositoryUrl(""),
				),
			},
			// change encryption externally (force new)
			{
				PreConfig: func() {
					testClient().ImageRepository.DropImageRepositoryFunc(t, id)()
					_, imageRepositoryCleanup := testClient().ImageRepository.CreateWithRequest(t, sdk.NewCreateImageRepositoryRequest(id).WithEncryption(*sdk.NewImageRepositoryEncryptionRequest(sdk.ImageRepositoryEncryptionTypeSnowflakeSse)))
					t.Cleanup(imageRepositoryCleanup)
				},
				Config: accconfig.FromModels(t, imageRepositoryModelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(imageRepositoryModelWithChangedComment.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.ImageRepositoryResource(t, imageRepositoryModelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasNoEncryption(),
					resourceshowoutputassert.ImageRepositoryShowOutput(t, imageRepositoryModelBasic.ResourceReference()).
						HasName(id.Name()).
						HasEncryption(sdk.ImageRepositoryEncryptionTypeSnowflakeFull),
				),
			},
			// set encryption to the current Snowflake value (expect no-op)
			{
				Config: accconfig.FromModels(t, imageRepositoryModelWithEncryption),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(imageRepositoryModelWithChangedComment.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(t,
					resourceassert.ImageRepositoryResource(t, imageRepositoryModelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasNoEncryption(),
					resourceshowoutputassert.ImageRepositoryShowOutput(t, imageRepositoryModelBasic.ResourceReference()).
						HasName(id.Name()).
						HasEncryption(sdk.ImageRepositoryEncryptionTypeSnowflakeFull),
				),
			},
			// set encryption to a different value (expect drop and recreate)
			{
				Config: accconfig.FromModels(t, imageRepositoryModelWithDifferentEncryption),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(imageRepositoryModelWithChangedComment.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.ImageRepositoryResource(t, imageRepositoryModelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasEncryptionString(string(sdk.ImageRepositoryEncryptionTypeSnowflakeSse)),
					resourceshowoutputassert.ImageRepositoryShowOutput(t, imageRepositoryModelBasic.ResourceReference()).
						HasName(id.Name()).
						HasEncryption(sdk.ImageRepositoryEncryptionTypeSnowflakeSse),
				),
			},
			// set encryption to current value lowercase (expect no-op)
			{
				Config: accconfig.FromModels(t, imageRepositoryModelWithDifferentEncryptionLowercase),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(imageRepositoryModelWithChangedComment.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(t,
					resourceassert.ImageRepositoryResource(t, imageRepositoryModelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasEncryptionString(string(sdk.ImageRepositoryEncryptionTypeSnowflakeSse)),
					resourceshowoutputassert.ImageRepositoryShowOutput(t, imageRepositoryModelBasic.ResourceReference()).
						HasName(id.Name()).
						HasEncryption(sdk.ImageRepositoryEncryptionTypeSnowflakeSse),
				),
			},
		},
	})
}

func TestAcc_ImageRepository_complete(t *testing.T) {
	// TODO(SNOW-2070746): We set up a separate database and schema with capitalized ids. Remove this after fix on snowflake side.
	db, dbCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, db.ID())
	t.Cleanup(schemaCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	comment := random.Comment()

	modelComplete := model.ImageRepository("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithComment(comment).
		WithEncryptionEnum(sdk.ImageRepositoryEncryptionTypeSnowflakeSse)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ImageRepository),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.ImageRepositoryResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString(comment).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasEncryptionString(string(sdk.ImageRepositoryEncryptionTypeSnowflakeSse)),
					resourceshowoutputassert.ImageRepositoryShowOutput(t, modelComplete.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasRepositoryUrlNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(comment).
						HasEncryption(sdk.ImageRepositoryEncryptionTypeSnowflakeSse).
						HasPrivatelinkRepositoryUrl(""),
				),
			},
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption"},
			},
		},
	})
}

func TestAcc_ImageRepository_importWithoutEncryptionSet(t *testing.T) {
	// TODO(SNOW-2070746): We set up a separate database and schema with capitalized ids. Remove this after fix on snowflake side.
	db, dbCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, db.ID())
	t.Cleanup(schemaCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	imageRepositoryModel := model.ImageRepository("test", id.DatabaseName(), id.SchemaName(), id.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ImageRepository),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, imageRepositoryCleanup := testClient().ImageRepository.CreateWithRequest(t, sdk.NewCreateImageRepositoryRequest(id))
					t.Cleanup(imageRepositoryCleanup)
				},
				Config:             accconfig.FromModels(t, imageRepositoryModel),
				ResourceName:       imageRepositoryModel.ResourceReference(),
				ImportState:        true,
				ImportStateId:      id.FullyQualifiedName(),
				ImportStatePersist: true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedImageRepositoryResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNoEncryption(),
					resourceshowoutputassert.ImportedImageRepositoryShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasEncryption(sdk.ImageRepositoryEncryptionTypeSnowflakeFull),
				),
			},
			// Plan to verify no diff
			{
				Config: accconfig.FromModels(t, imageRepositoryModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(imageRepositoryModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

func TestAcc_ImageRepository_importWithEncryptionSetToSnowflakeValue(t *testing.T) {
	// TODO(SNOW-2070746): We set up a separate database and schema with capitalized ids. Remove this after fix on snowflake side.
	db, dbCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, db.ID())
	t.Cleanup(schemaCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	imageRepositoryModel := model.ImageRepository("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithEncryptionEnum(sdk.ImageRepositoryEncryptionTypeSnowflakeFull)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ImageRepository),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, imageRepositoryCleanup := testClient().ImageRepository.CreateWithRequest(t, sdk.NewCreateImageRepositoryRequest(id))
					t.Cleanup(imageRepositoryCleanup)
				},
				Config:             accconfig.FromModels(t, imageRepositoryModel),
				ResourceName:       imageRepositoryModel.ResourceReference(),
				ImportState:        true,
				ImportStateId:      id.FullyQualifiedName(),
				ImportStatePersist: true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedImageRepositoryResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNoEncryption(),
					resourceshowoutputassert.ImportedImageRepositoryShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasEncryption(sdk.ImageRepositoryEncryptionTypeSnowflakeFull),
				),
			},
			// Plan to verify no diff
			{
				Config: accconfig.FromModels(t, imageRepositoryModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(imageRepositoryModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

func TestAcc_ImageRepository_importWithEncryptionSetToDifferentValue(t *testing.T) {
	// TODO(SNOW-2070746): We set up a separate database and schema with capitalized ids. Remove this after fix on snowflake side.
	db, dbCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, db.ID())
	t.Cleanup(schemaCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	imageRepositoryModel := model.ImageRepository("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithEncryptionEnum(sdk.ImageRepositoryEncryptionTypeSnowflakeSse)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ImageRepository),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, imageRepositoryCleanup := testClient().ImageRepository.CreateWithRequest(t, sdk.NewCreateImageRepositoryRequest(id))
					t.Cleanup(imageRepositoryCleanup)
				},
				Config:             accconfig.FromModels(t, imageRepositoryModel),
				ResourceName:       imageRepositoryModel.ResourceReference(),
				ImportState:        true,
				ImportStateId:      id.FullyQualifiedName(),
				ImportStatePersist: true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedImageRepositoryResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNoEncryption(),
					resourceshowoutputassert.ImportedImageRepositoryShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasEncryption(sdk.ImageRepositoryEncryptionTypeSnowflakeFull),
				),
			},
			// Plan to verify drop and recreate need
			{
				Config: accconfig.FromModels(t, imageRepositoryModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(imageRepositoryModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
		},
	})
}

func TestAcc_ImageRepository_migrateFromV2_14_1(t *testing.T) {
	// TODO(SNOW-2070746): We set up a separate database and schema with capitalized ids. Remove this after fix on snowflake side.
	db, dbCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, db.ID())
	t.Cleanup(schemaCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	imageRepositoryModel := model.ImageRepository("test", id.DatabaseName(), id.SchemaName(), id.Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ImageRepository),
		Steps: []resource.TestStep{
			// create with old provider (no encryption field)
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.14.1"),
				Config:            accconfig.FromModels(t, imageRepositoryModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(imageRepositoryModel.ResourceReference(), "name", id.Name())),
				),
			},
			// upgrade to current provider - encryption field not in config, no-op expected
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, imageRepositoryModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(imageRepositoryModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(t,
					resourceassert.ImageRepositoryResource(t, imageRepositoryModel.ResourceReference()).
						HasNameString(id.Name()).
						HasNoEncryption(),
					resourceshowoutputassert.ImageRepositoryShowOutput(t, imageRepositoryModel.ResourceReference()).
						HasEncryption(sdk.ImageRepositoryEncryptionTypeSnowflakeFull),
				),
			},
		},
	})
}
