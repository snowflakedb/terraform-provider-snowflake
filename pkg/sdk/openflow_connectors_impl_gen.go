package sdk

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ OpenflowConnectors = (*openflowConnectors)(nil)

var (
	_ convertibleRow[OpenflowConnector]        = new(openflowConnectorRow)
	_ convertibleRow[OpenflowConnectorDetails] = new(openflowConnectorDetailsRow)
)

type openflowConnectors struct {
	client *Client
}

func (v *openflowConnectors) Create(ctx context.Context, request *CreateOpenflowConnectorRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *openflowConnectors) Alter(ctx context.Context, request *AlterOpenflowConnectorRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *openflowConnectors) Drop(ctx context.Context, request *DropOpenflowConnectorRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *openflowConnectors) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error {
		return v.Drop(ctx, NewDropOpenflowConnectorRequest(id).WithIfExists(true))
	}, ctx, id)
}

func (v *openflowConnectors) Show(ctx context.Context, request *ShowOpenflowConnectorRequest) ([]OpenflowConnector, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[openflowConnectorRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[openflowConnectorRow, OpenflowConnector](dbRows)
}

func (v *openflowConnectors) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*OpenflowConnector, error) {
	request := NewShowOpenflowConnectorRequest().WithLike(Like{Pattern: String(id.Name())})
	connectors, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(connectors, func(r OpenflowConnector) bool {
		return r.Name == id.Name() && r.DatabaseName == id.DatabaseName() && r.SchemaName == id.SchemaName()
	})
}

func (v *openflowConnectors) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*OpenflowConnector, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *openflowConnectors) Describe(ctx context.Context, id SchemaObjectIdentifier) (*OpenflowConnectorDetails, error) {
	opts := &DescribeOpenflowConnectorOptions{name: id}
	result, err := validateAndQueryOne[openflowConnectorDetailsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return conversionErrorWrapped(result.convert())
}

func (r *CreateOpenflowConnectorRequest) toOpts() *CreateOpenflowConnectorOptions {
	return &CreateOpenflowConnectorOptions{
		IfNotExists:    r.IfNotExists,
		name:           r.name,
		InRuntime:      r.InRuntime,
		FromDefinition: r.FromDefinition,
		FromStage:      r.FromStage,
		DisplayName:    r.DisplayName,
		Comment:        r.Comment,
	}
}

func (r *AlterOpenflowConnectorRequest) toOpts() *AlterOpenflowConnectorOptions {
	opts := &AlterOpenflowConnectorOptions{
		name:  r.name,
		Start: r.Start,
		Stop:  r.Stop,
	}
	if r.Set != nil {
		opts.Set = &OpenflowConnectorSet{
			DisplayName: r.Set.DisplayName,
			Comment:     r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &OpenflowConnectorUnset{
			DisplayName: r.Unset.DisplayName,
			Comment:     r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropOpenflowConnectorRequest) toOpts() *DropOpenflowConnectorOptions {
	return &DropOpenflowConnectorOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
}

func (r *ShowOpenflowConnectorRequest) toOpts() *ShowOpenflowConnectorOptions {
	return &ShowOpenflowConnectorOptions{Like: r.Like}
}

func (r *DescribeOpenflowConnectorRequest) toOpts() *DescribeOpenflowConnectorOptions {
	return &DescribeOpenflowConnectorOptions{name: r.name}
}

func (r openflowConnectorRow) convert() (*OpenflowConnector, error) {
	status, err := ToOpenflowConnectorStatus(r.Status)
	if err != nil {
		return nil, fmt.Errorf("error converting openflow connector status: %w", err)
	}
	c := &OpenflowConnector{
		Name:         r.Name,
		DatabaseName: r.DatabaseName,
		SchemaName:   r.SchemaName,
		Runtime:      r.Runtime,
		Owner:        r.Owner,
		CreatedOn:    r.CreatedOn,
		UpdatedOn:    r.UpdatedOn,
		Status:       status,
		// Started is derived from status: connector is started when RUNNING or STARTING
		Started: status == OpenflowConnectorStatusRunning || status == OpenflowConnectorStatusStarting,
	}
	if r.ConnectorDefinition.Valid {
		c.ConnectorDefinition = &r.ConnectorDefinition.String
	}
	if r.DisplayName.Valid {
		c.DisplayName = &r.DisplayName.String
	}
	if r.Comment.Valid {
		c.Comment = &r.Comment.String
	}
	return c, nil
}

func (r openflowConnectorDetailsRow) convert() (*OpenflowConnectorDetails, error) {
	status, err := ToOpenflowConnectorStatus(r.Status)
	if err != nil {
		return nil, fmt.Errorf("error converting openflow connector status: %w", err)
	}
	c := &OpenflowConnectorDetails{
		Name:         r.Name,
		DatabaseName: r.DatabaseName,
		SchemaName:   r.SchemaName,
		Runtime:      r.Runtime,
		Owner:        r.Owner,
		CreatedOn:    r.CreatedOn,
		UpdatedOn:    r.UpdatedOn,
		Status:       status,
		// Started is derived from status: connector is started when RUNNING or STARTING
		Started: status == OpenflowConnectorStatusRunning || status == OpenflowConnectorStatusStarting,
	}
	if r.ConnectorDefinition.Valid {
		c.ConnectorDefinition = &r.ConnectorDefinition.String
	}
	if r.DisplayName.Valid {
		c.DisplayName = &r.DisplayName.String
	}
	if r.Comment.Valid {
		c.Comment = &r.Comment.String
	}
	if r.ErrorCode.Valid {
		c.ErrorCode = &r.ErrorCode.String
	}
	if r.StatusMessage.Valid {
		c.StatusMessage = &r.StatusMessage.String
	}
	return c, nil
}
