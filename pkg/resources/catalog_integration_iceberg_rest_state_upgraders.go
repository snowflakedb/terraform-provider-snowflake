package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// v2_24_0_CatalogIntegrationIcebergRestStateUpgrader fills describe_output.0.sigv4_rest_authentication.0.sigv4_external_id,
// which older versions did not store.
func v2_24_0_CatalogIntegrationIcebergRestStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}
	sigV4DescribeOutput, ok := sigV4RestAuthenticationInDescribeOutput(rawState)
	if !ok {
		return rawState, nil
	}

	idRaw, ok := rawState["id"].(string)
	if !ok {
		return nil, fmt.Errorf("cannot read Iceberg REST catalog integration id from the state")
	}
	id, err := sdk.ParseAccountObjectIdentifier(idRaw)
	if err != nil {
		return nil, err
	}

	client := meta.(*provider.Context).Client
	details, err := client.CatalogIntegrations.DescribeIcebergRestDetails(ctx, id)
	if err != nil {
		// Read removes a missing integration from the state.
		if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) || errors.Is(err, sdk.ErrObjectNotFound) {
			return rawState, nil
		}
		return nil, err
	}
	if details.Sigv4RestAuthentication != nil {
		sigV4DescribeOutput["sigv4_external_id"] = details.Sigv4RestAuthentication.Sigv4ExternalId
	}
	return rawState, nil
}

func sigV4RestAuthenticationInDescribeOutput(rawState map[string]any) (map[string]any, bool) {
	describeOutput, ok := rawState[DescribeOutputAttributeName].([]any)
	if !ok || len(describeOutput) != 1 {
		return nil, false
	}
	describeOutputMap, ok := describeOutput[0].(map[string]any)
	if !ok {
		return nil, false
	}
	sigV4, ok := describeOutputMap["sigv4_rest_authentication"].([]any)
	if !ok || len(sigV4) != 1 {
		return nil, false
	}
	sigV4Map, ok := sigV4[0].(map[string]any)
	return sigV4Map, ok
}
