//go:build !account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Service_basic_fromSpecification(t *testing.T) {
	computePool, computePoolCleanup := testClient().ComputePool.Create(t)
	t.Cleanup(computePoolCleanup)

	networkRule, networkRuleCleanup := testClient().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)

	externalAccessIntegrationId, externalAccessIntegration1Cleanup := testClient().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
	t.Cleanup(externalAccessIntegration1Cleanup)

	comment, changedComment := random.Comment(), random.Comment()
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	spec := testClient().Service.SampleSpec(t)

	modelBasic := model.ServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), spec)

	modelComplete := model.ServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), spec).
		WithAutoSuspendSecs(6767).
		WithExternalAccessIntegrations(externalAccessIntegrationId).
		WithAutoResume("true").
		WithMinInstances(2).
		WithMinReadyInstances(2).
		WithMaxInstances(2).
		WithQueryWarehouse(testClient().Ids.WarehouseId().FullyQualifiedName()).
		WithComment(comment)

	modelCompleteWithDifferentValues := model.ServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), spec).
		WithAutoSuspendSecs(2222).
		WithAutoResume("false").
		WithMinInstances(1).
		WithMinReadyInstances(1).
		WithMaxInstances(1).
		WithQueryWarehouse(testClient().Ids.WarehouseId().FullyQualifiedName()).
		WithComment(changedComment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Service),
		Steps: []resource.TestStep{
			// create without optionals
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.ServiceResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTextNotEmpty().
						HasAutoSuspendSecsString(r.IntDefaultString).
						HasExternalAccessIntegrationsEmpty().
						HasAutoResumeString(r.BooleanDefault).
						HasNoMinInstances().
						HasNoMinReadyInstances().
						HasNoMaxInstances().
						HasNoQueryWarehouse().
						HasServiceTypeString(string(sdk.ServiceTypeService)).
						HasCommentString(""),
					resourceshowoutputassert.ServiceShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
						HasCurrentInstances(1).
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
						HasIsJob(false).
						HasIsAsyncJob(false).
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
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.current_instances", "1")),
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
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_job", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_async_job", "false")),
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
					resourceassert.ImportedServiceResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTemplateEmpty().
						HasFromSpecificationEmpty().
						HasAutoSuspendSecsString("0").
						HasExternalAccessIntegrationsEmpty().
						HasAutoResumeString("true").
						HasMinInstancesString("1").
						HasMinReadyInstancesString("1").
						HasMaxInstancesString("1").
						HasNoQueryWarehouse().
						HasServiceTypeString(string(sdk.ServiceTypeService)).
						HasCommentString(""),
					resourceshowoutputassert.ImportedServiceShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
						HasCurrentInstances(1).
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
						HasIsJob(false).
						HasIsAsyncJob(false).
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
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.current_instances", "1")),
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
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.is_job", "false")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.is_async_job", "false")),
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
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ServiceResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTextNotEmpty().
						HasAutoSuspendSecsString("6767").
						HasExternalAccessIntegrations(externalAccessIntegrationId).
						HasAutoResumeString(r.BooleanTrue).
						HasMinInstancesString("2").
						HasMinReadyInstancesString("2").
						HasMaxInstancesString("2").
						HasQueryWarehouseString(testClient().Ids.WarehouseId().FullyQualifiedName()).
						HasServiceTypeString(string(sdk.ServiceTypeService)).
						HasCommentString(comment),
					resourceshowoutputassert.ServiceShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
						HasCurrentInstances(1).
						HasTargetInstances(2).
						HasMinReadyInstances(2).
						HasMinInstances(2).
						HasMaxInstances(2).
						HasAutoResume(true).
						HasExternalAccessIntegrations(externalAccessIntegrationId).
						HasCreatedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasResumedOnEmpty().
						HasSuspendedOnEmpty().
						HasAutoSuspendSecs(6767).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasQueryWarehouse(testClient().Ids.WarehouseId()).
						HasIsJob(false).
						HasIsAsyncJob(false).
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
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.current_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.target_instances", "2")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.min_ready_instances", "2")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.min_instances", "2")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.max_instances", "2")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.auto_resume", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.external_access_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.external_access_integrations.0", externalAccessIntegrationId.Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.resumed_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.suspended_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.auto_suspend_secs", "6767")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.query_warehouse", testClient().Ids.WarehouseId().Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_job", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_async_job", "false")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.spec_digest")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_upgrading", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.managing_object_domain", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.managing_object_name", "")),
				),
			},
			// import complete object
			{
				Config:       accconfig.FromModels(t, modelComplete),
				ResourceName: modelComplete.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedServiceResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationEmpty().
						HasAutoSuspendSecsString("6767").
						HasExternalAccessIntegrations(externalAccessIntegrationId).
						HasAutoResumeString("true").
						HasMinInstancesString("2").
						HasMinReadyInstancesString("2").
						HasMaxInstancesString("2").
						HasQueryWarehouseString(testClient().Ids.WarehouseId().FullyQualifiedName()).
						HasServiceTypeString(string(sdk.ServiceTypeService)).
						HasCommentString(comment),
					resourceshowoutputassert.ImportedServiceShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
						HasCurrentInstances(1).
						HasTargetInstances(2).
						HasMinReadyInstances(2).
						HasMinInstances(2).
						HasMaxInstances(2).
						HasAutoResume(true).
						HasExternalAccessIntegrations(externalAccessIntegrationId).
						HasCreatedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasResumedOnEmpty().
						HasSuspendedOnEmpty().
						HasAutoSuspendSecs(6767).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasQueryWarehouse(testClient().Ids.WarehouseId()).
						HasIsJob(false).
						HasIsAsyncJob(false).
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
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.current_instances", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.target_instances", "2")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.min_ready_instances", "2")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.min_instances", "2")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.max_instances", "2")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.auto_resume", "true")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.external_access_integrations.#", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.external_access_integrations.0", externalAccessIntegrationId.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateSet(helpers.EncodeResourceIdentifier(id), "describe_output.0.created_on")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateSet(helpers.EncodeResourceIdentifier(id), "describe_output.0.updated_on")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.resumed_on", "")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.suspended_on", "")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.auto_suspend_secs", "6767")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.comment", comment)),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.owner_role_type", "ROLE")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.query_warehouse", testClient().Ids.WarehouseId().Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.is_job", "false")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.is_async_job", "false")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateSet(helpers.EncodeResourceIdentifier(id), "describe_output.0.spec_digest")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.is_upgrading", "false")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.managing_object_domain", "")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.managing_object_name", ""))),
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
					resourceassert.ServiceResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTextNotEmpty().
						HasAutoSuspendSecsString("2222").
						HasExternalAccessIntegrationsEmpty().
						HasAutoResumeString(r.BooleanFalse).
						HasMinInstancesString("1").
						HasMinReadyInstancesString("1").
						HasMaxInstancesString("1").
						HasQueryWarehouseString(testClient().Ids.WarehouseId().FullyQualifiedName()).
						HasServiceTypeString(string(sdk.ServiceTypeService)).
						HasCommentString(changedComment),
					resourceshowoutputassert.ServiceShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
						HasCurrentInstances(1).
						HasTargetInstances(1).
						HasMinReadyInstances(1).
						HasMinInstances(1).
						HasMaxInstances(1).
						HasAutoResume(false).
						HasNoExternalAccessIntegrations().
						HasCreatedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasResumedOnEmpty().
						HasSuspendedOnEmpty().
						HasAutoSuspendSecs(2222).
						HasComment(changedComment).
						HasOwnerRoleType("ROLE").
						HasQueryWarehouse(testClient().Ids.WarehouseId()).
						HasIsJob(false).
						HasIsAsyncJob(false).
						HasSpecDigestNotEmpty().
						HasIsUpgrading(false).
						HasManagingObjectDomainEmpty().
						HasManagingObjectNameEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.current_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.target_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.min_ready_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.min_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.max_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.auto_resume", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.external_access_integrations.#", "0")),
					assert.Check(resource.TestCheckResourceAttrSet(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.resumed_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.suspended_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.auto_suspend_secs", "2222")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.comment", changedComment)),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.query_warehouse", testClient().Ids.WarehouseId().Name())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.is_job", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.is_async_job", "false")),
					assert.Check(resource.TestCheckResourceAttrSet(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.spec_digest")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.is_upgrading", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.managing_object_domain", "")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.managing_object_name", "")),
				),
			},
			// change externally
			{
				PreConfig: func() {
					testClient().Service.Alter(t, sdk.NewAlterServiceRequest(id).WithSet(
						*sdk.NewServiceSetRequest().
							WithMinReadyInstances(1).
							WithMinInstances(2).
							WithMaxInstances(3).
							WithAutoSuspendSecs(4242).
							WithAutoResume(true).
							WithQueryWarehouse(testClient().Ids.SnowflakeWarehouseId()).
							WithExternalAccessIntegrations(*sdk.NewServiceExternalAccessIntegrationsRequest([]sdk.AccountObjectIdentifier{externalAccessIntegrationId})).
							WithComment(random.Comment()),
					))
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ServiceResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTextNotEmpty().
						HasAutoSuspendSecsString("2222").
						HasExternalAccessIntegrationsEmpty().
						HasAutoResumeString(r.BooleanFalse).
						HasMinInstancesString("1").
						HasMinReadyInstancesString("1").
						HasMaxInstancesString("1").
						HasQueryWarehouseString(testClient().Ids.WarehouseId().FullyQualifiedName()).
						HasServiceTypeString(string(sdk.ServiceTypeService)).
						HasCommentString(changedComment),
					resourceshowoutputassert.ServiceShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
						// Current and target instances are skipped because the value is not consistent and depends on provisioning the compute pool instances.
						HasMinReadyInstances(1).
						HasMinInstances(1).
						HasMaxInstances(1).
						HasAutoResume(false).
						HasNoExternalAccessIntegrations().
						HasCreatedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasResumedOnEmpty().
						HasSuspendedOnEmpty().
						HasAutoSuspendSecs(2222).
						HasComment(changedComment).
						HasOwnerRoleType("ROLE").
						HasQueryWarehouse(testClient().Ids.WarehouseId()).
						HasIsJob(false).
						HasIsAsyncJob(false).
						HasSpecDigestNotEmpty().
						HasIsUpgrading(false).
						HasManagingObjectDomainEmpty().
						HasManagingObjectNameEmpty(),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.dns_name")),
					// Current and target instances are skipped because the value is not consistent and depends on provisioning the compute pool instances.
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.min_ready_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.min_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.max_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.auto_resume", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.external_access_integrations.#", "0")),
					assert.Check(resource.TestCheckResourceAttrSet(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.resumed_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.suspended_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.auto_suspend_secs", "2222")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.comment", changedComment)),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.query_warehouse", testClient().Ids.WarehouseId().Name())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.is_job", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.is_async_job", "false")),
					assert.Check(resource.TestCheckResourceAttrSet(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.spec_digest")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.is_upgrading", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.managing_object_domain", "")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.managing_object_name", "")),
				),
			},
			// unset
			{
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.ServiceResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTextNotEmpty().
						HasAutoSuspendSecsString(r.IntDefaultString).
						HasExternalAccessIntegrationsEmpty().
						HasAutoResumeString(r.BooleanDefault).
						HasMinInstancesString("0").
						HasMinReadyInstancesString("0").
						HasMaxInstancesString("0").
						HasQueryWarehouseString("").
						HasServiceTypeString(string(sdk.ServiceTypeService)).
						HasCommentString(""),
					resourceshowoutputassert.ServiceShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
						// Current and target instances are skipped because the value is not consistent and depends on provisioning the compute pool instances.
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
						HasIsJob(false).
						HasIsAsyncJob(false).
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
					// Current and target instances are skipped because the value is not consistent and depends on provisioning the compute pool instances.
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
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_job", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_async_job", "false")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.spec_digest")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_upgrading", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.managing_object_domain", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.managing_object_name", "")),
				),
			},
		},
	})
}

