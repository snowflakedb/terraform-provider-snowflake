//go:build non_account_level_tests

package testacc

import (
	"context"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-3151661]: test creation with DECFLOAT column

func TestAcc_Table_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}
	comment := random.Comment()
	commentChanged := random.Comment()
	externalComment := random.Comment()

	modelBasic := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns)
	modelWithAlterableOptional := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).
		WithComment(comment).
		WithChangeTracking(true).
		WithDataRetentionTimeInDays(5).
		WithClusterBy("ID")
	modelWithAlterableOptionalChanged := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).
		WithComment(commentChanged).
		WithChangeTracking(true).
		WithDataRetentionTimeInDays(3).
		WithClusterBy("NAME")
	modelWithAlterableOptionalUnset := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns)
	modelWithClusterByIdName := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).
		WithClusterBy("ID", "NAME")
	modelWithClusterByNameId := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).
		WithClusterBy("NAME", "ID")

	ref := modelBasic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.TableResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasCommentEmpty().
			HasClusterByEmpty().
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.TableShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasChangeTracking(false),
		resourceparametersassert.TableResourceParameters(t, ref).
			HasDataRetentionTimeInDays(1).
			HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeDatabase),
		resourceshowoutputassert.TableColumnDescribeOutputOnIdx(t, ref, 0).
			HasName("ID").
			HasKind(sdk.TableColumnKindColumn),
	}

	// change_tracking is absent from the state when it was never set. After it is removed from the configuration, the
	// external change detection writes the value now reported by Snowflake, and the diff suppressor hides it from the plan.
	basicNeverSetAssertions := append([]assert.TestCheckFuncProvider{resourceassert.TableResource(t, ref).HasNoChangeTracking()}, basicAssertions...)
	basicUnsetAssertions := append([]assert.TestCheckFuncProvider{resourceassert.TableResource(t, ref).HasChangeTracking(false)}, basicAssertions...)

	alterableOptionalAssertions := []assert.TestCheckFuncProvider{
		resourceassert.TableResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasComment(comment).
			HasChangeTracking(true).
			HasClusterBy("ID").
			HasDataRetentionTimeInDays(5).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.TableShowOutput(t, ref).
			HasName(id.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasChangeTracking(true),
		resourceparametersassert.TableResourceParameters(t, ref).
			HasDataRetentionTimeInDays(5).
			HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeTable),
	}

	clusterByIdNameAssertions := []assert.TestCheckFuncProvider{
		resourceassert.TableResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasClusterBy("ID", "NAME").
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
	}
	clusterByNameIdAssertions := []assert.TestCheckFuncProvider{
		resourceassert.TableResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasClusterBy("NAME", "ID").
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
	}

	changedOptionalAssertions := []assert.TestCheckFuncProvider{
		resourceassert.TableResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasComment(commentChanged).
			HasChangeTracking(true).
			HasClusterBy("NAME").
			HasDataRetentionTimeInDays(3).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.TableShowOutput(t, ref).
			HasComment(commentChanged).
			HasChangeTracking(true),
		resourceparametersassert.TableResourceParameters(t, ref).
			HasDataRetentionTimeInDays(3).
			HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeTable),
	}

	providerModelImportPermissiveMode := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.ImportPermissiveMode)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			// Create with only required fields
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, modelBasic),
				Check:                    assertThat(t, basicNeverSetAssertions...),
			},
			// Import - without optionals, permissive mode
			{
				ProtoV6ProviderFactories: importPermissiveModeProviderFactory,
				Config:                   accconfig.FromModels(t, providerModelImportPermissiveMode, modelBasic),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateKind:          resource.ImportBlockWithID,
			},
			// Change only alterable fields - expect an in-place update
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, modelWithAlterableOptional),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, alterableOptionalAssertions...),
			},
			// Import with all alterable fields set, permissive mode
			{
				ProtoV6ProviderFactories: importPermissiveModeProviderFactory,
				Config:                   accconfig.FromModels(t, providerModelImportPermissiveMode, modelWithAlterableOptional),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateKind:          resource.ImportBlockWithID,
			},
			// Change the alterable fields externally and detect drift - expect Terraform to fix it with an update
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig: func() {
					testClient().Table.AlterWithRequest(t, sdk.NewAlterTableRequest(id).WithSet(*sdk.NewTableSetRequest().
						WithComment(externalComment).
						WithChangeTracking(false).
						WithDataRetentionTimeInDays(1)))
					testClient().Table.AlterWithRequest(t, sdk.NewAlterTableRequest(id).WithClusteringAction(
						*sdk.NewTableClusteringActionRequest().WithDropClusteringKey(true),
					))
				},
				Config: accconfig.FromModels(t, modelWithAlterableOptionalChanged),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, changedOptionalAssertions...),
			},
			// Unset all alterable optional fields - expect an in-place update
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, modelWithAlterableOptionalUnset),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, basicUnsetAssertions...),
			},
			// Set cluster_by - expect an in-place update
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, modelWithAlterableOptional),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, alterableOptionalAssertions...),
			},
			// Change cluster_by - expect an in-place update
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, modelWithAlterableOptionalChanged),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, changedOptionalAssertions...),
			},
			// Unset cluster_by (and the other alterable fields) - expect an in-place update
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, modelWithAlterableOptionalUnset),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, basicUnsetAssertions...),
			},
			// Set two cluster_by keys - expect an in-place update
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, modelWithClusterByIdName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, clusterByIdNameAssertions...),
			},
			// Reorder cluster_by keys - expect an in-place update
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, modelWithClusterByNameId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, clusterByNameIdAssertions...),
			},
		},
	})
}

func TestAcc_Table_BasicUseCase_Import(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	resourceId := helpers.EncodeResourceIdentifier(id)
	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}

	modelBasic := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns)
	ref := modelBasic.ResourceReference()

	_, cleanup := testClient().Table.CreateWithRequest(t, sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns([]sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("ID", testdatatypes.DataTypeNumber_38_0),
		*sdk.NewTableColumnRequest("NAME", testdatatypes.DataTypeVarchar),
	})))
	t.Cleanup(cleanup)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			// Import - without optionals, enforcing mode
			{
				Config:             accconfig.FromModels(t, modelBasic),
				ResourceName:       ref,
				ImportState:        true,
				ImportStateId:      resourceId,
				ImportStatePersist: true,
				ImportStateCheck: assertThatImport(
					t,
					resourceassert.ImportedTableResource(t, resourceId).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCommentEmpty().
						HasChangeTracking(false).
						HasEnableSchemaEvolution("false").
						HasErrorLogging("false").
						HasRowTimestamp("false").
						HasIsTransient("false"),
				),
			},
			// Expect empty plan after the plan after import is applied.
			// is_transient forces a new resource, so resetting it to the unset value replaces the table.
			{
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectChange(ref, "enable_schema_evolution", tfjson.ActionDelete, new("false"), new(r.BooleanDefault)),
						planchecks.ExpectChange(ref, "error_logging", tfjson.ActionDelete, new("false"), new(r.BooleanDefault)),
						planchecks.ExpectChange(ref, "row_timestamp", tfjson.ActionDelete, new("false"), new(r.BooleanDefault)),
						planchecks.ExpectChange(ref, "is_transient", tfjson.ActionDelete, new("false"), new(r.BooleanDefault)),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(
					t,
					resourceassert.TableResource(t, ref).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasCommentEmpty().
						HasNoChangeTracking().
						HasEnableSchemaEvolution(r.BooleanDefault).
						HasErrorLogging(r.BooleanDefault).
						HasRowTimestamp(r.BooleanDefault).
						HasIsTransient(r.BooleanDefault),
				),
			},
		},
	})
}

func TestAcc_Table_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}
	comment := random.Comment()

	modelComplete := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).
		WithComment(comment).
		WithChangeTracking(true).
		WithDataRetentionTimeInDays(5).
		WithClusterBy("ID", "NAME")

	ref := modelComplete.ResourceReference()

	providerModelImportPermissiveMode := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.ImportPermissiveMode)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, modelComplete),
				Check: assertThat(
					t,
					resourceassert.TableResource(t, ref).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasComment(comment).
						HasChangeTracking(true).
						HasClusterBy("ID", "NAME").
						HasDataRetentionTimeInDays(5).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.TableShowOutput(t, ref).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasChangeTracking(true),
					resourceparametersassert.TableResourceParameters(t, ref).
						HasDataRetentionTimeInDays(5).
						HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeTable),
				),
			},
			// import - complete, enforcing mode
			// The configuration omits these fields, so the state after the create holds the unset value ("default"),
			// while the enforcing import writes the value reported by Snowflake.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, modelComplete),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateVerify:        true,
				ImportStateVerifyIgnore:  []string{"is_transient", "enable_schema_evolution", "error_logging", "row_timestamp"},
			},
			// import - complete, permissive mode
			{
				ProtoV6ProviderFactories: importPermissiveModeProviderFactory,
				Config:                   accconfig.FromModels(t, providerModelImportPermissiveMode, modelComplete),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateKind:          resource.ImportBlockWithID,
			},
		},
	})
}

