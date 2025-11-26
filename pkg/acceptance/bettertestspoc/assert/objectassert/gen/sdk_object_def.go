package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type SdkObjectDef struct {
	IdType             string
	ObjectStruct       any
	IsDataSourceOutput bool
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
}

func GetSdkObjectDetails() []genhelpers.SdkObjectDetails {
	allSdkObjectsDetails := make([]genhelpers.SdkObjectDetails, len(allStructs))
	for idx, d := range allStructs {
		structDetails := genhelpers.ExtractStructDetails(d.ObjectStruct)
		allSdkObjectsDetails[idx] = genhelpers.SdkObjectDetails{
			IdType:             d.IdType,
			StructDetails:      structDetails,
			IsDataSourceOutput: d.IsDataSourceOutput,
		}
	}
	return allSdkObjectsDetails
}
