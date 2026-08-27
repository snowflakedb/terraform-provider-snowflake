package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (c *CortexSearchServiceModel) WithAttributes(attrs ...string) *CortexSearchServiceModel {
	return c.WithAttributesValue(
		tfconfig.SetVariable(
			collections.Map(attrs, func(attr string) tfconfig.Variable { return tfconfig.StringVariable(attr) })...,
		),
	)
}
