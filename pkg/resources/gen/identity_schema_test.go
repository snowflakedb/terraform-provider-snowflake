package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

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

func TestIdentitySchema_hierarchyForceNew(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.AuthenticationPolicy}, nil, nil)
	require.True(t, model.Hooks.IdentitySchema.Generate)

	got := renderSchemaTemplate(t, "identitySchema", model)
	require.Contains(t, got, "blocklistedCharactersFieldDescription(\"Specifies the identifier for the authentication policy; must be unique for the schema in which the authentication policy is created.\")")
	require.Contains(t, got, "\"database\"")
	require.Contains(t, got, "\"schema\"")
	require.Contains(t, got, "ForceNew: true")
	require.Contains(t, got, "DiffSuppressFunc: suppressIdentifierQuoting")
}

func TestIdentitySchema_plainHistoricalSentence(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.Pipe}, nil, nil)
	require.True(t, model.Hooks.IdentitySchema.Generate)

	got := renderSchemaTemplate(t, "identitySchema", model)
	require.Contains(t, got, "\"Specifies the identifier for the pipe; must be unique for the schema in which the pipe is created.\"")
	require.Contains(t, got, "\"The database in which to create the pipe.\"")
	require.Contains(t, got, "\"The schema in which to create the pipe.\"")
	require.Contains(t, got, "ForceNew: true")
	require.NotContains(t, got, "blocklistedCharactersFieldDescription")
	require.NotContains(t, got, "DiffSuppressFunc")
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
