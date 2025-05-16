package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
)

func (c *PrimaryConnectionModel) WithEnableFailover(toAccount ...sdk.AccountIdentifier) *PrimaryConnectionModel {
	variables := make([]config.Variable, 0)
	for _, v := range toAccount {
		variables = append(variables, config.StringVariable(v.Name()))
	}

	c.EnableFailoverToAccounts = config.ListVariable(variables...)

	return c
}
