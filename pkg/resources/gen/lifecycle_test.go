package gen

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func lifecycleModel(opts ...func(*ssot.Object)) ResourceModel {
	obj := ssot.Object{Name: "Pool", NamePlural: "Pools", IdType: "SchemaObjectIdentifier"}
	for _, opt := range opts {
		opt(&obj)
	}
	return ModelFromInputObject(ResourceObject{Object: obj}, nil, nil)
}

func withHooks(hooks ...ssot.ResourceHookOption) func(*ssot.Object) {
	return func(obj *ssot.Object) { obj.ResourceHooks = hooks }
}

func renderLifecycle(t *testing.T, model ResourceModel) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, LifecycleTemplate.Execute(&b, model))
	return b.String()
}

func Test_lifecycle_setIdEncoding(t *testing.T) {
	t.Run("current id", func(t *testing.T) {
		got := renderLifecycle(t, lifecycleModel())
		require.Contains(t, got, "d.SetId(helpers.EncodeResourceIdentifier(id))")
		require.NotContains(t, got, "EncodeSnowflakeID")
	})

	t.Run("legacy id", func(t *testing.T) {
		got := renderLifecycle(t, lifecycleModel(func(obj *ssot.Object) { obj.LegacyIdEncoding = true }))
		require.Contains(t, got, "d.SetId(helpers.EncodeSnowflakeID(id))")
		require.NotContains(t, got, "EncodeResourceIdentifier")
	})
}

func Test_lifecycle_beforeCreate(t *testing.T) {
	adopt := func(change func(*ssot.Object)) ResourceModel {
		return lifecycleModel(withHooks(ssot.WithBeforeCreate(ssot.Generated())), func(obj *ssot.Object) {
			obj.IdentityNamesCreatedBySnowflake = []string{"PUBLIC"}
			change(obj)
		})
	}

	t.Run("omitted", func(t *testing.T) {
		require.Empty(t, renderTemplate(t, "beforeCreate", lifecycleModel()))
	})

	t.Run("ext call", func(t *testing.T) {
		require.Equal(t, `
	done, diags := poolBeforeCreateExt(ctx, d, meta, id)
	if done {
		return diags
	}

`, renderTemplate(t, "beforeCreate", lifecycleModel(withHooks(ssot.WithBeforeCreate(ssot.ExtCall())))))
	})

	t.Run("adopt an object that Snowflake creates", func(t *testing.T) {
		require.Equal(t, `
	if strings.EqualFold(strings.TrimSpace(d.Get("name").(string)), "PUBLIC") {
		_, err := meta.(*provider.Context).Client.Pools.ShowByID(ctx, id)
		if err == nil {
			log.Printf("[DEBUG] found %s pool during creation, updating...", id.Name())
			d.SetId(helpers.EncodeResourceIdentifier(id))
			return UpdatePool(ctx, d, meta)
		}
		if !errors.Is(err, sdk.ErrObjectNotFound) {
			return diag.FromErr(err)
		}
	}

`, renderTemplate(t, "beforeCreate", adopt(func(*ssot.Object) {})))
	})

	t.Run("any of the names", func(t *testing.T) {
		got := renderTemplate(t, "beforeCreate", adopt(func(obj *ssot.Object) { obj.IdentityNamesCreatedBySnowflake = []string{"PUBLIC", "DEFAULT"} }))
		require.Contains(t, got, `if strings.EqualFold(strings.TrimSpace(d.Get("name").(string)), "PUBLIC") || strings.EqualFold(strings.TrimSpace(d.Get("name").(string)), "DEFAULT") {`)
	})

	t.Run("legacy id", func(t *testing.T) {
		got := renderTemplate(t, "beforeCreate", adopt(func(obj *ssot.Object) { obj.LegacyIdEncoding = true }))
		require.Contains(t, got, "d.SetId(helpers.EncodeSnowflakeID(id))")
	})

	t.Run("the model holds the adopt data", func(t *testing.T) {
		require.Equal(t, BeforeCreateSteps{AdoptExisting: &BeforeCreateAdoptExisting{Names: []string{"PUBLIC"}, ObjectWords: "pool"}}, adopt(func(*ssot.Object) {}).BeforeCreateSteps)
		require.Equal(t, BeforeCreateSteps{}, lifecycleModel().BeforeCreateSteps)
	})
}

