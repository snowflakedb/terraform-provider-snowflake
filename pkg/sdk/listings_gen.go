package sdk

import (
	"context"
	"database/sql"
)

type Listings interface {
	Create(ctx context.Context, request *CreateListingRequest) error
	Alter(ctx context.Context, request *AlterListingRequest) error
	Drop(ctx context.Context, request *DropListingRequest) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
	Show(ctx context.Context, request *ShowListingRequest) ([]Listing, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Listing, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Listing, error)
	Describe(ctx context.Context, request *DescribeListingRequest) (*ListingDetails, error)
	ShowVersions(ctx context.Context, request *ShowVersionsListingRequest) ([]ListingVersion, error)
}

// CreateListingOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-listing.
type CreateListingOptions struct {
	create          bool                    `ddl:"static" sql:"CREATE"`
	externalListing bool                    `ddl:"static" sql:"EXTERNAL LISTING"`
	IfNotExists     *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name            AccountObjectIdentifier `ddl:"identifier"`
	With            *ListingWith            `ddl:"keyword"`
	As              *string                 `ddl:"parameter,double_dollar_quotes,no_equals" sql:"AS"`
	From            *Location               `ddl:"parameter,no_quotes,no_equals" sql:"FROM"`
	Publish         *bool                   `ddl:"parameter" sql:"PUBLISH"`
	Review          *bool                   `ddl:"parameter" sql:"REVIEW"`
	Comment         *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ListingWith struct {
	Share              *AccountObjectIdentifier `ddl:"identifier" sql:"SHARE"`
	ApplicationPackage *AccountObjectIdentifier `ddl:"identifier" sql:"APPLICATION PACKAGE"`
}

// AlterListingOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-listing.
type AlterListingOptions struct {
	alter          bool                     `ddl:"static" sql:"ALTER"`
	listing        bool                     `ddl:"static" sql:"LISTING"`
	IfExists       *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name           AccountObjectIdentifier  `ddl:"identifier"`
	Publish        *bool                    `ddl:"keyword" sql:"PUBLISH"`
	Unpublish      *bool                    `ddl:"keyword" sql:"UNPUBLISH"`
	Review         *bool                    `ddl:"keyword" sql:"REVIEW"`
	AlterListingAs *AlterListingAs          `ddl:"keyword" sql:"AS"`
	AddVersion     *AddListingVersion       `ddl:"keyword" sql:"ADD VERSION"`
	RenameTo       *AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	Set            *ListingSet              `ddl:"keyword" sql:"SET"`
	Unset          *ListingUnset            `ddl:"keyword" sql:"UNSET"`
}

type AlterListingAs struct {
	As      string  `ddl:"keyword,double_dollar_quotes"`
	Publish *bool   `ddl:"parameter" sql:"PUBLISH"`
	Review  *bool   `ddl:"parameter" sql:"REVIEW"`
	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type AddListingVersion struct {
	IfNotExists *bool    `ddl:"keyword" sql:"IF NOT EXISTS"`
	VersionName string   `ddl:"keyword,double_quotes"`
	From        Location `ddl:"parameter,no_quotes,no_equals" sql:"FROM"`
	Comment     *string  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ListingSet struct {
	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ListingUnset struct {
	Comment *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropListingOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-listing.
type DropListingOptions struct {
	drop     bool                    `ddl:"static" sql:"DROP"`
	listing  bool                    `ddl:"static" sql:"LISTING"`
	IfExists *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier `ddl:"identifier"`
}

// ShowListingOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-listings.
type ShowListingOptions struct {
	show       bool       `ddl:"static" sql:"SHOW"`
	listings   bool       `ddl:"static" sql:"LISTINGS"`
	Like       *Like      `ddl:"keyword" sql:"LIKE"`
	StartsWith *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	Limit      *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

type listingDBRow struct {
	GlobalName              string         `db:"global_name"`
	Name                    string         `db:"name"`
	Title                   string         `db:"title"`
	Subtitle                sql.NullString `db:"subtitle"`
	Profile                 string         `db:"profile"`
	CreatedOn               string         `db:"created_on"`
	UpdatedOn               string         `db:"updated_on"`
	PublishedOn             sql.NullString `db:"published_on"`
	State                   string         `db:"state"`
	ReviewState             sql.NullString `db:"review_state"`
	Comment                 sql.NullString `db:"comment"`
	Owner                   string         `db:"owner"`
	OwnerRoleType           string         `db:"owner_role_type"`
	Regions                 sql.NullString `db:"regions"`
	TargetAccounts          string         `db:"target_accounts"`
	IsMonetized             bool           `db:"is_monetized"`
	IsApplication           bool           `db:"is_application"`
	IsTargeted              bool           `db:"is_targeted"`
	IsLimitedTrial          sql.NullBool   `db:"is_limited_trial"`
	IsByRequest             sql.NullBool   `db:"is_by_request"`
	Distribution            sql.NullString `db:"distribution"`
	IsMountlessQueryable    sql.NullBool   `db:"is_mountless_queryable"`
	RejectedOn              sql.NullString `db:"rejected_on"`
	OrganizationProfileName sql.NullString `db:"organization_profile_name"`
	UniformListingLocator   sql.NullString `db:"uniform_listing_locator"`
	DetailedTargetAccounts  sql.NullString `db:"detailed_target_accounts"`
}

type Listing struct {
	GlobalName              string
	Name                    string
	Title                   string
	Subtitle                *string
	Profile                 string
	CreatedOn               string
	UpdatedOn               string
	PublishedOn             *string
	State                   ListingState
	ReviewState             *string
	Comment                 *string
	Owner                   string
	OwnerRoleType           string
	Regions                 *string
	TargetAccounts          string
	IsMonetized             bool
	IsApplication           bool
	IsTargeted              bool
	IsLimitedTrial          *bool
	IsByRequest             *bool
	Distribution            *string
	IsMountlessQueryable    *bool
	RejectedOn              *string
	OrganizationProfileName *string
	UniformListingLocator   *string
	DetailedTargetAccounts  *string
}

func (v *Listing) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}
func (v *Listing) ObjectType() ObjectType {
	return ObjectTypeListing
}

// DescribeListingOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-listing.
type DescribeListingOptions struct {
	describe bool                    `ddl:"static" sql:"DESCRIBE"`
	listing  bool                    `ddl:"static" sql:"LISTING"`
	name     AccountObjectIdentifier `ddl:"identifier"`
	Revision *ListingRevision        `ddl:"parameter,no_quotes" sql:"REVISION"`
}

type listingDetailsDBRow struct {
	GlobalName                   string         `db:"global_name"`
	Name                         string         `db:"name"`
	Owner                        string         `db:"owner"`
	OwnerRoleType                string         `db:"owner_role_type"`
	CreatedOn                    string         `db:"created_on"`
	UpdatedOn                    string         `db:"updated_on"`
	PublishedOn                  sql.NullString `db:"published_on"`
	Title                        string         `db:"title"`
	Subtitle                     sql.NullString `db:"subtitle"`
	Description                  sql.NullString `db:"description"`
	ListingTerms                 sql.NullString `db:"listing_terms"`
	State                        string         `db:"state"`
	Share                        sql.NullString `db:"share"`
	ApplicationPackage           sql.NullString `db:"application_package"`
	BusinessNeeds                sql.NullString `db:"business_needs"`
	UsageExamples                sql.NullString `db:"usage_examples"`
	DataAttributes               sql.NullString `db:"data_attributes"`
	Categories                   sql.NullString `db:"categories"`
	Resources                    sql.NullString `db:"resources"`
	Profile                      sql.NullString `db:"profile"`
	CustomizedContactInfo        sql.NullString `db:"customized_contact_info"`
	DataDictionary               sql.NullString `db:"data_dictionary"`
	DataPreview                  sql.NullString `db:"data_preview"`
	Comment                      sql.NullString `db:"comment"`
	Revisions                    string         `db:"revisions"`
	TargetAccounts               sql.NullString `db:"target_accounts"`
	Regions                      sql.NullString `db:"regions"`
	RefreshSchedule              sql.NullString `db:"refresh_schedule"`
	RefreshType                  sql.NullString `db:"refresh_type"`
	ReviewState                  sql.NullString `db:"review_state"`
	RejectionReason              sql.NullString `db:"rejection_reason"`
	UnpublishedByAdminReasons    sql.NullString `db:"unpublished_by_admin_reasons"`
	IsMonetized                  bool           `db:"is_monetized"`
	IsApplication                bool           `db:"is_application"`
	IsTargeted                   bool           `db:"is_targeted"`
	IsLimitedTrial               sql.NullBool   `db:"is_limited_trial"`
	IsByRequest                  sql.NullBool   `db:"is_by_request"`
	LimitedTrialPlan             sql.NullString `db:"limited_trial_plan"`
	RetriedOn                    sql.NullString `db:"retried_on"`
	ScheduledDropTime            sql.NullString `db:"scheduled_drop_time"`
	ManifestYaml                 string         `db:"manifest_yaml"`
	Distribution                 sql.NullString `db:"distribution"`
	IsMountlessQueryable         sql.NullBool   `db:"is_mountless_queryable"`
	OrganizationProfileName      sql.NullString `db:"organization_profile_name"`
	UniformListingLocator        sql.NullString `db:"uniform_listing_locator"`
	TrialDetails                 sql.NullString `db:"trial_details"`
	ApproverContact              sql.NullString `db:"approver_contact"`
	SupportContact               sql.NullString `db:"support_contact"`
	LiveVersionUri               sql.NullString `db:"live_version_uri"`
	LastCommittedVersionUri      sql.NullString `db:"last_committed_version_uri"`
	LastCommittedVersionName     sql.NullString `db:"last_committed_version_name"`
	LastCommittedVersionAlias    sql.NullString `db:"last_committed_version_alias"`
	PublishedVersionUri          sql.NullString `db:"published_version_uri"`
	PublishedVersionName         sql.NullString `db:"published_version_name"`
	PublishedVersionAlias        sql.NullString `db:"published_version_alias"`
	IsShare                      sql.NullBool   `db:"is_share"`
	RequestApprovalType          sql.NullString `db:"request_approval_type"`
	MonetizationDisplayOrder     sql.NullString `db:"monetization_display_order"`
	LegacyUniformListingLocators sql.NullString `db:"legacy_uniform_listing_locators"`
}

type ListingDetails struct {
	GlobalName                   string
	Name                         string
	Owner                        string
	OwnerRoleType                string
	CreatedOn                    string
	UpdatedOn                    string
	PublishedOn                  *string
	Title                        string
	Subtitle                     *string
	Description                  *string
	ListingTerms                 *string
	State                        ListingState
	Share                        *AccountObjectIdentifier
	ApplicationPackage           *AccountObjectIdentifier
	BusinessNeeds                *string
	UsageExamples                *string
	DataAttributes               *string
	Categories                   *string
	Resources                    *string
	Profile                      *string
	CustomizedContactInfo        *string
	DataDictionary               *string
	DataPreview                  *string
	Comment                      *string
	Revisions                    string
	TargetAccounts               *string
	Regions                      *string
	RefreshSchedule              *string
	RefreshType                  *string
	ReviewState                  *string
	RejectionReason              *string
	UnpublishedByAdminReasons    *string
	IsMonetized                  bool
	IsApplication                bool
	IsTargeted                   bool
	IsLimitedTrial               *bool
	IsByRequest                  *bool
	LimitedTrialPlan             *string
	RetriedOn                    *string
	ScheduledDropTime            *string
	ManifestYaml                 string
	Distribution                 *string
	IsMountlessQueryable         *bool
	OrganizationProfileName      *string
	UniformListingLocator        *string
	TrialDetails                 *string
	ApproverContact              *string
	SupportContact               *string
	LiveVersionUri               *string
	LastCommittedVersionUri      *string
	LastCommittedVersionName     *string
	LastCommittedVersionAlias    *string
	PublishedVersionUri          *string
	PublishedVersionName         *string
	PublishedVersionAlias        *string
	IsShare                      *bool
	RequestApprovalType          *string
	MonetizationDisplayOrder     *string
	LegacyUniformListingLocators *string
}

// ShowVersionsListingOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-listing.
type ShowVersionsListingOptions struct {
	show              bool                    `ddl:"static" sql:"SHOW"`
	versionsInListing bool                    `ddl:"static" sql:"VERSIONS IN LISTING"`
	name              AccountObjectIdentifier `ddl:"identifier"`
	Limit             *LimitFrom              `ddl:"keyword" sql:"LIMIT"`
}

type listingVersionDBRow struct {
	CreatedOn         string         `db:"created_on"`
	Name              string         `db:"name"`
	Alias             sql.NullString `db:"alias"`
	LocationUrl       string         `db:"location_url"`
	IsDefault         bool           `db:"is_default"`
	IsLive            bool           `db:"is_live"`
	IsFirst           bool           `db:"is_first"`
	IsLast            bool           `db:"is_last"`
	Comment           sql.NullString `db:"comment"`
	SourceLocationUrl string         `db:"source_location_url"`
	GitCommitHash     sql.NullString `db:"git_commit_hash"`
}

type ListingVersion struct {
	CreatedOn         string
	Name              string
	Alias             *string
	LocationUrl       string
	IsDefault         bool
	IsLive            bool
	IsFirst           bool
	IsLast            bool
	Comment           *string
	SourceLocationUrl string
	GitCommitHash     *string
}
