package ssot

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

// ParameterLevel identifies a catalog consumer set used by ParameterDefsForLevel. It is not
// the level returned by Snowflake from SHOW PARAMETERS.
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

// AccountLevels is where an account parameter can be set: the account, including both provider
// consumer sets (typed current_account fields and the generic account_parameter resource).
var AccountLevels = []ParameterLevel{ParameterLevelAccount, ParameterLevelAccountExt}

// SessionLevels is where a session parameter can be set: account, user, and session.
var SessionLevels = append(append([]ParameterLevel{}, AccountLevels...), ParameterLevelSession, ParameterLevelUser)

// ObjectLevels is where an object parameter can be set: the account levels plus the given object
// levels. Object parameters can apply to more than one object type.
func ObjectLevels(objects ...ParameterLevel) []ParameterLevel {
	return append(append([]ParameterLevel{}, AccountLevels...), objects...)
}

// ParameterDef is one parameter declaration consumed by generators, resources, and data sources.
type ParameterDef struct {
	SqlName      string
	Kind         string
	Levels       []ParameterLevel
	Description  string
	DefaultValue string
	DefaultLevel string
}

// FieldName is the Terraform attribute name: the SQL name in lowercase.
func (p ParameterDef) FieldName() string {
	return strings.ToLower(p.SqlName)
}

// StructFieldName is the exported SDK struct field name. It uses the same English
// title casing as the SDK generator, so TIMESTAMP_DAY_IS_ALWAYS_24H is
// TimestampDayIsAlways24H and UNSUPPORTED_DDL_ACTION is UnsupportedDdlAction.
func (p ParameterDef) StructFieldName() string {
	return genhelpers.SqlToGoName(p.SqlName, true)
}
