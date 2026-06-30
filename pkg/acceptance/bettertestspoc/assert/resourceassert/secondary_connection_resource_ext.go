package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *SecondaryConnectionResourceAssert) HasAsReplicaOfIdentifier(expected sdk.ExternalObjectIdentifier) *SecondaryConnectionResourceAssert {
	s.ValueSet("as_replica_of", expected.FullyQualifiedName())
	return s
}
