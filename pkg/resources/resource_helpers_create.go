package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func stringAttributeCreate(d *schema.ResourceData, key string, createField **string) error {
	if v, ok := d.GetOk(key); ok {
		*createField = sdk.String(v.(string))
	}
	return nil
}

func stringAttributeCreateBuilder[T any](d *schema.ResourceData, key string, setValue func(string) T) error {
	if v, ok := d.GetOk(key); ok {
		setValue(v.(string))
	}
	return nil
}

func intAttributeCreate(d *schema.ResourceData, key string, createField **int) error {
	if v, ok := d.GetOk(key); ok {
		*createField = sdk.Int(v.(int))
	}
	return nil
}

func intAttributeWithSpecialDefaultCreate(d *schema.ResourceData, key string, createField **int) error {
	if v := d.Get(key).(int); v != IntDefault {
		*createField = sdk.Int(v)
	}
	return nil
}

func intAttributeWithSpecialDefaultCreateBuilder[T any](d *schema.ResourceData, key string, setValue func(int) T) error {
	if v := d.Get(key).(int); v != IntDefault {
		setValue(v)
	}
	return nil
}

func booleanStringAttributeCreate(d *schema.ResourceData, key string, createField **bool) error {
	if v := d.Get(key).(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return err
		}
		*createField = sdk.Bool(parsed)
	}
	return nil
}

func booleanStringAttributeCreateBuilder[T any](d *schema.ResourceData, key string, setValue func(bool) T) error {
	if v := d.Get(key).(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return err
		}
		setValue(parsed)
	}
	return nil
}

func schemaObjectIdentifierAttributeCreate(d *schema.ResourceData, key string, createField **sdk.SchemaObjectIdentifier) error {
	if v, ok := d.GetOk(key); ok {
		*createField = sdk.Pointer(sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(v.(string)))
	}
	return nil
}

func accountObjectIdentifierAttributeCreate(d *schema.ResourceData, key string, createField **sdk.AccountObjectIdentifier) error {
	if v, ok := d.GetOk(key); ok {
		*createField = sdk.Pointer(sdk.NewAccountObjectIdentifier(v.(string)))
	}
	return nil
}

func objectIdentifierAttributeCreate(d *schema.ResourceData, key string, createField **sdk.ObjectIdentifier) error {
	if v, ok := d.GetOk(key); ok {
		objectIdentifier, err := sdk.ParseObjectIdentifierString(v.(string))
		if err != nil {
			return err
		}
		*createField = sdk.Pointer(objectIdentifier)
	}
	return nil
}

func attributeDirectValueCreate[T any](d *schema.ResourceData, key string, createField **T, value *T) error {
	if _, ok := d.GetOk(key); ok {
		*createField = value
	}
	return nil
}

func attributeMappedValueCreate[T any](d *schema.ResourceData, key string, createField **T, mapper func(value any) (*T, error)) error {
	if v, ok := d.GetOk(key); ok {
		value, err := mapper(v)
		if err != nil {
			return err
		}
		*createField = value
	}
	return nil
}

func attributeMappedValueCreateBuilder[InputType any, MappedType any, RequestBuilder any](d *schema.ResourceData, key string, setValue func(MappedType) RequestBuilder, mapper func(value InputType) (MappedType, error)) error {
	if v, ok := d.GetOk(key); ok {
		value, err := mapper(v.(InputType))
		if err != nil {
			return err
		}
		setValue(value)
	}
	return nil
}

func copyGrantsAttributeCreate(d *schema.ResourceData, isOrReplace bool, orReplaceField, copyGrantsField **bool) error {
	if isOrReplace {
		*orReplaceField = sdk.Bool(true)
		if d.GetRawConfig().AsValueMap()["copy_grants"].True() {
			*copyGrantsField = sdk.Bool(true)
		}
	}
	return nil
}
