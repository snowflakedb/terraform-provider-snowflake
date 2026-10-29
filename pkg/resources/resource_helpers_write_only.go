package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Write-only attributes (schema.Schema.WriteOnly) are never stored in the plan or the state; Terraform sends them
// only in the config, and only during apply. That has two consequences the helpers below deal with:
//   - the value has to be read from the raw config (d.GetOk always returns the zero value for such an attribute);
//   - d.HasChange on the attribute is always false, so updates have to be triggered by a companion version attribute
//     that is stored in the state.

// writeOnlyStringValue extracts a write-only string from a raw config value. An unknown value can be observed during
// plan (e.g. when the value comes from an ephemeral resource);
func writeOnlyStringValue(configValue cty.Value) *string {
	if configValue.IsNull() || !configValue.IsKnown() {
		return nil
	}
	value := configValue.AsString()
	if value == "" {
		return nil
	}
	return sdk.String(value)
}

// getWriteOnlyStringFromConfig reads the write-only attribute under key from the raw config.
func getWriteOnlyStringFromConfig(d *schema.ResourceData, key string) (*string, error) {
	configValue, diags := d.GetRawConfigAt(cty.GetAttrPath(key))
	if diags.HasError() {
		return nil, fmt.Errorf("reading raw config for write-only attribute %s: %v", key, diags)
	}
	return writeOnlyStringValue(configValue), nil
}

// isWriteOnlyAttributeInUse reports whether the write-only attribute under woKey is managed by the configuration.
func isWriteOnlyAttributeInUse(d *schema.ResourceData, woKey string, versionKey string) bool {
	if configValue, diags := d.GetRawConfigAt(cty.GetAttrPath(woKey)); !diags.HasError() && !configValue.IsNull() {
		return true
	}
	if version, ok := d.GetOk(versionKey); ok && version.(int) != 0 {
		return true
	}
	return false
}

// setStringUnlessWriteOnly behaves like d.Set, but keeps the value out of the state when
// it is managed through the write-only counterpart of the field.
func setStringUnlessWriteOnly(d *schema.ResourceData, key string, woKey string, versionKey string, value string) error {
	if isWriteOnlyAttributeInUse(d, woKey, versionKey) {
		return d.Set(key, "")
	}
	return d.Set(key, value)
}

// writeOnlyStringAttributeCreate sets createField to the write-only value found in the raw config.
func writeOnlyStringAttributeCreate(d *schema.ResourceData, key string, createField **string) error {
	value, err := getWriteOnlyStringFromConfig(d, key)
	if err != nil {
		return err
	}
	if value != nil {
		*createField = value
	}
	return nil
}

// stringAttributeWithWriteOnlyVariantUpdate handles a pair of attributes (a regular one and its write-only variant)
// that map to the same Snowflake property.
func stringAttributeWithWriteOnlyVariantUpdate(d *schema.ResourceData, key string, woKey string, versionKey string, setField **string, unsetField **bool) error {
	writeOnlyValue, err := getWriteOnlyStringFromConfig(d, woKey)
	if err != nil {
		return err
	}
	switch {
	case writeOnlyValue != nil:
		if d.HasChange(versionKey) {
			*setField = writeOnlyValue
		}
		return nil
	case d.Get(key).(string) != "":
		return stringAttributeUpdate(d, key, setField, unsetField)
	case d.HasChange(key) || d.HasChange(versionKey):
		*unsetField = sdk.Bool(true)
		return nil
	default:
		return nil
	}
}
