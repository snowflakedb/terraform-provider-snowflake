//go:build sdk_generation

package defs

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
)

func init() {
	fmt.Printf("katarakta %v\n", gen.AllSdkObjectDefinitions)
	gen.AllSdkObjectDefinitions = append(gen.AllSdkObjectDefinitions, SequencesDef)
	fmt.Printf("katarakta %v\n", gen.AllSdkObjectDefinitions)
}
