package ssot

import (
	"fmt"
	"go/token"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

// TODO [next PRs]: replace panics with errors and handle them higher up;
// genhelpers does not yet have a uniform way to surface generator validation failures.
func ValidateObject(obj Object) {
	if obj.Name == "" {
		panic("ssot object has empty Name")
	}
	attrs := make(map[string]struct{})
	sqlOnly := make(map[string]struct{})
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
		if a.SqlOnly && (!a.CreateSql.Present || a.CreateSql.Kind != SqlKindStatic || a.SetSql.Present || a.UnsetSql.Present) {
			panic(fmt.Sprintf("object %q: SQL-only attribute %q needs a static CREATE clause", obj.Name, a.Name))
		}
		if _, dup := seen[a.Name]; dup {
			panic(fmt.Sprintf("object %q: duplicate attribute %q", obj.Name, a.Name))
		}
		seen[a.Name] = struct{}{}
		attrs[a.Name] = struct{}{}
		if a.SqlOnly {
			sqlOnly[a.Name] = struct{}{}
		}
		validateCollectionBounds(obj, a)
		validateNestedAttribute(obj, a)
		validateAttribute(obj, a)
		validateSqlClauses(obj, a)
		validateAttributeOptions(obj.Name, a, false)
	}
	validateCreateSpec(obj, attrs)
	switch {
	case obj.VariantOf != nil:
		if obj.Create.Sql != "" {
			panic(fmt.Sprintf("object %q: variant CREATE words are VariantEntityCreateName", obj.Name))
		}
	case obj.VariantEntityCreateName != "":
		panic(fmt.Sprintf("object %q: VariantEntityCreateName is for a variant", obj.Name))
	case obj.VariantEntityAlterName != "":
		panic(fmt.Sprintf("object %q: VariantEntityAlterName is for a variant", obj.Name))
	}
	validateSqlSequence(obj, "SET", obj.Set, attrs, func(attr Attribute) SqlClause { return attr.SetSql })
	validateSqlSequence(obj, "UNSET", obj.Unset, attrs, func(attr Attribute) SqlClause { return attr.UnsetSql })
	validateAlterSpec(obj)
	validateCustomOperations(obj, attrs)
	validateCustomInterfaceMethods(obj)
	for _, a := range obj.Attributes {
		if a.CreateOnly && (a.Type != AttrBool || a.Required) {
			panic(fmt.Sprintf("object %q: CreateOnly attribute %q must be an optional bool", obj.Name, a.Name))
		}
		if a.NativeBool && a.Type != AttrBool {
			panic(fmt.Sprintf("object %q: NativeBool attribute %q must be a bool", obj.Name, a.Name))
		}
		validateSkipCreateBlock(obj, a)
		for _, other := range a.ConflictsWith {
			if other.Name == a.Name {
				panic(fmt.Sprintf("object %q: attribute %q conflicts with itself", obj.Name, a.Name))
			}
			if _, ok := attrs[other.Name]; !ok {
				panic(fmt.Sprintf("object %q: attribute %q conflicts with unknown attribute %q", obj.Name, a.Name, other.Name))
			}
		}
	}
	nested := nestedAttributeNames(obj)
	validateColumnList(obj, "SHOW", obj.Show.Columns, attrs, sqlOnly, nested)
	validateColumnList(obj, "SHOW generic", obj.ShowGeneric.Columns, attrs, sqlOnly, nested)
	validateColumnList(obj, "DESCRIBE", obj.Describe.Columns, attrs, sqlOnly, nested)
	validateColumnList(obj, "DESCRIBE rows", obj.DescribeRows.Columns, attrs, sqlOnly, nested)
	validateDescribeRows(obj)
	validateDescribeDetails(obj)
	g := ObjectGenerators(obj.Generators...)
	validateGeneratorSkips(obj, "object-assertion", "SHOW", g.ObjectAssert.Show, obj.PairedShowColumns(), g.ObjectAssert.ShowSkipFields)
	validateGeneratorSkips(obj, "object-assertion", "DESCRIBE", g.ObjectAssert.Describe, obj.Describe.Columns, g.ObjectAssert.DescribeSkipFields)
	validateGeneratorSkips(obj, "object-assertion nested", "SHOW", g.ObjectAssert.Show, obj.Show.Columns, nestedAttributes(g.ObjectAssert.ShowNestedFields))
	validateGeneratorSkips(obj, "object-assertion nested", "DESCRIBE", g.ObjectAssert.Describe, obj.Describe.Columns, nestedAttributes(g.ObjectAssert.DescribeNestedFields))
	validateGeneratorSkips(obj, "show-output", "SHOW", g.ShowOutputAssert.Show, obj.PairedShowColumns(), g.ShowOutputAssert.ShowSkipFields)
	validateGeneratorSkips(obj, "show-output", "DESCRIBE", g.ShowOutputAssert.Describe, obj.Describe.Columns, g.ShowOutputAssert.DescribeSkipFields)
	obj.ShowFromBase()
	obj.ParameterAssertionNames()
	validateSchemaStateUpgraders(obj)
	validateCustomizeDiff(obj)
}

