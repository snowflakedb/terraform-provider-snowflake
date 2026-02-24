package sdk

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// CopySentinelStorageLocation creates a copy of the given storage location with a
// sentinel name used for Terraform provider operations. This is useful for managing
// storage location state without affecting user-visible names.
func CopySentinelStorageLocation(
	storageLocationItem ExternalVolumeStorageLocationItem,
) (ExternalVolumeStorageLocationItem, error) {
	storageLocation := storageLocationItem.ExternalVolumeStorageLocation
	storageProvider, err := GetStorageLocationStorageProvider(storageLocationItem)
	if err != nil {
		return ExternalVolumeStorageLocationItem{}, err
	}

	newName := "terraform_provider_sentinel_storage_location"
	var tempNameStorageLocation ExternalVolumeStorageLocation
	switch storageProvider {
	case StorageProviderS3, StorageProviderS3GOV:
		tempNameStorageLocation = ExternalVolumeStorageLocation{
			Name: newName,
			S3StorageLocationParams: &S3StorageLocationParams{
				StorageProvider:          storageLocation.S3StorageLocationParams.StorageProvider,
				StorageBaseUrl:           storageLocation.S3StorageLocationParams.StorageBaseUrl,
				StorageAwsRoleArn:        storageLocation.S3StorageLocationParams.StorageAwsRoleArn,
				StorageAwsExternalId:     storageLocation.S3StorageLocationParams.StorageAwsExternalId,
				StorageAwsAccessPointArn: storageLocation.S3StorageLocationParams.StorageAwsAccessPointArn,
				UsePrivatelinkEndpoint:   storageLocation.S3StorageLocationParams.UsePrivatelinkEndpoint,
				Encryption:               storageLocation.S3StorageLocationParams.Encryption,
			},
		}
	case StorageProviderGCS:
		tempNameStorageLocation = ExternalVolumeStorageLocation{
			Name: newName,
			GCSStorageLocationParams: &GCSStorageLocationParams{
				StorageBaseUrl: storageLocation.GCSStorageLocationParams.StorageBaseUrl,
				Encryption:     storageLocation.GCSStorageLocationParams.Encryption,
			},
		}
	case StorageProviderAzure:
		tempNameStorageLocation = ExternalVolumeStorageLocation{
			Name: newName,
			AzureStorageLocationParams: &AzureStorageLocationParams{
				StorageBaseUrl:         storageLocation.AzureStorageLocationParams.StorageBaseUrl,
				AzureTenantId:          storageLocation.AzureStorageLocationParams.AzureTenantId,
				UsePrivatelinkEndpoint: storageLocation.AzureStorageLocationParams.UsePrivatelinkEndpoint,
			},
		}
	case StorageProviderS3Compatible:
		tempNameStorageLocation = ExternalVolumeStorageLocation{
			Name: newName,
			S3CompatStorageLocationParams: &S3CompatStorageLocationParams{
				StorageBaseUrl:  storageLocation.S3CompatStorageLocationParams.StorageBaseUrl,
				StorageEndpoint: storageLocation.S3CompatStorageLocationParams.StorageEndpoint,
				Credentials:     storageLocation.S3CompatStorageLocationParams.Credentials,
			},
		}
	default:
		return ExternalVolumeStorageLocationItem{}, fmt.Errorf("unsupported storage provider type: %s", storageProvider)
	}

	return ExternalVolumeStorageLocationItem{
		ExternalVolumeStorageLocation: tempNameStorageLocation,
	}, nil
}

func GetStorageLocationStorageProvider(i ExternalVolumeStorageLocationItem) (StorageProvider, error) {
	s := i.ExternalVolumeStorageLocation
	switch {
	case s.S3StorageLocationParams != nil && *s.S3StorageLocationParams != S3StorageLocationParams{}:
		return ToStorageProvider(string(s.S3StorageLocationParams.StorageProvider))
	case s.GCSStorageLocationParams != nil && *s.GCSStorageLocationParams != GCSStorageLocationParams{}:
		return StorageProviderGCS, nil
	case s.AzureStorageLocationParams != nil && *s.AzureStorageLocationParams != AzureStorageLocationParams{}:
		return StorageProviderAzure, nil
	case s.S3CompatStorageLocationParams != nil && *s.S3CompatStorageLocationParams != S3CompatStorageLocationParams{}:
		return StorageProviderS3Compatible, nil
	default:
		return "", fmt.Errorf("Invalid storage location")
	}
}

