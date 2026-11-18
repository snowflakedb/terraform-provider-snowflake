//go:build account_level_tests

// Account object types other than SHARES and DATABASES (e.g. ROLES, INTEGRATIONS, WAREHOUSES) can be a part of only one
// failover or replication group per account. Snowflake rejects the second group with 500070 "Account object types, other than
// SHARES and DATABASES, can only be within one failover group or replication group" (documented as "Objects other than databases,
// external volumes, and shares must be in the same failover group"). The restriction is on the type, so unique names do not help.
// The tests below use such types (ROLES, INTEGRATIONS, WAREHOUSES, ACCOUNT PARAMETERS).
// A database can also be added to only one failover or replication group, and these tests put the shared test database into groups.
// Failover and replication groups share these restrictions, so the tests of both are account-level tests (they cannot run in parallel
// with each other, other PRs, or the integration/acceptance suites).

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// SHOW prints the ACCOUNT PARAMETERS object type as PARAMETERS (GH #2517), so the provider has to map it back.
func TestAcc_FailoverGroup_issue2517_AccountParametersObjectType(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	resourceId := helpers.EncodeResourceIdentifier(id)
	sourceAccount := testClient().Account.GetAccountIdentifier(t)
	targetAccount := secondaryTestClient().Account.GetAccountIdentifier(t)
	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)

	withoutAccountParameters := model.FailoverGroup("t", id.Name(), []string{targetAccount.FullyQualifiedName()}, []string{"DATABASES"}).
		WithAllowedDatabases([]string{database.ID().Name()})
	withAccountParameters := model.FailoverGroup("t", id.Name(), []string{targetAccount.FullyQualifiedName()}, []string{"ACCOUNT PARAMETERS", "DATABASES"}).
		WithAllowedDatabases([]string{database.ID().Name()})
	ref := withoutAccountParameters.ResourceReference()

	assertions := func(resourceObjectTypes []string, objectTypes ...sdk.FailoverGroupObjectType) []assert.TestCheckFuncProvider {
		return []assert.TestCheckFuncProvider{
			resourceassert.FailoverGroupResource(t, ref).
				HasNameString(id.Name()).
				HasObjectTypes(resourceObjectTypes...).
				HasAllowedDatabases(database.ID().Name()),
			objectassert.FailoverGroup(t, id).
				HasName(id.Name()).
				HasObjectTypesInAnyOrder(objectTypes...).
				HasAllowedAccountsInAnyOrder(sourceAccount, targetAccount),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			// 1. Create without the object type
			{
				Config: config.FromModels(t, withoutAccountParameters),
				Check:  assertThat(t, assertions([]string{"DATABASES"}, sdk.FailoverGroupObjectTypeDatabases)...),
			},
			// 2. Add the object type
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectSetElementsChange(ref, "object_types", tfjson.ActionUpdate, []string{"ACCOUNT PARAMETERS"}, nil),
					},
				},
				Config: config.FromModels(t, withAccountParameters),
				Check:  assertThat(t, assertions([]string{"ACCOUNT PARAMETERS", "DATABASES"}, sdk.FailoverGroupObjectTypeAccountParameters, sdk.FailoverGroupObjectTypeDatabases)...),
			},
			// 3. Import
			{
				Config:            config.FromModels(t, withAccountParameters),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				// TODO [SNOW-1648997]: show_output is not recomputed when a set changes, so these two lists are stale after the update.
				ImportStateVerifyIgnore: []string{"ignore_edition_check", "show_output.0.object_types", "show_output.0.allowed_integration_types"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedFailoverGroupResource(t, resourceId).
						HasNameString(id.Name()).
						HasObjectTypes("ACCOUNT PARAMETERS", "DATABASES"),
				),
			},
			// 4. Remove the object type
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectSetElementsChange(ref, "object_types", tfjson.ActionUpdate, nil, []string{"ACCOUNT PARAMETERS"}),
					},
				},
				Config: config.FromModels(t, withoutAccountParameters),
				Check:  assertThat(t, assertions([]string{"DATABASES"}, sdk.FailoverGroupObjectTypeDatabases)...),
			},
		},
	})
}

