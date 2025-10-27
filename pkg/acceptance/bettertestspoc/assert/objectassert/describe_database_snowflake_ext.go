package objectassert

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// TODO [SNOW-1501905]: this file should be fully regenerated when adding and option to assert the results of describe
type DatabaseDescribeAssert struct {
	*assert.SnowflakeObjectAssert[sdk.DatabaseDetails, sdk.AccountObjectIdentifier]
}

func DatabaseDescribe(t *testing.T, id sdk.AccountObjectIdentifier) *DatabaseDescribeAssert {
	t.Helper()

	return &DatabaseDescribeAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectType("DATABASE_DETAILS"), id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.DatabaseDetails, sdk.AccountObjectIdentifier] {
			return testClient.Database.Describe
		}),
	}
}

func (d *DatabaseDescribeAssert) DoesNotContainPublicSchema() *DatabaseDescribeAssert {
	d.AddAssertion(func(t *testing.T, o *sdk.DatabaseDetails) error {
		t.Helper()
		if slices.ContainsFunc(o.Rows, func(row sdk.DatabaseDetailsRow) bool { return row.Name == "PUBLIC" && row.Kind == "SCHEMA" }) {
			return fmt.Errorf("expected database %s to not contain public schema", d.GetId())
		}
		return nil
	})
	return d
}

func (d *DatabaseDescribeAssert) ContainsPublicSchema() *DatabaseDescribeAssert {
	d.AddAssertion(func(t *testing.T, o *sdk.DatabaseDetails) error {
		t.Helper()
		if !slices.ContainsFunc(o.Rows, func(row sdk.DatabaseDetailsRow) bool { return row.Name == "PUBLIC" && row.Kind == "SCHEMA" }) {
			return fmt.Errorf("expected database %s to contain public schema", d.GetId())
		}
		return nil
	})
	return d
}

func DatabaseIsMissing(t *testing.T, id sdk.AccountObjectIdentifier) *DatabaseAssert {
	t.Helper()
	return &DatabaseAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectTypeDatabase, id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.Database, sdk.AccountObjectIdentifier] {
			return func(t *testing.T, identifier sdk.AccountObjectIdentifier) (*sdk.Database, error) {
				db, err := testClient.Database.Show(t, id)
				if err != nil {
					if errors.Is(err, sdk.ErrObjectNotFound) {
						return db, nil
					}
					return db, err
				}
				return db, fmt.Errorf("expected database %s to be missing, but it exists", id)
			}
		}),
	}
}
