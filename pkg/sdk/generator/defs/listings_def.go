package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var (
	ListingRevisionEnumDef = g.NewEnum(
		"ListingRevision", "ListingRevisions",
		"DRAFT", "PUBLISHED",
	)
	ListingStateEnumDef = g.NewEnum(
		"ListingState", "ListingStates",
		"DRAFT", "PUBLISHED", "UNPUBLISHED",
	)
)

var listingWithDef = g.NewQueryStruct("ListingWith").
	OptionalIdentifier("Share", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("SHARE")).
	OptionalIdentifier("ApplicationPackage", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("APPLICATION PACKAGE")).
	WithValidation(g.ExactlyOneValueSet, "Share", "ApplicationPackage")

// There are more fields listed than in https://docs.snowflake.com/en/sql-reference/sql/show-listings.
// They are mapped straight from the SHOW LISTINGS output.
var listingPairs = g.StructPair("listingDBRow", "Listing").
	Text("global_name").
	Text("name").
	Text("title").
	OptionalText("subtitle").
	Text("profile").
	Text("created_on").
	Text("updated_on").
	OptionalText("published_on").
	Enum("state", ListingStateEnumDef).
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

// There are more fields listed than in https://docs.snowflake.com/en/sql-reference/sql/desc-listing
// They are mapped straight from the DESC LISTING output.
var listingDetailsPairs = g.StructPair("listingDetailsDBRow", "ListingDetails").
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
	Enum("state", ListingStateEnumDef).
	OptionalAccountObjectIdentifier("share", g.WithPlainFieldName("Share")).
	OptionalAccountObjectIdentifier("application_package", g.WithPlainFieldName("ApplicationPackage")).
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

var listingVersionPairs = g.StructPair("listingVersionDBRow", "ListingVersion").
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

var listingsDef = g.NewInterface(
	"Listings",
	"Listing",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
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
			PredefinedQueryStructField("From", g.KindOfTPointer[sdkcommons.Location](), g.ParameterOptions().NoQuotes().NoEquals().SQL("FROM")).
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
			OptionalIdentifier("RenameTo", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
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
		g.WithDropSafelyHook(),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-listings",
		listingPairs,
		g.NewQueryStruct("ShowListings").
			Show().
			SQL("LISTINGS").
			OptionalLike().
			OptionalStartsWith().
			OptionalLimitFrom(),
	).
	CustomShowOperationWithPairedStructs(
		"Describe",
		g.ShowMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-listing",
		listingDetailsPairs,
		g.NewQueryStruct("DescribeListing").
			Describe().
			SQL("LISTING").
			Name().
			OptionalEnumAssignment("REVISION", ListingRevisionEnumDef, g.ParameterOptions().NoQuotes()).
			WithValidation(g.ValidIdentifier, "name"),
	).
	CustomShowOperationWithPairedStructs(
		"ShowVersions",
		g.ShowMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-listing",
		listingVersionPairs,
		g.NewQueryStruct("ShowListings").
			Show().
			SQL("VERSIONS IN LISTING").
			Name().
			OptionalLimit().
			WithValidation(g.ValidIdentifier, "name"),
	).
	WithEnums(
		ListingRevisionEnumDef,
		ListingStateEnumDef,
	)

	// TODO [SNOW-2236968]: Organization listing may have its interface, but most of the operations would be pass through functions to the Listings interface
	// TODO [SNOW-2236968]: Show available listings
	// TODO [SNOW-2236968]: Describe available listing
	// TODO [SNOW-2236968]: Listing manifest builder - https://docs.snowflake.com/en/progaccess/listing-manifest-reference
	// TODO [SNOW-2236968]: Test mapping functions (ToListingRevision and ToListingState)
