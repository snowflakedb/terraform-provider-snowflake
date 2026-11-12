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
	// NoKey leaves the attribute key out of the call, for a function that takes the requests and decides itself.
	NoKey bool
}

func attachUpdateMappings(obj ssot.Object, src []ssot.Attribute, dst []ResourceAttribute) {
	objectName := obj.Name
	setGroups, unsetGroups := obj.AlterSetGroups(), obj.AlterUnsetGroups()
	inSet, inUnset := attributeNames(obj.SetAttributes()), attributeNames(obj.UnsetAttributes())
	// An object whose ALTER spec has a SET or UNSET block updates the attributes of those blocks, and no others. An object
	// without such a block, whose ALTER is not described by the spec, updates every attribute that has an update helper.
	restricted := obj.HasAlterSetBlock() || obj.HasAlterUnsetBlock()
	if len(src) != len(dst) {
		panic(fmt.Sprintf("object %q: %d source attributes and %d resource attributes", objectName, len(src), len(dst)))
	}
	for i := range src {
		// TODO (next PRs): Support nested attributes whose fields are updated one by one. A nested attribute is updated as one
		// value only when its mapping says so.
		if len(src[i].NestedFields()) > 0 && !ssot.UpdatesWholeValue(src[i].UpdateMapping) {
			continue
		}
		name := src[i].Name
		if restricted && !inSet[name] && !inUnset[name] {
			if src[i].UpdateMapping != nil {
				panic(fmt.Sprintf("object %q: attribute %q has an UpdateMapping, but it is in neither the SET nor the UNSET block", objectName, name))
			}
			continue
		}
		// An attribute that the object cannot unset is only ever set.
		setOnly := restricted && !inUnset[name]
		call, ok := updateCallFor(objectName, src[i], dst[i], setOnly)
		if !ok {
			if restricted && updatable(src[i]) {
				panic(fmt.Sprintf("object %q: attribute %q is in the SET or UNSET block, but no update helper can be deduced for it, so give it an UpdateMapping or SkipUpdate", objectName, name))
			}
			continue
		}
		if call.Set && obj.HasAlterSpec() && !obj.HasAlterSetBlock() {
			panic(fmt.Sprintf("object %q: attribute %q is updated with SET, but the ALTER spec has no SET block", objectName, src[i].Name))
		}
		if call.Unset && obj.HasAlterSpec() && !obj.HasAlterUnsetBlock() {
			panic(fmt.Sprintf("object %q: attribute %q is updated with UNSET, but the ALTER spec has no UNSET block, so give it a set-only UpdateMapping or SkipUpdate", objectName, src[i].Name))
		}
		setRequest := "set" + groupOf(setGroups, src[i])
		unsetRequest := "unset" + groupOf(unsetGroups, src[i])
		dst[i].UpdateCall = &UpdateCall{Helper: call.Helper, Target: updateTarget(call, src[i], setRequest, unsetRequest), NoKey: call.Requests}
	}
}

func attributeNames(attributes []*ssot.Attribute) map[string]bool {
	names := make(map[string]bool, len(attributes))
	for _, attribute := range attributes {
		if attribute != nil {
			names[attribute.Name] = true
		}
	}
	return names
}

