package sdk

import (
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/pelletier/go-toml/v2"
	"github.com/snowflakedb/gosnowflake"
	"github.com/youmark/pkcs8"
	"golang.org/x/crypto/ssh"
)

// ConfigProvider is an interface that allows us to use the same code to parse both the new and legacy config formats.
type ConfigProvider interface {
	*ConfigDTO | *LegacyConfigDTO
	DriverConfig() (gosnowflake.Config, error)
}

// FileReaderConfig is a struct that holds the configuration for the file reader.
type FileReaderConfig struct {
	verifyPermissions   bool
	useLegacyTomlFormat bool
}

func WithVerifyPermissions(verifyPermissions bool) func(*FileReaderConfig) {
	return func(c *FileReaderConfig) {
		c.verifyPermissions = verifyPermissions
	}
}

func WithUseLegacyTomlFormat(useLegacyTomlFormat bool) func(*FileReaderConfig) {
	return func(c *FileReaderConfig) {
		c.useLegacyTomlFormat = useLegacyTomlFormat
	}
}

func DefaultConfig(opts ...func(*FileReaderConfig)) *gosnowflake.Config {
	config, err := ProfileConfig("default", opts...)
	if err != nil || config == nil {
		log.Printf("[DEBUG] No Snowflake config file found, proceeding with empty config, err = %v", err)
		config = EmptyDriverConfig()
	}
	return config
}

