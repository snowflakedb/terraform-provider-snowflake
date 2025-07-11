package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_AllValuesSliceName(t *testing.T) {
	tests := []struct {
		name              string
		enumName          string
		plural            string
		expectedSliceName string
	}{
		{
			name:              "default plural",
			enumName:          "Value",
			expectedSliceName: "allValues",
		},
		{
			name:              "custom plural",
			enumName:          "Policy",
			plural:            "Policies",
			expectedSliceName: "allPolicies",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enum := NewEnum(tt.enumName)
			if tt.plural != "" {
				enum.WithPlural(tt.plural)
			}
			result := enum.AllValuesSliceName()
			require.Equal(t, tt.expectedSliceName, result)
		})
	}
}

func TestEnum_ValueRepresentations(t *testing.T) {
	tests := []struct {
		name                         string
		enumName                     string
		values                       []string
		expectedValueRepresentations []EnumValueRepresentation
	}{
		{
			name:     "single value",
			enumName: "Status",
			values:   []string{"ACTIVE"},
			expectedValueRepresentations: []EnumValueRepresentation{
				{Name: "StatusActive", Value: "ACTIVE"},
			},
		},
		{
			name:     "multiple values",
			enumName: "TokenStatus",
			values:   []string{"ACTIVE_VALUE", "INACTIVE_VALUE", "EXPIRED_VALUE"},
			expectedValueRepresentations: []EnumValueRepresentation{
				{Name: "TokenStatusActiveValue", Value: "ACTIVE_VALUE"},
				{Name: "TokenStatusInactiveValue", Value: "INACTIVE_VALUE"},
				{Name: "TokenStatusExpiredValue", Value: "EXPIRED_VALUE"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enum := NewEnum(tt.enumName, tt.values...)
			result := enum.ValueRepresentations()

			require.Len(t, result, len(tt.expectedValueRepresentations))

			require.Equal(t, tt.expectedValueRepresentations, result)
		})
	}
}
