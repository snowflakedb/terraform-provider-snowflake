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
	Text("published_on").
	Text("state").
	Text("review_state").
	Text("comment").
	Text("owner").
	Text("owner_role_type").
	Text("regions").
	Text("target_accounts").
	Text("is_monetized").
	Text("is_application").
	Text("is_targeted")

var listing = g.PlainStruct("Listing").
	Text("GlobalName").
	Text("Name").
	Text("Title").
	Text("Subtitle").
	Text("Profile").
	Text("CreatedOn").
	Text("UpdatedOn").
	Text("PublishedOn").
	Text("State").
	Text("ReviewState").
	Text("Comment").
	Text("Owner").
	Text("OwnerRoleType").
	Text("Regions").
	Text("TargetAccounts").
	Text("IsMonetized").
	Text("IsApplication").
	Text("IsTargeted")

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
	Text("data_attributes"). // TODO: Not documented
	Text("categories").
	Text("resources").
	Text("profile").
	Text("customized_contact_info").
	Text("data_dictionary").
	Text("data_preview"). // TODO: Not documented
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
	Text("is_limited_trial").    // TODO: Not documented
	Text("is_by_request").       // TODO: Not documented
	Text("limited_trial_plan").  // TODO: Not documented
	Text("retried_on").          // TODO: Not documented
	Text("scheduled_drop_time"). // TODO: Not documented
	Text("manifest_yaml").
	Text("distribution").                   // TODO: Not documented
	Text("is_mountless_queryable").         // TODO: Not documented
	Text("organization_profile_name").      // TODO: Not documented
	Text("uniform_listing_locator").        // TODO: Not documented
	Text("trial_details").                  // TODO: Not documented
	Text("trial_details").                  // TODO: Not documented
	Text("approver_contact").               // TODO: Not documented
	Text("support_contact").                // TODO: Not documented
	Text("live_version_uri").               // TODO: Not documented
	Text("last_committed_version_uri").     // TODO: Not documented
	Text("last_committed_version_name").    // TODO: Not documented
	Text("last_committed_version_alias").   // TODO: Not documented
	Text("published_version_uri").          // TODO: Not documented
	Text("published_version_name").         // TODO: Not documented
	Text("published_version_alias").        // TODO: Not documented
	Text("is_share").                       // TODO: Not documented
	Text("request_approval_type").          // TODO: Not documented
	Text("monetization_display_order").     // TODO: Not documented
	Text("legacy_uniform_listing_locators") // TODO: Not documented

var listingDetails = g.PlainStruct("ListingDetails").
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
	Text("data_attributes").
	Text("categories").
	Text("resources").
	Text("profile").
	Text("customized_contact_info").
	Text("data_dictionary").
	Text("data_preview").
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
	Text("is_limited_trial").
	Text("is_by_request").
	Text("limited_trial_plan").
	Text("retried_on").
	Text("scheduled_drop_time").
	Text("manifest_yaml").
	Text("distribution").
	Text("is_mountless_queryable").
	Text("organization_profile_name").
	Text("uniform_listing_locator").
	Text("trial_details").
	Text("trial_details").
	Text("approver_contact").
	Text("support_contact").
	Text("live_version_uri").
	Text("last_committed_version_uri").
	Text("last_committed_version_name").
	Text("last_committed_version_alias").
	Text("published_version_uri").
	Text("published_version_name").
	Text("published_version_alias").
	Text("is_share").
	Text("request_approval_type").
	Text("monetization_display_order").
	Text("legacy_uniform_listing_locators")

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