// OBJECT_TYPES with ALLOWED_INTEGRATION_TYPES and REPLICATION_SCHEDULE cannot share one SET statement (GH #2544), so one apply that
// changes all of them is sent as separate statements.
func TestAcc_FailoverGroup_issue2544_SeparateAlterStatements(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	sourceAccount := testClient().Account.GetAccountIdentifier(t)
	targetAccount := secondaryTestClient().Account.GetAccountIdentifier(t)
	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)

	both := model.FailoverGroup("t", id.Name(), []string{targetAccount.FullyQualifiedName()}, []string{"DATABASES", "INTEGRATIONS"}).
		WithAllowedDatabases([]string{database.ID().Name()}).
		WithAllowedIntegrationTypes([]string{"API INTEGRATIONS", "NOTIFICATION INTEGRATIONS"}).
		WithReplicationSchedule("10 MINUTE")
	narrowed := model.FailoverGroup("t", id.Name(), []string{targetAccount.FullyQualifiedName()}, []string{"INTEGRATIONS"}).
		WithAllowedIntegrationTypes([]string{"API INTEGRATIONS"}).
		WithReplicationSchedule("20 MINUTE")
	ref := both.ResourceReference()

	assertBoth := []assert.TestCheckFuncProvider{
		resourceassert.FailoverGroupResource(t, ref).
			HasNameString(id.Name()).
			HasObjectTypes("DATABASES", "INTEGRATIONS").
			HasAllowedDatabases(database.ID().Name()).
			HasAllowedIntegrationTypes("API INTEGRATIONS", "NOTIFICATION INTEGRATIONS").
			HasReplicationSchedule("10 MINUTE"),
		objectassert.FailoverGroup(t, id).
			HasName(id.Name()).
			HasObjectTypesInAnyOrder(sdk.FailoverGroupObjectTypeDatabases, sdk.FailoverGroupObjectTypeIntegrations).
			HasAllowedIntegrationTypesInAnyOrder(sdk.FailoverGroupIntegrationTypeApiIntegrations, sdk.FailoverGroupIntegrationTypeNotificationIntegrations).
			HasAllowedAccountsInAnyOrder(sourceAccount, targetAccount).
			HasReplicationSchedule("10 MINUTE"),
	}
	assertNarrowed := []assert.TestCheckFuncProvider{
		resourceassert.FailoverGroupResource(t, ref).
			HasNameString(id.Name()).
			HasObjectTypes("INTEGRATIONS").
			HasAllowedDatabasesEmpty().
			HasAllowedIntegrationTypes("API INTEGRATIONS").
			HasReplicationSchedule("20 MINUTE"),
		objectassert.FailoverGroup(t, id).
			HasName(id.Name()).
			HasObjectTypesInAnyOrder(sdk.FailoverGroupObjectTypeIntegrations).
			HasAllowedIntegrationTypesInAnyOrder(sdk.FailoverGroupIntegrationTypeApiIntegrations).
			HasAllowedAccountsInAnyOrder(sourceAccount, targetAccount).
			HasReplicationSchedule("20 MINUTE"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			// 1. Create
			{
				Config: config.FromModels(t, both),
				Check:  assertThat(t, assertBoth...),
			},
			// 2. Narrow the object types and the integration types and change the schedule in one apply
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectSetElementsChange(ref, "object_types", tfjson.ActionUpdate, nil, []string{"DATABASES"}),
						planchecks.ExpectSetElementsChange(ref, "allowed_integration_types", tfjson.ActionUpdate, nil, []string{"NOTIFICATION INTEGRATIONS"}),
						planchecks.ExpectSetElementsChange(ref, "allowed_databases", tfjson.ActionUpdate, nil, []string{database.ID().Name()}),
						planchecks.ExpectChange(ref, "replication_schedule", tfjson.ActionUpdate, sdk.String("10 MINUTE"), sdk.String("20 MINUTE")),
					},
				},
				Config: config.FromModels(t, narrowed),
				Check:  assertThat(t, assertNarrowed...),
			},
			// 3. Widen them again in one apply
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectSetElementsChange(ref, "object_types", tfjson.ActionUpdate, []string{"DATABASES"}, nil),
						planchecks.ExpectSetElementsChange(ref, "allowed_integration_types", tfjson.ActionUpdate, []string{"NOTIFICATION INTEGRATIONS"}, nil),
						planchecks.ExpectSetElementsChange(ref, "allowed_databases", tfjson.ActionUpdate, []string{database.ID().Name()}, nil),
						planchecks.ExpectChange(ref, "replication_schedule", tfjson.ActionUpdate, sdk.String("20 MINUTE"), sdk.String("10 MINUTE")),
					},
				},
				Config: config.FromModels(t, both),
				Check:  assertThat(t, assertBoth...),
			},
		},
	})
}

func TestAcc_FailoverGroup_UpdateAllowedAccounts(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	// We need to use the secondary account to test the failover group. Otherwise, we will get an error:
	// Error: error removing allowed accounts for failover group TEST err = 003909 (55000): Disabling replication of the replication group to the account in which the primary currently resides is not allowed.
	sourceAccount := testClient().Context.CurrentAccountIdentifier(t)
	secondaryAccount := secondaryTestClient().Account.GetAccountIdentifier(t)
	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.FailoverGroupResource))

	newModel := func(accounts ...sdk.AccountIdentifier) *model.FailoverGroupModel {
		return model.FailoverGroup("t", id.Name(), collections.Map(accounts, sdk.AccountIdentifier.Name), []string{"DATABASES"}).
			WithAllowedDatabases([]string{database.ID().Name()})
	}
	// The old provider version lists the source account in allowed_accounts, the current one leaves it out.
	sourceAndTarget := newModel(sourceAccount, secondaryAccount)
	onlySource := newModel(sourceAccount)
	onlyTarget := newModel(secondaryAccount)
	ref := sourceAndTarget.ResourceReference()

	assertions := func(inState []sdk.AccountIdentifier, inSnowflake []sdk.AccountIdentifier) []assert.TestCheckFuncProvider {
		return []assert.TestCheckFuncProvider{
			resourceassert.FailoverGroupResource(t, ref).
				HasNameString(id.Name()).
				HasAllowedAccounts(collections.Map(inState, sdk.AccountIdentifier.Name)...),
			objectassert.FailoverGroup(t, id).
				HasName(id.Name()).
				HasAllowedAccountsInAnyOrder(inSnowflake...),
		}
	}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			// Create with allowed_accounts using old provider version
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
				Config:            config.FromModels(t, providerModel, sourceAndTarget),
				Check:             assertThat(t, assertions([]sdk.AccountIdentifier{sourceAccount, secondaryAccount}, []sdk.AccountIdentifier{sourceAccount, secondaryAccount})...),
			},
			// Remove an allowed account with old provider version - proves the bug
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
				Config:            config.FromModels(t, providerModel, onlySource),
				ExpectError:       regexp.MustCompile("error removing allowed accounts"),
			},
			// Remove it externally
			{
				PreConfig: func() {
					testClient().FailoverGroup.RemoveAllowedAccounts(t, id, secondaryAccount)
				},
				ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
				Config:            config.FromModels(t, providerModel, onlySource),
				Check:             assertThat(t, assertions([]sdk.AccountIdentifier{sourceAccount}, []sdk.AccountIdentifier{sourceAccount})...),
			},
			// Add it back again - proves the bug
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
				Config:            config.FromModels(t, providerModel, sourceAndTarget),
				ExpectError:       regexp.MustCompile("error adding allowed accounts"),
			},
			// Upgrade to the current provider version. The source account is no longer a part of allowed_accounts, so the config lists only the target account.
			// The secondary account was removed externally above, so this step adds it back with the current provider - proves the fix.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectSetElementsChange(ref, "allowed_accounts", tfjson.ActionUpdate, []string{secondaryAccount.Name()}, nil),
					},
				},
				Config: config.FromModels(t, onlyTarget),
				Check:  assertThat(t, assertions([]sdk.AccountIdentifier{secondaryAccount}, []sdk.AccountIdentifier{sourceAccount, secondaryAccount})...),
			},
			// Removing an allowed account needs a second target account, because allowed_accounts requires at least one.
			// TestAcc_FailoverGroup_AllowedAccounts_ThreeAccounts covers it where a third account is available.
		},
	})
}

