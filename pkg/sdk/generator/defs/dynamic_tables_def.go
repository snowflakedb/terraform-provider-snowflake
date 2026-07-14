package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var DynamicTableRefreshModeEnumDef = g.NewEnum(
	"DynamicTableRefreshMode", "DynamicTableRefreshModes",
	"AUTO", "INCREMENTAL", "FULL",
)

var DynamicTableInitializeEnumDef = g.NewEnum(
	"DynamicTableInitialize", "DynamicTableInitializes",
	"ON_CREATE", "ON_SCHEDULE",
)

var DynamicTableSchedulingStateEnumDef = g.NewEnum(
	"DynamicTableSchedulingState", "DynamicTableSchedulingStates",
	"ACTIVE", "SUSPENDED",
)

func targetLagDef() *g.QueryStruct {
	return g.NewQueryStruct("TargetLag").
		OptionalText("MaximumDuration", g.KeywordOptions().SingleQuotes()).
		OptionalSQL("DOWNSTREAM").
		WithValidation(g.ConflictingFields, "MaximumDuration", "Downstream")
}

func dynamicTableSetDef() *g.QueryStruct {
	return g.NewQueryStruct("DynamicTableSet").
		OptionalQueryStructField("TargetLag", targetLagDef(), g.ParameterOptions().SQL("TARGET_LAG").NoQuotes()).
		OptionalIdentifier("Warehouse", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("WAREHOUSE").Equals()).
		WithAdditionalValidations()
}

func dynamicTableAddStorageLifecyclePolicyDef() *g.QueryStruct {
	return g.NewQueryStruct("DynamicTableAddStorageLifecyclePolicy").
		SQL("ADD").
		Identifier("StorageLifecyclePolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("STORAGE LIFECYCLE POLICY").Required()).
		ListAssignment("ON", "Column", g.ParameterOptions().NoEquals().Parentheses().Required())
}

var dynamicTablesDef = g.NewInterface(
	"DynamicTables",
	"DynamicTable",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-dynamic-table",
		g.NewQueryStruct("CreateDynamicTable").
			Create().
			OrReplace().
			SQL("DYNAMIC TABLE").
			Name().
			QueryStructField("TargetLag", targetLagDef(), g.ParameterOptions().SQL("TARGET_LAG").NoQuotes().Required()).
			OptionalEnumAssignment("INITIALIZE", DynamicTableInitializeEnumDef, g.ParameterOptions().NoQuotes()).
			OptionalEnumAssignment("REFRESH_MODE", DynamicTableRefreshModeEnumDef, g.ParameterOptions().NoQuotes()).
			Identifier("Warehouse", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("WAREHOUSE").Equals().Required()).
			OptionalComment().
			SQL("AS").
			Text("Query", g.KeywordOptions().NoQuotes().Required()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidIdentifier, "Warehouse"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-dynamic-table",
		g.NewQueryStruct("AlterDynamicTable").
			Alter().
			SQL("DYNAMIC TABLE").
			Name().
			OptionalSQL("SUSPEND").
			OptionalSQL("RESUME").
			OptionalSQL("REFRESH").
			OptionalQueryStructField("Set", dynamicTableSetDef(), g.KeywordOptions().SQL("SET")).
			OptionalTextAssignment("SET COMMENT", g.ParameterOptions().SingleQuotes()).
			OptionalQueryStructField("AddStorageLifecyclePolicy", dynamicTableAddStorageLifecyclePolicyDef(), g.KeywordOptions()).
			OptionalSQL("DROP STORAGE LIFECYCLE POLICY").
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "Suspend", "Resume", "Refresh", "Set", "SetComment", "AddStorageLifecyclePolicy", "DropStorageLifecyclePolicy").
			WithAdditionalValidations(),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-dynamic-table",
		g.NewQueryStruct("DropDynamicTable").
			Drop().
			SQL("DYNAMIC TABLE").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-dynamic-tables",
		g.StructPair("dynamicTableRow", "DynamicTable").
			Time("created_on").
			Text("name").
			Text("reserved").
			Text("database_name").
			Text("schema_name").
			Text("cluster_by").
			Number("rows").
			Number("bytes").
			Text("owner").
			Text("target_lag").
			Enum("refresh_mode", DynamicTableRefreshModeEnumDef).
			OptionalText("refresh_mode_reason", g.WithRequiredInPlain()).
			Text("warehouse").
			Text("comment").
			Text("text", g.WithValueAdjuster("tracking.TrimMetadata")).
			Field("automatic_clustering", "string", "bool", g.WithBoolTrueValue("ON")).
			Enum("scheduling_state", DynamicTableSchedulingStateEnumDef).
			OptionalTime("last_suspended_on").
			Field("is_clone", "bool", "bool").
			Field("is_replica", "bool", "bool").
			OptionalTime("data_timestamp").
			OptionalText("owner_role_type", g.WithRequiredInPlain()),
		g.NewQueryStruct("ShowDynamicTables").
			Show().
			SQL("DYNAMIC TABLES").
			OptionalLike().
			OptionalIn().
			OptionalStartsWith().
			OptionalLimitFrom().
			WithAdditionalValidations(),
		g.ShowByIDInFiltering,
		g.ShowByIDLikeFiltering,
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-dynamic-table",
		g.StructPair("dynamicTableDetailsRow", "DynamicTableDetails").
			Text("name").
			Field("type", "string", "DataType", g.WithManualConvert()).
			Text("kind").
			Field("null?", "string", "bool", g.WithBoolTrueValue("Y"), g.WithDbFieldName("IsNull"), g.WithPlainFieldName("IsNull")).
			OptionalText("default", g.WithRequiredInPlain()).
			Text("primary key", g.WithPlainFieldName("PrimaryKey")).
			Text("unique key", g.WithPlainFieldName("UniqueKey")).
			OptionalText("check", g.WithRequiredInPlain()).
			OptionalText("expression", g.WithRequiredInPlain()).
			OptionalText("comment", g.WithRequiredInPlain()).
			OptionalText("policy name", g.WithPlainFieldName("PolicyName"), g.WithRequiredInPlain()),
		g.NewQueryStruct("DescribeDynamicTable").
			Describe().
			SQL("DYNAMIC TABLE").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	WithEnums(DynamicTableRefreshModeEnumDef, DynamicTableInitializeEnumDef, DynamicTableSchedulingStateEnumDef)
