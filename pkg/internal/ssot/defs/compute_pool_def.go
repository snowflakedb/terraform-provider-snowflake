package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

var ComputePoolInstanceFamily = ssot.NewEnum(
	"ComputePoolInstanceFamily", "ComputePoolInstanceFamilies",
	// Previous generation (AWS, Azure, GCP)
	"CPU_X64_XS",
	"CPU_X64_S",
	"CPU_X64_M",
	"CPU_X64_SL",
	"CPU_X64_L",
	"HIGHMEM_X64_S",
	// Note: Currently the list of instance families in https://docs.snowflake.com/en/sql-reference/sql/create-compute-pool
	// has two entries for HIGHMEM_X64_M. They have the same name, but have different values depending on the region.
	"HIGHMEM_X64_M",
	"HIGHMEM_X64_L",
	// Azure and GCP only
	"HIGHMEM_X64_SL",
	// AWS and Azure
	"GPU_NV_S",
	"GPU_NV_M",
	"GPU_NV_L",
	// Azure only
	"GPU_NV_XS",
	"GPU_NV_SM",
	"GPU_NV_2M",
	"GPU_NV_3M",
	"GPU_NV_SL",
	// GCP only
	"GPU_GCP_NV_L4_1_24G",
	"GPU_GCP_NV_L4_4_24G",
	"GPU_GCP_NV_A100_8_40G",
	// Current generation - AWS only
	"GEN_ARM_G1_2",
	"GEN_ARM_G1_4",
	"GEN_ARM_G1_8",
	"GEN_ARM_G1_16",
	"GEN_ARM_G1_32",
	// Current generation - AWS and Azure
	"GEN_X64_G2_2",
	"GEN_X64_G2_4",
	"GEN_X64_G2_8",
	"GEN_X64_G2_32",
	// Azure only
	"GEN_X64_G2_16",
	// Current generation - AWS and Azure
	"MEM_X64_G2_8",
	"MEM_X64_G2_32",
	"MEM_X64_G2_64",
	"MEM_X64_G2_192",
	// Azure only
	"MEM_X64_G2_96",
	// Current generation - AWS only
	"GPU_L40S_G1_8",
	"GPU_L40S_G1_16",
	"GPU_L40S_G1_48",
	"GPU_L40S_G1_192",
	"GPU_R6K_G1_8",
	"GPU_R6K_G1_16",
	"GPU_R6K_G1_32",
	"GPU_R6K_G1_48",
	"GPU_R6K_G1_96",
	"GPU_R6K_G1_192",
	// GCP only
	"GPU_A100_G1_12",
	"GPU_A100_G1_48",
)

var ComputePool = ssot.Object{
	Name:                   "ComputePool",
	NamePlural:             "ComputePools",
	DescribeFailRead:       true,
	IdType:                 "AccountObjectIdentifier",
	IdentityNameForceNew:   true,
	AttributesModification: true,
	Attributes: []ssot.Attribute{
		{
			Name:        "for_application",
			Type:        ssot.AttrIdentifier,
			ForceNew:    true,
			Description: "Specifies the Snowflake Native App name.",
		},
		{
			Name:        "min_nodes",
			Type:        ssot.AttrInt,
			Required:    true,
			Description: "Specifies the minimum number of nodes for the compute pool.",
			Min:         new(1),
		},
		{
			Name:        "max_nodes",
			Type:        ssot.AttrInt,
			Required:    true,
			Description: "Specifies the maximum number of nodes for the compute pool.",
			Min:         new(1),
		},
		{
			Name:                  "instance_family",
			Type:                  ssot.AttrEnum,
			Required:              true,
			ForceNew:              true,
			Description:           "Identifies the type of machine you want to provision for the nodes in the compute pool.",
			AdditionalDescription: "Not all instance families are supported in all regions. Run `SHOW COMPUTE POOL INSTANCE FAMILIES` to see the list of supported instance families in your region.",
			Enum:                  ComputePoolInstanceFamily,
		},
		{
			Name:        "backup_instance_families",
			Type:        ssot.AttrList,
			Description: "Specifies an ordered list of instance families to fall back on when the primary `instance_family` is unavailable. The order determines the fallback priority.",
			Elem: &ssot.Attribute{
				Type: ssot.AttrEnum,
				Enum: ComputePoolInstanceFamily,
			},
		},
		{
			Name:        "auto_resume",
			Type:        ssot.AttrBool,
			Description: "Specifies whether to automatically resume a compute pool when a service or job is submitted to it.",
		},
		{
			Name:           "initially_suspended",
			Type:           ssot.AttrBool,
			Description:    "Specifies whether the compute pool is created initially in the suspended state. This field is used only when creating a compute pool. Changes on this field are ignored after creation.",
			RawDescription: true,
		},
		{
			Name:        "auto_suspend_secs",
			Type:        ssot.AttrInt,
			Description: "Number of seconds of inactivity after which you want Snowflake to automatically suspend the compute pool.",
			Min:         new(0),
		},
		{
			Name:        "comment",
			Type:        ssot.AttrString,
			Description: "Specifies a comment for the compute pool.",
		},
	},
}
