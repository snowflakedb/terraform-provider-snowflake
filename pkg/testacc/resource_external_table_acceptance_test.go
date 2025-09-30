//go:build non_account_level_tests

package testacc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_ExternalTable_basic(t *testing.T) {
	// TODO [SNOW-1423486]: unskip
	t.Skipf("Skip because error %s; will be fixed in SNOW-1423486", "Error: 000606 (57P03): No active warehouse selected in the current session.  Select an active warehouse with the 'use warehouse' command.")

	awsBucketURL, awsKeyId, awsSecretKey := getExternalTableTestEnvsOrSkipTest(t)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()

	innerDirectory := "/external_tables_test_data/"
	configVariables := map[string]config.Variable{
		"name":           config.StringVariable(name),
		"location":       config.StringVariable(awsBucketURL),
		"aws_key_id":     config.StringVariable(awsKeyId),
		"aws_secret_key": config.StringVariable(awsSecretKey),
		"database":       config.StringVariable(TestDatabaseName),
		"schema":         config.StringVariable(TestSchemaName),
	}

	data, err := json.Marshal([]struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		{
			Name: "one",
			Age:  11,
		},
		{
			Name: "two",
			Age:  22,
		},
		{
			Name: "three",
			Age:  33,
		},
	})
	require.NoError(t, err)

	resourceName := "snowflake_external_table.test_table"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: configVariables,
			},
			{
				PreConfig: func() {
					testClient().ExternalTable.PublishDataToStage(t, id, data)
				},
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "database", TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "location", fmt.Sprintf(`@"%s"."%s"."%s"%s`, TestDatabaseName, TestSchemaName, name, innerDirectory)),
					resource.TestCheckResourceAttr(resourceName, "file_format", "TYPE = JSON, STRIP_OUTER_ARRAY = TRUE"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "column.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "column.0.name", "name"),
					resource.TestCheckResourceAttr(resourceName, "column.0.type", "string"),
					resource.TestCheckResourceAttr(resourceName, "column.0.as", "value:name::string"),
					resource.TestCheckResourceAttr(resourceName, "column.1.name", "age"),
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "number"),
					resource.TestCheckResourceAttr(resourceName, "column.1.as", "value:age::number"),
				),
			},
			{
				ConfigDirectory: ConfigurationSameAsStepN(2),
				ConfigVariables: configVariables,
				Check: externalTableContainsData(id, func(rows []map[string]*any) bool {
					expectedNames := []string{"one", "two", "three"}
					names := make([]string, 3)
					for i, row := range rows {
						nameValue, ok := row["NAME"]
						if !ok {
							return false
						}

						if nameValue == nil {
							return false
						}

						nameStringValue, ok := (*nameValue).(string)
						if !ok {
							return false
						}

						names[i] = nameStringValue
					}

					return !slices.ContainsFunc(expectedNames, func(expectedName string) bool {
						return !slices.Contains(names, expectedName)
					})
				}),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2310 is fixed
func TestAcc_ExternalTable_CorrectDataTypes(t *testing.T) {
	awsBucketURL, awsKeyId, awsSecretKey := getExternalTableTestEnvsOrSkipTest(t)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()

	innerDirectory := "/external_tables_test_data/"
	configVariables := map[string]config.Variable{
		"name":           config.StringVariable(name),
		"location":       config.StringVariable(awsBucketURL),
		"aws_key_id":     config.StringVariable(awsKeyId),
		"aws_secret_key": config.StringVariable(awsSecretKey),
		"database":       config.StringVariable(TestDatabaseName),
		"schema":         config.StringVariable(TestSchemaName),
	}

	resourceName := "snowflake_external_table.test_table"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "location", fmt.Sprintf(`@"%s"."%s"."%s"%s`, TestDatabaseName, TestSchemaName, name, innerDirectory)),
					resource.TestCheckResourceAttr(resourceName, "file_format", "TYPE = JSON, STRIP_OUTER_ARRAY = TRUE"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "column.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "column.0.name", "name"),
					resource.TestCheckResourceAttr(resourceName, "column.0.type", "varchar(200)"),
					resource.TestCheckResourceAttr(resourceName, "column.0.as", "value:name::string"),
					resource.TestCheckResourceAttr(resourceName, "column.1.name", "age"),
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "number(2, 2)"),
					resource.TestCheckResourceAttr(resourceName, "column.1.as", "value:age::number"),
					expectTableToHaveColumnDataTypes(id, []sdk.DataType{
						sdk.DataTypeVariant,
						"VARCHAR(200)",
						"NUMBER(2,2)",
					}),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2293 is fixed
func TestAcc_ExternalTable_CanCreateWithPartitions(t *testing.T) {
	awsBucketURL, awsKeyId, awsSecretKey := getExternalTableTestEnvsOrSkipTest(t)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()

	innerDirectory := "/external_tables_test_data/"
	configVariables := map[string]config.Variable{
		"name":           config.StringVariable(name),
		"location":       config.StringVariable(awsBucketURL),
		"aws_key_id":     config.StringVariable(awsKeyId),
		"aws_secret_key": config.StringVariable(awsSecretKey),
		"database":       config.StringVariable(TestDatabaseName),
		"schema":         config.StringVariable(TestSchemaName),
	}

	resourceName := "snowflake_external_table.test_table"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "location", fmt.Sprintf(`@"%s"."%s"."%s"%s`, TestDatabaseName, TestSchemaName, name, innerDirectory)),
					resource.TestCheckResourceAttr(resourceName, "file_format", "TYPE = JSON, STRIP_OUTER_ARRAY = TRUE"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "partition_by.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "partition_by.0", "filename"),
					resource.TestCheckResourceAttr(resourceName, "column.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "column.0.name", "filename"),
					resource.TestCheckResourceAttr(resourceName, "column.0.type", "string"),
					resource.TestCheckResourceAttr(resourceName, "column.0.as", "metadata$filename"),
					resource.TestCheckResourceAttr(resourceName, "column.1.name", "name"),
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "varchar(200)"),
					resource.TestCheckResourceAttr(resourceName, "column.1.as", "value:name::string"),
					resource.TestCheckResourceAttr(resourceName, "column.2.name", "age"),
					resource.TestCheckResourceAttr(resourceName, "column.2.type", "number(2, 2)"),
					resource.TestCheckResourceAttr(resourceName, "column.2.as", "value:age::number"),
					expectTableDDLContains(id, "partition by (FILENAME)"),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1564 is implemented
func TestAcc_ExternalTable_DeltaLake(t *testing.T) {
	awsBucketURL, awsKeyId, awsSecretKey := getExternalTableTestEnvsOrSkipTest(t)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()

	innerDirectory := "/external_tables_test_data/"
	configVariables := map[string]config.Variable{
		"name":           config.StringVariable(name),
		"location":       config.StringVariable(awsBucketURL),
		"aws_key_id":     config.StringVariable(awsKeyId),
		"aws_secret_key": config.StringVariable(awsSecretKey),
		"database":       config.StringVariable(TestDatabaseName),
		"schema":         config.StringVariable(TestSchemaName),
	}

	resourceName := "snowflake_external_table.test_table"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ExternalTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "location", fmt.Sprintf(`@"%s"."%s"."%s"%s`, TestDatabaseName, TestSchemaName, name, innerDirectory)),
					resource.TestCheckResourceAttr(resourceName, "file_format", "TYPE = PARQUET"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "table_format", "delta"),
					resource.TestCheckResourceAttr(resourceName, "partition_by.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "partition_by.0", "filename"),
					resource.TestCheckResourceAttr(resourceName, "column.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "column.0.name", "filename"),
					resource.TestCheckResourceAttr(resourceName, "column.0.type", "string"),
					resource.TestCheckResourceAttr(resourceName, "column.0.as", "metadata$filename"),
					resource.TestCheckResourceAttr(resourceName, "column.1.name", "name"),
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "string"),
					resource.TestCheckResourceAttr(resourceName, "column.1.as", "value:name::string"),
					func(state *terraform.State) error {
						client := TestAccProvider.Meta().(*provider.Context).Client
						ctx := context.Background()
						result, err := client.ExternalTables.ShowByID(ctx, id)
						if err != nil {
							return err
						}
						if result.TableFormat != "DELTA" {
							return fmt.Errorf("expeted table_format: DELTA, got: %s", result.TableFormat)
						}
						return nil
					},
				),
			},
		},
	})
}

