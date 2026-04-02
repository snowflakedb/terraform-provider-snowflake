package gen

import (
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type SdkObjectDef struct {
	IdType             string
	ObjectStruct       any
	IsDataSourceOutput bool
	// NameOverride, when non-empty, replaces the struct type name used for code generation.
	// Useful when multiple resources share the same SDK struct (e.g. WarehouseAdaptive reuses sdk.Warehouse).
	NameOverride string
	// FieldsToInclude, when non-empty, restricts code generation to only these struct field names.
	FieldsToInclude []string
}

var allStructs = []SdkObjectDef{
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.Database{},
	},
	{
		IdType:       "sdk.DatabaseObjectIdentifier",
		ObjectStruct: sdk.Schema{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.Role{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.Connection{},
	},
	{
		IdType:       "sdk.DatabaseObjectIdentifier",
		ObjectStruct: sdk.DatabaseRole{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.RowAccessPolicy{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.User{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.View{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.Warehouse{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.Warehouse{},
		NameOverride: "sdk.WarehouseAdaptive",
		FieldsToInclude: []string{
			"Name", "State", "Type", "Running", "Queued", "IsDefault", "IsCurrent",
			"AutoResume", "Available", "Provisioning", "Quiescing", "Other",
			"CreatedOn", "ResumedOn", "UpdatedOn", "Owner", "Comment",
			"ResourceMonitor", "OwnerRoleType", "MaxQueryPerformanceLevel", "QueryThroughputMultiplier",
		},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.ResourceMonitor{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.NetworkPolicy{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.MaskingPolicy{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.AuthenticationPolicy{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.Task{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.ExternalVolume{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.Secret{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.Stream{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.Tag{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.Account{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifierWithArguments",
		ObjectStruct: sdk.Function{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifierWithArguments",
		ObjectStruct: sdk.Procedure{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.ImageRepository{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.ComputePool{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.GitRepository{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.Service{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.ProgrammaticAccessToken{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.OrganizationAccount{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.Listing{},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.ListingDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.SemanticView{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.StorageIntegration{},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.StorageIntegrationAwsDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.StorageIntegrationAzureDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.StorageIntegrationGcsDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.Notebook{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.SecurityIntegration{},
	},
	{
		IdType:       "sdk.DatabaseObjectIdentifier",
		ObjectStruct: sdk.Schema{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.Streamlit{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.UserWorkloadIdentityAuthenticationMethod{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.Stage{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.NetworkRule{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.ExternalVolumeStorageLocationDetails{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.StorageLocationS3Details{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.StorageLocationGcsDetails{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.StorageLocationAzureDetails{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.StorageLocationS3CompatDetails{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.CatalogIntegration{},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.CatalogIntegrationAwsGlueDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.CatalogIntegrationObjectStorageDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.CatalogIntegrationOpenCatalogDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.CatalogIntegrationIcebergRestDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.CatalogIntegrationSapBdcDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.OpenCatalogRestConfigDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.IcebergRestRestConfigDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.OAuthRestAuthenticationDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.SigV4RestAuthenticationDetails{},
		IsDataSourceOutput: true,
	},
}

func GetSdkObjectDetails() []genhelpers.SdkObjectDetails {
	allSdkObjectsDetails := make([]genhelpers.SdkObjectDetails, len(allStructs))
	for idx, d := range allStructs {
		structDetails := genhelpers.ExtractStructDetails(d.ObjectStruct)
		if d.NameOverride != "" {
			structDetails.Name = d.NameOverride
		}
		if len(d.FieldsToInclude) > 0 {
			filtered := make([]genhelpers.Field, 0, len(d.FieldsToInclude))
			for _, f := range structDetails.Fields {
				if slices.Contains(d.FieldsToInclude, f.Name) {
					filtered = append(filtered, f)
				}
			}
			structDetails.Fields = filtered
		}
		allSdkObjectsDetails[idx] = genhelpers.SdkObjectDetails{
			IdType:             d.IdType,
			StructDetails:      structDetails,
			IsDataSourceOutput: d.IsDataSourceOutput,
		}
	}
	return allSdkObjectsDetails
}
