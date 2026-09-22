package telemetry

import (
	"context"
	"runtime"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
)

func EmitProviderInit(ctx context.Context, providerCtx *internalprovider.Context) {
	if providerCtx == nil {
		return
	}
	emit(ctx, providerCtx.Client, typeProviderInit, providerCtx.SpanID, providerInitFields(providerCtx))
}

func providerInitFields(providerCtx *internalprovider.Context) map[string]string {
	authType := "UNKNOWN"
	if providerCtx.Client != nil {
		if cfg := providerCtx.Client.GetConfig(); cfg != nil {
			authType = cfg.Authenticator.String()
		}
	}
	return map[string]string{
		"experimental_features_enabled": collections.SortedJoinStrings(providerCtx.Experiments.UserEnabled(), ","),
		"preview_features_enabled":      collections.SortedJoinStrings(providerCtx.EnabledFeatures, ","),
		"os":                            runtime.GOOS,
		"arch":                          runtime.GOARCH,
		"ci_environment":                string(detectCIEnvironmentFromOS()),
		"agent_environment":             string(detectAgentEnvironmentFromOS()),
		"auth_type":                     authType,
		"terraform_host":                string(detectTerraformHostFromOS()),
	}
}
