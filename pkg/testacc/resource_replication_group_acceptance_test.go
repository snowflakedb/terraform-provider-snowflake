//go:build account_level_tests

// Account object types other than SHARES and DATABASES (e.g. ROLES, INTEGRATIONS, WAREHOUSES) can be a part of only one
// failover or replication group per account. Snowflake rejects the second group with 500070 "Account object types, other than
// SHARES and DATABASES, can only be within one failover group or replication group" (documented as "Objects other than databases,
// external volumes, and shares must be in the same failover group"). The restriction is on the type, so unique names do not help.
// The tests below use such types (ROLES, INTEGRATIONS).
// Failover and replication groups share these restrictions, so the tests of both are account-level tests (they cannot run in parallel
// with each other, other PRs, or the integration/acceptance suites).

package testacc

import (
	"testing"

	tfresources "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/invokeactionassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ReplicationGroup_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	allowedAccount := secondaryTestClient().Account.GetAccountIdentifier(t)
	sourceAccount := testClient().Account.GetAccountIdentifier(t)
	comment := random.Comment()
	externalComment := random.Comment()
	const schedule = "10 MINUTE"
	const externalSchedule = "20 MINUTE"
	// Config names the target account. Snowflake also adds the source account; the resource suppresses that member.
	configAccounts := []string{allowedAccount.FullyQualifiedName()}

	database2, cleanupDatabase2 := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase2)
	share, cleanupShare := testClient().Share.Create(t)
	t.Cleanup(cleanupShare)
	externalVolume, cleanupExternalVolume := testClient().ExternalVolume.Create(t)
	t.Cleanup(cleanupExternalVolume)
	errorIntegration, cleanupErrorIntegration := testClient().NotificationIntegration.Create(t)
	t.Cleanup(cleanupErrorIntegration)

	// SHOW returns these in its own order. The resource stores that order.
	objectTypes := []sdk.ReplicationGroupObjectType{
		sdk.ReplicationGroupObjectTypeDatabases,
		sdk.ReplicationGroupObjectTypeRoles,
		sdk.ReplicationGroupObjectTypeExternalVolumes,
		sdk.ReplicationGroupObjectTypeShares,
		sdk.ReplicationGroupObjectTypeIntegrations,
	}
	integrationTypes := []sdk.ReplicationGroupIntegrationType{
		sdk.ReplicationGroupIntegrationTypeApiIntegrations,
		sdk.ReplicationGroupIntegrationTypeExternalAccessIntegrations,
		sdk.ReplicationGroupIntegrationTypeNotificationIntegrations,
		sdk.ReplicationGroupIntegrationTypeSecurityIntegrations,
		sdk.ReplicationGroupIntegrationTypeStorageIntegrations,
	}

	basic := model.ReplicationGroup("t", id.Name(), configAccounts, []string{string(sdk.ReplicationGroupObjectTypeDatabases)}).
		WithAllowedDatabases([]string{database.ID().FullyQualifiedName()})
	// One apply sets both exclusive SET groups: object types and allowed lists, then comment, schedule, and error integration.
	allFields := model.ReplicationGroup("t", id.Name(), configAccounts, replicationGroupEnumStrings(objectTypes)).
		WithAllowedDatabases([]string{database.ID().FullyQualifiedName(), database2.ID().FullyQualifiedName()}).
		WithAllowedExternalVolumes([]string{externalVolume.FullyQualifiedName()}).
		WithAllowedShares([]string{share.ID().FullyQualifiedName()}).
		WithAllowedIntegrationTypes(replicationGroupEnumStrings(integrationTypes)).
		WithComment(comment).
		WithReplicationSchedule(schedule).
		WithErrorIntegration(errorIntegration.ID().FullyQualifiedName())
	objectLists := model.ReplicationGroup("t", id.Name(), configAccounts, replicationGroupEnumStrings(objectTypes)).
		WithAllowedDatabases([]string{database.ID().FullyQualifiedName(), database2.ID().FullyQualifiedName()}).
		WithAllowedExternalVolumes([]string{externalVolume.FullyQualifiedName()}).
		WithAllowedShares([]string{share.ID().FullyQualifiedName()}).
		WithAllowedIntegrationTypes(replicationGroupEnumStrings(integrationTypes))
	renamedId := testClient().Ids.RandomAccountObjectIdentifier()
	renamed := model.ReplicationGroup("t", renamedId.Name(), configAccounts, replicationGroupEnumStrings(objectTypes)).
		WithAllowedDatabases([]string{database.ID().FullyQualifiedName(), database2.ID().FullyQualifiedName()}).
		WithAllowedExternalVolumes([]string{externalVolume.FullyQualifiedName()}).
		WithAllowedShares([]string{share.ID().FullyQualifiedName()}).
		WithAllowedIntegrationTypes(replicationGroupEnumStrings(integrationTypes))
	ref := basic.ResourceReference()

	basicState := replicationGroupExpected{
		databases:   []sdk.AccountObjectIdentifier{database.ID()},
		objectTypes: []sdk.ReplicationGroupObjectType{sdk.ReplicationGroupObjectTypeDatabases},
	}
	allState := replicationGroupExpected{
		databases:        []sdk.AccountObjectIdentifier{database.ID(), database2.ID()},
		externalVolumes:  []sdk.AccountObjectIdentifier{externalVolume},
		shares:           []sdk.AccountObjectIdentifier{share.ID()},
		objectTypes:      objectTypes,
		integrationTypes: integrationTypes,
		comment:          comment,
		schedule:         schedule,
		errorIntegration: errorIntegration.ID().FullyQualifiedName(),
	}
	objectListsState := allState
	objectListsState.comment = ""
	objectListsState.schedule = ""
	objectListsState.errorIntegration = ""
	objectListsState.errorIntegrationCleared = true

	assertBasic := replicationGroupAssertions(t, ref, id, configAccounts, sourceAccount, basicState)
	assertAll := replicationGroupAssertions(t, ref, id, configAccounts, sourceAccount, allState)
	assertObjectLists := replicationGroupAssertions(t, ref, id, configAccounts, sourceAccount, objectListsState)
	assertRenamed := append(replicationGroupAssertions(t, ref, renamedId, configAccounts, sourceAccount, objectListsState),
		invokeactionassert.ReplicationGroupDoesNotExist(t, id),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ReplicationGroup),
		Steps: []resource.TestStep{
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				// SHOW does not return the allowed-object lists or error_integration.
				// ignore_edition_check is a create-only default.
				ImportStateVerifyIgnore: replicationGroupImportIgnores(),
				ImportStateCheck:        replicationGroupImportChecks(t, id, configAccounts, basicState),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, nil, sdk.String(comment)),
						planchecks.ExpectChange(ref, "replication_schedule", tfjson.ActionUpdate, nil, sdk.String(schedule)),
						planchecks.ExpectChange(ref, "error_integration", tfjson.ActionUpdate, nil, sdk.String(errorIntegration.ID().FullyQualifiedName())),
						planchecks.ExpectNoChangeOnField(ref, "name"),
					},
				},
				Config: config.FromModels(t, allFields),
				Check:  assertThat(t, assertAll...),
			},
			{
				Config:                  config.FromModels(t, allFields),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: replicationGroupImportIgnores(),
				ImportStateCheck:        replicationGroupImportChecks(t, id, configAccounts, allState),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(comment), nil),
						planchecks.ExpectChange(ref, "replication_schedule", tfjson.ActionUpdate, sdk.String(schedule), nil),
						planchecks.ExpectChange(ref, "error_integration", tfjson.ActionUpdate, sdk.String(errorIntegration.ID().FullyQualifiedName()), nil),
					},
				},
				Config: config.FromModels(t, objectLists),
				Check:  assertThat(t, assertObjectLists...),
			},
			{
				PreConfig: func() {
					testClient().ReplicationGroup.Alter(t, sdk.NewAlterReplicationGroupRequest(id).
						WithSet(*sdk.NewReplicationGroupSetRequest().
							WithComment(externalComment).
							WithReplicationSchedule(externalSchedule)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "comment", nil, sdk.String(externalComment)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), nil),
						planchecks.ExpectDrift(ref, "replication_schedule", nil, sdk.String(externalSchedule)),
						planchecks.ExpectChange(ref, "replication_schedule", tfjson.ActionUpdate, sdk.String(externalSchedule), nil),
					},
				},
				Config: config.FromModels(t, objectLists),
				Check:  assertThat(t, assertObjectLists...),
			},
			// rename
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "name", tfjson.ActionUpdate, sdk.String(id.Name()), sdk.String(renamedId.Name())),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: config.FromModels(t, renamed),
				Check:  assertThat(t, assertRenamed...),
			},
		},
	})
}

