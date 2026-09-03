package resources

import (
	"context"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func openflowDeploymentParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) ([]*sdk.Parameter, error) {
	return parametersProvider(ctx, d, meta.(*provider.Context), openflowDeploymentParametersProviderFunc, sdk.ParseAccountObjectIdentifier)
}

func openflowDeploymentParametersProviderFunc(c *sdk.Client) showParametersFunc[sdk.AccountObjectIdentifier] {
	return c.OpenflowDeployments.ShowParameters
}

var openflowDeploymentParametersCustomDiff = ParametersCustomDiff(
	openflowDeploymentParametersProvider,
	parameter[sdk.OpenflowDeploymentParameter]{sdk.OpenflowDeploymentParameterEventTable, valueTypeString, sdk.ParameterTypeOpenflowDeployment},
)

func handleOpenflowDeploymentParameterCreate(d *schema.ResourceData, request *sdk.CreateOpenflowDeploymentRequest) diag.Diagnostics {
	eventTable := sdk.NewOpenflowDeploymentEventTableRequest()
	if diags := handleParameterCreateWithMapping(d, sdk.OpenflowDeploymentParameterEventTable, &eventTable.EventTable, sdk.ParseSchemaObjectIdentifier); diags.HasError() {
		return diags
	}
	if eventTable.EventTable != nil {
		request.WithEventTable(*eventTable)
	}
	return nil
}

func handleOpenflowDeploymentParameterUpdate(d *schema.ResourceData, set *sdk.OpenflowDeploymentSetRequest, unset *sdk.OpenflowDeploymentUnsetRequest) diag.Diagnostics {
	eventTable := sdk.NewOpenflowDeploymentEventTableRequest()
	if diags := handleParameterUpdateWithMapping(d, sdk.OpenflowDeploymentParameterEventTable, &eventTable.EventTable, &unset.EventTable, sdk.ParseSchemaObjectIdentifier); diags.HasError() {
		return diags
	}
	if eventTable.EventTable != nil {
		set.WithEventTable(*eventTable)
	}
	return nil
}

func handleOpenflowDeploymentParameterRead(d *schema.ResourceData, parameters []*sdk.Parameter) error {
	for _, p := range parameters {
		if strings.ToUpper(p.Key) == string(sdk.OpenflowDeploymentParameterEventTable) {
			if err := d.Set(strings.ToLower(p.Key), p.Value); err != nil {
				return err
			}
		}
	}
	return nil
}
