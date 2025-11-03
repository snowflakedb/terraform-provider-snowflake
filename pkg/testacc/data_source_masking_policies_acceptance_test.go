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

func TestAcc_MaskingPolicies_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	maskingPolicyId1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	maskingPolicyId2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	maskingPolicyId3 := testClient().Ids.RandomSchemaObjectIdentifier()
	body := "case when current_role() in ('ANALYST') then 'true' else 'false' end"

	maskingPolicyModel1 := model.MaskingPolicyDynamicArguments("test_1", maskingPolicyId1, body, sdk.DataTypeVARCHAR)
	maskingPolicyModel2 := model.MaskingPolicyDynamicArguments("test_2", maskingPolicyId2, body, sdk.DataTypeVARCHAR)
	maskingPolicyModel3 := model.MaskingPolicyDynamicArguments("test_3", maskingPolicyId3, body, sdk.DataTypeVARCHAR)
	variableModel := accconfig.SetMapStringVariable("arguments")

	commonVariables := config.Variables{
		"arguments": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("a"),
				"type": config.StringVariable("VARCHAR"),
			}),
		),
	}

	datasourceModelLikePrefix := datasourcemodel.MaskingPolicies("test").
		WithLike(prefix+"%").
		WithDependsOn(maskingPolicyModel1.ResourceReference(), maskingPolicyModel2.ResourceReference(), maskingPolicyModel3.ResourceReference())

	datasourceModelInSchema := datasourcemodel.MaskingPolicies("test").
		WithInSchema(maskingPolicyId1.SchemaId()).
		WithDependsOn(maskingPolicyModel1.ResourceReference(), maskingPolicyModel2.ResourceReference(), maskingPolicyModel3.ResourceReference())

	datasourceModelLimitFrom := datasourcemodel.MaskingPolicies("test").
		WithInSchema(maskingPolicyId1.SchemaId()).
		WithRowsAndFrom(1, maskingPolicyId1.Name()).
		WithDependsOn(maskingPolicyModel1.ResourceReference(), maskingPolicyModel2.ResourceReference(), maskingPolicyModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.MaskingPolicy),
		Steps: []resource.TestStep{
			// like (prefix)
			{
				Config:          accconfig.FromModels(t, variableModel, maskingPolicyModel1, maskingPolicyModel2, maskingPolicyModel3, datasourceModelLikePrefix),
				ConfigVariables: commonVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikePrefix.DatasourceReference(), "masking_policies.#", "2"),
				),
			},
			// in (schema)
			{
				Config:          accconfig.FromModels(t, variableModel, maskingPolicyModel1, maskingPolicyModel2, maskingPolicyModel3, datasourceModelInSchema),
				ConfigVariables: commonVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelInSchema.DatasourceReference(), "masking_policies.#", "3"),
				),
			},
			// limit rows from (scoped to schema)
			{
				Config:          accconfig.FromModels(t, variableModel, maskingPolicyModel1, maskingPolicyModel2, maskingPolicyModel3, datasourceModelLimitFrom),
				ConfigVariables: commonVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLimitFrom.DatasourceReference(), "masking_policies.#", "1"),
				),
			},
		},
	})
}

func TestAcc_MaskingPolicies_CompleteUseCase(t *testing.T) {
	maskingPolicyId := testClient().Ids.RandomSchemaObjectIdentifier()
	body := "case when current_role() in ('ANALYST') then 'true' else 'false' end"
	comment := random.Comment()

	maskingPolicyModel := model.MaskingPolicy("test", maskingPolicyId.DatabaseName(), maskingPolicyId.SchemaName(), maskingPolicyId.Name(), []sdk.TableColumnSignature{
		{
			Name: "a",
			Type: testdatatypes.DataTypeVarchar,
		},
		{
			Name: "b",
			Type: testdatatypes.DataTypeVarchar,
		},
	}, body, testdatatypes.DataTypeVarchar.ToSqlWithoutUnknowns()).WithComment(comment)

	withoutDescribe := datasourcemodel.MaskingPolicies("test").
		WithWithDescribe(false).
		WithLike(maskingPolicyId.Name()).
		WithDependsOn(maskingPolicyModel.ResourceReference())

	withDescribe := datasourcemodel.MaskingPolicies("test").
		WithWithDescribe(true).
		WithLike(maskingPolicyId.Name()).
		WithDependsOn(maskingPolicyModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.MaskingPolicy),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, maskingPolicyModel, withoutDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.MaskingPoliciesDatasourceShowOutput(t, withoutDescribe.DatasourceReference()).
						HasCreatedOnNotEmpty().
						HasDatabaseName(maskingPolicyId.DatabaseName()).
						HasKind(string(sdk.PolicyKindMaskingPolicy)).
						HasName(maskingPolicyId.Name()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleTypeNotEmpty().
						HasSchemaName(maskingPolicyId.SchemaName()).
						HasExemptOtherPolicies(false).
						HasComment(comment),

					assert.Check(resource.TestCheckResourceAttr("data.snowflake_masking_policies.test", "masking_policies.0.describe_output.#", "0")),
				),
			},
			{
				Config: accconfig.FromModels(t, maskingPolicyModel, withDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.MaskingPoliciesDatasourceShowOutput(t, withDescribe.DatasourceReference()).
						HasCreatedOnNotEmpty().
						HasDatabaseName(maskingPolicyId.DatabaseName()).
						HasKind(string(sdk.PolicyKindMaskingPolicy)).
						HasName(maskingPolicyId.Name()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleTypeNotEmpty().
						HasSchemaName(maskingPolicyId.SchemaName()).
						HasExemptOtherPolicies(false).
						HasComment(comment),

					assert.Check(resource.TestCheckResourceAttr("data.snowflake_masking_policies.test", "masking_policies.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_masking_policies.test", "masking_policies.0.describe_output.0.body", body)),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_masking_policies.test", "masking_policies.0.describe_output.0.name", maskingPolicyId.Name())),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_masking_policies.test", "masking_policies.0.describe_output.0.return_type", testdatatypes.DefaultVarcharAsString)),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_masking_policies.test", "masking_policies.0.describe_output.0.signature.#", "2")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_masking_policies.test", "masking_policies.0.describe_output.0.signature.0.name", "a")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_masking_policies.test", "masking_policies.0.describe_output.0.signature.0.type", testdatatypes.DefaultVarcharAsString)),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_masking_policies.test", "masking_policies.0.describe_output.0.signature.1.name", "b")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_masking_policies.test", "masking_policies.0.describe_output.0.signature.1.type", testdatatypes.DefaultVarcharAsString)),
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
