package schemas

import (
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// additionalSchemaMapper is implemented by generated mapper types when AdditionalMapping is set.
// The *_ext.go files provide the implementation.
type additionalSchemaMapper[T any] interface {
	additionalSchema() map[string]*schema.Schema
	additionalToSchema(src *T, dst map[string]any)
}

func mergeSchema(generated, additional map[string]*schema.Schema) map[string]*schema.Schema {
	merged := make(map[string]*schema.Schema, len(generated)+len(additional))
	maps.Copy(merged, generated)
	for key, value := range additional {
		if _, exists := merged[key]; exists {
			panic(fmt.Sprintf("additionalSchema key %q collides with generated schema", key))
		}
		merged[key] = value
	}
	return merged
}
