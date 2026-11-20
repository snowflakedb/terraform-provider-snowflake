package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FilterResourceParameters_emptyAllowlistReturnsCatalog(t *testing.T) {
	catalog := []ParameterDef{
		{SqlName: "DATA_RETENTION_TIME_IN_DAYS"},
		{SqlName: "CATALOG"},
	}
	require.Equal(t, catalog, FilterResourceParameters(catalog, nil))
}

func Test_FilterResourceParameters_keepsAllowlistInCatalogOrder(t *testing.T) {
	catalog := []ParameterDef{
		{SqlName: "CATALOG"},
		{SqlName: "DATA_RETENTION_TIME_IN_DAYS"},
		{SqlName: "LOG_LEVEL"},
	}
	got := FilterResourceParameters(catalog, []ParameterDef{{SqlName: "DATA_RETENTION_TIME_IN_DAYS"}})
	require.Equal(t, []ParameterDef{{SqlName: "DATA_RETENTION_TIME_IN_DAYS"}}, got)
}
