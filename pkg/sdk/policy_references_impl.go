package sdk

import "context"

var _ PolicyReferences = new(policyReference)
var _ convertibleRow[PolicyReference] = new(policyReferenceDBRow)

type policyReference struct {
	client *Client
}

func (v *policyReference) GetForEntity(ctx context.Context, request *GetForEntityPolicyReferenceRequest) ([]PolicyReference, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[policyReferenceDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRowsErr[policyReferenceDBRow, PolicyReference](dbRows)
}
