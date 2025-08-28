//go:build !account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_UserProgrammaticAccessToken_basic(t *testing.T) {
	currentUser := testClient().Context.CurrentUser(t)

	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	testClient().Role.GrantRoleToUser(t, role.ID(), user.ID())

	id := testClient().Ids.RandomAccountObjectIdentifier()
	resourceId := helpers.EncodeResourceIdentifier(user.ID().FullyQualifiedName(), id.FullyQualifiedName())
	comment, changedComment := random.Comment(), random.Comment()

	modelBasic := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name())
	modelComplete := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name()).
		WithRoleRestriction(role.ID().Name()).
		WithDaysToExpiry(10).
		WithMinsToBypassNetworkPolicyRequirement(10).
		WithDisabled(r.BooleanTrue).
		WithComment(comment)
	modelCompleteWithDifferentValues := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name()).
		WithRoleRestriction(role.ID().Name()).
		WithDaysToExpiry(10).
		WithMinsToBypassNetworkPolicyRequirement(20).
		WithDisabled(r.BooleanFalse).
		WithComment(changedComment)
	modelWithForceNewOptionals := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name()).
		WithDaysToExpiry(10).
		WithRoleRestriction(role.ID().Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckUserProgrammaticAccessTokenDestroy(t),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasUserString(user.ID().Name()).
						HasRoleRestrictionString("").
						HasNoDaysToExpiry().
						HasNoMinsToBypassNetworkPolicyRequirement().
						HasDisabledString(r.BooleanDefault).
						HasCommentString("").
						HasTokenNotEmpty(),
					resourceshowoutputassert.ProgrammaticAccessTokenShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasUserName(user.ID()).
						HasRoleRestrictionEmpty().
						HasExpiresAtNotEmpty().
						HasStatus(sdk.ProgrammaticAccessTokenStatusActive).
						HasComment("").
						HasCreatedOnNotEmpty().
						HasCreatedBy(currentUser.Name()).
						HasMinsToBypassNetworkPolicyRequirementNotEmpty().
						HasRotatedTo(""),
				),
			},
			// import - without optionals
			{
				Config:       accconfig.FromModels(t, modelBasic),
				ResourceName: modelBasic.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedUserProgrammaticAccessTokenResource(t, resourceId).
						HasNameString(id.Name()).
						HasUserString(user.ID().Name()).
						HasRoleRestrictionString("").
						HasNoDaysToExpiry().
						HasNoMinsToBypassNetworkPolicyRequirement().
						HasDisabledString(r.BooleanFalse).
						HasCommentString("").
						HasNoToken(),
					resourceshowoutputassert.ImportedProgrammaticAccessTokenShowOutput(t, resourceId).
						HasName(id.Name()).
						HasUserName(user.ID()).
						HasRoleRestrictionEmpty().
						HasExpiresAtNotEmpty().
						HasStatus(sdk.ProgrammaticAccessTokenStatusActive).
						HasComment("").
						HasCreatedOnNotEmpty().
						HasCreatedBy(currentUser.Name()).
						HasMinsToBypassNetworkPolicyRequirementNotEmpty().
						HasRotatedTo(""),
				),
			},
			// set optionals
			{
				Config: accconfig.FromModels(t, modelComplete),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectChange(modelCompleteWithDifferentValues.ResourceReference(), "days_to_expiry", tfjson.ActionCreate, nil, sdk.String("10")),
						planchecks.ExpectChange(modelCompleteWithDifferentValues.ResourceReference(), "role_restriction", tfjson.ActionCreate, nil, sdk.String(role.ID().Name())),
					},
				},
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasUserString(user.ID().Name()).
						HasRoleRestrictionString(role.ID().Name()).
						HasDaysToExpiryString("10").
						HasMinsToBypassNetworkPolicyRequirementString("10").
						HasDisabledString("true").
						HasCommentString(comment).
						HasTokenNotEmpty(),
					resourceshowoutputassert.ProgrammaticAccessTokenShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasUserName(user.ID()).
						HasRoleRestriction(role.ID()).
						HasExpiresAtNotEmpty().
						HasStatus(sdk.ProgrammaticAccessTokenStatusDisabled).
						HasComment(comment).
						HasCreatedOnNotEmpty().
						HasCreatedBy(currentUser.Name()).
						HasMinsToBypassNetworkPolicyRequirementNotEmpty().
						HasRotatedTo(""),
				),
			},
			// import - complete
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"days_to_expiry", "expire_rotated_token_after_hours", "mins_to_bypass_network_policy_requirement", "token"},
			},
			// alter
			{
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasUserString(user.ID().Name()).
						HasRoleRestrictionString(role.ID().Name()).
						HasDaysToExpiryString("10").
						HasMinsToBypassNetworkPolicyRequirementString("20").
						HasDisabledString(r.BooleanFalse).
						HasCommentString(changedComment).
						HasTokenNotEmpty(),
					resourceshowoutputassert.ProgrammaticAccessTokenShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasName(id.Name()).
						HasUserName(user.ID()).
						HasRoleRestriction(role.ID()).
						HasExpiresAtNotEmpty().
						HasStatus(sdk.ProgrammaticAccessTokenStatusActive).
						HasComment(changedComment).
						HasCreatedOnNotEmpty().
						HasCreatedBy(currentUser.Name()).
						HasMinsToBypassNetworkPolicyRequirementNotEmpty().
						HasRotatedTo(""),
				),
			},
			// change externally
			{
				PreConfig: func() {
					setRequest := sdk.NewModifyUserProgrammaticAccessTokenRequest(user.ID(), id).
						WithSet(*sdk.NewModifyProgrammaticAccessTokenSetRequest().
							WithDisabled(true).
							WithMinsToBypassNetworkPolicyRequirement(30).
							WithComment(comment),
						)
					testClient().User.ModifyProgrammaticAccessToken(t, setRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasUserString(user.ID().Name()).
						HasRoleRestrictionString(role.ID().Name()).
						HasDaysToExpiryString("10").
						HasMinsToBypassNetworkPolicyRequirementString("20").
						HasDisabledString(r.BooleanFalse).
						HasCommentString(changedComment).
						HasTokenNotEmpty(),
					resourceshowoutputassert.ProgrammaticAccessTokenShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasName(id.Name()).
						HasUserName(user.ID()).
						HasRoleRestriction(role.ID()).
						HasExpiresAtNotEmpty().
						HasStatus(sdk.ProgrammaticAccessTokenStatusActive).
						HasComment(changedComment).
						HasCreatedOnNotEmpty().
						HasCreatedBy(currentUser.Name()).
						HasMinsToBypassNetworkPolicyRequirementNotEmpty().
						HasRotatedTo(""),
				),
			},
			// external changes on days_to_expiry and mins_to_bypass_network_policy_requirement are not detected
			{
				PreConfig: func() {
					testClient().User.RemoveProgrammaticAccessTokenFunc(t, user.ID(), id)()
					request := sdk.NewAddUserProgrammaticAccessTokenRequest(user.ID(), id).
						WithRoleRestriction(role.ID()).
						WithDaysToExpiry(42).
						WithComment(changedComment).
						WithMinsToBypassNetworkPolicyRequirement(22)
					testClient().User.AddProgrammaticAccessTokenWithRequest(t, user.ID(), request)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
			},
			// unset
			{
				Config: accconfig.FromModels(t, modelWithForceNewOptionals),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithForceNewOptionals.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelWithForceNewOptionals.ResourceReference()).
						HasNameString(id.Name()).
						HasUserString(user.ID().Name()).
						HasRoleRestrictionString(role.ID().Name()).
						HasDaysToExpiryString("10").
						HasMinsToBypassNetworkPolicyRequirementString("0").
						HasDisabledString(r.BooleanDefault).
						HasCommentString("").
						HasTokenNotEmpty(),
					resourceshowoutputassert.ProgrammaticAccessTokenShowOutput(t, modelWithForceNewOptionals.ResourceReference()).
						HasName(id.Name()).
						HasUserName(user.ID()).
						HasRoleRestriction(role.ID()).
						HasExpiresAtNotEmpty().
						HasStatus(sdk.ProgrammaticAccessTokenStatusActive).
						HasComment("").
						HasCreatedOnNotEmpty().
						HasCreatedBy(currentUser.Name()).
						HasMinsToBypassNetworkPolicyRequirementNotEmpty().
						HasRotatedTo(""),
				),
			},
			// forcenew - unset all
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasUserString(user.ID().Name()).
						HasRoleRestrictionString("").
						HasNoDaysToExpiry().
						HasNoMinsToBypassNetworkPolicyRequirement().
						HasDisabledString(r.BooleanDefault).
						HasCommentString("").
						HasTokenNotEmpty(),
					resourceshowoutputassert.ProgrammaticAccessTokenShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasUserName(user.ID()).
						HasRoleRestrictionEmpty().
						HasExpiresAtNotEmpty().
						HasStatus(sdk.ProgrammaticAccessTokenStatusActive).
						HasComment("").
						HasCreatedOnNotEmpty().
						HasCreatedBy(currentUser.Name()).
						HasMinsToBypassNetworkPolicyRequirementNotEmpty().
						HasRotatedTo(""),
				),
			},
		},
	})
}

