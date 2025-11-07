package sdk

import (
	"fmt"
	"strings"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

var AllViewDataMetricScheduleMinutes = []int{5, 15, 30, 60, 720, 1440}

type ViewDataMetricScheduleStatusOperationOption string

const (
	ViewDataMetricScheduleStatusOperationResume  ViewDataMetricScheduleStatusOperationOption = "RESUME"
	ViewDataMetricScheduleStatusOperationSuspend ViewDataMetricScheduleStatusOperationOption = "SUSPEND"
)

var AllViewDataMetricScheduleStatusOperationOptions = []ViewDataMetricScheduleStatusOperationOption{
	ViewDataMetricScheduleStatusOperationResume,
	ViewDataMetricScheduleStatusOperationSuspend,
}

func ToViewDataMetricScheduleStatusOperationOption(s string) (ViewDataMetricScheduleStatusOperationOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(ViewDataMetricScheduleStatusOperationResume):
		return ViewDataMetricScheduleStatusOperationResume, nil
	case string(ViewDataMetricScheduleStatusOperationSuspend):
		return ViewDataMetricScheduleStatusOperationSuspend, nil
	default:
		return "", fmt.Errorf("invalid ViewDataMetricScheduleStatusOperationOption: %s", s)
	}
}

var viewDbRow = g.DbStruct("viewDBRow").
	Text("created_on").
	Text("name").
	OptionalText("kind").
	OptionalText("reserved").
	Text("database_name").
	Text("schema_name").
	OptionalText("owner").
	OptionalText("comment").
	OptionalText("text").
	OptionalBool("is_secure").
	OptionalBool("is_materialized").
	OptionalText("owner_role_type").
	OptionalText("change_tracking")

var view = g.PlainStruct("View").
	Text("CreatedOn").
	Text("Name").
	Text("Kind").
	Text("Reserved").
	Text("DatabaseName").
	Text("SchemaName").
	Text("Owner").
	Text("Comment").
	Text("Text").
	Bool("IsSecure").
	Bool("IsMaterialized").
	Text("OwnerRoleType").
	Text("ChangeTracking")

// TODO [SNOW-965322]: extract common type for describe
var viewDetailsDbRow = g.DbStruct("viewDetailsRow").
	Text("name").
	Field("type", "DataType").
	Text("kind").
	Text("null").
	OptionalText("default").
	Text("primary key").
	Text("unique key").
	OptionalText("check").
	OptionalText("expression").
	OptionalText("comment").
	OptionalText("policy name").
	OptionalText("privacy domain")

var viewDetails = g.PlainStruct("ViewDetails").
	Text("Name").
	Field("Type", "DataType").
	Text("Kind").
	Bool("IsNullable").
	OptionalText("Default").
	Bool("IsPrimary").
	Bool("IsUnique").
	OptionalBool("Check").
	OptionalText("Expression").
	OptionalText("Comment").
	OptionalText("PolicyName").
	OptionalText("PrivacyDomain")

var columnDef = g.NewQueryStruct("Column").
	Text("Value", g.KeywordOptions().Required().DoubleQuotes())

var viewMinute = g.NewQueryStruct("ViewMinute").
	Number("Minutes", g.KeywordOptions().Required()).
	SQL("MINUTE")

var viewUsingCron = g.NewQueryStruct("ViewUsingCron").
	SQL("USING CRON").
	Text("Cron", g.KeywordOptions().Required())

