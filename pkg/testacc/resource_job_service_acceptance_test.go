//go:build !account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
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

func TestAcc_JobService_basic_fromSpecification(t *testing.T) {
	computePool, computePoolCleanup := testClient().ComputePool.Create(t)
	t.Cleanup(computePoolCleanup)

	networkRule, networkRuleCleanup := testClient().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)

	externalAccessIntegration1Id, externalAccessIntegration1Cleanup := testClient().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
	t.Cleanup(externalAccessIntegration1Cleanup)

	externalAccessIntegration2Id, externalAccessIntegration2Cleanup := testClient().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
	t.Cleanup(externalAccessIntegration2Cleanup)

	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)

	comment, changedComment := random.Comment(), random.Comment()
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	spec := testClient().Service.SampleSpec(t)

	modelBasic := model.JobServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), spec)

	// TODO(SNOW-2138932): Test without async option. This probably requires a custom no-op image in the image registry.
	modelComplete := model.JobServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), spec).
		WithExternalAccessIntegrations(externalAccessIntegration1Id).
		WithQueryWarehouse(testClient().Ids.WarehouseId().FullyQualifiedName()).
		WithComment(comment)

	modelCompleteWithDifferentValues := model.JobServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), spec).
		WithExternalAccessIntegrations(externalAccessIntegration2Id).
		WithQueryWarehouse(warehouse.ID().FullyQualifiedName()).
		WithComment(changedComment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.JobService),
		Steps: []resource.TestStep{
			// create without optionals
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.JobServiceResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTextNotEmpty().
						HasExternalAccessIntegrationsEmpty().
						HasNoQueryWarehouse().
						HasCommentString(""),
					resourceshowoutputassert.ServiceShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
						HasCurrentInstances(0).
						HasTargetInstances(1).
						HasMinReadyInstances(1).
						HasMinInstances(1).
						HasMaxInstances(1).
						HasAutoResume(true).
						HasNoExternalAccessIntegrations().
						HasCreatedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasResumedOnEmpty().
						HasSuspendedOnEmpty().
						HasAutoSuspendSecs(0).
						HasComment("").
						HasOwnerRoleType("ROLE").
						HasQueryWarehouseEmpty().
						HasIsJob(true).
						HasIsAsyncJob(true).
						HasSpecDigestNotEmpty().
						HasIsUpgrading(false).
						HasManagingObjectDomainEmpty().
						HasManagingObjectNameEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.current_instances", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.target_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.min_ready_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.min_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.max_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.auto_resume", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.external_access_integrations.#", "0")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.resumed_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.suspended_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.auto_suspend_secs", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.query_warehouse", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_job", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_async_job", "true")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.spec_digest")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_upgrading", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.managing_object_domain", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.managing_object_name", "")),
				),
			},
			// import without optionals
			{
				Config:       accconfig.FromModels(t, modelBasic),
				ResourceName: modelBasic.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedJobServiceResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTemplateEmpty().
						HasFromSpecificationEmpty().
						HasExternalAccessIntegrationsEmpty().
						HasNoQueryWarehouse().
						HasCommentString(""),
					resourceshowoutputassert.ImportedServiceShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
						HasCurrentInstances(0).
						HasTargetInstances(1).
						HasMinReadyInstances(1).
						HasMinInstances(1).
						HasMaxInstances(1).
						HasAutoResume(true).
						HasNoExternalAccessIntegrations().
						HasCreatedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasResumedOnEmpty().
						HasSuspendedOnEmpty().
						HasAutoSuspendSecs(0).
						HasComment("").
						HasOwnerRoleType("ROLE").
						HasQueryWarehouseEmpty().
						HasIsJob(true).
						HasIsAsyncJob(true).
						HasSpecDigestNotEmpty().
						HasIsUpgrading(false).
						HasManagingObjectDomainEmpty().
						HasManagingObjectNameEmpty(),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.name", id.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.database_name", id.DatabaseName())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.schema_name", id.SchemaName())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateSet(helpers.EncodeResourceIdentifier(id), "describe_output.0.spec")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateSet(helpers.EncodeResourceIdentifier(id), "describe_output.0.dns_name")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.current_instances", "0")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.target_instances", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.min_ready_instances", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.min_instances", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.max_instances", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.auto_resume", "true")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.external_access_integrations.#", "0")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateSet(helpers.EncodeResourceIdentifier(id), "describe_output.0.created_on")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateSet(helpers.EncodeResourceIdentifier(id), "describe_output.0.updated_on")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.resumed_on", "")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.suspended_on", "")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.auto_suspend_secs", "0")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.comment", "")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.owner_role_type", "ROLE")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.query_warehouse", "")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.is_job", "true")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.is_async_job", "true")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateSet(helpers.EncodeResourceIdentifier(id), "describe_output.0.spec_digest")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.is_upgrading", "false")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.managing_object_domain", "")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.managing_object_name", ""))),
			},
			// set optionals
			{
				Config: accconfig.FromModels(t, modelComplete),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.JobServiceResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTextNotEmpty().
						HasExternalAccessIntegrations(externalAccessIntegration1Id).
						HasQueryWarehouseString(testClient().Ids.WarehouseId().FullyQualifiedName()).
						HasCommentString(comment),
					resourceshowoutputassert.ServiceShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
						HasCurrentInstances(0).
						HasTargetInstances(1).
						HasMinReadyInstances(1).
						HasMinInstances(1).
						HasMaxInstances(1).
						HasAutoResume(true).
						HasExternalAccessIntegrations(externalAccessIntegration1Id).
						HasCreatedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasResumedOnEmpty().
						HasSuspendedOnEmpty().
						HasAutoSuspendSecs(0).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasQueryWarehouse(testClient().Ids.WarehouseId()).
						HasIsJob(true).
						HasIsAsyncJob(true).
						HasSpecDigestNotEmpty().
						HasIsUpgrading(false).
						HasManagingObjectDomainEmpty().
						HasManagingObjectNameEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.current_instances", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.target_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.min_ready_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.min_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.max_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.auto_resume", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.external_access_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.external_access_integrations.0", externalAccessIntegration1Id.Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.resumed_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.suspended_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.auto_suspend_secs", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.query_warehouse", testClient().Ids.WarehouseId().Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_job", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_async_job", "true")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.spec_digest")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_upgrading", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.managing_object_domain", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.managing_object_name", "")),
				),
			},
			// import complete object
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"from_specification"},
			},
			// alter
			{
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.JobServiceResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTextNotEmpty().
						HasExternalAccessIntegrations(externalAccessIntegration2Id).
						HasQueryWarehouseString(warehouse.ID().FullyQualifiedName()).
						HasCommentString(changedComment),
					resourceshowoutputassert.ServiceShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
						HasCurrentInstances(0).
						HasTargetInstances(1).
						HasMinReadyInstances(1).
						HasMinInstances(1).
						HasMaxInstances(1).
						HasAutoResume(true).
						HasExternalAccessIntegrations(externalAccessIntegration2Id).
						HasCreatedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasResumedOnEmpty().
						HasSuspendedOnEmpty().
						HasAutoSuspendSecs(0).
						HasComment(changedComment).
						HasOwnerRoleType("ROLE").
						HasQueryWarehouse(warehouse.ID()).
						HasIsJob(true).
						HasIsAsyncJob(true).
						HasSpecDigestNotEmpty().
						HasIsUpgrading(false).
						HasManagingObjectDomainEmpty().
						HasManagingObjectNameEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.current_instances", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.target_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.min_ready_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.min_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.max_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.auto_resume", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.external_access_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.external_access_integrations.0", externalAccessIntegration2Id.Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.resumed_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.suspended_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.auto_suspend_secs", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.comment", changedComment)),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.query_warehouse", warehouse.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_job", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_async_job", "true")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.spec_digest")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_upgrading", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.managing_object_domain", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.managing_object_name", "")),
				),
			},
			// change externally
			{
				PreConfig: func() {
					testClient().Service.DropFunc(t, id)()
					_, serviceCleanup := testClient().Service.ExecuteJobService(t, computePool.ID(), id)
					t.Cleanup(serviceCleanup)
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.JobServiceResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTextNotEmpty().
						HasExternalAccessIntegrations(externalAccessIntegration2Id).
						HasQueryWarehouseString(warehouse.ID().FullyQualifiedName()).
						HasCommentString(changedComment),
					resourceshowoutputassert.ServiceShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
						HasCurrentInstances(0).
						HasTargetInstances(1).
						HasMinReadyInstances(1).
						HasMinInstances(1).
						HasMaxInstances(1).
						HasAutoResume(true).
						HasExternalAccessIntegrations(externalAccessIntegration2Id).
						HasCreatedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasResumedOnEmpty().
						HasSuspendedOnEmpty().
						HasAutoSuspendSecs(0).
						HasComment(changedComment).
						HasOwnerRoleType("ROLE").
						HasQueryWarehouse(warehouse.ID()).
						HasIsJob(true).
						HasIsAsyncJob(true).
						HasSpecDigestNotEmpty().
						HasIsUpgrading(false).
						HasManagingObjectDomainEmpty().
						HasManagingObjectNameEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.current_instances", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.target_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.min_ready_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.min_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.max_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.auto_resume", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.external_access_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.external_access_integrations.0", externalAccessIntegration2Id.Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.resumed_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.suspended_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.auto_suspend_secs", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.comment", changedComment)),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.query_warehouse", warehouse.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_job", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_async_job", "true")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.spec_digest")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_upgrading", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.managing_object_domain", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.managing_object_name", "")),
				),
			},
			// unset
			{
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.JobServiceResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTextNotEmpty().
						HasExternalAccessIntegrationsEmpty().
						HasNoQueryWarehouse().
						HasCommentString(""),
					resourceshowoutputassert.ServiceShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
						HasCurrentInstances(0).
						HasTargetInstances(1).
						HasMinReadyInstances(1).
						HasMinInstances(1).
						HasMaxInstances(1).
						HasAutoResume(true).
						HasNoExternalAccessIntegrations().
						HasCreatedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasResumedOnEmpty().
						HasSuspendedOnEmpty().
						HasAutoSuspendSecs(0).
						HasComment("").
						HasOwnerRoleType("ROLE").
						HasQueryWarehouseEmpty().
						HasIsJob(true).
						HasIsAsyncJob(true).
						HasSpecDigestNotEmpty().
						HasIsUpgrading(false).
						HasManagingObjectDomainEmpty().
						HasManagingObjectNameEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.current_instances", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.target_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.min_ready_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.min_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.max_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.auto_resume", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.external_access_integrations.#", "0")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.resumed_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.suspended_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.auto_suspend_secs", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.query_warehouse", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_job", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_async_job", "true")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.spec_digest")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_upgrading", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.managing_object_domain", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.managing_object_name", "")),
				),
			},
		},
	})
}

