package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (n *NetworkPolicyAttachmentModel) WithUsers(users ...string) *NetworkPolicyAttachmentModel {
	return n.WithUsersValue(
		tfconfig.SetVariable(
			collections.Map(users, func(u string) tfconfig.Variable { return tfconfig.StringVariable(u) })...,
		),
	)
}
