package gen

import "strings"

// OperationTestKind classifies an operation for the purpose of deriving unit_tests sql cases.
type OperationTestKind string

const (
	OperationTestKindCreate   OperationTestKind = "create"
	OperationTestKindAlter    OperationTestKind = "alter"
	OperationTestKindShow     OperationTestKind = "show"
	OperationTestKindDrop     OperationTestKind = "drop"
	OperationTestKindDescribe OperationTestKind = "describe"
	OperationTestKindOther    OperationTestKind = "other"
)

// TestKind classifies this operation for unit_tests sql case generation.
// ShowByID (and any other operation with no OptsField, i.e. no SQL to build) is always OperationTestKindOther.
func (s *Operation) TestKind() OperationTestKind {
	if s.OptsField == nil {
		return OperationTestKindOther
	}
	switch {
	case strings.HasPrefix(s.Name, string(OperationKindCreate)):
		return OperationTestKindCreate
	case strings.HasPrefix(s.Name, string(OperationKindAlter)):
		return OperationTestKindAlter
	case strings.HasPrefix(s.Name, string(OperationKindShow)):
		return OperationTestKindShow
	case strings.HasPrefix(s.Name, string(OperationKindDrop)):
		return OperationTestKindDrop
	case strings.HasPrefix(s.Name, string(OperationKindDescribe)):
		return OperationTestKindDescribe
	default:
		return OperationTestKindOther
	}
}

// MutuallyExclusiveOptions returns the child fields listed in the root OptsField's ExactlyOneValueSet or MoreThanOneValueSet validation, if any.
// E.g. for Alter<object>Options this can be [RenameTo, Set, Unset, SetSecure, UnsetSecure, SetTags, UnsetTags].
// Returns nil when the root has no such validation.
func (s *Operation) MutuallyExclusiveOptions() []*Field {
	if s.OptsField == nil {
		return nil
	}
	for _, v := range s.OptsField.Validations {
		if v.Type != ExactlyOneValueSet && v.Type != MoreThanOneValueSet {
			continue
		}
		branches := make([]*Field, 0, len(v.FieldNames))
		for _, name := range v.FieldNames {
			if child := s.OptsField.FindChild(name); child != nil {
				branches = append(branches, child)
			}
		}
		if len(branches) > 0 {
			return branches
		}
	}
	return nil
}

// ShowOptionalFields returns the root-level optional, non-flag fields of a Show operation (e.g. Like, In, StartsWith, Limit).
// Boolean flags (*bool, e.g. a hypothetical Terse) are excluded.
func (s *Operation) ShowOptionalFields() []*Field {
	if s.OptsField == nil {
		return nil
	}
	result := make([]*Field, 0)
	for idx := range s.OptsField.Fields {
		f := &s.OptsField.Fields[idx]
		if f.IsPointer() && f.KindNoPtr() != "bool" {
			result = append(result, f)
		}
	}
	return result
}
