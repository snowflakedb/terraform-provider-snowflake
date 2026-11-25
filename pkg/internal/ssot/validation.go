package ssot

import (
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

// ValidationKind is how many of the fields of a Validation can be set together.
type ValidationKind string

const (
	// ValidationAtLeastOne requires at least one of the fields to be set.
	ValidationAtLeastOne ValidationKind = "at_least_one"
	// ValidationExactlyOne requires exactly one of the fields to be set.
	ValidationExactlyOne ValidationKind = "exactly_one"
	// ValidationConflicting allows at most one of the fields to be set.
	ValidationConflicting ValidationKind = "conflicting"
)

// Validation is a rule over the fields of one parent: the top-level attributes of an object (Object.Validations), or the
// Children of an object attribute (Attribute.Validations). Fields are matched by Name. Each generator applies the rule
// where it can express it, see the README.
type Validation struct {
	Kind ValidationKind
	// Fields are the fields the rule covers. Empty covers every candidate field of an object attribute.
	Fields []Attribute
	// NotInSdk leaves the rule out of the SDK, for a field whose empty value is still sent, such as an empty list that
	// clears it. The Terraform schema still gets the rule.
	NotInSdk bool
}

// SchemaOnly leaves the rule out of the SDK and keeps it in the Terraform schema.
func (v Validation) SchemaOnly() Validation {
	v.NotInSdk = true
	return v
}

// AtLeastOneOf requires at least one of the fields to be set.
func AtLeastOneOf(fields ...Attribute) Validation {
	return Validation{Kind: ValidationAtLeastOne, Fields: fields}
}

// ExactlyOneOf requires exactly one of the fields to be set.
func ExactlyOneOf(fields ...Attribute) Validation {
	return Validation{Kind: ValidationExactlyOne, Fields: fields}
}

// Conflicting allows at most one of the fields to be set.
func Conflicting(fields ...Attribute) Validation {
	return Validation{Kind: ValidationConflicting, Fields: fields}
}

// Names reports whether the validation names the field.
func (v Validation) Names(name string) bool {
	return slices.ContainsFunc(v.Fields, func(field Attribute) bool { return field.Name == name })
}

// ValidationFields are the top-level attributes of one Validation of the object, in Attributes order.
func (o Object) ValidationFields(v Validation) []Attribute {
	return collections.Filter(o.Attributes, func(attr Attribute) bool { return v.Names(attr.Name) })
}
