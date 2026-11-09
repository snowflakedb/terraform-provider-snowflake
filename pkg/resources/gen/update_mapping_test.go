package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

func Test_ComputePoolUpdateMappings_renderApplySetUnset(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.ComputePool}, nil, nil)
	require.True(t, model.Hooks.ApplySetUnset.Generate)
	require.Equal(t, "intAttributeUpdateSetOnlyBuilder", model.Attributes[1].UpdateCall.Helper)
	require.Equal(t, `setUnset, err := computePoolApplySetUnset(d)
	if err != nil {
		return diag.FromErr(err)
	}
`, renderTemplate(t, "applySetUnset", model))

	require.Equal(t, `func computePoolApplySetUnset(d *tfschema.ResourceData) (*computePoolAlterRequestsExt, error) {
	set, unset := sdk.NewComputePoolSetRequest(), sdk.NewComputePoolUnsetRequest()
	errs := errors.Join(
		intAttributeUpdateSetOnlyBuilder(d, "min_nodes", set.WithMinNodes),
		intAttributeUpdateSetOnlyBuilder(d, "max_nodes", set.WithMaxNodes),
		listValueUpdateBuilder(d, "backup_instance_families", set.WithBackupInstanceFamilies, unset.WithBackupInstanceFamilies, toComputePoolBackupInstanceFamilies),
		booleanStringAttributeUpdateBuilder(d, "auto_resume", set.WithAutoResume, unset.WithAutoResume),
		intAttributeWithSpecialDefaultUpdateBuilder(d, "auto_suspend_secs", set.WithAutoSuspendSecs, unset.WithAutoSuspendSecs),
		stringAttributeUpdateBuilder(d, "comment", set.WithComment, unset.WithComment),
	)
	if errs != nil {
		return nil, errs
	}
	return &computePoolAlterRequestsExt{set: set, unset: unset}, nil
}
`, renderTemplate(t, "generatedApplySetUnset", model))
}

func Test_attachUpdateMappings_skipCreateStays(t *testing.T) {
	src := []ssot.Attribute{{Name: "enabled", Type: ssot.AttrBool, NativeBool: true, SkipCreate: true}}
	dst := applyResourceAttributeRules("", src)
	attachUpdateMappings(ssot.Object{Name: "Widget"}, src, dst)
	require.Equal(t, "booleanAttributeUpdateSetOnlyBuilder", dst[0].UpdateCall.Helper)
	require.Equal(t, "set.WithEnabled", dst[0].UpdateCall.Target)
}

func Test_attachUpdateMappings_nativeBoolIsSetOnly(t *testing.T) {
	src := []ssot.Attribute{{Name: "change_tracking", Type: ssot.AttrBool, NativeBool: true}}
	dst := applyResourceAttributeRules("", src)
	attachUpdateMappings(ssot.Object{Name: "Table"}, src, dst)
	require.Equal(t, "booleanAttributeUpdateSetOnlyBuilder", dst[0].UpdateCall.Helper)
	require.Equal(t, "set.WithChangeTracking", dst[0].UpdateCall.Target)
}

func Test_attachUpdateMappings_skipsNested(t *testing.T) {
	src := []ssot.Attribute{{
		Name:          "client_policy",
		Type:          ssot.AttrSet,
		UpdateMapping: ssot.StringAttributeUpdateBuilder,
		Elem: ssot.ObjectElem([]ssot.Attribute{{
			Name:          "minimum_version",
			Type:          ssot.AttrString,
			UpdateMapping: ssot.StringAttributeUpdateBuilder,
		}}),
	}}
	dst := applyResourceAttributeRules("", src)
	attachUpdateMappings(ssot.Object{Name: "AuthenticationPolicy"}, src, dst)
	require.Nil(t, dst[0].UpdateCall)
	require.Nil(t, dst[0].Children[0].UpdateCall)
}

func Test_SchemaUpdateMappings_stayExt(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.Schema}, nil, nil)
	require.False(t, model.Hooks.ApplySetUnset.Generate)
	require.True(t, model.Hooks.ApplySetUnset.Present)
	require.Equal(t, `setUnset, err := schemaApplySetUnsetExt(d)
	if err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}
`, renderTemplate(t, "applySetUnset", model))
}

func Test_attachUpdateMappings_routesToStatements(t *testing.T) {
	src := []ssot.Attribute{
		{Name: "objects", Type: ssot.AttrString, UpdateMapping: ssot.StringAttributeUpdateBuilder},
		{Name: "comment", Type: ssot.AttrString},
		{Name: "schedule", Type: ssot.AttrString, UpdateMapping: ssot.StringAttributeUpdateBuilder},
	}
	obj := ssot.Object{
		Name:           "Widget",
		UnsetConflicts: [][]*ssot.Attribute{{&src[0], &src[2]}},
		Alter: ssot.AlterSpec{
			Clauses: []ssot.AlterClause{
				ssot.AlterSet(ssot.AlterGroup("Objects", &src[0]), ssot.AlterGroup(ssot.MainAlterGroup, &src[1], &src[2])),
				ssot.AlterUnset(ssot.AlterGroup("Objects", &src[0]), ssot.AlterGroup(ssot.MainAlterGroup, &src[1], &src[2])),
			},
		},
	}
	dst := applyResourceAttributeRules("", src)
	attachUpdateMappings(obj, src, dst)
	// An attribute in a named group routes each block to its own statement, and the others stay in the main ones.
	require.Equal(t, "setObjects.WithObjects, unsetObjects.WithObjects", dst[0].UpdateCall.Target)
	require.Equal(t, "set.WithComment, unset.WithComment", dst[1].UpdateCall.Target)
	require.Equal(t, "set.WithSchedule, unset.WithSchedule", dst[2].UpdateCall.Target)
}

func Test_WidgetApplySetUnset_withStatements(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	objects := ssot.Attribute{Name: "objects", Type: ssot.AttrString, UpdateMapping: ssot.IntAttributeUpdateSetOnlyBuilder}
	model := ModelFromInputObject(ResourceObject{Object: ssot.Object{
		Name:       "Widget",
		NamePlural: "Widgets",
		IdType:     "AccountObjectIdentifier",
		Attributes: []ssot.Attribute{objects, comment},
		ResourceHooks: []ssot.ResourceHookOption{
			ssot.WithApplySetUnset(ssot.Generated()),
		},
		Alter: ssot.AlterSpec{
			Clauses: []ssot.AlterClause{ssot.AlterSet(ssot.AlterGroup("Objects", &objects), ssot.AlterGroup(ssot.MainAlterGroup, &comment))},
		},
	}}, nil, nil)
	require.Equal(t, `func widgetApplySetUnset(d *tfschema.ResourceData) (*widgetAlterRequestsExt, error) {
	setObjects, set, unset := sdk.NewWidgetSetRequest(), sdk.NewWidgetSetRequest(), sdk.NewWidgetUnsetRequest()
	errs := errors.Join(
		intAttributeUpdateSetOnlyBuilder(d, "objects", setObjects.WithObjects),
		stringAttributeUpdateBuilder(d, "comment", set.WithComment, unset.WithComment),
	)
	if errs != nil {
		return nil, errs
	}
	return &widgetAlterRequestsExt{setObjects: setObjects, set: set, unset: unset}, nil
}
`, renderTemplate(t, "generatedApplySetUnset", model))
}
