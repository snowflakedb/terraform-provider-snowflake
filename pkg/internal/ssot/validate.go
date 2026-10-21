package ssot

import (
	"fmt"
	"go/token"
)

// TODO [next PRs]: replace panics with errors and handle them higher up;
// genhelpers does not yet have a uniform way to surface generator validation failures.
func ValidateObject(obj Object) {
	if obj.Name == "" {
		panic("ssot object has empty Name")
	}
	attrs := make(map[string]struct{})
	seen := make(map[string]struct{})
	for _, a := range obj.Attributes {
		if a.IsTags() {
			panic(fmt.Sprintf("object %q: tags is not an attribute", obj.Name))
		}
		if a.IsParameters() {
			panic(fmt.Sprintf("object %q: parameters is not an attribute", obj.Name))
		}
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
		validateSqlClauses(obj, a)
	}
	validateCreateSpec(obj, attrs)
	validateSqlSequence(obj, "SET", obj.Set, attrs, func(attr Attribute) SqlClause { return attr.SetSql })
	validateSqlSequence(obj, "UNSET", obj.Unset, attrs, func(attr Attribute) SqlClause { return attr.UnsetSql })
	validateAlterSpec(obj)
	for _, a := range obj.Attributes {
		if a.CreateOnly && (a.Type != AttrBool || a.Required) {
			panic(fmt.Sprintf("object %q: CreateOnly attribute %q must be an optional bool", obj.Name, a.Name))
		}
		if a.NativeBool && a.Type != AttrBool {
			panic(fmt.Sprintf("object %q: NativeBool attribute %q must be a bool", obj.Name, a.Name))
		}
		for _, other := range a.ConflictsWith {
			if other.Name == a.Name {
				panic(fmt.Sprintf("object %q: attribute %q conflicts with itself", obj.Name, a.Name))
			}
			if _, ok := attrs[other.Name]; !ok {
				panic(fmt.Sprintf("object %q: attribute %q conflicts with unknown attribute %q", obj.Name, a.Name, other.Name))
			}
		}
	}
	validateColumnList(obj, "SHOW", obj.Show.Columns, attrs)
	validateColumnList(obj, "DESCRIBE", obj.Describe.Columns, attrs)
	validateSchemaStateUpgraders(obj)
}

func validateSchemaStateUpgraders(obj Object) {
	hooks := DefaultResourceExtHooks(obj.ResourceHooks...)
	for _, upgrade := range obj.SchemaStateUpgraders {
		name := string(upgrade)
		if name == "" {
			panic(fmt.Sprintf("object %q: SchemaStateUpgraders entry is empty", obj.Name))
		}
		if !token.IsIdentifier(name) || token.IsKeyword(name) {
			panic(fmt.Sprintf("object %q: SchemaStateUpgraders entry %q is not a function name", obj.Name, name))
		}
	}
	if len(obj.SchemaStateUpgraders) > 0 && !hooks.StateUpgraders.Generate {
		panic(fmt.Sprintf("object %q: SchemaStateUpgraders require WithStateUpgraders(Generated())", obj.Name))
	}
	if hooks.StateUpgraders.Generate && len(obj.SchemaStateUpgraders) == 0 {
		panic(fmt.Sprintf("object %q: WithStateUpgraders(Generated()) needs upgrade functions", obj.Name))
	}
}

func validateCreateSpec(obj Object, attrs map[string]struct{}) {
	seen := make(map[string]struct{})
	clause := func(attr Attribute) SqlClause { return attr.CreateSql }
	walkSqlSequence(obj, "KeywordsBefore", obj.Create.KeywordsBefore, attrs, clause, seen, false)
	walkSqlSequence(obj, "CREATE", obj.Create.Clauses, attrs, clause, seen, true)
	requirePresentClausesListed(obj, "CREATE", clause, seen)
}

func validateAlterSpec(obj Object) {
	seen := make(map[AlterClauseKind]struct{})
	for _, clause := range obj.Alter.Clauses {
		switch clause.Kind {
		case AlterClauseKeyword:
			if clause.SQL == "" {
				panic(fmt.Sprintf("object %q: ALTER keyword is empty", obj.Name))
			}
		case AlterClauseRenameTo, AlterClauseSwapWith, AlterClauseSet, AlterClauseUnset, AlterClauseSetTags, AlterClauseUnsetTags:
			if _, dup := seen[clause.Kind]; dup {
				panic(fmt.Sprintf("object %q: duplicate ALTER %s", obj.Name, clause.Kind))
			}
			seen[clause.Kind] = struct{}{}
		default:
			panic(fmt.Sprintf("object %q: ALTER clause has kind %q", obj.Name, clause.Kind))
		}
	}
}

