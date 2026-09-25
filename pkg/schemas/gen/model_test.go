package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ModelFromStructDetails_SkipFields(t *testing.T) {
	details := ShowResultSchemaDetails{
		SkipFields: []string{"skipped_field", "also_skipped"},
		StructDetails: genhelpers.StructDetails{
			Name: "sdk.Example",
			Fields: []genhelpers.Field{
				{Name: "Kept", ConcreteType: "string", UnderlyingType: "string"},
				{Name: "SkippedField", ConcreteType: "string", UnderlyingType: "string"},
				{Name: "AlsoSkipped", ConcreteType: "int", UnderlyingType: "int"},
			},
		},
	}

	model := ModelFromStructDetails(details, nil)

	require.Len(t, model.SchemaFields, 3)
	assert.Equal(t, "kept", model.SchemaFields[0].Name)
	assert.Equal(t, "Kept", model.SchemaFields[0].OriginalName)
	assert.False(t, model.SchemaFields[0].Skipped)
	assert.Equal(t, "skipped_field", model.SchemaFields[1].Name)
	assert.True(t, model.SchemaFields[1].Skipped)
	assert.False(t, model.SchemaFields[1].Manual)
	assert.Equal(t, "also_skipped", model.SchemaFields[2].Name)
	assert.True(t, model.SchemaFields[2].Skipped)
	assert.False(t, model.SchemaFields[2].Manual)
	assert.False(t, model.AdditionalMapping)
}

func Test_ModelFromStructDetails_ManualFields(t *testing.T) {
	details := ShowResultSchemaDetails{
		ManualFields: []string{"manual_field"},
		StructDetails: genhelpers.StructDetails{
			Name: "sdk.Example",
			Fields: []genhelpers.Field{
				{Name: "Kept", ConcreteType: "string", UnderlyingType: "string"},
				{Name: "ManualField", ConcreteType: "string", UnderlyingType: "string"},
			},
		},
	}

	model := ModelFromStructDetails(details, nil)

	require.Len(t, model.SchemaFields, 2)
	assert.Equal(t, "kept", model.SchemaFields[0].Name)
	assert.False(t, model.SchemaFields[0].Manual)
	assert.False(t, model.SchemaFields[0].Skipped)
	assert.Equal(t, "manual_field", model.SchemaFields[1].Name)
	assert.True(t, model.SchemaFields[1].Manual)
	assert.False(t, model.SchemaFields[1].Skipped)
	assert.True(t, model.AdditionalMapping)
}

func Test_ModelFromStructDetails_SkipAndManualOverlap(t *testing.T) {
	details := ShowResultSchemaDetails{
		SkipFields:   []string{"shared"},
		ManualFields: []string{"shared"},
		StructDetails: genhelpers.StructDetails{
			Name: "sdk.Example",
			Fields: []genhelpers.Field{
				{Name: "Shared", ConcreteType: "string", UnderlyingType: "string"},
			},
		},
	}

	assert.PanicsWithValue(t, "SkipFields and ManualFields for sdk.Example overlap: shared", func() {
		ModelFromStructDetails(details, nil)
	})
}

func Test_ModelFromStructDetails_ManualFieldsUnknownKey(t *testing.T) {
	details := ShowResultSchemaDetails{
		ManualFields: []string{"not_a_field"},
		StructDetails: genhelpers.StructDetails{
			Name: "sdk.Example",
			Fields: []genhelpers.Field{
				{Name: "Kept", ConcreteType: "string", UnderlyingType: "string"},
			},
		},
	}

	assert.PanicsWithValue(t, "ManualFields for sdk.Example contain unknown schema keys: not_a_field", func() {
		ModelFromStructDetails(details, nil)
	})
}

func Test_ModelFromStructDetails_SkipFieldsUnknownKey(t *testing.T) {
	details := ShowResultSchemaDetails{
		SkipFields: []string{"not_a_field"},
		StructDetails: genhelpers.StructDetails{
			Name: "sdk.Example",
			Fields: []genhelpers.Field{
				{Name: "Kept", ConcreteType: "string", UnderlyingType: "string"},
			},
		},
	}

	assert.PanicsWithValue(t, "SkipFields for sdk.Example contain unknown schema keys: not_a_field", func() {
		ModelFromStructDetails(details, nil)
	})
}

func Test_ShowResultSchemaModel_Filename(t *testing.T) {
	assert.Equal(t, "warehouse_gen.go", ShowResultSchemaModel{Name: "Warehouse"}.Filename())
	assert.Equal(t, "password_policy_desc_gen.go", ShowResultSchemaModel{Name: "PasswordPolicyDetails", IsDescribe: true}.Filename())
	assert.Equal(t, "catalog_integration_aws_glue_desc_gen.go", ShowResultSchemaModel{Name: "CatalogIntegrationAwsGlueDetails", IsDescribe: true}.Filename())
}

func Test_ModelFromStructDetails_IsDescribe(t *testing.T) {
	details := ShowResultSchemaDetails{
		IsDescribe: true,
		StructDetails: genhelpers.StructDetails{
			Name: "sdk.PasswordPolicyDetails",
			Fields: []genhelpers.Field{
				{Name: "Name", ConcreteType: "string", UnderlyingType: "string"},
			},
		},
	}

	model := ModelFromStructDetails(details, nil)
	assert.True(t, model.IsDescribe)
	assert.Equal(t, "password_policy_desc_gen.go", model.Filename())
}

func Test_ModelFromStructDetails_UsedAsListEntry(t *testing.T) {
	details := ShowResultSchemaDetails{
		UsedAsListEntry: true,
		StructDetails: genhelpers.StructDetails{
			Name: "sdk.SecurityIntegrationProperty",
			Fields: []genhelpers.Field{
				{Name: "Name", ConcreteType: "string", UnderlyingType: "string"},
			},
		},
	}

	model := ModelFromStructDetails(details, nil)
	assert.True(t, model.UsedAsListEntry)
	assert.False(t, model.IsDescribe)
	assert.Equal(t, "security_integration_property_gen.go", model.Filename())
}

func Test_ModelFromStructDetails_AdditionalMappingInferredFromManualFields(t *testing.T) {
	withoutManual := ShowResultSchemaDetails{
		StructDetails: genhelpers.StructDetails{
			Name: "sdk.Example",
			Fields: []genhelpers.Field{
				{Name: "Name", ConcreteType: "string", UnderlyingType: "string"},
			},
		},
	}
	assert.False(t, ModelFromStructDetails(withoutManual, nil).AdditionalMapping)

	withManual := withoutManual
	withManual.ManualFields = []string{"name"}
	assert.True(t, ModelFromStructDetails(withManual, nil).AdditionalMapping)
}

func Test_ModelFromStructDetails_IsDescribeAndUsedAsListEntry(t *testing.T) {
	details := ShowResultSchemaDetails{
		IsDescribe:      true,
		UsedAsListEntry: true,
		StructDetails: genhelpers.StructDetails{
			Name: "sdk.Example",
		},
	}

	assert.PanicsWithValue(t, "sdk.Example: IsDescribe and UsedAsListEntry are mutually exclusive", func() {
		ModelFromStructDetails(details, nil)
	})
}
