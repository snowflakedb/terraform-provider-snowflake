package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateNotebookOptions]   = new(CreateNotebookRequest)
	_ optionsProvider[AlterNotebookOptions]    = new(AlterNotebookRequest)
	_ optionsProvider[DropNotebookOptions]     = new(DropNotebookRequest)
	_ optionsProvider[DescribeNotebookOptions] = new(DescribeNotebookRequest)
	_ optionsProvider[ShowNotebookOptions]     = new(ShowNotebookRequest)
)

type CreateNotebookRequest struct {
	OrReplace                   *bool
	IfNotExists                 *bool
	name                        SchemaObjectIdentifier // required
	From                        *Location
	MainFile                    *string
	Comment                     *string
	QueryWarehouse              *AccountObjectIdentifier
	IdleAutoShutdownTimeSeconds *int
	Warehouse                   *AccountObjectIdentifier
	RuntimeName                 *string
	ComputePool                 *AccountObjectIdentifier
	ExternalAccessIntegrations  []AccountObjectIdentifier
	RuntimeEnvironmentVersion   *string
	DefaultVersion              *string
}

type AlterNotebookRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
	RenameTo *SchemaObjectIdentifier
	Set      *NotebookSetRequest
	Unset    *NotebookUnsetRequest
}

type NotebookSetRequest struct {
	Comment                     *string
	QueryWarehouse              *AccountObjectIdentifier
	IdleAutoShutdownTimeSeconds *int
	SecretsList                 *SecretsListRequest
	MainFile                    *string
	Warehouse                   *AccountObjectIdentifier
	RuntimeName                 *string
	ComputePool                 *AccountObjectIdentifier
	ExternalAccessIntegrations  []AccountObjectIdentifier
	RuntimeEnvironmentVersion   *string
}

type NotebookUnsetRequest struct {
	Comment                    *bool
	QueryWarehouse             *bool
	Secrets                    *bool
	Warehouse                  *bool
	RuntimeName                *bool
	ComputePool                *bool
	ExternalAccessIntegrations *bool
	RuntimeEnvironmentVersion  *bool
}

type DropNotebookRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type DescribeNotebookRequest struct {
	name SchemaObjectIdentifier // required
}

type ShowNotebookRequest struct {
	Like       *Like
	In         *In
	Limit      *LimitFrom
	StartsWith *string
}