func TestAcc_Table_BasicUseCase_Columns(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	placeholder := []sdk.TableColumnSignature{
		{Name: "COL1", Type: testdatatypes.DataTypeVarchar},
	}
	commentA := random.Comment()
	commentB := random.Comment()

	columnsCreate := []model.TableColumnConfig{
		{Name: "COL1", Type: testdatatypes.DataTypeVarchar, Comment: commentA},
		{Name: "COL2", Type: testdatatypes.DataTypeFloat},
	}
	columnsCommentsAltered := []model.TableColumnConfig{
		{Name: "COL1", Type: testdatatypes.DataTypeVarchar},
		{Name: "COL2", Type: testdatatypes.DataTypeFloat, Comment: commentB},
	}
	columnsReplaced := []model.TableColumnConfig{
		{Name: "COL2", Type: testdatatypes.DataTypeFloat, Comment: commentB},
		{Name: "COLUMN3", Type: testdatatypes.DataTypeNumber_38_0},
	}
	columnsNotNull := []model.TableColumnConfig{
		{Name: "COL2", Type: testdatatypes.DataTypeFloat, Comment: commentB},
		{Name: "COLUMN3", Type: testdatatypes.DataTypeNumber_38_0, NotNull: new(true)},
	}

	// The masking policies take VARCHAR arguments. The conditional one has two of them.
	maskingPolicy1, maskingPolicy1Cleanup := testClient().MaskingPolicy.CreateMaskingPolicyIdentity(t, testdatatypes.DataTypeVarchar)
	t.Cleanup(maskingPolicy1Cleanup)
	maskingPolicy2, maskingPolicy2Cleanup := testClient().MaskingPolicy.CreateMaskingPolicyIdentity(t, testdatatypes.DataTypeVarchar)
	t.Cleanup(maskingPolicy2Cleanup)
	conditionalMaskingPolicy, conditionalMaskingPolicyCleanup := testClient().MaskingPolicy.CreateMaskingPolicy(t)
	t.Cleanup(conditionalMaskingPolicyCleanup)

	projectionPolicy1, projectionPolicy1Cleanup := testClient().ProjectionPolicy.CreateProjectionPolicy(t)
	t.Cleanup(projectionPolicy1Cleanup)
	projectionPolicy2, projectionPolicy2Cleanup := testClient().ProjectionPolicy.CreateProjectionPolicy(t)
	t.Cleanup(projectionPolicy2Cleanup)

	columnsVarcharAdded := append(slices.Clone(columnsNotNull),
		model.TableColumnConfig{Name: "COL4", Type: testdatatypes.DataTypeVarchar},
		model.TableColumnConfig{Name: "COL5", Type: testdatatypes.DataTypeVarchar},
	)
	// withCol4 builds the columns of the policy steps: the not null columns, COL4 as given, COL5, and the extra columns.
	withCol4 := func(col4 model.TableColumnConfig, extra ...model.TableColumnConfig) []model.TableColumnConfig {
		return append(slices.Clone(columnsNotNull), append([]model.TableColumnConfig{col4, {Name: "COL5", Type: testdatatypes.DataTypeVarchar}}, extra...)...)
	}
	columnsMaskingPolicy := withCol4(model.TableColumnConfig{
		Name: "COL4", Type: testdatatypes.DataTypeVarchar, MaskingPolicy: new(maskingPolicy1.ID()),
	})
	// An unquoted policy name does not drift (https://github.com/snowflakedb/terraform-provider-snowflake/issues/4679,
	// https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2186).
	columnsMaskingPolicyChanged := withCol4(model.TableColumnConfig{
		Name: "COL4", Type: testdatatypes.DataTypeVarchar, MaskingPolicy: new(maskingPolicy2.ID()), MaskingPolicyUnquoted: true,
	})
	columnsConditionalMaskingPolicy := withCol4(model.TableColumnConfig{
		Name: "COL4", Type: testdatatypes.DataTypeVarchar, MaskingPolicy: new(conditionalMaskingPolicy.ID()), MaskingPolicyUsing: []string{"COL4", "COL5"},
	})

	columnsProjectionPolicy := withCol4(model.TableColumnConfig{
		Name: "COL4", Type: testdatatypes.DataTypeVarchar, ProjectionPolicy: &projectionPolicy1,
	})
	// changes the policy of COL4 and adds COL6 with a policy
	columnsProjectionPolicyChangedAndAdded := withCol4(
		model.TableColumnConfig{Name: "COL4", Type: testdatatypes.DataTypeVarchar, ProjectionPolicy: &projectionPolicy2},
		model.TableColumnConfig{Name: "COL6", Type: testdatatypes.DataTypeVarchar, ProjectionPolicy: &projectionPolicy1},
	)
	columnsProjectionPolicyUnset := withCol4(
		model.TableColumnConfig{Name: "COL4", Type: testdatatypes.DataTypeVarchar},
		model.TableColumnConfig{Name: "COL6", Type: testdatatypes.DataTypeVarchar},
	)

	modelCreate := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).
		WithColumnConfigs(columnsCreate...)
	modelCommentsAltered := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).
		WithColumnConfigs(columnsCommentsAltered...)
	modelReplaced := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).
		WithColumnConfigs(columnsReplaced...)
	modelNotNull := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).
		WithColumnConfigs(columnsNotNull...)
	modelVarcharAdded := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).
		WithColumnConfigs(columnsVarcharAdded...)
	modelMaskingPolicy := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).
		WithColumnConfigs(columnsMaskingPolicy...)
	modelMaskingPolicyChanged := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).
		WithColumnConfigs(columnsMaskingPolicyChanged...)
	modelConditionalMaskingPolicy := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).
		WithColumnConfigs(columnsConditionalMaskingPolicy...)
	modelProjectionPolicy := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).
		WithColumnConfigs(columnsProjectionPolicy...)
	modelProjectionPolicyChangedAndAdded := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).
		WithColumnConfigs(columnsProjectionPolicyChangedAndAdded...)
	modelProjectionPolicyUnset := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).
		WithColumnConfigs(columnsProjectionPolicyUnset...)
	ref := modelCreate.ResourceReference()

	updateWithoutDrift := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
		},
		PostApplyPostRefresh: []plancheck.PlanCheck{
			plancheck.ExpectEmptyPlan(),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelCreate),
				Check:  assertThat(t, tableColumnListChecks(t, ref, columnsCreate)...),
			},
			{
				Config: accconfig.FromModels(t, modelCommentsAltered),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, tableColumnListChecks(t, ref, columnsCommentsAltered)...),
			},
			{
				Config: accconfig.FromModels(t, modelReplaced),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, tableColumnListChecks(t, ref, columnsReplaced)...),
			},
			{
				Config: accconfig.FromModels(t, modelNotNull),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, tableColumnListChecks(t, ref, columnsNotNull)...),
			},
			{
				Config:           accconfig.FromModels(t, modelVarcharAdded),
				ConfigPlanChecks: updateWithoutDrift,
				Check:            assertThat(t, tableColumnListChecks(t, ref, columnsVarcharAdded)...),
			},
			// set the masking policy on an existing column
			{
				Config:           accconfig.FromModels(t, modelMaskingPolicy),
				ConfigPlanChecks: updateWithoutDrift,
				Check:            assertThat(t, tableColumnListChecks(t, ref, columnsMaskingPolicy)...),
			},
			// change it to another policy, written without quotes
			{
				Config:           accconfig.FromModels(t, modelMaskingPolicyChanged),
				ConfigPlanChecks: updateWithoutDrift,
				Check:            assertThat(t, tableColumnListChecks(t, ref, columnsMaskingPolicyChanged)...),
			},
			// change it to a conditional policy
			{
				Config:           accconfig.FromModels(t, modelConditionalMaskingPolicy),
				ConfigPlanChecks: updateWithoutDrift,
				Check:            assertThat(t, tableColumnListChecks(t, ref, columnsConditionalMaskingPolicy)...),
			},
			// unset it
			{
				Config:           accconfig.FromModels(t, modelVarcharAdded),
				ConfigPlanChecks: updateWithoutDrift,
				Check:            assertThat(t, tableColumnListChecks(t, ref, columnsVarcharAdded)...),
			},
			// set the projection policy on an existing column
			{
				Config:           accconfig.FromModels(t, modelProjectionPolicy),
				ConfigPlanChecks: updateWithoutDrift,
				Check:            assertThat(t, tableColumnListChecks(t, ref, columnsProjectionPolicy)...),
			},
			// change it and add a column with a projection policy
			{
				Config:           accconfig.FromModels(t, modelProjectionPolicyChangedAndAdded),
				ConfigPlanChecks: updateWithoutDrift,
				Check:            assertThat(t, tableColumnListChecks(t, ref, columnsProjectionPolicyChangedAndAdded)...),
			},
			// unset them
			{
				Config:           accconfig.FromModels(t, modelProjectionPolicyUnset),
				ConfigPlanChecks: updateWithoutDrift,
				Check:            assertThat(t, tableColumnListChecks(t, ref, columnsProjectionPolicyUnset)...),
			},
		},
	})
}