func TestAcc_UserProgrammaticAccessToken_rename(t *testing.T) {
	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	comment := random.Comment()

	id1 := testClient().Ids.RandomAccountObjectIdentifier()
	id2 := testClient().Ids.RandomAccountObjectIdentifier()
	id3 := testClient().Ids.RandomAccountObjectIdentifier()

	modelCompleteId1 := model.UserProgrammaticAccessToken("test", id1.Name(), user.ID().Name())
	modelCompleteId2 := model.UserProgrammaticAccessToken("test", id2.Name(), user.ID().Name())
	modelCompleteId3WithOptionalField := model.UserProgrammaticAccessToken("test", id3.Name(), user.ID().Name()).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckUserProgrammaticAccessTokenDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelCompleteId1),
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelCompleteId1.ResourceReference()).
						HasNameString(id1.Name()).
						HasUserString(user.ID().Name()),
				),
			},
			{
				Config: accconfig.FromModels(t, modelCompleteId2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteId2.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelCompleteId2.ResourceReference()).
						HasNameString(id2.Name()).
						HasUserString(user.ID().Name()),
				),
			},
			{
				Config: accconfig.FromModels(t, modelCompleteId3WithOptionalField),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteId3WithOptionalField.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(modelCompleteId3WithOptionalField.ResourceReference(), "name", tfjson.ActionUpdate, sdk.String(id2.Name()), sdk.String(id3.Name())),
						planchecks.ExpectChange(modelCompleteId3WithOptionalField.ResourceReference(), "comment", tfjson.ActionUpdate, nil, sdk.String(comment)),
					},
				},
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelCompleteId3WithOptionalField.ResourceReference()).
						HasNameString(id3.Name()).
						HasUserString(user.ID().Name()).
						HasCommentString(comment),
				),
			},
		},
	})
}

