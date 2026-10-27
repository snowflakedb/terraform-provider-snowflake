package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// See https://docs.snowflake.com/en/developer-guide/snowpark-container-services/working-with-compute-pool#compute-pool-lifecycle.
var ComputePoolState = ssot.NewEnum(
	"ComputePoolState", "ComputePoolStates",
	"IDLE", "ACTIVE", "SUSPENDED", "STARTING", "STOPPING", "RESIZING",
)

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

// ComputePoolAttrs holds reusable Compute Pool fields so SHOW Columns (and later
// object variants) can point at the same definitions instead of restating names.
var ComputePoolAttrs = struct {
	ForApplication         ssot.Attribute
	MinNodes               ssot.Attribute
	MaxNodes               ssot.Attribute
	InstanceFamily         ssot.Attribute
	BackupInstanceFamilies ssot.Attribute
	AutoResume             ssot.Attribute
	InitiallySuspended     ssot.Attribute
	AutoSuspendSecs        ssot.Attribute
	Comment                ssot.Attribute
	Name                   ssot.Attribute
	State                  ssot.Attribute
	NumServices            ssot.Attribute
	NumJobs                ssot.Attribute
	ActiveNodes            ssot.Attribute
	IdleNodes              ssot.Attribute
	TargetNodes            ssot.Attribute
	CreatedOn              ssot.Attribute
	ResumedOn              ssot.Attribute
	UpdatedOn              ssot.Attribute
	Owner                  ssot.Attribute
	IsExclusive            ssot.Attribute
	ErrorCode              ssot.Attribute
	StatusMessage          ssot.Attribute
}{
	ForApplication: ssot.Attribute{
		Name:           "for_application",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "AccountObjectIdentifier",
		ForceNew:       true,
		Description:    "Specifies the Snowflake Native App name.",
		OptionalInShow: true,
		OptionalInDb:   true,
		ShowName:       "application",
		CreateSql:      ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Sql("FOR APPLICATION")),
		DiffSuppress:   ssot.SuppressShow.WithKey("application"),
	},
	MinNodes: ssot.Attribute{
		Name:        "min_nodes",
		Type:        ssot.AttrInt,
		Required:    true,
		Description: "Specifies the minimum number of nodes for the compute pool.",
		Min:         new(1),
	},
	MaxNodes: ssot.Attribute{
		Name:        "max_nodes",
		Type:        ssot.AttrInt,
		Required:    true,
		Description: "Specifies the maximum number of nodes for the compute pool.",
		Min:         new(1),
	},
	InstanceFamily: ssot.Attribute{
		Name:                  "instance_family",
		Type:                  ssot.AttrEnum,
		Required:              true,
		ForceNew:              true,
		Description:           "Identifies the type of machine you want to provision for the nodes in the compute pool.",
		AdditionalDescription: "Not all instance families are supported in all regions. Run `SHOW COMPUTE POOL INSTANCE FAMILIES` to see the list of supported instance families in your region.",
		Enum:                  ComputePoolInstanceFamily,
	},
	BackupInstanceFamilies: ssot.Attribute{
		Name:        "backup_instance_families",
		Type:        ssot.AttrList,
		Description: "Specifies an ordered list of instance families to fall back on when the primary `instance_family` is unavailable. The order determines the fallback priority.",
		Elem: &ssot.Attribute{
			Type: ssot.AttrEnum,
			Enum: ComputePoolInstanceFamily,
		},
	},
	AutoResume: ssot.Attribute{
		Name:         "auto_resume",
		Type:         ssot.AttrBool,
		Description:  "Specifies whether to automatically resume a compute pool when a service or job is submitted to it.",
		DiffSuppress: ssot.SuppressShow,
	},
	InitiallySuspended: ssot.Attribute{
		Name:           "initially_suspended",
		Type:           ssot.AttrBool,
		CreateOnly:     true,
		Description:    "Specifies whether the compute pool is created initially in the suspended state. This field is used only when creating a compute pool. Changes on this field are ignored after creation.",
		RawDescription: true,
	},
	AutoSuspendSecs: ssot.Attribute{
		Name:         "auto_suspend_secs",
		Type:         ssot.AttrInt,
		Description:  "Number of seconds of inactivity after which you want Snowflake to automatically suspend the compute pool.",
		Min:          new(0),
		DiffSuppress: ssot.SuppressShow,
	},
	Comment: ssot.Attribute{
		Name:           "comment",
		Type:           ssot.AttrString,
		Description:    "Specifies a comment for the compute pool.",
		OptionalInShow: true,
		OptionalInDb:   true,
	},
	Name:          ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	State:         ssot.Attribute{Name: "state", Type: ssot.AttrEnum, Enum: ComputePoolState, Readonly: true},
	NumServices:   ssot.Attribute{Name: "num_services", Type: ssot.AttrInt, Readonly: true},
	NumJobs:       ssot.Attribute{Name: "num_jobs", Type: ssot.AttrInt, Readonly: true},
	ActiveNodes:   ssot.Attribute{Name: "active_nodes", Type: ssot.AttrInt, Readonly: true},
	IdleNodes:     ssot.Attribute{Name: "idle_nodes", Type: ssot.AttrInt, Readonly: true},
	TargetNodes:   ssot.Attribute{Name: "target_nodes", Type: ssot.AttrInt, Readonly: true},
	CreatedOn:     ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true},
	ResumedOn:     ssot.Attribute{Name: "resumed_on", Type: ssot.AttrTime, Readonly: true},
	UpdatedOn:     ssot.Attribute{Name: "updated_on", Type: ssot.AttrTime, Readonly: true},
	Owner:         ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true},
	IsExclusive:   ssot.Attribute{Name: "is_exclusive", Type: ssot.AttrBool, Readonly: true},
	ErrorCode:     ssot.Attribute{Name: "error_code", Type: ssot.AttrString, Readonly: true},
	StatusMessage: ssot.Attribute{Name: "status_message", Type: ssot.AttrString, Readonly: true},
}

