package sdk

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

type CatalogIntegrationAwsGlueDetails struct {
	Id                     AccountObjectIdentifier
	CatalogSource          CatalogIntegrationCatalogSourceType
	TableFormat            CatalogIntegrationTableFormat
	Enabled                bool
	RefreshIntervalSeconds int
	Comment                string
	GlueAwsRoleArn         string
	GlueCatalogId          string
	GlueRegion             string
	CatalogNamespace       string
}

type CatalogIntegrationObjectStorageDetails struct {
	Id                     AccountObjectIdentifier
	CatalogSource          CatalogIntegrationCatalogSourceType
	TableFormat            CatalogIntegrationTableFormat
	Enabled                bool
	RefreshIntervalSeconds int
	Comment                string
}

type CatalogIntegrationOpenCatalogDetails struct {
	Id                     AccountObjectIdentifier
	CatalogSource          CatalogIntegrationCatalogSourceType
	TableFormat            CatalogIntegrationTableFormat
	Enabled                bool
	RefreshIntervalSeconds int
	Comment                string
	CatalogNamespace       string
	RestConfig             OpenCatalogRestConfig
	RestAuthentication     OAuthRestAuthentication
}

type CatalogIntegrationIcebergRestDetails struct {
	Id                       AccountObjectIdentifier
	CatalogSource            CatalogIntegrationCatalogSourceType
	TableFormat              CatalogIntegrationTableFormat
	Enabled                  bool
	RefreshIntervalSeconds   int
	Comment                  string
	CatalogNamespace         string
	RestConfig               IcebergRestRestConfig
	OAuthRestAuthentication  *OAuthRestAuthentication
	BearerRestAuthentication *BearerRestAuthentication
	SigV4RestAuthentication  *SigV4RestAuthentication
}

type CatalogIntegrationSapBdcDetails struct {
	Id                     AccountObjectIdentifier
	CatalogSource          CatalogIntegrationCatalogSourceType
	TableFormat            CatalogIntegrationTableFormat
	Enabled                bool
	RefreshIntervalSeconds int
	Comment                string
}

func (r *CreateCatalogIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (v *catalogIntegrations) DescribeAwsGlueDetails(ctx context.Context, id AccountObjectIdentifier) (*CatalogIntegrationAwsGlueDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseAwsGlueProperties(properties, id)
}

func (v *catalogIntegrations) DescribeObjectStorageDetails(ctx context.Context, id AccountObjectIdentifier) (*CatalogIntegrationObjectStorageDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseObjectStorageProperties(properties, id)
}

func (v *catalogIntegrations) DescribeOpenCatalogDetails(ctx context.Context, id AccountObjectIdentifier) (*CatalogIntegrationOpenCatalogDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseOpenCatalogProperties(properties, id)
}

func (v *catalogIntegrations) DescribeIcebergRestDetails(ctx context.Context, id AccountObjectIdentifier) (*CatalogIntegrationIcebergRestDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseIcebergRestProperties(properties, id)
}

func (v *catalogIntegrations) DescribeSapBdcDetails(ctx context.Context, id AccountObjectIdentifier) (*CatalogIntegrationSapBdcDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseSapBdcProperties(properties, id)
}

type commonDetails struct {
	CatalogSource          CatalogIntegrationCatalogSourceType
	TableFormat            CatalogIntegrationTableFormat
	Enabled                bool
	RefreshIntervalSeconds int
	Comment                string
}

func parseCommonProperties(properties []CatalogIntegrationProperty) (*commonDetails, error) {
	commons := &commonDetails{}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "CATALOG_SOURCE":
			if catalogSource, err := ToCatalogIntegrationCatalogSourceType(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				commons.CatalogSource = catalogSource
			}
		case "TABLE_FORMAT":
			if tableFormat, err := ToCatalogIntegrationTableFormat(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				commons.TableFormat = tableFormat
			}
		case "ENABLED":
			if val, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				commons.Enabled = val
			}
		case "REFRESH_INTERVAL_SECONDS":
			if val, err := strconv.ParseInt(prop.Value, 10, 64); err != nil {
				errs = append(errs, err)
			} else {
				commons.RefreshIntervalSeconds = int(val)
			}
		case "COMMENT":
			commons.Comment = prop.Value
		}
	}
	return commons, errors.Join(errs...)
}

func parseAwsGlueProperties(properties []CatalogIntegrationProperty, id AccountObjectIdentifier) (*CatalogIntegrationAwsGlueDetails, error) {
	commons, err := parseCommonProperties(properties)
	if err != nil {
		return nil, err
	}
	details := &CatalogIntegrationAwsGlueDetails{
		Id:                     id,
		CatalogSource:          commons.CatalogSource,
		TableFormat:            commons.TableFormat,
		Enabled:                commons.Enabled,
		RefreshIntervalSeconds: commons.RefreshIntervalSeconds,
		Comment:                commons.Comment,
	}
	for _, prop := range properties {
		switch prop.Name {
		case "GLUE_AWS_ROLE_ARN":
			details.GlueAwsRoleArn = prop.Value
		case "GLUE_CATALOG_ID":
			details.GlueCatalogId = prop.Value
		case "GLUE_REGION":
			details.GlueRegion = prop.Value
		case "CATALOG_NAMESPACE":
			details.CatalogNamespace = prop.Value
		}
	}
	return details, nil
}

