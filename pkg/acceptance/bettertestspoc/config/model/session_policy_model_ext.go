package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (s *SessionPolicyModel) WithAllowedSecondaryRoles(roles ...string) *SessionPolicyModel {
	return s.WithAllowedSecondaryRolesValue(
		tfconfig.SetVariable(
			collections.Map(roles, func(role string) tfconfig.Variable {
				return tfconfig.StringVariable(role)
			})...,
		))
}

func (s *SessionPolicyModel) WithBlockedSecondaryRoles(roles ...string) *SessionPolicyModel {
	return s.WithBlockedSecondaryRolesValue(
		tfconfig.SetVariable(
			collections.Map(roles, func(role string) tfconfig.Variable {
				return tfconfig.StringVariable(role)
			})...,
		))
}
