package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateGitRepositoryOptions]   = new(CreateGitRepositoryRequest)
	_ optionsProvider[AlterGitRepositoryOptions]    = new(AlterGitRepositoryRequest)
	_ optionsProvider[DropGitRepositoryOptions]     = new(DropGitRepositoryRequest)
	_ optionsProvider[DescribeGitRepositoryOptions] = new(DescribeGitRepositoryRequest)
	_ optionsProvider[ShowGitRepositoryOptions]     = new(ShowGitRepositoryRequest)
)

type CreateGitRepositoryRequest struct {
	OrReplace      *bool
	IfNotExists    *bool
	name           SchemaObjectIdentifier  // required
	Origin         string                  // required
	ApiIntegration AccountObjectIdentifier // required
	GitCredentials *AccountObjectIdentifier
	Comment        *string
	Tag            []TagAssociation
}

type AlterGitRepositoryRequest struct {
	IfExists  *bool
	name      SchemaObjectIdentifier // required
	Set       *GitRepositorySetRequest
	Unset     *GitRepositoryUnsetRequest
	Fetch     *bool
	SetTags   []TagAssociation
	UnsetTags []ObjectIdentifier
}

type GitRepositorySetRequest struct {
	ApiIntegration *AccountObjectIdentifier // required
	GitCredentials *AccountObjectIdentifier
	Comment        *string
}

type GitRepositoryUnsetRequest struct {
	GitCredentials *bool
	Comment        *bool
}

type DropGitRepositoryRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type DescribeGitRepositoryRequest struct {
	name SchemaObjectIdentifier // required
}

type ShowGitRepositoryRequest struct {
	Like *Like
	In   *In
}