// ExternalVolumeStorageLocationDetails is the typed representation of a storage location
// returned by DESCRIBE EXTERNAL VOLUME. Common fields live on the wrapper; provider-specific
// fields live on exactly one of the provider sub-structs (the rest are nil).
type ExternalVolumeStorageLocationDetails struct {
	Name                    string
	StorageProvider         string
	StorageBaseUrl          string
	StorageAllowedLocations []string
	EncryptionType          string
	S3StorageLocation       *StorageLocationS3Details
	GCSStorageLocation      *StorageLocationGcsDetails
	AzureStorageLocation    *StorageLocationAzureDetails
	S3CompatStorageLocation *StorageLocationS3CompatDetails
}

type StorageLocationS3Details struct {
	StorageAwsRoleArn        string
	StorageAwsIamUserArn     string
	StorageAwsExternalId     string
	StorageAwsAccessPointArn string
	UsePrivatelinkEndpoint   string
	EncryptionKmsKeyId       string
}

type StorageLocationGcsDetails struct {
	StorageGcpServiceAccount string
	EncryptionKmsKeyId       string
}

type StorageLocationAzureDetails struct {
	AzureTenantId           string
	AzureMultiTenantAppName string
	AzureConsentUrl         string
}

type StorageLocationS3CompatDetails struct {
	Endpoint           string
	AwsAccessKeyId     string
	EncryptionKmsKeyId string
}

type ExternalVolumeDetails struct {
	StorageLocations []ExternalVolumeStorageLocationDetails
	Active           string
	Comment          string
	AllowWrites      string
}

// externalVolumeStorageLocationJsonRaw is the internal JSON deserialization type
// for storage location properties.
type externalVolumeStorageLocationJsonRaw struct {
	Name                     string   `json:"NAME"`
	StorageProvider          string   `json:"STORAGE_PROVIDER"`
	StorageBaseUrl           string   `json:"STORAGE_BASE_URL"`
	StorageAllowedLocations  []string `json:"STORAGE_ALLOWED_LOCATIONS"`
	StorageAwsRoleArn        string   `json:"STORAGE_AWS_ROLE_ARN"`
	StorageAwsIamUserArn     string   `json:"STORAGE_AWS_IAM_USER_ARN"`
	StorageAwsExternalId     string   `json:"STORAGE_AWS_EXTERNAL_ID"`
	StorageAwsAccessPointArn string   `json:"STORAGE_AWS_ACCESS_POINT_ARN"`
	Endpoint                 string   `json:"ENDPOINT"`
	UsePrivatelinkEndpoint   string   `json:"USE_PRIVATELINK_ENDPOINT"`
	EncryptionType           string   `json:"ENCRYPTION_TYPE"`
	EncryptionKmsKeyId       string   `json:"ENCRYPTION_KMS_KEY_ID"`
	AzureTenantId            string   `json:"AZURE_TENANT_ID"`
	AzureMultiTenantAppName  string   `json:"AZURE_MULTI_TENANT_APP_NAME"`
	AzureConsentUrl          string   `json:"AZURE_CONSENT_URL"`
	StorageGcpServiceAccount string   `json:"STORAGE_GCP_SERVICE_ACCOUNT"`
	AwsAccessKeyId           string   `json:"AWS_ACCESS_KEY_ID"`
}

