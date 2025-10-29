package invokeactionassert

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TODO [SNOW-1501905]: generalize this type of assertion
type queryHistoryCheck struct {
	expectedSql string
	limit       int

	// TODO [SNOW-1501905]: test client passed here temporarily to be able to check secondary (by default our assertions use the default one)
	testClient *helpers.TestClient
}

func (w *queryHistoryCheck) ToTerraformTestCheckFunc(t *testing.T, _ *helpers.TestClient) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		return w.checkQueryHistoryEntry(t)
	}
}

func (w *queryHistoryCheck) ToTerraformImportStateCheckFunc(t *testing.T, _ *helpers.TestClient) resource.ImportStateCheckFunc {
	t.Helper()
	return func(s []*terraform.InstanceState) error {
		return w.checkQueryHistoryEntry(t)
	}
}

func (w *queryHistoryCheck) checkQueryHistoryEntry(t *testing.T) error {
	t.Helper()
	if w.testClient == nil {
		return errors.New("testClient must not be nil")
	}
	queryHistory := w.testClient.InformationSchema.GetQueryHistory(t, w.limit)
	if _, err := collections.FindFirst(queryHistory, func(history helpers.QueryHistory) bool {
		if strings.Contains(history.QueryText, w.expectedSql) {
			return true
		}
		return false
	}); err != nil {
		return fmt.Errorf("query history does not contain query containing: %v", w.expectedSql)
	}
	return nil
}

func QueryHistoryEntry(t *testing.T, testClient *helpers.TestClient, sql string, limit int) assert.TestCheckFuncProvider {
	t.Helper()
	return &queryHistoryCheck{expectedSql: sql, limit: limit, testClient: testClient}
}

func QueryHistoryEntryInImport(t *testing.T, testClient *helpers.TestClient, sql string, limit int) assert.ImportStateCheckFuncProvider {
	t.Helper()
	return &queryHistoryCheck{expectedSql: sql, limit: limit, testClient: testClient}
}
