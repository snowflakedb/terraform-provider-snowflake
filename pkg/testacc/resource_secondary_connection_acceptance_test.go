//go:build non_account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// Recreation when promoting secondary to primary cannot be tested, because of the Terraform testing framework limitations.
// For the test that checks behavior for promoting secondary to primary, see `secondary_connection_promotion` manual test.

func TestAcc_SecondaryConnection_Basic(t *testing.T) {
	testenvs.SkipTestIfValueIn(t, testenvs.SnowflakeTestingEnvironment, []string{
		string(testenvs.SnowflakeProdEnvironment),
		string(testenvs.SnowflakePreProdGovEnvironment),
	}, "SNOW-3198924: Missing azure configuration on all testing environments")

	// create primary connection
	connection, connectionCleanup := azureTestClient().Connection.Create(t)
	t.Cleanup(connectionCleanup)

	accountId := testClient().Account.GetAccountIdentifier(t)
	azureTestClient().Connection.Alter(t, sdk.NewAlterConnectionRequest(connection.ID()).
		WithEnableConnectionFailover(
			*sdk.NewEnableConnectionFailoverRequest([]sdk.AccountIdentifier{accountId}),
		),
	)

	primaryConnectionAsExternalId := sdk.NewExternalObjectIdentifier(azureTestClient().Account.GetAccountIdentifier(t), connection.ID())
	comment := random.Comment()

	secondaryConnectionModel := model.SecondaryConnection("t", connection.ID().Name(), primaryConnectionAsExternalId.FullyQualifiedName())
	secondaryConnectionModelWithComment := model.SecondaryConnection("t", connection.ID().Name(), primaryConnectionAsExternalId.FullyQualifiedName()).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
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
							HasSnowflakeRegion(testClient().Context.CurrentRegion(t)).
							HasAccountLocator(testClient().GetAccountLocator()).
							HasAccountName(accountId.AccountName()).
							HasOrganizationName(accountId.OrganizationName()).
							HasComment("").
							HasIsPrimary(false).
							HasPrimaryIdentifier(primaryConnectionAsExternalId).
							HasFailoverAllowedToAccounts().
							HasConnectionUrl(azureTestClient().Connection.GetConnectionUrl(accountId.OrganizationName(), connection.ID().Name())),
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
			{
				PreConfig: func() {
					testClient().Connection.Alter(t, sdk.NewAlterConnectionRequest(connection.ID()).WithSet(*sdk.NewConnectionSetRequest().WithComment(comment)))
				},
				Config: config.FromModels(t, secondaryConnectionModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secondaryConnectionModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.SecondaryConnectionResource(t, secondaryConnectionModel.ResourceReference()).
							HasCommentString(""),
						resourceshowoutputassert.ConnectionShowOutput(t, secondaryConnectionModel.ResourceReference()).
							HasComment(""),
					),
				),
			},
		},
	})
}
