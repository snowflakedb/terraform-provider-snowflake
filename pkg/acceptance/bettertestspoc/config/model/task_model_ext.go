package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func TaskWithId(resourceName string, id sdk.SchemaObjectIdentifier, started bool, sqlStatement string) *TaskModel {
	t := &TaskModel{ResourceModelMeta: config.Meta(resourceName, resources.Task)}
	t.WithDatabase(id.DatabaseName())
	t.WithSchema(id.SchemaName())
	t.WithName(id.Name())
	t.WithStarted(started)
	t.WithSqlStatement(sqlStatement)
	return t
}

func (t *TaskModel) WithBinaryInputFormatEnum(binaryInputFormat sdk.BinaryInputFormat) *TaskModel {
	t.BinaryInputFormat = tfconfig.StringVariable(string(binaryInputFormat))
	return t
}

func (t *TaskModel) WithBinaryOutputFormatEnum(binaryOutputFormat sdk.BinaryOutputFormat) *TaskModel {
	t.BinaryOutputFormat = tfconfig.StringVariable(string(binaryOutputFormat))
	return t
}

func (t *TaskModel) WithClientTimestampTypeMappingEnum(clientTimestampTypeMapping sdk.ClientTimestampTypeMapping) *TaskModel {
	t.ClientTimestampTypeMapping = tfconfig.StringVariable(string(clientTimestampTypeMapping))
	return t
}

func (t *TaskModel) WithGeographyOutputFormatEnum(geographyOutputFormat sdk.GeographyOutputFormat) *TaskModel {
	t.GeographyOutputFormat = tfconfig.StringVariable(string(geographyOutputFormat))
	return t
}

func (t *TaskModel) WithGeometryOutputFormatEnum(geometryOutputFormat sdk.GeometryOutputFormat) *TaskModel {
	t.GeometryOutputFormat = tfconfig.StringVariable(string(geometryOutputFormat))
	return t
}

func (t *TaskModel) WithLogLevelEnum(logLevel sdk.LogLevel) *TaskModel {
	t.LogLevel = tfconfig.StringVariable(string(logLevel))
	return t
}

func (t *TaskModel) WithTimestampTypeMappingEnum(timestampTypeMapping sdk.TimestampTypeMapping) *TaskModel {
	t.TimestampTypeMapping = tfconfig.StringVariable(string(timestampTypeMapping))
	return t
}

func (t *TaskModel) WithTraceLevelEnum(traceLevel sdk.TraceLevel) *TaskModel {
	t.TraceLevel = tfconfig.StringVariable(string(traceLevel))
	return t
}

func (t *TaskModel) WithTransactionDefaultIsolationLevelEnum(transactionDefaultIsolationLevel sdk.TransactionDefaultIsolationLevel) *TaskModel {
	t.TransactionDefaultIsolationLevel = tfconfig.StringVariable(string(transactionDefaultIsolationLevel))
	return t
}

func (t *TaskModel) WithUnsupportedDdlActionEnum(unsupportedDdlAction sdk.UnsupportedDDLAction) *TaskModel {
	t.UnsupportedDdlAction = tfconfig.StringVariable(string(unsupportedDdlAction))
	return t
}

func (t *TaskModel) WithUserTaskManagedInitialWarehouseSizeEnum(warehouseSize sdk.WarehouseSize) *TaskModel {
	t.UserTaskManagedInitialWarehouseSize = tfconfig.StringVariable(string(warehouseSize))
	return t
}

func (t *TaskModel) WithServerlessTaskMinStatementSizeEnum(warehouseSize sdk.WarehouseSize) *TaskModel {
	t.ServerlessTaskMinStatementSize = tfconfig.StringVariable(string(warehouseSize))
	return t
}

func (t *TaskModel) WithServerlessTaskMaxStatementSizeEnum(warehouseSize sdk.WarehouseSize) *TaskModel {
	t.ServerlessTaskMaxStatementSize = tfconfig.StringVariable(string(warehouseSize))
	return t
}

func (t *TaskModel) WithScheduleMinutes(minutes int) *TaskModel {
	t.Schedule = tfconfig.MapVariable(map[string]tfconfig.Variable{
		"minutes": tfconfig.IntegerVariable(minutes),
	})
	return t
}

func (t *TaskModel) WithScheduleCron(cron string) *TaskModel {
	t.Schedule = tfconfig.MapVariable(map[string]tfconfig.Variable{
		"cron": tfconfig.StringVariable(cron),
	})
	return t
}
