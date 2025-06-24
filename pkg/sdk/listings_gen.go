package sdk

import "context"

type Listings interface {
	Create(ctx context.Context, request *CreateListingRequest) error
	CreateFromStage(ctx context.Context, request *CreateFromStageListingRequest) error
	Alter(ctx context.Context, request *AlterListingRequest) error
	Drop(ctx context.Context, request *DropListingRequest) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
	Show(ctx context.Context, request *ShowListingRequest) ([]Listing, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Listing, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Listing, error)
	Describe(ctx context.Context, id AccountObjectIdentifier) (*ListingDetails, error)
}

// CreateListingOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-listing.
type CreateListingOptions struct {
	create          bool                    `ddl:"static" sql:"CREATE"`
	externalListing bool                    `ddl:"static" sql:"EXTERNAL LISTING"`
	IfNotExists     *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name            AccountObjectIdentifier `ddl:"identifier"`
	With            *ListingWith            `ddl:"keyword"`
	As              string                  `ddl:"parameter,double_dollar_quotes,no_equals" sql:"AS"`
	Publish         *bool                   `ddl:"parameter" sql:"PUBLISH"`
	Review          *bool                   `ddl:"parameter" sql:"REVIEW"`
	Comment         *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ListingWith struct {
	Share              *AccountObjectIdentifier `ddl:"identifier" sql:"SHARE"`
	ApplicationPackage *AccountObjectIdentifier `ddl:"identifier" sql:"APPLICATION PACKAGE"`
}

// CreateFromStageListingOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-listing.
type CreateFromStageListingOptions struct {
	create          bool                    `ddl:"static" sql:"CREATE"`
	externalListing bool                    `ddl:"static" sql:"EXTERNAL LISTING"`
	IfNotExists     *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name            AccountObjectIdentifier `ddl:"identifier"`
	With            *ListingWith            `ddl:"keyword"`
	From            Location                `ddl:"parameter,no_quotes,no_equals" sql:"FROM"`
	Publish         *bool                   `ddl:"parameter" sql:"PUBLISH"`
	Review          *bool                   `ddl:"parameter" sql:"REVIEW"`
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
}

type AlterListingAs struct {
	As      string  `ddl:"keyword,double_dollar_quotes"`
	Publish *bool   `ddl:"parameter" sql:"PUBLISH"`
	Review  *bool   `ddl:"parameter" sql:"REVIEW"`
	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type AddListingVersion struct {
	IfNotExists *bool    `ddl:"keyword" sql:"IF NOT EXISTS"`
	VersionName string   `ddl:"keyword"`
	From        Location `ddl:"parameter,no_quotes,no_equals" sql:"FROM"`
	Comment     *string  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ListingSet struct {
	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
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
	GlobalName     string `db:"global_name"`
	Name           string `db:"name"`
	Title          string `db:"title"`
	Subtitle       string `db:"subtitle"`
	Profile        string `db:"profile"`
	CreatedOn      string `db:"created_on"`
	UpdatedOn      string `db:"updated_on"`
	PublishedOn    string `db:"published_on"`
	State          string `db:"state"`
	ReviewState    string `db:"review_state"`
	Comment        string `db:"comment"`
	Owner          string `db:"owner"`
	OwnerRoleType  string `db:"owner_role_type"`
	Regions        string `db:"regions"`
	TargetAccounts string `db:"target_accounts"`
	IsMonetized    string `db:"is_monetized"`
	IsApplication  string `db:"is_application"`
	IsTargeted     string `db:"is_targeted"`
}

type Listing struct {
	GlobalName     string
	Name           string
	Title          string
	Subtitle       string
	Profile        string
	CreatedOn      string
	UpdatedOn      string
	PublishedOn    string
	State          string
	ReviewState    string
	Comment        string
	Owner          string
	OwnerRoleType  string
	Regions        string
	TargetAccounts string
	IsMonetized    string
	IsApplication  string
	IsTargeted     string
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
	GlobalName                string `db:"global_name"`
	Name                      string `db:"name"`
	Owner                     string `db:"owner"`
	OwnerRoleType             string `db:"owner_role_type"`
	CreatedOn                 string `db:"created_on"`
	UpdatedOn                 string `db:"updated_on"`
	PublishedOn               string `db:"published_on"`
	Title                     string `db:"title"`
	Subtitle                  string `db:"subtitle"`
	Description               string `db:"description"`
	TargetAccounts            string `db:"target_accounts"`
	IsMonetized               string `db:"is_monetized"`
	IsApplication             string `db:"is_application"`
	IsTargeted                string `db:"is_targeted"`
	State                     string `db:"state"`
	Revisions                 string `db:"revisions"`
	Comment                   string `db:"comment"`
	RefreshedSchedule         string `db:"refreshed_schedule"`
	RefreshType               string `db:"refresh_type"`
	BusinessNeeds             string `db:"business_needs"`
	UsageExamples             string `db:"usage_examples"`
	ListingTerms              string `db:"listing_terms"`
	Profile                   string `db:"profile"`
	CustomizedContactInfo     string `db:"customized_contact_info"`
	ApplicationPackage        string `db:"application_package"`
	DataDictionary            string `db:"data_dictionary"`
	Regions                   string `db:"regions"`
	ManifestYaml              string `db:"manifest_yaml"`
	ReviewState               string `db:"review_state"`
	RejectionReason           string `db:"rejection_reason"`
	Categories                string `db:"categories"`
	Resources                 string `db:"resources"`
	UnpublishedByAdminReasons string `db:"unpublished_by_admin_reasons"`
}

type ListingDetails struct {
	GlobalName                string
	Name                      string
	Owner                     string
	OwnerRoleType             string
	CreatedOn                 string
	UpdatedOn                 string
	PublishedOn               string
	Title                     string
	Subtitle                  string
	Description               string
	TargetAccounts            string
	IsMonetized               string
	IsApplication             string
	IsTargeted                string
	State                     string
	Revisions                 string
	Comment                   string
	RefreshedSchedule         string
	RefreshType               string
	BusinessNeeds             string
	UsageExamples             string
	ListingTerms              string
	Profile                   string
	CustomizedContactInfo     string
	ApplicationPackage        string
	DataDictionary            string
	Regions                   string
	ManifestYaml              string
	ReviewState               string
	RejectionReason           string
	Categories                string
	Resources                 string
	UnpublishedByAdminReasons string
}
