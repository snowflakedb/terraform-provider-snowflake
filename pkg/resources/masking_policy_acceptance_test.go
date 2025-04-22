//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	tfconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_MaskingPolicy_basic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

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
	policyModel := model.MaskingPolicy("test", argument, body, id.DatabaseName(), id.Name(), testdatatypes.DataTypeVarchar.ToSqlWithoutUnknowns(), id.SchemaName())

	resourceName := "snowflake_masking_policy.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.MaskingPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_MaskingPolicy/basic"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel),
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasReturnDataTypeString(string(sdk.DataTypeVARCHAR)).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasBodyString(body).
					HasExemptOtherPoliciesString(r.BooleanDefault).
					HasArguments(argument),
				),
			},
			// set all fields
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_MaskingPolicy/complete"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel.WithBody(body).WithComment("Terraform acceptance test").WithExemptOtherPolicies(r.BooleanTrue)),
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasExemptOtherPoliciesString(r.BooleanTrue).
					HasReturnDataTypeString(string(sdk.DataTypeVARCHAR)).
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
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.return_type", string(sdk.DataTypeVARCHAR))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.0.name", "A")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.0.type", string(sdk.DataTypeVARCHAR))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.1.name", "B")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.1.type", string(sdk.DataTypeVARCHAR))),
				),
			},
			// change fields
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_MaskingPolicy/complete"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel.WithBody(bodyWithBooleanReturnType).WithReturnDataType(string(sdk.DataTypeBoolean)).WithArgument(argumentWithChangedFirstArgumentType).WithComment("Terraform acceptance test - changed comment").WithExemptOtherPolicies(r.BooleanFalse)),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_MaskingPolicy/complete"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel.WithBody(body).WithReturnDataType(string(sdk.DataTypeVARCHAR)).WithArgument(changedArgument).WithExemptOtherPolicies(r.BooleanTrue)),
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasReturnDataTypeString(string(sdk.DataTypeVARCHAR)).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasExemptOtherPoliciesString(r.BooleanTrue).
					HasCommentString("Terraform acceptance test - changed comment").
					HasBodyString(body).
					HasArguments(changedArgument),
				),
			},
			// external change on signature
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_MaskingPolicy/complete"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel),
				PreConfig: func() {
					acc.TestClient().MaskingPolicy.CreateOrReplaceMaskingPolicyWithOptions(t, id, argument, testdatatypes.DataTypeVarchar, body, &sdk.CreateMaskingPolicyOptions{
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_MaskingPolicy/complete"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel),
				PreConfig: func() {
					acc.TestClient().MaskingPolicy.Alter(t, id, &sdk.AlterMaskingPolicyOptions{
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
				ConfigVariables:   tfconfig.ConfigVariablesFromModel(t, policyModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// unset comment
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_MaskingPolicy/complete"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel.WithComment("")),
				PreConfig: func() {
					acc.TestClient().MaskingPolicy.Alter(t, id, &sdk.AlterMaskingPolicyOptions{
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
				ConfigVariables:   tfconfig.ConfigVariablesFromModel(t, policyModel),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_MaskingPolicy_complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

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
	policyModel := model.MaskingPolicy("test", argument, body, id.DatabaseName(), id.Name(), testdatatypes.DataTypeVarchar.ToSqlWithoutUnknowns(), id.SchemaName()).WithComment("foo").WithExemptOtherPolicies(r.BooleanTrue)

	resourceName := "snowflake_masking_policy.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.MaskingPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_MaskingPolicy/complete"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel),
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasExemptOtherPoliciesString(r.BooleanTrue).
					HasReturnDataTypeString(string(sdk.DataTypeVARCHAR)).
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
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.return_type", string(sdk.DataTypeVARCHAR))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.0.name", "A")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.0.type", string(sdk.DataTypeVARCHAR))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.1.name", "B")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.signature.1.type", string(sdk.DataTypeVARCHAR))),
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
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.MaskingPolicy),
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
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "body", "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "return_data_type", string(sdk.DataTypeVARCHAR)),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "argument.#", "2"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "argument.0.name", "val"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "argument.0.type", string(sdk.DataTypeVARCHAR)),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "argument.1.name", "val2"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "argument.1.type", string(sdk.DataTypeVARCHAR)),
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
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	body := "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
	policyModel := model.MaskingPolicy("test", []sdk.TableColumnSignature{
		{
			Name: "val",
			Type: testdatatypes.DataTypeVarchar,
		},
	}, body, id.DatabaseName(), id.Name(), testdatatypes.DataTypeVarchar.ToSqlWithoutUnknowns(), id.SchemaName())

	resourceName := "snowflake_masking_policy.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.94.1"),
				Config:            maskingPolicyConfig(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "qualified_name", id.FullyQualifiedName()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				ConfigDirectory:          acc.ConfigurationDirectory("TestAcc_MaskingPolicy/basic"),
				ConfigVariables:          tfconfig.ConfigVariablesFromModel(t, policyModel),
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
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	newId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then 'true' else 'false' end"
	policyModel := model.MaskingPolicy("test", []sdk.TableColumnSignature{
		{
			Name: "a",
			Type: testdatatypes.DataTypeVarchar,
		},
	}, body, id.DatabaseName(), id.Name(), testdatatypes.DataTypeVarchar.ToSqlWithoutUnknowns(), id.SchemaName())

	resourceName := "snowflake_masking_policy.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.MaskingPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_MaskingPolicy/basic"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel),
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			// rename
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_MaskingPolicy/basic"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel.WithName(newId.Name())),
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

// TODO [this PR]: fix setup of this test
func TestAcc_MaskingPolicy_InvalidDataType(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.MaskingPolicy("test", []sdk.TableColumnSignature{
		{
			Name: "a",
			Type: "invalid-type",
		},
	}, body, id.DatabaseName(), id.Name(), string(sdk.DataTypeVARCHAR), id.SchemaName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_MaskingPolicy/basic"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel),
				ExpectError:     regexp.MustCompile(`invalid data type: invalid-type`),
			},
		},
	})
}

func TestAcc_MaskingPolicy_DataTypeAliases(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then 'ok' else '***' end"
	policyModel := model.MaskingPolicy("test", []sdk.TableColumnSignature{
		{
			Name: "a",
			Type: testdatatypes.DataTypeText,
		},
	}, body, id.DatabaseName(), id.Name(), testdatatypes.DataTypeVarchar.ToSqlWithoutUnknowns(), id.SchemaName())

	resourceName := "snowflake_masking_policy.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_MaskingPolicy/basic"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel),
				Check: assertThat(t, resourceassert.MaskingPolicyResource(t, resourceName).
					HasNameString(id.Name()).
					HasReturnDataTypeString(string(sdk.DataTypeVARCHAR)).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "a",
							Type: testdatatypes.DataTypeVarchar,
						},
					}),
				),
			},
		},
	})
}

func TestAcc_MaskingPolicy_migrateFromVersion_0_95_0(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	comment := random.Comment()
	body := "case when current_role() in ('ANALYST') then 'true' else 'false' end"
	policyModel := model.MaskingPolicy("test", []sdk.TableColumnSignature{
		{
			Name: "A",
			Type: testdatatypes.DataTypeVarchar,
		},
		{
			Name: "b",
			Type: testdatatypes.DataTypeVarchar,
		},
	}, body, id.DatabaseName(), id.Name(), testdatatypes.DataTypeVarchar.ToSqlWithoutUnknowns(), id.SchemaName()).WithComment(comment).WithExemptOtherPolicies(r.BooleanTrue)

	resourceName := "snowflake_masking_policy.test"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() { acc.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.95.0"),
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
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				ConfigDirectory:          acc.ConfigurationDirectory("TestAcc_MaskingPolicy/complete"),
				ConfigVariables:          tfconfig.ConfigVariablesFromModel(t, policyModel),
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
					HasBodyString(body).
					HasArguments([]sdk.TableColumnSignature{
						{
							Name: "A",
							Type: testdatatypes.DataTypeVarchar,
						},
						{
							Name: "b",
							Type: testdatatypes.DataTypeVarchar,
						},
					}),
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
