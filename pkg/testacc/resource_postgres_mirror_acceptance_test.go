//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/stretchr/testify/require"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_PostgresMirror_BasicUseCase(t *testing.T) {
	instance := testClient().PostgresMirror.FixtureInstance(t)

	name := testClient().Ids.RandomPostgresMirrorName()
	targetDatabase := testClient().Ids.RandomAccountObjectIdentifier()

	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)
	t.Cleanup(testClient().Grant.GrantUsageOnWarehouseToSnowflakeApplication(t, warehouse.ID()))

	newModel := func() *model.PostgresMirrorModel {
		return model.PostgresMirror("test", name, helpers.PostgresMirrorFixtureDatabase, instance, targetDatabase.Name())
	}
	basic := newModel().
		WithPostgresTables(helpers.PostgresMirrorFixtureTable1, helpers.PostgresMirrorFixtureTable2)
	withOptionals := newModel().
		WithPostgresTables(helpers.PostgresMirrorFixtureTable2, helpers.PostgresMirrorFixtureTable3).
		WithPostgresSchemas(helpers.PostgresMirrorFixtureSchemaA).
		WithRefreshIntervalMinutes(1).
		WithWarehouse(warehouse.ID().Name())
	complete := newModel().
		WithPostgresTables(helpers.PostgresMirrorFixtureTable1, helpers.PostgresMirrorFixtureMixedCaseTable).
		WithPostgresSchemas(helpers.PostgresMirrorFixtureSchemaB).
		WithRefreshIntervalSeconds(30).
		WithTargetTableType(string(sdk.PostgresMirrorTargetTableTypeIceberg)).
		WithWarehouse(warehouse.ID().Name()).
		WithDropTargetDatabaseOnDestroy(true)

	ref := basic.ResourceReference()

	basicResourceAssert := func() *resourceassert.PostgresMirrorResourceAssert {
		return resourceassert.PostgresMirrorResource(t, ref).
			HasNameString(name).
			HasPostgresInstanceString(instance).
			HasPostgresDatabaseString(helpers.PostgresMirrorFixtureDatabase).
			HasTargetDatabaseString(targetDatabase.Name()).
			HasPostgresTables(helpers.PostgresMirrorFixtureTable1, helpers.PostgresMirrorFixtureTable2).
			HasPostgresSchemasEmpty().
			HasNoTargetTableType().
			HasNoExternalVolume().
			HasDropTargetDatabaseOnDestroyString("true")
	}
	basicRemoteAsserts := func(refreshInterval string) []assert.TestCheckFuncProvider {
		return []assert.TestCheckFuncProvider{
			resourceshowoutputassert.PostgresMirrorDescribeOutput(t, ref).
				HasName(name).
				HasPostgresInstance(instance).
				HasPostgresDatabase(helpers.PostgresMirrorFixtureDatabase).
				HasTargetDatabase(targetDatabase.Name()).
				HasStatusNotEmpty().
				HasRefreshInterval(refreshInterval).
				HasTargetTableType(sdk.PostgresMirrorTargetTableTypeSnowflake).
				HasExternalVolume("").
				HasWarehouse("").
				HasPostgresTables(helpers.PostgresMirrorFixtureTable1, helpers.PostgresMirrorFixtureTable2).
				HasNoPostgresSchemas(),
			objectassert.PostgresMirrorDetails(t, name).
				HasPostgresInstance(instance).
				HasPostgresDatabase(helpers.PostgresMirrorFixtureDatabase).
				HasTargetDatabase(targetDatabase.Name()).
				HasRefreshInterval(refreshInterval).
				HasTargetTableType(sdk.PostgresMirrorTargetTableTypeSnowflake).
				HasNoExternalVolume().
				HasNoWarehouse().
				HasPostgresTables(helpers.PostgresMirrorFixtureTable1, helpers.PostgresMirrorFixtureTable2).
				HasNoPostgresSchemas(),
		}
	}
	assertBasic := append([]assert.TestCheckFuncProvider{
		basicResourceAssert().
			HasNoWarehouse(),
	}, basicRemoteAsserts(sdk.PostgresMirrorDefaultRefreshInterval)...)
	// A value removed from the configuration stays in state as an empty string, unlike one that was never set.
	// Removing the refresh interval keeps the current interval, because ALTER_MIRROR cannot reset it.
	assertAfterUnset := append([]assert.TestCheckFuncProvider{
		basicResourceAssert().
			HasRefreshIntervalMinutes(1).
			HasWarehouseEmpty(),
	}, basicRemoteAsserts("1 minute")...)

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.PostgresMirrorResource(t, ref).
			HasNameString(name).
			HasPostgresTables(helpers.PostgresMirrorFixtureTable2, helpers.PostgresMirrorFixtureTable3).
			HasPostgresSchemas(helpers.PostgresMirrorFixtureSchemaA).
			HasRefreshIntervalMinutes(1).
			HasWarehouseString(warehouse.ID().Name()).
			HasNoTargetTableType().
			HasDropTargetDatabaseOnDestroyString("true"),
		resourceshowoutputassert.PostgresMirrorDescribeOutput(t, ref).
			HasName(name).
			HasRefreshInterval("1 minute").
			HasWarehouse(warehouse.ID().Name()).
			HasPostgresTables(helpers.PostgresMirrorFixtureTable2, helpers.PostgresMirrorFixtureTable3).
			HasPostgresSchemas(helpers.PostgresMirrorFixtureSchemaA),
		objectassert.PostgresMirrorDetails(t, name).
			HasRefreshInterval("1 minute").
			HasWarehouse(warehouse.ID().Name()).
			HasPostgresTables(helpers.PostgresMirrorFixtureTable2, helpers.PostgresMirrorFixtureTable3).
			HasPostgresSchemas(helpers.PostgresMirrorFixtureSchemaA),
	}

	assertComplete := []assert.TestCheckFuncProvider{
		resourceassert.PostgresMirrorResource(t, ref).
			HasNameString(name).
			HasPostgresInstanceString(instance).
			HasPostgresDatabaseString(helpers.PostgresMirrorFixtureDatabase).
			HasTargetDatabaseString(targetDatabase.Name()).
			HasPostgresTables(helpers.PostgresMirrorFixtureTable1, helpers.PostgresMirrorFixtureMixedCaseTable).
			HasPostgresSchemas(helpers.PostgresMirrorFixtureSchemaB).
			HasRefreshIntervalSeconds(30).
			HasTargetTableTypeString(string(sdk.PostgresMirrorTargetTableTypeIceberg)).
			HasNoExternalVolume().
			HasWarehouseString(warehouse.ID().Name()).
			HasDropTargetDatabaseOnDestroyString("true"),
		resourceshowoutputassert.PostgresMirrorDescribeOutput(t, ref).
			HasName(name).
			HasTargetDatabase(targetDatabase.Name()).
			HasStatusNotEmpty().
			HasRefreshInterval("30 seconds").
			HasTargetTableType(sdk.PostgresMirrorTargetTableTypeIceberg).
			HasExternalVolume("").
			HasWarehouse(warehouse.ID().Name()).
			HasPostgresTables(helpers.PostgresMirrorFixtureMixedCaseTable, helpers.PostgresMirrorFixtureTable1).
			HasPostgresSchemas(helpers.PostgresMirrorFixtureSchemaB),
		objectassert.PostgresMirrorDetails(t, name).
			HasTargetDatabase(targetDatabase.Name()).
			HasRefreshInterval("30 seconds").
			HasTargetTableType(sdk.PostgresMirrorTargetTableTypeIceberg).
			HasNoExternalVolume().
			HasWarehouse(warehouse.ID().Name()).
			HasPostgresTables(helpers.PostgresMirrorFixtureMixedCaseTable, helpers.PostgresMirrorFixtureTable1).
			HasPostgresSchemas(helpers.PostgresMirrorFixtureSchemaB),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresMirror),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: config.FromModels(t, basic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Check: assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				// The describe output carries the mirror's runtime state (status, lag, last apply), which changes between reads.
				ImportStateVerifyIgnore: []string{r.DescribeOutputAttributeName},
			},
			// Update - set optionals and change the tables and schemas in place
			{
				Config: config.FromModels(t, withOptionals),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, assertWithOptionals...),
			},
			// Import - with optionals
			{
				Config:                  config.FromModels(t, withOptionals),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{r.DescribeOutputAttributeName},
			},
			// Update - unset optionals
			{
				Config: config.FromModels(t, basic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "warehouse", tfjson.ActionUpdate, sdk.String(warehouse.ID().Name()), nil),
					},
				},
				Check: assertThat(t, assertAfterUnset...),
			},
			// Destroy
			{
				Destroy: true,
				Config:  config.FromModels(t, basic),
			},
			// Create - with all fields, reusing the target database name the default destroy freed
			{
				PreConfig: func() {
					_, err := testClient().Database.Show(t, targetDatabase)
					require.ErrorIs(t, err, sdk.ErrObjectNotFound)
				},
				Config: config.FromModels(t, complete),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Check: assertThat(t, assertComplete...),
			},
			// Import - with all fields
			{
				Config:                  config.FromModels(t, complete),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{r.DescribeOutputAttributeName},
			},
			// Update - external changes
			{
				PreConfig: func() {
					testClient().PostgresMirror.Alter(t, sdk.NewAlterPostgresMirrorRequest(*sdk.NewAlterPostgresMirrorArgumentsRequest(name).
						WithRefreshInterval("5 minutes").
						WithAddTables(helpers.PostgresMirrorNames(helpers.PostgresMirrorFixtureTable2)).
						WithUnsetWarehouse(true),
					))
				},
				Config: config.FromModels(t, complete),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "refresh_interval.0.minutes", sdk.String("0"), sdk.String("5")),
						planchecks.ExpectChange(ref, "refresh_interval.0.seconds", tfjson.ActionUpdate, sdk.String("0"), sdk.String("30")),
					},
				},
				Check: assertThat(t, assertComplete...),
			},
		},
	})
}

