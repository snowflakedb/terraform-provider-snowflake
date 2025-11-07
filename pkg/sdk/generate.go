package sdk

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/tmp"
)

//go:generate go run poc/generator/main/main.go $SF_TF_GENERATOR_ARGS

func init() {
	fmt.Printf("katarakta %v\n", tmp.AllSdkObjectDefinitions)
	tmp.AllSdkObjectDefinitions = append(tmp.AllSdkObjectDefinitions, defs.SequencesDef)
	fmt.Printf("katarakta %v\n", tmp.AllSdkObjectDefinitions)
}
