package parameterdefs

import "strings"

// ParameterLevel identifies a catalog consumer set used by ParameterDefsForLevel. It is not
// the level returned by Snowflake from SHOW PARAMETERS.
//
// TODO [SNOW-4160077]: Remove or relocate this temporary public API when the
// parameter generator becomes the single source of truth for SDK, resource,
// and data source parameter handling.
type ParameterLevel string

const (
	ParameterLevelAccount              ParameterLevel = "ACCOUNT"
	ParameterLevelDatabase             ParameterLevel = "DATABASE"
	ParameterLevelSchema               ParameterLevel = "SCHEMA"
	ParameterLevelTable                ParameterLevel = "TABLE"
	ParameterLevelIcebergTable         ParameterLevel = "ICEBERG_TABLE"
	ParameterLevelTask                 ParameterLevel = "TASK"
	ParameterLevelFunction             ParameterLevel = "FUNCTION"
	ParameterLevelProcedure            ParameterLevel = "PROCEDURE"
	ParameterLevelProject              ParameterLevel = "PROJECT"
	ParameterLevelService              ParameterLevel = "SERVICE"
	ParameterLevelUser                 ParameterLevel = "USER"
	ParameterLevelSession              ParameterLevel = "SESSION"
	ParameterLevelOpenflowDeployment   ParameterLevel = "OPENFLOW_DEPLOYMENT"
	ParameterLevelWarehouse            ParameterLevel = "WAREHOUSE"
	ParameterLevelWarehouseAdaptive    ParameterLevel = "WAREHOUSE_ADAPTIVE"
	ParameterLevelWarehouseInteractive ParameterLevel = "WAREHOUSE_INTERACTIVE"
	// ParameterLevelAccountExt has no Snowflake counterpart. It marks the broader set of account
	// parameters consumed by the generic account_parameter resource, while ParameterLevelAccount
	// marks the narrower set exposed as typed fields on current_account.
	ParameterLevelAccountExt ParameterLevel = "ACCOUNT_EXT"
	// ParameterLevelHybridTable has no Snowflake counterpart. It marks the narrower subset of
	// ParameterLevelTable parameters that are actually settable on HYBRID TABLE objects (verified
	// against ALTER HYBRID TABLE SET/UNSET and CREATE HYBRID TABLE support), mirroring the
	// ParameterLevelAccount/ParameterLevelAccountExt narrow-vs-broad split above.
	ParameterLevelHybridTable ParameterLevel = "HYBRID_TABLE"
)

// ParameterDef is one parameter declaration consumed by SDK generator definitions.
type ParameterDef struct {
	SqlName      string
	Kind         string
	Levels       []ParameterLevel
	Description  string
	DefaultValue string
	DefaultLevel string
}

func (p ParameterDef) FieldName() string {
	return strings.ToLower(p.SqlName)
}
