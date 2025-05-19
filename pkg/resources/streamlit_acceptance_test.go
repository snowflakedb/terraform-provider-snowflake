//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Streamlit_basic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stage, stageCleanup := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	// warehouse is needed because default warehouse uses lowercase, and it fails in snowflake.
	// TODO(SNOW-1541938): use a default warehouse after fix on snowflake side
	warehouse, warehouseCleanup := acc.TestClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)

	networkRule, networkRuleCleanup := acc.TestClient().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)

	externalAccessIntegrationId, externalAccessIntegrationCleanup := acc.TestClient().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
	t.Cleanup(externalAccessIntegrationCleanup)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	directoryLocation := "abc"
	rootLocationWithCatalog := fmt.Sprintf("%s/%s", stage.Location(), directoryLocation)
	comment := random.Comment()
	mainFile := "foo"
	title := "foo"

	streamlitModelBasic := model.StreamlitWithIds("test", id, mainFile, stage.ID())
	streamlitModelComplete := model.StreamlitWithIds("test", id, mainFile, stage.ID()).
		WithComment(comment).
		WithTitle(title).
		WithDirectoryLocation(directoryLocation).
		WithQueryWarehouse(warehouse.ID().Name()).
		WithExternalAccessIntegrations(externalAccessIntegrationId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Streamlit),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				Config: accconfig.FromModels(t, streamlitModelBasic),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "database", id.DatabaseName()),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "schema", id.SchemaName()),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "stage", stage.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "main_file", "foo"),

					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttrSet(streamlitModelBasic.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "show_output.0.database_name", id.DatabaseName()),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "show_output.0.schema_name", id.SchemaName()),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "show_output.0.title", ""),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "show_output.0.owner", acc.TestClient().Context.CurrentRole(t).Name()),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "show_output.0.comment", ""),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "show_output.0.query_warehouse", ""),
					resource.TestCheckResourceAttrSet(streamlitModelBasic.ResourceReference(), "show_output.0.url_id"),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "show_output.0.owner_role_type", "ROLE"),

					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "describe_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "describe_output.0.title", ""),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "describe_output.0.root_location", stage.Location()),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "describe_output.0.main_file", "foo"),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "describe_output.0.query_warehouse", ""),
					resource.TestCheckResourceAttrSet(streamlitModelBasic.ResourceReference(), "describe_output.0.url_id"),
					resource.TestCheckResourceAttrSet(streamlitModelBasic.ResourceReference(), "describe_output.0.default_packages"),
					resource.TestCheckResourceAttrSet(streamlitModelBasic.ResourceReference(), "describe_output.0.user_packages.#"),
					resource.TestCheckResourceAttrSet(streamlitModelBasic.ResourceReference(), "describe_output.0.import_urls.#"),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "describe_output.0.external_access_integrations.#", "0"),
					resource.TestCheckResourceAttrSet(streamlitModelBasic.ResourceReference(), "describe_output.0.external_access_secrets"),
				),
			},
			// import - without optionals
			{
				Config:       accconfig.FromModels(t, streamlitModelBasic),
				ResourceName: streamlitModelBasic.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "database", id.DatabaseName()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "schema", id.SchemaName()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "stage", stage.ID().FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "main_file", "foo"),
				),
			},
			// set optionals
			{
				Config: accconfig.FromModels(t, streamlitModelComplete),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "database", id.DatabaseName()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "schema", id.SchemaName()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "main_file", "foo"),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "external_access_integrations.#", "1"),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "external_access_integrations.0", externalAccessIntegrationId.Name()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "title", title),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttrSet(streamlitModelComplete.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.0.database_name", id.DatabaseName()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.0.schema_name", id.SchemaName()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.0.title", title),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.0.owner", acc.TestClient().Context.CurrentRole(t).Name()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.0.query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttrSet(streamlitModelComplete.ResourceReference(), "show_output.0.url_id"),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.0.owner_role_type", "ROLE"),

					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "describe_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "describe_output.0.title", title),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "describe_output.0.root_location", rootLocationWithCatalog),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "describe_output.0.main_file", "foo"),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "describe_output.0.query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttrSet(streamlitModelComplete.ResourceReference(), "describe_output.0.url_id"),
					resource.TestCheckResourceAttrSet(streamlitModelComplete.ResourceReference(), "describe_output.0.default_packages"),
					resource.TestCheckResourceAttrSet(streamlitModelComplete.ResourceReference(), "describe_output.0.user_packages.#"),
					resource.TestCheckResourceAttrSet(streamlitModelComplete.ResourceReference(), "describe_output.0.import_urls.#"),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "describe_output.0.external_access_integrations.#", "1"),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "describe_output.0.external_access_integrations.0", externalAccessIntegrationId.Name()),
					resource.TestCheckResourceAttrSet(streamlitModelComplete.ResourceReference(), "describe_output.0.external_access_secrets")),
			},
			// import - complete
			{
				Config:       accconfig.FromModels(t, streamlitModelComplete),
				ResourceName: streamlitModelComplete.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "database", id.DatabaseName()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "schema", id.SchemaName()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "stage", stage.ID().FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "directory_location", directoryLocation),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "main_file", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "query_warehouse", warehouse.ID().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "external_access_integrations.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "external_access_integrations.0", externalAccessIntegrationId.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "title", title),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", comment)),
			},
			// change externally
			{
				PreConfig: func() {
					acc.TestClient().Streamlit.Update(t, sdk.NewAlterStreamlitRequest(id).WithSet(
						*sdk.NewStreamlitSetRequest().
							WithRootLocation(fmt.Sprintf("%s/bar", stage.Location())).
							WithComment("bar"),
					))
				},
				Config: accconfig.FromModels(t, streamlitModelComplete),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(streamlitModelComplete.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(streamlitModelComplete.ResourceReference(), "directory_location", sdk.String(directoryLocation), sdk.String("bar")),
						planchecks.ExpectChange(streamlitModelComplete.ResourceReference(), "directory_location", tfjson.ActionUpdate, sdk.String("bar"), sdk.String(directoryLocation)),
						planchecks.ExpectDrift(streamlitModelComplete.ResourceReference(), "comment", sdk.String(comment), sdk.String("bar")),
						planchecks.ExpectChange(streamlitModelComplete.ResourceReference(), "comment", tfjson.ActionUpdate, sdk.String("bar"), sdk.String(comment)),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "database", id.DatabaseName()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "schema", id.SchemaName()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "main_file", "foo"),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "external_access_integrations.#", "1"),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "external_access_integrations.0", externalAccessIntegrationId.Name()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "title", title),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttrSet(streamlitModelComplete.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.0.database_name", id.DatabaseName()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.0.schema_name", id.SchemaName()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.0.title", title),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.0.owner", acc.TestClient().Context.CurrentRole(t).Name()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.0.query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttrSet(streamlitModelComplete.ResourceReference(), "show_output.0.url_id"),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "show_output.0.owner_role_type", "ROLE"),

					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "describe_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "describe_output.0.title", title),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "describe_output.0.root_location", rootLocationWithCatalog),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "describe_output.0.main_file", "foo"),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "describe_output.0.query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttrSet(streamlitModelComplete.ResourceReference(), "describe_output.0.url_id"),
					resource.TestCheckResourceAttrSet(streamlitModelComplete.ResourceReference(), "describe_output.0.default_packages"),
					resource.TestCheckResourceAttrSet(streamlitModelComplete.ResourceReference(), "describe_output.0.user_packages.#"),
					resource.TestCheckResourceAttrSet(streamlitModelComplete.ResourceReference(), "describe_output.0.import_urls.#"),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "describe_output.0.external_access_integrations.#", "1"),
					resource.TestCheckResourceAttr(streamlitModelComplete.ResourceReference(), "describe_output.0.external_access_integrations.0", externalAccessIntegrationId.Name()),
					resource.TestCheckResourceAttrSet(streamlitModelComplete.ResourceReference(), "describe_output.0.external_access_secrets")),
			},
			// unset
			{
				Config: accconfig.FromModels(t, streamlitModelBasic),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "database", id.DatabaseName()),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "schema", id.SchemaName()),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "stage", stage.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "directory_location", ""),
					resource.TestCheckResourceAttr(streamlitModelBasic.ResourceReference(), "main_file", "foo")),
			},
		},
	})
}