func TestAcc_PostgresMirror_DefaultTargetTableType(t *testing.T) {
	instance := testClient().PostgresMirror.FixtureInstance(t)

	name := testClient().Ids.RandomPostgresMirrorName()
	targetDatabase := testClient().Ids.RandomAccountObjectIdentifier()

	newModel := func() *model.PostgresMirrorModel {
		return model.PostgresMirror("test", name, helpers.PostgresMirrorFixtureDatabase, instance, targetDatabase.Name()).
			WithPostgresTables(helpers.PostgresMirrorFixtureTable1)
	}
	withDefaultTargetTableType := newModel().WithTargetTableType(string(sdk.PostgresMirrorTargetTableTypeSnowflake))
	withoutTargetTableType := newModel()

	ref := withDefaultTargetTableType.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresMirror),
		Steps: []resource.TestStep{
			// Create - with the default target table type set explicitly
			{
				Config: config.FromModels(t, withDefaultTargetTableType),
				Check: assertThat(t,
					resourceassert.PostgresMirrorResource(t, ref).
						HasTargetTableTypeString(string(sdk.PostgresMirrorTargetTableTypeSnowflake)),
				),
			},
			// Update - removing the explicit default replaces the mirror, because target_table_type is ForceNew and an omitted value is not compared with the Snowflake default
			{
				Config: config.FromModels(t, withoutTargetTableType),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
		},
	})
}

