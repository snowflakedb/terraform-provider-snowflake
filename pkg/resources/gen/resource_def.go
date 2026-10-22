package gen

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"

type ResourceDef struct {
	name                                string
	describeFailRead                    bool
	describeUsesShowId                  bool
	hasPartialOnUpdate                  bool
	additionalDescribeOutputDescription string
	idType                              string
	identityNameForceNew                bool
	additionalIdentityNameDescription   string
	attributes                          []Attribute
	attributesModification              bool
	hooks                               []HookOption
	*genhelpers.ObjectGenerationSettings
}

func (d ResourceDef) ObjectName() string { return d.name }

func GetAllObjects() []ResourceDef {
	return []ResourceDef{
		computePoolDef,
		schemaDef,
	}
}

var computePoolDef = ResourceDef{
	name:                   "ComputePool",
	describeFailRead:       true,
	idType:                 "AccountObjectIdentifier",
	identityNameForceNew:   true,
	attributesModification: true,
	attributes: []Attribute{
		{
			Name:        "for_application",
			Type:        AttrIdentifier,
			ForceNew:    true,
			Description: "Specifies the Snowflake Native App name.",
		},
		{
			Name:        "min_nodes",
			Type:        AttrInt,
			Required:    true,
			Description: "Specifies the minimum number of nodes for the compute pool.",
			Min:         new(1),
		},
		{
			Name:        "max_nodes",
			Type:        AttrInt,
			Required:    true,
			Description: "Specifies the maximum number of nodes for the compute pool.",
			Min:         new(1),
		},
		{
			Name:                  "instance_family",
			Type:                  AttrEnum,
			Required:              true,
			ForceNew:              true,
			Description:           "Identifies the type of machine you want to provision for the nodes in the compute pool.",
			AdditionalDescription: "Not all instance families are supported in all regions. Run `SHOW COMPUTE POOL INSTANCE FAMILIES` to see the list of supported instance families in your region.",
			EnumPlural:            "InstanceFamilies",
		},
		{
			Name:        "backup_instance_families",
			Type:        AttrList,
			Description: "Specifies an ordered list of instance families to fall back on when the primary `instance_family` is unavailable. The order determines the fallback priority.",
			Elem: &Attribute{
				Type:       AttrEnum,
				Enum:       "InstanceFamily",
				EnumPlural: "InstanceFamilies",
			},
		},
		{
			Name:        "auto_resume",
			Type:        AttrBool,
			Description: "Specifies whether to automatically resume a compute pool when a service or job is submitted to it.",
		},
		{
			Name:           "initially_suspended",
			Type:           AttrBool,
			Description:    "Specifies whether the compute pool is created initially in the suspended state. This field is used only when creating a compute pool. Changes on this field are ignored after creation.",
			RawDescription: true,
		},
		{
			Name:        "auto_suspend_secs",
			Type:        AttrInt,
			Description: "Number of seconds of inactivity after which you want Snowflake to automatically suspend the compute pool.",
			Min:         new(0),
		},
		{
			Name:        "comment",
			Type:        AttrString,
			Description: "Specifies a comment for the compute pool.",
		},
	},
}

var schemaDef = ResourceDef{
	name:                                "Schema",
	describeUsesShowId:                  true,
	hasPartialOnUpdate:                  true,
	additionalDescribeOutputDescription: "In order to handle this output, one must grant sufficient privileges, e.g. [grant_ownership](./grant_ownership) on all objects in the schema.",
	idType:                              "DatabaseObjectIdentifier",
	additionalIdentityNameDescription:   "When the name is `PUBLIC`, during creation the provider checks if this schema has already been created and, in such case, `ALTER` is used to match the desired state.",
	attributesModification:              true,
	attributes: []Attribute{
		{
			Name:        "with_managed_access",
			Type:        AttrBool,
			Description: "Specifies a managed schema. Managed access schemas centralize privilege management with the schema owner.",
		},
		{
			Name:        "is_transient",
			Type:        AttrBool,
			ForceNew:    true,
			Description: "Specifies the schema as transient. Transient schemas do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss.",
		},
		{
			Name:        "comment",
			Type:        AttrString,
			Description: "Specifies a comment for the schema.",
		},
	},
	hooks: []HookOption{
		WithBeforeCreate(extCall()),
		WithHierarchyRename(extCall()),
		WithBeforeAlter(extCall()),
		WithDescribeOutputSet(omitted()),
		WithSetDescribeOutput(extCall()),
		WithSchemaVersion(extCall()),
		WithApplyParametersCreate(extCall()),
		WithShowParametersInSdk(extCall()),
		WithParametersDetailsFromRaw(extCall()),
		WithSetParametersFields(extCall()),
		WithParametersOutputSet(extCall()),
		WithApplyParametersChanges(extCall()),
		WithParametersAttributesSchema(generated()),
		WithParametersOutputSchema(generated()),
		WithParametersFieldNames(generated()),
		WithParametersCustomDiff(generated()),
	},
	ObjectGenerationSettings: &genhelpers.ObjectGenerationSettings{
		AllowedGenerationParts: []genhelpers.GenerationPartNamer{PartDefault, PartSchema, PartParameters},
		EnabledGenerationParts: []genhelpers.GenerationPartNamer{PartParameters},
	},
}
