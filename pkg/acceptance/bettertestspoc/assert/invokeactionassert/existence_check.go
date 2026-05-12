package invokeactionassert

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type showByIDFunc[T any, ID sdk.ObjectIdentifier] func(t *testing.T, id ID) (T, error)

type nonExistenceCheck[T any, ID sdk.ObjectIdentifier] struct {
	objectType  sdk.ObjectType
	id          ID
	ShowCommand func(testClient *helpers.TestClient) showByIDFunc[T, ID]
}

func (w *nonExistenceCheck[T, ID]) ToTerraformTestCheckFunc(t *testing.T, testClient *helpers.TestClient) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		if testClient == nil {
			return errors.New("testClient must not be nil")
		}
		_, err := w.ShowCommand(testClient)(t, w.id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				return nil
			}
			return err
		}
		return fmt.Errorf("expected %s %s to be missing, but it exists", w.objectType, w.id.FullyQualifiedName())
	}
}

func newNonExistenceCheck[T any, ID sdk.ObjectIdentifier](
	objectType sdk.ObjectType,
	id ID,
	showCommand func(testClient *helpers.TestClient) showByIDFunc[T, ID],
) *nonExistenceCheck[T, ID] {
	return &nonExistenceCheck[T, ID]{
		objectType:  objectType,
		id:          id,
		ShowCommand: showCommand,
	}
}
