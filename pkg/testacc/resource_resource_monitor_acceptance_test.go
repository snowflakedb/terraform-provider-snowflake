//go:build !account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	configvariable "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ResourceMonitor_Basic(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	configModel := model.ResourceMonitor("test", id.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configModel),
				Check: assertThat(t,
					resourceassert.ResourceMonitorResource(t, "snowflake_resource_monitor.test").
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasNoCreditQuota().
						HasNotifyUsersLen(0).
						HasNoFrequency().
						HasNoStartTimestamp().
						HasNoEndTimestamp().
						HasNotifyTriggersEmpty().
						HasNoSuspendTrigger().
						HasNoSuspendImmediateTrigger(),
					resourceshowoutputassert.ResourceMonitorShowOutput(t, "snowflake_resource_monitor.test").
						HasName(id.Name()).
						HasCreditQuota(0).
						HasUsedCredits(0).
						HasRemainingCredits(0).
						HasLevel("").
						HasFrequency(sdk.FrequencyMonthly).
						HasStartTimeNotEmpty().
						HasEndTime("").
						HasSuspendAt(0).
						HasSuspendImmediateAt(0).
						HasCreatedOnNotEmpty().
						HasOwnerNotEmpty().
						HasComment(""),
				),
			},
			{
				ResourceName: "snowflake_resource_monitor.test",
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedResourceMonitorResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCreditQuotaString("0").
						HasNotifyUsersLen(0).
						HasFrequencyString(string(sdk.FrequencyMonthly)).
						HasStartTimestampNotEmpty().
						HasEndTimestampString("").
						HasNotifyTriggersEmpty().
						HasSuspendTriggerString("0").
						HasSuspendImmediateTriggerString("0"),
				),
			},
		},
	})
}

func TestAcc_ResourceMonitor_Complete(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	configModel := model.ResourceMonitor("test", id.Name()).
		WithNotifyUsersValue(configvariable.SetVariable(configvariable.StringVariable("JAN_CIESLAK"))).
		WithCreditQuota(10).
		WithFrequency(string(sdk.FrequencyWeekly)).
		WithStartTimestamp(time.Now().Add(time.Hour * 24 * 30).Format("2006-01-02 15:01")).
		WithEndTimestamp(time.Now().Add(time.Hour * 24 * 60).Format("2006-01-02 15:01")).
		WithNotifyTriggersValue(configvariable.SetVariable(
			configvariable.IntegerVariable(100),
			configvariable.IntegerVariable(110),
		)).
		WithSuspendTrigger(120).
		WithSuspendImmediateTrigger(150)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configModel),
				Check: assertThat(t,
					resourceassert.ResourceMonitorResource(t, "snowflake_resource_monitor.test").
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCreditQuotaString("10").
						HasNotifyUsersLen(1).
						HasNotifyUser(0, "JAN_CIESLAK").
						HasFrequencyString(string(sdk.FrequencyWeekly)).
						HasStartTimestampString(time.Now().Add(time.Hour*24*30).Format("2006-01-02 15:01")).
						HasEndTimestampString(time.Now().Add(time.Hour*24*60).Format("2006-01-02 15:01")).
						HasNotifyTriggersLen(2).
						HasNotifyTrigger(0, 100).
						HasNotifyTrigger(1, 110).
						HasSuspendTriggerString("120").
						HasSuspendImmediateTriggerString("150"),
					resourceshowoutputassert.ResourceMonitorShowOutput(t, "snowflake_resource_monitor.test").
						HasName(id.Name()).
						HasCreditQuota(10).
						HasUsedCredits(0).
						HasRemainingCredits(10).
						HasLevel("").
						HasFrequency(sdk.FrequencyWeekly).
						HasStartTimeNotEmpty().
						HasEndTimeNotEmpty().
						HasSuspendAt(120).
						HasSuspendImmediateAt(150).
						HasCreatedOnNotEmpty().
						HasOwnerNotEmpty().
						HasComment(""),
				),
			},
			{
				ResourceName: "snowflake_resource_monitor.test",
				ImportState:  true,
				Config:       config.FromModels(t, configModel),
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedResourceMonitorResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCreditQuotaString("10").
						HasNotifyUsersLen(1).
						HasNotifyUser(0, "JAN_CIESLAK").
						HasFrequencyString(string(sdk.FrequencyWeekly)).
						HasStartTimestampNotEmpty().
						HasEndTimestampNotEmpty().
						HasNotifyTriggersLen(2).
						HasNotifyTrigger(0, 100).
						HasNotifyTrigger(1, 110).
						HasSuspendTriggerString("120").
						HasSuspendImmediateTriggerString("150"),
				),
			},
		},
	})
}

