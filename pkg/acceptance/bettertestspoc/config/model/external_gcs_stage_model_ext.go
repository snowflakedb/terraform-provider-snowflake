package model

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func ExternalGcsStageWithId(id sdk.SchemaObjectIdentifier, storageIntegration, url string) *ExternalGcsStageModel {
	return ExternalGcsStage("test", id.DatabaseName(), id.SchemaName(), id.Name(), storageIntegration, url)
}

func (e *ExternalGcsStageModel) WithDirectoryEnabled(enable string) *ExternalGcsStageModel {
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable": tfconfig.StringVariable(enable),
		},
	))
	return e
}

func (e *ExternalGcsStageModel) WithDirectoryEnabledAndOptions(opts sdk.ExternalGCSDirectoryTableOptionsRequest) *ExternalGcsStageModel {
	directoryMap := map[string]tfconfig.Variable{
		"enable": tfconfig.BoolVariable(opts.Enable),
	}
	if opts.RefreshOnCreate != nil {
		directoryMap["refresh_on_create"] = tfconfig.StringVariable(strconv.FormatBool(*opts.RefreshOnCreate))
	}
	if opts.AutoRefresh != nil {
		directoryMap["auto_refresh"] = tfconfig.StringVariable(strconv.FormatBool(*opts.AutoRefresh))
	}
	if opts.NotificationIntegration != nil {
		directoryMap["notification_integration"] = tfconfig.StringVariable(*opts.NotificationIntegration)
	}
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(directoryMap))
	return e
}

func (e *ExternalGcsStageModel) WithInvalidAutoRefresh() *ExternalGcsStageModel {
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable":       tfconfig.BoolVariable(true),
			"auto_refresh": tfconfig.StringVariable("invalid"),
		},
	))
	return e
}

func (e *ExternalGcsStageModel) WithInvalidRefreshOnCreate() *ExternalGcsStageModel {
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable":            tfconfig.BoolVariable(true),
			"refresh_on_create": tfconfig.StringVariable("invalid"),
		},
	))
	return e
}

// WithFileFormatName sets a named file format reference.
func (e *ExternalGcsStageModel) WithFileFormatName(formatName string) *ExternalGcsStageModel {
	return e.WithFileFormatValue(stageFileFormatName(formatName))
}

// WithFileFormatCsv sets inline CSV file format with the provided options.
func (e *ExternalGcsStageModel) WithFileFormatCsv(opts sdk.FileFormatCsvOptions) *ExternalGcsStageModel {
	return e.WithFileFormatValue(stageFileFormatCsv(opts))
}

func (e *ExternalGcsStageModel) WithEncryptionGcsSseKms(kmsKeyId string) *ExternalGcsStageModel {
	encryptionMap := map[string]tfconfig.Variable{
		"gcs_sse_kms": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"kms_key_id": tfconfig.StringVariable(kmsKeyId),
		})),
	}
	return e.WithEncryptionValue(tfconfig.ListVariable(tfconfig.ObjectVariable(encryptionMap)))
}

func (e *ExternalGcsStageModel) WithEncryptionGcsSseKmsNoKey() *ExternalGcsStageModel {
	encryptionMap := map[string]tfconfig.Variable{
		"gcs_sse_kms": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		})),
	}
	return e.WithEncryptionValue(tfconfig.ListVariable(tfconfig.ObjectVariable(encryptionMap)))
}

func (e *ExternalGcsStageModel) WithEncryptionNone() *ExternalGcsStageModel {
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

func (e *ExternalGcsStageModel) WithEncryptionBothTypes() *ExternalGcsStageModel {
	return e.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"gcs_sse_kms": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"kms_key_id": tfconfig.StringVariable("foo"),
				})),
				"none": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
				})),
			},
		)),
	)
}

func (e *ExternalGcsStageModel) WithEncryptionNoneTypeSpecified() *ExternalGcsStageModel {
	return e.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
			},
		)),
	)
}

func (e *ExternalGcsStageModel) WithFileFormatAvro(opts sdk.FileFormatAvroOptions) *ExternalGcsStageModel {
	return e.WithFileFormatValue(stageFileFormatAvro(opts))
}
