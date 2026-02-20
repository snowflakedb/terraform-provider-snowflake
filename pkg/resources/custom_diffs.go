package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/sdkv2enhancements"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func StringParameterValueComputedIf[T ~string](key string, params []*sdk.Parameter, parameterLevel sdk.ParameterType, parameterName T) schema.CustomizeDiffFunc {
	return ParameterValueComputedIf(key, params, parameterLevel, parameterName, func(value any) string { return value.(string) })
}

func IntParameterValueComputedIf[T ~string](key string, params []*sdk.Parameter, parameterLevel sdk.ParameterType, parameterName T) schema.CustomizeDiffFunc {
	return ParameterValueComputedIf(key, params, parameterLevel, parameterName, func(value any) string { return strconv.Itoa(value.(int)) })
}

func BoolParameterValueComputedIf[T ~string](key string, params []*sdk.Parameter, parameterLevel sdk.ParameterType, parameterName T) schema.CustomizeDiffFunc {
	return ParameterValueComputedIf(key, params, parameterLevel, parameterName, func(value any) string { return strconv.FormatBool(value.(bool)) })
}

func ParameterValueComputedIf[T ~string](key string, parameters []*sdk.Parameter, objectParameterLevel sdk.ParameterType, parameterName T, valueToString func(v any) string) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
		providerCtx := meta.(*provider.Context)
		foundParameter, err := collections.FindFirst(parameters, func(parameter *sdk.Parameter) bool { return parameter.Key == string(parameterName) })
		if err != nil {
			log.Printf("[WARN] failed to find parameter: %s", parameterName)
			return nil
		}
		parameter := *foundParameter

		configValue, ok := d.GetRawConfig().AsValueMap()[key]

		// For cases where currently set value (in the config) is equal to the parameter, but not set on the right level.
		// The parameter is set somewhere higher in the hierarchy, and we need to "forcefully" set the value to
		// perform the actual set on Snowflake (and set the parameter on the correct level).
		if ok && !configValue.IsNull() && parameter.Level != objectParameterLevel && parameter.Value == valueToString(d.Get(key)) {
			return d.SetNewComputed(key)
		}

		// For all other cases, if a parameter is set in the configuration, we can ignore parts needed for Computed fields.
		if ok && !configValue.IsNull() {
			return nil
		}

		if experimentalfeatures.IsExperimentEnabled(experimentalfeatures.ParametersIgnoreValueChangesIfNotOnObjectLevel, providerCtx.EnabledExperiments) {
			// If the configuration is not set, perform SetNewComputed only for parameter being set on the object level (if so, it means that it was set externally, and we have to unset it).
			// The value change is handled through read.
			if parameter.Level == objectParameterLevel {
				return d.SetNewComputed(key)
			}
		} else {
			// If the configuration is not set, perform SetNewComputed for cases like:
			// 1. Check if the parameter value differs from the one saved in state (if they differ, we'll update the computed value).
			// 2. Check if the parameter is set on the object level (if so, it means that it was set externally, and we have to unset it).
			if parameter.Value != valueToString(d.Get(key)) || parameter.Level == objectParameterLevel {
				return d.SetNewComputed(key)
			}
		}

		return nil
	}
}

// ForceNewIfChangeToEmptySlice sets a ForceNew for a list field which was set to an empty value.
func ForceNewIfChangeToEmptySlice[T any](key string) schema.CustomizeDiffFunc {
	return customdiff.ForceNewIfChange(key, func(ctx context.Context, oldValue, newValue, meta any) bool {
		oldList, newList := oldValue.([]T), newValue.([]T)
		return len(oldList) > 0 && len(newList) == 0
	})
}

// ForceNewIfChangeToEmptySet sets a ForceNew for a list field which was set to an empty value.
func ForceNewIfChangeToEmptySet(key string) schema.CustomizeDiffFunc {
	return customdiff.ForceNewIfChange(key, func(ctx context.Context, oldValue, newValue, meta any) bool {
		oldList, newList := oldValue.(*schema.Set).List(), newValue.(*schema.Set).List()
		return len(oldList) > 0 && len(newList) == 0
	})
}

// ForceNewIfChangeToEmptyString sets a ForceNew for a string field which was set to an empty value.
func ForceNewIfChangeToEmptyString(key string) schema.CustomizeDiffFunc {
	return customdiff.ForceNewIfChange(key, func(ctx context.Context, oldValue, newValue, meta any) bool {
		oldString, newString := oldValue.(string), newValue.(string)
		return len(oldString) > 0 && len(newString) == 0
	})
}

