//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_HybridTable_basic(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with basic configuration
			{
				Config: hybridTableConfigBasic(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_hybrid_table.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_hybrid_table.test", "database", id.DatabaseName()),
					resource.TestCheckResourceAttr("snowflake_hybrid_table.test", "schema", id.SchemaName()),
					resource.TestCheckResourceAttr("snowflake_hybrid_table.test", "column.#", "3"),
					resource.TestCheckResourceAttr("snowflake_hybrid_table.test", "column.0.name", "id"),
					resource.TestCheckResourceAttr("snowflake_hybrid_table.test", "column.0.type", "NUMBER(38,0)"),
					resource.TestCheckResourceAttr("snowflake_hybrid_table.test", "column.0.nullable", "false"),
					resource.TestCheckResourceAttr("snowflake_hybrid_table.test", "column.0.primary_key", "true"),
					resource.TestCheckResourceAttr("snowflake_hybrid_table.test", "column.1.name", "name"),
					resource.TestCheckResourceAttr("snowflake_hybrid_table.test", "column.1.type", "VARCHAR(100)"),
					resource.TestCheckResourceAttr("snowflake_hybrid_table.test", "column.2.name", "created_at"),
					resource.TestCheckResourceAttr("snowflake_hybrid_table.test", "column.2.type", "TIMESTAMP_NTZ"),
					resource.TestCheckResourceAttrSet("snowflake_hybrid_table.test", "show_output.0.created_on"),
				),
			},
			// Import
			{
				ResourceName:      "snowflake_hybrid_table.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"or_replace",
					"column", // Columns are not imported back to state
					"index",
					"primary_key",
					"unique_constraint",
					"foreign_key",
				},
			},
			// Update comment
			{
				Config: hybridTableConfigWithComment(id, "updated comment"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_hybrid_table.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_hybrid_table.test", "comment", "updated comment"),
				),
			},
		},
	})
}

// Test 2: Complete table with all features
func TestAcc_HybridTable_complete(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	configModel := model.HybridTable("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithComment("complete test table").
		WithDataRetentionTimeInDays(7).
		WithColumnDescs([]model.ColumnDesc{
			{Name: "id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false)},
			{Name: "name", DataType: "VARCHAR(100)", Comment: "user name"},
			{Name: "created_at", DataType: "TIMESTAMP_NTZ", Comment: "creation timestamp"},
		}).
		WithPrimaryKeyColumns("id").
		WithIndexes([]model.IndexDesc{
			{Name: "idx_name", Columns: []string{"name"}},
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, configModel),
				Check: assert.AssertThat(t,
					resourceassert.HybridTableResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString("complete test table").
						HasDataRetentionTimeInDaysString("7").
						HasColumnCount(3).
						HasPrimaryKeyNotEmpty().
						HasIndexCount(1),
					objectassert.HybridTable(t, id).
						HasName(id.Name()).
						HasComment("complete test table"),
				),
			},
			{
				ResourceName:      configModel.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"or_replace",
					"column",
					"index",
					"primary_key",
					"unique_constraint",
					"foreign_key",
				},
			},
		},
	})
}

// Test 3: Composite primary key (out-of-line)
func TestAcc_HybridTable_compositePrimaryKey(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	configModel := model.HybridTable("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithComment("table with composite primary key").
		WithColumnDescs([]model.ColumnDesc{
			{Name: "tenant_id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false)},
			{Name: "user_id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false)},
			{Name: "email", DataType: "VARCHAR(255)"},
		}).
		WithPrimaryKeyColumns("tenant_id", "user_id")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, configModel),
				Check: assert.AssertThat(t,
					resourceassert.HybridTableResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()).
						HasColumnCount(3).
						HasPrimaryKeyNotEmpty(),
					objectassert.HybridTable(t, id).
						HasName(id.Name()).
						HasComment("table with composite primary key"),
				),
			},
			{
				ResourceName:      configModel.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"or_replace",
					"column",
					"index",
					"primary_key",
					"unique_constraint",
					"foreign_key",
				},
			},
		},
	})
}