func TestAcc_UserProgrammaticAccessToken_complete(t *testing.T) {
	currentUser := testClient().Context.CurrentUser(t)

	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	testClient().Role.GrantRoleToUser(t, role.ID(), user.ID())

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	modelComplete := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name()).
		WithRoleRestriction(role.ID().Name()).
		WithDaysToExpiry(30).
		WithMinsToBypassNetworkPolicyRequirement(10).
		WithDisabled("true").
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckUserProgrammaticAccessTokenDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasUserString(user.ID().Name()).
						HasRoleRestrictionString(role.ID().Name()).
						HasDaysToExpiryString("30").
						HasMinsToBypassNetworkPolicyRequirementString("10").
						HasDisabledString("true").
						HasCommentString(comment).
						HasTokenNotEmpty(),
					resourceshowoutputassert.ProgrammaticAccessTokenShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasUserName(user.ID()).
						HasRoleRestriction(role.ID()).
						HasExpiresAtNotEmpty().
						HasStatus(sdk.ProgrammaticAccessTokenStatusDisabled).
						HasComment(comment).
						HasCreatedOnNotEmpty().
						HasCreatedBy(currentUser.Name()).
						HasMinsToBypassNetworkPolicyRequirementNotEmpty().
						HasRotatedTo(""),
				),
			},
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"days_to_expiry", "expire_rotated_token_after_hours", "mins_to_bypass_network_policy_requirement", "token"},
			},
		},
	})
}

