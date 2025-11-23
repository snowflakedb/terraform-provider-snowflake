package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var ComputePoolsDef = g.NewInterface(
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
		WithValidation(g.ValidIdentifier, "name"),
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
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-compute-pools",
	g.DbStruct("computePoolsRow").
		Text("name").
		Text("state").
		Number("min_nodes").
		Number("max_nodes").
		Text("instance_family").
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
		OptionalText("application"),
	g.PlainStruct("ComputePool").
		Text("Name").
		Field("State", "ComputePoolState").
		Number("MinNodes").
		Number("MaxNodes").
		Field("InstanceFamily", "ComputePoolInstanceFamily").
		Number("NumServices").
		Number("NumJobs").
		Number("AutoSuspendSecs").
		Bool("AutoResume").
		Number("ActiveNodes").
		Number("IdleNodes").
		Number("TargetNodes").
		Time("CreatedOn").
		Time("ResumedOn").
		Time("UpdatedOn").
		Text("Owner").
		OptionalText("Comment").
		Bool("IsExclusive").
		Field("Application", "*AccountObjectIdentifier"),
	g.NewQueryStruct("ShowComputePools").
		Show().
		SQL("COMPUTE POOLS").
		OptionalLike().
		OptionalStartsWith().
		OptionalLimitFrom(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
).DescribeOperation(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-compute-pool",
	g.DbStruct("computePoolDescRow").
		Text("name").
		Text("state").
		Number("min_nodes").
		Number("max_nodes").
		Text("instance_family").
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
		OptionalText("application").
		Text("error_code").
		Text("status_message"),
	g.PlainStruct("ComputePoolDetails").
		Text("Name").
		Field("State", "ComputePoolState").
		Number("MinNodes").
		Number("MaxNodes").
		Field("InstanceFamily", "ComputePoolInstanceFamily").
		Number("NumServices").
		Number("NumJobs").
		Number("AutoSuspendSecs").
		Bool("AutoResume").
		Number("ActiveNodes").
		Number("IdleNodes").
		Number("TargetNodes").
		Time("CreatedOn").
		Time("ResumedOn").
		Time("UpdatedOn").
		Text("Owner").
		OptionalText("Comment").
		Bool("IsExclusive").
		Field("Application", "*AccountObjectIdentifier").
		Text("ErrorCode").
		Text("StatusMessage"),
	g.NewQueryStruct("DescComputePool").
		Describe().
		SQL("COMPUTE POOL").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