func TestAcc_ResourceMonitor_Updates(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	configModelNothingSet := model.ResourceMonitor("test", id.Name())

	configModelEverythingSet := model.ResourceMonitor("test", id.Name()).
		WithNotifyUsersValue(configvariable.SetVariable(configvariable.StringVariable("JAN_CIESLAK"))).
		WithCreditQuota(10).
		WithFrequency(string(sdk.FrequencyWeekly)).
		WithStartTimestamp(time.Now().Add(time.Hour * 24 * 30).Format("2006-01-02 15:01")).
		WithEndTimestamp(time.Now().Add(time.Hour * 24 * 60).Format("2006-01-02 15:01")).
		WithNotifyTriggersValue(configvariable.SetVariable(
			configvariable.IntegerVariable(100),
			configvariable.IntegerVariable(110),
		)).
		WithSuspendTrigger(120).
		WithSuspendImmediateTrigger(150)

	configModelUpdated := model.ResourceMonitor("test", id.Name()).
		WithNotifyUsersValue(configvariable.SetVariable(configvariable.StringVariable("JAN_CIESLAK"), configvariable.StringVariable("ARTUR_SAWICKI"))).
		WithCreditQuota(20).
		WithFrequency(string(sdk.FrequencyMonthly)).
		WithStartTimestamp(time.Now().Add(time.Hour * 24 * 40).Format("2006-01-02 15:01")).
		WithEndTimestamp(time.Now().Add(time.Hour * 24 * 70).Format("2006-01-02 15:01")).
		WithNotifyTriggersValue(configvariable.SetVariable(
			configvariable.IntegerVariable(110),
			configvariable.IntegerVariable(120),
		)).
		WithSuspendTrigger(130).
		WithSuspendImmediateTrigger(160)

	configModelEverythingUnset := model.ResourceMonitor("test", id.Name()).
		WithSuspendTrigger(130) // cannot completely remove all triggers (Snowflake limitation; tested below)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configModelNothingSet),
				Check: assertThat(t,
					resourceassert.ResourceMonitorResource(t, "snowflake_resource_monitor.test").
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasNoCreditQuota().
						HasNotifyUsersLen(0).
						HasNoFrequency().
						HasNoStartTimestamp().
						HasNoEndTimestamp().
						HasNotifyTriggersLen(0).
						HasNoSuspendTrigger().
						HasNoSuspendImmediateTrigger(),
					resourceshowoutputassert.ResourceMonitorShowOutput(t, "snowflake_resource_monitor.test").
						HasName(id.Name()).
						HasCreditQuota(0).
						HasUsedCredits(0).
						HasRemainingCredits(0).
						HasLevel("").
						HasFrequency(sdk.FrequencyMonthly).
						HasStartTimeNotEmpty().
						HasEndTime("").
						HasSuspendAt(0).
						HasSuspendImmediateAt(0).
						HasCreatedOnNotEmpty().
						HasOwnerNotEmpty().
						HasComment(""),
				),
			},
			// Set
			{
				Config: config.FromModels(t, configModelEverythingSet),
				Check: assertThat(t,
					resourceassert.ResourceMonitorResource(t, "snowflake_resource_monitor.test").
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCreditQuotaString("10").
						HasNotifyUsersLen(1).
						HasNotifyUser(0, "JAN_CIESLAK").
						HasFrequencyString(string(sdk.FrequencyWeekly)).
						HasStartTimestampString(time.Now().Add(time.Hour*24*30).Format("2006-01-02 15:01")).
						HasEndTimestampString(time.Now().Add(time.Hour*24*60).Format("2006-01-02 15:01")).
						HasNotifyTriggersLen(2).
						HasNotifyTrigger(0, 100).
						HasNotifyTrigger(1, 110).
						HasSuspendTriggerString("120").
						HasSuspendImmediateTriggerString("150"),
					resourceshowoutputassert.ResourceMonitorShowOutput(t, "snowflake_resource_monitor.test").
						HasName(id.Name()).
						HasCreditQuota(10).
						HasUsedCredits(0).
						HasRemainingCredits(10).
						HasLevel("").
						HasFrequency(sdk.FrequencyWeekly).
						HasStartTimeNotEmpty().
						HasEndTimeNotEmpty().
						HasSuspendAt(120).
						HasSuspendImmediateAt(150).
						HasCreatedOnNotEmpty().
						HasOwnerNotEmpty().
						HasComment(""),
				),
			},
			// Update
			{
				Config: config.FromModels(t, configModelUpdated),
				Check: assertThat(t,
					resourceassert.ResourceMonitorResource(t, "snowflake_resource_monitor.test").
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCreditQuotaString("20").
						HasNotifyUsersLen(2).
						HasNotifyUser(0, "ARTUR_SAWICKI").
						HasNotifyUser(1, "JAN_CIESLAK").
						HasFrequencyString(string(sdk.FrequencyMonthly)).
						HasStartTimestampString(time.Now().Add(time.Hour*24*40).Format("2006-01-02 15:01")).
						HasEndTimestampString(time.Now().Add(time.Hour*24*70).Format("2006-01-02 15:01")).
						HasNotifyTriggersLen(2).
						HasNotifyTrigger(0, 110).
						HasNotifyTrigger(1, 120).
						HasSuspendTriggerString("130").
						HasSuspendImmediateTriggerString("160"),
					resourceshowoutputassert.ResourceMonitorShowOutput(t, "snowflake_resource_monitor.test").
						HasName(id.Name()).
						HasCreditQuota(20).
						HasUsedCredits(0).
						HasRemainingCredits(20).
						HasLevel("").
						HasFrequency(sdk.FrequencyMonthly).
						HasStartTimeNotEmpty().
						HasEndTimeNotEmpty().
						HasSuspendAt(130).
						HasSuspendImmediateAt(160).
						HasCreatedOnNotEmpty().
						HasOwnerNotEmpty().
						HasComment(""),
				),
			},
			// Unset
			{
				Config: config.FromModels(t, configModelEverythingUnset),
				Check: assertThat(t,
					resourceassert.ResourceMonitorResource(t, "snowflake_resource_monitor.test").
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCreditQuotaString("0").
						HasNotifyUsersLen(0).
						HasFrequencyString("").
						HasStartTimestampString("").
						HasEndTimestampString("").
						HasSuspendTriggerString("130"),
					resourceshowoutputassert.ResourceMonitorShowOutput(t, "snowflake_resource_monitor.test").
						HasName(id.Name()).
						HasCreditQuota(0).
						HasUsedCredits(0).
						HasRemainingCredits(0).
						HasLevel("").
						HasFrequency(sdk.FrequencyMonthly).
						HasStartTimeNotEmpty().
						HasEndTime("").
						HasSuspendAt(130).
						HasSuspendImmediateAt(0).
						HasCreatedOnNotEmpty().
						HasOwnerNotEmpty().
						HasComment(""),
				),
			},
		},
	})
}

