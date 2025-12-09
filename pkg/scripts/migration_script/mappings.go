package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

var resourceIdDisallowedCharacters = regexp.MustCompile("[^a-zA-Z0-9\\-_]")

func NormalizeResourceId(resourceId string) string {
	resourceId = strings.ReplaceAll(resourceId, `.`, "_")
	return "snowflake_generated_" + string(resourceIdDisallowedCharacters.ReplaceAll([]byte(resourceId), []byte("")))
}

func ResourceId(resource resources.Resource, id string) string {
	return NormalizeResourceId(fmt.Sprintf("%s_%s", strings.TrimPrefix(resource.String(), "snowflake_"), id))
}

// ResourceFromModel is a copy of config.ResourceFromModel function, but it doesn't use testing.T internally.
func ResourceFromModel(model config.ResourceModel) (string, error) {
	resourceJson, err := config.DefaultJsonConfigProvider.ResourceJsonFromModel(model)
	if err != nil {
		return "", err
	}

	return config.DefaultHclConfigProvider.HclFromJson(resourceJson)
}

func HandleResources[T ConvertibleCsvRow[R], R any](
	config *Config,
	csvInput [][]string,
	mapObjToModel func(obj R) (accconfig.ResourceModel, *ImportModel, error),
) (string, error) {
	objects, err := ConvertCsvInput[T, R](csvInput)
	if err != nil {
		return "", err
	}

	resourceModels := make([]accconfig.ResourceModel, 0)
	importModels := make([]ImportModel, 0)

	for _, object := range objects {
		mappedModel, importModel, err := mapObjToModel(object)
		if err != nil {
			log.Printf("Error converting object of type %T to model: %v. Skipping object and continuing with other mappings.", object, err)
		} else {
			resourceModels = append(resourceModels, mappedModel)
			importModels = append(importModels, *importModel)
		}
	}

	mappedModels, err := collections.MapErr(resourceModels, ResourceFromModel)
	if err != nil {
		return "", fmt.Errorf("errors from resource model to HCL conversion: %w", err)
	}

	mappedImports, err := collections.MapErr(importModels, func(importModel ImportModel) (string, error) {
		return TransformImportModel(config, importModel)
	})
	if err != nil {
		return "", fmt.Errorf("errors during import transformations: %w", err)
	}

	outputBuilder := new(strings.Builder)
	outputBuilder.WriteString(collections.JoinStrings(mappedModels, "\n"))
	outputBuilder.WriteString(collections.JoinStrings(mappedImports, ""))

	return outputBuilder.String(), nil
}
