//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ProcedureSql_InlineBasic(t *testing.T) {
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	id := testClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)
	idWithChangedNameButTheSameDataType := testClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	definition := testClient().Procedure.SampleSqlDefinitionWithArgument(t)

	procedureModel := model.ProcedureSqlBasicInline("w", id, dataType, definition).
		WithArgument(argName, dataType)
	procedureModelRenamed := model.ProcedureSqlBasicInline("w", idWithChangedNameButTheSameDataType, dataType, definition).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: functionsAndProceduresProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ProcedureSql),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(
					t,
					resourceassert.ProcedureSqlResource(t, procedureModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanDefault).
						HasCommentString(sdk.DefaultProcedureComment).
						HasProcedureDefinitionString(definition).
						HasProcedureLanguageString("SQL").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ProcedureShowOutput(t, procedureModel.ResourceReference()).
						HasIsSecure(false),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "arguments.0.arg_name", argName)),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "arguments.0.arg_data_type", dataType.ToSql())),
					assert.Check(resource.TestCheckResourceAttr(procedureModel.ResourceReference(), "arguments.0.arg_default_value", "")),
				),
			},
			// REMOVE EXTERNALLY (CHECK RECREATION)
			{
				PreConfig: func() {
					testClient().Procedure.DropProcedureFunc(t, id)()
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(procedureModel.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(
					t,
					resourceassert.ProcedureSqlResource(t, procedureModel.ResourceReference()).
						HasNameString(id.Name()),
				),
			},
			// IMPORT
			{
				ResourceName:            procedureModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_secure", "arguments.0.arg_data_type", "null_input_behavior", "execute_as"},
				ImportStateCheck: assertThatImport(
					t,
					resourceassert.ImportedProcedureSqlResource(t, id.FullyQualifiedName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_name", argName)),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_data_type", "VARCHAR(16777216)")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_default_value", "")),
				),
			},
			// RENAME
			{
				Config: config.FromModels(t, procedureModelRenamed),
				Check: assertThat(
					t,
					resourceassert.ProcedureSqlResource(t, procedureModelRenamed.ResourceReference()).
						HasNameString(idWithChangedNameButTheSameDataType.Name()).
						HasFullyQualifiedNameString(idWithChangedNameButTheSameDataType.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_ProcedureSql_InlineFull(t *testing.T) {
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	id := testClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	definition := testClient().Procedure.SampleSqlDefinitionWithArgument(t)

	procedureModel := model.ProcedureSqlBasicInline("w", id, dataType, definition).
		WithArgument(argName, dataType).
		WithIsSecure("false").
		WithNullInputBehavior(string(sdk.NullInputBehaviorCalledOnNullInput)).
		WithExecuteAs(string(sdk.ExecuteAsCaller)).
		WithComment("some comment")

	procedureModelUpdateWithoutRecreation := model.ProcedureSqlBasicInline("w", id, dataType, definition).
		WithArgument(argName, dataType).
		WithIsSecure("false").
		WithNullInputBehavior(string(sdk.NullInputBehaviorCalledOnNullInput)).
		WithExecuteAs(string(sdk.ExecuteAsOwner)).
		WithComment("some other comment")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: functionsAndProceduresProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ProcedureSql),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(
					t,
					resourceassert.ProcedureSqlResource(t, procedureModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanFalse).
						HasProcedureDefinitionString(definition).
						HasCommentString("some comment").
						HasProcedureLanguageString("SQL").
						HasExecuteAsString(string(sdk.ExecuteAsCaller)).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ProcedureShowOutput(t, procedureModel.ResourceReference()).
						HasIsSecure(false),
				),
			},
			// IMPORT
			{
				ResourceName:            procedureModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"arguments.0.arg_data_type", "null_input_behavior"},
				ImportStateCheck: assertThatImport(
					t,
					resourceassert.ImportedProcedureSqlResource(t, id.FullyQualifiedName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_name", argName)),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_data_type", "VARCHAR(16777216)")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_default_value", "")),
				),
			},
			// UPDATE WITHOUT RECREATION
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(procedureModelUpdateWithoutRecreation.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, procedureModelUpdateWithoutRecreation),
				Check: assertThat(
					t,
					resourceassert.ProcedureSqlResource(t, procedureModelUpdateWithoutRecreation.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanFalse).
						HasProcedureDefinitionString(definition).
						HasCommentString("some other comment").
						HasProcedureLanguageString("SQL").
						HasExecuteAsString(string(sdk.ExecuteAsOwner)).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ProcedureShowOutput(t, procedureModelUpdateWithoutRecreation.ResourceReference()).
						HasIsSecure(false),
				),
			},
		},
	})
}

