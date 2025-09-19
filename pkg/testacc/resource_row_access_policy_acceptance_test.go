//go:build !account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_RowAccessPolicy(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	newComment := random.Comment()

	body := "case when current_role() in ('ANALYST') then true else false end"
	changedBody := "case when current_role() in ('CHANGED') then true else false end"
	argument := []sdk.TableColumnSignature{
		{
			Name: "A",
			Type: testdatatypes.DataTypeVarchar,
		},
		{
			Name: "B",
			Type: testdatatypes.DataTypeVarchar,
		},
	}
	changedArgument := []sdk.TableColumnSignature{
		{
			Name: "C",
			Type: testdatatypes.DataTypeBoolean,
		},
		{
			Name: "D",
			Type: testdatatypes.DataTypeTimestampNTZ,
		},
	}

	policyModel := model.RowAccessPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), argument, body).
		WithComment(comment)
	changedPolicyModel := model.RowAccessPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), argument, changedBody).
		WithComment(newComment)
	changedArgumentPolicyModel := model.RowAccessPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), argument, changedBody).
		WithComment(newComment).
		WithArgument(changedArgument)
	noCommentPolicyModel := model.RowAccessPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), argument, changedBody).
		WithComment("").
		WithArgument(changedArgument)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.RowAccessPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_RowAccessPolicy/complete"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, policyModel),
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString(comment).
					HasBodyString(body).
					HasArguments(argument),
					resourceshowoutputassert.RowAccessPolicyShowOutput(t, policyModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasKind(string(sdk.PolicyKindRowAccessPolicy)).
						HasName(id.Name()).
						HasOptions("").
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasSchemaName(id.SchemaName()).
						HasComment(comment),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "describe_output.0.body", body)),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "describe_output.0.return_type", "BOOLEAN")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "describe_output.0.signature.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "describe_output.0.signature.0.name", "A")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "describe_output.0.signature.0.type", testdatatypes.DefaultVarcharAsString)),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "describe_output.0.signature.1.name", "B")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "describe_output.0.signature.1.type", testdatatypes.DefaultVarcharAsString)),
				),
			},
			// change comment and expression
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_RowAccessPolicy/complete"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, changedPolicyModel),
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, changedPolicyModel.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString(newComment).
					HasBodyString(changedBody).
					HasArguments(argument),
				),
			},
			// change signature
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_RowAccessPolicy/complete"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, changedArgumentPolicyModel),
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, changedArgumentPolicyModel.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString(newComment).
					HasBodyString(changedBody).
					HasArguments(changedArgument),
				),
			},
			// external change on signature
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_RowAccessPolicy/complete"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, changedArgumentPolicyModel),
				PreConfig: func() {
					arg := sdk.NewCreateRowAccessPolicyArgsRequest("A", testdatatypes.DataTypeBoolean)
					createRequest := sdk.NewCreateRowAccessPolicyRequest(id, []sdk.CreateRowAccessPolicyArgsRequest{*arg}, "case when current_role() in ('ANALYST') then false else true end")
					testClient().RowAccessPolicy.CreateRowAccessPolicyWithRequest(t, *createRequest.WithOrReplace(sdk.Pointer(true)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, changedArgumentPolicyModel.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString(newComment).
					HasBodyString(changedBody).
					HasArguments(changedArgument),
				),
			},
			// external change on body
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_RowAccessPolicy/complete"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, changedArgumentPolicyModel),
				PreConfig: func() {
					testClient().RowAccessPolicy.Alter(t, *sdk.NewAlterRowAccessPolicyRequest(id).WithSetBody(sdk.Pointer("case when current_role() in ('EXTERNAL') then false else true end")))
				},
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, changedArgumentPolicyModel.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString(newComment).
					HasBodyString(changedBody).
					HasArguments(changedArgument),
				),
			},
			{
				ConfigVariables:   accconfig.ConfigVariablesFromModel(t, changedArgumentPolicyModel),
				ResourceName:      policyModel.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// unset comment
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_RowAccessPolicy/complete"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, noCommentPolicyModel),
				PreConfig: func() {
					testClient().RowAccessPolicy.Alter(t, *sdk.NewAlterRowAccessPolicyRequest(id).WithSetBody(sdk.Pointer("case when current_role() in ('EXTERNAL') then false else true end")))
				},
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, noCommentPolicyModel.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("").
					HasBodyString(changedBody).
					HasArguments(changedArgument),
				),
			},
			// IMPORT
			{
				ConfigVariables:   accconfig.ConfigVariablesFromModel(t, noCommentPolicyModel),
				ResourceName:      noCommentPolicyModel.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2053 is fixed
func TestAcc_RowAccessPolicy_Issue2053(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), []sdk.TableColumnSignature{
		{
			Name: "A",
			Type: testdatatypes.DataTypeVarchar,
		},
	}, body)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.95.0"),
				// these configs have "weird" format on purpose - to test against handling new lines during diff correctly
				Config: rowAccessPolicyV0950WithHeredoc(id, `    case
      when current_role() in ('ANALYST') then true
      else false
    end
`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				ExpectNonEmptyPlan: true,
			},
			{
				PreConfig:         func() { SetLegacyConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("1.0.0"),
				Config:            accconfig.FromModels(t, policyModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasNameString(id.Name()).
					HasBodyString(`case
  when current_role() in ('ANALYST') then true
  else false
end`),
				),
			},
		},
	})
}