func TestAcc_ReplicationGroup_RecreateWhenExternallyNotPrimary(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	database, cleanupDatabase := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	sourceAccount := testClient().Account.GetAccountIdentifier(t)
	targetAccount := secondaryTestClient().Account.GetAccountIdentifier(t)

	// A primary in the other account, so this account can replace the managed group with a replica of the same name.
	otherDatabase, cleanupOtherDatabase := secondaryTestClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupOtherDatabase)
	otherPrimary := secondaryTestClient().Ids.RandomAccountObjectIdentifier()
	secondaryTestClient().ReplicationGroup.CreateWithRequest(t, sdk.NewCreateReplicationGroupRequest(
		otherPrimary,
		[]sdk.ReplicationGroupObjectType{sdk.ReplicationGroupObjectTypeDatabases},
		[]sdk.AccountIdentifier{sourceAccount},
	).WithAllowedDatabases([]sdk.AccountObjectIdentifier{otherDatabase.ID()}))
	t.Cleanup(func() {
		// The brief replica copies this database into the current account. Drop that copy after the group is gone.
		_ = testClient().Database.DropDatabase(t, otherDatabase.ID())
		secondaryTestClient().ReplicationGroup.Drop(t, otherPrimary)
	})

	group := model.ReplicationGroup("t", id.Name(), []string{targetAccount.FullyQualifiedName()}, []string{string(sdk.ReplicationGroupObjectTypeDatabases)}).
		WithAllowedDatabases([]string{database.ID().FullyQualifiedName()})
	ref := group.ResourceReference()
	otherOf := sdk.NewExternalObjectIdentifier(targetAccount, otherPrimary)
	basicState := replicationGroupExpected{
		databases:   []sdk.AccountObjectIdentifier{database.ID()},
		objectTypes: []sdk.ReplicationGroupObjectType{sdk.ReplicationGroupObjectTypeDatabases},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ReplicationGroup),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, group),
				Check:  assertThat(t, replicationGroupAssertions(t, ref, id, []string{targetAccount.FullyQualifiedName()}, sourceAccount, basicState)...),
			},
			{
				// Replace the primary with a replica of the same name. The resource must destroy it and create the primary again.
				PreConfig: func() {
					testClient().ReplicationGroup.Drop(t, id)
					testClient().ReplicationGroup.WaitForPrimaryToBeVisible(t, otherOf)
					testClient().ReplicationGroup.CreateSecondary(t, id, otherOf)
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
				Check: assertThat(t, replicationGroupAssertions(t, ref, id, []string{targetAccount.FullyQualifiedName()}, sourceAccount, basicState)...),
			},
		},
	})
}

