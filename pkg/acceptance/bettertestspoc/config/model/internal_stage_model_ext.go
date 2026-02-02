package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func InternalStageWithId(id sdk.SchemaObjectIdentifier) *InternalStageModel {
	return InternalStage("test", id.DatabaseName(), id.SchemaName(), id.Name())
}

// WithDirectoryEnabled adds directory block with enable only
func (i *InternalStageModel) WithDirectoryEnabled(enable string) *InternalStageModel {
	i.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable": tfconfig.StringVariable(enable),
		},
	))
	return i
}

// WithDirectoryEnabledAndAutoRefresh adds directory block with enable and auto_refresh
func (i *InternalStageModel) WithDirectoryEnabledAndAutoRefresh(enable bool, autoRefresh string) *InternalStageModel {
	i.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable":       tfconfig.BoolVariable(enable),
			"auto_refresh": tfconfig.StringVariable(autoRefresh),
		},
	))
	return i
}

// WithEncryptionSnowflakeFull sets encryption to snowflake_full
func (i *InternalStageModel) WithEncryptionSnowflakeFull() *InternalStageModel {
	return i.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"snowflake_full": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
				})),
			},
		)),
	)
}

// WithEncryptionSnowflakeSse sets encryption to snowflake_sse
func (i *InternalStageModel) WithEncryptionSnowflakeSse() *InternalStageModel {
	return i.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"snowflake_sse": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
				})),
			},
		)),
	)
}

// WithEncryptionBothTypes sets both encryption types to trigger ExactlyOneOf validation error
func (i *InternalStageModel) WithEncryptionBothTypes() *InternalStageModel {
	return i.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"snowflake_full": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
				})),
				"snowflake_sse": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
				})),
			},
		)),
	)
}
