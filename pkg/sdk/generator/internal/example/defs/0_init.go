//go:build sdk_generation_examples

package defs

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

func init() {
	fmt.Printf("katarakta %v\n", gen.AllSdkObjectDefinitions)
	gen.AllSdkObjectDefinitions = append(gen.AllSdkObjectDefinitions, DatabaseRole, ToOptsOptionalExample)
	fmt.Printf("katarakta %v\n", gen.AllSdkObjectDefinitions)
}
