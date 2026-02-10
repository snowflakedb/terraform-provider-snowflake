package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func ExternalS3CompatStageWithId(id sdk.SchemaObjectIdentifier, url string, endpoint string) *ExternalS3CompatStageModel {
	return ExternalS3CompatStage("test", id.DatabaseName(), id.SchemaName(), id.Name(), endpoint, url)
}

func (e *ExternalS3CompatStageModel) WithDirectoryEnabled(enable string) *ExternalS3CompatStageModel {
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable": tfconfig.StringVariable(enable),
		},
	))
	return e
}

func (e *ExternalS3CompatStageModel) WithDirectoryEnabledAndOptions(opts sdk.StageS3CommonDirectoryTableOptionsRequest) *ExternalS3CompatStageModel {
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

func (e *ExternalS3CompatStageModel) WithInvalidAutoRefresh() *ExternalS3CompatStageModel {
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable":       tfconfig.BoolVariable(true),
			"auto_refresh": tfconfig.StringVariable("invalid"),
		},
	))
	return e
}

func (e *ExternalS3CompatStageModel) WithInvalidRefreshOnCreate() *ExternalS3CompatStageModel {
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable":            tfconfig.BoolVariable(true),
			"refresh_on_create": tfconfig.StringVariable("invalid"),
		},
	))
	return e
}

func (e *ExternalS3CompatStageModel) WithCredentials(awsKeyId string, awsSecretKey string) *ExternalS3CompatStageModel {
	e.Credentials = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"aws_key_id":     tfconfig.StringVariable(awsKeyId),
			"aws_secret_key": tfconfig.StringVariable(awsSecretKey),
		},
	))
	return e
}

func (e *ExternalS3CompatStageModel) WithEmptyCredentials() *ExternalS3CompatStageModel {
	e.Credentials = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		},
	))
	return e
}
