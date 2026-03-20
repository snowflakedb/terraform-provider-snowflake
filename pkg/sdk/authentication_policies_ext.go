package sdk

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

// ClientPolicyDriverType is the client/driver type for CLIENT_POLICY minimum version (distinct from ClientTypesOption).
type ClientPolicyDriverType string

const (
	ClientPolicyDriverTypeJdbcDriver                 ClientPolicyDriverType = "JDBC_DRIVER"
	ClientPolicyDriverTypeOdbcDriver                 ClientPolicyDriverType = "ODBC_DRIVER"
	ClientPolicyDriverTypePythonDriver               ClientPolicyDriverType = "PYTHON_DRIVER"
	ClientPolicyDriverTypeJavascriptDriver           ClientPolicyDriverType = "JAVASCRIPT_DRIVER"
	ClientPolicyDriverTypeCDriver                    ClientPolicyDriverType = "C_DRIVER"
	ClientPolicyDriverTypeGoDriver                   ClientPolicyDriverType = "GO_DRIVER"
	ClientPolicyDriverTypePhpDriver                  ClientPolicyDriverType = "PHP_DRIVER"
	ClientPolicyDriverTypeDotnetDriver               ClientPolicyDriverType = "DOTNET_DRIVER"
	ClientPolicyDriverTypeSQLApi                     ClientPolicyDriverType = "SQL_API"
	ClientPolicyDriverTypeSnowpipeStreamingClientSdk ClientPolicyDriverType = "SNOWPIPE_STREAMING_CLIENT_SDK"
	ClientPolicyDriverTypePyCore                     ClientPolicyDriverType = "PY_CORE"
	ClientPolicyDriverTypeSprocPython                ClientPolicyDriverType = "SPROC_PYTHON"
	ClientPolicyDriverTypePythonSnowpark             ClientPolicyDriverType = "PYTHON_SNOWPARK"
	ClientPolicyDriverTypeSQLAlchemy                 ClientPolicyDriverType = "SQL_ALCHEMY"
	ClientPolicyDriverTypeSnowpark                   ClientPolicyDriverType = "SNOWPARK"
	ClientPolicyDriverTypeSnowflakeClient            ClientPolicyDriverType = "SNOWFLAKE_CLIENT"
)

var AllClientPolicyDriverTypes = []ClientPolicyDriverType{
	ClientPolicyDriverTypeJdbcDriver,
	ClientPolicyDriverTypeOdbcDriver,
	ClientPolicyDriverTypePythonDriver,
	ClientPolicyDriverTypeJavascriptDriver,
	ClientPolicyDriverTypeCDriver,
	ClientPolicyDriverTypeGoDriver,
	ClientPolicyDriverTypePhpDriver,
	ClientPolicyDriverTypeDotnetDriver,
	ClientPolicyDriverTypeSQLApi,
	ClientPolicyDriverTypeSnowpipeStreamingClientSdk,
	ClientPolicyDriverTypePyCore,
	ClientPolicyDriverTypeSprocPython,
	ClientPolicyDriverTypePythonSnowpark,
	ClientPolicyDriverTypeSQLAlchemy,
	ClientPolicyDriverTypeSnowpark,
	ClientPolicyDriverTypeSnowflakeClient,
}

// ToClientPolicyDriverType validates and returns the client policy driver type.
func ToClientPolicyDriverType(s string) (ClientPolicyDriverType, error) {
	u := strings.ToUpper(s)
	if !slices.Contains(AllClientPolicyDriverTypes, ClientPolicyDriverType(u)) {
		return "", fmt.Errorf("invalid client policy driver type: %s", s)
	}
	return ClientPolicyDriverType(u), nil
}

type AuthenticationPolicyDetails []AuthenticationPolicyDescription

func (v AuthenticationPolicyDetails) GetAuthenticationMethods() ([]AuthenticationMethodsOption, error) {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == "AUTHENTICATION_METHODS" })
	if err != nil {
		return nil, err
	}
	return collections.MapErr(ParseCommaSeparatedStringArray(raw.Value, false), ToAuthenticationMethodsOption)
}

func (v AuthenticationPolicyDetails) Raw(key string) string {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == key })
	if err != nil {
		return ""
	}
	return raw.Value
}

func (v AuthenticationPolicyDetails) GetMfaEnrollment() (MfaEnrollmentReadOption, error) {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == "MFA_ENROLLMENT" })
	if err != nil {
		return "", err
	}
	return ToMfaEnrollmentReadOption(raw.Value)
}

func (v AuthenticationPolicyDetails) GetClientTypes() ([]ClientTypesOption, error) {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == "CLIENT_TYPES" })
	if err != nil {
		return nil, err
	}
	return collections.MapErr(ParseCommaSeparatedStringArray(raw.Value, false), ToClientTypesOption)
}

func (v AuthenticationPolicyDetails) GetSecurityIntegrations() ([]AccountObjectIdentifier, error) {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == "SECURITY_INTEGRATIONS" })
	if err != nil {
		return nil, err
	}
	return ParseCommaSeparatedAccountObjectIdentifierArray(raw.Value)
}

func (v AuthenticationPolicyDetails) GetMfaAuthenticationMethods() ([]MfaAuthenticationMethodsReadOption, error) {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == "MFA_AUTHENTICATION_METHODS" })
	if err != nil {
		return nil, err
	}
	return collections.MapErr(ParseCommaSeparatedStringArray(raw.Value, false), ToMfaAuthenticationMethodsReadOption)
}
