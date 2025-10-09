//go:build non_account_level_tests

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

func TestAcc_SecretWithGenericString_BasicFlow(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	comment := random.Comment()

	secretModel := model.SecretWithGenericString("s", id.DatabaseName(), id.SchemaName(), name, "foo")
	secretModelWithComment := model.SecretWithGenericString("s", id.DatabaseName(), id.SchemaName(), name, "bar").
		WithComment(comment)
	secretModelEmptySecretString := model.SecretWithGenericString("s", id.DatabaseName(), id.SchemaName(), name, "")

	resourceReference := secretModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecretWithGenericString),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModels(t, secretModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithGenericStringResource(t, secretModel.ResourceReference()).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasSecretStringString("foo").
							HasCommentString(""),

						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasName(name).
							HasDatabaseName(id.DatabaseName()).
							HasSecretType(string(sdk.SecretTypeGenericString)).
							HasSchemaName(id.SchemaName()).
							HasComment(""),
					),

					resource.TestCheckResourceAttr(resourceReference, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttrSet(resourceReference, "describe_output.0.created_on"),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.name", name),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.database_name", id.DatabaseName()),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.schema_name", id.SchemaName()),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.secret_type", string(sdk.SecretTypeGenericString)),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.username", ""),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.comment", ""),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.oauth_access_token_expiry_time", ""),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.oauth_refresh_token_expiry_time", ""),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.integration_name", ""),
					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.oauth_scopes.#", "0"),
				),
			},
			// set secret_string and comment
			{
				Config: config.FromModels(t, secretModelWithComment),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithGenericStringResource(t, resourceReference).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasSecretStringString("bar").
							HasCommentString(comment),

						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasSecretType(string(sdk.SecretTypeGenericString)).
							HasComment(comment),
					),

					resource.TestCheckResourceAttr(resourceReference, "describe_output.0.comment", comment),
				),
			},
			// set comment externally, external changes for secret_string are not being detected
			{
				PreConfig: func() {
					testClient().Secret.Alter(t, sdk.NewAlterSecretRequest(id).
						WithSet(*sdk.NewSecretSetRequest().
							WithComment("test_comment"),
						),
					)
				},
				Config: config.FromModels(t, secretModelWithComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(resourceReference, "comment", sdk.String(comment), sdk.String("test_comment")),
						planchecks.ExpectChange(resourceReference, "comment", tfjson.ActionUpdate, sdk.String("test_comment"), sdk.String(comment)),
					},
				},
				Check: assertThat(t,
					resourceassert.SecretWithGenericStringResource(t, resourceReference).
						HasNameString(name).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasSecretStringString("bar").
						HasCommentString(comment),
				),
			},
			// import
			{
				ResourceName:            resourceReference,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_string"},
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "database", id.DatabaseId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "schema", id.SchemaId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// unset comment
			{
				Config: config.FromModels(t, secretModelEmptySecretString),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(secretModel.ResourceReference(), "comment", tfjson.ActionUpdate, sdk.String(comment), nil),
					},
				},
				Check: assertThat(t,
					resourceassert.SecretWithClientCredentialsResource(t, secretModelEmptySecretString.ResourceReference()).
						HasCommentString(""),
				),
			},
			// import with no fields set
			{
				ResourceName:            secretModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_string"},
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "database", id.DatabaseId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "schema", id.SchemaId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", ""),
				),
			},
			// destroy
			{
				Config:  config.FromModels(t, secretModel),
				Destroy: true,
			},
			// create with empty secret_string
			{
				Config: config.FromModels(t, secretModelEmptySecretString),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithGenericStringResource(t, secretModelEmptySecretString.ResourceReference()).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasSecretStringString("").
							HasCommentString(""),
					),
				),
			},
		},
	})
}

func TestAcc_SecretWithGenericString_ExternalSecretTypeChange(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()

	secretModel := model.SecretWithGenericString("s", id.DatabaseName(), id.SchemaName(), name, "test_usr")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecretWithGenericString),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModels(t, secretModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecretWithGenericStringResource(t, secretModel.ResourceReference()).
							HasSecretTypeString(string(sdk.SecretTypeGenericString)),
						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasSecretType(string(sdk.SecretTypeGenericString)),
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
						resourceassert.SecretWithGenericStringResource(t, secretModel.ResourceReference()).
							HasSecretTypeString(string(sdk.SecretTypeGenericString)),
						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasSecretType(string(sdk.SecretTypeGenericString)),
					),
				),
			},
		},
	})
}
