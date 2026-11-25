package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// TableDataMetricScheduleStatusOperation is SUSPEND or RESUME on ALTER TABLE MODIFY DATA METRIC FUNCTION.
var TableDataMetricScheduleStatusOperation = ssot.NewEnum(
	"TableDataMetricScheduleStatusOperation", "TableDataMetricScheduleStatusOperations",
	"RESUME", "SUSPEND",
)

// TableDataMetricSensitivity is SENSITIVITY on ADD and MODIFY SET of a data metric function.
var TableDataMetricSensitivity = ssot.NewEnum(
	"TableDataMetricSensitivity", "TableDataMetricSensitivities",
	"LOW", "MEDIUM", "HIGH",
)

var tableDataMetricFunctionOn = ssot.ListParameter("Column", ssot.Sql("ON"), ssot.NoEquals(), ssot.Required(), ssot.FieldName("On"))

var tableDataMetricExpectationFields = []ssot.Attribute{
	{Name: "name", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.TextKeyword(ssot.DoubleQuotes(), ssot.Required())},
	{Name: "expression", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.PredefinedType("string", ssot.MustParentheses(), ssot.NoQuotes(), ssot.NoEquals(), ssot.Required(), ssot.FieldName("Expression"))},
}

func tableDataMetricFunctionIdentifier() ssot.Attribute {
	return ssot.Attribute{
		Name:           "data_metric_function",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "SchemaObjectIdentifier",
		Required:       true,
		SqlOnly:        true,
		// IdentifierKeyword always writes SQLText. Sql(" ") keeps the identifier with no DATA_METRIC_FUNCTION prefix;
		// that keyword already sits on the list field.
		CreateSql: ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Required(), ssot.FieldName("DataMetricFunction"), ssot.Sql(" ")),
	}
}

func tableDataMetricFunctionOnAttr() ssot.Attribute {
	return ssot.Attribute{
		Name:      "on",
		Type:      ssot.AttrList,
		Required:  true,
		SqlOnly:   true,
		Elem:      &ssot.Attribute{Type: ssot.AttrString},
		CreateSql: tableDataMetricFunctionOn,
	}
}

func tableAddDataMetricFunctionItemChildren() []ssot.Attribute {
	return []ssot.Attribute{
		tableDataMetricFunctionIdentifier(),
		tableDataMetricFunctionOnAttr(),
		{
			Name:      "within_group",
			Type:      ssot.AttrList,
			SqlOnly:   true,
			Elem:      &ssot.Attribute{Type: ssot.AttrString},
			CreateSql: ssot.ListParameter("Column", ssot.Sql("WITHIN GROUP"), ssot.NoEquals()),
		},
		{Name: "group_limit", Type: ssot.AttrInt, SqlOnly: true, CreateSql: ssot.IntParameter(ssot.Sql("GROUP LIMIT"), ssot.NoEquals())},
		{Name: "filter", Type: ssot.AttrString, SqlOnly: true, CreateSql: ssot.PredefinedType("*string", ssot.MustParentheses(), ssot.NoQuotes(), ssot.NoEquals(), ssot.Sql("FILTER"), ssot.FieldName("Filter"))},
		{
			Name:      "expectations",
			Type:      ssot.AttrList,
			SqlOnly:   true,
			CreateSql: ssot.ValueKeyword(ssot.Sql("EXPECTATION")),
			Elem: &ssot.Attribute{
				Type:      ssot.AttrObject,
				SqlOnly:   true,
				CreateSql: ssot.ObjectParameter(),
				Children:  tableDataMetricExpectationFields,
			},
		},
		{
			Name:           "execute_as_role",
			Type:           ssot.AttrIdentifier,
			IdentifierType: "AccountObjectIdentifier",
			SqlOnly:        true,
			CreateSql:      ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Sql("EXECUTE AS ROLE")),
		},
		{Name: "anomaly_detection", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolParameter()},
		{Name: "sensitivity", Type: ssot.AttrEnum, Enum: TableDataMetricSensitivity, SqlOnly: true, CreateSql: ssot.EnumParameter(ssot.SingleQuotes())},
	}
}

// TableAddDataMetricFunction is ALTER TABLE ADD DATA METRIC FUNCTION.
var TableAddDataMetricFunction = ssot.Attribute{
	Name:      "add_data_metric_function",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(ssot.StructName("TableAddDataMetricFunction")),
	Children: []ssot.Attribute{
		{Name: "add", SqlOnly: true, CreateSql: ssot.StaticSQL("add", "ADD")},
		{
			Name:      "added_functions",
			Type:      ssot.AttrList,
			Required:  true,
			SqlOnly:   true,
			CreateSql: ssot.ValueKeyword(ssot.Sql("DATA METRIC FUNCTION"), ssot.Required()),
			Elem: &ssot.Attribute{
				Type:      ssot.AttrObject,
				SqlOnly:   true,
				CreateSql: ssot.ObjectParameter(ssot.StructName("TableDataMetricFunction")),
				Children:  tableAddDataMetricFunctionItemChildren(),
			},
		},
	},
}

// TableDropDataMetricFunction is ALTER TABLE DROP DATA METRIC FUNCTION.
var TableDropDataMetricFunction = ssot.Attribute{
	Name:      "drop_data_metric_function",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(ssot.StructName("TableDropDataMetricFunction")),
	Children: []ssot.Attribute{
		{Name: "drop", SqlOnly: true, CreateSql: ssot.StaticSQL("drop", "DROP")},
		{
			Name:      "dropped_functions",
			Type:      ssot.AttrList,
			Required:  true,
			SqlOnly:   true,
			CreateSql: ssot.ValueKeyword(ssot.Sql("DATA METRIC FUNCTION"), ssot.Required()),
			Elem: &ssot.Attribute{
				Type:      ssot.AttrObject,
				SqlOnly:   true,
				CreateSql: ssot.ObjectParameter(ssot.StructName("TableDropDataMetricFunctionItem")),
				Children: []ssot.Attribute{
					tableDataMetricFunctionIdentifier(),
					tableDataMetricFunctionOnAttr(),
				},
			},
		},
	},
}

