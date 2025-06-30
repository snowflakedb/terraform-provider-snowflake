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

func (v *listings) Describe(ctx context.Context, id AccountObjectIdentifier) (*ListingDetails, error) {
	opts := &DescribeListingOptions{
		name: id,
	}
	result, err := validateAndQueryOne[listingDetailsDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert(), nil
}

func (r *CreateListingRequest) toOpts() *CreateListingOptions {
	opts := &CreateListingOptions{
		IfNotExists: r.IfNotExists,
		name:        r.name,

		From:    r.From,
		Publish: r.Publish,
		Review:  r.Review,
		Comment: r.Comment,
	}
	if r.As != nil {
		opts.As = *r.As
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
	if r.Unset != nil {
		opts.Unset = &ListingUnset{
			Comment: r.Unset.Comment,
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
	l := &Listing{
		GlobalName:     r.GlobalName,
		Name:           r.Name,
		Title:          r.Title,
		Profile:        r.Profile,
		CreatedOn:      r.CreatedOn,
		UpdatedOn:      r.UpdatedOn,
		ReviewState:    r.ReviewState,
		Owner:          r.Owner,
		OwnerRoleType:  r.OwnerRoleType,
		TargetAccounts: r.TargetAccounts,
		IsMonetized:    r.IsMonetized,
		IsApplication:  r.IsApplication,
		IsTargeted:     r.IsTargeted,
	}
	if state, err := ToListingState(r.State); err == nil {
		l.State = state
	}
	mapNullString(&l.Subtitle, r.Subtitle)
	mapNullString(&l.PublishedOn, r.PublishedOn)
	mapNullString(&l.Comment, r.Comment)
	mapNullString(&l.Regions, r.Regions)
	mapNullBool(&l.IsLimitedTrial, r.IsLimitedTrial)
	mapNullBool(&l.IsByRequest, r.IsByRequest)
	mapNullString(&l.Distribution, r.Distribution)
	mapNullBool(&l.IsMountlessQueryable, r.IsMountlessQueryable)
	mapNullString(&l.RejectedOn, r.RejectedOn)
	mapNullString(&l.OrganizationProfileName, r.OrganizationProfileName)
	mapNullString(&l.UniformListingLocator, r.UniformListingLocator)
	mapNullString(&l.DetailedTargetAccounts, r.DetailedTargetAccounts)

	return l
}

func (r *DescribeListingRequest) toOpts() *DescribeListingOptions {
	opts := &DescribeListingOptions{
		name:     r.name,
		Revision: r.Revision,
	}
	return opts
}

func (r listingDetailsDBRow) convert() *ListingDetails {
	ld := &ListingDetails{
		GlobalName:    r.GlobalName,
		Name:          r.Name,
		Owner:         r.Owner,
		OwnerRoleType: r.OwnerRoleType,
		CreatedOn:     r.CreatedOn,
		UpdatedOn:     r.UpdatedOn,
		Title:         r.Title,
		Revisions:     r.Revisions,
		ReviewState:   r.ReviewState,
		ManifestYaml:  r.ManifestYaml,
		IsMonetized:   r.IsMonetized,
		IsApplication: r.IsApplication,
		IsTargeted:    r.IsTargeted,
	}

	mapNullString(&ld.PublishedOn, r.PublishedOn)
	mapNullString(&ld.Subtitle, r.Subtitle)
	mapNullString(&ld.Description, r.Description)
	mapNullString(&ld.ListingTerms, r.ListingTerms)
	mapStringWithMapping(&ld.State, r.State, ToListingState)
	mapNullStringWithMapping(&ld.Share, r.Share, ParseAccountObjectIdentifier)
	mapNullStringWithMapping(&ld.ApplicationPackage, r.ApplicationPackage, ParseAccountObjectIdentifier)
	mapNullString(&ld.BusinessNeeds, r.BusinessNeeds)
	mapNullString(&ld.UsageExamples, r.UsageExamples)
	mapNullString(&ld.DataAttributes, r.DataAttributes)
	mapNullString(&ld.Categories, r.Categories)
	mapNullString(&ld.Resources, r.Resources)
	mapNullString(&ld.Profile, r.Profile)
	mapNullString(&ld.CustomizedContactInfo, r.CustomizedContactInfo)
	mapNullString(&ld.DataDictionary, r.DataDictionary)
	mapNullString(&ld.DataPreview, r.DataPreview)
	mapNullString(&ld.Comment, r.Comment)
	mapNullString(&ld.TargetAccounts, r.TargetAccounts)
	mapNullString(&ld.Regions, r.Regions)
	mapNullString(&ld.RefreshSchedule, r.RefreshSchedule)
	mapNullString(&ld.RefreshType, r.RefreshType)
	mapNullString(&ld.RejectionReason, r.RejectionReason)
	mapNullString(&ld.UnpublishedByAdminReasons, r.UnpublishedByAdminReasons)
	mapNullBool(&ld.IsLimitedTrial, r.IsLimitedTrial)
	mapNullBool(&ld.IsByRequest, r.IsByRequest)
	mapNullString(&ld.LimitedTrialPlan, r.LimitedTrialPlan)
	mapNullString(&ld.RetriedOn, r.RetriedOn)
	mapNullString(&ld.ScheduledDropTime, r.ScheduledDropTime)
	mapNullString(&ld.Distribution, r.Distribution)
	mapNullBool(&ld.IsMountlessQueryable, r.IsMountlessQueryable)
	mapNullString(&ld.OrganizationProfileName, r.OrganizationProfileName)
	mapNullString(&ld.UniformListingLocator, r.UniformListingLocator)
	mapNullString(&ld.TrialDetails, r.TrialDetails)
	mapNullString(&ld.ApproverContact, r.ApproverContact)
	mapNullString(&ld.SupportContact, r.SupportContact)
	mapNullString(&ld.LiveVersionUri, r.LiveVersionUri)
	mapNullString(&ld.LastCommittedVersionUri, r.LastCommittedVersionUri)
	mapNullString(&ld.LastCommittedVersionName, r.LastCommittedVersionName)
	mapNullString(&ld.LastCommittedVersionAlias, r.LastCommittedVersionAlias)
	mapNullString(&ld.PublishedVersionUri, r.PublishedVersionUri)
	mapNullString(&ld.PublishedVersionName, r.PublishedVersionName)
	mapNullString(&ld.PublishedVersionAlias, r.PublishedVersionAlias)
	mapNullBool(&ld.IsShare, r.IsShare)
	mapNullString(&ld.RequestApprovalType, r.RequestApprovalType)
	mapNullString(&ld.MonetizationDisplayOrder, r.MonetizationDisplayOrder)
	mapNullString(&ld.LegacyUniformListingLocators, r.LegacyUniformListingLocators)

	return ld
}