func Test_lifecycle_afterCreate(t *testing.T) {
	t.Run("omitted", func(t *testing.T) {
		require.Equal(t, "return ReadPoolFunc(false)(ctx, d, meta)\n", renderTemplate(t, "afterCreate", lifecycleModel()))
	})

	t.Run("ext call", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithAfterCreate(ssot.ExtCall())))
		require.Equal(t, `afterCreateDiags := poolAfterCreateExt(ctx, d, meta, id)
	return append(afterCreateDiags, ReadPoolFunc(false)(ctx, d, meta)...)
`, renderTemplate(t, "afterCreate", model))
	})
}

func Test_lifecycle_afterCreateSetAfterCreate(t *testing.T) {
	toggle := func(name string) ssot.Attribute {
		return ssot.Attribute{Name: name, Type: ssot.AttrBool, NativeBool: true, SkipCreate: true, SetAfterCreate: true}
	}
	enabled, managed := toggle("enabled"), toggle("managed")
	model := func(attributes []ssot.Attribute, clauses ...ssot.AlterClause) ResourceModel {
		return lifecycleModel(withHooks(ssot.WithAfterCreate(ssot.Generated())), func(obj *ssot.Object) {
			obj.Attributes = attributes
			obj.Alter = ssot.AlterSpec{Clauses: clauses}
		})
	}

	t.Run("calls the keyword linked to true and warns", func(t *testing.T) {
		got := renderTemplate(t, "afterCreate", model([]ssot.Attribute{enabled},
			ssot.AlterKeywordFor(&enabled, false, "SUSPEND"),
			ssot.AlterKeywordFor(&enabled, true, "RESUME"),
		))
		require.Equal(t, `var afterCreateDiags diag.Diagnostics
	if d.Get("enabled").(bool) {
		if err := meta.(*provider.Context).Client.Pools.Alter(ctx, sdk.NewAlterPoolRequest(id).WithResume(true)); err != nil {
			afterCreateDiags = append(afterCreateDiags, diag.Diagnostic{Severity: diag.Warning, Summary: fmt.Sprintf("RESUME after create: %s", err.Error())})
		}
	}
	return append(afterCreateDiags, ReadPoolFunc(false)(ctx, d, meta)...)
`, got)
	})

	t.Run("attributes run in attribute order", func(t *testing.T) {
		got := renderTemplate(t, "afterCreate", model([]ssot.Attribute{enabled, managed},
			ssot.AlterKeywordFor(&managed, true, "ENABLE MANAGED ACCESS"),
			ssot.AlterKeywordFor(&enabled, true, "RESUME"),
		))
		require.Less(t, strings.Index(got, "WithResume(true)"), strings.Index(got, "WithEnableManagedAccess(true)"))
	})

	t.Run("attributes without SetAfterCreate are left out", func(t *testing.T) {
		plain := enabled
		plain.SetAfterCreate = false
		require.Empty(t, model([]ssot.Attribute{plain}, ssot.AlterKeywordFor(&plain, true, "RESUME")).AfterCreateSteps.AlterKeywords)
	})

	t.Run("the model holds the keyword data", func(t *testing.T) {
		got := model([]ssot.Attribute{enabled}, ssot.AlterKeywordFor(&enabled, true, "ENABLE MANAGED ACCESS")).AfterCreateSteps
		require.Equal(t, AfterCreateSteps{AlterKeywords: []AfterCreateAlterKeyword{
			{Attribute: "enabled", Keyword: "ENABLE MANAGED ACCESS", Field: "EnableManagedAccess"},
		}}, got)
	})
}

func Test_lifecycle_afterAlter(t *testing.T) {
	t.Run("omitted", func(t *testing.T) {
		require.Equal(t, "return ReadPoolFunc(false)(ctx, d, meta)\n", renderTemplate(t, "afterAlter", lifecycleModel()))
	})

	t.Run("ext call", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithAfterAlter(ssot.ExtCall())))
		require.Equal(t, `afterAlterDiags := poolAfterAlterExt(ctx, d, meta, id)
	return append(afterAlterDiags, ReadPoolFunc(false)(ctx, d, meta)...)
`, renderTemplate(t, "afterAlter", model))
	})
}

