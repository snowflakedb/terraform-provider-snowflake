package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowResultSchemaDef is the input definition for show/describe schema generation.
type ShowResultSchemaDef struct {
	ObjectStruct any
	// IsDescribe generates DescribeXSchema and {snake}_desc_gen.go (trims a trailing _details).
	IsDescribe bool
	// SkipFields omits snake_case schema keys from both the schema map and mapping.
	SkipFields []string
	// TypeOverrides maps snake_case schema keys to a Terraform schema type.
	// Currently, no logic is implemented.
	TypeOverrides map[string]schema.ValueType
	// UsedAsListEntry generates NameSchema (no Show/Describe prefix) for a property-row list Elem.
	UsedAsListEntry bool
	// AdditionalMapping generates a mapper type that must implement additionalSchemaMapper[T] in *_ext.go. Not inferred from SkipFields.
	AdditionalMapping bool
}

// ShowResultSchemaDetails is the extracted generator input (struct fields + definition metadata).
type ShowResultSchemaDetails struct {
	IsDescribe        bool
	SkipFields        []string
	TypeOverrides     map[string]schema.ValueType
	UsedAsListEntry   bool
	AdditionalMapping bool
	genhelpers.StructDetails
}

var SdkShowResultStructs = []ShowResultSchemaDef{
	{ObjectStruct: sdk.Account{}},
	{ObjectStruct: sdk.Alert{}},
	{ObjectStruct: sdk.ApiIntegration{}},
	{ObjectStruct: sdk.ApplicationPackage{}},
	{ObjectStruct: sdk.ApplicationRole{}},
	{ObjectStruct: sdk.Application{}},
	{ObjectStruct: sdk.AuthenticationPolicy{}, SkipFields: []string{"target_scopes"}},
	{ObjectStruct: sdk.CatalogIntegration{}},
	{ObjectStruct: sdk.ComputePool{}, SkipFields: []string{"backup_instance_families"}, AdditionalMapping: true},
	{ObjectStruct: sdk.Connection{}, SkipFields: []string{"failover_allowed_to_accounts"}, AdditionalMapping: true},
	{ObjectStruct: sdk.CortexAgent{}, SkipFields: []string{"profile"}, AdditionalMapping: true},
	{ObjectStruct: sdk.DatabaseRole{}},
	{ObjectStruct: sdk.Database{}},
	{ObjectStruct: sdk.DynamicTable{}},
	{ObjectStruct: sdk.EventTable{}},
	{ObjectStruct: sdk.ExternalAccessIntegration{}},
	{ObjectStruct: sdk.ExternalFunction{}},
	{ObjectStruct: sdk.ExternalTable{}},
	{ObjectStruct: sdk.ExternalVolume{}},
	{ObjectStruct: sdk.FailoverGroup{}},
	{ObjectStruct: sdk.FileFormat{}},
	{ObjectStruct: sdk.FileFormatLegacy{}},
	{ObjectStruct: sdk.Function{}, SkipFields: []string{"arguments_old", "return_type_old"}},
	{ObjectStruct: sdk.GitRepository{}},
	{ObjectStruct: sdk.Grant{}, SkipFields: []string{"grant_on", "grant_to"}},
	{ObjectStruct: sdk.HybridTable{}},
	{ObjectStruct: sdk.HybridTableConstraint{}, SkipFields: []string{"columns", "referenced_table", "referenced_columns", "delete_rule", "update_rule"}, AdditionalMapping: true},
	{ObjectStruct: sdk.HybridTableIndex{}},
	{ObjectStruct: sdk.IcebergTable{}, SkipFields: []string{"auto_refresh_status", "partition_specs"}, AdditionalMapping: true},
	{ObjectStruct: sdk.ImageRepository{}},
	{ObjectStruct: sdk.Listing{}},
	{ObjectStruct: sdk.ManagedAccount{}},
	{ObjectStruct: sdk.MaskingPolicy{}, SkipFields: []string{"options"}},
	{ObjectStruct: sdk.MaterializedView{}},
	{ObjectStruct: sdk.McpServer{}},
	{ObjectStruct: sdk.NetworkPolicy{}},
	{ObjectStruct: sdk.NetworkRule{}},
	{ObjectStruct: sdk.Notebook{}},
	{ObjectStruct: sdk.NotificationIntegration{}},
	{ObjectStruct: sdk.OpenflowConnectorDefinition{}, SkipFields: []string{"categories"}, AdditionalMapping: true},
	{ObjectStruct: sdk.OpenflowConnector{}},
	{ObjectStruct: sdk.OpenflowDeployment{}},
	{ObjectStruct: sdk.OpenflowRuntime{}, SkipFields: []string{"external_access_integrations"}, AdditionalMapping: true},
	{ObjectStruct: sdk.OrganizationAccount{}},
	{ObjectStruct: sdk.Parameter{}},
	{ObjectStruct: sdk.PasswordPolicy{}},
	{ObjectStruct: sdk.Pipe{}},
	{ObjectStruct: sdk.PolicyReference{}},
	{ObjectStruct: sdk.PostgresInstance{}, SkipFields: []string{"is_highly_available"}, AdditionalMapping: true},
	{ObjectStruct: sdk.Procedure{}, SkipFields: []string{"arguments_old", "return_type_old"}},
	{ObjectStruct: sdk.ReplicationAccount{}, SkipFields: []string{"comment"}, AdditionalMapping: true},
	{ObjectStruct: sdk.ReplicationDatabase{}},
	{ObjectStruct: sdk.Region{}},
	{ObjectStruct: sdk.ResourceMonitor{}, SkipFields: []string{"notify_at", "notify_users", "suspend_immediately_at"}, AdditionalMapping: true},
	{ObjectStruct: sdk.Role{}},
	{ObjectStruct: sdk.RowAccessPolicy{}},
	{ObjectStruct: sdk.Schema{}},
	{ObjectStruct: sdk.Secret{}, SkipFields: []string{"oauth_scopes"}, AdditionalMapping: true},
	{ObjectStruct: sdk.SecurityIntegration{}},
	{ObjectStruct: sdk.SemanticView{}},
	{ObjectStruct: sdk.Service{}, SkipFields: []string{"external_access_integrations"}, AdditionalMapping: true},
	{ObjectStruct: sdk.Sequence{}},
	{ObjectStruct: sdk.SessionPolicy{}, SkipFields: []string{"target_scopes"}},
	{ObjectStruct: sdk.Share{}, SkipFields: []string{"owner_account", "name"}, AdditionalMapping: true},
	{ObjectStruct: sdk.Stage{}},
	{ObjectStruct: sdk.StorageIntegration{}},
	{ObjectStruct: sdk.StorageLifecyclePolicy{}},
	{ObjectStruct: sdk.Streamlit{}},
	{ObjectStruct: sdk.Stream{}, SkipFields: []string{"base_tables"}, AdditionalMapping: true},
	{ObjectStruct: sdk.Table{}},
	{ObjectStruct: sdk.Tag{}, SkipFields: []string{"allowed_values"}, AdditionalMapping: true},
	{ObjectStruct: sdk.Task{}, SkipFields: []string{"predecessors", "task_relations", "target_completion_interval"}, AdditionalMapping: true},
	{ObjectStruct: sdk.User{}},
	{ObjectStruct: sdk.ProgrammaticAccessToken{}},
	{ObjectStruct: sdk.View{}},
	{ObjectStruct: sdk.Warehouse{}},
	{ObjectStruct: sdk.CatalogIntegrationAwsGlueDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.CatalogIntegrationObjectStorageDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.CortexAgentDetails{}, IsDescribe: true, SkipFields: []string{"profile"}, AdditionalMapping: true},
	{ObjectStruct: sdk.IcebergTableDetails{}, IsDescribe: true, SkipFields: []string{"type", "data_type_raw"}, AdditionalMapping: true},
	{ObjectStruct: sdk.McpServerDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.PasswordPolicyDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.SecurityIntegrationProperty{}, UsedAsListEntry: true},
	{ObjectStruct: sdk.SessionPolicyDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.StorageIntegrationAllDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.StorageIntegrationAwsDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.StorageIntegrationAzureDetails{}, IsDescribe: true},
	{ObjectStruct: sdk.StorageIntegrationGcsDetails{}, IsDescribe: true},
}

func GetShowResultSchemaDetails() []ShowResultSchemaDetails {
	allDetails := make([]ShowResultSchemaDetails, len(SdkShowResultStructs))
	for idx, d := range SdkShowResultStructs {
		allDetails[idx] = ShowResultSchemaDetails{
			IsDescribe:        d.IsDescribe,
			SkipFields:        d.SkipFields,
			TypeOverrides:     d.TypeOverrides,
			UsedAsListEntry:   d.UsedAsListEntry,
			AdditionalMapping: d.AdditionalMapping,
			StructDetails:     genhelpers.ExtractStructDetails(d.ObjectStruct),
		}
	}
	return allDetails
}
