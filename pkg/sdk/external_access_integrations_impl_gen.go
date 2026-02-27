package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ ExternalAccessIntegrations = (*externalAccessIntegrations)(nil)

var (
	_ convertibleRow[ExternalAccessIntegration]         = new(showExternalAccessIntegrationsDbRow)
	_ convertibleRow[ExternalAccessIntegrationProperty] = new(descExternalAccessIntegrationsDbRow)
)

type externalAccessIntegrations struct {
	client *Client
}

func (v *externalAccessIntegrations) Create(ctx context.Context, request *CreateExternalAccessIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *externalAccessIntegrations) Alter(ctx context.Context, request *AlterExternalAccessIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *externalAccessIntegrations) Drop(ctx context.Context, request *DropExternalAccessIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *externalAccessIntegrations) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error {
		return v.Drop(ctx, NewDropExternalAccessIntegrationRequest(id).WithIfExists(true))
	}, ctx, id)
}

func (v *externalAccessIntegrations) Show(ctx context.Context, request *ShowExternalAccessIntegrationRequest) ([]ExternalAccessIntegration, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[showExternalAccessIntegrationsDbRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[showExternalAccessIntegrationsDbRow, ExternalAccessIntegration](dbRows)
}

func (v *externalAccessIntegrations) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ExternalAccessIntegration, error) {
	request := NewShowExternalAccessIntegrationRequest().
		WithLike(Like{Pattern: String(id.Name())})
	integrations, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(integrations, func(r ExternalAccessIntegration) bool { return r.Name == id.Name() })
}

func (v *externalAccessIntegrations) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*ExternalAccessIntegration, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *externalAccessIntegrations) Describe(ctx context.Context, id AccountObjectIdentifier) ([]ExternalAccessIntegrationProperty, error) {
	opts := &DescribeExternalAccessIntegrationOptions{
		name: id,
	}
	rows, err := validateAndQuery[descExternalAccessIntegrationsDbRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[descExternalAccessIntegrationsDbRow, ExternalAccessIntegrationProperty](rows)
}

func (r *CreateExternalAccessIntegrationRequest) toOpts() *CreateExternalAccessIntegrationOptions {
	opts := &CreateExternalAccessIntegrationOptions{
		OrReplace:                    r.OrReplace,
		IfNotExists:                  r.IfNotExists,
		name:                         r.name,
		AllowedNetworkRules:          r.AllowedNetworkRules,
		AllowedAuthenticationSecrets: r.AllowedAuthenticationSecrets,
		Enabled:                      r.Enabled,
		Comment:                      r.Comment,
	}
	return opts
}

func (r *AlterExternalAccessIntegrationRequest) toOpts() *AlterExternalAccessIntegrationOptions {
	opts := &AlterExternalAccessIntegrationOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	if r.Set != nil {
		opts.Set = &ExternalAccessIntegrationSet{
			AllowedNetworkRules:          r.Set.AllowedNetworkRules,
			AllowedAuthenticationSecrets: r.Set.AllowedAuthenticationSecrets,
			Enabled:                      r.Set.Enabled,
			Comment:                      r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &ExternalAccessIntegrationUnset{
			AllowedAuthenticationSecrets: r.Unset.AllowedAuthenticationSecrets,
			Comment:                      r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropExternalAccessIntegrationRequest) toOpts() *DropExternalAccessIntegrationOptions {
	return &DropExternalAccessIntegrationOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
}

func (r *ShowExternalAccessIntegrationRequest) toOpts() *ShowExternalAccessIntegrationOptions {
	return &ShowExternalAccessIntegrationOptions{
		Like: r.Like,
	}
}

func (r *DescribeExternalAccessIntegrationRequest) toOpts() *DescribeExternalAccessIntegrationOptions {
	return &DescribeExternalAccessIntegrationOptions{
		name: r.name,
	}
}

func (r showExternalAccessIntegrationsDbRow) convert() (*ExternalAccessIntegration, error) {
	s := &ExternalAccessIntegration{
		Name:      r.Name,
		IntType:   r.Type,
		Category:  r.Category,
		Enabled:   r.Enabled,
		CreatedOn: r.CreatedOn,
	}
	if r.Comment.Valid {
		s.Comment = r.Comment.String
	}
	return s, nil
}

func (r descExternalAccessIntegrationsDbRow) convert() (*ExternalAccessIntegrationProperty, error) {
	return &ExternalAccessIntegrationProperty{
		Name:    r.Property,
		Type:    r.PropertyType,
		Value:   r.PropertyValue,
		Default: r.PropertyDefault,
	}, nil
}