func validateAttribute(obj Object, a Attribute) {
	if a.Type == AttrIdentifier && a.IdentifierType == "" {
		name := a.Name
		if name == "" {
			name = "<elem>"
		}
		panic(fmt.Sprintf("object %q: identifier attribute %q has no IdentifierType", obj.Name, name))
	}
	if a.Type == AttrEnum {
		enum := a.ColumnEnum()
		if enum == nil || enum.Name == "" {
			name := a.Name
			if name == "" {
				name = "<elem>"
			}
			panic(fmt.Sprintf("object %q: enum attribute %q has no Enum", obj.Name, name))
		}
	}
	if a.Elem != nil {
		validateAttribute(obj, *a.Elem)
	}
	for _, child := range a.Children {
		validateAttribute(obj, child)
	}
}

func validateCustomizeDiff(obj Object) {
	hooks := DefaultResourceExtHooks(obj.ResourceHooks...)
	for _, diff := range obj.AdditionalCustomizeDiffs {
		name := string(diff)
		if name == "" {
			panic(fmt.Sprintf("object %q: AdditionalCustomizeDiffs entry is empty", obj.Name))
		}
		if !token.IsIdentifier(name) || token.IsKeyword(name) {
			panic(fmt.Sprintf("object %q: AdditionalCustomizeDiffs entry %q is not a function name", obj.Name, name))
		}
	}
	if len(obj.AdditionalCustomizeDiffs) > 0 && !hooks.CustomizeDiff.Generate {
		panic(fmt.Sprintf("object %q: AdditionalCustomizeDiffs require WithCustomizeDiff(Generated())", obj.Name))
	}
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

func validateNestedAttribute(obj Object, attr Attribute) {
	fields := attr.NestedFields()
	if len(fields) == 0 {
		return
	}
	if attr.Type != AttrList && attr.Type != AttrSet {
		panic(fmt.Sprintf("object %q: attribute %q Nested requires a list or set", obj.Name, attr.Name))
	}
	if attr.Elem != nil && !attr.IsObjectList() {
		panic(fmt.Sprintf("object %q: attribute %q has both Elem and Nested", obj.Name, attr.Name))
	}
	if attr.Elem != nil && len(attr.Children) > 0 {
		panic(fmt.Sprintf("object %q: attribute %q has both Children and an object Elem", obj.Name, attr.Name))
	}
	// The fields of an object element carry their own CREATE clause when the element has an object clause.
	sqlFields := attr.IsObjectList() && attr.Elem.CreateSql.Kind == SqlKindObject
	seen := make(map[string]struct{}, len(fields))
	for _, child := range fields {
		if child.Name == "" {
			panic(fmt.Sprintf("object %q: attribute %q has a nested attribute with empty Name", obj.Name, attr.Name))
		}
		if _, dup := seen[child.Name]; dup {
			panic(fmt.Sprintf("object %q: attribute %q has duplicate nested attribute %q", obj.Name, attr.Name, child.Name))
		}
		seen[child.Name] = struct{}{}
		if statement, ok := sqlClauseStatement(child); ok {
			if !sqlFields || statement != "CREATE" || child.SetSql.Present || child.UnsetSql.Present {
				panic(fmt.Sprintf("object %q: nested attribute %q.%q has a %s clause", obj.Name, attr.Name, child.Name, statement))
			}
		}
		validateCollectionBounds(obj, child)
		validateSkipCreateBlock(obj, child)
		validateNestedAttribute(obj, child)
	}
	if attr.IsObjectList() {
		validateObjectVariants(obj, attr)
	}
}

func validateSkipCreateBlock(obj Object, attr Attribute) {
	if !attr.SkipCreateBlock {
		return
	}
	if attr.Type != AttrList || len(attr.Children) == 0 {
		panic(fmt.Sprintf("object %q: SkipCreateBlock attribute %q must be a list block", obj.Name, attr.Name))
	}
}

// validateObjectVariants checks the variants of the Elem object of a list: the attributes each one refers to.
func validateObjectVariants(obj Object, list Attribute) {
	clause := list.Elem.CreateSql
	if len(clause.Variants) == 0 {
		return
	}
	if clause.Kind != SqlKindObject {
		panic(fmt.Sprintf("object %q: attribute %q has variants but its Elem has no object clause", obj.Name, list.Name))
	}
	keys := make(map[string]struct{}, len(clause.Variants))
	for _, variant := range clause.Variants {
		if _, dup := keys[variant.Key]; dup || variant.Key == "" {
			panic(fmt.Sprintf("object %q: attribute %q has an empty or duplicate variant %q", obj.Name, list.Name, variant.Key))
		}
		keys[variant.Key] = struct{}{}
		if len(variant.Fields) == 0 {
			panic(fmt.Sprintf("object %q: attribute %q variant %q has no fields", obj.Name, list.Name, variant.Key))
		}
		fields := make(map[string]struct{}, len(variant.Fields))
		for _, field := range variant.Fields {
			if field == nil {
				panic(fmt.Sprintf("object %q: attribute %q variant %q has a nil field", obj.Name, list.Name, variant.Key))
			}
			if _, dup := fields[field.Name]; dup {
				panic(fmt.Sprintf("object %q: attribute %q variant %q has duplicate field %q", obj.Name, list.Name, variant.Key, field.Name))
			}
			fields[field.Name] = struct{}{}
			validateVariantField(obj, list, variant, *field)
		}
	}
}

// validateVariantField checks one attribute a variant refers to: a child of the element, or a SqlOnly static or nested struct.
func validateVariantField(obj Object, list Attribute, variant SqlVariant, field Attribute) {
	if !field.CreateSql.Present {
		panic(fmt.Sprintf("object %q: attribute %q variant %q field %q has no CREATE clause", obj.Name, list.Name, variant.Key, field.Name))
	}
	if field.SqlOnly {
		switch field.CreateSql.Kind {
		case SqlKindStatic:
		case SqlKindObject:
			if field.Type != AttrObject || len(field.Children) == 0 {
				panic(fmt.Sprintf("object %q: attribute %q variant %q field %q needs Children", obj.Name, list.Name, variant.Key, field.Name))
			}
			for _, child := range field.Children {
				if !child.CreateSql.Present {
					panic(fmt.Sprintf("object %q: attribute %q variant %q struct %q field %q has no CREATE clause", obj.Name, list.Name, variant.Key, field.Name, child.Name))
				}
			}
		default:
			panic(fmt.Sprintf("object %q: attribute %q variant %q field %q is SQL-only but not static or an object", obj.Name, list.Name, variant.Key, field.Name))
		}
		return
	}
	if !slices.ContainsFunc(list.Elem.Children, func(c Attribute) bool { return c.Name == field.Name }) {
		panic(fmt.Sprintf("object %q: attribute %q variant %q field %q is not a child of the element", obj.Name, list.Name, variant.Key, field.Name))
	}
}

func validateCollectionBounds(obj Object, attr Attribute) {
	if attr.MinItems == nil && attr.MaxItems == nil {
		return
	}
	if attr.Type != AttrList && attr.Type != AttrSet {
		panic(fmt.Sprintf("object %q: attribute %q item bounds require a list or set", obj.Name, attr.Name))
	}
	if attr.MinItems != nil && attr.MaxItems != nil && *attr.MinItems > *attr.MaxItems {
		panic(fmt.Sprintf("object %q: attribute %q MinItems %d is greater than MaxItems %d", obj.Name, attr.Name, *attr.MinItems, *attr.MaxItems))
	}
}

func nestedAttributeNames(obj Object) map[string]struct{} {
	names := make(map[string]struct{})
	for _, attr := range obj.Attributes {
		if len(attr.NestedFields()) > 0 {
			names[attr.Name] = struct{}{}
		}
	}
	return names
}

func sqlClauseStatement(attr Attribute) (string, bool) {
	switch {
	case attr.CreateSql.Present:
		return "CREATE", true
	case attr.SetSql.Present:
		return "SET", true
	case attr.UnsetSql.Present:
		return "UNSET", true
	default:
		return "", false
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
	seenList := make(map[string]struct{})
	for _, clause := range obj.Alter.Clauses {
		switch clause.Kind {
		case AlterClauseKeyword:
			if clause.SQL == "" {
				panic(fmt.Sprintf("object %q: ALTER keyword is empty", obj.Name))
			}
		case AlterClauseList:
			if clause.NotExactlyOne {
				panic(fmt.Sprintf("object %q: ALTER list %q is an exactly-one branch", obj.Name, clause.SQL))
			}
			if clause.SQL == "" {
				panic(fmt.Sprintf("object %q: ALTER list is empty", obj.Name))
			}
			if clause.ListItem == "" {
				panic(fmt.Sprintf("object %q: ALTER list %q has no item type", obj.Name, clause.SQL))
			}
			if _, dup := seenList[clause.SQL]; dup {
				panic(fmt.Sprintf("object %q: duplicate ALTER list %q", obj.Name, clause.SQL))
			}
			seenList[clause.SQL] = struct{}{}
		case AlterClauseText:
			if clause.SQL == "" {
				panic(fmt.Sprintf("object %q: ALTER text is empty", obj.Name))
			}
			if _, dup := seenList[clause.SQL]; dup {
				panic(fmt.Sprintf("object %q: duplicate ALTER clause %q", obj.Name, clause.SQL))
			}
			seenList[clause.SQL] = struct{}{}
		case AlterClauseBlock:
			if clause.Field == "" || clause.SQL == "" || clause.Block == nil {
				panic(fmt.Sprintf("object %q: ALTER block needs a field, SQL, and an attribute", obj.Name))
			}
			if !clause.Block.IsObjectList() && (clause.Block.Type != AttrObject || len(clause.Block.Children) == 0) {
				panic(fmt.Sprintf("object %q: ALTER block %q needs a list of objects or an object with Children", obj.Name, clause.SQL))
			}
			if _, dup := seenList[clause.SQL]; dup {
				panic(fmt.Sprintf("object %q: duplicate ALTER clause %q", obj.Name, clause.SQL))
			}
			seenList[clause.SQL] = struct{}{}
		case AlterClauseRenameTo, AlterClauseSwapWith, AlterClauseSet, AlterClauseUnset, AlterClauseSetTags, AlterClauseUnsetTags:
			if clause.NotExactlyOne {
				panic(fmt.Sprintf("object %q: ALTER %s cannot leave exactly-one", obj.Name, clause.Kind))
			}
			if _, dup := seen[clause.Kind]; dup {
				panic(fmt.Sprintf("object %q: duplicate ALTER %s", obj.Name, clause.Kind))
			}
			seen[clause.Kind] = struct{}{}
		default:
			panic(fmt.Sprintf("object %q: ALTER clause has kind %q", obj.Name, clause.Kind))
		}
	}
}

// TODO [next PRs]: CustomOperationsFromObject and CustomInterfaceMethodsFromObject repeat these
// checks. make generate-sdk does not call ValidateObject, and the generator tests call the emitters
// on a fake object without it. Share one checker and delete the copies. See that TODO.

func validateCustomInterfaceMethods(obj Object) {
	seen := make(map[string]struct{})
	for _, spec := range obj.CustomInterfaceMethods {
		if spec.Name == "" {
			panic(fmt.Sprintf("object %q: custom interface method name is empty", obj.Name))
		}
		if _, dup := seen[spec.Name]; dup {
			panic(fmt.Sprintf("object %q: duplicate custom interface method %q", obj.Name, spec.Name))
		}
		seen[spec.Name] = struct{}{}
		if spec.Parameters == nil && obj.IdType == "" {
			panic(fmt.Sprintf("object %q: custom interface method %q needs IdType", obj.Name, spec.Name))
		}
		for _, param := range spec.Parameters {
			if param.Name == "" {
				panic(fmt.Sprintf("object %q: custom interface method %q parameter name is empty", obj.Name, spec.Name))
			}
			if param.Kind == "" && obj.IdType == "" {
				panic(fmt.Sprintf("object %q: custom interface method %q parameter %q needs a type", obj.Name, spec.Name, param.Name))
			}
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
		if _, dup := seen[name]; dup {
			panic(fmt.Sprintf("object %q: duplicate %s field %q", obj.Name, statement, name))
		}
		seen[name] = struct{}{}
		if _, ok := attrs[name]; !ok {
			validateAttribute(obj, *field)
			validateSqlClauses(obj, *field)
		}
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
		case AttrIdentifier:
			// An identifier with a parameter clause is quoted text.
		case AttrEnum:
			if enum := attr.ColumnEnum(); enum == nil || enum.Name == "" {
				panic(fmt.Sprintf("object %q: attribute %q %s enum parameter has no Enum", obj.Name, attr.Name, statement))
			}
		case AttrList, AttrSet:
			if !attr.IsObjectList() && clause.ItemStruct == "" && (attr.Elem == nil || attr.Elem.Type != AttrIdentifier || attr.Elem.IdentifierType == "") {
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
	case SqlKindObject:
		// The shape of an element or of a nested struct. Its fields and variants are checked with the list.
	case SqlKindStatic:
		if statement != "CREATE" {
			panic(fmt.Sprintf("object %q: attribute %q static clause is CREATE only", obj.Name, attr.Name))
		}
		if !attr.SqlOnly {
			panic(fmt.Sprintf("object %q: attribute %q static clause is SQL-only", obj.Name, attr.Name))
		}
		if clause.Field == "" || clause.SQL == "" {
			panic(fmt.Sprintf("object %q: attribute %q static clause needs a field name and SQL", obj.Name, attr.Name))
		}
		if !token.IsIdentifier(clause.Field) || token.IsExported(clause.Field) {
			panic(fmt.Sprintf("object %q: attribute %q static field %q is not private", obj.Name, attr.Name, clause.Field))
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

func validateColumnList(obj Object, kind string, list []*Attribute, attrs map[string]struct{}, sqlOnly map[string]struct{}, nested map[string]struct{}) {
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
		if _, ok := sqlOnly[a.Name]; ok {
			panic(fmt.Sprintf("object %q: %s entry %q is SQL-only", obj.Name, kind, a.Name))
		}
		// A variant SHOW struct reuses the parent's show attributes, not this object's resource fields.
		if obj.VariantOf == nil {
			if a.Readonly {
				if _, w := attrs[a.Name]; w {
					panic(fmt.Sprintf("object %q: %s Readonly %q duplicates a CREATE/ALTER attribute", obj.Name, kind, a.Name))
				}
			} else if _, w := attrs[a.Name]; !w {
				panic(fmt.Sprintf("object %q: %s entry %q is writable but not in Attributes", obj.Name, kind, a.Name))
			}
		}
		key := a.OutputKey()
		if _, dup := seenKey[key]; dup {
			panic(fmt.Sprintf("object %q: duplicate %s output key %q", obj.Name, kind, key))
		}
		if _, dup := seenName[a.Name]; dup {
			panic(fmt.Sprintf("object %q: duplicate %s attribute %q", obj.Name, kind, a.Name))
		}
		// TODO (next PRs): Support nested attributes.
		if _, nestedAttr := nested[a.Name]; nestedAttr || len(a.NestedFields()) > 0 {
			panic(fmt.Sprintf("object %q: nested attribute %q is not a %s column", obj.Name, a.Name, kind))
		}
		seenKey[key] = struct{}{}
		seenName[a.Name] = struct{}{}
		validateAttribute(obj, *a)
	}
}

// validateDescribeRows rejects a DescribeRows row that nothing reads or that cannot be a raw row.
func validateDescribeRows(obj Object) {
	if len(obj.DescribeRows.Columns) == 0 {
		return
	}
	if !ObjectGenerators(obj.Generators...).Sdk.Describe {
		panic(fmt.Sprintf("object %q: DescribeRows needs SdkDescribe", obj.Name))
	}
	for _, column := range obj.DescribeRows.Columns {
		if !column.Readonly {
			panic(fmt.Sprintf("object %q: DESCRIBE rows entry %q must be Readonly", obj.Name, column.Name))
		}
	}
}

func validateDescribeDetails(obj Object) {
	sdk := ObjectGenerators(obj.Generators...).Sdk
	if len(sdk.DescribeDetailsParse) > 0 {
		validateDescribeDetailsParse(obj, sdk)
	}
	if sdk.DescribeDetailsMethod && !sdk.DescribeDetails {
		panic(fmt.Sprintf("object %q: SdkDescribeDetailsMethod needs SdkDescribeDetails", obj.Name))
	}
	if !sdk.DescribeDetails {
		validateGeneratorSkips(obj, "sdk-describe-details nested", "DESCRIBE", false, obj.Describe.Columns, sdk.DescribeDetailsNestedFields)
		return
	}
	if !sdk.Describe || len(obj.DescribeRows.Columns) == 0 {
		panic(fmt.Sprintf("object %q: SdkDescribeDetails needs SdkDescribe and DescribeRows", obj.Name))
	}
	validateGeneratorSkips(obj, "sdk-describe-details nested", "DESCRIBE", true, obj.Describe.Columns, sdk.DescribeDetailsNestedFields)
}

// validateDescribeDetailsParse checks that the mappings cover every Describe column except id exactly once.
func validateDescribeDetailsParse(obj Object, sdk SdkGenerator) {
	if !sdk.DescribeDetails {
		panic(fmt.Sprintf("object %q: SdkDescribeDetailsParse needs SdkDescribeDetails", obj.Name))
	}
	mapped := make(map[*Attribute]struct{}, len(sdk.DescribeDetailsParse))
	for _, m := range sdk.DescribeDetailsParse {
		if m.Attr == nil || m.Property == "" {
			panic(fmt.Sprintf("object %q: a DESCRIBE details mapping needs an attribute and a property", obj.Name))
		}
		if !slices.Contains(obj.Describe.Columns, m.Attr) {
			panic(fmt.Sprintf("object %q: DESCRIBE details mapping %q is not a DESCRIBE column", obj.Name, m.Attr.Name))
		}
		if _, dup := mapped[m.Attr]; dup {
			panic(fmt.Sprintf("object %q: DESCRIBE details mapping %q is set twice", obj.Name, m.Attr.Name))
		}
		mapped[m.Attr] = struct{}{}
		isNested := slices.Contains(sdk.DescribeDetailsNestedFields, m.Attr)
		switch {
		case m.Prefix && !isNested:
			panic(fmt.Sprintf("object %q: DESCRIBE details prefix mapping %q must be a nested field", obj.Name, m.Attr.Name))
		case !m.Prefix && (isNested || m.Attr.Type != AttrString || m.Attr.OptionalInShow):
			panic(fmt.Sprintf("object %q: DESCRIBE details property mapping %q must be a required string column", obj.Name, m.Attr.Name))
		}
	}
	for _, column := range obj.Describe.Columns {
		if _, ok := mapped[column]; !ok && column.OutputKey() != "id" {
			panic(fmt.Sprintf("object %q: DESCRIBE column %q has no details mapping", obj.Name, column.Name))
		}
	}
}

func nestedAttributes(fields []ObjectAssertNestedField) []*Attribute {
	return collections.Map(fields, func(f ObjectAssertNestedField) *Attribute { return f.Attr })
}

func validateGeneratorSkips(obj Object, generator, kind string, enabled bool, columns []*Attribute, skipped []*Attribute) {
	if len(skipped) == 0 {
		return
	}
	if !enabled {
		panic(fmt.Sprintf("object %q: %s %s skip fields are set while that side is off", obj.Name, kind, generator))
	}
	inColumns := make(map[*Attribute]struct{}, len(columns))
	for _, col := range columns {
		if col != nil {
			inColumns[col] = struct{}{}
		}
	}
	seen := make(map[*Attribute]struct{}, len(skipped))
	for _, attr := range skipped {
		if attr == nil {
			panic(fmt.Sprintf("object %q: nil %s %s skip field", obj.Name, kind, generator))
		}
		if _, ok := inColumns[attr]; !ok {
			panic(fmt.Sprintf("object %q: %s %s skip %q is not a %s column", obj.Name, kind, generator, attr.Name, kind))
		}
		if _, dup := seen[attr]; dup {
			panic(fmt.Sprintf("object %q: duplicate %s %s skip %q", obj.Name, kind, generator, attr.Name))
		}
		seen[attr] = struct{}{}
	}
}
