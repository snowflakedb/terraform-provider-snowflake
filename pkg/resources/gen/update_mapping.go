package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// UpdateCall is one helper call inside ApplySetUnset.
// Target is the request expression passed after the attribute key.
type UpdateCall struct {
	Helper string
	Target string
}

func attachUpdateMappings(objectName string, src []ssot.Attribute, dst []ResourceAttribute) {
	if len(src) != len(dst) {
		panic(fmt.Sprintf("object %q: %d source attributes and %d resource attributes", objectName, len(src), len(dst)))
	}
	for i := range src {
		// TODO (next PRs): Support nested attributes.
		if len(src[i].Children) > 0 {
			continue
		}
		call, ok := updateCallFor(objectName, src[i], dst[i])
		if !ok {
			continue
		}
		dst[i].UpdateCall = &UpdateCall{Helper: call.Helper, Target: call.Target}
	}
}

// updateCallFor uses an explicit UpdateMapping when one is set.
// Otherwise it picks the helper from the attribute type and the schema default the attribute rules already chose.
// ForceNew and CreateOnly attributes are create-only, so they are left out of set/unset.
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
	case src.Type == ssot.AttrBool && resolved.DefaultExpr == "BooleanDefault":
		return ssot.BooleanStringAttributeUpdateBuilder.UpdateCall(objectName, src), true
	case src.Type == ssot.AttrString:
		return ssot.StringAttributeUpdateBuilder.UpdateCall(objectName, src), true
	case src.Type == ssot.AttrList && isCollectionOfEnum(src):
		return ssot.ListValueUpdate(collectionMapperName(objectName, src)).UpdateCall(objectName, src), true
	default:
		return ssot.UpdateCall{}, false
	}
}
