package objectassert

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func DatabaseDetails(t *testing.T, id sdk.AccountObjectIdentifier) *DatabaseDetailsAssert {
	t.Helper()
	return &DatabaseDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectType("DatabaseDetails"), id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.DatabaseDetails, sdk.AccountObjectIdentifier] {
			return testClient.Database.Describe
		}),
	}
}

func (d *DatabaseDetailsAssert) DoesNotContainPublicSchema() *DatabaseDetailsAssert {
	d.AddAssertion(func(t *testing.T, o *sdk.DatabaseDetails) error {
		t.Helper()
		if slices.ContainsFunc(o.Rows, func(row sdk.DatabaseDetailsRow) bool { return row.Name == "PUBLIC" && row.Kind == "SCHEMA" }) {
			return fmt.Errorf("expected database %s to not contain public schema", d.GetId())
		}
		return nil
	})
	return d
}

func (d *DatabaseDetailsAssert) ContainsPublicSchema() *DatabaseDetailsAssert {
	d.AddAssertion(func(t *testing.T, o *sdk.DatabaseDetails) error {
		t.Helper()
		if !slices.ContainsFunc(o.Rows, func(row sdk.DatabaseDetailsRow) bool { return row.Name == "PUBLIC" && row.Kind == "SCHEMA" }) {
			return fmt.Errorf("expected database %s to contain public schema", d.GetId())
		}
		return nil
	})
	return d
}
