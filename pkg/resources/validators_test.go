package resources

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/stretchr/testify/assert"
)

func Test_sdkValidation(t *testing.T) {
	genericNormalize := func(value string) (any, error) {
		if value == "ok" {
			return "ok", nil
		} else {
			return nil, fmt.Errorf("incorrect value %s", value)
		}
	}

	t.Run("valid generic normalize", func(t *testing.T) {
		valid := "ok"

		diag := sdkValidation(genericNormalize)(valid, cty.IndexStringPath("path"))

		assert.Empty(t, diag)
	})

	t.Run("invalid generic normalize", func(t *testing.T) {
		invalid := "nok"

		diag := sdkValidation(genericNormalize)(invalid, cty.IndexStringPath("path"))

		assert.Len(t, diag, 1)
		assert.Contains(t, diag[0].Summary, fmt.Sprintf("incorrect value %s", invalid))
	})

	t.Run("valid warehouse size", func(t *testing.T) {
		valid := string(sdk.WarehouseSizeSmall)

		diag := sdkValidation(sdk.ToWarehouseSize)(valid, cty.IndexStringPath("path"))

		assert.Empty(t, diag)
	})

	t.Run("invalid warehouse size", func(t *testing.T) {
		invalid := "SMALLa"

		diag := sdkValidation(sdk.ToWarehouseSize)(invalid, cty.IndexStringPath("path"))

		assert.Len(t, diag, 1)
		assert.Contains(t, diag[0].Summary, fmt.Sprintf("invalid warehouse size: %s", invalid))
	})
}

func Test_IsValidAccountIdentifier(t *testing.T) {
	testCases := []struct {
		Name  string
		Value any
		Error string
	}{
		{
			Name:  "validation: invalid value type",
			Value: 123,
			Error: "Expected schema string type, but got: int",
		},
		{
			Name:  "validation: account locator",
			Value: "ABC12345",
			Error: "Unable to parse the account identifier: ABC12345. Make sure you are using the correct form of the fully qualified account name: <organization_name>.<account_name>.",
		},
		{
			Name:  "validation: identifier too long",
			Value: "a.b.c",
			Error: "Unable to parse the account identifier: a.b.c. Make sure you are using the correct form of the fully qualified account name: <organization_name>.<account_name>.",
		},
		{
			Name:  "correct account object identifier",
			Value: "a.b",
		},
		{
			Name:  "correct account object identifier - quoted",
			Value: `"a"."b"`,
		},
		{
			Name:  "correct account object identifier - mixed quotes",
			Value: `a."b"`,
		},
		{
			Name:  "correct account object identifier - dot inside",
			Value: `a."b.c"`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			diag := IsValidAccountIdentifier()(tt.Value, cty.IndexStringPath("path"))
			if tt.Error != "" {
				assert.Len(t, diag, 1)
				assert.Contains(t, diag[0].Detail, tt.Error)
			} else {
				assert.Empty(t, diag)
			}
		})
	}
}

func Test_isNotEqualTo(t *testing.T) {
	testCases := []struct {
		Name             string
		Value            any
		NotExpectedValue string
		ExpectedError    string
	}{
		{
			Name:             "nil value",
			Value:            nil,
			NotExpectedValue: "123",
		},
		{
			Name:             "int value",
			Value:            123,
			NotExpectedValue: "123",
			ExpectedError:    "isNotEqualTo validator: expected string type, got int",
		},
		{
			Name:             "value equal to invalid one",
			Value:            "123",
			NotExpectedValue: "123",
			ExpectedError:    "invalid value (123) set for a field [{{} path}]. error message.",
		},
		{
			Name:             "value equal to valid one",
			Value:            "456",
			NotExpectedValue: "123",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			diag := isNotEqualTo(tt.NotExpectedValue, "error message.")(tt.Value, cty.GetAttrPath("path"))
			if tt.ExpectedError != "" {
				assert.Len(t, diag, 1)
				if diag[0].Detail != "" {
					assert.Contains(t, diag[0].Detail, tt.ExpectedError)
				} else {
					assert.Contains(t, diag[0].Summary, tt.ExpectedError)
				}
			} else {
				assert.Empty(t, diag)
			}
		})
	}
}

func Test_IsValidListingName(t *testing.T) {
	testCases := []struct {
		ListingName string
	}{
		{"listing"},
		{"\"listing\""},
		{"listing_name"},
		{"listing_name_123"},
		{"LISTING_NAME_123"},
	}

	invalidCases := []struct {
		ListingName string
	}{
		{"_invalid"},
		{"_INVALID"},
		{"1_invalid"},
		{"*_invalid"},
		{" _invalid"},
		{"invalid_name()"},
		{"INVALID_NAME()"},
		{"invalid name"},
		{"INVALID NAME"},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("valid name: %s", tt.ListingName), func(t *testing.T) {
			diag := IsValidListingName(tt.ListingName, cty.IndexStringPath("path"))
			assert.False(t, diag.HasError())
		})
	}

	for _, tt := range invalidCases {
		t.Run(fmt.Sprintf("invalid name: %s", tt.ListingName), func(t *testing.T) {
			diag := IsValidListingName(tt.ListingName, cty.IndexStringPath("path"))
			assert.True(t, diag.HasError())
		})
	}
}
