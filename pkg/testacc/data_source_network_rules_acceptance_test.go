//go:build non_account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_NetworkRules_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)

	nr1, nr1Cleanup := testClient().NetworkRule.CreateWithIdentifier(t, id1)
	t.Cleanup(nr1Cleanup)

	_, nr2Cleanup := testClient().NetworkRule.CreateWithIdentifier(t, id2)
	t.Cleanup(nr2Cleanup)

	likePrefix := datasourcemodel.NetworkRules("test").
		WithLike(prefix + "%").
		WithInSchema(id1.SchemaId())

	likeExact := datasourcemodel.NetworkRules("test").
		WithLike(nr1.Name).
		WithInSchema(id1.SchemaId())

	startsWith := datasourcemodel.NetworkRules("test").
		WithStartsWith(prefix).
		WithInSchema(id1.SchemaId())

	limitOne := datasourcemodel.NetworkRules("test").
		WithLike(prefix + "%").
		WithInSchema(id1.SchemaId()).
		WithRows(1)
	inDatabase := datasourcemodel.NetworkRules("test").
		WithInDatabase(id1.DatabaseId())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.NetworkRule),
		Steps: []resource.TestStep{
			// like (prefix)
			{
				Config: config.DatasourceFromModel(t, likePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(likePrefix.DatasourceReference(), "network_rules.#", "2"),
				),
			},
			// like (exact)
			{
				Config: config.DatasourceFromModel(t, likeExact),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(likeExact.DatasourceReference(), "network_rules.#", "1"),
					resource.TestCheckResourceAttr(likeExact.DatasourceReference(), "network_rules.0.show_output.0.name", nr1.Name),
				),
			},
			// starts_with
			{
				Config: config.DatasourceFromModel(t, startsWith),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(startsWith.DatasourceReference(), "network_rules.#", "2"),
				),
			},
			// limit
			{
				Config: config.DatasourceFromModel(t, limitOne),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(limitOne.DatasourceReference(), "network_rules.#", "1"),
				),
			},
			// in database
			{
				Config: config.DatasourceFromModel(t, inDatabase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(inDatabase.DatasourceReference(), "network_rules.#"),
				),
			},
		},
	})
}

func TestAcc_NetworkRules_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	nr, nrCleanup := testClient().NetworkRule.CreateWithRequest(t,
		sdk.NewCreateNetworkRuleRequest(id, sdk.NetworkRuleTypeIpv4, []sdk.NetworkRuleValue{{Value: "1.2.3.4"}, {Value: "5.6.7.8"}}, sdk.NetworkRuleModeIngress).
			WithComment(comment),
	)
	t.Cleanup(nrCleanup)

	withoutDescribe := datasourcemodel.NetworkRules("test").
		WithWithDescribe(false).
		WithLike(nr.Name).
		WithInSchema(id.SchemaId())

	withDescribe := datasourcemodel.NetworkRules("test").
		WithWithDescribe(true).
		WithLike(nr.Name).
		WithInSchema(id.SchemaId())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.NetworkRule),
		Steps: []resource.TestStep{
			{
				Config: config.DatasourceFromModel(t, withoutDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.NetworkRulesDatasourceShowOutput(t, withoutDescribe.DatasourceReference()).
						HasCreatedOnNotEmpty().
						HasName(nr.Name).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasType(sdk.NetworkRuleTypeIpv4).
						HasMode(sdk.NetworkRuleModeIngress).
						HasEntriesInValueList(2).
						HasOwnerRoleType("ROLE"),

					assert.Check(resource.TestCheckResourceAttr(withoutDescribe.DatasourceReference(), "network_rules.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(withoutDescribe.DatasourceReference(), "network_rules.0.describe_output.#", "0")),
				),
			},
			{
				Config: config.DatasourceFromModel(t, withDescribe),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(withDescribe.DatasourceReference(), "network_rules.#", "1")),
					resourceshowoutputassert.NetworkRulesDatasourceShowOutput(t, withDescribe.DatasourceReference()).
						HasCreatedOnNotEmpty().
						HasName(nr.Name).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasType(sdk.NetworkRuleTypeIpv4).
						HasMode(sdk.NetworkRuleModeIngress).
						HasEntriesInValueList(2).
						HasOwnerRoleType("ROLE"),
					resourceshowoutputassert.NetworkRulesDatasourceDescribeOutput(t, withDescribe.DatasourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.NetworkRuleTypeIpv4).
						HasMode(sdk.NetworkRuleModeIngress).
						HasComment(comment).
						HasValueList([]string{"1.2.3.4", "5.6.7.8"}).
						HasOwner(snowflakeroles.Accountadmin.Name()),
				),
			},
		},
	})
}
