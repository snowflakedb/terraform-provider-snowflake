package sdk

import (
	"context"
	"slices"
)

func (s *Schema) IsTransient() bool {
	if s.Options == nil {
		return false
	}
	return slices.Contains(ParseCommaSeparatedStringArray(*s.Options, false), "TRANSIENT")
}

func (s *Schema) IsManagedAccess() bool {
	if s.Options == nil {
		return false
	}
	return slices.Contains(ParseCommaSeparatedStringArray(*s.Options, false), "MANAGED ACCESS")
}

func (v *schemas) Use(ctx context.Context, id DatabaseObjectIdentifier) error {
	return v.client.Sessions.UseSchema(ctx, id)
}

func (v *schemas) ShowParameters(ctx context.Context, id DatabaseObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Schema: id,
		},
	})
}
