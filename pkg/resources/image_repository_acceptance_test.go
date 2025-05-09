//go:build !account_level_tests

package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ImageRepository_basic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	// TODO(SNOW-2070746): We set up a separate database and schema with capitalized ids. Remove this after fix on snowflake side.
	db, dbCleanup := acc.TestClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	schema, schemaCleanup := acc.TestClient().Schema.CreateSchemaInDatabase(t, db.ID())
	t.Cleanup(schemaCleanup)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	comment := random.Comment()
	changedComment := random.Comment()

	imageRepositoryModelBasic := model.ImageRepository("test", id.DatabaseName(), id.SchemaName(), id.Name())
	imageRepositoryModelWithComment := model.ImageRepository("test", id.DatabaseName(), id.SchemaName(), id.Name()).WithComment(comment)
	imageRepositoryModelWithChangedComment := model.ImageRepository("test", id.DatabaseName(), id.SchemaName(), id.Name()).WithComment(changedComment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ImageRepository),
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
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ImageRepositoryShowOutput(t, imageRepositoryModelBasic.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasRepositoryUrlNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment("").
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
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ImageRepositoryShowOutput(t, imageRepositoryModelBasic.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasRepositoryUrlNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(comment).
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
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ImageRepositoryShowOutput(t, imageRepositoryModelWithChangedComment.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasRepositoryUrlNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(changedComment).
						HasPrivatelinkRepositoryUrl(""),
				),
			},
			// change externally
			{
				PreConfig: func() {
					acc.TestClient().ImageRepository.Alter(t, sdk.NewAlterImageRepositoryRequest(id).WithSet(
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
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ImageRepositoryShowOutput(t, imageRepositoryModelWithChangedComment.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasRepositoryUrlNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(changedComment).
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
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ImageRepositoryShowOutput(t, imageRepositoryModelBasic.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasRepositoryUrlNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment("").
						HasPrivatelinkRepositoryUrl(""),
				),
			},
		},
	})
}

func TestAcc_ImageRepository_complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	// TODO(SNOW-2070746): We set up a separate database and schema with capitalized ids. Remove this after fix on snowflake side.
	db, dbCleanup := acc.TestClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	schema, schemaCleanup := acc.TestClient().Schema.CreateSchemaInDatabase(t, db.ID())
	t.Cleanup(schemaCleanup)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	comment := random.Comment()

	modelComplete := model.ImageRepository("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ImageRepository),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.ImageRepositoryResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString(comment).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ImageRepositoryShowOutput(t, modelComplete.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasRepositoryUrlNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(comment).
						HasPrivatelinkRepositoryUrl(""),
				),
			},
			{
				Config:            accconfig.FromModels(t, modelComplete),
				ResourceName:      modelComplete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
