package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

type ExternalVolumes interface {
	Create(ctx context.Context, request *CreateExternalVolumeRequest) error
	Alter(ctx context.Context, request *AlterExternalVolumeRequest) error
	Drop(ctx context.Context, request *DropExternalVolumeRequest) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
	Describe(ctx context.Context, id AccountObjectIdentifier) ([]ExternalVolumeProperty, error)
	Show(ctx context.Context, request *ShowExternalVolumeRequest) ([]ExternalVolume, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ExternalVolume, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*ExternalVolume, error)
}

// CreateExternalVolumeOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-external-volume.
type CreateExternalVolumeOptions struct {
	create           bool                            `ddl:"static" sql:"CREATE"`
	OrReplace        *bool                           `ddl:"keyword" sql:"OR REPLACE"`
	externalVolume   bool                            `ddl:"static" sql:"EXTERNAL VOLUME"`
	IfNotExists      *bool                           `ddl:"keyword" sql:"IF NOT EXISTS"`
	name             AccountObjectIdentifier         `ddl:"identifier"`
	StorageLocations []ExternalVolumeStorageLocation `ddl:"parameter,parentheses" sql:"STORAGE_LOCATIONS"`
	AllowWrites      *bool                           `ddl:"parameter" sql:"ALLOW_WRITES"`
	Comment          *string                         `ddl:"parameter,single_quotes" sql:"COMMENT"`
}
type ExternalVolumeStorageLocation struct {
	S3StorageLocationParams    *S3StorageLocationParams    `ddl:"list,parentheses,no_comma"`
	GCSStorageLocationParams   *GCSStorageLocationParams   `ddl:"list,parentheses,no_comma"`
	AzureStorageLocationParams *AzureStorageLocationParams `ddl:"list,parentheses,no_comma"`
}
type S3StorageLocationParams struct {
	Name                 string                      `ddl:"parameter,single_quotes" sql:"NAME"`
	StorageProvider      S3StorageProvider           `ddl:"parameter,single_quotes" sql:"STORAGE_PROVIDER"`
	StorageAwsRoleArn    string                      `ddl:"parameter,single_quotes" sql:"STORAGE_AWS_ROLE_ARN"`
	StorageBaseUrl       string                      `ddl:"parameter,single_quotes" sql:"STORAGE_BASE_URL"`
	StorageAwsExternalId *string                     `ddl:"parameter,single_quotes" sql:"STORAGE_AWS_EXTERNAL_ID"`
	Encryption           *ExternalVolumeS3Encryption `ddl:"list,parentheses,no_comma" sql:"ENCRYPTION ="`
}
type ExternalVolumeS3Encryption struct {
	Type     S3EncryptionType `ddl:"parameter,single_quotes" sql:"TYPE"`
	KmsKeyId *string          `ddl:"parameter,single_quotes" sql:"KMS_KEY_ID"`
}
type GCSStorageLocationParams struct {
	Name               string                       `ddl:"parameter,single_quotes" sql:"NAME"`
	StorageProviderGcs string                       `ddl:"static" sql:"STORAGE_PROVIDER = 'GCS'"`
	StorageBaseUrl     string                       `ddl:"parameter,single_quotes" sql:"STORAGE_BASE_URL"`
	Encryption         *ExternalVolumeGCSEncryption `ddl:"list,parentheses,no_comma" sql:"ENCRYPTION ="`
}
type ExternalVolumeGCSEncryption struct {
	Type     GCSEncryptionType `ddl:"parameter,single_quotes" sql:"TYPE"`
	KmsKeyId *string           `ddl:"parameter,single_quotes" sql:"KMS_KEY_ID"`
}
type AzureStorageLocationParams struct {
	Name                 string `ddl:"parameter,single_quotes" sql:"NAME"`
	StorageProviderAzure string `ddl:"static" sql:"STORAGE_PROVIDER = 'AZURE'"`
	AzureTenantId        string `ddl:"parameter,single_quotes" sql:"AZURE_TENANT_ID"`
	StorageBaseUrl       string `ddl:"parameter,single_quotes" sql:"STORAGE_BASE_URL"`
}

// AlterExternalVolumeOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-external-volume.
type AlterExternalVolumeOptions struct {
	alter                 bool                           `ddl:"static" sql:"ALTER"`
	externalVolume        bool                           `ddl:"static" sql:"EXTERNAL VOLUME"`
	IfExists              *bool                          `ddl:"keyword" sql:"IF EXISTS"`
	name                  AccountObjectIdentifier        `ddl:"identifier"`
	RemoveStorageLocation *string                        `ddl:"parameter,single_quotes,no_equals" sql:"REMOVE STORAGE_LOCATION"`
	Set                   *AlterExternalVolumeSet        `ddl:"keyword" sql:"SET"`
	AddStorageLocation    *ExternalVolumeStorageLocation `ddl:"parameter" sql:"ADD STORAGE_LOCATION"`
}
type AlterExternalVolumeSet struct {
	AllowWrites *bool   `ddl:"parameter" sql:"ALLOW_WRITES"`
	Comment     *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// DropExternalVolumeOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-external-volume.
type DropExternalVolumeOptions struct {
	drop           bool                    `ddl:"static" sql:"DROP"`
	externalVolume bool                    `ddl:"static" sql:"EXTERNAL VOLUME"`
	IfExists       *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name           AccountObjectIdentifier `ddl:"identifier"`
}

// DescribeExternalVolumeOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-external-volume.
type DescribeExternalVolumeOptions struct {
	describe       bool                    `ddl:"static" sql:"DESCRIBE"`
	externalVolume bool                    `ddl:"static" sql:"EXTERNAL VOLUME"`
	name           AccountObjectIdentifier `ddl:"identifier"`
}
type externalVolumeDescRow struct {
	ParentProperty  string `db:"parent_property"`
	Property        string `db:"property"`
	PropertyType    string `db:"property_type"`
	PropertyValue   string `db:"property_value"`
	PropertyDefault string `db:"property_default"`
}
type ExternalVolumeProperty struct {
	Parent  string
	Name    string
	Type    string
	Value   string
	Default string
}

// ShowExternalVolumeOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-external-volumes.
type ShowExternalVolumeOptions struct {
	show            bool  `ddl:"static" sql:"SHOW"`
	externalVolumes bool  `ddl:"static" sql:"EXTERNAL VOLUMES"`
	Like            *Like `ddl:"keyword" sql:"LIKE"`
}
type externalVolumeShowRow struct {
	Name        string         `db:"name"`
	AllowWrites bool           `db:"allow_writes"`
	Comment     sql.NullString `db:"comment"`
}
type ExternalVolume struct {
	Name        string
	AllowWrites bool
	Comment     string
}

// Returns a copy of the given storage location with a set name
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
				Name:                 newName,
				StorageProvider:      storageLocation.S3StorageLocationParams.StorageProvider,
				StorageBaseUrl:       storageLocation.S3StorageLocationParams.StorageBaseUrl,
				StorageAwsRoleArn:    storageLocation.S3StorageLocationParams.StorageAwsRoleArn,
				StorageAwsExternalId: storageLocation.S3StorageLocationParams.StorageAwsExternalId,
				Encryption:           storageLocation.S3StorageLocationParams.Encryption,
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
				Name:           newName,
				StorageBaseUrl: storageLocation.AzureStorageLocationParams.StorageBaseUrl,
				AzureTenantId:  storageLocation.AzureStorageLocationParams.AzureTenantId,
			},
		}
	}

	return tempNameStorageLocation, nil
}

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
	default:
		return "", fmt.Errorf("Invalid storage location")
	}
}

// Returns the index of the last matching elements in the list
// e.g. [1,2,3] [1,3,2] -> 0, [1,2,3] [1,2,4] -> 1
// -1 is returned if there are no common prefixes in the list
func CommonPrefixLastIndex(a []ExternalVolumeStorageLocation, b []ExternalVolumeStorageLocation) (int, error) {
	commonPrefixLastIndex := 0

	if len(a) == 0 || len(b) == 0 {
		return -1, nil
	}

	if !reflect.DeepEqual(a[0], b[0]) {
		return -1, nil
	}

	for i := 1; i < min(len(a), len(b)); i++ {
		if !reflect.DeepEqual(a[i], b[i]) {
			break
		}

		commonPrefixLastIndex = i
	}

	return commonPrefixLastIndex, nil
}

func (v *ExternalVolume) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}