func TestAcc_UserProgrammaticAccessToken_rotating(t *testing.T) {
	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	modelBasic := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name())
	modelWithKeeper := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name()).
		WithKeeper("key1=value1")
	modelWithExpireRotatedTokenAfterHours := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name()).
		WithKeeper("key1=value1").
		WithExpireRotatedTokenAfterHours(0)
	modelWithKeeperDifferentValueAndExpireRotatedTokenAfterHours := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name()).
		WithKeeper("key3=value3").
		WithExpireRotatedTokenAfterHours(0)
	modelWithKeeperDifferentValueAndDifferentExpireRotatedTokenAfterHours := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name()).
		WithKeeper("key4=value4").
		WithExpireRotatedTokenAfterHours(3)

	var token string
	tokenAssertion := func(f resource.CheckResourceAttrWithFunc) assert.TestCheckFuncProvider {
		return assert.Check(resource.TestCheckResourceAttrWith(modelBasic.ResourceReference(), "token", f))
	}
	assertTokenNotEmpty := tokenAssertion(func(value string) error {
		if value == "" {
			return fmt.Errorf("token is empty")
		}
		token = value
		return nil
	})

	assertTokenNotRotated := tokenAssertion(func(value string) error {
		if value != token {
			return fmt.Errorf("token was rotated, but should not be")
		}
		token = value
		return nil
	})

	assertTokenRotated := tokenAssertion(func(value string) error {
		if value == token {
			return fmt.Errorf("token was not rotated, but should be")
		}
		token = value
		return nil
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckUserProgrammaticAccessTokenDestroy(t),
		Steps: []resource.TestStep{
			// create the token
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasNoRotatedTokenName().
						HasUserString(user.ID().Name()),
					assertTokenNotEmpty,
				),
			},
			// do not rotate the token with added keeper
			{
				Config: accconfig.FromModels(t, modelWithKeeper),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithKeeper.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.PrintPlanDetails(modelWithKeeper.ResourceReference(), "keeper"),
						planchecks.ExpectChange(modelWithKeeper.ResourceReference(), "keeper", tfjson.ActionUpdate, nil, sdk.String("key1=value1")),
						planchecks.ExpectComputed(modelWithKeeper.ResourceReference(), "token", false),
						planchecks.ExpectComputed(modelWithKeeper.ResourceReference(), "rotated_token_name", false),
					},
				},
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelWithKeeper.ResourceReference()).
						HasNameString(id.Name()).
						HasNoRotatedTokenName().
						HasUserString(user.ID().Name()),
					assertTokenNotRotated,
				),
			},
			// do not rotate when only expire_rotated_token_after_hours is changed
			{
				Config: accconfig.FromModels(t, modelWithExpireRotatedTokenAfterHours),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithExpireRotatedTokenAfterHours.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectComputed(modelWithExpireRotatedTokenAfterHours.ResourceReference(), "token", false),
						planchecks.ExpectComputed(modelWithExpireRotatedTokenAfterHours.ResourceReference(), "rotated_token_name", false),
					},
				},
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelWithExpireRotatedTokenAfterHours.ResourceReference()).
						HasNameString(id.Name()).
						HasNoRotatedTokenName().
						HasUserString(user.ID().Name()),
					assertTokenNotRotated,
				),
			},
			// rotate the token with a different keeper and check that the token is updated
			{
				Config: accconfig.FromModels(t, modelWithKeeperDifferentValueAndExpireRotatedTokenAfterHours),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithKeeperDifferentValueAndExpireRotatedTokenAfterHours.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectComputed(modelWithKeeperDifferentValueAndExpireRotatedTokenAfterHours.ResourceReference(), "token", true),
						planchecks.ExpectComputed(modelWithKeeperDifferentValueAndExpireRotatedTokenAfterHours.ResourceReference(), "rotated_token_name", true),
					},
				},
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelWithKeeperDifferentValueAndExpireRotatedTokenAfterHours.ResourceReference()).
						HasNameString(id.Name()).
						HasRotatedTokenNameNotEmpty().
						HasUserString(user.ID().Name()),
					assertTokenRotated,
					// assert that the rotated token is expired
					assert.Check(resource.TestCheckResourceAttrWith(modelBasic.ResourceReference(), "rotated_token_name", func(value string) error {
						if value == "" {
							return fmt.Errorf("rotated_token_name is empty")
						}
						rotatedTokenId := sdk.NewAccountObjectIdentifier(value)
						token := testClient().User.ShowProgrammaticAccessToken(t, user.ID(), rotatedTokenId)
						if token.RotatedTo == nil {
							return fmt.Errorf("the rotated token is not found")
						}
						if token.Status != sdk.ProgrammaticAccessTokenStatusExpired {
							return fmt.Errorf("the rotated token is not expired")
						}
						return nil
					})),
				),
			},
			// rotate the token with a different keeper and different expire_rotated_token_after_hours and check that the token is updated
			{
				Config: accconfig.FromModels(t, modelWithKeeperDifferentValueAndDifferentExpireRotatedTokenAfterHours),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithKeeperDifferentValueAndDifferentExpireRotatedTokenAfterHours.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectComputed(modelWithKeeperDifferentValueAndDifferentExpireRotatedTokenAfterHours.ResourceReference(), "token", true),
						planchecks.ExpectComputed(modelWithKeeperDifferentValueAndDifferentExpireRotatedTokenAfterHours.ResourceReference(), "rotated_token_name", true),
					},
				},
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelWithKeeperDifferentValueAndDifferentExpireRotatedTokenAfterHours.ResourceReference()).
						HasNameString(id.Name()).
						HasRotatedTokenNameNotEmpty().
						HasUserString(user.ID().Name()),
					assertTokenRotated,
					// assert that the rotated token is expired
					assert.Check(resource.TestCheckResourceAttrWith(modelBasic.ResourceReference(), "rotated_token_name", func(value string) error {
						if value == "" {
							return fmt.Errorf("rotated_token_name is empty")
						}
						rotatedTokenId := sdk.NewAccountObjectIdentifier(value)
						token := testClient().User.ShowProgrammaticAccessToken(t, user.ID(), rotatedTokenId)
						if token.RotatedTo == nil {
							return fmt.Errorf("the rotated token is not found")
						}
						if token.Status != sdk.ProgrammaticAccessTokenStatusActive {
							return fmt.Errorf("the rotated token is expired")
						}
						return nil
					})),
				),
			},
			// do not rotate the token if the keeper is the same
			{
				Config: accconfig.FromModels(t, modelWithKeeperDifferentValueAndDifferentExpireRotatedTokenAfterHours),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithKeeperDifferentValueAndDifferentExpireRotatedTokenAfterHours.ResourceReference(), plancheck.ResourceActionNoop),
						planchecks.ExpectComputed(modelWithKeeperDifferentValueAndDifferentExpireRotatedTokenAfterHours.ResourceReference(), "token", false),
						planchecks.ExpectComputed(modelWithKeeperDifferentValueAndDifferentExpireRotatedTokenAfterHours.ResourceReference(), "rotated_token_name", false),
					},
				},
			},
			// do not rotate the token when the keeper is removed
			{
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.PrintPlanDetails(modelBasic.ResourceReference(), "keeper"),
						planchecks.ExpectChange(modelBasic.ResourceReference(), "keeper", tfjson.ActionUpdate, sdk.String("key4=value4"), nil),
						planchecks.ExpectComputed(modelBasic.ResourceReference(), "token", false),
						planchecks.ExpectComputed(modelBasic.ResourceReference(), "rotated_token_name", false),
					},
				},
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasRotatedTokenNameNotEmpty().
						HasUserString(user.ID().Name()),
					assertTokenNotRotated,
				),
			},
		},
	})
}

