//go:build !account_level_tests

package testacc

import (
	"testing"

	tfjson "github.com/hashicorp/terraform-json"

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

func TestAcc_SecretWithBasicAuthentication_BasicFlow(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	comment := random.Comment()

	secretModel := model.SecretWithBasicAuthentication("s", id.DatabaseName(), id.SchemaName(), name, "foo", "foo")
	secretModelDifferentCredentialsWithComment := model.SecretWithBasicAuthentication("s", id.DatabaseName(), id.SchemaName(), name, "bar", "bar").WithComment(comment)
	secretModelWithoutComment := model.SecretWithBasicAuthentication("s", id.DatabaseName(), id.SchemaName(), name, "bar", "bar")
	secretModelEmptyCredentials := model.SecretWithBasicAuthentication("s", id.DatabaseName(), id.SchemaName(), name, "", "")

	resourceReference := secretModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecretWithBasicAuthentication),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModels(t, secretModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithBasicAuthenticationResource(t, resourceReference).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasUsernameString("foo").
							HasPasswordString("foo").
							HasCommentString(""),

						resourceshowoutputassert.SecretShowOutput(t, resourceReference).
							HasName(name).
							HasDatabaseName(id.DatabaseName()).
							HasSecretType(string(sdk.SecretTypePassword)).
							HasSchemaName(id.SchemaName()).
							HasComment(""),
					),

					resource.TestCheckResourceAttr(resourceReference, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttrSet(resourceReference, "describe_output.0.created_on"),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.name", name),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.database_name", id.DatabaseName()),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.schema_name", id.SchemaName()),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.secret_type", string(sdk.SecretTypePassword)),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.username", "foo"),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.comment", ""),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.oauth_access_token_expiry_time", ""),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.oauth_refresh_token_expiry_time", ""),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.integration_name", ""),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.oauth_scopes.#", "0"),
				),
			},
			// set username, password and comment
			{
				Config: config.FromModels(t, secretModelDifferentCredentialsWithComment),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithBasicAuthenticationResource(t, resourceReference).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasUsernameString("bar").
							HasPasswordString("bar").
							HasCommentString(comment),

						resourceshowoutputassert.SecretShowOutput(t, resourceReference).
							HasSecretType(string(sdk.SecretTypePassword)).
							HasComment(comment),
					),

					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.username", "bar"),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.comment", comment),
				),
			},
			// set username and comment externally
			{
				PreConfig: func() {
					testClient().Secret.Alter(t, sdk.NewAlterSecretRequest(id).
						WithSet(*sdk.NewSecretSetRequest().
							WithComment("test_comment").
							WithSetForFlow(*sdk.NewSetForFlowRequest().
								WithSetForBasicAuthentication(*sdk.NewSetForBasicAuthenticationRequest().
									WithUsername("test_username"),
								),
							),
						),
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(resourceReference, "comment", sdk.String(comment), sdk.String("test_comment")),
						planchecks.ExpectDrift(resourceReference, "username", sdk.String("bar"), sdk.String("test_username")),

						planchecks.ExpectChange(resourceReference, "comment", tfjson.ActionUpdate, sdk.String("test_comment"), sdk.String(comment)),
						planchecks.ExpectChange(resourceReference, "username", tfjson.ActionUpdate, sdk.String("test_username"), sdk.String("bar")),
					},
				},
				Config: config.FromModels(t, secretModelDifferentCredentialsWithComment),
				Check: assertThat(t,
					resourceassert.SecretWithBasicAuthenticationResource(t, resourceReference).
						HasNameString(name).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUsernameString("bar").
						HasPasswordString("bar").
						HasCommentString(comment),
				),
			},
			// import
			{
				ResourceName:            resourceReference,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "database", id.DatabaseId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "schema", id.SchemaId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "username", "bar"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// unset comment
			{
				Config: config.FromModels(t, secretModelWithoutComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, "comment", tfjson.ActionUpdate, sdk.String(comment), nil),
					},
				},
				Check: assertThat(t,
					resourceassert.SecretWithClientCredentialsResource(t, resourceReference).
						HasCommentString(""),
				),
			},
			// import with no fields set
			{
				ResourceName:            resourceReference,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "database", id.DatabaseId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "schema", id.SchemaId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "username", "bar"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", ""),
				),
			},
			// set empty username and password
			{
				Config: config.FromModels(t, secretModelEmptyCredentials),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithBasicAuthenticationResource(t, resourceReference).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasUsernameString("").
							HasPasswordString("").
							HasCommentString(""),
					),
				),
			},
		},
	})
}

func TestAcc_SecretWithBasicAuthentication_CreateWithEmptyCredentials(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	secretModelEmptyCredentials := model.SecretWithBasicAuthentication("s", id.DatabaseName(), id.SchemaName(), name, "", "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecretWithBasicAuthentication),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, secretModelEmptyCredentials),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithBasicAuthenticationResource(t, secretModelEmptyCredentials.ResourceReference()).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasUsernameString("").
							HasPasswordString("").
							HasCommentString(""),
					),
				),
			},
		},
	})
}

func TestAcc_SecretWithBasicAuthentication_ExternalSecretTypeChange(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	secretModel := model.SecretWithBasicAuthentication("s", id.DatabaseName(), id.SchemaName(), name, "test_pswd", "test_usr")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecretWithBasicAuthentication),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModels(t, secretModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithBasicAuthenticationResource(t, secretModel.ResourceReference()).
							HasSecretTypeString(string(sdk.SecretTypePassword)),
						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasSecretType(string(sdk.SecretTypePassword)),
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
						resourceassert.SecretWithBasicAuthenticationResource(t, secretModel.ResourceReference()).
							HasSecretTypeString(string(sdk.SecretTypePassword)),
						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasSecretType(string(sdk.SecretTypePassword)),
					),
				),
			},
		},
	})
}
