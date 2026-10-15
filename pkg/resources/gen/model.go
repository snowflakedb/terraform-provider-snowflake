package gen

import (
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
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

func (m ResourceModel) Prefix() string {
	return genhelpers.FirstLetterLowercase(m.Name)
}

func (m ResourceModel) HasCreateOptionalCalls() bool {
	for _, attr := range m.Attributes {
		if attr.CreateOptionalCall != nil {
			return true
		}
	}
	return false
}

func ModelFromInputObject(input ResourceObject, preamble *genhelpers.PreambleModel, _ *log.Logger) ResourceModel {
	hooks := ssot.DefaultResourceExtHooks(input.ResourceHooks...)
	attributes := applyResourceAttributeRules(input.Attributes)
	if hooks.NewCreateRequest.Generate || hooks.ApplyCreateOptionals.Generate {
		attachCreateMappings(input.Name, input.Attributes, attributes)
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
