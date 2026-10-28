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

func importSetAttribute[T any](d *schema.ResourceData, key string, importPermissiveMode bool, snowflakeValue []T) error {
	if importPermissiveMode {
		return d.Set(key, []T{})
	}
	return d.Set(key, snowflakeValue)
}

// importBlockAttribute writes nothing in permissive mode. Snowflake may report optional blocks even when they aren't set,
// so writing them would make a config that omits a block show `- block`.
func importBlockAttribute(d *schema.ResourceData, key string, importPermissiveMode bool, snowflakeValue []map[string]any) error {
	if importPermissiveMode {
		return nil
	}
	return d.Set(key, snowflakeValue)
}