func Test_lifecycle_afterShowById(t *testing.T) {
	t.Run("omitted", func(t *testing.T) {
		require.Empty(t, renderTemplate(t, "afterShowById", lifecycleModel()))
	})

	t.Run("ext call", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithAfterShowById(ssot.ExtCall())))
		require.Equal(t, `afterShowByIdDiags := poolAfterShowByIdExt(ctx, d, meta, id, pool)
		if afterShowByIdDiags.HasError() {
			return afterShowByIdDiags
		}
		diags = append(diags, afterShowByIdDiags...)
`, renderTemplate(t, "afterShowById", model))
	})
}

func Test_lifecycle_importPermissiveMode(t *testing.T) {
	got := renderLifecycle(t, lifecycleModel())
	require.Contains(t, got, "importPermissiveMode := meta.(*provider.Context).Experiments.IsEnabled(experimentalfeatures.ImportPermissiveMode)")
	require.Contains(t, got, "poolSetFieldsNotSetByReadExt(d, id, pool, importPermissiveMode)")

	omitted := renderLifecycle(t, lifecycleModel(withHooks(ssot.WithSetFieldsNotSetByRead(ssot.Omitted()))))
	require.NotContains(t, omitted, "importPermissiveMode")
	require.NotContains(t, omitted, "SetFieldsNotSetByReadExt")
}

func Test_lifecycle_importOmittedWithoutShowByIdInSdk(t *testing.T) {
	require.Contains(t, renderLifecycle(t, lifecycleModel()), "func ImportPool(")

	got := renderLifecycle(t, lifecycleModel(withHooks(ssot.WithShowByIdInSdk(ssot.Omitted()))))
	require.NotContains(t, got, "func ImportPool(")
}

func Test_lifecycle_mappingHooksOmitted(t *testing.T) {
	tests := []struct {
		name     string
		hook     ssot.ResourceHookOption
		function string
	}{
		{"parseIdFromConfig", ssot.WithParseIdFromConfig(ssot.Omitted()), "func poolParseIdFromConfig("},
		{"newCreateRequest", ssot.WithNewCreateRequest(ssot.Omitted()), "func poolNewCreateRequest("},
		{"applyCreateOptionals", ssot.WithApplyCreateOptionals(ssot.Omitted()), "func poolApplyCreateOptionals("},
		{"applySetUnset", ssot.WithApplySetUnset(ssot.Omitted()), "func poolApplySetUnset("},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := lifecycleModel(withHooks(tt.hook))
			require.Empty(t, renderTemplate(t, tt.name, model))
			require.NotContains(t, renderLifecycle(t, model), tt.function)
		})
	}
}

func Test_lifecycle_parseIdFromConfigExt(t *testing.T) {
	model := lifecycleModel()
	require.Equal(t, "id, err := poolParseIdFromConfigExt(d)\n", renderTemplate(t, "parseIdFromConfig", model))
	require.NotContains(t, renderLifecycle(t, model), "func poolParseIdFromConfig(")
}

func Test_lifecycle_dropSafelyInSdk(t *testing.T) {
	t.Run("generated", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithDropSafelyInSdk(ssot.Generated())))
		got := renderLifecycle(t, model)
		require.Contains(t, got, "DeleteContext: TrackingDeleteWrapper(resources.Pool, DeletePool),")
		require.NotContains(t, got, "poolDeleteExt")
		require.Equal(t, `func DeletePool(ctx context.Context, d *tfschema.ResourceData, meta any) diag.Diagnostics {
	return ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.Pools.DropSafely
		},
	)(ctx, d, meta)
}
`, renderTemplate(t, "generatedDelete", model))
		require.Contains(t, got, "func DeletePool(")
	})

	t.Run("legacy id", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithDropSafelyInSdk(ssot.Generated())), func(obj *ssot.Object) {
			obj.LegacyIdEncoding = true
		})
		require.Equal(t, `func DeletePool(ctx context.Context, d *tfschema.ResourceData, meta any) diag.Diagnostics {
	return ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErrLegacy[sdk.SchemaObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.Pools.DropSafely
		},
	)(ctx, d, meta)
}
`, renderTemplate(t, "generatedDelete", model))
	})

	t.Run("preview feature", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithDropSafelyInSdk(ssot.Generated())), func(obj *ssot.Object) {
			obj.PreviewFeature = true
		})
		got := renderLifecycle(t, model)
		require.Contains(t, got, "DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.PoolResource), TrackingDeleteWrapper(resources.Pool, DeletePool)),")
		require.NotContains(t, got, "poolDeleteExt")
	})

	for name, hook := range map[string]ssot.ResourceHook{
		"ext call": ssot.ExtCall(),
		"omitted":  ssot.Omitted(),
	} {
		t.Run(name, func(t *testing.T) {
			model := lifecycleModel(withHooks(ssot.WithDropSafelyInSdk(hook)))
			got := renderLifecycle(t, model)
			require.Contains(t, got, "DeleteContext: poolDeleteExt,")
			require.NotContains(t, got, "func DeletePool(")
		})
	}
}

