//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_MaskingPolicies(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then 'true' else 'false' end"
	policyModel := model.MaskingPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), []sdk.TableColumnSignature{
		{
			Name: "a",
			Type: testdatatypes.DataTypeVarchar,
		},
		{
			Name: "b",
			Type: testdatatypes.DataTypeVarchar,
		},
	}, body, testdatatypes.DataTypeVarchar.ToSqlWithoutUnknowns()).WithComment("foo")

	dsName := "data.snowflake_masking_policies.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.MaskingPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_MaskingPolicies/optionals_set"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, policyModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsName, "masking_policies.#", "1")),

					resourceshowoutputassert.MaskingPoliciesDatasourceShowOutput(t, "snowflake_masking_policies.test").
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasKind(string(sdk.PolicyKindMaskingPolicy)).
						HasName(id.Name()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasSchemaName(id.SchemaName()).
						HasExemptOtherPolicies(false).
						HasComment("foo"),

					assert.Check(resource.TestCheckResourceAttr(dsName, "masking_policies.0.describe_output.0.body", body)),
					assert.Check(resource.TestCheckResourceAttr(dsName, "masking_policies.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "masking_policies.0.describe_output.0.return_type", testdatatypes.DefaultVarcharAsString)),
					assert.Check(resource.TestCheckResourceAttr(dsName, "masking_policies.0.describe_output.0.signature.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "masking_policies.0.describe_output.0.signature.0.name", "a")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "masking_policies.0.describe_output.0.signature.0.type", testdatatypes.DefaultVarcharAsString)),
					assert.Check(resource.TestCheckResourceAttr(dsName, "masking_policies.0.describe_output.0.signature.1.name", "b")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "masking_policies.0.describe_output.0.signature.1.type", testdatatypes.DefaultVarcharAsString)),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_MaskingPolicies/optionals_unset"),
				ConfigVariables: accconfig.ConfigVariablesFromModel(t, policyModel),

				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsName, "masking_policies.#", "1")),

					resourceshowoutputassert.MaskingPoliciesDatasourceShowOutput(t, "snowflake_masking_policies.test").
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasKind(string(sdk.PolicyKindMaskingPolicy)).
						HasName(id.Name()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasSchemaName(id.SchemaName()).
						HasExemptOtherPolicies(false).
						HasComment("foo"),
					assert.Check(resource.TestCheckResourceAttr(dsName, "masking_policies.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_MaskingPolicies_Filtering(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	idTwo := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	idThree := testClient().Ids.RandomSchemaObjectIdentifier()
	body := "case when current_role() in ('ANALYST') then 'true' else 'false' end"

	maskingPolicyModel1 := model.MaskingPolicyDynamicArguments("test_1", idOne, body, sdk.DataTypeVARCHAR)
	maskingPolicyModel2 := model.MaskingPolicyDynamicArguments("test_2", idTwo, body, sdk.DataTypeVARCHAR)
	maskingPolicyModel3 := model.MaskingPolicyDynamicArguments("test_3", idThree, body, sdk.DataTypeVARCHAR)
	maskingPoliciesModelLikeFirstOne := datasourcemodel.MaskingPolicies("test").
		WithLike(idOne.Name()).
		WithDependsOn(maskingPolicyModel1.ResourceReference(), maskingPolicyModel2.ResourceReference(), maskingPolicyModel3.ResourceReference())
	maskingPoliciesModelLikePrefix := datasourcemodel.MaskingPolicies("test").
		WithLike(prefix+"%").
		WithDependsOn(maskingPolicyModel1.ResourceReference(), maskingPolicyModel2.ResourceReference(), maskingPolicyModel3.ResourceReference())
	variableModel := accconfig.SetMapStringVariable("arguments")

	commonVariables := config.Variables{
		"arguments": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("a"),
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
				Config:          accconfig.FromModels(t, variableModel, maskingPolicyModel1, maskingPolicyModel2, maskingPolicyModel3, maskingPoliciesModelLikeFirstOne),
				ConfigVariables: commonVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_masking_policies.test", "masking_policies.#", "1"),
				),
			},
			{
				Config:          accconfig.FromModels(t, variableModel, maskingPolicyModel1, maskingPolicyModel2, maskingPolicyModel3, maskingPoliciesModelLikePrefix),
				ConfigVariables: commonVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_masking_policies.test", "masking_policies.#", "2"),
				),
			},
		},
	})
}

func TestAcc_MaskingPolicies_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      maskingPoliciesDatasourceEmptyIn(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func maskingPoliciesDatasourceEmptyIn() string {
	return `
data "snowflake_masking_policies" "test" {
  in {
  }
}
`
}

func TestAcc_MaskingPolicies_NotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_MaskingPolicies/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one masking policy"),
			},
		},
	})
}
