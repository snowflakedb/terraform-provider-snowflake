package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateServiceOptions]     = new(CreateServiceRequest)
	_ optionsProvider[AlterServiceOptions]      = new(AlterServiceRequest)
	_ optionsProvider[ExecuteJobServiceOptions] = new(ExecuteJobServiceRequest)
	_ optionsProvider[DropServiceOptions]       = new(DropServiceRequest)
	_ optionsProvider[ShowServiceOptions]       = new(ShowServiceRequest)
	_ optionsProvider[DescribeServiceOptions]   = new(DescribeServiceRequest)
)

type CreateServiceRequest struct {
	IfNotExists                *bool
	name                       SchemaObjectIdentifier  // required
	InComputePool              AccountObjectIdentifier // required
	FromSpecification          *ServiceFromSpecificationRequest
	FromSpecificationTemplate  *ServiceFromSpecificationTemplateRequest
	AutoSuspendSecs            *int
	ExternalAccessIntegrations *ServiceExternalAccessIntegrationsRequest
	AutoResume                 *bool
	MinInstances               *int
	MinReadyInstances          *int
	MaxInstances               *int
	QueryWarehouse             *AccountObjectIdentifier
	Tag                        []TagAssociation
	Comment                    *string
}

type ServiceFromSpecificationRequest struct {
	Location          Location
	SpecificationFile *string
	Specification     *string
}

type ServiceFromSpecificationTemplateRequest struct {
	Location                  Location
	SpecificationTemplateFile *string
	SpecificationTemplate     *string
	Using                     []ListItem // required
}

type ServiceExternalAccessIntegrationsRequest struct {
	ExternalAccessIntegrations []AccountObjectIdentifier // required
}

type AlterServiceRequest struct {
	IfExists                  *bool
	name                      SchemaObjectIdentifier // required
	Resume                    *bool
	Suspend                   *bool
	FromSpecification         *ServiceFromSpecificationRequest
	FromSpecificationTemplate *ServiceFromSpecificationTemplateRequest
	Restore                   *RestoreRequest
	Set                       *ServiceSetRequest
	Unset                     *ServiceUnsetRequest
	SetTags                   []TagAssociation
	UnsetTags                 []ObjectIdentifier
}

type RestoreRequest struct {
	Volume       string                 // required
	Instances    []int                  // required
	FromSnapshot SchemaObjectIdentifier // required
}

type ServiceSetRequest struct {
	MinInstances               *int
	MaxInstances               *int
	AutoSuspendSecs            *int
	MinReadyInstances          *int
	QueryWarehouse             *AccountObjectIdentifier
	AutoResume                 *bool
	ExternalAccessIntegrations *ServiceExternalAccessIntegrationsRequest
	Comment                    *string
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

type ExecuteJobServiceRequest struct {
	InComputePool                       AccountObjectIdentifier // required
	Name                                SchemaObjectIdentifier  // required
	Async                               *bool
	QueryWarehouse                      *AccountObjectIdentifier
	Comment                             *string
	ExternalAccessIntegrations          *ServiceExternalAccessIntegrationsRequest
	JobServiceFromSpecification         *JobServiceFromSpecificationRequest
	JobServiceFromSpecificationTemplate *JobServiceFromSpecificationTemplateRequest
	Tag                                 []TagAssociation
}

type JobServiceFromSpecificationRequest struct {
	Location          Location
	SpecificationFile *string
	Specification     *string
}

type JobServiceFromSpecificationTemplateRequest struct {
	Location                  Location
	SpecificationTemplateFile *string
	SpecificationTemplate     *string
	Using                     []ListItem // required
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