// ComputedIfAnyAttributeChanged marks the given fields as computed if any of the listed fields changes.
// It takes field-level diffSuppress into consideration based on the schema passed.
// If the field is not found in the given schema, it continues without error. Only top level schema fields should be used.
func ComputedIfAnyAttributeChanged(resourceSchema map[string]*schema.Schema, key string, changedAttributeKeys ...string) schema.CustomizeDiffFunc {
	return customdiff.ComputedIf(key, func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
		var result bool
		for _, changedKey := range changedAttributeKeys {
			if diff.HasChange(changedKey) || !diff.NewValueKnown(changedKey) {
				oldValue, newValue := diff.GetChange(changedKey)
				log.Printf("[DEBUG] ComputedIfAnyAttributeChanged: changed key: %s", changedKey)

				if v, ok := resourceSchema[changedKey]; ok {
					if diffSuppressFunc := v.DiffSuppressFunc; diffSuppressFunc != nil {
						resourceData, resourceDataOk := sdkv2enhancements.CreateResourceDataFromResourceDiff(resourceSchema, diff)
						if !resourceDataOk {
							log.Printf("[DEBUG] ComputedIfAnyAttributeChanged: did not create resource data correctly, skipping")
							continue
						}
						if !diffSuppressFunc(key, fmt.Sprintf("%v", oldValue), fmt.Sprintf("%v", newValue), resourceData) {
							log.Printf("[DEBUG] ComputedIfAnyAttributeChanged: key %s was changed and the diff is not suppressed", changedKey)
							result = true
						} else {
							log.Printf("[DEBUG] ComputedIfAnyAttributeChanged: key %s was changed but the diff is suppressed", changedKey)
						}
					} else {
						log.Printf("[DEBUG] ComputedIfAnyAttributeChanged: key %s was changed and it does not have a diff suppressor", changedKey)
						result = true
					}
				}
			}
		}
		return result
	})
}

type parameter[T ~string] struct {
	parameterName T
	valueType     valueType
	parameterType sdk.ParameterType
}

type valueType string

const (
	valueTypeInt    valueType = "int"
	valueTypeBool   valueType = "bool"
	valueTypeString valueType = "string"
)

type ResourceIdProvider interface {
	Id() string
}

func ParametersCustomDiff[T ~string](parametersProvider func(context.Context, ResourceIdProvider, any) ([]*sdk.Parameter, error), parameters ...parameter[T]) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
		if d.Id() == "" {
			return nil
		}

		params, err := parametersProvider(ctx, d, meta)
		if err != nil {
			return err
		}

		diffFunctions := make([]schema.CustomizeDiffFunc, len(parameters))
		for idx, p := range parameters {
			var diffFunction schema.CustomizeDiffFunc
			switch p.valueType {
			case valueTypeInt:
				diffFunction = IntParameterValueComputedIf(strings.ToLower(string(p.parameterName)), params, p.parameterType, p.parameterName)
			case valueTypeBool:
				diffFunction = BoolParameterValueComputedIf(strings.ToLower(string(p.parameterName)), params, p.parameterType, p.parameterName)
			case valueTypeString:
				diffFunction = StringParameterValueComputedIf(strings.ToLower(string(p.parameterName)), params, p.parameterType, p.parameterName)
			}
			diffFunctions[idx] = diffFunction
		}

		return customdiff.All(diffFunctions...)(ctx, d, meta)
	}
}

func ForceNewIfAllKeysAreNotSet(key string, keys ...string) schema.CustomizeDiffFunc {
	return customdiff.ForceNewIf(key, func(ctx context.Context, d *schema.ResourceDiff, meta any) bool {
		allUnset := true
		for _, k := range keys {
			if _, ok := d.GetOk(k); ok {
				allUnset = false
			}
		}
		return allUnset
	})
}

func RecreateWhenUserTypeChangedExternally(userType sdk.UserType) schema.CustomizeDiffFunc {
	return func(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
		if n := diff.Get("user_type"); n != nil {
			log.Printf("[DEBUG] new external value for user type: %s", n.(string))
			if acceptableUserTypes, ok := sdk.AcceptableUserTypes[userType]; ok && !slices.Contains(acceptableUserTypes, strings.ToUpper(n.(string))) {
				// we have to set here a value instead of just SetNewComputed
				// because with empty value (default snowflake behavior for type) ForceNew fails
				// because there are no changes (at least from the SDKv2 point of view) for "user_type"
				return errors.Join(diff.SetNew("user_type", "<changed externally>"), diff.ForceNew("user_type"))
			}
		}
		return nil
	}
}

