//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAcc_Experimental_ObjectParameter_UnsetOnDelete verifies that with the
// OBJECT_PARAMETER_UNSET_ON_DELETE experiment enabled, deleting the resource
// uses UNSET instead of resetting to the default value.
func TestAcc_Experimental_ObjectParameter_UnsetOnDelete(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.ObjectParameterResource)).
		WithExperimentalFeaturesEnabled(experimentalfeatures.ObjectParameterUnsetOnDelete)

	config := accconfig.FromModels(t, providerModel) + fmt.Sprintf(`
resource "snowflake_object_parameter" "p" {
	key = "%[1]s"
	value = "123456"
	object_type = "DATABASE"
	object_identifier {
		name = "%[2]s"
	}
}
`, sdk.DatabaseParameterUserTaskTimeoutMs, database.ID().Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: objectParameterUnsetOnDeleteProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_object_parameter.p", "key", string(sdk.DatabaseParameterUserTaskTimeoutMs))),
					assert.Check(resource.TestCheckResourceAttr("snowflake_object_parameter.p", "value", "123456")),
					objectparametersassert.DatabaseParameters(t, database.ID()).
						HasUserTaskTimeoutMs(123456),
				),
			},
			{
				Destroy: true,
				Config:  config,
				Check: assertThat(
					t,
					objectparametersassert.DatabaseParameters(t, database.ID()).
						HasDefaultUserTaskTimeoutMsValue(),
				),
			},
		},
	})
}

// SCHEMA is supported by the parameter catalog write path, but ICEBERG_DEFAULT_DDL_COLLATION is not, so the value is sent as written
// (the behavior before value_type), which Snowflake rejects for a string, and the value rendered by value_type is accepted.
func TestAcc_ObjectParameter_ForwardCompatible_SupportedObjectTypeUnsupportedParameter(t *testing.T) {
	const parameter = sdk.SchemaParameter("ICEBERG_DEFAULT_DDL_COLLATION")
	schemaObject, schemaCleanup := testClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	withoutValueType := model.ObjectParameter("p", string(parameter), "en-ci").WithObject(sdk.ObjectTypeSchema, schemaObject.ID())
	withValueType := model.ObjectParameter("p", string(parameter), "en-ci").WithObject(sdk.ObjectTypeSchema, schemaObject.ID()).WithValueType("STRING")
	withValueTypeInOtherCase := model.ObjectParameter("p", string(parameter), "en-ci").WithObject(sdk.ObjectTypeSchema, schemaObject.ID()).WithValueType("string")
	changed := model.ObjectParameter("p", string(parameter), "en-cs").WithObject(sdk.ObjectTypeSchema, schemaObject.ID()).WithValueType("STRING")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: assertThat(t, objectparametersassert.SchemaParameters(t, schemaObject.ID()).HasDefaultParameterValueOnDefaultLevel(parameter)),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, withoutValueType),
				ExpectError: regexp.MustCompile("SQL compilation error"),
			},
			{
				Config: accconfig.FromModels(t, withValueType),
				Check: assertThat(t,
					resourceassert.ObjectParameterResource(t, withValueType.ResourceReference()).
						HasKeyString(string(parameter)).
						HasValueString("en-ci").
						HasValueTypeString("STRING"),
					objectparametersassert.SchemaParameters(t, schemaObject.ID()).HasStringParameterValue(parameter, "en-ci"),
				),
			},
			// value_type is case-insensitive, so changing only its case shows no difference.
			{
				Config: accconfig.FromModels(t, withValueTypeInOtherCase),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				Config: accconfig.FromModels(t, changed),
				Check: assertThat(t,
					resourceassert.ObjectParameterResource(t, changed.ResourceReference()).HasValueString("en-cs"),
					objectparametersassert.SchemaParameters(t, schemaObject.ID()).HasStringParameterValue(parameter, "en-cs"),
				),
			},
		},
	})
}

