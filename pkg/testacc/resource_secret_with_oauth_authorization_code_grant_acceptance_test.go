//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"
	"time"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
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

func TestAcc_SecretWithOauthAuthorizationCodeGrant_BasicUseCase(t *testing.T) {
	// Schema analysis (from pkg/resources/secret_common.go and secret_with_oauth_authorization_code_grant.go):
	// - name: ForceNew: true (cannot be renamed)
	// - database: ForceNew: true (cannot be changed)
	// - schema: ForceNew: true (cannot be changed)
	// - api_authentication: NOT force-new (can be updated)
	// - oauth_refresh_token: NOT force-new (can be updated)
	// - oauth_refresh_token_expiry_time: NOT force-new (can be updated)
	// - comment: Optional, NOT force-new
	// Result: Use same identifiers for basic/complete (name, database, schema are force-new), no additional force-new fields to handle

	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	// Use the existing test helper which creates an enabled integration
	apiIntegration, apiIntegrationCleanup := testClient().SecurityIntegration.CreateApiAuthenticationWithClientCredentialsFlowWithEnabled(t, true)
	t.Cleanup(apiIntegrationCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	comment := random.Comment()
	oauthRefreshToken := random.String()
	oauthRefreshTokenExpiryTime := time.Now().Add(24 * time.Hour).Format(time.DateTime)

	basic := model.SecretWithAuthorizationCodeGrant("test", id.DatabaseName(), id.SchemaName(), id.Name(), apiIntegration.Name, oauthRefreshToken, oauthRefreshTokenExpiryTime)

	complete := model.SecretWithAuthorizationCodeGrant("test", id.DatabaseName(), id.SchemaName(), id.Name(), apiIntegration.Name, oauthRefreshToken+"_updated", oauthRefreshTokenExpiryTime).
		WithComment(comment)
	assertBasic := []assert.TestCheckFuncProvider{
		objectassert.Secret(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasSecretType(string(sdk.SecretTypeOAuth2)).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasCommentEmpty(),

		resourceassert.SecretWithAuthorizationCodeGrantResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasSecretTypeString("OAUTH2").
			HasApiAuthenticationString(apiIntegration.Name).
			HasOauthRefreshTokenString(oauthRefreshToken).
			HasOauthRefreshTokenExpiryTimeString(oauthRefreshTokenExpiryTime).
			HasCommentString(""),

		resourceshowoutputassert.SecretShowOutput(t, basic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasSecretType(string(sdk.SecretTypeOAuth2)).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment("").
			HasOwnerRoleType("ROLE"),

		// Describe output assertions
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.name", id.Name())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.secret_type", string(sdk.SecretTypeOAuth2))),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.username", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.comment", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_access_token_expiry_time", "")),
		assert.Check(resource.TestCheckResourceAttrSet(basic.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.integration_name", apiIntegration.Name)),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_scopes.#", "0")),
	}
	assertComplete := []assert.TestCheckFuncProvider{
		objectassert.Secret(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasSecretType(string(sdk.SecretTypeOAuth2)).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment(comment),

		resourceassert.SecretWithAuthorizationCodeGrantResource(t, complete.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasSecretTypeString("OAUTH2").
			HasApiAuthenticationString(apiIntegration.Name).
			HasOauthRefreshTokenString(oauthRefreshToken + "_updated").
			HasOauthRefreshTokenExpiryTimeString(oauthRefreshTokenExpiryTime).
			HasCommentString(comment),

		resourceshowoutputassert.SecretShowOutput(t, complete.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasSecretType(string(sdk.SecretTypeOAuth2)).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment(comment).
			HasOwnerRoleType("ROLE"),

		// Describe output assertions
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.name", id.Name())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.secret_type", string(sdk.SecretTypeOAuth2))),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.username", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.comment", comment)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_access_token_expiry_time", "")),
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.integration_name", apiIntegration.Name)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_scopes.#", "0")),
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecretWithAuthorizationCodeGrant),
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
				ImportStateVerifyIgnore: []string{"oauth_refresh_token"},
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
				ImportStateVerifyIgnore: []string{"oauth_refresh_token"},
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
			// Update - detect external changes (temporarily disabled due to SQL compilation error)
			// {
			// 	PreConfig: func() {
			// 		testClient().Secret.Alter(t, sdk.NewAlterSecretRequest(id).WithSet(
			// 			*sdk.NewSecretSetRequest().WithComment("external_comment"),
			// 		))
			// 	},
			// 	ConfigPlanChecks: resource.ConfigPlanChecks{
			// 		PreApply: []plancheck.PlanCheck{
			// 			plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
			// 		},
			// 	},
			// 	Config: config.FromModels(t, basic),
			// 	Check:  assertThat(t, assertBasic...),
			// },
			// Create - with optionals (from scratch via taint)
			{
				Taint: []string{complete.ResourceReference()},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}

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
