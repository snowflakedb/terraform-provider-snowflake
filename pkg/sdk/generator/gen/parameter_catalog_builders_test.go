package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/stretchr/testify/require"
)

func TestParameterCatalogBuilders(t *testing.T) {
	boolParameter := parameterdefs.ParameterDef{SqlName: "PIPE_EXECUTION_PAUSED", Kind: KindBool}
	intParameter := parameterdefs.ParameterDef{SqlName: "DATA_RETENTION_TIME_IN_DAYS", Kind: KindInt}
	textParameter := parameterdefs.ParameterDef{SqlName: "QUERY_TAG", Kind: KindString}
	emptyParameter := parameterdefs.ParameterDef{
		SqlName: "DEFAULT_DDL_COLLATION",
		Kind:    KindOfT[sdkcommons.StringAllowEmpty](),
	}
	enumParameter := parameterdefs.ParameterDef{
		SqlName: "LOG_LEVEL",
		Kind:    KindOfT[sdkcommons.LogLevel](),
	}
	identifierParameter := parameterdefs.ParameterDef{
		SqlName: "CATALOG",
		Kind:    KindOfT[sdkcommons.AccountObjectIdentifier](),
	}

	t.Run("WithParameters", func(t *testing.T) {
		root := NewQueryStruct("ParameterSet").
			WithParameters(
				boolParameter,
				intParameter,
				textParameter,
				emptyParameter,
				enumParameter,
				identifierParameter,
			).
			IntoField()

		byName := fieldsByName(root.Fields)
		tests := []struct {
			fieldName string
			wantKind  string
		}{
			{"PipeExecutionPaused", "*bool"},
			{"DataRetentionTimeInDays", "*int"},
			{"QueryTag", "*string"},
			{"DefaultDdlCollation", "*StringAllowEmpty"},
			{"LogLevel", "*LogLevel"},
			{"Catalog", "*AccountObjectIdentifier"},
		}
		for _, tt := range tests {
			t.Run(tt.fieldName, func(t *testing.T) {
				require.Equal(t, tt.wantKind, byName[tt.fieldName].Kind)
			})
		}
	})

	t.Run("WithParametersUnset", func(t *testing.T) {
		root := NewQueryStruct("ParameterUnset").
			WithParametersUnset(intParameter, enumParameter).
			IntoField()

		byName := fieldsByName(root.Fields)
		require.Contains(t, byName, "DataRetentionTimeInDays")
		require.Contains(t, byName, "LogLevel")
	})

	t.Run("ShowParametersDetails", func(t *testing.T) {
		iface := NewInterface("Schemas", "Schema", "DatabaseObjectIdentifier")
		iface.ShowParametersDetails(
			boolParameter,
			intParameter,
			textParameter,
			emptyParameter,
			enumParameter,
			identifierParameter,
		)

		require.NotNil(t, iface.ParametersDetails)
		require.Equal(t, []ParameterDetailField{
			{FieldName: "PipeExecutionPaused", Key: "PIPE_EXECUTION_PAUSED", GoType: "bool", Parser: "strconv.ParseBool"},
			{FieldName: "DataRetentionTimeInDays", Key: "DATA_RETENTION_TIME_IN_DAYS", GoType: "int", Parser: "strconv.Atoi"},
			{FieldName: "QueryTag", Key: "QUERY_TAG", GoType: "string", Parser: "identityParse"},
			{FieldName: "DefaultDdlCollation", Key: "DEFAULT_DDL_COLLATION", GoType: "string", Parser: "identityParse"},
			{FieldName: "LogLevel", Key: "LOG_LEVEL", GoType: "LogLevel", Parser: "ToLogLevel"},
			{FieldName: "Catalog", Key: "CATALOG", GoType: "AccountObjectIdentifier", Parser: "ParseAccountObjectIdentifier"},
		}, iface.ParametersDetails.Fields)

		require.Len(t, iface.CustomMethods, 1)
		require.Equal(t, "ShowParametersDetails", iface.CustomMethods[0].Name)
		require.Equal(t, []*MethodParameter{NewMethodParameter("id", "DatabaseObjectIdentifier")}, iface.CustomMethods[0].Parameters)
		require.Equal(t, []string{"*SchemaParametersDetails", "error"}, iface.CustomMethods[0].ReturnTypes)
	})

	t.Run("ParameterSqlToFieldName", func(t *testing.T) {
		tests := []struct {
			sql  string
			want string
		}{
			{"DATA_RETENTION_TIME_IN_DAYS", "DataRetentionTimeInDays"},
			{"DEFAULT_DDL_COLLATION", "DefaultDdlCollation"},
		}
		for _, tt := range tests {
			t.Run(tt.sql, func(t *testing.T) {
				require.Equal(t, tt.want, ParameterSqlToFieldName(parameterdefs.ParameterDef{SqlName: tt.sql}))
			})
		}
	})
}

func fieldsByName(fields []Field) map[string]Field {
	out := make(map[string]Field, len(fields))
	for _, f := range fields {
		out[f.Name] = f
	}
	return out
}