type replicationGroupExpected struct {
	databases               []sdk.AccountObjectIdentifier
	externalVolumes         []sdk.AccountObjectIdentifier
	shares                  []sdk.AccountObjectIdentifier
	objectTypes             []sdk.ReplicationGroupObjectType
	integrationTypes        []sdk.ReplicationGroupIntegrationType
	comment                 string
	schedule                string
	errorIntegration        string
	errorIntegrationCleared bool
}

func replicationGroupAssertions(
	t *testing.T,
	ref string,
	id sdk.AccountObjectIdentifier,
	allowedAccounts []string,
	sourceAccount sdk.AccountIdentifier,
	expected replicationGroupExpected,
) []assert.TestCheckFuncProvider {
	t.Helper()

	primary := sdk.NewExternalObjectIdentifier(sourceAccount, id)
	resourceCheck := resourceassert.ReplicationGroupResource(t, ref).
		HasNameString(id.Name()).
		HasFullyQualifiedNameString(id.FullyQualifiedName()).
		HasObjectTypes(replicationGroupEnumStrings(expected.objectTypes)...).
		HasAllowedDatabases(replicationGroupIdentifierNames(expected.databases)...).
		HasAllowedAccounts(allowedAccounts...).
		HasIgnoreEditionCheckString(tfresources.BooleanDefault).
		HasIsPrimary(true)
	showCheck := resourceshowoutputassert.ReplicationGroupShowOutput(t, ref).
		HasName(id.Name()).
		HasType(sdk.ReplicationGroupTypeReplication).
		HasIsPrimary(true).
		HasPrimary(primary).
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasAccountName(sourceAccount.AccountName()).
		HasAccountLocator(testClient().GetAccountLocator()).
		HasOrganizationName(sourceAccount.OrganizationName()).
		HasIsListingAutoFulfillmentGroup(false)
	objectCheck := objectassert.ReplicationGroup(t, id).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasType(sdk.ReplicationGroupTypeReplication).
		HasIsPrimary(true).
		HasPrimary(primary).
		HasObjectTypes(expected.objectTypes...).
		HasNoSecondaryState().
		HasNoReadOnly().
		HasIsListingAutoFulfillmentGroup(false).
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasAccountName(sourceAccount.AccountName()).
		HasAccountLocator(testClient().GetAccountLocator()).
		HasOrganizationName(sourceAccount.OrganizationName()).
		HasSnowflakeRegionNotEmpty()
	if len(expected.externalVolumes) > 0 {
		resourceCheck = resourceCheck.HasAllowedExternalVolumes(replicationGroupIdentifierNames(expected.externalVolumes)...)
	}
	if len(expected.shares) > 0 {
		resourceCheck = resourceCheck.HasAllowedShares(replicationGroupIdentifierNames(expected.shares)...)
	}
	if len(expected.integrationTypes) == 0 {
		objectCheck = objectCheck.HasNoAllowedIntegrationTypes()
	} else {
		resourceCheck = resourceCheck.HasAllowedIntegrationTypes(replicationGroupEnumStrings(expected.integrationTypes)...)
		objectCheck = objectCheck.HasAllowedIntegrationTypes(expected.integrationTypes...)
	}
	if expected.comment == "" {
		resourceCheck = resourceCheck.HasCommentEmpty()
		showCheck = showCheck.HasComment("")
		objectCheck = objectCheck.HasNoComment()
	} else {
		resourceCheck = resourceCheck.HasCommentString(expected.comment)
		showCheck = showCheck.HasComment(expected.comment)
		objectCheck = objectCheck.HasComment(expected.comment)
	}
	if expected.schedule == "" {
		resourceCheck = resourceCheck.HasReplicationScheduleEmpty()
		showCheck = showCheck.HasReplicationSchedule("")
		objectCheck = objectCheck.HasNoReplicationSchedule()
	} else {
		resourceCheck = resourceCheck.HasReplicationScheduleString(expected.schedule)
		showCheck = showCheck.HasReplicationSchedule(expected.schedule)
		objectCheck = objectCheck.HasReplicationSchedule(expected.schedule)
	}
	if expected.errorIntegration == "" {
		if expected.errorIntegrationCleared {
			resourceCheck = resourceCheck.HasErrorIntegrationEmpty()
		} else {
			resourceCheck = resourceCheck.HasNoErrorIntegration()
		}
	} else {
		resourceCheck = resourceCheck.HasErrorIntegrationString(expected.errorIntegration)
	}
	return []assert.TestCheckFuncProvider{resourceCheck, showCheck, objectCheck}
}

