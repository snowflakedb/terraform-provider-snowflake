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
	mapStringIfNotNil(&l.Subtitle, r.Subtitle)
	mapStringIfNotNil(&l.PublishedOn, r.PublishedOn)
	mapStringIfNotNil(&l.Comment, r.Comment)
	mapStringIfNotNil(&l.Regions, r.Regions)
	mapBoolIfNotNil(&l.IsLimitedTrial, r.IsLimitedTrial)
	mapBoolIfNotNil(&l.IsByRequest, r.IsByRequest)
	mapStringIfNotNil(&l.Distribution, r.Distribution)
	mapBoolIfNotNil(&l.IsMountlessQueryable, r.IsMountlessQueryable)
	mapStringIfNotNil(&l.RejectedOn, r.RejectedOn)
	mapStringIfNotNil(&l.OrganizationProfileName, r.OrganizationProfileName)
	mapStringIfNotNil(&l.UniformListingLocator, r.UniformListingLocator)
	mapStringIfNotNil(&l.DetailedTargetAccounts, r.DetailedTargetAccounts)

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

	mapStringIfNotNil(&ld.PublishedOn, r.PublishedOn)
	mapStringIfNotNil(&ld.Subtitle, r.Subtitle)
	mapStringIfNotNil(&ld.Description, r.Description)
	mapStringIfNotNil(&ld.ListingTerms, r.ListingTerms)
	mapStringWithMapping(&ld.State, r.State, ToListingState)
	mapStringWithMappingIfNotNil(&ld.Share, r.Share, ParseAccountObjectIdentifier)
	mapStringWithMappingIfNotNil(&ld.ApplicationPackage, r.ApplicationPackage, ParseAccountObjectIdentifier)
	mapStringIfNotNil(&ld.BusinessNeeds, r.BusinessNeeds)
	mapStringIfNotNil(&ld.UsageExamples, r.UsageExamples)
	mapStringIfNotNil(&ld.DataAttributes, r.DataAttributes)
	mapStringIfNotNil(&ld.Categories, r.Categories)
	mapStringIfNotNil(&ld.Resources, r.Resources)
	mapStringIfNotNil(&ld.Profile, r.Profile)
	mapStringIfNotNil(&ld.CustomizedContactInfo, r.CustomizedContactInfo)
	mapStringIfNotNil(&ld.DataDictionary, r.DataDictionary)
	mapStringIfNotNil(&ld.DataPreview, r.DataPreview)
	mapStringIfNotNil(&ld.Comment, r.Comment)
	mapStringIfNotNil(&ld.TargetAccounts, r.TargetAccounts)
	mapStringIfNotNil(&ld.Regions, r.Regions)
	mapStringIfNotNil(&ld.RefreshSchedule, r.RefreshSchedule)
	mapStringIfNotNil(&ld.RefreshType, r.RefreshType)
	mapStringIfNotNil(&ld.RejectionReason, r.RejectionReason)
	mapStringIfNotNil(&ld.UnpublishedByAdminReasons, r.UnpublishedByAdminReasons)
	mapBoolIfNotNil(&ld.IsLimitedTrial, r.IsLimitedTrial)
	mapBoolIfNotNil(&ld.IsByRequest, r.IsByRequest)
	mapStringIfNotNil(&ld.LimitedTrialPlan, r.LimitedTrialPlan)
	mapStringIfNotNil(&ld.RetriedOn, r.RetriedOn)
	mapStringIfNotNil(&ld.ScheduledDropTime, r.ScheduledDropTime)
	mapStringIfNotNil(&ld.Distribution, r.Distribution)
	mapBoolIfNotNil(&ld.IsMountlessQueryable, r.IsMountlessQueryable)
	mapStringIfNotNil(&ld.OrganizationProfileName, r.OrganizationProfileName)
	mapStringIfNotNil(&ld.UniformListingLocator, r.UniformListingLocator)
	mapStringIfNotNil(&ld.TrialDetails, r.TrialDetails)
	mapStringIfNotNil(&ld.ApproverContact, r.ApproverContact)
	mapStringIfNotNil(&ld.SupportContact, r.SupportContact)
	mapStringIfNotNil(&ld.LiveVersionUri, r.LiveVersionUri)
	mapStringIfNotNil(&ld.LastCommittedVersionUri, r.LastCommittedVersionUri)
	mapStringIfNotNil(&ld.LastCommittedVersionName, r.LastCommittedVersionName)
	mapStringIfNotNil(&ld.LastCommittedVersionAlias, r.LastCommittedVersionAlias)
	mapStringIfNotNil(&ld.PublishedVersionUri, r.PublishedVersionUri)
	mapStringIfNotNil(&ld.PublishedVersionName, r.PublishedVersionName)
	mapStringIfNotNil(&ld.PublishedVersionAlias, r.PublishedVersionAlias)
	mapBoolIfNotNil(&ld.IsShare, r.IsShare)
	mapStringIfNotNil(&ld.RequestApprovalType, r.RequestApprovalType)
	mapStringIfNotNil(&ld.MonetizationDisplayOrder, r.MonetizationDisplayOrder)
	mapStringIfNotNil(&ld.LegacyUniformListingLocators, r.LegacyUniformListingLocators)

	return ld
}
