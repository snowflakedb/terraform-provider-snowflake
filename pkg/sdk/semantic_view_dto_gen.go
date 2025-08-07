package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateSemanticViewOptions] = new(CreateSemanticViewRequest)
	_ optionsProvider[DropSemanticViewOptions]   = new(DropSemanticViewRequest)
)

type CreateSemanticViewRequest struct {
	OrReplace   *bool
	IfNotExists *bool
	name        SchemaObjectIdentifier // required
	Comment     *string
	CopyGrants  *bool
}

type DropSemanticViewRequest struct {
	OrReplace   *bool
	IfNotExists *bool
	name        SchemaObjectIdentifier // required
	Comment     *string
	CopyGrants  *bool
}