// TestAcc_FailoverGroup_AllowedAccounts_ThreeAccounts removes a target account from allowed_accounts and adds it back.
func TestAcc_FailoverGroup_AllowedAccounts_ThreeAccounts(t *testing.T) {
	// The prod environment has no setup with three accounts on the same deployment, which is needed to have two target accounts.
	if testenvs.GetSnowflakeEnvironmentWithProdDefault() == testenvs.SnowflakeProdEnvironment {
		t.Skip("Skipping for the prod Snowflake environment: three accounts on the same deployment are required")
	}

	id := testClient().Ids.RandomAccountObjectIdentifier()
	sourceAccount := testClient().Account.GetAccountIdentifier(t)
	secondaryAccount := secondaryTestClient().Account.GetAccountIdentifier(t)
	thirdAccount := thirdTestClient().Account.GetAccountIdentifier(t)
	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)

	newModel := func(accounts ...sdk.AccountIdentifier) *model.FailoverGroupModel {
		return model.FailoverGroup("t", id.Name(), collections.Map(accounts, sdk.AccountIdentifier.FullyQualifiedName), []string{"DATABASES"}).
			WithAllowedDatabases([]string{database.ID().Name()})
	}
	twoTargets := newModel(secondaryAccount, thirdAccount)
	oneTarget := newModel(secondaryAccount)
	ref := twoTargets.ResourceReference()

	// SHOW lists the source account too; the resource leaves it out of allowed_accounts.
	assertions := func(targets ...sdk.AccountIdentifier) []assert.TestCheckFuncProvider {
		return []assert.TestCheckFuncProvider{
			resourceassert.FailoverGroupResource(t, ref).
				HasNameString(id.Name()).
				HasObjectTypes("DATABASES").
				HasAllowedDatabases(database.ID().Name()).
				HasAllowedAccounts(collections.Map(targets, sdk.AccountIdentifier.FullyQualifiedName)...),
			objectassert.FailoverGroup(t, id).
				HasName(id.Name()).
				HasIsPrimary(true).
				HasAllowedAccountsInAnyOrder(append([]sdk.AccountIdentifier{sourceAccount}, targets...)...),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			// 1. Create with two target accounts
			{
				Config: config.FromModels(t, twoTargets),
				Check:  assertThat(t, assertions(secondaryAccount, thirdAccount)...),
			},
			// 2. Remove an allowed account
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectSetElementsChange(ref, "allowed_accounts", tfjson.ActionUpdate, nil, []string{thirdAccount.FullyQualifiedName()}),
					},
				},
				Config: config.FromModels(t, oneTarget),
				Check:  assertThat(t, assertions(secondaryAccount)...),
			},
			// 3. Add it back
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectSetElementsChange(ref, "allowed_accounts", tfjson.ActionUpdate, []string{thirdAccount.FullyQualifiedName()}, nil),
					},
				},
				Config: config.FromModels(t, twoTargets),
				Check:  assertThat(t, assertions(secondaryAccount, thirdAccount)...),
			},
		},
	})
}

// The configuration is validated before anything is created in Snowflake.
func TestAcc_FailoverGroup_Validations(t *testing.T) {
	accounts := []string{"MYORG.MYACCOUNT"}
	valid := func() *model.FailoverGroupModel {
		return model.FailoverGroup("t", "name", accounts, []string{"ROLES"})
	}

	invalidObjectType := model.FailoverGroup("t", "name", accounts, []string{"DATABASES", "INVALID OBJECT TYPE;"})
	invalidIntegrationType := valid().WithAllowedIntegrationTypes([]string{"INVALID INTEGRATION TYPE"})
	invalidAllowedAccount := model.FailoverGroup("t", "name", []string{"not_an_account_identifier"}, []string{"ROLES"})
	rpoAssuranceWithSchedule := valid().WithRpoAssurance(true).WithReplicationSchedule("10 MINUTE")
	rpoAssuranceWithOptimizedRefresh := valid().WithRpoAssurance(true).WithOptimizedRefresh(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, invalidObjectType),
				ExpectError: regexp.MustCompile("invalid failover group object type: INVALID OBJECT TYPE;"),
			},
			{
				Config:      config.FromModels(t, invalidIntegrationType),
				ExpectError: regexp.MustCompile("invalid failover group integration type: INVALID INTEGRATION TYPE"),
			},
			{
				Config:      config.FromModels(t, invalidAllowedAccount),
				ExpectError: regexp.MustCompile("Unable to parse the account identifier"),
			},
			{
				Config:      config.FromModels(t, rpoAssuranceWithSchedule),
				ExpectError: regexp.MustCompile("conflicts with"),
			},
			{
				Config:      config.FromModels(t, rpoAssuranceWithOptimizedRefresh),
				ExpectError: regexp.MustCompile("conflicts with"),
			},
		},
	})
}

