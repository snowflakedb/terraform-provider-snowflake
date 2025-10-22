package sdk

import (
	"fmt"
	"slices"
	"strings"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

type ComputePoolInstanceFamily string

const (
	ComputePoolInstanceFamilyCpuX64XS    ComputePoolInstanceFamily = "CPU_X64_XS"
	ComputePoolInstanceFamilyCpuX64S     ComputePoolInstanceFamily = "CPU_X64_S"
	ComputePoolInstanceFamilyCpuX64M     ComputePoolInstanceFamily = "CPU_X64_M"
	ComputePoolInstanceFamilyCpuX64L     ComputePoolInstanceFamily = "CPU_X64_L"
	ComputePoolInstanceFamilyHighMemX64S ComputePoolInstanceFamily = "HIGHMEM_X64_S"
	// Note: Currently the list of instance families in https://docs.snowflake.com/en/sql-reference/sql/create-compute-pool
	// has two entries for HIGHMEM_X64_M. They have the same name, but have different values depending on the region.
	ComputePoolInstanceFamilyHighMemX64M  ComputePoolInstanceFamily = "HIGHMEM_X64_M"
	ComputePoolInstanceFamilyHighMemX64L  ComputePoolInstanceFamily = "HIGHMEM_X64_L"
	ComputePoolInstanceFamilyHighMemX64SL ComputePoolInstanceFamily = "HIGHMEM_X64_SL"
	ComputePoolInstanceFamilyGpuNvS       ComputePoolInstanceFamily = "GPU_NV_S"
	ComputePoolInstanceFamilyGpuNvM       ComputePoolInstanceFamily = "GPU_NV_M"
	ComputePoolInstanceFamilyGpuNvL       ComputePoolInstanceFamily = "GPU_NV_L"
	ComputePoolInstanceFamilyGpuNvXS      ComputePoolInstanceFamily = "GPU_NV_XS"
	ComputePoolInstanceFamilyGpuNvSM      ComputePoolInstanceFamily = "GPU_NV_SM"
	ComputePoolInstanceFamilyGpuNv2M      ComputePoolInstanceFamily = "GPU_NV_2M"
	ComputePoolInstanceFamilyGpuNv3M      ComputePoolInstanceFamily = "GPU_NV_3M"
	ComputePoolInstanceFamilyGpuNvSL      ComputePoolInstanceFamily = "GPU_NV_SL"
)

var AllComputePoolInstanceFamilies = []ComputePoolInstanceFamily{
	ComputePoolInstanceFamilyCpuX64XS,
	ComputePoolInstanceFamilyCpuX64S,
	ComputePoolInstanceFamilyCpuX64M,
	ComputePoolInstanceFamilyCpuX64L,
	ComputePoolInstanceFamilyHighMemX64S,
	ComputePoolInstanceFamilyHighMemX64M,
	ComputePoolInstanceFamilyHighMemX64L,
	ComputePoolInstanceFamilyHighMemX64SL,
	ComputePoolInstanceFamilyGpuNvS,
	ComputePoolInstanceFamilyGpuNvM,
	ComputePoolInstanceFamilyGpuNvL,
	ComputePoolInstanceFamilyGpuNvXS,
	ComputePoolInstanceFamilyGpuNvSM,
	ComputePoolInstanceFamilyGpuNv2M,
	ComputePoolInstanceFamilyGpuNv3M,
	ComputePoolInstanceFamilyGpuNvSL,
}

func ToComputePoolInstanceFamily(s string) (ComputePoolInstanceFamily, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllComputePoolInstanceFamilies, ComputePoolInstanceFamily(s)) {
		return "", fmt.Errorf("invalid compute pool instance family: %s", s)
	}
	return ComputePoolInstanceFamily(s), nil
}

// See https://docs.snowflake.com/en/developer-guide/snowpark-container-services/working-with-compute-pool#compute-pool-lifecycle.
type ComputePoolState string

const (
	ComputePoolStateIdle      ComputePoolState = "IDLE"
	ComputePoolStateActive    ComputePoolState = "ACTIVE"
	ComputePoolStateSuspended ComputePoolState = "SUSPENDED"

	ComputePoolStateStarting ComputePoolState = "STARTING"
	ComputePoolStateStopping ComputePoolState = "STOPPING"
	ComputePoolStateResizing ComputePoolState = "RESIZING"
)

var allComputePoolStates = []ComputePoolState{
	ComputePoolStateIdle,
	ComputePoolStateActive,
	ComputePoolStateSuspended,
	ComputePoolStateStarting,
	ComputePoolStateStopping,
	ComputePoolStateResizing,
}

func ToComputePoolState(s string) (ComputePoolState, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(allComputePoolStates, ComputePoolState(s)) {
		return "", fmt.Errorf("invalid compute pool state: %s", s)
	}
	return ComputePoolState(s), nil
}

var ComputePoolsDef = g.NewInterface(
	"ComputePools",
	"ComputePool",
	g.KindOfT[AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-compute-pool",
	g.NewQueryStruct("CreateComputePool").
		Create().
		SQL("COMPUTE POOL").
		// Note: Currently, OR REPLACE is not supported for compute pools.
		IfNotExists().
		Name().
		OptionalIdentifier("ForApplication", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("FOR APPLICATION")).
		NumberAssignment("MIN_NODES", g.ParameterOptions().Required()).
		NumberAssignment("MAX_NODES", g.ParameterOptions().Required()).
		Assignment(
			"INSTANCE_FAMILY",
			g.KindOfT[ComputePoolInstanceFamily](),
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
