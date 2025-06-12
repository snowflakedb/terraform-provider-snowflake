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

var listingFromDef = g.NewQueryStruct("ListingFrom").
	OptionalIdentifier("Share", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("SHARE")).
	OptionalIdentifier("ApplicationPackage", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("APPLICATION PACKAGE")).
	WithValidation(g.ExactlyOneValueSet, "Share", "ApplicationPackage")

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
			PredefinedQueryStructField("From", "*ApplicationPackage", g.KeywordOptions()).
			Text("As", g.KeywordOptions().DoubleDollarQuotes().Required().SQL("AS")).
			OptionalBooleanAssignment("PUBLISH", g.ParameterOptions()).
			OptionalBooleanAssignment("REVIEW", g.ParameterOptions()).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name"),
		listingFromDef,
	).
	CustomOperation(
		"CreateFromStage",
		"https://docs.snowflake.com/en/sql-reference/sql/create-listing",
		g.NewQueryStruct("CreateListingFromStage").
			Create().
			SQL("EXTERNAL LISTING").
			IfNotExists().
			Name().
			PredefinedQueryStructField("From", "Location", g.ParameterOptions().NoQuotes().NoEquals()).
			OptionalBooleanAssignment("PUBLISH", g.ParameterOptions()).
			OptionalBooleanAssignment("REVIEW", g.ParameterOptions()).
			WithValidation(g.ValidIdentifier, "name"),
		listingFromDef,
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
					Text("As", g.KeywordOptions().Required().DoubleDollarQuotes().SQL("AS")).
					OptionalBooleanAssignment("PUBLISH", g.ParameterOptions()).
					OptionalBooleanAssignment("REVIEW", g.ParameterOptions()).
					OptionalComment(),
				g.KeywordOptions(),
			).
			QueryStructField(
				"ADD VERSION",
				g.NewQueryStruct("AddListingVersion").
					IfNotExists().
					OptionalText("VersionName", g.KeywordOptions()).
					PredefinedQueryStructField("From", "Location", g.ParameterOptions().Required().NoQuotes().NoEquals()).
					OptionalComment(),
				g.KeywordOptions(),
			).
			OptionalIdentifier("RenameTo", g.KindOfTPointer[AccountObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("ListingSet").
					OptionalComment(),
				g.KeywordOptions().SQL("SET"),
			).
			//TODO validations
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfExists", "AddVersion").
			WithValidation(g.ExactlyOneValueSet, "Set", "AddVersion"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-listing",
		g.NewQueryStruct("DropTask").
			Drop().
			SQL("LISTING").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-listings",
		taskDbRow, //todo
		task,      //todo
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
		taskDbRow, //todo
		task,      //todo
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
