package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type DataMetricFunctionReferencesClient struct {
	context *TestClientContext
}

func NewDataMetricFunctionReferencesClient(context *TestClientContext) *DataMetricFunctionReferencesClient {
	return &DataMetricFunctionReferencesClient{
		context: context,
	}
}

// GetDataMetricFunctionReferences is based on https://docs.snowflake.com/en/sql-reference/functions/data_metric_function_references.
func (c *DataMetricFunctionReferencesClient) GetDataMetricFunctionReferences(t *testing.T, id sdk.SchemaObjectIdentifier, domain sdk.DataMetricFunctionRefEntityDomainOption) []sdk.DataMetricFunctionReference {
	t.Helper()
	ctx := context.Background()

	refs, err := c.context.client.DataMetricFunctionReferences.GetForEntity(ctx, sdk.NewGetForEntityDataMetricFunctionReferenceRequest(id, domain))
	require.NoError(t, err)

	return refs
}
