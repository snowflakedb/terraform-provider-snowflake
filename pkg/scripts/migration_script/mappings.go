package main

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

// ResourceFromModel is a copy of config.ResourceFromModel function, but it doesn't use testing.T internally.
func ResourceFromModel(model config.ResourceModel) (string, error) {
	resourceJson, err := config.DefaultJsonConfigProvider.ResourceJsonFromModel(model)
	if err != nil {
		return "", err
	}

	return config.DefaultHclConfigProvider.HclFromJson(resourceJson)
}
