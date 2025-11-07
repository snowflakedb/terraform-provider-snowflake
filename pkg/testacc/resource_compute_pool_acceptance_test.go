//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ComputePool_basic(t *testing.T) {
	application := createApp(t)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment, changedComment := random.Comment(), random.Comment()

	modelBasic := model.ComputePool("test", id.Name(), string(sdk.ComputePoolInstanceFamilyCpuX64S), 2, 1)

	modelComplete := model.ComputePool("test", id.Name(), string(sdk.ComputePoolInstanceFamilyCpuX64S), 2, 1).
		WithForApplication(application.ID().FullyQualifiedName()).
		WithAutoResume("true").
		WithInitiallySuspended("true").
		WithAutoSuspendSecs(6767).
		WithComment(comment)

	modelCompleteWithDifferentValues := model.ComputePool("test", id.Name(), string(sdk.ComputePoolInstanceFamilyCpuX64S), 4, 3).
		WithForApplication(application.ID().FullyQualifiedName()).
		WithAutoResume("true").
		WithInitiallySuspended("true").
		WithAutoSuspendSecs(2222).
		WithComment(changedComment)

	modelBasicWithApp := model.ComputePool("test", id.Name(), string(sdk.ComputePoolInstanceFamilyCpuX64S), 4, 3).
		WithInitiallySuspended("true").
		WithForApplication(application.ID().FullyQualifiedName())
	modelCompleteWithDifferentInstanceFamily := model.ComputePool("test", id.Name(), string(sdk.ComputePoolInstanceFamilyCpuX64M), 4, 3).
		WithInitiallySuspended("true").
		WithForApplication(application.ID().FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ComputePool),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.ComputePoolResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasAutoResumeString(r.BooleanDefault).
						HasAutoSuspendSecsString(r.IntDefaultString).
						HasCommentString("").
						HasNoForApplication().
						HasInitiallySuspendedString(r.BooleanDefault).
						HasInstanceFamilyString(string(sdk.ComputePoolInstanceFamilyCpuX64S)).
						HasMaxNodesString("2").
						HasMinNodesString("1").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ComputePoolShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasState(sdk.ComputePoolStateStarting).
						HasMinNodes(1).
						HasMaxNodes(2).
						HasInstanceFamily(sdk.ComputePoolInstanceFamilyCpuX64S).
						HasNumServices(0).
						HasNumJobs(0).
						HasAutoSuspendSecs(3600).
						HasAutoResume(true).
						HasActiveNodes(0).
						HasIdleNodes(0).
						HasTargetNodes(1).
						HasCreatedOnNotEmpty().
						HasResumedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("").
						HasIsExclusive(false).
						HasApplicationEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.state", string(sdk.ComputePoolStateStarting))),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.min_nodes", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.max_nodes", "2")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.instance_family", string(sdk.ComputePoolInstanceFamilyCpuX64S))),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.num_services", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.num_jobs", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.auto_suspend_secs", "3600")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.auto_resume", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.active_nodes", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.idle_nodes", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.target_nodes", "1")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.resumed_on")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_exclusive", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.application", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.error_code", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.status_message", "Compute pool is starting for last 0 minutes")),
				),
			},
			// import - without optionals
			{
				Config:       accconfig.FromModels(t, modelBasic),
				ResourceName: modelBasic.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedComputePoolResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasAutoResumeString("true").
						HasAutoSuspendSecsString("3600").
						HasCommentString("").
						HasNoForApplication().
						HasNoInitiallySuspended().
						HasInstanceFamilyString(string(sdk.ComputePoolInstanceFamilyCpuX64S)).
						HasMaxNodesString("2").
						HasMinNodesString("1").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ImportedComputePoolShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasState(sdk.ComputePoolStateStarting).
						HasMinNodes(1).
						HasMaxNodes(2).
						HasInstanceFamily(sdk.ComputePoolInstanceFamilyCpuX64S).
						HasNumServices(0).
						HasNumJobs(0).
						HasAutoSuspendSecs(3600).
						HasAutoResume(true).
						HasActiveNodes(0).
						HasIdleNodes(0).
						HasTargetNodes(1).
						HasCreatedOnNotEmpty().
						HasResumedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("").
						HasIsExclusive(false).
						HasApplicationEmpty(),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.name", id.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.state", string(sdk.ComputePoolStateStarting))),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.min_nodes", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.max_nodes", "2")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.instance_family", string(sdk.ComputePoolInstanceFamilyCpuX64S))),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.num_services", "0")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.num_jobs", "0")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.auto_suspend_secs", "3600")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.auto_resume", "true")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.active_nodes", "0")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.idle_nodes", "0")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.target_nodes", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateSet(helpers.EncodeResourceIdentifier(id), "describe_output.0.created_on")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateSet(helpers.EncodeResourceIdentifier(id), "describe_output.0.resumed_on")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateSet(helpers.EncodeResourceIdentifier(id), "describe_output.0.updated_on")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.comment", "")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.is_exclusive", "false")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.application", "")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.error_code", "")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.status_message", "Compute pool is starting for last 0 minutes")),
				),
			},
			// set optionals
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.ComputePoolResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasAutoResumeString("true").
						HasAutoSuspendSecsString("6767").
						HasCommentString(comment).
						HasForApplicationString(application.ID().FullyQualifiedName()).
						HasInitiallySuspendedString("true").
						HasInstanceFamilyString(string(sdk.ComputePoolInstanceFamilyCpuX64S)).
						HasMaxNodesString("2").
						HasMinNodesString("1").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ComputePoolShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasState(sdk.ComputePoolStateSuspended).
						HasMinNodes(1).
						HasMaxNodes(2).
						HasInstanceFamily(sdk.ComputePoolInstanceFamilyCpuX64S).
						HasNumServices(0).
						HasNumJobs(0).
						HasAutoSuspendSecs(6767).
						HasAutoResume(true).
						HasActiveNodes(0).
						HasIdleNodes(0).
						HasTargetNodes(0).
						HasCreatedOnNotEmpty().
						HasResumedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasIsExclusive(true).
						HasApplication(application.ID()),
				),
			},
			// import - complete
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initially_suspended"},
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
					resourceassert.ComputePoolResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasAutoResumeString("true").
						HasAutoSuspendSecsString("2222").
						HasCommentString(changedComment).
						HasForApplicationString(application.ID().FullyQualifiedName()).
						HasInitiallySuspendedString("true").
						HasInstanceFamilyString(string(sdk.ComputePoolInstanceFamilyCpuX64S)).
						HasMaxNodesString("4").
						HasMinNodesString("3").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ComputePoolShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasState(sdk.ComputePoolStateSuspended).
						HasMinNodes(3).
						HasMaxNodes(4).
						HasInstanceFamily(sdk.ComputePoolInstanceFamilyCpuX64S).
						HasNumServices(0).
						HasNumJobs(0).
						HasAutoSuspendSecs(2222).
						HasAutoResume(true).
						HasActiveNodes(0).
						HasIdleNodes(0).
						HasTargetNodes(0).
						HasCreatedOnNotEmpty().
						HasResumedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(changedComment).
						HasIsExclusive(true).
						HasApplication(application.ID()),
				),
			},
			// change externally
			{
				PreConfig: func() {
					testClient().ComputePool.Alter(t, sdk.NewAlterComputePoolRequest(id).WithSet(
						*sdk.NewComputePoolSetRequest().
							WithMinNodes(4).
							WithMaxNodes(5).
							WithAutoResume(true).
							WithAutoSuspendSecs(3600).
							WithComment(comment),
					))
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				Check: assertThat(t,
					resourceassert.ComputePoolResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasAutoResumeString("true").
						HasAutoSuspendSecsString("2222").
						HasCommentString(changedComment).
						HasForApplicationString(application.ID().FullyQualifiedName()).
						HasInitiallySuspendedString("true").
						HasInstanceFamilyString(string(sdk.ComputePoolInstanceFamilyCpuX64S)).
						HasMaxNodesString("4").
						HasMinNodesString("3").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ComputePoolShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasName(id.Name()).
						HasState(sdk.ComputePoolStateSuspended).
						HasMinNodes(3).
						HasMaxNodes(4).
						HasInstanceFamily(sdk.ComputePoolInstanceFamilyCpuX64S).
						HasNumServices(0).
						HasNumJobs(0).
						HasAutoSuspendSecs(2222).
						HasAutoResume(true).
						HasActiveNodes(0).
						HasIdleNodes(0).
						HasTargetNodes(0).
						HasCreatedOnNotEmpty().
						HasResumedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(changedComment).
						HasIsExclusive(true).
						HasApplication(application.ID()),
				),
			},
			// ignore_after_creation does not cause plans
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				Check: assertThat(t,
					resourceassert.ComputePoolResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ComputePoolShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasComment(changedComment),
				),
			},
			// unset
			{
				Config: accconfig.FromModels(t, modelBasicWithApp),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasicWithApp.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ComputePoolResource(t, modelBasicWithApp.ResourceReference()).
						HasNameString(id.Name()).
						HasAutoResumeString(r.BooleanDefault).
						HasAutoSuspendSecsString(r.IntDefaultString).
						HasCommentString("").
						HasForApplicationString(application.ID().FullyQualifiedName()).
						HasInitiallySuspendedString(r.BooleanTrue).
						HasInstanceFamilyString(string(sdk.ComputePoolInstanceFamilyCpuX64S)).
						HasMaxNodesString("4").
						HasMinNodesString("3").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ComputePoolShowOutput(t, modelBasicWithApp.ResourceReference()).
						HasName(id.Name()).
						HasState(sdk.ComputePoolStateSuspended).
						HasMinNodes(3).
						HasMaxNodes(4).
						HasInstanceFamily(sdk.ComputePoolInstanceFamilyCpuX64S).
						HasNumServices(0).
						HasNumJobs(0).
						HasAutoSuspendSecs(3600).
						HasAutoResume(true).
						HasActiveNodes(0).
						HasIdleNodes(0).
						HasTargetNodes(0).
						HasCreatedOnNotEmpty().
						HasResumedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("").
						HasIsExclusive(true).
						HasApplication(application.ID()),
				),
			},
			// forcenew - instance family
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentInstanceFamily.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentInstanceFamily),
				Check: assertThat(t,
					resourceassert.ComputePoolResource(t, modelCompleteWithDifferentInstanceFamily.ResourceReference()).
						HasNameString(id.Name()).
						HasAutoResumeString(r.BooleanDefault).
						HasAutoSuspendSecsString(r.IntDefaultString).
						HasCommentString("").
						HasForApplicationString(application.ID().FullyQualifiedName()).
						HasInitiallySuspendedString(r.BooleanTrue).
						HasInstanceFamilyString(string(sdk.ComputePoolInstanceFamilyCpuX64M)).
						HasMaxNodesString("4").
						HasMinNodesString("3").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ComputePoolShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasState(sdk.ComputePoolStateSuspended).
						HasMinNodes(3).
						HasMaxNodes(4).
						HasInstanceFamily(sdk.ComputePoolInstanceFamilyCpuX64M).
						HasNumServices(0).
						HasNumJobs(0).
						HasAutoSuspendSecs(3600).
						HasAutoResume(true).
						HasActiveNodes(0).
						HasIdleNodes(0).
						HasTargetNodes(0).
						HasCreatedOnNotEmpty().
						HasResumedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("").
						HasIsExclusive(true).
						HasApplication(application.ID()),
				),
			},
		},
	})
}

