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

var serviceParametersSchema = make(map[string]*schema.Schema)

func init() {
	for _, p := range defs.ParameterDefsForLevel(parameterdefs.ParameterLevelService) {
		serviceParametersSchema[p.FieldName()] = parameterSchema(p)
	}
}

func handleServiceParametersUpdate(d *schema.ResourceData, set *sdk.ServiceSetRequest, unset *sdk.ServiceUnsetRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, defs.ServiceCallerTokenValiditySecs.FieldName(), &set.ServiceCallerTokenValiditySecs, &unset.ServiceCallerTokenValiditySecs),
	)
}

func handleServiceParameterRead(d *schema.ResourceData, parameters *sdk.ServiceParametersDetails) diag.Diagnostics {
	return setResourceData(d, defs.ServiceCallerTokenValiditySecs.FieldName(), parameters.ServiceCallerTokenValiditySecs.Value)
}

var serviceParametersCustomDiff = ParametersCustomDiffFromTypedParameters(
	serviceParametersProvider,
	serviceParameterDiffFunctions,
)

func serviceParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) (*sdk.ServiceParametersDetails, error) {
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	return meta.(*provider.Context).Client.Services.ShowParametersDetails(ctx, id)
}

func serviceParameterDiffFunctions(parameters *sdk.ServiceParametersDetails) []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		IntTypedParameterValueComputedIf(defs.ServiceCallerTokenValiditySecs.FieldName(), parameters.ServiceCallerTokenValiditySecs, sdk.ParameterTypeService),
	}
}
