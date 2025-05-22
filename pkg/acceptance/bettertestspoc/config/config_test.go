package config_test

import (
	"strings"
	"testing"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/stretchr/testify/require"
)

func Test_ResourceFromModelPoc(t *testing.T) {
	t.Run("test basic", func(t *testing.T) {
		someModel := Some("test", "Some Name")
		expectedOutput := strings.TrimPrefix(`
resource "snowflake_share" "test" {
  name = "Some Name"
}
`, "\n")
		result := config.ResourceFromModel(t, someModel)

		require.Equal(t, expectedOutput, result)
	})

	// TODO [SNOW-1501905]: replace \t characters with actual tabs
	t.Run("test tabs in multiline", func(t *testing.T) {
		someModel := Some("test", "Some Name").
			WithMultilineField("some\n\tmulti\tline\n\t\t\tcontent")
		expectedOutput := strings.TrimPrefix(`
resource "snowflake_share" "test" {
  name = "Some Name"
  multiline_field = <<EOT
some
\tmulti\tline
\t\t\tcontent
EOT
}
`, "\n")
		result := config.ResourceFromModel(t, someModel)

		require.Equal(t, expectedOutput, result)
	})

	t.Run("test full", func(t *testing.T) {
		someModel := Some("test", "Some Name").
			WithComment("Some Comment").
			WithStringList("a", "b", "a").
			WithStringSet("a", "b", "c").
			WithObjectList(
				Item{IntField: 1, StringField: "first item"},
				Item{IntField: 2, StringField: "second item"},
			).
			WithSingleObject("one", 2).
			WithTextFieldExplicitNull().
			WithListFieldEmpty().
			WithMultilineField("some\nmultiline\ncontent").
			WithDependsOn("some_other_resource.some_name", "other_resource.some_other_name", "third_resource.third_name")
		expectedOutput := strings.TrimPrefix(`
resource "snowflake_share" "test" {
  comment = "Some Comment"
  name = "Some Name"
  string_list = ["a", "b", "a"]
  string_set = ["a", "b", "c"]
  object_list {
    int_field = 1
    string_field = "first item"
  }
  object_list {
    int_field = 2
    string_field = "second item"
  }
  single_object {
    a = "one"
    b = 2
  }
  text_field = null
  list_field = []
  multiline_field = <<EOT
some
multiline
content
EOT
  depends_on = [some_other_resource.some_name, other_resource.some_other_name, third_resource.third_name]
}
`, "\n")

		result := config.ResourceFromModel(t, someModel)

		require.Equal(t, expectedOutput, result)
	})

	t.Run("test dynamic block", func(t *testing.T) {
		model := DynamicBlockExample("test", "abc").
			WithDynamicBlock(config.NewDynamicBlock("argument", "arguments", []string{"name", "type"}))
		expectedOutput := strings.TrimPrefix(`
resource "snowflake_share" "test" {
  name = "abc"
  dynamic "argument" {
    for_each = var.arguments
    content {
      name = argument.value["name"]
      type = argument.value["type"]
    }
  }
}
`, "\n")
		result := config.ResourceFromModel(t, model)

		require.Equal(t, expectedOutput, result)
	})

	t.Run("test no dynamic block", func(t *testing.T) {
		model := DynamicBlockExample("test", "abc")
		expectedOutput := strings.TrimPrefix(`
resource "snowflake_share" "test" {
  name = "abc"
}
`, "\n")
		result := config.ResourceFromModel(t, model)

		require.Equal(t, expectedOutput, result)
	})
}

func Test_DatasourceFromModelPoc(t *testing.T) {
	t.Run("test basic", func(t *testing.T) {
		datasourceModel := datasourcemodel.Databases("test")
		expectedOutput := strings.TrimPrefix(`
data "snowflake_databases" "test" {
}
`, "\n")
		result := config.DatasourceFromModel(t, datasourceModel)

		require.Equal(t, expectedOutput, result)
	})

	t.Run("test with some arguments", func(t *testing.T) {
		datasourceModel := datasourcemodel.Databases("test").WithLike("some").WithLimit(1)
		expectedOutput := strings.TrimPrefix(`
data "snowflake_databases" "test" {
  like = "some"
  limit {
    rows = 1
  }
}
`, "\n")
		result := config.DatasourceFromModel(t, datasourceModel)

		require.Equal(t, expectedOutput, result)
	})

	t.Run("test with depends on", func(t *testing.T) {
		datasourceModel := datasourcemodel.Databases("test").
			WithDependsOn("some_other_resource.some_name", "other_resource.some_other_name", "third_resource.third_name")
		expectedOutput := strings.TrimPrefix(`
data "snowflake_databases" "test" {
  depends_on = [some_other_resource.some_name, other_resource.some_other_name, third_resource.third_name]
}
`, "\n")
		result := config.DatasourceFromModel(t, datasourceModel)

		require.Equal(t, expectedOutput, result)
	})
}

func Test_ProviderFromModelPoc(t *testing.T) {
	t.Run("test basic", func(t *testing.T) {
		providerModel := providermodel.SnowflakeProvider()
		expectedOutput := strings.TrimPrefix(`
provider "snowflake" {}
`, "\n")
		result := config.ProviderFromModel(t, providerModel)

		require.Equal(t, expectedOutput, result)
	})

	t.Run("test with alias", func(t *testing.T) {
		providerModel := providermodel.SnowflakeProviderAlias("other_name")
		expectedOutput := strings.TrimPrefix(`
provider "snowflake" {
  alias = "other_name"
}
`, "\n")
		result := config.ProviderFromModel(t, providerModel)

		require.Equal(t, expectedOutput, result)
	})

	t.Run("test with some attributes", func(t *testing.T) {
		providerModel := providermodel.SnowflakeProvider().WithProfile("some_profile").WithUser("some user")
		expectedOutput := strings.TrimPrefix(`
provider "snowflake" {
  profile = "some_profile"
  user = "some user"
}
`, "\n")
		result := config.ProviderFromModel(t, providerModel)

		require.Equal(t, expectedOutput, result)
	})

	t.Run("test with parameters map", func(t *testing.T) {
		providerModel := providermodel.SnowflakeProvider().WithProfile("some_profile").WithParamsValue(
			tfconfig.MapVariable(map[string]tfconfig.Variable{
				"statement_timeout_in_seconds": tfconfig.IntegerVariable(31337),
			}),
		)
		expectedOutput := strings.TrimPrefix(`
provider "snowflake" {
  params = {
    statement_timeout_in_seconds = 31337
  }
  profile = "some_profile"
}
`, "\n")
		result := config.ProviderFromModel(t, providerModel)

		require.Equal(t, expectedOutput, result)
	})
}

func Test_ConfigFromModelsPoc(t *testing.T) {
	t.Run("test basic", func(t *testing.T) {
		providerModel := providermodel.SnowflakeProvider()
		someModel := Some("test", "Some Name")
		datasourceModel := datasourcemodel.Databases("test").WithDependsOn(someModel.ResourceReference())
		someOtherModel := Some("test2", "Some Name 2").WithDependsOn(datasourceModel.DatasourceReference())
		expectedOutput := strings.TrimPrefix(`
provider "snowflake" {}

resource "snowflake_share" "test" {
  name = "Some Name"
}

data "snowflake_databases" "test" {
  depends_on = [snowflake_share.test]
}

resource "snowflake_share" "test2" {
  name = "Some Name 2"
  depends_on = [data.snowflake_databases.test]
}
`, "\n")
		result := config.FromModels(t, providerModel, someModel, datasourceModel, someOtherModel)

		require.Equal(t, expectedOutput, result)
	})
}
