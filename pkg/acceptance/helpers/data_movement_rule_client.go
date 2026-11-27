package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

const dataMovementRuleMaxRows = "(100)"

type DataMovementRuleClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewDataMovementRuleClient(context *TestClientContext, idsGenerator *IdsGenerator) *DataMovementRuleClient {
	return &DataMovementRuleClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *DataMovementRuleClient) client() sdk.DataMovementRules {
	return c.context.client.DataMovementRules
}

func (c *DataMovementRuleClient) Create(t *testing.T) (*sdk.DataMovementRule, func()) {
	t.Helper()
	return c.CreateWithRequest(t, sdk.NewCreateDataMovementRuleRequest(c.ids.RandomSchemaObjectIdentifier(), sdk.DataMovementRuleTypeCopyIntoExternalStage, dataMovementRuleMaxRows))
}

func (c *DataMovementRuleClient) CreateWithRequest(t *testing.T, request *sdk.CreateDataMovementRuleRequest) (*sdk.DataMovementRule, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	id := request.GetName()
	rule, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return rule, c.DropFunc(t, id)
}

func (c *DataMovementRuleClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.DataMovementRule, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}

func (c *DataMovementRuleClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.DataMovementRuleDetails, error) {
	t.Helper()
	return c.client().Describe(context.Background(), id)
}

func (c *DataMovementRuleClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropDataMovementRuleRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
