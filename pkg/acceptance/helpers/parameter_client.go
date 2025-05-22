package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ParameterClient struct {
	context *TestClientContext
}

func NewParameterClient(context *TestClientContext) *ParameterClient {
	return &ParameterClient{
		context: context,
	}
}

func (c *ParameterClient) client() sdk.Parameters {
	return c.context.client.Parameters
}

func (c *ParameterClient) ShowAccountParameters(t *testing.T) []*sdk.Parameter {
	t.Helper()
	params, err := c.client().ShowParameters(context.Background(), &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Account: sdk.Bool(true),
		},
	})
	require.NoError(t, err)
	return params
}

func (c *ParameterClient) ShowDatabaseParameters(t *testing.T, id sdk.AccountObjectIdentifier) []*sdk.Parameter {
	t.Helper()
	params, err := c.client().ShowParameters(context.Background(), &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Database: id,
		},
	})
	require.NoError(t, err)
	return params
}

func (c *ParameterClient) ShowWarehouseParameters(t *testing.T, id sdk.AccountObjectIdentifier) []*sdk.Parameter {
	t.Helper()
	params, err := c.client().ShowParameters(context.Background(), &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Warehouse: id,
		},
	})
	require.NoError(t, err)
	return params
}

func (c *ParameterClient) ShowSchemaParameters(t *testing.T, id sdk.DatabaseObjectIdentifier) []*sdk.Parameter {
	t.Helper()
	params, err := c.client().ShowParameters(context.Background(), &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Schema: id,
		},
	})
	require.NoError(t, err)
	return params
}

func (c *ParameterClient) ShowUserParameters(t *testing.T, id sdk.AccountObjectIdentifier) []*sdk.Parameter {
	t.Helper()
	params, err := c.client().ShowParameters(context.Background(), &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			User: id,
		},
	})
	require.NoError(t, err)
	return params
}

func (c *ParameterClient) ShowUserParametersOrError(t *testing.T, id sdk.AccountObjectIdentifier) ([]*sdk.Parameter, error) {
	t.Helper()
	return c.client().ShowParameters(context.Background(), &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			User: id,
		},
	})
}

func (c *ParameterClient) ShowTaskParameters(t *testing.T, id sdk.SchemaObjectIdentifier) []*sdk.Parameter {
	t.Helper()
	params, err := c.client().ShowParameters(context.Background(), &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Task: id,
		},
	})
	require.NoError(t, err)
	return params
}

func (c *ParameterClient) ShowFunctionParameters(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) []*sdk.Parameter {
	t.Helper()
	params, err := c.client().ShowParameters(context.Background(), &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Function: id,
		},
	})
	require.NoError(t, err)
	return params
}

func (c *ParameterClient) ShowProcedureParameters(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) []*sdk.Parameter {
	t.Helper()
	params, err := c.client().ShowParameters(context.Background(), &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Procedure: id,
		},
	})
	require.NoError(t, err)
	return params
}

func (c *ParameterClient) ShowAccountParameter(t *testing.T, parameter sdk.AccountParameter) *sdk.Parameter {
	t.Helper()
	ctx := context.Background()

	param, err := c.client().ShowAccountParameter(ctx, parameter)
	require.NoError(t, err)

	return param
}

// TODO [SNOW-1473408]: add unset account parameter to sdk.Parameters
func (c *ParameterClient) UnsetAccountParameter(t *testing.T, parameter sdk.AccountParameter) {
	t.Helper()
	ctx := context.Background()

	_, err := c.context.client.ExecForTests(ctx, fmt.Sprintf("ALTER ACCOUNT UNSET %s", parameter))
	require.NoError(t, err)
}

func FindParameter[T ~string](t *testing.T, parameters []*sdk.Parameter, parameter T) *sdk.Parameter {
	t.Helper()
	param, err := collections.FindFirst(parameters, func(p *sdk.Parameter) bool { return p.Key == string(parameter) })
	require.NoError(t, err)
	return *param
}