// tableColumnListChecks asserts the columns in the state. A masking policy without USING is reported with the column
// itself as the only USING argument.
func tableColumnListChecks(t *testing.T, ref string, columns []model.TableColumnConfig) []assert.TestCheckFuncProvider {
	t.Helper()
	expected := make([]resourceassert.ExpectedTableColumn, len(columns))
	for i, c := range columns {
		expected[i] = resourceassert.ExpectedTableColumn{
			Name:               c.Name,
			Type:               c.Type.Canonical(),
			Comment:            c.Comment,
			NotNull:            c.NotNull != nil && *c.NotNull,
			MaskingPolicy:      c.MaskingPolicy,
			MaskingPolicyUsing: c.MaskingPolicyUsing,
			ProjectionPolicy:   c.ProjectionPolicy,
		}
		if c.MaskingPolicy != nil && len(c.MaskingPolicyUsing) == 0 {
			expected[i].MaskingPolicyUsing = []string{c.Name}
		}
	}
	return []assert.TestCheckFuncProvider{resourceassert.TableResource(t, ref).HasColumns(expected...)}
}

func TestAcc_Table_ColumnCollate_Inherited(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	placeholder := []sdk.TableColumnSignature{
		{Name: "COL1", Type: testdatatypes.DataTypeVarchar},
	}
	columnOmitted := model.TableColumnConfig{Name: "COL1", Type: testdatatypes.DataTypeVarchar}

	// The column is created without collate, so it inherits the table default_ddl_collation. DESCRIBE reports it,
	// but an omitted collate is not managed, so it must not plan a change, also after the table default changes.
	modelInherited := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).
		WithColumnConfigs(columnOmitted).
		WithDefaultDdlCollation("en-ci")
	modelInheritedTableDefaultChanged := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).
		WithColumnConfigs(columnOmitted).
		WithDefaultDdlCollation("en-cs")
	ref := modelInherited.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelInherited),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(ref, "column.0.collate", "en-ci")),
				),
			},
			{
				Config:   accconfig.FromModels(t, modelInherited),
				PlanOnly: true,
			},
			{
				Config: accconfig.FromModels(t, modelInheritedTableDefaultChanged),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.TableResource(t, ref).HasDefaultDdlCollation("en-cs"),
					assert.Check(resource.TestCheckResourceAttr(ref, "column.0.collate", "en-ci")),
				),
			},
			{
				Config:   accconfig.FromModels(t, modelInheritedTableDefaultChanged),
				PlanOnly: true,
			},
		},
	})
}

func TestAcc_Table_ColumnCollate_ExplicitAndCleared(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	placeholder := []sdk.TableColumnSignature{
		{Name: "COL1", Type: testdatatypes.DataTypeVarchar},
	}
	columnOmitted := model.TableColumnConfig{Name: "COL1", Type: testdatatypes.DataTypeVarchar}
	columnEn := model.TableColumnConfig{Name: "COL1", Type: testdatatypes.DataTypeVarchar, Collate: new("en")}
	columnCleared := model.TableColumnConfig{Name: "COL1", Type: testdatatypes.DataTypeVarchar, Collate: new("")}

	modelEn := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).WithColumnConfigs(columnEn)
	modelOmitted := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).WithColumnConfigs(columnOmitted)
	modelCleared := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).WithColumnConfigs(columnCleared)
	ref := modelEn.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelEn),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(ref, "column.0.collate", "en")),
				),
			},
			// Removing collate from the config means "not managed": no change is planned.
			{
				Config:   accconfig.FromModels(t, modelOmitted),
				PlanOnly: true,
			},
			// An explicit empty collate is managed, so it is a planned change. Snowflake cannot change the collation
			// of an existing column, so the apply fails with its error instead of the diff being hidden.
			{
				Config: accconfig.FromModels(t, modelCleared),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				ExpectError: regexp.MustCompile(`because they have incompatible collations`),
			},
		},
	})
}

func TestAcc_Table_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).WithEnableSchemaEvolution("INVALID")),
				ExpectError: regexp.MustCompile(`expected .*enable_schema_evolution.* to be one of .*, got INVALID`),
			},
			// An expression that contains only whitespace is rejected, unlike "" which is an empty string default.
			{
				Config: accconfig.FromModels(t, model.TableWithId("test", id, nil).WithColumnConfigs(
					model.TableColumnConfig{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0, DefaultExpression: new("  ")},
				)),
				ExpectError: regexp.MustCompile(`The default expression cannot be only whitespace`),
			},
		},
	})
}

func tableConfig(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	data_retention_time_in_days = 1
	comment  = "Terraform acceptance test"
	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func TestAcc_TableDefaults(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableColumnWithDefaults(tableId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_table.test_table", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", TestSchemaName),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "change_tracking"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "5"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARCHAR(16)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.default.0.expression", "'hello'"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "TIMESTAMP_NTZ(9)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.default.0.expression", "CURRENT_TIMESTAMP()"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.type", "NUMBER(38,0)"),
					resource.TestMatchResourceAttr("snowflake_table.test_table", "column.2.default.0.expression", regexp.MustCompile(`(?i)\.NEXTVAL$`)),
					// https://github.com/snowflakedb/terraform-provider-snowflake/issues/2406
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.name", "column4"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.default.0.expression", "''"),
					// https://github.com/snowflakedb/terraform-provider-snowflake/issues/5167
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.4.name", "column5"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.4.default.0.expression", "'NUMBER(38,15)'"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableColumnWithoutDefaults(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", TestSchemaName),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "change_tracking"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARCHAR(16)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.default.#", "0"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "TIMESTAMP_NTZ(9)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.default.#", "0"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.type", "NUMBER(38,0)"),
					resource.TestMatchResourceAttr("snowflake_table.test_table", "column.2.default.0.expression", regexp.MustCompile(`(?i)\.NEXTVAL$`)),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
		},
	})
}

func tableColumnWithDefaults(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_sequence" "test_seq" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
}

resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment             = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARCHAR(16)"
		default {
			expression = "'hello'"
		}
	}
	column {
		name = "column2"
		type = "TIMESTAMP_NTZ(9)"
		default {
			expression = "CURRENT_TIMESTAMP()"
		}
	}
	column {
		name = "column3"
		type = "NUMBER(38,0)"
		default {
			expression = "${snowflake_sequence.test_seq.fully_qualified_name}.NEXTVAL"
		}
	}
	column {
		name = "column4"
		type = "VARCHAR(22)"
		default {
			expression = "''"
		}
	}
	column {
		name = "column5"
		type = "STRING"
		default {
			expression = "'NUMBER(38,15)'"
		}
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableColumnWithoutDefaults(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_sequence" "test_seq" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
}

resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	comment             = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARCHAR(16)"
	}
	column {
		name = "column2"
		type = "TIMESTAMP_NTZ(9)"
	}
	column {
		name = "column3"
		type = "NUMBER(38,0)"
		default {
			expression = "${snowflake_sequence.test_seq.fully_qualified_name}.NEXTVAL"
		}
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

// proves https://github.com/snowflakedb/terraform-provider-snowflake/issues/2406 is fixed
// An empty column.default.expression is an empty string default. It is sent as the empty string literal (a bare DEFAULT
// is not valid SQL) and does not drift against the literal returned by DESCRIBE TABLE, so both spellings are interchangeable.
func TestAcc_Table_EmptyStringColumnDefault_issue2406(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()

	emptyExpression := model.TableColumnConfig{Name: "EMPTY_EXPRESSION", Type: testdatatypes.DataTypeVarchar_50, DefaultExpression: new("")}
	emptyLiteral := model.TableColumnConfig{Name: "EMPTY_LITERAL", Type: testdatatypes.DataTypeVarchar_50, DefaultExpression: new("''")}
	emptyExpressionAsLiteral := model.TableColumnConfig{Name: emptyExpression.Name, Type: emptyExpression.Type, DefaultExpression: emptyLiteral.DefaultExpression}
	emptyLiteralAsExpression := model.TableColumnConfig{Name: emptyLiteral.Name, Type: emptyLiteral.Type, DefaultExpression: emptyExpression.DefaultExpression}
	text := model.TableColumnConfig{Name: "TEXT", Type: testdatatypes.DataTypeVarchar_50, DefaultExpression: new("'text'")}
	addedEmptyExpression := model.TableColumnConfig{Name: "ADDED_EMPTY_EXPRESSION", Type: testdatatypes.DataTypeVarchar_50, DefaultExpression: new("")}

	tableModel := func(columns ...model.TableColumnConfig) *model.TableModel {
		return model.TableWithId("test", tableId, nil).WithColumnConfigs(columns...)
	}
	modelBasic := tableModel(emptyExpression, emptyLiteral, text)
	modelSwapped := tableModel(emptyExpressionAsLiteral, emptyLiteralAsExpression, text)
	modelWithAddedColumn := tableModel(emptyExpression, emptyLiteral, text, addedEmptyExpression)
	ref := modelBasic.ResourceReference()

	// DESCRIBE TABLE returns the empty string default as '' regardless of how it was configured
	assertSnowflakeColumnDefaults := func(expectedDefaults map[string]string) assert.TestCheckFuncProvider {
		return assert.Check(func(*terraform.State) error {
			details := testClient().Table.DescribeColumns(t, tableId)
			require.Len(t, details, len(expectedDefaults))
			for _, column := range details {
				assertThatObject(t, objectassert.TableColumnDetailsFromObject(t, &column).
					HasName(column.Name).
					HasDefault(expectedDefaults[column.Name]),
				)
			}
			return nil
		})
	}

	expectedColumn := func(column model.TableColumnConfig) resourceassert.ExpectedTableColumn {
		return resourceassert.ExpectedTableColumn{Name: column.Name, Type: column.Type.ToSql(), DefaultExpression: *column.DefaultExpression}
	}
	// the empty string default is reported as '' whichever way it was configured
	expectedEmptyDefault := func(column model.TableColumnConfig) resourceassert.ExpectedTableColumn {
		return resourceassert.ExpectedTableColumn{Name: column.Name, Type: column.Type.ToSql(), DefaultExpression: "''"}
	}

	emptyPlanAfterApply := resource.ConfigPlanChecks{
		PostApplyPostRefresh: []plancheck.PlanCheck{
			plancheck.ExpectEmptyPlan(),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			// Create: an empty expression and the explicit literal produce the same default.
			{
				Config:           accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: emptyPlanAfterApply,
				Check: assertThat(t,
					resourceassert.TableResource(t, ref).
						HasNameString(tableId.Name()).
						HasColumns(expectedEmptyDefault(emptyExpression), expectedEmptyDefault(emptyLiteral), expectedColumn(text)),
					assertSnowflakeColumnDefaults(map[string]string{
						emptyExpression.Name: "''",
						emptyLiteral.Name:    "''",
						text.Name:            "'text'",
					}),
				),
			},
			// Swap the two spellings: the default is the same, so there is nothing to change.
			{
				Config: accconfig.FromModels(t, modelSwapped),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(t,
					resourceassert.TableResource(t, ref).
						HasColumns(expectedEmptyDefault(emptyExpressionAsLiteral), expectedEmptyDefault(emptyLiteralAsExpression), expectedColumn(text)),
					assertSnowflakeColumnDefaults(map[string]string{
						emptyExpression.Name: "''",
						emptyLiteral.Name:    "''",
						text.Name:            "'text'",
					}),
				),
			},
			// Add a column with an empty expression to the existing table (ALTER TABLE ... ADD COLUMN).
			{
				Config:           accconfig.FromModels(t, modelWithAddedColumn),
				ConfigPlanChecks: emptyPlanAfterApply,
				Check: assertThat(t,
					resourceassert.TableResource(t, ref).
						HasColumns(expectedEmptyDefault(emptyExpression), expectedEmptyDefault(emptyLiteral), expectedColumn(text), expectedEmptyDefault(addedEmptyExpression)),
					assertSnowflakeColumnDefaults(map[string]string{
						emptyExpression.Name:      "''",
						emptyLiteral.Name:         "''",
						text.Name:                 "'text'",
						addedEmptyExpression.Name: "''",
					}),
				),
			},
		},
	})
}