func TestAcc_ComputePool_complete(t *testing.T) {
	application := createApp(t)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	modelComplete := model.ComputePool("test", id.Name(), string(sdk.ComputePoolInstanceFamilyCpuX64S), 2, 1).
		WithForApplication(application.ID().FullyQualifiedName()).
		WithAutoResume("true").
		WithInitiallySuspended("true").
		WithAutoSuspendSecs(6767).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ComputePool),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.ComputePoolResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasAutoResumeString("true").
						HasAutoSuspendSecsString("6767").
						HasCommentString(comment).
						HasForApplicationString(application.ID().FullyQualifiedName()).
						HasInitiallySuspendedString("true").
						HasInstanceFamilyString(string(sdk.ComputePoolInstanceFamilyCpuX64S)).
						HasMaxNodesString("2").
						HasMinNodesString("1").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ComputePoolShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasState(sdk.ComputePoolStateSuspended).
						HasMinNodes(1).
						HasMaxNodes(2).
						HasInstanceFamily(sdk.ComputePoolInstanceFamilyCpuX64S).
						HasNumServices(0).
						HasNumJobs(0).
						HasAutoSuspendSecs(6767).
						HasAutoResume(true).
						HasActiveNodes(0).
						HasIdleNodes(0).
						HasTargetNodes(0).
						HasCreatedOnNotEmpty().
						HasResumedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasIsExclusive(true).
						HasApplication(application.ID()),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.state", string(sdk.ComputePoolStateSuspended))),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.min_nodes", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.max_nodes", "2")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.instance_family", string(sdk.ComputePoolInstanceFamilyCpuX64S))),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.num_services", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.num_jobs", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.auto_suspend_secs", "6767")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.auto_resume", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.active_nodes", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.idle_nodes", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.target_nodes", "0")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.resumed_on")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_exclusive", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.application", application.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.error_code", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.status_message", "")),
				),
			},
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initially_suspended"},
			},
		},
	})
}

