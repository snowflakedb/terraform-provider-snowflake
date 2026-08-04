package sdk

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

// testCaseName identifies a single generated unit test case.
// Generated constants live in the object's _gen_test.go file;
// Use these references to modify/add tests in the _ext_test.go file.
type testCaseName string

// caseNilOptions is shared by every operation — the case is emitted by the runner, not generated.
const caseNilOptions testCaseName = "validation_nilOptions"

// sdkTestCase is the part of a test case that the shared skip/default/modify handling needs.
type sdkTestCase[PT validatable] interface {
	caseName() testCaseName
	modification() func(PT)
	noModifyNeeded() bool
}

// validationCase is a single "these opts must not validate" case.
type validationCase[PT validatable] struct {
	// Name is the case name and the registry key.
	Name testCaseName
	// ExpectedErr is matched as a substring against the joined validation error, see assertOptsInvalidJoinedErrors.
	ExpectedErr error
	// DefaultModify is the automatically generated modification, or nil when the generator could not derive one.
	// A nil DefaultModify with no ext registration fails the test with instructions.
	DefaultModify func(PT)
}

// sqlCase is a single "these opts validate and render to this SQL" case.
// The expected SQL is never derivable and always comes from the ext file.
type sqlCase[PT validatable] struct {
	Name testCaseName
	// NoModifyNeeded is true for the `basic` case, which uses the default opts as-is.
	// All other SQL cases require ext to register a modification via withModify.
	NoModifyNeeded bool
}

func (c validationCase[PT]) caseName() testCaseName { return c.Name }
func (c validationCase[PT]) modification() func(PT) { return c.DefaultModify }
func (c validationCase[PT]) noModifyNeeded() bool   { return false }

func (c sqlCase[PT]) caseName() testCaseName { return c.Name }
func (c sqlCase[PT]) modification() func(PT) { return nil }
func (c sqlCase[PT]) noModifyNeeded() bool   { return c.NoModifyNeeded }

// sdkTestCtx carries everything needed to run the generated unit tests of one SDK operation.
//
// The generated file builds it during package variable initialization (generated cases + generated default opts).
// The _ext_test.go file refines it from init(), which the Go runtime always runs after all package-level variables are initialized.
type sdkTestCtx[PT validatable] struct {
	objectName    string // e.g. "Functions"    — the Interface.Name; source of all derived names
	operationName string // e.g. "CreateForJava" — the operation field name

	defaultProvider func() PT

	validationCases []validationCase[PT]
	sqlCases        []sqlCase[PT]

	// Keeping modify and expectedSqls currently not as parts of validationCase or sqlCase makes the implementation easier
	// (check for presence in one map instead of 4); we can iterate on this implementation as we go.
	modify       map[testCaseName]func(PT)
	expectedSqls map[testCaseName]string

	extraValidationCases []validationCase[PT]
	extraSqlCases        []sqlCase[PT]
}

// newSdkTestCtx creates an empty test context for one SDK operation.
//
//   - objectName is the Interface.Name (e.g. "Functions"), the source from which ctxPath,
//     constPfx and extFile are all derived.
//   - operationName is the operation field name (e.g. "CreateForJava"), used in ctxPath.
func newSdkTestCtx[PT validatable](objectName, operationName string) *sdkTestCtx[PT] {
	return &sdkTestCtx[PT]{
		objectName:    objectName,
		operationName: operationName,

		modify:       make(map[testCaseName]func(PT)),
		expectedSqls: make(map[testCaseName]string),
	}
}

// ctxPath returns the Go expression that reaches this ctx from a test, e.g. "functionsTests.CreateForJava".
func (c *sdkTestCtx[PT]) ctxPath() string {
	return strings.ToLower(c.objectName[:1]) + c.objectName[1:] + "Tests." + c.operationName
}

// constPfx returns the constant name prefix, e.g. "case_Functions_".
func (c *sdkTestCtx[PT]) constPfx() string {
	return "case_" + c.objectName + "_"
}

// extFile returns the companion ext test file name, e.g. "functions_ext_test.go".
func (c *sdkTestCtx[PT]) extFile() string {
	return genhelpers.ToSnakeCase(c.objectName) + "_ext_test.go"
}

// optsTypeName returns the unqualified type name of PT for use in failure messages,
// e.g. "*CreateForJavaFunctionOptions" for PT = *CreateForJavaFunctionOptions.
func (c *sdkTestCtx[PT]) optsTypeName() string {
	t := reflect.TypeFor[PT]()
	if t.Kind() == reflect.Pointer {
		return "*" + t.Elem().Name()
	}
	return t.Name()
}

