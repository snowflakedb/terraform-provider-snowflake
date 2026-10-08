package sdk

import (
	"context"
	"errors"
)

func (r *CreateSnowflakeIntelligenceRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (d *SnowflakeIntelligenceDetails) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(d.Name)
}

func (a *SnowflakeIntelligenceAgent) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(a.DatabaseName, a.SchemaName, a.Name)
}

// DropAgentSafely detaches the agent from the Snowflake Intelligence object.
// DROP AGENT has no IF EXISTS clause for the agent. A missing agent returns
// ErrDoesNotExistOrOperationCannotBePerformed. An existing agent that is not attached returns
// ErrObjectWasNotFoundIn. Both mean the attachment is already gone.
func (v *snowflakeIntelligences) DropAgentSafely(ctx context.Context, id AccountObjectIdentifier, agentId SchemaObjectIdentifier) error {
	err := v.Alter(ctx, NewAlterSnowflakeIntelligenceRequest(id).WithIfExists(true).WithDropAgent(agentId))
	if errors.Is(err, ErrDoesNotExistOrOperationCannotBePerformed) || errors.Is(err, ErrObjectWasNotFoundIn) {
		return nil
	}
	return err
}
