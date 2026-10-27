package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func configFieldsObject(hooks ...ssot.ResourceHookOption) ssot.Object {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	autoResume := ssot.Attribute{Name: "auto_resume", Type: ssot.AttrBool}
	minNodes := ssot.Attribute{Name: "min_nodes", Type: ssot.AttrInt, Required: true}
	size := ssot.Attribute{Name: "warehouse_size", Type: ssot.AttrEnum, Required: true, Enum: ssot.NewEnum("WarehouseSize", "WarehouseSizes", "XSMALL")}
	backup := ssot.Attribute{
		Name: "backup_instance_families",
		Type: ssot.AttrList,
		Elem: &ssot.Attribute{Type: ssot.AttrString},
	}
	nested := ssot.Attribute{
		Name: "storage_location",
		Type: ssot.AttrList,
		Children: []ssot.Attribute{
			{Name: "url", Type: ssot.AttrString},
		},
	}
	monitor := ssot.Attribute{Name: "resource_monitor", Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier", Required: true}
	notShown := ssot.Attribute{Name: "catalog_namespace", Type: ssot.AttrString}
	createOnly := ssot.Attribute{Name: "initially_suspended", Type: ssot.AttrBool, CreateOnly: true, NativeBool: true}
	enabled := ssot.Attribute{Name: "enabled", Type: ssot.AttrBool, NativeBool: true}
	return ssot.Object{
		Name:   "Pool",
		IdType: "AccountObjectIdentifier",
		Attributes: []ssot.Attribute{
			comment,
			autoResume,
			minNodes,
			size,
			backup,
			nested,
			monitor,
			notShown,
			createOnly,
			enabled,
		},
		Show: ssot.OutputSpec{Columns: []*ssot.Attribute{
			{Name: "owner", Type: ssot.AttrString},
			{Name: "enabled", Type: ssot.AttrBool, NativeBool: true},
			{Name: "initially_suspended", Type: ssot.AttrBool, CreateOnly: true, NativeBool: true},
			{Name: "resource_monitor", Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier", Required: true},
			{Name: "storage_location", Type: ssot.AttrList, Children: nested.Children},
			{Name: "backup_instance_families", Type: ssot.AttrList, Elem: backup.Elem},
			{Name: "warehouse_size", Type: ssot.AttrEnum, Required: true, Enum: size.Enum, ShowName: "size"},
			{Name: "min_nodes", Type: ssot.AttrInt, Required: true},
			{Name: "auto_resume", Type: ssot.AttrBool},
			{Name: "comment", Type: ssot.AttrString, OptionalInShow: true},
		}},
		ResourceHooks: hooks,
	}
}

func Test_configFieldSets(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: configFieldsObject()}, nil, nil)
	require.Equal(t, []ConfigFieldSet{
		{Name: "comment", Expr: "obj.Comment"},
		{Name: "min_nodes", Expr: "obj.MinNodes"},
		{Name: "warehouse_size", Expr: "obj.Size"},
		{Name: "backup_instance_families", Expr: "obj.BackupInstanceFamilies"},
		{Name: "resource_monitor", Expr: "obj.ResourceMonitor.FullyQualifiedName()"},
		{Name: "enabled", Expr: "obj.Enabled"},
	}, model.ConfigFieldSets())
}

func Test_configFieldCode(t *testing.T) {
	obj := ssot.Object{Name: "Pool"}
	autoResume := ssot.Attribute{Name: "auto_resume", Type: ssot.AttrBool}
	require.Equal(t, "booleanStringFromBool(obj.AutoResume)", configFieldCode(obj, autoResume))
	autoResume.OptionalInShow = true
	require.Equal(t, "obj.AutoResume", configFieldCode(obj, autoResume))

	monitor := ssot.Attribute{Name: "resource_monitor", Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier"}
	require.Equal(t, "obj.ResourceMonitor.FullyQualifiedName()", configFieldCode(obj, monitor))
	monitor.OptionalInShow = true
	require.Equal(t, "obj.ResourceMonitor", configFieldCode(obj, monitor))
}

func Test_setConfigFields_template(t *testing.T) {
	t.Run("generated", func(t *testing.T) {
		model := ModelFromInputObject(ResourceObject{Object: configFieldsObject(ssot.WithSetConfigFields(ssot.Generated()))}, nil, nil)
		got := renderTemplate(t, "generatedSetConfigFields", model)
		require.Contains(t, got, "func poolSetConfigFields(d *tfschema.ResourceData, obj *sdk.Pool) error {")
		require.Contains(t, got, "return errors.Join(")
		require.Contains(t, got, `d.Set("comment", obj.Comment),`)
		require.Contains(t, got, `d.Set("warehouse_size", obj.Size),`)
		require.Contains(t, got, `d.Set("resource_monitor", obj.ResourceMonitor.FullyQualifiedName()),`)
		require.NotContains(t, got, "var resourceMonitor")
		call := renderTemplate(t, "setConfigFields", model)
		require.Contains(t, call, "poolSetConfigFields(d, pool),")
		require.NotContains(t, call, "SetConfigFieldsExt")
		lifecycle := renderLifecycle(t, model)
		require.Contains(t, lifecycle, "func poolSetConfigFields(d *tfschema.ResourceData, obj *sdk.Pool) error {")
		require.NotContains(t, lifecycle, "poolSetConfigFieldsExt")
	})

	t.Run("generated variant show type", func(t *testing.T) {
		model := lifecycleModel(func(obj *ssot.Object) {
			obj.VariantOf = &ssot.Object{Name: "Warehouse"}
			obj.ResourceHooks = []ssot.ResourceHookOption{ssot.WithSetConfigFields(ssot.Generated())}
		})
		require.Contains(t, renderTemplate(t, "generatedSetConfigFields", model), "obj *sdk.Warehouse")
		require.Contains(t, renderTemplate(t, "generatedSetConfigFields", model), "return nil")
	})

	t.Run("generated empty", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithSetConfigFields(ssot.Generated())))
		got := renderTemplate(t, "generatedSetConfigFields", model)
		require.Contains(t, got, "func poolSetConfigFields(d *tfschema.ResourceData, obj *sdk.Pool) error {")
		require.Contains(t, got, "return nil")
		require.NotContains(t, got, "errors.Join")
	})

	t.Run("ext call", func(t *testing.T) {
		model := lifecycleModel()
		require.Contains(t, renderTemplate(t, "setConfigFields", model), "poolSetConfigFieldsExt(d, pool),")
		require.NotContains(t, renderLifecycle(t, model), "func poolSetConfigFields(")
	})

	t.Run("omitted", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithSetConfigFields(ssot.Omitted())))
		require.NotContains(t, renderTemplate(t, "setConfigFields", model), "SetConfigFields")
		require.NotContains(t, renderLifecycle(t, model), "SetConfigFields")
	})
}