func TestAcc_ResourceMonitor_ExternalChanges(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	startTimestamp := time.Now().Add(time.Hour * 24 * 40).Format("2006-01-02 15:01")
	endTimestamp := time.Now().Add(time.Hour * 24 * 70).Format("2006-01-02 15:01")
	configModelEverythingSet := model.ResourceMonitor("test", id.Name()).
		WithNotifyUsersValue(configvariable.SetVariable(configvariable.StringVariable("JAN_CIESLAK"))).
		WithCreditQuota(10).
		WithFrequency(string(sdk.FrequencyWeekly)).
		WithStartTimestamp(startTimestamp).
		WithEndTimestamp(endTimestamp).
		WithNotifyTriggersValue(configvariable.SetVariable(
			configvariable.IntegerVariable(100),
			configvariable.IntegerVariable(110),
		)).
		WithSuspendTrigger(120).
		WithSuspendImmediateTrigger(150)

	configModelUpdated := model.ResourceMonitor("test", id.Name()).
		WithNotifyUsersValue(configvariable.SetVariable(configvariable.StringVariable("JAN_CIESLAK"), configvariable.StringVariable("ARTUR_SAWICKI"))).
		WithCreditQuota(20).
		WithFrequency(string(sdk.FrequencyMonthly)).
		WithStartTimestamp(startTimestamp).
		WithEndTimestamp(endTimestamp).
		WithNotifyTriggersValue(configvariable.SetVariable(
			configvariable.IntegerVariable(110),
			configvariable.IntegerVariable(120),
		)).
		WithSuspendTrigger(130).
		WithSuspendImmediateTrigger(160)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configModelEverythingSet),
			},
			// Update externally, but match the updated configuration (expected updates to the same values)
			{
				PreConfig: func() {
					testClient().ResourceMonitor.Alter(t, id, &sdk.AlterResourceMonitorOptions{
						Set: &sdk.ResourceMonitorSet{
							NotifyUsers: &sdk.NotifyUsers{
								Users: []sdk.NotifiedUser{
									{Name: sdk.NewAccountObjectIdentifier("JAN_CIESLAK")},
									{Name: sdk.NewAccountObjectIdentifier("ARTUR_SAWICKI")},
								},
							},
							CreditQuota:    sdk.Int(20),
							Frequency:      sdk.Pointer(sdk.FrequencyMonthly),
							StartTimestamp: sdk.String(startTimestamp),
							EndTimestamp:   sdk.String(endTimestamp),
						},
						Triggers: []sdk.TriggerDefinition{
							{
								Threshold:     110,
								TriggerAction: sdk.TriggerActionNotify,
							},
							{
								Threshold:     120,
								TriggerAction: sdk.TriggerActionNotify,
							},
							{
								Threshold:     130,
								TriggerAction: sdk.TriggerActionSuspend,
							},
							{
								Threshold:     160,
								TriggerAction: sdk.TriggerActionSuspendImmediate,
							},
						},
					})
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						planchecks.PrintPlanDetails(configModelUpdated.ResourceReference(), "credit_quota", "end_timestamp", "frequency", "fully_qualified_name", "name", "notify_triggers", "notify_users", "start_timestamp", "suspend_immediate_trigger", "suspend_trigger", r.ShowOutputAttributeName),
					},
				},
				Config: config.FromModels(t, configModelUpdated),
			},
		},
	})
}

