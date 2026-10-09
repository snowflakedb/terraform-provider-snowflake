package gen

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"

type ResourceModel struct {
	Name                                string
	DescribeFailRead                    bool
	DescribeUsesShowId                  bool
	HasPartialOnUpdate                  bool
	AdditionalDescribeOutputDescription string
	IdType                              string
	IdentityNameForceNew                bool
	AdditionalIdentityNameDescription   string
	Hooks                               ResourceHooks
	*genhelpers.PreambleModel
}

func (m ResourceModel) Prefix() string {
	return genhelpers.FirstLetterLowercase(m.Name)
}

func ModelFromInputObject(input ResourceDef, preamble *genhelpers.PreambleModel) ResourceModel {
	return ResourceModel{
		Name:                                input.name,
		DescribeFailRead:                    input.describeFailRead,
		DescribeUsesShowId:                  input.describeUsesShowId,
		HasPartialOnUpdate:                  input.hasPartialOnUpdate,
		AdditionalDescribeOutputDescription: input.additionalDescribeOutputDescription,
		IdType:                              input.idType,
		IdentityNameForceNew:                input.identityNameForceNew,
		AdditionalIdentityNameDescription:   input.additionalIdentityNameDescription,
		Hooks:                               defaultExtHooks(input.hooks...),
		PreambleModel:                       preamble,
	}
}
