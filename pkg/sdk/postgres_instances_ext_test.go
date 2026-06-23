package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresInstances_ParseDetails(t *testing.T) {
	t.Run("parse network policy into AccountObjectIdentifier", func(t *testing.T) {
		properties := []PostgresInstanceProperty{
			{Property: "name", Value: "test_instance"},
			{Property: "network_policy", Value: "my_network_policy"},
		}
		details, err := ParsePostgresInstanceDetails(properties)
		require.NoError(t, err)
		require.NotNil(t, details.NetworkPolicy)
		assert.Equal(t, NewAccountObjectIdentifier("my_network_policy"), *details.NetworkPolicy)
	})

	t.Run("parse storage integration into AccountObjectIdentifier", func(t *testing.T) {
		properties := []PostgresInstanceProperty{
			{Property: "name", Value: "test_instance"},
			{Property: "storage_integration", Value: "my_storage_integration"},
		}
		details, err := ParsePostgresInstanceDetails(properties)
		require.NoError(t, err)
		require.NotNil(t, details.StorageIntegration)
		assert.Equal(t, NewAccountObjectIdentifier("my_storage_integration"), *details.StorageIntegration)
	})

	t.Run("parse mixed-case property keys", func(t *testing.T) {
		properties := []PostgresInstanceProperty{
			{Property: "Name", Value: "test_instance"},
			{Property: "COMPUTE_FAMILY", Value: "STANDARD_M"},
			{Property: "Storage_Size_Gb", Value: "100"},
		}
		details, err := ParsePostgresInstanceDetails(properties)
		require.NoError(t, err)
		assert.Equal(t, "test_instance", details.Name)
		assert.Equal(t, "STANDARD_M", details.ComputeFamily)
		assert.Equal(t, 100, details.StorageSizeGb)
	})
}

func TestNormalizePostgresSettings(t *testing.T) {
	t.Run("empty and whitespace only", func(t *testing.T) {
		for _, s := range []string{"", "  ", "\t\n"} {
			got, err := NormalizePostgresSettings(s)
			require.NoError(t, err)
			require.Equal(t, "", got)
		}
	})

	t.Run("empty JSON object", func(t *testing.T) {
		got, err := NormalizePostgresSettings("{}")
		require.NoError(t, err)
		require.Equal(t, "", got)
	})

	t.Run("equivalent JSON with different formatting", func(t *testing.T) {
		want, err := NormalizePostgresSettings(`{"max_connections":"100","shared_buffers":"256MB"}`)
		require.NoError(t, err)

		equivalentForms := []string{
			`{"shared_buffers":"256MB","max_connections":"100"}`,
			`{  "max_connections"  :  "100"  ,  "shared_buffers"  :  "256MB"  }`,
			"{\n  \"max_connections\": \"100\",\n  \"shared_buffers\": \"256MB\"\n}",
		}
		for _, s := range equivalentForms {
			got, err := NormalizePostgresSettings(s)
			require.NoError(t, err)
			require.Equal(t, want, got)
		}
	})

	t.Run("non-equivalent JSON", func(t *testing.T) {
		want, err := NormalizePostgresSettings(`{"max_connections":"100"}`)
		require.NoError(t, err)

		got, err := NormalizePostgresSettings(`{"max_connections":"200"}`)
		require.NoError(t, err)
		require.NotEqual(t, want, got)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		_, err := NormalizePostgresSettings("{broken")
		require.Error(t, err)
	})
}

func TestNormalizePostgresSettingsPtr(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		require.Nil(t, NormalizePostgresSettingsPtr(nil))
	})

	t.Run("empty string returns nil", func(t *testing.T) {
		s := ""
		require.Nil(t, NormalizePostgresSettingsPtr(&s))
	})

	t.Run("empty JSON object returns nil", func(t *testing.T) {
		s := "{}"
		require.Nil(t, NormalizePostgresSettingsPtr(&s))
	})

	t.Run("valid JSON returns normalized pointer", func(t *testing.T) {
		s := `{"shared_buffers":"256MB","max_connections":"100"}`
		got := NormalizePostgresSettingsPtr(&s)
		require.NotNil(t, got)
		want, err := NormalizePostgresSettings(s)
		require.NoError(t, err)
		require.Equal(t, want, *got)
	})

	t.Run("invalid JSON returns nil", func(t *testing.T) {
		s := "{broken"
		require.Nil(t, NormalizePostgresSettingsPtr(&s))
	})
}
