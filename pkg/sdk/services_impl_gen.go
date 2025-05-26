package sdk

import (
	"context"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ Services = (*services)(nil)

type services struct {
	client *Client
}

func (v *services) Create(ctx context.Context, request *CreateServiceRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *services) Alter(ctx context.Context, request *AlterServiceRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *services) Drop(ctx context.Context, request *DropServiceRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *services) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropServiceRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *services) Show(ctx context.Context, request *ShowServiceRequest) ([]Service, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[servicesRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[servicesRow, Service](dbRows)
	return resultList, nil
}

func (v *services) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Service, error) {
	request := NewShowServiceRequest().
		WithLike(Like{Pattern: String(id.Name())})
	services, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(services, func(r Service) bool { return r.Name == id.Name() })
}

func (v *services) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*Service, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *services) Describe(ctx context.Context, id SchemaObjectIdentifier) (*ServiceDetails, error) {
	opts := &DescribeServiceOptions{
		name: id,
	}
	result, err := validateAndQueryOne[serviceDescRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert(), nil
}

func (r *CreateServiceRequest) toOpts() *CreateServiceOptions {
	opts := &CreateServiceOptions{
		IfNotExists:               r.IfNotExists,
		name:                      r.name,
		InComputePool:             r.InComputePool,
		FromSpecification:         r.FromSpecification,
		FromSpecificationTemplate: r.FromSpecificationTemplate,
		AutoSuspendSecs:           r.AutoSuspendSecs,

		AutoResume:        r.AutoResume,
		MinInstances:      r.MinInstances,
		MinReadyInstances: r.MinReadyInstances,
		MaxInstances:      r.MaxInstances,
		QueryWarehouse:    r.QueryWarehouse,
		Tag:               r.Tag,
		Comment:           r.Comment,
	}
	if r.ExternalAccessIntegrations != nil {
		opts.ExternalAccessIntegrations = &ServiceExternalAccessIntegrations{
			ExternalAccessIntegrations: r.ExternalAccessIntegrations.ExternalAccessIntegrations,
		}
	}
	return opts
}

func (r *AlterServiceRequest) toOpts() *AlterServiceOptions {
	opts := &AlterServiceOptions{
		IfExists:                  r.IfExists,
		name:                      r.name,
		Resume:                    r.Resume,
		Suspend:                   r.Suspend,
		FromSpecification:         r.FromSpecification,
		FromSpecificationTemplate: r.FromSpecificationTemplate,

		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	if r.Restore != nil {
		opts.Restore = &Restore{
			Volume:       r.Restore.Volume,
			Instances:    r.Restore.Instances,
			FromSnapshot: r.Restore.FromSnapshot,
		}
	}
	if r.Set != nil {
		opts.Set = &ServiceSet{
			MinInstances:      r.Set.MinInstances,
			MaxInstances:      r.Set.MaxInstances,
			AutoSuspendSecs:   r.Set.AutoSuspendSecs,
			MinReadyInstances: r.Set.MinReadyInstances,
			QueryWarehouse:    r.Set.QueryWarehouse,
			AutoResume:        r.Set.AutoResume,

			Comment: r.Set.Comment,
		}
		if r.Set.ExternalAccessIntegrations != nil {
			opts.Set.ExternalAccessIntegrations = &ServiceExternalAccessIntegrations{
				ExternalAccessIntegrations: r.Set.ExternalAccessIntegrations.ExternalAccessIntegrations,
			}
		}
	}
	if r.Unset != nil {
		opts.Unset = &ServiceUnset{
			MinInstances:               r.Unset.MinInstances,
			AutoSuspendSecs:            r.Unset.AutoSuspendSecs,
			MaxInstances:               r.Unset.MaxInstances,
			MinReadyInstances:          r.Unset.MinReadyInstances,
			QueryWarehouse:             r.Unset.QueryWarehouse,
			AutoResume:                 r.Unset.AutoResume,
			ExternalAccessIntegrations: r.Unset.ExternalAccessIntegrations,
			Comment:                    r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropServiceRequest) toOpts() *DropServiceOptions {
	opts := &DropServiceOptions{
		IfExists: r.IfExists,
		name:     r.name,
		Force:    r.Force,
	}
	return opts
}

func (r *ShowServiceRequest) toOpts() *ShowServiceOptions {
	opts := &ShowServiceOptions{
		Job:         r.Job,
		ExcludeJobs: r.ExcludeJobs,
		Like:        r.Like,
		In:          r.In,
		StartsWith:  r.StartsWith,
		Limit:       r.Limit,
	}
	return opts
}

func (r servicesRow) convert() *Service {
	service := &Service{
		Name:              r.Name,
		DatabaseName:      r.DatabaseName,
		SchemaName:        r.SchemaName,
		Owner:             r.Owner,
		DnsName:           r.DnsName,
		CurrentInstances:  r.CurrentInstances,
		TargetInstances:   r.TargetInstances,
		MinReadyInstances: r.MinReadyInstances,
		MinInstances:      r.MinInstances,
		MaxInstances:      r.MaxInstances,
		AutoResume:        r.AutoResume,
		CreatedOn:         r.CreatedOn,
		UpdatedOn:         r.UpdatedOn,
		AutoSuspendSecs:   r.AutoSuspendSecs,
		OwnerRoleType:     r.OwnerRoleType,
		IsJob:             r.IsJob,
		IsAsyncJob:        r.IsAsyncJob,
		SpecDigest:        r.SpecDigest,
		IsUpgrading:       r.IsUpgrading,
	}
	serviceStatus, err := ToServiceStatus(r.Status)
	if err != nil {
		log.Printf("[DEBUG] error converting service status: %v", err)
	} else {
		service.Status = serviceStatus
	}
	if r.Comment.Valid {
		service.Comment = &r.Comment.String
	}
	if r.ManagingObjectDomain.Valid {
		service.ManagingObjectDomain = &r.ManagingObjectDomain.String
	}
	if r.ManagingObjectName.Valid {
		service.ManagingObjectName = &r.ManagingObjectName.String
	}
	computePoolId, err := ParseAccountObjectIdentifier(r.ComputePool)
	if err != nil {
		log.Printf("[DEBUG] failed to parse compute pool in service: %v", err)
	} else {
		service.ComputePool = computePoolId
	}
	if r.QueryWarehouse.Valid {
		id, err := ParseAccountObjectIdentifier(r.QueryWarehouse.String)
		if err != nil {
			log.Printf("[DEBUG] failed to parse query warehouse in service: %v", err)
		} else {
			service.QueryWarehouse = &id
		}
	}
	if r.SuspendedOn.Valid {
		service.SuspendedOn = &r.SuspendedOn.Time
	}
	if r.ResumedOn.Valid {
		service.ResumedOn = &r.ResumedOn.Time
	}
	if r.ExternalAccessIntegrations.Valid {
		eaiIds, err := ParseCommaSeparatedAccountObjectIdentifierArray(r.ExternalAccessIntegrations.String)
		if err != nil {
			log.Printf("[DEBUG] failed to parse external access integrations in service: %v", err)
		} else {
			service.ExternalAccessIntegrations = eaiIds
		}
	}
	return service
}

func (r *DescribeServiceRequest) toOpts() *DescribeServiceOptions {
	opts := &DescribeServiceOptions{
		name: r.name,
	}
	return opts
}

func (r serviceDescRow) convert() *ServiceDetails {
	service := &ServiceDetails{
		Name:              r.Name,
		DatabaseName:      r.DatabaseName,
		SchemaName:        r.SchemaName,
		Owner:             r.Owner,
		Spec:              r.Spec,
		DnsName:           r.DnsName,
		CurrentInstances:  r.CurrentInstances,
		TargetInstances:   r.TargetInstances,
		MinReadyInstances: r.MinReadyInstances,
		MinInstances:      r.MinInstances,
		MaxInstances:      r.MaxInstances,
		AutoResume:        r.AutoResume,
		CreatedOn:         r.CreatedOn,
		UpdatedOn:         r.UpdatedOn,
		AutoSuspendSecs:   r.AutoSuspendSecs,
		OwnerRoleType:     r.OwnerRoleType,
		IsJob:             r.IsJob,
		IsAsyncJob:        r.IsAsyncJob,
		SpecDigest:        r.SpecDigest,
		IsUpgrading:       r.IsUpgrading,
	}
	serviceStatus, err := ToServiceStatus(r.Status)
	if err != nil {
		log.Printf("[DEBUG] error converting service status: %v", err)
	} else {
		service.Status = serviceStatus
	}
	if r.Comment.Valid {
		service.Comment = &r.Comment.String
	}
	if r.ManagingObjectDomain.Valid {
		service.ManagingObjectDomain = &r.ManagingObjectDomain.String
	}
	if r.ManagingObjectName.Valid {
		service.ManagingObjectName = &r.ManagingObjectName.String
	}
	computePoolId, err := ParseAccountObjectIdentifier(r.ComputePool)
	if err != nil {
		log.Printf("[DEBUG] failed to parse compute pool in service: %v", err)
	} else {
		service.ComputePool = computePoolId
	}
	if r.QueryWarehouse.Valid {
		id, err := ParseAccountObjectIdentifier(r.QueryWarehouse.String)
		if err != nil {
			log.Printf("[DEBUG] failed to parse query warehouse in service: %v", err)
		} else {
			service.QueryWarehouse = &id
		}
	}
	if r.SuspendedOn.Valid {
		service.SuspendedOn = &r.SuspendedOn.Time
	}
	if r.ResumedOn.Valid {
		service.ResumedOn = &r.ResumedOn.Time
	}
	if r.ExternalAccessIntegrations.Valid {
		eaiIds, err := ParseCommaSeparatedAccountObjectIdentifierArray(r.ExternalAccessIntegrations.String)
		if err != nil {
			log.Printf("[DEBUG] failed to parse external access integrations in service: %v", err)
		} else {
			service.ExternalAccessIntegrations = eaiIds
		}
	}
	return service
}
