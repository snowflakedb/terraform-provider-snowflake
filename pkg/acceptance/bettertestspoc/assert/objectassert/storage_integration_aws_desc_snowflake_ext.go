package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageIntegrationAwsDetailsAssert) HasIamUserArnNotEmpty() *StorageIntegrationAwsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAwsDetails) error {
		t.Helper()
		if o.IamUserArn == "" {
			return fmt.Errorf("expected iam user arn not empty; got empty")
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAwsDetailsAssert) HasExternalIdNotEmpty() *StorageIntegrationAwsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAwsDetails) error {
		t.Helper()
		if o.ExternalId == "" {
			return fmt.Errorf("expected external id not empty; got empty")
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAwsDetailsAssert) HasExternalIdNotEqualTo(externalId string) *StorageIntegrationAwsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAwsDetails) error {
		t.Helper()
		if o.ExternalId == externalId {
			return fmt.Errorf("expected external id to differ from %q; got the same", externalId)
		}
		return nil
	})
	return s
}
