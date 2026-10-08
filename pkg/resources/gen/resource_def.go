package gen

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"

type ResourceDef struct {
	name               string
	describeFailRead   bool
	describeUsesShowId bool
	hasPartialOnUpdate bool
	hooks              []HookOption
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
	name:             "ComputePool",
	describeFailRead: true,
}

var schemaDef = ResourceDef{
	name:               "Schema",
	describeUsesShowId: true,
	hasPartialOnUpdate: true,
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
		WithParametersAttributesSchema(extCall()),
		WithParametersOutputSchema(extCall()),
	},
	ObjectGenerationSettings: &genhelpers.ObjectGenerationSettings{
		AllowedGenerationParts: []genhelpers.GenerationPartNamer{PartDefault, PartSchema, PartParameters},
		EnabledGenerationParts: []genhelpers.GenerationPartNamer{PartParameters},
	},
}
