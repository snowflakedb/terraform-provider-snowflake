package gen

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"

// NewInstanceMethodCallStruct creates a QueryStruct for a Snowflake instance method invocation.
// The identifier field kind is set to "<will be replaced>" so that InstanceMethodOperation
// can substitute the interface's identifier kind automatically (same mechanism as Name()).
// Pass nil argsQueryStruct for no-arg methods; an empty struct will be created automatically,
// named after the method and interface (e.g. Budgets "GET_SPENDING_LIMIT" → "budgetGetSpendingLimitArgs").
func (i *Interface) newInstanceMethodCallStruct(structName string, methodName string, argsQueryStruct *QueryStruct) *QueryStruct {
	qs := NewQueryStruct(structName)
	qs.Call()
	identifier := NewField("name", "<will be replaced>", Tags().Identifier().InstanceMethod().SQL(methodName), IdentifierOptions().Required())
	qs.identifierField = identifier
	qs.fields = append(qs.fields, identifier)
	if argsQueryStruct == nil {
		// TODO [this PR]: improve naming handling
		qs.QueryStructField("args", NewQueryStruct(sqlToFieldName(genhelpers.ToSnakeCase(i.NameSingular)+"_"+methodName, false)+"Args"), ListOptions().MustParentheses())
	} else {
		argsQueryStruct.name = i.NameSingular + argsQueryStruct.name
		qs.QueryStructField("args", argsQueryStruct, ListOptions().MustParentheses().Required())
	}
	return qs
}

// InstanceMethodOperation adds a Snowflake instance method call operation to the interface.
// The operation name (e.g. "AddNotificationIntegration") is derived from methodName
// (e.g. "ADD_NOTIFICATION_INTEGRATION").
// Pass kind to control the result:
//   - InstanceMethodKindSingleValue: single-value struct result, requires non-nil pairedStructs
//   - InstanceMethodKindSlice: slice struct result, requires non-nil pairedStructs
//   - InstanceMethodKind("<goType>"): scalar return (e.g. "int", "string"), pairedStructs must be nil
//   - "" (empty): void operation, pairedStructs must be nil
//
// TODO [this PR]: introduce a dedicated field and helper method for scalar return type
func (i *Interface) InstanceMethodOperation(doc string, methodName string, argsQueryStruct *QueryStruct, pairedStructs *PairedStructs, kind InstanceMethodKind, helperStructs ...IntoField) *Interface {
	operationName := sqlToFieldName(methodName, true)
	qs := i.newInstanceMethodCallStruct(operationName+"Options", methodName, argsQueryStruct).
		WithValidation(ValidIdentifier, "name")
	switch kind {
	case InstanceMethodKindSingleValue, InstanceMethodKindSlice:
		i.newOperationWithDBMapping(operationName, doc, pairedStructs.asDbStruct(), pairedStructs.asPlainStruct(), qs, instanceMethodMappingForKind(kind), helperStructs...)
	default:
		i.newSimpleOperation(operationName, doc, qs, helperStructs...)
		if kind != "" {
			// scalar: store the return type in InstanceMethodKind; no InstanceMethodMapping
			op := i.Operations[len(i.Operations)-1]
			op.InstanceMethodKind = &kind
		}
	}
	return i
}
