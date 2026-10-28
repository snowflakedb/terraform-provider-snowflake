//go:build non_account_level_tests

package testacc

import (
	"testing"

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
	"github.com/stretchr/testify/require"
)

func TestAcc_BackupPolicy_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()
	const expireAfterDays = 7
	const schedule = "2 HOUR"

	// CREATE requires a schedule, an expiration, or both. expire_after_days is the one required-by-Snowflake field.
	// with_retention_lock is ForceNew, so the locked policy is created only after the unlocked one is destroyed.
	basic := model.BackupPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithExpireAfterDays(expireAfterDays).
		WithWithRetentionLockBool(false)
	withOptionals := model.BackupPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithExpireAfterDays(expireAfterDays).
		WithSchedule(schedule).
		WithComment(comment).
		WithWithRetentionLockBool(false)
	complete := model.BackupPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithExpireAfterDays(expireAfterDays).
		WithSchedule(schedule).
		WithComment(comment).
		WithWithRetentionLockBool(true)
	ref := basic.ResourceReference()

	assertBasic := backupPolicyAssertions(t, ref, id, expireAfterDays, "", "", false, false)
	// SDKv2 keeps a string that was set and then unset as "" instead of removing it.
	assertAfterUnset := backupPolicyAssertions(t, ref, id, expireAfterDays, "", "", false, true)
	assertWithOptionals := backupPolicyAssertions(t, ref, id, expireAfterDays, schedule, comment, false, false)
	assertComplete := backupPolicyAssertions(t, ref, id, expireAfterDays, schedule, comment, true, false)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.BackupPolicy),
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
				// SHOW leaves schedule nil when it was not set. Read keeps the attribute absent
				// because the config omits it. Import writes that nil as "".
				ImportStateVerifyIgnore: []string{"schedule"},
				ImportStateCheck:        backupPolicyImportChecks(t, id, expireAfterDays, "", "", false),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, nil, sdk.String(comment)),
						planchecks.ExpectChange(ref, "schedule", tfjson.ActionUpdate, nil, sdk.String(schedule)),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			{
				Config:            config.FromModels(t, withOptionals),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  backupPolicyImportChecks(t, id, expireAfterDays, schedule, comment, false),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(comment), nil),
						planchecks.ExpectChange(ref, "schedule", tfjson.ActionUpdate, sdk.String(schedule), nil),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertAfterUnset...),
			},
			{
				Config:  config.FromModels(t, basic),
				Destroy: true,
			},
			{
				PreConfig: func() {
					_, err := testClient().BackupPolicy.Show(t, id)
					require.ErrorIs(t, err, sdk.ErrObjectNotFound)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
			{
				Config:            config.FromModels(t, complete),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  backupPolicyImportChecks(t, id, expireAfterDays, schedule, comment, true),
			},
			{
				PreConfig: func() {
					testClient().BackupPolicy.Alter(t, sdk.NewAlterBackupPolicyRequest(id).
						WithSet(*sdk.NewBackupPolicySetRequest().WithComment(externalComment)),
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externalComment)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), sdk.String(comment)),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}

func backupPolicyAssertions(t *testing.T, ref string, id sdk.SchemaObjectIdentifier, expireAfterDays int, schedule, comment string, retentionLock bool, scheduleLeftEmpty bool) []assert.TestCheckFuncProvider {
	t.Helper()

	resourceCheck := resourceassert.BackupPolicyResource(t, ref).
		HasNameString(id.Name()).
		HasDatabaseString(id.DatabaseName()).
		HasSchemaString(id.SchemaName()).
		HasExpireAfterDays(expireAfterDays).
		HasWithRetentionLockBool(retentionLock).
		HasFullyQualifiedNameString(id.FullyQualifiedName())
	showCheck := resourceshowoutputassert.BackupPolicyShowOutput(t, ref).
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasExpireAfterDays(expireAfterDays).
		HasHasRetentionLock(retentionLock).
		HasOwnerRoleType("ROLE")
	describeCheck := resourceshowoutputassert.BackupPolicyDescribeOutput(t, ref).
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasExpireAfterDays(expireAfterDays).
		HasHasRetentionLock(retentionLock).
		HasOwnerRoleType("ROLE")
	showObject := objectassert.BackupPolicy(t, id).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasExpireAfterDays(expireAfterDays).
		HasHasRetentionLock(retentionLock).
		HasOwnerRoleType("ROLE")
	details := objectassert.BackupPolicyDetails(t, id).
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasExpireAfterDays(expireAfterDays).
		HasHasRetentionLock(retentionLock).
		HasOwnerRoleType("ROLE")

	if comment == "" {
		resourceCheck = resourceCheck.HasCommentEmpty()
		showCheck = showCheck.HasComment("")
		describeCheck = describeCheck.HasComment("")
		showObject = showObject.HasComment("")
		details = details.HasComment("")
	} else {
		resourceCheck = resourceCheck.HasCommentString(comment)
		showCheck = showCheck.HasComment(comment)
		describeCheck = describeCheck.HasComment(comment)
		showObject = showObject.HasComment(comment)
		details = details.HasComment(comment)
	}
	if schedule == "" {
		if scheduleLeftEmpty {
			resourceCheck = resourceCheck.HasScheduleEmpty()
		} else {
			resourceCheck = resourceCheck.HasNoSchedule()
		}
		showCheck = showCheck.HasSchedule("")
		describeCheck = describeCheck.HasSchedule("")
		showObject = showObject.HasNoSchedule()
		details = details.HasNoSchedule()
	} else {
		resourceCheck = resourceCheck.HasScheduleString(schedule)
		showCheck = showCheck.HasSchedule(schedule)
		describeCheck = describeCheck.HasSchedule(schedule)
		showObject = showObject.HasSchedule(schedule)
		details = details.HasSchedule(schedule)
	}

	return []assert.TestCheckFuncProvider{resourceCheck, showCheck, describeCheck, showObject, details}
}

func backupPolicyImportChecks(t *testing.T, id sdk.SchemaObjectIdentifier, expireAfterDays int, schedule, comment string, retentionLock bool) resource.ImportStateCheckFunc {
	t.Helper()
	encoded := helpers.EncodeResourceIdentifier(id)
	resourceCheck := resourceassert.ImportedBackupPolicyResource(t, encoded).
		HasNameString(id.Name()).
		HasDatabaseString(id.DatabaseName()).
		HasSchemaString(id.SchemaName()).
		HasExpireAfterDays(expireAfterDays).
		HasWithRetentionLockBool(retentionLock)
	showCheck := resourceshowoutputassert.ImportedBackupPolicyShowOutput(t, encoded).
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasExpireAfterDays(expireAfterDays).
		HasHasRetentionLock(retentionLock).
		HasOwnerRoleType("ROLE")
	describeCheck := resourceshowoutputassert.ImportedBackupPolicyDescribeOutput(t, encoded).
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasExpireAfterDays(expireAfterDays).
		HasHasRetentionLock(retentionLock).
		HasOwnerRoleType("ROLE")
	if comment == "" {
		resourceCheck = resourceCheck.HasCommentEmpty()
		showCheck = showCheck.HasComment("")
		describeCheck = describeCheck.HasComment("")
	} else {
		resourceCheck = resourceCheck.HasCommentString(comment)
		showCheck = showCheck.HasComment(comment)
		describeCheck = describeCheck.HasComment(comment)
	}
	if schedule == "" {
		resourceCheck = resourceCheck.HasScheduleEmpty()
		showCheck = showCheck.HasSchedule("")
		describeCheck = describeCheck.HasSchedule("")
	} else {
		resourceCheck = resourceCheck.HasScheduleString(schedule)
		showCheck = showCheck.HasSchedule(schedule)
		describeCheck = describeCheck.HasSchedule(schedule)
	}
	return assertThatImport(t, resourceCheck, showCheck, describeCheck)
}
