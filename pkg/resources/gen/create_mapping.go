package gen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// CreateConstructorArg is one constructor argument of a create request, after the id for the object request.
// Var is the local variable, Get is the d.Get expression, and a non-empty EnumName converts the value with sdk.To<EnumName>.
type CreateConstructorArg struct {
	Var      string
	Get      string
	EnumName string
}

// CreateOptionalCall is one create helper call on a request.
// Key is the Go expression for the config key. Target is the request expression passed after the key.
type CreateOptionalCall struct {
	Helper string
	Key    string
	Target string
}

// CreateRequestBuild is how one create request is built: the object request in NewCreateRequest and
// ApplyCreateOptionals, or a block request in its parse function.
// ErrZero is the request value returned with a constructor argument error.
type CreateRequestBuild struct {
	Args    []CreateConstructorArg
	Calls   []CreateOptionalCall
	ErrZero string
}

// CreateRequestFunc is one generated function that builds an SDK request.
// The object function sets IdType, takes id sdk.<IdType>, returns the constructor pointer, and leaves optionals to ApplyCreateOptionals.
// A block function leaves IdType empty, takes the config prefix, applies optionals, and returns the request value.
type CreateRequestFunc struct {
	Name           string
	IdType         string
	ResultType     string
	NewFunc        string
	Leading        []string
	ReturnExpr     string
	ApplyOptionals bool
	CreateRequestBuild
}

// CallArgs is the constructor argument list, including any leading values such as the object id.
func (f CreateRequestFunc) CallArgs() string {
	return f.ConstructorCallArgs(f.Leading...)
}

// IdConfigFields are the config keys for a generated ParseIdFromConfig, in New<IdType> argument order.
// Account objects use name. Database objects also use database. Schema objects also use schema.
func (m ResourceModel) IdConfigFields() []string {
	return collections.Map(ssot.IdentityFields(ssot.Object{Name: m.Name, IdType: m.IdType}), func(field ssot.ClassifiedField) string {
		return field.Name
	})
}

func (m ResourceModel) CreateRequestBuild() CreateRequestBuild {
	return newCreateRequestBuild(m.Attributes, "nil")
}

// CreateRequestFunc is the generated object create function. Optionals stay in ApplyCreateOptionals.
func (m ResourceModel) CreateRequestFunc() CreateRequestFunc {
	name := m.EffectiveSdkCreateRequestObjectName()
	return CreateRequestFunc{
		Name:               m.Prefix() + "NewCreateRequest",
		IdType:             m.IdType,
		ResultType:         "*sdk.Create" + name + "Request",
		NewFunc:            "sdk.NewCreate" + name + "Request",
		Leading:            []string{"id"},
		CreateRequestBuild: m.CreateRequestBuild(),
	}
}

// CreateBlockFuncs is every generated block parse function, each parent before the blocks under it.
func (m ResourceModel) CreateBlockFuncs() []CreateRequestFunc {
	funcs := make([]CreateRequestFunc, 0)
	var walk func([]ResourceAttribute)
	walk = func(attrs []ResourceAttribute) {
		for _, attr := range attrs {
			if attr.CreateBlock != nil {
				funcs = append(funcs, *attr.CreateBlock)
			}
			walk(attr.Children)
		}
	}
	walk(m.Attributes)
	return funcs
}

// ConstructorCallArgs is the constructor argument list: leading, then each argument variable.
func (b CreateRequestBuild) ConstructorCallArgs(leading ...string) string {
	return strings.Join(append(leading, collections.Map(b.Args, func(arg CreateConstructorArg) string { return arg.Var })...), ", ")
}

func newCreateRequestBuild(fields []ResourceAttribute, errZero string) CreateRequestBuild {
	build := CreateRequestBuild{ErrZero: errZero}
	for _, field := range fields {
		if field.CreateConstructorArg != nil {
			build.Args = append(build.Args, *field.CreateConstructorArg)
		}
		if field.CreateOptionalCall != nil {
			build.Calls = append(build.Calls, *field.CreateOptionalCall)
		}
	}
	return build
}

// createScope is the request a field belongs to. stem is the object name, or the block request stem inside a parse function.
// A field inside a block is read through the parse function prefix.
type createScope struct {
	stem    string
	inBlock bool
}

func (s createScope) key(name string) string {
	if s.inBlock {
		return fmt.Sprintf("prefix+%q", name)
	}
	return strconv.Quote(name)
}

func attachCreateMappings(objectName string, src []ssot.Attribute, dst []ResourceAttribute) {
	if len(src) != len(dst) {
		panic(fmt.Sprintf("object %q: %d source attributes and %d resource attributes", objectName, len(src), len(dst)))
	}
	scope := createScope{stem: objectName}
	for i := range src {
		attachFieldCreate(scope, src[i], &dst[i])
	}
}

// attachFieldCreate records how one field is created: a constructor argument, a helper call, or a parsed block.
// The bool is false when the field has no create mapping. An object field is then left out, and a block holding it is not parsed.
func attachFieldCreate(scope createScope, src ssot.Attribute, dst *ResourceAttribute) bool {
	if src.SkipCreateBlock {
		return false
	}
	if len(src.Children) > 0 && src.CreateMapping == nil && src.Type == ssot.AttrList {
		return attachBlockCreate(scope, src, dst)
	}
	if arg, ok := constructorArg(scope, src, *dst); ok {
		dst.CreateConstructorArg = arg
		return true
	}
	call, ok := createCallFor(scope.stem, src, *dst)
	if !ok {
		return false
	}
	dst.CreateOptionalCall = &CreateOptionalCall{Helper: call.Helper, Key: scope.key(src.Name), Target: call.Target}
	return true
}

