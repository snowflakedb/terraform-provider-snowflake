package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageLocationAzureDetailsAssert) HasUsePrivatelinkEndpointEmpty() *StorageLocationAzureDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageLocationAzureDetails) error {
		t.Helper()
		if o.UsePrivatelinkEndpoint != nil {
			return fmt.Errorf("expected use privatelink endpoint to be nil; got: %v", *o.UsePrivatelinkEndpoint)
		}
		return nil
	})
	return s
}
