package sdk

import (
	"context"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ ComputePools = (*computePools)(nil)

type computePools struct {
	client *Client
}

func (v *computePools) Create(ctx context.Context, request *CreateComputePoolRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *computePools) Alter(ctx context.Context, request *AlterComputePoolRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *computePools) Drop(ctx context.Context, request *DropComputePoolRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *computePools) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropComputePoolRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *computePools) Show(ctx context.Context, request *ShowComputePoolRequest) ([]ComputePool, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[computePoolsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[computePoolsRow, ComputePool](dbRows)
	return resultList, nil
}

func (v *computePools) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ComputePool, error) {
	request := NewShowComputePoolRequest().
		WithLike(Like{Pattern: String(id.Name())})
	computePools, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(computePools, func(r ComputePool) bool { return r.Name == id.Name() })
}

func (v *computePools) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*ComputePool, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *computePools) Describe(ctx context.Context, id AccountObjectIdentifier) (*ComputePoolDetails, error) {
	opts := &DescribeComputePoolOptions{
		name: id,
	}
	result, err := validateAndQueryOne[computePoolDescRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert(), nil
}

func (r *CreateComputePoolRequest) toOpts() *CreateComputePoolOptions {
	opts := &CreateComputePoolOptions{
		IfNotExists:        r.IfNotExists,
		name:               r.name,
		ForApplication:     r.ForApplication,
		MinNodes:           r.MinNodes,
		MaxNodes:           r.MaxNodes,
		InstanceFamily:     r.InstanceFamily,
		AutoResume:         r.AutoResume,
		InitiallySuspended: r.InitiallySuspended,
		AutoSuspendSecs:    r.AutoSuspendSecs,
		Tag:                r.Tag,
		Comment:            r.Comment,
	}
	return opts
}

func (r *AlterComputePoolRequest) toOpts() *AlterComputePoolOptions {
	opts := &AlterComputePoolOptions{
		IfExists: r.IfExists,
		name:     r.name,
		Resume:   r.Resume,
		Suspend:  r.Suspend,
		StopAll:  r.StopAll,

		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	if r.Set != nil {
		opts.Set = &ComputePoolSet{
			MinNodes:        r.Set.MinNodes,
			MaxNodes:        r.Set.MaxNodes,
			AutoResume:      r.Set.AutoResume,
			AutoSuspendSecs: r.Set.AutoSuspendSecs,
			Comment:         r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &ComputePoolUnset{
			AutoResume:      r.Unset.AutoResume,
			AutoSuspendSecs: r.Unset.AutoSuspendSecs,
			Comment:         r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropComputePoolRequest) toOpts() *DropComputePoolOptions {
	opts := &DropComputePoolOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowComputePoolRequest) toOpts() *ShowComputePoolOptions {
	opts := &ShowComputePoolOptions{
		Like:       r.Like,
		StartsWith: r.StartsWith,
		Limit:      r.Limit,
	}
	return opts
}

func (r computePoolsRow) convert() *ComputePool {
	cp := &ComputePool{
		Name:            r.Name,
		MinNodes:        r.MinNodes,
		MaxNodes:        r.MaxNodes,
		NumServices:     r.NumServices,
		NumJobs:         r.NumJobs,
		AutoSuspendSecs: r.AutoSuspendSecs,
		AutoResume:      r.AutoResume,
		ActiveNodes:     r.ActiveNodes,
		IdleNodes:       r.IdleNodes,
		TargetNodes:     r.TargetNodes,
		CreatedOn:       r.CreatedOn,
		ResumedOn:       r.ResumedOn,
		UpdatedOn:       r.UpdatedOn,
		Owner:           r.Owner,
		IsExclusive:     r.IsExclusive,
	}
	if r.Comment.Valid {
		cp.Comment = &r.Comment.String
	}
	if r.Application.Valid {
		id, err := ParseAccountObjectIdentifier(r.Application.String)
		if err != nil {
			log.Printf("[DEBUG] failed to parse application in compute pool: %v", err)
		} else {
			cp.Application = &id
		}
	}
	instanceFamily, err := ToComputePoolInstanceFamily(r.InstanceFamily)
	if err != nil {
		log.Printf("[DEBUG] error converting compute pool instance family: %v", err)
	} else {
		cp.InstanceFamily = instanceFamily
	}
	state, err := ToComputePoolState(r.State)
	if err != nil {
		log.Printf("[DEBUG] error converting compute pool state: %v", err)
	} else {
		cp.State = state
	}
	return cp
}

func (r *DescribeComputePoolRequest) toOpts() *DescribeComputePoolOptions {
	opts := &DescribeComputePoolOptions{
		name: r.name,
	}
	return opts
}

func (r computePoolDescRow) convert() *ComputePoolDetails {
	cp := &ComputePoolDetails{
		Name:            r.Name,
		MinNodes:        r.MinNodes,
		MaxNodes:        r.MaxNodes,
		NumServices:     r.NumServices,
		NumJobs:         r.NumJobs,
		AutoSuspendSecs: r.AutoSuspendSecs,
		AutoResume:      r.AutoResume,
		ActiveNodes:     r.ActiveNodes,
		IdleNodes:       r.IdleNodes,
		TargetNodes:     r.TargetNodes,
		CreatedOn:       r.CreatedOn,
		ResumedOn:       r.ResumedOn,
		UpdatedOn:       r.UpdatedOn,
		Owner:           r.Owner,
		IsExclusive:     r.IsExclusive,
		ErrorCode:       r.ErrorCode,
		StatusMessage:   r.StatusMessage,
	}
	if r.Comment.Valid {
		cp.Comment = &r.Comment.String
	}
	if r.Application.Valid {
		id, err := ParseAccountObjectIdentifier(r.Application.String)
		if err != nil {
			log.Printf("[DEBUG] failed to parse application in compute pool: %v", err)
		} else {
			cp.Application = &id
		}
	}
	instanceFamily, err := ToComputePoolInstanceFamily(r.InstanceFamily)
	if err != nil {
		log.Printf("[DEBUG] error converting compute pool instance family: %v", err)
	} else {
		cp.InstanceFamily = instanceFamily
	}
	state, err := ToComputePoolState(r.State)
	if err != nil {
		log.Printf("[DEBUG] error converting compute pool state: %v", err)
	} else {
		cp.State = state
	}
	return cp
}