func ProfileConfig(profile string, opts ...func(*FileReaderConfig)) (*gosnowflake.Config, error) {
	cfg := FileReaderConfig{
		verifyPermissions:   true,
		useLegacyTomlFormat: false,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	path, err := GetConfigFileName()
	if err != nil {
		return nil, err
	}

	if profile == "" {
		profile = "default"
	}
	log.Printf("[DEBUG] Retrieving %s profile from a TOML file", profile)
	var config *gosnowflake.Config
	if cfg.useLegacyTomlFormat {
		config, err = LoadProfileFromFile[*LegacyConfigDTO](profile, path, cfg.verifyPermissions)
	} else {
		config, err = LoadProfileFromFile[*ConfigDTO](profile, path, cfg.verifyPermissions)
	}
	if err != nil {
		return nil, err
	}

	if config == nil {
		log.Printf("[DEBUG] No config found for profile: \"%s\"", profile)
		return nil, nil
	}

	// us-west-2 is Snowflake's default region, but if you actually specify that it won't trigger the default code
	//  https://github.com/snowflakedb/gosnowflake/blob/52137ce8c32eaf93b0bd22fc5c7297beff339812/dsn.go#L61
	if config.Region == "us-west-2" {
		config.Region = ""
	}

	return config, nil
}

func (c *ConfigDTO) DriverConfig() (gosnowflake.Config, error) {
	driverCfg := EmptyDriverConfig()
	if c.AccountName != nil && c.OrganizationName != nil {
		driverCfg.Account = fmt.Sprintf("%s-%s", *c.OrganizationName, *c.AccountName)
	}
	pointerAttributeSet(c.User, &driverCfg.User)
	pointerAttributeSet(c.Username, &driverCfg.User)
	pointerAttributeSet(c.Password, &driverCfg.Password)
	pointerAttributeSet(c.Host, &driverCfg.Host)
	pointerAttributeSet(c.Warehouse, &driverCfg.Warehouse)
	pointerAttributeSet(c.Role, &driverCfg.Role)
	pointerAttributeSet(c.Params, &driverCfg.Params)
	pointerIpAttributeSet(c.ClientIp, &driverCfg.ClientIP)
	pointerAttributeSet(c.Protocol, &driverCfg.Protocol)
	pointerAttributeSet(c.Passcode, &driverCfg.Passcode)
	pointerAttributeSet(c.Port, &driverCfg.Port)
	pointerAttributeSet(c.PasscodeInPassword, &driverCfg.PasscodeInPassword)
	err := pointerUrlAttributeSet(c.OktaUrl, &driverCfg.OktaURL)
	if err != nil {
		return *EmptyDriverConfig(), err
	}
	pointerTimeInSecondsAttributeSet(c.ClientTimeout, &driverCfg.ClientTimeout)
	pointerTimeInSecondsAttributeSet(c.JwtClientTimeout, &driverCfg.JWTClientTimeout)
	pointerTimeInSecondsAttributeSet(c.LoginTimeout, &driverCfg.LoginTimeout)
	pointerTimeInSecondsAttributeSet(c.RequestTimeout, &driverCfg.RequestTimeout)
	pointerTimeInSecondsAttributeSet(c.JwtExpireTimeout, &driverCfg.JWTExpireTimeout)
	pointerTimeInSecondsAttributeSet(c.ExternalBrowserTimeout, &driverCfg.ExternalBrowserTimeout)
	pointerAttributeSet(c.MaxRetryCount, &driverCfg.MaxRetryCount)
	if c.Authenticator != nil {
		authenticator, err := ToAuthenticatorType(*c.Authenticator)
		if err != nil {
			return *EmptyDriverConfig(), err
		}
		driverCfg.Authenticator = authenticator
	}
	pointerAttributeSet(c.InsecureMode, &driverCfg.InsecureMode) //nolint:staticcheck
	if c.OcspFailOpen != nil {
		if *c.OcspFailOpen {
			driverCfg.OCSPFailOpen = gosnowflake.OCSPFailOpenTrue
		} else {
			driverCfg.OCSPFailOpen = gosnowflake.OCSPFailOpenFalse
		}
	}
	pointerAttributeSet(c.Token, &driverCfg.Token)
	pointerAttributeSet(c.KeepSessionAlive, &driverCfg.KeepSessionAlive)
	if c.PrivateKey != nil {
		passphrase := make([]byte, 0)
		if c.PrivateKeyPassphrase != nil {
			passphrase = []byte(*c.PrivateKeyPassphrase)
		}
		privKey, err := ParsePrivateKey([]byte(*c.PrivateKey), passphrase)
		if err != nil {
			return *EmptyDriverConfig(), err
		}
		driverCfg.PrivateKey = privKey
	}
	pointerAttributeSet(c.DisableTelemetry, &driverCfg.DisableTelemetry)
	pointerConfigBoolAttributeSet(c.ValidateDefaultParameters, &driverCfg.ValidateDefaultParameters)
	pointerConfigBoolAttributeSet(c.ClientRequestMfaToken, &driverCfg.ClientRequestMfaToken)
	pointerConfigBoolAttributeSet(c.ClientStoreTemporaryCredential, &driverCfg.ClientStoreTemporaryCredential)
	pointerAttributeSet(c.DriverTracing, &driverCfg.Tracing)
	pointerAttributeSet(c.TmpDirPath, &driverCfg.TmpDirPath)
	pointerAttributeSet(c.DisableQueryContextCache, &driverCfg.DisableQueryContextCache)
	pointerConfigBoolAttributeSet(c.IncludeRetryReason, &driverCfg.IncludeRetryReason)
	pointerConfigBoolAttributeSet(c.DisableConsoleLogin, &driverCfg.DisableConsoleLogin)
	pointerAttributeSet(c.OauthClientID, &driverCfg.OauthClientID)
	pointerAttributeSet(c.OauthClientSecret, &driverCfg.OauthClientSecret)
	pointerAttributeSet(c.OauthTokenRequestURL, &driverCfg.OauthTokenRequestURL)
	pointerAttributeSet(c.OauthAuthorizationURL, &driverCfg.OauthAuthorizationURL)
	pointerAttributeSet(c.OauthRedirectURI, &driverCfg.OauthRedirectURI)
	pointerAttributeSet(c.OauthScope, &driverCfg.OauthScope)
	pointerAttributeSet(c.EnableSingleUseRefreshTokens, &driverCfg.EnableSingleUseRefreshTokens)

	return *driverCfg, nil
}

func MergeConfig(baseConfig *gosnowflake.Config, mergeConfig *gosnowflake.Config) *gosnowflake.Config {
	if baseConfig == nil {
		return mergeConfig
	}
	if baseConfig.Account == "" {
		baseConfig.Account = mergeConfig.Account
	}
	if baseConfig.User == "" {
		baseConfig.User = mergeConfig.User
	}
	if baseConfig.Password == "" {
		baseConfig.Password = mergeConfig.Password
	}
	if baseConfig.Warehouse == "" {
		baseConfig.Warehouse = mergeConfig.Warehouse
	}
	if baseConfig.Role == "" {
		baseConfig.Role = mergeConfig.Role
	}
	if baseConfig.Region == "" {
		baseConfig.Region = mergeConfig.Region
	}
	if baseConfig.Host == "" {
		baseConfig.Host = mergeConfig.Host
	}
	if !configBoolSet(baseConfig.ValidateDefaultParameters) {
		baseConfig.ValidateDefaultParameters = mergeConfig.ValidateDefaultParameters
	}
	if mergedMap := collections.MergeMaps(mergeConfig.Params, baseConfig.Params); len(mergedMap) > 0 {
		baseConfig.Params = mergedMap
	}
	if baseConfig.ClientIP == nil {
		baseConfig.ClientIP = mergeConfig.ClientIP
	}
	if baseConfig.Protocol == "" {
		baseConfig.Protocol = mergeConfig.Protocol
	}
	if baseConfig.Host == "" {
		baseConfig.Host = mergeConfig.Host
	}
	if baseConfig.Port == 0 {
		baseConfig.Port = mergeConfig.Port
	}
	if baseConfig.Authenticator == GosnowflakeAuthTypeEmpty {
		baseConfig.Authenticator = mergeConfig.Authenticator
	}
	if baseConfig.Passcode == "" {
		baseConfig.Passcode = mergeConfig.Passcode
	}
	if !baseConfig.PasscodeInPassword {
		baseConfig.PasscodeInPassword = mergeConfig.PasscodeInPassword
	}
	if baseConfig.OktaURL == nil {
		baseConfig.OktaURL = mergeConfig.OktaURL
	}
	if baseConfig.LoginTimeout == 0 {
		baseConfig.LoginTimeout = mergeConfig.LoginTimeout
	}
	if baseConfig.RequestTimeout == 0 {
		baseConfig.RequestTimeout = mergeConfig.RequestTimeout
	}
	if baseConfig.JWTExpireTimeout == 0 {
		baseConfig.JWTExpireTimeout = mergeConfig.JWTExpireTimeout
	}
	if baseConfig.ClientTimeout == 0 {
		baseConfig.ClientTimeout = mergeConfig.ClientTimeout
	}
	if baseConfig.JWTClientTimeout == 0 {
		baseConfig.JWTClientTimeout = mergeConfig.JWTClientTimeout
	}
	if baseConfig.ExternalBrowserTimeout == 0 {
		baseConfig.ExternalBrowserTimeout = mergeConfig.ExternalBrowserTimeout
	}
	if baseConfig.MaxRetryCount == 0 {
		baseConfig.MaxRetryCount = mergeConfig.MaxRetryCount
	}
	if !baseConfig.InsecureMode { //nolint:staticcheck
		baseConfig.InsecureMode = mergeConfig.InsecureMode //nolint:staticcheck
	}
	if baseConfig.OCSPFailOpen == 0 {
		baseConfig.OCSPFailOpen = mergeConfig.OCSPFailOpen
	}
	if baseConfig.Token == "" {
		baseConfig.Token = mergeConfig.Token
	}
	if !baseConfig.KeepSessionAlive {
		baseConfig.KeepSessionAlive = mergeConfig.KeepSessionAlive
	}
	if baseConfig.PrivateKey == nil {
		baseConfig.PrivateKey = mergeConfig.PrivateKey
	}
	if !baseConfig.DisableTelemetry {
		baseConfig.DisableTelemetry = mergeConfig.DisableTelemetry
	}
	if baseConfig.Tracing == "" {
		baseConfig.Tracing = mergeConfig.Tracing
	}
	if baseConfig.TmpDirPath == "" {
		baseConfig.TmpDirPath = mergeConfig.TmpDirPath
	}
	if !configBoolSet(baseConfig.ClientRequestMfaToken) {
		baseConfig.ClientRequestMfaToken = mergeConfig.ClientRequestMfaToken
	}
	if !configBoolSet(baseConfig.ClientStoreTemporaryCredential) {
		baseConfig.ClientStoreTemporaryCredential = mergeConfig.ClientStoreTemporaryCredential
	}
	if !baseConfig.DisableQueryContextCache {
		baseConfig.DisableQueryContextCache = mergeConfig.DisableQueryContextCache
	}
	if !configBoolSet(baseConfig.IncludeRetryReason) {
		baseConfig.IncludeRetryReason = mergeConfig.IncludeRetryReason
	}
	if !configBoolSet(baseConfig.DisableConsoleLogin) {
		baseConfig.DisableConsoleLogin = mergeConfig.DisableConsoleLogin
	}
	if baseConfig.OauthClientID == "" {
		baseConfig.OauthClientID = mergeConfig.OauthClientID
	}
	if baseConfig.OauthClientSecret == "" {
		baseConfig.OauthClientSecret = mergeConfig.OauthClientSecret
	}
	if baseConfig.OauthAuthorizationURL == "" {
		baseConfig.OauthAuthorizationURL = mergeConfig.OauthAuthorizationURL
	}
	if baseConfig.OauthTokenRequestURL == "" {
		baseConfig.OauthTokenRequestURL = mergeConfig.OauthTokenRequestURL
	}
	if baseConfig.OauthRedirectURI == "" {
		baseConfig.OauthRedirectURI = mergeConfig.OauthRedirectURI
	}
	if baseConfig.OauthScope == "" {
		baseConfig.OauthScope = mergeConfig.OauthScope
	}
	if !baseConfig.EnableSingleUseRefreshTokens {
		baseConfig.EnableSingleUseRefreshTokens = mergeConfig.EnableSingleUseRefreshTokens
	}
	return baseConfig
}

func configBoolSet(v gosnowflake.ConfigBool) bool {
	// configBoolNotSet is unexported in the driver, so we check if it's neither true nor false
	return slices.Contains([]gosnowflake.ConfigBool{gosnowflake.ConfigBoolFalse, gosnowflake.ConfigBoolTrue}, v)
}

func boolToConfigBool(v bool) gosnowflake.ConfigBool {
	if v {
		return gosnowflake.ConfigBoolTrue
	}
	return gosnowflake.ConfigBoolFalse
}

func GetConfigFileName() (string, error) {
	// has the user overridden the default config path?
	if configPath, ok := oswrapper.LookupEnv("SNOWFLAKE_CONFIG_PATH"); ok {
		if configPath != "" {
			return configPath, nil
		}
	}
	dir, err := oswrapper.UserHomeDir()
	if err != nil {
		return "", err
	}
	// default config path is ~/.snowflake/config.
	return filepath.Join(dir, ".snowflake", "config"), nil
}

func pointerAttributeSet[T any](src, dst *T) {
	if src != nil {
		*dst = *src
	}
}

func pointerTimeInSecondsAttributeSet(src *int, dst *time.Duration) {
	if src != nil {
		*dst = time.Second * time.Duration(*src)
	}
}

func pointerConfigBoolAttributeSet(src *bool, dst *gosnowflake.ConfigBool) {
	if src != nil {
		*dst = boolToConfigBool(*src)
	}
}

func pointerIpAttributeSet(src *string, dst *net.IP) {
	if src != nil {
		*dst = net.ParseIP(*src)
	}
}

func pointerUrlAttributeSet(src *string, dst **url.URL) error {
	if src != nil {
		url, err := url.Parse(*src)
		if err != nil {
			return err
		}
		*dst = url
	}
	return nil
}

// LoadProfileFromFile loads a config file from the path and returns a ready configuration.
func LoadProfileFromFile[T ConfigProvider](profile string, path string, verifyPermissions bool) (*gosnowflake.Config, error) {
	configs, err := LoadConfigFile[T](path, verifyPermissions)
	if err != nil {
		return nil, fmt.Errorf("could not load config file: %w", err)
	}
	if cfg, ok := configs[profile]; ok {
		log.Printf("[DEBUG] Loading config for profile: \"%s\"", profile)
		driverCfg, err := cfg.DriverConfig()
		if err != nil {
			return nil, fmt.Errorf("converting profile \"%s\" in file %s failed: %w", profile, path, err)
		}
		return &driverCfg, nil
	}
	return nil, nil
}

// LoadConfigFile loads a config file from the path and returns a map of profiles to one of the TOML formats.
func LoadConfigFile[T ConfigProvider](path string, verifyPermissions bool) (map[string]T, error) {
	data, err := oswrapper.ReadFileSafe(path, verifyPermissions)
	if err != nil {
		return nil, err
	}
	var s map[string]T

	err = toml.Unmarshal(data, &s)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling config file %s: %w", path, err)
	}
	return s, nil
}

