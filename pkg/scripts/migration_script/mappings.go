package main

import (
	"regexp"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
)

var resourceIdDisallowedCharacters = regexp.MustCompile("[^a-zA-Z0-9\\-_]")

func NormalizeResourceId(resourceId string) string {
	resourceId = strings.ReplaceAll(resourceId, `.`, "_")
	return "snowflake_generated_" + string(resourceIdDisallowedCharacters.ReplaceAll([]byte(resourceId), []byte("")))
}

// ResourceFromModel is a copy of config.ResourceFromModel function, but it doesn't use testing.T internally.
func ResourceFromModel(model config.ResourceModel) (string, error) {
	resourceJson, err := config.DefaultJsonConfigProvider.ResourceJsonFromModel(model)
	if err != nil {
		return "", err
	}

	return config.DefaultHclConfigProvider.HclFromJson(resourceJson)
}
