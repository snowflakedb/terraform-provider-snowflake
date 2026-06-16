package sdk

import "strings"

func (r dataMetricFunctionReferencesRow) additionalConvert(result *DataMetricFunctionReference) error {
	result.MetricDatabaseName = strings.Trim(r.MetricDatabaseName, `"`)
	result.MetricSchemaName = strings.Trim(r.MetricSchemaName, `"`)
	result.MetricName = strings.Trim(r.MetricName, `"`)
	result.RefEntityDatabaseName = strings.Trim(r.RefEntityDatabaseName, `"`)
	result.RefEntitySchemaName = strings.Trim(r.RefEntitySchemaName, `"`)
	result.RefEntityName = strings.Trim(r.RefEntityName, `"`)
	return nil
}

type DataMetricFunctionRefArgument struct {
	Domain string `json:"domain"`
	Id     string `json:"id"`
	Name   string `json:"name"`
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
