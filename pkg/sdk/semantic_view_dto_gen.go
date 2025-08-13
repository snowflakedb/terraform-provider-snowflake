package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateSemanticViewOptions]   = new(CreateSemanticViewRequest)
	_ optionsProvider[DropSemanticViewOptions]     = new(DropSemanticViewRequest)
	_ optionsProvider[DescribeSemanticViewOptions] = new(DescribeSemanticViewRequest)
	_ optionsProvider[ShowSemanticViewOptions]     = new(ShowSemanticViewRequest)
)

type CreateSemanticViewRequest struct {
	OrReplace     *bool
	IfNotExists   *bool
	name          SchemaObjectIdentifier // required
	logicalTables []LogicalTableRequest  // required
	Comment       *string
	CopyGrants    *bool
}

type LogicalTableRequest struct {
	logicalTableAlias *LogicalTableAliasRequest
	TableName         SchemaObjectIdentifier // required
	primaryKeys       *PrimaryKeysRequest
	uniqueKeys        []UniqueKeysRequest
	synonyms          *SynonymsRequest
	Comment           *string
}

type LogicalTableAliasRequest struct {
	LogicalTableAlias string
}

type PrimaryKeysRequest struct {
	PrimaryKey []SemanticViewColumn
}

type UniqueKeysRequest struct {
	Unique []SemanticViewColumn
}

type SynonymsRequest struct {
	WithSynonyms []string
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