func TestAcc_TableIdentity(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	ref := "snowflake_table.test_table"

	noDriftAfterApply := resource.ConfigPlanChecks{
		PostApplyPostRefresh: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
		},
	}
	tableChecks := func(columns ...resourceassert.ExpectedTableColumn) assert.TestCheckFuncProvider {
		return resourceassert.TableResource(t, ref).
			HasNameString(tableId.Name()).
			HasDatabaseString(TestDatabaseName).
			HasSchemaString(TestSchemaName).
			HasNoChangeTracking().
			HasCommentString("Terraform acceptance test").
			HasPrimaryKeyEmpty().
			HasColumns(columns...)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableColumnWithIdentityDefault(tableId),
				// order is not set, so the value Snowflake reports does not plan a change
				ConfigPlanChecks: noDriftAfterApply,
				Check: assertThat(t, tableChecks(
					resourceassert.ExpectedTableColumn{Name: "column1", Type: "NUMBER(38,0)"},
					resourceassert.ExpectedTableColumn{Name: "column2", Type: "TIMESTAMP_NTZ(9)"},
					resourceassert.ExpectedTableColumn{Name: "column3", Type: "NUMBER(38,0)", Identity: &resourceassert.ExpectedTableColumnIdentity{StartNum: 1, StepNum: 1}},
				)),
			},
			// https://github.com/snowflakedb/terraform-provider-snowflake/issues/3221
			// the previous identity column is dropped, and adding a new column as an identity works
			{
				Config:           tableColumnWithIdentity(tableId),
				ConfigPlanChecks: noDriftAfterApply,
				Check: assertThat(t, tableChecks(
					resourceassert.ExpectedTableColumn{Name: "column1", Type: "NUMBER(38,0)"},
					resourceassert.ExpectedTableColumn{Name: "column2", Type: "TIMESTAMP_NTZ(9)"},
					resourceassert.ExpectedTableColumn{Name: "column4", Type: "NUMBER(38,0)", Identity: &resourceassert.ExpectedTableColumnIdentity{StartNum: 2, StepNum: 4, Order: new(true)}},
				)),
			},
		},
	})
}

func tableColumnWithIdentityDefault(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_sequence" "test_seq" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
}

resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
	comment             = "Terraform acceptance test"

	column {
		name = "column1"
		type = "NUMBER(38,0)"
	}
	column {
		name = "column2"
		type = "TIMESTAMP_NTZ(9)"
	}
	column {
		name = "column3"
		type = "NUMBER(38,0)"
		identity {
		}
	}
	depends_on = [snowflake_sequence.test_seq]
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableColumnWithIdentity(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_sequence" "test_seq" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
}

resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
	comment             = "Terraform acceptance test"

	column {
		name = "column1"
		type = "NUMBER(38,0)"
	}
	column {
		name = "column2"
		type = "TIMESTAMP_NTZ(9)"
	}

	column {
		name = "column4"
		type = "NUMBER(38,0)"
		identity {
			start_num = 2
			step_num = 4
			order = "true"
		}
	}
	depends_on = [snowflake_sequence.test_seq]
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func TestAcc_TableCollate(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableColumnWithCollate(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.collate", "en"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.collate", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.collate", ""),
				),
			},
			{
				Config: addColumnWithCollate(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "4"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.name", "column4"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.collate", "utf8"),
				),
			},
			{
				Config:      alterTableColumnWithIncompatibleCollate(tableId),
				ExpectError: regexp.MustCompile("\"VARCHAR\\(100\\) COLLATE 'fr'\" because they have incompatible collations\\."),
			},
		},
	})
}

func tableColumnWithCollate(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
		name = "column1"
		type = "VARCHAR(100)"
		collate = "en"
	}
	column {
		name = "column2"
		type = "VARCHAR(100)"
		collate = ""
	}
	column {
		name = "column3"
		type = "VARCHAR(100)"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func addColumnWithCollate(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
		name = "column1"
		type = "VARCHAR(100)"
		collate = "en"
	}
	column {
		name = "column2"
		type = "VARCHAR(100)"
		collate = ""
	}
	column {
		name = "column3"
		type = "VARCHAR(100)"
	}
	column {
		name = "column4"
		type = "VARCHAR"
		collate = "utf8"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func alterTableColumnWithIncompatibleCollate(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
		name = "column1"
		type = "VARCHAR(100)"
		collate = "fr"
	}
	column {
		name = "column2"
		type = "VARCHAR(100)"
		collate = ""
	}
	column {
		name = "column3"
		type = "VARCHAR(100)"
	}
	column {
		name = "column4"
		type = "VARCHAR"
		collate = "utf8"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func TestAcc_TableRename(t *testing.T) {
	oldId := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	oldComment := random.Comment()
	newComment := random.Comment()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfigWithName(oldId, oldComment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", oldId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", oldId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", oldComment),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARIANT"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableConfigWithName(newId, newComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_table.test_table", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "id", newId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", newId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", newComment),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARIANT"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
		},
	})
}

func tableConfigWithName(tableId sdk.SchemaObjectIdentifier, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
    comment  = "%[4]s"
	column {
		name = "column1"
		type = "VARIANT"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), comment)
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Table_DefaultDataRetentionTime(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testClient().Database.UpdateDataRetentionTime(t, database.ID(), 5)
				},
				Config: tableConfigWithoutDataRetentionTimeInDays(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "5"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 5, 5, 5),
				),
			},
			{
				PreConfig: func() {
					testClient().Schema.UpdateDataRetentionTime(t, schema.ID(), 10)
				},
				Config: tableConfigWithoutDataRetentionTimeInDays(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "10"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 5, 10, 10),
				),
			},
			{
				PreConfig: func() {
					testClient().Database.UpdateDataRetentionTime(t, database.ID(), 10)
					testClient().Schema.UpdateDataRetentionTime(t, schema.ID(), 3)
				},
				Config: tableConfigWithDataRetentionTimeInDays(tableId, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "5"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 5),
				),
			},
			{
				Config: tableConfigWithDataRetentionTimeInDays(tableId, 15),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "15"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 15),
				),
			},
			{
				Config: tableConfigWithoutDataRetentionTimeInDays(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "3"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 3),
				),
			},
			{
				PreConfig: func() {
					testClient().Schema.UnsetDataRetentionTime(t, schema.ID())
				},
				Config: tableConfigWithoutDataRetentionTimeInDays(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "10"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 10, 10),
				),
			},
			{
				PreConfig: func() {
					testClient().Schema.UpdateDataRetentionTime(t, schema.ID(), 5)
				},
				Config: tableConfigWithDataRetentionTimeInDays(tableId, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "0"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 5, 0),
				),
			},
			{
				Config: tableConfigWithDataRetentionTimeInDays(tableId, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "3"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 5, 3),
				),
			},
		},
	})
}