func RecreateWhenSecretTypeChangedExternally(secretType sdk.SecretType) schema.CustomizeDiffFunc {
	return func(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
		if n := diff.Get("secret_type"); n != nil {
			log.Printf("[DEBUG] new external value for secret type: %s", n.(string))

			diffSecretType, _ := sdk.ToSecretType(n.(string))
			if acceptableSecretTypes, ok := sdk.AcceptableSecretTypes[secretType]; ok && !slices.Contains(acceptableSecretTypes, diffSecretType) {
				return errors.Join(diff.SetNew("secret_type", "<changed externally>"), diff.ForceNew("secret_type"))
			}
			// both client_credentials and authorization_code_grant secrets have the same type: "OAUTH2"
			// to detect the external type change we need to check fields that are required in one, but should be absent in the other
			// we will check if the 'oauth_refresh_token_expiry_time' is present in the describe_output
			// since it is required in authorization_code_grant flow and should be empty in client_credentials flow
			if diffSecretType == sdk.SecretTypeOAuth2 {
				var isRefreshTokenExpiryTimeEmpty bool
				rt := diff.Get("describe_output.0.oauth_refresh_token_expiry_time").(string)

				switch secretType {
				case sdk.SecretTypeOAuth2AuthorizationCodeGrant:
					isRefreshTokenExpiryTimeEmpty = rt == ""
				case sdk.SecretTypeOAuth2ClientCredentials:
					isRefreshTokenExpiryTimeEmpty = rt != ""
				}
				if isRefreshTokenExpiryTimeEmpty {
					return errors.Join(diff.SetNew("secret_type", "<changed externally>"), diff.ForceNew("secret_type"))
				}
			}
		}
		return nil
	}
}

// RecreateWhenServiceTypeChangedExternally recreates a service when argument serviceType is different than in the state.
func RecreateWhenServiceTypeChangedExternally(serviceType sdk.ServiceType) schema.CustomizeDiffFunc {
	return RecreateWhenResourceTypeChangedExternally("service_type", serviceType, sdk.ToServiceType)
}

// RecreateWhenStreamTypeChangedExternally recreates a stream when argument streamType is different than in the state.
func RecreateWhenStreamTypeChangedExternally(streamType sdk.StreamSourceType) schema.CustomizeDiffFunc {
	return RecreateWhenResourceTypeChangedExternally("stream_type", streamType, sdk.ToStreamSourceType)
}

// RecreateWhenResourceTypeChangedExternally recreates a resource when argument wantType is different than the value in typeField.
func RecreateWhenResourceTypeChangedExternally[T ~string](typeField string, wantType T, toType func(string) (T, error)) schema.CustomizeDiffFunc {
	return func(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
		if n := diff.Get(typeField); n != nil {
			log.Printf("[DEBUG] new external value for %s", typeField)

			gotTypeRaw := n.(string)
			// if the type is empty, the state is empty - do not recreate
			if gotTypeRaw == "" {
				return nil
			}

			gotType, err := toType(gotTypeRaw)
			if err != nil {
				return fmt.Errorf("unknown type: %w", err)
			}
			if gotType != wantType {
				// we have to set here a value instead of just SetNewComputed
				// because with empty value (default snowflake behavior for type) ForceNew fails
				// because there are no changes (at least from the SDKv2 point of view) for typeField
				return errors.Join(diff.SetNew(typeField, "<changed externally>"), diff.ForceNew(typeField))
			}
		}
		return nil
	}
}

// RecreateWhenStreamIsStale detects when the stream is stale, and sets a `false` value for `stale` field.
// This means that the provider can detect that change in `stale` from `true` to `false`, where `false` is our desired state.
func RecreateWhenStreamIsStale() schema.CustomizeDiffFunc {
	return func(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
		if old, _ := diff.GetChange("stale"); old.(bool) {
			return diff.SetNew("stale", false)
		}
		return nil
	}
}

// RecreateWhenResourceBoolFieldChangedExternally recreates a resource when wantValue is different than value in boolField.
func RecreateWhenResourceBoolFieldChangedExternally(boolField string, wantValue bool) schema.CustomizeDiffFunc {
	return func(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
		if n := diff.Get(boolField); n != nil {
			log.Printf("[DEBUG] new external value for %v", boolField)
			if n.(bool) != wantValue {
				log.Printf("[DEBUG] recreating the resource, because %v doesn't equal to the expected value...", boolField)
				return errors.Join(diff.SetNew(boolField, wantValue), diff.ForceNew(boolField))
			}
		}
		return nil
	}
}

// RecreateWhenResourceStringFieldChangedExternally recreates a resource when wantValue is different from value in field.
// TODO [SNOW-1850370]: merge with above? test.
func RecreateWhenResourceStringFieldChangedExternally(field string, expectedValue string) schema.CustomizeDiffFunc {
	return func(_ context.Context, diff *schema.ResourceDiff, _ any) error {
		if o, n := diff.GetChange(field); n != nil && o != nil && o != "" && n.(string) != expectedValue {
			log.Printf("[DEBUG] new unexpected external value for: %s, was expecting: %s, recreating the resource...", field, expectedValue)
			return errors.Join(diff.SetNew(field, expectedValue), diff.ForceNew(field))
		}
		return nil
	}
}
