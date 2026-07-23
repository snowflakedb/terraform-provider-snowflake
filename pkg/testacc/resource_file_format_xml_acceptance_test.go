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

func TestAcc_FileFormatXml_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	renamedId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()

	basicModel := model.FileFormatXml("test", id.DatabaseName(), id.SchemaName(), id.Name())
	ref := basicModel.ResourceReference()

	renamedModel := model.FileFormatXml("test", id.DatabaseName(), id.SchemaName(), renamedId.Name())

	completeModel := model.FileFormatXml("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression("GZIP").
		WithPreserveSpace("true").
		WithStripOuterElement("true").
		WithDisableSnowflakeData("true").
		WithDisableAutoConvert("true").
		// ReplaceInvalidCharacters is incompatible with IgnoreUtf8Errors
		WithIgnoreUtf8Errors("true").
		WithSkipByteOrderMark("false").
		WithComment(comment)

	alteredModel := model.FileFormatXml("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression("BZ2").
		WithPreserveSpace("false").
		WithStripOuterElement("false").
		WithDisableSnowflakeData("false").
		WithDisableAutoConvert("false").
		// ReplaceInvalidCharacters is incompatible with IgnoreUtf8Errors, so only the latter is altered here
		WithIgnoreUtf8Errors("false").
		WithSkipByteOrderMark("true").
		WithComment(externalComment)

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatXmlResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasPreserveSpace(r.BooleanDefault).
			HasStripOuterElement(r.BooleanDefault).
			HasDisableSnowflakeData(r.BooleanDefault).
			HasDisableAutoConvert(r.BooleanDefault).
			HasReplaceInvalidCharacters(r.BooleanDefault).
			HasIgnoreUtf8Errors(r.BooleanDefault).
			HasSkipByteOrderMark(r.BooleanDefault),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeXml).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(""),
		resourceshowoutputassert.FileFormatXmlDescribeOutput(t, ref).
			HasId(id).
			HasCompression("AUTO").
			HasPreserveSpace(false).
			HasStripOuterElement(false).
			HasDisableSnowflakeData(false).
			HasDisableAutoConvert(false).
			HasReplaceInvalidCharacters(false).
			HasIgnoreUtf8Errors(false).
			HasSkipByteOrderMark(true),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatXmlResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCompressionString("GZIP").
			HasPreserveSpace(r.BooleanTrue).
			HasStripOuterElement(r.BooleanTrue).
			HasDisableSnowflakeData(r.BooleanTrue).
			HasDisableAutoConvert(r.BooleanTrue).
			HasReplaceInvalidCharacters(r.BooleanDefault).
			HasIgnoreUtf8Errors(r.BooleanTrue).
			HasSkipByteOrderMark(r.BooleanFalse).
			HasCommentString(comment),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeXml).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment),
		resourceshowoutputassert.FileFormatXmlDescribeOutput(t, ref).
			HasId(id).
			HasCompression("GZIP").
			HasPreserveSpace(true).
			HasStripOuterElement(true).
			HasDisableSnowflakeData(true).
			HasDisableAutoConvert(true).
			HasReplaceInvalidCharacters(false).
			HasIgnoreUtf8Errors(true).
			HasSkipByteOrderMark(false),
	}

	alteredAssertions := []assert.TestCheckFuncProvider{
		resourceassert.FileFormatXmlResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasCompressionString("BZ2").
			HasPreserveSpace(r.BooleanFalse).
			HasStripOuterElement(r.BooleanFalse).
			HasDisableSnowflakeData(r.BooleanFalse).
			HasDisableAutoConvert(r.BooleanFalse).
			HasReplaceInvalidCharacters(r.BooleanDefault).
			HasIgnoreUtf8Errors(r.BooleanFalse).
			HasSkipByteOrderMark(r.BooleanTrue).
			HasCommentString(externalComment),
		resourceshowoutputassert.FileFormatShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasType(sdk.FileFormatTypeXml).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(externalComment),
		resourceshowoutputassert.FileFormatXmlDescribeOutput(t, ref).
			HasId(id).
			HasCompression("BZ2").
			HasPreserveSpace(false).
			HasStripOuterElement(false).
			HasDisableSnowflakeData(false).
			HasDisableAutoConvert(false).
			HasReplaceInvalidCharacters(false).
			HasIgnoreUtf8Errors(false).
			HasSkipByteOrderMark(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatXml),
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
					testClient().FileFormat.CreateXmlWithRequest(t, id, sdk.NewCreateXmlFileFormatRequest(id).WithOrReplace(true))
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
			// rename
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, renamedModel),
				Check: assertThat(
					t,
					resourceassert.FileFormatXmlResource(t, ref).
						HasNameString(renamedId.Name()).
						HasFullyQualifiedNameString(renamedId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_FileFormatXml_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	completeModel := model.FileFormatXmlWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression(string(sdk.XmlCompressionGzip)).
		WithPreserveSpace("true").
		WithStripOuterElement("true").
		WithDisableSnowflakeData("true").
		WithDisableAutoConvert("true").
		// ReplaceInvalidCharacters is incompatible with IgnoreUtf8Errors
		WithIgnoreUtf8Errors("true").
		WithSkipByteOrderMark("false").
		WithComment(comment)
	modelWithReplaceInvalidCharacters := model.FileFormatXmlWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name()).
		WithReplaceInvalidCharacters("true")
	ref := completeModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatXml),
		Steps: []resource.TestStep{
			// create with all fields set
			{
				Config: config.FromModels(t, completeModel),
				Check: assertThat(
					t,
					resourceassert.FileFormatXmlResource(t, ref).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCompressionString("GZIP").
						HasPreserveSpace("true").
						HasStripOuterElement("true").
						HasDisableSnowflakeData("true").
						HasDisableAutoConvert("true").
						HasReplaceInvalidCharacters("default").
						HasIgnoreUtf8Errors("true").
						HasSkipByteOrderMark("false").
						HasCommentString(comment),
					resourceshowoutputassert.FileFormatShowOutput(t, ref).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.FileFormatTypeXml).
						HasComment(comment),
					resourceshowoutputassert.FileFormatXmlDescribeOutput(t, ref).
						HasId(id).
						HasCompression("GZIP").
						HasPreserveSpace(true).
						HasStripOuterElement(true).
						HasDisableSnowflakeData(true).
						HasDisableAutoConvert(true).
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
					resourceassert.FileFormatXmlResource(t, ref).
						HasReplaceInvalidCharacters("true"),
					resourceshowoutputassert.FileFormatXmlDescribeOutput(t, ref).
						HasReplaceInvalidCharacters(true),
				),
			},
		},
	})
}

func TestAcc_FileFormatXml_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	invalidCompression := model.FileFormatXml("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression("INVALID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatXml),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, invalidCompression),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid xml compression: INVALID`),
			},
		},
	})
}
