package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func readStateObject(hooks ...ssot.ResourceHookOption) ssot.Object {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	autoResume := ssot.Attribute{Name: "auto_resume", Type: ssot.AttrBool}
	autoSuspend := ssot.Attribute{Name: "auto_suspend", Type: ssot.AttrInt}
	initiallySuspended := ssot.Attribute{Name: "initially_suspended", Type: ssot.AttrBool, CreateOnly: true, NativeBool: true}
	suspendOnCreate := ssot.Attribute{Name: "suspend_on_create", Type: ssot.AttrBool, CreateOnly: true}
	methods := ssot.Attribute{
		Name: "authentication_methods",
		Type: ssot.AttrSet,
		Elem: &ssot.Attribute{Type: ssot.AttrString},
	}
	intSet := ssot.Attribute{
		Name: "allowed_ports",
		Type: ssot.AttrSet,
		Elem: &ssot.Attribute{Type: ssot.AttrInt},
	}
	boolSet := ssot.Attribute{
		Name: "flags",
		Type: ssot.AttrSet,
		Elem: &ssot.Attribute{Type: ssot.AttrBool},
	}
	backup := ssot.Attribute{
		Name: "backup_instance_families",
		Type: ssot.AttrList,
		Elem: &ssot.Attribute{Type: ssot.AttrString},
	}
	nested := ssot.Attribute{
		Name: "mfa_policy",
		Type: ssot.AttrList,
		Children: []ssot.Attribute{
			{Name: "allowed_methods", Type: ssot.AttrSet, Elem: &ssot.Attribute{Type: ssot.AttrString}},
		},
	}
	minNodes := ssot.Attribute{Name: "min_nodes", Type: ssot.AttrInt, Required: true}
	nativeBool := ssot.Attribute{Name: "enabled", Type: ssot.AttrBool, NativeBool: true}
	return ssot.Object{
		Name:   "Pool",
		IdType: "AccountObjectIdentifier",
		Attributes: []ssot.Attribute{
			comment,
			autoResume,
			autoSuspend,
			initiallySuspended,
			suspendOnCreate,
			methods,
			intSet,
			boolSet,
			backup,
			nested,
			minNodes,
			nativeBool,
		},
		ResourceHooks: hooks,
	}
}

func Test_setStateToValueFromConfigFields(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: readStateObject()}, nil, nil)
	require.Equal(t, []string{
		"auto_resume",
		"auto_suspend",
		"authentication_methods",
		"allowed_ports",
		"flags",
	}, model.SetStateToValueFromConfigFields())
}

func Test_setStateToValuesFromConfig_template(t *testing.T) {
	t.Run("generated", func(t *testing.T) {
		model := ModelFromInputObject(ResourceObject{Object: readStateObject()}, nil, nil)
		got := renderTemplate(t, "setStateToValuesFromConfig", model)
		require.Contains(t, got, "setStateToValuesFromConfig(d, poolSchema, []string{")
		require.Contains(t, got, `"auto_resume",`)
		require.Contains(t, got, `"auto_suspend",`)
		require.Contains(t, got, "return diag.FromErr(err)")
		require.NotContains(t, got, "FieldsExt")
		lifecycle := renderLifecycle(t, model)
		require.Contains(t, lifecycle, `"auto_resume",`)
		require.NotContains(t, lifecycle, "poolSetStateToValueFromConfigFields")
	})

	t.Run("generated empty", func(t *testing.T) {
		model := lifecycleModel()
		require.Empty(t, model.SetStateToValueFromConfigFields())
		require.NotContains(t, renderTemplate(t, "setStateToValuesFromConfig", model), "setStateToValuesFromConfig(")
		require.NotContains(t, renderLifecycle(t, model), "setStateToValuesFromConfig(")
	})

	t.Run("generated describe diagnostics", func(t *testing.T) {
		model := ModelFromInputObject(ResourceObject{Object: readStateObject(
			ssot.WithHandleExternalChangesInDescribe(ssot.ExtCall()),
		)}, nil, nil)
		require.Contains(t, renderTemplate(t, "setStateToValuesFromConfig", model), "return append(diags, diag.FromErr(err)...)")
	})

	t.Run("ext call", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithSetStateToValuesFromConfig(ssot.ExtCall())))
		got := renderTemplate(t, "setStateToValuesFromConfig", model)
		require.Contains(t, got, "setStateToValuesFromConfig(d, poolSchema, poolSetStateToValueFromConfigFieldsExt)")
		require.Contains(t, got, "return diag.FromErr(err)")
		require.NotContains(t, got, "[]string{")
	})

	t.Run("omitted", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithSetStateToValuesFromConfig(ssot.Omitted())))
		require.NotContains(t, renderTemplate(t, "setStateToValuesFromConfig", model), "setStateToValuesFromConfig(")
		require.NotContains(t, renderLifecycle(t, model), "setStateToValuesFromConfig(")
	})
}
