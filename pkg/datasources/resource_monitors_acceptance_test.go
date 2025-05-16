//go:build !account_level_tests

package datasources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ResourceMonitors(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	prefix := "data_source_resource_monitor_"
	resourceMonitorId := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	resourceMonitorId2 := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)

	resourceMonitorModel1 := model.ResourceMonitor("rm1", resourceMonitorId.Name()).
		WithCreditQuota(5)
	resourceMonitorModel2 := model.ResourceMonitor("rm2", resourceMonitorId2.Name()).
		WithCreditQuota(15)
	resourceMonitorsModelLikePrefix := datasourcemodel.ResourceMonitors("test").
		WithLike(prefix+"%").
		WithDependsOn(resourceMonitorModel1.ResourceReference(), resourceMonitorModel2.ResourceReference())
	resourceMonitorsModelLikeFirstMonitorName := datasourcemodel.ResourceMonitors("test").
		WithLike(resourceMonitorId.Name()).
		WithDependsOn(resourceMonitorModel1.ResourceReference(), resourceMonitorModel2.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Filter by prefix pattern (expect 2 items)
			{
				Config: accconfig.FromModels(t, resourceMonitorModel1, resourceMonitorModel2, resourceMonitorsModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceMonitorsModelLikePrefix.DatasourceReference(), "resource_monitors.#", "2"),
				),
			},
			// Filter by exact name (expect 1 item)
			{
				Config: accconfig.FromModels(t, resourceMonitorModel1, resourceMonitorModel2, resourceMonitorsModelLikeFirstMonitorName),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(resourceMonitorsModelLikeFirstMonitorName.DatasourceReference(), "resource_monitors.#", "1")),
					resourceshowoutputassert.ResourceMonitorDatasourceShowOutput(t, "snowflake_resource_monitors.test").
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