// TODO (next PR): add tests for detecting is_job

func TestAcc_JobService_fromSpecificationOnStage(t *testing.T) {
	computePool, computePoolCleanup := testClient().ComputePool.Create(t)
	t.Cleanup(computePoolCleanup)

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	spec := `
spec:
  containers:
  - name: example-container
    image: /snowflake/images/snowflake_images/exampleimage:latest
`
	specFileName := "spec.yaml"
	testClient().Stage.PutInLocationWithContent(t, stage.Location(), specFileName, spec)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelBasic := model.JobServiceWithSpecOnStage("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), stage.ID(), specFileName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.JobService),
		Steps: []resource.TestStep{
			// create without optionals
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.JobServiceResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationOnStageNotEmpty().
						HasExternalAccessIntegrationsEmpty().
						HasNoQueryWarehouse().
						HasCommentString(""),
					resourceshowoutputassert.ServiceShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
						HasCurrentInstances(0).
						HasTargetInstances(1).
						HasMinReadyInstances(1).
						HasMinInstances(1).
						HasMaxInstances(1).
						HasAutoResume(true).
						HasNoExternalAccessIntegrations().
						HasCreatedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasResumedOnEmpty().
						HasSuspendedOnEmpty().
						HasAutoSuspendSecs(0).
						HasComment("").
						HasOwnerRoleType("ROLE").
						HasQueryWarehouseEmpty().
						HasIsJob(true).
						HasIsAsyncJob(true).
						HasSpecDigestNotEmpty().
						HasIsUpgrading(false).
						HasManagingObjectDomainEmpty().
						HasManagingObjectNameEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.current_instances", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.target_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.min_ready_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.min_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.max_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.auto_resume", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.external_access_integrations.#", "0")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.resumed_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.suspended_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.auto_suspend_secs", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.query_warehouse", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_job", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_async_job", "true")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.spec_digest")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_upgrading", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.managing_object_domain", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.managing_object_name", "")),
				),
			},
		},
	})
}

