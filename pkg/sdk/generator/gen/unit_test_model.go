package gen

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

// UnitTestsModel is the fully prepared model for the unit_tests generation part.
// All case names, constant names, default modifications, and expected-error expressions
// are precomputed here so the template is pure rendering with almost no logic.
type UnitTestsModel struct {
	// TestContextTypeName is the Go type name of the per-object context struct, e.g. "functionsTestsContext".
	TestContextTypeName string
	// TestDefinitionVarName is the package-level var name, e.g. "functionsTests".
	TestDefinitionVarName string
	// TestNamePrefix is prepended to every case name to form its constant identifier, e.g. "case_Functions_".
	TestNamePrefix string
	// IdVars holds one entry per distinct identifier kind used by this object's operations.
	// The generated file declares them as package-level vars so ext can close over them.
	IdVars []IdVar
	// Operations holds one entry per operation with an OptsField (ShowByID is skipped).
	Operations []*OperationTestModel

	*genhelpers.PreambleModel
}

// IdVar represents a generated "functionsTestId<Kind>" package-level var.
type IdVar struct {
	// Name is the Go identifier, e.g. "functionsTestIdSchemaObjectIdentifier".
	Name string
	// Kind is the identifier type, e.g. "SchemaObjectIdentifier".
	// The constructor call is derived in the template as random<Kind>().
	Kind string
}

// OperationTestModel is the fully prepared model for one operation's test function.
type OperationTestModel struct {
	// Field is the field name in the per-object context struct, e.g. "CreateForJava".
	Field string
	// OptsType is the pointer-to-opts type, e.g. "*CreateForJavaFunctionOptions".
	OptsType string
	// TestNamePrefix is the constant name prefix, e.g. "case_Functions_". Carried here so
	// sub-templates that receive only OperationTestModel can build constant names without needing access to the parent UnitTestsModel.
	TestNamePrefix string
	// DefaultOptsFields holds the field-initializer lines inside the generated defaultProvider,
	// e.g. ["name: functionsTestIdSchemaObject"]. The identifier var name is looked up from IdVars.
	DefaultOptsFields []string
	// ValidationCases is the ordered list of generated validation cases for this operation.
	ValidationCases []*ValidationTestCase
	// SqlCases is the ordered list of generated SQL cases for this operation.
	SqlCases []*SqlTestCase
}

// ValidationTestCase is a fully prepared validation case.
type ValidationTestCase struct {
	// Name is the case name used as the t.Run string.
	Name string
	// ExpectedErrLine is the error expression, pre-rendered from Validation.ReturnedError
	// so the template and the validation body cannot drift.
	ExpectedErrLine string
	// ModifyLines are the statements of the DefaultModify closure body, or nil when not derivable.
	ModifyLines []string
	// HasModify is false when the generator could not derive a modification;
	// the template then omits DefaultModify so the harness fails loudly asking ext to supply it.
	HasModify bool
}

// SqlTestCase is a fully prepared SQL case.
type SqlTestCase struct {
	// Name is the case name used as the t.Run string.
	Name string
	// IsBasic is true for the basic case, which uses the default opts as-is.
	// All other SQL cases require ext to register a modification via withModify.
	IsBasic bool
}

