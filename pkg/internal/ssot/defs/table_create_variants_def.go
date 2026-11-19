package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// CloneMoment is AT or BEFORE in CREATE TABLE ... CLONE ... AT|BEFORE (...).
var CloneMoment = ssot.NewEnum("CloneMoment", "CloneMoments", "AT", "BEFORE")

var tableCreateVariantAttrs = struct {
	AsSelectColumns    ssot.Attribute
	AsSelectQuery      ssot.Attribute
	UsingTemplateQuery ssot.Attribute
	LikeSourceTable    ssot.Attribute
	CloneSourceTable   ssot.Attribute
	ClonePoint         ssot.Attribute
}{
	AsSelectColumns: ssot.Attribute{
		Name:      "as_select_column",
		Type:      ssot.AttrList,
		Required:  true,
		SqlOnly:   true,
		CreateSql: ssot.ListOfObjects(ssot.Parentheses(), ssot.FieldName("Columns")),
		Elem: &ssot.Attribute{
			Type:      ssot.AttrObject,
			CreateSql: ssot.ObjectParameter(),
			Children: []ssot.Attribute{
				{Name: "name", Type: ssot.AttrString, Required: true, CreateSql: ssot.TextKeyword(ssot.DoubleQuotes())},
				{Name: "column_type", Type: ssot.AttrString, CreateSql: ssot.PredefinedType("datatypes.DataType", ssot.NoQuotes(), ssot.NoEquals())},
				{Name: "masking_policy", Type: ssot.AttrIdentifier, IdentifierType: "SchemaObjectIdentifier", CreateSql: ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("MASKING POLICY"))},
			},
		},
	},
	AsSelectQuery: ssot.Attribute{
		Name:      "as_select_query",
		Type:      ssot.AttrString,
		Required:  true,
		SqlOnly:   true,
		CreateSql: ssot.PlainParameter(ssot.NoEquals(), ssot.Sql("AS"), ssot.FieldName("Query")),
	},
	UsingTemplateQuery: ssot.Attribute{
		Name:      "using_template_query",
		Type:      ssot.AttrList,
		SqlOnly:   true,
		Elem:      &ssot.Attribute{Type: ssot.AttrString},
		CreateSql: ssot.ListParameter("string", ssot.NoEquals(), ssot.FieldName("Query"), ssot.Sql("USING TEMPLATE")),
	},
	LikeSourceTable: ssot.Attribute{
		Name:           "like_source_table",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "SchemaObjectIdentifier",
		Required:       true,
		SqlOnly:        true,
		CreateSql:      ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("LIKE"), ssot.FieldName("SourceTable")),
	},
	CloneSourceTable: ssot.Attribute{
		Name:           "clone_source_table",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "SchemaObjectIdentifier",
		Required:       true,
		SqlOnly:        true,
		CreateSql:      ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("CLONE"), ssot.FieldName("SourceTable")),
	},
	ClonePoint: ssot.Attribute{
		Name:      "clone_point",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(),
		Children: []ssot.Attribute{
			{Name: "moment", Type: ssot.AttrEnum, Enum: CloneMoment, Required: true, SqlOnly: true, CreateSql: ssot.EnumKeyword()},
			{Name: "at", SqlOnly: true, CreateSql: ssot.PredefinedType("TimeTravel", ssot.Parentheses(), ssot.NoComma(), ssot.Required())},
		},
	},
}

func tableCreateVariantClauses(rest ...ssot.CustomOperationClause) []ssot.CustomOperationClause {
	return append([]ssot.CustomOperationClause{
		ssot.CustomStaticSQL("CREATE"),
		ssot.CustomOrReplace(),
		ssot.CustomObjectWords(),
		ssot.CustomName(),
	}, rest...)
}

// CREATE AS SELECT, USING TEMPLATE, LIKE, and CLONE are extra CREATE shapes of Table.
var tableCreateVariantOperations = []ssot.CustomOperationSpec{
	{
		Kind:            "CreateAsSelect",
		Doc:             "https://docs.snowflake.com/en/sql-reference/sql/create-table",
		QueryStructName: "CreateTableAsSelect",
		Clauses: tableCreateVariantClauses(
			ssot.CustomAttribute(&tableCreateVariantAttrs.AsSelectColumns),
			ssot.CustomAttribute(&TableAttrs.ClusterBy),
			ssot.CustomAttribute(&TableCreateAttrs.Copy),
			ssot.CustomAttribute(&TableCreateAttrs.RowAccessPolicy),
			ssot.CustomAttribute(&tableCreateVariantAttrs.AsSelectQuery),
		),
	},
	{
		Kind:            "CreateUsingTemplate",
		Doc:             "https://docs.snowflake.com/en/sql-reference/sql/create-table#create-table-using-template",
		QueryStructName: "CreateTableUsingTemplate",
		Clauses: tableCreateVariantClauses(
			ssot.CustomAttribute(&TableCreateAttrs.Copy),
			ssot.CustomAttribute(&tableCreateVariantAttrs.UsingTemplateQuery),
		),
	},
	{
		Kind:            "CreateLike",
		Doc:             "https://docs.snowflake.com/en/sql-reference/sql/create-table#create-table-like",
		QueryStructName: "CreateTableLike",
		Clauses: tableCreateVariantClauses(
			ssot.CustomAttribute(&tableCreateVariantAttrs.LikeSourceTable),
			ssot.CustomAttribute(&TableAttrs.ClusterBy),
			ssot.CustomAttribute(&TableCreateAttrs.Copy),
		),
	},
	{
		Kind:            "CreateClone",
		Doc:             "https://docs.snowflake.com/en/sql-reference/sql/create-clone",
		QueryStructName: "CreateTableClone",
		Clauses: tableCreateVariantClauses(
			ssot.CustomAttribute(&tableCreateVariantAttrs.CloneSourceTable),
			ssot.CustomAttribute(&tableCreateVariantAttrs.ClonePoint),
			ssot.CustomAttribute(&TableCreateAttrs.Copy),
		),
	},
	// TODO [SNOW-4532966]: add the CREATE TABLE ... FROM BACKUP SET <backup_set> IDENTIFIER '<identifier>' variant (CreateFromBackupSet).
}
