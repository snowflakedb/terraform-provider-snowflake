//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FileFormats_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	idTwo := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	idThree := testClient().Ids.RandomSchemaObjectIdentifier()
	schemaId := testClient().Ids.SchemaId()

	fileFormatModel1 := model.FileFormatCsv("test1", idOne.DatabaseName(), idOne.SchemaName(), idOne.Name())
	fileFormatModel2 := model.FileFormatCsv("test2", idTwo.DatabaseName(), idTwo.SchemaName(), idTwo.Name())
	fileFormatModel3 := model.FileFormatCsv("test3", idThree.DatabaseName(), idThree.SchemaName(), idThree.Name())

	dependsOn := []string{fileFormatModel1.ResourceReference(), fileFormatModel2.ResourceReference(), fileFormatModel3.ResourceReference()}

	fileFormatsLikeFirst := datasourcemodel.FileFormats("test").
		WithLike(idOne.Name()).
		WithInSchema(schemaId).
		WithWithDescribe(false).
		WithDependsOn(dependsOn...)

	fileFormatsLikePrefix := datasourcemodel.FileFormats("test").
		WithLike(prefix + "%").
		WithInSchema(schemaId).
		WithWithDescribe(false).
		WithDependsOn(dependsOn...)

	fileFormatsInDatabase := datasourcemodel.FileFormats("test").
		WithLike(prefix + "%").
		WithInDatabase(schemaId.DatabaseId()).
		WithWithDescribe(false).
		WithDependsOn(dependsOn...)

	fileFormatsNonExisting := datasourcemodel.FileFormats("test").
		WithLike("non-existing").
		WithInSchema(schemaId).
		WithWithDescribe(false).
		WithDependsOn(dependsOn...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatCsv),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, fileFormatModel1, fileFormatModel2, fileFormatModel3, fileFormatsLikeFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fileFormatsLikeFirst.DatasourceReference(), "file_formats.#", "1"),
					resource.TestCheckResourceAttr(fileFormatsLikeFirst.DatasourceReference(), "file_formats.0.show_output.0.name", idOne.Name()),
					resource.TestCheckResourceAttr(fileFormatsLikeFirst.DatasourceReference(), "file_formats.0.describe_output.#", "0"),
				),
			},
			{
				Config: accconfig.FromModels(t, fileFormatModel1, fileFormatModel2, fileFormatModel3, fileFormatsLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fileFormatsLikePrefix.DatasourceReference(), "file_formats.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, fileFormatModel1, fileFormatModel2, fileFormatModel3, fileFormatsInDatabase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fileFormatsInDatabase.DatasourceReference(), "file_formats.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, fileFormatModel1, fileFormatModel2, fileFormatModel3, fileFormatsNonExisting),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fileFormatsNonExisting.DatasourceReference(), "file_formats.#", "0"),
				),
			},
		},
	})
}

