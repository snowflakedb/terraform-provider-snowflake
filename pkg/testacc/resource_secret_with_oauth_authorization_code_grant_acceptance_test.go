//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"
	"time"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SecretWithAuthorizationCodeGrant_BasicFlow(t *testing.T) {
	apiIntegration, apiIntegrationCleanup := testClient().SecurityIntegration.CreateApiAuthenticationWithClientCredentialsFlowWithEnabled(t, true)
	t.Cleanup(apiIntegrationCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	comment := random.Comment()
	newComment := random.Comment()
	refreshTokenExpiryDateTime := time.Now().Add(24 * time.Hour).Format(time.DateTime)
	newRefreshTokenExpiryDateOnly := time.Now().Add(4 * 24 * time.Hour).Format(time.DateOnly)
	refreshToken := "test_token"
	newRefreshToken := "new_test_token"

	secretModel := model.SecretWithAuthorizationCodeGrant("s", id.DatabaseName(), id.SchemaName(), name, apiIntegration.ID().Name(), refreshToken, refreshTokenExpiryDateTime)
	secretModelAllSet := model.SecretWithAuthorizationCodeGrant("s", id.DatabaseName(), id.SchemaName(), name, apiIntegration.ID().Name(), newRefreshToken, newRefreshTokenExpiryDateOnly).WithComment(comment)

	resourceReference := secretModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecretWithAuthorizationCodeGrant),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModels(t, secretModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithAuthorizationCodeGrantResource(t, secretModel.ResourceReference()).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasApiAuthenticationString(apiIntegration.ID().Name()).
							HasOauthRefreshTokenString(refreshToken).
							HasOauthRefreshTokenExpiryTimeString(refreshTokenExpiryDateTime).
							HasCommentString(""),

						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasName(name).
							HasDatabaseName(id.DatabaseName()).
							HasSecretType(string(sdk.SecretTypeOAuth2)).
							HasSchemaName(id.SchemaName()).
							HasComment(""),
					),

					resource.TestCheckResourceAttr(resourceReference, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttrSet(resourceReference, "describe_output.0.created_on"),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.name", name),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.database_name", id.DatabaseName()),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.schema_name", id.SchemaName()),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.secret_type", string(sdk.SecretTypeOAuth2)),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.integration_name", apiIntegration.ID().Name()),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.username", ""),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.oauth_access_token_expiry_time", ""),
					resource.TestCheckResourceAttrSet(resourceReference, "describe_output.0.oauth_refresh_token_expiry_time"),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.comment", ""),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.oauth_scopes.#", "0"),
				),
			},
			// set all
			{
				Config: config.FromModels(t, secretModelAllSet),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "comment", tfjson.ActionUpdate, nil, sdk.String(comment)),
						planchecks.ExpectChange(resourceReference, "oauth_refresh_token", tfjson.ActionUpdate, sdk.String(refreshToken), sdk.String(newRefreshToken)),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithAuthorizationCodeGrantResource(t, resourceReference).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasApiAuthenticationString(apiIntegration.ID().Name()).
							HasOauthRefreshTokenString(newRefreshToken).
							HasOauthRefreshTokenExpiryTimeString(newRefreshTokenExpiryDateOnly).
							HasCommentString(comment),

						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasSecretType(string(sdk.SecretTypeOAuth2)).
							HasComment(comment),
					),
					resource.TestCheckResourceAttrSet(resourceReference, "describe_output.0.oauth_refresh_token_expiry_time"),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.comment", comment),
				),
			},
			// set comment and refresh_token_expiry_time externally
			{
				PreConfig: func() {
					testClient().Secret.Alter(t, sdk.NewAlterSecretRequest(id).WithSet(*sdk.NewSecretSetRequest().
						WithComment(newComment).
						WithSetForFlow(*sdk.NewSetForFlowRequest().
							WithSetForOAuthAuthorization(*sdk.NewSetForOAuthAuthorizationRequest().
								WithOauthRefreshTokenExpiryTime(time.Now().Add(24 * time.Hour).Format(time.DateOnly)),
							),
						),
					))
				},
				Config: config.FromModels(t, secretModelAllSet),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "comment", tfjson.ActionUpdate, sdk.String(newComment), sdk.String(comment)),
						planchecks.ExpectComputed(resourceReference, r.DescribeOutputAttributeName, true),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithAuthorizationCodeGrantResource(t, resourceReference).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasApiAuthenticationString(apiIntegration.ID().Name()).
							HasOauthRefreshTokenString(newRefreshToken).
							HasOauthRefreshTokenExpiryTimeString(newRefreshTokenExpiryDateOnly).
							HasCommentString(comment),
						assert.Check(resource.TestCheckResourceAttrSet(resourceReference, "describe_output.0.oauth_refresh_token_expiry_time")),
					),
				),
			},
			// import
			{
				ResourceName:            resourceReference,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth_refresh_token"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedSecretWithAuthorizationCodeGrantResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasApiAuthenticationString(apiIntegration.ID().Name()).
						HasCommentString(comment).
						HasOauthRefreshTokenExpiryTimeNotEmpty(),
				),
			},
		},
	})
}