// TestAcc_ResourceMonitor_PartialUpdate covers a situation where alter fails. In the previous versions, the alter would
// fail, but invalid values would be saved in the state anyway. In the new version, the old values in state will be preserved
// because the old values are also stored on the Snowflake side (they weren't altered).
func TestAcc_ResourceMonitor_PartialUpdate(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	validTimestamp := time.Now().Add(time.Hour * 24 * 60).Format("2006-01-02 15:01")
	configModel := model.ResourceMonitor("test", id.Name()).
		WithEndTimestamp(validTimestamp)

	configModelInvalidUpdate := model.ResourceMonitor("test", id.Name()).
		WithEndTimestamp(time.Now().Add(time.Hour*24*70).Format("2006-01-02 15:01") + "abc")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configModel),
			},
			{
				Config:      config.FromModels(t, configModelInvalidUpdate),
				ExpectError: regexp.MustCompile("Invalid date/time format string"),
				Check: assertThat(t,
					resourceassert.ResourceMonitorResource(t, "snowflake_resource_monitor.test").
						HasEndTimestampString(validTimestamp),
				),
			},
			// Without the partials plan check failed.
			// The following was printed (indicating the invalid value was saved into the state):
			// ComputedIfAnyAttributeChanged: changed key: end_timestamp old: 2024-11-19 10:11abc new: 2024-11-09 10:11
			{
				Config: config.FromModels(t, configModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(t,
					resourceassert.ResourceMonitorResource(t, "snowflake_resource_monitor.test").
						HasEndTimestampString(validTimestamp),
				),
			},
		},
	})
}

// TestAcc_ResourceMonitor_issue2167 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2167 issue.
// Second step is purposely error, because tests TestAcc_ResourceMonitorUpdateNotifyUsers and TestAcc_ResourceMonitorNotifyUsers are still skipped.
// It can be fixed with them.
func TestAcc_ResourceMonitor_issue2167(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	configNoUsers := model.ResourceMonitor("test", id.Name()).WithNotifyUsersValue(config.EmptyListVariable())
	configWithNonExistingUser := model.ResourceMonitor("test", id.Name()).WithNotifyUsersValue(configvariable.SetVariable(configvariable.StringVariable("non_existing_user")))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configNoUsers),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "name", id.Name()),
				),
			},
			{
				Config:      config.FromModels(t, configWithNonExistingUser),
				ExpectError: regexp.MustCompile(`.*090268 \(22023\): User non_existing_user does not exist.*`),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1990 is fixed
func TestAcc_ResourceMonitor_Issue1990_RemovingResourceMonitorOutsideOfTerraform(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	configModel := model.ResourceMonitor("test", id.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			// Create resource monitor
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.69.0"),
				Config:            config.FromModels(t, configModel),
			},
			// Same configuration, but we drop resource monitor externally
			{
				ExternalProviders: ExternalProviderWithExactVersion("0.69.0"),
				PreConfig: func() {
					testClient().ResourceMonitor.DropResourceMonitorFunc(t, id)()
				},
				Config:      config.FromModels(t, configModel),
				ExpectError: regexp.MustCompile("object does not exist or not authorized"),
			},
			// Same configuration, but it's the last version where it's still not working
			{
				ExternalProviders: ExternalProviderWithExactVersion("0.95.0"),
				Config:            config.FromModels(t, configModel),
				ExpectError:       regexp.MustCompile("object does not exist or not authorized"),
			},
			// Same configuration, but it's the latest version of the provider (0.96.0 and above)
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModel),
			},
		},
	})
}

