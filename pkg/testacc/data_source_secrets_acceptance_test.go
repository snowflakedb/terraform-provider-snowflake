//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Secrets_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	secretId1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	secretId2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	secretId3 := testClient().Ids.RandomSchemaObjectIdentifier()

	secretModel1 := model.SecretWithGenericString("test", secretId1.DatabaseName(), secretId1.SchemaName(), secretId1.Name(), "test_secret_string1")
	secretModel2 := model.SecretWithGenericString("test1", secretId2.DatabaseName(), secretId2.SchemaName(), secretId2.Name(), "test_secret_string2")
	secretModel3 := model.SecretWithGenericString("test2", secretId3.DatabaseName(), secretId3.SchemaName(), secretId3.Name(), "test_secret_string3")

	datasourceModelLikeExact := datasourcemodel.Secrets("test").
		WithWithDescribe(false).
		WithLike(secretId1.Name()).
		WithInDatabase(secretId1.DatabaseId()).
		WithDependsOn(secretModel1.ResourceReference(), secretModel2.ResourceReference(), secretModel3.ResourceReference())

	datasourceModelLikePrefix := datasourcemodel.Secrets("test").
		WithWithDescribe(false).
		WithLike(prefix+"%").
		WithInDatabase(secretId1.DatabaseId()).
		WithDependsOn(secretModel1.ResourceReference(), secretModel2.ResourceReference(), secretModel3.ResourceReference())

	datasourceModelInDatabase := datasourcemodel.Secrets("test").
		WithWithDescribe(false).
		WithInDatabase(secretId1.DatabaseId()).
		WithDependsOn(secretModel1.ResourceReference(), secretModel2.ResourceReference(), secretModel3.ResourceReference())

	datasourceModelInSchema := datasourcemodel.Secrets("test").
		WithWithDescribe(false).
		WithInSchema(secretId1.SchemaId()).
		WithDependsOn(secretModel1.ResourceReference(), secretModel2.ResourceReference(), secretModel3.ResourceReference())

	datasourceModelLikeInDatabase := datasourcemodel.Secrets("test").
		WithWithDescribe(false).
		WithLike(prefix+"%").
		WithInDatabase(secretId1.DatabaseId()).
		WithDependsOn(secretModel1.ResourceReference(), secretModel2.ResourceReference(), secretModel3.ResourceReference())

	datasourceModelLikeInSchema := datasourcemodel.Secrets("test").
		WithWithDescribe(false).
		WithLike(prefix+"%").
		WithInSchema(secretId1.SchemaId()).
		WithDependsOn(secretModel1.ResourceReference(), secretModel2.ResourceReference(), secretModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecretWithGenericString),
		Steps: []resource.TestStep{
			// like (exact)
			{
				Config: config.FromModels(t, secretModel1, secretModel2, secretModel3, datasourceModelLikeExact),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikeExact.DatasourceReference(), "secrets.#", "1"),
					resource.TestCheckResourceAttr(datasourceModelLikeExact.DatasourceReference(), "secrets.0.show_output.0.name", secretId1.Name()),
				),
			},
			// like (prefix)
			{
				Config: config.FromModels(t, secretModel1, secretModel2, secretModel3, datasourceModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikePrefix.DatasourceReference(), "secrets.#", "2"),
				),
			},
			// in database
			{
				Config: config.FromModels(t, secretModel1, secretModel2, secretModel3, datasourceModelInDatabase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelInDatabase.DatasourceReference(), "secrets.#", "3"),
				),
			},
			// in schema
			{
				Config: config.FromModels(t, secretModel1, secretModel2, secretModel3, datasourceModelInSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelInSchema.DatasourceReference(), "secrets.#", "3"),
				),
			},
			// like + in database
			{
				Config: config.FromModels(t, secretModel1, secretModel2, secretModel3, datasourceModelLikeInDatabase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikeInDatabase.DatasourceReference(), "secrets.#", "2"),
				),
			},
			// like + in schema
			{
				Config: config.FromModels(t, secretModel1, secretModel2, secretModel3, datasourceModelLikeInSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikeInSchema.DatasourceReference(), "secrets.#", "2"),
				),
			},
		},
	})
}

