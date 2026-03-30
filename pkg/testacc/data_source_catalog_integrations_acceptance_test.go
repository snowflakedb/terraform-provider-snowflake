//go:build non_account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CatalogIntegrations_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := testClient().Ids.RandomAccountObjectIdentifier()

	model1 := model.CatalogIntegrationObjectStorage("test1", idOne.Name(), false, string(sdk.CatalogIntegrationTableFormatIceberg))
	model2 := model.CatalogIntegrationObjectStorage("test2", idTwo.Name(), false, string(sdk.CatalogIntegrationTableFormatDelta))
	model3 := model.CatalogIntegrationObjectStorage("test3", idThree.Name(), false, string(sdk.CatalogIntegrationTableFormatIceberg))

	catalogIntegrationsModelLikeFirst := datasourcemodel.CatalogIntegrations("test").
		WithWithDescribe(false).
		WithLike(idOne.Name()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	catalogIntegrationsModelLikePrefix := datasourcemodel.CatalogIntegrations("test").
		WithWithDescribe(false).
		WithLike(prefix+"%").
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationObjectStorage),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, catalogIntegrationsModelLikeFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(catalogIntegrationsModelLikeFirst.DatasourceReference(), "catalog_integrations.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, catalogIntegrationsModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(catalogIntegrationsModelLikePrefix.DatasourceReference(), "catalog_integrations.#", "2"),
				),
			},
		},
	})
}

