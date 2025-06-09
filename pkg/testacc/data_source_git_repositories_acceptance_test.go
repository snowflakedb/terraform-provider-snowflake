//go:build !account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GitRepositories(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	origin := "https://github.com/octocat/hello-world"
	comment := random.Comment()

	apiIntegrationId, apiCleanup := testClient().ApiIntegration.CreateApiIntegrationForGitRepository(t, origin)
	t.Cleanup(apiCleanup)

	secretId := testClient().Ids.RandomSchemaObjectIdentifier()
	_, secretCleanup := testClient().Secret.CreateWithBasicAuthenticationFlow(t, secretId, "username", "password")
	t.Cleanup(secretCleanup)

	gitRepositoryModel := model.
		GitRepository("test", id.DatabaseName(), id.SchemaName(), id.Name(), apiIntegrationId.FullyQualifiedName(), origin).
		WithGitCredentials(secretId.FullyQualifiedName()).
		WithComment(comment)

	dataSourceModel := datasourcemodel.GitRepositories("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(gitRepositoryModel.ResourceReference())

	dataSourceWithoutOptionals := datasourcemodel.
		GitRepositories("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithDependsOn(gitRepositoryModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, gitRepositoryModel, dataSourceModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "git_repositories.#", "1")),
					resourceshowoutputassert.GitRepositoriesDatasourceShowOutput(t, "snowflake_git_repositories.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOrigin(origin).
						HasApiIntegration(apiIntegrationId).
						HasGitCredentials(secretId).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment),
					assert.Check(resource.TestCheckResourceAttrSet(dataSourceModel.DatasourceReference(), "git_repositories.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "git_repositories.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "git_repositories.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "git_repositories.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "git_repositories.0.describe_output.0.origin", origin)),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "git_repositories.0.describe_output.0.api_integration", apiIntegrationId.Name())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "git_repositories.0.describe_output.0.git_credentials", secretId.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "git_repositories.0.describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "git_repositories.0.describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "git_repositories.0.describe_output.0.comment", comment)),
				),
			},
			{
				Config: accconfig.FromModels(t, gitRepositoryModel, dataSourceWithoutOptionals),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dataSourceWithoutOptionals.DatasourceReference(), "git_repositories.#", "1")),
					resourceshowoutputassert.GitRepositoriesDatasourceShowOutput(t, "snowflake_git_repositories.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOrigin(origin).
						HasApiIntegration(apiIntegrationId).
						HasGitCredentials(secretId).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment),
					assert.Check(resource.TestCheckResourceAttr(dataSourceWithoutOptionals.DatasourceReference(), "git_repositories.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_GitRepositories_Filtering(t *testing.T) {
	id1 := testClient().Ids.RandomSchemaObjectIdentifier()
	id2 := testClient().Ids.RandomSchemaObjectIdentifier()
	id3 := testClient().Ids.RandomSchemaObjectIdentifier()

	origin := "https://github.com/octocat/hello-world"
	apiIntegrationId, apiCleanup := testClient().ApiIntegration.CreateApiIntegrationForGitRepository(t, origin)
	t.Cleanup(apiCleanup)

	secretID := testClient().Ids.RandomSchemaObjectIdentifier()
	_, secretCleanup := testClient().Secret.CreateWithBasicAuthenticationFlow(t, secretID, "u", "p")
	t.Cleanup(secretCleanup)

	gitRepositoryModel1 := model.GitRepository("test1", id1.DatabaseName(), id1.SchemaName(), id1.Name(), apiIntegrationId.FullyQualifiedName(), origin).
		WithGitCredentials(secretID.FullyQualifiedName())
	gitRepositoryModel2 := model.GitRepository("test2", id2.DatabaseName(), id2.SchemaName(), id2.Name(), apiIntegrationId.FullyQualifiedName(), origin)
	gitRepositoryModel3 := model.GitRepository("test3", id3.DatabaseName(), id3.SchemaName(), id3.Name(), apiIntegrationId.FullyQualifiedName(), origin)

	gitRepositoriesWithLikeModel := datasourcemodel.GitRepositories("test").
		WithLike(id1.Name()).
		WithDependsOn(gitRepositoryModel1.ResourceReference(), gitRepositoryModel2.ResourceReference(), gitRepositoryModel3.ResourceReference())

	gitRepositoriesWithInModel := datasourcemodel.GitRepositories("test").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(gitRepositoryModel1.ResourceReference(), gitRepositoryModel2.ResourceReference(), gitRepositoryModel3.ResourceReference())

	gitRepositoriesWithLimitModel := datasourcemodel.GitRepositories("test").
		WithRowsAndFrom(2, "").
		WithDependsOn(gitRepositoryModel1.ResourceReference(), gitRepositoryModel2.ResourceReference(), gitRepositoryModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.GitRepository),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, gitRepositoryModel1, gitRepositoryModel2, gitRepositoryModel3, gitRepositoriesWithLikeModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(gitRepositoriesWithLikeModel.DatasourceReference(), "git_repositories.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, gitRepositoryModel1, gitRepositoryModel2, gitRepositoryModel3, gitRepositoriesWithInModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(gitRepositoriesWithInModel.DatasourceReference(), "git_repositories.#", "3"),
				),
			},
			{
				Config: accconfig.FromModels(t, gitRepositoryModel1, gitRepositoryModel2, gitRepositoryModel3, gitRepositoriesWithLimitModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(gitRepositoriesWithLimitModel.DatasourceReference(), "git_repositories.#", "2"),
				),
			},
		},
	})
}

func TestAcc_GitRepositories_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.GitRepositories("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_GitRepositories_NotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_GitRepositories/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one git repository"),
			},
		},
	})
}
