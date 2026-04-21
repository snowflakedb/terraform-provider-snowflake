package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (s *SessionPolicyModel) WithAllowedSecondaryRolesNone() *SessionPolicyModel {
	return s.WithAllowedSecondaryRolesValue(tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"none": tfconfig.BoolVariable(true),
	})))
}

func (s *SessionPolicyModel) WithAllowedSecondaryRolesAll() *SessionPolicyModel {
	return s.WithAllowedSecondaryRolesValue(tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"all": tfconfig.BoolVariable(true),
	})))
}

func (s *SessionPolicyModel) WithAllowedSecondaryRoles(roles ...string) *SessionPolicyModel {
	return s.WithAllowedSecondaryRolesValue(tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"roles": tfconfig.SetVariable(collections.Map(roles, func(role string) tfconfig.Variable {
			return tfconfig.StringVariable(role)
		})...),
	})))
}

func (s *SessionPolicyModel) WithBlockedSecondaryRolesNone() *SessionPolicyModel {
	return s.WithBlockedSecondaryRolesValue(tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"none": tfconfig.BoolVariable(true),
	})))
}

func (s *SessionPolicyModel) WithBlockedSecondaryRolesAll() *SessionPolicyModel {
	return s.WithBlockedSecondaryRolesValue(tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"all": tfconfig.BoolVariable(true),
	})))
}

func (s *SessionPolicyModel) WithBlockedSecondaryRoles(roles ...string) *SessionPolicyModel {
	return s.WithBlockedSecondaryRolesValue(tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"roles": tfconfig.SetVariable(collections.Map(roles, func(role string) tfconfig.Variable {
			return tfconfig.StringVariable(role)
		})...),
	})))
}
