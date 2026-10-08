package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var openflowDeploymentParametersSchema = make(map[string]*schema.Schema)

func init() {
	for _, p := range defs.ParameterDefsForLevel(parameterdefs.ParameterLevelOpenflowDeployment) {
		openflowDeploymentParametersSchema[p.FieldName()] = parameterSchema(p)
	}
}

// handleOpenflowDeploymentParametersCreate builds the EventTable request wrapper by hand because
// WithParameters' identifier case can only emit `EVENT_TABLE = <id>`, with no way to express the NONE
// alternative that openflowDeploymentEventTableDef() supports.
func handleOpenflowDeploymentParametersCreate(d *schema.ResourceData, request *sdk.CreateOpenflowDeploymentRequest) diag.Diagnostics {
	eventTable := sdk.NewOpenflowDeploymentEventTableRequest()
	if diags := handleParameterCreateWithMapping(d, defs.EventTable.FieldName(), &eventTable.EventTable, sdk.ParseSchemaObjectIdentifier); diags.HasError() {
		return diags
	}
	if eventTable.EventTable != nil {
		request.WithEventTable(*eventTable)
	}
	return nil
}

func handleOpenflowDeploymentParametersChanges(d *schema.ResourceData, set *sdk.OpenflowDeploymentSetRequest, unset *sdk.OpenflowDeploymentUnsetRequest) diag.Diagnostics {
	eventTable := sdk.NewOpenflowDeploymentEventTableRequest()
	if diags := handleParameterUpdateWithMapping(d, defs.EventTable.FieldName(), &eventTable.EventTable, &unset.EventTable, sdk.ParseSchemaObjectIdentifier); diags.HasError() {
		return diags
	}
	if eventTable.EventTable != nil {
		set.WithEventTable(*eventTable)
	}
	return nil
}

func handleOpenflowDeploymentParameterRead(d *schema.ResourceData, parameters *sdk.OpenflowDeploymentParametersDetails) diag.Diagnostics {
	return setResourceData(d, defs.EventTable.FieldName(), parameters.EventTable.Value.FullyQualifiedName())
}

var openflowDeploymentParametersCustomDiff = ParametersCustomDiffFromTypedParameters(
	openflowDeploymentParametersProvider,
	openflowDeploymentParameterDiffFunctions,
)

func openflowDeploymentParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) (*sdk.OpenflowDeploymentParametersDetails, error) {
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	return meta.(*provider.Context).Client.OpenflowDeployments.ShowParametersDetails(ctx, id)
}

func openflowDeploymentParameterDiffFunctions(parameters *sdk.OpenflowDeploymentParametersDetails) []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		IdentifierTypedParameterValueComputedIf(defs.EventTable.FieldName(), parameters.EventTable, sdk.ParameterTypeOpenflowDeployment),
	}
}
