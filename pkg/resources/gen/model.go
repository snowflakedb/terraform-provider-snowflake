package gen

import (
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
)

type ResourceModel struct {
	Name                                string
	NamePlural                          string
	DescribeFailRead                    bool
	DescribeUsesShowId                  bool
	HasPartialOnUpdate                  bool
	AdditionalDescribeOutputDescription string
	DescribeDetailsName                 string
	IdType                              string
	LegacyIdEncoding                    bool
	PreviewFeature                      bool
	IdentityNameForceNew                bool
	IdentityHierarchyForceNew           bool
	AdditionalIdentityNameDescription   string
	ShowColumns                         []*ssot.Attribute
	DescribeColumns                     []*ssot.Attribute
	SdkCreateRequestObjectName          string
	VariantOf                           *ssot.Object
	Attributes                          []ResourceAttribute
	AttributesModification              bool
	ObjectType                          string
	Hooks                               ssot.ResourceHooks
	SchemaStateUpgraders                []ssot.StateUpgraderFunc
	AdditionalCustomizeDiffs            []ssot.CustomizeDiffFunc
	*genhelpers.PreambleModel
}

// ParameterHooks returns the catalog parameters for this object's parameter level, in catalog order.
func (m ResourceModel) ParameterHooks() []ParameterHookModel {
	parameters := ssotdefs.ParameterDefsForLevel(ssot.Object{Name: m.Name}.ParameterLevel())
	hooks := make([]ParameterHookModel, 0, len(parameters))
	for _, parameter := range parameters {
		hooks = append(hooks, newParameterHookModel(parameter))
	}
	return hooks
}

func (m ResourceModel) Prefix() string {
	return genhelpers.FirstLetterLowercase(m.Name)
}

// ShowObjectName is the SDK type of the SHOW row passed to SetConfigFields.
// A variant's row is the parent type.
func (m ResourceModel) ShowObjectName() string {
	if m.VariantOf != nil && m.VariantOf.Name != "" {
		return m.VariantOf.Name
	}
	return m.Name
}

func (m ResourceModel) EffectiveSdkCreateRequestObjectName() string {
	if m.SdkCreateRequestObjectName != "" {
		return m.SdkCreateRequestObjectName
	}
	return m.Name
}

func (m ResourceModel) EffectiveSdkParametersObjectName() string {
	if m.VariantOf != nil {
		return m.VariantOf.Name
	}
	return m.Name
}

// EffectiveSdkAlterObjectName is the SET and UNSET request stem.
// A variant uses its own name, so WarehouseAdaptive is WarehouseAdaptiveSetRequest.
func (m ResourceModel) EffectiveSdkAlterObjectName() string {
	return m.Name
}

func (m ResourceModel) ParseIdFunc() string {
	if m.Hooks.ParseId.Generate {
		return m.Prefix() + "ParseId"
	}
	return m.Prefix() + "ParseIdExt"
}

func (m ResourceModel) HasUpdateCalls() bool {
	for _, attr := range m.Attributes {
		if attr.UpdateCall != nil {
			return true
		}
	}
	return false
}

func ModelFromInputObject(input ResourceObject, preamble *genhelpers.PreambleModel, _ *log.Logger) ResourceModel {
	hooks := ssot.DefaultResourceExtHooks(input.ResourceHooks...)
	resourceAttributes := input.ResourceAttributes()
	attributes := applyResourceAttributeRules(input.IdType, resourceAttributes)
	applyExactlyOneOf(input.ExactlyOneOf, attributes)
	if hooks.NewCreateRequest.Generate || hooks.ApplyCreateOptionals.Generate {
		attachCreateMappings(input.Name, resourceAttributes, attributes)
	}
	if hooks.ApplySetUnset.Generate {
		attachUpdateMappings(input.Name, resourceAttributes, attributes)
	}
	return ResourceModel{
		Name:                                input.Name,
		NamePlural:                          input.NamePlural,
		DescribeFailRead:                    input.DescribeFailRead,
		DescribeUsesShowId:                  input.DescribeUsesShowId,
		HasPartialOnUpdate:                  input.HasPartialOnUpdate,
		AdditionalDescribeOutputDescription: input.AdditionalDescribeOutputDescription,
		DescribeDetailsName:                 input.DescribeDetailsName,
		IdType:                              input.IdType,
		LegacyIdEncoding:                    input.LegacyIdEncoding,
		PreviewFeature:                      input.PreviewFeature,
		IdentityNameForceNew:                input.IdentityNameForceNew,
		IdentityHierarchyForceNew:           input.IdentityHierarchyForceNew,
		AdditionalIdentityNameDescription:   input.AdditionalIdentityNameDescription,
		ShowColumns:                         input.Show.Columns,
		DescribeColumns:                     input.Describe.Columns,
		SdkCreateRequestObjectName:          input.SdkCreateRequestObjectName,
		VariantOf:                           input.VariantOf,
		Attributes:                          attributes,
		AttributesModification:              input.AttributesModification,
		ObjectType:                          input.ObjectType(),
		Hooks:                               hooks,
		SchemaStateUpgraders:                input.SchemaStateUpgraders,
		AdditionalCustomizeDiffs:            input.AdditionalCustomizeDiffs,
		PreambleModel:                       preamble,
	}
}
