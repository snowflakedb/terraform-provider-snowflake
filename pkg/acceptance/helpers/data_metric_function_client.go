package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// TODO(SNOW-1564959): change raw sqls to proper client
type DataMetricFunctionClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewDataMetricFunctionClient(context *TestClientContext, idsGenerator *IdsGenerator) *DataMetricFunctionClient {
	return &DataMetricFunctionClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *DataMetricFunctionClient) client() *sdk.Client {
	return c.context.client
}

func (c *DataMetricFunctionClient) CreateDataMetricFunction(t *testing.T, viewID sdk.SchemaObjectIdentifier) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`CREATE DATA METRIC FUNCTION %s(arg_t TABLE (arg_c INT))
RETURNS NUMBER AS
'SELECT COUNT(*) FROM arg_t
   WHERE arg_c IN (SELECT id FROM %s)'`, id.Name(), viewID.Name()))
	require.NoError(t, err)
	return id, c.DropDataMetricFunctionFunc(t, id)
}

func (c *DataMetricFunctionClient) DropDataMetricFunctionFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`DROP FUNCTION IF EXISTS %s (TABLE (INT))`, id.Name()))
		require.NoError(t, err)
	}
}
