//go:build sdk_generation

package defs

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/tmp"
)

func init() {
	fmt.Printf("katarakta %v\n", tmp.AllSdkObjectDefinitions)
	tmp.AllSdkObjectDefinitions = append(tmp.AllSdkObjectDefinitions, SequencesDef)
	fmt.Printf("katarakta %v\n", tmp.AllSdkObjectDefinitions)
}
