//go:build !account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Services(t *testing.T) {
	computePool, computePoolCleanup := testClient().ComputePool.Create(t)
	t.Cleanup(computePoolCleanup)

	networkRule, networkRuleCleanup := testClient().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)

	externalAccessIntegrationId, externalAccessIntegration1Cleanup := testClient().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
	t.Cleanup(externalAccessIntegration1Cleanup)

	comment := random.Comment()
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	spec := testClient().Service.SampleSpec(t)

	serviceModel := model.ServiceWithSpec("test", id.DatabaseName(), id.SchemaName(), id.Name(), computePool.ID().FullyQualifiedName(), spec).
		WithAutoSuspendSecs(6767).
		WithExternalAccessIntegrations(externalAccessIntegrationId).
		WithAutoResume("true").
		WithMinInstances(1).
		WithMinReadyInstances(1).
		WithMaxInstances(2).
		WithQueryWarehouse(testClient().Ids.WarehouseId().FullyQualifiedName()).
		WithComment(comment)

	dataSourceModel := datasourcemodel.Services("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(serviceModel.ResourceReference())

	dataSourceModelWithoutOptionals := datasourcemodel.Services("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithDependsOn(serviceModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, serviceModel, dataSourceModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.#", "1")),

					resourceshowoutputassert.ServicesDatasourceShowOutput(t, dataSourceModel.DatasourceReference()).
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
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.status", string(sdk.ServiceStatusPending))),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.compute_pool", computePool.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.spec")),
					assert.Check(resource.TestCheckResourceAttrSet(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.dns_name")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.current_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.target_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.min_ready_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.min_instances", "1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.max_instances", "2")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.auto_resume", "true")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.external_access_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.external_access_integrations.0", externalAccessIntegrationId.Name())),
					assert.Check(resource.TestCheckResourceAttrSet(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.resumed_on", "")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.suspended_on", "")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.auto_suspend_secs", "6767")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.query_warehouse", testClient().Ids.WarehouseId().Name())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.is_job", "false")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.is_async_job", "false")),
					assert.Check(resource.TestCheckResourceAttrSet(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.spec_digest")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.is_upgrading", "false")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.managing_object_domain", "")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "services.0.describe_output.0.managing_object_name", "")),
				),
			},
			{
				Config: accconfig.FromModels(t, serviceModel, dataSourceModelWithoutOptionals),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dataSourceModelWithoutOptionals.DatasourceReference(), "services.#", "1")),

					resourceshowoutputassert.ServicesDatasourceShowOutput(t, dataSourceModelWithoutOptionals.DatasourceReference()).
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
					assert.Check(resource.TestCheckResourceAttr(dataSourceModelWithoutOptionals.DatasourceReference(), "services.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_Services_Filtering(t *testing.T) {
	computePool1, computePool1Cleanup := testClient().ComputePool.Create(t)
	t.Cleanup(computePool1Cleanup)

	computePool2, computePool2Cleanup := testClient().ComputePool.Create(t)
	t.Cleanup(computePool2Cleanup)

	spec := testClient().Service.SampleSpec(t)
	prefix := random.AlphaUpperN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifier()
	id4 := testClient().Ids.RandomSchemaObjectIdentifier()

	model1 := model.ServiceWithSpec("test1", id1.DatabaseName(), id1.SchemaName(), id1.Name(), computePool1.ID().FullyQualifiedName(), spec)
	model2 := model.ServiceWithSpec("test2", id2.DatabaseName(), id2.SchemaName(), id2.Name(), computePool2.ID().FullyQualifiedName(), spec)
	model3 := model.ServiceWithSpec("test3", id3.DatabaseName(), id3.SchemaName(), id3.Name(), computePool2.ID().FullyQualifiedName(), spec)
	jobModel := model.JobServiceWithSpec("test", id4.DatabaseName(), id4.SchemaName(), id4.Name(), computePool2.ID().FullyQualifiedName(), spec)

	dataSourceModelLikeFirstOne := datasourcemodel.Services("test").
		WithLike(id1.Name()).
		WithWithDescribe(false).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference(), jobModel.ResourceReference())
	dataSourceModelLikePrefix := datasourcemodel.Services("test").
		WithLike(prefix+"%").
		WithWithDescribe(false).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference(), jobModel.ResourceReference())
	dataSourceModelInComputePool := datasourcemodel.Services("test").
		WithWithDescribe(false).
		WithInComputePool(computePool2.ID()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference(), jobModel.ResourceReference())
	dataSourceModelInComputePoolJobsOnly := datasourcemodel.Services("test").
		WithWithDescribe(false).
		WithServiceType(string(datasources.ShowServicesTypeJobsOnly)).
		WithInComputePool(computePool2.ID()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference(), jobModel.ResourceReference())
	dataSourceModelInComputePoolExcludeJobs := datasourcemodel.Services("test").
		WithWithDescribe(false).
		WithServiceType(string(datasources.ShowServicesTypeServicesOnly)).
		WithInComputePool(computePool2.ID()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference(), jobModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, jobModel, dataSourceModelLikeFirstOne),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceModelLikeFirstOne.DatasourceReference(), "services.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, jobModel, dataSourceModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceModelLikePrefix.DatasourceReference(), "services.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, jobModel, dataSourceModelInComputePool),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceModelInComputePool.DatasourceReference(), "services.#", "3"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, jobModel, dataSourceModelInComputePoolJobsOnly),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceModelInComputePoolJobsOnly.DatasourceReference(), "services.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, jobModel, dataSourceModelInComputePoolExcludeJobs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceModelInComputePoolExcludeJobs.DatasourceReference(), "services.#", "2"),
				),
			},
		},
	})
}

func TestAcc_Services_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.Services("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Services_NotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Services/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one service"),
			},
		},
	})
}
