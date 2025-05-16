package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ExternalFunctionClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewExternalFunctionClient(context *TestClientContext, idsGenerator *IdsGenerator) *ExternalFunctionClient {
	return &ExternalFunctionClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ExternalFunctionClient) client() sdk.ExternalFunctions {
	return c.context.client.ExternalFunctions
}

func (c *ExternalFunctionClient) Create(t *testing.T, apiIntegrationId sdk.AccountObjectIdentifier, arguments ...sdk.DataType) *sdk.ExternalFunction {
	t.Helper()
	return c.CreateWithIdentifier(t, apiIntegrationId, c.ids.RandomSchemaObjectIdentifierWithArguments(arguments...))
}

func (c *ExternalFunctionClient) CreateWithIdentifier(t *testing.T, apiIntegrationId sdk.AccountObjectIdentifier, id sdk.SchemaObjectIdentifierWithArguments) *sdk.ExternalFunction {
	t.Helper()
	ctx := context.Background()
	argumentRequests := make([]sdk.ExternalFunctionArgumentRequest, len(id.ArgumentDataTypes()))
	for i, argumentDataType := range id.ArgumentDataTypes() {
		argumentRequests[i] = *sdk.NewExternalFunctionArgumentRequest(c.ids.Alpha(), argumentDataType)
	}
	err := c.client().Create(ctx,
		sdk.NewCreateExternalFunctionRequest(
			id.SchemaObjectId(),
			sdk.DataTypeVariant,
			&apiIntegrationId,
			"https://xyz.execute-api.us-west-2.amazonaws.com/production/remote_echo",
		).WithArguments(argumentRequests),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, c.context.client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id).WithIfExists(true)))
	})

	externalFunction, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return externalFunction
}
