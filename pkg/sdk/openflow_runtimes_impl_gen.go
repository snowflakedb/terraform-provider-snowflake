package sdk

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ OpenflowRuntimes = (*openflowRuntimes)(nil)

var (
	_ convertibleRow[OpenflowRuntime]        = new(openflowRuntimeRow)
	_ convertibleRow[OpenflowRuntimeDetails] = new(openflowRuntimeDetailsRow)
)

type openflowRuntimes struct {
	client *Client
}

func (v *openflowRuntimes) Create(ctx context.Context, request *CreateOpenflowRuntimeRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *openflowRuntimes) Alter(ctx context.Context, request *AlterOpenflowRuntimeRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *openflowRuntimes) Drop(ctx context.Context, request *DropOpenflowRuntimeRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *openflowRuntimes) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error {
		return v.Drop(ctx, NewDropOpenflowRuntimeRequest(id).WithIfExists(true))
	}, ctx, id)
}

func (v *openflowRuntimes) Show(ctx context.Context, request *ShowOpenflowRuntimeRequest) ([]OpenflowRuntime, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[openflowRuntimeRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[openflowRuntimeRow, OpenflowRuntime](dbRows)
}

func (v *openflowRuntimes) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*OpenflowRuntime, error) {
	request := NewShowOpenflowRuntimeRequest().WithLike(Like{Pattern: String(id.Name())})
	runtimes, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(runtimes, func(r OpenflowRuntime) bool {
		return r.Name == id.Name() && r.DatabaseName == id.DatabaseName() && r.SchemaName == id.SchemaName()
	})
}

func (v *openflowRuntimes) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*OpenflowRuntime, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *openflowRuntimes) Describe(ctx context.Context, id SchemaObjectIdentifier) (*OpenflowRuntimeDetails, error) {
	opts := &DescribeOpenflowRuntimeOptions{name: id}
	result, err := validateAndQueryOne[openflowRuntimeDetailsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return conversionErrorWrapped(result.convert())
}

func (r *CreateOpenflowRuntimeRequest) toOpts() *CreateOpenflowRuntimeOptions {
	return &CreateOpenflowRuntimeOptions{
		IfNotExists:                r.IfNotExists,
		name:                       r.name,
		InDeployment:               r.InDeployment,
		ExecuteAsRole:              r.ExecuteAsRole,
		NodeType:                   r.NodeType,
		MinNodes:                   r.MinNodes,
		MaxNodes:                   r.MaxNodes,
		ExternalAccessIntegrations: r.ExternalAccessIntegrations,
		DisplayName:                r.DisplayName,
		Comment:                    r.Comment,
	}
}

func (r *AlterOpenflowRuntimeRequest) toOpts() *AlterOpenflowRuntimeOptions {
	opts := &AlterOpenflowRuntimeOptions{
		name:             r.name,
		Suspend:          r.Suspend,
		Resume:           r.Resume,
		Terminate:        r.Terminate,
		TerminateCascade: r.TerminateCascade,
	}
	if r.Set != nil {
		opts.Set = &OpenflowRuntimeSet{
			MinNodes:                   r.Set.MinNodes,
			MaxNodes:                   r.Set.MaxNodes,
			ExecuteAsRole:              r.Set.ExecuteAsRole,
			ExternalAccessIntegrations: r.Set.ExternalAccessIntegrations,
			DisplayName:                r.Set.DisplayName,
			Comment:                    r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &OpenflowRuntimeUnset{
			DisplayName: r.Unset.DisplayName,
			Comment:     r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropOpenflowRuntimeRequest) toOpts() *DropOpenflowRuntimeOptions {
	return &DropOpenflowRuntimeOptions{
		IfExists: r.IfExists,
		name:     r.name,
		Cascade:  r.Cascade,
	}
}

func (r *ShowOpenflowRuntimeRequest) toOpts() *ShowOpenflowRuntimeOptions {
	return &ShowOpenflowRuntimeOptions{Like: r.Like}
}

func (r *DescribeOpenflowRuntimeRequest) toOpts() *DescribeOpenflowRuntimeOptions {
	return &DescribeOpenflowRuntimeOptions{name: r.name}
}

func parseExternalAccessIntegrations(s sql.NullString) []string {
	if !s.Valid || s.String == "" {
		return nil
	}
	// Snowflake returns a JSON array e.g. ["EAI1","EAI2"]; try that first.
	var parsed []string
	if err := json.Unmarshal([]byte(s.String), &parsed); err == nil {
		return parsed
	}
	// Fall back to comma-separated.
	parts := strings.Split(s.String, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func (r openflowRuntimeRow) convert() (*OpenflowRuntime, error) {
	rt := &OpenflowRuntime{
		Name:         r.Name,
		DatabaseName: r.DatabaseName,
		SchemaName:   r.SchemaName,
		Deployment:   r.Deployment,
		MinNodes:     r.MinNodes,
		MaxNodes:     r.MaxNodes,
		ExecuteAsRole: r.ExecuteAsRole,
		Owner:        r.Owner,
		CreatedOn:    r.CreatedOn,
	}
	rt.ExternalAccessIntegrations = parseExternalAccessIntegrations(r.ExternalAccessIntegrations)
	if r.DisplayName.Valid {
		rt.DisplayName = &r.DisplayName.String
	}
	if r.Comment.Valid {
		rt.Comment = &r.Comment.String
	}
	nodeType, err := ToOpenflowRuntimeNodeType(r.NodeType)
	if err != nil {
		return nil, fmt.Errorf("error converting openflow runtime node type: %w", err)
	}
	rt.NodeType = nodeType
	status, err := ToOpenflowRuntimeStatus(r.Status)
	if err != nil {
		return nil, fmt.Errorf("error converting openflow runtime status: %w", err)
	}
	rt.Status = status
	return rt, nil
}

func (r openflowRuntimeDetailsRow) convert() (*OpenflowRuntimeDetails, error) {
	rt := &OpenflowRuntimeDetails{
		Name:         r.Name,
		DatabaseName: r.DatabaseName,
		SchemaName:   r.SchemaName,
		Deployment:   r.Deployment,
		MinNodes:     r.MinNodes,
		MaxNodes:     r.MaxNodes,
		ExecuteAsRole: r.ExecuteAsRole,
		Owner:        r.Owner,
		CreatedOn:    r.CreatedOn,
	}
	rt.ExternalAccessIntegrations = parseExternalAccessIntegrations(r.ExternalAccessIntegrations)
	if r.DisplayName.Valid {
		rt.DisplayName = &r.DisplayName.String
	}
	if r.Comment.Valid {
		rt.Comment = &r.Comment.String
	}
	if r.ErrorCode.Valid {
		rt.ErrorCode = &r.ErrorCode.String
	}
	if r.StatusMessage.Valid {
		rt.StatusMessage = &r.StatusMessage.String
	}
	nodeType, err := ToOpenflowRuntimeNodeType(r.NodeType)
	if err != nil {
		return nil, fmt.Errorf("error converting openflow runtime node type: %w", err)
	}
	rt.NodeType = nodeType
	status, err := ToOpenflowRuntimeStatus(r.Status)
	if err != nil {
		return nil, fmt.Errorf("error converting openflow runtime status: %w", err)
	}
	rt.Status = status
	return rt, nil
}