func tableConfigWithoutDataRetentionTimeInDays(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
    	name = "id"
    	type = "NUMBER(38,0)"
  	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfigWithDataRetentionTimeInDays(tableId sdk.SchemaObjectIdentifier, dataRetentionTimeInDays int) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	data_retention_time_in_days = %[4]d

	column {
    	name = "id"
    	type = "NUMBER(38,0)"
  	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), dataRetentionTimeInDays)
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Table_DefaultDataRetentionTime_SetOutsideOfTerraform(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testClient().Database.UpdateDataRetentionTime(t, database.ID(), 5)
				},
				Config: tableConfigWithoutDataRetentionTimeInDays(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "5"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 5, 5, 5),
				),
			},
			{
				PreConfig: func() {
					testClient().Table.SetDataRetentionTime(t, tableId, 20)
				},
				Config: tableConfigWithoutDataRetentionTimeInDays(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "5"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 5, 5, 5),
				),
			},
			{
				PreConfig: func() {
					testClient().Schema.UpdateDataRetentionTime(t, schema.ID(), 10)
				},
				Config: tableConfigWithDataRetentionTimeInDays(tableId, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "3"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 5, 10, 3),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Table_DefaultDataRetentionTimeSettingUnsetting(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testClient().Database.UpdateDataRetentionTime(t, database.ID(), 10)
					testClient().Schema.UpdateDataRetentionTime(t, schema.ID(), 3)
				},
				Config: tableConfigWithDataRetentionTimeInDays(tableId, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "5"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 5),
				),
			},
			{
				Config: tableConfigWithoutDataRetentionTimeInDays(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "3"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 3),
				),
			},
			{
				Config: tableConfigWithDataRetentionTimeInDays(tableId, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "5"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 5),
				),
			},
		},
	})
}

func checkDatabaseSchemaAndTableDataRetentionTime(id sdk.SchemaObjectIdentifier, expectedDatabaseRetentionDays int, expectedSchemaRetentionDays int, expectedTableRetentionsDays int) func(state *terraform.State) error {
	return func(state *terraform.State) error {
		client := TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()

		database, err := client.Databases.ShowByID(ctx, id.DatabaseId())
		if err != nil {
			return err
		}

		if database.RetentionTime != expectedDatabaseRetentionDays {
			return fmt.Errorf("invalid database retention time, expected: %d, got: %d", expectedDatabaseRetentionDays, database.RetentionTime)
		}

		s, err := client.Schemas.ShowByID(ctx, id.SchemaId())
		if err != nil {
			return err
		}

		// "retention_time" may sometimes be an empty string instead of an integer
		var schemaRetentionTime int64
		{
			rt := s.RetentionTime
			if rt == "" {
				rt = "0"
			}

			schemaRetentionTime, err = strconv.ParseInt(rt, 10, 64)
			if err != nil {
				return err
			}
		}

		if schemaRetentionTime != int64(expectedSchemaRetentionDays) {
			return fmt.Errorf("invalid schema retention time, expected: %d, got: %d", expectedSchemaRetentionDays, schemaRetentionTime)
		}

		table, err := client.Tables.ShowByID(ctx, id)
		if err != nil {
			return err
		}

		if table.RetentionTime != expectedTableRetentionsDays {
			return fmt.Errorf("invalid table retention time, expected: %d, got: %d", expectedTableRetentionsDays, table.RetentionTime)
		}

		return nil
	}
}

// proves issues https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2110 and https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2495
func TestAcc_Table_ClusterBy(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfigWithComplexClusterBy(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "cluster_by.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "cluster_by.0", "date_trunc('month', LAST_LOAD_TIME)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "cluster_by.1", "COL1"),
				),
			},
		},
	})
}

func tableConfigWithComplexClusterBy(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database            = "%[1]s"
	schema              = "%[2]s"
	name                = "%[3]s"
	cluster_by = ["date_trunc('month', LAST_LOAD_TIME)", "COL1"]
	column {
		name = "COL1"
		type = "VARCHAR(16777216)"
	}
    column {
        name     = "LAST_LOAD_TIME"
        type     = "TIMESTAMP_LTZ(6)"
        not_null = false
    }
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

// proves https://github.com/snowflakedb/terraform-provider-snowflake/issues/2588 is fixed
func TestAcc_ColumnTypeChangeWithNonTextType(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	placeholder := []sdk.TableColumnSignature{{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0}}

	columnsWide := []model.TableColumnConfig{{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0}}
	columnsNarrow := []model.TableColumnConfig{{Name: "ID", Type: testdatatypes.DataTypeNumber_36_0}}

	modelWide := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).WithColumnConfigs(columnsWide...)
	modelNarrow := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).WithColumnConfigs(columnsNarrow...)
	ref := modelWide.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelWide),
				Check:  assertThat(t, tableColumnListChecks(t, ref, columnsWide)...),
			},
			// Decreasing the precision is an in-place ALTER without COLLATE.
			{
				Config: accconfig.FromModels(t, modelNarrow),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, tableColumnListChecks(t, ref, columnsNarrow)...),
			},
			{
				Config: accconfig.FromModels(t, modelWide),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, tableColumnListChecks(t, ref, columnsWide)...),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2733 is fixed
func TestAcc_Table_gh2733(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfigGh2733(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "MY_INT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "NUMBER(38,0)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "MY_STRING"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "VARCHAR(16777216)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.name", "MY_DATE"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.type", "TIMESTAMP_NTZ(9)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.name", "MY_DATE2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.type", "TIMESTAMP_NTZ(9)"),
				),
			},
		},
	})
}

func tableConfigGh2733(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
  database            = "%[1]s"
  schema              = "%[2]s"
  name                = "%[3]s"

  column {
    name = "MY_INT"
    type = "int"
    # type  = "NUMBER(38,0)" # Should be equivalent
  }

  column {
    name = "MY_STRING"
    type = "VARCHAR(16777216)"
    # type = "STRING" # Should be equivalent
  }

  column {
    name = "MY_DATE"
    type = "TIMESTAMP_NTZ"
    # type = "TIMESTAMP_NTZ(9)" # Should be equivalent
  }

  column {
    name = "MY_DATE2"
    type = "DATETIME"
    # type = "TIMESTAMP_NTZ" # Equivalent to TIMESTAMP_NTZ
  }
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

// proves https://github.com/snowflakedb/terraform-provider-snowflake/issues/2844 is fixed
// and documents https://github.com/snowflakedb/terraform-provider-snowflake/issues/1387 (column names are case-sensitive)
func TestAcc_Table_issue2844_apostropheInColumnName(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	columnNames := []string{"my'id", `My"Col`, "project_id"}
	tableModel := model.TableWithId("test", tableId, collections.Map(columnNames, func(name string) sdk.TableColumnSignature {
		return sdk.TableColumnSignature{Name: name, Type: testdatatypes.DataTypeInteger}
	}))
	ref := tableModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tableModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(
					t,
					resourceassert.TableResource(t, ref).
						HasNameString(tableId.Name()),
					assert.Check(resource.TestCheckResourceAttr(ref, "column.0.name", columnNames[0])),
					assert.Check(resource.TestCheckResourceAttr(ref, "column.1.name", columnNames[1])),
					assert.Check(resource.TestCheckResourceAttr(ref, "column.2.name", columnNames[2])),
					assert.Check(func(*terraform.State) error {
						cols := testClient().Table.GetTableColumnsFor(t, tableId)
						if len(cols) != len(columnNames) {
							return fmt.Errorf("expected %d Snowflake columns; got %d", len(columnNames), len(cols))
						}
						for i, name := range columnNames {
							if cols[i].ColumnName != name {
								return fmt.Errorf("expected Snowflake column name %q; got %q", name, cols[i].ColumnName)
							}
						}
						return nil
					}),
				),
			},
		},
	})
}

// proves https://github.com/snowflakedb/terraform-provider-snowflake/issues/4146 is fixed
func TestAcc_Table_StructuredArray_issue4146(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	placeholder := []sdk.TableColumnSignature{{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0}}
	withArrayType := func(arrayType string) *model.TableModel {
		return model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), placeholder).
			WithColumnValue(tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"name": tfconfig.StringVariable("SOME_STRUCTURED_ARRAY"),
				"type": tfconfig.StringVariable(arrayType),
			})))
	}
	modelExplicit := withArrayType("ARRAY(NUMBER(38,0))")
	modelShort := withArrayType("ARRAY(NUMBER)")
	ref := modelExplicit.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelExplicit),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(ref, "column.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(ref, "column.0.type", "ARRAY(NUMBER(38,0))")),
				),
			},
			{
				Config: accconfig.FromModels(t, modelShort),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(ref, "column.0.type", "ARRAY(NUMBER(38,0))")),
				),
			},
		},
	})
}

