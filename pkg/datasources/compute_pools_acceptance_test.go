//go:build !account_level_tests

package datasources_test

import (
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO [SNOW-1431726]: Move to helpers
func createApp(t *testing.T) *sdk.Application {
	t.Helper()

	stage, cleanupStage := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(cleanupStage)

	acc.TestClient().Stage.PutOnStage(t, stage.ID(), "TestAcc_Application/manifest.yml")
	acc.TestClient().Stage.PutOnStage(t, stage.ID(), "TestAcc_Application/setup.sql")

	applicationPackage, cleanupApplicationPackage := acc.TestClient().ApplicationPackage.CreateApplicationPackage(t)
	t.Cleanup(cleanupApplicationPackage)

	acc.TestClient().ApplicationPackage.AddApplicationPackageVersion(t, applicationPackage.ID(), stage.ID(), "v1")

	application, cleanupApplication := acc.TestClient().Application.CreateApplication(t, applicationPackage.ID(), "v1")
	t.Cleanup(cleanupApplication)
	return application
}

func TestAcc_ComputePools(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	application := createApp(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	computePoolModel := model.ComputePool("test", id.Name(), string(sdk.ComputePoolInstanceFamilyCpuX64S), 2, 1).
		WithForApplication(application.ID().FullyQualifiedName()).
		WithAutoResume("true").
		WithInitiallySuspended("true").
		WithAutoSuspendSecs(6767).
		WithComment(comment)

	dataSourceModel := datasourcemodel.ComputePools("test").
		WithLike(id.Name()).
		WithDependsOn(computePoolModel.ResourceReference())

	dataSourceModelWithoutOptionals := datasourcemodel.ComputePools("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithDependsOn(computePoolModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, computePoolModel, dataSourceModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.#", "1")),

					resourceshowoutputassert.ComputePoolsDatasourceShowOutput(t, "snowflake_compute_pools.test").
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
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.state", string(sdk.ComputePoolStateSuspended))),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.min_nodes", "1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.max_nodes", "2")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.instance_family", string(sdk.ComputePoolInstanceFamilyCpuX64S))),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.num_services", "0")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.num_jobs", "0")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.auto_suspend_secs", "6767")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.auto_resume", "true")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.active_nodes", "0")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.idle_nodes", "0")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.target_nodes", "0")),
					assert.Check(resource.TestCheckResourceAttrSet(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.resumed_on")),
					assert.Check(resource.TestCheckResourceAttrSet(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.is_exclusive", "true")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.application", application.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.error_code", "")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "compute_pools.0.describe_output.0.status_message", "")),
				),
			},
			{
				Config: accconfig.FromModels(t, computePoolModel, dataSourceModelWithoutOptionals),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dataSourceModelWithoutOptionals.DatasourceReference(), "compute_pools.#", "1")),

					resourceshowoutputassert.ComputePoolsDatasourceShowOutput(t, "snowflake_compute_pools.test").
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
					assert.Check(resource.TestCheckResourceAttr(dataSourceModelWithoutOptionals.DatasourceReference(), "compute_pools.0.describe.output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_ComputePools_Filtering(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	prefix := random.AlphaUpperN(4)
	idOne := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	cpModel1 := model.ComputePool("test", idOne.Name(), string(sdk.ComputePoolInstanceFamilyCpuX64XS), 1, 1)
	cpModel2 := model.ComputePool("test1", idTwo.Name(), string(sdk.ComputePoolInstanceFamilyCpuX64XS), 1, 1)
	cpModel3 := model.ComputePool("test2", idThree.Name(), string(sdk.ComputePoolInstanceFamilyCpuX64XS), 1, 1)
	computePoolsWithLikeModel := datasourcemodel.ComputePools("test").
		WithLike(idOne.Name()).
		WithDependsOn(cpModel1.ResourceReference(), cpModel2.ResourceReference(), cpModel3.ResourceReference())
	computePoolsWithStartsWithModel := datasourcemodel.ComputePools("test").
		WithStartsWith(prefix).
		WithDependsOn(cpModel1.ResourceReference(), cpModel2.ResourceReference(), cpModel3.ResourceReference())
	computePoolsWithLimitModel := datasourcemodel.ComputePools("test").
		WithRowsAndFrom(1, prefix).
		WithDependsOn(cpModel1.ResourceReference(), cpModel2.ResourceReference(), cpModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ComputePool),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, cpModel1, cpModel2, cpModel3, computePoolsWithLikeModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(computePoolsWithLikeModel.DatasourceReference(), "compute_pools.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, cpModel1, cpModel2, cpModel3, computePoolsWithStartsWithModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(computePoolsWithLikeModel.DatasourceReference(), "compute_pools.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, cpModel1, cpModel2, cpModel3, computePoolsWithLimitModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(computePoolsWithLikeModel.DatasourceReference(), "compute_pools.#", "1"),
				),
			},
		},
	})
}

func TestAcc_ComputePools_NotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ComputePools/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one compute pool"),
			},
		},
	})
}
