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
		validateAttributeOptions(obj.Name, a, false)
	}
	validateCreateSpec(obj, attrs)
	validateSqlSequence(obj, "SET", obj.Set, attrs, func(attr Attribute) SqlClause { return attr.SetSql })
	validateSqlSequence(obj, "UNSET", obj.Unset, attrs, func(attr Attribute) SqlClause { return attr.UnsetSql })
	validateAlterSpec(obj)
	validateCustomOperations(obj, attrs)
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
	walkSqlSequence(obj, "KeywordsBefore", obj.Create.KeywordsBefore, attrs, seen, false)
	walkSqlSequence(obj, "CREATE", obj.Create.Clauses, attrs, seen, true)
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

func validateCustomOperations(obj Object, attrs map[string]struct{}) {
	seenKind := make(map[string]struct{})
	for _, spec := range obj.CustomOperations {
		if spec.Kind == "" {
			panic(fmt.Sprintf("object %q: custom operation kind is empty", obj.Name))
		}
		if _, dup := seenKind[spec.Kind]; dup {
			panic(fmt.Sprintf("object %q: duplicate custom operation %q", obj.Name, spec.Kind))
		}
		seenKind[spec.Kind] = struct{}{}
		validateCustomOperationClauses(obj, spec, attrs)
	}
}

func validateCustomOperationClauses(obj Object, spec CustomOperationSpec, attrs map[string]struct{}) {
	seen := make(map[CustomOperationClauseKind]struct{})
	seenAttr := make(map[string]struct{})
	seenPredefined := make(map[string]struct{})
	singleton := func(kind CustomOperationClauseKind) {
		if _, dup := seen[kind]; dup {
			panic(fmt.Sprintf("object %q: custom operation %q has duplicate %s", obj.Name, spec.Kind, kind))
		}
		seen[kind] = struct{}{}
	}
	for _, clause := range spec.Clauses {
		switch clause.Kind {
		case CustomOperationClauseStaticSQL:
			if clause.SQL == "" {
				panic(fmt.Sprintf("object %q: custom operation %q static SQL is empty", obj.Name, spec.Kind))
			}
		case CustomOperationClauseOrReplace, CustomOperationClauseIfNotExists, CustomOperationClauseObjectWords, CustomOperationClauseName:
			singleton(clause.Kind)
		case CustomOperationClauseAttribute:
			if clause.Attribute == nil {
				panic(fmt.Sprintf("object %q: custom operation %q attribute is empty", obj.Name, spec.Kind))
			}
			if clause.Attribute.IsTags() {
				panic(fmt.Sprintf("object %q: custom operation %q attribute is tags", obj.Name, spec.Kind))
			}
			if clause.Attribute.IsParameters() {
				panic(fmt.Sprintf("object %q: custom operation %q attribute is parameters", obj.Name, spec.Kind))
			}
			name := clause.Attribute.Name
			if _, ok := attrs[name]; !ok {
				panic(fmt.Sprintf("object %q: custom operation %q attribute %q is not in Attributes", obj.Name, spec.Kind, name))
			}
			if _, dup := seenAttr[name]; dup {
				panic(fmt.Sprintf("object %q: custom operation %q has duplicate attribute %q", obj.Name, spec.Kind, name))
			}
			seenAttr[name] = struct{}{}
			validateSqlClause(obj, *clause.Attribute, "CREATE", clause.Attribute.EffectiveSqlClause(ClauseCreate))
		case CustomOperationClausePredefined:
			if clause.Predefined == "" {
				panic(fmt.Sprintf("object %q: custom operation %q predefined type is empty", obj.Name, spec.Kind))
			}
			if _, dup := seenPredefined[clause.Predefined]; dup {
				panic(fmt.Sprintf("object %q: custom operation %q has duplicate predefined %q", obj.Name, spec.Kind, clause.Predefined))
			}
			seenPredefined[clause.Predefined] = struct{}{}
		default:
			panic(fmt.Sprintf("object %q: custom operation %q clause has kind %q", obj.Name, spec.Kind, clause.Kind))
		}
	}
}

func validateSqlSequence(obj Object, statement string, list []*Attribute, attrs map[string]struct{}, clause func(Attribute) SqlClause) {
	seen := make(map[string]struct{})
	walkSqlSequence(obj, statement, list, attrs, seen, statement == "CREATE")
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

func walkSqlSequence(obj Object, statement string, list []*Attribute, attrs map[string]struct{}, seen map[string]struct{}, allowTags bool) {
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
		validateSqlClause(obj, *field, statement, field.EffectiveSqlClause(clauseUse(statement)))
	}
}

func clauseUse(statement string) ClauseUse {
	switch statement {
	case "CREATE", "KeywordsBefore":
		return ClauseCreate
	case "SET":
		return ClauseSet
	case "UNSET":
		return ClauseUnset
	default:
		panic(fmt.Sprintf("unknown sql statement %q", statement))
	}
}

func validateSqlClauses(obj Object, attr Attribute) {
	validateExplicitSqlClause(obj, attr, "CREATE", attr.CreateSql)
	validateExplicitSqlClause(obj, attr, "SET", attr.SetSql)
	validateExplicitSqlClause(obj, attr, "UNSET", attr.UnsetSql)
}

func validateExplicitSqlClause(obj Object, attr Attribute, statement string, clause SqlClause) {
	if !clause.Present || clause.fromType {
		return
	}
	validateSqlClause(obj, attr, statement, clause)
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

func validateAttributeOptions(objectName string, attr Attribute, asElem bool) {
	if attr.DiffSuppress != nil {
		for _, suppress := range attr.DiffSuppress.Suppressors() {
			switch suppress.Kind {
			case SuppressStatement:
				if attr.Type != AttrString || suppress.Key != "" {
					panic(fmt.Sprintf("object %q: invalid statement suppression for attribute %q", objectName, attr.Name))
				}
			case SuppressShow, SuppressDescribe:
				if asElem {
					panic(fmt.Sprintf("object %q: output suppression on collection elem %q", objectName, attr.Name))
				}
			default:
				panic(fmt.Sprintf("object %q: unknown suppression %q for attribute %q", objectName, suppress.Kind, attr.Name))
			}
		}
	}
	if attr.Elem != nil {
		validateAttributeOptions(objectName, *attr.Elem, true)
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