// TableModifyDataMetricFunction is ALTER TABLE MODIFY DATA METRIC FUNCTION. Exactly one of suspend/resume, expectation
// add/modify/drop, FILTER, and SET is present on each item. Identifier and ON sit beside that choice.
var TableModifyDataMetricFunction = ssot.Attribute{
	Name:      "modify_data_metric_function",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(ssot.StructName("TableModifyDataMetricFunction")),
	Children: []ssot.Attribute{
		{Name: "modify", SqlOnly: true, CreateSql: ssot.StaticSQL("modify", "MODIFY")},
		{
			Name:      "modified_functions",
			Type:      ssot.AttrList,
			Required:  true,
			SqlOnly:   true,
			CreateSql: ssot.ValueKeyword(ssot.Sql("DATA METRIC FUNCTION"), ssot.Required()),
			Elem: &ssot.Attribute{
				Type:    ssot.AttrObject,
				SqlOnly: true,
				Validations: []ssot.Validation{ssot.ExactlyOneOf(
					ssot.Attribute{Name: "status"},
					ssot.Attribute{Name: "add_expectation"},
					ssot.Attribute{Name: "modify_expectation"},
					ssot.Attribute{Name: "drop_expectation"},
					ssot.Attribute{Name: "filter"},
					ssot.Attribute{Name: "set"},
				)},
				CreateSql: ssot.ObjectParameter(ssot.StructName("TableModifyDataMetricFunctionItem")),
				Children: []ssot.Attribute{
					tableDataMetricFunctionIdentifier(),
					tableDataMetricFunctionOnAttr(),
					{Name: "status", Type: ssot.AttrEnum, Enum: TableDataMetricScheduleStatusOperation, SqlOnly: true, CreateSql: ssot.EnumKeyword()},
					{
						Name:      "add_expectation",
						Type:      ssot.AttrList,
						SqlOnly:   true,
						CreateSql: ssot.ValueKeyword(ssot.Sql("ADD EXPECTATION")),
						Elem: &ssot.Attribute{
							Type:      ssot.AttrObject,
							SqlOnly:   true,
							CreateSql: ssot.ObjectParameter(),
							Children:  tableDataMetricExpectationFields,
						},
					},
					{
						Name:      "modify_expectation",
						Type:      ssot.AttrList,
						SqlOnly:   true,
						CreateSql: ssot.ValueKeyword(ssot.Sql("MODIFY EXPECTATION")),
						Elem: &ssot.Attribute{
							Type:      ssot.AttrObject,
							SqlOnly:   true,
							CreateSql: ssot.ObjectParameter(),
							Children:  tableDataMetricExpectationFields,
						},
					},
					{
						Name:      "drop_expectation",
						Type:      ssot.AttrList,
						SqlOnly:   true,
						CreateSql: ssot.ValueKeyword(ssot.Sql("DROP EXPECTATION")),
						Elem: &ssot.Attribute{
							Type:      ssot.AttrObject,
							SqlOnly:   true,
							CreateSql: ssot.ObjectParameter(),
							Children: []ssot.Attribute{
								{Name: "name", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.TextKeyword(ssot.DoubleQuotes(), ssot.Required())},
							},
						},
					},
					{Name: "filter", Type: ssot.AttrString, SqlOnly: true, CreateSql: ssot.PredefinedType("*string", ssot.MustParentheses(), ssot.NoQuotes(), ssot.NoEquals(), ssot.Sql("FILTER"), ssot.FieldName("Filter"))},
					{
						Name:      "set",
						Type:      ssot.AttrObject,
						SqlOnly:   true,
						CreateSql: ssot.ObjectParameter(ssot.Sql("SET"), ssot.StructName("TableModifyDataMetricFunctionSet")),
						Children: []ssot.Attribute{
							{Name: "anomaly_detection", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolParameter()},
							{Name: "sensitivity", Type: ssot.AttrEnum, Enum: TableDataMetricSensitivity, SqlOnly: true, CreateSql: ssot.EnumParameter(ssot.SingleQuotes())},
							{Name: "data_quality_notification", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolParameter()},
						},
					},
				},
			},
		},
	},
}

// TableSetDataMetricSchedule is ALTER TABLE SET DATA_METRIC_SCHEDULE. It is a dataMetricFunctionAction, not a SET property.
var TableSetDataMetricSchedule = ssot.Attribute{
	Name:      "set_data_metric_schedule",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(ssot.StructName("TableSetDataMetricSchedule")),
	Children: []ssot.Attribute{
		{Name: "set", SqlOnly: true, CreateSql: ssot.StaticSQL("set", "SET")},
		{
			Name:      "data_metric_schedule",
			Type:      ssot.AttrString,
			Required:  true,
			SqlOnly:   true,
			CreateSql: ssot.StringParameter(ssot.Required()),
		},
	},
}

// TableUnsetDataMetricSchedule is ALTER TABLE UNSET DATA_METRIC_SCHEDULE.
var TableUnsetDataMetricSchedule = ssot.Attribute{
	Name:      "unset_data_metric_schedule",
	Type:      ssot.AttrBool,
	SqlOnly:   true,
	CreateSql: ssot.BoolKeyword(ssot.Sql("UNSET DATA_METRIC_SCHEDULE")),
}
