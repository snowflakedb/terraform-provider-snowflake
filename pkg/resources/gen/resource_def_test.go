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
		require.False(t, obj.VariantBaseOnly, obj.Name)
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
			object := model(ssot.Attribute{Name: "block", Type: ssot.AttrObject, Required: required, ForceNew: true, Children: []ssot.Attribute{leaf}})
			listAttr := ssot.Attribute{
				Name: "block", Type: ssot.AttrList, Required: required, ForceNew: true, MaxItems: new(1),
				Elem: ssot.ObjectElem([]ssot.Attribute{leaf}),
			}
			if required {
				listAttr.MinItems = new(1)
			}
			require.Equal(t, model(listAttr).Attributes, object.Attributes)
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