var dataMetricFunctionDef = g.NewQueryStruct("ViewDataMetricFunction").
	Identifier("DataMetricFunction", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
	ListAssignment("ON", "Column", g.ParameterOptions().Required().NoEquals().Parentheses()).
	WithValidation(g.ValidIdentifier, "DataMetricFunction")

var modifyDataMetricFunctionDef = g.NewQueryStruct("ViewModifyDataMetricFunction").
	Identifier("DataMetricFunction", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
	ListAssignment("ON", "Column", g.ParameterOptions().Required().NoEquals().Parentheses()).
	Assignment(
		"",
		g.KindOfT[ViewDataMetricScheduleStatusOperationOption](),
		g.ParameterOptions().NoEquals().NoQuotes(),
	).
	WithValidation(g.ValidIdentifier, "DataMetricFunction")

var viewColumn = g.NewQueryStruct("ViewColumn").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	OptionalQueryStructField("ProjectionPolicy", viewColumnProjectionPolicy, g.KeywordOptions()).
	OptionalQueryStructField("MaskingPolicy", viewColumnMaskingPolicy, g.KeywordOptions()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes().NoEquals()).
	OptionalTags()

var viewColumnMaskingPolicy = g.NewQueryStruct("ViewColumnMaskingPolicy").
	Identifier("MaskingPolicy", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("MASKING POLICY").Required()).
	ListAssignment("USING", "Column", g.ParameterOptions().NoEquals().Parentheses())

var viewColumnProjectionPolicy = g.NewQueryStruct("ViewColumnProjectionPolicy").
	Identifier("ProjectionPolicy", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("PROJECTION POLICY").Required())

var viewRowAccessPolicy = g.NewQueryStruct("ViewRowAccessPolicy").
	Identifier("RowAccessPolicy", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("ROW ACCESS POLICY").Required()).
	ListAssignment("ON", "Column", g.ParameterOptions().Required().NoEquals().Parentheses()).
	WithValidation(g.ValidIdentifier, "RowAccessPolicy")

var viewAggregationPolicy = g.NewQueryStruct("ViewAggregationPolicy").
	Identifier("AggregationPolicy", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("AGGREGATION POLICY").Required()).
	ListAssignment("ENTITY KEY", "Column", g.ParameterOptions().NoEquals().Parentheses()).
	WithValidation(g.ValidIdentifier, "AggregationPolicy")

var viewAddDataMetricFunction = g.NewQueryStruct("ViewAddDataMetricFunction").
	SQL("ADD").
	ListAssignment("DATA METRIC FUNCTION", "ViewDataMetricFunction", g.ParameterOptions().NoEquals().Required())

var viewDropDataMetricFunction = g.NewQueryStruct("ViewDropDataMetricFunction").
	SQL("DROP").
	ListAssignment("DATA METRIC FUNCTION", "ViewDataMetricFunction", g.ParameterOptions().NoEquals().Required())

var viewModifyDataMetricFunction = g.NewQueryStruct("ViewModifyDataMetricFunctions").
	SQL("MODIFY").
	ListAssignment("DATA METRIC FUNCTION", "ViewModifyDataMetricFunction", g.ParameterOptions().NoEquals().Required())

var viewSetDataMetricSchedule = g.NewQueryStruct("ViewSetDataMetricSchedule").
	SQL("SET").
	TextAssignment("DATA_METRIC_SCHEDULE", g.ParameterOptions().SingleQuotes())

var viewUnsetDataMetricSchedule = g.NewQueryStruct("ViewUnsetDataMetricSchedule").
	SQL("UNSET DATA_METRIC_SCHEDULE")

var viewAddRowAccessPolicy = g.NewQueryStruct("ViewAddRowAccessPolicy").
	SQL("ADD").
	Identifier("RowAccessPolicy", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("ROW ACCESS POLICY").Required()).
	ListAssignment("ON", "Column", g.ParameterOptions().Required().NoEquals().Parentheses()).
	WithValidation(g.ValidIdentifier, "RowAccessPolicy")

var viewDropRowAccessPolicy = g.NewQueryStruct("ViewDropRowAccessPolicy").
	SQL("DROP").
	Identifier("RowAccessPolicy", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("ROW ACCESS POLICY").Required()).
	WithValidation(g.ValidIdentifier, "RowAccessPolicy")

var viewDropAndAddRowAccessPolicy = g.NewQueryStruct("ViewDropAndAddRowAccessPolicy").
	QueryStructField("Drop", viewDropRowAccessPolicy, g.KeywordOptions().Required()).
	QueryStructField("Add", viewAddRowAccessPolicy, g.KeywordOptions().Required())

var viewSetAggregationPolicy = g.NewQueryStruct("ViewSetAggregationPolicy").
	SQL("SET").
	Identifier("AggregationPolicy", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("AGGREGATION POLICY").Required()).
	ListAssignment("ENTITY KEY", "Column", g.ParameterOptions().NoEquals().Parentheses()).
	OptionalSQL("FORCE").
	WithValidation(g.ValidIdentifier, "AggregationPolicy")

var viewUnsetAggregationPolicy = g.NewQueryStruct("ViewUnsetAggregationPolicy").
	SQL("UNSET AGGREGATION POLICY")

var viewSetColumnMaskingPolicy = g.NewQueryStruct("ViewSetColumnMaskingPolicy").
	// In the docs there is a MODIFY alternative, but for simplicity only one is supported here.
	SQL("ALTER").
	SQL("COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	SQL("SET").
	Identifier("MaskingPolicy", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("MASKING POLICY").Required()).
	ListAssignment("USING", "Column", g.ParameterOptions().NoEquals().Parentheses()).
	OptionalSQL("FORCE")

var viewUnsetColumnMaskingPolicy = g.NewQueryStruct("ViewUnsetColumnMaskingPolicy").
	// In the docs there is a MODIFY alternative, but for simplicity only one is supported here.
	SQL("ALTER").
	SQL("COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	SQL("UNSET").
	SQL("MASKING POLICY")

var viewSetProjectionPolicy = g.NewQueryStruct("ViewSetProjectionPolicy").
	// In the docs there is a MODIFY alternative, but for simplicity only one is supported here.
	SQL("ALTER").
	SQL("COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	SQL("SET").
	Identifier("ProjectionPolicy", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("PROJECTION POLICY").Required()).
	OptionalSQL("FORCE")

var viewUnsetProjectionPolicy = g.NewQueryStruct("ViewUnsetProjectionPolicy").
	// In the docs there is a MODIFY alternative, but for simplicity only one is supported here.
	SQL("ALTER").
	SQL("COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	SQL("UNSET").
	SQL("PROJECTION POLICY")

var viewSetColumnTags = g.NewQueryStruct("ViewSetColumnTags").
	// In the docs there is a MODIFY alternative, but for simplicity only one is supported here.
	SQL("ALTER").
	SQL("COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	SetTags()

var viewUnsetColumnTags = g.NewQueryStruct("ViewUnsetColumnTags").
	// In the docs there is a MODIFY alternative, but for simplicity only one is supported here.
	SQL("ALTER").
	SQL("COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	UnsetTags()

var ViewsDef = g.NewInterface(
	"Views",
	"View",
	g.KindOfT[SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-view",
		g.NewQueryStruct("CreateView").
			Create().
			OrReplace().
			OptionalSQL("SECURE").
			// There are multiple variants in the docs: { [ { LOCAL | GLOBAL } ] TEMP | TEMPORARY | VOLATILE }
			// but from description they are all the same. For the sake of simplicity only one option is used here.
			OptionalSQL("TEMPORARY").
			OptionalSQL("RECURSIVE").
			SQL("VIEW").
			IfNotExists().
			Name().
			ListQueryStructField("Columns", viewColumn, g.ListOptions().Parentheses()).
			OptionalSQL("COPY GRANTS").
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			// In the current docs ROW ACCESS POLICY and TAG are specified twice.
			// It is a mistake probably so here they are present only once.
			OptionalQueryStructField("RowAccessPolicy", viewRowAccessPolicy, g.KeywordOptions()).
			OptionalQueryStructField("AggregationPolicy", viewAggregationPolicy, g.KeywordOptions()).
			OptionalTags().
			SQL("AS").
			Text("sql", g.KeywordOptions().NoQuotes().Required()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	).
	CustomOperation(
		"Alter",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-view",
		g.NewQueryStruct("AlterView").
			Alter().
			SQL("VIEW").
			IfExists().
			Name().
			OptionalIdentifier("RenameTo", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
			OptionalTextAssignment("SET COMMENT", g.ParameterOptions().SingleQuotes()).
			OptionalSQL("UNSET COMMENT").
			OptionalSQL("SET SECURE").
			OptionalBooleanAssignment("SET CHANGE_TRACKING", nil).
			OptionalSQL("UNSET SECURE").
			OptionalSetTags().
			OptionalUnsetTags().
			OptionalQueryStructField("AddDataMetricFunction", viewAddDataMetricFunction, g.KeywordOptions()).
			OptionalQueryStructField("DropDataMetricFunction", viewDropDataMetricFunction, g.KeywordOptions()).
			OptionalQueryStructField("ModifyDataMetricFunction", viewModifyDataMetricFunction, g.KeywordOptions()).
			OptionalQueryStructField("SetDataMetricSchedule", viewSetDataMetricSchedule, g.KeywordOptions()).
			OptionalQueryStructField("UnsetDataMetricSchedule", viewUnsetDataMetricSchedule, g.KeywordOptions()).
			OptionalQueryStructField("AddRowAccessPolicy", viewAddRowAccessPolicy, g.KeywordOptions()).
			OptionalQueryStructField("DropRowAccessPolicy", viewDropRowAccessPolicy, g.KeywordOptions()).
			OptionalQueryStructField("DropAndAddRowAccessPolicy", viewDropAndAddRowAccessPolicy, g.ListOptions().NoParentheses()).
			OptionalSQL("DROP ALL ROW ACCESS POLICIES").
			OptionalQueryStructField("SetAggregationPolicy", viewSetAggregationPolicy, g.KeywordOptions()).
			OptionalQueryStructField("UnsetAggregationPolicy", viewUnsetAggregationPolicy, g.KeywordOptions()).
			OptionalQueryStructField("SetMaskingPolicyOnColumn", viewSetColumnMaskingPolicy, g.KeywordOptions()).
			OptionalQueryStructField("UnsetMaskingPolicyOnColumn", viewUnsetColumnMaskingPolicy, g.KeywordOptions()).
			OptionalQueryStructField("SetProjectionPolicyOnColumn", viewSetProjectionPolicy, g.KeywordOptions()).
			OptionalQueryStructField("UnsetProjectionPolicyOnColumn", viewUnsetProjectionPolicy, g.KeywordOptions()).
			OptionalQueryStructField("SetTagsOnColumn", viewSetColumnTags, g.KeywordOptions()).
			OptionalQueryStructField("UnsetTagsOnColumn", viewUnsetColumnTags, g.KeywordOptions()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "RenameTo", "SetComment", "UnsetComment", "SetSecure", "SetChangeTracking",
				"UnsetSecure", "SetTags", "UnsetTags", "AddDataMetricFunction", "DropDataMetricFunction", "ModifyDataMetricFunction", "SetDataMetricSchedule", "UnsetDataMetricSchedule",
				"AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy",
				"DropAllRowAccessPolicies", "SetAggregationPolicy", "UnsetAggregationPolicy", "SetMaskingPolicyOnColumn",
				"UnsetMaskingPolicyOnColumn", "SetProjectionPolicyOnColumn", "UnsetProjectionPolicyOnColumn", "SetTagsOnColumn",
				"UnsetTagsOnColumn").
			WithValidation(g.ConflictingFields, "IfExists", "SetSecure").
			WithValidation(g.ConflictingFields, "IfExists", "UnsetSecure"),
		columnDef,
		dataMetricFunctionDef,
		modifyDataMetricFunctionDef,
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-view",
		g.NewQueryStruct("DropView").
			Drop().
			SQL("VIEW").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-views",
		viewDbRow,
		view,
		g.NewQueryStruct("ShowViews").
			Show().
			Terse().
			SQL("VIEWS").
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
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-view",
		viewDetailsDbRow,
		viewDetails,
		g.NewQueryStruct("DescribeView").
			Describe().
			SQL("VIEW").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
