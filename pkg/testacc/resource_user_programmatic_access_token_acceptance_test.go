//go:build !account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_UserProgrammaticAccessToken_basic(t *testing.T) {
	currentUser := testClient().Context.CurrentUser(t)

	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	// resourceId := helpers.EncodeResourceIdentifier(user.ID().FullyQualifiedName(), id.FullyQualifiedName())
	comment, changedComment := random.Comment(), random.Comment()

	modelBasic := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name())
	modelWithRoleRestriction := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name()).
		WithRoleRestriction(snowflakeroles.Public.Name())
	// modelComplete := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name()).
	// 	WithRoleRestriction(snowflakeroles.Public.Name()).
	// 	WithDaysToExpiry(30).
	// 	WithMinsToBypassNetworkPolicyRequirement(10).
	// 	WithDisabled("true").
	// 	WithComment(comment)
	modelCompleteWithDifferentValues := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name()).
		WithRoleRestriction(snowflakeroles.Public.Name()).
		WithDaysToExpiry(40).
		WithMinsToBypassNetworkPolicyRequirement(20).
		WithDisabled("false").
		WithComment(changedComment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckUserProgrammaticAccessTokenDestroy(t),
		Steps: []resource.TestStep{
			// // create with empty optionals
			// {
			// 	Config: accconfig.FromModels(t, modelBasic),
			// 	Check: assertThat(t,
			// 		resourceassert.UserProgrammaticAccessTokenResource(t, modelBasic.ResourceReference()).
			// 			HasNameString(id.Name()).
			// 			HasUserString(user.ID().Name()).
			// 			HasRoleRestrictionString("").
			// 			HasNoDaysToExpiry().
			// 			// TODO: use tolerance
			// 			// HasMinsToBypassNetworkPolicyRequirementString("10").
			// 			HasDisabledString(r.BooleanDefault).
			// 			HasCommentString("").
			// 			HasTokenNotEmpty(),
			// 		resourceshowoutputassert.ProgrammaticAccessTokenShowOutput(t, modelBasic.ResourceReference()).
			// 			HasName(id.Name()).
			// 			HasUserName(user.ID()).
			// 			HasRoleRestrictionEmpty().
			// 			HasExpiresAtNotEmpty().
			// 			HasStatus(sdk.ProgrammaticAccessTokenStatusActive).
			// 			HasComment("").
			// 			HasCreatedOnNotEmpty().
			// 			HasCreatedBy(currentUser.Name()).
			// 			// TODO: use tolerance
			// 			// HasMinsToBypassNetworkPolicyRequirementNotEmpty().
			// 			HasRotatedTo(""),
			// 	),
			// },
			// // import - without optionals
			// {
			// 	Config:       accconfig.FromModels(t, modelBasic),
			// 	ResourceName: modelBasic.ResourceReference(),
			// 	ImportState:  true,
			// 	ImportStateCheck: assertThatImport(t,
			// 		resourceassert.ImportedUserProgrammaticAccessTokenResource(t, resourceId).
			// 			HasNameString(id.Name()).
			// 			HasUserString(user.ID().Name()).
			// 			HasRoleRestrictionString("").
			// 			HasNoDaysToExpiry().
			// 			// TODO: use tolerance
			// 			// HasMinsToBypassNetworkPolicyRequirementString("10").
			// 			HasDisabledString(r.BooleanFalse).
			// 			HasCommentString("").
			// 			HasNoToken(),
			// 		resourceshowoutputassert.ImportedProgrammaticAccessTokenShowOutput(t, resourceId).
			// 			HasName(id.Name()).
			// 			HasUserName(user.ID()).
			// 			HasRoleRestrictionEmpty().
			// 			HasExpiresAtNotEmpty().
			// 			HasStatus(sdk.ProgrammaticAccessTokenStatusActive).
			// 			HasComment("").
			// 			HasCreatedOnNotEmpty().
			// 			HasCreatedBy(currentUser.Name()).
			// 			// TODO: use tolerance
			// 			// HasMinsToBypassNetworkPolicyRequirementNotEmpty().
			// 			HasRotatedTo(""),
			// 	),
			// },
			// // set optionals
			// {
			// 	Config: accconfig.FromModels(t, modelComplete),
			// 	Check: assertThat(t,
			// 		resourceassert.UserProgrammaticAccessTokenResource(t, modelBasic.ResourceReference()).
			// 			HasNameString(id.Name()).
			// 			HasUserString(user.ID().Name()).
			// 			HasRoleRestrictionString(snowflakeroles.Public.Name()).
			// 			HasDaysToExpiryString("30").
			// 			// TODO: use tolerance
			// 			// HasMinsToBypassNetworkPolicyRequirementString("10").
			// 			HasDisabledString("true").
			// 			HasCommentString(comment).
			// 			HasTokenNotEmpty(),
			// 		resourceshowoutputassert.ProgrammaticAccessTokenShowOutput(t, modelBasic.ResourceReference()).
			// 			HasName(id.Name()).
			// 			HasUserName(user.ID()).
			// 			HasRoleRestriction(snowflakeroles.Public).
			// 			HasExpiresAtNotEmpty().
			// 			HasStatus(sdk.ProgrammaticAccessTokenStatusDisabled).
			// 			HasComment(comment).
			// 			HasCreatedOnNotEmpty().
			// 			HasCreatedBy(currentUser.Name()).
			// 			// TODO: use tolerance
			// 			// HasMinsToBypassNetworkPolicyRequirementNotEmpty().
			// 			HasRotatedTo(""),
			// 	),
			// },
			// // import - complete
			// {
			// 	Config:                  accconfig.FromModels(t, modelComplete),
			// 	ResourceName:            modelComplete.ResourceReference(),
			// 	ImportState:             true,
			// 	ImportStateVerify:       true,
			// 	ImportStateVerifyIgnore: []string{"days_to_expiry", "mins_to_bypass_network_policy_requirement", "token"},
			// },
			// alter
			{
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				// ConfigPlanChecks: resource.ConfigPlanChecks{
				// 	PreApply: []plancheck.PlanCheck{
				// 		plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionUpdate),
				// 	},
				// },
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasUserString(user.ID().Name()).
						HasRoleRestrictionString(snowflakeroles.Public.Name()).
						HasDaysToExpiryString("40").
						// TODO: use tolerance
						// HasMinsToBypassNetworkPolicyRequirementString("10").
						HasDisabledString(r.BooleanFalse).
						HasCommentString(changedComment).
						HasTokenNotEmpty(),
					resourceshowoutputassert.ProgrammaticAccessTokenShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasUserName(user.ID()).
						HasRoleRestriction(snowflakeroles.Public).
						HasExpiresAtNotEmpty().
						HasStatus(sdk.ProgrammaticAccessTokenStatusActive).
						HasComment(changedComment).
						HasCreatedOnNotEmpty().
						HasCreatedBy(currentUser.Name()).
						// TODO: use tolerance
						// HasMinsToBypassNetworkPolicyRequirementNotEmpty().
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
							WithComment("DUPA" + comment),
						)
					testClient().User.ModifyProgrammaticAccessToken(t, setRequest)
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasUserString(user.ID().Name()).
						HasRoleRestrictionString(snowflakeroles.Public.Name()).
						HasDaysToExpiryString("40").
						// TODO: use tolerance
						// HasMinsToBypassNetworkPolicyRequirementString("10").
						HasDisabledString(r.BooleanFalse).
						HasCommentString(changedComment).
						HasTokenNotEmpty(),
					resourceshowoutputassert.ProgrammaticAccessTokenShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasUserName(user.ID()).
						HasRoleRestriction(snowflakeroles.Public).
						HasExpiresAtNotEmpty().
						HasStatus(sdk.ProgrammaticAccessTokenStatusActive).
						HasComment(changedComment).
						HasCreatedOnNotEmpty().
						HasCreatedBy(currentUser.Name()).
						// TODO: use tolerance
						// HasMinsToBypassNetworkPolicyRequirementNotEmpty().
						HasRotatedTo(""),
				),
			},
			// mins_to_bypass_network_policy_requirement does not cause plans
			// {
			// 	ConfigPlanChecks: resource.ConfigPlanChecks{
			// 		PreApply: []plancheck.PlanCheck{
			// 			plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionNoop),
			// 		},
			// 	},
			// 	Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
			// 	Check: assertThat(t,
			// 		resourceassert.ComputePoolResource(t, modelCompleteWithDifferentValues.ResourceReference()).
			// 			HasNameString(id.Name()).
			// 			HasFullyQualifiedNameString(id.FullyQualifiedName()),
			// 		resourceshowoutputassert.ComputePoolShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
			// 			HasCreatedOnNotEmpty().
			// 			HasName(id.Name()).
			// 			HasComment(changedComment),
			// 	),
			// },
			// unset
			{
				Config: accconfig.FromModels(t, modelWithRoleRestriction),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithRoleRestriction.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelWithRoleRestriction.ResourceReference()).
						HasNameString(id.Name()).
						HasUserString(user.ID().Name()).
						HasRoleRestrictionString(snowflakeroles.Public.Name()).
						HasDaysToExpiryString("0").
						// TODO: use tolerance
						// HasMinsToBypassNetworkPolicyRequirementString("10").
						HasDisabledString(r.BooleanDefault).
						HasCommentString("").
						HasTokenNotEmpty(),
					resourceshowoutputassert.ProgrammaticAccessTokenShowOutput(t, modelWithRoleRestriction.ResourceReference()).
						HasName(id.Name()).
						HasUserName(user.ID()).
						HasRoleRestriction(snowflakeroles.Public).
						HasExpiresAtNotEmpty().
						HasStatus(sdk.ProgrammaticAccessTokenStatusActive).
						HasComment("").
						HasCreatedOnNotEmpty().
						HasCreatedBy(currentUser.Name()).
						// TODO: use tolerance
						// HasMinsToBypassNetworkPolicyRequirementNotEmpty().
						HasRotatedTo(""),
				),
			},
			// forcenew - unset role restriction
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
						// TODO: use tolerance
						// HasMinsToBypassNetworkPolicyRequirementString("10").
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
						// TODO: use tolerance
						// HasMinsToBypassNetworkPolicyRequirementNotEmpty().
						HasRotatedTo(""),
				),
			},
		},
	})
}

