package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (s *SecondaryConnectionResourceAssert) HasAsReplicaOfIdentifier(expected sdk.ExternalObjectIdentifier) *SecondaryConnectionResourceAssert {
	s.AddAssertion(assert.ValueSet("as_replica_of", expected.Name()))
	return s
}