func ParsePrivateKey(privateKeyBytes []byte, passphrase []byte) (*rsa.PrivateKey, error) {
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	if privateKeyBlock == nil {
		return nil, fmt.Errorf("could not parse private key, key is not in PEM format")
	}

	if privateKeyBlock.Type == "ENCRYPTED PRIVATE KEY" {
		if len(passphrase) == 0 {
			return nil, fmt.Errorf("private key requires a passphrase, but private_key_passphrase was not supplied")
		}
		privateKey, err := pkcs8.ParsePKCS8PrivateKeyRSA(privateKeyBlock.Bytes, passphrase)
		if err != nil {
			return nil, fmt.Errorf("could not parse encrypted private key with passphrase, only ciphers aes-128-cbc, aes-128-gcm, aes-192-cbc, aes-192-gcm, aes-256-cbc, aes-256-gcm, and des-ede3-cbc are supported err = %w", err)
		}
		return privateKey, nil
	}

	// TODO(SNOW-1754327): check if we can simply use ssh.ParseRawPrivateKeyWithPassphrase
	privateKey, err := ssh.ParseRawPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key err = %w", err)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("privateKey not of type RSA")
	}
	return rsaPrivateKey, nil
}

type AuthenticationType string

const (
	AuthenticationTypeSnowflake               AuthenticationType = "SNOWFLAKE"
	AuthenticationTypeOauth                   AuthenticationType = "OAUTH"
	AuthenticationTypeExternalBrowser         AuthenticationType = "EXTERNALBROWSER"
	AuthenticationTypeOkta                    AuthenticationType = "OKTA"
	AuthenticationTypeJwt                     AuthenticationType = "SNOWFLAKE_JWT"
	AuthenticationTypeTokenAccessor           AuthenticationType = "TOKENACCESSOR"
	AuthenticationTypeUsernamePasswordMfa     AuthenticationType = "USERNAMEPASSWORDMFA"
	AuthenticationTypeProgrammaticAccessToken AuthenticationType = "PROGRAMMATIC_ACCESS_TOKEN" //nolint:gosec
	AuthenticationTypeOauthClientCredentials  AuthenticationType = "OAUTH_CLIENT_CREDENTIALS"  //nolint:gosec
	AuthenticationTypeOauthAuthorizationCode  AuthenticationType = "OAUTH_AUTHORIZATION_CODE"

	AuthenticationTypeEmpty AuthenticationType = ""
)

