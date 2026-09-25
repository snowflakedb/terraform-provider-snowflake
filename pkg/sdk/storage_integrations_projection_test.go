package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorageIntegrationDescribeProjection(t *testing.T) {
	enabled := StorageIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "true"}
	provider := StorageIntegrationProperty{Name: "STORAGE_PROVIDER", Type: "String", Value: "S3", Default: ""}
	comment := StorageIntegrationProperty{Name: "COMMENT", Type: "String", Value: "c", Default: ""}
	unknown := StorageIntegrationProperty{Name: "AZURE_TENANT_ID", Type: "String", Value: "tid", Default: ""}

	details := AsStorageIntegrationDescribe([]StorageIntegrationProperty{enabled, provider, comment, unknown})
	require.Equal(t, &enabled, details.Enabled)
	require.Equal(t, &provider, details.StorageProvider)
	require.Equal(t, &comment, details.Comment)
	require.Nil(t, details.AzureConsentUrl)
	require.Nil(t, AsStorageIntegrationDescribe(nil))
}
