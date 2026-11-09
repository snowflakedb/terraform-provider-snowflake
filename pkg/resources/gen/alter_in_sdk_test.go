package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

func Test_AlterInSdk_computePool(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.ComputePool}, nil, nil)
	require.True(t, model.Hooks.AlterSetInSdk.Generate)
	require.True(t, model.Hooks.AlterUnsetInSdk.Generate)

	got := renderLifecycle(t, model)
	require.Contains(t, got, `if !reflect.DeepEqual(*setUnset.set, *sdk.NewComputePoolSetRequest()) {
		if err := meta.(*provider.Context).Client.ComputePools.Alter(ctx, sdk.NewAlterComputePoolRequest(id).WithSet(*setUnset.set)); err != nil {
			return diag.FromErr(err)
		}
	}`)
	require.Contains(t, got, `if !reflect.DeepEqual(*setUnset.unset, *sdk.NewComputePoolUnsetRequest()) {
		if err := meta.(*provider.Context).Client.ComputePools.Alter(ctx, sdk.NewAlterComputePoolRequest(id).WithUnset(*setUnset.unset)); err != nil {
			return diag.FromErr(err)
		}
	}`)
	require.NotContains(t, got, "computePoolAlterSetInSdkExt")
	require.NotContains(t, got, "computePoolAlterUnsetInSdkExt")
}