func TestAcc_FileFormats_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	schemaId := testClient().Ids.SchemaId()
	comment := random.Comment()

	csvModel := model.FileFormatCsv("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithCompression(string(sdk.CsvCompressionGzip)).
		WithFieldDelimiter(";").
		WithComment(comment)

	fileFormatsWithoutDescribe := datasourcemodel.FileFormats("test").
		WithLike(id.Name()).
		WithInSchema(schemaId).
		WithWithDescribe(false).
		WithDependsOn(csvModel.ResourceReference())

	fileFormatsWithDescribe := datasourcemodel.FileFormats("test").
		WithLike(id.Name()).
		WithInSchema(schemaId).
		WithDependsOn(csvModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.FileFormatCsv),
		Steps: []resource.TestStep{
			// without describe
			{
				Config: accconfig.FromModels(t, csvModel, fileFormatsWithoutDescribe),
				Check: assertThat(
					t,
					resourceshowoutputassert.FileFormatsDatasourceShowOutput(t, fileFormatsWithoutDescribe.DatasourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasType(sdk.FileFormatTypeCsv).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasFormatOptionsNotEmpty(),
					assert.Check(resource.TestCheckResourceAttr(fileFormatsWithoutDescribe.DatasourceReference(), "file_formats.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(fileFormatsWithoutDescribe.DatasourceReference(), "file_formats.0.describe_output.#", "0")),
				),
			},
			// with describe
			{
				Config: accconfig.FromModels(t, csvModel, fileFormatsWithDescribe),
				Check: assertThat(
					t,
					resourceshowoutputassert.FileFormatsDatasourceShowOutput(t, fileFormatsWithDescribe.DatasourceReference()).
						HasName(id.Name()).
						HasType(sdk.FileFormatTypeCsv).
						HasComment(comment),
					assert.Check(resource.TestCheckResourceAttr(fileFormatsWithDescribe.DatasourceReference(), "file_formats.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(fileFormatsWithDescribe.DatasourceReference(), "file_formats.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(fileFormatsWithDescribe.DatasourceReference(), "file_formats.0.describe_output.0.id", id.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(fileFormatsWithDescribe.DatasourceReference(), "file_formats.0.describe_output.0.type", string(sdk.FileFormatTypeCsv))),
					assert.Check(resource.TestCheckResourceAttr(fileFormatsWithDescribe.DatasourceReference(), "file_formats.0.describe_output.0.csv.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(fileFormatsWithDescribe.DatasourceReference(), "file_formats.0.describe_output.0.csv.0.compression", string(sdk.CsvCompressionGzip))),
					assert.Check(resource.TestCheckResourceAttr(fileFormatsWithDescribe.DatasourceReference(), "file_formats.0.describe_output.0.csv.0.field_delimiter", ";")),
					assert.Check(resource.TestCheckResourceAttr(fileFormatsWithDescribe.DatasourceReference(), "file_formats.0.describe_output.0.json.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(fileFormatsWithDescribe.DatasourceReference(), "file_formats.0.describe_output.0.avro.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(fileFormatsWithDescribe.DatasourceReference(), "file_formats.0.describe_output.0.orc.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(fileFormatsWithDescribe.DatasourceReference(), "file_formats.0.describe_output.0.parquet.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(fileFormatsWithDescribe.DatasourceReference(), "file_formats.0.describe_output.0.xml.#", "0")),
				),
			},
		},
	})
}

// TestAcc_FileFormats_AllTypes checks that the type-specific describe output is filled for every supported file format type.
func TestAcc_FileFormats_AllTypes(t *testing.T) {
	prefix := random.AlphaN(4)
	schemaId := testClient().Ids.SchemaId()
	csvId := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	jsonId := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	avroId := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	orcId := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	parquetId := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	xmlId := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)

	csvModel := model.FileFormatCsv("csv", csvId.DatabaseName(), csvId.SchemaName(), csvId.Name())
	jsonModel := model.FileFormatJson("json", jsonId.DatabaseName(), jsonId.SchemaName(), jsonId.Name())
	avroModel := model.FileFormatAvro("avro", avroId.DatabaseName(), avroId.SchemaName(), avroId.Name())
	orcModel := model.FileFormatOrc("orc", orcId.DatabaseName(), orcId.SchemaName(), orcId.Name())
	parquetModel := model.FileFormatParquet("parquet", parquetId.DatabaseName(), parquetId.SchemaName(), parquetId.Name())
	xmlModel := model.FileFormatXml("xml", xmlId.DatabaseName(), xmlId.SchemaName(), xmlId.Name())

	fileFormats := datasourcemodel.FileFormats("test").
		WithLike(prefix+"%").
		WithInSchema(schemaId).
		WithDependsOn(
			csvModel.ResourceReference(),
			jsonModel.ResourceReference(),
			avroModel.ResourceReference(),
			orcModel.ResourceReference(),
			parquetModel.ResourceReference(),
			xmlModel.ResourceReference(),
		)

	// SHOW FILE FORMATS returns the file formats ordered by name, so the indexes are resolved from the sorted names.
	sortedNames := collections.Map([]sdk.SchemaObjectIdentifier{csvId, jsonId, avroId, orcId, parquetId, xmlId}, func(id sdk.SchemaObjectIdentifier) string { return id.Name() })
	slices.Sort(sortedNames)
	indexOf := func(id sdk.SchemaObjectIdentifier) int {
		return slices.Index(sortedNames, id.Name())
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, csvModel, jsonModel, avroModel, orcModel, parquetModel, xmlModel, fileFormats),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fileFormats.DatasourceReference(), "file_formats.#", "6"),
					resource.TestCheckResourceAttr(fileFormats.DatasourceReference(), fmt.Sprintf("file_formats.%d.describe_output.0.csv.#", indexOf(csvId)), "1"),
					resource.TestCheckResourceAttr(fileFormats.DatasourceReference(), fmt.Sprintf("file_formats.%d.describe_output.0.json.#", indexOf(jsonId)), "1"),
					resource.TestCheckResourceAttr(fileFormats.DatasourceReference(), fmt.Sprintf("file_formats.%d.describe_output.0.avro.#", indexOf(avroId)), "1"),
					resource.TestCheckResourceAttr(fileFormats.DatasourceReference(), fmt.Sprintf("file_formats.%d.describe_output.0.orc.#", indexOf(orcId)), "1"),
					resource.TestCheckResourceAttr(fileFormats.DatasourceReference(), fmt.Sprintf("file_formats.%d.describe_output.0.parquet.#", indexOf(parquetId)), "1"),
					resource.TestCheckResourceAttr(fileFormats.DatasourceReference(), fmt.Sprintf("file_formats.%d.describe_output.0.xml.#", indexOf(xmlId)), "1"),
				),
			},
		},
	})
}
