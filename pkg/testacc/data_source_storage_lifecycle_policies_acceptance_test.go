//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_StorageLifecyclePolicies_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	arguments := []sdk.TableColumnSignature{
		{
			Name: "VAL",
			Type: testdatatypes.DataTypeVarchar_200,
		},
	}
	expectedSignature := []sdk.TableColumnSignature{
		{
			Name: "VAL",
			Type: testdatatypes.DataTypeVarchar,
		},
	}
	body := "LENGTH(VAL) > 0"
	archiveTier := string(sdk.StorageLifecyclePolicyArchiveTierCold)
	archiveForDays := 365

	policyModel := model.StorageLifecyclePolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), arguments, body).
		WithArchiveTier(archiveTier).
		WithArchiveForDays(archiveForDays).
		WithComment(comment)

	storageLifecyclePoliciesModel := datasourcemodel.StorageLifecyclePolicies("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(policyModel.ResourceReference())

	storageLifecyclePoliciesModelWithoutDescribe := datasourcemodel.StorageLifecyclePolicies("test").
		WithWithDescribe(false).
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(policyModel.ResourceReference())

	showOutputAssertions := resourceshowoutputassert.StorageLifecyclePoliciesDatasourceShowOutput(t, "snowflake_storage_lifecycle_policies.test").
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasKind("STORAGE_LIFECYCLE_POLICY").
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasComment(comment).
		HasOwnerRoleType("ROLE").
		HasOptions(`{"ARCHIVE_FOR_DAYS":365,"ARCHIVE_TIER":"COLD"}`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, policyModel, storageLifecyclePoliciesModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(storageLifecyclePoliciesModel.DatasourceReference(), "storage_lifecycle_policies.#", "1")),
					showOutputAssertions,
					resourceshowoutputassert.StorageLifecyclePoliciesDatasourceDescribeOutput(t, "snowflake_storage_lifecycle_policies.test").
						HasName(id.Name()).
						HasSignature(expectedSignature...).
						HasReturnType(testdatatypes.DataTypeBoolean).
						HasBody(body).
						HasArchiveTier(archiveTier).
						HasArchiveForDays(archiveForDays),
				),
			},
			{
				Config: accconfig.FromModels(t, policyModel, storageLifecyclePoliciesModelWithoutDescribe),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(storageLifecyclePoliciesModelWithoutDescribe.DatasourceReference(), "storage_lifecycle_policies.#", "1")),
					showOutputAssertions,
					assert.Check(resource.TestCheckResourceAttr(storageLifecyclePoliciesModelWithoutDescribe.DatasourceReference(), "storage_lifecycle_policies.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_StorageLifecyclePolicies_Filtering(t *testing.T) {
	secondSchema, secondSchemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, sdk.NewAccountObjectIdentifier(TestDatabaseName))
	t.Cleanup(secondSchemaCleanup)

	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifierInSchema(secondSchema.ID())

	arguments := []sdk.TableColumnSignature{
		{
			Name: "VAL",
			Type: testdatatypes.DataTypeVarchar_200,
		},
	}
	body := "LENGTH(VAL) > 0"

	model1 := model.StorageLifecyclePolicy("test1", id1.DatabaseName(), id1.SchemaName(), id1.Name(), arguments, body)
	model2 := model.StorageLifecyclePolicy("test2", id2.DatabaseName(), id2.SchemaName(), id2.Name(), arguments, body)
	model3 := model.StorageLifecyclePolicy("test3", id3.DatabaseName(), id3.SchemaName(), id3.Name(), arguments, body)

	policiesModelLikeFirst := datasourcemodel.StorageLifecyclePolicies("test").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	policiesModelLikePrefix := datasourcemodel.StorageLifecyclePolicies("test").
		WithLike(prefix+"%").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	policiesModelInSchema := datasourcemodel.StorageLifecyclePolicies("test").
		WithInSchema(id1.SchemaId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, policiesModelLikeFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(policiesModelLikeFirst.DatasourceReference(), "storage_lifecycle_policies.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, policiesModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(policiesModelLikePrefix.DatasourceReference(), "storage_lifecycle_policies.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, policiesModelInSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(policiesModelInSchema.DatasourceReference(), "storage_lifecycle_policies.#", "2"),
				),
			},
		},
	})
}

func TestAcc_StorageLifecyclePolicies_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.StorageLifecyclePolicies("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}
