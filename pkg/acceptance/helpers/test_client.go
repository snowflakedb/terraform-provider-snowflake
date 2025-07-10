package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type AccountInformation interface {
	GetAccountLocator() string
}

type TestClient struct {
	context *TestClientContext

	Ids *IdsGenerator
	AccountInformation

	Account                      *AccountClient
	AggregationPolicy            *AggregationPolicyClient
	Alert                        *AlertClient
	ApiIntegration               *ApiIntegrationClient
	Application                  *ApplicationClient
	ApplicationPackage           *ApplicationPackageClient
	AuthenticationPolicy         *AuthenticationPolicyClient
	BcrBundles                   *BcrBundlesClient
	ComputePool                  *ComputePoolClient
	Connection                   *ConnectionClient
	Context                      *ContextClient
	CortexSearchService          *CortexSearchServiceClient
	CatalogIntegration           *CatalogIntegrationClient
	Database                     *DatabaseClient
	DatabaseRole                 *DatabaseRoleClient
	DataMetricFunctionClient     *DataMetricFunctionClient
	DataMetricFunctionReferences *DataMetricFunctionReferencesClient
	DynamicTable                 *DynamicTableClient
	EventTable                   *EventTableClient
	ExternalAccessIntegration    *ExternalAccessIntegrationClient
	ExternalFunction             *ExternalFunctionClient
	ExternalTable                *ExternalTableClient
	ExternalVolume               *ExternalVolumeClient
	FailoverGroup                *FailoverGroupClient
	FeaturePolicy                *FeaturePolicyClient
	FileFormat                   *FileFormatClient
	Function                     *FunctionClient
	GitRepository                *GitRepositoryClient
	Grant                        *GrantClient
	HybridTable                  *HybridTableClient
	ImageRepository              *ImageRepositoryClient
	InformationSchema            *InformationSchemaClient
	MaskingPolicy                *MaskingPolicyClient
	MaterializedView             *MaterializedViewClient
	NetworkPolicy                *NetworkPolicyClient
	NetworkRule                  *NetworkRuleClient
	NotificationIntegration      *NotificationIntegrationClient
	PackagesPolicy               *PackagesPolicyClient
	Parameter                    *ParameterClient
	PasswordPolicy               *PasswordPolicyClient
	Pipe                         *PipeClient
	Procedure                    *ProcedureClient
	ProjectionPolicy             *ProjectionPolicyClient
	PolicyReferences             *PolicyReferencesClient
	ResourceMonitor              *ResourceMonitorClient
	Role                         *RoleClient
	RowAccessPolicy              *RowAccessPolicyClient
	Schema                       *SchemaClient
	Secret                       *SecretClient
	SecurityIntegration          *SecurityIntegrationClient
	Service                      *ServiceClient
	Sequence                     *SequenceClient
	SessionPolicy                *SessionPolicyClient
	Share                        *ShareClient
	Snapshot                     *SnapshotClient
	Stage                        *StageClient
	StorageIntegration           *StorageIntegrationClient
	Stream                       *StreamClient
	Streamlit                    *StreamlitClient
	Table                        *TableClient
	Tag                          *TagClient
	Task                         *TaskClient
	User                         *UserClient
	View                         *ViewClient
	Warehouse                    *WarehouseClient
}

