package gen

import "fmt"

// ShowByIDFromIdentifierConfig makes the generated ShowByID copy identifier fields onto the SHOW row it found.
// It is for a SHOW that does not return those columns, such as SHOW DATABASE ROLES, which has no database column.
type ShowByIDFromIdentifierConfig struct {
	// Fields are plain-only SHOW fields copied from id, in assignment order.
	// Each field name is also the identifier method: DatabaseName calls id.DatabaseName().
	Fields []string
}

// WithShowByIDFromIdentifier makes the generated ShowByID copy each field from the identifier after FindFirst.
func (i *Interface) WithShowByIDFromIdentifier(fields []string) *Interface {
	if len(fields) == 0 {
		panic(fmt.Sprintf("interface %q: ShowByID from identifier needs a field", i.Name))
	}
	if !i.hasShowByID() {
		panic(fmt.Sprintf("interface %q: ShowByID from identifier needs a generated ShowByID", i.Name))
	}
	i.ShowByIDFromIdentifier = &ShowByIDFromIdentifierConfig{Fields: fields}
	return i
}
