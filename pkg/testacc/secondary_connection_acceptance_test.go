//go:build !account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SecondaryConnection_Basic(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed; also, different regions needed
	t.Skipf("Skipped due to 003813 (23001): The connection cannot be failed over to an account in the same region")

	// create primary connection
	connection, connectionCleanup := secondaryTestClient().Connection.Create(t)
	t.Cleanup(connectionCleanup)

	accountId := testClient().Account.GetAccountIdentifier(t)
	secondaryTestClient().Connection.Alter(t, sdk.NewAlterConnectionRequest(connection.ID()).WithEnableConnectionFailover(*sdk.NewEnableConnectionFailoverRequest([]sdk.AccountIdentifier{accountId})))

	primaryConnectionAsExternalId := sdk.NewExternalObjectIdentifier(accountId, connection.ID())
	comment := random.Comment()

	secondaryConnectionModel := model.SecondaryConnection("t", connection.ID().Name(), primaryConnectionAsExternalId.FullyQualifiedName())
	secondaryConnectionModelWithComment := model.SecondaryConnection("t", connection.ID().Name(), primaryConnectionAsExternalId.FullyQualifiedName()).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SecondaryConnection),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, secondaryConnectionModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecondaryConnectionResource(t, secondaryConnectionModel.ResourceReference()).
							HasNameString(connection.ID().Name()).
							HasFullyQualifiedNameString(connection.ID().FullyQualifiedName()).
							HasAsReplicaOfIdentifier(primaryConnectionAsExternalId).
							HasIsPrimaryString("false").
							HasCommentString(""),

						resourceshowoutputassert.ConnectionShowOutput(t, secondaryConnectionModel.ResourceReference()).
							HasName(connection.ID().Name()).
							HasSnowflakeRegion(secondaryTestClient().Context.CurrentRegion(t)).
							HasAccountLocator(secondaryTestClient().GetAccountLocator()).
							HasAccountName(accountId.AccountName()).
							HasOrganizationName(accountId.OrganizationName()).
							HasComment("").
							HasIsPrimary(false).
							HasPrimaryIdentifier(primaryConnectionAsExternalId).
							HasFailoverAllowedToAccounts(accountId).
							HasConnectionUrl(
								secondaryTestClient().Connection.GetConnectionUrl(accountId.OrganizationName(), connection.ID().Name()),
							),
					),
				),
			},
			// set comment
			{
				Config: config.FromModels(t, secondaryConnectionModelWithComment),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecondaryConnectionResource(t, secondaryConnectionModelWithComment.ResourceReference()).
							HasNameString(connection.ID().Name()).
							HasFullyQualifiedNameString(connection.ID().FullyQualifiedName()).
							HasCommentString(comment),

						resourceshowoutputassert.ConnectionShowOutput(t, secondaryConnectionModelWithComment.ResourceReference()).
							HasComment(comment),
					),
				),
			},
			// import
			{
				ResourceName:      secondaryConnectionModelWithComment.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(connection.ID()), "name", connection.ID().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(connection.ID()), "comment", comment),
				),
			},
			// unset comment
			{
				Config: config.FromModels(t, secondaryConnectionModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecondaryConnectionResource(t, secondaryConnectionModel.ResourceReference()).
							HasCommentString(""),

						resourceshowoutputassert.ConnectionShowOutput(t, secondaryConnectionModel.ResourceReference()).
							HasComment(""),
					),
				),
			},
			// recreate when externally promoted to primary
			{
				PreConfig: func() {
					testClient().Connection.Alter(t, sdk.NewAlterConnectionRequest(connection.ID()).WithPrimary(true))
				},
				Config: config.FromModels(t, secondaryConnectionModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secondaryConnectionModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecondaryConnectionResource(t, secondaryConnectionModel.ResourceReference()).
							HasIsPrimaryString("false"),

						resourceshowoutputassert.ConnectionShowOutput(t, secondaryConnectionModel.ResourceReference()).
							HasIsPrimary(false),
					),
				),
			},
		},
	})
}
