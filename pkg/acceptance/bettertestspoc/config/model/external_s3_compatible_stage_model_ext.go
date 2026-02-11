package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func ExternalS3CompatibleStageWithId(id sdk.SchemaObjectIdentifier, url string, endpoint string) *ExternalS3CompatibleStageModel {
	return ExternalS3CompatibleStage("test", id.DatabaseName(), id.SchemaName(), id.Name(), endpoint, url)
}

func (e *ExternalS3CompatibleStageModel) WithDirectoryEnabled(enable string) *ExternalS3CompatibleStageModel {
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable": tfconfig.StringVariable(enable),
		},
	))
	return e
}

func (e *ExternalS3CompatibleStageModel) WithDirectoryEnabledAndOptions(opts sdk.StageS3CommonDirectoryTableOptionsRequest) *ExternalS3CompatibleStageModel {
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

func (e *ExternalS3CompatibleStageModel) WithInvalidAutoRefresh() *ExternalS3CompatibleStageModel {
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable":       tfconfig.BoolVariable(true),
			"auto_refresh": tfconfig.StringVariable("invalid"),
		},
	))
	return e
}

func (e *ExternalS3CompatibleStageModel) WithInvalidRefreshOnCreate() *ExternalS3CompatibleStageModel {
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable":            tfconfig.BoolVariable(true),
			"refresh_on_create": tfconfig.StringVariable("invalid"),
		},
	))
	return e
}

func (e *ExternalS3CompatibleStageModel) WithCredentials(awsKeyId string, awsSecretKey string) *ExternalS3CompatibleStageModel {
	e.Credentials = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"aws_key_id":     tfconfig.StringVariable(awsKeyId),
			"aws_secret_key": tfconfig.StringVariable(awsSecretKey),
		},
	))
	return e
}

func (e *ExternalS3CompatibleStageModel) WithEmptyCredentials() *ExternalS3CompatibleStageModel {
	e.Credentials = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		},
	))
	return e
}
