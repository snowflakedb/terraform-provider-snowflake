package sdk

import (
	"fmt"
	"strings"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

type StreamSourceType string

const (
	StreamSourceTypeTable         StreamSourceType = "TABLE"
	StreamSourceTypeExternalTable StreamSourceType = "EXTERNAL TABLE"
	StreamSourceTypeView          StreamSourceType = "VIEW"
	StreamSourceTypeStage         StreamSourceType = "STAGE"
)

func ToStreamSourceType(s string) (StreamSourceType, error) {
	switch streamSourceType := StreamSourceType(strings.ToUpper(s)); streamSourceType {
	case StreamSourceTypeTable,
		StreamSourceTypeExternalTable,
		StreamSourceTypeView,
		StreamSourceTypeStage:
		return streamSourceType, nil
	default:
		return "", fmt.Errorf("invalid stream source type: %s", s)
	}
}

type StreamMode string

const (
	StreamModeDefault    StreamMode = "DEFAULT"
	StreamModeAppendOnly StreamMode = "APPEND_ONLY"
	StreamModeInsertOnly StreamMode = "INSERT_ONLY"
)

func ToStreamMode(s string) (StreamMode, error) {
	switch streamMode := StreamMode(strings.ToUpper(s)); streamMode {
	case StreamModeDefault,
		StreamModeAppendOnly,
		StreamModeInsertOnly:
		return streamMode, nil
	default:
		return "", fmt.Errorf("invalid stream mode: %s", s)
	}
}

var (
	onStreamDef = g.NewQueryStruct("OnStream").
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

	showStreamDbRowDef = g.DbStruct("showStreamsDbRow").
				Field("created_on", "time.Time").
				Field("name", "string").
				Field("database_name", "string").
				Field("schema_name", "string").
				Field("owner", "sql.NullString").
				Field("comment", "sql.NullString").
				Field("table_name", "sql.NullString").
				Field("source_type", "sql.NullString").
				Field("base_tables", "sql.NullString").
				Field("type", "sql.NullString").
				Field("stale", "string").
				Field("mode", "sql.NullString").
				Field("stale_after", "sql.NullTime").
				Field("invalid_reason", "sql.NullString").
				Field("owner_role_type", "sql.NullString")

	streamPlainStructDef = g.PlainStruct("Stream").
				Field("CreatedOn", "time.Time").
				Field("Name", "string").
				Field("DatabaseName", "string").
				Field("SchemaName", "string").
				Field("Owner", "*string").
				Field("Comment", "*string").
				Field("TableName", "*string").
				Field("SourceType", "*StreamSourceType").
				Field("BaseTables", "[]string").
				Field("Type", "*string").
				Field("Stale", "bool").
				Field("Mode", "*StreamMode").
				Field("StaleAfter", "*time.Time").
				Field("InvalidReason", "*string").
				Field("OwnerRoleType", "*string")

	StreamsDef = g.NewInterface(
		"Streams",
		"Stream",
		g.KindOfT[SchemaObjectIdentifier](),
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
				Identifier("TableId", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
				OptionalQueryStructField("On", onStreamDef, g.KeywordOptions()).
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
				Identifier("ExternalTableId", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
				OptionalQueryStructField("On", onStreamDef, g.KeywordOptions()).
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
				Identifier("StageId", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
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
				Identifier("ViewId", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
				OptionalQueryStructField("On", onStreamDef, g.KeywordOptions()).
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
				Identifier("sourceStream", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("CLONE").Required()).
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
		ShowOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/show-streams",
			showStreamDbRowDef,
			streamPlainStructDef,
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
		DescribeOperation(
			g.DescriptionMappingKindSingleValue,
			"https://docs.snowflake.com/en/sql-reference/sql/desc-stream",
			showStreamDbRowDef,
			streamPlainStructDef,
			g.NewQueryStruct("DescribeStream").
				Describe().
				SQL("STREAM").
				Name().
				WithValidation(g.ValidIdentifier, "name"),
		)
)
