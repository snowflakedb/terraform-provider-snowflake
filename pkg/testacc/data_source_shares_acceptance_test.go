//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/ids"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Shares_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(6)
	id1 := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomAccountObjectIdentifier()

	share1 := model.Share("t1", id1.Name())
	share2 := model.Share("t2", id2.Name())
	share3 := model.Share("t3", id3.Name())
	dependsOn := []string{share1.ResourceReference(), share2.ResourceReference(), share3.ResourceReference()}

	dsLikeExact := datasourcemodel.Shares("test").
		WithLike(id1.Name()).
		WithDependsOn(dependsOn...)
	dsLikePrefix := datasourcemodel.Shares("test").
		WithLike(prefix + "%").
		WithDependsOn(dependsOn...)
	dsStartsWith := datasourcemodel.Shares("test").
		WithStartsWith(prefix).
		WithDependsOn(dependsOn...)
	dsLimit := datasourcemodel.Shares("test").
		WithLike(prefix + "%").
		WithLimit(1).
		WithDependsOn(dependsOn...)
	dsNoMatch := datasourcemodel.Shares("test").
		WithLike(prefix + "_does_not_exist").
		WithDependsOn(dependsOn...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// like only (exact match)
			{
				Config: accconfig.FromModels(t, share1, share2, share3, dsLikeExact),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsLikeExact.DatasourceReference(), "shares.#", "1")),
					resourceshowoutputassert.SharesDatasourceShowOutput(t, dsLikeExact.DatasourceReference()).HasName(id1.Name()),
				),
			},
			// like only (prefix pattern)
			{
				Config: accconfig.FromModels(t, share1, share2, share3, dsLikePrefix),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsLikePrefix.DatasourceReference(), "shares.#", "2")),
				),
			},
			// starts_with only
			{
				Config: accconfig.FromModels(t, share1, share2, share3, dsStartsWith),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsStartsWith.DatasourceReference(), "shares.#", "2")),
				),
			},
			// like + limit
			{
				Config: accconfig.FromModels(t, share1, share2, share3, dsLimit),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsLimit.DatasourceReference(), "shares.#", "1")),
				),
			},
			// like that matches nothing
			{
				Config: accconfig.FromModels(t, share1, share2, share3, dsNoMatch),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsNoMatch.DatasourceReference(), "shares.#", "0")),
				),
			},
		},
	})
}

