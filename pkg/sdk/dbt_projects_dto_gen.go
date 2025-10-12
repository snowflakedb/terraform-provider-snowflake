package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateDbtProjectOptions]   = new(CreateDbtProjectRequest)
	_ optionsProvider[AlterDbtProjectOptions]    = new(AlterDbtProjectRequest)
	_ optionsProvider[DropDbtProjectOptions]     = new(DropDbtProjectRequest)
	_ optionsProvider[ShowDbtProjectOptions]     = new(ShowDbtProjectRequest)
	_ optionsProvider[DescribeDbtProjectOptions] = new(DescribeDbtProjectRequest)
)

type CreateDbtProjectRequest struct {
	OrReplace      *bool
	IfNotExists    *bool
	name           SchemaObjectIdentifier // required
	From           *string
	DefaultArgs    *string
	DefaultVersion *DbtProjectDefaultVersion
	Comment        *string
}

type AlterDbtProjectRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
	Set      *DbtProjectSetRequest
	Unset    *DbtProjectUnsetRequest
}

type DbtProjectSetRequest struct {
	DefaultArgs    *string
	DefaultVersion *DbtProjectDefaultVersion
	Comment        *string
}

type DbtProjectUnsetRequest struct {
	DefaultArgs    *bool
	DefaultVersion *bool
	Comment        *bool
}

type DropDbtProjectRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type ShowDbtProjectRequest struct {
	Like *Like
	In   *In
}

type DescribeDbtProjectRequest struct {
	name SchemaObjectIdentifier // required
}
