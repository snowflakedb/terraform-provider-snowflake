package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO [SNOW-2054208]: decide its fate during packages cleanup
// Exports in this file are used by test resources in the functional tests.
// It will be handled during the packages cleanup.
// Temporarily it's done like this to not update ~40 production files using these functions.

func SdkValidation[T any](normalize func(string) (T, error)) schema.SchemaValidateDiagFunc {
	return sdkValidation(normalize)
}

func SetStateToValuesFromConfig(d *schema.ResourceData, resourceSchema map[string]*schema.Schema, fields []string) error {
	return setStateToValuesFromConfig(d, resourceSchema, fields)
}
