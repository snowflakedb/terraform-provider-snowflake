package resources

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var accountParameterSupportedParameters = buildAccountParameterSupportedParameters()

func buildAccountParameterSupportedParameters() []sdk.AccountParameter {
	catalogParams := defs.ParameterDefsForLevel(parameterdefs.ParameterLevelAccountExt)
	params := make([]sdk.AccountParameter, len(catalogParams))
	for i, p := range catalogParams {
		params[i] = sdk.AccountParameter(p.SqlName)
	}
	return params
}

func ToAccountParameter(s string) (sdk.AccountParameter, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(accountParameterSupportedParameters, sdk.AccountParameter(s)) {
		return "", fmt.Errorf("invalid account parameter: %s", s)
	}
	return sdk.AccountParameter(s), nil
}

var accountParameterSchema = map[string]*schema.Schema{
	"key": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		ValidateDiagFunc: sdkValidation(ToAccountParameter),
		DiffSuppressFunc: NormalizeAndCompare(ToAccountParameter),
		Description:      fmt.Sprintf("Name of account parameter. Valid values are (case-insensitive): %s. Deprecated parameters are not supported in the provider.", possibleValuesListed(sdk.AsStringList(accountParameterSupportedParameters))),
	},
	"value": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Value of account parameter, as a string. Constraints are the same as those for the parameters in Snowflake documentation. The parameter values are validated in Snowflake.",
	},
}

func AccountParameter() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.AccountParameter, CreateAccountParameter),
		ReadContext:   TrackingReadWrapper(resources.AccountParameter, ReadAccountParameter),
		UpdateContext: TrackingUpdateWrapper(resources.AccountParameter, UpdateAccountParameter),
		DeleteContext: TrackingDeleteWrapper(resources.AccountParameter, DeleteAccountParameter),

		Description: "Resource used to manage current account parameters. For more information, check [parameters documentation](https://docs.snowflake.com/en/sql-reference/parameters).",

		Schema: accountParameterSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

// CreateAccountParameter implements schema.CreateFunc.
func CreateAccountParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
	key := strings.ToUpper(d.Get("key").(string))
	value := d.Get("value").(string)

	if providerCtx.Experiments.IsEnabled(experimentalfeatures.AccountParameterCatalogWritePath) {
		req := sdk.NewAccountParametersRequest()
		if err := req.SetParameterFromRaw(key, value); err != nil {
			return diag.FromErr(err)
		}
		if err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().WithParameters(*req))); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := client.Parameters.SetAccountParameter(ctx, sdk.AccountParameter(key), value); err != nil {
			return diag.FromErr(fmt.Errorf("error creating account parameter err = %w", err))
		}
	}

	d.SetId(helpers.EncodeResourceIdentifier(key))
	return ReadAccountParameter(ctx, d, meta)
}

// ReadAccountParameter implements schema.ReadFunc.
func ReadAccountParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	parameterNameRaw := d.Id()
	parameterName, err := ToAccountParameter(parameterNameRaw)
	if err != nil {
		return diag.FromErr(err)
	}
	parameter, err := client.Parameters.ShowAccountParameter(ctx, parameterName)
	if err != nil {
		return diag.FromErr(fmt.Errorf("reading account parameter: %w", err))
	}
	errs := errors.Join(
		d.Set("value", parameter.Value),
		d.Set("key", parameter.Key),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	return nil
}

// UpdateAccountParameter implements schema.UpdateFunc.
func UpdateAccountParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return CreateAccountParameter(ctx, d, meta)
}

// DeleteAccountParameter implements schema.DeleteFunc.
func DeleteAccountParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
	key := strings.ToUpper(d.Get("key").(string))

	if providerCtx.Experiments.IsEnabled(experimentalfeatures.AccountParameterCatalogWritePath) {
		req := sdk.NewAccountParametersUnsetRequest()
		if err := req.UnsetParameterFromRaw(key); err != nil {
			return diag.FromErr(err)
		}
		if err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithUnset(*sdk.NewAccountUnsetRequest().WithParameters(*req))); err != nil {
			return diag.FromErr(fmt.Errorf("unsetting account parameter: %w", err))
		}
	} else {
		if err := client.Parameters.UnsetAccountParameter(ctx, sdk.AccountParameter(key)); err != nil {
			return diag.FromErr(fmt.Errorf("unsetting account parameter: %w", err))
		}
	}

	d.SetId("")
	return nil
}
