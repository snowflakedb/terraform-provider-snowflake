package sdk

import (
	"fmt"
	"strings"
)

type (
	SecretType string
)

func ToSecretType(s string) (SecretType, error) {
	switch strings.ToUpper(s) {
	case string(SecretTypePassword):
		return SecretTypePassword, nil
	case string(SecretTypeOAuth2):
		return SecretTypeOAuth2, nil
	case string(SecretTypeGenericString):
		return SecretTypeGenericString, nil
	case string(SecretTypeOAuth2ClientCredentials):
		return SecretTypeOAuth2ClientCredentials, nil
	case string(SecretTypeOAuth2AuthorizationCodeGrant):
		return SecretTypeOAuth2AuthorizationCodeGrant, nil
	default:
		return "", fmt.Errorf("invalid secret type: %s", s)
	}
}

const (
	SecretTypePassword                     SecretType = "PASSWORD"
	SecretTypeOAuth2                       SecretType = "OAUTH2"
	SecretTypeGenericString                SecretType = "GENERIC_STRING"
	SecretTypeOAuth2ClientCredentials      SecretType = "OAUTH2_CLIENT_CREDENTIALS"       // #nosec G101
	SecretTypeOAuth2AuthorizationCodeGrant SecretType = "OAUTH2_AUTHORIZATION_CODE_GRANT" // #nosec G101
)

var AcceptableSecretTypes = map[SecretType][]SecretType{
	SecretTypePassword:                     {SecretTypePassword},
	SecretTypeOAuth2:                       {SecretTypeOAuth2},
	SecretTypeGenericString:                {SecretTypeGenericString},
	SecretTypeOAuth2ClientCredentials:      {SecretTypeOAuth2ClientCredentials, SecretTypeOAuth2},
	SecretTypeOAuth2AuthorizationCodeGrant: {SecretTypeOAuth2ClientCredentials, SecretTypeOAuth2},
}