// The state written by the last release with the hand-written resource (nested replication_schedule, boolean ignore_edition_check,
// the source account in allowed_accounts) is upgraded without a diff once the configuration is adjusted.
func TestAcc_FailoverGroup_UpgradeFromV2_23_0(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	sourceAccount := testClient().Context.CurrentAccountIdentifier(t)
	targetAccount := secondaryTestClient().Account.GetAccountIdentifier(t)
	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.FailoverGroupResource))

	// The nested replication_schedule block cannot be written with the current model builders.
	configInV2_23_0 := fmt.Sprintf(`
resource "snowflake_failover_group" "t" {
	name                 = "%s"
	object_types         = ["DATABASES"]
	allowed_accounts     = ["%s", "%s"]
	allowed_databases    = ["%s"]
	ignore_edition_check = true
	optimized_refresh    = true
	replication_schedule {
		interval = 10
	}
}
`, id.Name(), sourceAccount.Name(), targetAccount.Name(), database.ID().Name())
	current := model.FailoverGroup("t", id.Name(), []string{targetAccount.FullyQualifiedName()}, []string{"DATABASES"}).
		WithAllowedDatabases([]string{database.ID().Name()}).
		WithIgnoreEditionCheck(true).
		WithOptimizedRefresh(true).
		WithReplicationSchedule("10 MINUTE")
	comment := random.Comment()
	updated := model.FailoverGroup("t", id.Name(), []string{targetAccount.FullyQualifiedName()}, []string{"DATABASES"}).
		WithAllowedDatabases([]string{database.ID().Name()}).
		WithIgnoreEditionCheck(true).
		WithOptimizedRefresh(true).
		WithReplicationSchedule("20 MINUTE").
		WithComment(comment)
	ref := current.ResourceReference()

	// The same bundle describes the group after the upgrade and after the update: only the schedule and the comment differ.
	assertions := func(schedule string, comment *string) []assert.TestCheckFuncProvider {
		resourceAssert := resourceassert.FailoverGroupResource(t, ref).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasObjectTypes("DATABASES").
			HasAllowedAccounts(targetAccount.FullyQualifiedName()).
			HasAllowedDatabases(database.ID().Name()).
			HasReplicationSchedule(schedule).
			HasOptimizedRefresh(true).
			HasIgnoreEditionCheck(true).
			HasIsPrimary(true)
		showOutputAssert := resourceshowoutputassert.FailoverGroupShowOutput(t, ref).
			HasName(id.Name()).
			HasType("FAILOVER").
			HasIsPrimary(true).
			HasReplicationSchedule(schedule).
			HasIsOptimizedRefreshEnabled(true).
			HasRpoAssurance(false).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasObjectTypes("DATABASES").
			HasNoAllowedIntegrationTypes()
		objectAssert := objectassert.FailoverGroup(t, id).
			HasName(id.Name()).
			HasType("FAILOVER").
			HasIsPrimary(true).
			HasObjectTypesInAnyOrder(sdk.FailoverGroupObjectTypeDatabases).
			HasAllowedAccountsInAnyOrder(sourceAccount, targetAccount).
			HasReplicationSchedule(schedule).
			HasIsOptimizedRefreshEnabled(true).
			HasRpoAssurance(false)
		if comment == nil {
			resourceAssert.HasCommentEmpty()
			objectAssert.HasNoComment()
		} else {
			resourceAssert.HasCommentString(*comment)
			showOutputAssert.HasComment(*comment)
			objectAssert.HasComment(*comment)
		}
		return []assert.TestCheckFuncProvider{resourceAssert, showOutputAssert, objectAssert}
	}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			// Create with the last release
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.23.0"),
				Config:            config.FromModels(t, providerModel) + configInV2_23_0,
				Check: assertThat(t,
					resourceassert.FailoverGroupResource(t, ref).
						HasNameString(id.Name()).
						HasAllowedAccounts(sourceAccount.Name(), targetAccount.Name()),
					objectassert.FailoverGroup(t, id).
						HasName(id.Name()).
						HasReplicationSchedule("10 MINUTE").
						HasIsOptimizedRefreshEnabled(true),
				),
			},
			// Upgrade: the state is converted, the configuration is rewritten (string schedule, no source account), and nothing is left to apply
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: config.FromModels(t, current),
				Check:  assertThat(t, assertions("10 MINUTE", nil)...),
			},
			// Update: the upgraded state is managed like any other
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "replication_schedule", tfjson.ActionUpdate, sdk.String("10 MINUTE"), sdk.String("20 MINUTE")),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, nil, sdk.String(comment)),
					},
				},
				Config: config.FromModels(t, updated),
				Check:  assertThat(t, assertions("20 MINUTE", &comment)...),
			},
		},
	})
}

