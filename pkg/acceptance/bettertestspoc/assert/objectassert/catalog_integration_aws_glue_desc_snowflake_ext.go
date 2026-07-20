package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *CatalogIntegrationAwsGlueDetailsAssert) HasGlueAwsIamUserArnNotEmpty() *CatalogIntegrationAwsGlueDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CatalogIntegrationAwsGlueDetails) error {
		t.Helper()
		if o.GlueAwsIamUserArn == "" {
			return fmt.Errorf("expected glue aws iam user arn not empty; got empty")
		}
		return nil
	})
	return c
}

func (c *CatalogIntegrationAwsGlueDetailsAssert) HasGlueAwsExternalIdNotEmpty() *CatalogIntegrationAwsGlueDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CatalogIntegrationAwsGlueDetails) error {
		t.Helper()
		if o.GlueAwsExternalId == "" {
			return fmt.Errorf("expected glue aws external id not empty; got empty")
		}
		return nil
	})
	return c
}