func getExternalTableTestEnvsOrSkipTest(t *testing.T) (string, string, string) {
	t.Helper()
	awsBucketURL := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)
	return awsBucketURL, awsKeyId, awsSecretKey
}

func externalTableContainsData(id sdk.SchemaObjectIdentifier, contains func(rows []map[string]*any) bool) func(state *terraform.State) error {
	return func(state *terraform.State) error {
		client := TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()
		rows, err := client.QueryUnsafe(ctx, fmt.Sprintf("select * from %s", id.FullyQualifiedName()))
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Number of retrieved rows for %s: %d", id.FullyQualifiedName(), len(rows))

		if !contains(rows) {
			return fmt.Errorf("unexpected data returned by external table %s", id.FullyQualifiedName())
		}

		return nil
	}
}

func expectTableToHaveColumnDataTypes(id sdk.SchemaObjectIdentifier, expectedDataTypes []sdk.DataType) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		client := TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()
		columnsDesc, err := client.ExternalTables.DescribeColumns(ctx, sdk.NewDescribeExternalTableColumnsRequest(id))
		if err != nil {
			return err
		}

		actualTableDataTypes := make([]sdk.DataType, len(columnsDesc))
		for i, desc := range columnsDesc {
			actualTableDataTypes[i] = desc.Type
		}

		slices.SortFunc(expectedDataTypes, func(a, b sdk.DataType) int {
			return strings.Compare(string(a), string(b))
		})
		slices.SortFunc(actualTableDataTypes, func(a, b sdk.DataType) int {
			return strings.Compare(string(a), string(b))
		})

		if !slices.Equal(expectedDataTypes, actualTableDataTypes) {
			return fmt.Errorf("expected table %s to have columns with data types: %v, got: %v", id.FullyQualifiedName(), expectedDataTypes, actualTableDataTypes)
		}

		return nil
	}
}

func expectTableDDLContains(id sdk.SchemaObjectIdentifier, substr string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		client := TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()

		rows, err := client.QueryUnsafe(ctx, fmt.Sprintf("select get_ddl('table', '%s')", id.FullyQualifiedName()))
		if err != nil {
			return err
		}

		if len(rows) != 1 {
			return fmt.Errorf("unexpectedly returned more than one row: %d", len(rows))
		}

		row := rows[0]

		if len(row) != 1 {
			return fmt.Errorf("unexpectedly returned more than one columns: %d", len(row))
		}

		for _, v := range row {
			if v == nil {
				return fmt.Errorf("unexpectedly row value of ddl is nil")
			}

			ddl, ok := (*v).(string)

			if !ok {
				return fmt.Errorf("unexpectedly ddl is not type string")
			}

			if !strings.Contains(ddl, substr) {
				return fmt.Errorf("expected '%s' to be a substring of '%s'", substr, ddl)
			}
		}

		return nil
	}
}
