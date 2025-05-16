package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
)

var _ Streamlits = (*streamlits)(nil)

type streamlits struct {
	client *Client
}

func (v *streamlits) Create(ctx context.Context, request *CreateStreamlitRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *streamlits) Alter(ctx context.Context, request *AlterStreamlitRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *streamlits) Drop(ctx context.Context, request *DropStreamlitRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *streamlits) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropStreamlitRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *streamlits) Show(ctx context.Context, request *ShowStreamlitRequest) ([]Streamlit, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[streamlitsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[streamlitsRow, Streamlit](dbRows)
	return resultList, nil
}

func (v *streamlits) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Streamlit, error) {
	request := NewShowStreamlitRequest().
		WithLike(Like{Pattern: String(id.Name())}).
		WithIn(In{Schema: id.SchemaId()})
	streamlits, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(streamlits, func(r Streamlit) bool { return r.Name == id.Name() })
}

func (v *streamlits) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*Streamlit, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *streamlits) Describe(ctx context.Context, id SchemaObjectIdentifier) (*StreamlitDetail, error) {
	opts := &DescribeStreamlitOptions{
		name: id,
	}
	result, err := validateAndQueryOne[streamlitsDetailRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert(), nil
}

func (r *CreateStreamlitRequest) toOpts() *CreateStreamlitOptions {
	opts := &CreateStreamlitOptions{
		OrReplace:      r.OrReplace,
		IfNotExists:    r.IfNotExists,
		name:           r.name,
		RootLocation:   r.RootLocation,
		MainFile:       r.MainFile,
		QueryWarehouse: r.QueryWarehouse,

		Title:   r.Title,
		Comment: r.Comment,
	}

	if r.ExternalAccessIntegrations != nil {
		opts.ExternalAccessIntegrations = &ExternalAccessIntegrations{
			ExternalAccessIntegrations: r.ExternalAccessIntegrations.ExternalAccessIntegrations,
		}
	}

	return opts
}

func (r *AlterStreamlitRequest) toOpts() *AlterStreamlitOptions {
	opts := &AlterStreamlitOptions{
		IfExists: r.IfExists,
		name:     r.name,

		RenameTo: r.RenameTo,
	}

	if r.Set != nil {
		opts.Set = &StreamlitSet{
			RootLocation:   r.Set.RootLocation,
			MainFile:       r.Set.MainFile,
			QueryWarehouse: r.Set.QueryWarehouse,
			Comment:        r.Set.Comment,
			Title:          r.Set.Title,
		}

		if r.Set.ExternalAccessIntegrations != nil {
			opts.Set.ExternalAccessIntegrations = &ExternalAccessIntegrations{
				ExternalAccessIntegrations: r.Set.ExternalAccessIntegrations.ExternalAccessIntegrations,
			}
		}
	}

	if r.Unset != nil {
		opts.Unset = &StreamlitUnset{
			QueryWarehouse: r.Unset.QueryWarehouse,
			Comment:        r.Unset.Comment,
			Title:          r.Unset.Title,
		}
	}

	return opts
}

func (r *DropStreamlitRequest) toOpts() *DropStreamlitOptions {
	opts := &DropStreamlitOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowStreamlitRequest) toOpts() *ShowStreamlitOptions {
	opts := &ShowStreamlitOptions{
		Terse: r.Terse,
		Like:  r.Like,
		In:    r.In,
		Limit: r.Limit,
	}
	return opts
}

func (r streamlitsRow) convert() *Streamlit {
	e := &Streamlit{
		CreatedOn:     r.CreatedOn,
		Name:          r.Name,
		DatabaseName:  r.DatabaseName,
		SchemaName:    r.SchemaName,
		Owner:         r.Owner,
		UrlId:         r.UrlId,
		OwnerRoleType: r.OwnerRoleType,
	}
	if r.Title.Valid {
		e.Title = r.Title.String
	}
	if r.Comment.Valid {
		e.Comment = r.Comment.String
	}
	if r.QueryWarehouse.Valid {
		e.QueryWarehouse = r.QueryWarehouse.String
	}
	return e
}

func (r *DescribeStreamlitRequest) toOpts() *DescribeStreamlitOptions {
	opts := &DescribeStreamlitOptions{
		name: r.name,
	}
	return opts
}

func (r streamlitsDetailRow) convert() *StreamlitDetail {
	e := &StreamlitDetail{
		Name:                       r.Name,
		RootLocation:               r.RootLocation,
		MainFile:                   r.MainFile,
		UrlId:                      r.UrlId,
		DefaultPackages:            r.DefaultPackages,
		UserPackages:               ParseCommaSeparatedStringArray(r.UserPackages, false),
		ImportUrls:                 ParseCommaSeparatedStringArray(r.ImportUrls, false),
		ExternalAccessIntegrations: ParseCommaSeparatedStringArray(r.ExternalAccessIntegrations, false),
		ExternalAccessSecrets:      r.ExternalAccessSecrets,
	}
	if r.Title.Valid {
		e.Title = r.Title.String
	}
	if r.QueryWarehouse.Valid {
		e.QueryWarehouse = r.QueryWarehouse.String
	}
	integrationsRaw := ParseCommaSeparatedStringArray(r.ExternalAccessIntegrations, false)
	externalAccessIntegrations := make([]string, len(integrationsRaw))
	for i, v := range integrationsRaw {
		externalAccessIntegrations[i] = NewObjectIdentifierFromFullyQualifiedName(v).Name()
	}
	e.ExternalAccessIntegrations = externalAccessIntegrations
	return e
}
