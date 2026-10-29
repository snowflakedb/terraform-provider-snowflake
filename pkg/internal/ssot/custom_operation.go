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
)

// CustomOperationClause is one custom-operation entry, in emission order.
type CustomOperationClause struct {
	Kind       CustomOperationClauseKind
	SQL        string
	Attribute  *Attribute
	Predefined string
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

// CustomOperationSpec is one custom operation.
// Kind is the method name. The query struct name is Kind plus Name.
// Doc overrides the default doc URL. AdditionalValidations adds that validation call.
type CustomOperationSpec struct {
	Kind                  string
	Doc                   string
	AdditionalValidations bool
	Clauses               []CustomOperationClause
}
