package sdk

import (
	"context"
)

var _ DataMetricFunctionReferences = (*dataMetricFunctionReferences)(nil)
var _ convertibleRow[DataMetricFunctionReference] = new(dataMetricFunctionReferencesRow)

type dataMetricFunctionReferences struct {
	client *Client
}

func (v *dataMetricFunctionReferences) GetForEntity(ctx context.Context, request *GetForEntityDataMetricFunctionReferenceRequest) ([]DataMetricFunctionReference, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[dataMetricFunctionReferencesRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[dataMetricFunctionReferencesRow, DataMetricFunctionReference](dbRows)
}

func (r *GetForEntityDataMetricFunctionReferenceRequest) toOpts() *GetForEntityDataMetricFunctionReferenceOptions {
	opts := &GetForEntityDataMetricFunctionReferenceOptions{
		parameters: &dataMetricFunctionReferenceParameters{
			arguments: &dataMetricFunctionReferenceFunctionArguments{
				refEntityName:   []ObjectIdentifier{r.refEntityName},
				refEntityDomain: Pointer(r.RefEntityDomain),
			},
		},
	}
	return opts
}
