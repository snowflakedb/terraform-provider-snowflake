package sdk

import (
	"fmt"
	"slices"
	"strings"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

type ListingRevision string

const (
	ListingRevisionDraft     ListingRevision = "DRAFT"
	ListingRevisionPublished ListingRevision = "PUBLISHED"
)

type ListingState string

const (
	ListingStateDraft       ListingState = "DRAFT"
	ListingStatePublished   ListingState = "PUBLISHED"
	ListingStateUnpublished ListingState = "UNPUBLISHED"
)

var AllListingStates = []ListingState{
	ListingStateDraft,
	ListingStatePublished,
	ListingStateUnpublished,
}

func ToListingState(s string) (ListingState, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllListingStates, ListingState(s)) {
		return "", fmt.Errorf("invalid listing state: %s", s)
	}
	return ListingState(s), nil
}

var listingWithDef = g.NewQueryStruct("ListingWith").
	OptionalIdentifier("Share", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("SHARE")).
	OptionalIdentifier("ApplicationPackage", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("APPLICATION PACKAGE")).
	WithValidation(g.ExactlyOneValueSet, "Share", "ApplicationPackage")

// There are more fields listed than in https://docs.snowflake.com/en/sql-reference/sql/show-listings.
// They are mapped straight from the SHOW LISTINGS output.
var listingDbRow = g.DbStruct("listingDBRow").
	Text("global_name").
	Text("name").
	Text("title").
	OptionalText("subtitle").
	Text("profile").
	Text("created_on").
	Text("updated_on").
	OptionalText("published_on").
	Text("state").
	OptionalText("review_state").
	OptionalText("comment").
	Text("owner").
	Text("owner_role_type").
	OptionalText("regions").
	Text("target_accounts").
	Bool("is_monetized").
	Bool("is_application").
	Bool("is_targeted").
	OptionalBool("is_limited_trial").
	OptionalBool("is_by_request").
	OptionalText("distribution").
	OptionalBool("is_mountless_queryable").
	OptionalText("rejected_on").
	OptionalText("organization_profile_name").
	OptionalText("uniform_listing_locator").
	OptionalText("detailed_target_accounts")

var listing = g.PlainStruct("Listing").
	Text("GlobalName").
	Text("Name").
	Text("Title").
	OptionalText("Subtitle").
	Text("Profile").
	Text("CreatedOn").
	Text("UpdatedOn").
	OptionalText("PublishedOn").
	Field("State", g.KindOfT[ListingState]()).
	OptionalText("ReviewState").
	OptionalText("Comment").
	Text("Owner").
	Text("OwnerRoleType").
	OptionalText("Regions").
	Text("TargetAccounts").
	Bool("IsMonetized").
	Bool("IsApplication").
	Bool("IsTargeted").
	OptionalBool("IsLimitedTrial").
	OptionalBool("IsByRequest").
	OptionalText("Distribution").
	OptionalBool("IsMountlessQueryable").
	OptionalText("RejectedOn").
	OptionalText("OrganizationProfileName").
	OptionalText("UniformListingLocator").
	OptionalText("DetailedTargetAccounts")

