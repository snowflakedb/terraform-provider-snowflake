//go:build !account_level_tests

package testacc

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GitRepository_basic(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment, changedComment := random.Comment(), random.Comment()
	origin := testvars.ExampleGitRepositoryOrigin

	apiIntegrationId, apiIntegrationCleanup := testClient().ApiIntegration.CreateApiIntegrationForGitRepository(t, origin)
	t.Cleanup(apiIntegrationCleanup)
	changedApiIntegrationId, changedApiIntegrationCleanup := testClient().ApiIntegration.CreateApiIntegrationForGitRepository(t, origin)
	t.Cleanup(changedApiIntegrationCleanup)

	secretId, secretCleanup := testClient().Secret.CreateRandomPasswordSecret(t)
	t.Cleanup(secretCleanup)

	changedSecretId, changedSecretCleanup := testClient().Secret.CreateRandomPasswordSecret(t)
	t.Cleanup(changedSecretCleanup)

	modelBasic := model.GitRepository(
		"test",
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		apiIntegrationId.FullyQualifiedName(),
		origin,
	)

	modelComplete := model.GitRepository(
		"test",
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		apiIntegrationId.FullyQualifiedName(),
		origin,
	).WithGitCredentials(secretId.FullyQualifiedName()).
		WithComment(comment)

	modelCompleteWithDifferentValues := model.GitRepository(
		"test",
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		changedApiIntegrationId.FullyQualifiedName(),
		origin,
	).WithGitCredentials(changedSecretId.FullyQualifiedName()).
		WithComment(changedComment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.GitRepository),

		Steps: []resource.TestStep{
			// create with only required attributes
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.GitRepositoryResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasOriginString(origin).
						HasApiIntegrationString(apiIntegrationId.FullyQualifiedName()).
						HasGitCredentialsString("").
						HasCommentString("").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.GitRepositoryShowOutput(t, modelBasic.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOrigin(origin).
						HasApiIntegration(apiIntegrationId).
						HasGitCredentialsEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(""),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.origin", origin)),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.api_integration", apiIntegrationId.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.git_credentials", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.comment", "")),
				),
			},
			// import minimal state
			{
				Config:       accconfig.FromModels(t, modelBasic),
				ResourceName: modelBasic.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedGitRepositoryResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasOriginString(origin).
						HasApiIntegrationString(apiIntegrationId.FullyQualifiedName()).
						HasGitCredentialsString("").
						HasCommentString("").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ImportedGitRepositoryShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOrigin(origin).
						HasApiIntegration(apiIntegrationId).
						HasGitCredentialsEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(""),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateSet(helpers.EncodeResourceIdentifier(id), "describe_output.0.created_on")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.name", id.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.database_name", id.DatabaseName())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.schema_name", id.SchemaName())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.origin", origin)),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.api_integration", apiIntegrationId.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.git_credentials", "")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.owner_role_type", "ROLE")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.comment", "")),
				),
			},
			// add optional attributes
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.GitRepositoryResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasOriginString(origin).
						HasApiIntegrationString(apiIntegrationId.FullyQualifiedName()).
						HasGitCredentialsString(secretId.FullyQualifiedName()).
						HasCommentString(comment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.GitRepositoryShowOutput(t, modelComplete.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOrigin(origin).
						HasApiIntegration(apiIntegrationId).
						HasGitCredentials(secretId).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(comment),
				),
			},
			// import complete state
			{
				Config:            accconfig.FromModels(t, modelComplete),
				ResourceName:      modelComplete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// alter
			{
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.GitRepositoryResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasOriginString(origin).
						HasApiIntegrationString(changedApiIntegrationId.FullyQualifiedName()).
						HasGitCredentialsString(changedSecretId.FullyQualifiedName()).
						HasCommentString(changedComment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.GitRepositoryShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOrigin(origin).
						HasApiIntegration(changedApiIntegrationId).
						HasGitCredentials(changedSecretId).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(changedComment),
				),
			},
			// change externally alter
			{
				PreConfig: func() {
					testClient().GitRepository.Alter(t, sdk.NewAlterGitRepositoryRequest(id).WithSet(
						*sdk.NewGitRepositorySetRequest().
							WithApiIntegration(apiIntegrationId).
							WithGitCredentials(secretId).
							WithComment(comment)))
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.GitRepositoryResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasOriginString(origin).
						HasApiIntegrationString(changedApiIntegrationId.FullyQualifiedName()).
						HasGitCredentialsString(changedSecretId.FullyQualifiedName()).
						HasCommentString(changedComment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.GitRepositoryShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOrigin(origin).
						HasApiIntegration(changedApiIntegrationId).
						HasGitCredentials(changedSecretId).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(changedComment),
				),
			},
			// change externally create
			{
				PreConfig: func() {
					newOrigin := "https://github.com/octocat/git-consortium"
					newApiIntegrationId, newApiIntegrationCleanup := testClient().ApiIntegration.CreateApiIntegrationForGitRepository(t, newOrigin)
					t.Cleanup(newApiIntegrationCleanup)
					testClient().GitRepository.CreateWithRequest(t, sdk.NewCreateGitRepositoryRequest(id, newOrigin, newApiIntegrationId).
						WithGitCredentials(secretId).WithComment(comment).WithOrReplace(true))
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.GitRepositoryResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasOriginString(origin).
						HasApiIntegrationString(changedApiIntegrationId.FullyQualifiedName()).
						HasGitCredentialsString(changedSecretId.FullyQualifiedName()).
						HasCommentString(changedComment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.GitRepositoryShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOrigin(origin).
						HasApiIntegration(changedApiIntegrationId).
						HasGitCredentials(changedSecretId).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(changedComment),
				),
			},
			// unset
			{
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.GitRepositoryResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasOriginString(origin).
						HasApiIntegrationString(apiIntegrationId.FullyQualifiedName()).
						HasGitCredentialsString("").
						HasCommentString("").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.GitRepositoryShowOutput(t, modelBasic.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOrigin(origin).
						HasApiIntegration(apiIntegrationId).
						HasGitCredentialsEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(""),
				),
			},
		},
	})
}

func TestAcc_GitRepository_complete(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	origin := testvars.ExampleGitRepositoryOrigin

	apiIntegrationId, apiIntegrationCleanup := testClient().ApiIntegration.CreateApiIntegrationForGitRepository(t, origin)
	t.Cleanup(apiIntegrationCleanup)

	secretId, secretCleanup := testClient().Secret.CreateRandomPasswordSecret(t)
	t.Cleanup(secretCleanup)

	modelComplete := model.
		GitRepository("test", id.DatabaseName(), id.SchemaName(), id.Name(), apiIntegrationId.FullyQualifiedName(), origin).
		WithGitCredentials(secretId.FullyQualifiedName()).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.GitRepository),
		Steps: []resource.TestStep{
			// create with all attributes set
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.GitRepositoryResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasOriginString(origin).
						HasApiIntegrationString(apiIntegrationId.FullyQualifiedName()).
						HasGitCredentialsString(secretId.FullyQualifiedName()).
						HasCommentString(comment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.GitRepositoryShowOutput(t, modelComplete.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOrigin(origin).
						HasApiIntegration(apiIntegrationId).
						HasGitCredentials(secretId).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(comment),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.origin", origin)),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.api_integration", apiIntegrationId.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.git_credentials", secretId.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.comment", comment)),
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