// The replication schedule is an interval or a cron expression; OPTIMIZED_REFRESH needs a schedule of 6 hours or less and
// RPO_ASSURANCE needs no schedule, so the transitions between the two (and the implicit Optimized Refresh of RPO assurance) are the point of this test.
func TestAcc_FailoverGroup_CompleteUseCase_ReplicationScheduleTransitions(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	sourceAccount := testClient().Account.GetAccountIdentifier(t)
	targetAccount := secondaryTestClient().Account.GetAccountIdentifier(t)
	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)

	const (
		interval10 = "10 MINUTE"
		interval20 = "20 MINUTE"
		cronA      = "USING CRON 0 0 10-20 * TUE,THU UTC"
		cronB      = "USING CRON 0 0 5-20 * TUE,THU UTC"
	)
	newModel := func(schedule string, optimizedRefresh, rpoAssurance bool) *model.FailoverGroupModel {
		m := model.FailoverGroup("t", id.Name(), []string{targetAccount.FullyQualifiedName()}, []string{"DATABASES"}).
			WithAllowedDatabases([]string{database.ID().Name()})
		if schedule != "" {
			m.WithReplicationSchedule(schedule)
		}
		// optimized_refresh and rpo_assurance conflict, so a false value is left out of the configuration.
		if optimizedRefresh {
			m.WithOptimizedRefresh(true)
		}
		if rpoAssurance {
			m.WithRpoAssurance(true)
		}
		return m
	}
	intervalWithOptimizedRefresh := newModel(interval10, true, false)
	longerInterval := newModel(interval20, true, false)
	cronWithoutOptimizedRefresh := newModel(cronA, false, false)
	otherCron := newModel(cronB, false, false)
	rpo := newModel("", false, true)
	intervalAfterRpoWithOptimizedRefresh := newModel(interval10, true, false)
	intervalAfterRpoWithoutOptimizedRefresh := newModel(interval10, false, false)
	ref := intervalWithOptimizedRefresh.ResourceReference()

	// optimizedRefreshInSnowflake differs from the state of optimized_refresh while RPO assurance is on: Snowflake enables Optimized Refresh with it.
	assertions := func(schedule string, optimizedRefresh, rpoAssurance bool) []assert.TestCheckFuncProvider {
		optimizedRefreshInSnowflake := optimizedRefresh || rpoAssurance
		resourceAssert := resourceassert.FailoverGroupResource(t, ref).
			HasNameString(id.Name()).
			HasObjectTypes("DATABASES").
			HasAllowedAccounts(targetAccount.FullyQualifiedName()).
			HasOptimizedRefresh(optimizedRefresh).
			HasRpoAssurance(rpoAssurance)
		showOutputAssert := resourceshowoutputassert.FailoverGroupShowOutput(t, ref).
			HasName(id.Name()).
			HasIsPrimary(true).
			HasObjectTypes("DATABASES").
			HasIsOptimizedRefreshEnabled(optimizedRefreshInSnowflake).
			HasRpoAssurance(rpoAssurance)
		objectAssert := objectassert.FailoverGroup(t, id).
			HasName(id.Name()).
			HasIsPrimary(true).
			HasAllowedAccountsInAnyOrder(sourceAccount, targetAccount).
			HasIsOptimizedRefreshEnabled(optimizedRefreshInSnowflake).
			HasRpoAssurance(rpoAssurance)
		if schedule == "" {
			resourceAssert.HasReplicationScheduleEmpty()
			objectAssert.HasNoReplicationSchedule()
		} else {
			resourceAssert.HasReplicationSchedule(schedule)
			showOutputAssert.HasReplicationSchedule(schedule)
			objectAssert.HasReplicationSchedule(schedule)
		}
		return []assert.TestCheckFuncProvider{resourceAssert, showOutputAssert, objectAssert}
	}

	update := func(checks ...plancheck.PlanCheck) resource.ConfigPlanChecks {
		return resource.ConfigPlanChecks{
			PreApply: append([]plancheck.PlanCheck{
				plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
				planchecks.ExpectNoChangeOnField(ref, "ignore_edition_check"),
			}, checks...),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			// 1. Create - interval with Optimized Refresh
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate)},
				},
				Config: config.FromModels(t, intervalWithOptimizedRefresh),
				Check:  assertThat(t, assertions(interval10, true, false)...),
			},
			// 2. Interval -> a longer interval
			{
				ConfigPlanChecks: update(
					planchecks.ExpectChange(ref, "replication_schedule", tfjson.ActionUpdate, sdk.String(interval10), sdk.String(interval20)),
					planchecks.ExpectNoChangeOnField(ref, "optimized_refresh"),
				),
				Config: config.FromModels(t, longerInterval),
				Check:  assertThat(t, assertions(interval20, true, false)...),
			},
			// 3. Interval -> cron; Optimized Refresh is switched off first, because it needs a schedule of 6 hours or less
			{
				ConfigPlanChecks: update(
					planchecks.ExpectChange(ref, "replication_schedule", tfjson.ActionUpdate, sdk.String(interval20), sdk.String(cronA)),
					planchecks.ExpectChange(ref, "optimized_refresh", tfjson.ActionUpdate, sdk.String("true"), nil),
				),
				Config: config.FromModels(t, cronWithoutOptimizedRefresh),
				Check:  assertThat(t, assertions(cronA, false, false)...),
			},
			// 4. Cron -> another cron
			{
				ConfigPlanChecks: update(
					planchecks.ExpectChange(ref, "replication_schedule", tfjson.ActionUpdate, sdk.String(cronA), sdk.String(cronB)),
				),
				Config: config.FromModels(t, otherCron),
				Check:  assertThat(t, assertions(cronB, false, false)...),
			},
			// 5. Cron -> no schedule with RPO assurance
			{
				ConfigPlanChecks: update(
					planchecks.ExpectChange(ref, "replication_schedule", tfjson.ActionUpdate, sdk.String(cronB), nil),
					planchecks.ExpectChange(ref, "rpo_assurance", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
				),
				Config: config.FromModels(t, rpo),
				Check:  assertThat(t, assertions("", false, true)...),
			},
			// 6. RPO assurance -> interval with Optimized Refresh (Snowflake keeps Optimized Refresh on after RPO assurance is unset)
			{
				ConfigPlanChecks: update(
					planchecks.ExpectChange(ref, "replication_schedule", tfjson.ActionUpdate, nil, sdk.String(interval10)),
					planchecks.ExpectChange(ref, "rpo_assurance", tfjson.ActionUpdate, sdk.String("true"), nil),
					planchecks.ExpectChange(ref, "optimized_refresh", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
				),
				Config: config.FromModels(t, intervalAfterRpoWithOptimizedRefresh),
				Check:  assertThat(t, assertions(interval10, true, false)...),
			},
			// 7. Interval with Optimized Refresh -> no schedule with RPO assurance
			{
				ConfigPlanChecks: update(
					planchecks.ExpectChange(ref, "replication_schedule", tfjson.ActionUpdate, sdk.String(interval10), nil),
					planchecks.ExpectChange(ref, "rpo_assurance", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
					planchecks.ExpectChange(ref, "optimized_refresh", tfjson.ActionUpdate, sdk.String("true"), nil),
				),
				Config: config.FromModels(t, rpo),
				Check:  assertThat(t, assertions("", false, true)...),
			},
			// 8. RPO assurance -> interval without Optimized Refresh (the Optimized Refresh that RPO assurance enabled is unset explicitly)
			{
				ConfigPlanChecks: update(
					planchecks.ExpectChange(ref, "replication_schedule", tfjson.ActionUpdate, nil, sdk.String(interval10)),
					planchecks.ExpectChange(ref, "rpo_assurance", tfjson.ActionUpdate, sdk.String("true"), nil),
				),
				Config: config.FromModels(t, intervalAfterRpoWithoutOptimizedRefresh),
				Check:  assertThat(t, assertions(interval10, false, false)...),
			},
			// 9. Import
			{
				Config:                  config.FromModels(t, intervalAfterRpoWithoutOptimizedRefresh),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_edition_check"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedFailoverGroupResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasReplicationSchedule(interval10).
						HasOptimizedRefresh(false).
						HasRpoAssurance(false),
				),
			},
		},
	})
}

