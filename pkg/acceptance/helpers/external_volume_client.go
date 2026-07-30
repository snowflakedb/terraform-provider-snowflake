package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ExternalVolumeClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewExternalVolumeClient(context *TestClientContext, idsGenerator *IdsGenerator) *ExternalVolumeClient {
	return &ExternalVolumeClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ExternalVolumeClient) client() sdk.ExternalVolumes {
	return c.context.client.ExternalVolumes
}

// TODO(SNOW-999142): Switch to returning *sdk.ExternalVolume. Need to update existing acceptance tests for this.
func (c *ExternalVolumeClient) Create(t *testing.T) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()
	kmsKeyId := "1234abcd-12ab-34cd-56ef-1234567890ab"
	storageLocations := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: "my-s3-us-west-2",
			S3StorageLocationParams: &sdk.S3StorageLocationParams{
				StorageProvider:   sdk.S3StorageProviderS3,
				StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole",
				StorageBaseUrl:    "s3://my-example-bucket/",
				Encryption: &sdk.ExternalVolumeS3Encryption{
					EncryptionType: sdk.S3EncryptionTypeAwsSseKms,
					KmsKeyId:       &kmsKeyId,
				},
			},
		}},
	}

	req := sdk.NewCreateExternalVolumeRequest(id, storageLocations)
	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	_, showErr := c.client().ShowByID(ctx, id)
	require.NoError(t, showErr)

	return id, c.DropFunc(t, id)
}

func (c *ExternalVolumeClient) CreateWithRequest(t *testing.T, req *sdk.CreateExternalVolumeRequest) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()
	err := c.client().Create(ctx, req)
	require.NoError(t, err)
	return req.GetName(), c.DropFunc(t, req.GetName())
}

func (c *ExternalVolumeClient) CreateS3Compat(t *testing.T, s3CompatBaseUrl string, s3CompatEndpoint string, awsKeyId string, awsSecretKey string) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	id := c.ids.RandomAccountObjectIdentifier()
	storageLocations := []sdk.ExternalVolumeStorageLocationItem{
		{ExternalVolumeStorageLocation: sdk.ExternalVolumeStorageLocation{
			Name: "my-s3compat-loc",
			S3CompatStorageLocationParams: &sdk.S3CompatStorageLocationParams{
				StorageBaseUrl:  s3CompatBaseUrl,
				StorageEndpoint: s3CompatEndpoint,
				Credentials: sdk.ExternalVolumeS3CompatCredentials{
					AwsKeyId:     awsKeyId,
					AwsSecretKey: awsSecretKey,
				},
			},
		}},
	}
	req := sdk.NewCreateExternalVolumeRequest(id, storageLocations)
	return c.CreateWithRequest(t, req)
}

func (c *ExternalVolumeClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ExternalVolume, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *ExternalVolumeClient) Describe(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ExternalVolumeDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().DescribeDetails(ctx, id)
}

func (c *ExternalVolumeClient) Alter(t *testing.T, req *sdk.AlterExternalVolumeRequest) {
	t.Helper()
	ctx := context.Background()
	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}

// ReplaceLocation replaces a storage location with a new one.
// It creates a temporary storage location, performs the update, and ensures
// the temporary storage location is always cleaned up.
// This is needed because Snowflake does not allow external volumes to have no storage locations.
func (c *ExternalVolumeClient) ReplaceLocation(t *testing.T, id sdk.AccountObjectIdentifier, req sdk.ExternalVolumeStorageLocationRequest) {
	t.Helper()
	tempLocation := "temp_location"
	c.Alter(t, sdk.NewAlterExternalVolumeRequest(id).WithAddStorageLocation(
		*sdk.NewExternalVolumeStorageLocationItemRequest(
			*sdk.NewExternalVolumeStorageLocationRequest(tempLocation).WithS3StorageLocationParams(*sdk.NewS3StorageLocationParamsRequest(sdk.S3StorageProviderS3, "arn:aws:iam::123456789012:role/myrole", "s3://my-example-bucket/")),
		),
	))
	c.Alter(t, sdk.NewAlterExternalVolumeRequest(id).WithRemoveStorageLocation(req.Name))
	c.Alter(t, sdk.NewAlterExternalVolumeRequest(id).WithAddStorageLocation(*sdk.NewExternalVolumeStorageLocationItemRequest(req)))
	c.Alter(t, sdk.NewAlterExternalVolumeRequest(id).WithRemoveStorageLocation(tempLocation))
}

func (c *ExternalVolumeClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropExternalVolumeRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
