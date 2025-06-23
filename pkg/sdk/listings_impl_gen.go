package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ Listings = (*listings)(nil)

type listings struct {
	client *Client
}

func (v *listings) Create(ctx context.Context, request *CreateListingRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *listings) CreateFromStage(ctx context.Context, request *CreateFromStageListingRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *listings) Alter(ctx context.Context, request *AlterListingRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *listings) Drop(ctx context.Context, request *DropListingRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *listings) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropListingRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *listings) Show(ctx context.Context, request *ShowListingRequest) ([]Listing, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[listingDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[listingDBRow, Listing](dbRows)
	return resultList, nil
}

func (v *listings) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Listing, error) {
	request := NewShowListingRequest().
		WithLike(Like{Pattern: String(id.Name())})
	listings, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(listings, func(r Listing) bool { return r.Name == id.Name() })
}

func (v *listings) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Listing, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *listings) Describe(ctx context.Context, id AccountObjectIdentifier) (*Listing, error) {
	opts := &DescribeListingOptions{
		name: id,
	}
	result, err := validateAndQueryOne[listingDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert(), nil
}

func (r *CreateListingRequest) toOpts() *CreateListingOptions {
	opts := &CreateListingOptions{
		IfNotExists: r.IfNotExists,
		name:        r.name,

		As:      r.As,
		Publish: r.Publish,
		Review:  r.Review,
		Comment: r.Comment,
	}
	if r.With != nil {
		opts.With = &ListingWith{
			Share:              r.With.Share,
			ApplicationPackage: r.With.ApplicationPackage,
		}
	}
	return opts
}

func (r *CreateFromStageListingRequest) toOpts() *CreateFromStageListingOptions {
	opts := &CreateFromStageListingOptions{
		IfNotExists: r.IfNotExists,
		name:        r.name,

		From:    r.From,
		Publish: r.Publish,
		Review:  r.Review,
	}
	if r.With != nil {
		opts.With = &ListingWith{
			Share:              r.With.Share,
			ApplicationPackage: r.With.ApplicationPackage,
		}
	}
	return opts
}

func (r *AlterListingRequest) toOpts() *AlterListingOptions {
	opts := &AlterListingOptions{
		IfExists:  r.IfExists,
		name:      r.name,
		Publish:   r.Publish,
		Unpublish: r.Unpublish,
		Review:    r.Review,

		RenameTo: r.RenameTo,
	}
	if r.AlterListingAs != nil {
		opts.AlterListingAs = &AlterListingAs{
			As:      r.AlterListingAs.As,
			Publish: r.AlterListingAs.Publish,
			Review:  r.AlterListingAs.Review,
			Comment: r.AlterListingAs.Comment,
		}
	}
	if r.AddVersion != nil {
		opts.AddVersion = &AddListingVersion{
			IfNotExists: r.AddVersion.IfNotExists,
			VersionName: r.AddVersion.VersionName,
			From:        r.AddVersion.From,
			Comment:     r.AddVersion.Comment,
		}
	}
	if r.Set != nil {
		opts.Set = &ListingSet{
			Comment: r.Set.Comment,
		}
	}
	return opts
}

func (r *DropListingRequest) toOpts() *DropListingOptions {
	opts := &DropListingOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowListingRequest) toOpts() *ShowListingOptions {
	opts := &ShowListingOptions{
		Like:       r.Like,
		StartsWith: r.StartsWith,
		Limit:      r.Limit,
	}
	return opts
}

func (r listingDBRow) convert() *Listing {
	return &Listing{
		GlobalName:     "",
		Name:           "",
		Title:          "",
		Subtitle:       "",
		Profile:        "",
		CreatedOn:      "",
		UpdatedOn:      "",
		PublishedOn:    "",
		State:          "",
		ReviewState:    "",
		Comment:        "",
		Owner:          "",
		OwnerRoleType:  "",
		Regions:        "",
		TargetAccounts: "",
		IsMonetized:    "",
		IsApplication:  "",
		IsTargeted:     "",
	}
}

func (r *DescribeListingRequest) toOpts() *DescribeListingOptions {
	opts := &DescribeListingOptions{
		name:     r.name,
		Revision: r.Revision,
	}
	return opts
}
