//go:build !account_level_tests

package testacc

import (
	"fmt"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
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

func TestAcc_PrimaryConnection_Basic(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	accountId := testClient().Account.GetAccountIdentifier(t)
	secondaryAccountId := secondaryTestClient().Account.GetAccountIdentifier(t)
	primaryConnectionAsExternalId := sdk.NewExternalObjectIdentifier(accountId, id)

	connectionModel := model.PrimaryConnection("t", id.Name())
	connectionModelWithComment := model.PrimaryConnection("t", id.Name()).WithComment(comment)
	connectionModelWithFailover := model.PrimaryConnection("t", id.Name()).WithEnableFailover(secondaryAccountId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PrimaryConnection),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModels(t, connectionModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.PrimaryConnectionResource(t, connectionModel.ResourceReference()).
							HasNameString(id.Name()).
							HasFullyQualifiedNameString(id.FullyQualifiedName()).
							HasEnableFailoverToAccountsEmpty().
							HasCommentString(""),

						resourceshowoutputassert.ConnectionShowOutput(t, connectionModel.ResourceReference()).
							HasName(id.Name()).
							HasSnowflakeRegion(testClient().Context.CurrentRegion(t)).
							HasAccountLocator(testClient().GetAccountLocator()).
							HasAccountName(accountId.AccountName()).
							HasOrganizationName(accountId.OrganizationName()).
							HasComment("").
							HasIsPrimary(true).
							HasPrimaryIdentifier(primaryConnectionAsExternalId).
							HasFailoverAllowedToAccounts(accountId).
							HasConnectionUrl(
								testClient().Connection.GetConnectionUrl(accountId.OrganizationName(), id.Name()),
							),
					),
				),
			},
			// set comment
			{
				Config: config.FromModels(t, connectionModelWithComment),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.PrimaryConnectionResource(t, connectionModelWithComment.ResourceReference()).
							HasNameString(id.Name()).
							HasFullyQualifiedNameString(id.FullyQualifiedName()).
							HasEnableFailoverToAccountsEmpty().
							HasCommentString(comment),

						resourceshowoutputassert.ConnectionShowOutput(t, connectionModelWithComment.ResourceReference()).
							HasComment(comment),
					),
				),
			},
			// import
			{
				ResourceName:      connectionModelWithComment.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// unset comment
			{
				Config: config.FromModels(t, connectionModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.PrimaryConnectionResource(t, connectionModel.ResourceReference()).
							HasCommentString(""),

						resourceshowoutputassert.ConnectionShowOutput(t, connectionModel.ResourceReference()).
							HasComment(""),
					),
				),
			},
			// enable failover to second account
			{
				Config: config.FromModels(t, connectionModelWithFailover),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.PrimaryConnectionResource(t, connectionModelWithFailover.ResourceReference()).
							HasNameString(id.Name()).
							HasFullyQualifiedNameString(id.FullyQualifiedName()).
							HasExactlyFailoverToAccountsInOrder(secondaryAccountId).
							HasCommentString(""),

						resourceshowoutputassert.ConnectionShowOutput(t, connectionModelWithFailover.ResourceReference()).
							HasFailoverAllowedToAccounts(accountId, secondaryAccountId),
					),
				),
			},
			// disable failover to second account
			{
				Config: config.FromModels(t, connectionModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.PrimaryConnectionResource(t, connectionModel.ResourceReference()).
							HasNameString(id.Name()).
							HasFullyQualifiedNameString(id.FullyQualifiedName()).
							HasEnableFailoverToAccountsEmpty().
							HasCommentString(""),

						resourceshowoutputassert.ConnectionShowOutput(t, connectionModel.ResourceReference()).
							HasFailoverAllowedToAccounts(accountId),
					),
				),
			},
		},
	})
}

