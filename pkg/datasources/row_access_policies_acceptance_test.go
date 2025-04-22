//go:build !account_level_tests

package datasources_test

import (
	"maps"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	tfconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_RowAccessPolicies(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	body := "case when current_role() in ('ANALYST') then true else false end"
	policyModel := model.RowAccessPolicy("test", []sdk.TableColumnSignature{
		{
			Name: "a",
			Type: testdatatypes.DataTypeVarchar,
		},
		{
			Name: "b",
			Type: testdatatypes.DataTypeVarchar,
		},
	}, body, id.DatabaseName(), id.Name(), id.SchemaName()).WithComment("foo")

	dsName := "data.snowflake_row_access_policies.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.RowAccessPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicies/optionals_set"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsName, "row_access_policies.#", "1")),

					resourceshowoutputassert.RowAccessPoliciesDatasourceShowOutput(t, "snowflake_row_access_policies.test").
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasKind(string(sdk.PolicyKindRowAccessPolicy)).
						HasName(id.Name()).
						HasOptions("").
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasSchemaName(id.SchemaName()).
						HasComment("foo"),

					assert.Check(resource.TestCheckResourceAttr(dsName, "row_access_policies.0.describe_output.0.body", "case when current_role() in ('ANALYST') then true else false end")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "row_access_policies.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "row_access_policies.0.describe_output.0.return_type", "BOOLEAN")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "row_access_policies.0.describe_output.0.signature.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "row_access_policies.0.describe_output.0.signature.0.name", "a")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "row_access_policies.0.describe_output.0.signature.0.type", testdatatypes.DefaultVarcharAsString)),
					assert.Check(resource.TestCheckResourceAttr(dsName, "row_access_policies.0.describe_output.0.signature.1.name", "b")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "row_access_policies.0.describe_output.0.signature.1.type", testdatatypes.DefaultVarcharAsString)),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicies/optionals_unset"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, policyModel),

				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsName, "row_access_policies.#", "1")),

					resourceshowoutputassert.RowAccessPoliciesDatasourceShowOutput(t, "snowflake_row_access_policies.test").
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasKind(string(sdk.PolicyKindRowAccessPolicy)).
						HasName(id.Name()).
						HasOptions("").
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasSchemaName(id.SchemaName()).
						HasComment("foo"),
					assert.Check(resource.TestCheckResourceAttr(dsName, "row_access_policies.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_RowAccessPolicies_Filtering(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	prefix := random.AlphaN(4)
	idOne := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseId := acc.TestClient().Ids.DatabaseId()
	schemaId := acc.TestClient().Ids.SchemaId()
	commonVariables := config.Variables{
		"name_1":   config.StringVariable(idOne.Name()),
		"name_2":   config.StringVariable(idTwo.Name()),
		"name_3":   config.StringVariable(idThree.Name()),
		"schema":   config.StringVariable(schemaId.Name()),
		"database": config.StringVariable(databaseId.Name()),
		"arguments": config.SetVariable(
			config.MapVariable(map[string]config.Variable{
				"name": config.StringVariable("a"),
				"type": config.StringVariable(string(sdk.DataTypeVARCHAR)),
			}),
		),
		"body": config.StringVariable("case when current_role() in ('ANALYST') then true else false end"),
	}

	likeConfig := config.Variables{
		"like": config.StringVariable(idOne.Name()),
	}
	maps.Copy(likeConfig, commonVariables)

	likeConfig2 := config.Variables{
		"like": config.StringVariable(prefix + "%"),
	}
	maps.Copy(likeConfig2, commonVariables)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.RowAccessPolicy),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicies/like"),
				ConfigVariables: likeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_row_access_policies.test", "row_access_policies.#", "1"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicies/like"),
				ConfigVariables: likeConfig2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_row_access_policies.test", "row_access_policies.#", "2"),
				),
			},
		},
	})
}

func TestAcc_RowAccessPolicies_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      rowAccessPoliciesDatasourceEmptyIn(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func rowAccessPoliciesDatasourceEmptyIn() string {
	return `
data "snowflake_row_access_policies" "test" {
  in {
  }
}
`
}

func TestAcc_RowAccessPolicies_NotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_RowAccessPolicies/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one row access policy"),
			},
		},
	})
}