// There are more fields listed than in https://docs.snowflake.com/en/sql-reference/sql/desc-listing
// They are mapped straight from the DESC LISTING output.
var listingDetailsDbRow = g.DbStruct("listingDetailsDBRow").
	Text("global_name").
	Text("name").
	Text("owner").
	Text("owner_role_type").
	Text("created_on").
	Text("updated_on").
	OptionalText("published_on").
	Text("title").
	OptionalText("subtitle").
	OptionalText("description").
	OptionalText("listing_terms").
	Text("state").
	OptionalText("share").
	OptionalText("application_package").
	OptionalText("business_needs").
	OptionalText("usage_examples").
	OptionalText("data_attributes").
	OptionalText("categories").
	OptionalText("resources").
	OptionalText("profile").
	OptionalText("customized_contact_info").
	OptionalText("data_dictionary").
	OptionalText("data_preview").
	OptionalText("comment").
	Text("revisions").
	OptionalText("target_accounts").
	OptionalText("regions").
	OptionalText("refresh_schedule").
	OptionalText("refresh_type").
	OptionalText("review_state").
	OptionalText("rejection_reason").
	OptionalText("unpublished_by_admin_reasons").
	Bool("is_monetized").
	Bool("is_application").
	Bool("is_targeted").
	OptionalBool("is_limited_trial").
	OptionalBool("is_by_request").
	OptionalText("limited_trial_plan").
	OptionalText("retried_on").
	OptionalText("scheduled_drop_time").
	Text("manifest_yaml").
	OptionalText("distribution").
	OptionalBool("is_mountless_queryable").
	OptionalText("organization_profile_name").
	OptionalText("uniform_listing_locator").
	OptionalText("trial_details").
	OptionalText("approver_contact").
	OptionalText("support_contact").
	OptionalText("live_version_uri").
	OptionalText("last_committed_version_uri").
	OptionalText("last_committed_version_name").
	OptionalText("last_committed_version_alias").
	OptionalText("published_version_uri").
	OptionalText("published_version_name").
	OptionalText("published_version_alias").
	OptionalBool("is_share").
	OptionalText("request_approval_type").
	OptionalText("monetization_display_order").
	OptionalText("legacy_uniform_listing_locators")

var listingDetails = g.PlainStruct("ListingDetails").
	Text("GlobalName").
	Text("Name").
	Text("Owner").
	Text("OwnerRoleType").
	Text("CreatedOn").
	Text("UpdatedOn").
	OptionalText("PublishedOn").
	Text("Title").
	OptionalText("Subtitle").
	OptionalText("Description").
	OptionalText("ListingTerms").
	Field("State", g.KindOfT[ListingState]()).
	Field("Share", g.KindOfTPointer[AccountObjectIdentifier]()).
	Field("ApplicationPackage", g.KindOfTPointer[AccountObjectIdentifier]()).
	OptionalText("BusinessNeeds").
	OptionalText("UsageExamples").
	OptionalText("DataAttributes").
	OptionalText("Categories").
	OptionalText("Resources").
	OptionalText("Profile").
	OptionalText("CustomizedContactInfo").
	OptionalText("DataDictionary").
	OptionalText("DataPreview").
	OptionalText("Comment").
	Text("Revisions").
	OptionalText("TargetAccounts").
	OptionalText("Regions").
	OptionalText("RefreshSchedule").
	OptionalText("RefreshType").
	OptionalText("ReviewState").
	OptionalText("RejectionReason").
	OptionalText("UnpublishedByAdminReasons").
	Bool("IsMonetized").
	Bool("IsApplication").
	Bool("IsTargeted").
	OptionalBool("IsLimitedTrial").
	OptionalBool("IsByRequest").
	OptionalText("LimitedTrialPlan").
	OptionalText("RetriedOn").
	OptionalText("ScheduledDropTime").
	Text("ManifestYaml").
	OptionalText("Distribution").
	OptionalBool("IsMountlessQueryable").
	OptionalText("OrganizationProfileName").
	OptionalText("UniformListingLocator").
	OptionalText("TrialDetails").
	OptionalText("ApproverContact").
	OptionalText("SupportContact").
	OptionalText("LiveVersionUri").
	OptionalText("LastCommittedVersionUri").
	OptionalText("LastCommittedVersionName").
	OptionalText("LastCommittedVersionAlias").
	OptionalText("PublishedVersionUri").
	OptionalText("PublishedVersionName").
	OptionalText("PublishedVersionAlias").
	OptionalBool("IsShare").
	OptionalText("RequestApprovalType").
	OptionalText("MonetizationDisplayOrder").
	OptionalText("LegacyUniformListingLocators")

var listingVersionDbRow = g.DbStruct("listingVersionDBRow").
	Text("created_on").
	Text("name").
	OptionalText("alias").
	Text("location_url").
	Bool("is_default").
	Bool("is_live").
	Bool("is_first").
	Bool("is_last").
	OptionalText("comment").
	Text("source_location_url").
	OptionalText("git_commit_hash")

var listingVersion = g.PlainStruct("ListingVersion").
	Text("CreatedOn").
	Text("Name").
	OptionalText("Alias").
	Text("LocationUrl").
	Bool("IsDefault").
	Bool("IsLive").
	Bool("IsFirst").
	Bool("IsLast").
	OptionalText("Comment").
	Text("SourceLocationUrl").
	OptionalText("GitCommitHash")

