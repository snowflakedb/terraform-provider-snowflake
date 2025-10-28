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

func (t *TagAssert) HasAllowedValuesUnordered(expected ...string) *TagAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Tag) error {
		t.Helper()
		if len(o.AllowedValues) != len(expected) {
			return fmt.Errorf("expected allowed values length: %v; got: %v", len(expected), len(o.AllowedValues))
		}
		var errs []error
		for _, wantElem := range expected {
			if !slices.ContainsFunc(o.AllowedValues, func(gotElem string) bool {
				return wantElem == gotElem
			}) {
				errs = append(errs, fmt.Errorf("expected value: %s, to be in the value list: %v", wantElem, o.AllowedValues))
			}
		}
		return errors.Join(errs...)
	})
	return t
}
