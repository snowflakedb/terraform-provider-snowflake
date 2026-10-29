package gen

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

func Test_AuthenticationPolicyRead_preservesPropertyReadWarnings(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.AuthenticationPolicy}, nil, nil)
	require.True(t, model.Hooks.HandleExternalChangesInDescribe.Present)
	require.True(t, model.Hooks.DescribeInSdk.Present)

	hook := renderTemplate(t, "handleExternalChangesInDescribe", model)
	require.Contains(t, hook, "diags = authenticationPolicyHandleExternalChangesInDescribeExt(d, authenticationPolicyDetails, withExternalChangesMarking)")
	require.Contains(t, hook, "if diags.HasError() {")
	require.Contains(t, hook, "return diags")
	require.NotContains(t, hook, "diag.FromErr")

	setState := renderTemplate(t, "setStateToValuesFromConfig", model)
	require.Contains(t, setState, "return append(diags, diag.FromErr(err)...)")

	read := renderTemplate(t, "lifecycleTemplate", model)
	require.Contains(t, read, "var diags diag.Diagnostics")
	require.Contains(t, read, "return append(diags, diag.FromErr(errs)...)")
	require.Contains(t, read, "return diags")
}

func Test_ComputePoolRead_doesNotThreadDescribeWarnings(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.ComputePool}, nil, nil)
	require.False(t, model.Hooks.HandleExternalChangesInDescribe.Present)

	require.Empty(t, strings.TrimSpace(renderTemplate(t, "handleExternalChangesInDescribe", model)))

	setState := renderTemplate(t, "setStateToValuesFromConfig", model)
	require.Contains(t, setState, "return diag.FromErr(err)")
	require.NotContains(t, setState, "append(diags")

	read := renderTemplate(t, "lifecycleTemplate", model)
	require.Contains(t, read, "return diag.FromErr(errs)")
	require.NotContains(t, read, "var diags diag.Diagnostics")
	require.NotContains(t, read, "append(diags")
	require.NotContains(t, read, "HandleExternalChangesInDescribeExt")
}
