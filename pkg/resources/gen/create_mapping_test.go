package gen

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

func Test_ComputePoolCreateMappings_renderCreateFuncs(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.ComputePool}, nil, nil)
	require.True(t, model.Hooks.NewCreateRequest.Generate)
	require.True(t, model.Hooks.ApplyCreateOptionals.Generate)
	require.False(t, model.Hooks.ApplySetUnset.Generate)
	require.Equal(t, "for_application", model.Attributes[0].Name)
	require.Equal(t, "accountObjectIdentifierAttributeCreateBuilder", model.Attributes[0].CreateOptionalCall.Helper)

	require.Equal(t, "req, err := computePoolNewCreateRequest(d, id)\n", renderTemplate(t, "newCreateRequest", model))
	require.Equal(t, `if err := computePoolApplyCreateOptionals(d, req); err != nil {
		return diag.FromErr(err)
	}
`, renderTemplate(t, "applyCreateOptionals", model))

	require.Equal(t, `func computePoolNewCreateRequest(d *tfschema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateComputePoolRequest, error) {
	minNodes := d.Get("min_nodes").(int)
	maxNodes := d.Get("max_nodes").(int)
	instanceFamily, err := sdk.ToComputePoolInstanceFamily(d.Get("instance_family").(string))
	if err != nil {
		return nil, err
	}
	return sdk.NewCreateComputePoolRequest(id, minNodes, maxNodes, instanceFamily), nil
}
`, renderTemplate(t, "generatedNewCreateRequest", model))

	require.Equal(t, `func computePoolApplyCreateOptionals(d *tfschema.ResourceData, request *sdk.CreateComputePoolRequest) error {
	return errors.Join(
		accountObjectIdentifierAttributeCreateBuilder(d, "for_application", request.WithForApplication),
		listValueCreate(d, "backup_instance_families", &request.BackupInstanceFamilies, toComputePoolBackupInstanceFamily),
		booleanStringAttributeCreateBuilder(d, "auto_resume", request.WithAutoResume),
		booleanStringAttributeCreateBuilder(d, "initially_suspended", request.WithInitiallySuspended),
		intAttributeWithSpecialDefaultCreateBuilder(d, "auto_suspend_secs", request.WithAutoSuspendSecs),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
}
`, renderTemplate(t, "generatedApplyCreateOptionals", model))
}

func Test_SchemaCreateMappings_stayExt(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.Schema}, nil, nil)
	require.False(t, model.Hooks.NewCreateRequest.Generate)
	require.False(t, model.Hooks.ApplyCreateOptionals.Generate)
	require.Equal(t, "req, err := schemaNewCreateRequestExt(d, id)\n", renderTemplate(t, "newCreateRequest", model))
	require.Equal(t, `if err := schemaApplyCreateOptionalsExt(d, req); err != nil {
		return diag.FromErr(err)
	}
`, renderTemplate(t, "applyCreateOptionals", model))
	for _, attr := range model.Attributes {
		require.Nil(t, attr.CreateConstructorArg)
		require.Nil(t, attr.CreateOptionalCall)
	}
}

func renderTemplate(t *testing.T, name string, model ResourceModel) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, LifecycleTemplate.ExecuteTemplate(&b, name, model))
	return b.String()
}
