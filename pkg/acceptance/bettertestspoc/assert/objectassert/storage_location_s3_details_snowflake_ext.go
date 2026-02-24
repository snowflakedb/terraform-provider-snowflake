package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageLocationS3DetailsAssert) HasStorageAwsIamUserArnNotEmpty() *StorageLocationS3DetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageLocationS3Details) error {
		t.Helper()
		if o.StorageAwsIamUserArn == "" {
			return fmt.Errorf("expected storage aws iam user arn not empty; got empty")
		}
		return nil
	})
	return s
}