func TestAcc_SecretWithAuthorizationCodeGrant_DifferentTimeFormats(t *testing.T) {
	apiIntegration, apiIntegrationCleanup := testClient().SecurityIntegration.CreateApiAuthenticationWithClientCredentialsFlowWithEnabled(t, true)
	t.Cleanup(apiIntegrationCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()

	refreshTokenExpiryDateOnly := time.Now().Add(4 * 24 * time.Hour).Format(time.DateOnly)
	refreshTokenExpiryWithoutSeconds := time.Now().Add(4 * 24 * time.Hour).Format("2006-01-02 15:04")
	refreshTokenExpiryDateTime := time.Now().Add(4 * 24 * time.Hour).Format(time.DateTime)
	refreshTokenExpiryWithPDT := fmt.Sprintf("%s %s", time.Now().Add(4*24*time.Hour).Format("2006-01-02 15:04"), "PDT")

	secretModelDateOnly := model.SecretWithAuthorizationCodeGrant("s", id.DatabaseName(), id.SchemaName(), name, apiIntegration.ID().Name(), "test_token", refreshTokenExpiryDateOnly)
	secretModelWithoutSeconds := model.SecretWithAuthorizationCodeGrant("s", id.DatabaseName(), id.SchemaName(), name, apiIntegration.ID().Name(), "test_token", refreshTokenExpiryWithoutSeconds)
	secretModelDateTime := model.SecretWithAuthorizationCodeGrant("s", id.DatabaseName(), id.SchemaName(), name, apiIntegration.ID().Name(), "test_token", refreshTokenExpiryDateTime)
	secretModelWithPDT := model.SecretWithAuthorizationCodeGrant("s", id.DatabaseName(), id.SchemaName(), name, apiIntegration.ID().Name(), "test_token", refreshTokenExpiryWithPDT)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecretWithAuthorizationCodeGrant),
		Steps: []resource.TestStep{
			// create with DateOnly
			{
				Config: config.FromModels(t, secretModelDateOnly),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithAuthorizationCodeGrantResource(t, secretModelDateOnly.ResourceReference()).
							HasOauthRefreshTokenExpiryTimeString(refreshTokenExpiryDateOnly),
						assert.Check(resource.TestCheckResourceAttrSet(secretModelDateOnly.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time")),
					),
				),
			},
			// update with DateTime without seconds
			{
				Config: config.FromModels(t, secretModelWithoutSeconds),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithAuthorizationCodeGrantResource(t, secretModelWithoutSeconds.ResourceReference()).
							HasOauthRefreshTokenExpiryTimeString(refreshTokenExpiryWithoutSeconds),
						assert.Check(resource.TestCheckResourceAttrSet(secretModelWithoutSeconds.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time")),
					),
				),
			},
			// update with DateTime
			{
				Config: config.FromModels(t, secretModelDateTime),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithAuthorizationCodeGrantResource(t, secretModelDateTime.ResourceReference()).
							HasOauthRefreshTokenExpiryTimeString(refreshTokenExpiryDateTime),
						assert.Check(resource.TestCheckResourceAttrSet(secretModelDateTime.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time")),
					),
				),
			},
			// update with DateTime with PDT timezone
			{
				Config: config.FromModels(t, secretModelWithPDT),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithAuthorizationCodeGrantResource(t, secretModelWithPDT.ResourceReference()).
							HasOauthRefreshTokenExpiryTimeString(refreshTokenExpiryWithPDT),
						assert.Check(resource.TestCheckResourceAttrSet(secretModelWithPDT.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time")),
					),
				),
			},
		},
	})
}

