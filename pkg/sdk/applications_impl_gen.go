package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ Applications = (*applications)(nil)
var _ convertibleRow[Application] = new(applicationRow)
var _ convertibleRow[ApplicationProperty] = new(applicationPropertyRow)

type applications struct {
	client *Client
}

func (v *applications) Create(ctx context.Context, request *CreateApplicationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *applications) Drop(ctx context.Context, request *DropApplicationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *applications) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropApplicationRequest(id).WithIfExists(Bool(true))) }, ctx, id)
}

func (v *applications) Alter(ctx context.Context, request *AlterApplicationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *applications) Show(ctx context.Context, request *ShowApplicationRequest) ([]Application, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[applicationRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[applicationRow, Application](dbRows)
}

func (v *applications) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Application, error) {
	request := NewShowApplicationRequest().
		WithLike(Like{Pattern: String(id.Name())})
	applications, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(applications, func(r Application) bool { return r.Name == id.Name() })
}

func (v *applications) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Application, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *applications) Describe(ctx context.Context, id AccountObjectIdentifier) ([]ApplicationProperty, error) {
	opts := &DescribeApplicationOptions{
		name: id,
	}
	rows, err := validateAndQuery[applicationPropertyRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[applicationPropertyRow, ApplicationProperty](rows)
}

func (r *CreateApplicationRequest) toOpts() *CreateApplicationOptions {
	opts := &CreateApplicationOptions{
		name:        r.name,
		PackageName: r.PackageName,

		DebugMode: r.DebugMode,
		Comment:   r.Comment,
		Tag:       r.Tag,
	}
	if r.Version != nil {
		opts.Version = &ApplicationVersion{
			VersionDirectory: r.Version.VersionDirectory,
		}
		if r.Version.VersionAndPatch != nil {
			opts.Version.VersionAndPatch = &VersionAndPatch{
				Version: r.Version.VersionAndPatch.Version,
				Patch:   r.Version.VersionAndPatch.Patch,
			}
		}
	}
	return opts
}

func (r *DropApplicationRequest) toOpts() *DropApplicationOptions {
	opts := &DropApplicationOptions{
		IfExists: r.IfExists,
		name:     r.name,
		Cascade:  r.Cascade,
	}
	return opts
}

func (r *AlterApplicationRequest) toOpts() *AlterApplicationOptions {
	opts := &AlterApplicationOptions{
		IfExists: r.IfExists,
		name:     r.name,

		Upgrade: r.Upgrade,

		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	if r.Set != nil {
		opts.Set = &ApplicationSet{
			Comment:                 r.Set.Comment,
			ShareEventsWithProvider: r.Set.ShareEventsWithProvider,
			DebugMode:               r.Set.DebugMode,
		}
	}
	if r.Unset != nil {
		opts.Unset = &ApplicationUnset{
			Comment:                 r.Unset.Comment,
			ShareEventsWithProvider: r.Unset.ShareEventsWithProvider,
			DebugMode:               r.Unset.DebugMode,
		}
	}
	if r.UpgradeVersion != nil {
		opts.UpgradeVersion = &ApplicationVersion{
			VersionDirectory: r.UpgradeVersion.VersionDirectory,
		}
		if r.UpgradeVersion.VersionAndPatch != nil {
			opts.UpgradeVersion.VersionAndPatch = &VersionAndPatch{
				Version: r.UpgradeVersion.VersionAndPatch.Version,
				Patch:   r.UpgradeVersion.VersionAndPatch.Patch,
			}
		}
	}
	if r.UnsetReferences != nil {
		opts.UnsetReferences = &ApplicationReferences{}
		if r.UnsetReferences.References != nil {
			s := make([]ApplicationReference, len(r.UnsetReferences.References))
			for i, v := range r.UnsetReferences.References {
				s[i] = ApplicationReference(v)
			}
			opts.UnsetReferences.References = s
		}
	}
	return opts
}

func (r *ShowApplicationRequest) toOpts() *ShowApplicationOptions {
	opts := &ShowApplicationOptions{
		Like:       r.Like,
		StartsWith: r.StartsWith,
		Limit:      r.Limit,
	}
	return opts
}

func (r applicationRow) convert() (*Application, error) {
	return &Application{
		CreatedOn:     r.CreatedOn,
		Name:          r.Name,
		IsDefault:     r.IsDefault == "Y",
		IsCurrent:     r.IsCurrent == "Y",
		SourceType:    r.SourceType,
		Source:        r.Source,
		Owner:         r.Owner,
		Comment:       r.Comment,
		Version:       r.Version,
		Label:         r.Label,
		Patch:         r.Patch,
		Options:       r.Options,
		RetentionTime: r.RetentionTime,
	}, nil
}

func (r *DescribeApplicationRequest) toOpts() *DescribeApplicationOptions {
	opts := &DescribeApplicationOptions{
		name: r.name,
	}
	return opts
}

func (r applicationPropertyRow) convert() (*ApplicationProperty, error) {
	e := &ApplicationProperty{
		Property: r.Property,
	}
	if r.Value.Valid {
		e.Value = r.Value.String
	}
	return e, nil
}