func (e externalVolumeStorageLocationJsonRaw) toStorageLocationDetails() (ExternalVolumeStorageLocationDetails, error) {
	details := ExternalVolumeStorageLocationDetails{
		Name:                    e.Name,
		StorageProvider:         e.StorageProvider,
		StorageBaseUrl:          e.StorageBaseUrl,
		StorageAllowedLocations: e.StorageAllowedLocations,
		EncryptionType:          e.EncryptionType,
	}

	storageProvider, err := ToStorageProviderInDescribe(e.StorageProvider)
	if err != nil {
		return ExternalVolumeStorageLocationDetails{}, err
	}

	switch storageProvider {
	case StorageProviderS3, StorageProviderS3GOV:
		details.S3StorageLocation = &StorageLocationS3Details{
			StorageAwsRoleArn:        e.StorageAwsRoleArn,
			StorageAwsIamUserArn:     e.StorageAwsIamUserArn,
			StorageAwsExternalId:     e.StorageAwsExternalId,
			StorageAwsAccessPointArn: e.StorageAwsAccessPointArn,
			UsePrivatelinkEndpoint:   e.UsePrivatelinkEndpoint,
			EncryptionKmsKeyId:       e.EncryptionKmsKeyId,
		}
	case StorageProviderGCS:
		details.GCSStorageLocation = &StorageLocationGcsDetails{
			StorageGcpServiceAccount: e.StorageGcpServiceAccount,
			EncryptionKmsKeyId:       e.EncryptionKmsKeyId,
		}
	case StorageProviderAzure:
		details.AzureStorageLocation = &StorageLocationAzureDetails{
			AzureTenantId:           e.AzureTenantId,
			AzureMultiTenantAppName: e.AzureMultiTenantAppName,
			AzureConsentUrl:         e.AzureConsentUrl,
		}
	case StorageProviderS3Compatible:
		details.S3CompatStorageLocation = &StorageLocationS3CompatDetails{
			Endpoint:           e.Endpoint,
			AwsAccessKeyId:     e.AwsAccessKeyId,
			EncryptionKmsKeyId: e.EncryptionKmsKeyId,
		}
	default:
		return ExternalVolumeStorageLocationDetails{}, fmt.Errorf("unsupported storage provider type: %s", e.StorageProvider)
	}

	return details, nil
}

func ParseExternalVolumeDescribed(props []ExternalVolumeProperty) (ExternalVolumeDetails, error) {
	externalVolumeDetails := ExternalVolumeDetails{}
	var storageLocations []ExternalVolumeStorageLocationDetails
	for _, p := range props {
		switch {
		case p.Name == "COMMENT":
			externalVolumeDetails.Comment = p.Value
		case p.Name == "ACTIVE":
			externalVolumeDetails.Active = p.Value
		case p.Name == "ALLOW_WRITES":
			externalVolumeDetails.AllowWrites = p.Value
		// TODO: don't assume the order is correct. Parse the location index from the name.
		case strings.Contains(p.Name, "STORAGE_LOCATION_"):
			var raw externalVolumeStorageLocationJsonRaw
			err := json.Unmarshal([]byte(p.Value), &raw)
			if err != nil {
				return ExternalVolumeDetails{}, err
			}
			details, err := raw.toStorageLocationDetails()
			if err != nil {
				return ExternalVolumeDetails{}, err
			}
			storageLocations = append(storageLocations, details)
		default:
			return ExternalVolumeDetails{}, fmt.Errorf("Unrecognized external volume property: %s", p.Name)
		}
	}

	externalVolumeDetails.StorageLocations = storageLocations
	err := validateExternalVolumeDetails(externalVolumeDetails)
	if err != nil {
		return ExternalVolumeDetails{}, err
	}

	return externalVolumeDetails, nil
}

func validateExternalVolumeDetails(p ExternalVolumeDetails) error {
	if len(p.StorageLocations) == 0 {
		return fmt.Errorf("No storage locations could be parsed from the external volume.")
	}
	if len(p.AllowWrites) == 0 {
		return fmt.Errorf("The external volume AllowWrites property could not be parsed.")
	}

	for _, s := range p.StorageLocations {
		if len(s.Name) == 0 {
			return fmt.Errorf("A storage location's Name in this volume could not be parsed.")
		}
		if !slices.Contains(AsStringList(AllStorageProviderValuesInDescribe), s.StorageProvider) {
			return fmt.Errorf("invalid storage provider parsed: %v", s)
		}
		if len(s.StorageBaseUrl) == 0 {
			return fmt.Errorf("A storage location's StorageBaseUrl in this volume could not be parsed.")
		}

		switch {
		case s.S3StorageLocation != nil:
			if len(s.S3StorageLocation.StorageAwsRoleArn) == 0 {
				return fmt.Errorf("An S3 storage location's StorageAwsRoleArn in this volume could not be parsed.")
			}
		case s.AzureStorageLocation != nil:
			if len(s.AzureStorageLocation.AzureTenantId) == 0 {
				return fmt.Errorf("An Azure storage location's AzureTenantId in this volume could not be parsed.")
			}
		case s.S3CompatStorageLocation != nil:
			if len(s.S3CompatStorageLocation.Endpoint) == 0 {
				return fmt.Errorf("An S3Compatible storage location's StorageEndpoint in this volume could not be parsed.")
			}
		}
	}

	return nil
}

func (r *CreateExternalVolumeRequest) GetName() AccountObjectIdentifier {
	return r.name
}
