package gen

// NewInstanceMethodCallStruct creates a QueryStruct for a Snowflake instance method invocation.
// The identifier field kind is set to "<will be replaced>" so that InstanceMethodOperation
// can substitute the interface's identifier kind automatically (same mechanism as Name()).
// Pass nil argsQueryStruct for no-arg methods; an empty struct will be created automatically,
// named after the method (e.g. "GET_SPENDING_LIMIT" → "GetSpendingLimitArgs").
func NewInstanceMethodCallStruct(structName string, methodName string, argsQueryStruct *QueryStruct) *QueryStruct {
	qs := NewQueryStruct(structName)
	qs.Call()
	identifier := NewField("name", "<will be replaced>", Tags().Identifier().InstanceMethod().SQL(methodName), IdentifierOptions().Required())
	qs.identifierField = identifier
	qs.fields = append(qs.fields, identifier)
	if argsQueryStruct == nil {
		qs.QueryStructField("args", NewQueryStruct(sqlToFieldName(methodName, false)+"Args"), ListOptions().MustParentheses())
	} else {
		qs.QueryStructField("args", argsQueryStruct, ListOptions().MustParentheses().Required())
	}
	return qs
}

// InstanceMethodOperation adds a Snowflake instance method call operation to the interface.
// The operation name (e.g. "AddNotificationIntegration") is derived from methodName
// (e.g. "ADD_NOTIFICATION_INTEGRATION").
// Pass a non-nil pairedStructs to also generate db row / plain result structs with a convert mapping;
// pass nil when the method produces no result to map.
func (i *Interface) InstanceMethodOperation(doc string, methodName string, argsQueryStruct *QueryStruct, pairedStructs *PairedStructs, helperStructs ...IntoField) *Interface {
	operationName := sqlToFieldName(methodName, true)
	qs := NewInstanceMethodCallStruct(operationName+"Options", methodName, argsQueryStruct).
		WithValidation(ValidIdentifier, "name")
	if pairedStructs == nil {
		return i.newSimpleOperation(operationName, doc, qs, helperStructs...)
	}
	i.newOperationWithDBMapping(operationName, doc, pairedStructs.asDbStruct(), pairedStructs.asPlainStruct(), qs, addInstanceMethodMapping, helperStructs...)
	return i
}