// withDefaultOpts sets the default opts provider.
// Called from the generated file with a minimal-valid provider;
// the _ext_test.go init() may call it again to replace it when additional required fields must be set.
func (c *sdkTestCtx[PT]) withDefaultOpts(f func() PT) *sdkTestCtx[PT] {
	c.defaultProvider = f
	return c
}

func (c *sdkTestCtx[PT]) withValidationCases(cases ...validationCase[PT]) *sdkTestCtx[PT] {
	c.validationCases = cases
	return c
}

func (c *sdkTestCtx[PT]) withSqlCases(cases ...sqlCase[PT]) *sdkTestCtx[PT] {
	c.sqlCases = cases
	return c
}

// withModify registers the modification for a generated case, overriding its DefaultModify.
// Panics immediately if name is not a case registered on this context, so a stale reference
// left behind after regeneration is caught at package init instead of silently doing nothing.
func (c *sdkTestCtx[PT]) withModify(name testCaseName, f func(PT)) *sdkTestCtx[PT] {
	c.assertCaseExists(name, "withModify")
	c.modify[name] = f
	return c
}

// withExpectedSql registers the SQL a sql case must render to.
func (c *sdkTestCtx[PT]) withExpectedSql(name testCaseName, sql string) *sdkTestCtx[PT] {
	c.assertSqlCaseExists(name, "withExpectedSql")
	c.expectedSqls[name] = sql
	return c
}

// withExpectedSqlf is withExpectedSql with fmt.Sprintf applied.
func (c *sdkTestCtx[PT]) withExpectedSqlf(name testCaseName, format string, args ...any) *sdkTestCtx[PT] {
	return c.withExpectedSql(name, fmt.Sprintf(format, args...))
}

// withModifyAndExpectedSqlf registers both the modification and the expected SQL for a generated
// SQL case in a single call. It is the common case for non-basic SQL cases where the ext file
// must supply both.
func (c *sdkTestCtx[PT]) withModifyAndExpectedSqlf(name testCaseName, modify func(PT), format string, args ...any) *sdkTestCtx[PT] {
	return c.withModify(name, modify).
		withExpectedSqlf(name, format, args...)
}

// withAdditionalValidationCase adds a case with no generated counterpart — typically one produced
// by additionalValidations(), which the generator deliberately does not emit.
// Panics if a generated validation or SQL case with the same name already exists, to prevent
// silently shadowing a generated case with an extra one.
func (c *sdkTestCtx[PT]) withAdditionalValidationCase(name string, modify func(PT), expectedErr error) *sdkTestCtx[PT] {
	c.assertCaseDoesNotExist(testCaseName(name), "withAdditionalValidationCase")
	c.extraValidationCases = append(c.extraValidationCases, validationCase[PT]{
		Name:          testCaseName(name),
		ExpectedErr:   expectedErr,
		DefaultModify: modify,
	})
	return c
}

// withAdditionalSqlCasef adds a SQL case with no generated counterpart.
// Panics if a generated validation or SQL case with the same name already exists, to prevent
// silently shadowing a generated case with an extra one.
func (c *sdkTestCtx[PT]) withAdditionalSqlCasef(name string, modify func(PT), format string, args ...any) *sdkTestCtx[PT] {
	n := testCaseName(name)
	c.assertCaseDoesNotExist(n, "withAdditionalSqlCasef")
	c.extraSqlCases = append(c.extraSqlCases, sqlCase[PT]{Name: n})
	c.modify[n] = modify
	c.expectedSqls[n] = fmt.Sprintf(format, args...)
	return c
}

