package generator

import (
	"fmt"
	"log"
)

type ShowByIDFilteringKind uint

const (
	ShowByIDLikeFiltering ShowByIDFilteringKind = iota
	ShowByIDInFiltering
	ShowByIDExtendedInFiltering
	ShowByIDApplicationNameFiltering
	ShowByIDServiceInFiltering
)

type idPrefix string

const (
	AccountIdentifierPrefix  idPrefix = "Account"
	DatabaseIdentifierPrefix idPrefix = "Database"
	SchemaIdentifierPrefix   idPrefix = "Schema"
)

func identifierStringToPrefix(s string) idPrefix {
	switch s {
	case "AccountObjectIdentifier":
		return AccountIdentifierPrefix
	case "DatabaseObjectIdentifier":
		return DatabaseIdentifierPrefix
	case "SchemaObjectIdentifier":
		return SchemaIdentifierPrefix
	case "SchemaObjectIdentifierWithArguments":
		return SchemaIdentifierPrefix
	default:
		return ""
	}
}

type ShowByIDFiltering interface {
	WithFiltering() string
}

type showByIDFilter struct {
	Name string
	Kind string
	Args string
}

func (s *showByIDFilter) WithFiltering() string {
	return fmt.Sprintf("With%s(%s{%s})", s.Name, s.Kind, s.Args)
}

func newShowByIDFiltering(name, kind, args string) ShowByIDFiltering {
	return &showByIDFilter{
		Name: name,
		Kind: kind,
		Args: args,
	}
}

func newShowByIDLikeFiltering() ShowByIDFiltering {
	return newShowByIDFiltering("Like", "Like", "Pattern: String(id.Name())")
}

func newShowByIDInFiltering(identifierKind idPrefix) ShowByIDFiltering {
	return newShowByIDFiltering("In", "In", fmt.Sprintf("%[1]v: id.%[1]vId()", identifierKind))
}

func newShowByIDExtendedInFiltering(identifierKind idPrefix) ShowByIDFiltering {
	return newShowByIDFiltering("In", "ExtendedIn", fmt.Sprintf("In: In{%[1]v: id.%[1]vId()}", identifierKind))
}

func newShowByIDServiceInFiltering(identifierKind idPrefix) ShowByIDFiltering {
	return newShowByIDFiltering("In", "ServiceIn", fmt.Sprintf("In: In{%[1]v: id.%[1]vId()}", identifierKind))
}

type showByIDApplicationFilter struct {
	showByIDFilter
}

func (s *showByIDApplicationFilter) WithFiltering() string {
	return fmt.Sprintf("With%s(%s)", s.Name, s.Args)
}

func newShowByIDApplicationFiltering() ShowByIDFiltering {
	return &showByIDApplicationFilter{
		showByIDFilter: showByIDFilter{
			Name: "ApplicationName",
			Kind: "",
			Args: "id.DatabaseId()",
		},
	}
}

func (s *Operation) withFiltering(filtering ...ShowByIDFilteringKind) *Operation {
	for _, filteringKind := range filtering {
		switch filteringKind {
		case ShowByIDInFiltering:
			s.ShowByIDFiltering = append(s.ShowByIDFiltering, newShowByIDInFiltering(s.ObjectInterface.ObjectIdentifierPrefix()))
		case ShowByIDExtendedInFiltering:
			s.ShowByIDFiltering = append(s.ShowByIDFiltering, newShowByIDExtendedInFiltering(s.ObjectInterface.ObjectIdentifierPrefix()))
		case ShowByIDServiceInFiltering:
			s.ShowByIDFiltering = append(s.ShowByIDFiltering, newShowByIDServiceInFiltering(s.ObjectInterface.ObjectIdentifierPrefix()))
		case ShowByIDLikeFiltering:
			s.ShowByIDFiltering = append(s.ShowByIDFiltering, newShowByIDLikeFiltering())
		case ShowByIDApplicationNameFiltering:
			s.ShowByIDFiltering = append(s.ShowByIDFiltering, newShowByIDApplicationFiltering())
		default:
			log.Println("[DEBUG] No showByID filtering found for kind:", filteringKind)
		}
	}
	return s
}