func replicationGroupImportChecks(t *testing.T, id sdk.AccountObjectIdentifier, allowedAccounts []string, expected replicationGroupExpected) resource.ImportStateCheckFunc {
	t.Helper()
	encoded := helpers.EncodeResourceIdentifier(id)
	resourceCheck := resourceassert.ImportedReplicationGroupResource(t, encoded).
		HasNameString(id.Name()).
		HasObjectTypes(replicationGroupEnumStrings(expected.objectTypes)...).
		HasAllowedAccounts(allowedAccounts...).
		HasIsPrimary(true)
	showCheck := resourceshowoutputassert.ImportedReplicationGroupShowOutput(t, encoded).
		HasName(id.Name()).
		HasType(sdk.ReplicationGroupTypeReplication).
		HasIsPrimary(true)
	if len(expected.integrationTypes) == 0 {
		resourceCheck = resourceCheck.HasAllowedIntegrationTypesEmpty()
	} else {
		resourceCheck = resourceCheck.HasAllowedIntegrationTypes(replicationGroupEnumStrings(expected.integrationTypes)...)
	}
	if expected.comment == "" {
		resourceCheck = resourceCheck.HasCommentEmpty()
		showCheck = showCheck.HasComment("")
	} else {
		resourceCheck = resourceCheck.HasCommentString(expected.comment)
		showCheck = showCheck.HasComment(expected.comment)
	}
	if expected.schedule == "" {
		resourceCheck = resourceCheck.HasReplicationScheduleEmpty()
		showCheck = showCheck.HasReplicationSchedule("")
	} else {
		resourceCheck = resourceCheck.HasReplicationScheduleString(expected.schedule)
		showCheck = showCheck.HasReplicationSchedule(expected.schedule)
	}
	return assertThatImport(t, resourceCheck, showCheck)
}

func replicationGroupImportIgnores() []string {
	return []string{
		"allowed_databases",
		"allowed_external_volumes",
		"allowed_shares",
		"ignore_edition_check",
		"error_integration",
	}
}

func replicationGroupEnumStrings[T ~string](values []T) []string {
	out := make([]string, len(values))
	for i, value := range values {
		out[i] = string(value)
	}
	return out
}

func replicationGroupIdentifierNames(values []sdk.AccountObjectIdentifier) []string {
	out := make([]string, len(values))
	for i, value := range values {
		out[i] = value.FullyQualifiedName()
	}
	return out
}