func TestAcc_ProcedureSql_DecfloatUnsupported(t *testing.T) {
	argName := "x"
	dataType := testdatatypes.DataTypeDecfloat

	id := testClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	definition := testClient().Procedure.SampleSqlDefinitionWithArgument(t)

	procedureModel := model.ProcedureSqlBasicInline("w", id, dataType, definition).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: functionsAndProceduresProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ProcedureSql),
		Steps: []resource.TestStep{
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(procedureModel.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModels(t, procedureModel),
				ExpectError: regexp.MustCompile(`Language SQL does not support type 'DECFLOAT\(38\)' for argument or return type`),
			},
		},
	})
}

func TestAcc_ProcedureSql_tableReturnTypeWithParametrizedColumnsNonDefaults(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes()

	returnType := "TABLE(O_ERR_CODE NUMBER(24, 2), O_ERR_SEVERITY VARCHAR(150))"
	definition := testClient().Procedure.SampleSqlDefinitionReturningTable(t)

	procedureModel := model.ProcedureSql("test", id.DatabaseName(), id.SchemaName(), id.Name(), definition, returnType)
	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.ProcedureSqlResource))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ProcedureSql),
		Steps: []resource.TestStep{
			// Step 1: v2.15.0 fails because TABLE with parametrized columns cannot be parsed.
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.15.0"),
				Config:            config.FromModels(t, providerModel, procedureModel),
				ExpectError:       regexp.MustCompile(`number NUMBER\(24 could not be parsed, use "NUMBER\(precision, scale\)" format`),
			},
			// Step 2: Current version correctly parses TABLE with parametrized columns.
			{
				ProtoV6ProviderFactories: functionsAndProceduresProviderFactory,
				Config:                   config.FromModels(t, procedureModel),
				Check: assertThat(
					t,
					resourceassert.ProcedureSqlResource(t, procedureModel.ResourceReference()).
						HasNameString(id.Name()).
						HasReturnTypeString(returnType),
				),
			},
		},
	})
}

// When a procedure with TABLE return type using non-parametrized column types was created in
// v2.15.0, the current version must not produce spurious plan drift after migration.
func TestAcc_ProcedureSql_defaultParamsNoDriftAfterMigration(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes()

	returnType := "TABLE(O_ERR_CODE NUMBER, O_ERR_SEVERITY VARCHAR)"
	definition := testClient().Procedure.SampleSqlDefinitionReturningTable(t)

	procedureModel := model.ProcedureSql("test", id.DatabaseName(), id.SchemaName(), id.Name(), definition, returnType)
	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.ProcedureSqlResource))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ProcedureSql),
		Steps: []resource.TestStep{
			// v2.15.0 creates the procedure with NUMBER return type successfully.
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.15.0"),
				Config:            config.FromModels(t, providerModel, procedureModel),
				Check: assertThat(
					t,
					resourceassert.ProcedureSqlResource(t, procedureModel.ResourceReference()).
						HasNameString(id.Name()).
						HasReturnTypeString(returnType),
				),
			},
			// Current version: no drift even though state now contains the full type..
			{
				ProtoV6ProviderFactories: functionsAndProceduresProviderFactory,
				Config:                   config.FromModels(t, procedureModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(
					t,
					resourceassert.ProcedureSqlResource(t, procedureModel.ResourceReference()).
						HasNameString(id.Name()).
						HasReturnTypeString(returnType),
				),
			},
		},
	})
}
