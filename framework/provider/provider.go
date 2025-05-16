package provider

import (
	"context"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/snowflakedb/gosnowflake"
)

// Ensure SnowflakeProvider satisfies various provider interfaces.
var _ provider.Provider = new(SnowflakeProvider)

// SnowflakeProvider defines the provider implementation.
type SnowflakeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// SnowflakeProviderModel describes the provider data model.
type snowflakeProviderModelV0 struct {
	Account                        types.String `tfsdk:"account"`
	User                           types.String `tfsdk:"user"`
	Password                       types.String `tfsdk:"password"`
	Warehouse                      types.String `tfsdk:"warehouse"`
	Role                           types.String `tfsdk:"role"`
	Region                         types.String `tfsdk:"region"`
	ValidateDefaultParameters      types.Bool   `tfsdk:"validate_default_parameters"`
	Params                         types.Map    `tfsdk:"params"`
	ClientIP                       types.String `tfsdk:"client_ip"`
	Protocol                       types.String `tfsdk:"protocol"`
	Host                           types.String `tfsdk:"host"`
	Port                           types.Int64  `tfsdk:"port"`
	Authenticator                  types.String `tfsdk:"authenticator"`
	Passcode                       types.String `tfsdk:"passcode"`
	PasscodeInPassword             types.Bool   `tfsdk:"passcode_in_password"`
	OktaURL                        types.String `tfsdk:"okta_url"`
	LoginTimeout                   types.Int64  `tfsdk:"login_timeout"`
	RequestTimeout                 types.Int64  `tfsdk:"request_timeout"`
	JWTExpireTimeout               types.Int64  `tfsdk:"jwt_expire_timeout"`
	ClientTimeout                  types.Int64  `tfsdk:"client_timeout"`
	JWTClientTimeout               types.Int64  `tfsdk:"jwt_client_timeout"`
	ExternalBrowserTimeout         types.Int64  `tfsdk:"external_browser_timeout"`
	InsecureMode                   types.Bool   `tfsdk:"insecure_mode"`
	OCSPFailOpen                   types.Bool   `tfsdk:"ocsp_fail_open"`
	Token                          types.String `tfsdk:"token"`
	TokenAccessor                  types.List   `tfsdk:"token_accessor"`
	KeepSessionAlive               types.Bool   `tfsdk:"keep_session_alive"`
	PrivateKey                     types.String `tfsdk:"private_key"`
	PrivateKeyPassphrase           types.String `tfsdk:"private_key_passphrase"`
	DisableTelemetry               types.Bool   `tfsdk:"disable_telemetry"`
	ClientRequestMFAToken          types.Bool   `tfsdk:"client_request_mfa_token"`
	ClientStoreTemporaryCredential types.Bool   `tfsdk:"client_store_temporary_credential"`
	DisableQueryContextCache       types.Bool   `tfsdk:"disable_query_context_cache"`
	Profile                        types.String `tfsdk:"profile"`
	// Deprecated Attributes
	Username          types.String `tfsdk:"username"`
	OauthAccessToken  types.String `tfsdk:"oauth_access_token"`
	OauthRefreshToken types.String `tfsdk:"oauth_refresh_token"`
	OauthClientID     types.String `tfsdk:"oauth_client_id"`
	OauthClientSecret types.String `tfsdk:"oauth_client_secret"`
	OauthEndpoint     types.String `tfsdk:"oauth_endpoint"`
	OauthRedirectURL  types.String `tfsdk:"oauth_redirect_url"`
	BrowserAuth       types.Bool   `tfsdk:"browser_auth"`
	SessionParams     types.Map    `tfsdk:"session_params"`
}

type RefreshTokenAccesor struct {
	TokenEndpoint types.String `tfsdk:"token_endpoint"`
	RefreshToken  types.String `tfsdk:"refresh_token"`
	ClientID      types.String `tfsdk:"client_id"`
	ClientSecret  types.String `tfsdk:"client_secret"`
	RedirectURI   types.String `tfsdk:"redirect_uri"`
}

func (p *SnowflakeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "snowflake"
	resp.Version = p.version
}