func NewTestClient(c *sdk.Client, database string, schema string, warehouse string, testObjectSuffix string) *TestClient {
	context := &TestClientContext{
		client:           c,
		database:         database,
		schema:           schema,
		warehouse:        warehouse,
		testObjectSuffix: testObjectSuffix,
	}
	idsGenerator := NewIdsGenerator(context)
	return &TestClient{
		context: context,

		Ids:                idsGenerator,
		AccountInformation: context.client,

		Account:                      NewAccountClient(context, idsGenerator),
		AggregationPolicy:            NewAggregationPolicyClient(context, idsGenerator),
		Alert:                        NewAlertClient(context, idsGenerator),
		ApiIntegration:               NewApiIntegrationClient(context, idsGenerator),
		Application:                  NewApplicationClient(context, idsGenerator),
		ApplicationPackage:           NewApplicationPackageClient(context, idsGenerator),
		AuthenticationPolicy:         NewAuthenticationPolicyClient(context, idsGenerator),
		BcrBundles:                   NewBcrBundlesClient(context),
		ComputePool:                  NewComputePoolClient(context, idsGenerator),
		Connection:                   NewConnectionClient(context, idsGenerator),
		Context:                      NewContextClient(context),
		CortexSearchService:          NewCortexSearchServiceClient(context, idsGenerator),
		CatalogIntegration:           NewCatalogIntegrationClient(context, idsGenerator),
		Database:                     NewDatabaseClient(context, idsGenerator),
		DatabaseRole:                 NewDatabaseRoleClient(context, idsGenerator),
		DataMetricFunctionClient:     NewDataMetricFunctionClient(context, idsGenerator),
		DataMetricFunctionReferences: NewDataMetricFunctionReferencesClient(context),
		DynamicTable:                 NewDynamicTableClient(context, idsGenerator),
		EventTable:                   NewEventTableClient(context, idsGenerator),
		ExternalAccessIntegration:    NewExternalAccessIntegrationClient(context, idsGenerator),
		ExternalFunction:             NewExternalFunctionClient(context, idsGenerator),
		ExternalTable:                NewExternalTableClient(context, idsGenerator),
		ExternalVolume:               NewExternalVolumeClient(context, idsGenerator),
		FailoverGroup:                NewFailoverGroupClient(context, idsGenerator),
		FeaturePolicy:                NewFeaturePolicyClient(context, idsGenerator),
		FileFormat:                   NewFileFormatClient(context, idsGenerator),
		Function:                     NewFunctionClient(context, idsGenerator),
		GitRepository:                NewGitRepositoryClient(context, idsGenerator),
		Grant:                        NewGrantClient(context, idsGenerator),
		HybridTable:                  NewHybridTableClient(context, idsGenerator),
		ImageRepository:              NewImageRepositoryClient(context, idsGenerator),
		InformationSchema:            NewInformationSchemaClient(context, idsGenerator),
		MaskingPolicy:                NewMaskingPolicyClient(context, idsGenerator),
		MaterializedView:             NewMaterializedViewClient(context, idsGenerator),
		NetworkPolicy:                NewNetworkPolicyClient(context, idsGenerator),
		NetworkRule:                  NewNetworkRuleClient(context, idsGenerator),
		NotificationIntegration:      NewNotificationIntegrationClient(context, idsGenerator),
		PackagesPolicy:               NewPackagesPolicyClient(context, idsGenerator),
		Parameter:                    NewParameterClient(context),
		PasswordPolicy:               NewPasswordPolicyClient(context, idsGenerator),
		Pipe:                         NewPipeClient(context, idsGenerator),
		Procedure:                    NewProcedureClient(context, idsGenerator),
		ProjectionPolicy:             NewProjectionPolicyClient(context, idsGenerator),
		PolicyReferences:             NewPolicyReferencesClient(context),
		ResourceMonitor:              NewResourceMonitorClient(context, idsGenerator),
		Role:                         NewRoleClient(context, idsGenerator),
		RowAccessPolicy:              NewRowAccessPolicyClient(context, idsGenerator),
		Schema:                       NewSchemaClient(context, idsGenerator),
		Secret:                       NewSecretClient(context, idsGenerator),
		SecurityIntegration:          NewSecurityIntegrationClient(context, idsGenerator),
		Snapshot:                     NewSnapshotClient(context, idsGenerator),
		Service:                      NewServiceClient(context, idsGenerator),
		Sequence:                     NewSequenceClient(context, idsGenerator),
		SessionPolicy:                NewSessionPolicyClient(context, idsGenerator),
		Share:                        NewShareClient(context, idsGenerator),
		Stage:                        NewStageClient(context, idsGenerator),
		StorageIntegration:           NewStorageIntegrationClient(context, idsGenerator),
		Stream:                       NewStreamClient(context, idsGenerator),
		Streamlit:                    NewStreamlitClient(context, idsGenerator),
		Table:                        NewTableClient(context, idsGenerator),
		Tag:                          NewTagClient(context, idsGenerator),
		Task:                         NewTaskClient(context, idsGenerator),
		User:                         NewUserClient(context, idsGenerator),
		View:                         NewViewClient(context, idsGenerator),
		Warehouse:                    NewWarehouseClient(context, idsGenerator),
	}
}

type TestClientContext struct {
	client           *sdk.Client
	database         string
	schema           string
	warehouse        string
	testObjectSuffix string
}