func TestAcc_CatalogIntegrations_CompleteUseCase(t *testing.T) {
	prefix := random.AlphaN(4)
	comment := random.Comment()
	refreshIntervalSeconds := random.IntRange(30, 86400)

	glueAwsRoleArn := "arn:aws:iam::123456789012:role/sqsAccess"
	glueCatalogId := random.NumericN(15)
	glueId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix + "glue")

	objectStorageId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix + "obj")

	catalogUri := "https://testorg-testacc.snowflakecomputing.com/polaris/api/catalog"
	catalogName := random.AlphanumericN(15)
	oAuthClientId := random.AlphanumericN(15)
	oAuthClientSecret := random.AlphanumericN(15)
	oAuthAllowedScope := "PRINCIPAL_ROLE:ALL"
	basicRestAuth := []sdk.OAuthRestAuthenticationRequest{
		*sdk.NewOAuthRestAuthenticationRequest(oAuthClientId, oAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}}),
	}
	basicRestConfig := []sdk.OpenCatalogRestConfigRequest{
		*sdk.NewOpenCatalogRestConfigRequest(catalogUri, catalogName),
	}
	openCatalogId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix + "oc")

	icebergCatalogUri := "https://api.tabular.io/ws"
	icebergRestConfig := *sdk.NewIcebergRestRestConfigRequest(icebergCatalogUri).
		WithCatalogApiType(sdk.CatalogIntegrationCatalogApiTypeAwsApiGateway)
	sigV4IamRole := "arn:aws:iam::123456789012:role/sigv4-role"
	icebergId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix + "ice")

	catalogIntegrationAwsGlue := model.CatalogIntegrationAwsGlue("w", glueId.Name(), false, glueAwsRoleArn, glueCatalogId)
	catalogIntegrationObjectStorage := model.CatalogIntegrationObjectStorage("w", objectStorageId.Name(), true, string(sdk.CatalogIntegrationTableFormatIceberg))
	catalogIntegrationOpenCatalog := model.CatalogIntegrationOpenCatalog("w", openCatalogId.Name(), false, basicRestAuth, basicRestConfig).
		WithComment(comment)
	catalogIntegrationIcebergRestBearer := model.CatalogIntegrationIcebergRestSigV4("w", icebergId.Name(), false, icebergRestConfig, *sdk.NewSigV4RestAuthenticationRequest(sigV4IamRole)).
		WithRefreshIntervalSeconds(refreshIntervalSeconds)

	glueNoDescribe := datasourcemodel.CatalogIntegrations("test").
		WithLike(glueId.Name()).
		WithWithDescribe(false).
		WithDependsOn(catalogIntegrationAwsGlue.ResourceReference())

	glueWithDescribe := datasourcemodel.CatalogIntegrations("test").
		WithLike(glueId.Name()).
		WithWithDescribe(true).
		WithDependsOn(catalogIntegrationAwsGlue.ResourceReference())

	objectStorageWithDescribe := datasourcemodel.CatalogIntegrations("test").
		WithLike(objectStorageId.Name()).
		WithWithDescribe(true).
		WithDependsOn(catalogIntegrationObjectStorage.ResourceReference())

	openCatalogWithDescribe := datasourcemodel.CatalogIntegrations("test").
		WithLike(openCatalogId.Name()).
		WithWithDescribe(true).
		WithDependsOn(catalogIntegrationOpenCatalog.ResourceReference())

	icebergBearerWithDescribe := datasourcemodel.CatalogIntegrations("test").
		WithLike(icebergId.Name()).
		WithWithDescribe(true).
		WithDependsOn(catalogIntegrationIcebergRestBearer.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: ComposeCheckDestroy(t,
			resources.CatalogIntegrationAwsGlue,
			resources.CatalogIntegrationObjectStorage,
			resources.CatalogIntegrationOpenCatalog,
			resources.CatalogIntegrationIcebergRest,
		),
		Steps: []resource.TestStep{
			// AWS Glue without describe
			{
				Config: accconfig.FromModels(t, catalogIntegrationAwsGlue, glueNoDescribe),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(glueNoDescribe.DatasourceReference(), "catalog_integrations.#", "1")),
					resourceshowoutputassert.CatalogIntegrationsDatasourceShowOutput(t, glueNoDescribe.DatasourceReference(), 0).
						HasName(glueId.Name()).
						HasType("CATALOG").
						HasCategory("CATALOG").
						HasEnabled(false).
						HasComment(""),
					assert.Check(resource.TestCheckResourceAttr(glueNoDescribe.DatasourceReference(), "catalog_integrations.0.describe_output.#", "0")),
				),
			},
			// AWS Glue with describe
			{
				Config: accconfig.FromModels(t, catalogIntegrationAwsGlue, glueWithDescribe),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(glueWithDescribe.DatasourceReference(), "catalog_integrations.#", "1")),
					resourceshowoutputassert.CatalogIntegrationsDatasourceShowOutput(t, glueWithDescribe.DatasourceReference(), 0).
						HasName(glueId.Name()).
						HasType("CATALOG").
						HasCategory("CATALOG").
						HasEnabled(false).
						HasComment(""),
					resourceshowoutputassert.CatalogIntegrationsDatasourceAwsGlueDescribeOutput(t, glueWithDescribe.DatasourceReference()).
						HasId(glueId).
						HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeAWSGlue).
						HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
						HasEnabled(false).
						HasRefreshIntervalSeconds(30).
						HasComment("").
						HasGlueAwsRoleArn(glueAwsRoleArn).
						HasGlueCatalogId(glueCatalogId).
						HasCatalogNamespace(""),
				),
			},
			// Object Storage with describe
			{
				Config: accconfig.FromModels(t, catalogIntegrationObjectStorage, objectStorageWithDescribe),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(objectStorageWithDescribe.DatasourceReference(), "catalog_integrations.#", "1")),
					resourceshowoutputassert.CatalogIntegrationsDatasourceShowOutput(t, objectStorageWithDescribe.DatasourceReference(), 0).
						HasName(objectStorageId.Name()).
						HasType("CATALOG").
						HasCategory("CATALOG").
						HasEnabled(true).
						HasComment(""),
					resourceshowoutputassert.CatalogIntegrationsDatasourceObjectStorageDescribeOutput(t, objectStorageWithDescribe.DatasourceReference()).
						HasId(objectStorageId).
						HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeObjectStorage).
						HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
						HasEnabled(true).
						HasRefreshIntervalSeconds(30).
						HasComment(""),
				),
			},
			// Open Catalog with describe
			{
				Config: accconfig.FromModels(t, catalogIntegrationOpenCatalog, openCatalogWithDescribe),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(openCatalogWithDescribe.DatasourceReference(), "catalog_integrations.#", "1")),
					resourceshowoutputassert.CatalogIntegrationsDatasourceShowOutput(t, openCatalogWithDescribe.DatasourceReference(), 0).
						HasName(openCatalogId.Name()).
						HasType("CATALOG").
						HasCategory("CATALOG").
						HasEnabled(false).
						HasComment(comment),
					resourceshowoutputassert.CatalogIntegrationsDatasourceOpenCatalogDescribeOutput(t, openCatalogWithDescribe.DatasourceReference()).
						HasId(openCatalogId).
						HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypePolaris).
						HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
						HasEnabled(false).
						HasRefreshIntervalSeconds(30).
						HasComment(comment).
						HasCatalogNamespace(""),
					resourceshowoutputassert.OpenCatalogRestConfigDatasourceDescribeOutput(t, openCatalogWithDescribe.DatasourceReference()).
						HasCatalogUri(catalogUri).
						HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
						HasCatalogName(catalogName).
						HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials),
					resourceshowoutputassert.OAuthRestAuthenticationDatasourceDescribeOutput(t, openCatalogWithDescribe.DatasourceReference(), "oauth_rest_authentication").
						HasOauthTokenUri(catalogUri+"/v1/oauth/tokens").
						HasOauthClientId(oAuthClientId).
						HasOauthAllowedScopes(oAuthAllowedScope),
				),
			},
			// Iceberg REST Bearer with describe
			{
				Config: accconfig.FromModels(t, catalogIntegrationIcebergRestBearer, icebergBearerWithDescribe),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(icebergBearerWithDescribe.DatasourceReference(), "catalog_integrations.#", "1")),
					resourceshowoutputassert.CatalogIntegrationsDatasourceShowOutput(t, icebergBearerWithDescribe.DatasourceReference(), 0).
						HasName(icebergId.Name()).
						HasType("CATALOG").
						HasCategory("CATALOG").
						HasEnabled(false).
						HasComment(""),
					resourceshowoutputassert.CatalogIntegrationsDatasourceIcebergRestDescribeOutput(t, icebergBearerWithDescribe.DatasourceReference()).
						HasId(icebergId).
						HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergREST).
						HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
						HasEnabled(false).
						HasRefreshIntervalSeconds(refreshIntervalSeconds).
						HasComment("").
						HasCatalogNamespace(""),
					resourceshowoutputassert.IcebergRestRestConfigDatasourceDescribeOutput(t, icebergBearerWithDescribe.DatasourceReference()).
						HasCatalogUri(icebergCatalogUri).
						HasPrefix("").
						HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypeAwsApiGateway).
						HasCatalogName("").
						HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials),
					resourceshowoutputassert.SigV4RestAuthenticationDatasourceDescribeOutput(t, icebergBearerWithDescribe.DatasourceReference()).
						// Don't check sigv4_signing_region, as its default value depends on the current region name
						HasSigv4IamRole(sigV4IamRole),
				),
			},
		},
	})
}