func rowAccessPolicyV0950WithHeredoc(id sdk.SchemaObjectIdentifier, expr string) string {
	return fmt.Sprintf(`
resource "snowflake_row_access_policy" "test" {
  name     = "%s"
  database = "%s"
  schema   = "%s"
  signature = {
    A = "VARCHAR",
  }
  row_access_expression = <<-EOT
%s
EOT
}`, id.Name(), id.DatabaseName(), id.SchemaName(), expr)
}

func TestAcc_RowAccessPolicy_Rename(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	body := "case when current_role() in ('ANALYST') then true else false end"

	policyModel := model.RowAccessPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), []sdk.TableColumnSignature{
		{
			Name: "a",
			Type: testdatatypes.DataTypeVarchar,
		},
	}, body)
	renamedPolicyModel := model.RowAccessPolicy("test", newId.DatabaseName(), newId.SchemaName(), newId.Name(), []sdk.TableColumnSignature{
		{
			Name: "a",
			Type: testdatatypes.DataTypeVarchar,
		},
	}, body)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.RowAccessPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_RowAccessPolicy/basic"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, policyModel),
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasNameString(id.Name()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			// rename
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_RowAccessPolicy/basic"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, renamedPolicyModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(renamedPolicyModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, renamedPolicyModel.ResourceReference()).
					HasNameString(newId.Name()).
					HasFullyQualifiedNameString(newId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_RowAccessPolicy_InvalidDataType(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicyDynamicArguments("test", id, body)
	variableModel := accconfig.SetMapStringVariable("arguments")

	commonVariables := config.Variables{
		"arguments": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("A"),
				"type": config.StringVariable("invalid-type"),
			}),
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile(`invalid data type: invalid-type`),
			},
		},
	})
}

func TestAcc_RowAccessPolicy_DataTypeAliases(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), []sdk.TableColumnSignature{
		{
			Name: "A",
			Type: testdatatypes.DataTypeText,
		},
	}, body)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_RowAccessPolicy/basic"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, policyModel),
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasNameString(id.Name()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeText,
						},
					}),
				),
			},
		},
	})
}

func TestAcc_RowAccessPolicy_migrateFromVersion_0_95_0_LowercaseArgName(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), []sdk.TableColumnSignature{
		{
			Name: "A",
			Type: testdatatypes.DataTypeVarchar,
		},
		{
			Name: "b",
			Type: testdatatypes.DataTypeVarchar,
		},
	}, body)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.95.0"),
				Config:            rowAccessPolicyV0950(id, body),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						// expect change - arg name is lower case which causes a diff
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ExpectNonEmptyPlan: true,
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "row_access_expression", body)),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "signature.A", string(sdk.DataTypeVARCHAR))),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "signature.B", string(sdk.DataTypeVARCHAR))),
				),
			},
			{
				PreConfig:         func() { SetLegacyConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("1.0.0"),
				Config:            accconfig.FromModels(t, policyModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasBodyString(body),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.0.name", "A")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.0.type", string(sdk.DataTypeVARCHAR))),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.1.name", "b")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.1.type", string(sdk.DataTypeVARCHAR))),
				),
			},
		},
	})
}

