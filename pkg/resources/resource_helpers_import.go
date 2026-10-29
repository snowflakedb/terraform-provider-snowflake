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

func importStringAttribute[T ~string](d *schema.ResourceData, key string, importPermissiveMode bool, snowflakeValue T) error {
	if importPermissiveMode {
		return d.Set(key, "")
	}
	return d.Set(key, snowflakeValue)
}

func importOptionalBooleanStringAttribute(d *schema.ResourceData, key string, importPermissiveMode bool, snowflakeValue *bool) error {
	if snowflakeValue == nil {
		return d.Set(key, BooleanDefault)
	}
	return importBooleanStringAttribute(d, key, importPermissiveMode, *snowflakeValue)
}

func importOptionalIntWithIntDefaultAttribute(d *schema.ResourceData, key string, importPermissiveMode bool, snowflakeValue *int) error {
	if snowflakeValue == nil {
		return d.Set(key, IntDefault)
	}
	return importIntWithIntDefaultAttribute(d, key, importPermissiveMode, *snowflakeValue)
}

func importOptionalIntAttribute(d *schema.ResourceData, key string, importPermissiveMode bool, snowflakeValue *int) error {
	if snowflakeValue == nil {
		return nil
	}
	return importIntAttribute(d, key, importPermissiveMode, *snowflakeValue)
}

func importOptionalStringAttribute[T ~string](d *schema.ResourceData, key string, importPermissiveMode bool, snowflakeValue *T) error {
	if snowflakeValue == nil {
		return nil
	}
	return importStringAttribute(d, key, importPermissiveMode, string(*snowflakeValue))
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