func TestAcc_Secrets_CompleteUseCase(t *testing.T) {
	prefix := random.AlphaN(6)
	comment := random.Comment()

	apiIntegrationId := testClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := testClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(apiIntegrationId, true, "test_oauth_client_id", "test_oauth_client_secret").
			WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "scope1"}, {Scope: "scope2"}}),
	)
	t.Cleanup(apiIntegrationCleanup)

	genericId := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix + "_gen")
	basicId := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix + "_bas")
	clientCredsId := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix + "_cli")
	authCodeId := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix + "_aut")

	genericModel := model.SecretWithGenericString("gen", genericId.DatabaseName(), genericId.SchemaName(), genericId.Name(), "generic_value").WithComment(comment)
	basicModel := model.SecretWithBasicAuthentication("bas", basicId.DatabaseName(), basicId.SchemaName(), basicId.Name(), "pwd", "user1")
	clientCredsModel := model.SecretWithClientCredentials("cli", clientCredsId.DatabaseName(), clientCredsId.SchemaName(), clientCredsId.Name(), apiIntegrationId.Name(), []string{"scope1", "scope2"})
	authCodeModel := model.SecretWithAuthorizationCodeGrant("aut", authCodeId.DatabaseName(), authCodeId.SchemaName(), authCodeId.Name(), apiIntegrationId.Name(), "refresh_token_value", time.Now().Add(24*time.Hour).Format(time.DateTime)).WithComment(comment)

	genericSecretNoDescribe := datasourcemodel.Secrets("test").
		WithLike(genericId.Name()).
		WithInDatabase(genericId.DatabaseId()).
		WithWithDescribe(false).
		WithDependsOn(genericModel.ResourceReference())

	genericSecretWithDescribe := datasourcemodel.Secrets("test").
		WithLike(genericId.Name()).
		WithInDatabase(genericId.DatabaseId()).
		WithWithDescribe(true).
		WithDependsOn(genericModel.ResourceReference())

	basicSecretWithDescribe := datasourcemodel.Secrets("test").
		WithLike(basicId.Name()).
		WithInDatabase(basicId.DatabaseId()).
		WithWithDescribe(true).
		WithDependsOn(basicModel.ResourceReference())

	clientCredsSecretWithDescribe := datasourcemodel.Secrets("test").
		WithLike(clientCredsId.Name()).
		WithInDatabase(clientCredsId.DatabaseId()).
		WithWithDescribe(true).
		WithDependsOn(clientCredsModel.ResourceReference())

	authCodeSecretWithDescribe := datasourcemodel.Secrets("test").
		WithLike(authCodeId.Name()).
		WithInDatabase(authCodeId.DatabaseId()).
		WithWithDescribe(true).
		WithDependsOn(authCodeModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: ComposeCheckDestroy(t,
			resources.SecretWithGenericString,
			resources.SecretWithBasicAuthentication,
			resources.SecretWithClientCredentials,
			resources.SecretWithAuthorizationCodeGrant,
		),
		Steps: []resource.TestStep{
			// Generic string without describe
			{
				Config: config.FromModels(t, genericModel, genericSecretNoDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.SecretsDatasourceShowOutput(t, genericSecretNoDescribe.DatasourceReference()).
						HasName(genericId.Name()).
						HasDatabaseName(genericId.DatabaseName()).
						HasSchemaName(genericId.SchemaName()).
						HasComment(comment).
						HasSecretType(string(sdk.SecretTypeGenericString)),
					assert.Check(resource.TestCheckResourceAttr(genericSecretNoDescribe.DatasourceReference(), "secrets.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(genericSecretNoDescribe.DatasourceReference(), "secrets.0.describe_output.#", "0")),
				),
			},
			// Generic string with describe
			{
				Config: config.FromModels(t, genericModel, genericSecretWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.SecretsDatasourceShowOutput(t, genericSecretWithDescribe.DatasourceReference()).
						HasName(genericId.Name()).
						HasDatabaseName(genericId.DatabaseName()).
						HasSchemaName(genericId.SchemaName()).
						HasComment(comment).
						HasSecretType(string(sdk.SecretTypeGenericString)),
					assert.Check(resource.TestCheckResourceAttr(genericSecretWithDescribe.DatasourceReference(), "secrets.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(genericSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(genericSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.0.name", genericId.Name())),
					assert.Check(resource.TestCheckResourceAttr(genericSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.0.database_name", genericId.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(genericSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.0.schema_name", genericId.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(genericSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.0.secret_type", string(sdk.SecretTypeGenericString))),
					assert.Check(resource.TestCheckResourceAttr(genericSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.0.username", "")),
					assert.Check(resource.TestCheckResourceAttr(genericSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(genericSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.0.oauth_scopes.#", "0")),
				),
			},
			// Basic authentication with describe
			{
				Config: config.FromModels(t, basicModel, basicSecretWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.SecretsDatasourceShowOutput(t, basicSecretWithDescribe.DatasourceReference()).
						HasName(basicId.Name()).
						HasDatabaseName(basicId.DatabaseName()).
						HasSchemaName(basicId.SchemaName()).
						HasComment("").
						HasSecretType(string(sdk.SecretTypePassword)),
					assert.Check(resource.TestCheckResourceAttr(basicSecretWithDescribe.DatasourceReference(), "secrets.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(basicSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.0.secret_type", string(sdk.SecretTypePassword))),
					assert.Check(resource.TestCheckResourceAttr(basicSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.0.username", "user1")),
				),
			},
			// OAuth2 client credentials with describe
			{
				Config: config.FromModels(t, clientCredsModel, clientCredsSecretWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.SecretsDatasourceShowOutput(t, clientCredsSecretWithDescribe.DatasourceReference()).
						HasName(clientCredsId.Name()).
						HasDatabaseName(clientCredsId.DatabaseName()).
						HasSchemaName(clientCredsId.SchemaName()).
						HasComment("").
						HasSecretType(string(sdk.SecretTypeOAuth2)),
					assert.Check(resource.TestCheckResourceAttr(clientCredsSecretWithDescribe.DatasourceReference(), "secrets.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(clientCredsSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(clientCredsSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.0.secret_type", string(sdk.SecretTypeOAuth2))),
					assert.Check(resource.TestCheckResourceAttr(clientCredsSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.0.oauth_scopes.#", "2")),
					assert.Check(resource.TestCheckTypeSetElemAttr(clientCredsSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.0.oauth_scopes.*", "scope1")),
					assert.Check(resource.TestCheckTypeSetElemAttr(clientCredsSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.0.oauth_scopes.*", "scope2")),
				),
			},
			// OAuth2 authorization code grant with describe
			{
				Config: config.FromModels(t, authCodeModel, authCodeSecretWithDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.SecretsDatasourceShowOutput(t, authCodeSecretWithDescribe.DatasourceReference()).
						HasName(authCodeId.Name()).
						HasDatabaseName(authCodeId.DatabaseName()).
						HasSchemaName(authCodeId.SchemaName()).
						HasComment(comment).
						HasSecretType(string(sdk.SecretTypeOAuth2)),
					assert.Check(resource.TestCheckResourceAttr(authCodeSecretWithDescribe.DatasourceReference(), "secrets.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(authCodeSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(authCodeSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.0.secret_type", string(sdk.SecretTypeOAuth2))),
					assert.Check(resource.TestCheckResourceAttr(authCodeSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.0.oauth_scopes.#", "0")),
					assert.Check(resource.TestCheckResourceAttrSet(authCodeSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.0.oauth_refresh_token_expiry_time")),
					assert.Check(resource.TestCheckResourceAttr(authCodeSecretWithDescribe.DatasourceReference(), "secrets.0.describe_output.0.comment", comment)),
				),
			},
		},
	})
}

func TestAcc_Secrets_EmptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      secretDatasourceEmptyIn(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func secretDatasourceEmptyIn() string {
	return `
    data "snowflake_secrets" "test" {
        in {
        }
    }
`
}
