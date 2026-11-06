package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// Describe and select methods stay custom operations. SdkDescribe is one DESCRIBE per object.
// CREATE AS SELECT, CREATE USING TEMPLATE, CREATE LIKE, and CREATE CLONE stay handwritten,
// so SdkDefineFromSsot stays off. Follow-up: convert those four and turn the setting on.

func tableDescribeText(name string) *ssot.Attribute {
	return &ssot.Attribute{Name: name, Type: ssot.AttrString, Readonly: true}
}

func tableDescribeOptionalText(name string) *ssot.Attribute {
	return &ssot.Attribute{Name: name, Type: ssot.AttrString, OptionalInShow: true, OptionalInDb: true, Readonly: true}
}

func tableDescribeBoolText(name, trueValue string) *ssot.Attribute {
	return &ssot.Attribute{
		Name: name, Type: ssot.AttrBool, StoredAsText: true, BoolTrueValue: trueValue, Readonly: true,
	}
}

func tableDescribeField(name, plainKind string, manual bool, parser string) *ssot.Attribute {
	return &ssot.Attribute{
		Name:            name,
		Readonly:        true,
		PairedPlainKind: plainKind,
		ManualConvert:   manual,
		CustomParser:    parser,
	}
}

var tableDescribeOperations = []ssot.CustomOperationSpec{
	{
		Kind:            "DescribeColumns",
		Doc:             "https://docs.snowflake.com/en/sql-reference/sql/desc-table",
		QueryStructName: "DescribeTableColumns",
		Clauses: []ssot.CustomOperationClause{
			ssot.CustomStaticSQL("DESCRIBE"),
			ssot.CustomStaticSQL("TABLE"),
			ssot.CustomName(),
			ssot.CustomStaticField("columnsType", "TYPE = COLUMNS"),
		},
		Paired: &ssot.CustomPairedShow{
			DbRow: "tableColumnDetailsRow",
			Plain: "TableColumnDetails",
			Columns: []*ssot.Attribute{
				tableDescribeText("name"),
				tableDescribeField("type", "DataType", true, ""),
				tableDescribeText("kind"),
				{Name: "null?", Type: ssot.AttrBool, StoredAsText: true, BoolTrueValue: "Y", Readonly: true, SchemaName: "is_nullable"},
				tableDescribeOptionalText("default"),
				{Name: "primary key", Type: ssot.AttrBool, StoredAsText: true, BoolTrueValue: "Y", Readonly: true, SchemaName: "is_primary"},
				{Name: "unique key", Type: ssot.AttrBool, StoredAsText: true, BoolTrueValue: "Y", Readonly: true, SchemaName: "is_unique"},
				{Name: "check", Type: ssot.AttrBool, StoredAsText: true, OptionalInShow: true, OptionalInDb: true, BoolTrueValue: "Y", Readonly: true},
				tableDescribeOptionalText("expression"),
				tableDescribeOptionalText("comment"),
				tableDescribeOptionalText("policy name"),
				{Name: "Collation", PlainOnlyKind: "*string"},
				tableDescribeOptionalText("schema evolution record"),
			},
		},
	},
	{
		Kind:            "DescribeStage",
		Doc:             "https://docs.snowflake.com/en/sql-reference/sql/desc-table",
		QueryStructName: "DescribeTableStage",
		Clauses: []ssot.CustomOperationClause{
			ssot.CustomStaticSQL("DESCRIBE"),
			ssot.CustomStaticSQL("TABLE"),
			ssot.CustomName(),
			ssot.CustomStaticField("stageType", "TYPE = STAGE"),
		},
		Paired: &ssot.CustomPairedShow{
			DbRow: "tableStageDetailsRow",
			Plain: "TableStageDetails",
			Columns: []*ssot.Attribute{
				tableDescribeText("parent_property"),
				tableDescribeText("property"),
				tableDescribeText("property_type"),
				tableDescribeText("property_value"),
				tableDescribeText("property_default"),
			},
		},
	},
	{
		Kind:            "DescribeSearchOptimization",
		Doc:             "https://docs.snowflake.com/en/sql-reference/sql/desc-search-optimization",
		QueryStructName: "DescribeSearchOptimization",
		Clauses: []ssot.CustomOperationClause{
			ssot.CustomStaticSQL("DESCRIBE"),
			ssot.CustomStaticSQL("SEARCH OPTIMIZATION"),
			ssot.CustomStaticSQL("ON"),
			ssot.CustomName(),
		},
		Paired: &ssot.CustomPairedShow{
			DbRow: "tableSearchOptimizationDetailsRow",
			Plain: "TableSearchOptimizationDetails",
			Columns: []*ssot.Attribute{
				{Name: "expression_id", Type: ssot.AttrInt, Readonly: true},
				tableDescribeText("method"),
				tableDescribeText("target"),
				tableDescribeField("target_data_type", "datatypes.DataType", false, "datatypes.ParseDataType"),
				tableDescribeBoolText("active", "true"),
			},
		},
	},
	{
		Kind:            "SelectTableConstraints",
		Doc:             "https://docs.snowflake.com/en/sql-reference/info-schema/table_constraints",
		QueryStructName: "ShowTableConstraints",
		Clauses: []ssot.CustomOperationClause{
			ssot.CustomStaticField("selectAll", "SELECT * FROM"),
			ssot.CustomIdentifier("Database", "AccountObjectIdentifier", true),
			ssot.CustomStaticField("dot", "."),
			ssot.CustomStaticField("informationSchemaTableConstraints", "INFORMATION_SCHEMA.TABLE_CONSTRAINTS"),
			ssot.CustomStaticField("where", "WHERE"),
			ssot.CustomText("TABLE_SCHEMA"),
			ssot.CustomStaticField("and", "AND"),
			ssot.CustomText("TABLE_NAME"),
		},
		Paired: &ssot.CustomPairedShow{
			DbRow: "tableConstraintDetailsRow",
			Plain: "TableConstraintDetails",
			Columns: []*ssot.Attribute{
				tableDescribeText("CONSTRAINT_CATALOG"),
				tableDescribeText("CONSTRAINT_SCHEMA"),
				tableDescribeText("CONSTRAINT_NAME"),
				tableDescribeText("TABLE_CATALOG"),
				tableDescribeText("TABLE_SCHEMA"),
				tableDescribeText("TABLE_NAME"),
				{Name: "CONSTRAINT_TYPE", Type: ssot.AttrEnum, Enum: TableConstraintType, Readonly: true},
				tableDescribeBoolText("IS_DEFERRABLE", "YES"),
				tableDescribeBoolText("INITIALLY_DEFERRED", "YES"),
				tableDescribeOptionalText("COMMENT"),
				{Name: "CREATED", Type: ssot.AttrTime, Readonly: true},
				{Name: "LAST_ALTERED", Type: ssot.AttrTime, Readonly: true},
				tableDescribeBoolText("ENFORCED", "YES"),
				tableDescribeBoolText("RELY", "YES"),
			},
		},
	},
	{
		Kind:            "SelectCheckConstraints",
		Doc:             "https://docs.snowflake.com/en/sql-reference/info-schema/check_constraints",
		QueryStructName: "SelectCheckConstraints",
		Clauses: []ssot.CustomOperationClause{
			ssot.CustomStaticField("selectAll", "SELECT * FROM"),
			ssot.CustomIdentifier("Database", "AccountObjectIdentifier", true),
			ssot.CustomStaticField("dot", "."),
			ssot.CustomStaticField("informationSchemaCheckConstraints", "INFORMATION_SCHEMA.CHECK_CONSTRAINTS"),
			ssot.CustomStaticField("where", "WHERE"),
			ssot.CustomText("CONSTRAINT_SCHEMA"),
			ssot.CustomStaticField("and", "AND"),
			ssot.CustomText("CONSTRAINT_TABLE"),
		},
		Paired: &ssot.CustomPairedShow{
			DbRow: "tableCheckConstraintDetailsRow",
			Plain: "TableCheckConstraintDetails",
			Columns: []*ssot.Attribute{
				tableDescribeText("CONSTRAINT_CATALOG"),
				tableDescribeText("CONSTRAINT_SCHEMA"),
				tableDescribeText("CONSTRAINT_TABLE"),
				tableDescribeText("CONSTRAINT_NAME"),
				tableDescribeText("CHECK_CLAUSE"),
			},
		},
	},
}