// TODO (next PR): fill
// func TestAcc_Service_fromSpecificationTemplate(t *testing.T) {
// }

// TODO (next PR): fill
// func TestAcc_Service_fromSpecificationTemplateOnStage(t *testing.T) {
// }

func TestAcc_JobService_changingSpec(t *testing.T) {
	computePool, computePoolCleanup := testClient().ComputePool.Create(t)
	t.Cleanup(computePoolCleanup)

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	spec := testClient().Service.SampleSpec(t)
	specFileName := "spec.yaml"
	testClient().Stage.PutInLocationWithContent(t, stage.Location(), specFileName, spec)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelBasicOnStage := model.JobServiceWithSpecOnStage("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), stage.ID(), specFileName)
	modelBasic := model.JobServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), spec)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.JobService),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.JobServiceResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTextNotEmpty(),
					resourceshowoutputassert.ServiceShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.name", id.Name())),
				),
			},
			// update
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasicOnStage.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, modelBasicOnStage),
				Check: assertThat(t,
					resourceassert.JobServiceResource(t, modelBasicOnStage.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationOnStageNotEmpty(),
					resourceshowoutputassert.ServiceShowOutput(t, modelBasicOnStage.ResourceReference()).
						HasName(id.Name()),
					assert.Check(resource.TestCheckResourceAttr(modelBasicOnStage.ResourceReference(), "describe_output.0.name", id.Name())),
				),
			},
			// external changed are not detected
			{
				PreConfig: func() {
					testClient().Service.DropFunc(t, id)()
					_, serviceCleanup := testClient().Service.ExecuteJobService(t, computePool.ID(), id)
					t.Cleanup(serviceCleanup)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasicOnStage.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: accconfig.FromModels(t, modelBasicOnStage),
				Check: assertThat(t,
					resourceassert.JobServiceResource(t, modelBasicOnStage.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationOnStageNotEmpty(),
					resourceshowoutputassert.ServiceShowOutput(t, modelBasicOnStage.ResourceReference()).
						HasName(id.Name()),
					assert.Check(resource.TestCheckResourceAttr(modelBasicOnStage.ResourceReference(), "describe_output.0.name", id.Name())),
				),
			},
		},
	})
}

