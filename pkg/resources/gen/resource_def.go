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
	name:                 "ComputePool",
	describeFailRead:     true,
	idType:               "AccountObjectIdentifier",
	identityNameForceNew: true,
}

var schemaDef = ResourceDef{
	name:                                "Schema",
	describeUsesShowId:                  true,
	hasPartialOnUpdate:                  true,
	additionalDescribeOutputDescription: "In order to handle this output, one must grant sufficient privileges, e.g. [grant_ownership](./grant_ownership) on all objects in the schema.",
	idType:                              "DatabaseObjectIdentifier",
	additionalIdentityNameDescription:   "When the name is `PUBLIC`, during creation the provider checks if this schema has already been created and, in such case, `ALTER` is used to match the desired state.",
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
	},
	ObjectGenerationSettings: &genhelpers.ObjectGenerationSettings{
		AllowedGenerationParts: []genhelpers.GenerationPartNamer{PartDefault, PartSchema, PartParameters},
		EnabledGenerationParts: []genhelpers.GenerationPartNamer{PartParameters},
	},
}