func TestAcc_PostgresMirror_KeptTargetDatabase(t *testing.T) {
	instance := testClient().PostgresMirror.FixtureInstance(t)

	name := testClient().Ids.RandomPostgresMirrorName()
	targetDatabase := testClient().Ids.RandomAccountObjectIdentifier()
	newTargetDatabase := testClient().Ids.RandomAccountObjectIdentifier()
	t.Cleanup(testClient().PostgresMirror.DropKeptTargetDatabaseFunc(t, targetDatabase))
	t.Cleanup(testClient().PostgresMirror.DropKeptTargetDatabaseFunc(t, newTargetDatabase))

	newModel := func(targetDatabase sdk.AccountObjectIdentifier) *model.PostgresMirrorModel {
		return model.PostgresMirror("test", name, helpers.PostgresMirrorFixtureDatabase, instance, targetDatabase.Name()).
			WithPostgresTables(helpers.PostgresMirrorFixtureTable1).
			WithDropTargetDatabaseOnDestroy(false)
	}
	kept := newModel(targetDatabase)
	replaced := newModel(newTargetDatabase)

	ref := kept.ResourceReference()

	assertTargetDatabase := func(targetDatabase sdk.AccountObjectIdentifier) resource.TestCheckFunc {
		return assertThat(t,
			resourceassert.PostgresMirrorResource(t, ref).
				HasTargetDatabaseString(targetDatabase.Name()).
				HasDropTargetDatabaseOnDestroyString("false"),
			objectassert.PostgresMirrorDetails(t, name).
				HasTargetDatabase(targetDatabase.Name()),
		)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresMirror),
		Steps: []resource.TestStep{
			// Create - keep the target database on destroy
			{
				Config: config.FromModels(t, kept),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Check: assertTargetDatabase(targetDatabase),
			},
			// Update - a new target database replaces the mirror and keeps the old database
			{
				Config: config.FromModels(t, replaced),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assertTargetDatabase(newTargetDatabase),
					assertThat(t, objectassert.Database(t, targetDatabase)),
				),
			},
			// Destroy - both target databases are kept
			{
				Destroy: true,
				Config:  config.FromModels(t, replaced),
				Check:   assertThat(t, objectassert.Database(t, targetDatabase), objectassert.Database(t, newTargetDatabase)),
			},
		},
	})
}

