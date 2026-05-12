//go:build non_account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ResourceMonitors_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	resourceMonitorId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)

	resourceMonitorModel := model.ResourceMonitor("rm", resourceMonitorId.Name()).
		WithCreditQuota(5)

	resourceMonitorsModelLikePrefix := datasourcemodel.ResourceMonitors("test").
		WithLike(prefix + "%").
		WithDependsOn(resourceMonitorModel.ResourceReference())

	resourceMonitorsModelLikeExact := datasourcemodel.ResourceMonitors("test").
		WithLike(resourceMonitorId.Name()).
		WithDependsOn(resourceMonitorModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// like (prefix)
			{
				Config: accconfig.FromModels(t, resourceMonitorModel, resourceMonitorsModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceMonitorsModelLikePrefix.DatasourceReference(), "resource_monitors.#", "1"),
				),
			},
			// like (exact)
			{
				Config: accconfig.FromModels(t, resourceMonitorModel, resourceMonitorsModelLikeExact),
				Check: assertThat(t,
					resourceshowoutputassert.ResourceMonitorDatasourceShowOutput(t, resourceMonitorsModelLikeExact.DatasourceReference()).
						HasName(resourceMonitorId.Name()).
						HasCreditQuota(5).
						HasUsedCredits(0).
						HasRemainingCredits(5).
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
		},
	})
}

func TestAcc_ResourceMonitors_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resourceMonitorModel1 := model.ResourceMonitor("test", id.Name()).
		WithCreditQuota(5)

	datasourceModel := datasourcemodel.ResourceMonitors("test").
		WithLike(id.Name()).
		WithDependsOn(resourceMonitorModel1.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ResourceMonitor),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, resourceMonitorModel1, datasourceModel),
				Check: assertThat(t,
					resourceshowoutputassert.ResourceMonitorDatasourceShowOutput(t, datasourceModel.DatasourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasCreditQuota(5).
						HasUsedCredits(0).
						HasRemainingCredits(5).
						HasLevel("").
						HasFrequency(sdk.FrequencyMonthly).
						HasStartTimeNotEmpty().
						HasEndTime("").
						HasSuspendAt(0).
						HasSuspendImmediateAt(0).
						HasOwnerNotEmpty().
						HasComment(""),
					assert.Check(resource.TestCheckResourceAttr(datasourceModel.DatasourceReference(), "resource_monitors.#", "1")),
				),
			},
		},
	})
}
