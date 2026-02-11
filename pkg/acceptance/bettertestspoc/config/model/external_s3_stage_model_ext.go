package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func ExternalS3StageWithId(id sdk.SchemaObjectIdentifier, url string) *ExternalS3StageModel {
	return ExternalS3Stage("test", id.DatabaseName(), id.SchemaName(), id.Name(), url)
}

func (e *ExternalS3StageModel) WithDirectoryEnabled(enable string) *ExternalS3StageModel {
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable": tfconfig.StringVariable(enable),
		},
	))
	return e
}

func (e *ExternalS3StageModel) WithDirectoryEnabledAndOptions(opts sdk.StageS3CommonDirectoryTableOptionsRequest) *ExternalS3StageModel {
	directoryMap := map[string]tfconfig.Variable{
		"enable": tfconfig.BoolVariable(opts.Enable),
	}
	if opts.RefreshOnCreate != nil {
		directoryMap["refresh_on_create"] = tfconfig.BoolVariable(*opts.RefreshOnCreate)
	}
	if opts.AutoRefresh != nil {
		directoryMap["auto_refresh"] = tfconfig.BoolVariable(*opts.AutoRefresh)
	}
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(directoryMap))
	return e
}

func (e *ExternalS3StageModel) WithInvalidAutoRefresh() *ExternalS3StageModel {
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable":       tfconfig.BoolVariable(true),
			"auto_refresh": tfconfig.StringVariable("invalid"),
		},
	))
	return e
}

func (e *ExternalS3StageModel) WithInvalidRefreshOnCreate() *ExternalS3StageModel {
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable":            tfconfig.BoolVariable(true),
			"refresh_on_create": tfconfig.StringVariable("invalid"),
		},
	))
	return e
}

func (e *ExternalS3StageModel) WithEncryptionAwsCse(masterKey string) *ExternalS3StageModel {
	return e.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"aws_cse": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"master_key": tfconfig.StringVariable(masterKey),
				})),
			},
		)),
	)
}

func (e *ExternalS3StageModel) WithEncryptionAwsSseS3() *ExternalS3StageModel {
	return e.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"aws_sse_s3": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
				})),
			},
		)),
	)
}

func (e *ExternalS3StageModel) WithEncryptionAwsSseKms(kmsKeyId string) *ExternalS3StageModel {
	return e.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"aws_sse_kms": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"kms_key_id": tfconfig.StringVariable(kmsKeyId),
				})),
			},
		)),
	)
}

func (e *ExternalS3StageModel) WithEncryptionAwsSseKmsEmpty() *ExternalS3StageModel {
	return e.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"aws_sse_kms": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
				})),
			},
		)),
	)
}

func (e *ExternalS3StageModel) WithEncryptionNone() *ExternalS3StageModel {
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

func (e *ExternalS3StageModel) WithEncryptionBothTypes() *ExternalS3StageModel {
	return e.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"aws_cse": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"master_key": tfconfig.StringVariable("foo"),
				})),
				"none": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
				})),
			},
		)),
	)
}

func (e *ExternalS3StageModel) WithEncryptionNoneTypeSpecified() *ExternalS3StageModel {
	return e.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
			},
		)),
	)
}

// WithFileFormatName sets a named file format reference.
func (e *ExternalS3StageModel) WithFileFormatName(formatName string) *ExternalS3StageModel {
	return e.WithFileFormatValue(stageFileFormatName(formatName))
}

// WithFileFormatCsv sets inline CSV file format with the provided options.
func (e *ExternalS3StageModel) WithFileFormatCsv(opts sdk.FileFormatCsvOptions) *ExternalS3StageModel {
	return e.WithFileFormatValue(stageFileFormatCsv(opts))
}

func (e *ExternalS3StageModel) WithCredentialsAwsKey(keyId, secretKey string) *ExternalS3StageModel {
	e.Credentials = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"aws_key_id":     tfconfig.StringVariable(keyId),
			"aws_secret_key": tfconfig.StringVariable(secretKey),
		},
	))
	return e
}

func (e *ExternalS3StageModel) WithCredentialsAwsKeyWithToken(keyId, secretKey, token string) *ExternalS3StageModel {
	e.Credentials = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"aws_key_id":     tfconfig.StringVariable(keyId),
			"aws_secret_key": tfconfig.StringVariable(secretKey),
			"aws_token":      tfconfig.StringVariable(token),
		},
	))
	return e
}

func (e *ExternalS3StageModel) WithCredentialsAwsRole(roleArn string) *ExternalS3StageModel {
	e.Credentials = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"aws_role": tfconfig.StringVariable(roleArn),
		},
	))
	return e
}
