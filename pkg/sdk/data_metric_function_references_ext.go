package sdk

import (
	"encoding/json"
	"fmt"
	"strings"
)

var _ convertibleRow[DataMetricFunctionReference] = new(dataMetricFunctionReferencesRow)

type DataMetricFunctionRefArgument struct {
	Domain string `json:"domain"`
	Id     string `json:"id"`
	Name   string `json:"name"`
}

func (row dataMetricFunctionReferencesRow) convert() (*DataMetricFunctionReference, error) {
	x := &DataMetricFunctionReference{
		MetricDatabaseName:    strings.Trim(row.MetricDatabaseName, `"`),
		MetricSchemaName:      strings.Trim(row.MetricSchemaName, `"`),
		MetricName:            strings.Trim(row.MetricName, `"`),
		ArgumentSignature:     row.MetricSignature,
		DataType:              row.MetricDataType,
		RefEntityDatabaseName: strings.Trim(row.RefEntityDatabaseName, `"`),
		RefEntitySchemaName:   strings.Trim(row.RefEntitySchemaName, `"`),
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

func NewGetForEntityDataMetricFunctionReferenceRequestCustom(
	refEntityName ObjectIdentifier,
	refEntityDomain DataMetricFunctionRefEntityDomainOption,
) *GetForEntityDataMetricFunctionReferenceRequest {
	return NewGetForEntityDataMetricFunctionReferenceRequest(
		NewdataMetricFunctionReferenceParametersRequest(
			NewdataMetricFunctionReferenceFunctionArgumentsRequest(
				[]ObjectIdentifier{refEntityName},
				&refEntityDomain,
			),
		),
	)
}
