package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateCortexSearchServiceOptions]   = new(CreateCortexSearchServiceRequest)
	_ optionsProvider[AlterCortexSearchServiceOptions]    = new(AlterCortexSearchServiceRequest)
	_ optionsProvider[ShowCortexSearchServiceOptions]     = new(ShowCortexSearchServiceRequest)
	_ optionsProvider[DescribeCortexSearchServiceOptions] = new(DescribeCortexSearchServiceRequest)
	_ optionsProvider[DropCortexSearchServiceOptions]     = new(DropCortexSearchServiceRequest)
)

type CreateCortexSearchServiceRequest struct {
	OrReplace       *bool
	IfNotExists     *bool
	name            SchemaObjectIdentifier // required
	On              string                 // required
	Attributes      *AttributesRequest
	Warehouse       AccountObjectIdentifier // required
	TargetLag       string                  // required
	EmbeddingModel  *string
	Comment         *string
	QueryDefinition string // required
}

type AttributesRequest struct {
	Columns []string
}

type AlterCortexSearchServiceRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
	Set      *CortexSearchServiceSetRequest
}

type CortexSearchServiceSetRequest struct {
	TargetLag *string
	Warehouse *AccountObjectIdentifier
	Comment   *string
}

type ShowCortexSearchServiceRequest struct {
	Like       *Like
	In         *In
	StartsWith *string
	Limit      *LimitFrom
}

type DescribeCortexSearchServiceRequest struct {
	name SchemaObjectIdentifier // required
}

type DropCortexSearchServiceRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}
