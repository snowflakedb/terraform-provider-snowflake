//go:build !account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_MaskingPolicy_basic(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then 'true' else 'false' end"
	changedBody := "case when current_role() in ('CHANGED') then 'foo' else 'bar' end"
	bodyWithBooleanReturnType := "case when current_role() in ('ANALYST') then true else false end"
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
	argumentWithChangedFirstArgumentType := []sdk.TableColumnSignature{
		{
			Name: "A",
			Type: testdatatypes.DataTypeBoolean,
		},
		{
			Name: "B",
			Type: testdatatypes.DataTypeVarchar,
		},
	}
	changedArgument := []sdk.TableColumnSignature{
		{
			Name: "C",
			Type: testdatatypes.DataTypeVarchar,
		},
		{
			Name: "D",
			Type: testdatatypes.DataTypeTimestampNTZ,
		},
	}
	policyModel := model.MaskingPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), argument, body, testdatatypes.DataTypeVarchar.ToSqlWithoutUnknowns())

	resourceName := "snowflake_masking_policy.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.MaskingPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_MaskingPolicy/basic"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, policyModel),
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasReturnDataTypeString(testdatatypes.DefaultVarcharAsString).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasBodyString(body).
					HasExemptOtherPoliciesString(r.BooleanDefault).
					HasArguments(argument),
				),
			},
			// set all fields
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_MaskingPolicy/complete"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, policyModel.WithBody(body).WithComment("Terraform acceptance test").WithExemptOtherPolicies(r.BooleanTrue)),
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasExemptOtherPoliciesString(r.BooleanTrue).
					HasReturnDataTypeString(testdatatypes.DefaultVarcharAsString).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("Terraform acceptance test").
					HasBodyString(body).
					HasArguments(argument),
					resourceshowoutputassert.MaskingPolicyShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasKind(string(sdk.PolicyKindMaskingPolicy)).
						HasName(id.Name()).
						HasExemptOtherPolicies(true).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasSchemaName(id.SchemaName()),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.body", body)),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.return_type", testdatatypes.DefaultVarcharAsString)),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.0.name", "A")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.0.type", testdatatypes.DefaultVarcharAsString)),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.1.name", "B")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.1.type", testdatatypes.DefaultVarcharAsString)),
				),
			},
			// change fields
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_MaskingPolicy/complete"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, policyModel.WithBody(bodyWithBooleanReturnType).WithReturnDataType(string(sdk.DataTypeBoolean)).WithArgument(argumentWithChangedFirstArgumentType).WithComment("Terraform acceptance test - changed comment").WithExemptOtherPolicies(r.BooleanFalse)),
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasReturnDataTypeString(string(sdk.DataTypeBoolean)).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasExemptOtherPoliciesString(r.BooleanFalse).
					HasCommentString("Terraform acceptance test - changed comment").
					HasBodyString(bodyWithBooleanReturnType).
					HasArguments(argumentWithChangedFirstArgumentType),
				),
			},
			// restore previous types - first argument type, return_type, and returned value in `body` must be the same type
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_MaskingPolicy/complete"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, policyModel.WithBody(body).WithReturnDataType(string(sdk.DataTypeVARCHAR)).WithArgument(changedArgument).WithExemptOtherPolicies(r.BooleanTrue)),
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasReturnDataTypeString(testdatatypes.DefaultVarcharAsString).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasExemptOtherPoliciesString(r.BooleanTrue).
					HasCommentString("Terraform acceptance test - changed comment").
					HasBodyString(body).
					HasArguments(changedArgument),
				),
			},
			// external change on signature
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_MaskingPolicy/complete"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, policyModel),
				PreConfig: func() {
					testClient().MaskingPolicy.CreateOrReplaceMaskingPolicyWithOptions(t, id, argument, testdatatypes.DataTypeVarchar, body, &sdk.CreateMaskingPolicyOptions{
						ExemptOtherPolicies: sdk.Pointer(false),
						Comment:             sdk.Pointer("Terraform acceptance test - changed comment"),
						OrReplace:           sdk.Pointer(true),
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("Terraform acceptance test - changed comment").
					HasBodyString(body).
					HasArguments(changedArgument),
				),
			},
			// external change on body and exempt other policies
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_MaskingPolicy/complete"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, policyModel),
				PreConfig: func() {
					testClient().MaskingPolicy.Alter(t, id, &sdk.AlterMaskingPolicyOptions{
						Set: &sdk.MaskingPolicySet{
							Body: sdk.Pointer(changedBody),
						},
					})
				},
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("Terraform acceptance test - changed comment").
					HasBodyString(body).
					HasArguments(changedArgument),
				),
			},
			{
				ConfigVariables:   accconfig.ConfigVariablesFromModel(t, policyModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// unset comment
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_MaskingPolicy/complete"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, policyModel.WithComment("")),
				PreConfig: func() {
					testClient().MaskingPolicy.Alter(t, id, &sdk.AlterMaskingPolicyOptions{
						Unset: &sdk.MaskingPolicyUnset{
							Comment: sdk.Pointer(true),
						},
					})
				},
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("").
					HasBodyString(body).
					HasArguments(changedArgument),
				),
			},
			// IMPORT
			{
				ConfigVariables:   accconfig.ConfigVariablesFromModel(t, policyModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_MaskingPolicy_complete(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then 'true' else 'false' end"
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
	policyModel := model.MaskingPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), argument, body, testdatatypes.DataTypeVarchar.ToSqlWithoutUnknowns()).WithComment("foo").WithExemptOtherPolicies(r.BooleanTrue)

	resourceName := "snowflake_masking_policy.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.MaskingPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_MaskingPolicy/complete"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, policyModel),
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasExemptOtherPoliciesString(r.BooleanTrue).
					HasReturnDataTypeString(testdatatypes.DefaultVarcharAsString).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasCommentString("foo").
					HasBodyString(body).
					HasArguments(argument),
					resourceshowoutputassert.MaskingPolicyShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasKind(string(sdk.PolicyKindMaskingPolicy)).
						HasName(id.Name()).
						HasExemptOtherPolicies(true).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasSchemaName(id.SchemaName()),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.body", body)),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.return_type", testdatatypes.DefaultVarcharAsString)),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.0.name", "A")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.0.type", testdatatypes.DefaultVarcharAsString)),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.1.name", "B")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.1.type", testdatatypes.DefaultVarcharAsString)),
				),
			},
		},
	})
}

