package sdk

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

func (r *CreateCatalogIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (v *catalogIntegrations) DescribeAwsGlueParams(ctx context.Context, id AccountObjectIdentifier) (*AwsGlueParams, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseAwsGlueProperties(properties)
}

func (v *catalogIntegrations) DescribeObjectStorageParams(ctx context.Context, id AccountObjectIdentifier) (*ObjectStorageParams, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseObjectStorageProperties(properties)
}

func (v *catalogIntegrations) DescribeOpenCatalogParams(ctx context.Context, id AccountObjectIdentifier) (*OpenCatalogParams, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseOpenCatalogProperties(properties)
}

func (v *catalogIntegrations) DescribeIcebergRestParams(ctx context.Context, id AccountObjectIdentifier) (*IcebergRestParams, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseIcebergRestProperties(properties)
}

func (v *catalogIntegrations) DescribeSapBdcParams(ctx context.Context, id AccountObjectIdentifier) (*SapBdcParams, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseSapBdcProperties(properties)
}

func parseAwsGlueProperties(properties []CatalogIntegrationProperty) (*AwsGlueParams, error) {
	params := &AwsGlueParams{}
	for _, prop := range properties {
		switch prop.Name {
		case "GLUE_AWS_ROLE_ARN":
			params.GlueAwsRoleArn = prop.Value
		case "GLUE_CATALOG_ID":
			params.GlueCatalogId = prop.Value
		case "GLUE_REGION":
			params.GlueRegion = String(prop.Value)
		case "CATALOG_NAMESPACE":
			params.CatalogNamespace = String(prop.Value)
		}
	}
	return params, nil
}

func parseObjectStorageProperties(properties []CatalogIntegrationProperty) (*ObjectStorageParams, error) {
	params := &ObjectStorageParams{}
	for _, prop := range properties {
		if prop.Name == "TABLE_FORMAT" {
			tableFormat, err := ToCatalogIntegrationTableFormat(prop.Value)
			if err != nil {
				return nil, err
			}
			params.TableFormat = tableFormat
		}
	}
	return params, nil
}

func parseOpenCatalogProperties(properties []CatalogIntegrationProperty) (*OpenCatalogParams, error) {
	params := &OpenCatalogParams{}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "CATALOG_NAMESPACE":
			params.CatalogNamespace = String(prop.Value)
		case "REST_CONFIG":
			if restConfig, err := parseOpenCatalogRestConfigProperty(prop); err != nil {
				errs = append(errs, err)
			} else {
				params.RestConfig = restConfig
			}
		case "REST_AUTHENTICATION":
			if oAuthRestAuth, _, _, err := parseRestAuthenticationProperty(prop); err != nil {
				errs = append(errs, err)
			} else if oAuthRestAuth == nil {
				errs = append(errs, errors.New("REST_AUTHENTICATION property is not of OAUTH type"))
			} else {
				params.RestAuthentication = *oAuthRestAuth
			}
		}
	}
	return params, errors.Join(errs...)
}

func parseIcebergRestProperties(properties []CatalogIntegrationProperty) (*IcebergRestParams, error) {
	params := &IcebergRestParams{}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "CATALOG_NAMESPACE":
			params.CatalogNamespace = String(prop.Value)
		case "REST_CONFIG":
			if restConfig, err := parseIcebergRestRestConfigProperty(prop); err != nil {
				errs = append(errs, err)
			} else {
				params.RestConfig = restConfig
			}
		case "REST_AUTHENTICATION":
			if oAuthRestAuth, bearerRestAuth, sigV4RestAuth, err := parseRestAuthenticationProperty(prop); err != nil {
				errs = append(errs, err)
			} else {
				params.OAuthRestAuthentication = oAuthRestAuth
				params.BearerRestAuthentication = bearerRestAuth
				params.SigV4RestAuthentication = sigV4RestAuth
			}
		}
	}
	return params, errors.Join(errs...)
}

func parseSapBdcProperties(properties []CatalogIntegrationProperty) (*SapBdcParams, error) {
	params := &SapBdcParams{RestConfig: SapBdcRestConfig{}}
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
				if restAuth, err := parseBearerRestAuthenticationProperty(parts); err != nil {
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

func parseBearerRestAuthenticationProperty(parts []string) (*BearerRestAuthentication, error) {
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
