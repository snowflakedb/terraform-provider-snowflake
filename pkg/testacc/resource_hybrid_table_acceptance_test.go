//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_HybridTable_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
	}
	pk := []sdk.TableColumnSignature{
		{Name: "ID"},
	}

	modelBasic := model.HybridTableFromId("test", id, columns, pk)

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.HybridTableResource(t, modelBasic.ResourceReference()).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCommentString("").
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasColumns(columns).
			HasPrimaryKeyKeys("ID"),
		objectparametersassert.HybridTableParameters(t, id).
			HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeSnowflakeDefault).
			HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeSnowflakeDefault),
		resourceshowoutputassert.HybridTableShowOutput(t, modelBasic.ResourceReference()).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""),
		objectassert.HybridTable(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(""),
	}

	modelComplete := model.HybridTableFromId("test", id, columns, pk).
		WithComment(comment).
		WithDataRetentionTimeInDays(1).
		WithMaxDataExtensionTimeInDays(10)

	assertComplete := []assert.TestCheckFuncProvider{
		resourceassert.HybridTableResource(t, modelComplete.ResourceReference()).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCommentString(comment).
			HasDataRetentionTimeInDaysString("1").
			HasMaxDataExtensionTimeInDaysString("10").
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasColumns(columns).
			HasPrimaryKeyKeys("ID"),
		objectparametersassert.HybridTableParameters(t, id).
			HasDataRetentionTimeInDays(1).
			HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeHybridTable).
			HasMaxDataExtensionTimeInDays(10).
			HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeHybridTable),
		resourceshowoutputassert.HybridTableShowOutput(t, modelComplete.ResourceReference()).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(comment),
		objectassert.HybridTable(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment),
	}

	assertAfterUnset := []assert.TestCheckFuncProvider{
		resourceassert.HybridTableResource(t, modelBasic.ResourceReference()).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCommentString("").
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasColumns(columns).
			HasPrimaryKeyKeys("ID"),
		objectparametersassert.HybridTableParameters(t, id).
			HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeSnowflakeDefault).
			HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeSnowflakeDefault),
		resourceshowoutputassert.HybridTableShowOutput(t, modelBasic.ResourceReference()).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""),
		objectassert.HybridTable(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(""),
	}

	importStateVerifyIgnore := []string{
		// DESCRIBE normalizes types (e.g. INTEGER -> NUMBER(38,0)); DiffSuppressDataTypes
		// handles this at plan time, but the raw state values differ after import.
		"column",
		// Constraint name may differ between config (empty) and what DESCRIBE returns.
		"primary_key",
		// Computed + Optional parameter fields: import reads the live Snowflake value via
		// ShowParameters (e.g. the account default), which can differ from a config that
		// omits the field. Ignoring during ImportStateVerify is standard for
		// Computed + Optional fields.
		"data_retention_time_in_days",
		"max_data_extension_time_in_days",
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:                  accconfig.FromModels(t, modelBasic),
				ResourceName:            modelBasic.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importStateVerifyIgnore,
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelComplete),
				Check:  assertThat(t, assertComplete...),
			},
			// Import - with optionals
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importStateVerifyIgnore,
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, assertAfterUnset...),
			},
			// External change detection
			{
				PreConfig: func() {
					testClient().HybridTable.Alter(t, sdk.NewAlterHybridTableRequest(id).WithSet(
						*sdk.NewHybridTableSetPropertiesRequest().WithComment("external comment"),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, assertAfterUnset...),
			},
			// External deletion — resource dropped outside Terraform; Read detects absence
			// (ErrObjectNotFound → d.SetId("")) and the next plan recreates it.
			{
				PreConfig: func() {
					testClient().HybridTable.DropFunc(t, id)()
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, assertBasic...),
			},
			// Destroy
			{
				Destroy: true,
				Config:  accconfig.FromModels(t, modelBasic),
			},
			// Create - with optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, modelComplete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}

func TestAcc_HybridTable_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment, changedComment := random.Comment(), random.Comment()

	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar},
	}
	pk := []sdk.TableColumnSignature{
		{Name: "ID"},
	}

	modelComplete := model.HybridTableFromId("test", id, columns, pk).
		WithComment(comment).
		WithDataRetentionTimeInDays(5).
		WithMaxDataExtensionTimeInDays(10)

	modelChanged := model.HybridTableFromId("test", id, columns, pk).
		WithComment(changedComment).
		WithDataRetentionTimeInDays(10).
		WithMaxDataExtensionTimeInDays(20)

	importStateVerifyIgnore := []string{
		// DESCRIBE normalizes types (e.g. INTEGER -> NUMBER(38,0)); DiffSuppressDataTypes
		// handles this at plan time, but the raw state values differ after import.
		"column",
		// Constraint name may differ between config (empty) and what DESCRIBE returns.
		"primary_key",
		// Computed + Optional parameter fields: import reads the live Snowflake value via
		// ShowParameters (e.g. the account default), which can differ from a config that
		// omits the field. Ignoring during ImportStateVerify is standard for
		// Computed + Optional fields.
		"data_retention_time_in_days",
		"max_data_extension_time_in_days",
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create - with all attributes
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(comment).
						HasDataRetentionTimeInDaysString("5").
						HasMaxDataExtensionTimeInDaysString("10").
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasColumns(columns).
						HasPrimaryKeyKeys("ID"),
					objectparametersassert.HybridTableParameters(t, id).
						HasDataRetentionTimeInDays(5).
						HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeHybridTable).
						HasMaxDataExtensionTimeInDays(10).
						HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeHybridTable),
					resourceshowoutputassert.HybridTableShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(comment),
					objectassert.HybridTable(t, id).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment),
				),
			},
			// Import
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importStateVerifyIgnore,
			},
			// Update - change mutable properties
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelChanged.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelChanged),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, modelChanged.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(changedComment).
						HasDataRetentionTimeInDaysString("10").
						HasMaxDataExtensionTimeInDaysString("20").
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasColumns(columns).
						HasPrimaryKeyKeys("ID"),
					objectparametersassert.HybridTableParameters(t, id).
						HasDataRetentionTimeInDays(10).
						HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeHybridTable).
						HasMaxDataExtensionTimeInDays(20).
						HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeHybridTable),
					resourceshowoutputassert.HybridTableShowOutput(t, modelChanged.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(changedComment),
					objectassert.HybridTable(t, id).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(changedComment),
				),
			},
		},
	})
}