func maskingPolicyConfig(maskingPolicyId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_masking_policy" "test" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
	signature {
		column {
			name = "val"
			type = "VARCHAR"
		}
	}
	masking_expression = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
	return_data_type = "VARCHAR"
}
`, maskingPolicyId.DatabaseName(), maskingPolicyId.SchemaName(), maskingPolicyId.Name())
}

func TestAcc_MaskingPolicyMultiColumns(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.MaskingPolicy),
		Steps: []resource.TestStep{
			{
				Config: maskingPolicyConfigMultiColumn(id),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "body", "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "return_data_type", testdatatypes.DefaultVarcharAsString),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "argument.#", "2"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "argument.0.name", "val"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "argument.0.type", testdatatypes.DefaultVarcharAsString),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "argument.1.name", "val2"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "argument.1.type", testdatatypes.DefaultVarcharAsString),
				),
			},
		},
	})
}

func maskingPolicyConfigMultiColumn(maskingPolicyId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_masking_policy" "test" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
	argument {
		name = "val"
		type = "VARCHAR"
	}

	argument {
		name = "val2"
		type = "VARCHAR"
	}
	body = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
	return_data_type = "VARCHAR"
}
`, maskingPolicyId.DatabaseName(), maskingPolicyId.SchemaName(), maskingPolicyId.Name())
}

func TestAcc_MaskingPolicy_migrateFromVersion_0_94_1(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	body := "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
	policyModel := model.MaskingPolicyDynamicArguments("test", id, body, sdk.DataTypeVARCHAR)
	variableModel := accconfig.SetMapStringVariable("arguments")

	commonVariables := config.Variables{
		"arguments": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("val"),
				"type": config.StringVariable("VARCHAR"),
			}),
		),
	}

	resourceName := "snowflake_masking_policy.test"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            maskingPolicyConfig(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "qualified_name", id.FullyQualifiedName()),
				),
			},
			{
				PreConfig:         func() { SetLegacyConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("1.0.0"),
				Config:            accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables:   commonVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceName, "qualified_name"),
				),
			},
		},
	})
}

