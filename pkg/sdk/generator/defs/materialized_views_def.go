package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var materializedViewColumn = g.NewQueryStruct("MaterializedViewColumn").
	Text("Name", g.KeywordOptions().DoubleQuotes().Required()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes().NoEquals())

var materializedViewColumnMaskingPolicy = g.NewQueryStruct("MaterializedViewColumnMaskingPolicy").
	Text("Name", g.KeywordOptions().Required()).
	Identifier("MaskingPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("MASKING POLICY").Required()).
	NamedListWithParens("USING", g.KindOfT[string](), nil). // TODO: double quotes here?
	OptionalTags()

var materializedViewRowAccessPolicy = g.NewQueryStruct("MaterializedViewRowAccessPolicy").
	Identifier("RowAccessPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("ROW ACCESS POLICY").Required()).
	NamedListWithParens("ON", g.KindOfT[string](), g.KeywordOptions().Required()). // TODO: double quotes here?
	WithValidation(g.ValidIdentifier, "RowAccessPolicy").
	WithValidation(g.ValidateValueSet, "On")

var materializedViewClusterByExpression = g.NewQueryStruct("MaterializedViewClusterByExpression").
	Text("Name", g.KeywordOptions().DoubleQuotes().Required())

var materializedViewClusterBy = func() *g.QueryStruct {
	return g.NewQueryStruct("MaterializedViewClusterBy").
		SQL("CLUSTER BY").
		ListQueryStructField("Expressions", materializedViewClusterByExpression, g.ListOptions().Parentheses()).
		WithValidation(g.ValidateValueSet, "Expressions")
}

var materializedViewSet = g.NewQueryStruct("MaterializedViewSet").
	OptionalSQL("SECURE").
	OptionalComment().
	WithValidation(g.ExactlyOneValueSet, "Secure", "Comment")

var materializedViewUnset = g.NewQueryStruct("MaterializedViewUnset").
	OptionalSQL("SECURE").
	OptionalSQL("COMMENT").
	WithValidation(g.ExactlyOneValueSet, "Secure", "Comment")

var materializedViewPairs = g.StructPair("materializedViewDBRow", "MaterializedView").
	Text("created_on").
	Text("name").
	OptionalText("reserved").
	Text("database_name").
	Text("schema_name").
	OptionalText("cluster_by", g.WithRequiredInPlain()).
	Number("rows").
	Number("bytes").
	Text("source_database_name").
	Text("source_schema_name").
	Text("source_table_name").
	Time("refreshed_on").
	Time("compacted_on").
	Text("owner").
	Bool("invalid").
	OptionalText("invalid_reason", g.WithRequiredInPlain()).
	Text("behind_by").
	OptionalText("comment", g.WithRequiredInPlain()).
	Text("text").
	Bool("is_secure").
	Field("automatic_clustering", "string", "bool", g.WithBoolTrueValue("ON")).
	OptionalText("owner_role_type", g.WithRequiredInPlain()).
	OptionalText("budget", g.WithRequiredInPlain()).
	WithConvertGeneration()

var materializedViewDetailsPairs = g.StructPair("materializedViewDetailsRow", "MaterializedViewDetails").
	Text("name").
	Field("type", "DataType", "DataType").
	Text("kind").
	Field("null?", "string", "bool", g.WithDbFieldName("Null"), g.WithPlainFieldName("IsNullable")).
	OptionalText("default").
	Field("primary key", "string", "bool", g.WithPlainFieldName("IsPrimary")).
	Field("unique key", "string", "bool", g.WithPlainFieldName("IsUnique")).
	Field("check", "sql.NullString", "*bool").
	OptionalText("expression").
	OptionalText("comment").
	WithConvertGeneration()

var materializedViewsDef = g.NewInterface(
	"MaterializedViews",
	"MaterializedView",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-materialized-view",
		g.NewQueryStruct("CreateMaterializedView").
			Create().
			OrReplace().
			OptionalSQL("SECURE").
			SQL("MATERIALIZED VIEW").
			IfNotExists().
			Name().
			OptionalCopyGrants().
			ListQueryStructField("Columns", materializedViewColumn, g.ListOptions().Parentheses()).
			ListQueryStructField("ColumnsMaskingPolicies", materializedViewColumnMaskingPolicy, g.ListOptions().NoParentheses().NoEquals()).
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			OptionalQueryStructField("RowAccessPolicy", materializedViewRowAccessPolicy, g.KeywordOptions()).
			OptionalTags().
			OptionalQueryStructField("ClusterBy", materializedViewClusterBy(), g.KeywordOptions()).
			SQL("AS").
			Text("sql", g.KeywordOptions().NoQuotes().Required()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-materialized-view",
		g.NewQueryStruct("AlterMaterializedView").
			Alter().
			SQL("MATERIALIZED VIEW").
			Name().
			OptionalIdentifier("RenameTo", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
			OptionalQueryStructField("ClusterBy", materializedViewClusterBy(), g.KeywordOptions()).
			OptionalSQL("DROP CLUSTERING KEY").
			OptionalSQL("SUSPEND RECLUSTER").
			OptionalSQL("RESUME RECLUSTER").
			OptionalSQL("SUSPEND").
			OptionalSQL("RESUME").
			OptionalQueryStructField("Set", materializedViewSet, g.KeywordOptions().SQL("SET")).
			OptionalQueryStructField("Unset", materializedViewUnset, g.KeywordOptions().SQL("UNSET")).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "RenameTo", "ClusterBy", "DropClusteringKey", "SuspendRecluster", "ResumeRecluster", "Suspend", "Resume", "Set", "Unset"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-materialized-view",
		g.NewQueryStruct("DropMaterializedView").
			Drop().
			SQL("MATERIALIZED VIEW").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-materialized-views",
		materializedViewPairs,
		g.NewQueryStruct("ShowMaterializedViews").
			Show().
			SQL("MATERIALIZED VIEWS").
			OptionalLike().
			OptionalIn(),
		g.ShowByIDInFiltering,
		g.ShowByIDLikeFiltering,
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-materialized-view",
		materializedViewDetailsPairs,
		g.NewQueryStruct("DescribeMaterializedView").
			Describe().
			SQL("MATERIALIZED VIEW").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
