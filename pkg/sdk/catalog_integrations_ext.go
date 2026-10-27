package sdk

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (r *CreateAwsGlueCatalogIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *CreateIcebergRestCatalogIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *CreateObjectStorageCatalogIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *CreateOpenCatalogCatalogIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (d *CatalogIntegrationOpenCatalogDetails) ID() AccountObjectIdentifier {
	return d.Id
}

func (d *CatalogIntegrationIcebergRestDetails) ID() AccountObjectIdentifier {
	return d.Id
}

func (d *CatalogIntegrationAwsGlueDetails) ID() AccountObjectIdentifier {
	return d.Id
}

func (d *CatalogIntegrationObjectStorageDetails) ID() AccountObjectIdentifier {
	return d.Id
}

func (d *CatalogIntegrationDetails) ID() AccountObjectIdentifier {
	return d.Id
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
		switch prop.Property {
		case "CATALOG_SOURCE":
			if catalogSource, err := ToCatalogIntegrationCatalogSourceType(prop.PropertyValue); err != nil {
				errs = append(errs, err)
			} else {
				commons.CatalogSource = catalogSource
			}
		case "TABLE_FORMAT":
			if tableFormat, err := ToCatalogIntegrationTableFormat(prop.PropertyValue); err != nil {
				errs = append(errs, err)
			} else {
				commons.TableFormat = tableFormat
			}
		case "ENABLED":
			if val, err := strconv.ParseBool(prop.PropertyValue); err != nil {
				errs = append(errs, err)
			} else {
				commons.Enabled = val
			}
		case "REFRESH_INTERVAL_SECONDS":
			if val, err := strconv.ParseInt(prop.PropertyValue, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				commons.RefreshIntervalSeconds = int(val)
			}
		case "COMMENT":
			commons.Comment = prop.PropertyValue
		}
	}
	return commons, errors.Join(errs...)
}

type awsGlueSpecificDetails struct {
	glueAwsRoleArn    string
	glueCatalogId     string
	glueRegion        string
	catalogNamespace  string
	glueAwsIamUserArn string
	glueAwsExternalId string
}

func ParseCatalogIntegrationAwsGlueDescribed(properties []CatalogIntegrationProperty) (CatalogIntegrationAwsGlueDetails, error) {
	commons, err := parseCommonProperties(properties)
	if err != nil {
		return CatalogIntegrationAwsGlueDetails{}, err
	}
	details := CatalogIntegrationAwsGlueDetails{
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
	details.GlueAwsIamUserArn = awsGlueDetails.glueAwsIamUserArn
	details.GlueAwsExternalId = awsGlueDetails.glueAwsExternalId
	return details, nil
}

func parseAwsGlueSpecificProperties(properties []CatalogIntegrationProperty) *awsGlueSpecificDetails {
	details := &awsGlueSpecificDetails{}
	for _, prop := range properties {
		switch prop.Property {
		case "GLUE_AWS_ROLE_ARN":
			details.glueAwsRoleArn = prop.PropertyValue
		case "GLUE_CATALOG_ID":
			details.glueCatalogId = prop.PropertyValue
		case "GLUE_REGION":
			details.glueRegion = prop.PropertyValue
		case "CATALOG_NAMESPACE":
			details.catalogNamespace = prop.PropertyValue
		case "GLUE_AWS_IAM_USER_ARN":
			details.glueAwsIamUserArn = prop.PropertyValue
		case "GLUE_AWS_EXTERNAL_ID":
			details.glueAwsExternalId = prop.PropertyValue
		}
	}
	return details
}

func ParseCatalogIntegrationObjectStorageDescribed(properties []CatalogIntegrationProperty) (CatalogIntegrationObjectStorageDetails, error) {
	commons, err := parseCommonProperties(properties)
	if err != nil {
		return CatalogIntegrationObjectStorageDetails{}, err
	}
	details := CatalogIntegrationObjectStorageDetails{
		CatalogSource:          commons.CatalogSource,
		TableFormat:            commons.TableFormat,
		Enabled:                commons.Enabled,
		RefreshIntervalSeconds: commons.RefreshIntervalSeconds,
		Comment:                commons.Comment,
	}
	return details, nil
}

func ParseCatalogIntegrationOpenCatalogDescribed(properties []CatalogIntegrationProperty) (CatalogIntegrationOpenCatalogDetails, error) {
	commons, err := parseCommonProperties(properties)
	if err != nil {
		return CatalogIntegrationOpenCatalogDetails{}, err
	}
	details := CatalogIntegrationOpenCatalogDetails{
		CatalogSource:          commons.CatalogSource,
		TableFormat:            commons.TableFormat,
		Enabled:                commons.Enabled,
		RefreshIntervalSeconds: commons.RefreshIntervalSeconds,
		Comment:                commons.Comment,
	}
	var errs []error
	for _, prop := range properties {
		switch prop.Property {
		case "CATALOG_NAMESPACE":
			details.CatalogNamespace = prop.PropertyValue
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
				details.RestAuthentication = CatalogIntegrationOpenCatalogRestAuthenticationDetails(*oAuthRestAuth)
			}
		}
	}
	return details, errors.Join(errs...)
}

