package genhelpers

import "fmt"

type Mapper func(string) string

var (
	Identity           = func(field string) string { return field }
	ToString           = func(field string) string { return fmt.Sprintf("%s.String()", field) }
	FullyQualifiedName = func(field string) string { return fmt.Sprintf("%s.FullyQualifiedName()", field) }
	Name               = func(field string) string { return fmt.Sprintf("%s.Name()", field) }
	CastToString       = func(field string) string { return fmt.Sprintf("string(%s)", field) }
	CastToInt          = func(field string) string { return fmt.Sprintf("int(%s)", field) }
	Dereference        = func(field string) string { return fmt.Sprintf("*%s", field) }
	Parentheses        = func(field string) string { return fmt.Sprintf("(%s)", field) }
)