var ComputePool = ssot.Object{
	Name:                 "ComputePool",
	NamePlural:           "ComputePools",
	DescribeFailRead:     true,
	IdType:               "AccountObjectIdentifier",
	IdentityNameForceNew: true,
	Generators: []ssot.GeneratorOption{
		ssot.WithSchema(ssot.SchemaShow(), ssot.SchemaDescribe()),
		ssot.WithModel(),
		ssot.WithResourceAssert(),
		ssot.WithShowOutputAssert(ssot.ShowOutputShow(), ssot.ShowOutputDescribe()),
		ssot.WithObjectAssert(ssot.ObjectAssertShow(), ssot.ObjectAssertDescribe()),
		ssot.WithSdk(
			ssot.SdkDefineFromSsot(),
			ssot.SdkShow(
				ssot.SdkShowOptionalLike(),
				ssot.SdkShowOptionalStartsWith(),
				ssot.SdkShowOptionalLimitFrom(),
			),
			ssot.SdkDescribe(),
			ssot.SdkDrop(ssot.SdkDropSafelyHook()),
			ssot.SdkCreate(
				ssot.SdkCreateWithoutOrReplace(),
				ssot.SdkCreateAdditionalValidations(),
			),
			ssot.SdkAlter(
				ssot.SdkAlterSetAdditionalValidations(),
			),
		),
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithParseId(ssot.Generated()),
		ssot.WithNewCreateRequest(ssot.Generated()),
		ssot.WithApplyCreateOptionals(ssot.Generated()),
		ssot.WithCreateInSdk(ssot.Generated()),
		ssot.WithApplySetUnset(ssot.Generated()),
		ssot.WithAlterSetInSdk(ssot.Generated()),
		ssot.WithAlterUnsetInSdk(ssot.Generated()),
		ssot.WithDropSafelyInSdk(ssot.Generated()),
		ssot.WithShowByIdInSdk(ssot.Generated()),
		ssot.WithShowByIdSafelyInSdk(ssot.Generated()),
		ssot.WithDescribeInSdk(ssot.Generated()),
	},
	Attributes: []ssot.Attribute{
		ComputePoolAttrs.ForApplication,
		ComputePoolAttrs.MinNodes,
		ComputePoolAttrs.MaxNodes,
		ComputePoolAttrs.InstanceFamily,
		ComputePoolAttrs.BackupInstanceFamilies,
		ComputePoolAttrs.AutoResume,
		ComputePoolAttrs.InitiallySuspended,
		ComputePoolAttrs.AutoSuspendSecs,
		ComputePoolAttrs.Comment,
	},
	Create: ssot.CreateSpec{
		Clauses: []*ssot.Attribute{
			&ComputePoolAttrs.ForApplication,
			&ComputePoolAttrs.MinNodes,
			&ComputePoolAttrs.MaxNodes,
			&ComputePoolAttrs.InstanceFamily,
			&ComputePoolAttrs.AutoResume,
			&ComputePoolAttrs.InitiallySuspended,
			&ComputePoolAttrs.AutoSuspendSecs,
			ssot.Tags,
			&ComputePoolAttrs.Comment,
			&ComputePoolAttrs.BackupInstanceFamilies,
		},
	},
	Set: []*ssot.Attribute{
		&ComputePoolAttrs.MinNodes,
		&ComputePoolAttrs.MaxNodes,
		&ComputePoolAttrs.AutoResume,
		&ComputePoolAttrs.AutoSuspendSecs,
		&ComputePoolAttrs.BackupInstanceFamilies,
		&ComputePoolAttrs.Comment,
	},
	Unset: []*ssot.Attribute{
		&ComputePoolAttrs.AutoResume,
		&ComputePoolAttrs.AutoSuspendSecs,
		&ComputePoolAttrs.BackupInstanceFamilies,
		&ComputePoolAttrs.Comment,
	},
	Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
		ssot.AlterKeyword("RESUME"),
		ssot.AlterKeyword("SUSPEND"),
		ssot.AlterKeyword("STOP ALL"),
		ssot.AlterSet(),
		ssot.AlterUnset(),
		ssot.AlterSetTags(),
		ssot.AlterUnsetTags(),
	}},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&ComputePoolAttrs.Name,
			&ComputePoolAttrs.State,
			&ComputePoolAttrs.MinNodes,
			&ComputePoolAttrs.MaxNodes,
			&ComputePoolAttrs.InstanceFamily,
			&ComputePoolAttrs.NumServices,
			&ComputePoolAttrs.NumJobs,
			&ComputePoolAttrs.AutoSuspendSecs,
			&ComputePoolAttrs.AutoResume,
			&ComputePoolAttrs.ActiveNodes,
			&ComputePoolAttrs.IdleNodes,
			&ComputePoolAttrs.TargetNodes,
			&ComputePoolAttrs.CreatedOn,
			&ComputePoolAttrs.ResumedOn,
			&ComputePoolAttrs.UpdatedOn,
			&ComputePoolAttrs.Owner,
			&ComputePoolAttrs.Comment,
			&ComputePoolAttrs.IsExclusive,
			&ComputePoolAttrs.ForApplication,
			&ComputePoolAttrs.BackupInstanceFamilies,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&ComputePoolAttrs.Name,
			&ComputePoolAttrs.State,
			&ComputePoolAttrs.MinNodes,
			&ComputePoolAttrs.MaxNodes,
			&ComputePoolAttrs.InstanceFamily,
			&ComputePoolAttrs.NumServices,
			&ComputePoolAttrs.NumJobs,
			&ComputePoolAttrs.AutoSuspendSecs,
			&ComputePoolAttrs.AutoResume,
			&ComputePoolAttrs.ActiveNodes,
			&ComputePoolAttrs.IdleNodes,
			&ComputePoolAttrs.TargetNodes,
			&ComputePoolAttrs.CreatedOn,
			&ComputePoolAttrs.ResumedOn,
			&ComputePoolAttrs.UpdatedOn,
			&ComputePoolAttrs.Owner,
			&ComputePoolAttrs.Comment,
			&ComputePoolAttrs.IsExclusive,
			&ComputePoolAttrs.ForApplication,
			&ComputePoolAttrs.BackupInstanceFamilies,
			&ComputePoolAttrs.ErrorCode,
			&ComputePoolAttrs.StatusMessage,
		},
	},
}
