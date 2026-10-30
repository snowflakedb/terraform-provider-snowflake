//go:build non_account_level_tests

package testacc

import (
	"testing"

	tfresources "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
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
	// Config names the target account. Snowflake also adds the source account; the resource suppresses that member.
	configAccounts := []string{allowedAccount.FullyQualifiedName()}

	database2, cleanupDatabase2 := testClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase2)
	share, cleanupShare := testClient().Share.CreateShare(t)
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
	objectListsState.scheduleCleared = true
	objectListsState.errorIntegrationCleared = true

	assertBasic := replicationGroupAssertions(t, ref, id, configAccounts, sourceAccount, basicState)
	assertAll := replicationGroupAssertions(t, ref, id, configAccounts, sourceAccount, allState)
	assertObjectLists := replicationGroupAssertions(t, ref, id, configAccounts, sourceAccount, objectListsState)

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
				// ignore_edition_check is a create-only default. An omitted schedule comes back nil; import writes that nil as "".
				ImportStateVerifyIgnore: replicationGroupImportIgnores(true),
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
				ImportStateVerifyIgnore: replicationGroupImportIgnores(false),
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
						WithSet(*sdk.NewReplicationGroupSetRequest().WithComment(externalComment)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "comment", nil, sdk.String(externalComment)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), nil),
					},
				},
				Config: config.FromModels(t, objectLists),
				Check:  assertThat(t, assertObjectLists...),
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
	scheduleCleared         bool
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
		HasIgnoreEditionCheckString(tfresources.BooleanDefault)
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
		if expected.scheduleCleared {
			resourceCheck = resourceCheck.HasReplicationScheduleEmpty()
		} else {
			resourceCheck = resourceCheck.HasNoReplicationSchedule()
		}
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
		HasAllowedAccounts(allowedAccounts...)
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
	if expected.schedule != "" {
		resourceCheck = resourceCheck.HasReplicationScheduleString(expected.schedule)
		showCheck = showCheck.HasReplicationSchedule(expected.schedule)
	}
	return assertThatImport(t, resourceCheck, showCheck)
}

func replicationGroupImportIgnores(scheduleOmitted bool) []string {
	ignores := []string{
		"allowed_databases",
		"allowed_external_volumes",
		"allowed_shares",
		"ignore_edition_check",
		"error_integration",
	}
	if scheduleOmitted {
		ignores = append(ignores, "replication_schedule")
	}
	return ignores
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
