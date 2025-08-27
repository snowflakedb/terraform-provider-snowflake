//go:build !account_level_tests

package testacc

import (
	"testing"
	"time"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SecretWithClientCredentials_BasicFlow(t *testing.T) {
	integrationId := testClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := testClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "test_client_id", "test_client_secret").
			WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "foo"}, {Scope: "bar"}, {Scope: "test"}}),
	)
	t.Cleanup(apiIntegrationCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	comment := random.Comment()
	newComment := random.Comment()

	secretModel := model.SecretWithClientCredentials("s", id.DatabaseName(), id.SchemaName(), name, integrationId.Name(), []string{"foo", "bar"}).WithComment(comment)
	secretModelTestInScopes := model.SecretWithClientCredentials("s", id.DatabaseName(), id.SchemaName(), name, integrationId.Name(), []string{"test"}).WithComment(newComment)
	secretModelFooInScopesWithComment := model.SecretWithClientCredentials("s", id.DatabaseName(), id.SchemaName(), name, integrationId.Name(), []string{"foo"}).WithComment(newComment)
	secretModelFooInScopes := model.SecretWithClientCredentials("s", id.DatabaseName(), id.SchemaName(), name, integrationId.Name(), []string{"foo"})
	secretModelWithoutComment := model.SecretWithClientCredentials("s", id.DatabaseName(), id.SchemaName(), name, integrationId.Name(), []string{"foo", "bar"})
	secretModelWithoutCommentWithOauthScopes := model.SecretWithClientCredentials("s", id.DatabaseName(), id.SchemaName(), name, integrationId.Name(), []string{"foo", "bar"}).
		WithOauthScopes([]string{"foo", "bar"})
	secretName := secretModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecretWithClientCredentials),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, secretModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithClientCredentialsResource(t, secretModel.ResourceReference()).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasApiAuthenticationString(integrationId.Name()).
							HasOauthScopesLength(len([]string{"foo", "bar"})).
							HasCommentString(comment),

						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasName(name).
							HasDatabaseName(id.DatabaseName()).
							HasSecretType(string(sdk.SecretTypeOAuth2)).
							HasSchemaName(id.SchemaName()),
					),
					resource.TestCheckResourceAttr(secretName, "oauth_scopes.#", "2"),
					resource.TestCheckTypeSetElemAttr(secretName, "oauth_scopes.*", "foo"),
					resource.TestCheckTypeSetElemAttr(secretName, "oauth_scopes.*", "bar"),

					resource.TestCheckResourceAttr(secretName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttrSet(secretName, "describe_output.0.created_on"),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.name", name),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.database_name", id.DatabaseName()),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.schema_name", id.SchemaName()),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.username", ""),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.oauth_access_token_expiry_time", ""),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.oauth_refresh_token_expiry_time", ""),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.integration_name", integrationId.Name()),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.oauth_scopes.#", "2"),
					resource.TestCheckTypeSetElemAttr(secretName, "describe_output.0.oauth_scopes.*", "foo"),
					resource.TestCheckTypeSetElemAttr(secretName, "describe_output.0.oauth_scopes.*", "bar"),
				),
			},
			// set oauth_scopes and comment in config
			{
				Config: config.FromModels(t, secretModelTestInScopes),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretName, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(secretName, "oauth_scopes", tfjson.ActionUpdate, sdk.String("[bar foo]"), sdk.String("[test]")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithClientCredentialsResource(t, "snowflake_secret_with_client_credentials.s").
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasApiAuthenticationString(integrationId.Name()).
							HasOauthScopesLength(len([]string{"test"})).
							HasCommentString(newComment),
						assert.Check(resource.TestCheckResourceAttr(secretName, "oauth_scopes.#", "1")),
						assert.Check(resource.TestCheckTypeSetElemAttr(secretName, "oauth_scopes.*", "test")),

						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasSecretType(string(sdk.SecretTypeOAuth2)).
							HasComment(newComment),
					),

					resource.TestCheckResourceAttr(secretName, "describe_output.0.comment", newComment),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.oauth_scopes.#", "1"),
					resource.TestCheckTypeSetElemAttr(secretName, "describe_output.0.oauth_scopes.*", "test"),
				),
			},
			// set oauth_scopes and comment externally
			{
				PreConfig: func() {
					req := sdk.NewAlterSecretRequest(id).WithSet(*sdk.NewSecretSetRequest().
						WithSetForFlow(*sdk.NewSetForFlowRequest().
							WithSetForOAuthClientCredentials(
								*sdk.NewSetForOAuthClientCredentialsRequest().WithOauthScopes(
									*sdk.NewOauthScopesListRequest([]sdk.ApiIntegrationScope{{Scope: "bar"}}),
								),
							),
						),
					)
					testClient().Secret.Alter(t, req)
				},
				Config: config.FromModels(t, secretModelFooInScopesWithComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(secretModel.ResourceReference(), "oauth_scopes", sdk.String("[test]"), sdk.String("[bar]")),
						planchecks.ExpectChange(secretModel.ResourceReference(), "oauth_scopes", tfjson.ActionUpdate, sdk.String("[bar]"), sdk.String("[foo]")),
					},
				},
				Check: assertThat(t,
					resourceassert.SecretWithClientCredentialsResource(t, "snowflake_secret_with_client_credentials.s").
						HasNameString(name).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasApiAuthenticationString(integrationId.Name()).
						HasOauthScopesLength(len([]string{"foo"})).
						HasCommentString(newComment),
					assert.Check(resource.TestCheckResourceAttr(secretName, "oauth_scopes.#", "1")),
					assert.Check(resource.TestCheckTypeSetElemAttr(secretName, "oauth_scopes.*", "foo")),
				),
			},
			// unset comment
			{
				Config: config.FromModels(t, secretModelFooInScopes),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretModelFooInScopes.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(secretModelFooInScopes.ResourceReference(), "comment", tfjson.ActionUpdate, sdk.String(newComment), nil),
					},
				},
				Check: assertThat(t,
					resourceassert.SecretWithClientCredentialsResource(t, secretModelFooInScopes.ResourceReference()).
						HasCommentString(""),
				),
			},
			// set comment externally
			{
				PreConfig: func() {
					req := sdk.NewAlterSecretRequest(id).WithSet(*sdk.NewSecretSetRequest().WithComment(comment))
					testClient().Secret.Alter(t, req)
				},
				Config: config.FromModels(t, secretModelWithoutComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(secretModelWithoutComment.ResourceReference(), "comment", tfjson.ActionUpdate, sdk.String(comment), nil),
					},
				},
				Check: assertThat(t,
					resourceassert.SecretWithClientCredentialsResource(t, secretModelWithoutComment.ResourceReference()).
						HasCommentString(""),
				),
			},
			// create without comment
			{
				Config: config.FromModels(t, secretModelWithoutCommentWithOauthScopes),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithClientCredentialsResource(t, "snowflake_secret_with_client_credentials.s").
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasApiAuthenticationString(integrationId.Name()).
							HasOauthScopesLength(len([]string{"foo", "bar"})).
							HasCommentString(""),
					),
					resource.TestCheckResourceAttr(secretName, "oauth_scopes.#", "2"),
					resource.TestCheckTypeSetElemAttr(secretName, "oauth_scopes.*", "foo"),
					resource.TestCheckTypeSetElemAttr(secretName, "oauth_scopes.*", "bar"),
				),
			},
			// import
			{
				ResourceName:      secretModelWithoutComment.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "database", id.DatabaseId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "schema", id.SchemaId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "api_authentication", integrationId.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "oauth_scopes.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", ""),
				),
			},
		},
	})
}

func TestAcc_SecretWithClientCredentials_EmptyScopesList(t *testing.T) {
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

func TestAcc_SecretWithClientCredentials_ExternalSecretTypeChange(t *testing.T) {
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

func TestAcc_SecretWithClientCredentials_ExternalSecretTypeChangeToOAuthAuthCodeGrant(t *testing.T) {
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
