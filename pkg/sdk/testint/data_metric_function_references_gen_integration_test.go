//go:build non_account_level_tests

package testint

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_DataMetricFunctionReferences(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("view domain", func(t *testing.T) {
		functionId := sdk.NewSchemaObjectIdentifier("SNOWFLAKE", "CORE", "BLANK_COUNT")
		statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
		view, viewCleanup := testClientHelper().View.CreateView(t, statement)
		t.Cleanup(viewCleanup)

		// when we specify schedule by a number of minutes, a cron is returned from Snowflake - see SNOW-1640024
		err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(view.ID()).WithSetDataMetricSchedule(*sdk.NewViewSetDataMetricScheduleRequest("5 MINUTE")))
		require.NoError(t, err)
		err = client.Views.Alter(ctx, sdk.NewAlterViewRequest(view.ID()).WithAddDataMetricFunction(*sdk.NewViewAddDataMetricFunctionRequest([]sdk.ViewDataMetricFunction{{
			DataMetricFunction: functionId,
			On:                 []sdk.Column{{Value: "ROLE_NAME"}},
		}})))
		require.NoError(t, err)

		dmfs, err := client.DataMetricFunctionReferences.GetForEntity(ctx, sdk.NewGetForEntityDataMetricFunctionReferenceRequest(view.ID(), sdk.DataMetricFunctionRefEntityDomainView))
		require.NoError(t, err)
		require.Len(t, dmfs, 1)
		dmf := dmfs[0]
		assert.Equal(t, string(sdk.DataMetricFunctionRefEntityDomainView), strings.ToUpper(dmf.RefEntityDomain))
		assert.Equal(t, functionId.DatabaseName(), dmf.MetricDatabaseName)
		assert.Equal(t, functionId.SchemaName(), dmf.MetricSchemaName)
		assert.Equal(t, functionId.Name(), dmf.MetricName)
		assert.Equal(t, view.ID().DatabaseName(), dmf.RefEntityDatabaseName)
		assert.Equal(t, view.ID().SchemaName(), dmf.RefEntitySchemaName)
		assert.Equal(t, view.ID().Name(), dmf.RefEntityName)
		assert.Equal(t, fmt.Sprintf("TABLE(VARCHAR(%d))", datatypes.DefaultVarcharLength), dmf.ArgumentSignature)
		assert.Equal(t, "NUMBER(38,0)", dmf.DataType)
		assert.NotEmpty(t, dmf.RefArguments)
		assert.NotEmpty(t, dmf.RefId)
		assert.Equal(t, "*/5 * * * * UTC", dmf.Schedule)
		assert.Equal(t, string(sdk.DataMetricScheduleStatusStarted), dmf.ScheduleStatus)
	})
}
