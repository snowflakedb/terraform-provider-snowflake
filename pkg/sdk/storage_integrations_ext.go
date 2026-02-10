package sdk

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type S3Protocol string

const (
	RegularS3Protocol S3Protocol = "S3"
	GovS3Protocol     S3Protocol = "S3GOV"
	ChinaS3Protocol   S3Protocol = "S3CHINA"
)

var (
	AllS3Protocols      = []S3Protocol{RegularS3Protocol, GovS3Protocol, ChinaS3Protocol}
	AllStorageProviders = append(AsStringList(AllS3Protocols), "GCS", "AZURE")
)

func ToS3Protocol(s string) (S3Protocol, error) {
	switch protocol := S3Protocol(strings.ToUpper(s)); protocol {
	case RegularS3Protocol, GovS3Protocol, ChinaS3Protocol:
		return protocol, nil
	default:
		return "", fmt.Errorf("invalid S3 protocol: %s", s)
	}
}

func (d *StorageIntegrationAwsDetails) ID() AccountObjectIdentifier {
	return d.Id
}

func (d *StorageIntegrationAzureDetails) ID() AccountObjectIdentifier {
	return d.Id
}

func (d *StorageIntegrationGcsDetails) ID() AccountObjectIdentifier {
	return d.Id
}

func (v *storageIntegrations) DescribeAwsDetails(ctx context.Context, id AccountObjectIdentifier) (*StorageIntegrationAwsDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseAwsProperties(properties, id)
}

func (v *storageIntegrations) DescribeAzureDetails(ctx context.Context, id AccountObjectIdentifier) (*StorageIntegrationAzureDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseAzureProperties(properties, id)
}

func (v *storageIntegrations) DescribeGcsDetails(ctx context.Context, id AccountObjectIdentifier) (*StorageIntegrationGcsDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseGcsProperties(properties, id)
}

func (v *storageIntegrations) DescribeDetails(ctx context.Context, id AccountObjectIdentifier) (*StorageIntegrationAllDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseAllProperties(properties, id)
}

// TODO [next PRs]: extract common mapping logic
func parseAwsProperties(properties []StorageIntegrationProperty, id AccountObjectIdentifier) (*StorageIntegrationAwsDetails, error) {
	details := &StorageIntegrationAwsDetails{
		Id: id,
	}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "ENABLED":
			if val, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.Enabled = val
			}
		case "STORAGE_PROVIDER":
			details.Provider = prop.Value
		case "STORAGE_ALLOWED_LOCATIONS":
			details.AllowedLocations = ParseCommaSeparatedStringArray(prop.Value, false)
		case "STORAGE_BLOCKED_LOCATIONS":
			details.BlockedLocations = ParseCommaSeparatedStringArray(prop.Value, false)
		case "COMMENT":
			details.Comment = prop.Value
		case "USE_PRIVATELINK_ENDPOINT":
			if val, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.UsePrivatelinkEndpoint = val
			}
		case "STORAGE_AWS_IAM_USER_ARN":
			details.IamUserArn = prop.Value
		case "STORAGE_AWS_ROLE_ARN":
			details.RoleArn = prop.Value
		case "STORAGE_AWS_OBJECT_ACL":
			details.ObjectAcl = prop.Value
		case "STORAGE_AWS_EXTERNAL_ID":
			details.ExternalId = prop.Value
		}
	}
	return details, errors.Join(errs...)
}

func parseAzureProperties(properties []StorageIntegrationProperty, id AccountObjectIdentifier) (*StorageIntegrationAzureDetails, error) {
	details := &StorageIntegrationAzureDetails{
		Id: id,
	}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "ENABLED":
			if val, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.Enabled = val
			}
		case "STORAGE_PROVIDER":
			details.Provider = prop.Value
		case "STORAGE_ALLOWED_LOCATIONS":
			details.AllowedLocations = ParseCommaSeparatedStringArray(prop.Value, false)
		case "STORAGE_BLOCKED_LOCATIONS":
			details.BlockedLocations = ParseCommaSeparatedStringArray(prop.Value, false)
		case "COMMENT":
			details.Comment = prop.Value
		case "USE_PRIVATELINK_ENDPOINT":
			if val, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.UsePrivatelinkEndpoint = val
			}
		case "AZURE_TENANT_ID":
			details.TenantId = prop.Value
		case "AZURE_CONSENT_URL":
			details.ConsentUrl = prop.Value
		case "AZURE_MULTI_TENANT_APP_NAME":
			details.MultiTenantAppName = prop.Value
		}
	}
	return details, errors.Join(errs...)
}

func parseGcsProperties(properties []StorageIntegrationProperty, id AccountObjectIdentifier) (*StorageIntegrationGcsDetails, error) {
	details := &StorageIntegrationGcsDetails{
		Id: id,
	}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "ENABLED":
			if val, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.Enabled = val
			}
		case "STORAGE_PROVIDER":
			details.Provider = prop.Value
		case "STORAGE_ALLOWED_LOCATIONS":
			details.AllowedLocations = ParseCommaSeparatedStringArray(prop.Value, false)
		case "STORAGE_BLOCKED_LOCATIONS":
			details.BlockedLocations = ParseCommaSeparatedStringArray(prop.Value, false)
		case "COMMENT":
			details.Comment = prop.Value
		case "USE_PRIVATELINK_ENDPOINT":
			if val, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.UsePrivatelinkEndpoint = val
			}
		case "STORAGE_GCP_SERVICE_ACCOUNT":
			details.ServiceAccount = prop.Value
		}
	}
	return details, errors.Join(errs...)
}

func parseAllProperties(properties []StorageIntegrationProperty, id AccountObjectIdentifier) (*StorageIntegrationAllDetails, error) {
	details := &StorageIntegrationAllDetails{
		Id: id,
	}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "ENABLED":
			if val, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.Enabled = val
			}
		case "STORAGE_PROVIDER":
			details.Provider = prop.Value
		case "STORAGE_ALLOWED_LOCATIONS":
			details.AllowedLocations = ParseCommaSeparatedStringArray(prop.Value, false)
		case "STORAGE_BLOCKED_LOCATIONS":
			details.BlockedLocations = ParseCommaSeparatedStringArray(prop.Value, false)
		case "COMMENT":
			details.Comment = prop.Value
		case "USE_PRIVATELINK_ENDPOINT":
			if val, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.UsePrivatelinkEndpoint = val
			}
		case "STORAGE_AWS_IAM_USER_ARN":
			details.IamUserArn = prop.Value
		case "STORAGE_AWS_ROLE_ARN":
			details.RoleArn = prop.Value
		case "STORAGE_AWS_OBJECT_ACL":
			details.ObjectAcl = prop.Value
		case "STORAGE_AWS_EXTERNAL_ID":
			details.ExternalId = prop.Value
		case "AZURE_TENANT_ID":
			details.TenantId = prop.Value
		case "AZURE_CONSENT_URL":
			details.ConsentUrl = prop.Value
		case "AZURE_MULTI_TENANT_APP_NAME":
			details.MultiTenantAppName = prop.Value
		case "STORAGE_GCP_SERVICE_ACCOUNT":
			details.ServiceAccount = prop.Value
		}
	}
	return details, errors.Join(errs...)
}
