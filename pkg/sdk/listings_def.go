package sdk

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

//go:generate go run ./poc/main.go

type ListingRevision string

const (
	ListingRevisionDraft     ListingRevision = "DRAFT"
	ListingRevisionPublished ListingRevision = "PUBLISHED"
)

var listingWithDef = g.NewQueryStruct("ListingWith").
	OptionalIdentifier("Share", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("SHARE")).
	OptionalIdentifier("ApplicationPackage", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("APPLICATION PACKAGE")).
	WithValidation(g.ExactlyOneValueSet, "Share", "ApplicationPackage")

var listingDbRow = g.DbStruct("listingDBRow").
	Text("global_name").
	Text("name").
	Text("title").
	Text("subtitle").
	Text("profile").
	Text("created_on").
	Text("updated_on").
	OptionalText("published_on").
	Text("state").
	Text("review_state").
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
	Text("Subtitle").
	Text("Profile").
	Text("CreatedOn").
	Text("UpdatedOn").
	OptionalText("PublishedOn").
	Text("State").
	Text("ReviewState").
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

var listingDetailsDbRow = g.DbStruct("listingDetailsDBRow").
	Text("global_name").
	Text("name").
	Text("owner").
	Text("owner_role_type").
	Text("created_on").
	Text("updated_on").
	Text("published_on").
	Text("title").
	Text("subtitle").
	Text("description").
	Text("listing_terms").
	Text("state").
	Text("share").
	Text("application_package").
	Text("business_needs").
	Text("usage_examples").
	OptionalText("data_attributes"). // TODO: Not documented
	Text("categories").
	Text("resources").
	Text("profile").
	Text("customized_contact_info").
	Text("data_dictionary").
	OptionalText("data_preview"). // TODO: Not documented
	Text("comment").
	Text("revisions").
	Text("target_accounts").
	Text("regions").
	Text("refresh_schedule").
	Text("refresh_type").
	Text("review_state").
	Text("rejection_reason").
	Text("unpublished_by_admin_reasons").
	Text("is_monetized").
	Text("is_application").
	Text("is_targeted").
	OptionalText("is_limited_trial").    // TODO: Not documented
	OptionalText("is_by_request").       // TODO: Not documented
	OptionalText("limited_trial_plan").  // TODO: Not documented
	OptionalText("retried_on").          // TODO: Not documented
	OptionalText("scheduled_drop_time"). // TODO: Not documented
	Text("manifest_yaml").
	OptionalText("distribution").                   // TODO: Not documented
	OptionalText("is_mountless_queryable").         // TODO: Not documented
	OptionalText("organization_profile_name").      // TODO: Not documented
	OptionalText("uniform_listing_locator").        // TODO: Not documented
	OptionalText("trial_details").                  // TODO: Not documented
	OptionalText("trial_details").                  // TODO: Not documented
	OptionalText("approver_contact").               // TODO: Not documented
	OptionalText("support_contact").                // TODO: Not documented
	OptionalText("live_version_uri").               // TODO: Not documented
	OptionalText("last_committed_version_uri").     // TODO: Not documented
	OptionalText("last_committed_version_name").    // TODO: Not documented
	OptionalText("last_committed_version_alias").   // TODO: Not documented
	OptionalText("published_version_uri").          // TODO: Not documented
	OptionalText("published_version_name").         // TODO: Not documented
	OptionalText("published_version_alias").        // TODO: Not documented
	OptionalText("is_share").                       // TODO: Not documented
	OptionalText("request_approval_type").          // TODO: Not documented
	OptionalText("monetization_display_order").     // TODO: Not documented
	OptionalText("legacy_uniform_listing_locators") // TODO: Not documented

var listingDetails = g.PlainStruct("ListingDetails").
	Text("GlobalName").
	Text("Name").
	Text("Owner").
	Text("OwnerRoleType").
	Text("CreatedOn").
	Text("UpdatedOn").
	Text("PublishedOn").
	Text("Title").
	Text("Subtitle").
	Text("Description").
	Text("ListingTerms").
	Text("State").
	Text("Share").
	Text("ApplicationPackage").
	Text("BusinessNeeds").
	Text("UsageExamples").
	OptionalText("DataAttributes").
	Text("Categories").
	Text("Resources").
	Text("Profile").
	Text("CustomizedContactInfo").
	Text("DataDictionary").
	OptionalText("DataPreview").
	Text("Comment").
	Text("Revisions").
	Text("TargetAccounts").
	Text("Regions").
	Text("RefreshSchedule").
	Text("RefreshType").
	Text("ReviewState").
	Text("RejectionReason").
	Text("UnpublishedByAdminReasons").
	Text("IsMonetized").
	Text("IsApplication").
	Text("IsTargeted").
	OptionalText("IsLimitedTrial").
	OptionalText("IsByRequest").
	OptionalText("LimitedTrialPlan").
	OptionalText("RetriedOn").
	OptionalText("ScheduledDropTime").
	Text("ManifestYaml").
	OptionalText("Distribution").
	OptionalText("IsMountlessQueryable").
	OptionalText("OrganizationProfileName").
	OptionalText("UniformListingLocator").
	OptionalText("TrialDetails").
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
	OptionalText("IsShare").
	OptionalText("RequestApprovalType").
	OptionalText("MonetizationDisplayOrder").
	OptionalText("LegacyUniformListingLocators")

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
			TextAssignment("AS", g.ParameterOptions().NoEquals().DoubleDollarQuotes().Required()).
			OptionalBooleanAssignment("PUBLISH", g.ParameterOptions()).
			OptionalBooleanAssignment("REVIEW", g.ParameterOptions()).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name"),
	).
	CustomOperation(
		"CreateFromStage",
		"https://docs.snowflake.com/en/sql-reference/sql/create-listing",
		g.NewQueryStruct("CreateListingFromStage").
			Create().
			SQL("EXTERNAL LISTING").
			IfNotExists().
			Name().
			OptionalQueryStructField("With", listingWithDef, g.KeywordOptions()).
			PredefinedQueryStructField("From", "Location", g.ParameterOptions().Required().NoQuotes().NoEquals().SQL("FROM")).
			OptionalBooleanAssignment("PUBLISH", g.ParameterOptions()).
			OptionalBooleanAssignment("REVIEW", g.ParameterOptions()).
			WithValidation(g.ValidIdentifier, "name"),
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
					Text("VersionName", g.KeywordOptions().Required()).
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
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfExists", "AddVersion").
			WithValidation(g.ExactlyOneValueSet, "Publish", "Unpublish", "Review", "AlterListingAs", "AddVersion", "RenameTo", "Set"),
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
	DescribeOperation(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-listing",
		listingDetailsDbRow,
		listingDetails,
		g.NewQueryStruct("DescribeListing").
			Describe().
			SQL("LISTING").
			Name().
			OptionalAssignment("REVISION", g.KindOfT[ListingRevision](), g.ParameterOptions().NoQuotes()).
			WithValidation(g.ValidIdentifier, "name"),
	)

	// TODO: Organization listing will have its interface, but most of the operations will be pass through functions to the Listings interface
	// TODO: Show available listings
	// TODO: Show versions in listing
	// TODO: Describe available listing
	// TODO: Listing manifest builder - https://docs.snowflake.com/en/progaccess/listing-manifest-reference
