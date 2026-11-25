package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCommaSeparatedFailoverGroupObjectTypeArray(t *testing.T) {
	t.Run("known values, any case, with spaces", func(t *testing.T) {
		got, err := ParseCommaSeparatedFailoverGroupObjectTypeArray("databases, External Volumes ,ROLES")
		require.NoError(t, err)
		assert.Equal(t, []FailoverGroupObjectType{FailoverGroupObjectTypeDatabases, FailoverGroupObjectTypeExternalVolumes, FailoverGroupObjectTypeRoles}, got)
	})

	t.Run("PARAMETERS is ACCOUNT PARAMETERS", func(t *testing.T) {
		got, err := ParseCommaSeparatedFailoverGroupObjectTypeArray("parameters,DATABASES")
		require.NoError(t, err)
		assert.Equal(t, []FailoverGroupObjectType{FailoverGroupObjectTypeAccountParameters, FailoverGroupObjectTypeDatabases}, got)
	})

	t.Run("an unknown value is kept in upper case", func(t *testing.T) {
		got, err := ParseCommaSeparatedFailoverGroupObjectTypeArray("databases,new thing")
		require.NoError(t, err)
		assert.Equal(t, []FailoverGroupObjectType{FailoverGroupObjectTypeDatabases, FailoverGroupObjectType("NEW THING")}, got)
	})

	t.Run("empty value", func(t *testing.T) {
		got, err := ParseCommaSeparatedFailoverGroupObjectTypeArray(" , ")
		require.NoError(t, err)
		assert.Empty(t, got)
	})
}

func TestParseCommaSeparatedFailoverGroupIntegrationTypeArray(t *testing.T) {
	t.Run("short SHOW tokens", func(t *testing.T) {
		got, err := ParseCommaSeparatedFailoverGroupIntegrationTypeArray("API, external_access,NOTIFICATION,security,STORAGE")
		require.NoError(t, err)
		assert.Equal(t, []FailoverGroupIntegrationType{
			FailoverGroupIntegrationTypeApiIntegrations,
			FailoverGroupIntegrationTypeExternalAccessIntegrations,
			FailoverGroupIntegrationTypeNotificationIntegrations,
			FailoverGroupIntegrationTypeSecurityIntegrations,
			FailoverGroupIntegrationTypeStorageIntegrations,
		}, got)
	})

	t.Run("an unknown value is kept in upper case", func(t *testing.T) {
		got, err := ParseCommaSeparatedFailoverGroupIntegrationTypeArray("api,new_kind")
		require.NoError(t, err)
		assert.Equal(t, []FailoverGroupIntegrationType{FailoverGroupIntegrationTypeApiIntegrations, FailoverGroupIntegrationType("NEW KIND INTEGRATIONS")}, got)
	})
}