// proves https://github.com/snowflakedb/terraform-provider-snowflake/issues/2418 is fixed
func TestAcc_Table_VirtualColumn_issue2418(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()

	id := model.TableColumnConfig{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0}
	doubleId := func(expression string) model.TableColumnConfig {
		return model.TableColumnConfig{Name: "DOUBLE_ID", Type: testdatatypes.DataTypeNumber_38_0, As: expression, Comment: "computed"}
	}
	plusOne := model.TableColumnConfig{Name: "PLUS_ONE", Type: testdatatypes.DataTypeNumber_38_0, As: `"ID" + 1`}
	invalid := model.TableColumnConfig{Name: "INVALID", Type: testdatatypes.DataTypeNumber_38_0, As: "ID + 2", DefaultExpression: new("1")}

	tableModel := func(columns ...model.TableColumnConfig) *model.TableModel {
		return model.TableWithId("test_table", tableId, nil).WithColumnConfigs(columns...)
	}
	modelSnowflakeNormalized := tableModel(id, doubleId("id*2"))
	modelChangedExpression := tableModel(id, doubleId("ID * 3"))
	modelWithAddedColumn := tableModel(id, doubleId("ID * 3"), plusOne)
	modelWithDefault := tableModel(id, doubleId("ID * 3"), invalid)
	ref := modelSnowflakeNormalized.ResourceReference()

	noDriftAfterApply := []plancheck.PlanCheck{
		plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			// Snowflake reports the expression as ID * 2, which does not plan a change.
			{
				Config:           accconfig.FromModels(t, modelSnowflakeNormalized),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: noDriftAfterApply},
				Check: assertThat(t,
					resourceassert.TableResource(t, ref).
						HasColumns(
							resourceassert.ExpectedTableColumn{Name: "ID", Type: "NUMBER(38,0)"},
							resourceassert.ExpectedTableColumn{Name: "DOUBLE_ID", Type: "NUMBER(38,0)", Comment: "computed", Virtual: "ID * 2"},
						),
				),
			},
			{
				ResourceName: ref,
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(tableId.FullyQualifiedName(), "column.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(tableId.FullyQualifiedName(), "column.1.as", "ID * 2"),
				),
			},
			// Changing the expression drops the column and adds it again.
			{
				Config: accconfig.FromModels(t, modelChangedExpression),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
					PostApplyPostRefresh: noDriftAfterApply,
				},
				Check: assertThat(t,
					resourceassert.TableResource(t, ref).
						HasColumns(
							resourceassert.ExpectedTableColumn{Name: "ID", Type: "NUMBER(38,0)"},
							resourceassert.ExpectedTableColumn{Name: "DOUBLE_ID", Type: "NUMBER(38,0)", Comment: "computed", Virtual: "ID * 3"},
						),
				),
			},
			// Adding a virtual column uses ALTER TABLE ADD COLUMN.
			{
				Config: accconfig.FromModels(t, modelWithAddedColumn),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
					PostApplyPostRefresh: noDriftAfterApply,
				},
				Check: assertThat(t,
					resourceassert.TableResource(t, ref).
						HasColumns(
							resourceassert.ExpectedTableColumn{Name: "ID", Type: "NUMBER(38,0)"},
							resourceassert.ExpectedTableColumn{Name: "DOUBLE_ID", Type: "NUMBER(38,0)", Comment: "computed", Virtual: "ID * 3"},
							resourceassert.ExpectedTableColumn{Name: "PLUS_ONE", Type: "NUMBER(38,0)", Virtual: "ID + 1"},
						),
				),
			},
			// A virtual column cannot have a default.
			{
				Config:      accconfig.FromModels(t, modelWithDefault),
				ExpectError: regexp.MustCompile("a virtual column \\(as\\) cannot have default or identity"),
			},
		},
	})
}

// TestAcc_Table_migrateFromV2_23_0 checks the state upgrader.
func TestAcc_Table_migrateFromV2_23_0(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	sequenceId, sequenceCleanup := testClient().Sequence.Create(t)
	t.Cleanup(sequenceCleanup)
	maskingPolicy, maskingPolicyCleanup := testClient().MaskingPolicy.CreateMaskingPolicyIdentity(t, testdatatypes.DataTypeVarchar)
	t.Cleanup(maskingPolicyCleanup)

	providerModel := providermodel.SnowflakeProvider().WithPreviewFeaturesEnabled(string(previewfeatures.TableResource))
	tableModelV2_23_0 := model.TableWithId("test", id, tableColumns).WithComment(comment).
		WithColumnValue(tableMigrationV2_23_0Columns(sequenceId, maskingPolicy.ID()))
	// v2.23.0 did not read virtual columns, so the virtual column is added outside of Terraform before the upgrade.
	virtualColumn := tableMigrationColumn("VIRTUAL_COLUMN", "NUMBER(38,0)", map[string]tfconfig.Variable{"as": tfconfig.StringVariable("ID * 2")})
	tableModel := model.TableWithId("test", id, tableColumns).WithComment(comment).
		WithColumnValue(tableMigrationCurrentColumns(sequenceId, maskingPolicy.ID(), virtualColumn))
	tableModelWithoutVirtualColumn := model.TableWithId("test", id, tableColumns).WithComment(comment).
		WithColumnValue(tableMigrationCurrentColumns(sequenceId, maskingPolicy.ID()))
	ref := tableModel.ResourceReference()
	// DESCRIBE TABLE returns the sequence name without the double quotes, unlike its database and schema.
	sequenceNextval := fmt.Sprintf(`"%s"."%s".%s.NEXTVAL`, sequenceId.DatabaseName(), sequenceId.SchemaName(), sequenceId.Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			// Create - the previous provider with the pipe-separated id. The unquoted masking policy name plans an update
			// in v2.23.0 (https://github.com/snowflakedb/terraform-provider-snowflake/issues/4679).
			{
				ExternalProviders:  ExternalProviderWithExactVersion("2.23.0"),
				ExpectNonEmptyPlan: true,
				Config:             accconfig.FromModels(t, providerModel, tableModelV2_23_0),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(ref, "id", helpers.EncodeSnowflakeID(id))),
				),
			},
			// Upgrade - the state upgraders rewrite the id and the column fields, and Read now reports the virtual column
			{
				PreConfig: func() {
					testClient().Table.AlterWithRequest(t, sdk.NewAlterTableRequest(id).WithColumnAction(*sdk.NewTableColumnActionRequest().WithAdd(
						*sdk.NewTableColumnAddActionRequest("VIRTUAL_COLUMN", testdatatypes.DataTypeNumber_38_0).WithVirtual(*sdk.NewTableColumnVirtualRequest("(ID * 2)")),
					)))
				},
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
					},
				},
				Config: accconfig.FromModels(t, providerModel, tableModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(ref, "id", id.FullyQualifiedName())),
					resourceassert.TableResource(t, ref).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(comment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceassert.TableResource(t, ref).
						HasColumns(
							resourceassert.ExpectedTableColumn{Name: "ID", Type: "NUMBER(38,0)"},
							resourceassert.ExpectedTableColumn{Name: "TEXT_DEFAULT", Type: "VARCHAR(16)", DefaultExpression: "'hello'"},
							resourceassert.ExpectedTableColumn{Name: "BOOLEAN_DEFAULT", Type: "BOOLEAN", DefaultExpression: "FALSE"},
							resourceassert.ExpectedTableColumn{Name: "EXPRESSION_DEFAULT", Type: "TIMESTAMP_NTZ(9)", DefaultExpression: "CURRENT_TIMESTAMP()"},
							resourceassert.ExpectedTableColumn{Name: "SEQUENCE_DEFAULT", Type: "NUMBER(38,0)", DefaultExpression: sequenceNextval},
							resourceassert.ExpectedTableColumn{Name: "MASKED", Type: "VARCHAR(16)", MaskingPolicy: new(maskingPolicy.ID()), MaskingPolicyUsing: []string{"MASKED"}},
							resourceassert.ExpectedTableColumn{Name: "IDENTITY_COLUMN", Type: "NUMBER(38,0)", Identity: &resourceassert.ExpectedTableColumnIdentity{StartNum: 2, StepNum: 3}},
							resourceassert.ExpectedTableColumn{Name: "REQUIRED", Type: "VARCHAR(16)", NotNull: true},
							resourceassert.ExpectedTableColumn{Name: "VIRTUAL_COLUMN", Type: "NUMBER(38,0)", Virtual: "ID * 2"},
						),
				),
			},
			// A virtual column that is not in the configuration is dropped.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, providerModel, tableModelWithoutVirtualColumn),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(ref, "column.#", "8")),
				),
			},
		},
	})
}