func parseObjectStorageProperties(properties []CatalogIntegrationProperty, id AccountObjectIdentifier) (*CatalogIntegrationObjectStorageDetails, error) {
	commons, err := parseCommonProperties(properties)
	if err != nil {
		return nil, err
	}
	details := &CatalogIntegrationObjectStorageDetails{
		Id:                     id,
		CatalogSource:          commons.CatalogSource,
		TableFormat:            commons.TableFormat,
		Enabled:                commons.Enabled,
		RefreshIntervalSeconds: commons.RefreshIntervalSeconds,
		Comment:                commons.Comment,
	}
	return details, nil
}

func parseOpenCatalogProperties(properties []CatalogIntegrationProperty, id AccountObjectIdentifier) (*CatalogIntegrationOpenCatalogDetails, error) {
	commons, err := parseCommonProperties(properties)
	if err != nil {
		return nil, err
	}
	details := &CatalogIntegrationOpenCatalogDetails{
		Id:                     id,
		CatalogSource:          commons.CatalogSource,
		TableFormat:            commons.TableFormat,
		Enabled:                commons.Enabled,
		RefreshIntervalSeconds: commons.RefreshIntervalSeconds,
		Comment:                commons.Comment,
	}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "CATALOG_NAMESPACE":
			details.CatalogNamespace = prop.Value
		case "REST_CONFIG":
			if restConfig, err := parseOpenCatalogRestConfigProperty(prop); err != nil {
				errs = append(errs, err)
			} else {
				details.RestConfig = restConfig
			}
		case "REST_AUTHENTICATION":
			if oAuthRestAuth, _, _, err := parseRestAuthenticationProperty(prop); err != nil {
				errs = append(errs, err)
			} else if oAuthRestAuth == nil {
				errs = append(errs, errors.New("REST_AUTHENTICATION property is not of OAUTH type"))
			} else {
				details.RestAuthentication = *oAuthRestAuth
			}
		}
	}
	return details, errors.Join(errs...)
}

func parseIcebergRestProperties(properties []CatalogIntegrationProperty, id AccountObjectIdentifier) (*CatalogIntegrationIcebergRestDetails, error) {
	commons, err := parseCommonProperties(properties)
	if err != nil {
		return nil, err
	}
	details := &CatalogIntegrationIcebergRestDetails{
		Id:                     id,
		CatalogSource:          commons.CatalogSource,
		TableFormat:            commons.TableFormat,
		Enabled:                commons.Enabled,
		RefreshIntervalSeconds: commons.RefreshIntervalSeconds,
		Comment:                commons.Comment,
	}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "CATALOG_NAMESPACE":
			details.CatalogNamespace = prop.Value
		case "REST_CONFIG":
			if restConfig, err := parseIcebergRestRestConfigProperty(prop); err != nil {
				errs = append(errs, err)
			} else {
				details.RestConfig = restConfig
			}
		case "REST_AUTHENTICATION":
			if oAuthRestAuth, bearerRestAuth, sigV4RestAuth, err := parseRestAuthenticationProperty(prop); err != nil {
				errs = append(errs, err)
			} else {
				details.OAuthRestAuthentication = oAuthRestAuth
				details.BearerRestAuthentication = bearerRestAuth
				details.SigV4RestAuthentication = sigV4RestAuth
			}
		}
	}
	return details, errors.Join(errs...)
}

func parseSapBdcProperties(properties []CatalogIntegrationProperty, id AccountObjectIdentifier) (*CatalogIntegrationSapBdcDetails, error) {
	commons, err := parseCommonProperties(properties)
	if err != nil {
		return nil, err
	}
	params := &CatalogIntegrationSapBdcDetails{
		Id:                     id,
		CatalogSource:          commons.CatalogSource,
		TableFormat:            commons.TableFormat,
		Enabled:                commons.Enabled,
		RefreshIntervalSeconds: commons.RefreshIntervalSeconds,
		Comment:                commons.Comment,
	}
	return params, nil
}

func parseOpenCatalogRestConfigProperty(property CatalogIntegrationProperty) (OpenCatalogRestConfig, error) {
	restConfig := OpenCatalogRestConfig{}
	var errs []error
	parts := parseCommaSeparatedEnumMap(property)
	for _, part := range parts {
		k, v, _ := strings.Cut(part, "=")
		switch k {
		case "CATALOG_URI":
			restConfig.CatalogUri = v
		case "CATALOG_API_TYPE":
			if catalogApiType, err := ToCatalogIntegrationCatalogApiType(v); err != nil {
				errs = append(errs, err)
			} else {
				restConfig.CatalogApiType = Pointer(catalogApiType)
			}
		case "CATALOG_NAME":
			restConfig.CatalogName = v
		case "ACCESS_DELEGATION_MODE":
			if accessDelegationMode, err := ToCatalogIntegrationAccessDelegationMode(v); err != nil {
				errs = append(errs, err)
			} else {
				restConfig.AccessDelegationMode = Pointer(accessDelegationMode)
			}
		}
	}
	return restConfig, errors.Join(errs...)
}

