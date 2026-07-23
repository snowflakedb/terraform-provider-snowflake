package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type SdkObjectDef struct {
	IdType             string
	ObjectStruct       any
	IsDataSourceOutput bool
	IsSubStruct        bool
	ObjectTypeName     string
	NoShowById         bool
	// NoIdentifiableObject marks objects that have no meaningful identifier (e.g. list-returned
	// describe results). Suppresses the constructor and New*Assert(), uses placeholder in FromObject.
	NoIdentifiableObject bool
	// ShowByParentId groups the fields needed to generate a constructor that fetches the object
	// via a parent identifier (e.g. userId for ProgrammaticAccessToken). All three fields must be set together.
	ShowByParentId *genhelpers.ShowByParentIdDef
	// NestedAssertFields lists struct-type fields that get the nested assertion adapter pattern instead
	// of a simple Has* assertion. The sub-struct type must have a corresponding asserter in allStructs.
	NestedAssertFields []string
	// SkipFields lists fields to suppress Has* generation entirely. Use for fields that have a
	// custom implementation in an _ext.go file that should not be overwritten.
	SkipFields []string
	// DescribeOverride overrides the default test client and method in the IsDataSourceOutput constructor
	// when the naming convention does not match the actual helper.
	DescribeOverride *genhelpers.DescribeOverrideDef
	// FromObjectIDExpr overrides the default `<object>.ID()` expression in the FromObject constructor.
	// Use when the SDK ID() return type doesn't match IdType (e.g. AccountIdentifier vs AccountObjectIdentifier).
	FromObjectIDExpr string
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
		SkipFields:   []string{"Tables"},
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
		SkipFields:   []string{"TaskRelations"},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.ExternalVolume{},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.ExternalVolumeDetails{},
		IsDataSourceOutput: true,
		NestedAssertFields: []string{"StorageLocations"},
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
		ObjectStruct: sdk.HybridTable{},
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
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.ComputePoolDetails{},
		IsDataSourceOutput: true,
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
		IdType:             "sdk.SchemaObjectIdentifier",
		ObjectStruct:       sdk.ServiceDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.ProgrammaticAccessToken{},
		ShowByParentId: &genhelpers.ShowByParentIdDef{
			ParentIdType:   "sdk.AccountObjectIdentifier",
			ClientName:     "User",
			ShowMethodName: "ShowProgrammaticAccessToken",
		},
	},
	{
		IdType:           "sdk.AccountObjectIdentifier",
		ObjectStruct:     sdk.OrganizationAccount{},
		ObjectTypeName:   "Account",
		FromObjectIDExpr: "organizationAccount.ID().AsAccountObjectIdentifier()",
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
		ObjectStruct: sdk.StorageLifecyclePolicy{},
	},
	{
		IdType:             "sdk.SchemaObjectIdentifier",
		ObjectStruct:       sdk.StorageLifecyclePolicyDetails{},
		IsDataSourceOutput: true,
		SkipFields:         []string{"Signature"},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.StorageIntegrationAwsDetails{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "StorageIntegration", MethodName: "DescribeAws"},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.StorageIntegrationAzureDetails{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "StorageIntegration", MethodName: "DescribeAzure"},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.StorageIntegrationGcsDetails{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "StorageIntegration", MethodName: "DescribeGcs"},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.Notebook{},
	},
	{
		IdType:             "sdk.SchemaObjectIdentifier",
		ObjectStruct:       sdk.NotebookDetails{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "Notebook", MethodName: "Describe"},
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
		ShowByParentId: &genhelpers.ShowByParentIdDef{
			ParentIdType:   "sdk.AccountObjectIdentifier",
			ClientName:     "User",
			ShowMethodName: "ShowUserWorkloadIdentityAuthenticationMethod",
		},
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
		IdType:             "sdk.SchemaObjectIdentifier",
		ObjectStruct:       sdk.NetworkRuleDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.Pipe{},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.ExternalVolumeStorageLocationDetails{},
		IsSubStruct:        true,
		NestedAssertFields: []string{"S3StorageLocation", "GCSStorageLocation", "AzureStorageLocation", "S3CompatStorageLocation"},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.StorageLocationS3Details{},
		IsSubStruct:  true,
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.StorageLocationGcsDetails{},
		IsSubStruct:  true,
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.StorageLocationAzureDetails{},
		IsSubStruct:  true,
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.StorageLocationS3CompatDetails{},
		IsSubStruct:  true,
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.ApiIntegration{},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.ApiIntegrationAwsDetails{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "ApiIntegration", MethodName: "DescribeAws"},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.ApiIntegrationAzureDetails{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "ApiIntegration", MethodName: "DescribeAzure"},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.ApiIntegrationExternalMcpDetails{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "ApiIntegration", MethodName: "DescribeExternalMcp"},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.ApiIntegrationGitHttpsApiDetails{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "ApiIntegration", MethodName: "DescribeGitHttpsApi"},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.ApiIntegrationGoogleDetails{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "ApiIntegration", MethodName: "DescribeGoogle"},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.ApiIntegrationAllDetails{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "ApiIntegration", MethodName: "DescribeAllDetails"},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.CatalogIntegration{},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.CatalogIntegrationAwsGlueDetails{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "CatalogIntegration", MethodName: "DescribeAwsGlue"},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.CatalogIntegrationObjectStorageDetails{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "CatalogIntegration", MethodName: "DescribeObjectStorage"},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.CatalogIntegrationOpenCatalogDetails{},
		IsDataSourceOutput: true,
		NestedAssertFields: []string{"RestConfig", "RestAuthentication"},
		DescribeOverride: &genhelpers.DescribeOverrideDef{
			ClientName: "CatalogIntegration",
			MethodName: "DescribeOpenCatalog",
		},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.CatalogIntegrationIcebergRestDetails{},
		IsDataSourceOutput: true,
		NestedAssertFields: []string{"RestConfig", "OAuthRestAuthentication", "SigV4RestAuthentication"},
		SkipFields:         []string{"BearerRestAuthentication"},
		DescribeOverride: &genhelpers.DescribeOverrideDef{
			ClientName: "CatalogIntegration",
			MethodName: "DescribeIcebergRest",
		},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.CatalogIntegrationAllDetails{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "CatalogIntegration", MethodName: "DescribeDetails"},
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.OpenCatalogRestConfigDetails{},
		IsDataSourceOutput: true,
		IsSubStruct:        true,
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.IcebergRestRestConfigDetails{},
		IsDataSourceOutput: true,
		IsSubStruct:        true,
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.OAuthRestAuthenticationDetails{},
		IsDataSourceOutput: true,
		IsSubStruct:        true,
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.SigV4RestAuthenticationDetails{},
		IsDataSourceOutput: true,
		IsSubStruct:        true,
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.SessionPolicy{},
	},
	{
		IdType:             "sdk.SchemaObjectIdentifier",
		ObjectStruct:       sdk.SessionPolicyDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:             "sdk.SchemaObjectIdentifier",
		ObjectStruct:       sdk.PasswordPolicyDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:               "sdk.SchemaObjectIdentifier",
		ObjectStruct:         sdk.IcebergTableDetails{},
		IsDataSourceOutput:   true,
		NoIdentifiableObject: true,
	},
	{
		IdType:               "sdk.SchemaObjectIdentifier",
		ObjectStruct:         sdk.TableSearchOptimizationDetails{},
		IsDataSourceOutput:   true,
		NoIdentifiableObject: true,
	},
	{
		IdType:               "sdk.SchemaObjectIdentifier",
		ObjectStruct:         sdk.TableConstraintDetails{},
		IsDataSourceOutput:   true,
		NoIdentifiableObject: true,
	},
	{
		IdType:               "sdk.SchemaObjectIdentifier",
		ObjectStruct:         sdk.TableCheckConstraintDetails{},
		IsDataSourceOutput:   true,
		NoIdentifiableObject: true,
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.PasswordPolicy{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.HybridTable{},
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.IcebergTable{},
		SkipFields:   []string{"AutoRefreshStatus"},
	},
	{
		IdType:       "sdk.AccountObjectIdentifier",
		ObjectStruct: sdk.PostgresInstance{},
	},
	{
		IdType:         "sdk.SchemaObjectIdentifier",
		ObjectStruct:   sdk.CortexAgent{},
		ObjectTypeName: "Agent",
	},
	{
		IdType:             "sdk.SchemaObjectIdentifier",
		ObjectStruct:       sdk.CortexAgentDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.McpServer{},
	},
	{
		IdType:             "sdk.SchemaObjectIdentifier",
		ObjectStruct:       sdk.McpServerDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.TagReference{},
		NoShowById:   true,
	},
	{
		IdType:             "sdk.AccountObjectIdentifier",
		ObjectStruct:       sdk.PostgresInstanceDetails{},
		IsDataSourceOutput: true,
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.PolicyReference{},
		NoShowById:   true,
	},
	{
		IdType:       "sdk.SchemaObjectIdentifier",
		ObjectStruct: sdk.FileFormat{},
	},
	{
		IdType:             "sdk.SchemaObjectIdentifier",
		ObjectStruct:       sdk.FileFormatCsv{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "FileFormat", MethodName: "DescribeCsvDetails"},
	},
	{
		IdType:             "sdk.SchemaObjectIdentifier",
		ObjectStruct:       sdk.FileFormatJson{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "FileFormat", MethodName: "DescribeJsonDetails"},
	},
	{
		IdType:             "sdk.SchemaObjectIdentifier",
		ObjectStruct:       sdk.FileFormatAvro{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "FileFormat", MethodName: "DescribeAvroDetails"},
	},
	{
		IdType:             "sdk.SchemaObjectIdentifier",
		ObjectStruct:       sdk.FileFormatOrc{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "FileFormat", MethodName: "DescribeOrcDetails"},
	},
	{
		IdType:             "sdk.SchemaObjectIdentifier",
		ObjectStruct:       sdk.FileFormatParquet{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "FileFormat", MethodName: "DescribeParquetDetails"},
	},
	{
		IdType:             "sdk.SchemaObjectIdentifier",
		ObjectStruct:       sdk.FileFormatXml{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "FileFormat", MethodName: "DescribeXmlDetails"},
	},
	{
		IdType:             "sdk.SchemaObjectIdentifier",
		ObjectStruct:       sdk.FileFormatAllDetails{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "FileFormat", MethodName: "DescribeAllDetails"},
	},
	{
		IdType:               "sdk.AccountObjectIdentifier",
		ObjectStruct:         sdk.DatabaseDetails{},
		IsDataSourceOutput:   true,
		NoIdentifiableObject: true,
	},
	{
		IdType:             "sdk.SchemaObjectIdentifier",
		ObjectStruct:       sdk.SemanticViewDetails{},
		IsDataSourceOutput: true,
		SkipFields:         []string{"Id", "Tables", "Relationships", "Dimensions", "Facts", "Metrics"},
	},
	{
		IdType:             "sdk.SchemaObjectIdentifier",
		ObjectStruct:       sdk.StageDetails{},
		IsDataSourceOutput: true,
		DescribeOverride:   &genhelpers.DescribeOverrideDef{ClientName: "Stage", MethodName: "DescribeDetails"},
	},
}

func GetSdkObjectDetails() []genhelpers.SdkObjectDetails {
	allSdkObjectsDetails := make([]genhelpers.SdkObjectDetails, len(allStructs))
	for idx, d := range allStructs {
		structDetails := genhelpers.ExtractStructDetails(d.ObjectStruct)
		allSdkObjectsDetails[idx] = genhelpers.SdkObjectDetails{
			IdType:               d.IdType,
			StructDetails:        structDetails,
			IsDataSourceOutput:   d.IsDataSourceOutput,
			IsSubStruct:          d.IsSubStruct,
			ObjectTypeName:       d.ObjectTypeName,
			NoShowById:           d.NoShowById,
			NoIdentifiableObject: d.NoIdentifiableObject,
			ShowByParentId:       d.ShowByParentId,
			DescribeOverride:     d.DescribeOverride,
			FromObjectIDExpr:     d.FromObjectIDExpr,
			NestedAssertFields:   d.NestedAssertFields,
			SkipFields:           d.SkipFields,
		}
	}
	return allSdkObjectsDetails
}