func Test_lifecycle_outputSetHooks(t *testing.T) {
	tests := []struct {
		name      string
		template  string
		withHook  func(ssot.ResourceHook) ssot.ResourceHookOption
		generated string
		extCall   string
	}{
		{
			name:      "show output set",
			template:  "showOutputSet",
			withHook:  ssot.WithShowOutputSet,
			generated: "d.Set(ShowOutputAttributeName, []map[string]any{schemas.PoolToSchema(pool)}),\n",
			extCall:   "poolShowOutputSetExt(d, pool),\n",
		},
		{
			name:      "describe output set",
			template:  "describeOutputSet",
			withHook:  ssot.WithDescribeOutputSet,
			generated: "d.Set(DescribeOutputAttributeName, []map[string]any{schemas.PoolDetailsToSchema(poolDetails)}),\n",
			extCall:   "poolDescribeOutputSetExt(d, poolDetails),\n",
		},
		{
			name:      "fqn set",
			template:  "fqnSet",
			withHook:  ssot.WithFqnSet,
			generated: "d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),\n",
			extCall:   "poolFqnSetExt(d, id),\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name+" generated", func(t *testing.T) {
			model := lifecycleModel(withHooks(tt.withHook(ssot.Generated())))
			require.Equal(t, tt.generated, renderTemplate(t, tt.template, model))
		})

		t.Run(tt.name+" is generated by default", func(t *testing.T) {
			require.Equal(t, tt.generated, renderTemplate(t, tt.template, lifecycleModel()))
		})

		t.Run(tt.name+" ext call", func(t *testing.T) {
			model := lifecycleModel(withHooks(tt.withHook(ssot.ExtCall())))
			require.Equal(t, tt.extCall, renderTemplate(t, tt.template, model))
		})

		t.Run(tt.name+" omitted", func(t *testing.T) {
			model := lifecycleModel(withHooks(tt.withHook(ssot.Omitted())))
			require.Empty(t, renderTemplate(t, tt.template, model))
		})
	}
}

