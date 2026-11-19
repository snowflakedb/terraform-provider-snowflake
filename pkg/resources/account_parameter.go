package resources

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var accountParameterSupportedParameters = buildAccountParameterSupportedParameters()

func buildAccountParameterSupportedParameters() []sdk.AccountParameter {
	catalogParams := defs.ParameterDefsForLevel(ssot.ParameterLevelAccountExt)
	params := make([]sdk.AccountParameter, len(catalogParams))
	for i, p := range catalogParams {
		params[i] = sdk.AccountParameter(p.SqlName)
	}
	return params
}

// suppressValueTypeForSupportedAccountParameter hides value_type diffs once the provider supports the parameter.
func suppressValueTypeForSupportedAccountParameter(_, _, _ string, d *schema.ResourceData) bool {
	return slices.Contains(accountParameterSupportedParameters, sdk.AccountParameter(strings.ToUpper(d.Get("key").(string))))
}

var accountParameterSchema = map[string]*schema.Schema{
	"key": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: ignoreCaseSuppressFunc,
		Description:      fmt.Sprintf("Name of account parameter (case-insensitive). The provider supports the following parameters, and it parses and quotes their values itself: %s. Any other account parameter is sent to Snowflake as is, and requires `value_type`. Deprecated parameters are not supported in the provider.", possibleValuesListed(sdk.AsStringList(accountParameterSupportedParameters))),
	},
	"value_type": parameterValueTypeSchema("the `key` field description", experimentalfeatures.AccountParameterCatalogWritePath, suppressValueTypeForSupportedAccountParameter),
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

// withExperimentHint tells how to send a parameter that the legacy write path does not support.
func withExperimentHint(err error) error {
	if errors.Is(err, sdk.ErrParameterNotSupported) || errors.Is(err, sdk.ErrLegacyParameterNotSupported) {
		return fmt.Errorf("%w; to send it, remove %s from `experimental_features_disabled` in the provider configuration (and set `value_type` if the provider does not support the parameter)", err, experimentalfeatures.AccountParameterCatalogWritePath)
	}
	return err
}

// CreateAccountParameter implements schema.CreateFunc.
func CreateAccountParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
	key := strings.ToUpper(d.Get("key").(string))
	value := d.Get("value").(string)

	// Empty when value_type is not set.
	var kind sdk.ParameterValueKind
	if valueType := d.Get("value_type").(string); valueType != "" {
		var err error
		if kind, err = sdk.ToParameterValueKind(valueType); err != nil {
			return diag.FromErr(err)
		}
	}

	switch {
	case !providerCtx.Experiments.IsEnabled(experimentalfeatures.AccountParameterCatalogWritePath):
		if err := client.Parameters.SetAccountParameterLegacy(ctx, sdk.AccountParameter(key), value); err != nil {
			return diag.FromErr(fmt.Errorf("error creating account parameter err = %w", withExperimentHint(err)))
		}
	case kind != "":
		if err := client.Parameters.SetAccountParameterWithGenericFallback(ctx, key, kind, value); err != nil {
			return diag.FromErr(fmt.Errorf("error creating account parameter err = %w", err))
		}
	default:
		handled, err := client.Parameters.SetObjectParameterOnAccountViaCatalog(ctx, sdk.ObjectParameter(key), value)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error creating account parameter err = %w", err))
		}
		if !handled {
			return diag.FromErr(fmt.Errorf("error creating account parameter err = %w: %s; to send it anyway, set the `value_type` field (one of: %s)", sdk.ErrParameterNotSupported, key, strings.Join(sdk.AsStringList(sdk.AllParameterValueKinds), ", ")))
		}
	}

	d.SetId(helpers.EncodeResourceIdentifier(key))
	return ReadAccountParameter(ctx, d, meta)
}

// ReadAccountParameter implements schema.ReadFunc.
func ReadAccountParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	parameter, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameter(strings.ToUpper(d.Id())))
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
		if err := client.Parameters.UnsetAccountParameterWithGenericFallback(ctx, key); err != nil {
			return diag.FromErr(fmt.Errorf("unsetting account parameter: %w", err))
		}
	} else if err := client.Parameters.UnsetAccountParameterLegacy(ctx, sdk.AccountParameter(key)); err != nil {
		return diag.FromErr(fmt.Errorf("unsetting account parameter: %w", withExperimentHint(err)))
	}

	d.SetId("")
	return nil
}