func TestAcc_Service_changingSpec(t *testing.T) {
	computePool, computePoolCleanup := testClient().ComputePool.Create(t)
	t.Cleanup(computePoolCleanup)

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	spec := testClient().Service.SampleSpec(t)
	specFileName := "spec.yaml"
	testClient().Stage.PutInLocationWithContent(t, stage.Location(), specFileName, spec)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelBasicOnStage := model.ServiceWithSpecOnStage("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), stage.ID(), specFileName)
	modelBasic := model.ServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), spec)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Service),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.ServiceResource(t, modelBasic.ResourceReference()).
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
						plancheck.ExpectResourceAction(modelBasicOnStage.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelBasicOnStage),
				Check: assertThat(t,
					resourceassert.ServiceResource(t, modelBasicOnStage.ResourceReference()).
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
					testClient().Service.Alter(t, sdk.NewAlterServiceRequest(id).WithFromSpecification(*sdk.NewServiceFromSpecificationRequest().WithSpecification(testClient().Service.SampleSpecWithContainerName(t, "external-changed"))))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasicOnStage.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: accconfig.FromModels(t, modelBasicOnStage),
				Check: assertThat(t,
					resourceassert.ServiceResource(t, modelBasicOnStage.ResourceReference()).
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

func TestAcc_Service_changeServiceTypeExternally(t *testing.T) {
	computePool, computePoolCleanup := testClient().ComputePool.Create(t)
	t.Cleanup(computePoolCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	spec := testClient().Service.SampleSpec(t)

	modelBasic := model.ServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), spec)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Service),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.ServiceResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasServiceTypeString(string(sdk.ServiceTypeService)),
					resourceshowoutputassert.ServiceShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasIsJob(false).
						HasIsAsyncJob(false),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_job", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_async_job", "false")),
				),
			},
			{
				PreConfig: func() {
					testClient().Service.DropFunc(t, id)()
					_, serviceCleanup := testClient().Service.ExecuteJobService(t, computePool.ID(), id)
					t.Cleanup(serviceCleanup)
				},
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t,
					resourceassert.ServiceResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasServiceTypeString(string(sdk.ServiceTypeService)),
					resourceshowoutputassert.ServiceShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasIsJob(false).
						HasIsAsyncJob(false),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_job", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_async_job", "false")),
				),
			},
		},
	})
}

