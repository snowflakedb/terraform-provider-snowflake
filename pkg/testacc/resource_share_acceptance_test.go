//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO [next PRs]: Add resourceshowoutputassert, the describe output assertions and the parameters assertions to the use case tests
// once the share resource has show_output, describe_output and parameters (none of them exists in the current resource shape).
// TODO [next PRs]: The share does not import or detect accounts that are not in the configuration (see the import step below).
// TODO [next PRs]: Assert the out-of-band change of accounts with a plan check once planchecks support list fields (drift is only
// asserted on the comment in the plan, and on the accounts after the apply).
func TestAcc_Share_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	account := secondaryTestClient().Account.GetAccountIdentifier(t)
	comment := random.Comment()
	externalComment := random.Comment()

	basic := model.Share("t", id.Name())
	complete := model.Share("t", id.Name()).
		WithComment(comment).
		WithAccounts(account.Name())
	// TODO [next PRs]: The comment cannot be unset by the resource: it sends SET COMMENT = '' and the SDK rejects an empty SET. The unset
	// step removes only the accounts until the update uses UNSET COMMENT.
	withoutAccounts := model.Share("t", id.Name()).WithComment(comment)
	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ShareResource(t, ref).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasCommentEmpty().
			HasAccountsEmpty(),
		objectassert.Share(t, id).
			HasName(id.Name()).
			HasKind(sdk.ShareKindOutbound).
			HasOwnerAccountNotEmpty().
			HasComment("").
			HasNoTo().
			HasToTotal(0),
	}
	withoutAccountsAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ShareResource(t, ref).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasCommentString(comment).
			HasAccountsEmpty(),
		objectassert.Share(t, id).
			HasName(id.Name()).
			HasKind(sdk.ShareKindOutbound).
			HasOwnerAccountNotEmpty().
			HasComment(comment).
			HasNoTo().
			HasToTotal(0),
	}
	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ShareResource(t, ref).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasCommentString(comment).
			HasAccounts(account.Name()),
		objectassert.Share(t, id).
			HasName(id.Name()).
			HasKind(sdk.ShareKindOutbound).
			HasOwnerAccountNotEmpty().
			HasComment(comment).
			HasTo(account).
			HasToTotal(0),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import - without optionals
			{
				Config:            accconfig.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedShareResource(t, id.Name()).
						HasNameString(id.Name()).
						HasCommentEmpty().
						HasAccountsEmpty(),
				),
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			// Import - with optionals
			{
				Config:            accconfig.FromModels(t, complete),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				// TODO [next PRs]: accounts are not imported: reorderStringList drops the accounts that are not configured, and an
				// import has none configured. It also hides accounts that were added in Snowflake outside of the configuration.
				ImportStateVerifyIgnore: []string{"accounts"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedShareResource(t, id.Name()).
						HasNameString(id.Name()).
						HasCommentString(comment),
				),
			},
			// Update - unset optionals (only the accounts, see the TODO above)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectNoChangeOnField(ref, "comment"),
					},
				},
				Config: accconfig.FromModels(t, withoutAccounts),
				Check:  assertThat(t, withoutAccountsAssertions...),
			},
			// Update - external changes (set the optionals again, then change them in Snowflake)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			{
				PreConfig: func() {
					testClient().Share.Alter(t, sdk.NewAlterShareRequest(id).WithSet(sdk.ShareSetRequest{Comment: sdk.String(externalComment)}))
					testClient().Share.Alter(t, sdk.NewAlterShareRequest(id).WithRemove(sdk.ShareRemoveRequest{Accounts: []sdk.AccountIdentifier{account}}))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externalComment)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), sdk.String(comment)),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
		},
	})
}

