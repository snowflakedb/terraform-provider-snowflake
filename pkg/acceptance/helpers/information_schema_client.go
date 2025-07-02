package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type InformationSchemaClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewInformationSchemaClient(context *TestClientContext, idsGenerator *IdsGenerator) *InformationSchemaClient {
	return &InformationSchemaClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *InformationSchemaClient) client() *sdk.Client {
	return c.context.client
}

type QueryHistory struct {
	QueryId   string
	QueryText string
	QueryTag  string
}

func (c *InformationSchemaClient) GetQueryHistory(t *testing.T, limit int) []QueryHistory {
	t.Helper()

	result, err := c.client().QueryUnsafe(context.Background(), fmt.Sprintf("SELECT * FROM TABLE(INFORMATION_SCHEMA.QUERY_HISTORY(RESULT_LIMIT => %d))", limit))
	require.NoError(t, err)

	return collections.Map(result, func(queryResult map[string]*any) QueryHistory {
		return c.mapQueryHistory(t, queryResult)
	})
}

func (c *InformationSchemaClient) GetQueryHistoryByQueryId(t *testing.T, limit int, queryId string) QueryHistory {
	t.Helper()

	result, err := c.client().QueryUnsafe(context.Background(), fmt.Sprintf("SELECT * FROM TABLE(INFORMATION_SCHEMA.QUERY_HISTORY(RESULT_LIMIT => %d)) WHERE QUERY_ID = '%s'", limit, queryId))
	require.NoError(t, err)
	require.Len(t, result, 1)

	return c.mapQueryHistory(t, result[0])
}

func (c *InformationSchemaClient) mapQueryHistory(t *testing.T, queryResult map[string]*any) QueryHistory {
	t.Helper()
	var queryHistory QueryHistory
	if v, ok := queryResult["QUERY_ID"]; ok && v != nil {
		queryHistory.QueryId = (*v).(string)
	}
	if v, ok := queryResult["QUERY_TEXT"]; ok && v != nil {
		queryHistory.QueryText = (*v).(string)
	}
	if v, ok := queryResult["QUERY_TAG"]; ok && v != nil {
		queryHistory.QueryTag = (*v).(string)
	}
	return queryHistory
}

func (c *InformationSchemaClient) GetFunctionDataType(t *testing.T, functionId sdk.SchemaObjectIdentifierWithArguments) string {
	t.Helper()

	viewId := sdk.NewSchemaObjectIdentifier(functionId.DatabaseName(), "INFORMATION_SCHEMA", "FUNCTIONS")
	rows, err := c.client().QueryUnsafe(context.Background(), fmt.Sprintf(`SELECT DATA_TYPE FROM %s WHERE function_name = '%s'`, viewId.FullyQualifiedName(), functionId.Name()))
	require.NoError(t, err)
	require.Len(t, rows, 1)

	row := rows[0]
	require.NotNil(t, row["DATA_TYPE"])

	dataType, ok := (*rows[0]["DATA_TYPE"]).(string)
	require.True(t, ok)

	return dataType
}