func TestAcc_JobService_complete(t *testing.T) {
	computePool, computePoolCleanup := testClient().ComputePool.Create(t)
	t.Cleanup(computePoolCleanup)

	networkRule, networkRuleCleanup := testClient().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)

	externalAccessIntegrationId, externalAccessIntegrationCleanup := testClient().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
	t.Cleanup(externalAccessIntegrationCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	spec := testClient().Service.SampleSpec(t)

	modelComplete := model.JobServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), spec).
		WithExternalAccessIntegrations(externalAccessIntegrationId).
		WithQueryWarehouse(testClient().Ids.WarehouseId().FullyQualifiedName()).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.JobService),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.JobServiceResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTextNotEmpty().
						HasExternalAccessIntegrations(externalAccessIntegrationId).
						HasQueryWarehouseString(testClient().Ids.WarehouseId().FullyQualifiedName()).
						HasCommentString(comment),
					resourceshowoutputassert.ServiceShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
						HasCurrentInstances(0).
						HasTargetInstances(1).
						HasMinReadyInstances(1).
						HasMinInstances(1).
						HasMaxInstances(1).
						HasAutoResume(true).
						HasExternalAccessIntegrations(externalAccessIntegrationId).
						HasCreatedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasResumedOnEmpty().
						HasSuspendedOnEmpty().
						HasAutoSuspendSecs(0).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasQueryWarehouse(testClient().Ids.WarehouseId()).
						HasIsJob(true).
						HasIsAsyncJob(true).
						HasSpecDigestNotEmpty().
						HasIsUpgrading(false).
						HasManagingObjectDomainEmpty().
						HasManagingObjectNameEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.current_instances", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.target_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.min_ready_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.min_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.max_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.auto_resume", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.external_access_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.external_access_integrations.0", externalAccessIntegrationId.Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.resumed_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.suspended_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.auto_suspend_secs", "0")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.query_warehouse", testClient().Ids.WarehouseId().Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_job", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_async_job", "true")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.spec_digest")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_upgrading", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.managing_object_domain", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.managing_object_name", "")),
				),
			},
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"from_specification"},
			},
		},
	})
}

// TODO (next PR): Implement validations and add tests for them.
// func TestAcc_JobService_Validations(t *testing.T) {
// }