func TestAcc_PostgresMirror_RecreatedOutsideTerraform(t *testing.T) {
	instance := testClient().PostgresMirror.FixtureInstance(t)

	name := testClient().Ids.RandomPostgresMirrorName()
	targetDatabase := testClient().Ids.RandomAccountObjectIdentifier()
	otherTargetDatabase := testClient().Ids.RandomAccountObjectIdentifier()

	basic := model.PostgresMirror("test", name, helpers.PostgresMirrorFixtureDatabase, instance, targetDatabase.Name()).
		WithPostgresTables(helpers.PostgresMirrorFixtureTable1)
	ref := basic.ResourceReference()

	assertConfiguredTargetDatabase := assertThat(t,
		resourceassert.PostgresMirrorResource(t, ref).
			HasTargetDatabaseString(targetDatabase.Name()),
		objectassert.PostgresMirrorDetails(t, name).
			HasTargetDatabase(targetDatabase.Name()),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresMirror),
		Steps: []resource.TestStep{
			// Create - basic
			{
				Config: config.FromModels(t, basic),
				Check:  assertConfiguredTargetDatabase,
			},
			// Update - the mirror recreated outside of Terraform under the same name with another target database
			{
				PreConfig: func() {
					testClient().PostgresMirror.DropFunc(t, name)()
					testClient().PostgresMirror.CreateWithRequest(t, sdk.NewCreatePostgresMirrorRequest(
						*sdk.NewCreatePostgresMirrorArgumentsRequest(name, instance, helpers.PostgresMirrorFixtureDatabase, otherTargetDatabase.Name()).
							WithPostgresTables(helpers.PostgresMirrorNames(helpers.PostgresMirrorFixtureTable1)),
					))
				},
				Config: config.FromModels(t, basic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(ref, "target_database", sdk.String(targetDatabase.Name()), sdk.String(otherTargetDatabase.Name())),
					},
				},
				Check: assertConfiguredTargetDatabase,
			},
		},
	})
}