func TestAcc_MaskingPolicy_Rename(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then 'true' else 'false' end"
	policyModel := model.MaskingPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), []sdk.TableColumnSignature{
		{
			Name: "a",
			Type: testdatatypes.DataTypeVarchar,
		},
	}, body, testdatatypes.DataTypeVarchar.ToSqlWithoutUnknowns())

	resourceName := "snowflake_masking_policy.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.MaskingPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_MaskingPolicy/basic"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, policyModel),
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			// rename
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_MaskingPolicy/basic"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, policyModel.WithName(newId.Name())),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(newId.Name()).
					HasFullyQualifiedNameString(newId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_MaskingPolicy_InvalidDataType(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.MaskingPolicyDynamicArguments("test", id, body, sdk.DataTypeVARCHAR)
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

func TestAcc_MaskingPolicy_DataTypeAliases(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then 'ok' else '***' end"
	policyModel := model.MaskingPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), []sdk.TableColumnSignature{
		{
			Name: "a",
			Type: testdatatypes.DataTypeText,
		},
	}, body, testdatatypes.DataTypeVarchar.ToSqlWithoutUnknowns())

	resourceName := "snowflake_masking_policy.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_MaskingPolicy/basic"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, policyModel),
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasReturnDataTypeString(testdatatypes.DefaultVarcharAsString).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "a",
							Type: testdatatypes.DataTypeText,
						},
					}),
				),
			},
		},
	})
}

func TestAcc_MaskingPolicy_migrateFromVersion_0_95_0(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	comment := random.Comment()
	body := "case when current_role() in ('ANALYST') then 'true' else 'false' end"
	policyModel := model.MaskingPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), []sdk.TableColumnSignature{
		{
			Name: "A",
			Type: testdatatypes.DataTypeVarchar,
		},
		{
			Name: "b",
			Type: testdatatypes.DataTypeVarchar,
		},
	}, body, testdatatypes.DataTypeVarchar.ToSqlWithoutUnknowns()).WithComment(comment).WithExemptOtherPolicies(r.BooleanTrue)

	resourceName := "snowflake_masking_policy.test"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.95.0"),
				Config:            maskingPolicyV0950(id, body, comment),
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString(comment).
					HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(id))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "masking_expression", body)),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "signature.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "signature.0.column.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "signature.0.column.0.name", "A")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "signature.0.column.0.type", string(sdk.DataTypeVARCHAR))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "signature.0.column.1.name", "b")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "signature.0.column.1.type", string(sdk.DataTypeVARCHAR))),
				),
			},
			{
				PreConfig:         func() { SetLegacyConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("1.0.0"),
				Config:            accconfig.FromModels(t, policyModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString(comment).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasBodyString(body),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "argument.0.name", "A")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "argument.0.type", string(sdk.DataTypeVARCHAR))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "argument.1.name", "b")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "argument.1.type", string(sdk.DataTypeVARCHAR))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "id", id.FullyQualifiedName())),
					assert.Check(resource.TestCheckNoResourceAttr(resourceName, "masking_expression")),
				),
			},
		},
	})
}