func (p *SnowflakeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"account": schema.StringAttribute{
				Description: "Specifies your Snowflake account identifier assigned, by Snowflake. For information about account identifiers, see the [Snowflake documentation](https://docs.snowflake.com/en/user-guide/admin-account-identifier.html). Can also be sourced from the `SNOWFLAKE_ACCOUNT` environment variable. Required unless using `profile`.",
				Optional:    true,
			},
			"user": schema.StringAttribute{
				Description: "Username. Can also be sourced from the `SNOWFLAKE_USER` environment variable. Required unless using `profile`.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description:        "Username for username+password authentication. Can also be sourced from the `SNOWFLAKE_USER` environment variable. Required unless using `profile`.",
				Optional:           true,
				DeprecationMessage: "Use `user` instead",
			},
			"password": schema.StringAttribute{
				Description: "Password for username+password auth. Cannot be used with `browser_auth`. Can also be sourced from the `SNOWFLAKE_PASSWORD` environment variable.",
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("browser_auth"), path.MatchRoot("private_key"), path.MatchRoot("private_key_passphrase"), path.MatchRoot("oauth_access_token"), path.MatchRoot("oauth_refresh_token")),
				},
			},
			"warehouse": schema.StringAttribute{
				Description: "Specifies the virtual warehouse to use by default for queries, loading, etc. in the client session. Can also be sourced from the `SNOWFLAKE_WAREHOUSE` environment variable.",
				Optional:    true,
			},
			"role": schema.StringAttribute{
				Description: "Specifies the role to use by default for accessing Snowflake objects in the client session. Can also be sourced from the `SNOWFLAKE_ROLE` environment variable. .",
				Optional:    true,
			},
			"validate_default_parameters": schema.BoolAttribute{
				Description: "True by default. If false, disables the validation checks for Database, Schema, Warehouse and Role at the time a connection is established. Can also be sourced from the `SNOWFLAKE_VALIDATE_DEFAULT_PARAMETERS` environment variable.",
				Optional:    true,
			},
			"params": schema.MapAttribute{
				Description: "Sets other connection (i.e. session) parameters. [Parameters](https://docs.snowflake.com/en/sql-reference/parameters)",
				Optional:    true,
				ElementType: types.StringType,
			},
			"client_ip": schema.StringAttribute{
				Description: "IP address for network checks. Can also be sourced from the `SNOWFLAKE_CLIENT_IP` environment variable.",
				Optional:    true,
			},
			"protocol": schema.StringAttribute{
				Description: "Either http or https, defaults to https. Can also be sourced from the `SNOWFLAKE_PROTOCOL` environment variable.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("http", "https"),
				},
			},
			"host": schema.StringAttribute{
				Description: "Supports passing in a custom host value to the snowflake go driver for use with privatelink. Can also be sourced from the `SNOWFLAKE_HOST` environment variable. ",
				Optional:    true,
			},
			"port": schema.Int64Attribute{
				Description: "Support custom port values to snowflake go driver for use with privatelink. Can also be sourced from the `SNOWFLAKE_PORT` environment variable. ",
				Optional:    true,
			},
			"authenticator": schema.StringAttribute{
				Description: "Specifies the [authentication type](https://pkg.go.dev/github.com/snowflakedb/gosnowflake#AuthType) to use when connecting to Snowflake. Valid values include: Snowflake, OAuth, ExternalBrowser, Okta, JWT, TokenAccessor, UsernamePasswordMFA. Can also be sourced from the `SNOWFLAKE_AUTHENTICATOR` environment variable.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Snowflake", "OAuth", "ExternalBrowser", "Okta", "JWT", "TokenAccessor", "UsernamePasswordMFA"),
				},
			},
			"passcode": schema.StringAttribute{
				Description: "Specifies the passcode provided by Duo when using multi-factor authentication (MFA) for login. Can also be sourced from the `SNOWFLAKE_PASSCODE` environment variable. ",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("passcode_in_password")),
				},
			},
			"passcode_in_password": schema.BoolAttribute{
				Description: "False by default. Set to true if the MFA passcode is embedded in the login password. Appends the MFA passcode to the end of the password. Can also be sourced from the `SNOWFLAKE_PASSCODE_IN_PASSWORD` environment variable. ",
				Optional:    true,
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(path.MatchRoot("passcode")),
				},
			},
			"okta_url": schema.StringAttribute{
				Description: "The URL of the Okta server. e.g. https://example.okta.com. Can also be sourced from the `SNOWFLAKE_OKTA_URL` environment variable.",
				Optional:    true,
			},
			"login_timeout": schema.Int64Attribute{
				Description: "Login retry timeout EXCLUDING network roundtrip and read out http response. Can also be sourced from the `SNOWFLAKE_LOGIN_TIMEOUT` environment variable.",
				Optional:    true,
			},
			"request_timeout": schema.Int64Attribute{
				Description: "request retry timeout EXCLUDING network roundtrip and read out http response. Can also be sourced from the `SNOWFLAKE_REQUEST_TIMEOUT` environment variable.",
				Optional:    true,
			},
			"jwt_expire_timeout": schema.Int64Attribute{
				Description: "JWT expire after timeout in seconds. Can also be sourced from the `SNOWFLAKE_JWT_EXPIRE_TIMEOUT` environment variable.",
				Optional:    true,
			},
			"client_timeout": schema.Int64Attribute{
				Description: "The timeout in seconds for the client to complete the authentication. Default is 900 seconds. Can also be sourced from the `SNOWFLAKE_CLIENT_TIMEOUT` environment variable.",
				Optional:    true,
			},
			"jwt_client_timeout": schema.Int64Attribute{
				Description: "The timeout in seconds for the JWT client to complete the authentication. Default is 10 seconds. Can also be sourced from the `SNOWFLAKE_JWT_CLIENT_TIMEOUT` environment variable.",
				Optional:    true,
			},
			"external_browser_timeout": schema.Int64Attribute{
				Description: "The timeout in seconds for the external browser to complete the authentication. Default is 120 seconds. Can also be sourced from the `SNOWFLAKE_EXTERNAL_BROWSER_TIMEOUT` environment variable.",
				Optional:    true,
			},
			"insecure_mode": schema.BoolAttribute{
				Description: "If true, bypass the Online Certificate Status Protocol (OCSP) certificate revocation check. IMPORTANT: Change the default value for testing or emergency situations only. Can also be sourced from the `SNOWFLAKE_INSECURE_MODE` environment variable.",
				Optional:    true,
			},
			"ocsp_fail_open": schema.BoolAttribute{
				Description: "True represents OCSP fail open mode. False represents OCSP fail closed mode. Fail open true by default. Can also be sourced from the `SNOWFLAKE_OCSP_FAIL_OPEN` environment variable.",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "Token to use for OAuth and other forms of token based auth. Can also be sourced from the `SNOWFLAKE_TOKEN` environment variable.",
				Sensitive:   true,
				Optional:    true,
			},
			"keep_session_alive": schema.BoolAttribute{
				Optional:    true,
				Description: "Enables the session to persist even after the connection is closed. Can also be sourced from the `SNOWFLAKE_KEEP_SESSION_ALIVE` environment variable.",
			},
			"private_key": schema.StringAttribute{
				Description: "Private Key for username+private-key auth. Cannot be used with `browser_auth` or `password`. Can also be sourced from `SNOWFLAKE_PRIVATE_KEY` environment variable.",
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("browser_auth"), path.MatchRoot("password"), path.MatchRoot("oauth_access_token"), path.MatchRoot("oauth_refresh_token")),
				},
			},
			"private_key_passphrase": schema.StringAttribute{
				Description: "Supports the encryption ciphers aes-128-cbc, aes-128-gcm, aes-192-cbc, aes-192-gcm, aes-256-cbc, aes-256-gcm, and des-ede3-cbc. Can also be sourced from `SNOWFLAKE_PRIVATE_KEY_PASSPHRASE` environment variable.",
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("browser_auth"), path.MatchRoot("password"), path.MatchRoot("oauth_access_token"), path.MatchRoot("oauth_refresh_token")),
				},
			},
			"disable_telemetry": schema.BoolAttribute{
				Description: "Indicates whether to disable telemetry. Can also be sourced from the `SNOWFLAKE_DISABLE_TELEMETRY` environment variable.",
				Optional:    true,
			},
			"client_request_mfa_token": schema.BoolAttribute{
				Description: "When true the MFA token is cached in the credential manager. True by default in Windows/OSX. False for Linux. Can also be sourced from the `SNOWFLAKE_CLIENT_REQUEST_MFA_TOKEN` environment variable.",
				Optional:    true,
			},
			"client_store_temporary_credential": schema.BoolAttribute{
				Description: "When true the ID token is cached in the credential manager. True by default in Windows/OSX. False for Linux. Can also be sourced from the `SNOWFLAKE_CLIENT_STORE_TEMPORARY_CREDENTIAL` environment variable.",
				Optional:    true,
			},
			"disable_query_context_cache": schema.BoolAttribute{
				Description: "Should HTAP query context cache be disabled. Can also be sourced from the `SNOWFLAKE_DISABLE_QUERY_CONTEXT_CACHE` environment variable.",
				Optional:    true,
			},
			"profile": schema.StringAttribute{
				Description: "Sets the profile to read from ~/.snowflake/config file. Can also be sourced from the `SNOWFLAKE_PROFILE` environment variable.",
				Optional:    true,
			},
			/*
					Feature not yet released as of latest gosnowflake release
					https://github.com/snowflakedb/gosnowflake/blob/master/dsn.go#L103
				"include_retry_reason": schema.BoolAttribute {
					Description: "Should retried request contain retry reason. Can also be sourced from the `SNOWFLAKE_INCLUDE_RETRY_REASON` environment variable.",
					Optional:    true,
				},
			*/
			// Deprecated Attributes
			"region": schema.StringAttribute{
				Description:        "Snowflake region, such as \"eu-central-1\", with this parameter. However, since this parameter is deprecated, it is best to specify the region as part of the account parameter. For details, see the description of the account parameter. [Snowflake region](https://docs.snowflake.com/en/user-guide/intro-regions.html) to use.  Required if using the [legacy format for the `account` identifier](https://docs.snowflake.com/en/user-guide/admin-account-identifier.html#format-2-legacy-account-locator-in-a-region) in the form of `<cloud_region_id>.<cloud>`. Can also be sourced from the `SNOWFLAKE_REGION` environment variable. ",
				Optional:           true,
				DeprecationMessage: "Specify the region as part of the account parameter",
			},
			"session_params": schema.MapAttribute{
				Description:        "Sets session parameters. [Parameters](https://docs.snowflake.com/en/sql-reference/parameters)",
				Optional:           true,
				ElementType:        types.StringType,
				DeprecationMessage: "Use `params` instead",
			},
			"oauth_access_token": schema.StringAttribute{
				Description: "Token for use with OAuth. Generating the token is left to other tools. Cannot be used with `browser_auth`, `oauth_refresh_token` or `password`. Can also be sourced from `SNOWFLAKE_OAUTH_ACCESS_TOKEN` environment variable.",
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("browser_auth"), path.MatchRoot("private_key"), path.MatchRoot("private_key_passphrase"), path.MatchRoot("password"), path.MatchRoot("oauth_refresh_token")),
				},
				DeprecationMessage: "Use `token` instead",
			},
			"oauth_refresh_token": schema.StringAttribute{
				Description: "Token for use with OAuth. Setup and generation of the token is left to other tools. Should be used in conjunction with `oauth_client_id`, `oauth_client_secret`, `oauth_endpoint`, `oauth_redirect_url`. Cannot be used with `browser_auth`, `oauth_access_token` or `password`. Can also be sourced from `SNOWFLAKE_OAUTH_REFRESH_TOKEN` environment variable.",
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("browser_auth"), path.MatchRoot("private_key"), path.MatchRoot("private_key_passphrase"), path.MatchRoot("password"), path.MatchRoot("oauth_access_token")),
					stringvalidator.AlsoRequires(path.MatchRoot("oauth_client_id"), path.MatchRoot("oauth_client_secret"), path.MatchRoot("oauth_endpoint"), path.MatchRoot("oauth_redirect_url")),
				},
				DeprecationMessage: "Use `token_accessor.0.refresh_token` instead",
			},
			"oauth_client_id": schema.StringAttribute{
				Description: "Required when `oauth_refresh_token` is used. Can also be sourced from `SNOWFLAKE_OAUTH_CLIENT_ID` environment variable.",
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("browser_auth"), path.MatchRoot("private_key"), path.MatchRoot("private_key_passphrase"), path.MatchRoot("password"), path.MatchRoot("oauth_access_token")),
					stringvalidator.AlsoRequires(path.MatchRoot("oauth_refresh_token"), path.MatchRoot("oauth_client_secret"), path.MatchRoot("oauth_endpoint"), path.MatchRoot("oauth_redirect_url")),
				},
				DeprecationMessage: "Use `token_accessor.0.client_id` instead",
			},
			"oauth_client_secret": schema.StringAttribute{
				Description: "Required when `oauth_refresh_token` is used. Can also be sourced from `SNOWFLAKE_OAUTH_CLIENT_SECRET` environment variable.",
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("browser_auth"), path.MatchRoot("private_key"), path.MatchRoot("private_key_passphrase"), path.MatchRoot("password"), path.MatchRoot("oauth_access_token")),
					stringvalidator.AlsoRequires(path.MatchRoot("oauth_refresh_token"), path.MatchRoot("oauth_client_id"), path.MatchRoot("oauth_endpoint"), path.MatchRoot("oauth_redirect_url")),
				},
				DeprecationMessage: "Use `token_accessor.0.client_secret` instead",
			},
			"oauth_endpoint": schema.StringAttribute{
				Description: "Required when `oauth_refresh_token` is used. Can also be sourced from `SNOWFLAKE_OAUTH_ENDPOINT` environment variable.",
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("browser_auth"), path.MatchRoot("private_key"), path.MatchRoot("private_key_passphrase"), path.MatchRoot("password"), path.MatchRoot("oauth_access_token")),
					stringvalidator.AlsoRequires(path.MatchRoot("oauth_refresh_token"), path.MatchRoot("oauth_client_id"), path.MatchRoot("oauth_client_secret"), path.MatchRoot("oauth_redirect_url")),
				},
				DeprecationMessage: "Use `token_accessor.0.token_endpoint` instead",
			},
			"oauth_redirect_url": schema.StringAttribute{
				Description: "Required when `oauth_refresh_token` is used. Can also be sourced from `SNOWFLAKE_OAUTH_REDIRECT_URL` environment variable.",
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("browser_auth"), path.MatchRoot("private_key"), path.MatchRoot("private_key_passphrase"), path.MatchRoot("password"), path.MatchRoot("oauth_access_token")),
					stringvalidator.AlsoRequires(path.MatchRoot("oauth_refresh_token"), path.MatchRoot("oauth_client_id"), path.MatchRoot("oauth_client_secret"), path.MatchRoot("oauth_endpoint")),
				},
				DeprecationMessage: "Use `token_accessor.0.redirect_uri` instead",
			},
			"browser_auth": schema.BoolAttribute{
				Description:        "Required when `oauth_refresh_token` is used. Can also be sourced from `SNOWFLAKE_USE_BROWSER_AUTH` environment variable.",
				Optional:           true,
				Sensitive:          false,
				DeprecationMessage: "Use `authenticator` instead",
			},
		},
		Blocks: map[string]schema.Block{
			"token_accessor": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"token_endpoint": schema.StringAttribute{
							Description: "The token endpoint for the OAuth provider e.g. https://{yourDomain}/oauth/token when using a refresh token to renew access token. Can also be sourced from the `SNOWFLAKE_TOKEN_ACCESSOR_TOKEN_ENDPOINT` environment variable.",
							Required:    true,
							Sensitive:   true,
						},
						"refresh_token": schema.StringAttribute{
							Description: "The refresh token for the OAuth provider when using a refresh token to renew access token. Can also be sourced from the `SNOWFLAKE_TOKEN_ACCESSOR_REFRESH_TOKEN` environment variable.",
							Required:    true,
							Sensitive:   true,
						},
						"client_id": schema.StringAttribute{
							Description: "The client ID for the OAuth provider when using a refresh token to renew access token. Can also be sourced from the `SNOWFLAKE_TOKEN_ACCESSOR_CLIENT_ID` environment variable.",
							Required:    true,
							Sensitive:   true,
						},
						"client_secret": schema.StringAttribute{
							Description: "The client secret for the OAuth provider when using a refresh token to renew access token. Can also be sourced from the `SNOWFLAKE_TOKEN_ACCESSOR_CLIENT_SECRET` environment variable.",
							Required:    true,
							Sensitive:   true,
						},
						"redirect_uri": schema.StringAttribute{
							Description: "The redirect URI for the OAuth provider when using a refresh token to renew access token. Can also be sourced from the `SNOWFLAKE_TOKEN_ACCESSOR_REDIRECT_URI` environment variable.",
							Required:    true,
							Sensitive:   true,
						},
					},
				},
			},
		},
	}
}

