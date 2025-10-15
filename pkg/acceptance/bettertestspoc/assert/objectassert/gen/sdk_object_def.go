package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type SdkObjectDef struct {
	IdType       string
	ObjectType   sdk.ObjectType
	ObjectStruct any
}

// TODO [SNOW-2324252]: remove object type?
var allStructs = []SdkObjectDef{
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectType:   sdk.ObjectTypeDatabase,
		ObjectStruct: sdk.Database{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectType:   sdk.ObjectTypeConnection,
		ObjectStruct: sdk.Connection{},
	},
	{
		IdType:       "sdk.DatabaseObjectIdentifier",
		ObjectType:   sdk.ObjectTypeDatabaseRole,
		ObjectStruct: sdk.DatabaseRole{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeRowAccessPolicy,
		ObjectStruct: sdk.RowAccessPolicy{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectType:   sdk.ObjectTypeUser,
		ObjectStruct: sdk.User{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeView,
		ObjectStruct: sdk.View{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectType:   sdk.ObjectTypeWarehouse,
		ObjectStruct: sdk.Warehouse{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectType:   sdk.ObjectTypeResourceMonitor,
		ObjectStruct: sdk.ResourceMonitor{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeMaskingPolicy,
		ObjectStruct: sdk.MaskingPolicy{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeAuthenticationPolicy,
		ObjectStruct: sdk.AuthenticationPolicy{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeTask,
		ObjectStruct: sdk.Task{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectType:   sdk.ObjectTypeExternalVolume,
		ObjectStruct: sdk.ExternalVolume{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeSecret,
		ObjectStruct: sdk.Secret{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeStream,
		ObjectStruct: sdk.Stream{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeTag,
		ObjectStruct: sdk.Tag{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectType:   sdk.ObjectTypeAccount,
		ObjectStruct: sdk.Account{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifierWithArguments",
		ObjectType:   sdk.ObjectTypeFunction,
		ObjectStruct: sdk.Function{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifierWithArguments",
		ObjectType:   sdk.ObjectTypeProcedure,
		ObjectStruct: sdk.Procedure{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeImageRepository,
		ObjectStruct: sdk.ImageRepository{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectType:   sdk.ObjectTypeComputePool,
		ObjectStruct: sdk.ComputePool{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeGitRepository,
		ObjectStruct: sdk.GitRepository{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeService,
		ObjectStruct: sdk.Service{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectType:   sdk.ObjectTypeProgrammaticAccessToken,
		ObjectStruct: sdk.ProgrammaticAccessToken{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectType:   sdk.ObjectTypeAccount,
		ObjectStruct: sdk.OrganizationAccount{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectType:   sdk.ObjectTypeListing,
		ObjectStruct: sdk.Listing{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeSemanticView,
		ObjectStruct: sdk.SemanticView{},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectType:   sdk.ObjectTypeStorageIntegration,
		ObjectStruct: sdk.StorageIntegration{},
	},
	{
		IdType:       "sdkSchemaObjectIdentifier",
		ObjectType:   sdk.ObjectTypeNotebook,
		ObjectStruct: sdk.Notebook{},
	},
}

func GetSdkObjectDetails() []genhelpers.SdkObjectDetails {
	allSdkObjectsDetails := make([]genhelpers.SdkObjectDetails, len(allStructs))
	for idx, d := range allStructs {
		structDetails := genhelpers.ExtractStructDetails(d.ObjectStruct)
		allSdkObjectsDetails[idx] = genhelpers.SdkObjectDetails{
			IdType:        d.IdType,
			StructDetails: structDetails,
		}
	}
	return allSdkObjectsDetails
}
