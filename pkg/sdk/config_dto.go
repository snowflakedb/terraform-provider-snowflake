package sdk

import "github.com/pelletier/go-toml/v2"

//go:generate go run ./dto-builder-generator/main.go

type ConfigFile struct {
	Profiles map[string]ConfigDTO
}

func (c *ConfigFile) MarshalToml() ([]byte, error) {
	return toml.Marshal(c.Profiles)
}

// TODO(SNOW-1787920): improve TOML parsing
type ConfigDTO struct {
	AccountName            *string             `toml:"account_name"`
	OrganizationName       *string             `toml:"organization_name"`
	User                   *string             `toml:"user"`
	Username               *string             `toml:"username"`
	Password               *string             `toml:"password"`
	Host                   *string             `toml:"host"`
	Warehouse              *string             `toml:"warehouse"`
	Role                   *string             `toml:"role"`
	Params                 *map[string]*string `toml:"params"`
	ClientIp               *string             `toml:"client_ip"`
	Protocol               *string             `toml:"protocol"`
	Passcode               *string             `toml:"passcode"`
	Port                   *int                `toml:"port"`
	PasscodeInPassword     *bool               `toml:"passcode_in_password"`
	OktaUrl                *string             `toml:"okta_url"`
	ClientTimeout          *int                `toml:"client_timeout"`
	JwtClientTimeout       *int                `toml:"jwt_client_timeout"`
	LoginTimeout           *int                `toml:"login_timeout"`
	RequestTimeout         *int                `toml:"request_timeout"`
	JwtExpireTimeout       *int                `toml:"jwt_expire_timeout"`
	ExternalBrowserTimeout *int                `toml:"external_browser_timeout"`
	MaxRetryCount          *int                `toml:"max_retry_count"`
	Authenticator          *string             `toml:"authenticator"`
	InsecureMode           *bool               `toml:"insecure_mode"`
	OcspFailOpen           *bool               `toml:"ocsp_fail_open"`
	Token                  *string             `toml:"token"`
	KeepSessionAlive       *bool               `toml:"keep_session_alive"`
	PrivateKey             *string             `toml:"private_key,multiline"`
	PrivateKeyPassphrase   *string             `toml:"private_key_passphrase"`
	DisableTelemetry       *bool               `toml:"disable_telemetry"`
	// TODO [SNOW-1827312]: handle and test 3-value booleans properly from TOML
	ValidateDefaultParameters      *bool   `toml:"validate_default_parameters"`
	ClientRequestMfaToken          *bool   `toml:"client_request_mfa_token"`
	ClientStoreTemporaryCredential *bool   `toml:"client_store_temporary_credential"`
	DriverTracing                  *string `toml:"driver_tracing"`
	TmpDirPath                     *string `toml:"tmp_dir_path"`
	DisableQueryContextCache       *bool   `toml:"disable_query_context_cache"`
	IncludeRetryReason             *bool   `toml:"include_retry_reason"`
	DisableConsoleLogin            *bool   `toml:"disable_console_login"`
}
