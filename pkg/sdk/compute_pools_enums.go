package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type ComputePoolInstanceFamily string

const (
	ComputePoolInstanceFamilyCpuX64XS    ComputePoolInstanceFamily = "CPU_X64_XS"
	ComputePoolInstanceFamilyCpuX64S     ComputePoolInstanceFamily = "CPU_X64_S"
	ComputePoolInstanceFamilyCpuX64M     ComputePoolInstanceFamily = "CPU_X64_M"
	ComputePoolInstanceFamilyCpuX64SL    ComputePoolInstanceFamily = "CPU_X64_SL"
	ComputePoolInstanceFamilyCpuX64L     ComputePoolInstanceFamily = "CPU_X64_L"
	ComputePoolInstanceFamilyHighMemX64S ComputePoolInstanceFamily = "HIGHMEM_X64_S"
	// Note: Currently the list of instance families in https://docs.snowflake.com/en/sql-reference/sql/create-compute-pool
	// has two entries for HIGHMEM_X64_M. They have the same name, but have different values depending on the region.
	ComputePoolInstanceFamilyHighMemX64M      ComputePoolInstanceFamily = "HIGHMEM_X64_M"
	ComputePoolInstanceFamilyHighMemX64L      ComputePoolInstanceFamily = "HIGHMEM_X64_L"
	ComputePoolInstanceFamilyHighMemX64SL     ComputePoolInstanceFamily = "HIGHMEM_X64_SL"
	ComputePoolInstanceFamilyGpuNvS           ComputePoolInstanceFamily = "GPU_NV_S"
	ComputePoolInstanceFamilyGpuNvM           ComputePoolInstanceFamily = "GPU_NV_M"
	ComputePoolInstanceFamilyGpuNvL           ComputePoolInstanceFamily = "GPU_NV_L"
	ComputePoolInstanceFamilyGpuNvXS          ComputePoolInstanceFamily = "GPU_NV_XS"
	ComputePoolInstanceFamilyGpuNvSM          ComputePoolInstanceFamily = "GPU_NV_SM"
	ComputePoolInstanceFamilyGpuNv2M          ComputePoolInstanceFamily = "GPU_NV_2M"
	ComputePoolInstanceFamilyGpuNv3M          ComputePoolInstanceFamily = "GPU_NV_3M"
	ComputePoolInstanceFamilyGpuNvSL          ComputePoolInstanceFamily = "GPU_NV_SL"
	ComputePoolInstanceFamilyGpuGcpNvL4_1_24G ComputePoolInstanceFamily = "GPU_GCP_NV_L4_1_24G"
	ComputePoolInstanceFamilyGpuGcpNvL4_4_24G ComputePoolInstanceFamily = "GPU_GCP_NV_L4_4_24G"
	ComputePoolInstanceFamilyGpuGcpNvA100840G ComputePoolInstanceFamily = "GPU_GCP_NV_A100_8_40G"
)

var AllComputePoolInstanceFamilies = []ComputePoolInstanceFamily{
	ComputePoolInstanceFamilyCpuX64XS,
	ComputePoolInstanceFamilyCpuX64S,
	ComputePoolInstanceFamilyCpuX64M,
	ComputePoolInstanceFamilyCpuX64SL,
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
	ComputePoolInstanceFamilyGpuGcpNvL4_1_24G,
	ComputePoolInstanceFamilyGpuGcpNvL4_4_24G,
	ComputePoolInstanceFamilyGpuGcpNvA100840G,
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
