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
	attachUpdateMappings("AuthenticationPolicy", src, dst)
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