func Test_AlterInSdk_hookModes(t *testing.T) {
	tests := []struct {
		name             string
		hook             func(ssot.ResourceHook) ssot.ResourceHookOption
		generated        string
		generatedPartial string
		ext              string
		extPartial       string
	}{
		{
			name: "alterSetInSdk", hook: ssot.WithAlterSetInSdk,
			generated: `if !reflect.DeepEqual(*setUnset.set, *sdk.NewPoolSetRequest()) {
		if err := meta.(*provider.Context).Client.Pools.Alter(ctx, sdk.NewAlterPoolRequest(id).WithSet(*setUnset.set)); err != nil {
			return diag.FromErr(err)
		}
	}
`,
			generatedPartial: `if !reflect.DeepEqual(*setUnset.set, *sdk.NewPoolSetRequest()) {
		if err := meta.(*provider.Context).Client.Pools.Alter(ctx, sdk.NewAlterPoolRequest(id).WithSet(*setUnset.set)); err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}
`,
			ext: `if err := poolAlterSetInSdkExt(ctx, d, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
`,
			extPartial: `if err := poolAlterSetInSdkExt(ctx, d, meta, id, setUnset); err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}
`,
		},
		{
			name: "alterUnsetInSdk", hook: ssot.WithAlterUnsetInSdk,
			generated: `if !reflect.DeepEqual(*setUnset.unset, *sdk.NewPoolUnsetRequest()) {
		if err := meta.(*provider.Context).Client.Pools.Alter(ctx, sdk.NewAlterPoolRequest(id).WithUnset(*setUnset.unset)); err != nil {
			return diag.FromErr(err)
		}
	}
`,
			generatedPartial: `if !reflect.DeepEqual(*setUnset.unset, *sdk.NewPoolUnsetRequest()) {
		if err := meta.(*provider.Context).Client.Pools.Alter(ctx, sdk.NewAlterPoolRequest(id).WithUnset(*setUnset.unset)); err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}
`,
			ext: `if err := poolAlterUnsetInSdkExt(ctx, meta, id, setUnset); err != nil {
		return diag.FromErr(err)
	}
`,
			extPartial: `if err := poolAlterUnsetInSdkExt(ctx, meta, id, setUnset); err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, mode := range []struct {
				name        string
				hook        ssot.ResourceHook
				want        string
				wantPartial string
			}{
				{"generated", ssot.Generated(), tt.generated, tt.generatedPartial},
				{"ext call", ssot.ExtCall(), tt.ext, tt.extPartial},
			} {
				t.Run(mode.name, func(t *testing.T) {
					model := lifecycleModel(withHooks(tt.hook(mode.hook)))
					require.Equal(t, mode.want, renderTemplate(t, tt.name, model))
					t.Run("partial state", func(t *testing.T) {
						model.HasPartialOnUpdate = true
						require.Equal(t, mode.wantPartial, renderTemplate(t, tt.name, model))
					})
				})
			}
			t.Run("omitted", func(t *testing.T) {
				model := lifecycleModel(withHooks(tt.hook(ssot.Omitted())))
				require.Empty(t, renderTemplate(t, tt.name, model))
			})
		})
	}
}

func Test_AlterInSdk_statementsAreSentBeforeTheMainOne(t *testing.T) {
	model := lifecycleModel(withHooks(ssot.WithAlterSetInSdk(ssot.Generated()), ssot.WithAlterUnsetInSdk(ssot.Generated())))
	model.HasPartialOnUpdate = true
	model.AlterStatements = AlterStatementsModel{Set: []string{"setBody", "set"}, Unset: []string{"unsetSchedule", "unset"}}

	require.Equal(t, `if !reflect.DeepEqual(*setUnset.setBody, *sdk.NewPoolSetRequest()) {
		if err := meta.(*provider.Context).Client.Pools.Alter(ctx, sdk.NewAlterPoolRequest(id).WithSet(*setUnset.setBody)); err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}
if !reflect.DeepEqual(*setUnset.set, *sdk.NewPoolSetRequest()) {
		if err := meta.(*provider.Context).Client.Pools.Alter(ctx, sdk.NewAlterPoolRequest(id).WithSet(*setUnset.set)); err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}
`, renderTemplate(t, "alterSetInSdk", model))
	require.Equal(t, `if !reflect.DeepEqual(*setUnset.unsetSchedule, *sdk.NewPoolUnsetRequest()) {
		if err := meta.(*provider.Context).Client.Pools.Alter(ctx, sdk.NewAlterPoolRequest(id).WithUnset(*setUnset.unsetSchedule)); err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}
if !reflect.DeepEqual(*setUnset.unset, *sdk.NewPoolUnsetRequest()) {
		if err := meta.(*provider.Context).Client.Pools.Alter(ctx, sdk.NewAlterPoolRequest(id).WithUnset(*setUnset.unset)); err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}
`, renderTemplate(t, "alterUnsetInSdk", model))
}

func Test_ModelFromInputObject_alterStatements(t *testing.T) {
	body := ssot.Attribute{Name: "body", Type: ssot.AttrString}
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	schedule := ssot.Attribute{Name: "schedule", Type: ssot.AttrString}
	expire := ssot.Attribute{Name: "expire_after_days", Type: ssot.AttrInt}
	model := ModelFromInputObject(ResourceObject{Object: ssot.Object{
		Name:           "Widget",
		NamePlural:     "Widgets",
		IdType:         "AccountObjectIdentifier",
		Attributes:     []ssot.Attribute{body, comment, schedule, expire},
		UnsetConflicts: [][]*ssot.Attribute{{&schedule, &expire}},
		Alter: ssot.AlterSpec{
			Clauses: []ssot.AlterClause{
				ssot.AlterSet(
					ssot.AlterGroup("Schedule", &schedule),
					ssot.AlterGroup("Body", &body),
					ssot.AlterGroup(ssot.MainAlterGroup, &comment),
				),
				ssot.AlterUnset(
					ssot.AlterGroup("Schedule", &schedule),
					ssot.AlterGroup(ssot.MainAlterGroup, &expire),
				),
			},
			SetConflictingFields: [][]*ssot.Attribute{{&body, &comment}},
		},
	}}, nil, nil)
	require.Equal(t, AlterStatementsModel{
		Set:   []string{"setSchedule", "setBody", "set"},
		Unset: []string{"unsetSchedule", "unset"},
		Requests: []AlterStatementRequest{
			{"setSchedule", "Set"}, {"setBody", "Set"}, {"set", "Set"}, {"unsetSchedule", "Unset"}, {"unset", "Unset"},
		},
	}, model.AlterStatements)
}

func Test_ModelFromInputObject_withoutStatementsSendsTheMainOnes(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.ComputePool}, nil, nil)
	require.Equal(t, []string{"set"}, model.AlterStatements.Set)
	require.Equal(t, []string{"unset"}, model.AlterStatements.Unset)
}

func Test_ModelFromInputObject_alterStatementsFollowDeclarationOrder(t *testing.T) {
	body := ssot.Attribute{Name: "body", Type: ssot.AttrString}
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	newModel := func(clauses ...ssot.AlterClause) AlterStatementsModel {
		return ModelFromInputObject(ResourceObject{Object: ssot.Object{
			Name:       "Widget",
			NamePlural: "Widgets",
			IdType:     "AccountObjectIdentifier",
			Attributes: []ssot.Attribute{body, comment},
			Alter:      ssot.AlterSpec{Clauses: clauses},
		}}, nil, nil).AlterStatements
	}

	// The main group is sent where it is declared.
	model := newModel(ssot.AlterSet(ssot.AlterGroup(ssot.MainAlterGroup, &comment), ssot.AlterGroup("Body", &body)))
	require.Equal(t, []string{"set", "setBody"}, model.Set)
	// An object whose ALTER spec has no UNSET block has no unset request, because the SDK has no such type.
	require.Empty(t, model.Unset)

	// A block of named groups only still has the main request, last, for the parameters.
	model = newModel(ssot.AlterSet(ssot.AlterGroup("Body", &body), ssot.AlterGroup("Comment", &comment)))
	require.Equal(t, []string{"setBody", "setComment", "set"}, model.Set)

	// Attributes passed directly are one main group.
	model = newModel(ssot.AlterSet(&body, &comment), ssot.AlterUnset(&comment))
	require.Equal(t, []string{"set"}, model.Set)
	require.Equal(t, []string{"unset"}, model.Unset)
}
