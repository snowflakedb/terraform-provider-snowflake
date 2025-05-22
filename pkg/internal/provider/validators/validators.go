package validators

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func NormalizeValidation[T any](normalize func(string) (T, error)) schema.SchemaValidateDiagFunc {
	return func(val interface{}, _ cty.Path) diag.Diagnostics {
		_, err := normalize(val.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
}

// IsValidIdentifier is a validator that can be used for validating identifiers passed in resources and data sources.
//
// Typically, we expect passed identifiers to be a variation of sdk.ObjectIdentifier.
// For now, we're expecting implementations of sdk.ObjectIdentifier, because we won't support sdk.ExternalObjectIdentifiers.
// The reason behind it is that the functions that parse identifiers are not able to differentiate between
// sdk.ExternalObjectIdentifiers and sdk.DatabaseObjectIdentifier or sdk.SchemaObjectIdentifier.
// That's because sdk.ExternalObjectIdentifiers has varying parts count (2 or 3).
//
// To use this function, pass it as a validation function on identifier field with generic parameter set to the desired sdk.ObjectIdentifier implementation.
func IsValidIdentifier[T sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier | sdk.TableColumnIdentifier]() schema.SchemaValidateDiagFunc {
	return func(value any, path cty.Path) diag.Diagnostics {
		if _, ok := value.(string); !ok {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "Invalid object identifier type",
					Detail:        fmt.Sprintf("Expected object string type, but got: %T. This is a provider error please file a report: https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/new/choose", value),
					AttributePath: path,
				},
			}
		}

		// TODO(SNOW-1495079): Right now we have to skip validation for AccountObjectIdentifier to handle a case where identifier contains dots
		// TODO(SNOW-1495079): with sdk.AccountObjectIdentifier{} (or a new type of identifier) we should be able to validate individual part of the identifier field (e.g. "database" or "schema" field)
		if _, ok := any(sdk.AccountObjectIdentifier{}).(T); ok {
			return nil
		}

		stringValue := value.(string)
		id, err := helpers.DecodeSnowflakeParameterID(stringValue)
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Unable to parse the identifier",
					Detail: fmt.Sprintf(
						"Unable to parse the identifier: %s. Make sure you are using the correct form of the fully qualified name for this field: %s.\nOriginal Error: %s",
						stringValue,
						getExpectedIdentifierRepresentationFromGeneric[T](),
						err.Error(),
					),
					AttributePath: path,
				},
			}
		}

		if _, ok := id.(T); !ok {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Invalid identifier type",
					Detail: fmt.Sprintf(
						"Expected %s identifier type, but got: %T. The correct form of the fully qualified name for this field is: %s, but was %s",
						reflect.TypeOf(new(T)).Elem().Name(),
						id,
						getExpectedIdentifierRepresentationFromGeneric[T](),
						getExpectedIdentifierRepresentationFromParam(id),
					),
					AttributePath: path,
				},
			}
		}

		return nil
	}
}

func getExpectedIdentifierRepresentationFromGeneric[T sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier | sdk.TableColumnIdentifier]() string {
	return getExpectedIdentifierForm(new(T))
}

func getExpectedIdentifierRepresentationFromParam(id sdk.ObjectIdentifier) string {
	return getExpectedIdentifierForm(id)
}

func getExpectedIdentifierForm(id any) string {
	switch id.(type) {
	case sdk.AccountObjectIdentifier, *sdk.AccountObjectIdentifier:
		return "<name>"
	case sdk.DatabaseObjectIdentifier, *sdk.DatabaseObjectIdentifier:
		return "<database_name>.<name>"
	case sdk.SchemaObjectIdentifier, *sdk.SchemaObjectIdentifier:
		return "<database_name>.<schema_name>.<name>"
	case sdk.TableColumnIdentifier, *sdk.TableColumnIdentifier:
		return "<database_name>.<schema_name>.<table_name>.<column_name>"
	}
	return ""
}

// StringInSlice has the same implementation as validation.StringInSlice, but adapted to schema.SchemaValidateDiagFunc
func StringInSlice(valid []string, ignoreCase bool) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		v, ok := i.(string)
		if !ok {
			return diag.Errorf("expected type of %v to be string", path)
		}

		for _, str := range valid {
			if v == str || (ignoreCase && strings.EqualFold(v, str)) {
				return nil
			}
		}

		return diag.Errorf("expected %v to be one of %q, got %s", path, valid, v)
	}
}

var ValidateBooleanString = StringInSlice([]string{provider.BooleanTrue, provider.BooleanFalse}, false)

var ValidateBooleanStringWithDefault = StringInSlice([]string{provider.BooleanTrue, provider.BooleanFalse, provider.BooleanDefault}, false)