// proves
// https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1821
// https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1832
// https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1624
// https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1716
// https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1754
// are fixed and errors are more meaningful for the user
func TestAcc_ResourceMonitor_Issue_TimestampInfinitePlan(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	configModel := model.ResourceMonitor("test", id.Name())
	configModelWithDateStartTimestamp := model.ResourceMonitor("test", id.Name()).
		WithFrequency(string(sdk.FrequencyWeekly)).
		WithStartTimestamp(time.Now().Add(time.Hour * 24 * 30).Format("2006-01-02")).
		WithEndTimestamp(time.Now().Add(time.Hour * 24 * 60).Format("2006-01-02"))
	configModelWithDateTimeFormat := model.ResourceMonitor("test", id.Name()).
		WithFrequency(string(sdk.FrequencyWeekly)).
		WithStartTimestamp(time.Now().Add(time.Hour * 24 * 30).Format("2006-01-02 15:04")).
		WithEndTimestamp(time.Now().Add(time.Hour * 24 * 60).Format("2006-01-02 15:04"))

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			// Create resource monitor without the timestamps
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.90.0"),
				Config:            config.FromModels(t, configModel),
			},
			// Alter resource timestamps to have the following format: 2006-01-02 (produces a plan because of the format difference)
			{
				ExternalProviders:  ExternalProviderWithExactVersion("0.90.0"),
				Config:             config.FromModels(t, configModelWithDateStartTimestamp),
				ExpectNonEmptyPlan: true,
			},
			// Alter resource timestamps to have the following format: 2006-01-02 15:04 (won't produce plan because of the internal format mapping to this exact format)
			{
				ExternalProviders: ExternalProviderWithExactVersion("0.90.0"),
				Config:            config.FromModels(t, configModelWithDateTimeFormat),
			},
			// Destroy the resource
			{
				ExternalProviders: ExternalProviderWithExactVersion("0.90.0"),
				Config:            config.FromModels(t, configModelWithDateTimeFormat),
				Destroy:           true,
			},
			// Create resource monitor without the timestamps
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModel),
			},
			// Alter resource timestamps to have the following format: 2006-01-02 (no plan produced)
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModelWithDateStartTimestamp),
			},
			// Alter resource timestamps to have the following format: 2006-01-02 15:04 (no plan produced and the internal mapping is not applied in this version)
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModelWithDateTimeFormat),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1500 is fixed and errors are more meaningful for the user
func TestAcc_ResourceMonitor_Issue1500_CreatingWithOnlyTriggers(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	configModel := model.ResourceMonitor("test", id.Name()).
		WithNotifyTriggersValue(configvariable.SetVariable(
			configvariable.IntegerVariable(100),
			configvariable.IntegerVariable(110),
		)).
		WithSuspendTrigger(120).
		WithSuspendImmediateTrigger(150)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			// Create resource monitor with only triggers (old version)
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.90.0"),
				Config:            config.FromModels(t, configModel),
				ExpectError:       regexp.MustCompile("SQL compilation error"),
			},
			// Create resource monitor with only triggers (the latest version)
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModel),
				ExpectError:              regexp.MustCompile("due to Snowflake limitations you cannot create Resource Monitor with only triggers set"),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1500 is fixed and errors are more meaningful for the user
