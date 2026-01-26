package objectassert

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type StageDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.StageDetails, sdk.SchemaObjectIdentifier]
}

func StageDetails(t *testing.T, id sdk.SchemaObjectIdentifier) *StageDetailsAssert {
	t.Helper()
	return &StageDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(
			sdk.ObjectType("STAGE_DETAILS"),
			id,
			func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.StageDetails, sdk.SchemaObjectIdentifier] {
				return testClient.Stage.DescribeDetails
			}),
	}
}

func (s *StageDetailsAssert) HasFileFormatCsv(expected sdk.FileFormatCsv) *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.FileFormatCsv == nil {
			return fmt.Errorf("expected file format to be CSV; got: nil")
		}
		if !reflect.DeepEqual(*o.FileFormatCsv, expected) {
			return fmt.Errorf("expected file format csv:\n%+v\ngot:\n%+v", expected, *o.FileFormatCsv)
		}
		return nil
	})
	return s
}
