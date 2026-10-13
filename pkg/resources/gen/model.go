package gen

import (
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

func ModelFromInputObject(input ResourceObject, preamble *genhelpers.PreambleModel) ResourceModel {
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
		Attributes:                          applyResourceAttributeRules(input.Attributes),
		AttributesModification:              input.AttributesModification,
		Hooks:                               ssot.DefaultResourceExtHooks(input.ResourceHooks...),
		PreambleModel:                       preamble,
	}
}
