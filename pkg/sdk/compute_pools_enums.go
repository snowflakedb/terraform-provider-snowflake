package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type ComputePoolInstanceFamily string

const (
	// Previous generation (AWS, Azure, GCP)
	ComputePoolInstanceFamilyCpuX64XS    ComputePoolInstanceFamily = "CPU_X64_XS"
	ComputePoolInstanceFamilyCpuX64S     ComputePoolInstanceFamily = "CPU_X64_S"
	ComputePoolInstanceFamilyCpuX64M     ComputePoolInstanceFamily = "CPU_X64_M"
	ComputePoolInstanceFamilyCpuX64SL    ComputePoolInstanceFamily = "CPU_X64_SL"
	ComputePoolInstanceFamilyCpuX64L     ComputePoolInstanceFamily = "CPU_X64_L"
	ComputePoolInstanceFamilyHighMemX64S ComputePoolInstanceFamily = "HIGHMEM_X64_S"
	ComputePoolInstanceFamilyHighMemX64M ComputePoolInstanceFamily = "HIGHMEM_X64_M"
	ComputePoolInstanceFamilyHighMemX64L ComputePoolInstanceFamily = "HIGHMEM_X64_L"
	// Azure and GCP only
	ComputePoolInstanceFamilyHighMemX64SL ComputePoolInstanceFamily = "HIGHMEM_X64_SL"
	// AWS and Azure
	ComputePoolInstanceFamilyGpuNvS ComputePoolInstanceFamily = "GPU_NV_S"
	ComputePoolInstanceFamilyGpuNvM ComputePoolInstanceFamily = "GPU_NV_M"
	ComputePoolInstanceFamilyGpuNvL ComputePoolInstanceFamily = "GPU_NV_L"
	// Azure only
	ComputePoolInstanceFamilyGpuNvXS ComputePoolInstanceFamily = "GPU_NV_XS"
	ComputePoolInstanceFamilyGpuNvSM ComputePoolInstanceFamily = "GPU_NV_SM"
	ComputePoolInstanceFamilyGpuNv2M ComputePoolInstanceFamily = "GPU_NV_2M"
	ComputePoolInstanceFamilyGpuNv3M ComputePoolInstanceFamily = "GPU_NV_3M"
	ComputePoolInstanceFamilyGpuNvSL ComputePoolInstanceFamily = "GPU_NV_SL"
	// GCP only
	ComputePoolInstanceFamilyGpuGcpNvL4_1_24G ComputePoolInstanceFamily = "GPU_GCP_NV_L4_1_24G"
	ComputePoolInstanceFamilyGpuGcpNvL4_4_24G ComputePoolInstanceFamily = "GPU_GCP_NV_L4_4_24G"
	ComputePoolInstanceFamilyGpuGcpNvA100840G ComputePoolInstanceFamily = "GPU_GCP_NV_A100_8_40G"
	// Current generation - AWS only
	ComputePoolInstanceFamilyGenArmG1_2  ComputePoolInstanceFamily = "GEN_ARM_G1_2"
	ComputePoolInstanceFamilyGenArmG1_4  ComputePoolInstanceFamily = "GEN_ARM_G1_4"
	ComputePoolInstanceFamilyGenArmG1_8  ComputePoolInstanceFamily = "GEN_ARM_G1_8"
	ComputePoolInstanceFamilyGenArmG1_16 ComputePoolInstanceFamily = "GEN_ARM_G1_16"
	ComputePoolInstanceFamilyGenArmG1_32 ComputePoolInstanceFamily = "GEN_ARM_G1_32"
	// Current generation - AWS and Azure
	ComputePoolInstanceFamilyGenX64G2_2  ComputePoolInstanceFamily = "GEN_X64_G2_2"
	ComputePoolInstanceFamilyGenX64G2_4  ComputePoolInstanceFamily = "GEN_X64_G2_4"
	ComputePoolInstanceFamilyGenX64G2_8  ComputePoolInstanceFamily = "GEN_X64_G2_8"
	ComputePoolInstanceFamilyGenX64G2_32 ComputePoolInstanceFamily = "GEN_X64_G2_32"
	// Azure only
	ComputePoolInstanceFamilyGenX64G2_16 ComputePoolInstanceFamily = "GEN_X64_G2_16"
	// Current generation - AWS and Azure
	ComputePoolInstanceFamilyMemX64G2_8   ComputePoolInstanceFamily = "MEM_X64_G2_8"
	ComputePoolInstanceFamilyMemX64G2_32  ComputePoolInstanceFamily = "MEM_X64_G2_32"
	ComputePoolInstanceFamilyMemX64G2_64  ComputePoolInstanceFamily = "MEM_X64_G2_64"
	ComputePoolInstanceFamilyMemX64G2_192 ComputePoolInstanceFamily = "MEM_X64_G2_192"
	// Azure only
	ComputePoolInstanceFamilyMemX64G2_96 ComputePoolInstanceFamily = "MEM_X64_G2_96"
	// Current generation - AWS only
	ComputePoolInstanceFamilyGpuL40SG1_8   ComputePoolInstanceFamily = "GPU_L40S_G1_8"
	ComputePoolInstanceFamilyGpuL40SG1_16  ComputePoolInstanceFamily = "GPU_L40S_G1_16"
	ComputePoolInstanceFamilyGpuL40SG1_48  ComputePoolInstanceFamily = "GPU_L40S_G1_48"
	ComputePoolInstanceFamilyGpuL40SG1_192 ComputePoolInstanceFamily = "GPU_L40S_G1_192"
	ComputePoolInstanceFamilyGpuR6KG1_8    ComputePoolInstanceFamily = "GPU_R6K_G1_8"
	ComputePoolInstanceFamilyGpuR6KG1_16   ComputePoolInstanceFamily = "GPU_R6K_G1_16"
	ComputePoolInstanceFamilyGpuR6KG1_32   ComputePoolInstanceFamily = "GPU_R6K_G1_32"
	ComputePoolInstanceFamilyGpuR6KG1_48   ComputePoolInstanceFamily = "GPU_R6K_G1_48"
	ComputePoolInstanceFamilyGpuR6KG1_96   ComputePoolInstanceFamily = "GPU_R6K_G1_96"
	ComputePoolInstanceFamilyGpuR6KG1_192  ComputePoolInstanceFamily = "GPU_R6K_G1_192"
	// GCP only
	ComputePoolInstanceFamilyGpuA100G1_12 ComputePoolInstanceFamily = "GPU_A100_G1_12"
	ComputePoolInstanceFamilyGpuA100G1_48 ComputePoolInstanceFamily = "GPU_A100_G1_48"
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
	ComputePoolInstanceFamilyGenArmG1_2,
	ComputePoolInstanceFamilyGenArmG1_4,
	ComputePoolInstanceFamilyGenArmG1_8,
	ComputePoolInstanceFamilyGenArmG1_16,
	ComputePoolInstanceFamilyGenArmG1_32,
	ComputePoolInstanceFamilyGenX64G2_2,
	ComputePoolInstanceFamilyGenX64G2_4,
	ComputePoolInstanceFamilyGenX64G2_8,
	ComputePoolInstanceFamilyGenX64G2_16,
	ComputePoolInstanceFamilyGenX64G2_32,
	ComputePoolInstanceFamilyMemX64G2_8,
	ComputePoolInstanceFamilyMemX64G2_32,
	ComputePoolInstanceFamilyMemX64G2_64,
	ComputePoolInstanceFamilyMemX64G2_96,
	ComputePoolInstanceFamilyMemX64G2_192,
	ComputePoolInstanceFamilyGpuL40SG1_8,
	ComputePoolInstanceFamilyGpuL40SG1_16,
	ComputePoolInstanceFamilyGpuL40SG1_48,
	ComputePoolInstanceFamilyGpuL40SG1_192,
	ComputePoolInstanceFamilyGpuR6KG1_8,
	ComputePoolInstanceFamilyGpuR6KG1_16,
	ComputePoolInstanceFamilyGpuR6KG1_32,
	ComputePoolInstanceFamilyGpuR6KG1_48,
	ComputePoolInstanceFamilyGpuR6KG1_96,
	ComputePoolInstanceFamilyGpuR6KG1_192,
	ComputePoolInstanceFamilyGpuA100G1_12,
	ComputePoolInstanceFamilyGpuA100G1_48,
}

func ToComputePoolInstanceFamily(s string) (ComputePoolInstanceFamily, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllComputePoolInstanceFamilies, ComputePoolInstanceFamily(s)) {
		return "", fmt.Errorf("invalid compute pool instance family: %s", s)
	}
	return ComputePoolInstanceFamily(s), nil
}
