package gen

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_SchemaStateUpgraders_versionFollowsSliceOrder(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: ssot.Object{
		Name: "Widget",
		SchemaStateUpgraders: []ssot.StateUpgraderFunc{
			"firstUpgrader",
			"secondUpgrader",
		},
		ResourceHooks: []ssot.ResourceHookOption{
			ssot.WithStateUpgraders(ssot.Generated()),
		},
	}}, nil, nil)

	got := renderTemplate(t, "constructor", model)
	require.Contains(t, got, `		SchemaVersion: 2,
		// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
		StateUpgraders: []tfschema.StateUpgrader{
			{
				Version: 0,
				Type:    cty.EmptyObject,
				Upgrade: firstUpgrader,
			},
			{
				Version: 1,
				Type:    cty.EmptyObject,
				Upgrade: secondUpgrader,
			},
		},
`)
	require.NotContains(t, got, "StateUpgradersExt")
	require.NotContains(t, got, "SchemaVersionExt")
}

func Test_SchemaVersionHook_staysExtWhenUpgradersUnset(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: ssot.Object{
		Name: "Widget",
		ResourceHooks: []ssot.ResourceHookOption{
			ssot.WithStateUpgraders(ssot.ExtCall()),
		},
	}}, nil, nil)

	got := renderTemplate(t, "constructor", model)
	require.Contains(t, got, "SchemaVersion:  widgetSchemaVersionExt,")
	require.Contains(t, got, "StateUpgraders: widgetStateUpgradersExt,")
	require.Less(t, strings.Index(got, "SchemaVersion:"), strings.Index(got, "StateUpgraders:"))
}