func TestAcc_PrimaryConnection_ExternalChanges(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	accountId := testClient().Account.GetAccountIdentifier(t)
	secondaryAccountId := secondaryTestClient().Account.GetAccountIdentifier(t)
	primaryConnectionAsExternalId := sdk.NewExternalObjectIdentifier(accountId, id)

	connectionModel := model.PrimaryConnection("t", id.Name()).WithComment("config comment")
	connectionModelWithFailover := model.PrimaryConnection("t", id.Name()).WithEnableFailover(secondaryAccountId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PrimaryConnection),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModels(t, connectionModel),
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.PrimaryConnectionResource(t, connectionModel.ResourceReference()).
							HasNameString(id.Name()).
							HasFullyQualifiedNameString(id.FullyQualifiedName()).
							HasEnableFailoverToAccountsEmpty().
							HasCommentString("config comment"),

						resourceshowoutputassert.ConnectionShowOutput(t, connectionModel.ResourceReference()).
							HasName(id.Name()).
							HasSnowflakeRegion(testClient().Context.CurrentRegion(t)).
							HasAccountLocator(testClient().GetAccountLocator()).
							HasAccountName(accountId.AccountName()).
							HasOrganizationName(accountId.OrganizationName()).
							HasComment("config comment").
							HasIsPrimary(true).
							HasPrimaryIdentifier(primaryConnectionAsExternalId).
							HasFailoverAllowedToAccounts(accountId),
					),
				),
			},
			// change comment externally
			{
				PreConfig: func() {
					testClient().Connection.Alter(t, sdk.NewAlterConnectionRequest(id).
						WithSet(*sdk.NewConnectionSetRequest().
							WithComment("external comment")))
				},
				Config: config.FromModels(t, connectionModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(connectionModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(connectionModel.ResourceReference(), "comment", tfjson.ActionUpdate, sdk.String("external comment"), sdk.String("config comment")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.PrimaryConnectionResource(t, connectionModel.ResourceReference()).
							HasCommentString("config comment"),
						resourceshowoutputassert.ConnectionShowOutput(t, connectionModel.ResourceReference()).
							HasComment("config comment"),
					),
				),
			},
			// externally change enable failover accounts list - add second account
			{
				PreConfig: func() {
					testClient().Connection.Alter(t, sdk.NewAlterConnectionRequest(id).
						WithEnableConnectionFailover(*sdk.NewEnableConnectionFailoverRequest(
							[]sdk.AccountIdentifier{secondaryAccountId})),
					)
				},
				Config: config.FromModels(t, connectionModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(connectionModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(
							connectionModel.ResourceReference(),
							"enable_failover_to_accounts",
							tfjson.ActionUpdate,
							sdk.String(fmt.Sprintf("[%s.%s]", secondaryAccountId.OrganizationName(), secondaryAccountId.AccountName())),
							sdk.String("[]"),
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.PrimaryConnectionResource(t, connectionModel.ResourceReference()).
							HasEnableFailoverToAccountsEmpty(),
						resourceshowoutputassert.ConnectionShowOutput(t, connectionModel.ResourceReference()).
							HasFailoverAllowedToAccounts(accountId),
					),
				),
			},
			// externally change disable failover
			{
				PreConfig: func() {
					testClient().Connection.Alter(t, sdk.NewAlterConnectionRequest(id).WithDisableConnectionFailover(*sdk.NewDisableConnectionFailoverRequest()))
				},
				Config: config.FromModels(t, connectionModelWithFailover),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(connectionModelWithFailover.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(
							connectionModel.ResourceReference(),
							"enable_failover_to_accounts",
							tfjson.ActionUpdate,
							sdk.String("[]"),
							sdk.String(fmt.Sprintf("[%s.%s]", secondaryAccountId.OrganizationName(), secondaryAccountId.AccountName())),
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assertThat(t,
						resourceassert.PrimaryConnectionResource(t, connectionModelWithFailover.ResourceReference()).
							HasExactlyFailoverToAccountsInOrder(secondaryAccountId),
						resourceshowoutputassert.ConnectionShowOutput(t, connectionModelWithFailover.ResourceReference()).
							HasFailoverAllowedToAccounts(accountId, secondaryAccountId),
					),
				),
			},
		},
	})
}
