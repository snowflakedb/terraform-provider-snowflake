package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateServiceOptions]   = new(CreateServiceRequest)
	_ optionsProvider[AlterServiceOptions]    = new(AlterServiceRequest)
	_ optionsProvider[DropServiceOptions]     = new(DropServiceRequest)
	_ optionsProvider[ShowServiceOptions]     = new(ShowServiceRequest)
	_ optionsProvider[DescribeServiceOptions] = new(DescribeServiceRequest)
)

type CreateServiceRequest struct {
	IfNotExists                       *bool
	name                              SchemaObjectIdentifier  // required
	InComputePool                     AccountObjectIdentifier // required
	FromSpecification                 *ServiceFromSpecification
	FromSpecificationTemplate         *ServiceFromSpecificationTemplate
	AutoSuspendSecs                   *int
	ServiceExternalAccessIntegrations *ServiceExternalAccessIntegrationsRequest
	AutoResume                        *bool
	MinInstances                      *int
	MinReadyInstances                 *int
	MaxInstances                      *int
	QueryWarehouse                    *AccountObjectIdentifier
	Tag                               []TagAssociation
	Comment                           *string
}

type ServiceExternalAccessIntegrationsRequest struct {
	ServiceExternalAccessIntegrations []AccountObjectIdentifier // required
}

type AlterServiceRequest struct {
	IfExists                  *bool
	name                      SchemaObjectIdentifier // required
	Resume                    *bool
	Suspend                   *bool
	FromSpecification         *ServiceFromSpecification
	FromSpecificationTemplate *ServiceFromSpecificationTemplate
	Restore                   *RestoreRequest
	Set                       *ServiceSetRequest
	Unset                     *ServiceUnsetRequest
	SetTags                   []TagAssociation
	UnsetTags                 []ObjectIdentifier
}

type RestoreRequest struct {
	Volume       string // required
	Instances    string // required
	FromSnapshot string // required
}

type ServiceSetRequest struct {
	MinInstances                      *int
	MaxInstances                      *int
	AutoSuspendSecs                   *int
	MinReadyInstances                 *int
	QueryWarehouse                    *AccountObjectIdentifier
	AutoResume                        *bool
	ServiceExternalAccessIntegrations *ServiceExternalAccessIntegrationsRequest
	Comment                           *string
}

type ServiceUnsetRequest struct {
	MinInstances               *bool
	AutoSuspendSecs            *bool
	MaxInstances               *bool
	MinReadyInstances          *bool
	QueryWarehouse             *bool
	AutoResume                 *bool
	ExternalAccessIntegrations *bool
	Comment                    *bool
}

type DropServiceRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
	Force    *bool
}

type ShowServiceRequest struct {
	Job         *bool
	ExcludeJobs *bool
	Like        *Like
	In          *ServiceIn
	StartsWith  *string
	Limit       *LimitFrom
}

type DescribeServiceRequest struct {
	name SchemaObjectIdentifier // required
}
