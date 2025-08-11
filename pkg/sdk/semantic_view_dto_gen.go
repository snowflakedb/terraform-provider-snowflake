package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateSemanticViewOptions]   = new(CreateSemanticViewRequest)
	_ optionsProvider[DropSemanticViewOptions]     = new(DropSemanticViewRequest)
	_ optionsProvider[DescribeSemanticViewOptions] = new(DescribeSemanticViewRequest)
	_ optionsProvider[ShowSemanticViewsOptions]    = new(ShowSemanticViewsRequest)
)

type CreateSemanticViewRequest struct {
	OrReplace     *bool
	IfNotExists   *bool
	name          SchemaObjectIdentifier   // required
	tables        []SchemaObjectIdentifier // required
	relationships []SchemaObjectIdentifier
	facts         []SchemaObjectIdentifier
	dimensions    []SchemaObjectIdentifier
	metrics       []SchemaObjectIdentifier
	Comment       *string
	CopyGrants    *bool
}

type DropSemanticViewRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type DescribeSemanticViewRequest struct {
	name SchemaObjectIdentifier // required
}

type ShowSemanticViewsRequest struct {
	Like       *Like
	Terse      *bool
	In         *In
	StartsWith *string
	Limit      *LimitFrom
}