func TestAcc_Streamlit_complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stage, stageCleanup := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	// warehouse is needed because default warehouse uses lowercase, and it fails in snowflake.
	// TODO(SNOW-1541938): use a default warehouse after fix on snowflake side
	warehouse, warehouseCleanup := acc.TestClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)

	networkRule, networkRuleCleanup := acc.TestClient().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)

	externalAccessIntegrationId, externalAccessIntegrationCleanup := acc.TestClient().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
	t.Cleanup(externalAccessIntegrationCleanup)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	directoryLocation := "abc"
	rootLocationWithCatalog := fmt.Sprintf("%s/%s", stage.Location(), directoryLocation)
	comment := random.Comment()
	mainFile := "foo"
	title := "foo"

	streamlitModelComplete := model.StreamlitWithIds("test", id, mainFile, stage.ID()).
		WithComment(comment).
		WithTitle(title).
		WithDirectoryLocation(directoryLocation).
		WithQueryWarehouse(warehouse.ID().Name()).
		WithExternalAccessIntegrations(externalAccessIntegrationId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Streamlit),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, streamlitModelComplete),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "database", id.DatabaseName()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "schema", id.SchemaName()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "stage", stage.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "directory_location", directoryLocation),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "main_file", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "external_access_integrations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "external_access_integrations.0", externalAccessIntegrationId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "title", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "fully_qualified_name", id.FullyQualifiedName()),

					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.database_name", id.DatabaseName()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.schema_name", id.SchemaName()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.title", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.owner", acc.TestClient().Context.CurrentRole(t).Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.comment", comment),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "show_output.0.url_id"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.owner_role_type", "ROLE"),

					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.title", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.root_location", rootLocationWithCatalog),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.main_file", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.url_id"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.default_packages"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.user_packages.#"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.import_urls.#"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.external_access_integrations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.external_access_integrations.0", externalAccessIntegrationId.Name()),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.external_access_secrets"),
				),
			},
			{
				Config:            accconfig.FromModels(t, streamlitModelComplete),
				ResourceName:      "snowflake_streamlit.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_Streamlit_Rename(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stage, stageCleanup := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	newId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	comment := random.Comment()
	newComment := random.Comment()
	mainFile := "foo"

	streamlitModel := model.StreamlitWithIds("test", id, mainFile, stage.ID()).
		WithComment(comment)
	streamlitModelRenamedAndUpdated := model.StreamlitWithIds("test", newId, mainFile, stage.ID()).
		WithComment(newComment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, streamlitModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(streamlitModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(streamlitModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(streamlitModel.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttr(streamlitModel.ResourceReference(), "fully_qualified_name", id.FullyQualifiedName()),
				),
			},
			{
				Config: accconfig.FromModels(t, streamlitModelRenamedAndUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(streamlitModelRenamedAndUpdated.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(streamlitModelRenamedAndUpdated.ResourceReference(), "name", newId.Name()),
					resource.TestCheckResourceAttr(streamlitModelRenamedAndUpdated.ResourceReference(), "show_output.0.name", newId.Name()),
					resource.TestCheckResourceAttr(streamlitModelRenamedAndUpdated.ResourceReference(), "show_output.0.comment", newComment),
					resource.TestCheckResourceAttr(streamlitModelRenamedAndUpdated.ResourceReference(), "fully_qualified_name", newId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_Streamlit_InvalidStage(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	streamlitModel := model.Streamlit("test", id.DatabaseId().FullyQualifiedName(), id.SchemaId().FullyQualifiedName(), id.Name(), "some", "some")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
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
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stage, stageCleanup := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	streamlitModel := model.StreamlitWithIds("test", id, "main_file", stage.ID())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Streamlit),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.94.1"),
				Config:            accconfig.FromModels(t, streamlitModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(streamlitModel.ResourceReference(), "id", helpers.EncodeSnowflakeID(id)),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, streamlitModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(streamlitModel.ResourceReference(), "id", id.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_Streamlit_IdentifierQuotingDiffSuppression(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stage, stageCleanup := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	quotedDatabaseName := fmt.Sprintf(`"%s"`, id.DatabaseName())
	quotedSchemaName := fmt.Sprintf(`"%s"`, id.SchemaName())
	quotedName := fmt.Sprintf(`"%s"`, id.Name())

	streamlitModel := model.Streamlit("test", quotedDatabaseName, quotedSchemaName, quotedName, "main_file", stage.ID().FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Streamlit),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  acc.ExternalProviderWithExactVersion("0.94.1"),
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
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
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
