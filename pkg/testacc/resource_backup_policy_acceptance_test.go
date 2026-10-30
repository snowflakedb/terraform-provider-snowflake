//go:build non_account_level_tests

package testacc

import (
	"regexp"
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
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

// The schedule block with 2 hours is the schedule string of SHOW.
const (
	scheduleHours = 2
	schedule      = "2 HOUR"
)

func TestAcc_BackupPolicy_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()
	const expireAfterDays = 7

	// CREATE requires a schedule, an expiration, or both. expire_after_days is the one required-by-Snowflake field.
	// with_retention_lock is ForceNew, so the locked policy is created only after the unlocked one is destroyed.
	basic := model.BackupPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithExpireAfterDays(expireAfterDays).
		WithWithRetentionLockBool(false)
	withOptionals := model.BackupPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithExpireAfterDays(expireAfterDays).
		WithScheduleHours(scheduleHours).
		WithComment(comment).
		WithWithRetentionLockBool(false)
	complete := model.BackupPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithExpireAfterDays(expireAfterDays).
		WithScheduleHours(scheduleHours).
		WithComment(comment).
		WithWithRetentionLockBool(true)
	ref := basic.ResourceReference()

	assertBasic := backupPolicyAssertions(t, ref, id, expireAfterDays, "", "", false)
	assertAfterUnset := backupPolicyAssertions(t, ref, id, expireAfterDays, "", "", false)
	assertWithOptionals := backupPolicyAssertions(t, ref, id, expireAfterDays, schedule, comment, false)
	assertComplete := backupPolicyAssertions(t, ref, id, expireAfterDays, schedule, comment, true)

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
				ImportStateCheck:  backupPolicyImportChecks(t, id, expireAfterDays, "", "", false),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, nil, sdk.String(comment)),
						planchecks.ExpectChange(ref, "schedule.0.hours", tfjson.ActionUpdate, nil, sdk.String("2")),
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
						planchecks.ExpectChange(ref, "schedule.0.hours", tfjson.ActionUpdate, sdk.String("2"), nil),
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

func backupPolicyAssertions(t *testing.T, ref string, id sdk.SchemaObjectIdentifier, expireAfterDays int, schedule, comment string, retentionLock bool) []assert.TestCheckFuncProvider {
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
		resourceCheck = resourceCheck.HasNoSchedule()
		showCheck = showCheck.HasSchedule("")
		describeCheck = describeCheck.HasSchedule("")
		showObject = showObject.HasNoSchedule()
		details = details.HasNoSchedule()
	} else {
		resourceCheck = resourceCheck.HasScheduleHours(scheduleHours)
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
		resourceCheck = resourceCheck.HasNoSchedule()
		showCheck = showCheck.HasSchedule("")
		describeCheck = describeCheck.HasSchedule("")
	} else {
		resourceCheck = resourceCheck.HasScheduleHours(scheduleHours)
		showCheck = showCheck.HasSchedule(schedule)
		describeCheck = describeCheck.HasSchedule(schedule)
	}
	return assertThatImport(t, resourceCheck, showCheck, describeCheck)
}

func TestAcc_BackupPolicy_ScheduleVariants(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	const cron = "0 0 * * * UTC"

	newModel := func() *model.BackupPolicyModel {
		return model.BackupPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name()).WithExpireAfterDays(7).WithWithRetentionLockBool(false)
	}
	hours := newModel().WithScheduleHours(2)
	moreHours := newModel().WithScheduleHours(3)
	moreMinutes := newModel().WithScheduleMinutes(120)
	minutes := newModel().WithScheduleMinutes(90)
	usingCron := newModel().WithScheduleCron(cron)
	ref := hours.ResourceReference()

	step := func(m *model.BackupPolicyModel, resourceCheck func(*resourceassert.BackupPolicyResourceAssert) *resourceassert.BackupPolicyResourceAssert, showSchedule string) resource.TestStep {
		return resource.TestStep{
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
				},
			},
			Config: config.FromModels(t, m),
			Check: assertThat(t,
				resourceCheck(resourceassert.BackupPolicyResource(t, ref)),
				resourceshowoutputassert.BackupPolicyShowOutput(t, ref).HasSchedule(showSchedule),
				resourceshowoutputassert.BackupPolicyDescribeOutput(t, ref).HasSchedule(showSchedule),
				objectassert.BackupPolicy(t, id).HasSchedule(showSchedule),
			),
		}
	}

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
				Config: config.FromModels(t, hours),
				Check: assertThat(t,
					resourceassert.BackupPolicyResource(t, ref).HasScheduleHours(2),
					resourceshowoutputassert.BackupPolicyShowOutput(t, ref).HasSchedule("2 HOUR"),
					objectassert.BackupPolicy(t, id).HasSchedule("2 HOUR"),
				),
			},
			{
				Config:            config.FromModels(t, hours),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			step(minutes, func(a *resourceassert.BackupPolicyResourceAssert) *resourceassert.BackupPolicyResourceAssert {
				return a.HasScheduleMinutes(90)
			}, "90 MINUTE"),
			step(usingCron, func(a *resourceassert.BackupPolicyResourceAssert) *resourceassert.BackupPolicyResourceAssert {
				return a.HasScheduleCron(cron)
			}, "USING CRON "+cron),
			{
				Config:            config.FromModels(t, usingCron),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			step(hours, func(a *resourceassert.BackupPolicyResourceAssert) *resourceassert.BackupPolicyResourceAssert {
				return a.HasScheduleHours(2)
			}, "2 HOUR"),
			// The value of the same sub-attribute changes in place.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "schedule.0.hours", tfjson.ActionUpdate, sdk.String("2"), sdk.String("3")),
					},
				},
				Config: config.FromModels(t, moreHours),
				Check: assertThat(t,
					resourceassert.BackupPolicyResource(t, ref).HasScheduleHours(3),
					resourceshowoutputassert.BackupPolicyShowOutput(t, ref).HasSchedule("3 HOUR"),
					objectassert.BackupPolicy(t, id).HasSchedule("3 HOUR"),
				),
			},
			step(moreMinutes, func(a *resourceassert.BackupPolicyResourceAssert) *resourceassert.BackupPolicyResourceAssert {
				return a.HasScheduleMinutes(120)
			}, "120 MINUTE"),
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "schedule.0.minutes", tfjson.ActionUpdate, sdk.String("120"), sdk.String("90")),
					},
				},
				Config: config.FromModels(t, minutes),
				Check:  assertThat(t, resourceassert.BackupPolicyResource(t, ref).HasScheduleMinutes(90)),
			},
		},
	})
}

