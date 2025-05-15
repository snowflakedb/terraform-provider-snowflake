package sdk

import (
	"fmt"
	"slices"
	"strings"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

//go:generate go run ./poc/main.go

type ComputePoolInstanceFamily string

const (
	ComputePoolInstanceFamilyCPUX64XS    ComputePoolInstanceFamily = "CPU_X64_XS"
	ComputePoolInstanceFamilyCPUX64S     ComputePoolInstanceFamily = "CPU_X64_S"
	ComputePoolInstanceFamilyCPUX64M     ComputePoolInstanceFamily = "CPU_X64_M"
	ComputePoolInstanceFamilyCPUX64L     ComputePoolInstanceFamily = "CPU_X64_L"
	ComputePoolInstanceFamilyHIGHMEMX64S ComputePoolInstanceFamily = "HIGHMEM_X64_S"
	// Note: Currently the list of instance families in https://docs.snowflake.com/en/sql-reference/sql/create-compute-pool
	// has two entries for HIGHMEM_X64_M. They have the same name, but have different values depending on the region.
	ComputePoolInstanceFamilyHIGHMEMX64M  ComputePoolInstanceFamily = "HIGHMEM_X64_M"
	ComputePoolInstanceFamilyHIGHMEMX64L  ComputePoolInstanceFamily = "HIGHMEM_X64_L"
	ComputePoolInstanceFamilyHIGHMEMX64SL ComputePoolInstanceFamily = "HIGHMEM_X64_SL"
	ComputePoolInstanceFamilyGPUNVS       ComputePoolInstanceFamily = "GPU_NV_S"
	ComputePoolInstanceFamilyGPUNVM       ComputePoolInstanceFamily = "GPU_NV_M"
	ComputePoolInstanceFamilyGPUNVL       ComputePoolInstanceFamily = "GPU_NV_L"
	ComputePoolInstanceFamilyGPUNVXS      ComputePoolInstanceFamily = "GPU_NV_XS"
	ComputePoolInstanceFamilyGPUNVSM      ComputePoolInstanceFamily = "GPU_NV_SM"
	ComputePoolInstanceFamilyGPUNV2M      ComputePoolInstanceFamily = "GPU_NV_2M"
	ComputePoolInstanceFamilyGPUNV3M      ComputePoolInstanceFamily = "GPU_NV_3M"
	ComputePoolInstanceFamilyGPUNVSL      ComputePoolInstanceFamily = "GPU_NV_SL"
)

var allComputePoolInstanceFamilies = []ComputePoolInstanceFamily{
	ComputePoolInstanceFamilyCPUX64XS,
	ComputePoolInstanceFamilyCPUX64S,
	ComputePoolInstanceFamilyCPUX64M,
	ComputePoolInstanceFamilyCPUX64L,
	ComputePoolInstanceFamilyHIGHMEMX64S,
	ComputePoolInstanceFamilyHIGHMEMX64M,
	ComputePoolInstanceFamilyHIGHMEMX64L,
	ComputePoolInstanceFamilyHIGHMEMX64SL,
	ComputePoolInstanceFamilyGPUNVS,
	ComputePoolInstanceFamilyGPUNVM,
	ComputePoolInstanceFamilyGPUNVL,
	ComputePoolInstanceFamilyGPUNVXS,
	ComputePoolInstanceFamilyGPUNVSM,
	ComputePoolInstanceFamilyGPUNV2M,
	ComputePoolInstanceFamilyGPUNV3M,
	ComputePoolInstanceFamilyGPUNVSL,
}

func ToComputePoolInstanceFamily(s string) (ComputePoolInstanceFamily, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(allComputePoolInstanceFamilies, ComputePoolInstanceFamily(s)) {
		return "", fmt.Errorf("invalid compute pool instance family: %s", s)
	}
	return ComputePoolInstanceFamily(s), nil
}

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
		WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags"),
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
		OptionalText("application").
		OptionalText("budget"),
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
		OptionalText("Application").
		OptionalText("Budget"),
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
		OptionalText("budget").
		OptionalText("error_code").
		OptionalText("status_message"),
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
		OptionalText("Application").
		OptionalText("Budget").
		OptionalText("ErrorCode").
		OptionalText("StatusMessage"),
	g.NewQueryStruct("DescComputePool").
		Describe().
		SQL("COMPUTE POOL").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
