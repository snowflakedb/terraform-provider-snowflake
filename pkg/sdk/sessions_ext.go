package sdk

import "context"

func (v *sessions) ShowParameters(ctx context.Context, opts *ShowParametersOptions) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, opts)
}
