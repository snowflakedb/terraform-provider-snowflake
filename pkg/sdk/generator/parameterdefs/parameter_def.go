package parameterdefs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

type (
	ParameterDef   = ssot.ParameterDef
	ParameterLevel = ssot.ParameterLevel
)

const (
	ParameterLevelAccount              = ssot.ParameterLevelAccount
	ParameterLevelDatabase             = ssot.ParameterLevelDatabase
	ParameterLevelSchema               = ssot.ParameterLevelSchema
	ParameterLevelTable                = ssot.ParameterLevelTable
	ParameterLevelIcebergTable         = ssot.ParameterLevelIcebergTable
	ParameterLevelTask                 = ssot.ParameterLevelTask
	ParameterLevelFunction             = ssot.ParameterLevelFunction
	ParameterLevelProcedure            = ssot.ParameterLevelProcedure
	ParameterLevelProject              = ssot.ParameterLevelProject
	ParameterLevelService              = ssot.ParameterLevelService
	ParameterLevelUser                 = ssot.ParameterLevelUser
	ParameterLevelSession              = ssot.ParameterLevelSession
	ParameterLevelOpenflowDeployment   = ssot.ParameterLevelOpenflowDeployment
	ParameterLevelWarehouse            = ssot.ParameterLevelWarehouse
	ParameterLevelWarehouseAdaptive    = ssot.ParameterLevelWarehouseAdaptive
	ParameterLevelWarehouseInteractive = ssot.ParameterLevelWarehouseInteractive
	ParameterLevelAccountExt           = ssot.ParameterLevelAccountExt
	ParameterLevelHybridTable          = ssot.ParameterLevelHybridTable
)
