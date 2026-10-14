package ssot

import "fmt"

// TODO [next PRs]: replace panics with errors and handle them higher up;
// genhelpers does not yet have a uniform way to surface generator validation failures.
func ValidateObject(obj Object) {
	if obj.Name == "" {
		panic("ssot object has empty Name")
	}
	attrs := make(map[string]struct{})
	seen := make(map[string]struct{})
	for _, a := range obj.Attributes {
		if a.Name == "" {
			panic(fmt.Sprintf("object %q: attribute has empty Name", obj.Name))
		}
		if a.Readonly {
			panic(fmt.Sprintf("object %q: CREATE/ALTER attribute %q is marked Readonly", obj.Name, a.Name))
		}
		if _, dup := seen[a.Name]; dup {
			panic(fmt.Sprintf("object %q: duplicate attribute %q", obj.Name, a.Name))
		}
		seen[a.Name] = struct{}{}
		attrs[a.Name] = struct{}{}
	}
	validateColumnList(obj, "SHOW", obj.Show.Columns, attrs)
	validateColumnList(obj, "DESCRIBE", obj.Describe.Columns, attrs)
}

func validateColumnList(obj Object, kind string, list []*Attribute, attrs map[string]struct{}) {
	seenKey := make(map[string]struct{})
	seenName := make(map[string]struct{})
	for _, a := range list {
		if a == nil {
			panic(fmt.Sprintf("object %q: nil %s entry", obj.Name, kind))
		}
		if a.Name == "" {
			panic(fmt.Sprintf("object %q: %s entry has empty Name", obj.Name, kind))
		}
		if a.Readonly {
			if _, w := attrs[a.Name]; w {
				panic(fmt.Sprintf("object %q: %s Readonly %q duplicates a CREATE/ALTER attribute", obj.Name, kind, a.Name))
			}
		} else if _, w := attrs[a.Name]; !w {
			panic(fmt.Sprintf("object %q: %s entry %q is writable but not in Attributes", obj.Name, kind, a.Name))
		}
		key := a.OutputKey()
		if _, dup := seenKey[key]; dup {
			panic(fmt.Sprintf("object %q: duplicate %s output key %q", obj.Name, kind, key))
		}
		if _, dup := seenName[a.Name]; dup {
			panic(fmt.Sprintf("object %q: duplicate %s attribute %q", obj.Name, kind, a.Name))
		}
		seenKey[key] = struct{}{}
		seenName[a.Name] = struct{}{}
	}
}