var ListingsDef = g.NewInterface(
	"Listings",
	"Listing",
	g.KindOfT[AccountObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-listing",
		g.NewQueryStruct("CreateListing").
			Create().
			SQL("EXTERNAL LISTING").
			IfNotExists().
			Name().
			OptionalQueryStructField("With", listingWithDef, g.KeywordOptions()).
			OptionalTextAssignment("AS", g.ParameterOptions().NoEquals().DoubleDollarQuotes()).
			PredefinedQueryStructField("From", g.KindOfTPointer[Location](), g.ParameterOptions().NoQuotes().NoEquals().SQL("FROM")).
			OptionalBooleanAssignment("PUBLISH", g.ParameterOptions()).
			OptionalBooleanAssignment("REVIEW", g.ParameterOptions()).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "As", "From"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-listing",
		g.NewQueryStruct("AlterListing").
			Alter().
			SQL("LISTING").
			IfExists().
			Name().
			OptionalSQL("PUBLISH").
			OptionalSQL("UNPUBLISH").
			OptionalSQL("REVIEW").
			OptionalQueryStructField(
				"AlterListingAs",
				g.NewQueryStruct("AlterListingAs").
					Text("As", g.KeywordOptions().Required().DoubleDollarQuotes()).
					OptionalBooleanAssignment("PUBLISH", g.ParameterOptions()).
					OptionalBooleanAssignment("REVIEW", g.ParameterOptions()).
					OptionalComment(),
				g.KeywordOptions().SQL("AS"),
			).
			OptionalQueryStructField(
				"AddVersion",
				g.NewQueryStruct("AddListingVersion").
					IfNotExists().
					Text("VersionName", g.KeywordOptions().DoubleQuotes()).
					PredefinedQueryStructField("From", "Location", g.ParameterOptions().Required().NoQuotes().NoEquals().SQL("FROM")).
					OptionalComment(),
				g.KeywordOptions().SQL("ADD VERSION"),
			).
			OptionalIdentifier("RenameTo", g.KindOfTPointer[AccountObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("ListingSet").
					OptionalComment(),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("ListingUnset").
					OptionalSQL("COMMENT"),
				g.KeywordOptions().SQL("UNSET"),
			).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfExists", "AddVersion").
			WithValidation(g.ExactlyOneValueSet, "Publish", "Unpublish", "Review", "AlterListingAs", "AddVersion", "RenameTo", "Set", "Unset"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-listing",
		g.NewQueryStruct("DropListing").
			Drop().
			SQL("LISTING").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-listings",
		listingDbRow,
		listing,
		g.NewQueryStruct("ShowListings").
			Show().
			SQL("LISTINGS").
			OptionalLike().
			OptionalStartsWith().
			OptionalLimitFrom(),
	).
	ShowByIdOperationWithFiltering(g.ShowByIDLikeFiltering).
	CustomShowOperation(
		"Describe",
		g.ShowMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-listing",
		listingDetailsDbRow,
		listingDetails,
		g.NewQueryStruct("DescribeListing").
			Describe().
			SQL("LISTING").
			Name().
			OptionalAssignment("REVISION", g.KindOfT[ListingRevision](), g.ParameterOptions().NoQuotes()).
			WithValidation(g.ValidIdentifier, "name"),
	).
	CustomShowOperation(
		"ShowVersions",
		g.ShowMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-listing",
		listingVersionDbRow,
		listingVersion,
		g.NewQueryStruct("ShowListings").
			Show().
			SQL("VERSIONS IN LISTING").
			Name().
			OptionalLimit().
			WithValidation(g.ValidIdentifier, "name"),
	)

	// TODO [SNOW-2236968]: Organization listing may have its interface, but most of the operations would be pass through functions to the Listings interface
	// TODO [SNOW-2236968]: Show available listings
	// TODO [SNOW-2236968]: Describe available listing
	// TODO [SNOW-2236968]: Listing manifest builder - https://docs.snowflake.com/en/progaccess/listing-manifest-reference
	// TODO [SNOW-2236968]: Test mapping functions (ToListingRevision and ToListingState)
