package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var (
	StreamSourceTypeEnumDef = g.NewEnum(
		"StreamSourceType", "StreamSourceTypes",
		"TABLE", "EXTERNAL TABLE", "VIEW", "STAGE",
	)
	StreamModeEnumDef = g.NewEnum(
		"StreamMode", "StreamModes",
		"DEFAULT", "APPEND_ONLY", "INSERT_ONLY",
	)

	onStreamDef = func() *g.QueryStruct {
		return g.NewQueryStruct("OnStream").
			OptionalSQL("AT").
			OptionalSQL("BEFORE").
			QueryStructField(
				"Statement",
				g.NewQueryStruct("OnStreamStatement").
					OptionalTextAssignment("TIMESTAMP", g.ParameterOptions().ArrowEquals().SingleQuotes()).
					OptionalTextAssignment("OFFSET", g.ParameterOptions().ArrowEquals()).
					OptionalTextAssignment("STATEMENT", g.ParameterOptions().ArrowEquals().SingleQuotes()).
					OptionalTextAssignment("STREAM", g.ParameterOptions().ArrowEquals().SingleQuotes()).
					WithValidation(g.ExactlyOneValueSet, "Timestamp", "Offset", "Statement", "Stream"),
				g.ListOptions().Parentheses(),
			).
			WithValidation(g.ExactlyOneValueSet, "At", "Before")
	}

	streamPairs = g.StructPair("showStreamsDbRow", "Stream").
			Time("created_on").
			Text("name").
			Text("database_name").
			Text("schema_name").
			OptionalText("owner").
			OptionalText("comment").
			OptionalText("table_name").
			Field("source_type", "sql.NullString", StreamSourceTypeEnumDef.KindPtr()).
			Field("base_tables", "sql.NullString", "[]string").
			OptionalText("type").
			Field("stale", "string", "bool").
			Field("mode", "sql.NullString", StreamModeEnumDef.KindPtr()).
			OptionalTime("stale_after").
			OptionalText("invalid_reason").
			OptionalText("owner_role_type")

	streamsDef = g.NewInterface(
		"Streams",
		"Stream",
		g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
	).
		CustomOperation(
			"CreateOnTable",
			"https://docs.snowflake.com/en/sql-reference/sql/create-stream",
			g.NewQueryStruct("CreateStreamOnTable").
				Create().
				OrReplace().
				SQL("STREAM").
				IfNotExists().
				Name().
				OptionalTags().
				OptionalCopyGrants().
				SQL("ON TABLE").
				Identifier("TableId", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
				OptionalQueryStructField("On", onStreamDef(), g.KeywordOptions()).
				OptionalBooleanAssignment("APPEND_ONLY", nil).
				OptionalBooleanAssignment("SHOW_INITIAL_ROWS", nil).
				OptionalComment().
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ValidIdentifier, "TableId").
				WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
		).
		CustomOperation(
			"CreateOnExternalTable",
			"https://docs.snowflake.com/en/sql-reference/sql/create-stream",
			g.NewQueryStruct("CreateStreamOnExternalTable").
				Create().
				OrReplace().
				SQL("STREAM").
				IfNotExists().
				Name().
				OptionalTags().
				OptionalCopyGrants().
				SQL("ON EXTERNAL TABLE").
				Identifier("ExternalTableId", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
				OptionalQueryStructField("On", onStreamDef(), g.KeywordOptions()).
				OptionalBooleanAssignment("INSERT_ONLY", nil).
				OptionalComment().
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ValidIdentifier, "ExternalTableId").
				WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
		).
		CustomOperation(
			"CreateOnDirectoryTable",
			"https://docs.snowflake.com/en/sql-reference/sql/create-stream",
			g.NewQueryStruct("CreateStreamOnDirectoryTable").
				Create().
				OrReplace().
				SQL("STREAM").
				IfNotExists().
				Name().
				OptionalTags().
				OptionalCopyGrants().
				SQL("ON STAGE").
				Identifier("StageId", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
				OptionalComment().
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ValidIdentifier, "StageId").
				WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
		).
		CustomOperation(
			"CreateOnView",
			"https://docs.snowflake.com/en/sql-reference/sql/create-stream",
			g.NewQueryStruct("CreateStreamOnView").
				Create().
				OrReplace().
				SQL("STREAM").
				IfNotExists().
				Name().
				OptionalTags().
				OptionalCopyGrants().
				SQL("ON VIEW").
				Identifier("ViewId", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
				OptionalQueryStructField("On", onStreamDef(), g.KeywordOptions()).
				OptionalBooleanAssignment("APPEND_ONLY", nil).
				OptionalBooleanAssignment("SHOW_INITIAL_ROWS", nil).
				OptionalComment().
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ValidIdentifier, "ViewId").
				WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
		).
		CustomOperation(
			"Clone",
			"https://docs.snowflake.com/en/sql-reference/sql/create-stream#variant-syntax",
			g.NewQueryStruct("CloneStream").
				Create().
				OrReplace().
				SQL("STREAM").
				Name().
				Identifier("sourceStream", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("CLONE").Required()).
				OptionalCopyGrants().
				WithValidation(g.ValidIdentifier, "name"),
		).
		AlterOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/alter-stream",
			g.NewQueryStruct("AlterStream").
				Alter().
				SQL("STREAM").
				IfExists().
				Name().
				OptionalTextAssignment("SET COMMENT", g.ParameterOptions().SingleQuotes()).
				OptionalSQL("UNSET COMMENT").
				OptionalSetTags().
				OptionalUnsetTags().
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ConflictingFields, "IfExists", "UnsetTags").
				WithValidation(g.ExactlyOneValueSet, "SetComment", "UnsetComment", "SetTags", "UnsetTags"),
		).
		DropOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/drop-stream",
			g.NewQueryStruct("DropStream").
				Drop().
				SQL("STREAM").
				IfExists().
				Name().
				WithValidation(g.ValidIdentifier, "name"),
		).
		ShowOperationWithPairedStructs(
			"https://docs.snowflake.com/en/sql-reference/sql/show-streams",
			streamPairs,
			g.NewQueryStruct("ShowStreams").
				Show().
				Terse().
				SQL("STREAMS").
				OptionalLike().
				OptionalExtendedIn().
				OptionalStartsWith().
				OptionalLimit(),
		).
		ShowByIdOperationWithFiltering(
			g.ShowByIDExtendedInFiltering,
			g.ShowByIDLikeFiltering,
		).
		DescribeOperationWithPairedStructs(
			g.DescriptionMappingKindSingleValue,
			"https://docs.snowflake.com/en/sql-reference/sql/desc-stream",
			streamPairs,
			g.NewQueryStruct("DescribeStream").
				Describe().
				SQL("STREAM").
				Name().
				WithValidation(g.ValidIdentifier, "name"),
		).
		WithEnums(
			StreamSourceTypeEnumDef,
			StreamModeEnumDef,
		)
)
