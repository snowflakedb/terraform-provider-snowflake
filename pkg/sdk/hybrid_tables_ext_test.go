package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHybridTableDetailsRow_SplitTypeAndCollation(t *testing.T) {
	testCases := []struct {
		Name              string
		Value             string
		ExpectedType      string
		ExpectedCollation *string
	}{
		{
			Name:              "with utf8",
			Value:             "VARCHAR(10) COLLATE 'utf8'",
			ExpectedType:      "VARCHAR(10)",
			ExpectedCollation: String("utf8"),
		},
		{
			Name:              "with locale",
			Value:             "VARCHAR(10) COLLATE 'en_US'",
			ExpectedType:      "VARCHAR(10)",
			ExpectedCollation: String("en_US"),
		},
		{
			Name:              "with multiple specifiers",
			Value:             "VARCHAR(10) COLLATE 'fr_CA-ai-pi-trim'",
			ExpectedType:      "VARCHAR(10)",
			ExpectedCollation: String("fr_CA-ai-pi-trim"),
		},
		{
			Name:              "with empty collation",
			Value:             "VARCHAR(10) COLLATE ''",
			ExpectedType:      "VARCHAR(10)",
			ExpectedCollation: String(""),
		},
		{
			Name:              "without collation",
			Value:             "NUMBER(38, 0)",
			ExpectedType:      "NUMBER(38, 0)",
			ExpectedCollation: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			row := hybridTableDetailsRow{Type: tc.Value}
			actualType, actualCollation := row.splitTypeAndCollation()
			assert.Equal(t, tc.ExpectedType, actualType)
			if tc.ExpectedCollation == nil {
				assert.Nil(t, actualCollation)
			} else {
				assert.Equal(t, *tc.ExpectedCollation, *actualCollation)
			}
		})
	}
}
