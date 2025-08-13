//go:build !account_level_tests

package testacc

import (
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/assert"
)

// Recreation when promoting secondary to primary cannot be tested, because of the Terraform testing framework limitations.
// For the test that checks behavior for promoting secondary to primary, see `secondary_connection_promotion` manual test.

func TestAcc_SecondaryConnection_Basic(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed; also, different regions needed
	t.Skipf("Skipped due to 003813 (23001): The connection cannot be failed over to an account in the same region")

	// create primary connection
	connection, connectionCleanup := testClient().Connection.Create(t)
	t.Cleanup(connectionCleanup)

	secondaryAccountId := secondaryTestClient().Account.GetAccountIdentifier(t)
	testClient().Connection.Alter(t, sdk.NewAlterConnectionRequest(connection.ID()).
		WithEnableConnectionFailover(
			*sdk.NewEnableConnectionFailoverRequest([]sdk.AccountIdentifier{secondaryAccountId}),
		),
	)

	primaryConnectionAsExternalId := sdk.NewExternalObjectIdentifier(testClient().Account.GetAccountIdentifier(t), connection.ID())
	comment := random.Comment()

	provider := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)
	secondaryConnectionModel := model.SecondaryConnection("t", connection.ID().Name(), primaryConnectionAsExternalId.FullyQualifiedName())
	secondaryConnectionModelWithComment := model.SecondaryConnection("t", connection.ID().Name(), primaryConnectionAsExternalId.FullyQualifiedName()).
		WithComment(comment)

	assert.Eventually(t, func() bool {
		if _, err := secondaryTestClient().Connection.CreateReplication(t, connection.ID(), primaryConnectionAsExternalId); err == nil {
			secondaryTestClient().Connection.DropFunc(t, connection.ID())()
			return true
		}
		return false
	}, 10*time.Second, time.Second)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		// TODO: Check destroy is client dependent (which means we should be able to use dedicated client for checking: e.g. secondary in this case)
		//   Or just reverse it (create on secondary and use primary for resources)
		// CheckDestroy: CheckDestroy(t, resources.SecondaryConnection),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, provider, secondaryConnectionModel),
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
							HasAccountName(secondaryAccountId.AccountName()).
							HasOrganizationName(secondaryAccountId.OrganizationName()).
							HasComment("").
							HasIsPrimary(false).
							HasPrimaryIdentifier(primaryConnectionAsExternalId).
							HasFailoverAllowedToAccounts().
							HasConnectionUrl(testClient().Connection.GetConnectionUrl(secondaryAccountId.OrganizationName(), connection.ID().Name())),
					),
				),
			},
			// set comment
			{
				Config: config.FromModels(t, provider, secondaryConnectionModelWithComment),
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
				Config: config.FromModels(t, provider, secondaryConnectionModel),
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
					secondaryTestClient().Connection.Alter(t, sdk.NewAlterConnectionRequest(connection.ID()).WithSet(*sdk.NewConnectionSetRequest().WithComment(comment)))
				},
				Config: config.FromModels(t, provider, secondaryConnectionModel),
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
