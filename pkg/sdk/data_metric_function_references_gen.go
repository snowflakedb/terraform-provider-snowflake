package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type DataMetricFunctionReferences interface {
	GetForEntity(ctx context.Context, request *GetForEntityDataMetricFunctionReferenceRequest) ([]DataMetricFunctionReference, error)
}

// GetForEntityDataMetricFunctionReferenceOptions is based on https://docs.snowflake.com/en/sql-reference/functions/data_metric_function_references.
type GetForEntityDataMetricFunctionReferenceOptions struct {
	selectEverythingFrom bool                                   `ddl:"static" sql:"SELECT * FROM TABLE"`
	parameters           *dataMetricFunctionReferenceParameters `ddl:"list,parentheses,no_comma"`
}
type dataMetricFunctionReferenceParameters struct {
	functionFullyQualifiedName bool                                          `ddl:"static" sql:"SNOWFLAKE.INFORMATION_SCHEMA.DATA_METRIC_FUNCTION_REFERENCES"`
	arguments                  *dataMetricFunctionReferenceFunctionArguments `ddl:"list,parentheses"`
}
type dataMetricFunctionReferenceFunctionArguments struct {
	refEntityName   []ObjectIdentifier                       `ddl:"parameter,single_quotes,arrow_equals" sql:"REF_ENTITY_NAME"`
	refEntityDomain *DataMetricFunctionRefEntityDomainOption `ddl:"parameter,single_quotes,arrow_equals" sql:"REF_ENTITY_DOMAIN"`
}

type dataMetricFunctionReferencesRow struct {
	MetricDatabaseName string `db:"METRIC_DATABASE_NAME"`
	MetricSchemaName   string `db:"METRIC_SCHEMA_NAME"`
	MetricName         string `db:"METRIC_NAME"`
	ArgumentSignature  string `db:"METRIC_SIGNATURE"`
	DataType           string `db:"METRIC_DATA_TYPE"`
	RefDatabaseName    string `db:"REF_ENTITY_DATABASE_NAME"`
	RefSchemaName      string `db:"REF_ENTITY_SCHEMA_NAME"`
	RefEntityName      string `db:"REF_ENTITY_NAME"`
	RefEntityDomain    string `db:"REF_ENTITY_DOMAIN"`
	RefArguments       string `db:"REF_ARGUMENTS"`
	RefId              string `db:"REF_ID"`
	Schedule           string `db:"SCHEDULE"`
	ScheduleStatus     string `db:"SCHEDULE_STATUS"`
}

type DataMetricFunctionRefArgument struct {
	Domain string `json:"domain"`
	Id     string `json:"id"`
	Name   string `json:"name"`
}
type DataMetricFunctionReference struct {
	MetricDatabaseName    string
	MetricSchemaName      string
	MetricName            string
	ArgumentSignature     string
	DataType              string
	RefEntityDatabaseName string
	RefEntitySchemaName   string
	RefEntityName         string
	RefEntityDomain       string
	RefArguments          []DataMetricFunctionRefArgument
	RefId                 string
	Schedule              string
	ScheduleStatus        string
}

func (row dataMetricFunctionReferencesRow) convert() (*DataMetricFunctionReference, error) {
	x := &DataMetricFunctionReference{
		MetricDatabaseName:    strings.Trim(row.MetricDatabaseName, `"`),
		MetricSchemaName:      strings.Trim(row.MetricSchemaName, `"`),
		MetricName:            strings.Trim(row.MetricName, `"`),
		ArgumentSignature:     row.ArgumentSignature,
		DataType:              row.DataType,
		RefEntityDatabaseName: strings.Trim(row.RefDatabaseName, `"`),
		RefEntitySchemaName:   strings.Trim(row.RefSchemaName, `"`),
		RefEntityName:         strings.Trim(row.RefEntityName, `"`),
		RefEntityDomain:       row.RefEntityDomain,
		RefId:                 row.RefId,
		Schedule:              row.Schedule,
		ScheduleStatus:        row.ScheduleStatus,
	}
	err := json.Unmarshal([]byte(row.RefArguments), &x.RefArguments)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data metric function reference arguments: %w", err)
	}
	return x, nil
}
