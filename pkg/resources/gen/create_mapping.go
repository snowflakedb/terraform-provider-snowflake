package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// CreateConstructorArg is one argument of NewCreate<Object>Request after the id.
type CreateConstructorArg struct {
	Get string
}

// CreateOptionalCall is one helper call inside ApplyCreateOptionals.
// Target is the request expression passed after the attribute key.
type CreateOptionalCall struct {
	Helper string
	Target string
}

// IdConfigFields are the config keys for a generated ParseIdFromConfig, in New<IdType> argument order.
// Account objects use name. Database objects also use database. Schema objects also use schema.
func (m ResourceModel) IdConfigFields() []string {
	return collections.Map(ssot.IdentityFields(ssot.Object{Name: m.Name, IdType: m.IdType}), func(field ssot.ClassifiedField) string {
		return field.Name
	})
}

func attachCreateMappings(objectName, idType string, src []ssot.Attribute, dst []ResourceAttribute) {
	if len(src) != len(dst) {
		panic(fmt.Sprintf("object %q: %d source attributes and %d resource attributes", objectName, len(src), len(dst)))
	}
	for i := range src {
		// TODO (next PRs): Support nested attributes.
		if len(src[i].Children) > 0 {
			continue
		}
		if get, ok := src[i].ConstructorGet(); ok {
			dst[i].CreateConstructorArg = &CreateConstructorArg{Get: get}
			continue
		}
		call, ok := createCallFor(objectName, idType, src[i], dst[i])
		if !ok {
			continue
		}
		dst[i].CreateOptionalCall = &CreateOptionalCall{Helper: call.Helper, Target: call.Target}
	}
}

// createCallFor uses an explicit CreateMapping when one is set.
// Otherwise it picks the helper from the attribute type and the schema default the attribute rules already chose.
func createCallFor(objectName, idType string, src ssot.Attribute, resolved ResourceAttribute) (ssot.CreateCall, bool) {
	if src.CreateMapping != nil {
		return src.CreateMapping.CreateCall(objectName, src), true
	}
	return deduceCreateCall(objectName, idType, src, resolved)
}

func deduceCreateCall(objectName, idType string, src ssot.Attribute, resolved ResourceAttribute) (ssot.CreateCall, bool) {
	if src.Required || src.Readonly {
		return ssot.CreateCall{}, false
	}
	switch {
	case src.Type == ssot.AttrBool && resolved.DefaultExpr == "BooleanDefault":
		return ssot.BooleanStringAttributeCreateBuilder.CreateCall(objectName, src), true
	case src.Type == ssot.AttrInt && resolved.DefaultExpr == "IntDefault":
		return ssot.IntAttributeWithSpecialDefaultCreateBuilder.CreateCall(objectName, src), true
	case src.Type == ssot.AttrInt:
		return ssot.IntAttributeCreateBuilder.CreateCall(objectName, src), true
	case src.Type == ssot.AttrString:
		return ssot.StringAttributeCreateBuilder.CreateCall(objectName, src), true
	case src.Type == ssot.AttrIdentifier && idType == "AccountObjectIdentifier":
		return ssot.AccountObjectIdentifierAttributeCreateBuilder.CreateCall(objectName, src), true
	case src.Type == ssot.AttrList && isCollectionOfEnum(src):
		return ssot.ListValueCreate(listEnumMapperName(objectName, src)).CreateCall(objectName, src), true
	default:
		return ssot.CreateCall{}, false
	}
}

// listEnumMapperName is the hand-written per-element mapper for a list of enums.
// backup_instance_families on ComputePool becomes toComputePoolBackupInstanceFamilies.
func listEnumMapperName(objectName string, attr ssot.Attribute) string {
	return "to" + objectName + genhelpers.SnakeCaseToCamel(attr.Name)
}
