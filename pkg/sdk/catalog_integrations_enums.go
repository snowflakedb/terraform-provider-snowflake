package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type CatalogIntegrationCatalogSourceType string

const (
	CatalogIntegrationCatalogSourceTypeAWSGlue              CatalogIntegrationCatalogSourceType = "GLUE"
	CatalogIntegrationCatalogSourceTypeObjectStorage        CatalogIntegrationCatalogSourceType = "OBJECT_STORE"
	CatalogIntegrationCatalogSourceTypePolaris              CatalogIntegrationCatalogSourceType = "POLARIS"
	CatalogIntegrationCatalogSourceTypeIcebergREST          CatalogIntegrationCatalogSourceType = "ICEBERG_REST"
	CatalogIntegrationCatalogSourceTypeSAPBusinessDataCloud CatalogIntegrationCatalogSourceType = "SAP_BDC"
)

var AllCatalogIntegrationCatalogSourceTypes = []CatalogIntegrationCatalogSourceType{
	CatalogIntegrationCatalogSourceTypeAWSGlue,
	CatalogIntegrationCatalogSourceTypeObjectStorage,
	CatalogIntegrationCatalogSourceTypePolaris,
	CatalogIntegrationCatalogSourceTypeIcebergREST,
	CatalogIntegrationCatalogSourceTypeSAPBusinessDataCloud,
}

func ToCatalogIntegrationCatalogSourceType(s string) (CatalogIntegrationCatalogSourceType, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllCatalogIntegrationCatalogSourceTypes, CatalogIntegrationCatalogSourceType(s)) {
		return "", fmt.Errorf("invalid catalog source type: %s", s)
	}
	return CatalogIntegrationCatalogSourceType(s), nil
}

type CatalogIntegrationTableFormat string

const (
	CatalogIntegrationTableFormatIceberg CatalogIntegrationTableFormat = "ICEBERG"
	CatalogIntegrationTableFormatDelta   CatalogIntegrationTableFormat = "DELTA"
)

var AllCatalogIntegrationTableFormats = []CatalogIntegrationTableFormat{
	CatalogIntegrationTableFormatIceberg,
	CatalogIntegrationTableFormatDelta,
}

func ToCatalogIntegrationTableFormat(s string) (CatalogIntegrationTableFormat, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllCatalogIntegrationTableFormats, CatalogIntegrationTableFormat(s)) {
		return "", fmt.Errorf("invalid table format: %s", s)
	}
	return CatalogIntegrationTableFormat(s), nil
}

type CatalogIntegrationRestAuthenticationType string

const (
	CatalogIntegrationRestAuthenticationTypeOAuth  CatalogIntegrationRestAuthenticationType = "OAUTH"
	CatalogIntegrationRestAuthenticationTypeBearer CatalogIntegrationRestAuthenticationType = "BEARER"
	CatalogIntegrationRestAuthenticationTypeSigV4  CatalogIntegrationRestAuthenticationType = "SIGV4"
)

var AllCatalogIntegrationRestAuthenticationTypes = []CatalogIntegrationRestAuthenticationType{
	CatalogIntegrationRestAuthenticationTypeOAuth,
	CatalogIntegrationRestAuthenticationTypeBearer,
	CatalogIntegrationRestAuthenticationTypeSigV4,
}

func ToCatalogIntegrationRestAuthenticationType(s string) (CatalogIntegrationRestAuthenticationType, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllCatalogIntegrationRestAuthenticationTypes, CatalogIntegrationRestAuthenticationType(s)) {
		return "", fmt.Errorf("invalid rest authentication type: %s", s)
	}
	return CatalogIntegrationRestAuthenticationType(s), nil
}

type CatalogIntegrationAccessDelegationMode string

const (
	CatalogIntegrationAccessDelegationModeVendedCredentials         CatalogIntegrationAccessDelegationMode = "VENDED_CREDENTIALS"
	CatalogIntegrationAccessDelegationModeExternalVolumeCredentials CatalogIntegrationAccessDelegationMode = "EXTERNAL_VOLUME_CREDENTIALS"
)

var AllCatalogIntegrationAccessDelegationModes = []CatalogIntegrationAccessDelegationMode{
	CatalogIntegrationAccessDelegationModeVendedCredentials,
	CatalogIntegrationAccessDelegationModeExternalVolumeCredentials,
}

func ToCatalogIntegrationAccessDelegationMode(s string) (CatalogIntegrationAccessDelegationMode, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllCatalogIntegrationAccessDelegationModes, CatalogIntegrationAccessDelegationMode(s)) {
		return "", fmt.Errorf("invalid access delegation mode: %s", s)
	}
	return CatalogIntegrationAccessDelegationMode(s), nil
}

type CatalogIntegrationCatalogApiType string

const (
	CatalogIntegrationCatalogApiTypePublic               CatalogIntegrationCatalogApiType = "PUBLIC"
	CatalogIntegrationCatalogApiTypePrivate              CatalogIntegrationCatalogApiType = "PRIVATE"
	CatalogIntegrationCatalogApiTypeAwsApiGateway        CatalogIntegrationCatalogApiType = "AWS_API_GATEWAY"
	CatalogIntegrationCatalogApiTypeAwsPrivateApiGateway CatalogIntegrationCatalogApiType = "AWS_PRIVATE_API_GATEWAY"
	CatalogIntegrationCatalogApiTypeAwsGlue              CatalogIntegrationCatalogApiType = "AWS_GLUE"
	CatalogIntegrationCatalogApiTypeAwsPrivateGlue       CatalogIntegrationCatalogApiType = "AWS_PRIVATE_GLAUE"
)

var AllCatalogIntegrationCatalogApiTypes = []CatalogIntegrationCatalogApiType{
	CatalogIntegrationCatalogApiTypePublic,
	CatalogIntegrationCatalogApiTypePrivate,
	CatalogIntegrationCatalogApiTypeAwsApiGateway,
	CatalogIntegrationCatalogApiTypeAwsPrivateApiGateway,
	CatalogIntegrationCatalogApiTypeAwsGlue,
	CatalogIntegrationCatalogApiTypeAwsPrivateGlue,
}

func ToCatalogIntegrationCatalogApiType(s string) (CatalogIntegrationCatalogApiType, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllCatalogIntegrationCatalogApiTypes, CatalogIntegrationCatalogApiType(s)) {
		return "", fmt.Errorf("invalid catalog api type: %s", s)
	}
	return CatalogIntegrationCatalogApiType(s), nil
}
