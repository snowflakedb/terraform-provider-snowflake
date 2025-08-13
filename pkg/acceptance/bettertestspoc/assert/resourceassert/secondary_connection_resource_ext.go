package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *SecondaryConnectionResourceAssert) HasAsReplicaOfIdentifier(expected sdk.ExternalObjectIdentifier) *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueSet("as_replica_of", expected.FullyQualifiedName()))
	return s
}
