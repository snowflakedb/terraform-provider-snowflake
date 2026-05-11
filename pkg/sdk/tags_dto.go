package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[createTagOptions] = new(CreateTagRequest)
	_ optionsProvider[alterTagOptions]  = new(AlterTagRequest)
	_ optionsProvider[showTagOptions]   = new(ShowTagRequest)
	_ optionsProvider[dropTagOptions]   = new(DropTagRequest)
	_ optionsProvider[undropTagOptions] = new(UndropTagRequest)
	_ optionsProvider[setTagOptions]    = new(SetTagRequest)
)

type SetTagRequest struct {
	objectType ObjectType       // required
	objectName ObjectIdentifier // required

	SetTags []TagAssociation
}

type UnsetTagRequest struct {
	objectType ObjectType       // required
	objectName ObjectIdentifier // required

	IfExists  *bool
	UnsetTags []ObjectIdentifier
}

type SetTagOnCurrentAccountRequest struct {
	SetTags []TagAssociation
}

type UnsetTagOnCurrentAccountRequest struct {
	UnsetTags []ObjectIdentifier
}

type TagPropagateRequest struct {
	propagationMethod TagPropagation // required
	onConflict        *TagOnConflict
}

type CreateTagRequest struct {
	orReplace   *bool
	ifNotExists *bool

	name SchemaObjectIdentifier // required

	comment       *string
	allowedValues []string
	propagate     *TagPropagateRequest
}

func (r *CreateTagRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

type AlterTagRequest struct {
	ifExists *bool
	name     SchemaObjectIdentifier // required

	// One of
	add    []string
	drop   []string
	set    *TagSetRequest
	unset  *TagUnsetRequest
	rename *SchemaObjectIdentifier
}

type TagSetRequest struct {
	maskingPolicies []SchemaObjectIdentifier
	force           *bool
	allowedValues   []string
	propagate       *TagPropagateRequest
	comment         *string
}

type TagUnsetRequest struct {
	maskingPolicies []SchemaObjectIdentifier
	allowedValues   *bool
	propagate       *bool
	onConflict      *bool
	comment         *bool
}

type ShowTagRequest struct {
	Like *Like
	In   *ExtendedIn
}

type DropTagRequest struct {
	ifExists bool

	name SchemaObjectIdentifier // required
}

type UndropTagRequest struct {
	name SchemaObjectIdentifier // required
}
