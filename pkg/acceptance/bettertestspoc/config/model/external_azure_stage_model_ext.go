package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func ExternalAzureStageWithId(id sdk.SchemaObjectIdentifier, url string) *ExternalAzureStageModel {
	return ExternalAzureStage("test", id.DatabaseName(), id.SchemaName(), id.Name(), url)
}

func (e *ExternalAzureStageModel) WithDirectoryEnabled(enable string) *ExternalAzureStageModel {
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable": tfconfig.StringVariable(enable),
		},
	))
	return e
}

func (e *ExternalAzureStageModel) WithDirectoryEnabledAndOptions(enable bool, refreshOnCreate string, autoRefresh *bool) *ExternalAzureStageModel {
	directoryMap := map[string]tfconfig.Variable{
		"enable": tfconfig.BoolVariable(enable),
	}
	if refreshOnCreate != "" {
		directoryMap["refresh_on_create"] = tfconfig.StringVariable(refreshOnCreate)
	}
	if autoRefresh != nil {
		directoryMap["auto_refresh"] = tfconfig.BoolVariable(*autoRefresh)
	}
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(directoryMap))
	return e
}

func (e *ExternalAzureStageModel) WithDirectoryEnabledAndNotificationIntegration(enable bool, notificationIntegration string) *ExternalAzureStageModel {
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable":                   tfconfig.BoolVariable(enable),
			"notification_integration": tfconfig.StringVariable(notificationIntegration),
		},
	))
	return e
}

func (e *ExternalAzureStageModel) WithEncryptionAzureCse(masterKey string) *ExternalAzureStageModel {
	return e.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"azure_cse": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"master_key": tfconfig.StringVariable(masterKey),
				})),
			},
		)),
	)
}

func (e *ExternalAzureStageModel) WithEncryptionNone() *ExternalAzureStageModel {
	return e.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"none": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
				})),
			},
		)),
	)
}

func (e *ExternalAzureStageModel) WithEncryptionBothTypes() *ExternalAzureStageModel {
	return e.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"azure_cse": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"master_key": tfconfig.StringVariable("foo"),
				})),
				"none": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
				})),
			},
		)),
	)
}

func (e *ExternalAzureStageModel) WithEncryptionNoneTypeSpecified() *ExternalAzureStageModel {
	return e.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
			},
		)),
	)
}

func (e *ExternalAzureStageModel) WithCredentials(azureSasToken string) *ExternalAzureStageModel {
	e.Credentials = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"azure_sas_token": tfconfig.StringVariable(azureSasToken),
		},
	))
	return e
}