// CREATE has no COMMENT clause, so the comment is set with ALTER after the group is created.
// SHOW does not return allowed_external_volumes, so it cannot be read back or imported.
func TestAcc_FailoverGroup_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	resourceId := helpers.EncodeResourceIdentifier(id)
	sourceAccount := testClient().Account.GetAccountIdentifier(t)
	targetAccount := secondaryTestClient().Account.GetAccountIdentifier(t)
	comment := random.Comment()
	externalComment := random.Comment()

	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	share, cleanupShare := testClient().Share.Create(t)
	t.Cleanup(cleanupShare)
	externalVolume, cleanupExternalVolume := testClient().ExternalVolume.Create(t)
	t.Cleanup(cleanupExternalVolume)
	errorIntegration, cleanupErrorIntegration := testClient().NotificationIntegration.Create(t)
	t.Cleanup(cleanupErrorIntegration)

	objectTypes := []string{"ROLES", "DATABASES", "SHARES", "EXTERNAL VOLUMES", "INTEGRATIONS", "WAREHOUSES"}
	integrationTypes := []string{"API INTEGRATIONS", "EXTERNAL ACCESS INTEGRATIONS", "NOTIFICATION INTEGRATIONS", "SECURITY INTEGRATIONS", "STORAGE INTEGRATIONS"}

	basic := model.FailoverGroup("t", id.Name(), []string{targetAccount.FullyQualifiedName()}, []string{"ROLES"})
	withOptionals := model.FailoverGroup("t", id.Name(), []string{targetAccount.FullyQualifiedName()}, objectTypes).
		WithAllowedDatabases([]string{database.ID().Name()}).
		WithAllowedShares([]string{share.ID().Name()}).
		WithAllowedExternalVolumes([]string{externalVolume.Name()}).
		WithAllowedIntegrationTypes(integrationTypes).
		WithReplicationSchedule("10 MINUTE").
		WithOptimizedRefresh(true).
		WithErrorIntegration(errorIntegration.ID().Name()).
		WithReadOnly(true).
		WithComment(comment)
	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.FailoverGroupResource(t, ref).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasObjectTypes("ROLES").
			HasAllowedAccounts(targetAccount.FullyQualifiedName()).
			HasAllowedDatabasesEmpty().
			HasAllowedIntegrationTypesEmpty().
			HasReplicationScheduleEmpty().
			HasOptimizedRefreshString("false").
			HasRpoAssuranceString("false").
			HasCommentEmpty().
			HasReadOnly(false).
			HasErrorIntegrationEmpty().
			HasNoIgnoreEditionCheck().
			HasIsPrimaryString("true"),
		resourceshowoutputassert.FailoverGroupShowOutput(t, ref).
			HasName(id.Name()).
			HasType("FAILOVER").
			HasIsPrimary(true).
			HasIsPrimaryReadOnly(false).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasCreatedOnNotEmpty().
			HasSnowflakeRegionNotEmpty().
			HasAllowedAccounts(sourceAccount.Name(), targetAccount.Name()).
			HasObjectTypes("ROLES").
			HasNoAllowedIntegrationTypes().
			HasIsOptimizedRefreshEnabled(false).
			HasRpoAssurance(false),
		objectassert.FailoverGroup(t, id).
			HasName(id.Name()).
			HasType("FAILOVER").
			HasIsPrimary(true).
			HasPrimary(sdk.NewExternalObjectIdentifier(sourceAccount, id)).
			HasObjectTypes(sdk.FailoverGroupObjectTypeRoles).
			HasAllowedAccountsInAnyOrder(sourceAccount, targetAccount).
			HasNoAllowedIntegrationTypes().
			HasNoReplicationSchedule().
			HasNoComment().
			HasIsOptimizedRefreshEnabled(false).
			HasRpoAssurance(false).
			HasOwner(snowflakeroles.Accountadmin.Name()),
	}
	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.FailoverGroupResource(t, ref).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasObjectTypes(objectTypes...).
			HasAllowedAccounts(targetAccount.FullyQualifiedName()).
			HasAllowedDatabases(database.ID().Name()).
			HasAllowedShares(share.ID().Name()).
			HasAllowedExternalVolumes(externalVolume.Name()).
			HasAllowedIntegrationTypes(integrationTypes...).
			HasReplicationSchedule("10 MINUTE").
			HasOptimizedRefresh(true).
			HasErrorIntegration(errorIntegration.ID().Name()).
			HasReadOnly(true).
			HasCommentString(comment).
			HasNoIgnoreEditionCheck().
			HasIsPrimaryString("true"),
		resourceshowoutputassert.FailoverGroupShowOutput(t, ref).
			HasName(id.Name()).
			HasType("FAILOVER").
			HasIsPrimary(true).
			HasIsPrimaryReadOnly(true).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasReplicationSchedule("10 MINUTE").
			HasErrorIntegration(errorIntegration.ID()).
			HasObjectTypes(objectTypes...).
			HasAllowedIntegrationTypes(integrationTypes...).
			HasIsOptimizedRefreshEnabled(true).
			HasRpoAssurance(false),
		objectassert.FailoverGroup(t, id).
			HasName(id.Name()).
			HasIsPrimary(true).
			HasObjectTypesInAnyOrder(sdk.FailoverGroupObjectTypeRoles, sdk.FailoverGroupObjectTypeDatabases, sdk.FailoverGroupObjectTypeShares, sdk.FailoverGroupObjectTypeExternalVolumes, sdk.FailoverGroupObjectTypeIntegrations, sdk.FailoverGroupObjectTypeWarehouses).
			HasAllowedAccountsInAnyOrder(sourceAccount, targetAccount).
			HasAllowedIntegrationTypesInAnyOrder(sdk.AllFailoverGroupIntegrationTypes...).
			HasReplicationSchedule("10 MINUTE").
			HasComment(comment).
			HasErrorIntegration(errorIntegration.ID()).
			HasIsPrimaryReadOnly(true).
			HasIsOptimizedRefreshEnabled(true).
			HasRpoAssurance(false),
	}
	importIgnore := []string{
		"allowed_external_volumes", // SHOW does not return the allowed external volumes.
		"ignore_edition_check",     // The value is used only when the group is created, the import keeps the default.
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			// 1. Create - without optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// 2. Import - without optionals
			{
				Config:                  config.FromModels(t, basic),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importIgnore,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedFailoverGroupResource(t, resourceId).
						HasNameString(id.Name()).
						HasObjectTypes("ROLES").
						HasAllowedAccounts(targetAccount.FullyQualifiedName()).
						HasIsPrimaryString("true"),
					resourceshowoutputassert.ImportedFailoverGroupShowOutput(t, resourceId).
						HasName(id.Name()).
						HasType("FAILOVER").
						HasIsPrimary(true),
				),
			},
			// 3. Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, nil, sdk.String(comment)),
						planchecks.ExpectChange(ref, "replication_schedule", tfjson.ActionUpdate, nil, sdk.String("10 MINUTE")),
						planchecks.ExpectSetElementsChange(ref, "object_types", tfjson.ActionUpdate, []string{"DATABASES", "SHARES", "EXTERNAL VOLUMES", "INTEGRATIONS", "WAREHOUSES"}, nil),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// 4. Import - with optionals
			{
				Config:                  config.FromModels(t, withOptionals),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importIgnore,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedFailoverGroupResource(t, resourceId).
						HasNameString(id.Name()).
						HasObjectTypes(objectTypes...).
						HasAllowedAccounts(targetAccount.FullyQualifiedName()).
						HasAllowedDatabases(database.ID().Name()).
						HasAllowedShares(share.ID().Name()).
						HasAllowedIntegrationTypes(integrationTypes...).
						HasReplicationSchedule("10 MINUTE").
						HasOptimizedRefresh(true).
						HasCommentString(comment),
					resourceshowoutputassert.ImportedFailoverGroupShowOutput(t, resourceId).
						HasName(id.Name()).
						HasComment(comment).
						HasReplicationSchedule("10 MINUTE"),
				),
			},
			// 5. Update - external changes
			{
				PreConfig: func() {
					testClient().FailoverGroup.Alter(t, sdk.NewAlterFailoverGroupRequest(id).
						WithSet(*sdk.NewFailoverGroupSetRequest().WithComment(externalComment)))
					testClient().FailoverGroup.Alter(t, sdk.NewAlterFailoverGroupRequest(id).
						WithSet(*sdk.NewFailoverGroupSetRequest().WithReplicationSchedule("20 MINUTE")))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externalComment)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), sdk.String(comment)),
						planchecks.ExpectDrift(ref, "replication_schedule", sdk.String("10 MINUTE"), sdk.String("20 MINUTE")),
						planchecks.ExpectChange(ref, "replication_schedule", tfjson.ActionUpdate, sdk.String("20 MINUTE"), sdk.String("10 MINUTE")),
						planchecks.ExpectNoChangeOnField(ref, "optimized_refresh"),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// 6. Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(comment), nil),
						planchecks.ExpectChange(ref, "replication_schedule", tfjson.ActionUpdate, sdk.String("10 MINUTE"), nil),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
		},
	})
}