func TestAcc_UserProgrammaticAccessToken_rename(t *testing.T) {
	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	newId := testClient().Ids.RandomAccountObjectIdentifier()

	modelComplete := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name())
	modelCompleteNewId := model.UserProgrammaticAccessToken("test", newId.Name(), user.ID().Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
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
						HasUserString(user.ID().Name()),
				),
			},
			{
				Config: accconfig.FromModels(t, modelCompleteNewId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteNewId.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.UserProgrammaticAccessTokenResource(t, modelCompleteNewId.ResourceReference()).
						HasNameString(newId.Name()).
						HasUserString(user.ID().Name()),
				),
			},
		},
	})
}

func TestAcc_UserProgrammaticAccessToken_complete(t *testing.T) {
	currentUser := testClient().Context.CurrentUser(t)

	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	modelComplete := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name()).
		WithRoleRestriction(snowflakeroles.Public.Name()).
		WithDaysToExpiry(30).
		WithMinsToBypassNetworkPolicyRequirement(10).
		WithDisabled("true").
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
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
						HasRoleRestrictionString(snowflakeroles.Public.Name()).
						HasDaysToExpiryString("30").
						// TODO: use tolerance
						// HasMinsToBypassNetworkPolicyRequirementString("10").
						HasDisabledString("true").
						HasCommentString(comment).
						HasTokenNotEmpty(),
					resourceshowoutputassert.ProgrammaticAccessTokenShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasUserName(user.ID()).
						HasRoleRestriction(snowflakeroles.Public).
						HasExpiresAtNotEmpty().
						HasStatus(sdk.ProgrammaticAccessTokenStatusDisabled).
						HasComment(comment).
						HasCreatedOnNotEmpty().
						HasCreatedBy(currentUser.Name()).
						// TODO: use tolerance
						// HasMinsToBypassNetworkPolicyRequirementNotEmpty().
						HasRotatedTo(""),
				),
			},
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"days_to_expiry", "mins_to_bypass_network_policy_requirement", "token"},
			},
		},
	})
}

// TODO(next PR): add tests for rotating the token

func TestAcc_UserProgrammaticAccessToken_Validations(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()
	id := testClient().Ids.RandomAccountObjectIdentifier()

	modelInvalidDaysToExpiry := model.UserProgrammaticAccessToken("test", id.Name(), userId.Name()).
		WithDaysToExpiry(-1)
	modelInvalidMinsToBypassNetworkPolicyRequirement := model.UserProgrammaticAccessToken("test", id.Name(), userId.Name()).
		WithMinsToBypassNetworkPolicyRequirement(-1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
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