// TestAcc_HybridTable_InvalidConfig verifies that schema-level validators reject
// out-of-range values before a Snowflake connection is needed.
func TestAcc_HybridTable_InvalidConfig(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pk := []sdk.TableColumnSignature{{Name: "ID"}}
	cols := []sdk.TableColumnSignature{{Name: "ID", Type: testdatatypes.DataTypeInteger}}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t,
					model.HybridTableFromId("test", id, cols, pk).WithColumnConfigs([]model.HybridTableColumnConfig{
						{Name: "ID", Type: "INVALIDTYPE"},
					}),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid data type`),
			},
		},
	})
}

func TestAcc_HybridTable_UniqueConstraint(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	cols := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar},
	}
	pk := []sdk.TableColumnSignature{{Name: "ID"}}

	// Single-column unique constraint
	model1 := model.HybridTableFromId("test", id, cols, pk).
		WithUniqueConstraint([]string{"NAME"})

	// Change the unique constraint to span two columns — forces recreation
	model2 := model.HybridTableFromId("test", id, cols, pk).
		WithUniqueConstraint([]string{"NAME", "EMAIL"})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with a single-column unique constraint
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumns(cols).
						HasPrimaryKeyKeys("ID").
						HasUniqueConstraintCount(1),
				),
			},
			// Change the unique constraint columns — any diff on unique_constraint forces recreation
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(cols).
						HasPrimaryKeyKeys("ID").
						HasUniqueConstraintCount(1),
				),
			},
		},
	})
}

func TestAcc_HybridTable_ForeignKey(t *testing.T) {
	// Create parent hybrid table externally — it is not managed by Terraform in this test.
	parentId := testClient().Ids.RandomSchemaObjectIdentifier()
	testClient().HybridTable.CreateWithRequest(t, parentId, sdk.HybridTableColumnsConstraintsAndIndexesRequest{
		Columns: []sdk.HybridTableColumnRequest{
			*sdk.NewHybridTableColumnRequest("ID", sdk.DataType("INTEGER")).
				WithInlineConstraint(sdk.ColumnInlineConstraint{Type: sdk.ColumnConstraintTypePrimaryKey}),
		},
	})
	t.Cleanup(testClient().HybridTable.DropFunc(t, parentId))

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	cols := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "PARENT_ID", Type: testdatatypes.DataTypeInteger},
	}
	pk := []sdk.TableColumnSignature{{Name: "ID"}}

	// Child table with FK → parent.ID
	model1 := model.HybridTableFromId("test", id, cols, pk).
		WithForeignKey([]string{"PARENT_ID"}, parentId.FullyQualifiedName(), []string{"ID"})

	// Child table without FK
	model2 := model.HybridTableFromId("test", id, cols, pk)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with a foreign key referencing the parent table
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumns(cols).
						HasPrimaryKeyKeys("ID").
						HasForeignKeyCount(1),
				),
			},
			// Remove the foreign key — any diff on foreign_key forces recreation
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(cols).
						HasPrimaryKeyKeys("ID").
						HasForeignKeyCount(0),
				),
			},
		},
	})
}

func TestAcc_HybridTable_ColumnDefault(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pk := []sdk.TableColumnSignature{{Name: "ID"}}
	baseCols := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "SCORE", Type: testdatatypes.DataTypeInteger},
	}

	zero := "0"

	// SCORE has a constant default of 0
	model1 := model.HybridTableFromId("test", id, baseCols, pk).
		WithColumnConfigs([]model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
			{Name: "SCORE", Type: testdatatypes.DataTypeInteger.ToSql(), Default: &model.HybridTableColumnDefaultConfig{Constant: &zero}},
		})

	// SCORE has no default
	model2 := model.HybridTableFromId("test", id, baseCols, pk)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with a constant default on SCORE
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumns(baseCols).
						HasColumnDefaultConstant(1, "0"),
				),
			},
			// Drop the default → in-place update, no recreation
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(baseCols).
						HasColumnNoDefault(1),
				),
			},
		},
	})
}

func TestAcc_HybridTable_PrimaryKeyForceNew(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	cols := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}

	// Single-column PK
	model1 := model.HybridTableFromId("test", id, cols, []sdk.TableColumnSignature{{Name: "ID"}})

	// Composite PK — any change to primary_key forces recreation
	model2 := model.HybridTableFromId("test", id, cols, []sdk.TableColumnSignature{{Name: "ID"}, {Name: "NAME"}})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with single-column PK
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumns(cols).
						HasPrimaryKeyKeys("ID"),
				),
			},
			// Change to composite PK → ForceNew (DestroyBeforeCreate)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(cols).
						HasPrimaryKeyKeys("ID", "NAME"),
				),
			},
		},
	})
}

func TestAcc_HybridTable_ColumnAdd(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pk := []sdk.TableColumnSignature{{Name: "ID"}}

	cols1 := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
	}
	cols2 := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}
	cols3 := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar},
		{Name: "AGE", Type: testdatatypes.DataTypeInteger},
	}

	model1 := model.HybridTableFromId("test", id, cols1, pk)
	model2 := model.HybridTableFromId("test", id, cols2, pk)
	model3 := model.HybridTableFromId("test", id, cols3, pk)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with single column
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumns(cols1).
						HasPrimaryKeyKeys("ID"),
				),
			},
			// Add one column
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(cols2).
						HasPrimaryKeyKeys("ID"),
				),
			},
			// Add two more columns in one apply
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model3.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model3),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model3.ResourceReference()).
						HasColumns(cols3).
						HasPrimaryKeyKeys("ID"),
				),
			},
		},
	})
}

func TestAcc_HybridTable_ColumnDrop(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pk := []sdk.TableColumnSignature{{Name: "ID"}}

	cols1 := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar},
		{Name: "AGE", Type: testdatatypes.DataTypeInteger},
	}
	cols2 := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
		{Name: "EMAIL", Type: testdatatypes.DataTypeVarchar},
	}
	cols3 := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
	}

	model1 := model.HybridTableFromId("test", id, cols1, pk)
	model2 := model.HybridTableFromId("test", id, cols2, pk)
	model3 := model.HybridTableFromId("test", id, cols3, pk)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with four columns
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumns(cols1).
						HasPrimaryKeyKeys("ID"),
				),
			},
			// Drop one column
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(cols2).
						HasPrimaryKeyKeys("ID"),
				),
			},
			// Drop two more columns in one apply
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model3.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model3),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model3.ResourceReference()).
						HasColumns(cols3).
						HasPrimaryKeyKeys("ID"),
				),
			},
		},
	})
}

func TestAcc_HybridTable_ColumnAlterComment(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pk := []sdk.TableColumnSignature{{Name: "ID"}}
	baseCols := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}

	// No comment on NAME
	model1 := model.HybridTableFromId("test", id, baseCols, pk)

	// Set comment on NAME
	model2 := model.HybridTableFromId("test", id, baseCols, pk).
		WithColumnConfigs([]model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
			{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Comment: "this is a name column"},
		})

	// Change comment on NAME
	model3 := model.HybridTableFromId("test", id, baseCols, pk).
		WithColumnConfigs([]model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
			{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Comment: "updated comment"},
		})

	// Unset comment on NAME (back to no comment)
	model4 := model.HybridTableFromId("test", id, baseCols, pk)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with two columns, NAME has no comment
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumns(baseCols).
						HasColumnComment(1, ""),
				),
			},
			// Set comment on NAME
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(baseCols).
						HasColumnComment(1, "this is a name column"),
				),
			},
			// Change comment on NAME
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model3.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model3),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model3.ResourceReference()).
						HasColumns(baseCols).
						HasColumnComment(1, "updated comment"),
				),
			},
			// Unset comment on NAME
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model4.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, model4),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model4.ResourceReference()).
						HasColumns(baseCols).
						HasColumnComment(1, ""),
				),
			},
		},
	})
}

func TestAcc_HybridTable_ColumnNullableForceNew(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pk := []sdk.TableColumnSignature{{Name: "ID"}}
	baseCols := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}

	// NAME is explicitly nullable=true
	model1 := model.HybridTableFromId("test", id, baseCols, pk).
		WithColumnConfigs([]model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
			{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Nullable: sdk.Bool(true)},
		})

	// NAME changed to nullable=false — must force recreation
	model2 := model.HybridTableFromId("test", id, baseCols, pk).
		WithColumnConfigs([]model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
			{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Nullable: sdk.Bool(false)},
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with NAME nullable=true
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumns(baseCols).
						HasColumnNullable(1, true),
				),
			},
			// Change NAME nullable=false — expect DestroyBeforeCreate (ForceNew)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(baseCols).
						HasColumnNullable(1, false),
				),
			},
		},
	})
}

func TestAcc_HybridTable_ColumnCollateForceNew(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pk := []sdk.TableColumnSignature{{Name: "ID"}}
	baseCols := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}

	// NAME with collate='en'
	model1 := model.HybridTableFromId("test", id, baseCols, pk).
		WithColumnConfigs([]model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
			{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Collate: "en"},
		})

	// NAME collate changed to 'fr' — must force recreation
	model2 := model.HybridTableFromId("test", id, baseCols, pk).
		WithColumnConfigs([]model.HybridTableColumnConfig{
			{Name: "ID", Type: testdatatypes.DataTypeInteger.ToSql()},
			{Name: "NAME", Type: testdatatypes.DataTypeVarchar.ToSql(), Collate: "fr"},
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create with NAME collate='en'
			{
				Config: accconfig.FromModels(t, model1),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model1.ResourceReference()).
						HasColumns(baseCols),
				),
			},
			// Change NAME collate='fr' — expect DestroyBeforeCreate (ForceNew)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model2.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, model2),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, model2.ResourceReference()).
						HasColumns(baseCols),
				),
			},
		},
	})
}

func TestAcc_HybridTable_Rename(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	schemaId := sdk.NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName())
	newId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeInteger},
	}
	pk := []sdk.TableColumnSignature{
		{Name: "ID"},
	}

	renamedComment := random.Comment()
	modelBasic := model.HybridTableFromId("test", id, columns, pk)
	modelRenamed := model.HybridTableFromId("test", newId, columns, pk).WithComment(renamedComment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.HybridTable),
		Steps: []resource.TestStep{
			// Create
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					objectassert.HybridTable(t, id).
						HasName(id.Name()),
				),
			},
			// Rename
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelRenamed.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelRenamed),
				Check: assertThat(t,
					resourceassert.HybridTableResource(t, modelRenamed.ResourceReference()).
						HasNameString(newId.Name()).
						HasDatabaseString(newId.DatabaseName()).
						HasSchemaString(newId.SchemaName()).
						HasFullyQualifiedNameString(newId.FullyQualifiedName()).
						HasCommentString(renamedComment),
					objectassert.HybridTable(t, newId).
						HasName(newId.Name()).
						HasDatabaseName(newId.DatabaseName()).
						HasSchemaName(newId.SchemaName()).
						HasComment(renamedComment),
				),
			},
		},
	})
}
