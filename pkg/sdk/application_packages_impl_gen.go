package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
)

var _ ApplicationPackages = (*applicationPackages)(nil)

type applicationPackages struct {
	client *Client
}

func (v *applicationPackages) Create(ctx context.Context, request *CreateApplicationPackageRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *applicationPackages) Alter(ctx context.Context, request *AlterApplicationPackageRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *applicationPackages) Drop(ctx context.Context, request *DropApplicationPackageRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *applicationPackages) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropApplicationPackageRequest(id).WithIfExists(Bool(true))) }, ctx, id)
}

func (v *applicationPackages) Show(ctx context.Context, request *ShowApplicationPackageRequest) ([]ApplicationPackage, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[applicationPackageRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[applicationPackageRow, ApplicationPackage](dbRows)
	return resultList, nil
}

func (v *applicationPackages) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ApplicationPackage, error) {
	request := NewShowApplicationPackageRequest().
		WithLike(Like{Pattern: String(id.Name())})
	applicationPackages, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(applicationPackages, func(r ApplicationPackage) bool { return r.Name == id.Name() })
}

func (v *applicationPackages) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*ApplicationPackage, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (r *CreateApplicationPackageRequest) toOpts() *CreateApplicationPackageOptions {
	opts := &CreateApplicationPackageOptions{
		IfNotExists:                r.IfNotExists,
		name:                       r.name,
		DataRetentionTimeInDays:    r.DataRetentionTimeInDays,
		MaxDataExtensionTimeInDays: r.MaxDataExtensionTimeInDays,
		DefaultDdlCollation:        r.DefaultDdlCollation,
		Comment:                    r.Comment,
		Distribution:               r.Distribution,
		Tag:                        r.Tag,
	}
	return opts
}

func (r *AlterApplicationPackageRequest) toOpts() *AlterApplicationPackageOptions {
	opts := &AlterApplicationPackageOptions{
		IfExists: r.IfExists,
		name:     r.name,

		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	if r.Set != nil {
		opts.Set = &ApplicationPackageSet{
			DataRetentionTimeInDays:    r.Set.DataRetentionTimeInDays,
			MaxDataExtensionTimeInDays: r.Set.MaxDataExtensionTimeInDays,
			DefaultDdlCollation:        r.Set.DefaultDdlCollation,
			Comment:                    r.Set.Comment,
			Distribution:               r.Set.Distribution,
		}
	}
	if r.Unset != nil {
		opts.Unset = &ApplicationPackageUnset{
			DataRetentionTimeInDays:    r.Unset.DataRetentionTimeInDays,
			MaxDataExtensionTimeInDays: r.Unset.MaxDataExtensionTimeInDays,
			DefaultDdlCollation:        r.Unset.DefaultDdlCollation,
			Comment:                    r.Unset.Comment,
			Distribution:               r.Unset.Distribution,
		}
	}
	if r.ModifyReleaseDirective != nil {
		opts.ModifyReleaseDirective = &ModifyReleaseDirective{
			ReleaseDirective: r.ModifyReleaseDirective.ReleaseDirective,
			Version:          r.ModifyReleaseDirective.Version,
			Patch:            r.ModifyReleaseDirective.Patch,
		}
	}
	if r.SetDefaultReleaseDirective != nil {
		opts.SetDefaultReleaseDirective = &SetDefaultReleaseDirective{
			Version: r.SetDefaultReleaseDirective.Version,
			Patch:   r.SetDefaultReleaseDirective.Patch,
		}
	}
	if r.SetReleaseDirective != nil {
		opts.SetReleaseDirective = &SetReleaseDirective{
			ReleaseDirective: r.SetReleaseDirective.ReleaseDirective,
			Accounts:         r.SetReleaseDirective.Accounts,
			Version:          r.SetReleaseDirective.Version,
			Patch:            r.SetReleaseDirective.Patch,
		}
	}
	if r.UnsetReleaseDirective != nil {
		opts.UnsetReleaseDirective = &UnsetReleaseDirective{
			ReleaseDirective: r.UnsetReleaseDirective.ReleaseDirective,
		}
	}
	if r.AddVersion != nil {
		opts.AddVersion = &AddVersion{
			VersionIdentifier: r.AddVersion.VersionIdentifier,
			Using:             r.AddVersion.Using,
			Label:             r.AddVersion.Label,
		}
	}
	if r.DropVersion != nil {
		opts.DropVersion = &DropVersion{
			VersionIdentifier: r.DropVersion.VersionIdentifier,
		}
	}
	if r.AddPatchForVersion != nil {
		opts.AddPatchForVersion = &AddPatchForVersion{
			VersionIdentifier: r.AddPatchForVersion.VersionIdentifier,
			Using:             r.AddPatchForVersion.Using,
			Label:             r.AddPatchForVersion.Label,
		}
	}
	return opts
}

func (r *DropApplicationPackageRequest) toOpts() *DropApplicationPackageOptions {
	opts := &DropApplicationPackageOptions{
		name:     r.name,
		IfExists: r.IfExists,
	}
	return opts
}

func (r *ShowApplicationPackageRequest) toOpts() *ShowApplicationPackageOptions {
	opts := &ShowApplicationPackageOptions{
		Like:       r.Like,
		StartsWith: r.StartsWith,
		Limit:      r.Limit,
	}
	return opts
}

func (r applicationPackageRow) convert() *ApplicationPackage {
	e := &ApplicationPackage{
		CreatedOn:     r.CreatedOn,
		Name:          r.Name,
		IsDefault:     r.IsDefault == "Y",
		IsCurrent:     r.IsCurrent == "Y",
		Distribution:  r.Distribution,
		Owner:         r.Owner,
		Comment:       r.Comment,
		RetentionTime: r.RetentionTime,
		Options:       r.Options,
	}
	if r.DroppedOn.Valid {
		e.DroppedOn = r.DroppedOn.String
	}
	if r.ApplicationClass.Valid {
		e.ApplicationClass = r.ApplicationClass.String
	}
	return e
}