func parseIcebergRestRestConfigProperty(property CatalogIntegrationProperty) (IcebergRestRestConfig, error) {
	restConfig := IcebergRestRestConfig{}
	var errs []error
	parts := parseCommaSeparatedEnumMap(property)
	for _, part := range parts {
		k, v, _ := strings.Cut(part, "=")
		switch k {
		case "CATALOG_URI":
			restConfig.CatalogUri = v
		case "PREFIX":
			restConfig.Prefix = String(emptyIfNull(v))
		case "CATALOG_API_TYPE":
			if catalogApiType, err := ToCatalogIntegrationCatalogApiType(v); err != nil {
				errs = append(errs, err)
			} else {
				restConfig.CatalogApiType = Pointer(catalogApiType)
			}
		case "CATALOG_NAME":
			restConfig.CatalogName = String(emptyIfNull(v))
		case "ACCESS_DELEGATION_MODE":
			if accessDelegationMode, err := ToCatalogIntegrationAccessDelegationMode(v); err != nil {
				errs = append(errs, err)
			} else {
				restConfig.AccessDelegationMode = Pointer(accessDelegationMode)
			}
		}
	}
	return restConfig, errors.Join(errs...)
}

func parseRestAuthenticationProperty(property CatalogIntegrationProperty) (*OAuthRestAuthentication, *BearerRestAuthentication, *SigV4RestAuthentication, error) {
	var oAuthRestAuthentication *OAuthRestAuthentication
	var bearerRestAuthentication *BearerRestAuthentication
	var sigV4RestAuthentication *SigV4RestAuthentication
	var errs []error
	parts := parseCommaSeparatedEnumMap(property)
	for _, part := range parts {
		k, v, _ := strings.Cut(part, "=")
		if k == "TYPE" {
			switch v {
			case string(CatalogIntegrationRestAuthenticationTypeOAuth):
				if restAuth, err := parseOAuthRestAuthenticationProperty(parts); err != nil {
					errs = append(errs, err)
				} else {
					oAuthRestAuthentication = restAuth
				}
			case string(CatalogIntegrationRestAuthenticationTypeBearer):
				if restAuth, err := parseBearerRestAuthenticationProperty(); err != nil {
					errs = append(errs, err)
				} else {
					bearerRestAuthentication = restAuth
				}
			case string(CatalogIntegrationRestAuthenticationTypeSigV4):
				if restAuth, err := parseSigV4RestAuthenticationProperty(parts); err != nil {
					errs = append(errs, err)
				} else {
					sigV4RestAuthentication = restAuth
				}
			}
		}
	}
	return oAuthRestAuthentication, bearerRestAuthentication, sigV4RestAuthentication, errors.Join(errs...)
}

func parseOAuthRestAuthenticationProperty(parts []string) (*OAuthRestAuthentication, error) {
	restAuthentication := &OAuthRestAuthentication{}
	for _, part := range parts {
		k, v, _ := strings.Cut(part, "=")
		switch k {
		case "OAUTH_TOKEN_URI":
			// OAUTH_TOKEN_URI is always returned, even if unset
			restAuthentication.OauthTokenUri = String(v)
		case "OAUTH_CLIENT_ID":
			restAuthentication.OauthClientId = v
		case "OAUTH_ALLOWED_SCOPES":
			restAuthentication.OauthAllowedScopes = collections.Map(ParseCommaSeparatedStringArray(v, false), func(s string) StringListItemWrapper { return StringListItemWrapper{s} })
		}
		// OAUTH_CLIENT_SECRET not returned
	}
	return restAuthentication, nil
}

func parseBearerRestAuthenticationProperty() (*BearerRestAuthentication, error) {
	restAuthentication := &BearerRestAuthentication{}
	return restAuthentication, nil
}

func parseSigV4RestAuthenticationProperty(parts []string) (*SigV4RestAuthentication, error) {
	restAuthentication := &SigV4RestAuthentication{}
	for _, part := range parts {
		k, v, _ := strings.Cut(part, "=")
		switch k {
		case "SIGV4_IAM_ROLE":
			restAuthentication.Sigv4IamRole = v
		case "SIGV4_SIGNING_REGION":
			// SIGV4_SIGNING_REGION is always returned, even if unset
			restAuthentication.Sigv4SigningRegion = String(v)
		}
		// SIGV4_EXTERNAL_ID not returned
	}
	return restAuthentication, nil
}

func parseCommaSeparatedEnumMap(property CatalogIntegrationProperty) []string {
	s := strings.TrimPrefix(property.Value, "{")
	s = strings.TrimSuffix(s, "}")
	return ParseOuterCommaSeparatedStringArray(fmt.Sprintf("[%s]", s), false)
}

func emptyIfNull(s string) string {
	if s == "null" {
		return ""
	}
	return s
}
