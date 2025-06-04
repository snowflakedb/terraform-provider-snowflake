package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ServiceClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewServiceClient(context *TestClientContext, idsGenerator *IdsGenerator) *ServiceClient {
	return &ServiceClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ServiceClient) client() sdk.Services {
	return c.context.client.Services
}

func (c *ServiceClient) Create(t *testing.T, computePoolId sdk.AccountObjectIdentifier) (*sdk.Service, func()) {
	t.Helper()
	return c.CreateWithId(t, computePoolId, c.ids.RandomSchemaObjectIdentifier())
}

func (c *ServiceClient) CreateWithId(t *testing.T, computePoolId sdk.AccountObjectIdentifier, id sdk.SchemaObjectIdentifier) (*sdk.Service, func()) {
	t.Helper()
	spec := `
spec:
  containers:
  - name: example-container
    image: /snowflake/images/snowflake_images/exampleimage:latest
`
	return c.CreateWithRequest(t, sdk.NewCreateServiceRequest(id, computePoolId).WithFromSpecification(*sdk.NewServiceFromSpecificationRequest().WithSpecification(spec)))
}

func (c *ServiceClient) CreateWithIdWithBlockVolume(t *testing.T, computePoolId sdk.AccountObjectIdentifier, id sdk.SchemaObjectIdentifier) (*sdk.Service, func()) {
	t.Helper()
	spec := `
spec:
  containers:
  - name: example-container
    image: /snowflake/images/snowflake_images/exampleimage:latest
  volumes:
  - name: block-volume
    source: block
    size: 1Gi
`
	return c.CreateWithRequest(t, sdk.NewCreateServiceRequest(id, computePoolId).WithFromSpecification(*sdk.NewServiceFromSpecificationRequest().WithSpecification(spec)))
}

func (c *ServiceClient) CreateWithRequest(t *testing.T, req *sdk.CreateServiceRequest) (*sdk.Service, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)
	service, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)
	return service, c.DropFunc(t, req.GetName())
}

func (c *ServiceClient) ExecuteJobService(t *testing.T, computePoolId sdk.AccountObjectIdentifier, id sdk.SchemaObjectIdentifier) (*sdk.Service, func()) {
	t.Helper()
	ctx := context.Background()

	spec := `
spec:
  containers:
  - name: example-container
    image: /snowflake/images/snowflake_images/exampleimage:latest
`
	req := sdk.NewExecuteJobServiceServiceRequest(computePoolId, id).WithJobServiceFromSpecification(*sdk.NewJobServiceFromSpecificationRequest().WithSpecification(spec)).WithAsync(true)
	err := c.client().ExecuteJobService(ctx, req)
	require.NoError(t, err)
	service, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)
	return service, c.DropFunc(t, req.GetName())
}

func (c *ServiceClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropServiceRequest(id).WithIfExists(true).WithForce(true))
		require.NoError(t, err)
	}
}

func (c *ServiceClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.Service, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().ShowByID(ctx, id)
}

func (c *ServiceClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.ServiceDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().Describe(ctx, id)
}

func (c *ServiceClient) Alter(t *testing.T, req *sdk.AlterServiceRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}