func Test_hierarchyRename(t *testing.T) {
	t.Run("both identity fields force new emits nothing", func(t *testing.T) {
		model := lifecycleModel(func(obj *ssot.Object) {
			obj.IdentityNameForceNew = true
			obj.IdentityHierarchyForceNew = true
			obj.ResourceHooks = []ssot.ResourceHookOption{ssot.WithHierarchyRename(ssot.Generated())}
		})
		require.Empty(t, renderTemplate(t, "hierarchyRename", model))
	})

	t.Run("hierarchy force new renames name", func(t *testing.T) {
		model := lifecycleModel(func(obj *ssot.Object) {
			obj.IdentityHierarchyForceNew = true
			obj.ResourceHooks = []ssot.ResourceHookOption{ssot.WithHierarchyRename(ssot.Generated())}
		})
		got := renderTemplate(t, "hierarchyRename", model)
		require.Contains(t, got, `if d.HasChange("name") {`)
		require.Contains(t, got, `sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))`)
		require.Contains(t, got, `Client.Pools.Alter(ctx, sdk.NewAlterPoolRequest(id).WithRenameTo(newId))`)
		require.Contains(t, got, `d.SetId(helpers.EncodeResourceIdentifier(newId))`)
		require.NotContains(t, got, `d.HasChange("database")`)
		require.NotContains(t, got, `d.HasChange("schema")`)
		require.NotContains(t, got, "handleThreeLevelHierarchyRename")
		require.NotContains(t, got, "HierarchyRenameExt")
	})

	t.Run("legacy id encoding", func(t *testing.T) {
		model := lifecycleModel(func(obj *ssot.Object) {
			obj.LegacyIdEncoding = true
			obj.ResourceHooks = []ssot.ResourceHookOption{ssot.WithHierarchyRename(ssot.Generated())}
		})
		got := renderTemplate(t, "hierarchyRename", model)
		require.Contains(t, got, "helpers.EncodeSnowflakeID(newId)")
		require.Contains(t, got, "helpers.EncodeSnowflakeID(renamed)")
		require.NotContains(t, got, "EncodeResourceIdentifier")
	})

	t.Run("neither identity field force new handles database schema and name", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithHierarchyRename(ssot.Generated())))
		got := renderTemplate(t, "hierarchyRename", model)
		require.Contains(t, got, `providerCtx.Experiments.IsEnabled(experimentalfeatures.HierarchyRenames) && (d.HasChange("database") || d.HasChange("schema"))`)
		require.Contains(t, got, "handleThreeLevelHierarchyRename(")
		require.Contains(t, got, "client.Pools.ShowByID,")
		require.Contains(t, got, `"pool",`)
		require.Contains(t, got, `if d.HasChange("name") {`)
		require.Contains(t, got, `meta.(*provider.Context).Client.Pools.Alter(ctx, sdk.NewAlterPoolRequest(id).WithRenameTo(newId))`)
		require.NotContains(t, got, "HierarchyRenameExt")
	})

	t.Run("account object renames name", func(t *testing.T) {
		model := lifecycleModel(func(obj *ssot.Object) {
			obj.IdType = "AccountObjectIdentifier"
			obj.ResourceHooks = []ssot.ResourceHookOption{ssot.WithHierarchyRename(ssot.Generated())}
		})
		got := renderTemplate(t, "hierarchyRename", model)
		require.Contains(t, got, `sdk.NewAccountObjectIdentifier(d.Get("name").(string))`)
		require.Contains(t, got, `Client.Pools.Alter(ctx, sdk.NewAlterPoolRequest(id).WithRenameTo(newId))`)
		require.NotContains(t, got, "handleTwoLevelHierarchyRename")
		require.NotContains(t, got, "handleThreeLevelHierarchyRename")
		require.NotContains(t, got, `d.HasChange("database")`)
	})

	t.Run("account object with force new name emits nothing", func(t *testing.T) {
		model := lifecycleModel(func(obj *ssot.Object) {
			obj.IdType = "AccountObjectIdentifier"
			obj.IdentityNameForceNew = true
			obj.ResourceHooks = []ssot.ResourceHookOption{ssot.WithHierarchyRename(ssot.Generated())}
		})
		require.Empty(t, renderTemplate(t, "hierarchyRename", model))
	})

	t.Run("generated database object is unimplemented", func(t *testing.T) {
		model := lifecycleModel(func(obj *ssot.Object) {
			obj.IdType = "DatabaseObjectIdentifier"
			obj.ResourceHooks = []ssot.ResourceHookOption{ssot.WithHierarchyRename(ssot.Generated())}
		})
		err := LifecycleTemplate.ExecuteTemplate(&strings.Builder{}, "hierarchyRename", model)
		require.ErrorContains(t, err, "generated implementation for HierarchyRename is not implemented")
	})
}

func guardedPoolModel(enabled bool, hooks ...ssot.ResourceHookOption) ResourceModel {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	flag := ssot.Attribute{Name: "enabled", Type: ssot.AttrBool, NativeBool: true}
	suspend := &ssot.SuspendForAlter{When: []*ssot.Attribute{&flag, &comment}}
	if enabled {
		suspend.Enabled = &flag
	}
	return lifecycleModel(withHooks(append([]ssot.ResourceHookOption{ssot.WithSuspendForAlter(ssot.Generated())}, hooks...)...), func(obj *ssot.Object) {
		obj.SuspendForAlter = suspend
	})
}

