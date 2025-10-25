//go:build non_account_level_tests

package testacc

import (
	"testing"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
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

func TestAcc_SecretWithGenericString_BasicUseCase(t *testing.T) {
	// Schema analysis (from pkg/resources/secret_common.go and secret_with_generic_string.go):
	// - name: ForceNew: true (cannot be renamed)
	// - database: ForceNew: true (cannot be changed)
	// - schema: ForceNew: true (cannot be changed)
	// - secret_string: NOT force-new (can be updated)
	// - comment: Optional, NOT force-new
	// Result: Use same identifiers for basic/complete (name, database, schema are force-new), no additional force-new fields to handle

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	secretString := random.String()

	basic := model.SecretWithGenericString("test", id.DatabaseName(), id.SchemaName(), id.Name(), secretString)

	complete := model.SecretWithGenericString("test", id.DatabaseName(), id.SchemaName(), id.Name(), secretString+"_updated").
		WithComment(comment)

	assertBasic := []assert.TestCheckFuncProvider{
		objectassert.Secret(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasSecretType(string(sdk.SecretTypeGenericString)).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasCommentEmpty(),

		resourceassert.SecretWithGenericStringResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasSecretTypeString("GENERIC_STRING").
			HasSecretStringString(secretString).
			HasCommentString(""),

		resourceshowoutputassert.SecretShowOutput(t, basic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasSecretType(string(sdk.SecretTypeGenericString)).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment("").
			HasOwnerRoleType("ROLE"),

		// Describe output assertions
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.name", id.Name())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.secret_type", string(sdk.SecretTypeGenericString))),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.username", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.comment", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_access_token_expiry_time", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.integration_name", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.oauth_scopes.#", "0")),
	}

	assertComplete := []assert.TestCheckFuncProvider{
		objectassert.Secret(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasSecretType(string(sdk.SecretTypeGenericString)).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment(comment),

		resourceassert.SecretWithGenericStringResource(t, complete.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasSecretTypeString("GENERIC_STRING").
			HasSecretStringString(secretString + "_updated").
			HasCommentString(comment),

		resourceshowoutputassert.SecretShowOutput(t, complete.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasSecretType(string(sdk.SecretTypeGenericString)).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment(comment).
			HasOwnerRoleType("ROLE"),

		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.name", id.Name())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.secret_type", string(sdk.SecretTypeGenericString))),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.username", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.comment", comment)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_access_token_expiry_time", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_refresh_token_expiry_time", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.integration_name", "")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.oauth_scopes.#", "0")),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecretWithGenericString),
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
				ImportStateVerifyIgnore: []string{"secret_string"},
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
				ImportStateVerifyIgnore: []string{"secret_string"},
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
			// Create - with empty secret_string
			{
				Config: config.FromModels(t, model.SecretWithGenericString("test", id.DatabaseName(), id.SchemaName(), id.Name(), "")),
				Check: assertThat(t,
					objectassert.Secret(t, id).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasSecretType(string(sdk.SecretTypeGenericString)).
						HasOwner(testClient().Context.CurrentRole(t).Name()).
						HasCommentEmpty(),

					resourceassert.SecretWithGenericStringResource(t, "snowflake_secret_with_generic_string.test").
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasSecretTypeString("GENERIC_STRING").
						HasSecretStringString("").
						HasCommentString(""),

					resourceshowoutputassert.SecretShowOutput(t, "snowflake_secret_with_generic_string.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasSecretType(string(sdk.SecretTypeGenericString)).
						HasOwner(testClient().Context.CurrentRole(t).Name()).
						HasComment("").
						HasOwnerRoleType("ROLE"),
				),
			},
		},
	})
}

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
