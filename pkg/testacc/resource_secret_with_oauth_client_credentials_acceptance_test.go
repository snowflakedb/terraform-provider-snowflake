//go:build non_account_level_tests

package testacc

import (
	"testing"
	"time"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SecretWithOauthClientCredentials_BasicUseCase(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	integrationId := testClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := testClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "test_client_id", "test_client_secret").
			WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "scope1"}, {Scope: "scope2"}, {Scope: "scope3"}, {Scope: "scope4"}}),
	)
	t.Cleanup(apiIntegrationCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	comment := random.Comment()
	oauthScopes := []string{"scope1", "scope2"}
	currentRole := testClient().Context.CurrentRole(t).Name()

	basic := model.SecretWithClientCredentials("test", id.DatabaseName(), id.SchemaName(), id.Name(), integrationId.Name(), oauthScopes)

	complete := model.SecretWithClientCredentials("test", id.DatabaseName(), id.SchemaName(), id.Name(), integrationId.Name(), []string{"scope3", "scope4"}).
		WithComment(comment)

	assertBasic := []assert.TestCheckFuncProvider{
		objectassert.Secret(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasSecretType(string(sdk.SecretTypeOAuth2)).
			HasOwner(currentRole).
			HasNoComment(),

		resourceassert.SecretWithClientCredentialsResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasSecretTypeString("OAUTH2").
			HasApiAuthenticationString(integrationId.Name()).
			HasOauthScopes("scope1", "scope2").
			HasCommentString(""),

		resourceshowoutputassert.SecretShowOutput(t, basic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasSecretType(string(sdk.SecretTypeOAuth2)).
			HasOwner(currentRole).
			HasComment("").
			HasOwnerRoleType("ROLE"),

		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.name", id.Name())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.secret_type", string(sdk.SecretTypeOAuth2))),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.username", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.comment", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_access_token_expiry_time", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.integration_name", integrationId.Name())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_scopes.#", "2")),
		assert.Check(resource.TestCheckTypeSetElemAttr(basic.ResourceReference(), "describe_output.0.oauth_scopes.*", "scope1")),
		assert.Check(resource.TestCheckTypeSetElemAttr(basic.ResourceReference(), "describe_output.0.oauth_scopes.*", "scope2")),
	}
	assertComplete := []assert.TestCheckFuncProvider{
		objectassert.Secret(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasSecretType(string(sdk.SecretTypeOAuth2)).
			HasOwner(currentRole).
			HasComment(comment),

		resourceassert.SecretWithClientCredentialsResource(t, complete.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasSecretTypeString("OAUTH2").
			HasApiAuthenticationString(integrationId.Name()).
			HasOauthScopes("scope3", "scope4").
			HasCommentString(comment),

		resourceshowoutputassert.SecretShowOutput(t, complete.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasSecretType(string(sdk.SecretTypeOAuth2)).
			HasOwner(currentRole).
			HasComment(comment).
			HasOwnerRoleType("ROLE"),

		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.name", id.Name())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.secret_type", string(sdk.SecretTypeOAuth2))),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.username", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.comment", comment)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_access_token_expiry_time", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.integration_name", integrationId.Name())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_scopes.#", "2")),
		assert.Check(resource.TestCheckTypeSetElemAttr(complete.ResourceReference(), "describe_output.0.oauth_scopes.*", "scope3")),
		assert.Check(resource.TestCheckTypeSetElemAttr(complete.ResourceReference(), "describe_output.0.oauth_scopes.*", "scope4")),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecretWithClientCredentials),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:                  config.FromModels(t, basic),
				ResourceName:            basic.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_authentication"},
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
			// Import - with optionals
			{
				Config:                  config.FromModels(t, complete),
				ResourceName:            complete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_authentication"},
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// TODO(SNOW-2457144; next prs): see why is this step failing
			// Update - detect external changes (temporarily disabled due to SQL compilation error)
			// {
			//	PreConfig: func() {
			//		testClient().Secret.Alter(t, sdk.NewAlterSecretRequest(id).WithSet(
			//			*sdk.NewSecretSetRequest().WithComment(comment),
			//		))
			//	},
			//	ConfigPlanChecks: resource.ConfigPlanChecks{
			//		PreApply: []plancheck.PlanCheck{
			//			plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
			//		},
			//	},
			//	Config: config.FromModels(t, basic),
			//	Check:  assertThat(t, assertBasic...),
			// },
			// Destroy - ensure secret is destroyed before the next step
			{
				Destroy: true,
				Config:  config.FromModels(t, basic),
				Check: assertThat(t,
					objectassert.SecretDoesNotExist(t, id),
				),
			},
			// Create - with optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}

func TestAcc_SecretWithClientCredentials_CompleteUseCase_EmptyScopesList(t *testing.T) {
	integrationId := testClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := testClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "test_client_id", "test_client_secret").
			WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "foo"}, {Scope: "bar"}, {Scope: "test"}}),
	)
	t.Cleanup(apiIntegrationCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()

	secretModel := model.SecretWithClientCredentials("s", id.DatabaseName(), id.SchemaName(), name, integrationId.Name(), []string{})
	secretModelEmptyScopes := model.SecretWithClientCredentials("s", id.DatabaseName(), id.SchemaName(), name, integrationId.Name(), []string{})
	secretModelWithScope := model.SecretWithClientCredentials("s", id.DatabaseName(), id.SchemaName(), name, integrationId.Name(), []string{}).WithOauthScopes([]string{"foo"})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecretWithClientCredentials),
		Steps: []resource.TestStep{
			// create secret without providing oauth_scopes value
			{
				Config: config.FromModels(t, secretModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithClientCredentialsResource(t, secretModel.ResourceReference()).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasApiAuthenticationString(integrationId.Name()).
							HasCommentString(""),
					),
					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "oauth_scopes.#", "0"),
				),
			},
			// Set oauth_scopes
			{
				Config: config.FromModels(t, secretModelWithScope),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretModelWithScope.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(secretModelWithScope.ResourceReference(), "oauth_scopes", tfjson.ActionUpdate, sdk.String("[]"), sdk.String("[foo]")),
					},
				},
				Check: assertThat(t,
					resourceassert.SecretWithClientCredentialsResource(t, secretModel.ResourceReference()).
						HasNameString(name).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasApiAuthenticationString(integrationId.Name()),
					assert.Check(resource.TestCheckResourceAttr(secretModel.ResourceReference(), "oauth_scopes.#", "1")),
					assert.Check(resource.TestCheckTypeSetElemAttr(secretModel.ResourceReference(), "oauth_scopes.*", "foo")),
				),
			},
			// Set empty oauth_scopes
			{
				Config: config.FromModels(t, secretModelEmptyScopes),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(secretModel.ResourceReference(), "oauth_scopes", tfjson.ActionUpdate, sdk.String("[foo]"), sdk.String("[]")),
					},
				},
				Check: assertThat(t,
					resourceassert.SecretWithClientCredentialsResource(t, secretModel.ResourceReference()).
						HasNameString(name).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasApiAuthenticationString(integrationId.Name()),
					assert.Check(resource.TestCheckResourceAttr(secretModel.ResourceReference(), "oauth_scopes.#", "0")),
				),
			},
		},
	})
}

