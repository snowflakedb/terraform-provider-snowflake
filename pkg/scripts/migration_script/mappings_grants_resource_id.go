package main

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func CreateGrantPrivilegesToAccountRoleResourceIdOnAccount(grant sdk.Grant) string {
	return NormalizeResourceId(fmt.Sprintf("grant_on_account_to_%s_%s", grant.GranteeName.Name(), withGrantOptionString(grant)))
}

func CreateGrantPrivilegesToAccountRoleResourceIdOnAccountObject(grant sdk.Grant) string {
	return NormalizeResourceId(fmt.Sprintf("grant_on_%s_%s_to_%s_%s", grant.GrantedOn, grant.Name.Name(), grant.GranteeName.Name(), withGrantOptionString(grant)))
}

func CreateGrantPrivilegesToAccountRoleResourceIdOnSchema(grant sdk.Grant) string {
	return NormalizeResourceId(fmt.Sprintf("grant_on_schema_%s_to_%s_%s", grant.Name.FullyQualifiedName(), grant.GranteeName.Name(), withGrantOptionString(grant)))
}

func CreateGrantPrivilegesToAccountRoleResourceIdOnSchemaObject(grant sdk.Grant) string {
	return NormalizeResourceId(fmt.Sprintf("grant_on_%s_%s_to_%s_%s", grant.GrantedOn, grant.Name.FullyQualifiedName(), grant.GranteeName.Name(), withGrantOptionString(grant)))
}

func withGrantOptionString(grant sdk.Grant) string {
	if grant.GrantOption {
		return "with_grant_option"
	}
	return "without_grant_option"
}