// buildUnitTestsModel walks the interface definition and derives the fully prepared model.
func (i *Interface) buildUnitTestsModel() *UnitTestsModel {
	camelName := genhelpers.FirstLetterLowercase(i.Name) // e.g. "computePools" for Go identifiers
	testDefinitionVarName := camelName + "Tests"
	testContextTypeName := i.Name + "TestsContext"
	testNamePrefix := "case_" + i.Name + "_"

	// Collect the distinct identifier kinds actually used across all operations, in encounter order.
	// idVarsByKind is used for O(1) deduplication; idVars accumulates in insertion order directly.
	idVarsByKind := make(map[string]IdVar)
	idVars := make([]IdVar, 0)

	lookupOrAddIdVar := func(kind string) IdVar {
		if v, ok := idVarsByKind[kind]; ok {
			return v
		}
		name := camelName + "TestId" + kind // e.g. "computePoolsTestIdAccountObjectIdentifier"
		v := IdVar{Name: name, Kind: kind}
		idVarsByKind[kind] = v
		idVars = append(idVars, v)
		return v
	}

	operationTestModels := make([]*OperationTestModel, 0, len(i.Operations))

	for _, op := range i.Operations {
		if op.OptsField == nil {
			continue
		}

		optsType := "*" + op.OptsField.KindNoPtr()
		field := op.Name

		var defaultOptsFields []string
		nameField := op.OptsField.FindChild("name")
		if nameField != nil && nameField.IsIdentifier() {
			idKind := nameField.KindNoPtr()
			defaultOptsFields = []string{fmt.Sprintf("name: %s", lookupOrAddIdVar(idKind).Name)}
		}

		valCases := buildValidationCases(op)
		sqlCases := buildSqlCases(op)

		operationTestModels = append(operationTestModels, &OperationTestModel{
			Field:             field,
			OptsType:          optsType,
			TestNamePrefix:    testNamePrefix,
			DefaultOptsFields: defaultOptsFields,
			ValidationCases:   valCases,
			SqlCases:          sqlCases,
		})
	}

	return &UnitTestsModel{
		TestContextTypeName:   testContextTypeName,
		TestDefinitionVarName: testDefinitionVarName,
		TestNamePrefix:        testNamePrefix,
		IdVars:                idVars,
		Operations:            operationTestModels,
	}
}

// buildValidationCases derives all unit test validation cases for one operation recursively traversing its subtree.
func buildValidationCases(op *Operation) []*ValidationTestCase {
	cases := make([]*ValidationTestCase, 0)
	collectValidationCases(op.OptsField, op.Name, &cases)
	return cases
}

func collectValidationCases(f *Field, opName string, out *[]*ValidationTestCase) {
	for _, v := range f.Validations {
		if v.IsAdditionalValidations() {
			continue
		}

		casesForValidation := buildCasesForValidation(v, f, opName)
		*out = append(*out, casesForValidation...)
	}

	for idx := range f.Fields {
		child := &f.Fields[idx]
		if child.HasAnyValidationInSubtree() {
			collectValidationCases(child, opName, out)
		}
	}
}

func buildCasesForValidation(v *Validation, f *Field, opName string) []*ValidationTestCase {
	expectedErrLine := v.ReturnedError(f)
	if !v.IsMultiField() {
		return []*ValidationTestCase{buildSingleFieldValidationCase(v, f, opName, expectedErrLine)}
	}
	return buildMultiFieldValidationCases(v, f, opName, expectedErrLine)
}

func buildSingleFieldValidationCase(v *Validation, f *Field, opName, expectedErrLine string) *ValidationTestCase {
	// Use the field name, not the container slug, so two ValidateValueSet validations on
	// different fields of the same container produce distinct names
	// (e.g. Handler_ValidateValueSet vs RuntimeVersion_ValidateValueSet).
	fieldName := v.FieldNames[0]
	var fieldSlug string
	if f.IsRoot() {
		fieldSlug = fieldName
	} else {
		// container.SlugPath() = "opts_A_B_..._Container"; strip the leading "opts_"
		fieldSlug = strings.TrimPrefix(f.SlugPath(), "opts_") + "_" + fieldName
	}
	caseName := fmt.Sprintf("validation_%s_%s_%s", opName, fieldSlug, v.TypeName())
	lines, ok := v.DeriveModify(f)
	return &ValidationTestCase{Name: caseName, ExpectedErrLine: expectedErrLine, ModifyLines: lines, HasModify: ok}
}