// Snowflake rejects the membership changes of a failover group in the read-only mode, so the plan fails when read_only stays true.
// Switching read_only off in the same apply leaves the read-only mode before the change, and switching it on is done afterwards.
// SHOW returns the read-only mode of a primary failover group in is_primary_read_only, so the checks read it back.
// TODO [next PRs]: add an experiment that lifts and restores the read-only mode around the membership changes.
func TestAcc_FailoverGroup_ReadOnlyMode(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	resourceId := helpers.EncodeResourceIdentifier(id)
	sourceAccount := testClient().Account.GetAccountIdentifier(t)
	targetAccount := secondaryTestClient().Account.GetAccountIdentifier(t)

	databaseA, cleanupDatabaseA := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabaseA)
	databaseB, cleanupDatabaseB := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabaseB)

	newModel := func(readOnly *bool, databases ...string) *model.FailoverGroupModel {
		m := model.FailoverGroup("t", id.Name(), []string{targetAccount.FullyQualifiedName()}, []string{"DATABASES"}).WithAllowedDatabases(databases)
		if readOnly != nil {
			m.WithReadOnly(*readOnly)
		}
		return m
	}
	on, off := true, false
	readOnlyWithA := newModel(&on, databaseA.ID().Name())
	readOnlyWithAAndB := newModel(&on, databaseA.ID().Name(), databaseB.ID().Name())
	notReadOnlyWithAAndB := newModel(&off, databaseA.ID().Name(), databaseB.ID().Name())
	withoutReadOnlyWithA := newModel(nil, databaseA.ID().Name())
	ref := readOnlyWithA.ResourceReference()

	assertions := func(readOnly *bool, databases ...*sdk.AccountObjectIdentifier) []assert.TestCheckFuncProvider {
		databaseNames := collections.Map(databases, func(database *sdk.AccountObjectIdentifier) string { return database.Name() })
		resourceAssert := resourceassert.FailoverGroupResource(t, ref).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasObjectTypes("DATABASES").
			HasAllowedAccounts(targetAccount.FullyQualifiedName()).
			HasAllowedDatabases(databaseNames...).
			HasIsPrimaryString("true")
		switch {
		case readOnly == nil:
			resourceAssert.HasReadOnlyString("false")
		default:
			resourceAssert.HasReadOnly(*readOnly)
		}
		return []assert.TestCheckFuncProvider{
			resourceAssert,
			resourceshowoutputassert.FailoverGroupShowOutput(t, ref).
				HasName(id.Name()).
				HasType("FAILOVER").
				HasIsPrimary(true).
				HasIsPrimaryReadOnly(readOnly != nil && *readOnly).
				HasOwner(snowflakeroles.Accountadmin.Name()),
			objectassert.FailoverGroup(t, id).
				HasName(id.Name()).
				HasType("FAILOVER").
				HasIsPrimary(true).
				HasObjectTypesInAnyOrder(sdk.FailoverGroupObjectTypeDatabases).
				HasAllowedAccountsInAnyOrder(sourceAccount, targetAccount).
				HasIsPrimaryReadOnly(readOnly != nil && *readOnly).
				HasOwner(snowflakeroles.Accountadmin.Name()),
		}
	}
	idA, idB := databaseA.ID(), databaseB.ID()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			// 1. Create in the read-only mode
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, readOnlyWithA),
				Check:  assertThat(t, assertions(&on, &idA)...),
			},
			// 2. Changing the membership while read_only stays true fails at the plan
			{
				Config:      config.FromModels(t, readOnlyWithAAndB),
				ExpectError: regexp.MustCompile(`cannot change allowed_databases while read_only is true`),
			},
			// 3. Leaving the read-only mode and changing the membership in one apply
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "read_only", tfjson.ActionUpdate, sdk.String("true"), sdk.String("false")),
						planchecks.ExpectSetElementsChange(ref, "allowed_databases", tfjson.ActionUpdate, []string{databaseB.ID().Name()}, nil),
					},
				},
				Config: config.FromModels(t, notReadOnlyWithAAndB),
				Check:  assertThat(t, assertions(&off, &idA, &idB)...),
			},
			// 4. Import
			{
				Config:                  config.FromModels(t, notReadOnlyWithAAndB),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_edition_check"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedFailoverGroupResource(t, resourceId).
						HasNameString(id.Name()).
						HasObjectTypes("DATABASES").
						HasAllowedDatabases(databaseA.ID().Name(), databaseB.ID().Name()).
						HasReadOnly(false).
						HasIsPrimaryString("true"),
				),
			},
			// 5. Entering the read-only mode again
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "read_only", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
						planchecks.ExpectNoChangeOnField(ref, "ignore_edition_check"),
					},
				},
				Config: config.FromModels(t, readOnlyWithAAndB),
				Check:  assertThat(t, assertions(&on, &idA, &idB)...),
			},
			// 6. Removing read_only and changing the membership in one apply
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectSetElementsChange(ref, "allowed_databases", tfjson.ActionUpdate, nil, []string{databaseB.ID().Name()}),
					},
				},
				Config: config.FromModels(t, withoutReadOnlyWithA),
				Check:  assertThat(t, assertions(nil, &idA)...),
			},
		},
	})
}

