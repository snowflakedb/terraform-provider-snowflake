package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type AuthenticationMethodsOption string

const (
	AuthenticationMethodsAll                     AuthenticationMethodsOption = "ALL"
	AuthenticationMethodsSaml                    AuthenticationMethodsOption = "SAML"
	AuthenticationMethodsPassword                AuthenticationMethodsOption = "PASSWORD"
	AuthenticationMethodsOauth                   AuthenticationMethodsOption = "OAUTH"
	AuthenticationMethodsKeyPair                 AuthenticationMethodsOption = "KEYPAIR"
	AuthenticationMethodsProgrammaticAccessToken AuthenticationMethodsOption = "PROGRAMMATIC_ACCESS_TOKEN" //nolint:gosec
	AuthenticationMethodsWorkloadIdentity        AuthenticationMethodsOption = "WORKLOAD_IDENTITY"
)

var AllAuthenticationMethods = []AuthenticationMethodsOption{
	AuthenticationMethodsAll,
	AuthenticationMethodsSaml,
	AuthenticationMethodsPassword,
	AuthenticationMethodsOauth,
	AuthenticationMethodsKeyPair,
	AuthenticationMethodsProgrammaticAccessToken,
	AuthenticationMethodsWorkloadIdentity,
}

type MfaAuthenticationMethodsOption string

const (
	MfaAuthenticationMethodsAll      MfaAuthenticationMethodsOption = "ALL"
	MfaAuthenticationMethodsSaml     MfaAuthenticationMethodsOption = "SAML"
	MfaAuthenticationMethodsPassword MfaAuthenticationMethodsOption = "PASSWORD"
)

var AllMfaAuthenticationMethods = []MfaAuthenticationMethodsOption{
	MfaAuthenticationMethodsAll,
	MfaAuthenticationMethodsSaml,
	MfaAuthenticationMethodsPassword,
}

type MfaEnrollmentOption string

const (
	MfaEnrollmentRequired             MfaEnrollmentOption = "REQUIRED"
	MfaEnrollmentRequiredPasswordOnly MfaEnrollmentOption = "REQUIRED_PASSWORD_ONLY"
	MfaEnrollmentOptional             MfaEnrollmentOption = "OPTIONAL"
)

var AllMfaEnrollmentOptions = []MfaEnrollmentOption{
	MfaEnrollmentRequired,
	MfaEnrollmentRequiredPasswordOnly,
	MfaEnrollmentOptional,
}

type MfaEnrollmentReadOption string

const (
	MfaEnrollmentReadRequired                        MfaEnrollmentReadOption = "REQUIRED"
	MfaEnrollmentReadRequiredPasswordOnly            MfaEnrollmentReadOption = "REQUIRED_PASSWORD_ONLY"
	MfaEnrollmentReadOptional                        MfaEnrollmentReadOption = "OPTIONAL"
	MfaEnrollmentReadRequiredSnowflakeUiPasswordOnly MfaEnrollmentReadOption = "REQUIRED_SNOWFLAKE_UI_PASSWORD_ONLY"
)

var AllMfaEnrollmentReadOptions = []MfaEnrollmentReadOption{
	MfaEnrollmentReadRequired,
	MfaEnrollmentReadRequiredPasswordOnly,
	MfaEnrollmentReadOptional,
	MfaEnrollmentReadRequiredSnowflakeUiPasswordOnly,
}

type ClientTypesOption string

const (
	ClientTypesAll          ClientTypesOption = "ALL"
	ClientTypesSnowflakeUi  ClientTypesOption = "SNOWFLAKE_UI"
	ClientTypesDrivers      ClientTypesOption = "DRIVERS"
	ClientTypesSnowSql      ClientTypesOption = "SNOWSQL"
	ClientTypesSnowflakeCli ClientTypesOption = "SNOWFLAKE_CLI"
)

var AllClientTypes = []ClientTypesOption{
	ClientTypesAll,
	ClientTypesSnowflakeUi,
	ClientTypesDrivers,
	ClientTypesSnowSql,
	ClientTypesSnowflakeCli,
}

type MfaPolicyAllowedMethodsOption string

const (
	MfaPolicyAllowedMethodAll     MfaPolicyAllowedMethodsOption = "ALL"
	MfaPolicyAllowedMethodPassKey MfaPolicyAllowedMethodsOption = "PASSKEY"
	MfaPolicyAllowedMethodTotp    MfaPolicyAllowedMethodsOption = "TOTP"
	MfaPolicyAllowedMethodDuo     MfaPolicyAllowedMethodsOption = "DUO"
)

var AllMfaPolicyOptions = []MfaPolicyAllowedMethodsOption{
	MfaPolicyAllowedMethodAll,
	MfaPolicyAllowedMethodPassKey,
	MfaPolicyAllowedMethodTotp,
	MfaPolicyAllowedMethodDuo,
}

