package resources

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		t.Run(name, func(t *testing.T) {
			r := require.New(t)
			r.Equal(tc.expected, pipeCopyStatementDiffSuppress("", tc.declared, tc.showPipes, nil))
		})
	}
}

func TestNotificationChannelIsAwsSnsTopic(t *testing.T) {
	testCases := map[string]struct {
		notificationChannel string
		expected            bool
	}{
		"aws-iso SNS ARN is recognized": {
			notificationChannel: "arn:aws-iso:sns:us-iso-east-1:123456789012:topic",
			expected:            true,
		},
		"commercial SNS ARN is recognized": {
			notificationChannel: "arn:aws:sns:us-east-1:123456789012:topic",
			expected:            true,
		},
		"aws-us-gov SNS ARN is recognized": {
			notificationChannel: "arn:aws-us-gov:sns:us-gov-west-1:123456789012:topic",
			expected:            true,
		},
		"commercial SQS ARN is rejected": {
			notificationChannel: "arn:aws:sqs:us-east-1:123456789012:sf-snowpipe",
			expected:            false,
		},
		"aws-iso SQS ARN is rejected": {
			notificationChannel: "arn:aws-iso:sqs:us-iso-east-1:123456789012:sf-snowpipe",
			expected:            false,
		},
		"empty notification channel is rejected": {
			notificationChannel: "",
			expected:            false,
		},
		"sns substring outside the service field is rejected": {
			notificationChannel: "arn:aws:sqs:us-east-1:123456789012:sns-named-queue",
			expected:            false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			r := require.New(t)
			r.Equal(tc.expected, notificationChannelIsAwsSnsTopic(tc.notificationChannel))
		})
	}
}

func TestPipeHandleExternalChanges(t *testing.T) {
	newData := func(t *testing.T, configured string, inShow string) *schema.ResourceData {
		t.Helper()
		d := schema.TestResourceDataRaw(t, pipeSchema, map[string]any{
			"name":              "p",
			"database":          "d",
			"schema":            "s",
			"copy_statement":    "copy into t from @s",
			"integration":       configured,
			"error_integration": configured,
		})
		require.NoError(t, d.Set(ShowOutputAttributeName, []map[string]any{{"integration": inShow, "error_integration": inShow}}))
		return d
	}

	t.Run("integrations changed outside of Terraform are written to the state", func(t *testing.T) {
		d := newData(t, "OLD", "OLD")
		pipe := &sdk.Pipe{Integration: "NEW", ErrorIntegration: "NEW"}

		require.NoError(t, pipeHandleExternalChanges(d, pipe, true))

		require.Equal(t, "NEW", d.Get("integration"))
		require.Equal(t, "NEW", d.Get("error_integration"))
	})

	t.Run("unchanged integrations are left alone", func(t *testing.T) {
		d := newData(t, "SAME", "SAME")
		pipe := &sdk.Pipe{Integration: "SAME", ErrorIntegration: "SAME"}

		require.NoError(t, pipeHandleExternalChanges(d, pipe, true))

		require.Equal(t, "SAME", d.Get("integration"))
		require.Equal(t, "SAME", d.Get("error_integration"))
	})

	t.Run("nothing is marked without external changes marking", func(t *testing.T) {
		d := newData(t, "OLD", "OLD")
		pipe := &sdk.Pipe{Integration: "NEW", ErrorIntegration: "NEW"}

		require.NoError(t, pipeHandleExternalChanges(d, pipe, false))

		require.Equal(t, "OLD", d.Get("integration"))
		require.Equal(t, "OLD", d.Get("error_integration"))
	})
}
