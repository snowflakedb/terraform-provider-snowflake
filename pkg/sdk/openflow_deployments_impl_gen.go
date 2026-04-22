package sdk

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ OpenflowDeployments = (*openflowDeployments)(nil)

var (
	_ convertibleRow[OpenflowDeployment]        = new(openflowDeploymentRow)
	_ convertibleRow[OpenflowDeploymentDetails] = new(openflowDeploymentDetailsRow)
)

type openflowDeployments struct {
	client *Client
}

func (v *openflowDeployments) Create(ctx context.Context, request *CreateOpenflowDeploymentRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *openflowDeployments) Alter(ctx context.Context, request *AlterOpenflowDeploymentRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *openflowDeployments) Drop(ctx context.Context, request *DropOpenflowDeploymentRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *openflowDeployments) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error {
		return v.Drop(ctx, NewDropOpenflowDeploymentRequest(id).WithIfExists(true))
	}, ctx, id)
}

func (v *openflowDeployments) Show(ctx context.Context, request *ShowOpenflowDeploymentRequest) ([]OpenflowDeployment, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[openflowDeploymentRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[openflowDeploymentRow, OpenflowDeployment](dbRows)
}

func (v *openflowDeployments) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*OpenflowDeployment, error) {
	request := NewShowOpenflowDeploymentRequest().WithLike(Like{Pattern: String(id.Name())})
	deployments, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(deployments, func(r OpenflowDeployment) bool { return r.Name == id.Name() })
}

func (v *openflowDeployments) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*OpenflowDeployment, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *openflowDeployments) Describe(ctx context.Context, id AccountObjectIdentifier) (*OpenflowDeploymentDetails, error) {
	opts := &DescribeOpenflowDeploymentOptions{name: id}
	result, err := validateAndQueryOne[openflowDeploymentDetailsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return conversionErrorWrapped(result.convert())
}

func (r *CreateOpenflowDeploymentRequest) toOpts() *CreateOpenflowDeploymentOptions {
	return &CreateOpenflowDeploymentOptions{
		IfNotExists:               r.IfNotExists,
		name:                      r.name,
		DeploymentType:            r.DeploymentType,
		VpcType:                   r.VpcType,
		CustomIngressHostname:     r.CustomIngressHostname,
		UsePrivateLink:            r.UsePrivateLink,
		UseUserAuthOverPrivatelink: r.UseUserAuthOverPrivatelink,
		EventTable:                r.EventTable,
		DisplayName:               r.DisplayName,
		Comment:                   r.Comment,
	}
}

func (r *AlterOpenflowDeploymentRequest) toOpts() *AlterOpenflowDeploymentOptions {
	opts := &AlterOpenflowDeploymentOptions{name: r.name}
	if r.Set != nil {
		opts.Set = &OpenflowDeploymentSet{
			Comment:     r.Set.Comment,
			DisplayName: r.Set.DisplayName,
			EventTable:  r.Set.EventTable,
		}
	}
	if r.Unset != nil {
		opts.Unset = &OpenflowDeploymentUnset{
			Comment:     r.Unset.Comment,
			DisplayName: r.Unset.DisplayName,
			EventTable:  r.Unset.EventTable,
		}
	}
	return opts
}

func (r *DropOpenflowDeploymentRequest) toOpts() *DropOpenflowDeploymentOptions {
	return &DropOpenflowDeploymentOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
}

func (r *ShowOpenflowDeploymentRequest) toOpts() *ShowOpenflowDeploymentOptions {
	return &ShowOpenflowDeploymentOptions{Like: r.Like}
}

func (r *DescribeOpenflowDeploymentRequest) toOpts() *DescribeOpenflowDeploymentOptions {
	return &DescribeOpenflowDeploymentOptions{name: r.name}
}

func (r openflowDeploymentRow) convert() (*OpenflowDeployment, error) {
	d := &OpenflowDeployment{
		Name:                       r.Name,
		UsePrivateLink:             r.UsePrivateLink,
		UseUserAuthOverPrivatelink: r.UseUserAuthOverPrivatelink,
		Owner:                      r.Owner,
		CreatedOn:                  r.CreatedOn,
		UpdatedOn:                  r.UpdatedOn,
	}
	if r.VpcType.Valid {
		vpcType, err := ToOpenflowVpcType(r.VpcType.String)
		if err != nil {
			return nil, fmt.Errorf("error converting openflow deployment vpc type: %w", err)
		}
		d.VpcType = &vpcType
	}
	if r.CustomIngressHostname.Valid {
		d.CustomIngressHostname = &r.CustomIngressHostname.String
	}
	if r.Key.Valid {
		d.Key = &r.Key.String
	}
	if r.DisplayName.Valid {
		d.DisplayName = &r.DisplayName.String
	}
	if r.Comment.Valid {
		d.Comment = &r.Comment.String
	}
	// EventTable is not returned by SHOW OPENFLOW DEPLOYMENTS; it stays nil
	deploymentType, err := ToOpenflowDeploymentType(r.DeploymentType)
	if err != nil {
		return nil, fmt.Errorf("error converting openflow deployment type: %w", err)
	}
	d.DeploymentType = deploymentType
	status, err := ToOpenflowDeploymentStatus(r.Status)
	if err != nil {
		return nil, fmt.Errorf("error converting openflow deployment status: %w", err)
	}
	d.Status = status
	return d, nil
}

func (r openflowDeploymentDetailsRow) convert() (*OpenflowDeploymentDetails, error) {
	d := &OpenflowDeploymentDetails{
		Name:                       r.Name,
		UsePrivateLink:             r.UsePrivateLink,
		UseUserAuthOverPrivatelink: r.UseUserAuthOverPrivatelink,
		Owner:                      r.Owner,
		CreatedOn:                  r.CreatedOn,
		UpdatedOn:                  r.UpdatedOn,
	}
	if r.VpcType.Valid {
		vpcType, err := ToOpenflowVpcType(r.VpcType.String)
		if err != nil {
			return nil, fmt.Errorf("error converting openflow deployment vpc type: %w", err)
		}
		d.VpcType = &vpcType
	}
	if r.CustomIngressHostname.Valid {
		d.CustomIngressHostname = &r.CustomIngressHostname.String
	}
	if r.Key.Valid {
		d.Key = &r.Key.String
	}
	if r.DisplayName.Valid {
		d.DisplayName = &r.DisplayName.String
	}
	if r.Comment.Valid {
		d.Comment = &r.Comment.String
	}
	// EventTable is not returned by SHOW/DESCRIBE OPENFLOW DEPLOYMENT; it stays nil
	if r.ErrorCode.Valid {
		d.ErrorCode = &r.ErrorCode.String
	}
	if r.StatusMessage.Valid {
		d.StatusMessage = &r.StatusMessage.String
	}
	deploymentType, err := ToOpenflowDeploymentType(r.DeploymentType)
	if err != nil {
		return nil, fmt.Errorf("error converting openflow deployment type: %w", err)
	}
	d.DeploymentType = deploymentType
	status, err := ToOpenflowDeploymentStatus(r.Status)
	if err != nil {
		return nil, fmt.Errorf("error converting openflow deployment status: %w", err)
	}
	d.Status = status
	return d, nil
}
