//go:build !account_level_tests

package testacc

import (
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FunctionJavascript_InlineBasic(t *testing.T) {
	argName := "x"
	dataType := testdatatypes.DataTypeVariant

	id := testClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)
	idWithChangedNameButTheSameDataType := testClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	definition := testClient().Function.SampleJavascriptDefinition(t, argName)

	functionModel := model.FunctionJavascriptInline("test", id, definition, datatypes.VariantLegacyDataType).
		WithArgument(argName, dataType)
	functionModelRenamed := model.FunctionJavascriptInline("test", idWithChangedNameButTheSameDataType, definition, datatypes.VariantLegacyDataType).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FunctionJavascript),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, functionModel),
				Check: assertThat(t,
					resourceassert.FunctionJavascriptResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasReturnTypeString(datatypes.VariantLegacyDataType).
						HasIsSecureString(r.BooleanDefault).
						HasCommentString(sdk.DefaultFunctionComment).
						HasFunctionDefinitionString(definition).
						HasFunctionLanguageString("JAVASCRIPT").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.FunctionShowOutput(t, functionModel.ResourceReference()).
						HasIsSecure(false),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_name", argName)),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_data_type", dataType.ToSql())),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_default_value", "")),
				),
			},
			// REMOVE EXTERNALLY (CHECK RECREATION)
			{
				PreConfig: func() {
					testClient().Function.DropFunctionFunc(t, id)()
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(functionModel.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, functionModel),
				Check: assertThat(t,
					resourceassert.FunctionJavascriptResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()),
				),
			},
			// RENAME
			{
				Config: config.FromModels(t, functionModelRenamed),
				Check: assertThat(t,
					resourceassert.FunctionJavaResource(t, functionModelRenamed.ResourceReference()).
						HasNameString(idWithChangedNameButTheSameDataType.Name()).
						HasFullyQualifiedNameString(idWithChangedNameButTheSameDataType.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_FunctionJavascript_InlineEmptyArgs(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes()
	definition := testClient().Function.SampleJavascriptDefinitionNoArgs(t)
	functionModel := model.FunctionJavascriptInline("test", id, definition, datatypes.VariantLegacyDataType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FunctionJavascript),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, functionModel),
				Check: assertThat(t,
					resourceassert.FunctionJavaResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFunctionDefinitionString(definition).
						HasFunctionLanguageString("JAVASCRIPT").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_FunctionJavascript_InlineFull(t *testing.T) {
	argName := "x"
	dataType := testdatatypes.DataTypeVariant

	id := testClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)
	definition := testClient().Function.SampleJavascriptDefinition(t, argName)
	functionModel := model.FunctionJavascriptInline("test", id, definition, datatypes.VariantLegacyDataType).
		WithIsSecure(r.BooleanFalse).
		WithArgument(argName, dataType).
		WithNullInputBehavior(string(sdk.NullInputBehaviorReturnsNullInput)).
		WithReturnResultsBehavior(string(sdk.ReturnResultsBehaviorVolatile)).
		WithComment("some comment")

	functionModelUpdateWithoutRecreation := model.FunctionJavascriptInline("test", id, definition, datatypes.VariantLegacyDataType).
		WithArgument(argName, dataType).
		WithIsSecure(r.BooleanFalse).
		WithNullInputBehavior(string(sdk.NullInputBehaviorReturnsNullInput)).
		WithReturnResultsBehavior(string(sdk.ReturnResultsBehaviorVolatile)).
		WithComment("some other comment")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.FunctionJavascript),
		Steps: []resource.TestStep{
			// CREATE WITH ALL
			{
				Config: config.FromModels(t, functionModel),
				Check: assertThat(t,
					resourceassert.FunctionJavascriptResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanFalse).
						HasFunctionDefinitionString(definition).
						HasCommentString("some comment").
						HasFunctionLanguageString("JAVASCRIPT").
						HasNullInputBehaviorString(string(sdk.NullInputBehaviorReturnsNullInput)).
						HasReturnResultsBehaviorString(string(sdk.ReturnResultsBehaviorVolatile)).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.FunctionShowOutput(t, functionModel.ResourceReference()).
						HasIsSecure(false),
				),
			},
			// IMPORT
			{
				ResourceName:            functionModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"arguments.0.arg_data_type"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedFunctionJavaResource(t, id.FullyQualifiedName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_name", argName)),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_data_type", "VARIANT")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_default_value", "")),
				),
			},
			// UPDATE WITHOUT RECREATION
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(functionModelUpdateWithoutRecreation.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, functionModelUpdateWithoutRecreation),
				Check: assertThat(t,
					resourceassert.FunctionJavascriptResource(t, functionModelUpdateWithoutRecreation.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanFalse).
						HasFunctionDefinitionString(definition).
						HasCommentString("some other comment").
						HasFunctionLanguageString("JAVASCRIPT").
						HasNullInputBehaviorString(string(sdk.NullInputBehaviorReturnsNullInput)).
						HasReturnResultsBehaviorString(string(sdk.ReturnResultsBehaviorVolatile)).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.FunctionShowOutput(t, functionModelUpdateWithoutRecreation.ResourceReference()).
						HasIsSecure(false),
				),
			},
		},
	})
}
