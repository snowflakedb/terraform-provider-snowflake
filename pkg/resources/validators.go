package resources

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider/validators"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func IsValidIdentifier[T sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier | sdk.TableColumnIdentifier]() schema.SchemaValidateDiagFunc {
	return validators.IsValidIdentifier[T]()
}

// IsValidAccountIdentifier is a validator that can be used for validating account identifiers passed in resources and data sources.
//
// Provider supported both account locators and organization name + account name pairs.
// The account locators are not supported, so this function accepts only the new format.
func IsValidAccountIdentifier() schema.SchemaValidateDiagFunc {
	return func(value any, path cty.Path) diag.Diagnostics {
		if _, ok := value.(string); !ok {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "Invalid schema identifier type",
					Detail:        fmt.Sprintf("Expected schema string type, but got: %T. This is a provider error please file a report: https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/new/choose", value),
					AttributePath: path,
				},
			}
		}

		stringValue := value.(string)
		_, err := helpers.DecodeSnowflakeAccountIdentifier(stringValue)
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Unable to parse the account identifier",
					Detail: fmt.Sprintf(
						"Unable to parse the account identifier: %s. Make sure you are using the correct form of the fully qualified account name: <organization_name>.<account_name>.\nOriginal Error: %s",
						stringValue,
						err.Error(),
					),
					AttributePath: path,
				},
			}
		}

		return nil
	}
}

// StringInSlice has the same implementation as validation.StringInSlice, but adapted to schema.SchemaValidateDiagFunc
func StringInSlice(valid []string, ignoreCase bool) schema.SchemaValidateDiagFunc {
	return validators.StringInSlice(valid, ignoreCase)
}

// IntInSlice has the same implementation as validation.StringInSlice, but adapted to schema.SchemaValidateDiagFunc
func IntInSlice(valid []int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		v, ok := i.(int)
		if !ok {
			return diag.Errorf("expected type of %v to be integer", path)
		}

		for _, validInt := range valid {
			if v == validInt {
				return nil
			}
		}

		return diag.Errorf("expected %v to be one of %q, got %d", path, valid, v)
	}
}

func sdkValidation[T any](normalize func(string) (T, error)) schema.SchemaValidateDiagFunc {
	return validators.NormalizeValidation(normalize)
}

var IsDataTypeValid = sdkValidation(datatypes.ParseDataType)

func isNotEqualTo(notExpectedValue string, errorMessage string) schema.SchemaValidateDiagFunc {
	return func(value any, path cty.Path) diag.Diagnostics {
		if value != nil {
			if stringValue, ok := value.(string); ok {
				if stringValue == notExpectedValue {
					return diag.Diagnostics{
						{
							Severity: diag.Error,
							Summary:  "Invalid value set",
							Detail:   fmt.Sprintf("invalid value (%s) set for a field %v. %s", notExpectedValue, path, errorMessage),
						},
					}
				}
			} else {
				return diag.Errorf("isNotEqualTo validator: expected string type, got %T", value)
			}
		}

		return nil
	}
}

func isValidSecondaryRole() func(value any, path cty.Path) diag.Diagnostics {
	return func(value any, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		if secondaryRole, ok := value.(string); !ok || strings.ToUpper(secondaryRole) != "ALL" {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("Unsupported secondary role '%s'", secondaryRole),
				Detail:        `The only supported default secondary roles value is currently 'ALL'. For more details check: https://docs.snowflake.com/en/sql-reference/sql/create-user#optional-object-properties-objectproperties.`,
				AttributePath: nil,
			})
		}
		return diags
	}
}
