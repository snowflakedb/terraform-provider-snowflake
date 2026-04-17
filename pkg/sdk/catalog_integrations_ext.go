package sdk

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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

func (v *catalogIntegrations) DescribeDetails(ctx context.Context, id AccountObjectIdentifier) (*CatalogIntegrationAllDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseAllCatalogIntegrationProperties(properties, id)
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
			if val, err := strconv.ParseInt(prop.Value, 10, 32); err != nil {
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

type awsGlueSpecificDetails struct {
	glueAwsRoleArn   string
	glueCatalogId    string
	glueRegion       string
	catalogNamespace string
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
	awsGlueDetails := parseAwsGlueSpecificProperties(properties)
	details.GlueAwsRoleArn = awsGlueDetails.glueAwsRoleArn
	details.GlueCatalogId = awsGlueDetails.glueCatalogId
	details.GlueRegion = awsGlueDetails.glueRegion
	details.CatalogNamespace = awsGlueDetails.catalogNamespace
	return details, nil
}

func parseAwsGlueSpecificProperties(properties []CatalogIntegrationProperty) *awsGlueSpecificDetails {
	details := &awsGlueSpecificDetails{}
	for _, prop := range properties {
		switch prop.Name {
		case "GLUE_AWS_ROLE_ARN":
			details.glueAwsRoleArn = prop.Value
		case "GLUE_CATALOG_ID":
			details.glueCatalogId = prop.Value
		case "GLUE_REGION":
			details.glueRegion = prop.Value
		case "CATALOG_NAMESPACE":
			details.catalogNamespace = prop.Value
		}
	}
	return details
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

func parseAllCatalogIntegrationProperties(properties []CatalogIntegrationProperty, id AccountObjectIdentifier) (*CatalogIntegrationAllDetails, error) {
	commons, err := parseCommonProperties(properties)
	if err != nil {
		return nil, err
	}
	details := &CatalogIntegrationAllDetails{
		Id:                     id,
		CatalogSource:          commons.CatalogSource,
		TableFormat:            commons.TableFormat,
		Enabled:                commons.Enabled,
		RefreshIntervalSeconds: commons.RefreshIntervalSeconds,
		Comment:                commons.Comment,
	}

	awsGlueDetails := parseAwsGlueSpecificProperties(properties)
	details.GlueAwsRoleArn = awsGlueDetails.glueAwsRoleArn
	details.GlueCatalogId = awsGlueDetails.glueCatalogId
	details.GlueRegion = awsGlueDetails.glueRegion
	details.CatalogNamespace = awsGlueDetails.catalogNamespace

	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "REST_CONFIG":
			if restConfig, err := parseIcebergRestRestConfigProperty(prop); err != nil {
				errs = append(errs, err)
			} else {
				details.RestConfig = &restConfig
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

func parseOpenCatalogRestConfigProperty(property CatalogIntegrationProperty) (OpenCatalogRestConfigDetails, error) {
	restConfig := OpenCatalogRestConfigDetails{}
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
				restConfig.CatalogApiType = catalogApiType
			}
		case "CATALOG_NAME":
			restConfig.CatalogName = v
		case "ACCESS_DELEGATION_MODE":
			if accessDelegationMode, err := ToCatalogIntegrationAccessDelegationMode(v); err != nil {
				errs = append(errs, err)
			} else {
				restConfig.AccessDelegationMode = accessDelegationMode
			}
		}
	}
	return restConfig, errors.Join(errs...)
}

func parseIcebergRestRestConfigProperty(property CatalogIntegrationProperty) (IcebergRestRestConfigDetails, error) {
	restConfig := IcebergRestRestConfigDetails{}
	var errs []error
	parts := parseCommaSeparatedEnumMap(property)
	for _, part := range parts {
		k, v, _ := strings.Cut(part, "=")
		switch k {
		case "CATALOG_URI":
			restConfig.CatalogUri = v
		case "PREFIX":
			restConfig.Prefix = emptyIfNull(v)
		case "CATALOG_API_TYPE":
			if catalogApiType, err := ToCatalogIntegrationCatalogApiType(v); err != nil {
				errs = append(errs, err)
			} else {
				restConfig.CatalogApiType = catalogApiType
			}
		case "CATALOG_NAME":
			restConfig.CatalogName = emptyIfNull(v)
		case "ACCESS_DELEGATION_MODE":
			if accessDelegationMode, err := ToCatalogIntegrationAccessDelegationMode(v); err != nil {
				errs = append(errs, err)
			} else {
				restConfig.AccessDelegationMode = accessDelegationMode
			}
		}
	}
	return restConfig, errors.Join(errs...)
}

func parseRestAuthenticationProperty(property CatalogIntegrationProperty) (*OAuthRestAuthenticationDetails, *BearerRestAuthenticationDetails, *SigV4RestAuthenticationDetails, error) {
	var oAuthRestAuthentication *OAuthRestAuthenticationDetails
	var bearerRestAuthentication *BearerRestAuthenticationDetails
	var sigV4RestAuthentication *SigV4RestAuthenticationDetails
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

func parseOAuthRestAuthenticationProperty(parts []string) (*OAuthRestAuthenticationDetails, error) {
	restAuthentication := &OAuthRestAuthenticationDetails{}
	for _, part := range parts {
		k, v, _ := strings.Cut(part, "=")
		switch k {
		case "OAUTH_TOKEN_URI":
			// OAUTH_TOKEN_URI is always returned, even if unset
			restAuthentication.OauthTokenUri = v
		case "OAUTH_CLIENT_ID":
			restAuthentication.OauthClientId = v
		case "OAUTH_ALLOWED_SCOPES":
			restAuthentication.OauthAllowedScopes = ParseCommaSeparatedStringArray(v, false)
		}
		// OAUTH_CLIENT_SECRET not returned
	}
	return restAuthentication, nil
}

func parseBearerRestAuthenticationProperty() (*BearerRestAuthenticationDetails, error) {
	restAuthentication := &BearerRestAuthenticationDetails{}
	return restAuthentication, nil
}

func parseSigV4RestAuthenticationProperty(parts []string) (*SigV4RestAuthenticationDetails, error) {
	restAuthentication := &SigV4RestAuthenticationDetails{}
	for _, part := range parts {
		k, v, _ := strings.Cut(part, "=")
		switch k {
		case "SIGV4_IAM_ROLE":
			restAuthentication.Sigv4IamRole = v
		case "SIGV4_SIGNING_REGION":
			// SIGV4_SIGNING_REGION is always returned, even if unset
			restAuthentication.Sigv4SigningRegion = v
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
