package resources

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

func nestedStringAttributeCreate(config map[string]any, key string, createField **string) error {
	if v, ok := config[key]; ok && v.(string) != "" {
		*createField = sdk.String(v.(string))
	}
	return nil
}

func nestedBooleanStringAttributeCreate(config map[string]any, key string, createField **bool) error {
	if v, ok := config[key]; ok && v.(string) != BooleanDefault {
		parsed, err := booleanStringToBool(v.(string))
		if err != nil {
			return err
		}
		*createField = sdk.Bool(parsed)
	}
	return nil
}
