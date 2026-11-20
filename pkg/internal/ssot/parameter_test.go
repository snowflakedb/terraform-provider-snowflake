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

func TestFilterResourceParameters(t *testing.T) {
	retention := ParameterDef{SqlName: "DATA_RETENTION_TIME_IN_DAYS"}
	collation := ParameterDef{SqlName: "DEFAULT_DDL_COLLATION"}
	catalog := ParameterDef{SqlName: "CATALOG"}
	all := []ParameterDef{catalog, retention, collation}

	require.Equal(t, all, FilterResourceParameters(all, nil))
	require.Equal(t, all, FilterResourceParameters(all, []ParameterDef{}))
	require.Equal(t, []ParameterDef{retention, collation}, FilterResourceParameters(all, []ParameterDef{collation, retention}))
}
