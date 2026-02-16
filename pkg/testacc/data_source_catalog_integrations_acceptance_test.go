//go:build non_account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CatalogIntegrations_BasicFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := testClient().Ids.RandomAccountObjectIdentifier()

	catalogModel1 := model.CatalogIntegration("test1", idOne.Name(), "OBJECT_STORE", "ICEBERG", true)
	catalogModel2 := model.CatalogIntegration("test2", idTwo.Name(), "OBJECT_STORE", "ICEBERG", true)
	catalogModel3 := model.CatalogIntegration("test3", idThree.Name(), "OBJECT_STORE", "ICEBERG", true)

	catalogIntegrationsModelLikeFirst := datasourcemodel.CatalogIntegrations("test").
		WithWithDescribe(false).
		WithLike(idOne.Name()).
		WithDependsOn(catalogModel1.ResourceReference(), catalogModel2.ResourceReference(), catalogModel3.ResourceReference())

	catalogIntegrationsModelLikePrefix := datasourcemodel.CatalogIntegrations("test").
		WithWithDescribe(false).
		WithLike(prefix+"%").
		WithDependsOn(catalogModel1.ResourceReference(), catalogModel2.ResourceReference(), catalogModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegration),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, catalogModel1, catalogModel2, catalogModel3, catalogIntegrationsModelLikeFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(catalogIntegrationsModelLikeFirst.DatasourceReference(), "catalog_integrations.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, catalogModel1, catalogModel2, catalogModel3, catalogIntegrationsModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(catalogIntegrationsModelLikePrefix.DatasourceReference(), "catalog_integrations.#", "2"),
				),
			},
		},
	})
}

func TestAcc_CatalogIntegrations_WithDescribe(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	catalogModel := model.CatalogIntegration("w", id.Name(), "OBJECT_STORE", "ICEBERG", true).
		WithComment(comment)

	catalogIntegrationsNoDescribe := datasourcemodel.CatalogIntegrations("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithDependsOn(catalogModel.ResourceReference())

	catalogIntegrationsWithDescribe := datasourcemodel.CatalogIntegrations("test").
		WithLike(id.Name()).
		WithWithDescribe(true).
		WithDependsOn(catalogModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegration),
		Steps: []resource.TestStep{
			// Without describe
			{
				Config: accconfig.FromModels(t, catalogModel, catalogIntegrationsNoDescribe),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(catalogIntegrationsNoDescribe.DatasourceReference(), "catalog_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(catalogIntegrationsNoDescribe.DatasourceReference(), "catalog_integrations.0.describe_output.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(catalogIntegrationsNoDescribe.DatasourceReference(), "catalog_integrations.0.show_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(catalogIntegrationsNoDescribe.DatasourceReference(), "catalog_integrations.0.show_output.0.enabled", "true")),
					assert.Check(resource.TestCheckResourceAttr(catalogIntegrationsNoDescribe.DatasourceReference(), "catalog_integrations.0.show_output.0.comment", comment)),
				),
			},
			// With describe
			{
				Config: accconfig.FromModels(t, catalogModel, catalogIntegrationsWithDescribe),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(catalogIntegrationsWithDescribe.DatasourceReference(), "catalog_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(catalogIntegrationsWithDescribe.DatasourceReference(), "catalog_integrations.0.show_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(catalogIntegrationsWithDescribe.DatasourceReference(), "catalog_integrations.0.show_output.0.enabled", "true")),
					assert.Check(resource.TestCheckResourceAttr(catalogIntegrationsWithDescribe.DatasourceReference(), "catalog_integrations.0.show_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttrSet(catalogIntegrationsWithDescribe.DatasourceReference(), "catalog_integrations.0.describe_output.#")),
				),
			},
		},
	})
}

func TestAcc_CatalogIntegrations_MultipleTypes(t *testing.T) {
	prefix := random.AlphaN(4)
	idObjectStore := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix + "1")
	idDelta := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix + "2")

	commentObjectStore := random.Comment()
	commentDelta := random.Comment()

	catalogModelObjectStore := model.CatalogIntegration("w1", idObjectStore.Name(), "OBJECT_STORE", "ICEBERG", true).
		WithComment(commentObjectStore)

	catalogModelDelta := model.CatalogIntegration("w2", idDelta.Name(), "OBJECT_STORE", "DELTA", false).
		WithComment(commentDelta)

	catalogIntegrationsModel := datasourcemodel.CatalogIntegrations("test").
		WithLike(prefix+"%").
		WithDependsOn(catalogModelObjectStore.ResourceReference(), catalogModelDelta.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegration),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, catalogModelObjectStore, catalogModelDelta, catalogIntegrationsModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(catalogIntegrationsModel.DatasourceReference(), "catalog_integrations.#", "2"),

					resource.TestCheckResourceAttr(catalogIntegrationsModel.DatasourceReference(), "catalog_integrations.0.show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr(catalogIntegrationsModel.DatasourceReference(), "catalog_integrations.0.show_output.0.comment", commentObjectStore),
					resource.TestCheckResourceAttrSet(catalogIntegrationsModel.DatasourceReference(), "catalog_integrations.0.describe_output.#"),

					resource.TestCheckResourceAttr(catalogIntegrationsModel.DatasourceReference(), "catalog_integrations.1.show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr(catalogIntegrationsModel.DatasourceReference(), "catalog_integrations.1.show_output.0.comment", commentDelta),
					resource.TestCheckResourceAttrSet(catalogIntegrationsModel.DatasourceReference(), "catalog_integrations.1.describe_output.#"),
				),
			},
		},
	})
}

func TestAcc_CatalogIntegrations_DefaultWithDescribe(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	catalogModel := model.CatalogIntegration("w", id.Name(), "OBJECT_STORE", "ICEBERG", true)

	// Test that with_describe defaults to true
	catalogIntegrationsModel := datasourcemodel.CatalogIntegrations("test").
		WithLike(id.Name()).
		WithDependsOn(catalogModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegration),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, catalogModel, catalogIntegrationsModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(catalogIntegrationsModel.DatasourceReference(), "catalog_integrations.#", "1"),
					// Should have describe output since default is true
					resource.TestCheckResourceAttrSet(catalogIntegrationsModel.DatasourceReference(), "catalog_integrations.0.describe_output.#"),
				),
			},
		},
	})
}
