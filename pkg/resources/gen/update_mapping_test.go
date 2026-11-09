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

	require.Equal(t, `func computePoolApplySetUnset(d *tfschema.ResourceData) (*computePoolAlterRequests, error) {
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
	return &computePoolAlterRequests{set: set, unset: unset}, nil
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
			Clauses: []ssot.AlterClause{ssot.AlterSet(ssot.AlterGroup("Objects", &objects), ssot.AlterGroup(ssot.MainAlterGroup, &comment)), ssot.AlterUnset(&comment)},
		},
	}}, nil, nil)
	require.Equal(t, `func widgetApplySetUnset(d *tfschema.ResourceData) (*widgetAlterRequests, error) {
	setObjects, set, unset := sdk.NewWidgetSetRequest(), sdk.NewWidgetSetRequest(), sdk.NewWidgetUnsetRequest()
	errs := errors.Join(
		intAttributeUpdateSetOnlyBuilder(d, "objects", setObjects.WithObjects),
		stringAttributeUpdateBuilder(d, "comment", set.WithComment, unset.WithComment),
	)
	if errs != nil {
		return nil, errs
	}
	return &widgetAlterRequests{setObjects: setObjects, set: set, unset: unset}, nil
}
`, renderTemplate(t, "generatedApplySetUnset", model))
}

func Test_attachUpdateMappings_fieldsStyleFollowsTheGroups(t *testing.T) {
	src := []ssot.Attribute{
		{Name: "error_integration", Type: ssot.AttrIdentifier, UpdateMapping: ssot.AccountObjectIdentifierAttributeUpdate},
		{Name: "comment", Type: ssot.AttrString},
	}
	obj := ssot.Object{
		Name: "Widget",
		Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
			ssot.AlterSet(ssot.AlterGroup("Integration", &src[0]), ssot.AlterGroup(ssot.MainAlterGroup, &src[1])),
			ssot.AlterUnset(&src[0], &src[1]),
		}},
	}
	dst := applyResourceAttributeRules("", src)
	attachUpdateMappings(obj, src, dst)
	require.Equal(t, "&setIntegration.ErrorIntegration, &unset.ErrorIntegration", dst[0].UpdateCall.Target)
	require.Equal(t, "set.WithComment, unset.WithComment", dst[1].UpdateCall.Target)
}

func Test_attachUpdateMappings_nestedAttributeUpdatedAsOneValue(t *testing.T) {
	child := ssot.Attribute{Name: "hours", Type: ssot.AttrInt}
	src := []ssot.Attribute{
		{Name: "schedule", Type: ssot.AttrObject, Children: []ssot.Attribute{child}, UpdateMapping: ssot.AttributeMappedValueUpdate("toWidgetSchedule")},
		{Name: "window", Type: ssot.AttrObject, Children: []ssot.Attribute{child}, UpdateMapping: ssot.StringAttributeUpdateBuilder},
	}
	obj := ssot.Object{
		Name: "Widget",
		Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
			ssot.AlterSet(&src[0], &src[1]),
			ssot.AlterUnset(ssot.AlterGroup("Schedule", &src[0]), ssot.AlterGroup(ssot.MainAlterGroup, &src[1])),
		}},
	}
	dst := applyResourceAttributeRules("", src)
	attachUpdateMappings(obj, src, dst)
	// The whole-value mapping is generated, and its unset goes to the request of its group.
	require.Equal(t, "attributeMappedValueUpdate", dst[0].UpdateCall.Helper)
	require.Equal(t, "&set.Schedule, &unsetSchedule.Schedule, toWidgetSchedule", dst[0].UpdateCall.Target)
	// Any other nested attribute stays skipped, even with a mapping.
	require.Nil(t, dst[1].UpdateCall)
}

func Test_WidgetAlterRequests_oneFieldPerGroup(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	objects := ssot.Attribute{Name: "objects", Type: ssot.AttrString}
	newModel := func(applySetUnset ssot.ResourceHookOption, clauses ...ssot.AlterClause) ResourceModel {
		return ModelFromInputObject(ResourceObject{Object: ssot.Object{
			Name:          "Widget",
			NamePlural:    "Widgets",
			IdType:        "AccountObjectIdentifier",
			Attributes:    []ssot.Attribute{objects, comment},
			ResourceHooks: []ssot.ResourceHookOption{applySetUnset},
			Alter:         ssot.AlterSpec{Clauses: clauses},
		}}, nil, nil)
	}

	model := newModel(
		ssot.WithApplySetUnset(ssot.Generated()),
		ssot.AlterSet(ssot.AlterGroup("Objects", &objects), ssot.AlterGroup(ssot.MainAlterGroup, &comment)),
		ssot.AlterUnset(&objects, &comment),
	)
	require.Equal(t, "widgetAlterRequests", model.AlterRequestsType())
	require.Equal(t, `// widgetAlterRequests holds the SET and UNSET requests of the generated update, one per group of the ALTER spec.
type widgetAlterRequests struct {
	setObjects *sdk.WidgetSetRequest
	set *sdk.WidgetSetRequest
	unset *sdk.WidgetUnsetRequest
}
`, renderTemplate(t, "generatedAlterRequests", model))

	// An object that keeps ApplySetUnset in the ext file declares the struct there, and the name says so.
	model = newModel(ssot.WithApplySetUnset(ssot.ExtCall()), ssot.AlterSet(&comment))
	require.Equal(t, "widgetAlterRequestsExt", model.AlterRequestsType())
}

func Test_attachUpdateMappings_deducesEnumsIdentifiersAndLists(t *testing.T) {
	size := ssot.Enum{Name: "WidgetSize"}
	src := []ssot.Attribute{
		{Name: "size", Type: ssot.AttrEnum, Enum: &size},
		{Name: "monitor", Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier"},
		{Name: "databases", Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier"}},
		{Name: "kinds", Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrEnum, Enum: &size}},
	}
	obj := ssot.Object{
		Name: "Widget",
		Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
			ssot.AlterSet(&src[0], &src[1], &src[2], &src[3]),
			// databases and kinds are not in the UNSET block, so they can only be set.
			ssot.AlterUnset(&src[0], &src[1]),
		}},
	}
	dst := applyResourceAttributeRules("", src)
	attachUpdateMappings(obj, src, dst)

	require.Equal(t, "attributeMappedValueUpdate", dst[0].UpdateCall.Helper)
	require.Equal(t, "&set.Size, &unset.Size, sdk.ToWidgetSize", dst[0].UpdateCall.Target)
	require.Equal(t, "accountObjectIdentifierAttributeUpdate", dst[1].UpdateCall.Helper)
	require.Equal(t, "&set.Monitor, &unset.Monitor", dst[1].UpdateCall.Target)
	require.Equal(t, "listValueUpdateSetOnlyBuilder", dst[2].UpdateCall.Helper)
	require.Equal(t, "set.WithDatabases, toWidgetDatabases", dst[2].UpdateCall.Target)
	require.Equal(t, "listValueUpdateSetOnlyBuilder", dst[3].UpdateCall.Helper)
	require.Equal(t, "set.WithKinds, toWidgetKinds", dst[3].UpdateCall.Target)
}

func Test_attachUpdateMappings_updatesOnlyTheAttributesOfTheBlocks(t *testing.T) {
	src := []ssot.Attribute{
		{Name: "comment", Type: ssot.AttrString},
		{Name: "with_managed_access", Type: ssot.AttrBool},
	}
	obj := ssot.Object{
		Name:  "Widget",
		Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{ssot.AlterSet(&src[0]), ssot.AlterUnset(&src[0])}},
	}
	dst := applyResourceAttributeRules("", src)
	attachUpdateMappings(obj, src, dst)
	require.NotNil(t, dst[0].UpdateCall)
	// with_managed_access is in neither block, so the update leaves it to a hook and the model does not need a SkipUpdate.
	require.Nil(t, dst[1].UpdateCall)

	// An object without SET and UNSET blocks keeps updating every attribute that has a helper.
	dst = applyResourceAttributeRules("", src[:1])
	attachUpdateMappings(ssot.Object{Name: "Widget"}, src[:1], dst)
	require.NotNil(t, dst[0].UpdateCall)

	// A mapping on an attribute that is in no block is a mistake.
	src[1].UpdateMapping = ssot.StringAttributeUpdateBuilder
	dst = applyResourceAttributeRules("", src)
	require.PanicsWithValue(t, `object "Widget": attribute "with_managed_access" has an UpdateMapping, but it is in neither the SET nor the UNSET block`, func() {
		attachUpdateMappings(obj, src, dst)
	})
}

func Test_attachUpdateMappings_failsWhenAnAttributeOfABlockHasNoHelper(t *testing.T) {
	src := []ssot.Attribute{{Name: "size", Type: ssot.AttrEnum}}
	obj := ssot.Object{
		Name:  "Widget",
		Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{ssot.AlterSet(&src[0]), ssot.AlterUnset(&src[0])}},
	}
	dst := applyResourceAttributeRules("", src)
	// The enum has no SDK name to build the mapper from, so nothing can be deduced.
	require.PanicsWithValue(t, `object "Widget": attribute "size" is in the SET or UNSET block, but no update helper can be deduced for it, so give it an UpdateMapping or SkipUpdate`, func() {
		attachUpdateMappings(obj, src, dst)
	})

	// SkipUpdate says another hook updates it.
	src[0].SkipUpdate = true
	dst = applyResourceAttributeRules("", src)
	attachUpdateMappings(obj, src, dst)
	require.Nil(t, dst[0].UpdateCall)
}

func Test_attachUpdateMappings_setOnlyPerType(t *testing.T) {
	src := []ssot.Attribute{
		{Name: "enabled", Type: ssot.AttrBool},
		{Name: "count", Type: ssot.AttrInt},
	}
	obj := ssot.Object{
		Name: "Widget",
		// Neither attribute is in the UNSET block, so both are only ever set.
		Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{ssot.AlterSet(&src[0], &src[1]), ssot.AlterUnsetTags()}},
	}
	dst := applyResourceAttributeRules("", src)
	require.Equal(t, "BooleanDefault", dst[0].DefaultExpr)
	attachUpdateMappings(obj, src, dst)
	require.Equal(t, "booleanStringAttributeUpdateSetOnly", dst[0].UpdateCall.Helper)
	require.Equal(t, "&set.Enabled", dst[0].UpdateCall.Target)
	require.Equal(t, "intAttributeUpdateSetOnlyBuilder", dst[1].UpdateCall.Helper)
	require.Equal(t, "set.WithCount", dst[1].UpdateCall.Target)

	// A set-only identifier has no deduced helper, so generation asks for a mapping instead of writing an UNSET that does not exist.
	identifier := []ssot.Attribute{{Name: "monitor", Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier"}}
	setOnlyObject := ssot.Object{
		Name:  "Widget",
		Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{ssot.AlterSet(&identifier[0])}},
	}
	dst = applyResourceAttributeRules("", identifier)
	require.PanicsWithValue(t, `object "Widget": attribute "monitor" is in the SET or UNSET block, but no update helper can be deduced for it, so give it an UpdateMapping or SkipUpdate`, func() {
		attachUpdateMappings(setOnlyObject, identifier, dst)
	})
}
