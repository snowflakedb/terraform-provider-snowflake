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
	IdType                              string
	IdentityNameForceNew                bool
	AdditionalIdentityNameDescription   string
	Attributes                          []ResourceAttribute
	AttributesModification              bool
	Hooks                               ssot.ResourceHooks
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

func (m ResourceModel) ParseIdFunc() string {
	if m.Hooks.ParseId.Generate {
		return m.Prefix() + "ParseId"
	}
	return m.Prefix() + "ParseIdExt"
}

func (m ResourceModel) HasCreateOptionalCalls() bool {
	for _, attr := range m.Attributes {
		if attr.CreateOptionalCall != nil {
			return true
		}
	}
	return false
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
	attributes := applyResourceAttributeRules(input.IdType, input.Attributes)
	if hooks.NewCreateRequest.Generate || hooks.ApplyCreateOptionals.Generate {
		attachCreateMappings(input.Name, input.Attributes, attributes)
	}
	if hooks.ApplySetUnset.Generate {
		attachUpdateMappings(input.Name, input.Attributes, attributes)
	}
	return ResourceModel{
		Name:                                input.Name,
		NamePlural:                          input.NamePlural,
		DescribeFailRead:                    input.DescribeFailRead,
		DescribeUsesShowId:                  input.DescribeUsesShowId,
		HasPartialOnUpdate:                  input.HasPartialOnUpdate,
		AdditionalDescribeOutputDescription: input.AdditionalDescribeOutputDescription,
		IdType:                              input.IdType,
		IdentityNameForceNew:                input.IdentityNameForceNew,
		AdditionalIdentityNameDescription:   input.AdditionalIdentityNameDescription,
		Attributes:                          attributes,
		AttributesModification:              input.AttributesModification,
		Hooks:                               hooks,
		PreambleModel:                       preamble,
	}
}