func maskingPolicyV0950(id sdk.SchemaObjectIdentifier, expr, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_masking_policy" "test" {
  database = "%[1]s"
  schema   = "%[2]s"
  name     = "%[3]s"
  signature {
	column {
      name = "A"
      type = "VARCHAR"
    }
	column {
      name = "b"
      type = "VARCHAR"
    }
  }
  return_data_type = "VARCHAR"
  masking_expression = "%[4]s"
  exempt_other_policies = true
  comment = "%[5]s"
}`, id.DatabaseName(), id.SchemaName(), id.Name(), expr, comment)
}

func TestAcc_MaskingPolicy_migrateToV2_0_0(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then 'true' else 'false' end"
	policyModel := model.MaskingPolicyDynamicArguments("test", id, body, sdk.DataTypeVARCHAR)
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
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.0.type", "VARCHAR")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.0.name", "A")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "return_data_type", "VARCHAR")),
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
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
					}).HasReturnDataTypeString(testdatatypes.DefaultVarcharAsString),
				),
			},
		},
	})
}

func TestAcc_MaskingPolicy_migrateToV2_0_0_nonDefaultInConfig(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then 'true' else 'false' end"
	policyModel := model.MaskingPolicyDynamicArguments("test", id, body, sdk.DataType("VARCHAR(100)"))
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
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.0.type", "VARCHAR")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "argument.0.name", "A")),
					assert.Check(resource.TestCheckResourceAttr(policyModel.ResourceReference(), "return_data_type", "VARCHAR")),
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
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar_100,
						},
					}).HasReturnDataTypeString("VARCHAR(100)"),
				),
			},
		},
	})
}

func TestAcc_MaskingPolicy_dataType_argumentDefaultToSpecific(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then 'true' else 'false' end"
	policyModel := model.MaskingPolicyDynamicArguments("test", id, body, sdk.DataTypeVARCHAR)
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
		CheckDestroy: CheckDestroy(t, resources.MaskingPolicy),
		Steps: []resource.TestStep{
			{
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, policyModel.ResourceReference()).
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
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, policyModel.ResourceReference()).
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

func TestAcc_MaskingPolicy_dataType_returnTypeDefaultToSpecific(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then 'true' else 'false' end"
	policyModel := model.MaskingPolicyDynamicArguments("test", id, body, sdk.DataTypeVARCHAR)
	updatedPolicyModel := model.MaskingPolicyDynamicArguments("test", id, body, sdk.DataType("VARCHAR(100)"))
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
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.MaskingPolicy),
		Steps: []resource.TestStep{
			{
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasReturnDataTypeString(testdatatypes.DefaultVarcharAsString),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config:          accconfig.FromModels(t, variableModel, updatedPolicyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasReturnDataTypeString("VARCHAR(100)"),
				),
			},
		},
	})
}

func TestAcc_MaskingPolicy_dataType_externalChange(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then 'true' else 'false' end"
	updatedBody := "case when current_role() in ('ANALYST') then 1 else 2 end"
	policyModel := model.MaskingPolicyDynamicArguments("test", id, body, sdk.DataTypeVARCHAR)
	variableModel := accconfig.SetMapStringVariable("arguments")

	commonVariables := config.Variables{
		"arguments": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("A"),
				"type": config.StringVariable("VARCHAR"),
			}),
		),
	}

	externalSignature := []sdk.TableColumnSignature{
		{
			Name: "A",
			Type: testdatatypes.DataTypeNumber,
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.MaskingPolicy),
		Steps: []resource.TestStep{
			{
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
					}).
					HasReturnDataTypeString(testdatatypes.DefaultVarcharAsString),
				),
			},
			{
				PreConfig: func() {
					// had to update the body too, as otherwise error is returned
					// 090207 (42601): Declared return type 'NUMBER(38,0)' is incompatible with actual return type 'VARCHAR(5)'
					// we could later suppress the body changes to be sure what triggers the changes
					// TODO [SNOW-2298286]: suppress the body changes
					testClient().MaskingPolicy.CreateOrReplaceMaskingPolicyWithOptions(t, id, externalSignature, testdatatypes.DataTypeNumber, updatedBody, &sdk.CreateMaskingPolicyOptions{})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
					}).
					HasReturnDataTypeString(testdatatypes.DefaultVarcharAsString),
				),
			},
		},
	})
}

func TestAcc_MaskingPolicy_dataType_argumentExternalChangeSuppressed(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then 'true' else 'false' end"
	policyModel := model.MaskingPolicyDynamicArguments("test", id, body, sdk.DataTypeVARCHAR)
	variableModel := accconfig.SetMapStringVariable("arguments")

	commonVariables := config.Variables{
		"arguments": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("A"),
				"type": config.StringVariable("VARCHAR"),
			}),
		),
	}

	externalSignature := []sdk.TableColumnSignature{
		{
			Name: "A",
			Type: testdatatypes.DataTypeVarchar_100,
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.MaskingPolicy),
		Steps: []resource.TestStep{
			{
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
					}).
					HasReturnDataTypeString(testdatatypes.DefaultVarcharAsString),
				),
			},
			{
				PreConfig: func() {
					testClient().MaskingPolicy.CreateOrReplaceMaskingPolicyWithOptions(t, id, externalSignature, testdatatypes.DataTypeVarchar, body, &sdk.CreateMaskingPolicyOptions{})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
					}).
					HasReturnDataTypeString(testdatatypes.DefaultVarcharAsString),
				),
			},
		},
	})
}

func TestAcc_MaskingPolicy_dataType_returnTypeExternalChange(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then 'true' else 'false' end"
	policyModel := model.MaskingPolicyDynamicArguments("test", id, body, sdk.DataTypeVARCHAR)
	variableModel := accconfig.SetMapStringVariable("arguments")

	commonVariables := config.Variables{
		"arguments": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("A"),
				"type": config.StringVariable("VARCHAR"),
			}),
		),
	}

	externalSignature := []sdk.TableColumnSignature{
		{
			Name: "A",
			Type: testdatatypes.DataTypeVarchar,
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.MaskingPolicy),
		Steps: []resource.TestStep{
			{
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
					}).
					HasReturnDataTypeString(testdatatypes.DefaultVarcharAsString),
				),
			},
			{
				PreConfig: func() {
					testClient().MaskingPolicy.CreateOrReplaceMaskingPolicyWithOptions(t, id, externalSignature, testdatatypes.DataTypeVarchar_100, body, &sdk.CreateMaskingPolicyOptions{})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(policyModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config:          accconfig.FromModels(t, variableModel, policyModel),
				ConfigVariables: commonVariables,
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, policyModel.ResourceReference()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
					}).
					HasReturnDataTypeString(testdatatypes.DefaultVarcharAsString),
				),
			},
		},
	})
}