func TestAcc_RowAccessPolicy_migrateFromVersion_0_95_0_UppercaseArgName(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), []sdk.TableColumnSignature{
		{
			Name: "A",
			Type: testdatatypes.DataTypeVarchar,
		},
		{
			Name: "B",
			Type: testdatatypes.DataTypeVarchar,
		},
	}, body)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.95.0"),
				Config:            rowAccessPolicyV0950(id, body),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						// expect change - arg name is lower case which causes a diff
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ExpectNonEmptyPlan: true,
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "row_access_expression", body)),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "signature.A", string(sdk.DataTypeVARCHAR))),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "signature.B", string(sdk.DataTypeVARCHAR))),
				),
			},
			{
				PreConfig:         func() { SetLegacyConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("1.0.0"),
				Config:            accconfig.FromModels(t, policyModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasBodyString(body),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.0.name", "A")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.0.type", string(sdk.DataTypeVARCHAR))),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.1.name", "B")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.1.type", string(sdk.DataTypeVARCHAR))),
				),
			},
		},
	})
}

func rowAccessPolicyV0950(id sdk.SchemaObjectIdentifier, expr string) string {
	return fmt.Sprintf(`
resource "snowflake_row_access_policy" "test" {
  name     = "%s"
  database = "%s"
  schema   = "%s"
  signature = {
    A = "VARCHAR",
    b = "VARCHAR",
  }
  row_access_expression = "%s"
}`, id.Name(), id.DatabaseName(), id.SchemaName(), expr)
}

func TestAcc_RowAccessPolicy_migrateToV2_0_0(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicyDynamicArguments("test", id, body)
	variableModel := accconfig.SetMapStringVariable("arguments")

	commonVariables := config.Variables{
		"arguments": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("A"),
				"type": config.StringVariable("VARCHAR"),
			}),
		),
	}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetLegacyConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("1.2.1"),
				Config:            accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables:   commonVariables,
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.0.type", "VARCHAR")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.0.name", "A")),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables:          commonVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
					}),
				),
			},
		},
	})
}

func TestAcc_RowAccessPolicy_migrateToV2_0_0_nonDefaultInConfig(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicyDynamicArguments("test", id, body)
	variableModel := accconfig.SetMapStringVariable("arguments")

	commonVariables := config.Variables{
		"arguments": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("A"),
				"type": config.StringVariable("VARCHAR(100)"),
			}),
		),
	}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetLegacyConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("1.2.1"),
				Config:            accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables:   commonVariables,
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.0.type", "VARCHAR")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.0.name", "A")),
				),
				// this shows the incorrect behavior in the 1.2.1 version (permadiff); in the next step, after apply, the plan will be empty
				ExpectNonEmptyPlan: true,
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables:          commonVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar_100,
						},
					}),
				),
			},
		},
	})
}

func TestAcc_RowAccessPolicy_dataType_argumentDefaultToSpecific(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicyDynamicArguments("test", id, body)
	variableModel := accconfig.SetMapStringVariable("arguments")

	commonVariables := config.Variables{
		"arguments": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("A"),
				"type": config.StringVariable("VARCHAR"),
			}),
		),
	}

	updatedDataType := config.Variables{
		"arguments": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("A"),
				"type": config.StringVariable("VARCHAR(100)"),
			}),
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
					}),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: updatedDataType,
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar_100,
						},
					}),
				),
			},
		},
	})
}

