package objectassert

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TODO [SNOW-1501905]: generalize this type of assertion
type tagNonExistenceCheck struct {
	id sdk.SchemaObjectIdentifier
}

func (w *tagNonExistenceCheck) ToTerraformTestCheckFunc(t *testing.T, testClient *helpers.TestClient) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		if testClient == nil {
			return errors.New("testClient must not be nil")
		}
		_, err := testClient.Streamlit.Show(t, w.id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				return nil
			}
			return err
		}
		return fmt.Errorf("expected tag %s to be missing, but it exists", w.id.FullyQualifiedName())
	}
}

func TagDoesNotExist(t *testing.T, id sdk.SchemaObjectIdentifier) assert.TestCheckFuncProvider {
	t.Helper()
	return &tagNonExistenceCheck{id: id}
}

// HasAllowedValuesSet checks that the allowed_values field contains the expected values (order independent)
func (t *TagAssert) HasAllowedValuesSet(expected ...string) *TagAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Tag) error {
		t.Helper()
		if len(o.AllowedValues) != len(expected) {
			return fmt.Errorf("expected allowed values length: %d; got: %d", len(expected), len(o.AllowedValues))
		}

		// Sort both slices for comparison
		actualSorted := make([]string, len(o.AllowedValues))
		copy(actualSorted, o.AllowedValues)
		slices.Sort(actualSorted)

		expectedSorted := make([]string, len(expected))
		copy(expectedSorted, expected)
		slices.Sort(expectedSorted)

		if !slices.Equal(actualSorted, expectedSorted) {
			return fmt.Errorf("expected allowed values: %v; got: %v", expected, o.AllowedValues)
		}
		return nil
	})
	return t
}
