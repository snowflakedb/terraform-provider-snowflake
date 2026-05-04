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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_PasswordPolicies_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	completeModel := model.PasswordPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMinLength(10).
		WithMaxLength(30).
		WithMinUpperCaseChars(2).
		WithMinLowerCaseChars(3).
		WithMinNumericChars(1).
		WithMinSpecialChars(1).
		WithMinAgeDays(1).
		WithMaxAgeDays(30).
		WithMaxRetries(3).
		WithLockoutTimeMins(15).
		WithHistory(5).
		WithComment(comment)

	passwordPoliciesModel := datasourcemodel.PasswordPolicies("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(completeModel.ResourceReference())

	passwordPoliciesModelWithoutDescribe := datasourcemodel.PasswordPolicies("test").
		WithWithDescribe(false).
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(completeModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, completeModel, passwordPoliciesModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(passwordPoliciesModel.DatasourceReference(), "password_policies.#", "1")),
					resourceshowoutputassert.PasswordPoliciesDatasourceShowOutput(t, "snowflake_password_policies.test").
						HasName(id.Name()).
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindPasswordPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					resourceshowoutputassert.PasswordPoliciesDatasourceDescribeOutput(t, "snowflake_password_policies.test").
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasPasswordMinLength(10).
						HasPasswordMaxLength(30).
						HasPasswordMinUpperCaseChars(2).
						HasPasswordMinLowerCaseChars(3).
						HasPasswordMinNumericChars(1).
						HasPasswordMinSpecialChars(1).
						HasPasswordMinAgeDays(1).
						HasPasswordMaxAgeDays(30).
						HasPasswordMaxRetries(3).
						HasPasswordLockoutTimeMins(15).
						HasPasswordHistory(5),
				),
			},
			{
				Config: accconfig.FromModels(t, completeModel, passwordPoliciesModelWithoutDescribe),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(passwordPoliciesModelWithoutDescribe.DatasourceReference(), "password_policies.#", "1")),
					resourceshowoutputassert.PasswordPoliciesDatasourceShowOutput(t, "snowflake_password_policies.test").
						HasName(id.Name()).
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindPasswordPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(passwordPoliciesModelWithoutDescribe.DatasourceReference(), "password_policies.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_PasswordPolicies_Filtering(t *testing.T) {
	secondSchema, secondSchemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, sdk.NewAccountObjectIdentifier(TestDatabaseName))
	t.Cleanup(secondSchemaCleanup)

	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifierInSchema(secondSchema.ID())

	testUser, testUserCleanup := testClient().User.CreateUser(t)
	t.Cleanup(testUserCleanup)

	model1 := model.PasswordPolicy("test1", id1.DatabaseName(), id1.SchemaName(), id1.Name())
	model2 := model.PasswordPolicy("test2", id2.DatabaseName(), id2.SchemaName(), id2.Name())
	model3 := model.PasswordPolicy("test3", id3.DatabaseName(), id3.SchemaName(), id3.Name())

	passwordPoliciesModelLikeFirst := datasourcemodel.PasswordPolicies("test").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	passwordPoliciesModelLikePrefix := datasourcemodel.PasswordPolicies("test").
		WithLike(prefix+"%").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	passwordPoliciesModelStartsWith := datasourcemodel.PasswordPolicies("test").
		WithStartsWith(prefix).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	passwordPoliciesModelLimit := datasourcemodel.PasswordPolicies("test").
		WithRowsAndFrom(1, prefix).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	passwordPoliciesModelOnUser := datasourcemodel.PasswordPolicies("test").
		WithOnUser(testUser.ID()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	passwordPoliciesModelOnAccount := datasourcemodel.PasswordPolicies("test").
		WithOnAccount().
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	passwordPoliciesModelInSchema := datasourcemodel.PasswordPolicies("test").
		WithInSchema(id1.SchemaId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, passwordPoliciesModelLikeFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(passwordPoliciesModelLikeFirst.DatasourceReference(), "password_policies.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, passwordPoliciesModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(passwordPoliciesModelLikePrefix.DatasourceReference(), "password_policies.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, passwordPoliciesModelStartsWith),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(passwordPoliciesModelStartsWith.DatasourceReference(), "password_policies.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, passwordPoliciesModelLimit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(passwordPoliciesModelLimit.DatasourceReference(), "password_policies.#", "1"),
				),
			},
			{
				PreConfig: func() {
					testClient().User.Alter(t, testUser.ID(), &sdk.AlterUserOptions{
						Set: &sdk.UserSet{PasswordPolicy: sdk.Pointer(id3)},
					})
				},
				Config: accconfig.FromModels(t, model1, model2, model3, passwordPoliciesModelOnUser),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(passwordPoliciesModelOnUser.DatasourceReference(), "password_policies.#", "1"),
				),
				// Unset the password policy from the user.
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						planchecks.Execute(func() {
							testClient().User.Alter(t, testUser.ID(), &sdk.AlterUserOptions{
								Unset: &sdk.UserUnset{PasswordPolicy: sdk.Bool(true)},
							})
						}),
					},
				},
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, passwordPoliciesModelOnAccount),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(passwordPoliciesModelOnAccount.DatasourceReference(), "password_policies.#", "0"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, passwordPoliciesModelInSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(passwordPoliciesModelInSchema.DatasourceReference(), "password_policies.#", "2"),
				),
			},
		},
	})
}

func TestAcc_PasswordPolicies_emptyOn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.PasswordPolicies("test").WithEmptyOn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_PasswordPolicies_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.PasswordPolicies("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}