func TestAcc_ResourceMonitor_Issue1500_AlteringWithOnlyTriggers(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	configModelWithCreditQuota := model.ResourceMonitor("test", id.Name()).
		WithCreditQuota(100).
		WithNotifyTriggersValue(configvariable.SetVariable(
			configvariable.IntegerVariable(100),
			configvariable.IntegerVariable(110),
		)).
		WithSuspendTrigger(120).
		WithSuspendImmediateTrigger(150)

	configModelWithUpdatedTriggers := model.ResourceMonitor("test", id.Name()).
		WithCreditQuota(100).
		WithNotifyTriggersValue(configvariable.SetVariable(
			configvariable.IntegerVariable(110),
			configvariable.IntegerVariable(120),
		)).
		WithSuspendTrigger(130).
		WithSuspendImmediateTrigger(160)

	configModelWithoutTriggers := model.ResourceMonitor("test", id.Name()).
		WithCreditQuota(100)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.90.0"),
				Config:            config.FromModels(t, configModelWithCreditQuota),
			},
			// Update only triggers (not allowed in Snowflake)
			{
				ExternalProviders: ExternalProviderWithExactVersion("0.90.0"),
				Config:            config.FromModels(t, configModelWithUpdatedTriggers),
				// For some reason, not returning error (SQL compilation error should be returned in this case; most likely update was handled incorrectly, or it was handled similarly as in the current version)
			},
			// Remove all triggers (not allowed in Snowflake)
			{
				ExternalProviders: ExternalProviderWithExactVersion("0.90.0"),
				Config:            config.FromModels(t, configModelWithoutTriggers),
				// For some reason, not returning the correct error (SQL compilation error should be returned in this case; most likely update was processed incorrectly)
				ExpectError: regexp.MustCompile(`at least one of AlterResourceMonitorOptions fields \[Set Triggers] must be set`),
			},
			// Upgrade to the latest version
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModelWithCreditQuota),
			},
			// Update only triggers (not allowed in Snowflake)
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModelWithUpdatedTriggers),
			},
			// Update only triggers (not allowed in Snowflake; recreating)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_resource_monitor.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModelWithoutTriggers),
			},
		},
	})
}

func TestAcc_ResourceMonitor_RemovingAllTriggers(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	configModelWithNotifyTriggers := model.ResourceMonitor("test", id.Name()).
		WithCreditQuota(100).
		WithNotifyTriggersValue(configvariable.SetVariable(
			configvariable.IntegerVariable(100),
			configvariable.IntegerVariable(110),
		))

	configModelWithSuspendTrigger := model.ResourceMonitor("test", id.Name()).
		WithCreditQuota(100).
		WithSuspendTrigger(120)

	configModelWithSuspendImmediateTrigger := model.ResourceMonitor("test", id.Name()).
		WithCreditQuota(100).
		WithSuspendImmediateTrigger(120)

	configModelWithAllTriggers := model.ResourceMonitor("test", id.Name()).
		WithCreditQuota(100).
		WithNotifyTriggersValue(configvariable.SetVariable(
			configvariable.IntegerVariable(100),
			configvariable.IntegerVariable(110),
		)).
		WithSuspendTrigger(120).
		WithSuspendImmediateTrigger(150)

	configModelWithoutTriggers := model.ResourceMonitor("test", id.Name()).
		WithCreditQuota(100)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			// Config with all triggers
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModelWithAllTriggers),
			},
			// No triggers (force new expected)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_resource_monitor.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModelWithoutTriggers),
			},
			// Config with only notify triggers
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModelWithNotifyTriggers),
			},
			// No triggers (force new expected)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_resource_monitor.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModelWithoutTriggers),
			},
			// Config with only suspend trigger
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModelWithSuspendTrigger),
			},
			// No triggers (force new expected)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_resource_monitor.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModelWithoutTriggers),
			},
			// Config with only suspend immediate trigger
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModelWithSuspendImmediateTrigger),
			},
			// No triggers (force new expected)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_resource_monitor.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModelWithoutTriggers),
			},
		},
	})
}

// proves that fields that were present in the previous versions are not kept in the state after the upgrade
func TestAcc_ResourceMonitor_SetForWarehouse(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	newVersionModel := model.ResourceMonitor("test", id.Name()).
		WithCreditQuota(100).
		WithSuspendTrigger(100)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.90.0"),
				Config: fmt.Sprintf(`
resource "snowflake_resource_monitor" "test" {
	name = "%[1]s"
	credit_quota = 100
	suspend_trigger = 100
	warehouses = [ "%[2]s" ]
}
`, id.Name(), testClient().Ids.SnowflakeWarehouseId().Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_resource_monitor.test", "warehouses.#", "1"),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, newVersionModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(newVersionModel.ResourceReference(), "warehouses"),
					resource.TestCheckNoResourceAttr(newVersionModel.ResourceReference(), "warehouses.#"),
				),
			},
		},
	})
}
