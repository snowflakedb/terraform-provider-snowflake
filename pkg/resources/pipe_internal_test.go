package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const pipeCopyStatementTemplate string = "COPY INTO MY_DATABASE.MY_SCHEMA.%[4]s (%[1]s%[2]sID%[1]s%[2]s,VALUE%[1]s) FROM (%[1]s%[2]sSELECT%[1]s%[2]s%[2]sSRC.$1%[1]s%[2]s%[2]s,SRC.$2%[1]s%[2]sFROM @MY_DATABASE.MY_SCHEMA.MY_STAGE AS SRC%[1]s)%[1]sFILE_FORMAT = (%[1]s%[2]sFORMAT_NAME = MY_DATABASE.MY_SCHEMA.JSON%[1]s)%[1]sON_ERROR = CONTINUE%[3]s"

func generatePipeCopyStatement(lineEnding string, indent string, includeSemiColon bool, tableName string) string {
	semiColon := ""

	if includeSemiColon {
		semiColon = ";"
	}

	return fmt.Sprintf(
		pipeCopyStatementTemplate,
		lineEnding,
		indent,
		semiColon,
		tableName,
	)
}

func TestPipeCopyStatementDiffSuppress(t *testing.T) {
	type testCaseData struct {
		declared  string
		showPipes string
		expected  bool
	}

	testCases := map[string]testCaseData{
		"TestDiffSuppressSingleLine": {
			declared:  generatePipeCopyStatement("", " ", false, "MY_TABLE"),
			showPipes: generatePipeCopyStatement("", " ", false, "MY_TABLE"),
			expected:  true,
		},
		"TestDiffSuppressMultiLineUnix": {
			declared:  generatePipeCopyStatement("\n", "  ", false, "MY_TABLE"),
			showPipes: generatePipeCopyStatement("\n", "  ", false, "MY_TABLE"),
			expected:  true,
		},
		"TestDiffSuppressMultiLineWindows": {
			declared:  generatePipeCopyStatement("\n", "  ", false, "MY_TABLE"),
			showPipes: generatePipeCopyStatement("\r\n", "  ", false, "MY_TABLE"),
			expected:  true,
		},
		"TestDiffSuppressSingleLineWithEndStatement": {
			declared:  generatePipeCopyStatement("", " ", true, "MY_TABLE"),
			showPipes: generatePipeCopyStatement("", " ", false, "MY_TABLE"),
			expected:  true,
		},
		"TestDiffSuppressMultiLineUnixWithEndStatement": {
			declared:  generatePipeCopyStatement("\n", "  ", true, "MY_TABLE"),
			showPipes: generatePipeCopyStatement("\n", "  ", false, "MY_TABLE"),
			expected:  true,
		},
		"TestDiffSuppressMultiLineWindowsWithEndStatement": {
			declared:  generatePipeCopyStatement("\n", "  ", true, "MY_TABLE"),
			showPipes: generatePipeCopyStatement("\r\n", "  ", false, "MY_TABLE"),
			expected:  true,
		},
		"TestNoDiffSuppressSingleLine": {
			declared:  generatePipeCopyStatement("", " ", false, "MY_TABLE"),
			showPipes: generatePipeCopyStatement("", " ", false, "MY_OTHER_TABLE"),
			expected:  false,
		},
		"TestNoDiffSuppressMultiLineUnix": {
			declared:  generatePipeCopyStatement("\n", "  ", false, "MY_TABLE"),
			showPipes: generatePipeCopyStatement("\n", "  ", false, "MY_OTHER_TABLE"),
			expected:  false,
		},
		"TestNoDiffSuppressMultiLineWindows": {
			declared:  generatePipeCopyStatement("\n", "  ", false, "MY_TABLE"),
			showPipes: generatePipeCopyStatement("\r\n", "  ", false, "MY_OTHER_TABLE"),
			expected:  false,
		},
		"TestNoDiffSuppressSingleLineWithEndStatement": {
			declared:  generatePipeCopyStatement("", " ", true, "MY_TABLE"),
			showPipes: generatePipeCopyStatement("", " ", false, "MY_OTHER_TABLE"),
			expected:  false,
		},
		"TestNoDiffSuppressMultiLineUnixWithEndStatement": {
			declared:  generatePipeCopyStatement("\n", "  ", true, "MY_TABLE"),
			showPipes: generatePipeCopyStatement("\n", "  ", false, "MY_OTHER_TABLE"),
			expected:  false,
		},
		"TestNoDiffSuppressMultiLineWindowsWithEndStatement": {
			declared:  generatePipeCopyStatement("\n", "  ", true, "MY_TABLE"),
			showPipes: generatePipeCopyStatement("\r\n", "  ", false, "MY_OTHER_TABLE"),
			expected:  false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			r := require.New(t)
			r.Equal(tc.expected, pipeCopyStatementDiffSuppress("", tc.declared, tc.showPipes, nil))
		})
	}
}
