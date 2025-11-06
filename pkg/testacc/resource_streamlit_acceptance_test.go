//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/invokeactionassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Streamlit_BasicUseCase(t *testing.T) {
	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	// warehouse is needed because default warehouse uses lowercase, and it fails in snowflake.
	// TODO(SNOW-1541938): use a default warehouse after fix on snowflake side
	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)

	networkRule, networkRuleCleanup := testClient().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)

	externalAccessIntegrationId, externalAccessIntegrationCleanup := testClient().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
	t.Cleanup(externalAccessIntegrationCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	title := random.AlphaN(4)
	directoryLocation := "abc"
	rootLocationWithCatalog := fmt.Sprintf("%s/%s", stage.Location(), directoryLocation)
	mainFile := "foo"

	basic := model.StreamlitWithIds("test", id, mainFile, stage.ID())

	complete := model.StreamlitWithIds("test", newId, mainFile, stage.ID()).
		WithComment(comment).
		WithTitle(title).
		WithDirectoryLocation(directoryLocation).
		WithQueryWarehouse(warehouse.ID().Name()).
		WithExternalAccessIntegrations(externalAccessIntegrationId)

	assertBasic := []assert.TestCheckFuncProvider{
		objectassert.Streamlit(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasTitle("").
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment("").
			HasQueryWarehouse("").
			HasUrlIdNotEmpty(),

		resourceassert.StreamlitResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasStageString(stage.ID().FullyQualifiedName()).
			HasMainFileString(mainFile).
			HasDirectoryLocationString("").
			HasQueryWarehouseString("").
			HasTitleString("").
			HasCommentString("").
			HasExternalAccessIntegrationsEmpty(),

		resourceshowoutputassert.StreamlitShowOutput(t, basic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasTitle("").
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment("").
			HasQueryWarehouse("").
			HasUrlIdNotEmpty().
			HasOwnerRoleType("ROLE"),

		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.name", id.Name())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.title", "")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.root_location", stage.Location())),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.main_file", mainFile)),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.query_warehouse", "")),
		assert.Check(resource.TestCheckResourceAttrSet(basic.ResourceReference(), "describe_output.0.url_id")),
		assert.Check(resource.TestCheckResourceAttrSet(basic.ResourceReference(), "describe_output.0.default_packages")),
		assert.Check(resource.TestCheckResourceAttrSet(basic.ResourceReference(), "describe_output.0.user_packages.#")),
		assert.Check(resource.TestCheckResourceAttrSet(basic.ResourceReference(), "describe_output.0.import_urls.#")),
		assert.Check(resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.0.external_access_integrations.#", "0")),
		assert.Check(resource.TestCheckResourceAttrSet(basic.ResourceReference(), "describe_output.0.external_access_secrets")),
	}

	assertComplete := []assert.TestCheckFuncProvider{
		objectassert.Streamlit(t, newId).
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()).
			HasTitle(title).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment(comment).
			HasQueryWarehouse(warehouse.ID().Name()).
			HasUrlIdNotEmpty(),

		resourceassert.StreamlitResource(t, complete.ResourceReference()).
			HasNameString(newId.Name()).
			HasFullyQualifiedNameString(newId.FullyQualifiedName()).
			HasDatabaseString(newId.DatabaseName()).
			HasSchemaString(newId.SchemaName()).
			HasStageString(stage.ID().FullyQualifiedName()).
			HasMainFileString(mainFile).
			HasDirectoryLocationString(directoryLocation).
			HasQueryWarehouseString(warehouse.ID().Name()).
			HasTitleString(title).
			HasCommentString(comment).
			HasExternalAccessIntegrations([]string{externalAccessIntegrationId.Name()}),

		resourceshowoutputassert.StreamlitShowOutput(t, complete.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()).
			HasTitle(title).
			HasOwner(testClient().Context.CurrentRole(t).Name()).
			HasComment(comment).
			HasQueryWarehouse(warehouse.ID().Name()).
			HasUrlIdNotEmpty().
			HasOwnerRoleType("ROLE"),

		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.name", newId.Name())),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.title", title)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.root_location", rootLocationWithCatalog)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.main_file", mainFile)),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.query_warehouse", warehouse.ID().Name())),
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.url_id")),
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.default_packages")),
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.user_packages.#")),
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.import_urls.#")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.external_access_integrations.#", "1")),
		assert.Check(resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.0.external_access_integrations.0", externalAccessIntegrationId.Name())),
		assert.Check(resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.external_access_secrets")),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Streamlit),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:            accconfig.FromModels(t, basic),
				ResourceName:      basic.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
			// Import - with optionals
			{
				Config:            accconfig.FromModels(t, complete),
				ResourceName:      complete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Update - detect external changes
			{
				PreConfig: func() {
					testClient().Streamlit.Update(t, sdk.NewAlterStreamlitRequest(id).WithSet(
						*sdk.NewStreamlitSetRequest().
							WithRootLocation(rootLocationWithCatalog).
							WithTitle(title).
							WithQueryWarehouse(warehouse.ID()).
							WithExternalAccessIntegrations(*sdk.NewExternalAccessIntegrationsRequest([]sdk.AccountObjectIdentifier{externalAccessIntegrationId})).
							WithComment(comment),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Destroy - ensure streamlit is destroyed before the next step
			{
				Destroy: true,
				Config:  accconfig.FromModels(t, basic),
				Check: assertThat(t,
					invokeactionassert.StreamlitDoesNotExist(t, id),
				),
			},
			// Create - with optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}

func TestAcc_Streamlit_InvalidStage(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	streamlitModel := model.Streamlit("test", id.DatabaseId().FullyQualifiedName(), id.SchemaId().FullyQualifiedName(), id.Name(), "some", "some")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck: func() {
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, streamlitModel),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Invalid identifier type`),
			},
		},
	})
}

func TestAcc_Streamlit_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()

	streamlitModel := model.StreamlitWithIds("test", id, "main_file", stage.ID())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Streamlit),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            accconfig.FromModels(t, streamlitModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(streamlitModel.ResourceReference(), "id", helpers.EncodeSnowflakeID(id)),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, streamlitModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(streamlitModel.ResourceReference(), "id", id.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_Streamlit_IdentifierQuotingDiffSuppression(t *testing.T) {
	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	quotedDatabaseName := fmt.Sprintf(`"%s"`, id.DatabaseName())
	quotedSchemaName := fmt.Sprintf(`"%s"`, id.SchemaName())
	quotedName := fmt.Sprintf(`"%s"`, id.Name())

	streamlitModel := model.Streamlit("test", quotedDatabaseName, quotedSchemaName, quotedName, "main_file", stage.ID().FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Streamlit),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             accconfig.FromModels(t, streamlitModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(streamlitModel.ResourceReference(), "database", fmt.Sprintf("\"%s\"", id.DatabaseName())),
					resource.TestCheckResourceAttr(streamlitModel.ResourceReference(), "schema", id.SchemaName()),
					resource.TestCheckResourceAttr(streamlitModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(streamlitModel.ResourceReference(), "id", helpers.EncodeSnowflakeID(id)),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, streamlitModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(streamlitModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(streamlitModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(streamlitModel.ResourceReference(), "database", fmt.Sprintf("\"%s\"", id.DatabaseName())),
					resource.TestCheckResourceAttr(streamlitModel.ResourceReference(), "schema", id.SchemaName()),
					resource.TestCheckResourceAttr(streamlitModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(streamlitModel.ResourceReference(), "id", id.FullyQualifiedName()),
				),
			},
		},
	})
}
