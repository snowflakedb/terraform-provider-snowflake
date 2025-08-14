package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ Pipes = (*pipes)(nil)

type pipes struct {
	client *Client
}

func (v *pipes) Create(ctx context.Context, id SchemaObjectIdentifier, copyStatement string, opts *CreatePipeOptions) error {
	opts = createIfNil[CreatePipeOptions](opts)
	opts.name = id
	opts.copyStatement = copyStatement
	return validateAndExec(v.client, ctx, opts)
}

func (v *pipes) Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterPipeOptions) error {
	opts = createIfNil[AlterPipeOptions](opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *pipes) Drop(ctx context.Context, id SchemaObjectIdentifier, opts *DropPipeOptions) error {
	opts = createIfNil[DropPipeOptions](opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *pipes) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, id, nil) }, ctx, id)
}

func (v *pipes) Show(ctx context.Context, opts *ShowPipeOptions) ([]Pipe, error) {
	dbRows, err := validateAndQuery[pipeDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}

	return convertRowsErr[pipeDBRow, Pipe](dbRows)
}

func (v *pipes) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Pipe, error) {
	pipes, err := v.Show(ctx, &ShowPipeOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
		In: &In{
			Schema: id.SchemaId(),
		},
	})
	if err != nil {
		return nil, err
	}

	return collections.FindFirst(pipes, func(p Pipe) bool { return p.ID().name == id.Name() })
}

func (v *pipes) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*Pipe, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *pipes) Describe(ctx context.Context, id SchemaObjectIdentifier) (*Pipe, error) {
	opts := &describePipeOptions{
		name: id,
	}
	pipeRow, err := validateAndQueryOne[pipeDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return pipeRow.convert(), nil
}
