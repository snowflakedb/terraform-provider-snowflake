package main

import (
	"regexp"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
)

var resourceIdAllowedCharacters = regexp.MustCompile("[a-zA-Z0-9\\-_]+")

func NormalizeResourceId(resourceId string) string {
	mappedResourceId := strings.Map(func(r rune) rune {
		switch r {
		case '.':
			return '_'
		default:
			if resourceIdAllowedCharacters.MatchString(string(r)) {
				return r
			}
			return -1
		}
	}, resourceId)

	return "snowflake_generated_" + mappedResourceId
}

// ResourceFromModel is a copy of config.ResourceFromModel function, but it doesn't use testing.T internally.
func ResourceFromModel(model config.ResourceModel) (string, error) {
	resourceJson, err := config.DefaultJsonConfigProvider.ResourceJsonFromModel(model)
	if err != nil {
		return "", err
	}

	return config.DefaultHclConfigProvider.HclFromJson(resourceJson)
}