var AllAuthenticationTypes = []AuthenticationType{
	AuthenticationTypeSnowflake,
	AuthenticationTypeOauth,
	AuthenticationTypeExternalBrowser,
	AuthenticationTypeOkta,
	AuthenticationTypeJwt,
	AuthenticationTypeTokenAccessor,
	AuthenticationTypeUsernamePasswordMfa,
	AuthenticationTypeProgrammaticAccessToken,
	AuthenticationTypeOauthClientCredentials,
	AuthenticationTypeOauthAuthorizationCode,
}

func ToAuthenticatorType(s string) (gosnowflake.AuthType, error) {
	switch strings.ToUpper(s) {
	case string(AuthenticationTypeSnowflake):
		return gosnowflake.AuthTypeSnowflake, nil
	case string(AuthenticationTypeOauth):
		return gosnowflake.AuthTypeOAuth, nil
	case string(AuthenticationTypeExternalBrowser):
		return gosnowflake.AuthTypeExternalBrowser, nil
	case string(AuthenticationTypeOkta):
		return gosnowflake.AuthTypeOkta, nil
	case string(AuthenticationTypeJwt):
		return gosnowflake.AuthTypeJwt, nil
	case string(AuthenticationTypeTokenAccessor):
		return gosnowflake.AuthTypeTokenAccessor, nil
	case string(AuthenticationTypeUsernamePasswordMfa):
		return gosnowflake.AuthTypeUsernamePasswordMFA, nil
	case string(AuthenticationTypeProgrammaticAccessToken):
		return gosnowflake.AuthTypePat, nil
	case string(AuthenticationTypeOauthClientCredentials):
		return gosnowflake.AuthTypeOAuthClientCredentials, nil
	case string(AuthenticationTypeOauthAuthorizationCode):
		return gosnowflake.AuthTypeOAuthAuthorizationCode, nil
	default:
		return gosnowflake.AuthType(0), fmt.Errorf("invalid authenticator type: %s", s)
	}
}

