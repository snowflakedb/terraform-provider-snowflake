//go:build sdk_generation

package defs

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

func init() {
	fmt.Printf("katarakta %v\n", gen.AllSdkObjectDefinitions)
	gen.AllSdkObjectDefinitions = append(gen.AllSdkObjectDefinitions, SequencesDef)
	fmt.Printf("katarakta %v\n", gen.AllSdkObjectDefinitions)
}