// Test 4: Table with multiple indexes
func TestAcc_HybridTable_withIndexes(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	configModel := model.HybridTable("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithComment("table with multiple indexes").
		WithColumnDescs([]model.ColumnDesc{
			{Name: "id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false), PrimaryKey: true},
			{Name: "first_name", DataType: "VARCHAR(100)"},
			{Name: "last_name", DataType: "VARCHAR(100)"},
			{Name: "email", DataType: "VARCHAR(255)"},
		}).
		WithIndexes([]model.IndexDesc{
			{Name: "idx_email", Columns: []string{"email"}},
			{Name: "idx_name", Columns: []string{"first_name", "last_name"}},
			{Name: "idx_last_name", Columns: []string{"last_name"}},
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, configModel),
				Check: assert.AssertThat(t,
					resourceassert.HybridTableResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()).
						HasColumnCount(4).
						HasIndexCount(3),
					objectassert.HybridTable(t, id).
						HasName(id.Name()).
						HasComment("table with multiple indexes"),
				),
			},
			{
				ResourceName:      configModel.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"or_replace",
					"column",
					"index",
					"primary_key",
					"unique_constraint",
					"foreign_key",
				},
			},
		},
	})
}

// Test 5: Validation - No primary key should fail
func TestAcc_HybridTable_validation_noPrimaryKey(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	configModel := model.HybridTable("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithColumnDescs([]model.ColumnDesc{
			{Name: "col1", DataType: "NUMBER"},
			{Name: "col2", DataType: "VARCHAR(100)"},
		})
	// Intentionally no primary key defined!

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      config.FromModel(t, configModel),
				ExpectError: regexp.MustCompile("hybrid table requires a primary key"),
			},
		},
	})
}

// Test 6: Validation - Conflicting primary key definitions should fail
func TestAcc_HybridTable_validation_conflictingPK(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	configModel := model.HybridTable("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithColumnDescs([]model.ColumnDesc{
			{Name: "id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false), PrimaryKey: true}, // Inline PK
		}).
		WithPrimaryKeyColumns("id") // Out-of-line PK - conflicts!

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      config.FromModel(t, configModel),
				ExpectError: regexp.MustCompile("primary key cannot be defined both inline"),
			},
		},
	})
}

