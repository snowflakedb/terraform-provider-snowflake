//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_RowAccessPolicies_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	rowAccessPolicyId1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	rowAccessPolicyId2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	rowAccessPolicyId3 := testClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then true else false end"

	columnSignature := func(colName string) []sdk.TableColumnSignature {
		return []sdk.TableColumnSignature{
			{
				Name: colName,
				Type: testdatatypes.DataTypeVarchar,
			},
		}
	}
	rowAccessPolicyModel1 := model.RowAccessPolicy("test", rowAccessPolicyId1.DatabaseName(), rowAccessPolicyId1.SchemaName(), rowAccessPolicyId1.Name(), columnSignature("a"), body)
	rowAccessPolicyModel2 := model.RowAccessPolicy("test1", rowAccessPolicyId2.DatabaseName(), rowAccessPolicyId2.SchemaName(), rowAccessPolicyId2.Name(), columnSignature("b"), body)
	rowAccessPolicyModel3 := model.RowAccessPolicy("test2", rowAccessPolicyId3.DatabaseName(), rowAccessPolicyId3.SchemaName(), rowAccessPolicyId3.Name(), columnSignature("c"), body)

	datasourceModelLikeExact := datasourcemodel.RowAccessPolicies("test").
		WithLike(rowAccessPolicyId1.Name()).
		WithDependsOn(rowAccessPolicyModel1.ResourceReference(), rowAccessPolicyModel2.ResourceReference(), rowAccessPolicyModel3.ResourceReference())

	datasourceModelLikePrefix := datasourcemodel.RowAccessPolicies("test").
		WithLike(prefix+"%").
		WithDependsOn(rowAccessPolicyModel1.ResourceReference(), rowAccessPolicyModel2.ResourceReference(), rowAccessPolicyModel3.ResourceReference())

	datasourceModelInDatabase := datasourcemodel.RowAccessPolicies("test").
		WithInDatabase(rowAccessPolicyId1.DatabaseId()).
		WithDependsOn(rowAccessPolicyModel1.ResourceReference(), rowAccessPolicyModel2.ResourceReference(), rowAccessPolicyModel3.ResourceReference())

	datasourceModelInSchema := datasourcemodel.RowAccessPolicies("test").
		WithInSchema(rowAccessPolicyId1.SchemaId()).
		WithDependsOn(rowAccessPolicyModel1.ResourceReference(), rowAccessPolicyModel2.ResourceReference(), rowAccessPolicyModel3.ResourceReference())

	datasourceModelLikeInDatabase := datasourcemodel.RowAccessPolicies("test").
		WithLike(prefix+"%").
		WithInDatabase(rowAccessPolicyId1.DatabaseId()).
		WithDependsOn(rowAccessPolicyModel1.ResourceReference(), rowAccessPolicyModel2.ResourceReference(), rowAccessPolicyModel3.ResourceReference())

	datasourceModelLikeInSchema := datasourcemodel.RowAccessPolicies("test").
		WithLike(prefix+"%").
		WithInSchema(rowAccessPolicyId1.SchemaId()).
		WithDependsOn(rowAccessPolicyModel1.ResourceReference(), rowAccessPolicyModel2.ResourceReference(), rowAccessPolicyModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.RowAccessPolicy),
		Steps: []resource.TestStep{
			// like (exact)
			{
				Config: config.FromModels(t, rowAccessPolicyModel1, rowAccessPolicyModel2, rowAccessPolicyModel3, datasourceModelLikeExact),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikeExact.DatasourceReference(), "row_access_policies.#", "1"),
					resource.TestCheckResourceAttr(datasourceModelLikeExact.DatasourceReference(), "row_access_policies.0.show_output.0.name", rowAccessPolicyId1.Name()),
				),
			},
			// like (prefix)
			{
				Config: config.FromModels(t, rowAccessPolicyModel1, rowAccessPolicyModel2, rowAccessPolicyModel3, datasourceModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikePrefix.DatasourceReference(), "row_access_policies.#", "2"),
				),
			},
			// in database
			{
				Config: config.FromModels(t, rowAccessPolicyModel1, rowAccessPolicyModel2, rowAccessPolicyModel3, datasourceModelInDatabase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelInDatabase.DatasourceReference(), "row_access_policies.#", "3"),
				),
			},
			// in schema
			{
				Config: config.FromModels(t, rowAccessPolicyModel1, rowAccessPolicyModel2, rowAccessPolicyModel3, datasourceModelInSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelInSchema.DatasourceReference(), "row_access_policies.#", "3"),
				),
			},
			// like + in database
			{
				Config: config.FromModels(t, rowAccessPolicyModel1, rowAccessPolicyModel2, rowAccessPolicyModel3, datasourceModelLikeInDatabase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikeInDatabase.DatasourceReference(), "row_access_policies.#", "2"),
				),
			},
			// like + in schema
			{
				Config: config.FromModels(t, rowAccessPolicyModel1, rowAccessPolicyModel2, rowAccessPolicyModel3, datasourceModelLikeInSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikeInSchema.DatasourceReference(), "row_access_policies.#", "2"),
				),
			},
		},
	})
}

