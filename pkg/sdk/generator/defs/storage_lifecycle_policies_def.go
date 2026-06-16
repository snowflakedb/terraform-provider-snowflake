package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var StorageLifecyclePolicyArchiveTierEnumDef = g.NewEnum(
	"StorageLifecyclePolicyArchiveTier", "StorageLifecyclePolicyArchiveTiers",
	"COOL", "COLD",
)

var storageLifecyclePoliciesDef = g.NewInterface(
	"StorageLifecyclePolicies",
	"StorageLifecyclePolicy",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-storage-lifecycle-policy",
		g.NewQueryStruct("CreateStorageLifecyclePolicy").
			Create().
			OrReplace().
			SQL("STORAGE LIFECYCLE POLICY").
			IfNotExists().
			Name().
			ListQueryStructField(
				"args",
				g.NewQueryStruct("CreateStorageLifecyclePolicyArgs").
					Text("Name", g.KeywordOptions().DoubleQuotes().Required()).
					PredefinedQueryStructField("DataType", "datatypes.DataType", g.ParameterOptions().NoEquals().Required()),
				g.ParameterOptions().Parentheses().SQL("AS").NoEquals().Required(),
			).
			SQL("RETURNS BOOLEAN").
			BodyWithPrecedingArrow().
			OptionalEnumAssignment("ARCHIVE_TIER", StorageLifecyclePolicyArchiveTierEnumDef, g.ParameterOptions()).
			OptionalNumberAssignment("ARCHIVE_FOR_DAYS", g.ParameterOptions()).
			OptionalComment().
			OptionalTags().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidateValueSet, "args").
			WithValidation(g.ValidateValueSet, "body").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-storage-lifecycle-policy",
		g.NewQueryStruct("AlterStorageLifecyclePolicy").
			Alter().
			SQL("STORAGE LIFECYCLE POLICY").
			Name().
			OptionalIdentifier("RenameTo", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
			OptionalSetBodyWithPrecedingArrow().
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("StorageLifecyclePolicySet").
					OptionalEnumAssignment("ARCHIVE_TIER", StorageLifecyclePolicyArchiveTierEnumDef, g.ParameterOptions()).
					OptionalNumberAssignment("ARCHIVE_FOR_DAYS", g.ParameterOptions()).
					OptionalComment().
					WithValidation(g.AtLeastOneValueSet, "ArchiveTier", "ArchiveForDays", "Comment"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalSetTags().
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("StorageLifecyclePolicyUnset").
					OptionalSQL("ARCHIVE_FOR_DAYS").
					OptionalSQL("COMMENT").
					WithValidation(g.AtLeastOneValueSet, "ArchiveForDays", "Comment"),
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).
			OptionalUnsetTags().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "RenameTo", "SetBody", "Set", "SetTags", "Unset", "UnsetTags"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-storage-lifecycle-policy",
		g.NewQueryStruct("DropStorageLifecyclePolicy").
			Drop().
			SQL("STORAGE LIFECYCLE POLICY").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-storage-lifecycle-policies",
		g.StructPair("storageLifecyclePolicyDBRow", "StorageLifecyclePolicy").
			Time("created_on").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("kind").
			Text("owner").
			Text("comment").
			Text("owner_role_type").
			Text("options"),
		g.NewQueryStruct("ShowStorageLifecyclePolicies").
			Show().
			SQL("STORAGE LIFECYCLE POLICIES").
			OptionalLike().
			OptionalExtendedIn(),
		g.ShowByIDExtendedInFiltering,
		g.ShowByIDLikeFiltering,
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-storage-lifecycle-policy",
		g.StructPair("describeStorageLifecyclePolicyDBRow", "StorageLifecyclePolicyDetails").
			Text("name").
			Field("signature", "string", "[]TableColumnSignature", g.WithCustomParser("ParseTableColumnSignatureWithVectorSupport")).
			DataType("return_type").
			Text("body").
			OptionalNumber("archive_for_days").
			Text("archive_tier", g.WithValueAdjuster("normalizeStorageLifecyclePolicyArchiveTier")),
		g.NewQueryStruct("DescribeStorageLifecyclePolicy").
			Describe().
			SQL("STORAGE LIFECYCLE POLICY").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	WithEnums(StorageLifecyclePolicyArchiveTierEnumDef)