// Test 7: Update comment
func TestAcc_HybridTable_updateComment(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	configModelNoComment := model.HybridTable("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithColumnDescs([]model.ColumnDesc{
			{Name: "id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false), PrimaryKey: true},
		})

	configModelWithComment := model.HybridTable("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithComment("updated comment").
		WithColumnDescs([]model.ColumnDesc{
			{Name: "id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false), PrimaryKey: true},
		})

	configModelEmptyComment := model.HybridTable("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithComment("").
		WithColumnDescs([]model.ColumnDesc{
			{Name: "id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false), PrimaryKey: true},
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create without comment
			{
				Config: config.FromModel(t, configModelNoComment),
				Check: assert.AssertThat(t,
					resourceassert.HybridTableResource(t, configModelNoComment.ResourceReference()).
						HasNameString(id.Name()).
						HasNoComment(),
				),
			},
			// Update with comment
			{
				Config: config.FromModel(t, configModelWithComment),
				Check: assert.AssertThat(t,
					resourceassert.HybridTableResource(t, configModelWithComment.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString("updated comment"),
					objectassert.HybridTable(t, id).
						HasComment("updated comment"),
				),
			},
			// Remove comment (set to empty)
			{
				Config: config.FromModel(t, configModelEmptyComment),
				Check: assert.AssertThat(t,
					resourceassert.HybridTableResource(t, configModelEmptyComment.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentEmpty(),
				),
			},
		},
	})
}

// Test 8: Update data_retention_time_in_days
func TestAcc_HybridTable_updateDataRetention(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	configModelWithRetention7 := model.HybridTable("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithDataRetentionTimeInDays(7).
		WithColumnDescs([]model.ColumnDesc{
			{Name: "id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false), PrimaryKey: true},
		})

	configModelWithRetention14 := model.HybridTable("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithDataRetentionTimeInDays(14).
		WithColumnDescs([]model.ColumnDesc{
			{Name: "id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false), PrimaryKey: true},
		})

	configModelNoRetention := model.HybridTable("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithColumnDescs([]model.ColumnDesc{
			{Name: "id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false), PrimaryKey: true},
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with data_retention_time_in_days = 7
			{
				Config: config.FromModel(t, configModelWithRetention7),
				Check: assert.AssertThat(t,
					resourceassert.HybridTableResource(t, configModelWithRetention7.ResourceReference()).
						HasNameString(id.Name()).
						HasDataRetentionTimeInDaysString("7"),
				),
			},
			// Update to 14
			{
				Config: config.FromModel(t, configModelWithRetention14),
				Check: assert.AssertThat(t,
					resourceassert.HybridTableResource(t, configModelWithRetention14.ResourceReference()).
						HasNameString(id.Name()).
						HasDataRetentionTimeInDaysString("14"),
				),
			},
			// Remove (unset)
			{
				Config: config.FromModel(t, configModelNoRetention),
				Check: assert.AssertThat(t,
					resourceassert.HybridTableResource(t, configModelNoRetention.ResourceReference()).
						HasNameString(id.Name()).
						HasNoDataRetentionTimeInDays(),
				),
			},
		},
	})
}

// Test 9: Import test
func TestAcc_HybridTable_import(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	configModel := model.HybridTable("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithComment("import test table").
		WithDataRetentionTimeInDays(10).
		WithColumnDescs([]model.ColumnDesc{
			{Name: "id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false), PrimaryKey: true},
			{Name: "name", DataType: "VARCHAR(100)"},
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create
			{
				Config: config.FromModel(t, configModel),
				Check: assert.AssertThat(t,
					resourceassert.HybridTableResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString("import test table").
						HasDataRetentionTimeInDaysString("10"),
				),
			},
			// Import
			{
				ResourceName:      configModel.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"or_replace",
					"column",
					"index",
					"primary_key",
					"unique_constraint",
					"foreign_key",
				},
			},
		},
	})
}

// Test 10: Disappears test (drift handling)
func TestAcc_HybridTable_disappears(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	configModel := model.HybridTable("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithColumnDescs([]model.ColumnDesc{
			{Name: "id", DataType: "NUMBER(38,0)", Nullable: model.Bool(false), PrimaryKey: true},
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, configModel),
				Check: assert.AssertThat(t,
					resourceassert.HybridTableResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()),
					objectassert.HybridTable(t, id).
						HasName(id.Name()),
				),
			},
			{
				Config: config.FromModel(t, configModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					acc.TestClient().HybridTable.DropFunc(t, id)(),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func hybridTableConfigBasic(id sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
provider "snowflake" {
  preview_features_enabled = ["snowflake_hybrid_table_resource"]
}

resource "snowflake_hybrid_table" "test" {
  database = "%s"
  schema   = "%s"
  name     = "%s"

  column {
    name        = "id"
    type        = "NUMBER(38,0)"
    nullable    = false
    primary_key = true
  }

  column {
    name = "name"
    type = "VARCHAR(100)"
  }

  column {
    name = "created_at"
    type = "TIMESTAMP_NTZ"
  }
}
`, id.DatabaseName(), id.SchemaName(), id.Name())
}

func hybridTableConfigWithComment(id sdk.SchemaObjectIdentifier, comment string) string {
	return fmt.Sprintf(`
provider "snowflake" {
  preview_features_enabled = ["snowflake_hybrid_table_resource"]
}

resource "snowflake_hybrid_table" "test" {
  database = "%s"
  schema   = "%s"
  name     = "%s"
  comment  = "%s"

  column {
    name        = "id"
    type        = "NUMBER(38,0)"
    nullable    = false
    primary_key = true
  }

  column {
    name = "name"
    type = "VARCHAR(100)"
  }

  column {
    name = "created_at"
    type = "TIMESTAMP_NTZ"
  }
}
`, id.DatabaseName(), id.SchemaName(), id.Name(), comment)
}