func TestAcc_RowAccessPolicies_CompleteUseCase(t *testing.T) {
	objectNamePrefix := random.AlphaN(10)
	id := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(objectNamePrefix + "1")
	_ = testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(objectNamePrefix + "2")
	_ = testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	body := "case when current_role() in ('ANALYST') then true else false end"

	rowAccessPolicyModel := model.RowAccessPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), []sdk.TableColumnSignature{
		{
			Name: "a",
			Type: testdatatypes.DataTypeVarchar,
		},
	}, body).WithComment(comment)

	datasourceModelWithoutDescribe := datasourcemodel.RowAccessPolicies("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithDependsOn(rowAccessPolicyModel.ResourceReference())

	datasourceModelWithDescribe := datasourcemodel.RowAccessPolicies("test").
		WithLike(id.Name()).
		WithWithDescribe(true).
		WithDependsOn(rowAccessPolicyModel.ResourceReference())

	commonShowOutputAsserts := func(t *testing.T, datasourceReference string) *resourceshowoutputassert.RowAccessPolicyShowOutputAssert {
		return resourceshowoutputassert.RowAccessPoliciesDatasourceShowOutput(t, datasourceReference).
			HasCreatedOnNotEmpty().
			HasDatabaseName(id.DatabaseName()).
			HasKind(string(sdk.PolicyKindRowAccessPolicy)).
			HasName(id.Name()).
			HasOptions("").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasSchemaName(id.SchemaName()).
			HasComment(comment)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.RowAccessPolicy),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, rowAccessPolicyModel, datasourceModelWithoutDescribe),
				Check: assertThat(t,
					commonShowOutputAsserts(t, datasourceModelWithoutDescribe.DatasourceReference()),

					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutDescribe.DatasourceReference(), "row_access_policies.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutDescribe.DatasourceReference(), "row_access_policies.0.describe_output.#", "0")),
				),
			},
			{
				Config: config.FromModels(t, rowAccessPolicyModel, datasourceModelWithDescribe),
				Check: assertThat(t,
					commonShowOutputAsserts(t, datasourceModelWithDescribe.DatasourceReference()),

					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithDescribe.DatasourceReference(), "row_access_policies.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithDescribe.DatasourceReference(), "row_access_policies.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithDescribe.DatasourceReference(), "row_access_policies.0.describe_output.0.body", body)),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithDescribe.DatasourceReference(), "row_access_policies.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithDescribe.DatasourceReference(), "row_access_policies.0.describe_output.0.return_type", "BOOLEAN")),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithDescribe.DatasourceReference(), "row_access_policies.0.describe_output.0.signature.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithDescribe.DatasourceReference(), "row_access_policies.0.describe_output.0.signature.0.name", "a")),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithDescribe.DatasourceReference(), "row_access_policies.0.describe_output.0.signature.0.type", testdatatypes.DefaultVarcharAsString)),
				),
			},
		},
	})
}

func TestAcc_RowAccessPolicies_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      rowAccessPoliciesDatasourceEmptyIn(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func rowAccessPoliciesDatasourceEmptyIn() string {
	return `
data "snowflake_row_access_policies" "test" {
  in {
  }
}
`
}

func TestAcc_RowAccessPolicies_NotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_RowAccessPolicies/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one row access policy"),
			},
		},
	})
}
