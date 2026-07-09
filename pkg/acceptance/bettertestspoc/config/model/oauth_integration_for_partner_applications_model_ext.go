package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (t *OauthIntegrationForPartnerApplicationsModel) WithBlockedRolesList(blockedRoles ...string) *OauthIntegrationForPartnerApplicationsModel {
	blockedRolesListStringVariables := make([]tfconfig.Variable, len(blockedRoles))
	for i, v := range blockedRoles {
		blockedRolesListStringVariables[i] = tfconfig.StringVariable(v)
	}

	t.BlockedRolesList = tfconfig.SetVariable(blockedRolesListStringVariables...)
	return t
}

func (t *OauthIntegrationForPartnerApplicationsModel) WithAllowedRoles(roles ...sdk.AccountObjectIdentifier) *OauthIntegrationForPartnerApplicationsModel {
	t.AllowedRolesList = tfconfig.SetVariable(
		collections.Map(roles, func(role sdk.AccountObjectIdentifier) tfconfig.Variable {
			return tfconfig.StringVariable(role.Name())
		})...,
	)
	return t
}

func (t *OauthIntegrationForPartnerApplicationsModel) WithAllowedRolesEmpty() *OauthIntegrationForPartnerApplicationsModel {
	t.AllowedRolesList = config.EmptyListVariable()
	return t
}
