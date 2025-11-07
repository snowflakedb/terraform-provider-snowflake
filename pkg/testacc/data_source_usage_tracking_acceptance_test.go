//go:build non_account_level_tests

package testacc

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/tracking"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CompleteUsageTracking_Datasource(t *testing.T) {
	schemaId := testClient().Ids.RandomDatabaseObjectIdentifier()
	schemaModel := model.Schema("test", schemaId.DatabaseName(), schemaId.Name())
	schemasModel := datasourcemodel.Schemas("test").
		WithLike(schemaId.Name()).
		WithInDatabase(schemaId.DatabaseId()).
		WithDependsOn(schemaModel.ResourceReference())

	assertQueryMetadataExists := func(t *testing.T, queryHistory []helpers.QueryHistory, query string) resource.TestCheckFunc {
		t.Helper()
		return func(state *terraform.State) error {
			expectedMetadata := tracking.NewVersionedDatasourceMetadata(datasources.Schemas)
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, schemaModel, schemasModel),
				Check: assertThat(t,
					resourceassert.SchemaResource(t, schemaModel.ResourceReference()).HasNameString(schemaId.Name()),
					assert.Check(func(state *terraform.State) error {
						queryHistory := testClient().InformationSchema.GetQueryHistory(t, 100)
						return errors.Join(
							assertQueryMetadataExists(t, queryHistory, fmt.Sprintf(`SHOW SCHEMAS LIKE '%s' IN DATABASE "%s"`, schemaId.Name(), schemaId.DatabaseName()))(state),
							// SHOW PARAMETERS IN SCHEMA "acc_test_db_AT_1AB7E1DE_1A10_89C3_C13C_899754A250B6"."FPGDHEAT_1AB7E1DE_1A10_89C3_C13C_899754A250B6" --terraform_provider_usage_tracking {"json_schema_version":"1","version":"v0.99.0","datasource":"snowflake_schemas","operation":"read"}
							assertQueryMetadataExists(t, queryHistory, fmt.Sprintf(`SHOW PARAMETERS IN SCHEMA %s`, schemaId.FullyQualifiedName()))(state),
							assertQueryMetadataExists(t, queryHistory, fmt.Sprintf(`DESCRIBE SCHEMA %s`, schemaId.FullyQualifiedName()))(state),
						)
					}),
				),
			},
		},
	})
}