type NetworkPolicyEvaluationOption string

const (
	NetworkPolicyEvaluationEnforcedRequired    NetworkPolicyEvaluationOption = "ENFORCED_REQUIRED"
	NetworkPolicyEvaluationEnforcedNotRequired NetworkPolicyEvaluationOption = "ENFORCED_NOT_REQUIRED"
	NetworkPolicyEvaluationNotEnforced         NetworkPolicyEvaluationOption = "NOT_ENFORCED"
)

var AllNetworkPolicyEvaluationOptions = []NetworkPolicyEvaluationOption{
	NetworkPolicyEvaluationEnforcedRequired,
	NetworkPolicyEvaluationEnforcedNotRequired,
	NetworkPolicyEvaluationNotEnforced,
}

type AllowedProviderOption string

const (
	AllowedProviderAll   AllowedProviderOption = "ALL"
	AllowedProviderAws   AllowedProviderOption = "AWS"
	AllowedProviderAzure AllowedProviderOption = "AZURE"
	AllowedProviderGcp   AllowedProviderOption = "GCP"
	AllowedProviderOidc  AllowedProviderOption = "OIDC"
)

var AllAllowedProviderOptions = []AllowedProviderOption{
	AllowedProviderAll,
	AllowedProviderAws,
	AllowedProviderAzure,
	AllowedProviderGcp,
	AllowedProviderOidc,
}

type EnforceMfaOnExternalAuthenticationOption string

const (
	EnforceMfaOnExternalAuthenticationAll  EnforceMfaOnExternalAuthenticationOption = "ALL"
	EnforceMfaOnExternalAuthenticationNone EnforceMfaOnExternalAuthenticationOption = "NONE"
)

var AllEnforceMfaOnExternalAuthenticationOptions = []EnforceMfaOnExternalAuthenticationOption{
	EnforceMfaOnExternalAuthenticationAll,
	EnforceMfaOnExternalAuthenticationNone,
}

func ToAuthenticationMethodsOption(s string) (AuthenticationMethodsOption, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllAuthenticationMethods, AuthenticationMethodsOption(s)) {
		return "", fmt.Errorf("invalid authentication method: %s", s)
	}
	return AuthenticationMethodsOption(s), nil
}

func ToMfaAuthenticationMethodsOption(s string) (MfaAuthenticationMethodsOption, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllMfaAuthenticationMethods, MfaAuthenticationMethodsOption(s)) {
		return "", fmt.Errorf("invalid MFA authentication method: %s", s)
	}
	return MfaAuthenticationMethodsOption(s), nil
}

func ToMfaEnrollmentOption(s string) (MfaEnrollmentOption, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllMfaEnrollmentOptions, MfaEnrollmentOption(s)) {
		return "", fmt.Errorf("invalid MFA enrollment option: %s", s)
	}
	return MfaEnrollmentOption(s), nil
}

func ToMfaEnrollmentReadOption(s string) (MfaEnrollmentReadOption, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllMfaEnrollmentReadOptions, MfaEnrollmentReadOption(s)) {
		return "", fmt.Errorf("invalid MFA enrollment read option: %s", s)
	}
	return MfaEnrollmentReadOption(s), nil
}

func ToClientTypesOption(s string) (ClientTypesOption, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllClientTypes, ClientTypesOption(s)) {
		return "", fmt.Errorf("invalid client type: %s", s)
	}
	return ClientTypesOption(s), nil
}

func ToAllowedProviderOption(s string) (AllowedProviderOption, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllAllowedProviderOptions, AllowedProviderOption(s)) {
		return "", fmt.Errorf("invalid allowed provider: %s", s)
	}
	return AllowedProviderOption(s), nil
}

func ToNetworkPolicyEvaluationOption(s string) (NetworkPolicyEvaluationOption, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllNetworkPolicyEvaluationOptions, NetworkPolicyEvaluationOption(s)) {
		return "", fmt.Errorf("invalid network policy evaluation option: %s", s)
	}
	return NetworkPolicyEvaluationOption(s), nil
}

func ToEnforceMfaOnExternalAuthenticationOption(s string) (EnforceMfaOnExternalAuthenticationOption, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllEnforceMfaOnExternalAuthenticationOptions, EnforceMfaOnExternalAuthenticationOption(s)) {
		return "", fmt.Errorf("invalid enforce MFA on external authentication option: %s", s)
	}
	return EnforceMfaOnExternalAuthenticationOption(s), nil
}

func ToMfaPolicyAllowedMethodsOption(s string) (MfaPolicyAllowedMethodsOption, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllMfaPolicyOptions, MfaPolicyAllowedMethodsOption(s)) {
		return "", fmt.Errorf("invalid MFA policy allowed methods option: %s", s)
	}
	return MfaPolicyAllowedMethodsOption(s), nil
}