// TABLE is not supported by the parameter catalog write path, so DEFAULT_DDL_COLLATION is sent as written without value_type
// (the behavior before value_type), which Snowflake rejects for a string, and the value rendered by value_type is accepted.
func TestAcc_ObjectParameter_ForwardCompatible_UnsupportedObjectType(t *testing.T) {
	const parameter = sdk.SchemaParameter("DEFAULT_DDL_COLLATION")
	table, tableCleanup := testClient().Table.Create(t)
	t.Cleanup(tableCleanup)

	tableParameters := func(assertion assert.SnowflakeParameterAssertion) *assert.SnowflakeParametersAssert[sdk.SchemaObjectIdentifier] {
		tableAssert := assert.NewSnowflakeParametersAssertWithTestClientParametersProvider(table.ID(), sdk.ObjectTypeTable, func(testClient *helpers.TestClient) assert.ParametersProvider[sdk.SchemaObjectIdentifier] {
			return testClient.Parameter.ShowTableParameters
		})
		tableAssert.AddAssertion(assertion)
		return tableAssert
	}
	withoutValueType := model.ObjectParameter("p", string(parameter), "en-ci").WithObject(sdk.ObjectTypeTable, table.ID())
	withValueType := model.ObjectParameter("p", string(parameter), "en-ci").WithObject(sdk.ObjectTypeTable, table.ID()).WithValueType("STRING")
	withObjectTypeInOtherCase := model.ObjectParameter("p", string(parameter), "en-ci").WithObject(sdk.ObjectTypeTable, table.ID()).WithObjectType("table").WithValueType("STRING")
	changed := model.ObjectParameter("p", string(parameter), "en-cs").WithObject(sdk.ObjectTypeTable, table.ID()).WithValueType("STRING")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: assertThat(t, tableParameters(assert.SnowflakeParameterDefaultValueOnDefaultLevelSet(parameter))),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, withoutValueType),
				ExpectError: regexp.MustCompile("SQL compilation error"),
			},
			{
				Config: accconfig.FromModels(t, withValueType),
				Check: assertThat(t,
					resourceassert.ObjectParameterResource(t, withValueType.ResourceReference()).
						HasKeyString(string(parameter)).
						HasValueString("en-ci").
						HasValueTypeString("STRING"),
					tableParameters(assert.SnowflakeParameterValueSet(parameter, "en-ci")),
				),
			},
			// object_type is case-insensitive, so changing only its case shows no difference (and does not recreate the resource).
			{
				Config: accconfig.FromModels(t, withObjectTypeInOtherCase),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				Config: accconfig.FromModels(t, changed),
				Check: assertThat(t,
					resourceassert.ObjectParameterResource(t, changed.ResourceReference()).HasValueString("en-cs"),
					tableParameters(assert.SnowflakeParameterValueSet(parameter, "en-cs")),
				),
			},
		},
	})
}

// value_type is ignored for a parameter that the provider supports, so changing or removing it shows no difference.
func TestAcc_ObjectParameter_ForwardCompatible_ValueTypeIsIgnoredForSupportedParameter(t *testing.T) {
	const parameter = sdk.DatabaseParameterDataRetentionTimeInDays
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	withType := model.ObjectParameter("p", string(parameter), "2").WithObject(sdk.ObjectTypeDatabase, database.ID()).WithValueType("STRING")
	withOtherType := model.ObjectParameter("p", string(parameter), "2").WithObject(sdk.ObjectTypeDatabase, database.ID()).WithValueType("INTEGER")
	withoutType := model.ObjectParameter("p", string(parameter), "2").WithObject(sdk.ObjectTypeDatabase, database.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: assertThat(t, objectparametersassert.DatabaseParameters(t, database.ID()).HasDefaultParameterValue(parameter)),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, withType),
				Check: assertThat(t,
					resourceassert.ObjectParameterResource(t, withType.ResourceReference()).
						HasValueString("2").
						// The diff of value_type is suppressed for a supported parameter, so the field is not stored.
						HasNoValueType(),
					objectparametersassert.DatabaseParameters(t, database.ID()).HasDataRetentionTimeInDays(2),
				),
			},
			{
				Config: accconfig.FromModels(t, withOtherType),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply:             []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				Config: accconfig.FromModels(t, withoutType),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply:             []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}
