package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type PolicyReferencesClient struct {
	context *TestClientContext
}

func NewPolicyReferencesClient(context *TestClientContext) *PolicyReferencesClient {
	return &PolicyReferencesClient{
		context: context,
	}
}

func (c *PolicyReferencesClient) client() sdk.PolicyReferences {
	return c.context.client.PolicyReferences
}

// GetPolicyReferences is based on https://docs.snowflake.com/en/sql-reference/functions/policy_references.
func (c *PolicyReferencesClient) GetPolicyReferences(t *testing.T, objectId sdk.ObjectIdentifier, entity sdk.PolicyEntityDomain) ([]sdk.PolicyReference, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(
		objectId,
		entity,
	))
}

// GetPolicyReference is based on https://docs.snowflake.com/en/sql-reference/functions/policy_references.
func (c *PolicyReferencesClient) GetPolicyReference(t *testing.T, id sdk.ObjectIdentifier, entity sdk.PolicyEntityDomain) (*sdk.PolicyReference, error) {
	t.Helper()

	references, err := c.GetPolicyReferences(t, id, entity)
	require.NoError(t, err)
	require.Len(t, references, 1)

	return sdk.Pointer(references[0]), nil
}
