package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func TestValueConversionFromClassified(t *testing.T) {
	boolField := ssot.ClassifiedField{SdkValueType: genhelpers.SdkValueType{SdkType: "bool", SdkKind: "bool"}}
	boolConversion := valueConversionFromClassified(boolField)
	require.Equal(t, "booleanStringToBool", boolConversion.DecodeFunc)
	require.Equal(t, "booleanStringFromBool(computePool.AutoResume)", boolConversion.StateExpression("computePool.AutoResume"))

	identifier := ssot.ClassifyParameter("catalog", "AccountObjectIdentifier")
	identifierConversion := valueConversionFromClassified(identifier)
	require.Equal(t, "sdk.ParseAccountObjectIdentifier", identifierConversion.DecodeFunc)
	require.Equal(t, "parameters.Catalog.Value.FullyQualifiedName()", identifierConversion.StateExpression("parameters.Catalog.Value"))

	enumField := ssot.ClassifyParameter("log_level", "LogLevel")
	enumConversion := valueConversionFromClassified(enumField)
	require.Equal(t, "sdk.ToLogLevel", enumConversion.DecodeFunc)
	require.Equal(t, "parameters.LogLevel.Value", enumConversion.StateExpression("parameters.LogLevel.Value"))

	plain := ssot.ClassifyParameter("comment", string(ssot.AttrString))
	plainConversion := valueConversionFromClassified(plain)
	require.Empty(t, plainConversion.DecodeFunc)
	require.Equal(t, "comment", plainConversion.StateExpression("comment"))
}
