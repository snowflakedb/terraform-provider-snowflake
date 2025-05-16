package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ApplicationPackageClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewApplicationPackageClient(context *TestClientContext, idsGenerator *IdsGenerator) *ApplicationPackageClient {
	return &ApplicationPackageClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ApplicationPackageClient) client() sdk.ApplicationPackages {
	return c.context.client.ApplicationPackages
}

func (c *ApplicationPackageClient) CreateApplicationPackage(t *testing.T) (*sdk.ApplicationPackage, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()
	err := c.client().Create(ctx, sdk.NewCreateApplicationPackageRequest(id))
	require.NoError(t, err)

	applicationPackage, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return applicationPackage, c.DropApplicationPackageFunc(t, id)
}

func (c *ApplicationPackageClient) DropApplicationPackageFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropApplicationPackageRequest(id).WithIfExists(sdk.Bool(true)))
		require.NoError(t, err)
	}
}

func (c *ApplicationPackageClient) AddApplicationPackageVersion(t *testing.T, id sdk.AccountObjectIdentifier, stageId sdk.SchemaObjectIdentifier, versionName string) {
	t.Helper()
	ctx := context.Background()

	using := "@" + stageId.FullyQualifiedName()

	err := c.client().Alter(ctx, sdk.NewAlterApplicationPackageRequest(id).WithAddVersion(sdk.NewAddVersionRequest(using).WithVersionIdentifier(sdk.String(versionName))))
	require.NoError(t, err)
}

func (c *ApplicationPackageClient) ShowVersions(t *testing.T, id sdk.AccountObjectIdentifier) []ApplicationPackageVersion {
	t.Helper()

	var versions []ApplicationPackageVersion
	err := c.context.client.QueryForTests(context.Background(), &versions, fmt.Sprintf(`SHOW VERSIONS IN APPLICATION PACKAGE %s`, id.FullyQualifiedName()))
	require.NoError(t, err)
	return versions
}

type ApplicationPackageVersion struct {
	Version string `json:"version"`
	Patch   int    `json:"patch"`
}