func TestAcc_BackupPolicy_ScheduleSchemaValidation(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newModel := func() *model.BackupPolicyModel {
		return model.BackupPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name()).WithExpireAfterDays(7)
	}

	hoursAndMinutes := newModel().WithScheduleValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"hours":   tfconfig.IntegerVariable(2),
		"minutes": tfconfig.IntegerVariable(90),
	}))
	hoursAndCron := newModel().WithScheduleValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"hours":      tfconfig.IntegerVariable(2),
		"using_cron": tfconfig.StringVariable("0 0 * * * UTC"),
	}))
	minutesAndCron := newModel().WithScheduleValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"minutes":    tfconfig.IntegerVariable(90),
		"using_cron": tfconfig.StringVariable("0 0 * * * UTC"),
	}))
	allOfThem := newModel().WithScheduleValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"hours":      tfconfig.IntegerVariable(2),
		"minutes":    tfconfig.IntegerVariable(90),
		"using_cron": tfconfig.StringVariable("0 0 * * * UTC"),
	}))
	emptyBlock := newModel().WithScheduleValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
	}))
	zeroHours := newModel().WithScheduleHours(0)
	tooFewMinutes := newModel().WithScheduleMinutes(59)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.BackupPolicy),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, hoursAndMinutes),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Conflicting configuration arguments"),
			},
			{
				Config:      config.FromModels(t, hoursAndCron),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Conflicting configuration arguments"),
			},
			{
				Config:      config.FromModels(t, minutesAndCron),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Conflicting configuration arguments"),
			},
			{
				Config:      config.FromModels(t, allOfThem),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Conflicting configuration arguments"),
			},
			{
				Config:      config.FromModels(t, emptyBlock),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      config.FromModels(t, zeroHours),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected hours to be at least \(1\), got 0`),
			},
			{
				Config:      config.FromModels(t, tooFewMinutes),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected minutes to be at least \(60\), got 59`),
			},
		},
	})
}

// Snowflake rejects UNSET of the schedule and the expiration in one statement, and a policy always keeps at least one of them
// ("Cannot unset expire_after_days if schedule is also not set."). The only way to drop one is to set the other in the same update,
// so the resource sends SET before UNSET.
func TestAcc_BackupPolicy_SwitchBetweenScheduleAndExpiration(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newModel := func() *model.BackupPolicyModel {
		return model.BackupPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name()).WithWithRetentionLockBool(false)
	}
	scheduleOnly := newModel().WithScheduleHours(2)
	expirationOnly := newModel().WithExpireAfterDays(7)
	ref := scheduleOnly.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.BackupPolicy),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, scheduleOnly),
				Check: assertThat(t,
					resourceassert.BackupPolicyResource(t, ref).HasScheduleHours(2),
					objectassert.BackupPolicy(t, id).HasSchedule("2 HOUR").HasNoExpireAfterDays(),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, expirationOnly),
				Check: assertThat(t,
					resourceassert.BackupPolicyResource(t, ref).HasNoSchedule().HasExpireAfterDays(7),
					objectassert.BackupPolicy(t, id).HasNoSchedule().HasExpireAfterDays(7),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, scheduleOnly),
				Check: assertThat(t,
					resourceassert.BackupPolicyResource(t, ref).HasScheduleHours(2),
					objectassert.BackupPolicy(t, id).HasSchedule("2 HOUR").HasNoExpireAfterDays(),
				),
			},
		},
	})
}