// updatable reports whether the update can change the attribute at all. ForceNew and CreateOnly attributes are create-only,
// and a Readonly one is not written.
func updatable(src ssot.Attribute) bool {
	return !src.SkipUpdate && !src.ForceNew && !src.CreateOnly && !src.Readonly
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
// attribute, then the mapper if the helper takes one. The setters are With methods, or request fields for a helper that
// takes pointers.
func updateTarget(call ssot.UpdateCall, attr ssot.Attribute, setRequest, unsetRequest string) string {
	goName := genhelpers.SnakeCaseToCamel(attr.Name)
	if call.Field != "" {
		goName = call.Field
	}
	if call.Requests {
		var requests []string
		if call.Set {
			requests = append(requests, setRequest)
		}
		if call.Unset {
			requests = append(requests, unsetRequest)
		}
		return strings.Join(requests, ", ")
	}
	setter := func(request string) string {
		if call.Fields {
			return "&" + request + "." + goName
		}
		return request + ".With" + goName
	}
	var args []string
	if call.Set {
		args = append(args, setter(setRequest))
	}
	if call.Unset {
		args = append(args, setter(unsetRequest))
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
// SkipUpdate does.
// setOnly is true for an attribute that is in the SET block and not in the UNSET block.
func updateCallFor(objectName string, src ssot.Attribute, resolved ResourceAttribute, setOnly bool) (ssot.UpdateCall, bool) {
	if src.SkipUpdate {
		return ssot.UpdateCall{}, false
	}
	if src.UpdateMapping != nil {
		return src.UpdateMapping.UpdateCall(objectName, src), true
	}
	return deduceUpdateCall(objectName, src, resolved, setOnly)
}

func deduceUpdateCall(objectName string, src ssot.Attribute, resolved ResourceAttribute, setOnly bool) (ssot.UpdateCall, bool) {
	if !updatable(src) {
		return ssot.UpdateCall{}, false
	}
	switch {
	case src.Type == ssot.AttrInt && (src.Required || setOnly):
		return ssot.IntAttributeUpdateSetOnlyBuilder.UpdateCall(objectName, src), true
	case src.Type == ssot.AttrInt && resolved.DefaultExpr == "IntDefault":
		return ssot.IntAttributeWithSpecialDefaultUpdateBuilder.UpdateCall(objectName, src), true
	case src.Type == ssot.AttrInt:
		return ssot.IntAttributeUpdateBuilder.UpdateCall(objectName, src), true
	case src.Type == ssot.AttrBool && src.NativeBool:
		return ssot.BooleanAttributeUpdateSetOnlyBuilder.UpdateCall(objectName, src), true
	case src.Type == ssot.AttrBool && resolved.DefaultExpr == "BooleanDefault" && setOnly:
		return ssot.BooleanStringAttributeUpdateSetOnly.UpdateCall(objectName, src), true
	case src.Type == ssot.AttrBool && resolved.DefaultExpr == "BooleanDefault":
		return ssot.BooleanStringAttributeUpdateBuilder.UpdateCall(objectName, src), true
	case src.Type == ssot.AttrString && setOnly:
		return ssot.StringAttributeUpdateSetOnlyNotEmpty(src.SetSql.Field).UpdateCall(objectName, src), true
	case src.Type == ssot.AttrString:
		return ssot.StringAttributeUpdateBuilder.UpdateCall(objectName, src), true
	case src.Type == ssot.AttrEnum && resolved.EnumName != "" && !setOnly:
		return ssot.AttributeMappedValueUpdate("sdk.To"+resolved.EnumName).UpdateCall(objectName, src), true
	case src.Type == ssot.AttrIdentifier && src.IdentifierType == "AccountObjectIdentifier" && !setOnly:
		return ssot.AccountObjectIdentifierAttributeUpdate.UpdateCall(objectName, src), true
	case src.Type == ssot.AttrList && setOnly && isMappedCollection(src):
		return ssot.ListValueUpdateSetOnly(collectionMapperName(objectName, src)).UpdateCall(objectName, src), true
	case src.Type == ssot.AttrList && isMappedCollection(src):
		return ssot.ListValueUpdate(collectionMapperName(objectName, src)).UpdateCall(objectName, src), true
	default:
		return ssot.UpdateCall{}, false
	}
}

// isMappedCollection reports whether the elements of a list are enums or account object identifiers, which a mapper function
// named like the one of create converts.
func isMappedCollection(src ssot.Attribute) bool {
	return src.IsEnumCollection() || (src.Elem != nil && src.Elem.Type == ssot.AttrIdentifier && src.Elem.IdentifierType == "AccountObjectIdentifier")
}
