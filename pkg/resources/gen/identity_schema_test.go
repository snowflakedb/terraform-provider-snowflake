package gen

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

func renderIdentitySchema(t *testing.T, model ResourceModel) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, SchemaTemplate.Lookup("identitySchema").Execute(&b, model))
	return b.String()
}

func TestIdentitySchema_preservesGeneratedSentences(t *testing.T) {
	schemaModel := ModelFromInputObject(ResourceObject{Object: defs.Schema}, nil, nil)
	got := renderSchemaTemplate(t, "identitySchema", schemaModel)
	require.Contains(t, got, "blocklistedCharactersFieldDescription(\"Specifies the identifier for the schema; must be unique for the database in which the schema is created. When the name is `PUBLIC`, during creation the provider checks if this schema has already been created and, in such case, `ALTER` is used to match the desired state.\")")
	require.Contains(t, got, "DiffSuppressFunc: suppressIdentifierQuoting")
	require.NotContains(t, got, "ForceNew")

	pool := ModelFromInputObject(ResourceObject{Object: defs.ComputePool}, nil, nil)
	got = renderSchemaTemplate(t, "identitySchema", pool)
	require.Contains(t, got, "ForceNew: true")
	require.Contains(t, got, "blocklistedCharactersFieldDescription(\"Specifies the identifier for the compute pool; must be unique for the account.\")")
	require.NotContains(t, got, "\"database\"")
}

func Test_identitySchema_hierarchyForceNew(t *testing.T) {
	t.Run("schema object without the flag", func(t *testing.T) {
		got := renderIdentitySchema(t, lifecycleModel())
		require.NotContains(t, got, "ForceNew:")
		require.Contains(t, got, `"database"`)
		require.Contains(t, got, `"schema"`)
	})

	t.Run("schema object with hierarchy force new", func(t *testing.T) {
		got := renderIdentitySchema(t, lifecycleModel(func(obj *ssot.Object) {
			obj.IdentityHierarchyForceNew = true
		}))
		require.Contains(t, got, `"database": {
			Type:             tfschema.TypeString,
			Required:         true,
			ForceNew:         true,`)
		require.Contains(t, got, `"schema": {
			Type:             tfschema.TypeString,
			Required:         true,
			ForceNew:         true,`)
		require.NotContains(t, got, `"name": {
			Type:     tfschema.TypeString,
			Required: true,
			ForceNew: true,`)
	})

	t.Run("name and hierarchy force new", func(t *testing.T) {
		got := renderIdentitySchema(t, lifecycleModel(func(obj *ssot.Object) {
			obj.IdentityNameForceNew = true
			obj.IdentityHierarchyForceNew = true
		}))
		require.Contains(t, got, `"name": {
			Type:     tfschema.TypeString,
			Required: true,
			ForceNew: true,`)
		require.Contains(t, got, `"database": {
			Type:             tfschema.TypeString,
			Required:         true,
			ForceNew:         true,`)
		require.Contains(t, got, `"schema": {
			Type:             tfschema.TypeString,
			Required:         true,
			ForceNew:         true,`)
	})

	t.Run("database object", func(t *testing.T) {
		got := renderIdentitySchema(t, lifecycleModel(func(obj *ssot.Object) {
			obj.IdType = "DatabaseObjectIdentifier"
			obj.IdentityHierarchyForceNew = true
		}))
		require.Contains(t, got, `"database": {
			Type:             tfschema.TypeString,
			Required:         true,
			ForceNew:         true,`)
		require.NotContains(t, got, `"schema"`)
	})

	t.Run("account object", func(t *testing.T) {
		got := renderIdentitySchema(t, lifecycleModel(func(obj *ssot.Object) {
			obj.IdType = "AccountObjectIdentifier"
			obj.IdentityHierarchyForceNew = true
		}))
		require.NotContains(t, got, `"database"`)
		require.NotContains(t, got, `"schema"`)
	})
}

func TestIdentitySchema_hierarchyForceNew(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.AuthenticationPolicy}, nil, nil)
	require.True(t, model.Hooks.IdentitySchema.Generate)

	got := renderSchemaTemplate(t, "identitySchema", model)
	require.Contains(t, got, "blocklistedCharactersFieldDescription(\"Specifies the identifier for the authentication policy; must be unique for the schema in which the authentication policy is created.\")")
	require.Contains(t, got, "\"database\"")
	require.Contains(t, got, "\"schema\"")
	require.Contains(t, got, "ForceNew:         true")
	require.Contains(t, got, "DiffSuppressFunc: suppressIdentifierQuoting")
}

func TestIdentitySchema_legacyIdEncodingUsesQuotingSuppression(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.Pipe}, nil, nil)
	require.True(t, model.Hooks.IdentitySchema.Generate)

	got := renderSchemaTemplate(t, "identitySchema", model)
	require.Contains(t, got, "blocklistedCharactersFieldDescription(\"Specifies the identifier for the pipe; must be unique for the schema in which the pipe is created.\")")
	require.Contains(t, got, "blocklistedCharactersFieldDescription(\"The database in which to create the pipe.\")")
	require.Contains(t, got, "blocklistedCharactersFieldDescription(\"The schema in which to create the pipe.\")")
	require.Contains(t, got, "ForceNew: true")
	require.Contains(t, got, "DiffSuppressFunc: suppressIdentifierQuoting")
}

func TestIdentitySchema_nameDescription(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.Warehouse}, nil, nil)
	require.True(t, model.Hooks.IdentitySchema.Generate)

	got := renderSchemaTemplate(t, "identitySchema", model)
	require.Contains(t, got, "blocklistedCharactersFieldDescription(\"Specifies the identifier for the warehouse; must be unique for the account.\")")
	require.Contains(t, got, "DiffSuppressFunc: suppressIdentifierQuoting")
	require.NotContains(t, got, "ForceNew")
	require.NotContains(t, got, "\"database\"")
}

func TestIdentitySchema_legacyIdEncodingKeepsQuotingSuppression(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.Table}, nil, nil)
	require.True(t, model.Hooks.IdentitySchema.Generate)

	got := renderSchemaTemplate(t, "identitySchema", model)
	require.Contains(t, got, "blocklistedCharactersFieldDescription(\"Specifies the identifier for the table; must be unique for the schema in which the table is created.\")")
	require.Contains(t, got, "blocklistedCharactersFieldDescription(\"The database in which to create the table.\")")
	require.Contains(t, got, "blocklistedCharactersFieldDescription(\"The schema in which to create the table.\")")
	require.Contains(t, got, "DiffSuppressFunc: suppressIdentifierQuoting")
	require.NotContains(t, got, "ForceNew")
}
