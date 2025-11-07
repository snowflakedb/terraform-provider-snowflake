package testacc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/snowflakedb/gosnowflake"
)

// TODO [SNOW-2234579]: populate all the remaining fields of gosnowflake.Config
//   - validate_default_parameters
//   - params
//   - client_ip
//   - protocol
//   - host
//   - port
//   - okta_url
//   - login_timeout
//   - request_timeout
//   - jwt_expire_timeout
//   - client_timeout
//   - jwt_client_timeout
//   - external_browser_timeout
//   - insecure_mode
//   - ocsp_fail_open
//   - token
//   - keep_session_alive
//   - token_accessor
//   - disable_telemetry
//   - client_request_mfa_token
//   - client_store_temporary_credential
//   - disable_query_context_cache
//   - include_retry_reason
//   - max_retry_count
//   - tmp_directory_path
//   - disable_console_login
//   - DisableSamlURLCheck
func (p *pluginFrameworkPocProvider) getDriverConfigFromTerraform(configModel pluginFrameworkPocProviderModelV0) (*gosnowflake.Config, error) {
	config := sdk.EmptyDriverConfigWithApplication("terraform-provider-snowflake")
	if errs := errors.Join(
		setAccount(configModel, config),
		setStringAttribute(configModel.User, snowflakeenvs.User, &config.User),
		setStringAttribute(configModel.Password, snowflakeenvs.Password, &config.Password),
		setStringAttribute(configModel.Warehouse, snowflakeenvs.Warehouse, &config.Warehouse),
		setStringAttribute(configModel.Role, snowflakeenvs.Role, &config.Role),
		setEnumAttribute(configModel.Authenticator, snowflakeenvs.Authenticator, sdk.ToExtendedAuthenticatorType, &config.Authenticator),
		setStringAttribute(configModel.Passcode, snowflakeenvs.Passcode, &config.Passcode),
		setBoolAttribute(configModel.PasscodeInPassword, snowflakeenvs.PasscodeInPassword, &config.PasscodeInPassword),
		setEnumAttributeString(configModel.DriverTracing, snowflakeenvs.DriverTracing, sdk.ToDriverLogLevel, &config.Tracing),
		setPrivateKey(configModel, config),
		// profile is handled in the calling function
	); errs != nil {
		return nil, errs
	}

	return config, nil
}

func setAccount(configModel pluginFrameworkPocProviderModelV0, config *gosnowflake.Config) error {
	accountName := getStringAttribute(configModel.AccountName, snowflakeenvs.AccountName)
	organizationName := getStringAttribute(configModel.OrganizationName, snowflakeenvs.OrganizationName)
	if accountName != "" && organizationName != "" {
		config.Account = strings.Join([]string{organizationName, accountName}, "-")
	}
	return nil
}

func setPrivateKey(configModel pluginFrameworkPocProviderModelV0, config *gosnowflake.Config) error {
	privateKey := getStringAttribute(configModel.PrivateKey, snowflakeenvs.PrivateKey)
	privateKeyPassphrase := getStringAttribute(configModel.PrivateKeyPassphrase, snowflakeenvs.PrivateKeyPassphrase)
	v, err := provider.GetPrivateKey(privateKey, privateKeyPassphrase)
	if err != nil {
		return fmt.Errorf("could not retrieve private key: %w", err)
	}
	if v != nil {
		config.PrivateKey = v
	}
	return nil
}

func getProfile(configModel pluginFrameworkPocProviderModelV0) string {
	profile := getStringAttribute(configModel.Profile, snowflakeenvs.Profile)
	if profile == "" {
		// There are no default in plugin framework, so it needs to be assigned manually.
		// This is to achieve the same behavior as in the existing SDKv2 implementation.
		profile = "default"
	}
	return profile
}

func getStringAttribute(stringField types.String, envName string) string {
	var value string
	if !stringField.IsNull() {
		value = stringField.ValueString()
	} else {
		value = oswrapper.Getenv(envName)
	}
	return value
}

func setStringAttribute(stringField types.String, envName string, setInConfig *string) error {
	value := getStringAttribute(stringField, envName)
	if value != "" {
		*setInConfig = value
	}
	return nil
}

func setEnumAttribute[T any](stringField types.String, envName string, toEnumFunc func(string) (T, error), setInConfig *T) error {
	var value string
	if !stringField.IsNull() {
		value = stringField.ValueString()
	} else {
		value = oswrapper.Getenv(envName)
	}
	enumValue, err := toEnumFunc(value)
	if err != nil {
		return err
	}
	*setInConfig = enumValue
	return nil
}

func setEnumAttributeString[T ~string](stringField types.String, envName string, toEnumFunc func(string) (T, error), setInConfig *string) error {
	var value string
	if !stringField.IsNull() {
		value = stringField.ValueString()
	} else {
		value = oswrapper.Getenv(envName)
	}
	if value != "" {
		enumValue, err := toEnumFunc(value)
		if err != nil {
			return err
		}
		*setInConfig = string(enumValue)
	}
	return nil
}

func setBoolAttribute(boolField types.Bool, envName string, setInConfig *bool) error {
	var value bool
	if !boolField.IsNull() {
		value = boolField.ValueBool()
	} else {
		if v, err := oswrapper.GetenvBool(envName); err != nil {
			return err
		} else {
			value = v
		}
	}
	if !value {
		*setInConfig = value
	}
	return nil
}