// buildMultiFieldValidationCases generates per-shape cases for multi-field validations:
// ConflictingFields   -> one case (all fields set)
// AtLeastOneValueSet  -> one case (all fields zeroed)
// ExactlyOneValueSet  -> NoneSet + MoreThanOneSet (+ OneValidOneInvalid for slices)
// MoreThanOneValueSet -> MoreThanOneSet (+ OneValidOneInvalid for slices)
func buildMultiFieldValidationCases(v *Validation, f *Field, opName, expectedErrLine string) []*ValidationTestCase {
	baseName := fmt.Sprintf("validation_%s_%s_%s", opName, f.SlugPath(), v.TypeName())

	tc := func(name string, lines []string, ok bool) *ValidationTestCase {
		return &ValidationTestCase{Name: name, ExpectedErrLine: expectedErrLine, ModifyLines: lines, HasModify: ok}
	}

	var cases []*ValidationTestCase
	switch v.Type {
	case ConflictingFields:
		lines, ok := deriveConflictingFieldsModify(v, f)
		cases = append(cases, tc(baseName, lines, ok))

	case AtLeastOneValueSet:
		lines, ok := deriveAtLeastOneValueSetModify(v, f)
		cases = append(cases, tc(baseName, lines, ok))

	case ExactlyOneValueSet:
		noneSetLines, noneSetOk := deriveNoneSetModify(v, f)
		cases = append(cases, tc(baseName+"_NoneSet", noneSetLines, noneSetOk))
		moreThanOneSetLines, moreThanOneSetOk := deriveMoreThanOneSetModify(v, f)
		cases = append(cases, tc(baseName+"_MoreThanOneSet", moreThanOneSetLines, moreThanOneSetOk))
		if f.IsSlice() {
			cases = append(cases, tc(baseName+"_OneValidOneInvalid", nil, false))
		}

	case MoreThanOneValueSet:
		moreThanOneSetLines, moreThanOneSetOk := deriveMoreThanOneSetModify(v, f)
		cases = append(cases, tc(baseName+"_MoreThanOneSet", moreThanOneSetLines, moreThanOneSetOk))
		if f.IsSlice() {
			cases = append(cases, tc(baseName+"_OneValidOneInvalid", nil, false))
		}
	}
	return cases
}

// deriveConflictingFieldsModify sets ALL listed fields to a non-zero value (all set → conflict).
func deriveConflictingFieldsModify(v *Validation, f *Field) ([]string, bool) {
	prime := primeAncestors(f)
	stmts := make([]string, 0, len(v.FieldNames))
	for _, name := range v.FieldNames {
		child := f.FindChild(name)
		if child == nil {
			return nil, false
		}
		nz, ok := nonZeroValueFor(child)
		if !ok {
			return nil, false
		}
		stmts = append(stmts, fmt.Sprintf("opts%s.%s = %s", f.Path(), name, nz))
	}
	return append(prime, stmts...), true
}

// deriveAtLeastOneValueSetModify primes the container field and then zeroes all listed fields.
func deriveAtLeastOneValueSetModify(v *Validation, f *Field) ([]string, bool) {
	// The container itself must be non-nil for the validation to run at all. primeAncestors
	// handles ancestors up to f's parent; we prime f itself separately.
	prime := primeAncestors(f)
	if f.IsPointer() {
		prime = append(prime, fmt.Sprintf("opts%s = &%s{}", f.Path(), f.KindNoPtr()))
	}
	stmts := make([]string, 0, len(v.FieldNames))
	for _, name := range v.FieldNames {
		child := f.FindChild(name)
		if child == nil {
			return nil, false
		}
		stmts = append(stmts, fmt.Sprintf("opts%s.%s = %s", f.Path(), name, zeroValueFor(child)))
	}
	return append(prime, stmts...), true
}