func TestAcc_SecretWithClientCredentials_CompleteUseCase_ExternalSecretTypeChange(t *testing.T) {
	integrationId := testClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := testClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "test_client_id", "test_client_secret").
			WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "foo"}, {Scope: "bar"}, {Scope: "test"}}),
	)
	t.Cleanup(apiIntegrationCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()

	secretModel := model.SecretWithClientCredentials("s", id.DatabaseName(), id.SchemaName(), name, integrationId.Name(), []string{"foo", "bar"})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecretWithClientCredentials),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModels(t, secretModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithClientCredentialsResource(t, secretModel.ResourceReference()).
							HasSecretTypeString(string(sdk.SecretTypeOAuth2)),
						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasSecretType(string(sdk.SecretTypeOAuth2)),
					),
				),
			},
			// create or replace with different secret type
			{
				PreConfig: func() {
					testClient().Secret.DropFunc(t, id)()
					_, cleanup := testClient().Secret.CreateWithGenericString(t, id, "test_secret_string")
					t.Cleanup(cleanup)
				},
				Config: config.FromModels(t, secretModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithClientCredentialsResource(t, secretModel.ResourceReference()).
							HasSecretTypeString(string(sdk.SecretTypeOAuth2)),
						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasSecretType(string(sdk.SecretTypeOAuth2)),
					),
				),
			},
		},
	})
}

func TestAcc_SecretWithClientCredentials_CompleteUseCase_ExternalSecretTypeChangeToOAuthAuthCodeGrant(t *testing.T) {
	integrationId := testClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := testClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "test_client_id", "test_client_secret").
			WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "foo"}, {Scope: "bar"}, {Scope: "test"}}),
	)
	t.Cleanup(apiIntegrationCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()

	secretModel := model.SecretWithClientCredentials("s", id.DatabaseName(), id.SchemaName(), name, integrationId.Name(), []string{"foo", "bar"})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecretWithClientCredentials),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModels(t, secretModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithClientCredentialsResource(t, secretModel.ResourceReference()).
							HasSecretTypeString(string(sdk.SecretTypeOAuth2)),
						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasSecretType(string(sdk.SecretTypeOAuth2)),
					),
					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "describe_output.0.oauth_scopes.#", "2"),
					resource.TestCheckTypeSetElemAttr(secretModel.ResourceReference(), "describe_output.0.oauth_scopes.*", "foo"),
					resource.TestCheckTypeSetElemAttr(secretModel.ResourceReference(), "describe_output.0.oauth_scopes.*", "bar"),
					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time", ""),
				),
			},
			// create or replace with the same secret type but different flow
			{
				PreConfig: func() {
					testClient().Secret.DropFunc(t, id)()
					_, cleanup := testClient().Secret.CreateWithOAuthAuthorizationCodeFlow(t, id, integrationId, "test_refresh_token", time.Now().Add(24*time.Hour).Format(time.DateOnly))
					t.Cleanup(cleanup)
				},
				Config: config.FromModels(t, secretModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithClientCredentialsResource(t, secretModel.ResourceReference()).
							HasSecretTypeString(string(sdk.SecretTypeOAuth2)),
						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasSecretType(string(sdk.SecretTypeOAuth2)),
					),
					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "describe_output.0.oauth_scopes.#", "2"),
					resource.TestCheckTypeSetElemAttr(secretModel.ResourceReference(), "describe_output.0.oauth_scopes.*", "foo"),
					resource.TestCheckTypeSetElemAttr(secretModel.ResourceReference(), "describe_output.0.oauth_scopes.*", "bar"),
					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time", ""),
				),
			},
		},
	})
}
