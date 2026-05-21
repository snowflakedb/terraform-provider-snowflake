package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests validate that NetworkPolicy and StorageIntegration use AccountObjectIdentifier
// instead of *string, and render as double-quoted identifiers in SQL.

func TestPostgresInstances_Create_IdentifierTypes(t *testing.T) {
	id := randomAccountObjectIdentifier()

	t.Run("create with network policy as identifier", func(t *testing.T) {
		networkPolicyId := NewAccountObjectIdentifier("my_policy")
		opts := &CreatePostgresInstanceOptions{
			name:                    id,
			ComputeFamily:           "STANDARD_S",
			StorageSizeGb:           50,
			AuthenticationAuthority: PostgresInstanceAuthenticationAuthorityPostgres,
			NetworkPolicy:           &networkPolicyId,
		}
		assertOptsValidAndSQLEquals(t, opts,
			`CREATE POSTGRES INSTANCE %s COMPUTE_FAMILY = 'STANDARD_S' STORAGE_SIZE_GB = 50 AUTHENTICATION_AUTHORITY = POSTGRES NETWORK_POLICY = "my_policy"`,
			id.FullyQualifiedName())
	})

	t.Run("create with storage integration as identifier", func(t *testing.T) {
		storageIntegrationId := NewAccountObjectIdentifier("my_integration")
		opts := &CreatePostgresInstanceOptions{
			name:                    id,
			ComputeFamily:           "STANDARD_S",
			StorageSizeGb:           50,
			AuthenticationAuthority: PostgresInstanceAuthenticationAuthorityPostgres,
			StorageIntegration:      &storageIntegrationId,
		}
		assertOptsValidAndSQLEquals(t, opts,
			`CREATE POSTGRES INSTANCE %s COMPUTE_FAMILY = 'STANDARD_S' STORAGE_SIZE_GB = 50 AUTHENTICATION_AUTHORITY = POSTGRES STORAGE_INTEGRATION = "my_integration"`,
			id.FullyQualifiedName())
	})

	t.Run("create with all identifier options", func(t *testing.T) {
		comment := random.Comment()
		networkPolicyId := NewAccountObjectIdentifier("my_policy")
		storageIntegrationId := NewAccountObjectIdentifier("my_integration")
		tagId := NewAccountObjectIdentifier("tag1")
		opts := &CreatePostgresInstanceOptions{
			name:                    id,
			ComputeFamily:           "STANDARD_S",
			StorageSizeGb:           50,
			AuthenticationAuthority: PostgresInstanceAuthenticationAuthorityPostgres,
			PostgresVersion:         Pointer(17),
			NetworkPolicy:           &networkPolicyId,
			HighAvailability:        Pointer(true),
			StorageIntegration:      &storageIntegrationId,
			PostgresSettings:        Pointer("{}"),
			Comment:                 &comment,
			Tag: []TagAssociation{
				{
					Name:  tagId,
					Value: "value1",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts,
			`CREATE POSTGRES INSTANCE %s COMPUTE_FAMILY = 'STANDARD_S' STORAGE_SIZE_GB = 50`+
				` AUTHENTICATION_AUTHORITY = POSTGRES POSTGRES_VERSION = 17 NETWORK_POLICY = "my_policy"`+
				` HIGH_AVAILABILITY = true STORAGE_INTEGRATION = "my_integration" POSTGRES_SETTINGS = '{}'`+
				` COMMENT = '%s' TAG (%s = 'value1')`,
			id.FullyQualifiedName(), comment, tagId.FullyQualifiedName())
	})
}

func TestPostgresInstances_Alter_IdentifierTypes(t *testing.T) {
	id := randomAccountObjectIdentifier()

	t.Run("alter set network policy as identifier", func(t *testing.T) {
		networkPolicyId := NewAccountObjectIdentifier("my_policy")
		opts := &AlterPostgresInstanceOptions{
			name: id,
			Set: &PostgresInstanceSet{
				NetworkPolicy: &networkPolicyId,
			},
		}
		assertOptsValidAndSQLEquals(t, opts,
			`ALTER POSTGRES INSTANCE %s SET NETWORK_POLICY = "my_policy"`,
			id.FullyQualifiedName())
	})

	t.Run("alter set storage integration as identifier", func(t *testing.T) {
		storageIntegrationId := NewAccountObjectIdentifier("my_integration")
		opts := &AlterPostgresInstanceOptions{
			name: id,
			Set: &PostgresInstanceSet{
				StorageIntegration: &storageIntegrationId,
			},
		}
		assertOptsValidAndSQLEquals(t, opts,
			`ALTER POSTGRES INSTANCE %s SET STORAGE_INTEGRATION = "my_integration"`,
			id.FullyQualifiedName())
	})

	t.Run("alter set both identifiers", func(t *testing.T) {
		comment := random.Comment()
		networkPolicyId := NewAccountObjectIdentifier("my_policy")
		storageIntegrationId := NewAccountObjectIdentifier("my_integration")
		auth := PostgresInstanceAuthenticationAuthorityPostgresOrSnowflake
		opts := &AlterPostgresInstanceOptions{
			name: id,
			Set: &PostgresInstanceSet{
				NetworkPolicy:           &networkPolicyId,
				AuthenticationAuthority: &auth,
				Comment:                 &comment,
				HighAvailability:        Pointer(true),
				ComputeFamily:           Pointer("STANDARD_M"),
				StorageSizeGb:           Pointer(100),
				StorageIntegration:      &storageIntegrationId,
				PostgresVersion:         Pointer(18),
				PostgresSettings:        Pointer("{}"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts,
			`ALTER POSTGRES INSTANCE %s SET NETWORK_POLICY = "my_policy" AUTHENTICATION_AUTHORITY = POSTGRES_OR_SNOWFLAKE`+
				` COMMENT = '%s' HIGH_AVAILABILITY = true COMPUTE_FAMILY = 'STANDARD_M' STORAGE_SIZE_GB = 100`+
				` STORAGE_INTEGRATION = "my_integration" POSTGRES_VERSION = 18 POSTGRES_SETTINGS = '{}'`,
			id.FullyQualifiedName(), comment)
	})
}

func TestPostgresInstances_ParseDetails_IdentifierTypes(t *testing.T) {
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

	t.Run("parse empty network policy returns nil", func(t *testing.T) {
		properties := []PostgresInstanceProperty{
			{Property: "name", Value: "test_instance"},
			{Property: "network_policy", Value: ""},
		}
		details, err := ParsePostgresInstanceDetails(properties)
		require.NoError(t, err)
		// After the change, empty string should not create an identifier
		assert.Nil(t, details.NetworkPolicy)
	})

	t.Run("parse empty storage integration returns nil", func(t *testing.T) {
		properties := []PostgresInstanceProperty{
			{Property: "name", Value: "test_instance"},
			{Property: "storage_integration", Value: ""},
		}
		details, err := ParsePostgresInstanceDetails(properties)
		require.NoError(t, err)
		// After the change, empty string should not create an identifier
		assert.Nil(t, details.StorageIntegration)
	})
}

func TestPostgresInstances_DtoBuilders_IdentifierTypes(t *testing.T) {
	t.Run("WithNetworkPolicy accepts AccountObjectIdentifier", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		networkPolicyId := NewAccountObjectIdentifier("my_policy")
		req := NewCreatePostgresInstanceRequest(
			id,
			"STANDARD_S",
			50,
			PostgresInstanceAuthenticationAuthorityPostgres,
		).WithNetworkPolicy(networkPolicyId)
		require.NotNil(t, req.NetworkPolicy)
		assert.Equal(t, networkPolicyId, *req.NetworkPolicy)
	})

	t.Run("WithStorageIntegration accepts AccountObjectIdentifier", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		storageIntegrationId := NewAccountObjectIdentifier("my_integration")
		req := NewCreatePostgresInstanceRequest(
			id,
			"STANDARD_S",
			50,
			PostgresInstanceAuthenticationAuthorityPostgres,
		).WithStorageIntegration(storageIntegrationId)
		require.NotNil(t, req.StorageIntegration)
		assert.Equal(t, storageIntegrationId, *req.StorageIntegration)
	})

	t.Run("alter set WithNetworkPolicy accepts AccountObjectIdentifier", func(t *testing.T) {
		networkPolicyId := NewAccountObjectIdentifier("my_policy")
		req := NewPostgresInstanceSetRequest().WithNetworkPolicy(networkPolicyId)
		require.NotNil(t, req.NetworkPolicy)
		assert.Equal(t, networkPolicyId, *req.NetworkPolicy)
	})

	t.Run("alter set WithStorageIntegration accepts AccountObjectIdentifier", func(t *testing.T) {
		storageIntegrationId := NewAccountObjectIdentifier("my_integration")
		req := NewPostgresInstanceSetRequest().WithStorageIntegration(storageIntegrationId)
		require.NotNil(t, req.StorageIntegration)
		assert.Equal(t, storageIntegrationId, *req.StorageIntegration)
	})
}