// deriveNoneSetModify sets all listed fields to their zero values (none set → ExactlyOneValueSet fails).
func deriveNoneSetModify(v *Validation, f *Field) ([]string, bool) {
	prime := primeAncestors(f)
	if f.IsSlice() {
		// For a slice field the validation runs per-element. A single-element slice whose element
		// has all listed fields at zero value is the correct "none set" setup.
		// e.g. opts.Arguments = []FunctionArgument{{}}
		return append(prime, fmt.Sprintf("opts%s = []%s{{}}", f.Path(), f.KindNoPtr())), true
	}
	if f.IsPointer() {
		prime = append(prime, fmt.Sprintf("opts%s = &%s{}", f.Path(), f.KindNoPtr()))
	}
	stmts := make([]string, 0, len(v.FieldNames))
	for _, name := range v.FieldNames {
		child := f.FindChild(name)
		if child == nil {
			return nil, false
		}
		stmts = append(stmts, fmt.Sprintf("opts%s.%s = %s", f.Path(), name, zeroValueFor(child)))
	}
	return append(prime, stmts...), true
}

// deriveMoreThanOneSetModify sets the first two listed fields to non-zero values.
func deriveMoreThanOneSetModify(v *Validation, f *Field) ([]string, bool) {
	if f.IsSlice() {
		// Slice-element fields like DataType / datatypes.DataType require domain-specific values.
		// Ext must supply the modification.
		return nil, false
	}
	prime := primeAncestors(f)
	if f.IsPointer() {
		prime = append(prime, fmt.Sprintf("opts%s = &%s{}", f.Path(), f.KindNoPtr()))
	}
	if len(v.FieldNames) < 2 {
		return nil, false
	}
	stmts := make([]string, 0, 2)
	for _, name := range v.FieldNames[:2] {
		child := f.FindChild(name)
		if child == nil {
			return nil, false
		}
		nz, ok := nonZeroValueFor(child)
		if !ok {
			return nil, false
		}
		stmts = append(stmts, fmt.Sprintf("opts%s.%s = %s", f.Path(), name, nz))
	}
	return append(prime, stmts...), true
}

// buildSqlCases derives the SQL cases for one operation based on its test kind.
func buildSqlCases(op *Operation) []*SqlTestCase {
	cases := make([]*SqlTestCase, 0)
	opName := op.Name

	basic := func(name string) *SqlTestCase { return &SqlTestCase{Name: name, IsBasic: true} }
	extReq := func(name string) *SqlTestCase { return &SqlTestCase{Name: name} }

	switch op.TestKind() {
	case OperationTestKindCreate:
		cases = append(cases, basic(fmt.Sprintf("sql_%s_basic", opName)), extReq(fmt.Sprintf("sql_%s_all", opName)))

	case OperationTestKindAlter:
		exclusiveOptions := op.MutuallyExclusiveOptions()
		if len(exclusiveOptions) == 0 {
			cases = append(cases, basic(fmt.Sprintf("sql_%s_basic", opName)), extReq(fmt.Sprintf("sql_%s_all", opName)))
		} else {
			for _, branch := range exclusiveOptions {
				cases = append(cases, extReq(fmt.Sprintf("sql_%s_%s", opName, branch.Name)))
			}
		}

	case OperationTestKindShow:
		cases = append(cases, basic(fmt.Sprintf("sql_%s_basic", opName)), extReq(fmt.Sprintf("sql_%s_all", opName)))
		for _, field := range op.ShowOptionalFields() {
			cases = append(cases, extReq(fmt.Sprintf("sql_%s_%s", opName, field.Name)))
		}

	case OperationTestKindDrop:
		cases = append(cases, basic(fmt.Sprintf("sql_%s_basic", opName)), extReq(fmt.Sprintf("sql_%s_all", opName)))

	case OperationTestKindDescribe:
		cases = append(cases, basic(fmt.Sprintf("sql_%s_basic", opName)))

	case OperationTestKindOther:
		// No SQL cases for Grant/Revoke/etc.
	}

	return cases
}
