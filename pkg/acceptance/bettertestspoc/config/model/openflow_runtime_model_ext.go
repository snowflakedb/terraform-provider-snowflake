package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (o *OpenflowRuntimeModel) WithExternalAccessIntegrations(ids ...sdk.AccountObjectIdentifier) *OpenflowRuntimeModel {
	return o.WithExternalAccessIntegrationsValue(
		tfconfig.SetVariable(
			collections.Map(ids, func(id sdk.AccountObjectIdentifier) tfconfig.Variable {
				return tfconfig.StringVariable(id.Name())
			})...,
		),
	)
}
