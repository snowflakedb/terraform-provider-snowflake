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

func (v *storageIntegrations) DescribeAwsDetails(ctx context.Context, id AccountObjectIdentifier) (*StorageIntegrationAwsDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseAwsProperties(properties)
}

func (v *storageIntegrations) DescribeAzureDetails(ctx context.Context, id AccountObjectIdentifier) (*StorageIntegrationAzureDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseAzureProperties(properties)
}

func (v *storageIntegrations) DescribeGcsDetails(ctx context.Context, id AccountObjectIdentifier) (*StorageIntegrationGcsDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseGcsProperties(properties)
}

// TODO [next PRs]: extract common mapping logic
func parseAwsProperties(properties []StorageIntegrationProperty) (*StorageIntegrationAwsDetails, error) {
	details := &StorageIntegrationAwsDetails{}
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
			details.AllowedLocations = strings.Split(prop.Value, ",")
		case "STORAGE_BLOCKED_LOCATIONS":
			details.BlockedLocations = strings.Split(prop.Value, ",")
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

func parseAzureProperties(properties []StorageIntegrationProperty) (*StorageIntegrationAzureDetails, error) {
	details := &StorageIntegrationAzureDetails{}
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
			details.AllowedLocations = strings.Split(prop.Value, ",")
		case "STORAGE_BLOCKED_LOCATIONS":
			details.BlockedLocations = strings.Split(prop.Value, ",")
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

func parseGcsProperties(properties []StorageIntegrationProperty) (*StorageIntegrationGcsDetails, error) {
	details := &StorageIntegrationGcsDetails{}
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
			details.AllowedLocations = strings.Split(prop.Value, ",")
		case "STORAGE_BLOCKED_LOCATIONS":
			details.BlockedLocations = strings.Split(prop.Value, ",")
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