// tableMigrationV2_23_0Columns and tableMigrationCurrentColumns hold the same columns in the v2.23.0 schema and in
// the current one. Every v2.24.0 change of the column fields adds its columns to both.
func tableMigrationV2_23_0Columns(sequenceId sdk.SchemaObjectIdentifier, maskingPolicyId sdk.SchemaObjectIdentifier) tfconfig.Variable {
	oldDefault := func(kind string, value string) map[string]tfconfig.Variable {
		return map[string]tfconfig.Variable{"default": tfconfig.ObjectVariable(map[string]tfconfig.Variable{kind: tfconfig.StringVariable(value)})}
	}
	return tfconfig.ListVariable(
		tableMigrationColumn("ID", "NUMBER(38,0)", nil),
		tableMigrationColumn("TEXT_DEFAULT", "VARCHAR(16)", oldDefault("constant", "hello")),
		tableMigrationColumn("BOOLEAN_DEFAULT", "BOOLEAN", oldDefault("constant", "FALSE")),
		tableMigrationColumn("EXPRESSION_DEFAULT", "TIMESTAMP_NTZ(9)", oldDefault("expression", "CURRENT_TIMESTAMP()")),
		tableMigrationColumn("SEQUENCE_DEFAULT", "NUMBER(38,0)", oldDefault("sequence", sequenceId.FullyQualifiedName())),
		tableMigrationColumn("MASKED", "VARCHAR(16)", map[string]tfconfig.Variable{
			"masking_policy": tfconfig.StringVariable(fmt.Sprintf("%s.%s.%s", maskingPolicyId.DatabaseName(), maskingPolicyId.SchemaName(), maskingPolicyId.Name())),
		}),
		tableMigrationColumn("IDENTITY_COLUMN", "NUMBER(38,0)", tableMigrationIdentity()),
		tableMigrationColumn("REQUIRED", "VARCHAR(16)", map[string]tfconfig.Variable{"nullable": tfconfig.BoolVariable(false)}),
	)
}

func tableMigrationIdentity() map[string]tfconfig.Variable {
	return map[string]tfconfig.Variable{"identity": tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"start_num": tfconfig.IntegerVariable(2),
		"step_num":  tfconfig.IntegerVariable(3),
	})}
}

func tableMigrationCurrentColumns(sequenceId sdk.SchemaObjectIdentifier, maskingPolicyId sdk.SchemaObjectIdentifier, extraColumns ...tfconfig.Variable) tfconfig.Variable {
	expression := func(value string) map[string]tfconfig.Variable {
		return map[string]tfconfig.Variable{"default": tfconfig.ObjectVariable(map[string]tfconfig.Variable{"expression": tfconfig.StringVariable(value)})}
	}
	return tfconfig.ListVariable(append([]tfconfig.Variable{
		tableMigrationColumn("ID", "NUMBER(38,0)", nil),
		tableMigrationColumn("TEXT_DEFAULT", "VARCHAR(16)", expression("'hello'")),
		tableMigrationColumn("BOOLEAN_DEFAULT", "BOOLEAN", expression("FALSE")),
		tableMigrationColumn("EXPRESSION_DEFAULT", "TIMESTAMP_NTZ(9)", expression("CURRENT_TIMESTAMP()")),
		tableMigrationColumn("SEQUENCE_DEFAULT", "NUMBER(38,0)", expression(sequenceId.FullyQualifiedName()+".NEXTVAL")),
		tableMigrationColumn("MASKED", "VARCHAR(16)", map[string]tfconfig.Variable{
			"masking_policy": tfconfig.ObjectVariable(map[string]tfconfig.Variable{"policy_name": tfconfig.StringVariable(maskingPolicyId.FullyQualifiedName())}),
		}),
		tableMigrationColumn("IDENTITY_COLUMN", "NUMBER(38,0)", tableMigrationIdentity()),
		tableMigrationColumn("REQUIRED", "VARCHAR(16)", map[string]tfconfig.Variable{"not_null": tfconfig.BoolVariable(true)}),
	}, extraColumns...)...)
}

func tableMigrationColumn(name string, columnType string, fields map[string]tfconfig.Variable) tfconfig.Variable {
	column := map[string]tfconfig.Variable{
		"name": tfconfig.StringVariable(name),
		"type": tfconfig.StringVariable(columnType),
	}
	maps.Copy(column, fields)
	return tfconfig.ObjectVariable(column)
}

func TestAcc_Table_migrateFromVersion_0_94_1(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	providerConfig := providermodel.V097CompatibleProviderConfig(t)

	resourceName := "snowflake_table.test_table"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigWithServiceUserPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            providerConfig + tableConfig(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tableId)),
					resource.TestCheckResourceAttr(resourceName, "name", tableId.Name()),
					resource.TestCheckResourceAttr(resourceName, "qualified_name", tableId.FullyQualifiedName()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   tableConfig(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", tableId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "name", tableId.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", tableId.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceName, "qualified_name"),
				),
			},
		},
	})
}

func TestAcc_Table_SuppressQuotingOnDefaultSequence_issue2644(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, database.ID())
	t.Cleanup(schemaCleanup)

	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	providerConfig := providermodel.V097CompatibleProviderConfig(t)

	resourceName := "snowflake_table.test_table"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { SetV097CompatibleConfigWithServiceUserPathEnv(t) },
				ExternalProviders:  ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             providerConfig + tableConfigWithSequence(tableId, fmt.Sprintf(`sequence = "%s.%s.%s"`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   tableConfigWithSequence(tableId, fmt.Sprintf(`expression = %q`, tableId.FullyQualifiedName()+".NEXTVAL")),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "column.0.default.0.expression", regexp.MustCompile(`(?i)\.NEXTVAL$`)),
				),
			},
		},
	})
}

func tableConfigWithSequence(tableId sdk.SchemaObjectIdentifier, defaultValue string) string {
	return fmt.Sprintf(`
resource "snowflake_sequence" "test_sequence" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
}

resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
	data_retention_time_in_days = 1
	comment  = "Terraform acceptance test"
	column {
		name = "column1"
		type = "NUMBER"
		default {
			%[4]s
		}
	}
	depends_on = [snowflake_sequence.test_sequence]
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), defaultValue)
}

func TestAcc_Table_issue3007_textColumn(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()

	defaultVarchar := fmt.Sprintf("VARCHAR(%d)", datatypes.DefaultVarcharLength)

	resourceName := "snowflake_table.test_table"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "VARCHAR(3)"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.0.type", "NUMBER(11,2)"),
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "VARCHAR(3)"),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "VARCHAR(256)"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(resourceName, "column.1.type", tfjson.ActionUpdate, sdk.String("VARCHAR(3)"), sdk.String("VARCHAR(256)")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "VARCHAR(256)"),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "VARCHAR"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(resourceName, "column.1.type", tfjson.ActionUpdate, sdk.String("VARCHAR(256)"), sdk.String("VARCHAR")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultVarchar),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, defaultVarchar),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultVarchar),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "text"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultVarchar),
				),
			},
		},
	})
}

func TestAcc_Table_issue3007_numberColumn(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()

	defaultNumber := fmt.Sprintf("NUMBER(%d,%d)", datatypes.DefaultNumberPrecision, datatypes.DefaultNumberScale)

	resourceName := "snowflake_table.test_table"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "NUMBER"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.0.type", "NUMBER(11,2)"),
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "NUMBER(38,0)"),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "NUMBER(11)"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(resourceName, "column.1.type", tfjson.ActionUpdate, sdk.String("NUMBER(38,0)"), sdk.String("NUMBER(11)")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.0.type", "NUMBER(11,2)"),
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "NUMBER(11,0)"),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "NUMBER"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(resourceName, "column.1.type", tfjson.ActionUpdate, sdk.String("NUMBER(11,0)"), sdk.String("NUMBER")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultNumber),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, defaultNumber),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultNumber),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "decimal"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultNumber),
				),
			},
			// Snowflake cannot change the scale of a NUMBER column in place.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "NUMBER(38,2)"),
				ExpectError:              regexp.MustCompile("changing the scale of a number is not supported"),
			},
		},
	})
}

func tableConfigIssue3007(tableId sdk.SchemaObjectIdentifier, dataType string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
    database = "%[1]s"
    schema   = "%[2]s"
    name     = "%[3]s"
    comment  = "Issue 3007 confirmation"

    column {
        name = "ID"
        type = "NUMBER(11,2)"
    }

    column {
        name = "SOME_COLUMN"
        type = "%[4]s"
    }
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), dataType)
}

func TestAcc_Table_SchemaRemovedExternally(t *testing.T) {
	schema, schemaCleanup := testClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)
	tableId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetLegacyConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("1.0.5"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_table.test_table", plancheck.ResourceActionCreate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_table.test_table", plancheck.ResourceActionNoop),
					},
				},
				Config: tableConfig(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "id", helpers.EncodeSnowflakeID(tableId)),
				),
			},
			{
				ExternalProviders: ExternalProviderWithExactVersion("1.0.5"),
				PreConfig:         func() { schemaCleanup() },
				Config:            tableConfig(tableId),
				ExpectError:       regexp.MustCompile("object does not exist or not authorized"),
			},
			// The New version removes table from the state which ends up with create operation
			// and the "error creating table" error (because of the missing underlying schema).
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   tableConfig(tableId),
				ExpectError:              regexp.MustCompile("object does not exist or not authorized"),
			},
		},
	})
}

func TestAcc_Table_Decfloat(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	argName := "A"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfigDecfloat(tableId, argName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", tableId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "id"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", argName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "DECFLOAT(38)"),
				),
			},
		},
	})
}

func tableConfigDecfloat(tableId sdk.SchemaObjectIdentifier, argName string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
		name = "id"
		type = "NUMBER"
	}

	column {
		name = "%s"
		type = "DECFLOAT"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), argName)
}

func TestAcc_Table_AddColumnWithConstantDefault_issue4730(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfigIssue4730SingleColumn(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "COLUMN_1"),
				),
			},
			{
				Config: tableConfigIssue4730AddColumnWithConstantDefault(tableId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", tableId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "COLUMN_2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.default.0.expression", "FALSE"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.name", "column_3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.default.0.expression", "FALSE"),
				),
			},
		},
	})
}

func tableConfigIssue4730SingleColumn(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
		name = "COLUMN_1"
		type = "VARCHAR(16)"
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func tableConfigIssue4730AddColumnWithConstantDefault(tableId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"

	column {
		name = "COLUMN_1"
		type = "VARCHAR(16)"
	}
	column {
		name = "COLUMN_2"
		type = "BOOLEAN"
		default {
			expression = "FALSE"
		}
	}
	column {
		name = "column_3"
		type = "BOOLEAN"
		default {
			expression = "FALSE"
		}
	}
}
`, tableId.DatabaseName(), tableId.SchemaName(), tableId.Name())
}

