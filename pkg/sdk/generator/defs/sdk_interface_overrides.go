//go:build sdk_generation

package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

// sdkInterfaceOverrides adds to an ssot-built SDK interface. The key is the ssot object Name.
// An object that sets SdkDefineFromSsot does not use this map. Extend the SSOT spec instead.
var sdkInterfaceOverrides = map[string]gen.SdkInterfaceOverride{}