// TestAcc_Share_Regression_ThreeAccounts adds two recipients to a share and removes them again.
// TODO [next PRs]: Add resourceshowoutputassert and the describe output assertions once the share resource has show_output and describe_output.
func TestAcc_Share_Regression_ThreeAccounts(t *testing.T) {
	// The prod environment has no setup with three accounts on the same deployment, which is needed to add two recipients.
	if testenvs.GetSnowflakeEnvironmentWithProdDefault() == testenvs.SnowflakeProdEnvironment {
		t.Skip("Skipping for the prod Snowflake environment: three accounts on the same deployment are required")
	}

	account2 := secondaryTestClient().Account.GetAccountIdentifier(t)
	account3 := thirdTestClient().Account.GetAccountIdentifier(t)

	comment := random.Comment()
	id := testClient().Ids.RandomAccountObjectIdentifier()

	noAccounts := model.Share("t", id.Name()).WithComment(comment)
	twoAccounts := model.Share("t", id.Name()).WithComment(comment).WithAccounts(account2.Name(), account3.Name())
	oneAccount := model.Share("t", id.Name()).WithComment(comment).WithAccounts(account2.Name())
	ref := noAccounts.ResourceReference()

	noAccountsAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ShareResource(t, ref).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasCommentString(comment).
			HasAccountsEmpty(),
		objectassert.Share(t, id).
			HasName(id.Name()).
			HasComment(comment).
			HasNoTo().
			HasToTotal(0),
	}
	twoAccountsAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ShareResource(t, ref).
			HasNameString(id.Name()).
			HasCommentString(comment).
			HasAccounts(account2.Name(), account3.Name()),
		objectassert.Share(t, id).
			HasName(id.Name()).
			HasComment(comment).
			HasTo(account2, account3).
			HasToTotal(0),
	}
	// TODO [next PRs]: The update only adds accounts and the read drops the accounts that are not configured (reorderStringList), so
	// the account that is removed from the configuration stays in Snowflake and is not shown in the state. The two assertion bundles
	// below do not assert the recipients in Snowflake; they should assert exactly the configured accounts once this is fixed.
	oneAccountAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ShareResource(t, ref).
			HasNameString(id.Name()).
			HasCommentString(comment).
			HasAccounts(account2.Name()),
		objectassert.Share(t, id).
			HasName(id.Name()).
			HasComment(comment),
	}
	noAccountsAfterRemovalAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ShareResource(t, ref).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasCommentString(comment).
			HasAccountsEmpty(),
		objectassert.Share(t, id).
			HasName(id.Name()).
			HasComment(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			// Create - without accounts
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, noAccounts),
				Check:  assertThat(t, noAccountsAssertions...),
			},
			// Update - add two accounts
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, twoAccounts),
				Check:  assertThat(t, twoAccountsAssertions...),
			},
			// Update - remove one account
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, oneAccount),
				Check:  assertThat(t, oneAccountAssertions...),
			},
			// Update - remove all accounts
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, noAccounts),
				Check:  assertThat(t, noAccountsAfterRemovalAssertions...),
			},
			// Import
			{
				Config:            accconfig.FromModels(t, noAccounts),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedShareResource(t, id.Name()).
						HasNameString(id.Name()).
						HasCommentString(comment).
						HasAccountsEmpty(),
				),
			},
		},
	})
}

func TestAcc_Share_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	incorrectAccount := model.Share("t", id.Name()).WithComment(comment).WithAccounts("incorrect")
	correctAndIncorrectAccounts := model.Share("t", id.Name()).WithComment(comment).WithAccounts("correct.one", "incorrect")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, incorrectAccount),
				ExpectError: regexp.MustCompile("Unable to parse the account identifier"),
			},
			{
				Config:      accconfig.FromModels(t, correctAndIncorrectAccounts),
				ExpectError: regexp.MustCompile("Unable to parse the account identifier"),
			},
		},
	})
}

