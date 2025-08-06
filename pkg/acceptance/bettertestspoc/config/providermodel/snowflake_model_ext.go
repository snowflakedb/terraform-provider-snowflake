package providermodel

import (
	"encoding/json"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// Based on https://medium.com/picus-security-engineering/custom-json-marshaller-in-go-and-common-pitfalls-c43fa774db05.
func (m *SnowflakeModel) MarshalJSON() ([]byte, error) {
	type AliasModelType SnowflakeModel
	return json.Marshal(&struct {
		*AliasModelType
		Alias string `json:"alias,omitempty"`
	}{
		AliasModelType: (*AliasModelType)(m),
		Alias:          m.Alias(),
	})
}

func (m *SnowflakeModel) WithUserId(userId sdk.AccountObjectIdentifier) *SnowflakeModel {
	return m.WithUser(userId.Name())
}

func (m *SnowflakeModel) WithRoleId(roleId sdk.AccountObjectIdentifier) *SnowflakeModel {
	return m.WithRole(roleId.Name())
}

func (m *SnowflakeModel) WithWarehouseId(warehouseId sdk.AccountObjectIdentifier) *SnowflakeModel {
	return m.WithWarehouse(warehouseId.Name())
}

func (m *SnowflakeModel) WithAuthenticatorType(authenticationType sdk.AuthenticationType) *SnowflakeModel {
	return m.WithAuthenticator(string(authenticationType))
}

func (m *SnowflakeModel) WithPrivateKeyMultiline(privateKey string) *SnowflakeModel {
	return m.WithPrivateKeyValue(config.MultilineWrapperVariable(privateKey))
}

func (m *SnowflakeModel) WithValidateDefaultParametersBool(validateDefaultParameters bool) *SnowflakeModel {
	m.ValidateDefaultParameters = tfconfig.BoolVariable(validateDefaultParameters)
	return m
}

func (m *SnowflakeModel) WithClientRequestMfaTokenBool(clientRequestMfaToken bool) *SnowflakeModel {
	m.ClientRequestMfaToken = tfconfig.BoolVariable(clientRequestMfaToken)
	return m
}

func (m *SnowflakeModel) WithClientStoreTemporaryCredentialBool(clientStoreTemporaryCredential bool) *SnowflakeModel {
	m.ClientStoreTemporaryCredential = tfconfig.BoolVariable(clientStoreTemporaryCredential)
	return m
}

func (m *SnowflakeModel) WithIncludeRetryReasonBool(includeRetryReason bool) *SnowflakeModel {
	m.IncludeRetryReason = tfconfig.BoolVariable(includeRetryReason)
	return m
}

func (m *SnowflakeModel) WithDisableConsoleLoginBool(disableConsoleLogin bool) *SnowflakeModel {
	m.DisableConsoleLogin = tfconfig.BoolVariable(disableConsoleLogin)
	return m
}

func (m *SnowflakeModel) WithPreviewFeaturesEnabled(previewFeaturesEnabled ...string) *SnowflakeModel {
	previewFeaturesEnabledStringVariables := make([]tfconfig.Variable, len(previewFeaturesEnabled))
	for i, v := range previewFeaturesEnabled {
		previewFeaturesEnabledStringVariables[i] = tfconfig.StringVariable(v)
	}
	m.PreviewFeaturesEnabled = tfconfig.SetVariable(previewFeaturesEnabledStringVariables...)
	return m
}

func (m *SnowflakeModel) AllFields(tmpConfig *helpers.TmpTomlConfig, tmpUser *helpers.TmpServiceUser) *SnowflakeModel {
	return SnowflakeProvider().
		WithProfile(tmpConfig.Profile).
		WithOrganizationName(tmpUser.AccountId.OrganizationName()).
		WithAccountName(tmpUser.AccountId.AccountName()).
		WithUserId(tmpUser.UserId).
		WithPrivateKeyMultiline(tmpUser.PrivateKey).
		WithWarehouseId(tmpUser.WarehouseId).
		WithProtocol("https").
		WithPort(443).
		WithRoleId(tmpUser.RoleId).
		WithValidateDefaultParameters("true").
		WithClientIp("3.3.3.3").
		WithAuthenticatorType(sdk.AuthenticationTypeJwt).
		WithToken("correct token").
		WithOktaUrl(testvars.ExampleOktaUrlString).
		WithLoginTimeout(101).
		WithRequestTimeout(201).
		WithJwtExpireTimeout(301).
		WithClientTimeout(401).
		WithJwtClientTimeout(501).
		WithExternalBrowserTimeout(601).
		WithInsecureMode(true).
		WithOcspFailOpen("true").
		WithKeepSessionAlive(true).
		WithDisableTelemetry(true).
		WithClientRequestMfaToken("true").
		WithClientStoreTemporaryCredential("true").
		WithDisableQueryContextCache(true).
		WithIncludeRetryReason("true").
		WithMaxRetryCount(3).
		WithDriverTracing("warning").
		WithTmpDirectoryPath("../../").
		WithDisableConsoleLogin("true").
		WithParamsValue(
			tfconfig.ObjectVariable(
				map[string]tfconfig.Variable{
					"foo": tfconfig.StringVariable("piyo"),
				},
			),
		)
}

func PatConfig(h helpers.TmpServiceUserWithPat) *SnowflakeModel {
	return SnowflakeProvider().
		WithAuthenticatorType(sdk.AuthenticationTypeProgrammaticAccessToken).
		WithUserId(h.UserId).
		WithToken(h.Pat).
		WithRoleId(h.RoleId).
		WithOrganizationName(h.AccountId.OrganizationName()).
		WithAccountName(h.AccountId.AccountName()).
		WithWarehouseId(h.WarehouseId)
}
