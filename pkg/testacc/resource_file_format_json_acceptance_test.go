//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FileFormatJson_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()

	basicModel := model.FileFormatJson("test", id.DatabaseName(), id.SchemaName(), id.Name())
	ref := basicModel.ResourceReference()

	completeModel := model.FileFormatJson("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression("GZIP").
		WithDateFormat("YYYY-MM-DD").
		WithTimeFormat("HH24:MI:SS").
		WithTimestampFormat("YYYY-MM-DD HH24:MI:SS.FF3").
		WithBinaryFormat("BASE64").
		WithTrimSpace("true").
		WithMultiLine("false").
		WithFileExtension(".json").
		WithEnableOctal("true").
		WithAllowDuplicate("true").
		WithStripOuterArray("true").
		WithStripNullValues("true").
		// ReplaceInvalidCharacters is incompatible with IgnoreUtf8Errors
		WithIgnoreUtf8Errors("true").
		WithSkipByteOrderMark("false").
		WithComment(comment)

	alteredModel := model.FileFormatJson("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression("BZ2").
		WithDateFormat("MM-DD-YYYY").
		WithTimeFormat("HH24:MI").
		WithTimestampFormat("YYYY-MM-DD HH24:MI:SS.FF6").
		WithBinaryFormat("UTF8").
		WithTrimSpace("false").
		WithMultiLine("true").
		WithFileExtension(".jsonl").
		WithEnableOctal("false").
		WithAllowDuplicate("false").
		WithStripOuterArray("false").
		WithStripNullValues("false").
		// ReplaceInvalidCharacters is incompatible with IgnoreUtf8Errors, so only the latter is altered here
		WithIgnoreUtf8Errors("false").
		WithSkipByteOrderMark("true").
		WithComment(externalComment)

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatJsonResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasTrimSpace(r.BooleanDefault).
			HasMultiLine(r.BooleanDefault).
			HasEnableOctal(r.BooleanDefault).
			HasAllowDuplicate(r.BooleanDefault).
			HasStripOuterArray(r.BooleanDefault).
			HasStripNullValues(r.BooleanDefault).
			HasReplaceInvalidCharacters(r.BooleanDefault).
			HasIgnoreUtf8Errors(r.BooleanDefault).
			HasSkipByteOrderMark(r.BooleanDefault),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeJson).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(""),
		resourceshowoutputassert.FileFormatJsonDescribeOutput(t, ref).
			HasId(id).
			HasCompression("AUTO").
			HasBinaryFormat(string(sdk.BinaryFormatHex)).
			HasTrimSpace(false).
			HasMultiLine(true).
			HasEnableOctal(false).
			HasAllowDuplicate(false).
			HasStripOuterArray(false).
			HasStripNullValues(false).
			HasReplaceInvalidCharacters(false).
			HasIgnoreUtf8Errors(false).
			HasSkipByteOrderMark(true),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatJsonResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCompressionString(string(sdk.CsvCompressionGzip)).
			HasDateFormatString("YYYY-MM-DD").
			HasTimeFormatString("HH24:MI:SS").
			HasTimestampFormatString("YYYY-MM-DD HH24:MI:SS.FF3").
			HasBinaryFormatString(string(sdk.BinaryFormatBase64)).
			HasTrimSpace(r.BooleanTrue).
			HasMultiLine(r.BooleanFalse).
			HasFileExtensionString(".json").
			HasEnableOctal(r.BooleanTrue).
			HasAllowDuplicate(r.BooleanTrue).
			HasStripOuterArray(r.BooleanTrue).
			HasStripNullValues(r.BooleanTrue).
			HasReplaceInvalidCharacters(r.BooleanDefault).
			HasIgnoreUtf8Errors(r.BooleanTrue).
			HasSkipByteOrderMark(r.BooleanFalse).
			HasCommentString(comment),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeJson).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment),
		resourceshowoutputassert.FileFormatJsonDescribeOutput(t, ref).
			HasId(id).
			HasCompression("GZIP").
			HasDateFormat("YYYY-MM-DD").
			HasTimeFormat("HH24:MI:SS").
			HasTimestampFormat("YYYY-MM-DD HH24:MI:SS.FF3").
			HasBinaryFormat("BASE64").
			HasTrimSpace(true).
			HasMultiLine(false).
			HasFileExtension(".json").
			HasEnableOctal(true).
			HasAllowDuplicate(true).
			HasStripOuterArray(true).
			HasStripNullValues(true).
			HasReplaceInvalidCharacters(false).
			HasIgnoreUtf8Errors(true).
			HasSkipByteOrderMark(false),
	}

	alteredAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatJsonResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCompressionString("BZ2").
			HasDateFormatString("MM-DD-YYYY").
			HasTimeFormatString("HH24:MI").
			HasTimestampFormatString("YYYY-MM-DD HH24:MI:SS.FF6").
			HasBinaryFormatString("UTF8").
			HasTrimSpace(r.BooleanFalse).
			HasMultiLine(r.BooleanTrue).
			HasFileExtensionString(".jsonl").
			HasEnableOctal(r.BooleanFalse).
			HasAllowDuplicate(r.BooleanFalse).
			HasStripOuterArray(r.BooleanFalse).
			HasStripNullValues(r.BooleanFalse).
			HasReplaceInvalidCharacters(r.BooleanDefault).
			HasIgnoreUtf8Errors(r.BooleanFalse).
			HasSkipByteOrderMark(r.BooleanTrue).
			HasCommentString(externalComment),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeJson).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(externalComment),
		resourceshowoutputassert.FileFormatJsonDescribeOutput(t, ref).
			HasId(id).
			HasCompression("BZ2").
			HasDateFormat("MM-DD-YYYY").
			HasTimeFormat("HH24:MI").
			HasTimestampFormat("YYYY-MM-DD HH24:MI:SS.FF6").
			HasBinaryFormat("UTF8").
			HasTrimSpace(false).
			HasMultiLine(true).
			HasFileExtension(".jsonl").
			HasEnableOctal(false).
			HasAllowDuplicate(false).
			HasStripOuterArray(false).
			HasStripNullValues(false).
			HasReplaceInvalidCharacters(false).
			HasIgnoreUtf8Errors(false).
			HasSkipByteOrderMark(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatJson),
		Steps: []resource.TestStep{
			// create with only required fields
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basicModel),
				Check:  assertThat(t, basicAssertions...),
			},
			// import
			{
				Config:            config.FromModels(t, basicModel),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// set all optional fields
			{
				Config: config.FromModels(t, completeModel),
				Check:  assertThat(t, completeAssertions...),
			},
			// import with all attributes
			{
				Config:            config.FromModels(t, completeModel),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// alter all optional fields (non-recreating change)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, alteredModel),
				Check:  assertThat(t, alteredAssertions...),
			},
			// external non-type change is detected and corrected back to the config value with an update (non-recreating change)
			{
				PreConfig: func() {
					testClient().FileFormat.CreateJsonWithRequest(t, id, sdk.NewCreateJsonFileFormatRequest(id).WithOrReplace(true))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, alteredModel),
				Check:  assertThat(t, alteredAssertions...),
			},
			// external type change is detected (the object is recreated in Snowflake as a different file format type)
			{
				PreConfig: func() {
					testClient().FileFormat.CreateCsvWithRequest(t, id, sdk.NewCreateCsvFileFormatRequest(id).WithOrReplace(true))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, alteredModel),
				Check:  assertThat(t, alteredAssertions...),
			},
			// unset optional fields
			{
				Config: config.FromModels(t, basicModel),
				Check:  assertThat(t, basicAssertions...),
			},
		},
	})
}