const GosnowflakeAuthTypeEmpty = gosnowflake.AuthType(-1)

func ToExtendedAuthenticatorType(s string) (gosnowflake.AuthType, error) {
	switch strings.ToUpper(s) {
	case string(AuthenticationTypeEmpty):
		return GosnowflakeAuthTypeEmpty, nil
	default:
		return ToAuthenticatorType(s)
	}
}

const GosnowflakeBoolConfigDefault = gosnowflake.ConfigBool(0)

type DriverLogLevel string

const (
	// these values are lower case on purpose to match gosnowflake case
	DriverLogLevelTrace   DriverLogLevel = "trace"
	DriverLogLevelDebug   DriverLogLevel = "debug"
	DriverLogLevelInfo    DriverLogLevel = "info"
	DriverLogLevelPrint   DriverLogLevel = "print"
	DriverLogLevelWarning DriverLogLevel = "warning"
	DriverLogLevelError   DriverLogLevel = "error"
	DriverLogLevelFatal   DriverLogLevel = "fatal"
	DriverLogLevelPanic   DriverLogLevel = "panic"
)

var AllDriverLogLevels = []DriverLogLevel{
	DriverLogLevelTrace,
	DriverLogLevelDebug,
	DriverLogLevelInfo,
	DriverLogLevelPrint,
	DriverLogLevelWarning,
	DriverLogLevelError,
	DriverLogLevelFatal,
	DriverLogLevelPanic,
}

