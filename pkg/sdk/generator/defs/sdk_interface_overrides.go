//go:build sdk_generation

package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

// sdkInterfaceOverrides adds to an ssot-built SDK interface. The key is the ssot object Name.
var sdkInterfaceOverrides = map[string]gen.SdkInterfaceOverride{}
