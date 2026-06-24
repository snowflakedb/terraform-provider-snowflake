//go:build non_account_level_tests

package testacc

import (
	"fmt"
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

func TestAcc_SessionPolicies_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	role1, role1Cleanup := testClient().Role.CreateRole(t)
	t.Cleanup(role1Cleanup)
	role2, role2Cleanup := testClient().Role.CreateRole(t)
	t.Cleanup(role2Cleanup)
	role3, role3Cleanup := testClient().Role.CreateRole(t)
	t.Cleanup(role3Cleanup)

	completeModel := model.SessionPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithSessionIdleTimeoutMins(30).
		WithSessionUiIdleTimeoutMins(60).
		WithAllowedSecondaryRoles(role1.Name, role2.Name).
		WithBlockedSecondaryRoles(role3.Name).
		WithComment(comment)

	sessionPoliciesModel := datasourcemodel.SessionPolicies("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(completeModel.ResourceReference())

	sessionPoliciesModelWithoutDescribe := datasourcemodel.SessionPolicies("test").
		WithWithDescribe(false).
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(completeModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: policiesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, completeModel, sessionPoliciesModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(sessionPoliciesModel.DatasourceReference(), "session_policies.#", "1")),
					resourceshowoutputassert.SessionPoliciesDatasourceShowOutput(t, "snowflake_session_policies.test").
						HasName(id.Name()).
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindSessionPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					resourceshowoutputassert.SessionPoliciesDatasourceDescribeOutput(t, "snowflake_session_policies.test").
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(comment).
						HasSessionIdleTimeoutMins(30).
						HasSessionUiIdleTimeoutMins(60).
						HasAllowedSecondaryRoles(role1.Name, role2.Name).
						HasBlockedSecondaryRoles(role3.Name),
				),
			},
			{
				Config: accconfig.FromModels(t, completeModel, sessionPoliciesModelWithoutDescribe),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(sessionPoliciesModelWithoutDescribe.DatasourceReference(), "session_policies.#", "1")),
					resourceshowoutputassert.SessionPoliciesDatasourceShowOutput(t, "snowflake_session_policies.test").
						HasName(id.Name()).
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindSessionPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(sessionPoliciesModelWithoutDescribe.DatasourceReference(), "session_policies.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_SessionPolicies_Filtering(t *testing.T) {
	secondSchema, secondSchemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, sdk.NewAccountObjectIdentifier(TestDatabaseName))
	t.Cleanup(secondSchemaCleanup)

	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifierInSchema(secondSchema.ID())

	testUser, testUserCleanup := testClient().User.CreateUser(t)
	t.Cleanup(testUserCleanup)

	model1 := model.SessionPolicy("test1", id1.DatabaseName(), id1.SchemaName(), id1.Name())
	model2 := model.SessionPolicy("test2", id2.DatabaseName(), id2.SchemaName(), id2.Name())
	model3 := model.SessionPolicy("test3", id3.DatabaseName(), id3.SchemaName(), id3.Name())

	userSessionPolicyAttachment := model.UserSessionPolicyAttachment("attach", "", testUser.ID().Name()).
		WithSessionPolicyNameValue(accconfig.UnquotedWrapperVariable(fmt.Sprintf("%s.fully_qualified_name", model3.ResourceReference()))).
		WithDependsOn(model3.ResourceReference())

	sessionPoliciesModelLikeFirst := datasourcemodel.SessionPolicies("test").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	sessionPoliciesModelLikePrefix := datasourcemodel.SessionPolicies("test").
		WithLike(prefix+"%").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	sessionPoliciesModelStartsWith := datasourcemodel.SessionPolicies("test").
		WithStartsWith(prefix).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	sessionPoliciesModelLimit := datasourcemodel.SessionPolicies("test").
		WithRowsAndFrom(1, prefix).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	sessionPoliciesModelOnUser := datasourcemodel.SessionPolicies("test").
		WithOnUser(testUser.ID()).
		WithDependsOn(
			model1.ResourceReference(),
			model2.ResourceReference(),
			model3.ResourceReference(),
			userSessionPolicyAttachment.ResourceReference(),
		)

	sessionPoliciesModelOnAccount := datasourcemodel.SessionPolicies("test").
		WithOnAccount().
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	sessionPoliciesModelInSchema := datasourcemodel.SessionPolicies("test").
		WithInSchema(id1.SchemaId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: policiesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, sessionPoliciesModelLikeFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(sessionPoliciesModelLikeFirst.DatasourceReference(), "session_policies.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, sessionPoliciesModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(sessionPoliciesModelLikePrefix.DatasourceReference(), "session_policies.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, sessionPoliciesModelStartsWith),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(sessionPoliciesModelStartsWith.DatasourceReference(), "session_policies.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, sessionPoliciesModelLimit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(sessionPoliciesModelLimit.DatasourceReference(), "session_policies.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, userSessionPolicyAttachment, sessionPoliciesModelOnUser),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(sessionPoliciesModelOnUser.DatasourceReference(), "session_policies.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, sessionPoliciesModelOnAccount),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(sessionPoliciesModelOnAccount.DatasourceReference(), "session_policies.#", "0"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, sessionPoliciesModelInSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(sessionPoliciesModelInSchema.DatasourceReference(), "session_policies.#", "2"),
				),
			},
		},
	})
}

func TestAcc_SessionPolicies_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: policiesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.SessionPolicies("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_SessionPolicies_emptyOn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: policiesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.SessionPolicies("test").WithEmptyOn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}
