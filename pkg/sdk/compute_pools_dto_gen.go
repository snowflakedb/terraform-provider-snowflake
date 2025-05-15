package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateComputePoolOptions]   = new(CreateComputePoolRequest)
	_ optionsProvider[AlterComputePoolOptions]    = new(AlterComputePoolRequest)
	_ optionsProvider[DropComputePoolOptions]     = new(DropComputePoolRequest)
	_ optionsProvider[ShowComputePoolOptions]     = new(ShowComputePoolRequest)
	_ optionsProvider[DescribeComputePoolOptions] = new(DescribeComputePoolRequest)
)

type CreateComputePoolRequest struct {
	IfNotExists        *bool
	name               AccountObjectIdentifier // required
	ForApplication     *AccountObjectIdentifier
	MinNodes           int                       // required
	MaxNodes           int                       // required
	InstanceFamily     ComputePoolInstanceFamily // required
	AutoResume         *bool
	InitiallySuspended *bool
	AutoSuspendSecs    *int
	Tag                []TagAssociation
	Comment            *string
}

type AlterComputePoolRequest struct {
	IfExists  *bool
	name      AccountObjectIdentifier // required
	Set       *ComputePoolSetRequest
	Unset     *ComputePoolUnsetRequest
	SetTags   []TagAssociation
	UnsetTags []ObjectIdentifier
}

type ComputePoolSetRequest struct {
	MinNodes        *int
	MaxNodes        *int
	AutoResume      *bool
	AutoSuspendSecs *int
	Comment         *string
}

type ComputePoolUnsetRequest struct {
	AutoResume      *bool
	AutoSuspendSecs *bool
	Comment         *bool
}

type DropComputePoolRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type ShowComputePoolRequest struct {
	Like       *Like
	StartsWith *string
	Limit      *LimitFrom
}

type DescribeComputePoolRequest struct {
	name AccountObjectIdentifier // required
}
