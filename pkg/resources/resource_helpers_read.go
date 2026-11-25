package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// showByIdErrorDiagnostics turns the error of a failed ShowByIDSafely in Read into diagnostics.
// An object that no longer exists is removed from the state with a warning, and any other error is returned as it is.
func showByIdErrorDiagnostics(d *schema.ResourceData, objectType sdk.ObjectType, id sdk.ObjectIdentifier, err error) diag.Diagnostics {
	if !errors.Is(err, sdk.ErrObjectNotFound) {
		return diag.FromErr(err)
	}
	d.SetId("")
	objectDescription := strings.ToLower(string(objectType))
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("Failed to query %s. Marking the resource as removed.", objectDescription),
			Detail:   fmt.Sprintf("Resource %s id: %s, Err: %s", objectDescription, id.FullyQualifiedName(), err),
		},
	}
}

func setStringIfNotEmpty(d *schema.ResourceData, key, value string) error {
	if value != "" {
		return d.Set(key, value)
	}
	return nil
}

// setListKeepingNull skips d.Set when Snowflake and state hold no items, because d.Set stores an empty list instead of null.
func setListKeepingNull[T any](d *schema.ResourceData, key string, value []T) error {
	if len(value) == 0 && len(d.Get(key).([]any)) == 0 {
		return nil
	}
	return d.Set(key, value)
}

func attributeMappedValueReadOrDefault[T, R any](d *schema.ResourceData, key string, value *T, mapper func(*T) (R, error), defaultValue *R) error {
	if value != nil {
		mappedValue, err := mapper(value)
		if err != nil {
			return err
		}
		return d.Set(key, mappedValue)
	}
	if defaultValue != nil {
		return d.Set(key, *defaultValue)
	}
	return d.Set(key, nil)
}

func setOptionalValueWithMapping[T, R any](d *schema.ResourceData, key string, value *T, mapper func(*T) R) error {
	if value != nil {
		return d.Set(key, mapper(value))
	}
	return d.Set(key, nil)
}

func setOptionalFromStringPtr(d *schema.ResourceData, key string, ptr *string) error {
	if ptr != nil {
		if err := d.Set(key, *ptr); err != nil {
			return err
		}
	}
	return nil
}

func setOptionalFromPtr[T any](d *schema.ResourceData, key string, ptr *T) error {
	if ptr != nil {
		if err := d.Set(key, *ptr); err != nil {
			return err
		}
	}
	return nil
}

// TODO [SNOW-1348103]: return error if nil
func setRequiredFromStringPtr(d *schema.ResourceData, key string, ptr *string) error {
	if ptr != nil {
		if err := d.Set(key, *ptr); err != nil {
			return err
		}
	}
	return nil
}

func optionalStringOutputMapping[T ~string](value *T) (any, string) {
	if value != nil {
		return *value, string(*value)
	}
	return nil, ""
}

func optionalBooleanStringOutputMapping(value *bool) (any, string) {
	if value != nil {
		return *value, booleanStringFromBool(*value)
	}
	return nil, BooleanDefault
}

func optionalIntOutputMapping[T ~int](value *T) any {
	if value != nil {
		return int(*value)
	}
	return 0
}

func optionalIntOutputMappingIntDefault[T ~int](value *T) any {
	if value != nil {
		return int(*value)
	}
	return IntDefault
}

func enumStrings[T ~string](values []T) []string {
	out := make([]string, len(values))
	for i, value := range values {
		out[i] = string(value)
	}
	return out
}

func currentAccountIdentifier(ctx context.Context, meta any) (sdk.AccountIdentifier, error) {
	details, err := meta.(*provider.Context).Client.ContextFunctions.CurrentSessionDetails(ctx)
	if err != nil {
		return sdk.AccountIdentifier{}, err
	}
	return sdk.NewAccountIdentifier(details.OrganizationName, details.AccountName), nil
}

func allowedAccountsWithoutCurrent(ctx context.Context, meta any, accounts []sdk.AccountIdentifier) ([]sdk.AccountIdentifier, error) {
	current, err := currentAccountIdentifier(ctx, meta)
	if err != nil {
		return nil, err
	}
	return omitCurrentAccount(accounts, current), nil
}

func omitCurrentAccount(accounts []sdk.AccountIdentifier, current sdk.AccountIdentifier) []sdk.AccountIdentifier {
	filtered := make([]sdk.AccountIdentifier, 0, len(accounts))
	for _, account := range accounts {
		if account.FullyQualifiedName() == current.FullyQualifiedName() {
			continue
		}
		filtered = append(filtered, account)
	}
	return filtered
}

// SHOW returns NULL when no error integration is set, and on a secondary group. An attribute that is not in the configuration stays null.
func setErrorIntegrationFromShow(d *schema.ResourceData, id *sdk.AccountObjectIdentifier) error {
	if id != nil {
		return d.Set("error_integration", id.FullyQualifiedName())
	}
	if d.Get("error_integration").(string) != "" {
		return d.Set("error_integration", "")
	}
	return nil
}