func TestAcc_UserProgrammaticAccessToken_RotatingWithExternalProvider(t *testing.T) {
	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()

	timeRotating := `
	resource "time_rotating" "rotation_schedule" {
		rotation_minutes = 30
	}
	`
	timeRotatingWithTriggers := `
	resource "time_rotating" "rotation_schedule" {
		rotation_minutes = 30
		triggers = {
			foo = "bar"
		}
	}
	`

	modelBasic := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name())
	modelWithKeeper := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name()).
		WithKeeper("id ${time_rotating.rotation_schedule.rotation_rfc3339}")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				VersionConstraint: "0.13.1",
				Source:            "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			// create the token
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasNoRotatedTokenName().
						HasUserString(user.ID().Name()),
				),
			},
			// do not rotate the token with added a computed keeper
			{
				Config: timeRotating + accconfig.FromModels(t, modelWithKeeper),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithKeeper.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectComputed(modelWithKeeper.ResourceReference(), "token", false),
						planchecks.ExpectComputed(modelWithKeeper.ResourceReference(), "keeper", true),
						planchecks.ExpectComputed(modelWithKeeper.ResourceReference(), "rotated_token_name", false),
					},
				},
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasNoRotatedTokenName().
						HasUserString(user.ID().Name()),
				),
			},
			// rotate on a changed computed field
			{
				Config: timeRotatingWithTriggers + accconfig.FromModels(t, modelWithKeeper),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithKeeper.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectComputed(modelWithKeeper.ResourceReference(), "token", true),
						planchecks.ExpectComputed(modelWithKeeper.ResourceReference(), "keeper", true),
						planchecks.ExpectComputed(modelWithKeeper.ResourceReference(), "rotated_token_name", true),
					},
				},
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelWithKeeper.ResourceReference()).
						HasNameString(id.Name()).
						HasRotatedTokenNameNotEmpty().
						HasUserString(user.ID().Name()),
				),
			},
			// do not rotate on a removed computed field
			{
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.PrintPlanDetails(modelBasic.ResourceReference(), "keeper"),
						planchecks.ExpectComputed(modelBasic.ResourceReference(), "token", false),
						planchecks.ExpectComputed(modelBasic.ResourceReference(), "rotated_token_name", false),
					},
				},
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasRotatedTokenNameNotEmpty().
						HasUserString(user.ID().Name()),
				),
			},
		},
	})
}

func TestAcc_UserProgrammaticAccessToken_Validations(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()
	id := testClient().Ids.RandomAccountObjectIdentifier()

	modelInvalidDaysToExpiry := model.UserProgrammaticAccessToken("test", id.Name(), userId.Name()).
		WithDaysToExpiry(-1)
	modelInvalidMinsToBypassNetworkPolicyRequirement := model.UserProgrammaticAccessToken("test", id.Name(), userId.Name()).
		WithMinsToBypassNetworkPolicyRequirement(-1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckUserProgrammaticAccessTokenDestroy(t),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, modelInvalidDaysToExpiry),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected days_to_expiry to be at least \(1\), got -1`),
			},
			{
				Config:      config.FromModels(t, modelInvalidMinsToBypassNetworkPolicyRequirement),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected mins_to_bypass_network_policy_requirement to be at least \(1\), got -1`),
			},
		},
	})
}
