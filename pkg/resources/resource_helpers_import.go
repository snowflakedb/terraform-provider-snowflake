package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func importBooleanStringAttribute(d *schema.ResourceData, key string, importPermissiveMode bool, snowflakeValue bool) error {
	if importPermissiveMode {
		return d.Set(key, BooleanDefault)
	}
	return d.Set(key, booleanStringFromBool(snowflakeValue))
}

func importIntWithIntDefaultAttribute(d *schema.ResourceData, key string, importPermissiveMode bool, snowflakeValue int) error {
	if importPermissiveMode {
		return d.Set(key, IntDefault)
	}
	return d.Set(key, snowflakeValue)
}

func importIntAttribute(d *schema.ResourceData, key string, importPermissiveMode bool, snowflakeValue int) error {
	if importPermissiveMode {
		return d.Set(key, 0)
	}
	return d.Set(key, snowflakeValue)
}

func importStringAttribute(d *schema.ResourceData, key string, importPermissiveMode bool, snowflakeValue string) error {
	if importPermissiveMode {
		return d.Set(key, "")
	}
	return d.Set(key, snowflakeValue)
}