func TestAcc_RowAccessPolicy_dataType_externalChange(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicyDynamicArguments("test", id, body)
	variableModel := accconfig.SetMapStringVariable("arguments")

	commonVariables := config.Variables{
		"arguments": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("A"),
				"type": config.StringVariable("VARCHAR"),
			}),
		),
	}

	externalArgs := []sdk.CreateRowAccessPolicyArgsRequest{
		*sdk.NewCreateRowAccessPolicyArgsRequest("A", testdatatypes.DataTypeNumber),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.RowAccessPolicy),
		Steps: []resource.TestStep{
			{
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
					}),
				),
			},
			{
				PreConfig: func() {
					testClient().RowAccessPolicy.DropRowAccessPolicyFunc(t, id)()
					req := sdk.NewCreateRowAccessPolicyRequest(id, externalArgs, body)
					testClient().RowAccessPolicy.CreateRowAccessPolicyWithRequest(t, *req)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
					}),
				),
			},
		},
	})
}

func TestAcc_RowAccessPolicy_dataType_argumentExternalChangeSuppressed(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicyDynamicArguments("test", id, body)
	variableModel := accconfig.SetMapStringVariable("arguments")

	commonVariables := config.Variables{
		"arguments": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("A"),
				"type": config.StringVariable("VARCHAR"),
			}),
		),
	}

	externalArgs := []sdk.CreateRowAccessPolicyArgsRequest{
		*sdk.NewCreateRowAccessPolicyArgsRequest("A", testdatatypes.DataTypeVarchar_100),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.RowAccessPolicy),
		Steps: []resource.TestStep{
			{
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
					}),
				),
			},
			{
				PreConfig: func() {
					testClient().RowAccessPolicy.DropRowAccessPolicyFunc(t, id)()
					req := sdk.NewCreateRowAccessPolicyRequest(id, externalArgs, body)
					testClient().RowAccessPolicy.CreateRowAccessPolicyWithRequest(t, *req)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
					}),
				),
			},
		},
	})
}

func TestAcc_RowAccessPolicy_dataType_externalChangeMoreArguments(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicyDynamicArguments("test", id, body)
	variableModel := accconfig.SetMapStringVariable("arguments")

	commonVariables := config.Variables{
		"arguments": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("A"),
				"type": config.StringVariable("VARCHAR"),
			}),
		),
	}

	externalArgs := []sdk.CreateRowAccessPolicyArgsRequest{
		*sdk.NewCreateRowAccessPolicyArgsRequest("A", testdatatypes.DataTypeVarchar),
		*sdk.NewCreateRowAccessPolicyArgsRequest("B", testdatatypes.DataTypeVarchar),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.RowAccessPolicy),
		Steps: []resource.TestStep{
			{
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
					}),
				),
			},
			{
				PreConfig: func() {
					testClient().RowAccessPolicy.DropRowAccessPolicyFunc(t, id)()
					req := sdk.NewCreateRowAccessPolicyRequest(id, externalArgs, body)
					testClient().RowAccessPolicy.CreateRowAccessPolicyWithRequest(t, *req)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
					}),
				),
			},
		},
	})
}

func TestAcc_RowAccessPolicy_dataType_externalChangeFewerArguments(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicyDynamicArguments("test", id, body)
	variableModel := accconfig.SetMapStringVariable("arguments")

	commonVariables := config.Variables{
		"arguments": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("A"),
				"type": config.StringVariable("VARCHAR"),
			}),
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("B"),
				"type": config.StringVariable("VARCHAR"),
			}),
		),
	}

	externalArgs := []sdk.CreateRowAccessPolicyArgsRequest{
		*sdk.NewCreateRowAccessPolicyArgsRequest("A", testdatatypes.DataTypeVarchar),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.RowAccessPolicy),
		Steps: []resource.TestStep{
			{
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
						{
							Name: "B",
							Type: testdatatypes.DataTypeVarchar,
						},
					}),
				),
			},
			{
				PreConfig: func() {
					testClient().RowAccessPolicy.DropRowAccessPolicyFunc(t, id)()
					req := sdk.NewCreateRowAccessPolicyRequest(id, externalArgs, body)
					testClient().RowAccessPolicy.CreateRowAccessPolicyWithRequest(t, *req)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.RowAccessPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
						{
							Name: "B",
							Type: testdatatypes.DataTypeVarchar,
						},
					}),
				),
			},
		},
	})
}
