//go:build !account_level_tests

package resources_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/tracking"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CompleteUsageTracking(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	comment := random.Comment()

	schemaModel := model.Schema("test", id.DatabaseName(), id.Name())
	schemaModelWithComment := model.Schema("test", id.DatabaseName(), id.Name()).WithComment(comment)

	assertQueryMetadataExistsPrefetched := func(t *testing.T, queryHistory []helpers.QueryHistory, operation tracking.Operation, query string) resource.TestCheckFunc {
		t.Helper()
		return func(state *terraform.State) error {
			expectedMetadata := tracking.NewVersionedResourceMetadata(resources.Schema, operation)
			if _, err := collections.FindFirst(queryHistory, func(history helpers.QueryHistory) bool {
				metadata, err := tracking.ParseMetadata(history.QueryText)
				return err == nil &&
					expectedMetadata == metadata &&
					strings.Contains(history.QueryText, query)
			}); err != nil {
				return fmt.Errorf("query history does not contain query metadata: %v with query containing: %s", expectedMetadata, query)
			}
			return nil
		}
	}

	assertQueryMetadataExists := func(t *testing.T, operation tracking.Operation, query string) resource.TestCheckFunc {
		t.Helper()
		queryHistory := acc.TestClient().InformationSchema.GetQueryHistory(t, 100)
		return assertQueryMetadataExistsPrefetched(t, queryHistory, operation, query)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			// Create
			{
				Config: config.FromModels(t, schemaModel),
				Check: assertThat(t,
					resourceassert.SchemaResource(t, schemaModel.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString(""),
					assert.Check(assertQueryMetadataExists(t, tracking.CreateOperation, fmt.Sprintf(`CREATE SCHEMA %s`, id.FullyQualifiedName()))),
				),
			},
			// Import
			{
				ResourceName: schemaModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedSchemaResource(t, id.FullyQualifiedName()).
						HasCommentString(""),
					assert.CheckImport(func(states []*terraform.InstanceState) error {
						return assertQueryMetadataExists(t, tracking.ImportOperation, fmt.Sprintf(`SHOW SCHEMAS LIKE '%s'`, id.Name()))(nil)
					}),
				),
			},
			// Update + CustomDiff (parameters) + Read
			{
				Config: config.FromModels(t, schemaModelWithComment),
				Check: assertThat(t,
					resourceassert.SchemaResource(t, schemaModelWithComment.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString(comment),
					assert.Check(func(state *terraform.State) error {
						queryHistory := acc.TestClient().InformationSchema.GetQueryHistory(t, 100)
						return errors.Join(
							assertQueryMetadataExistsPrefetched(t, queryHistory, tracking.UpdateOperation, fmt.Sprintf(`ALTER SCHEMA %s SET COMMENT = '%s'`, id.FullyQualifiedName(), comment))(state),
							assertQueryMetadataExistsPrefetched(t, queryHistory, tracking.ReadOperation, fmt.Sprintf(`SHOW SCHEMAS LIKE '%s'`, id.Name()))(state),
							assertQueryMetadataExistsPrefetched(t, queryHistory, tracking.CustomDiffOperation, fmt.Sprintf(`SHOW PARAMETERS IN SCHEMA %s`, id.FullyQualifiedName()))(state),
						)
					}),
				),
			},
			// Delete
			{
				Config:  config.FromModels(t, schemaModelWithComment),
				Destroy: true,
				Check: assertThat(t,
					assert.Check(assertQueryMetadataExists(t, tracking.DeleteOperation, fmt.Sprintf(`DROP SCHEMA IF EXISTS %s`, id.FullyQualifiedName()))),
				),
			},
		},
	})
}