func TestAcc_FileFormatJson_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	completeModel := model.FileFormatJsonWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression(string(sdk.CsvCompressionGzip)).
		WithDateFormat("YYYY-MM-DD").
		WithTimeFormat("HH24:MI:SS").
		WithTimestampFormat("YYYY-MM-DD HH24:MI:SS.FF3").
		WithBinaryFormat(string(sdk.BinaryFormatBase64)).
		WithTrimSpace("true").
		WithMultiLine("false").
		WithFileExtension(".json").
		WithEnableOctal("true").
		WithAllowDuplicate("true").
		WithStripOuterArray("true").
		// ReplaceInvalidCharacters is incompatible with IgnoreUtf8Errors
		WithStripNullValues("true").
		WithIgnoreUtf8Errors("true").
		WithSkipByteOrderMark("false").
		WithComment(comment)
	modelWithReplaceInvalidCharacters := model.FileFormatJsonWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name()).
		WithReplaceInvalidCharacters("true")
	ref := completeModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatJson),
		Steps: []resource.TestStep{
			// create with all fields set
			{
				Config: config.FromModels(t, completeModel),
				Check: assertThat(
					t,
					resourceassert.FileFormatJsonResource(t, ref).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCompressionString("GZIP").
						HasDateFormatString("YYYY-MM-DD").
						HasTimeFormatString("HH24:MI:SS").
						HasTimestampFormatString("YYYY-MM-DD HH24:MI:SS.FF3").
						HasBinaryFormatString("BASE64").
						HasTrimSpace("true").
						HasMultiLine("false").
						HasFileExtensionString(".json").
						HasEnableOctal("true").
						HasAllowDuplicate("true").
						HasStripOuterArray("true").
						HasStripNullValues("true").
						HasReplaceInvalidCharacters("default").
						HasIgnoreUtf8Errors("true").
						HasSkipByteOrderMark("false").
						HasCommentString(comment),
					resourceshowoutputassert.FileFormatShowOutput(t, ref).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.FileFormatTypeJson).
						HasComment(comment),
					resourceshowoutputassert.FileFormatJsonDescribeOutput(t, ref).
						HasId(id).
						HasCompression("GZIP").
						HasDateFormat("YYYY-MM-DD").
						HasTimeFormat("HH24:MI:SS").
						HasTimestampFormat("YYYY-MM-DD HH24:MI:SS.FF3").
						HasBinaryFormat("BASE64").
						HasTrimSpace(true).
						HasMultiLine(false).
						HasFileExtension(".json").
						HasEnableOctal(true).
						HasAllowDuplicate(true).
						HasStripOuterArray(true).
						HasStripNullValues(true).
						HasReplaceInvalidCharacters(false).
						HasIgnoreUtf8Errors(true).
						HasSkipByteOrderMark(false),
				),
			},
			// import
			{
				Config:            config.FromModels(t, completeModel),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  config.FromModels(t, completeModel),
				Destroy: true,
			},
			{
				Config: config.FromModels(t, modelWithReplaceInvalidCharacters),
				Check: assertThat(
					t,
					resourceassert.FileFormatJsonResource(t, ref).
						HasReplaceInvalidCharacters("true"),
					resourceshowoutputassert.FileFormatJsonDescribeOutput(t, ref).
						HasReplaceInvalidCharacters(true),
				),
			},
		},
	})
}

func TestAcc_FileFormatJson_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	invalidCompression := model.FileFormatJson("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression("INVALID")
	invalidBinaryFormat := model.FileFormatJson("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithBinaryFormat("INVALID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatJson),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, invalidCompression),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid json compression: INVALID`),
			},
			{
				Config:      config.FromModels(t, invalidBinaryFormat),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid binary format: INVALID`),
			},
		},
	})
}