// TestAcc_Share_Regression_Issue4398 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/4398 is fixed.
// TODO [next PRs]: Add resourceshowoutputassert and the describe output assertions once the share resource has show_output and describe_output.
func TestAcc_Share_Regression_Issue4398(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	account := secondaryTestClient().Account.GetAccountIdentifier(t)
	id := testClient().Ids.RandomAccountObjectIdentifier()

	providerModel := providermodel.SnowflakeProvider().WithPreviewFeaturesEnabled(string(previewfeatures.ShareResource))
	noAccounts := model.Share("t", id.Name())
	oneAccount := model.Share("t", id.Name()).WithAccounts(account.Name())
	ref := noAccounts.ResourceReference()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			// Step 1: Share with no accounts
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
				Config:            accconfig.FromModels(t, providerModel, noAccounts),
				Check: assertThat(t,
					resourceassert.ShareResource(t, ref).
						HasNameString(id.Name()).
						HasAccountsEmpty(),
					objectassert.Share(t, id).
						HasName(id.Name()).
						HasNoTo(),
				),
			},
			// Step 2: Grant USAGE on database to the share externally, then add accounts.
			// v2.13.0 fails because setShareAccounts creates a temp database and tries to
			// grant USAGE on it, conflicting with the already-granted USAGE on the real database.
			{
				PreConfig: func() {
					testClient().Grant.GrantPrivilegeOnDatabaseToShare(t, database.ID(), id, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage})
				},
				ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
				Config:            accconfig.FromModels(t, providerModel, oneAccount),
				ExpectError:       regexp.MustCompile("does not belong to the database that is being shared"),
			},
			// Step 3: Succeeds after fix (skip temp database when share already has a database granted).
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, oneAccount),
				Check: assertThat(t,
					resourceassert.ShareResource(t, ref).
						HasNameString(id.Name()).
						HasAccounts(account.Name()),
					objectassert.Share(t, id).
						HasName(id.Name()).
						HasDatabaseName(database.ID()).
						HasTo(account),
				),
			},
		},
	})
}

// TestAcc_Share_Regression_WithGrantedDatabase proves that adding and updating accounts on a share
// that already has a database granted (via snowflake_grant_privileges_to_share) works correctly.
// This is the regression scenario reported after PR #4174 — previously, the provider would
// erroneously try to create a temp database even when a real DB was already in the share,
// causing a conflict on the USAGE grant.
// See: https://github.com/snowflakedb/terraform-provider-snowflake/pull/4174
// TODO [next PRs]: Add resourceshowoutputassert and the describe output assertions once the share resource has show_output and describe_output.
func TestAcc_Share_Regression_WithGrantedDatabase(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	account := secondaryTestClient().Account.GetAccountIdentifier(t)
	id := testClient().Ids.RandomAccountObjectIdentifier()

	shareModel := model.ShareWithDefaultMeta(id.Name())
	shareModelWithAccount := model.ShareWithDefaultMeta(id.Name()).WithAccounts(account.Name())
	grantModel := model.GrantPrivilegesToShareWithDefaultMeta([]string{"USAGE"}, id.Name()).
		WithOnDatabase(database.ID().Name()).
		WithDependsOn(shareModel.ResourceReference())
	ref := shareModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			// Step 1: Share with no accounts + grant DB to share (no accounts yet).
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, shareModel, grantModel),
				Check: assertThat(t,
					resourceassert.ShareResource(t, ref).
						HasNameString(id.Name()).
						HasAccountsEmpty(),
					objectassert.Share(t, id).
						HasName(id.Name()).
						HasDatabaseName(database.ID()).
						HasNoTo(),
				),
			},
			// Step 2: Add an account while DB is already on the share.
			// This exercises the UpdateShare path that directly adds accounts
			// (skipping the temp DB workaround) when share.DatabaseName is set.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, shareModelWithAccount, grantModel),
				Check: assertThat(t,
					resourceassert.ShareResource(t, ref).
						HasNameString(id.Name()).
						HasAccounts(account.Name()),
					objectassert.Share(t, id).
						HasName(id.Name()).
						HasDatabaseName(database.ID()).
						HasTo(account),
				),
			},
			// Step 3: Remove account (DB still granted).
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, shareModel, grantModel),
				Check: assertThat(t,
					resourceassert.ShareResource(t, ref).
						HasNameString(id.Name()).
						HasAccountsEmpty(),
					objectassert.Share(t, id).
						HasName(id.Name()).
						HasDatabaseName(database.ID()).
						HasNoTo(),
				),
			},
		},
	})
}
