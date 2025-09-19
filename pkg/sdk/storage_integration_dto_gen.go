package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateStorageIntegrationOptions]   = new(CreateStorageIntegrationRequest)
	_ optionsProvider[AlterStorageIntegrationOptions]    = new(AlterStorageIntegrationRequest)
	_ optionsProvider[DropStorageIntegrationOptions]     = new(DropStorageIntegrationRequest)
	_ optionsProvider[ShowStorageIntegrationOptions]     = new(ShowStorageIntegrationRequest)
	_ optionsProvider[DescribeStorageIntegrationOptions] = new(DescribeStorageIntegrationRequest)
)

type CreateStorageIntegrationRequest struct {
	OrReplace                  *bool
	IfNotExists                *bool
	name                       AccountObjectIdentifier // required
	S3StorageProviderParams    *S3StorageParamsRequest
	GCSStorageProviderParams   *GCSStorageParamsRequest
	AzureStorageProviderParams *AzureStorageParamsRequest
	Enabled                    bool              // required
	StorageAllowedLocations    []StorageLocation // required
	StorageBlockedLocations    []StorageLocation
	Comment                    *string
}

type S3StorageParamsRequest struct {
	Protocol             S3Protocol // required
	StorageAwsRoleArn    string     // required
	StorageAwsExternalId *string
	StorageAwsObjectAcl  *string
}

type GCSStorageParamsRequest struct{}

type AzureStorageParamsRequest struct {
	AzureTenantId *string // required
}

type AlterStorageIntegrationRequest struct {
	IfExists  *bool
	name      AccountObjectIdentifier // required
	Set       *StorageIntegrationSetRequest
	Unset     *StorageIntegrationUnsetRequest
	SetTags   []TagAssociation
	UnsetTags []ObjectIdentifier
}

type StorageIntegrationSetRequest struct {
	S3Params                *SetS3StorageParamsRequest
	AzureParams             *SetAzureStorageParamsRequest
	Enabled                 *bool
	StorageAllowedLocations []StorageLocation
	StorageBlockedLocations []StorageLocation
	Comment                 *string
}

type SetS3StorageParamsRequest struct {
	StorageAwsRoleArn    string // required
	StorageAwsExternalId *string
	StorageAwsObjectAcl  *string
}

type SetAzureStorageParamsRequest struct {
	AzureTenantId string // required
}

type StorageIntegrationUnsetRequest struct {
	StorageAwsExternalId    *bool
	StorageAwsObjectAcl     *bool
	Enabled                 *bool
	StorageBlockedLocations *bool
	Comment                 *bool
}

type DropStorageIntegrationRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type ShowStorageIntegrationRequest struct {
	Like *Like
}

type DescribeStorageIntegrationRequest struct {
	name AccountObjectIdentifier // required
}
