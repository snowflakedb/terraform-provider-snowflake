package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_RequiresUpdate_omitsUpdateWhenEveryFieldIsForceNew(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: ssot.Object{
		Name:                      "Widget",
		IdType:                    "SchemaObjectIdentifier",
		IdentityNameForceNew:      true,
		IdentityHierarchyForceNew: true,
		Attributes: []ssot.Attribute{{
			Name:     "from",
			Type:     ssot.AttrString,
			Required: true,
			ForceNew: true,
		}},
	}}, nil, nil)

	require.False(t, model.RequiresUpdate())
	constructor := renderTemplate(t, "constructor", model)
	require.NotContains(t, constructor, "UpdateContext")
	require.Contains(t, constructor, "ReadContext:   widgetReadExt,\nDeleteContext:")
	require.NotContains(t, renderLifecycle(t, model), "func UpdateWidget(")
}

func Test_RequiresUpdate_keepsUpdateWhenAFieldCanChangeInPlace(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: ssot.Object{
		Name:                      "Widget",
		IdType:                    "SchemaObjectIdentifier",
		IdentityNameForceNew:      true,
		IdentityHierarchyForceNew: true,
		Attributes: []ssot.Attribute{{
			Name: "comment",
			Type: ssot.AttrString,
		}},
	}}, nil, nil)

	require.True(t, model.RequiresUpdate())
	constructor := renderTemplate(t, "constructor", model)
	require.Contains(t, constructor, "UpdateContext: widgetUpdateExt,\nDeleteContext:")
	require.NotContains(t, constructor, "UpdateContext: widgetUpdateExt,\n\n")
	require.Contains(t, renderLifecycle(t, model), "func UpdateWidget(")
}
