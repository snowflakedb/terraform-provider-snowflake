package sdk

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMappingHelpers(t *testing.T) {
	t.Run("mapNullString", func(t *testing.T) {
		testCases := []struct {
			Input          sql.NullString
			ExpectedOutput *string
		}{
			{Input: sql.NullString{Valid: true, String: "test"}, ExpectedOutput: String("test")},
			{Input: sql.NullString{Valid: true, String: ""}, ExpectedOutput: String("")},
			{Input: sql.NullString{Valid: false, String: "test"}, ExpectedOutput: nil},
			{Input: sql.NullString{Valid: false, String: ""}, ExpectedOutput: nil},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("valid: %t, string: %s", tc.Input.Valid, tc.Input.String), func(t *testing.T) {
				var output *string
				mapNullString(&output, tc.Input)
				if tc.ExpectedOutput != nil {
					require.NotNil(t, output)
					assert.Equal(t, *tc.ExpectedOutput, *output)
				} else {
					assert.Nil(t, output)
				}
			})
		}
	})

	t.Run("mapNullBool", func(t *testing.T) {
		testCases := []struct {
			Input          sql.NullBool
			ExpectedOutput *bool
		}{
			{Input: sql.NullBool{Valid: true, Bool: true}, ExpectedOutput: Bool(true)},
			{Input: sql.NullBool{Valid: true, Bool: false}, ExpectedOutput: Bool(false)},
			{Input: sql.NullBool{Valid: false, Bool: true}, ExpectedOutput: nil},
			{Input: sql.NullBool{Valid: false, Bool: false}, ExpectedOutput: nil},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("valid: %t, bool: %t", tc.Input.Valid, tc.Input.Bool), func(t *testing.T) {
				var output *bool
				mapNullBool(&output, tc.Input)
				if tc.ExpectedOutput != nil {
					require.NotNil(t, output)
					assert.Equal(t, *tc.ExpectedOutput, *output)
				} else {
					assert.Nil(t, output)
				}
			})
		}
	})

	t.Run("mapNullStringWithMapping", func(t *testing.T) {
		testCases := []struct {
			Input          sql.NullString
			ExpectedOutput *ListingState
		}{
			{Input: sql.NullString{Valid: true, String: "DRAFT"}, ExpectedOutput: Pointer(ListingStateDraft)},
			{Input: sql.NullString{Valid: true, String: "test"}, ExpectedOutput: nil},
			{Input: sql.NullString{Valid: true, String: ""}, ExpectedOutput: nil},
			{Input: sql.NullString{Valid: false, String: "DRAFT"}, ExpectedOutput: nil},
			{Input: sql.NullString{Valid: false, String: "test"}, ExpectedOutput: nil},
			{Input: sql.NullString{Valid: false, String: ""}, ExpectedOutput: nil},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("valid: %t, string: %s", tc.Input.Valid, tc.Input.String), func(t *testing.T) {
				var output *ListingState
				mapNullStringWithMapping(&output, tc.Input, ToListingState)
				if tc.ExpectedOutput != nil {
					require.NotNil(t, output)
					assert.Equal(t, *tc.ExpectedOutput, *output)
				} else {
					assert.Nil(t, output)
				}
			})
		}
	})

	t.Run("mapStringWithMapping", func(t *testing.T) {
		testCases := []struct {
			Input          string
			ExpectedOutput *ListingState
		}{
			{Input: "DRAFT", ExpectedOutput: Pointer(ListingStateDraft)},
			{Input: "test", ExpectedOutput: Pointer(ListingState(""))},
			{Input: "", ExpectedOutput: Pointer(ListingState(""))},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("value: %s", tc.Input), func(t *testing.T) {
				var output ListingState
				mapStringWithMapping(&output, tc.Input, ToListingState)
				if tc.ExpectedOutput != nil {
					require.NotNil(t, output)
					assert.Equal(t, *tc.ExpectedOutput, output)
				} else {
					assert.Nil(t, output)
				}
			})
		}
	})
}
