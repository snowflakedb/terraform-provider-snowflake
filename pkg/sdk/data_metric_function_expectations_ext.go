package sdk

import "strings"

func (r dataMetricFunctionExpectationsRow) additionalConvert(result *DataMetricFunctionExpectation) error {
	result.MetricDatabaseName = strings.Trim(r.MetricDatabaseName, `"`)
	result.MetricSchemaName = strings.Trim(r.MetricSchemaName, `"`)
	result.MetricName = strings.Trim(r.MetricName, `"`)
	return nil
}

func NewGetForEntityDataMetricFunctionExpectationRequestCustom(
	refEntityName ObjectIdentifier,
	refEntityDomain DataMetricFunctionRefEntityDomainOption,
) *GetForEntityDataMetricFunctionExpectationRequest {
	return NewGetForEntityDataMetricFunctionExpectationRequest(
		NewdataMetricFunctionExpectationParametersRequest(
			NewdataMetricFunctionExpectationFunctionArgumentsRequest(
				[]ObjectIdentifier{refEntityName},
				&refEntityDomain,
			),
		),
	)
}
