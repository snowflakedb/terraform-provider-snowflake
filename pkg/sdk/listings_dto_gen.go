package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateListingOptions]          = new(CreateListingRequest)
	_ optionsProvider[CreateFromStageListingOptions] = new(CreateFromStageListingRequest)
	_ optionsProvider[AlterListingOptions]           = new(AlterListingRequest)
	_ optionsProvider[DropListingOptions]            = new(DropListingRequest)
	_ optionsProvider[ShowListingOptions]            = new(ShowListingRequest)
	_ optionsProvider[DescribeListingOptions]        = new(DescribeListingRequest)
)

type CreateListingRequest struct {
	IfNotExists *bool
	name        AccountObjectIdentifier // required
	With        *ListingWithRequest
	As          string // required
	Publish     *bool
	Review      *bool
	Comment     *string
}

type ListingWithRequest struct {
	Share              *AccountObjectIdentifier
	ApplicationPackage *AccountObjectIdentifier
}

type CreateFromStageListingRequest struct {
	IfNotExists *bool
	name        AccountObjectIdentifier // required
	With        *ListingWithRequest
	From        Location // required
	Publish     *bool
	Review      *bool
}

type AlterListingRequest struct {
	IfExists       *bool
	name           AccountObjectIdentifier // required
	Publish        *bool
	Unpublish      *bool
	Review         *bool
	AlterListingAs *AlterListingAsRequest
	AddVersion     *AddListingVersionRequest
	RenameTo       *AccountObjectIdentifier
	Set            *ListingSetRequest
}

type AlterListingAsRequest struct {
	As      string // required
	Publish *bool
	Review  *bool
	Comment *string
}

type AddListingVersionRequest struct {
	IfNotExists *bool
	VersionName string   // required
	From        Location // required
	Comment     *string
}

type ListingSetRequest struct {
	Comment *string
}

type DropListingRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type ShowListingRequest struct {
	Like       *Like
	StartsWith *string
	Limit      *LimitFrom
}

type DescribeListingRequest struct {
	name     AccountObjectIdentifier // required
	Revision *ListingRevision
}