// attachBlockCreate parses an optional list block when every field in it has a create mapping.
// Its request stem is the parent stem plus the block name, so the fields are mapped like the object's own.
func attachBlockCreate(scope createScope, src ssot.Attribute, dst *ResourceAttribute) bool {
	if src.Required {
		return false
	}
	if len(src.Children) != len(dst.Children) {
		panic(fmt.Sprintf("object %q: attribute %q has %d nested attributes and %d resource attributes", scope.stem, src.Name, len(src.Children), len(dst.Children)))
	}
	blockScope := createScope{stem: scope.stem + genhelpers.SnakeCaseToCamel(src.Name), inBlock: true}
	children := append([]ResourceAttribute(nil), dst.Children...)
	for i := range src.Children {
		if !attachFieldCreate(blockScope, src.Children[i], &children[i]) {
			return false
		}
	}
	dst.Children = children

	parseName := "parse" + blockScope.stem
	call := blockCreateCall(src, parseName)
	dst.CreateOptionalCall = &CreateOptionalCall{Helper: call.Helper, Key: scope.key(src.Name), Target: call.Target}
	requestType := "sdk." + blockScope.stem + "Request"
	dst.CreateBlock = &CreateRequestFunc{
		Name:               parseName,
		ResultType:         requestType,
		NewFunc:            "sdk.New" + blockScope.stem + "Request",
		ReturnExpr:         "*request",
		ApplyOptionals:     true,
		CreateRequestBuild: newCreateRequestBuild(children, requestType+"{}"),
	}
	return true
}

// constructorArg is the constructor argument for a required field with no CreateMapping.
// A type the constructor cannot take panics for an object field. Inside a block it returns false, so the block is not parsed.
func constructorArg(scope createScope, src ssot.Attribute, resolved ResourceAttribute) (*CreateConstructorArg, bool) {
	if src.CreateMapping != nil || !src.Required {
		return nil, false
	}
	var cast string
	switch src.Type {
	case ssot.AttrInt:
		cast = "int"
	case ssot.AttrString:
		cast = "string"
	case ssot.AttrBool:
		cast = "bool"
	case ssot.AttrEnum:
		if resolved.EnumName != "" {
			cast = "string"
		}
	}
	if cast == "" {
		if scope.inBlock {
			return nil, false
		}
		panic(fmt.Sprintf("object %q: required attribute %q has no CreateMapping", scope.stem, src.Name))
	}
	return &CreateConstructorArg{
		Var:      genhelpers.FirstLetterLowercase(genhelpers.SnakeCaseToCamel(src.Name)),
		Get:      fmt.Sprintf("d.Get(%s).(%s)", scope.key(src.Name), cast),
		EnumName: resolved.EnumName,
	}, true
}

// createCallFor uses an explicit CreateMapping when one is set.
// Otherwise it picks the helper from the attribute type and the schema default the attribute rules already chose.
func createCallFor(objectName string, src ssot.Attribute, resolved ResourceAttribute) (ssot.CreateCall, bool) {
	if src.CreateMapping != nil {
		return src.CreateMapping.CreateCall(objectName, src), true
	}
	return deduceCreateCall(objectName, src, resolved)
}

func deduceCreateCall(objectName string, src ssot.Attribute, resolved ResourceAttribute) (ssot.CreateCall, bool) {
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
	case src.Type == ssot.AttrEnum && resolved.EnumName != "":
		return withFieldMapperCall("attributeMappedValueCreateBuilder", src, "sdk.To"+resolved.EnumName), true
	case src.Type == ssot.AttrIdentifier && src.IdentifierType == "AccountObjectIdentifier":
		return ssot.AccountObjectIdentifierAttributeCreateBuilder.CreateCall(objectName, src), true
	case src.Type == ssot.AttrList && src.Elem != nil:
		return ssot.ListValueCreate(collectionMapperName(objectName, src)).CreateCall(objectName, src), true
	case src.Type == ssot.AttrSet:
		return withFieldMapperCall("setMappedValueCreateBuilder", src, collectionMapperName(objectName, src)), true
	default:
		return ssot.CreateCall{}, false
	}
}

// withFieldMapperCall is attributeMappedValueCreateBuilder or setMappedValueCreateBuilder.
// The target is request.With<Field> and the mapper function name.
func withFieldMapperCall(helper string, attr ssot.Attribute, mapper string) ssot.CreateCall {
	return ssot.CreateCall{
		Helper: helper,
		Target: fmt.Sprintf("request.With%s, %s", genhelpers.SnakeCaseToCamel(attr.Name), mapper),
	}
}

// blockCreateCall parses a list block. MaxItems 1 passes the value to request.With<Field>.
// Any other list assigns &request.<Field>.
func blockCreateCall(attr ssot.Attribute, parseName string) ssot.CreateCall {
	goName := genhelpers.SnakeCaseToCamel(attr.Name)
	if attr.MaxItems != nil && *attr.MaxItems == 1 {
		return ssot.CreateCall{
			Helper: "singleBlockAttributeCreate",
			Target: fmt.Sprintf("request.With%s, %s", goName, parseName),
		}
	}
	return ssot.CreateCall{
		Helper: "listBlockAttributeCreate",
		Target: fmt.Sprintf("&request.%s, %s", goName, parseName),
	}
}

// collectionMapperName is the hand-written per-element mapper for a list or set. A set element may be a block.
// objectName is the request stem, so backup_instance_families on ComputePool becomes toComputePoolBackupInstanceFamilies,
// and allowed_methods in the AuthenticationPolicy mfa_policy block becomes toAuthenticationPolicyMfaPolicyAllowedMethods.
func collectionMapperName(objectName string, attr ssot.Attribute) string {
	return "to" + objectName + genhelpers.SnakeCaseToCamel(attr.Name)
}
