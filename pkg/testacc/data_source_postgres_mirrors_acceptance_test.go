//go:build account_level_tests

package testacc

import (
	"maps"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_PostgresMirrors_CompleteUseCase(t *testing.T) {
	instance := testClient().PostgresMirror.FixtureInstance(t)
	mirror := testClient().PostgresMirror.CreateWithTables(t, instance, helpers.PostgresMirrorFixtureTable1)

	mirrorsModel := datasourcemodel.PostgresMirrors("test").
		WithPostgresInstance(instance)
	mirrorsModelWithoutDescribe := datasourcemodel.PostgresMirrors("test").
		WithPostgresInstance(instance).
		WithWithDescribe(false)

	// LIST_MIRRORS filters only by instance, so other mirrors on the fixture can be listed too; the mirror is matched by name.
	showOutput := map[string]string{
		"show_output.#":                   "1",
		"show_output.0.name":              mirror.Name,
		"show_output.0.postgres_instance": instance,
		"show_output.0.postgres_database": helpers.PostgresMirrorFixtureDatabase,
		"show_output.0.target_database":   mirror.TargetDatabase,
		"show_output.0.refresh_interval":  helpers.PostgresMirrorDefaultRefreshInterval,
	}
	withDescribeOutput := func(values map[string]string) map[string]string {
		result := map[string]string{
			"describe_output.#":                   "1",
			"describe_output.0.name":              mirror.Name,
			"describe_output.0.postgres_instance": instance,
			"describe_output.0.target_database":   mirror.TargetDatabase,
			"describe_output.0.target_table_type": string(sdk.PostgresMirrorTargetTableTypeSnowflake),
			"describe_output.0.postgres_tables.#": "1",
			"describe_output.0.postgres_tables.0": helpers.PostgresMirrorFixtureTable1,
			"describe_output.0.refresh_interval":  helpers.PostgresMirrorDefaultRefreshInterval,
		}
		maps.Copy(result, values)
		return result
	}
	withoutDescribeOutput := func(values map[string]string) map[string]string {
		result := map[string]string{"describe_output.#": "0"}
		maps.Copy(result, values)
		return result
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// With describe
			{
				Config: accconfig.FromModels(t, mirrorsModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckTypeSetElemNestedAttrs(mirrorsModel.DatasourceReference(), "postgres_mirrors.*", withDescribeOutput(showOutput))),
				),
			},
			// Without describe
			{
				Config: accconfig.FromModels(t, mirrorsModelWithoutDescribe),
				Check: assertThat(t,
					assert.Check(resource.TestCheckTypeSetElemNestedAttrs(mirrorsModelWithoutDescribe.DatasourceReference(), "postgres_mirrors.*", withoutDescribeOutput(showOutput))),
				),
			},
		},
	})
}

func TestAcc_PostgresMirrors_BasicUseCase_DifferentFiltering(t *testing.T) {
	instance := testClient().PostgresMirror.FixtureInstance(t)
	mirror := testClient().PostgresMirror.CreateWithTables(t, instance, helpers.PostgresMirrorFixtureTable1)
	otherInstance := testClient().Ids.Alpha()

	byInstance := datasourcemodel.PostgresMirrors("test").
		WithWithDescribe(false).
		WithPostgresInstance(instance)
	byQuotedInstance := datasourcemodel.PostgresMirrors("test").
		WithWithDescribe(false).
		WithPostgresInstance(sdk.NewAccountObjectIdentifier(instance).FullyQualifiedName())
	byOtherInstance := datasourcemodel.PostgresMirrors("test").
		WithWithDescribe(false).
		WithPostgresInstance(otherInstance)

	// LIST_MIRRORS filters only by instance, and other mirrors on the fixture can be listed too, so the mirror is matched by name.
	listed := map[string]string{"show_output.0.name": mirror.Name}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// instance
			{
				Config: accconfig.FromModels(t, byInstance),
				Check: assertThat(t,
					assert.Check(resource.TestCheckTypeSetElemNestedAttrs(byInstance.DatasourceReference(), "postgres_mirrors.*", listed)),
				),
			},
			// quoted instance
			{
				Config: accconfig.FromModels(t, byQuotedInstance),
				Check: assertThat(t,
					assert.Check(resource.TestCheckTypeSetElemNestedAttrs(byQuotedInstance.DatasourceReference(), "postgres_mirrors.*", listed)),
				),
			},
			// LIST_MIRRORS rejects an instance that does not exist instead of returning no mirrors, and names it in the error.
			{
				Config:      accconfig.FromModels(t, byOtherInstance),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta(otherInstance)),
			},
		},
	})
}