func TestAcc_ComputePool_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	modelInvalidMinNodes := model.ComputePool("test", id.Name(), string(sdk.ComputePoolInstanceFamilyCpuX64S), 2, 0)
	modelInvalidMaxNodes := model.ComputePool("test", id.Name(), string(sdk.ComputePoolInstanceFamilyCpuX64S), 0, 2)
	modelInvalidInstanceFamily := model.ComputePool("test", id.Name(), "invalid", 2, 1)
	modelInvalidAutoResume := model.ComputePool("test", id.Name(), string(sdk.ComputePoolInstanceFamilyCpuX64S), 1, 1).
		WithAutoResume("invalid")
	modelInvalidInitiallySuspended := model.ComputePool("test", id.Name(), string(sdk.ComputePoolInstanceFamilyCpuX64S), 1, 1).
		WithInitiallySuspended("invalid")
	modelInvalidAutoSuspendSecs := model.ComputePool("test", id.Name(), string(sdk.ComputePoolInstanceFamilyCpuX64S), 1, 1).
		WithAutoSuspendSecs(-1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ComputePool),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, modelInvalidMinNodes),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected min_nodes to be at least \(1\), got 0`),
			},
			{
				Config:      config.FromModels(t, modelInvalidMaxNodes),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected max_nodes to be at least \(1\), got 0`),
			},
			{
				Config:      config.FromModels(t, modelInvalidInstanceFamily),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid compute pool instance family: INVALID`),
			},
			{
				Config:      config.FromModels(t, modelInvalidAutoResume),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected \[\{\{} auto_resume}] to be one of \["true" "false"], got invalid`),
			},
			{
				Config:      config.FromModels(t, modelInvalidInitiallySuspended),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected \[\{\{} initially_suspended}] to be one of \["true" "false"], got invalid`),
			},
			{
				Config:      config.FromModels(t, modelInvalidAutoSuspendSecs),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected auto_suspend_secs to be at least \(0\), got -1`),
			},
		},
	})
}