func Test_lifecycle_suspendForAlter(t *testing.T) {
	t.Run("the guard suspends once and cleans up on every path", func(t *testing.T) {
		require.Equal(t, `
	suspension, err := SuspendForAlter(ctx, d, meta.(*provider.Context).Client.Pools.SuspendIfStarted, id, "enabled", "comment")
	defer suspension.CleanUp(ctx, d, &diags)
	if err != nil {
		return diag.FromErr(err)
	}

`, renderTemplate(t, "suspendForAlter", guardedPoolModel(true)))
	})

	t.Run("the guard asks the resource when it decides whether the change needs the suspension", func(t *testing.T) {
		model := guardedPoolModel(false)
		model.SuspendForAlter.NeededExt = true
		require.Equal(t, `
	suspension, err := SuspendForAlterIf(ctx, d, poolSuspendForAlterNeededExt(d), meta.(*provider.Context).Client.Pools.SuspendIfStarted, id, "enabled", "comment")
	defer suspension.CleanUp(ctx, d, &diags)
	if err != nil {
		return diag.FromErr(err)
	}

`, renderTemplate(t, "suspendForAlter", model))
	})

	t.Run("the resume follows the enabled attribute", func(t *testing.T) {
		require.Equal(t, `
	diags = append(diags, suspension.RestoreIfEnabled(ctx, d.Get("enabled").(bool), func(ctx context.Context) error {
		return meta.(*provider.Context).Client.Pools.Alter(ctx, sdk.NewAlterPoolRequest(id).WithResume(true))
	})...)

`, renderTemplate(t, "resumeAfterAlter", guardedPoolModel(true)))
	})

	t.Run("the resume restores the previous state without an enabled attribute", func(t *testing.T) {
		require.Equal(t, `
	diags = append(diags, suspension.Restore(ctx)...)

`, renderTemplate(t, "resumeAfterAlter", guardedPoolModel(false)))
	})

	t.Run("after alter returns the guard diagnostics", func(t *testing.T) {
		require.Equal(t, "return append(diags, ReadPoolFunc(false)(ctx, d, meta)...)\n", renderTemplate(t, "afterAlter", guardedPoolModel(true)))
		require.Equal(t, `afterAlterDiags := poolAfterAlterExt(ctx, d, meta, id)
	return append(append(diags, afterAlterDiags...), ReadPoolFunc(false)(ctx, d, meta)...)
`, renderTemplate(t, "afterAlter", guardedPoolModel(true, ssot.WithAfterAlter(ssot.ExtCall()))))
	})

	t.Run("ext before and after alter hooks run inside the guard", func(t *testing.T) {
		got := renderLifecycle(t, guardedPoolModel(true, ssot.WithBeforeAlter(ssot.ExtCall()), ssot.WithAfterAlter(ssot.ExtCall())))
		suspend := strings.Index(got, "SuspendIfStarted(ctx, id)")
		before := strings.Index(got, "poolBeforeAlterExt(ctx, d, meta, id)")
		resume := strings.Index(got, "suspension.RestoreIfEnabled(ctx")
		after := strings.Index(got, "poolAfterAlterExt(ctx, d, meta, id)")
		require.True(t, suspend < before && before < resume && resume < after, "order: suspend, before alter, resume, after alter")
	})

	t.Run("update has a named result only with the guard", func(t *testing.T) {
		require.Contains(t, renderLifecycle(t, guardedPoolModel(true)), "func UpdatePool(ctx context.Context, d *tfschema.ResourceData, meta any) (diags diag.Diagnostics) {")
		without := lifecycleModel(func(obj *ssot.Object) {
			obj.Attributes = []ssot.Attribute{{Name: "comment", Type: ssot.AttrString}}
		})
		require.Contains(t, renderLifecycle(t, without), "func UpdatePool(ctx context.Context, d *tfschema.ResourceData, meta any) diag.Diagnostics {")
	})

	t.Run("a generated before alter hook is still unimplemented", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithBeforeAlter(ssot.Generated())))
		var b strings.Builder
		require.ErrorContains(t, LifecycleTemplate.ExecuteTemplate(&b, "beforeAlter", model), "BeforeAlter")
	})

	t.Run("the model holds the trigger and the enabled attribute", func(t *testing.T) {
		require.Equal(t, &SuspendForAlterModel{When: []string{"enabled", "comment"}, Enabled: "enabled"}, guardedPoolModel(true).SuspendForAlter)
		require.Equal(t, &SuspendForAlterModel{When: []string{"enabled", "comment"}}, guardedPoolModel(false).SuspendForAlter)
		require.Nil(t, lifecycleModel().SuspendForAlter)
	})

	t.Run("without when the model holds every attribute that changes in place", func(t *testing.T) {
		model := lifecycleModel(withHooks(ssot.WithSuspendForAlter(ssot.Generated())), func(obj *ssot.Object) {
			obj.Attributes = []ssot.Attribute{
				{Name: "comment", Type: ssot.AttrString},
				{Name: "kind", Type: ssot.AttrString, ForceNew: true},
				{Name: "size", Type: ssot.AttrInt},
			}
			obj.SuspendForAlter = &ssot.SuspendForAlter{}
		})
		require.Equal(t, &SuspendForAlterModel{When: []string{"comment", "size"}}, model.SuspendForAlter)
	})
}
