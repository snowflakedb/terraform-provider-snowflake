package testacc

import (
	"crypto/rsa"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/snowflakedb/gosnowflake"
)

// TODO [mux-PR]: populate all the remaining fields of gosnowflake.Config
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
func (p *pluginFrameworkPocProvider) getDriverConfigFromTerraform(m pluginFrameworkPocProviderModelV0) (*gosnowflake.Config, error) {
	config := &gosnowflake.Config{
		Application: "terraform-provider-snowflake",
	}

	// TODO [mux-PR]: use os wrapper
	// todoFromEnv := os.Getenv("SNOWFLAKE_TODO")

	// account_name and organization_name
	// user
	// password
	// warehouse
	// role
	// authenticator
	// passcode
	// passcode_in_password
	// private_key and private_key_passphrase
	// driver_tracing
	// profile (is handled in the calling function)

	// account_name and organization_name override legacy account field
	//accountName := s.Get("account_name").(string)
	//organizationName := s.Get("organization_name").(string)
	//if accountName != "" && organizationName != "" {
	//	config.Account = strings.Join([]string{organizationName, accountName}, "-")
	//}

	// private_key and private_key_passphrase have additional logic
	//privateKey := s.Get("private_key").(string)
	//privateKeyPassphrase := s.Get("private_key_passphrase").(string)
	//v, err := getPrivateKey(privateKey, privateKeyPassphrase)
	//if err != nil {
	//	return nil, fmt.Errorf("could not retrieve private key: %w", err)
	//}
	//if v != nil {
	//	config.PrivateKey = v
	//}

	return config, nil
}

// this method was copied from SDKv2 provider_helpers.go
func getPrivateKey(privateKeyString, privateKeyPassphrase string) (*rsa.PrivateKey, error) {
	if privateKeyString == "" {
		return nil, nil
	}
	privateKeyBytes := []byte(privateKeyString)
	return sdk.ParsePrivateKey(privateKeyBytes, []byte(privateKeyPassphrase))
}