func TestAcc_Shares_CompleteUseCase(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	account := secondaryTestClient().Account.GetAccountIdentifier(t)
	comment := random.Comment()

	shareModel := model.Share("test", id.Name()).
		WithComment(comment).
		WithAccounts(account.Name()).
		WithSecureObjectsOnlyBool(false)
	grantModel := model.GrantPrivilegesToShare("test", []string{sdk.ObjectPrivilegeUsage.String()}, id.Name()).
		WithOnDatabase(database.ID().Name()).
		WithDependsOn(shareModel.ResourceReference())

	ds := datasourcemodel.Shares("test").
		WithLike(id.Name()).
		WithDependsOn(shareModel.ResourceReference(), grantModel.ResourceReference())
	dsWithoutDescribe := datasourcemodel.Shares("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithDependsOn(shareModel.ResourceReference(), grantModel.ResourceReference())

	showOutput := func(ref string) assert.TestCheckFuncProvider {
		return resourceshowoutputassert.SharesDatasourceShowOutput(t, ref).
			HasName(id.Name()).
			HasKind(sdk.ShareKindOutbound).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasDatabaseName(database.ID()).
			HasComment(comment).
			HasListingGlobalName("").
			HasSecureObjectsOnly(false).
			HasToTotal(0)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// with_describe = true (default)
			{
				Config: accconfig.FromModels(t, shareModel, grantModel, ds),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(ds.DatasourceReference(), "shares.#", "1")),
					showOutput(ds.DatasourceReference()),
					assert.Check(resource.TestCheckResourceAttr(ds.DatasourceReference(), "shares.0.show_output.0.to.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(ds.DatasourceReference(), "shares.0.show_output.0.to.0", account.Name())),
					// DESCRIBE returns the granted database first and can add system objects after it
					resourceshowoutputassert.SharesDatasourceDescribeOutputRow(t, ds.DatasourceReference(), 0, 0).
						HasKind(sdk.ObjectTypeDatabase).
						HasName(database.ID()),
					assert.Check(resource.TestCheckResourceAttrSet(ds.DatasourceReference(), "shares.0.describe_output.0.shared_on")),
				),
			},
			// with_describe = false
			{
				Config: accconfig.FromModels(t, shareModel, grantModel, dsWithoutDescribe),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsWithoutDescribe.DatasourceReference(), "shares.#", "1")),
					showOutput(dsWithoutDescribe.DatasourceReference()),
					assert.Check(resource.TestCheckResourceAttr(dsWithoutDescribe.DatasourceReference(), "shares.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_Shares_InboundAndOutbound(t *testing.T) {
	prefix := random.AlphaN(6)
	outboundId := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	inboundId := secondaryTestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	providerAccount := secondaryTestClient().Account.GetAccountIdentifier(t)
	mainAccount := testClient().Account.GetAccountIdentifier(t)
	comment := random.Comment()

	providerDatabase, providerDatabaseCleanup := secondaryTestClient().Database.CreateDatabase(t)
	t.Cleanup(providerDatabaseCleanup)
	_, inboundShareCleanup := secondaryTestClient().Share.CreateWithIdentifier(t, inboundId)
	t.Cleanup(inboundShareCleanup)
	t.Cleanup(secondaryTestClient().Grant.GrantPrivilegeOnDatabaseToShare(t, providerDatabase.ID(), inboundId, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))
	secondaryTestClient().Share.Alter(t, sdk.NewAlterShareRequest(inboundId).WithAdd(sdk.ShareAddRequest{
		Accounts:          []sdk.AccountIdentifier{mainAccount},
		ShareRestrictions: sdk.Bool(false),
	}))

	outboundShare := model.Share("outbound", outboundId.Name()).WithComment(comment)
	dsBoth := datasourcemodel.Shares("both").
		WithLike(prefix + "%").
		WithDependsOn(outboundShare.ResourceReference())
	dsOutbound := datasourcemodel.Shares("outbound").
		WithLike(outboundId.Name()).
		WithDependsOn(outboundShare.ResourceReference())
	dsInbound := datasourcemodel.Shares("inbound").
		WithLike(inboundId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, outboundShare, dsBoth, dsOutbound, dsInbound),
				Check: assertThat(t,
					assert.Check(checkSharesKindsByName(dsBoth.DatasourceReference(), map[string]string{
						outboundId.Name(): string(sdk.ShareKindOutbound),
						inboundId.Name():  string(sdk.ShareKindInbound),
					})),
					resourceshowoutputassert.SharesDatasourceShowOutput(t, dsOutbound.DatasourceReference()).
						HasName(outboundId.Name()).
						HasKind(sdk.ShareKindOutbound).
						HasOwnerAccount(mainAccount.Name()).
						HasComment(comment).
						HasToTotal(0),
					assert.Check(resource.TestCheckResourceAttr(dsOutbound.DatasourceReference(), "shares.0.describe_output.#", "0")),
					resourceshowoutputassert.SharesDatasourceShowOutput(t, dsInbound.DatasourceReference()).
						HasName(inboundId.Name()).
						HasKind(sdk.ShareKindInbound).
						HasOwnerAccount(providerAccount.Name()).
						HasComment("").
						HasToTotal(0),
					assert.Check(resource.TestCheckResourceAttr(dsInbound.DatasourceReference(), "shares.0.show_output.0.to.#", "0")),
					resourceshowoutputassert.SharesDatasourceDescribeOutputRow(t, dsInbound.DatasourceReference(), 0, 0).
						HasKind(sdk.ObjectTypeDatabase).
						HasName(ids.DatabasePlaceholder),
					assert.Check(resource.TestCheckResourceAttrSet(dsInbound.DatasourceReference(), "shares.0.describe_output.0.shared_on")),
				),
			},
		},
	})
}

func checkSharesKindsByName(datasourceReference string, expected map[string]string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		datasource, ok := state.RootModule().Resources[datasourceReference]
		if !ok {
			return fmt.Errorf("data source %s not found in the state", datasourceReference)
		}
		count, err := strconv.Atoi(datasource.Primary.Attributes["shares.#"])
		if err != nil {
			return fmt.Errorf("cannot read the number of shares: %w", err)
		}
		actual := make(map[string]string, count)
		for i := range count {
			actual[datasource.Primary.Attributes[fmt.Sprintf("shares.%d.show_output.0.name", i)]] = datasource.Primary.Attributes[fmt.Sprintf("shares.%d.show_output.0.kind", i)]
		}
		if len(actual) != len(expected) {
			return fmt.Errorf("expected shares %v, got %v", expected, actual)
		}
		for name, kind := range expected {
			if actual[name] != kind {
				return fmt.Errorf("expected shares %v, got %v", expected, actual)
			}
		}
		return nil
	}
}