func ToDriverLogLevel(s string) (DriverLogLevel, error) {
	lowerCase := strings.ToLower(s)
	switch lowerCase {
	case string(DriverLogLevelTrace),
		string(DriverLogLevelDebug),
		string(DriverLogLevelInfo),
		string(DriverLogLevelPrint),
		string(DriverLogLevelWarning),
		string(DriverLogLevelError),
		string(DriverLogLevelFatal),
		string(DriverLogLevelPanic):
		return DriverLogLevel(lowerCase), nil
	default:
		return "", fmt.Errorf("invalid driver log level: %s", s)
	}
}

// EmptyDriverConfig returns a default driver config with the Authenticator set to GosnowflakeAuthTypeEmpty.
// This is used when no config is found in the config file.
func EmptyDriverConfig() *gosnowflake.Config {
	return &gosnowflake.Config{
		Authenticator: GosnowflakeAuthTypeEmpty,
	}
}

// EmptyDriverConfigWithApplication returns a default driver config with the Authenticator set to GosnowflakeAuthTypeEmpty,
// and the passed application name.
func EmptyDriverConfigWithApplication(application string) *gosnowflake.Config {
	return &gosnowflake.Config{
		Application:   application,
		Authenticator: GosnowflakeAuthTypeEmpty,
	}
}