func TestAcc_CatalogIntegrations_MultipleTypes(t *testing.T) {
	glueAwsRoleArn := "arn:aws:iam::123456789012:role/sqsAccess"
	glueCatalogId := random.NumericN(15)

	prefix := random.AlphaN(4)
	idOne := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix + "1")
	idTwo := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix + "2")

	catalogIntegrationAwsGlueModel := model.CatalogIntegrationAwsGlue("w", idOne.Name(), false, glueAwsRoleArn, glueCatalogId)
	catalogIntegrationObjectStorageModel := model.CatalogIntegrationObjectStorage("w", idTwo.Name(), false, string(sdk.CatalogIntegrationTableFormatIceberg))

	catalogIntegrationsModel := datasourcemodel.CatalogIntegrations("test").
		WithLike(prefix+"%").
		WithDependsOn(catalogIntegrationAwsGlueModel.ResourceReference(), catalogIntegrationObjectStorageModel.ResourceReference())

	ref := catalogIntegrationsModel.DatasourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, catalogIntegrationAwsGlueModel, catalogIntegrationObjectStorageModel, catalogIntegrationsModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(ref, "catalog_integrations.#", "2")),
					resourceshowoutputassert.CatalogIntegrationsDatasourceShowOutput(t, ref, 0).
						// Don't check name, as the order of elements in SHOW output is unpredictable
						HasType("CATALOG").
						HasCategory("CATALOG").
						HasEnabled(false).
						HasComment(""),
					resourceshowoutputassert.CatalogIntegrationsDatasourceShowOutput(t, ref, 1).
						// Don't check name, as the order of elements in SHOW output is unpredictable
						HasType("CATALOG").
						HasCategory("CATALOG").
						HasEnabled(false).
						HasComment(""),
					assert.Check(resource.TestCheckResourceAttr(ref, "catalog_integrations.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(ref, "catalog_integrations.0.describe_output.0.enabled", "false")),
					assert.Check(resource.TestCheckResourceAttr(ref, "catalog_integrations.0.describe_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr(ref, "catalog_integrations.1.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(ref, "catalog_integrations.1.describe_output.0.enabled", "false")),
					assert.Check(resource.TestCheckResourceAttr(ref, "catalog_integrations.1.describe_output.0.comment", "")),
				),
			},
		},
	})
}