func (p *SnowflakeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data snowflakeProviderModelV0

	// Read configuration data into model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	config := &gosnowflake.Config{
		Application: "terraform-provider-snowflake",
	}

	account := os.Getenv("SNOWFLAKE_ACCOUNT")
	if data.Account.ValueString() != "" {
		account = data.Account.ValueString()
	}
	if account != "" {
		config.Account = account
	}

	user := os.Getenv("SNOWFLAKE_USER")
	if user == "" {
		user = os.Getenv("SNOWFLAKE_USERNAME")
	}
	if data.Username.ValueString() != "" {
		user = data.Username.ValueString()
	}
	if data.User.ValueString() != "" {
		user = data.User.ValueString()
	}
	if user != "" {
		config.User = user
	}

	password := os.Getenv("SNOWFLAKE_PASSWORD")
	if data.Password.ValueString() != "" {
		password = data.Password.ValueString()
	}
	if password != "" {
		config.Password = password
	}

	warehouse := os.Getenv("SNOWFLAKE_WAREHOUSE")
	if data.Warehouse.ValueString() != "" {
		warehouse = data.Warehouse.ValueString()
	}
	if warehouse != "" {
		config.Warehouse = warehouse
	}

	role := os.Getenv("SNOWFLAKE_ROLE")
	if data.Role.ValueString() != "" {
		role = data.Role.ValueString()
	}
	if role != "" {
		config.Role = role
	}

	validateDefaultParameters := getBoolEnv("SNOWFLAKE_VALIDATE_DEFAULT_PARAMETERS", true)
	if !data.ValidateDefaultParameters.IsNull() && !data.ValidateDefaultParameters.IsUnknown() {
		validateDefaultParameters = data.ValidateDefaultParameters.ValueBool()
	}
	if validateDefaultParameters {
		config.ValidateDefaultParameters = gosnowflake.ConfigBoolTrue
	} else {
		config.ValidateDefaultParameters = gosnowflake.ConfigBoolFalse
	}

	clientIP := os.Getenv("SNOWFLAKE_CLIENT_IP")
	if data.ClientIP.ValueString() != "" {
		clientIP = data.ClientIP.ValueString()
	}
	if clientIP != "" {
		config.ClientIP = net.ParseIP(clientIP)
	}

	protocol := os.Getenv("SNOWFLAKE_PROTOCOL")
	if data.Protocol.ValueString() != "" {
		protocol = data.Protocol.ValueString()
	}
	if protocol != "" {
		config.Protocol = protocol
	}

	host := os.Getenv("SNOWFLAKE_HOST")
	if data.Host.ValueString() != "" {
		host = data.Host.ValueString()
	}
	if host != "" {
		config.Host = host
	}

	port := getInt64Env("SNOWFLAKE_PORT", -1)
	if !data.Port.IsNull() && !data.Port.IsUnknown() {
		port = data.Port.ValueInt64()
	}
	if port > 0 {
		config.Port = int(port)
	}

	browserAuth := getBoolEnv("SNOWFLAKE_USE_BROWSER_AUTH", false)
	if !data.BrowserAuth.IsNull() && !data.BrowserAuth.IsUnknown() {
		browserAuth = data.BrowserAuth.ValueBool()
	}
	if browserAuth {
		config.Authenticator = gosnowflake.AuthTypeExternalBrowser
	}

	authenticator := os.Getenv("SNOWFLAKE_AUTHENTICATOR")
	if data.Authenticator.ValueString() != "" {
		authenticator = data.Authenticator.ValueString()
	}
	if authenticator != "" {
		config.Authenticator = toAuthenticatorType(authenticator)
	}

	passcode := os.Getenv("SNOWFLAKE_PASSCODE")
	if data.Passcode.ValueString() != "" {
		passcode = data.Passcode.ValueString()
	}
	if passcode != "" {
		config.Passcode = passcode
	}

	passcodeInPassword := getBoolEnv("SNOWFLAKE_PASSCODE_IN_PASSWORD", false)
	if !data.PasscodeInPassword.IsNull() && !data.PasscodeInPassword.IsUnknown() {
		passcodeInPassword = data.PasscodeInPassword.ValueBool()
	}
	config.PasscodeInPassword = passcodeInPassword

	oktaURL := os.Getenv("SNOWFLAKE_OKTA_URL")
	if data.OktaURL.ValueString() != "" {
		oktaURL = data.OktaURL.ValueString()
	}
	if oktaURL != "" {
		parsedOktaURL, err := url.Parse(oktaURL)
		if err != nil {
			resp.Diagnostics.AddError("Error parsing Okta URL", err.Error())
		}
		config.OktaURL = parsedOktaURL
	}

	loginTimeout := getInt64Env("SNOWFLAKE_LOGIN_TIMEOUT", -1)
	if !data.LoginTimeout.IsNull() && !data.LoginTimeout.IsUnknown() {
		loginTimeout = data.LoginTimeout.ValueInt64()
	}
	if loginTimeout > 0 {
		config.LoginTimeout = time.Second * time.Duration(loginTimeout)
	}

	requestTimeout := getInt64Env("SNOWFLAKE_REQUEST_TIMEOUT", -1)
	if !data.RequestTimeout.IsNull() && !data.RequestTimeout.IsUnknown() {
		requestTimeout = data.RequestTimeout.ValueInt64()
	}
	if requestTimeout > 0 {
		config.RequestTimeout = time.Second * time.Duration(requestTimeout)
	}

	jwtExpireTimeout := getInt64Env("SNOWFLAKE_JWT_EXPIRE_TIMEOUT", -1)
	if !data.JWTExpireTimeout.IsNull() && !data.JWTExpireTimeout.IsUnknown() {
		jwtExpireTimeout = data.JWTClientTimeout.ValueInt64()
	}
	if jwtExpireTimeout > 0 {
		config.JWTClientTimeout = time.Second * time.Duration(jwtExpireTimeout)
	}

	clientTimeout := getInt64Env("SNOWFLAKE_CLIENT_TIMEOUT", -1)
	if !data.ClientTimeout.IsNull() && !data.ClientTimeout.IsUnknown() {
		clientTimeout = data.ClientTimeout.ValueInt64()
	}
	if clientTimeout > 0 {
		config.ClientTimeout = time.Second * time.Duration(clientTimeout)
	}

	jwtClientTimeout := getInt64Env("SNOWFLAKE_JWT_CLIENT_TIMEOUT", -1)
	if !data.JWTClientTimeout.IsNull() && !data.JWTClientTimeout.IsUnknown() {
		jwtClientTimeout = data.JWTClientTimeout.ValueInt64()
	}
	if jwtClientTimeout > 0 {
		config.JWTClientTimeout = time.Second * time.Duration(jwtClientTimeout)
	}

	externalBrowserTimeout := getInt64Env("SNOWFLAKE_EXTERNAL_BROWSER_TIMEOUT", -1)
	if !data.ExternalBrowserTimeout.IsNull() && !data.ExternalBrowserTimeout.IsUnknown() {
		externalBrowserTimeout = data.ExternalBrowserTimeout.ValueInt64()
	}
	if externalBrowserTimeout > 0 {
		config.ExternalBrowserTimeout = time.Second * time.Duration(externalBrowserTimeout)
	}

	insecureMode := getBoolEnv("SNOWFLAKE_INSECURE_MODE", false)
	if !data.InsecureMode.IsNull() && !data.InsecureMode.IsUnknown() {
		insecureMode = data.InsecureMode.ValueBool()
	}
	config.InsecureMode = insecureMode //nolint:staticcheck

	ocspFailOpen := getBoolEnv("SNOWFLAKE_OCSP_FAIL_OPEN", true)
	if !data.OCSPFailOpen.IsNull() && !data.OCSPFailOpen.IsUnknown() {
		ocspFailOpen = data.OCSPFailOpen.ValueBool()
	}
	if ocspFailOpen {
		config.OCSPFailOpen = gosnowflake.OCSPFailOpenTrue
	} else {
		config.OCSPFailOpen = gosnowflake.OCSPFailOpenFalse
	}

	token := os.Getenv("SNOWFLAKE_TOKEN")
	if data.Token.ValueString() != "" {
		token = data.Token.ValueString()
	}
	if token != "" {
		config.Token = token
	}

	keepSessionAlive := getBoolEnv("SNOWFLAKE_KEEP_SESSION_ALIVE", false)
	if !data.KeepSessionAlive.IsNull() && !data.KeepSessionAlive.IsUnknown() {
		keepSessionAlive = data.KeepSessionAlive.ValueBool()
	}
	config.KeepSessionAlive = keepSessionAlive

	privateKey := os.Getenv("SNOWFLAKE_PRIVATE_KEY")
	if data.PrivateKey.ValueString() != "" {
		privateKey = data.PrivateKey.ValueString()
	}
	privateKeyPassphrase := os.Getenv("SNOWFLAKE_PRIVATE_KEY_PASSPHRASE")
	if data.PrivateKeyPassphrase.ValueString() != "" {
		privateKeyPassphrase = data.PrivateKeyPassphrase.ValueString()
	}
	if privateKey != "" {
		if v, err := getPrivateKey(privateKey, privateKeyPassphrase); err != nil && v != nil {
			config.PrivateKey = v
		}
	}
	disableTelemetry := getBoolEnv("SNOWFLAKE_DISABLE_TELEMETRY", false)
	if !data.DisableTelemetry.IsNull() && !data.DisableTelemetry.IsUnknown() {
		disableTelemetry = data.DisableTelemetry.ValueBool()
	}
	config.DisableTelemetry = disableTelemetry

	clientRequestMFAToken := getBoolEnv("SNOWFLAKE_CLIENT_REQUEST_MFA_TOKEN", true)
	if !data.ClientRequestMFAToken.IsNull() && !data.ClientRequestMFAToken.IsUnknown() {
		clientRequestMFAToken = data.ClientRequestMFAToken.ValueBool()
	}
	if clientRequestMFAToken {
		config.ClientRequestMfaToken = gosnowflake.ConfigBoolTrue
	} else {
		config.ClientRequestMfaToken = gosnowflake.ConfigBoolFalse
	}

	clientStoreTemporaryCredential := getBoolEnv("SNOWFLAKE_CLIENT_STORE_TEMPORARY_CREDENTIAL", true)
	if !data.ClientStoreTemporaryCredential.IsNull() && !data.ClientStoreTemporaryCredential.IsUnknown() {
		clientStoreTemporaryCredential = data.ClientStoreTemporaryCredential.ValueBool()
	}
	if clientStoreTemporaryCredential {
		config.ClientStoreTemporaryCredential = gosnowflake.ConfigBoolTrue
	} else {
		config.ClientStoreTemporaryCredential = gosnowflake.ConfigBoolFalse
	}

	disableQueryContextCache := getBoolEnv("SNOWFLAKE_DISABLE_QUERY_CONTEXT_CACHE", false)
	if !data.DisableQueryContextCache.IsNull() && !data.DisableQueryContextCache.IsUnknown() {
		disableQueryContextCache = data.DisableQueryContextCache.ValueBool()
	}
	config.DisableQueryContextCache = disableQueryContextCache

	tokenEndpoint := os.Getenv("SNOWFLAKE_TOKEN_ACCESSOR_TOKEN_ENDPOINT")
	if tokenEndpoint == "" {
		tokenEndpoint = os.Getenv("SNOWFLAKE_OAUTH_ENDPOINT")
	}
	if data.OauthEndpoint.ValueString() != "" {
		tokenEndpoint = data.OauthEndpoint.ValueString()
	}
	refreshToken := os.Getenv("SNOWFLAKE_TOKEN_ACCESSOR_REFRESH_TOKEN")
	if refreshToken == "" {
		refreshToken = os.Getenv("SNOWFLAKE_OAUTH_REFRESH_TOKEN")
	}
	if data.OauthRefreshToken.ValueString() != "" {
		refreshToken = data.OauthRefreshToken.ValueString()
	}
	clientID := os.Getenv("SNOWFLAKE_TOKEN_ACCESSOR_CLIENT_ID")
	if clientID == "" {
		clientID = os.Getenv("SNOWFLAKE_OAUTH_CLIENT_ID")
	}
	if data.OauthClientID.ValueString() != "" {
		clientID = data.OauthClientID.ValueString()
	}
	clientSecret := os.Getenv("SNOWFLAKE_TOKEN_ACCESSOR_CLIENT_SECRET")
	if clientSecret == "" {
		clientSecret = os.Getenv("SNOWFLAKE_OAUTH_CLIENT_SECRET")
	}
	if data.OauthClientSecret.ValueString() != "" {
		clientSecret = data.OauthClientSecret.ValueString()
	}
	redirectURI := os.Getenv("SNOWFLAKE_TOKEN_ACCESSOR_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = os.Getenv("SNOWFLAKE_OAUTH_REDIRECT_URL")
	}
	if data.OauthRedirectURL.ValueString() != "" {
		redirectURI = data.OauthRedirectURL.ValueString()
	}
	var tokenAccesors []RefreshTokenAccesor
	data.TokenAccessor.ElementsAs(ctx, &tokenAccesors, false)
	if len(tokenAccesors) > 0 {
		tokenAccessor := tokenAccesors[0]
		if tokenAccessor.TokenEndpoint.ValueString() != "" {
			tokenEndpoint = tokenAccessor.TokenEndpoint.ValueString()
		}
		if tokenAccessor.RefreshToken.ValueString() != "" {
			refreshToken = tokenAccessor.RefreshToken.ValueString()
		}
		if tokenAccessor.ClientID.ValueString() != "" {
			clientID = tokenAccessor.ClientID.ValueString()
		}
		if tokenAccessor.ClientSecret.ValueString() != "" {
			clientSecret = tokenAccessor.ClientSecret.ValueString()
		}
		if tokenAccessor.RedirectURI.ValueString() != "" {
			redirectURI = tokenAccessor.RedirectURI.ValueString()
		}
	}

	if tokenEndpoint != "" && refreshToken != "" && clientID != "" && clientSecret != "" && redirectURI != "" {
		accessToken, err := GetAccessTokenWithRefreshToken(tokenEndpoint, clientID, clientSecret, refreshToken, redirectURI)
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving access token from refresh token", err.Error())
		}
		config.Token = accessToken
		config.Authenticator = gosnowflake.AuthTypeOAuth
	}

	region := os.Getenv("SNOWFLAKE_REGION")
	if data.Region.ValueString() != "" {
		region = data.Region.ValueString()
	}
	if region != "" {
		config.Region = region
	}

	if !data.SessionParams.IsNull() && !data.SessionParams.IsUnknown() {
		var m map[string]interface{}
		params := make(map[string]*string, 0)
		data.SessionParams.ElementsAs(ctx, m, false)
		for k, v := range m {
			s := v.(string)
			params[k] = &s
		}
		config.Params = params
	}

	if !data.Params.IsNull() && !data.Params.IsUnknown() {
		var m map[string]interface{}
		params := make(map[string]*string, 0)
		data.Params.ElementsAs(ctx, m, false)
		for k, v := range m {
			s := v.(string)
			params[k] = &s
		}
		config.Params = params
	}

	profile := os.Getenv("SNOWFLAKE_PROFILE")
	if data.Profile.ValueString() != "" {
		profile = data.Profile.ValueString()
	}

	if profile != "" {
		if profile == "default" {
			defaultConfig := sdk.DefaultConfig()
			if defaultConfig.Account == "" || defaultConfig.User == "" {
				resp.Diagnostics.AddError("Error retrieving default profile config", "default profile not found in config file")
			}
			config = sdk.MergeConfig(config, defaultConfig)
		} else {
			profileConfig, err := sdk.ProfileConfig(profile)
			if err != nil {
				resp.Diagnostics.AddError("Error retrieving profile config", err.Error())
			}
			if profileConfig == nil {
				resp.Diagnostics.AddError("Error retrieving profile config", "profile with name: "+profile+" not found in config file")
			}
			// merge any credentials found in profile with config
			config = sdk.MergeConfig(config, profileConfig)
		}
	}

	client, err := sdk.NewClient(config)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Snowflake client", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}
	providerData := &ProviderData{
		client: client,
	}
	resp.DataSourceData = providerData
	resp.ResourceData = providerData
}

type ProviderData struct {
	client *sdk.Client
}

func (p *SnowflakeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// NewResourceMonitorResource,
	}
}

func (p *SnowflakeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SnowflakeProvider{
			version: version,
		}
	}
}