// RunValidationCases runs the shared nil-options case, every generated validation case, and
// finally the ext-only cases.
func (c *sdkTestCtx[PT]) RunValidationCases(t *testing.T) {
	t.Helper()

	t.Run(string(caseNilOptions), func(t *testing.T) {
		var opts PT // typed nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	cases := make([]validationCase[PT], 0, len(c.validationCases)+len(c.extraValidationCases))
	cases = append(cases, c.validationCases...)
	cases = append(cases, c.extraValidationCases...)

	for _, tc := range cases {
		t.Run(string(tc.Name), func(t *testing.T) {
			opts := c.prepareOpts(t, tc)
			assertOptsInvalidJoinedErrors(t, opts, tc.ExpectedErr)
		})
	}
}

// RunSqlCases runs every generated SQL case and then the ext-only ones.
func (c *sdkTestCtx[PT]) RunSqlCases(t *testing.T) {
	t.Helper()

	cases := make([]sqlCase[PT], 0, len(c.sqlCases)+len(c.extraSqlCases))
	cases = append(cases, c.sqlCases...)
	cases = append(cases, c.extraSqlCases...)

	for _, tc := range cases {
		t.Run(string(tc.Name), func(t *testing.T) {
			opts := c.prepareOpts(t, tc)
			sql, ok := c.expectedSqls[tc.Name]
			if !ok {
				t.Fatalf(
					"no expected SQL registered for case %[1]q.\nRegister it in %[2]s:\n\n\t%[3]s.\n\t\twithExpectedSqlf(%[4]s, `...`, ...)\n\nOr if a modification is also needed:\n\n\t%[3]s.\n\t\twithModifyAndExpectedSqlf(%[4]s, func(opts %[5]s) { ... }, `...`, ...)\n",
					tc.Name, c.extFile(), c.ctxPath(), c.constFor(tc.Name), c.optsTypeName(),
				)
			}
			// "%s" rather than sql-as-format so a literal % in the SQL is not interpreted.
			// TODO [SNOW-3882943]: with the current approach, the assertOptsValidAndSQLEquals should have a formatted but also explicit formats - address after all the tests are migrated
			assertOptsValidAndSQLEquals(t, opts, "%s", sql)
		})
	}
}

// prepareOpts builds the default opts and applies the modification —
// ext registration first, generated default second, loud failure last.
func (c *sdkTestCtx[PT]) prepareOpts(t *testing.T, tc sdkTestCase[PT]) PT {
	t.Helper()
	name := tc.caseName()

	if c.defaultProvider == nil {
		t.Fatalf("no default opts provider for %s; register one in %s via %s.withDefaultOpts(...)",
			c.ctxPath(), c.extFile(), c.ctxPath())
	}
	opts := c.defaultProvider()

	switch {
	case c.modify[name] != nil:
		c.modify[name](opts) // ext wins
	case tc.noModifyNeeded():
		// basic SQL case: use default opts as-is
	case tc.modification() != nil:
		tc.modification()(opts) // generated default (validation cases only)
	default:
		t.Fatalf(
			"no modification registered for case %[1]q — the generator could not derive one.\nRegister it in %[2]s:\n\n\t%[3]s.\n\t\twithModify(%[4]s, func(opts %[5]s) { ... })\n\nOr if the expected SQL also needs to be registered at the same time:\n\n\t%[3]s.\n\t\twithModifyAndExpectedSqlf(%[4]s, func(opts %[5]s) { ... }, `...`, ...)\n",
			name, c.extFile(), c.ctxPath(), c.constFor(name), c.optsTypeName(),
		)
	}
	return opts
}

// constFor names the generated constant for a case: prefix + case name verbatim.
func (c *sdkTestCtx[PT]) constFor(name testCaseName) string {
	return c.constPfx() + string(name)
}

func (c *sdkTestCtx[PT]) assertCaseExists(name testCaseName, method string) {
	if c.hasValidationCase(name) || c.hasSqlCase(name) {
		return
	}
	panic(fmt.Sprintf(
		"%s(%q) on %s: no such generated case. It may have been renamed or removed by regeneration — check %s.",
		method, name, c.ctxPath(), c.extFile(),
	))
}

func (c *sdkTestCtx[PT]) assertCaseDoesNotExist(name testCaseName, method string) {
	if c.hasValidationCase(name) || c.hasSqlCase(name) {
		panic(fmt.Sprintf(
			"%s(%q) on %s: a generated case with this name already exists. Use withModify or withModifyAndExpectedSqlf to configure it instead of adding a duplicate.",
			method, name, c.ctxPath(),
		))
	}
}

func (c *sdkTestCtx[PT]) assertSqlCaseExists(name testCaseName, method string) {
	if c.hasSqlCase(name) {
		return
	}
	panic(fmt.Sprintf(
		"%s(%q) on %s: no such generated sql case. It may have been renamed or removed by regeneration — check %s.",
		method, name, c.ctxPath(), c.extFile(),
	))
}

func (c *sdkTestCtx[PT]) hasValidationCase(name testCaseName) bool {
	for _, tc := range c.validationCases {
		if tc.Name == name {
			return true
		}
	}
	return false
}

func (c *sdkTestCtx[PT]) hasSqlCase(name testCaseName) bool {
	for _, tc := range c.sqlCases {
		if tc.Name == name {
			return true
		}
	}
	return false
}
