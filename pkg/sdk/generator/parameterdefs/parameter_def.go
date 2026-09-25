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
	ParameterLevelAccount   ParameterLevel = "ACCOUNT"
	ParameterLevelDatabase  ParameterLevel = "DATABASE"
	ParameterLevelSchema    ParameterLevel = "SCHEMA"
	ParameterLevelTable     ParameterLevel = "TABLE"
	ParameterLevelTask      ParameterLevel = "TASK"
	ParameterLevelFunction  ParameterLevel = "FUNCTION"
	ParameterLevelProcedure ParameterLevel = "PROCEDURE"
	ParameterLevelProject   ParameterLevel = "PROJECT"
	ParameterLevelService   ParameterLevel = "SERVICE"
	ParameterLevelUser      ParameterLevel = "USER"
	ParameterLevelSession   ParameterLevel = "SESSION"
	// ParameterLevelAccountExt has no Snowflake counterpart. It marks the broader set of account
	// parameters consumed by the generic account_parameter resource, while ParameterLevelAccount
	// marks the narrower set exposed as typed fields on current_account.
	ParameterLevelAccountExt ParameterLevel = "ACCOUNT_EXT"
)

// ParameterDef is one parameter declaration consumed by SDK generator definitions.
type ParameterDef struct {
	SqlName     string
	Kind        string
	Levels      []ParameterLevel
	Description string
}

func (p ParameterDef) FieldName() string {
	return strings.ToLower(p.SqlName)
}
