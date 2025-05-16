package config_test

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_JsonConfigProvider(t *testing.T) {
	t.Run("test resource json config", func(t *testing.T) {
		model := Some("some_name", "abc").WithDependsOn("abc.def")
		expectedResult := `{
    "resource": {
        "snowflake_share": {
            "some_name": {
                "name": "abc",
                "depends_on": [
                    "abc.def"
                ]
            }
        }
    }
}`

		result, err := config.DefaultJsonConfigProvider.ResourceJsonFromModel(model)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, string(result))
	})

	t.Run("test special variables", func(t *testing.T) {
		model := Some("some_name", "abc").WithTextFieldExplicitNull().WithListFieldEmpty()
		expectedResult := fmt.Sprintf(`{
    "resource": {
        "snowflake_share": {
            "some_name": {
                "name": "abc",
                "text_field": "%[1]s",
                "list_field": []
            }
        }
    }
}`, config.SnowflakeProviderConfigNull)

		result, err := config.DefaultJsonConfigProvider.ResourceJsonFromModel(model)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, string(result))
	})

	t.Run("test multiline variable", func(t *testing.T) {
		model := Some("some_name", "abc").WithMultilineField("some\nmultiline\ncontent")
		expectedResult := fmt.Sprintf(`{
    "resource": {
        "snowflake_share": {
            "some_name": {
                "name": "abc",
                "multiline_field": "%[1]s%[2]s%[1]s"
            }
        }
    }
}`, config.SnowflakeProviderConfigMultilineMarker, "some\\nmultiline\\ncontent")

		result, err := config.DefaultJsonConfigProvider.ResourceJsonFromModel(model)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, string(result))
	})

	t.Run("test resource json config when proper marshaller is absent", func(t *testing.T) {
		model := SomeOther("some_name", "abc").WithDependsOn("abc.def")
		expectedResult := `{
    "resource": {
        "snowflake_share": {
            "some_name": {
                "name": "abc"
            }
        }
    }
}`

		result, err := config.DefaultJsonConfigProvider.ResourceJsonFromModel(model)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, string(result))
	})

	t.Run("test datasource json config", func(t *testing.T) {
		model := datasourcemodel.Databases("some_name")
		expectedResult := `{
    "data": {
        "snowflake_databases": {
            "some_name": {
                "single_attribute_workaround": "SF_TF_TEST_SINGLE_ATTRIBUTE_WORKAROUND"
            }
        }
    }
}`

		result, err := config.DefaultJsonConfigProvider.DatasourceJsonFromModel(model)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, string(result))
	})

	t.Run("test provider json config", func(t *testing.T) {
		model := providermodel.SnowflakeProvider()
		expectedResult := `{
    "provider": {
        "snowflake": {}
    }
}`

		result, err := config.DefaultJsonConfigProvider.ProviderJsonFromModel(model)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, string(result))
	})

	t.Run("test dynamic block config", func(t *testing.T) {
		dynamicBlock := config.NewDynamicBlock("argument", "arguments", []string{"name", "type"})
		model := DynamicBlockExample("test", "abc").
			WithDynamicBlock(dynamicBlock)
		expectedResult := fmt.Sprintf(`{
    "resource": {
        "snowflake_share": {
            "test": {
                "name": "abc",
                "dynamic": {
                    "argument": {
                        "for_each": "%[1]svar.arguments%[1]s",
                        "content": {
                            "name": "%[1]sargument.value[%[2]sname%[2]s]%[1]s",
                            "type": "%[1]sargument.value[%[2]stype%[2]s]%[1]s"
                        }
                    }
                }
            }
        }
    }
}`, config.SnowflakeProviderConfigUnquoteMarker, config.SnowflakeProviderConfigQuoteMarker)

		result, err := config.DefaultJsonConfigProvider.ResourceJsonFromModel(model)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, string(result))
	})
}
