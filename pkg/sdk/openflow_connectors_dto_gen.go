package sdk

var (
	_ optionsProvider[CreateOpenflowConnectorOptions]   = new(CreateOpenflowConnectorRequest)
	_ optionsProvider[AlterOpenflowConnectorOptions]    = new(AlterOpenflowConnectorRequest)
	_ optionsProvider[DropOpenflowConnectorOptions]     = new(DropOpenflowConnectorRequest)
	_ optionsProvider[ShowOpenflowConnectorOptions]     = new(ShowOpenflowConnectorRequest)
	_ optionsProvider[DescribeOpenflowConnectorOptions] = new(DescribeOpenflowConnectorRequest)
)

type CreateOpenflowConnectorRequest struct {
	name           SchemaObjectIdentifier // required
	InRuntime      SchemaObjectIdentifier // required
	IfNotExists    *bool
	FromDefinition *string
	FromStage      *string
	DisplayName    *string
	Comment        *string
}

type AlterOpenflowConnectorRequest struct {
	name  SchemaObjectIdentifier // required
	Start *bool
	Stop  *bool
	Set   *OpenflowConnectorSetRequest
	Unset *OpenflowConnectorUnsetRequest
}

type OpenflowConnectorSetRequest struct {
	DisplayName *string
	Comment     *string
}

type OpenflowConnectorUnsetRequest struct {
	DisplayName *bool
	Comment     *bool
}

type DropOpenflowConnectorRequest struct {
	name     SchemaObjectIdentifier // required
	IfExists *bool
}

type ShowOpenflowConnectorRequest struct {
	Like *Like
}

type DescribeOpenflowConnectorRequest struct {
	name SchemaObjectIdentifier // required
}
