package ssot

// CustomOperationClauseKind is one entry in a custom operation.
type CustomOperationClauseKind string

const (
	CustomOperationClauseStaticSQL   CustomOperationClauseKind = "static_sql"
	CustomOperationClauseOrReplace   CustomOperationClauseKind = "or_replace"
	CustomOperationClauseIfNotExists CustomOperationClauseKind = "if_not_exists"
	CustomOperationClauseAttribute   CustomOperationClauseKind = "attribute"
	CustomOperationClauseObjectWords CustomOperationClauseKind = "object_words"
	CustomOperationClauseName        CustomOperationClauseKind = "name"
	CustomOperationClausePredefined  CustomOperationClauseKind = "predefined"
	// CustomOperationClauseStaticField is static SQL with an explicit Go field name.
	CustomOperationClauseStaticField CustomOperationClauseKind = "static_field"
	// CustomOperationClauseIdentifier is a required identifier field.
	CustomOperationClauseIdentifier CustomOperationClauseKind = "identifier"
	// CustomOperationClauseText is a required single-quoted text assignment.
	CustomOperationClauseText CustomOperationClauseKind = "text"
)

// CustomOperationClause is one custom-operation entry, in emission order.
type CustomOperationClause struct {
	Kind            CustomOperationClauseKind
	SQL             string
	Field           string
	IdType          string
	ValidIdentifier bool
	Attribute       *Attribute
	Predefined      string
}

// CustomStaticSQL is a static SQL word, such as CREATE or UNDROP.
func CustomStaticSQL(sql string) CustomOperationClause {
	return CustomOperationClause{Kind: CustomOperationClauseStaticSQL, SQL: sql}
}

// CustomOrReplace is optional OR REPLACE.
func CustomOrReplace() CustomOperationClause {
	return CustomOperationClause{Kind: CustomOperationClauseOrReplace}
}

// CustomIfNotExists is optional IF NOT EXISTS.
func CustomIfNotExists() CustomOperationClause {
	return CustomOperationClause{Kind: CustomOperationClauseIfNotExists}
}

// CustomAttribute emits the attribute's create clause.
func CustomAttribute(attr *Attribute) CustomOperationClause {
	return CustomOperationClause{Kind: CustomOperationClauseAttribute, Attribute: attr}
}

// CustomObjectWords is the uppercase words of Name.
func CustomObjectWords() CustomOperationClause {
	return CustomOperationClause{Kind: CustomOperationClauseObjectWords}
}

// CustomName is the object identifier. The type is the object's IdType.
func CustomName() CustomOperationClause {
	return CustomOperationClause{Kind: CustomOperationClauseName}
}

// CustomPredefined is a required keyword field.
// The field name and the predefined type use the same spelling.
func CustomPredefined(typeName string) CustomOperationClause {
	return CustomOperationClause{Kind: CustomOperationClausePredefined, Predefined: typeName}
}

// CustomStaticField is static SQL with an explicit Go field name, such as TYPE = COLUMNS.
func CustomStaticField(field, sql string) CustomOperationClause {
	return CustomOperationClause{Kind: CustomOperationClauseStaticField, Field: field, SQL: sql}
}

// CustomIdentifier is a required identifier. valid marks it with ValidIdentifier.
func CustomIdentifier(field, idType string, valid bool) CustomOperationClause {
	return CustomOperationClause{Kind: CustomOperationClauseIdentifier, Field: field, IdType: idType, ValidIdentifier: valid}
}

// CustomText is a required single-quoted assignment, such as TABLE_SCHEMA.
func CustomText(sql string) CustomOperationClause {
	return CustomOperationClause{Kind: CustomOperationClauseText, SQL: sql}
}

// CustomPairedShow is the row and plain pair of a custom operation that returns rows.
// The mapping is a slice. Columns are paired SHOW fields. PlainOnly is on the plain struct only.
type CustomPairedShow struct {
	DbRow     string
	Plain     string
	Columns   []*Attribute
	PlainOnly []SdkPlainOnlyField
}

// CustomOperationSpec is one custom operation.
// Kind is the method name. The query struct name is QueryStructName, or Kind plus Name.
// Doc overrides the default doc URL. AdditionalValidations adds that validation call.
// Paired emits CustomShowOperationWithPairedStructs instead of CustomOperation.
type CustomOperationSpec struct {
	Kind                  string
	Doc                   string
	QueryStructName       string
	AdditionalValidations bool
	Clauses               []CustomOperationClause
	Paired                *CustomPairedShow
}
