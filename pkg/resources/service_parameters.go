package resources

import (
	"context"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	serviceParametersSchema     = make(map[string]*schema.Schema)
	serviceParametersCustomDiff = ParametersCustomDiff(
		serviceParametersProvider,
		parameter[sdk.ServiceParameter]{sdk.ServiceParameterServiceCallerTokenValiditySecs, valueTypeInt, sdk.ParameterTypeService},
	)
)

func init() {
	serviceParameterFields := []parameterDef[sdk.ServiceParameter]{
		{Name: sdk.ServiceParameterServiceCallerTokenValiditySecs, Type: schema.TypeInt, Description: "Controls how long a caller's rights login token is valid for Snowpark Container Services."},
	}

	for _, field := range serviceParameterFields {
		fieldName := strings.ToLower(string(field.Name))

		serviceParametersSchema[fieldName] = &schema.Schema{
			Type:             field.Type,
			Description:      enrichWithReferenceToParameterDocs(field.Name, field.Description),
			Computed:         true,
			Optional:         true,
			ValidateDiagFunc: field.ValidateDiag,
			DiffSuppressFunc: field.DiffSuppress,
		}
	}
}

func serviceParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) ([]*sdk.Parameter, error) {
	return parametersProvider(ctx, d, meta.(*provider.Context), serviceParametersProviderFunc, sdk.ParseSchemaObjectIdentifier)
}

func serviceParametersProviderFunc(c *sdk.Client) showParametersFunc[sdk.SchemaObjectIdentifier] {
	return c.Services.ShowParameters
}

func handleServiceParameterRead(d *schema.ResourceData, serviceParameters []*sdk.Parameter) error {
	for _, p := range serviceParameters {
		if p.Key == string(sdk.ServiceParameterServiceCallerTokenValiditySecs) {
			value, err := strconv.Atoi(p.Value)
			if err != nil {
				return err
			}
			if err := d.Set(strings.ToLower(p.Key), value); err != nil {
				return err
			}
		}
	}

	return nil
}

func handleServiceParametersUpdate(d *schema.ResourceData, set *sdk.ServiceSetRequest, unset *sdk.ServiceUnsetRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, sdk.ServiceParameterServiceCallerTokenValiditySecs, &set.ServiceCallerTokenValiditySecs, &unset.ServiceCallerTokenValiditySecs),
	)
}
