package sdk

import "context"

// ShowParameters returns the parameters visible at the TABLE level for the given hybrid table.
// Mirrors pkg/sdk/functions_ext.go:155 (ParametersIn.Function) with ParametersIn.Table.
func (v *hybridTables) ShowParameters(ctx context.Context, id SchemaObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Table: id,
		},
	})
}
