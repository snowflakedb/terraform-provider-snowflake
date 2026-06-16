package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

type objectIdentifierKind string

const (
	AccountObjectIdentifier             objectIdentifierKind = "AccountObjectIdentifier"
	DatabaseObjectIdentifier            objectIdentifierKind = "DatabaseObjectIdentifier"
	SchemaObjectIdentifier              objectIdentifierKind = "SchemaObjectIdentifier"
	SchemaObjectIdentifierWithArguments objectIdentifierKind = "SchemaObjectIdentifierWithArguments"
)

// ShowByIDFindPredicateKind selects the predicate strategy used in the generated ShowByID FindFirst call.
// Use ResolvedShowByIDFindPredicateKind() to get the fully-resolved value (manual > auto-detect > default).
type ShowByIDFindPredicateKind string

const (
	ShowByIDFindPredicateName           ShowByIDFindPredicateKind = "name"
	ShowByIDFindPredicateFullID         ShowByIDFindPredicateKind = "full_id"
	ShowByIDFindPredicateAccountName    ShowByIDFindPredicateKind = "account_name"
	ShowByIDFindPredicateNameAndLocator ShowByIDFindPredicateKind = "name_and_locator"
)

func (k ShowByIDFindPredicateKind) IsName() bool        { return k == ShowByIDFindPredicateName }
func (k ShowByIDFindPredicateKind) IsFullID() bool      { return k == ShowByIDFindPredicateFullID }
func (k ShowByIDFindPredicateKind) IsAccountName() bool { return k == ShowByIDFindPredicateAccountName }
func (k ShowByIDFindPredicateKind) IsNameAndLocator() bool {
	return k == ShowByIDFindPredicateNameAndLocator
}

func ToObjectIdentifierKind(s string) (objectIdentifierKind, error) {
	switch s {
	case "AccountObjectIdentifier":
		return AccountObjectIdentifier, nil
	case "DatabaseObjectIdentifier":
		return DatabaseObjectIdentifier, nil
	case "SchemaObjectIdentifier":
		return SchemaObjectIdentifier, nil
	case "SchemaObjectIdentifierWithArguments":
		return SchemaObjectIdentifierWithArguments, nil
	default:
		return "", fmt.Errorf("invalid string identifier type: %s", s)
	}
}

// Interface groups operations for particular object or objects family (e.g. DATABASE ROLE)
type Interface struct {
	// Name is the interface's name, e.g. "DatabaseRoles"
	Name string
	// NameSingular is the prefix/suffix which can be used to create other structs and methods, e.g. "DatabaseRole"
	NameSingular string
	// Operations contains all operations for given interface
	Operations []*Operation
	// IdentifierKind keeps identifier of the underlying object (e.g. DatabaseObjectIdentifier)
	IdentifierKind string
	// Enums contains all enum definitions for this operation group.
	Enums []*Enum
	// CustomMethods holds interface methods that have no generated implementation.
	// They will appear in the generated interface but the user is responsible for implementing them.
	CustomMethods []*CustomInterfaceMethod
	// ShowObjectTypeName overrides the suffix used in the generated ObjectType() return value.
	// If empty, NameSingular is used (producing ObjectType<NameSingular>).
	ShowObjectTypeName string
	// ShowObjectName is the name of the main object returned from this interface through Show methods family.
	ShowObjectName string
	// ShowByIDFindPredicateKind selects the predicate strategy for the generated ShowByID FindFirst call.
	ShowByIDFindPredicateKind ShowByIDFindPredicateKind

	*genhelpers.PreambleModel
	*genhelpers.ObjectGenerationSettings
}

// WithAllowedGenerationParts restricts this object to only the specified generation parts.
// Parts not listed here will be skipped during generation, even if enabled globally.
func (i *Interface) WithAllowedGenerationParts(parts ...GenerationPartName) *Interface {
	if i.ObjectGenerationSettings == nil {
		i.ObjectGenerationSettings = &genhelpers.ObjectGenerationSettings{}
	}
	i.ObjectGenerationSettings.AllowedGenerationParts = generationPartNamesToNamers(parts)
	return i
}

// WithEnabledGenerationParts enables optional (disabled-by-default) generation parts for this object.
func (i *Interface) WithEnabledGenerationParts(parts ...GenerationPartName) *Interface {
	if i.ObjectGenerationSettings == nil {
		i.ObjectGenerationSettings = &genhelpers.ObjectGenerationSettings{}
	}
	i.ObjectGenerationSettings.EnabledGenerationParts = generationPartNamesToNamers(parts)
	return i
}

func generationPartNamesToNamers(parts []GenerationPartName) []genhelpers.GenerationPartNamer {
	return collections.Map(parts, func(p GenerationPartName) genhelpers.GenerationPartNamer { return p })
}

func (i *Interface) ObjectName() string {
	return i.Name
}

func NewInterface(name string, nameSingular string, identifierKind string, operations ...*Operation) *Interface {
	return &Interface{
		Name:           name,
		NameSingular:   nameSingular,
		IdentifierKind: identifierKind,
		Operations:     operations,
	}
}

// NameLowerCased returns interface name starting with a lower case letter
func (i *Interface) NameLowerCased() string {
	return startingWithLowerCase(i.Name)
}

// SharedToOptsFields returns all nested struct fields marked with GenerateSharedToOpts
// across all operations. Used by the impl template to emit standalone toOpts() methods.
func (i *Interface) SharedToOptsFields() []*Field {
	var result []*Field
	for _, op := range i.Operations {
		if op.OptsField == nil {
			continue
		}
		collectSharedToOpts(op.OptsField, &result)
	}
	return result
}

func collectSharedToOpts(f *Field, result *[]*Field) {
	if f.GenerateSharedToOpts && !f.IsShared {
		*result = append(*result, f)
	}
	for idx := range f.Fields {
		collectSharedToOpts(&f.Fields[idx], result)
	}
}

// ObjectIdentifierKind returns the level of the object identifier (e.g. for DatabaseObjectIdentifier, it returns the prefix "Database")
func (i *Interface) ObjectIdentifierPrefix() idPrefix {
	return identifierStringToPrefix(i.IdentifierKind)
}

func (i *Interface) WithEnums(enums ...*Enum) *Interface {
	i.Enums = append(i.Enums, enums...)
	return i
}

// WithShowObjectType overrides the ObjectType constant used in the generated ObjectType() method.
// By default the generator produces `return ObjectType<NameSingular>`. Use this when that constant
// does not exist and the real constant uses a different suffix (e.g. "Account" instead of "OrganizationAccount").
func (i *Interface) WithShowObjectType(name string) *Interface {
	i.ShowObjectTypeName = name
	return i
}

// WithShowByIDFindPredicateKind sets the predicate strategy used in the generated ShowByID FindFirst call.
// All comparison expressions live in the template; this field only selects which one to emit.
func (i *Interface) WithShowByIDFindPredicateKind(kind ShowByIDFindPredicateKind) *Interface {
	i.ShowByIDFindPredicateKind = kind
	return i
}

// ResolvedShowByIDFindPredicateKind returns the fully-resolved predicate strategy for ShowByID.
func (i *Interface) ResolvedShowByIDFindPredicateKind() ShowByIDFindPredicateKind {
	if i.ShowByIDFindPredicateKind != "" {
		return i.ShowByIDFindPredicateKind
	}
	if i.IdentifierKind == string(SchemaObjectIdentifierWithArguments) {
		return ShowByIDFindPredicateFullID
	}
	return ShowByIDFindPredicateName
}
