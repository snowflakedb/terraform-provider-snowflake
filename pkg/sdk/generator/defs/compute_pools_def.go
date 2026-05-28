package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

// See https://docs.snowflake.com/en/developer-guide/snowpark-container-services/working-with-compute-pool#compute-pool-lifecycle.
var ComputePoolStateEnumDef = g.NewEnum(
	"ComputePoolState", "ComputePoolStates",
	"IDLE", "ACTIVE", "SUSPENDED", "STARTING", "STOPPING", "RESIZING",
)

var computePoolsDef = g.NewInterface(
	"ComputePools",
	"ComputePool",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-compute-pool",
	g.NewQueryStruct("CreateComputePool").
		Create().
		SQL("COMPUTE POOL").
		// Note: Currently, OR REPLACE is not supported for compute pools.
		IfNotExists().
		Name().
		OptionalIdentifier("ForApplication", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("FOR APPLICATION")).
		NumberAssignment("MIN_NODES", g.ParameterOptions().Required()).
		NumberAssignment("MAX_NODES", g.ParameterOptions().Required()).
		Assignment(
			"INSTANCE_FAMILY",
			g.KindOfT[sdkcommons.ComputePoolInstanceFamily](),
			g.ParameterOptions().NoQuotes().Required(),
		).
		OptionalBooleanAssignment("AUTO_RESUME", g.ParameterOptions()).
		OptionalBooleanAssignment("INITIALLY_SUSPENDED", g.ParameterOptions()).
		OptionalNumberAssignment("AUTO_SUSPEND_SECS", g.ParameterOptions()).
		OptionalTags().
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.ValidIdentifier, "name").
		WithAdditionalValidations(),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-compute-pool",
	g.NewQueryStruct("AlterComputePool").
		Alter().
		SQL("COMPUTE POOL").
		IfExists().
		Name().
		OptionalSQL("RESUME").
		OptionalSQL("SUSPEND").
		OptionalSQL("STOP ALL").
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("ComputePoolSet").
				OptionalNumberAssignment("MIN_NODES", g.ParameterOptions()).
				OptionalNumberAssignment("MAX_NODES", g.ParameterOptions()).
				OptionalBooleanAssignment("AUTO_RESUME", g.ParameterOptions()).
				OptionalNumberAssignment("AUTO_SUSPEND_SECS", g.ParameterOptions()).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				WithAdditionalValidations().
				WithValidation(g.AtLeastOneValueSet, "MinNodes", "MaxNodes", "AutoResume", "AutoSuspendSecs", "Comment"),
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("ComputePoolUnset").
				OptionalSQL("AUTO_RESUME").
				OptionalSQL("AUTO_SUSPEND_SECS").
				OptionalSQL("COMMENT").
				WithValidation(g.AtLeastOneValueSet, "AutoResume", "AutoSuspendSecs", "Comment"),
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Resume", "Suspend", "StopAll", "Set", "Unset", "SetTags", "UnsetTags"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-compute-pool",
	g.NewQueryStruct("DropComputePool").
		Drop().
		SQL("COMPUTE POOL").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-compute-pools",
	g.StructPair("computePoolsRow", "ComputePool").
		Text("name").
		Enum("state", ComputePoolStateEnumDef).
		Number("min_nodes").
		Number("max_nodes").
		PlainField("instance_family", "ComputePoolInstanceFamily", g.WithCustomParser("ToComputePoolInstanceFamily")).
		Number("num_services").
		Number("num_jobs").
		Number("auto_suspend_secs").
		Bool("auto_resume").
		Number("active_nodes").
		Number("idle_nodes").
		Number("target_nodes").
		Time("created_on").
		Time("resumed_on").
		Time("updated_on").
		Text("owner").
		OptionalText("comment").
		Bool("is_exclusive").
		Field("application", "sql.NullString", "*AccountObjectIdentifier", g.WithPlainFieldName("Application")).
		WithConvertGeneration(),
	g.NewQueryStruct("ShowComputePools").
		Show().
		SQL("COMPUTE POOLS").
		OptionalLike().
		OptionalStartsWith().
		OptionalLimitFrom(),
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-compute-pool",
	g.StructPair("computePoolDescRow", "ComputePoolDetails").
		Text("name").
		Enum("state", ComputePoolStateEnumDef).
		Number("min_nodes").
		Number("max_nodes").
		PlainField("instance_family", "ComputePoolInstanceFamily", g.WithCustomParser("ToComputePoolInstanceFamily")).
		Number("num_services").
		Number("num_jobs").
		Number("auto_suspend_secs").
		Bool("auto_resume").
		Number("active_nodes").
		Number("idle_nodes").
		Number("target_nodes").
		Time("created_on").
		Time("resumed_on").
		Time("updated_on").
		Text("owner").
		OptionalText("comment").
		Bool("is_exclusive").
		Field("application", "sql.NullString", "*AccountObjectIdentifier", g.WithPlainFieldName("Application")).
		Text("error_code").
		Text("status_message").
		WithConvertGeneration(),
	g.NewQueryStruct("DescComputePool").
		Describe().
		SQL("COMPUTE POOL").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).
	WithEnums(
		ComputePoolStateEnumDef,
	)