// The error integration is read from SHOW, so changing it or unsetting it outside of Terraform is drift that the next apply restores.
func TestAcc_FailoverGroup_ErrorIntegration_ExternalChanges(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	allowedAccount := secondaryTestClient().Account.GetAccountIdentifier(t)
	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	errorIntegration, cleanupErrorIntegration := testClient().NotificationIntegration.Create(t)
	t.Cleanup(cleanupErrorIntegration)
	otherIntegration, cleanupOtherIntegration := testClient().NotificationIntegration.Create(t)
	t.Cleanup(cleanupOtherIntegration)

	configured := model.FailoverGroup("t", id.Name(), []string{allowedAccount.FullyQualifiedName()}, []string{"DATABASES"}).
		WithAllowedDatabases([]string{database.ID().Name()}).
		WithErrorIntegration(errorIntegration.ID().Name())
	ref := configured.ResourceReference()

	restored := []assert.TestCheckFuncProvider{
		resourceassert.FailoverGroupResource(t, ref).
			HasNameString(id.Name()).
			HasErrorIntegrationString(errorIntegration.ID().Name()),
		objectassert.FailoverGroup(t, id).
			HasName(id.Name()).
			HasErrorIntegration(errorIntegration.ID()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			// 1. Create with the error integration
			{
				Config: config.FromModels(t, configured),
				Check:  assertThat(t, restored...),
			},
			// 2. The integration is changed outside of Terraform
			{
				PreConfig: func() {
					testClient().FailoverGroup.Alter(t, sdk.NewAlterFailoverGroupRequest(id).
						WithSet(*sdk.NewFailoverGroupSetRequest().WithErrorIntegration(otherIntegration.ID())))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "error_integration", sdk.String(errorIntegration.ID().Name()), sdk.String(otherIntegration.ID().Name())),
						planchecks.ExpectChange(ref, "error_integration", tfjson.ActionUpdate, sdk.String(otherIntegration.ID().Name()), sdk.String(errorIntegration.ID().Name())),
					},
				},
				Config: config.FromModels(t, configured),
				Check:  assertThat(t, restored...),
			},
			// 3. The integration is unset outside of Terraform
			{
				PreConfig: func() {
					testClient().FailoverGroup.Alter(t, sdk.NewAlterFailoverGroupRequest(id).
						WithUnset(*sdk.NewFailoverGroupUnsetRequest().WithErrorIntegration(true)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "error_integration", sdk.String(errorIntegration.ID().Name()), sdk.String("")),
						planchecks.ExpectChange(ref, "error_integration", tfjson.ActionUpdate, sdk.String(""), sdk.String(errorIntegration.ID().Name())),
					},
				},
				Config: config.FromModels(t, configured),
				Check:  assertThat(t, restored...),
			},
		},
	})
}

// A failover (a promotion of a replica in another account) makes the managed group a replica. The resource is recreated as a primary then.
func TestAcc_FailoverGroup_RecreateWhenExternallyNotPrimary(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	sourceAccount := testClient().Account.GetAccountIdentifier(t)
	targetAccount := secondaryTestClient().Account.GetAccountIdentifier(t)
	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)

	// A primary in the other account, so this account can replace the managed group with a replica of the same name.
	otherDatabase, cleanupOtherDatabase := secondaryTestClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupOtherDatabase)
	otherPrimary := secondaryTestClient().Ids.RandomAccountObjectIdentifier()
	_, cleanupOtherPrimary := secondaryTestClient().FailoverGroup.CreateWithRequest(t, sdk.NewCreateFailoverGroupRequest(
		otherPrimary,
		[]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeDatabases},
		[]sdk.AccountIdentifier{sourceAccount},
	).WithAllowedDatabases([]sdk.AccountObjectIdentifier{otherDatabase.ID()}))
	t.Cleanup(cleanupOtherPrimary)
	t.Cleanup(func() {
		// The brief replica copies the database of the other primary into this account. Drop that copy before the other primary goes.
		_ = testClient().Database.DropDatabase(t, otherDatabase.ID())
	})

	group := model.FailoverGroup("t", id.Name(), []string{targetAccount.FullyQualifiedName()}, []string{"DATABASES"}).
		WithAllowedDatabases([]string{database.ID().Name()})
	ref := group.ResourceReference()
	otherOf := sdk.NewExternalObjectIdentifier(targetAccount, otherPrimary)
	assertPrimary := []assert.TestCheckFuncProvider{
		resourceassert.FailoverGroupResource(t, ref).
			HasNameString(id.Name()).
			HasObjectTypes("DATABASES").
			HasAllowedDatabases(database.ID().Name()).
			HasIsPrimary(true),
		objectassert.FailoverGroup(t, id).
			HasName(id.Name()).
			HasIsPrimary(true).
			HasPrimary(sdk.NewExternalObjectIdentifier(sourceAccount, id)),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FailoverGroup),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, group),
				Check:  assertThat(t, assertPrimary...),
			},
			{
				// Replace the primary with a replica of the same name. The resource must destroy it and create the primary again.
				PreConfig: func() {
					testClient().FailoverGroup.Drop(t, id)
					testClient().FailoverGroup.WaitForPrimaryToBeVisible(t, otherOf)
					testClient().FailoverGroup.CreateSecondary(t, id, otherOf)
				},
				Config: config.FromModels(t, group),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(t, assertPrimary...),
			},
		},
	})
}