func TestAcc_PostgresMirror_ExternalVolume(t *testing.T) {
	instance := testClient().PostgresMirror.FixtureInstance(t)

	name := testClient().Ids.RandomPostgresMirrorName()
	targetDatabase := testClient().Ids.RandomAccountObjectIdentifier()

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)
	t.Cleanup(testClient().Grant.GrantUsageOnExternalVolumeToSnowflakeApplication(t, externalVolumeId))

	iceberg := model.PostgresMirror("test", name, helpers.PostgresMirrorFixtureDatabase, instance, targetDatabase.Name()).
		WithPostgresTables(helpers.PostgresMirrorFixtureTable1).
		WithTargetTableType(string(sdk.PostgresMirrorTargetTableTypeIceberg)).
		WithExternalVolumeSdk(externalVolumeId)
	ref := iceberg.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresMirror),
		Steps: []resource.TestStep{
			// Create - Iceberg target tables on an external volume
			{
				Config: config.FromModels(t, iceberg),
				Check: assertThat(t,
					resourceassert.PostgresMirrorResource(t, ref).
						HasTargetTableTypeString(string(sdk.PostgresMirrorTargetTableTypeIceberg)).
						HasExternalVolumeString(externalVolumeId.Name()),
					// Create sends the external volume quoted, and DESCRIBE_MIRROR returns it as sent.
					resourceshowoutputassert.PostgresMirrorDescribeOutput(t, ref).
						HasTargetTableType(sdk.PostgresMirrorTargetTableTypeIceberg).
						HasExternalVolume(externalVolumeId.FullyQualifiedName()),
					objectassert.PostgresMirrorDetails(t, name).
						HasTargetTableType(sdk.PostgresMirrorTargetTableTypeIceberg).
						HasExternalVolume(externalVolumeId.FullyQualifiedName()),
				),
			},
			// Import - the external volume is imported quoted, as DESCRIBE_MIRROR returns it
			{
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"drop_target_database_on_destroy", "external_volume", r.DescribeOutputAttributeName},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedPostgresMirrorResource(t, name).
						HasExternalVolumeString(externalVolumeId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_PostgresMirror_ExternalVolumeWithoutIceberg(t *testing.T) {
	instance := testClient().PostgresMirror.FixtureInstance(t)

	name := testClient().Ids.RandomPostgresMirrorName()
	targetDatabase := testClient().Ids.RandomAccountObjectIdentifier()
	externalVolumeId := testClient().Ids.RandomAccountObjectIdentifier()

	withoutTargetTableType := model.PostgresMirror("test", name, helpers.PostgresMirrorFixtureDatabase, instance, targetDatabase.Name()).
		WithPostgresTables(helpers.PostgresMirrorFixtureTable1).
		WithExternalVolumeSdk(externalVolumeId)
	withSnowflakeTargetTableType := model.PostgresMirror("test", name, helpers.PostgresMirrorFixtureDatabase, instance, targetDatabase.Name()).
		WithPostgresTables(helpers.PostgresMirrorFixtureTable1).
		WithTargetTableType(string(sdk.PostgresMirrorTargetTableTypeSnowflake)).
		WithExternalVolumeSdk(externalVolumeId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresMirror),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, withoutTargetTableType),
				ExpectError: regexp.MustCompile("ExternalVolume can only be set when TargetTableType is ICEBERG"),
			},
			{
				Config:      config.FromModels(t, withSnowflakeTargetTableType),
				ExpectError: regexp.MustCompile("ExternalVolume can only be set when TargetTableType is ICEBERG"),
			},
		},
	})
}

func TestAcc_PostgresMirror_RefreshIntervalValidations(t *testing.T) {
	instance := testClient().PostgresMirror.FixtureInstance(t)

	name := testClient().Ids.RandomPostgresMirrorName()
	targetDatabase := testClient().Ids.RandomAccountObjectIdentifier()

	newModel := func(interval map[string]tfconfig.Variable) *model.PostgresMirrorModel {
		return model.PostgresMirror("test", name, helpers.PostgresMirrorFixtureDatabase, instance, targetDatabase.Name()).
			WithPostgresTables(helpers.PostgresMirrorFixtureTable1).
			WithRefreshIntervalValue(tfconfig.ObjectVariable(interval))
	}
	twoUnits := newModel(map[string]tfconfig.Variable{
		"minutes": tfconfig.IntegerVariable(1),
		"seconds": tfconfig.IntegerVariable(30),
	})
	noUnit := newModel(map[string]tfconfig.Variable{
		"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PostgresMirror),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, twoUnits),
				ExpectError: regexp.MustCompile("conflicts with refresh_interval.0.seconds"),
			},
			{
				Config:      config.FromModels(t, noUnit),
				ExpectError: regexp.MustCompile(`"refresh_interval.0.hours": one of\s+` + "`refresh_interval.0.hours,refresh_interval.0.minutes,refresh_interval.0.seconds`"),
			},
		},
	})
}
