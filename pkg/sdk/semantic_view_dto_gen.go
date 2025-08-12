package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateSemanticViewOptions]   = new(CreateSemanticViewRequest)
	_ optionsProvider[DropSemanticViewOptions]     = new(DropSemanticViewRequest)
	_ optionsProvider[DescribeSemanticViewOptions] = new(DescribeSemanticViewRequest)
	_ optionsProvider[ShowSemanticViewsOptions]    = new(ShowSemanticViewRequest)
)

type CreateSemanticViewRequest struct {
	OrReplace   *bool
	IfNotExists *bool
	name        SchemaObjectIdentifier // required
	Tables      []LogicalTable         // required
	Comment     *string
	CopyGrants  *bool
}

type DropSemanticViewRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type DescribeSemanticViewRequest struct {
	name SchemaObjectIdentifier // required
}

type ShowSemanticViewRequest struct {
	Terse      *bool
	Like       *Like
	In         *In
	StartsWith *string
	Limit      *LimitFrom
}
