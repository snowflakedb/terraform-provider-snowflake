package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_ResourceObject_implementsHasObjectGenerationSettings(t *testing.T) {
	o := ResourceObject{ObjectGenerationSettings: &genhelpers.ObjectGenerationSettings{}}
	_, ok := any(o).(genhelpers.HasObjectGenerationSettings)
	require.True(t, ok)
}

func Test_GetAllObjects_skipsVariantBaseOnlyObjects(t *testing.T) {
	for _, obj := range GetAllObjects() {
		require.False(t, obj.IsVariantBaseOnly(), obj.Name)
	}
}

func Test_ModelFromInputObject_singleBlockIsAListOfOne(t *testing.T) {
	leaf := ssot.Attribute{Name: "leaf", Type: ssot.AttrString, ForceNew: true}
	model := func(attr ssot.Attribute) ResourceModel {
		return ModelFromInputObject(ResourceObject{Object: ssot.Object{
			Name:       "Widget",
			NamePlural: "Widgets",
			IdType:     "AccountObjectIdentifier",
			Attributes: []ssot.Attribute{attr},
			ResourceHooks: []ssot.ResourceHookOption{
				ssot.WithApplyCreateOptionals(ssot.Generated()),
				ssot.WithApplySetUnset(ssot.Generated()),
			},
		}}, nil, nil)
	}
	for name, required := range map[string]bool{"optional": false, "required": true} {
		t.Run(name, func(t *testing.T) {
			got := model(ssot.Attribute{Name: "block", Type: ssot.AttrObject, Required: required, ForceNew: true, Children: []ssot.Attribute{leaf}})
			require.Len(t, got.Attributes, 1)
			attr := got.Attributes[0]
			require.Equal(t, SchemaTypeList, attr.TfType)
			require.Equal(t, required, attr.Required)
			require.Equal(t, !required, attr.Optional)
			require.Equal(t, new(1), attr.MaxItems)
			if required {
				require.Equal(t, new(1), attr.MinItems)
				require.Nil(t, attr.CreateOptionalCall)
			} else {
				require.Nil(t, attr.MinItems)
				require.Equal(t, "singleBlockAttributeCreate", attr.CreateOptionalCall.Helper)
			}
		})
	}
}

func Test_ModelFromInputObject_staticFieldsStayOutOfTheBlock(t *testing.T) {
	leaf := ssot.Attribute{Name: "leaf", Type: ssot.AttrString, ForceNew: true}
	static := ssot.Attribute{Name: "type_fixed", SqlOnly: true, CreateSql: ssot.StaticSQL("fixed", "TYPE = FIXED")}
	model := func(children ...ssot.Attribute) ResourceModel {
		return ModelFromInputObject(ResourceObject{Object: ssot.Object{
			Name:       "Widget",
			NamePlural: "Widgets",
			IdType:     "AccountObjectIdentifier",
			Attributes: []ssot.Attribute{{Name: "block", Type: ssot.AttrObject, ForceNew: true, Children: children}},
			ResourceHooks: []ssot.ResourceHookOption{
				ssot.WithApplyCreateOptionals(ssot.Generated()),
				ssot.WithApplySetUnset(ssot.Generated()),
			},
		}}, nil, nil)
	}
	require.Equal(t, model(leaf).Attributes, model(static, leaf).Attributes)
}