func ParseCatalogIntegrationIcebergRestDescribed(properties []CatalogIntegrationProperty) (CatalogIntegrationIcebergRestDetails, error) {
	commons, err := parseCommonProperties(properties)
	if err != nil {
		return CatalogIntegrationIcebergRestDetails{}, err
	}
	details := CatalogIntegrationIcebergRestDetails{
		CatalogSource:          commons.CatalogSource,
		TableFormat:            commons.TableFormat,
		Enabled:                commons.Enabled,
		RefreshIntervalSeconds: commons.RefreshIntervalSeconds,
		Comment:                commons.Comment,
	}
	var errs []error
	for _, prop := range properties {
		switch prop.Property {
		case "CATALOG_NAMESPACE":
			details.CatalogNamespace = prop.PropertyValue
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
				details.OauthRestAuthentication = oAuthRestAuth
				details.BearerRestAuthentication = bearerRestAuth
				details.Sigv4RestAuthentication = sigV4RestAuth
			}
		}
	}
	return details, errors.Join(errs...)
}

func ParseCatalogIntegrationDescribed(properties []CatalogIntegrationProperty) (CatalogIntegrationDetails, error) {
	commons, err := parseCommonProperties(properties)
	if err != nil {
		return CatalogIntegrationDetails{}, err
	}
	details := CatalogIntegrationDetails{
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
	details.GlueAwsIamUserArn = awsGlueDetails.glueAwsIamUserArn
	details.GlueAwsExternalId = awsGlueDetails.glueAwsExternalId

	var errs []error
	for _, prop := range properties {
		switch prop.Property {
		case "REST_CONFIG":
			// IcebergRestRestConfigDetails contains all properties of OpenCatalogRestConfigDetails
			if restConfig, err := parseIcebergRestRestConfigProperty(prop); err != nil {
				errs = append(errs, err)
			} else {
				converted := CatalogIntegrationRestConfigDetails(restConfig)
				details.RestConfig = &converted
			}
		case "REST_AUTHENTICATION":
			if oAuthRestAuth, bearerRestAuth, sigV4RestAuth, err := parseRestAuthenticationProperty(prop); err != nil {
				errs = append(errs, err)
			} else {
				if oAuthRestAuth != nil {
					converted := CatalogIntegrationOauthRestAuthenticationDetails(*oAuthRestAuth)
					details.OauthRestAuthentication = &converted
				}
				if bearerRestAuth != nil {
					converted := CatalogIntegrationBearerRestAuthenticationDetails(*bearerRestAuth)
					details.BearerRestAuthentication = &converted
				}
				if sigV4RestAuth != nil {
					converted := CatalogIntegrationSigv4RestAuthenticationDetails(*sigV4RestAuth)
					details.Sigv4RestAuthentication = &converted
				}
			}
		}
	}
	return details, errors.Join(errs...)
}

func parseOpenCatalogRestConfigProperty(property CatalogIntegrationProperty) (CatalogIntegrationOpenCatalogRestConfigDetails, error) {
	restConfig := CatalogIntegrationOpenCatalogRestConfigDetails{}
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

func parseIcebergRestRestConfigProperty(property CatalogIntegrationProperty) (CatalogIntegrationIcebergRestRestConfigDetails, error) {
	restConfig := CatalogIntegrationIcebergRestRestConfigDetails{}
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

func parseRestAuthenticationProperty(property CatalogIntegrationProperty) (*CatalogIntegrationIcebergRestOauthRestAuthenticationDetails, *CatalogIntegrationIcebergRestBearerRestAuthenticationDetails, *CatalogIntegrationIcebergRestSigv4RestAuthenticationDetails, error) {
	var oAuthRestAuthentication *CatalogIntegrationIcebergRestOauthRestAuthenticationDetails
	var bearerRestAuthentication *CatalogIntegrationIcebergRestBearerRestAuthenticationDetails
	var sigV4RestAuthentication *CatalogIntegrationIcebergRestSigv4RestAuthenticationDetails
	var errs []error
	parts := parseCommaSeparatedEnumMap(property)
	for _, part := range parts {
		k, v, _ := strings.Cut(part, "=")
		if k == "TYPE" {
			switch v {
			case string(CatalogIntegrationRestAuthenticationTypeOauth):
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
			case string(CatalogIntegrationRestAuthenticationTypeSigv4):
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

func parseOAuthRestAuthenticationProperty(parts []string) (*CatalogIntegrationIcebergRestOauthRestAuthenticationDetails, error) {
	restAuthentication := &CatalogIntegrationIcebergRestOauthRestAuthenticationDetails{}
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

func parseBearerRestAuthenticationProperty() (*CatalogIntegrationIcebergRestBearerRestAuthenticationDetails, error) {
	restAuthentication := &CatalogIntegrationIcebergRestBearerRestAuthenticationDetails{}
	return restAuthentication, nil
}

func parseSigV4RestAuthenticationProperty(parts []string) (*CatalogIntegrationIcebergRestSigv4RestAuthenticationDetails, error) {
	restAuthentication := &CatalogIntegrationIcebergRestSigv4RestAuthenticationDetails{}
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
	s := strings.TrimPrefix(property.PropertyValue, "{")
	s = strings.TrimSuffix(s, "}")
	return ParseOuterCommaSeparatedStringArray(fmt.Sprintf("[%s]", s), false)
}
