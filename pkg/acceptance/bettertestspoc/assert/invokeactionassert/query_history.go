package invokeactionassert

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/tracking"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TODO [SNOW-1501905]: generalize this type of assertion (extract as query history object that can have assertions run on it)
type queryHistoryCheck struct {
	expectedSql       string
	expectedOperation tracking.Operation
	limit             int

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
		expectedMetadata := tracking.NewVersionedResourceMetadata(resources.Warehouse, w.expectedOperation)
		if strings.Contains(history.QueryText, w.expectedSql) {
			metadata, err := tracking.ParseMetadata(history.QueryText)
			if err != nil {
				return false
			}
			return expectedMetadata == metadata
		}
		return false
	}); err != nil {
		return fmt.Errorf("query history does not contain query containing: %s with operation: %s", w.expectedSql, w.expectedOperation)
	}
	return nil
}

func QueryHistoryEntry(t *testing.T, testClient *helpers.TestClient, sql string, operation tracking.Operation, limit int) assert.TestCheckFuncProvider {
	t.Helper()
	return &queryHistoryCheck{expectedSql: sql, expectedOperation: operation, limit: limit, testClient: testClient}
}

func QueryHistoryEntryInImport(t *testing.T, testClient *helpers.TestClient, sql string, operation tracking.Operation, limit int) assert.ImportStateCheckFuncProvider {
	t.Helper()
	return &queryHistoryCheck{expectedSql: sql, expectedOperation: operation, limit: limit, testClient: testClient}
}
