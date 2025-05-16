//go:build !account_level_tests

package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ProcedureSql_InlineBasic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)
	idWithChangedNameButTheSameDataType := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	definition := acc.TestClient().Procedure.SampleSqlDefinitionWithArgument(t)

	procedureModel := model.ProcedureSqlBasicInline("w", id, dataType, definition).
		WithArgument(argName, dataType)
	procedureModelRenamed := model.ProcedureSqlBasicInline("w", idWithChangedNameButTheSameDataType, dataType, definition).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.ProcedureSql),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(t,
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
					acc.TestClient().Procedure.DropProcedureFunc(t, id)()
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(procedureModel.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(t,
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
				ImportStateCheck: assertThatImport(t,
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
				Check: assertThat(t,
					resourceassert.ProcedureSqlResource(t, procedureModelRenamed.ResourceReference()).
						HasNameString(idWithChangedNameButTheSameDataType.Name()).
						HasFullyQualifiedNameString(idWithChangedNameButTheSameDataType.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_ProcedureSql_InlineFull(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	definition := acc.TestClient().Procedure.SampleSqlDefinitionWithArgument(t)

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
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.ProcedureSql),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, procedureModel),
				Check: assertThat(t,
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
				ImportStateCheck: assertThatImport(t,
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
				Check: assertThat(t,
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
