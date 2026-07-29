package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (t *OauthIntegrationForPartnerApplicationsModel) WithBlockedRolesList(blockedRoles ...string) *OauthIntegrationForPartnerApplicationsModel {
	blockedRolesListStringVariables := make([]tfconfig.Variable, len(blockedRoles))
	for i, v := range blockedRoles {
		blockedRolesListStringVariables[i] = tfconfig.StringVariable(v)
	}

	t.BlockedRolesList = tfconfig.SetVariable(blockedRolesListStringVariables...)
	return t
}
