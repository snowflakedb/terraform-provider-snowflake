package objectassert

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TODO [SNOW-1501905]: generalize this type of assertion
type secretNonExistenceCheck struct {
	id sdk.SchemaObjectIdentifier
}

func (w *secretNonExistenceCheck) ToTerraformTestCheckFunc(t *testing.T, testClient *helpers.TestClient) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		if testClient == nil {
			return errors.New("testClient must not be nil")
		}
		_, err := testClient.Secret.Show(t, w.id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				return nil
			}
			return err
		}
		return fmt.Errorf("expected secret %s to be missing, but it exists", w.id.FullyQualifiedName())
	}
}

func SecretDoesNotExist(t *testing.T, id sdk.SchemaObjectIdentifier) assert.TestCheckFuncProvider {
	t.Helper()
	return &secretNonExistenceCheck{id: id}
}

func (s *SecretAssert) HasCreatedOnNotEmpty() *SecretAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Secret) error {
		t.Helper()
		if o.CreatedOn.IsZero() {
			return fmt.Errorf("expected created_on to be not empty; got: %v", o.CreatedOn)
		}
		return nil
	})
	return s
}

func (s *SecretAssert) HasNoComment() *SecretAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Secret) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to be nil; got: %s", *o.Comment)
		}
		return nil
	})
	return s
}
