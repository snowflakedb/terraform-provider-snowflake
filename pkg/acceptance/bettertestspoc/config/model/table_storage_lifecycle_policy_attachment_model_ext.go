package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (t *TableStorageLifecyclePolicyAttachmentModel) WithOn(on []string) *TableStorageLifecyclePolicyAttachmentModel {
	return t.WithOnValue(tfconfig.ListVariable(collections.Map(on, func(c string) tfconfig.Variable { return tfconfig.StringVariable(c) })...))
}
