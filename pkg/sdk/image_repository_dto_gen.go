package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateImageRepositoryOptions] = new(CreateImageRepositoryRequest)
	_ optionsProvider[AlterImageRepositoryOptions]  = new(AlterImageRepositoryRequest)
	_ optionsProvider[DropImageRepositoryOptions]   = new(DropImageRepositoryRequest)
	_ optionsProvider[ShowImageRepositoryOptions]   = new(ShowImageRepositoryRequest)
)

type CreateImageRepositoryRequest struct {
	OrReplace   *bool
	IfNotExists *bool
	name        SchemaObjectIdentifier // required
	Comment     *string
}

type AlterImageRepositoryRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
	Set      *ImageRepositorySetRequest
}

type ImageRepositorySetRequest struct {
	Comment *StringAllowEmpty
}

type DropImageRepositoryRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type ShowImageRepositoryRequest struct {
	Like *Like
	In   *In
}
