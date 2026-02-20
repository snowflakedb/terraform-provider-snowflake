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
	storageLocation ExternalVolumeStorageLocation,
) (ExternalVolumeStorageLocation, error) {
	storageProvider, err := GetStorageLocationStorageProvider(storageLocation)
	if err != nil {
		return ExternalVolumeStorageLocation{}, err
	}

	newName := "terraform_provider_sentinel_storage_location"
	var tempNameStorageLocation ExternalVolumeStorageLocation
	switch storageProvider {
	case StorageProviderS3, StorageProviderS3GOV:
		tempNameStorageLocation = ExternalVolumeStorageLocation{
			S3StorageLocationParams: &S3StorageLocationParams{
				Name:                     newName,
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
			GCSStorageLocationParams: &GCSStorageLocationParams{
				Name:           newName,
				StorageBaseUrl: storageLocation.GCSStorageLocationParams.StorageBaseUrl,
				Encryption:     storageLocation.GCSStorageLocationParams.Encryption,
			},
		}
	case StorageProviderAzure:
		tempNameStorageLocation = ExternalVolumeStorageLocation{
			AzureStorageLocationParams: &AzureStorageLocationParams{
				Name:                   newName,
				StorageBaseUrl:         storageLocation.AzureStorageLocationParams.StorageBaseUrl,
				AzureTenantId:          storageLocation.AzureStorageLocationParams.AzureTenantId,
				UsePrivatelinkEndpoint: storageLocation.AzureStorageLocationParams.UsePrivatelinkEndpoint,
			},
		}
	case StorageProviderS3Compatible:
		tempNameStorageLocation = ExternalVolumeStorageLocation{
			S3CompatStorageLocationParams: &S3CompatStorageLocationParams{
				Name:            newName,
				StorageBaseUrl:  storageLocation.S3CompatStorageLocationParams.StorageBaseUrl,
				StorageEndpoint: storageLocation.S3CompatStorageLocationParams.StorageEndpoint,
				Credentials:     storageLocation.S3CompatStorageLocationParams.Credentials,
			},
		}
	default:
		return ExternalVolumeStorageLocation{}, fmt.Errorf("unsupported storage provider type: %s", storageProvider)
	}

	return tempNameStorageLocation, nil
}

// GetStorageLocationName retrieves the name from a storage location configuration.
// Returns an error if the storage location is invalid or the name is empty.
func GetStorageLocationName(s ExternalVolumeStorageLocation) (string, error) {
	switch {
	case s.S3StorageLocationParams != nil && *s.S3StorageLocationParams != S3StorageLocationParams{}:
		if len(s.S3StorageLocationParams.Name) == 0 {
			return "", fmt.Errorf("Invalid S3 storage location - no name set")
		}

		return s.S3StorageLocationParams.Name, nil
	case s.GCSStorageLocationParams != nil && *s.GCSStorageLocationParams != GCSStorageLocationParams{}:
		if len(s.GCSStorageLocationParams.Name) == 0 {
			return "", fmt.Errorf("Invalid GCS storage location - no name set")
		}

		return s.GCSStorageLocationParams.Name, nil
	case s.AzureStorageLocationParams != nil && *s.AzureStorageLocationParams != AzureStorageLocationParams{}:
		if len(s.AzureStorageLocationParams.Name) == 0 {
			return "", fmt.Errorf("Invalid Azure storage location - no name set")
		}

		return s.AzureStorageLocationParams.Name, nil
	case s.S3CompatStorageLocationParams != nil && *s.S3CompatStorageLocationParams != S3CompatStorageLocationParams{}:
		if len(s.S3CompatStorageLocationParams.Name) == 0 {
			return "", fmt.Errorf("Invalid S3Compatible storage location - no name set")
		}

		return s.S3CompatStorageLocationParams.Name, nil
	default:
		return "", fmt.Errorf("Invalid storage location")
	}
}

func GetStorageLocationStorageProvider(s ExternalVolumeStorageLocation) (StorageProvider, error) {
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


type ExternalVolumeDetails struct {
	StorageLocations []ExternalVolumeStorageLocationDetails
	Active           string
	Comment          string
	AllowWrites      string
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
		case strings.Contains(p.Name, "STORAGE_LOCATION_"):
			storageLocation := ExternalVolumeStorageLocationDetails{}
			err := json.Unmarshal([]byte(p.Value), &storageLocation)
			if err != nil {
				return ExternalVolumeDetails{}, err
			}
			storageLocations = append(
				storageLocations,
				storageLocation,
			)
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

type ExternalVolumeStorageLocationDetails struct {
	Name                     string `json:"NAME"`
	StorageProvider          string `json:"STORAGE_PROVIDER"`
	StorageBaseUrl           string `json:"STORAGE_BASE_URL"`
	StorageAwsRoleArn        string `json:"STORAGE_AWS_ROLE_ARN"`
	StorageAwsExternalId     string `json:"STORAGE_AWS_EXTERNAL_ID"`
	StorageAwsAccessPointArn string `json:"STORAGE_AWS_ACCESS_POINT_ARN"`
	StorageEndpoint          string `json:"STORAGE_ENDPOINT"`
	UsePrivatelinkEndpoint   string `json:"USE_PRIVATELINK_ENDPOINT"`
	EncryptionType           string `json:"ENCRYPTION_TYPE"`
	EncryptionKmsKeyId       string `json:"ENCRYPTION_KMS_KEY_ID"`
	AzureTenantId            string `json:"AZURE_TENANT_ID"`
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
			return fmt.Errorf("invalid storage provider parsed: %s", s)
		}
		if len(s.StorageBaseUrl) == 0 {
			return fmt.Errorf("A storage location's StorageBaseUrl in this volume could not be parsed.")
		}

		storageProvider, err := ToStorageProviderInDescribe(s.StorageProvider)
		if err != nil {
			return err
		}

		switch storageProvider {
		case StorageProviderS3, StorageProviderS3GOV:
			if len(s.StorageAwsRoleArn) == 0 {
				return fmt.Errorf("An S3 storage location's StorageAwsRoleArn in this volume could not be parsed.")
			}
		case StorageProviderAzure:
			if len(s.AzureTenantId) == 0 {
				return fmt.Errorf("An Azure storage location's AzureTenantId in this volume could not be parsed.")
			}
		case StorageProviderS3Compatible:
			if len(s.StorageEndpoint) == 0 {
				return fmt.Errorf("An S3Compatible storage location's StorageEndpoint in this volume could not be parsed.")
			}
		}
	}

	return nil
}