func TestAcc_Table_Transient(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}

	tableModelTransient := model.TableWithId("test_table", tableId, columns).
		WithIsTransient("true")
	tableModelPermanent := model.TableWithId("test_table", tableId, columns).
		WithIsTransient("false")
	tableModelWithoutTransient := model.TableWithId("test_table", tableId, columns)
	ref := tableModelTransient.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tableModelTransient),
				Check: assertThat(
					t,
					resourceassert.TableResource(t, ref).
						HasIsTransient("true"),
					resourceshowoutputassert.TableShowOutput(t, ref).
						HasKind("TRANSIENT"),
				),
			},
			{
				Config: accconfig.FromModels(t, tableModelPermanent),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(
					t,
					resourceassert.TableResource(t, ref).
						HasIsTransient("false"),
					resourceshowoutputassert.TableShowOutput(t, ref).
						HasKind("TABLE"),
				),
			},
			{
				Config: accconfig.FromModels(t, tableModelWithoutTransient),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// is_transient forces a new resource, so removing it from the configuration replaces the table.
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
		},
	})
}

func TestAcc_Table_BasicUseCase_Policies(t *testing.T) {
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()

	rowAccessPolicy, rowAccessPolicyCleanup := testClient().RowAccessPolicy.CreateRowAccessPolicyWithDataType(t, testdatatypes.DataTypeNumber)
	t.Cleanup(rowAccessPolicyCleanup)
	aggregationPolicy, aggregationPolicyCleanup := testClient().AggregationPolicy.CreateAggregationPolicy(t)
	t.Cleanup(aggregationPolicyCleanup)
	joinPolicy, joinPolicyCleanup := testClient().JoinPolicy.Create(t)
	t.Cleanup(joinPolicyCleanup)

	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeNumber},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}
	tableWithoutPolicies := model.TableWithId("test_table", tableId, columns)
	tableWithPolicies := model.TableWithId("test_table", tableId, columns).
		WithRowAccessPolicy(rowAccessPolicy.ID(), "ID").
		WithAggregationPolicy(aggregationPolicy, "ID").
		WithJoinPolicy(joinPolicy, "ID")
	ref := tableWithoutPolicies.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tableWithPolicies),
				Check: assertThat(
					t,
					resourceassert.TableResource(t, ref).
						HasRowAccessPolicy(rowAccessPolicy.ID(), "ID").
						HasAggregationPolicy(aggregationPolicy, "ID").
						HasJoinPolicy(joinPolicy, "ID"),
				),
			},
			{
				PreConfig: func() {
					testClient().Table.AlterWithRequest(t, sdk.NewAlterTableRequest(tableId).WithDropAllRowAccessPolicies(true))
					testClient().Table.AlterWithRequest(t, sdk.NewAlterTableRequest(tableId).WithUnsetAggregationPolicy(*sdk.NewTableUnsetAggregationPolicyRequest()))
					testClient().Table.AlterWithRequest(t, sdk.NewAlterTableRequest(tableId).WithUnsetJoinPolicy(*sdk.NewTableUnsetJoinPolicyRequest()))
				},
				Config: accconfig.FromModels(t, tableWithPolicies),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(
					t,
					resourceassert.TableResource(t, ref).
						HasRowAccessPolicy(rowAccessPolicy.ID(), "ID").
						HasAggregationPolicy(aggregationPolicy, "ID").
						HasJoinPolicy(joinPolicy, "ID"),
				),
			},
			{
				Config: accconfig.FromModels(t, tableWithoutPolicies),
				Check: assertThat(
					t,
					resourceassert.TableResource(t, ref).
						HasRowAccessPolicyEmpty().
						HasAggregationPolicyEmpty().
						HasJoinPolicyEmpty(),
				),
			},
			{
				Config: accconfig.FromModels(t, tableWithPolicies),
				Check: assertThat(
					t,
					resourceassert.TableResource(t, ref).
						HasRowAccessPolicy(rowAccessPolicy.ID(), "ID").
						HasAggregationPolicy(aggregationPolicy, "ID").
						HasJoinPolicy(joinPolicy, "ID"),
				),
			},
		},
	})
}

func TestAcc_Table_SchemaEvolutionErrorLoggingRowTimestamp(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}

	modelBasic := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns)
	modelSet := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).
		WithEnableSchemaEvolution("true").
		WithErrorLogging("true").
		WithRowTimestamp("true")
	modelChanged := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).
		WithEnableSchemaEvolution("false").
		WithErrorLogging("false").
		WithRowTimestamp("false")
	ref := modelBasic.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelSet),
				Check: assertThat(
					t,
					resourceassert.TableResource(t, ref).
						HasEnableSchemaEvolution("true").
						HasErrorLogging("true").
						HasRowTimestamp("true"),
					resourceshowoutputassert.TableShowOutput(t, ref).
						HasEnableSchemaEvolution(true).
						HasErrorLogging(true).
						HasRowTimestamp(true),
				),
			},
			{
				Config: accconfig.FromModels(t, modelChanged),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(
					t,
					resourceassert.TableResource(t, ref).
						HasEnableSchemaEvolution("false").
						HasErrorLogging("false").
						HasRowTimestamp("false"),
					resourceshowoutputassert.TableShowOutput(t, ref).
						HasEnableSchemaEvolution(false).
						HasErrorLogging(false).
						HasRowTimestamp(false),
				),
			},
			{
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(
					t,
					resourceassert.TableResource(t, ref).
						HasEnableSchemaEvolutionString("default").
						HasErrorLoggingString("default").
						HasRowTimestampString("default"),
					resourceshowoutputassert.TableShowOutput(t, ref).
						HasEnableSchemaEvolution(false).
						HasErrorLogging(false).
						HasRowTimestamp(false),
				),
			},
		},
	})
}

func TestAcc_Table_ErrorLoggingAndCommentInOneApply(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0},
	}

	modelBasic := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns)
	// Snowflake silently ignores the other properties of a SET that enables ERROR_LOGGING,
	// so the provider has to apply error_logging separately from the comment.
	modelBoth := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).
		WithErrorLogging("true").
		WithComment(comment)
	ref := modelBasic.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelBasic),
			},
			{
				Config: accconfig.FromModels(t, modelBoth),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(
					t,
					resourceassert.TableResource(t, ref).
						HasErrorLogging("true").
						HasComment(comment),
					resourceshowoutputassert.TableShowOutput(t, ref).
						HasErrorLogging(true).
						HasComment(comment),
				),
			},
		},
	})
}

func TestAcc_Table_MaxDataExtensionTimeAndDefaultDdlCollation(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0},
		{Name: "NAME", Type: testdatatypes.DataTypeVarchar},
	}

	modelSet := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).
		WithMaxDataExtensionTimeInDays(20).
		WithDefaultDdlCollation("en-ci")
	modelChanged := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).
		WithMaxDataExtensionTimeInDays(30).
		WithDefaultDdlCollation("en-cs")
	modelInherit := model.TableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns)
	ref := modelSet.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelSet),
				Check: assertThat(
					t,
					resourceassert.TableResource(t, ref).
						HasMaxDataExtensionTimeInDays(20).
						HasDefaultDdlCollation("en-ci"),
					resourceparametersassert.TableResourceParameters(t, ref).
						HasMaxDataExtensionTimeInDays(20).
						HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeTable).
						HasDefaultDdlCollation("en-ci").
						HasDefaultDdlCollationLevel(sdk.ParameterTypeTable),
				),
			},
			{
				Config: accconfig.FromModels(t, modelChanged),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(
					t,
					resourceassert.TableResource(t, ref).
						HasMaxDataExtensionTimeInDays(30).
						HasDefaultDdlCollation("en-cs"),
					resourceparametersassert.TableResourceParameters(t, ref).
						HasMaxDataExtensionTimeInDays(30).
						HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeTable).
						HasDefaultDdlCollation("en-cs").
						HasDefaultDdlCollationLevel(sdk.ParameterTypeTable),
				),
			},
			{
				Config: accconfig.FromModels(t, modelInherit),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(
					t,
					resourceassert.TableResource(t, ref).
						HasMaxDataExtensionTimeInDaysNotEmpty().
						HasDefaultDdlCollationEmpty(),
					resourceparametersassert.TableResourceParameters(t, ref).
						HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeDatabase).
						HasDefaultDdlCollationLevel(sdk.ParameterTypeSnowflakeDefault),
				),
			},
		},
	})
}