func TestAcc_Service_fromSpecificationOnStage(t *testing.T) {
	computePool, computePoolCleanup := testClient().ComputePool.Create(t)
	t.Cleanup(computePoolCleanup)

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	spec := testClient().Service.SampleSpec(t)
	specFileName := "spec.yaml"
	testClient().Stage.PutInLocationWithContent(t, stage.Location(), specFileName, spec)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelBasic := model.ServiceWithSpecOnStage("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), stage.ID(), specFileName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Service),
		Steps: []resource.TestStep{
			// create without optionals
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.ServiceResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationOnStageNotEmpty().
						HasAutoSuspendSecsString(r.IntDefaultString).
						HasExternalAccessIntegrationsEmpty().
						HasAutoResumeString(r.BooleanDefault).
						HasNoMinInstances().
						HasNoMinReadyInstances().
						HasNoMaxInstances().
						HasNoQueryWarehouse().
						HasServiceTypeString(string(sdk.ServiceTypeService)).
						HasCommentString(""),
					resourceshowoutputassert.ServiceShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
						HasCurrentInstances(1).
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
						HasIsJob(false).
						HasIsAsyncJob(false).
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
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.current_instances", "1")),
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
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_job", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.is_async_job", "false")),
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

func TestAcc_Service_complete(t *testing.T) {
	computePool, computePoolCleanup := testClient().ComputePool.Create(t)
	t.Cleanup(computePoolCleanup)

	networkRule, networkRuleCleanup := testClient().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)

	externalAccessIntegrationId, externalAccessIntegrationCleanup := testClient().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
	t.Cleanup(externalAccessIntegrationCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	spec := testClient().Service.SampleSpec(t)

	modelComplete := model.ServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), spec).
		WithAutoSuspendSecs(6767).
		WithExternalAccessIntegrations(externalAccessIntegrationId).
		WithAutoResume("true").
		WithMinInstances(1).
		WithMinReadyInstances(1).
		WithMaxInstances(2).
		WithQueryWarehouse(testClient().Ids.WarehouseId().FullyQualifiedName()).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Service),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.ServiceResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTextNotEmpty().
						HasAutoSuspendSecsString("6767").
						HasExternalAccessIntegrations(externalAccessIntegrationId).
						HasAutoResumeString(r.BooleanTrue).
						HasMinInstancesString("1").
						HasMinReadyInstancesString("1").
						HasMaxInstancesString("2").
						HasQueryWarehouseString(testClient().Ids.WarehouseId().FullyQualifiedName()).
						HasServiceTypeString(string(sdk.ServiceTypeService)).
						HasCommentString(comment),
					resourceshowoutputassert.ServiceShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
						HasCurrentInstances(1).
						HasTargetInstances(1).
						HasMinReadyInstances(1).
						HasMinInstances(1).
						HasMaxInstances(2).
						HasAutoResume(true).
						HasExternalAccessIntegrations(externalAccessIntegrationId).
						HasCreatedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasResumedOnEmpty().
						HasSuspendedOnEmpty().
						HasAutoSuspendSecs(6767).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasQueryWarehouse(testClient().Ids.WarehouseId()).
						HasIsJob(false).
						HasIsAsyncJob(false).
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
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.current_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.target_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.min_ready_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.min_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.max_instances", "2")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.auto_resume", "true")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.external_access_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.external_access_integrations.0", externalAccessIntegrationId.Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.resumed_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.suspended_on", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.auto_suspend_secs", "6767")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.query_warehouse", testClient().Ids.WarehouseId().Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_job", "false")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.is_async_job", "false")),
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
func TestAcc_Service_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	computePoolId := testClient().Ids.RandomAccountObjectIdentifier()
	spec := testClient().Service.SampleSpec(t)

	modelCompleteWithInvalidAutoSuspendSecs := model.ServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePoolId.FullyQualifiedName(), spec).
		WithAutoSuspendSecs(-1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Service),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, modelCompleteWithInvalidAutoSuspendSecs),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected auto_suspend_secs to be at least \(0\), got -1`),
			},
		},
	})
}