func TestAcc_SecretWithAuthorizationCodeGrant_ExternalRefreshTokenExpiryTimeChange(t *testing.T) {
	apiIntegration, apiIntegrationCleanup := testClient().SecurityIntegration.CreateApiAuthenticationWithClientCredentialsFlowWithEnabled(t, true)
	t.Cleanup(apiIntegrationCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	comment := random.Comment()
	refreshTokenExpiryDateTime := time.Now().Add(24 * time.Hour).Format(time.DateTime)
	externalRefreshTokenExpiryTime := time.Now().Add(10 * 24 * time.Hour)
	refreshToken := "test_token"

	secretModel := model.SecretWithAuthorizationCodeGrant("s", id.DatabaseName(), id.SchemaName(), name, apiIntegration.ID().Name(), refreshToken, refreshTokenExpiryDateTime).WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecretWithAuthorizationCodeGrant),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModels(t, secretModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithAuthorizationCodeGrantResource(t, secretModel.ResourceReference()).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasApiAuthenticationString(apiIntegration.ID().Name()).
							HasOauthRefreshTokenString(refreshToken).
							HasOauthRefreshTokenExpiryTimeString(refreshTokenExpiryDateTime).
							HasCommentString(comment),

						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasName(name).
							HasDatabaseName(id.DatabaseName()).
							HasSecretType(string(sdk.SecretTypeOAuth2)).
							HasSchemaName(id.SchemaName()).
							HasComment(comment),
					),
				),
			},
			{
				PreConfig: func() {
					testClient().Secret.Alter(t, sdk.NewAlterSecretRequest(id).
						WithSet(*sdk.NewSecretSetRequest().
							WithSetForFlow(*sdk.NewSetForFlowRequest().
								WithSetForOAuthAuthorization(*sdk.NewSetForOAuthAuthorizationRequest().
									WithOauthRefreshTokenExpiryTime(externalRefreshTokenExpiryTime.Format(time.DateOnly)),
								),
							),
						),
					)
				},
				Config: config.FromModels(t, secretModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretModel.ResourceReference(), plancheck.ResourceActionUpdate),
						// cannot check before value due to snowflake timestamp format
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithAuthorizationCodeGrantResource(t, secretModel.ResourceReference()).
							HasOauthRefreshTokenExpiryTimeString(refreshTokenExpiryDateTime),
						assert.Check(resource.TestCheckResourceAttrSet(secretModel.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time")),
					),
				),
			},
		},
	})
}

func TestAcc_SecretWithAuthorizationCodeGrant_ExternalSecretTypeChange(t *testing.T) {
	apiIntegration, apiIntegrationCleanup := testClient().SecurityIntegration.CreateApiAuthenticationWithClientCredentialsFlowWithEnabled(t, true)
	t.Cleanup(apiIntegrationCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()

	secretModel := model.SecretWithAuthorizationCodeGrant("s", id.DatabaseName(), id.SchemaName(), name, apiIntegration.ID().Name(), "test_refresh_token", time.Now().Add(24*time.Hour).Format(time.DateOnly))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecretWithAuthorizationCodeGrant),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModels(t, secretModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithAuthorizationCodeGrantResource(t, secretModel.ResourceReference()).
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
					_, cleanup := testClient().Secret.CreateWithBasicAuthenticationFlow(t, id, "test_pswd", "test_usr")
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
						resourceassert.SecretWithAuthorizationCodeGrantResource(t, secretModel.ResourceReference()).
							HasSecretTypeString(string(sdk.SecretTypeOAuth2)),
						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasSecretType(string(sdk.SecretTypeOAuth2)),
					),
				),
			},
		},
	})
}

func TestAcc_SecretWithAuthorizationCodeGrant_ExternalSecretTypeChangeToOAuthClientCredentials(t *testing.T) {
	integrationId := testClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := testClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "test_client_id", "test_client_secret").
			WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "foo"}, {Scope: "bar"}, {Scope: "test"}}),
	)
	t.Cleanup(apiIntegrationCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()

	secretModel := model.SecretWithAuthorizationCodeGrant("s", id.DatabaseName(), id.SchemaName(), name, integrationId.Name(), "test_refresh_token", time.Now().Add(24*time.Hour).Format(time.DateOnly))

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
						resourceassert.SecretWithAuthorizationCodeGrantResource(t, secretModel.ResourceReference()).
							HasSecretTypeString(string(sdk.SecretTypeOAuth2)),
						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasSecretType(string(sdk.SecretTypeOAuth2)),
					),
					resource.TestCheckResourceAttrSet(secretModel.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time"),
					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "describe_output.0.oauth_scopes.#", "0"),
				),
			},
			// create or replace with same secret type, but different create flow
			{
				PreConfig: func() {
					testClient().Secret.DropFunc(t, id)()
					_, cleanup := testClient().Secret.CreateWithOAuthClientCredentialsFlow(t, id, integrationId, []sdk.ApiIntegrationScope{})
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
						resourceassert.SecretWithAuthorizationCodeGrantResource(t, secretModel.ResourceReference()).
							HasSecretTypeString(string(sdk.SecretTypeOAuth2)),
						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasSecretType(string(sdk.SecretTypeOAuth2)),
					),
					resource.TestCheckResourceAttrSet(secretModel.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time"),
					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "describe_output.0.oauth_scopes.#", "0"),
				),
			},
		},
	})
}
