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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_BackupPolicies_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	policyModel := model.BackupPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithExpireAfterDays(7).
		WithComment(comment).
		WithWithRetentionLockBool(false)

	backupPoliciesModel := datasourcemodel.BackupPolicies("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(policyModel.ResourceReference())

	backupPoliciesModelWithoutDescribe := datasourcemodel.BackupPolicies("test").
		WithWithDescribe(false).
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(policyModel.ResourceReference())

	showOutputAssertions := resourceshowoutputassert.BackupPoliciesDatasourceShowOutput(t, backupPoliciesModel.DatasourceReference()).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasComment(comment).
		HasSchedule("").
		HasExpireAfterDays(7).
		HasHasRetentionLock(false).
		HasOwnerRoleType("ROLE")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, policyModel, backupPoliciesModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(backupPoliciesModel.DatasourceReference(), "backup_policies.#", "1")),
					showOutputAssertions,
					resourceshowoutputassert.BackupPoliciesDatasourceDescribeOutput(t, backupPoliciesModel.DatasourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasSchedule("").
						HasExpireAfterDays(7).
						HasHasRetentionLock(false).
						HasOwnerRoleType("ROLE"),
				),
			},
			{
				Config: accconfig.FromModels(t, policyModel, backupPoliciesModelWithoutDescribe),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(backupPoliciesModelWithoutDescribe.DatasourceReference(), "backup_policies.#", "1")),
					showOutputAssertions,
					assert.Check(resource.TestCheckResourceAttr(backupPoliciesModelWithoutDescribe.DatasourceReference(), "backup_policies.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_BackupPolicies_Filtering(t *testing.T) {
	secondSchema, secondSchemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, sdk.NewAccountObjectIdentifier(TestDatabaseName))
	t.Cleanup(secondSchemaCleanup)

	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, secondSchema.ID())

	model1 := model.BackupPolicy("test1", id1.DatabaseName(), id1.SchemaName(), id1.Name()).
		WithExpireAfterDays(7).
		WithWithRetentionLockBool(false)
	model2 := model.BackupPolicy("test2", id2.DatabaseName(), id2.SchemaName(), id2.Name()).
		WithExpireAfterDays(7).
		WithWithRetentionLockBool(false)
	model3 := model.BackupPolicy("test3", id3.DatabaseName(), id3.SchemaName(), id3.Name()).
		WithExpireAfterDays(7).
		WithWithRetentionLockBool(false)

	likeFirst := datasourcemodel.BackupPolicies("test").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	likePrefix := datasourcemodel.BackupPolicies("test").
		WithLike(prefix+"%").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	inSchema := datasourcemodel.BackupPolicies("test").
		WithInSchema(id1.SchemaId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	inDatabase := datasourcemodel.BackupPolicies("test").
		WithStartsWith(prefix).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	withLimit := datasourcemodel.BackupPolicies("test").
		WithStartsWith(prefix).
		WithInDatabase(id1.DatabaseId()).
		WithLimit(1).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, likeFirst),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(likeFirst.DatasourceReference(), "backup_policies.#", "1")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, likePrefix),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(likePrefix.DatasourceReference(), "backup_policies.#", "3")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, inSchema),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(inSchema.DatasourceReference(), "backup_policies.#", "2")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, inDatabase),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(inDatabase.DatasourceReference(), "backup_policies.#", "3")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, withLimit),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(withLimit.DatasourceReference(), "backup_policies.#", "1")),
			},
		},
	})
}

func TestAcc_BackupPolicies_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.BackupPolicies("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_BackupPolicies_EmptyResult(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	backupPoliciesModel := datasourcemodel.BackupPolicies("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, backupPoliciesModel),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(backupPoliciesModel.DatasourceReference(), "backup_policies.#", "0")),
			},
		},
	})
}
