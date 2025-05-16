package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// handleParameterCreate calls internally handleParameterCreateWithMapping with identity mapping
func handleParameterCreate[T any, P ~string](d *schema.ResourceData, parameterName P, createField **T) diag.Diagnostics {
	return handleParameterCreateWithMapping[T, T](d, parameterName, createField, identityMapping[T])
}

// handleParameterCreateWithMapping gets the property pointer from raw config.
// If the value is set, createField is set to the new planned value applying mapping beforehand.
// Otherwise, there is no change to the createField and nil is returned.
func handleParameterCreateWithMapping[T, R any, P ~string](d *schema.ResourceData, parameterName P, createField **R, mapping func(value T) (R, error)) diag.Diagnostics {
	key := strings.ToLower(string(parameterName))
	if v := GetConfigPropertyAsPointerAllowingZeroValue[T](d, key); v != nil {
		mappedValue, err := mapping(*v)
		if err != nil {
			return diag.FromErr(err)
		}
		*createField = sdk.Pointer(mappedValue)
	}
	return nil
}

// handleParameterUpdate calls internally handleParameterUpdateWithMapping with identity mapping
func handleParameterUpdate[T any, P ~string](d *schema.ResourceData, parameterName P, setField **T, unsetField **bool) diag.Diagnostics {
	return handleParameterUpdateWithMapping[T, T](d, parameterName, setField, unsetField, identityMapping[T])
}

// handleParameterUpdateWithMapping checks schema.ResourceData for change in key's value. If there's a change detected
// (or unknown value that basically indicates diff.SetNewComputed was called on the key), it checks if the value is set in the configuration.
// If the value is set, setField (representing setter for a value) is set to the new planned value applying mapping beforehand in cases where enum values,
// identifiers, etc. have to be set. Otherwise, unsetField is populated.
func handleParameterUpdateWithMapping[T, R any, P ~string](d *schema.ResourceData, parameterName P, setField **R, unsetField **bool, mapping func(value T) (R, error)) diag.Diagnostics {
	key := strings.ToLower(string(parameterName))
	if d.HasChange(key) || !d.GetRawPlan().AsValueMap()[key].IsKnown() {
		if !d.GetRawConfig().AsValueMap()[key].IsNull() {
			mappedValue, err := mapping(d.Get(key).(T))
			if err != nil {
				return diag.FromErr(err)
			}
			*setField = sdk.Pointer(mappedValue)
		} else {
			*unsetField = sdk.Bool(true)
		}
	}
	return nil
}

func identityMapping[T any](value T) (T, error) {
	return value, nil
}

func stringToAccountObjectIdentifier(value string) (sdk.AccountObjectIdentifier, error) {
	return sdk.NewAccountObjectIdentifier(value), nil
}

func stringToStringEnumProvider[T ~string](mapper func(string) (T, error)) func(value string) (T, error) {
	return func(value string) (T, error) {
		return mapper(value)
	}
}

// enrichWithReferenceToParameterDocs adds the reference to the Snowflake Parameter docs at the end of the provided description.
func enrichWithReferenceToParameterDocs[T ~string](parameter T, description string) string {
	link := fmt.Sprintf("https://docs.snowflake.com/en/sql-reference/parameters#%s", strings.ReplaceAll(strings.ToLower(string(parameter)), "_", "-"))
	return fmt.Sprintf("%s For more information, check [%s docs](%s).", description, parameter, link)
}

type (
	showParametersFunc[T sdk.ObjectIdentifier]         func(ctx context.Context, id T) ([]*sdk.Parameter, error)
	showParametersFuncProvider[T sdk.ObjectIdentifier] func(client *sdk.Client) showParametersFunc[T]
)

// parametersProvider is a generic function that can be used with ParametersCustomDiff
func parametersProvider[T sdk.ObjectIdentifier](
	ctx context.Context,
	d ResourceIdProvider,
	meta *provider.Context,
	showParametersFuncProvider showParametersFuncProvider[T],
	idParser func(string) (T, error),
) ([]*sdk.Parameter, error) {
	client := meta.Client
	id, err := idParser(d.Id())
	if err != nil {
		return nil, err
	}
	return showParametersFuncProvider(client)(ctx, id)
}
