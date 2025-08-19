package config

import (
	"reflect"
	"strings"
	"testing"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/stretchr/testify/require"
)

// ResourceFromModel should be used in terraform acceptance tests for Config attribute to get string config from ResourceModel.
// Current implementation is an improved implementation using two steps:
// - .tf.json generation
// - conversion to HCL using hcl v1 lib
// It is still not ideal. HCL v2 should be considered.
func ResourceFromModel(t *testing.T, model ResourceModel) string {
	t.Helper()

	resourceJson, err := DefaultJsonConfigProvider.ResourceJsonFromModel(model)
	require.NoError(t, err)

	hcl, err := DefaultHclConfigProvider.HclFromJson(resourceJson)
	require.NoError(t, err)

	// TODO: Uncomment or change to log.Println as it breaks migration script due to incorrectly set up testing.T
	//t.Logf("Generated config:\n%s", hcl)

	return hcl
}

// DatasourceFromModel should be used in terraform acceptance tests for Config attribute to get string config from DatasourceModel.
// Current implementation is an improved implementation using two steps:
// - .tf.json generation
// - conversion to HCL using hcl v1 lib
// It is still not ideal. HCL v2 should be considered.
func DatasourceFromModel(t *testing.T, model DatasourceModel) string {
	t.Helper()

	datasourceJson, err := DefaultJsonConfigProvider.DatasourceJsonFromModel(model)
	require.NoError(t, err)

	hcl, err := DefaultHclConfigProvider.HclFromJson(datasourceJson)
	require.NoError(t, err)
	t.Logf("Generated config:\n%s", hcl)

	return hcl
}

// ProviderFromModel should be used in terraform acceptance tests for Config attribute to get string config from ProviderModel.
// Current implementation is an improved implementation using two steps:
// - .tf.json generation
// - conversion to HCL using hcl v1 lib
// It is still not ideal. HCL v2 should be considered.
func ProviderFromModel(t *testing.T, model ProviderModel) string {
	t.Helper()

	providerJson, err := DefaultJsonConfigProvider.ProviderJsonFromModel(model)
	require.NoError(t, err)

	hcl, err := DefaultHclConfigProvider.HclFromJson(providerJson)
	require.NoError(t, err)
	hcl, err = revertEqualSignForMapTypeAttributes(hcl)
	require.NoError(t, err)

	return hcl
}

// VariableFromModel should be used in terraform acceptance tests for Config attribute to get string config from VariableModel.
// Current implementation is an improved implementation using two steps:
// - .tf.json generation
// - conversion to HCL using hcl v1 lib
// It is still not ideal. HCL v2 should be considered.
func VariableFromModel(t *testing.T, model TerraformBlockModel) string {
	t.Helper()

	variableJson, err := DefaultJsonConfigProvider.TerraformBlockJsonFromModel(model)
	require.NoError(t, err)

	hcl, err := DefaultHclConfigProvider.HclFromJson(variableJson)
	require.NoError(t, err)
	t.Logf("Generated config:\n%s", hcl)

	return hcl
}

// FromModels should be used in terraform acceptance tests for Config attribute to get string config from all models.
// FromModels allows to combine multiple model types.
// TODO [SNOW-1501905]: introduce some common interface for all three existing models (ResourceModel, DatasourceModel, and ProviderModel)
func FromModels(t *testing.T, models ...any) string {
	t.Helper()

	var sb strings.Builder
	for i, model := range models {
		switch m := model.(type) {
		case ResourceModel:
			sb.WriteString(ResourceFromModel(t, m))
		case DatasourceModel:
			sb.WriteString(DatasourceFromModel(t, m))
		case ProviderModel:
			sb.WriteString(ProviderFromModel(t, m))
		case TerraformBlockModel:
			sb.WriteString(VariableFromModel(t, m))
		default:
			t.Fatalf("unknown model: %T", model)
		}
		if i < len(models)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// ConfigVariablesFromModel constructs config.Variables needed in acceptance tests that are using ConfigVariables in
// combination with ConfigDirectory. It's necessary for cases not supported by FromModel, like lists of objects.
// Use ResourceFromModel, DatasourceFromModel, ProviderFromModel, and FromModels instead.
func ConfigVariablesFromModel(t *testing.T, model ResourceModel) tfconfig.Variables {
	t.Helper()
	variables := make(tfconfig.Variables)
	rType := reflect.TypeOf(model).Elem()
	rValue := reflect.ValueOf(model).Elem()
	for i := 0; i < rType.NumField(); i++ {
		field := rType.Field(i)
		if jsonTag, ok := field.Tag.Lookup("json"); ok {
			name := strings.Split(jsonTag, ",")[0]
			if fieldValue, ok := rValue.Field(i).Interface().(tfconfig.Variable); ok {
				variables[name] = fieldValue
			}
		}
	}
	return variables
}

// ConfigVariablesFromModels can be used to create a list of objects that are referring to the same resource model.
// It's useful when there's a need to create associations between objects of the same type in Snowflake.
func ConfigVariablesFromModels(t *testing.T, variableName string, models ...ResourceModel) tfconfig.Variables {
	t.Helper()
	allVariables := make([]tfconfig.Variable, 0)
	for _, model := range models {
		allVariables = append(allVariables, tfconfig.ObjectVariable(ConfigVariablesFromModel(t, model)))
	}
	return tfconfig.Variables{
		variableName: tfconfig.ListVariable(allVariables...),
	}
}
