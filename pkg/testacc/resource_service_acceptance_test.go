//go:build !account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	acchelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/customassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
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
					assert.Check(resource.TestCheckResourceAttrWith(modelBasic.ResourceReference(), "show_output.0.current_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttrWith(modelBasic.ResourceReference(), "show_output.0.target_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttrWith(modelBasic.ResourceReference(), "describe_output.0.current_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttrWith(modelBasic.ResourceReference(), "describe_output.0.target_instances", customassert.BetweenFunc(0, 1))),
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
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateWith(helpers.EncodeResourceIdentifier(id), "show_output.0.current_instances", customassert.BetweenFunc(0, 2))),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateWith(helpers.EncodeResourceIdentifier(id), "show_output.0.target_instances", customassert.BetweenFunc(0, 2))),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.name", id.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.database_name", id.DatabaseName())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.schema_name", id.SchemaName())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateSet(helpers.EncodeResourceIdentifier(id), "describe_output.0.spec")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateSet(helpers.EncodeResourceIdentifier(id), "describe_output.0.dns_name")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateWith(helpers.EncodeResourceIdentifier(id), "describe_output.0.current_instances", customassert.BetweenFunc(0, 2))),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateWith(helpers.EncodeResourceIdentifier(id), "describe_output.0.target_instances", customassert.BetweenFunc(0, 2))),
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
					assert.Check(resource.TestCheckResourceAttrWith(modelComplete.ResourceReference(), "show_output.0.current_instances", customassert.BetweenFunc(0, 2))),
					assert.Check(resource.TestCheckResourceAttrWith(modelComplete.ResourceReference(), "show_output.0.target_instances", customassert.BetweenFunc(0, 2))),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttrWith(modelComplete.ResourceReference(), "describe_output.0.current_instances", customassert.BetweenFunc(0, 2))),
					assert.Check(resource.TestCheckResourceAttrWith(modelComplete.ResourceReference(), "describe_output.0.target_instances", customassert.BetweenFunc(0, 2))),
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
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateWith(helpers.EncodeResourceIdentifier(id), "show_output.0.current_instances", customassert.BetweenFunc(0, 2))),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateWith(helpers.EncodeResourceIdentifier(id), "show_output.0.target_instances", customassert.BetweenFunc(0, 2))),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.name", id.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.database_name", id.DatabaseName())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.schema_name", id.SchemaName())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateSet(helpers.EncodeResourceIdentifier(id), "describe_output.0.spec")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateSet(helpers.EncodeResourceIdentifier(id), "describe_output.0.dns_name")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateWith(helpers.EncodeResourceIdentifier(id), "describe_output.0.current_instances", customassert.BetweenFunc(0, 2))),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceStateWith(helpers.EncodeResourceIdentifier(id), "describe_output.0.target_instances", customassert.BetweenFunc(0, 2))),
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
					assert.Check(resource.TestCheckResourceAttrWith(modelCompleteWithDifferentValues.ResourceReference(), "show_output.0.current_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttrWith(modelCompleteWithDifferentValues.ResourceReference(), "show_output.0.target_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttrWith(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.current_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttrWith(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.target_instances", customassert.BetweenFunc(0, 1))),
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
					assert.Check(resource.TestCheckResourceAttrWith(modelCompleteWithDifferentValues.ResourceReference(), "show_output.0.current_instances", customassert.BetweenFunc(0, 2))),
					assert.Check(resource.TestCheckResourceAttrWith(modelCompleteWithDifferentValues.ResourceReference(), "show_output.0.target_instances", customassert.BetweenFunc(0, 2))),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttrWith(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.current_instances", customassert.BetweenFunc(0, 2))),
					assert.Check(resource.TestCheckResourceAttrWith(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.target_instances", customassert.BetweenFunc(0, 2))),
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
					assert.Check(resource.TestCheckResourceAttrWith(modelBasic.ResourceReference(), "show_output.0.current_instances", customassert.BetweenFunc(0, 2))),
					assert.Check(resource.TestCheckResourceAttrWith(modelBasic.ResourceReference(), "show_output.0.target_instances", customassert.BetweenFunc(0, 2))),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttrWith(modelBasic.ResourceReference(), "describe_output.0.current_instances", customassert.BetweenFunc(0, 2))),
					assert.Check(resource.TestCheckResourceAttrWith(modelBasic.ResourceReference(), "describe_output.0.target_instances", customassert.BetweenFunc(0, 2))),
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
	specTemplate, using := testClient().Service.SampleSpecTemplateWithUsingValue(t)
	specFileName := "spec.yaml"
	testClient().Stage.PutInLocationWithContent(t, stage.Location(), specFileName, spec)
	specTemplateFileName := "spec_template.yaml"
	testClient().Stage.PutInLocationWithContent(t, stage.Location(), specTemplateFileName, specTemplate)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelBasicOnStage := model.ServiceWithSpecOnStage("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), stage.ID(), specFileName)
	modelBasic := model.ServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), spec)
	modelBasicOnStageTemplate := model.ServiceWithSpecTemplateOnStage("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), stage.ID(), specTemplateFileName, using...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
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
			// update (text -> on stage)
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
						HasFromSpecificationOnStage(stage.ID(), "", specFileName),
					resourceshowoutputassert.ServiceShowOutput(t, modelBasicOnStage.ResourceReference()).
						HasName(id.Name()),
					assert.Check(resource.TestCheckResourceAttr(modelBasicOnStage.ResourceReference(), "describe_output.0.name", id.Name())),
				),
			},
			// update (on stage -> template on stage)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasicOnStageTemplate.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelBasicOnStageTemplate),
				Check: assertThat(t,
					resourceassert.ServiceResource(t, modelBasicOnStageTemplate.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTemplateOnStage(stage.ID(), "", specTemplateFileName, using...),
					resourceshowoutputassert.ServiceShowOutput(t, modelBasicOnStageTemplate.ResourceReference()).
						HasName(id.Name()),
					assert.Check(resource.TestCheckResourceAttr(modelBasicOnStageTemplate.ResourceReference(), "describe_output.0.name", id.Name())),
				),
			},
			// external changes are not detected
			{
				PreConfig: func() {
					testClient().Service.Alter(t, sdk.NewAlterServiceRequest(id).WithFromSpecification(*sdk.NewServiceFromSpecificationRequest().WithSpecification(testClient().Service.SampleSpecWithContainerName(t, "external-changed"))))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasicOnStageTemplate.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: accconfig.FromModels(t, modelBasicOnStageTemplate),
				Check: assertThat(t,
					resourceassert.ServiceResource(t, modelBasicOnStageTemplate.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()),
					resourceshowoutputassert.ServiceShowOutput(t, modelBasicOnStageTemplate.ResourceReference()).
						HasName(id.Name()),
					assert.Check(resource.TestCheckResourceAttr(modelBasicOnStageTemplate.ResourceReference(), "describe_output.0.name", id.Name())),
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
						planchecks.ExpectDrift(modelBasic.ResourceReference(), "service_type", sdk.Pointer(string(sdk.ServiceTypeService)), sdk.Pointer(string(sdk.ServiceTypeJobService))),
						planchecks.ExpectChange(modelBasic.ResourceReference(), "service_type", tfjson.ActionDelete, sdk.Pointer(string(sdk.ServiceTypeJobService)), nil),
						planchecks.ExpectChange(modelBasic.ResourceReference(), "service_type", tfjson.ActionCreate, sdk.Pointer(string(sdk.ServiceTypeJobService)), nil),
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
						HasFromSpecificationOnStage(stage.ID(), "", specFileName).
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
					assert.Check(resource.TestCheckResourceAttrWith(modelBasic.ResourceReference(), "show_output.0.current_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttrWith(modelBasic.ResourceReference(), "show_output.0.target_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(modelBasic.ResourceReference(), "describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttrWith(modelBasic.ResourceReference(), "describe_output.0.current_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttrWith(modelBasic.ResourceReference(), "describe_output.0.target_instances", customassert.BetweenFunc(0, 1))),
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

func TestAcc_Service_fromSpecificationTemplate(t *testing.T) {
	computePool, computePoolCleanup := testClient().ComputePool.Create(t)
	t.Cleanup(computePoolCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	specTemplate, using := testClient().Service.SampleSpecTemplateWithUsingValue(t)

	model := model.ServiceWithSpecTemplate("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), specTemplate, using...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Service),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model),
				Check: assertThat(t,
					resourceassert.ServiceResource(t, model.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTemplateTextNotEmpty(using...).
						HasAutoSuspendSecsString(r.IntDefaultString).
						HasExternalAccessIntegrationsEmpty().
						HasAutoResumeString(r.BooleanDefault).
						HasNoMinInstances().
						HasNoMinReadyInstances().
						HasNoMaxInstances().
						HasNoQueryWarehouse().
						HasCommentString(""),
					resourceshowoutputassert.ServiceShowOutput(t, model.ResourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
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
					assert.Check(resource.TestCheckResourceAttrWith(model.ResourceReference(), "show_output.0.current_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttrWith(model.ResourceReference(), "show_output.0.target_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(model.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(model.ResourceReference(), "describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttrWith(model.ResourceReference(), "describe_output.0.current_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttrWith(model.ResourceReference(), "describe_output.0.target_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.min_ready_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.min_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.max_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.auto_resume", "true")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.external_access_integrations.#", "0")),
					assert.Check(resource.TestCheckResourceAttrSet(model.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(model.ResourceReference(), "describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.resumed_on", "")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.suspended_on", "")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.auto_suspend_secs", "0")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.query_warehouse", "")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.is_job", "false")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.is_async_job", "false")),
					assert.Check(resource.TestCheckResourceAttrSet(model.ResourceReference(), "describe_output.0.spec_digest")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.is_upgrading", "false")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.managing_object_domain", "")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.managing_object_name", "")),
				),
			},
		},
	})
}

func TestAcc_Service_fromSpecificationTemplateOnStage(t *testing.T) {
	computePool, computePoolCleanup := testClient().ComputePool.Create(t)
	t.Cleanup(computePoolCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	specTemplate, using := testClient().Service.SampleSpecTemplateWithUsingValue(t)

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	specFileName := "spec.yaml"
	testClient().Stage.PutInLocationWithContent(t, stage.Location(), specFileName, specTemplate)

	model := model.ServiceWithSpecTemplateOnStage("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), stage.ID(), specFileName, using...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Service),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model),
				Check: assertThat(t,
					resourceassert.ServiceResource(t, model.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasComputePoolString(computePool.ID().FullyQualifiedName()).
						HasFromSpecificationTemplateOnStage(stage.ID(), "", specFileName, using...).
						HasAutoSuspendSecsString(r.IntDefaultString).
						HasExternalAccessIntegrationsEmpty().
						HasAutoResumeString(r.BooleanDefault).
						HasNoMinInstances().
						HasNoMinReadyInstances().
						HasNoMaxInstances().
						HasNoQueryWarehouse().
						HasCommentString(""),
					resourceshowoutputassert.ServiceShowOutput(t, model.ResourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.ServiceStatusPending).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComputePool(computePool.ID()).
						HasDnsNameNotEmpty().
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
					assert.Check(resource.TestCheckResourceAttrWith(model.ResourceReference(), "show_output.0.current_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttrWith(model.ResourceReference(), "show_output.0.target_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(model.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(model.ResourceReference(), "describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttrWith(model.ResourceReference(), "describe_output.0.current_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttrWith(model.ResourceReference(), "describe_output.0.target_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.min_ready_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.min_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.max_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.auto_resume", "true")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.external_access_integrations.#", "0")),
					assert.Check(resource.TestCheckResourceAttrSet(model.ResourceReference(), "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(model.ResourceReference(), "describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.resumed_on", "")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.suspended_on", "")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.auto_suspend_secs", "0")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.query_warehouse", "")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.is_job", "false")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.is_async_job", "false")),
					assert.Check(resource.TestCheckResourceAttrSet(model.ResourceReference(), "describe_output.0.spec_digest")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.is_upgrading", "false")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.managing_object_domain", "")),
					assert.Check(resource.TestCheckResourceAttr(model.ResourceReference(), "describe_output.0.managing_object_name", "")),
				),
			},
		},
	})
}

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
					assert.Check(resource.TestCheckResourceAttrWith(modelComplete.ResourceReference(), "show_output.0.current_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttrWith(modelComplete.ResourceReference(), "show_output.0.target_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(modelComplete.ResourceReference(), "describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttrWith(modelComplete.ResourceReference(), "describe_output.0.current_instances", customassert.BetweenFunc(0, 1))),
					assert.Check(resource.TestCheckResourceAttrWith(modelComplete.ResourceReference(), "describe_output.0.target_instances", customassert.BetweenFunc(0, 1))),
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

func TestAcc_Service_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	computePoolId := testClient().Ids.RandomAccountObjectIdentifier()
	spec := testClient().Service.SampleSpec(t)
	specTemplate := testClient().Service.SampleSpecTemplate(t)

	modelWithInvalidAutoSuspendSecs := model.ServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePoolId.FullyQualifiedName(), spec).
		WithAutoSuspendSecs(-1)
	modelWithInvalidAutoResume := model.ServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePoolId.FullyQualifiedName(), spec).
		WithAutoResume("invalid")
	modelWithInvalidMinInstances := model.ServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePoolId.FullyQualifiedName(), spec).
		WithMinInstances(0)
	modelWithInvalidMaxInstances := model.ServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePoolId.FullyQualifiedName(), spec).
		WithMaxInstances(0)
	modelWithInvalidMinReadyInstances := model.ServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePoolId.FullyQualifiedName(), spec).
		WithMinReadyInstances(0)

	modelWithSpecAndSpecTemplate := model.Service("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePoolId.FullyQualifiedName()).
		WithFromSpecification(spec).
		WithFromSpecificationTemplate(specTemplate, acchelpers.ServiceSpecUsing{Key: "key", Value: "value"})
	modelWithUsingMissingKey := model.Service("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePoolId.FullyQualifiedName()).
		WithFromSpecificationTemplateValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"text": config.MultilineWrapperVariable(spec),
			"using": tfconfig.SetVariable(
				tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"value": tfconfig.StringVariable("value"),
				}),
			),
		}))
	modelWithUsingMissingValue := model.Service("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePoolId.FullyQualifiedName()).
		WithFromSpecificationTemplateValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"text": config.MultilineWrapperVariable(spec),
			"using": tfconfig.SetVariable(
				tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"key": tfconfig.StringVariable("key"),
				}),
			),
		}))
	modelWithEmptyUsing := model.Service("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePoolId.FullyQualifiedName()).
		WithFromSpecificationTemplateValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"text": config.MultilineWrapperVariable(spec),
		}))
	modelWithNoSpecAndNoSpecTemplate := model.Service("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePoolId.FullyQualifiedName())
	modelWithEmptyExtAccessIntegrations := model.Service("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePoolId.FullyQualifiedName()).
		WithExternalAccessIntegrations()
	modelWithInvalidStage := model.Service("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePoolId.FullyQualifiedName()).
		WithFromSpecificationValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"stage": tfconfig.StringVariable("invalid"),
			"file":  tfconfig.StringVariable("file"),
		}))
	modelWithTextAndFile := model.Service("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePoolId.FullyQualifiedName()).
		WithFromSpecificationValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"file": tfconfig.StringVariable("file"),
			"text": config.MultilineWrapperVariable(spec),
		}))
	modelWithFileAndNoStage := model.Service("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePoolId.FullyQualifiedName()).
		WithFromSpecificationValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"file": tfconfig.StringVariable("file"),
		}))
	modelWithStageAndNoFile := model.Service("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePoolId.FullyQualifiedName()).
		WithFromSpecificationValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"stage": tfconfig.StringVariable("stage"),
		}))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Service),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, modelWithInvalidAutoSuspendSecs),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected auto_suspend_secs to be at least \(0\), got -1`),
			},
			{
				Config:      config.FromModels(t, modelWithInvalidAutoResume),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected \[\{\{} auto_resume}] to be one of \["true" "false"], got invalid`),
			},
			{
				Config:      config.FromModels(t, modelWithInvalidMinInstances),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected min_instances to be at least \(1\), got 0`),
			},
			{
				Config:      config.FromModels(t, modelWithInvalidMaxInstances),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected max_instances to be at least \(1\), got 0`),
			},
			{
				Config:      config.FromModels(t, modelWithInvalidMinReadyInstances),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected min_ready_instances to be at least \(1\), got 0`),
			},
			{
				Config:      config.FromModels(t, modelWithSpecAndSpecTemplate),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("`from_specification,from_specification_template` were specified"),
			},
			{
				Config:      config.FromModels(t, modelWithUsingMissingKey),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`The argument "key" is required, but no definition was found`),
			},
			{
				Config:      config.FromModels(t, modelWithUsingMissingValue),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`The argument "value" is required, but no definition was found`),
			},
			{
				Config:      config.FromModels(t, modelWithEmptyUsing),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`At least 1 "using" blocks are required.`),
			},
			{
				Config:      config.FromModels(t, modelWithNoSpecAndNoSpecTemplate),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("`from_specification,from_specification_template` must be specified"),
			},
			{
				Config:      config.FromModels(t, modelWithEmptyExtAccessIntegrations),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Attribute external_access_integrations requires 1 item minimum`),
			},
			{
				Config:      config.FromModels(t, modelWithInvalidStage),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Expected SchemaObjectIdentifier identifier type`),
			},
			{
				Config:      config.FromModels(t, modelWithTextAndFile),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("`from_specification.0.file,from_specification.0.text` were specified"),
			},
			{
				Config:      config.FromModels(t, modelWithFileAndNoStage),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("`from_specification.0.file,from_specification.0.stage` must be specified"),
			},
			{
				Config:      config.FromModels(t, modelWithStageAndNoFile),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("`from_specification.0.file,from_specification.0.stage` must be specified"),
			},
		},
	})
}