func validateSqlSequence(obj Object, statement string, list []*Attribute, attrs map[string]struct{}, clause func(Attribute) SqlClause) {
	seen := make(map[string]struct{})
	walkSqlSequence(obj, statement, list, attrs, clause, seen, statement == "CREATE")
	requirePresentClausesListed(obj, statement, clause, seen)
}

func requirePresentClausesListed(obj Object, statement string, clause func(Attribute) SqlClause, seen map[string]struct{}) {
	for _, attr := range obj.Attributes {
		if !clause(attr).Present {
			continue
		}
		if _, ok := seen[attr.Name]; !ok {
			panic(fmt.Sprintf("object %q: attribute %q %s clause is not in the %s sequence", obj.Name, attr.Name, statement, statement))
		}
	}
}

func walkSqlSequence(obj Object, statement string, list []*Attribute, attrs map[string]struct{}, clause func(Attribute) SqlClause, seen map[string]struct{}, allowTags bool) {
	for _, field := range list {
		if field == nil {
			panic(fmt.Sprintf("object %q: nil %s entry", obj.Name, statement))
		}
		if field.IsTags() {
			if !allowTags {
				panic(fmt.Sprintf("object %q: %s list contains tags", obj.Name, statement))
			}
			continue
		}
		if field.IsParameters() {
			if statement != "CREATE" && statement != "SET" && statement != "UNSET" {
				panic(fmt.Sprintf("object %q: %s list contains parameters", obj.Name, statement))
			}
			if _, dup := seen["parameters"]; dup {
				panic(fmt.Sprintf("object %q: duplicate parameters in %s", obj.Name, statement))
			}
			seen["parameters"] = struct{}{}
			continue
		}
		name := field.Name
		if _, ok := attrs[name]; !ok {
			panic(fmt.Sprintf("object %q: %s field %q is not in Attributes", obj.Name, statement, name))
		}
		if _, dup := seen[name]; dup {
			panic(fmt.Sprintf("object %q: duplicate %s field %q", obj.Name, statement, name))
		}
		seen[name] = struct{}{}
		if !clause(*field).Present {
			panic(fmt.Sprintf("object %q: attribute %q has no %s clause", obj.Name, name, statement))
		}
	}
}

func validateSqlClauses(obj Object, attr Attribute) {
	validateSqlClause(obj, attr, "CREATE", attr.CreateSql)
	validateSqlClause(obj, attr, "SET", attr.SetSql)
	validateSqlClause(obj, attr, "UNSET", attr.UnsetSql)
}

func validateSqlClause(obj Object, attr Attribute, statement string, clause SqlClause) {
	if !clause.Present {
		return
	}
	switch clause.Kind {
	case SqlKindParameter:
		switch attr.Type {
		case AttrInt, AttrBool, AttrString:
		case AttrEnum:
			if attr.Enum == nil || attr.Enum.Name == "" {
				panic(fmt.Sprintf("object %q: attribute %q %s enum parameter has no Enum", obj.Name, attr.Name, statement))
			}
		case AttrList:
			if clause.ItemStruct == "" {
				panic(fmt.Sprintf("object %q: attribute %q %s list parameter has no item struct", obj.Name, attr.Name, statement))
			}
		default:
			panic(fmt.Sprintf("object %q: attribute %q %s parameter has type %q", obj.Name, attr.Name, statement, attr.Type))
		}
	case SqlKindKeyword:
		// A keyword marks presence. Unset keeps the attribute's value type.
	case SqlKindIdentifier:
		if attr.Type != AttrIdentifier {
			panic(fmt.Sprintf("object %q: attribute %q %s identifier has type %q", obj.Name, attr.Name, statement, attr.Type))
		}
		if clause.IdType == "" {
			panic(fmt.Sprintf("object %q: attribute %q %s identifier has no IdType", obj.Name, attr.Name, statement))
		}
	default:
		panic(fmt.Sprintf("object %q: attribute %q %s has sql kind %q", obj.Name, attr.Name, statement, clause.Kind))
	}
}

func validateColumnList(obj Object, kind string, list []*Attribute, attrs map[string]struct{}) {
	seenKey := make(map[string]struct{})
	seenName := make(map[string]struct{})
	for _, a := range list {
		if a == nil {
			panic(fmt.Sprintf("object %q: nil %s entry", obj.Name, kind))
		}
		if a.IsTags() {
			panic(fmt.Sprintf("object %q: tags is not a %s column", obj.Name, kind))
		}
		if a.IsParameters() {
			panic(fmt.Sprintf("object %q: parameters is not a %s column", obj.Name, kind))
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
