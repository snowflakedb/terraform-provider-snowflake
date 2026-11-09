package gen

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// UpdateCall is one helper call inside ApplySetUnset.
// Target is the request expression passed after the attribute key.
type UpdateCall struct {
	Helper string
	Target string
}

func attachUpdateMappings(obj ssot.Object, src []ssot.Attribute, dst []ResourceAttribute) {
	objectName := obj.Name
	setGroups, unsetGroups := obj.AlterSetGroups(), obj.AlterUnsetGroups()
	if len(src) != len(dst) {
		panic(fmt.Sprintf("object %q: %d source attributes and %d resource attributes", objectName, len(src), len(dst)))
	}
	for i := range src {
		// TODO (next PRs): Support nested attributes.
		if len(src[i].NestedFields()) > 0 {
			continue
		}
		call, ok := updateCallFor(objectName, src[i], dst[i])
		if !ok {
			continue
		}
		setRequest := "set" + groupOf(setGroups, src[i])
		unsetRequest := "unset" + groupOf(unsetGroups, src[i])
		dst[i].UpdateCall = &UpdateCall{Helper: call.Helper, Target: updateTarget(call, src[i], setRequest, unsetRequest)}
	}
}

// groupOf is the name of the group that holds the attribute. It is empty for the main group, and for an attribute that no
// group holds, which the generated update sends with the main request.
func groupOf(groups []ssot.AlterGroupSpec, attr ssot.Attribute) string {
	for _, group := range groups {
		if slices.ContainsFunc(group.Attributes, func(member *ssot.Attribute) bool { return member.Name == attr.Name }) {
			return group.Name
		}
	}
	return ssot.MainAlterGroup
}

// updateTarget is the arguments of the update helper after the attribute key: the setters of the requests that hold the
// attribute, then the mapper if the helper takes one.
func updateTarget(call ssot.UpdateCall, attr ssot.Attribute, setRequest, unsetRequest string) string {
	setter := ".With" + genhelpers.SnakeCaseToCamel(attr.Name)
	var args []string
	if call.Set {
		args = append(args, setRequest+setter)
	}
	if call.Unset {
		args = append(args, unsetRequest+setter)
	}
	if call.Mapper != "" {
		args = append(args, call.Mapper)
	}
	return strings.Join(args, ", ")
}

// updateCallFor uses an explicit UpdateMapping when one is set.
// Otherwise it picks the helper from the attribute type and the schema default the attribute rules already chose.
// ForceNew and CreateOnly attributes are create-only, so they are left out of set/unset.
// SkipCreate does not leave the attribute out of set/unset.
func updateCallFor(objectName string, src ssot.Attribute, resolved ResourceAttribute) (ssot.UpdateCall, bool) {
	if src.UpdateMapping != nil {
		return src.UpdateMapping.UpdateCall(objectName, src), true
	}
	return deduceUpdateCall(objectName, src, resolved)
}

func deduceUpdateCall(objectName string, src ssot.Attribute, resolved ResourceAttribute) (ssot.UpdateCall, bool) {
	if src.ForceNew || src.CreateOnly || src.Readonly {
		return ssot.UpdateCall{}, false
	}
	switch {
	case src.Type == ssot.AttrInt && src.Required:
		return ssot.IntAttributeUpdateSetOnlyBuilder.UpdateCall(objectName, src), true
	case src.Type == ssot.AttrInt && resolved.DefaultExpr == "IntDefault":
		return ssot.IntAttributeWithSpecialDefaultUpdateBuilder.UpdateCall(objectName, src), true
	case src.Type == ssot.AttrInt:
		return ssot.IntAttributeUpdateBuilder.UpdateCall(objectName, src), true
	case src.Type == ssot.AttrBool && src.NativeBool:
		return ssot.BooleanAttributeUpdateSetOnlyBuilder.UpdateCall(objectName, src), true
	case src.Type == ssot.AttrBool && resolved.DefaultExpr == "BooleanDefault":
		return ssot.BooleanStringAttributeUpdateBuilder.UpdateCall(objectName, src), true
	case src.Type == ssot.AttrString:
		return ssot.StringAttributeUpdateBuilder.UpdateCall(objectName, src), true
	case src.Type == ssot.AttrList && src.IsEnumCollection():
		return ssot.ListValueUpdate(collectionMapperName(objectName, src)).UpdateCall(objectName, src), true
	default:
		return ssot.UpdateCall{}, false
	}
}
