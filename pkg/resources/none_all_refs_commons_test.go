package resources

import (
	"errors"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func threeBranchTestConfig() NoneAllRefsBlockConfig {
	return NoneAllRefsBlockConfig{
		AttrPath: "allowed_things",
		WithNone: true,
		WithAll:  true,
		RefsKey:  "things",
	}
}

func twoBranchTestConfig() NoneAllRefsBlockConfig {
	return NoneAllRefsBlockConfig{
		AttrPath: "allowed_things",
		WithNone: true,
		RefsKey:  "things",
	}
}

func refsSet(values ...string) *schema.Set {
	elems := make([]any, len(values))
	for i, v := range values {
		elems[i] = v
	}
	return schema.NewSet(schema.HashString, elems)
}

func stateBlock(t *testing.T, state []any) map[string]any {
	t.Helper()
	require.Len(t, state, 1)
	block, ok := state[0].(map[string]any)
	require.True(t, ok, "expected the block to be a map[string]any, got %T", state[0])
	return block
}

func stateRefs(t *testing.T, block map[string]any, refsKey string) []string {
	t.Helper()
	set, ok := block[refsKey].(*schema.Set)
	require.True(t, ok, "expected %s to be a *schema.Set, got %T", refsKey, block[refsKey])
	return expandStringList(set.List())
}

func Test_NoneAllRefsBlockConfig_branchNames(t *testing.T) {
	testCases := []struct {
		Name     string
		Config   NoneAllRefsBlockConfig
		Expected []string
	}{
		{
			Name:     "all three branches, in none-all-refs order",
			Config:   threeBranchTestConfig(),
			Expected: []string{"none", "all", "things"},
		},
		{
			Name:     "none and refs only",
			Config:   twoBranchTestConfig(),
			Expected: []string{"none", "things"},
		},
		{
			Name:     "all and refs only",
			Config:   NoneAllRefsBlockConfig{AttrPath: "allowed_things", WithAll: true, RefsKey: "things"},
			Expected: []string{"all", "things"},
		},
		{
			Name:     "refs is always present",
			Config:   NoneAllRefsBlockConfig{AttrPath: "allowed_things", RefsKey: "things"},
			Expected: []string{"things"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, tc.Config.branchNames())
		})
	}
}

func Test_NoneAllRefsBlockConfig_exactlyOneOf(t *testing.T) {
	t.Run("every branch is qualified with the block's first element", func(t *testing.T) {
		assert.Equal(
			t,
			[]string{"allowed_things.0.none", "allowed_things.0.all", "allowed_things.0.things"},
			threeBranchTestConfig().exactlyOneOf(),
		)
	})

	t.Run("absent branches are left out", func(t *testing.T) {
		assert.Equal(
			t,
			[]string{"allowed_things.0.none", "allowed_things.0.things"},
			twoBranchTestConfig().exactlyOneOf(),
		)
	})
}

func Test_NoneAllRefsBlockInnerSchema(t *testing.T) {
	t.Run("all three branches are declared and mutually exclusive", func(t *testing.T) {
		inner := NoneAllRefsBlockInnerSchema(threeBranchTestConfig())

		require.Len(t, inner, 3)
		exactlyOneOf := []string{"allowed_things.0.none", "allowed_things.0.all", "allowed_things.0.things"}
		for _, branch := range []string{"none", "all", "things"} {
			require.Contains(t, inner, branch)
			assert.True(t, inner[branch].Optional, "%s should be optional", branch)
			assert.Equal(t, exactlyOneOf, inner[branch].ExactlyOneOf, "%s should be mutually exclusive with every other branch", branch)
		}

		assert.Equal(t, schema.TypeBool, inner["none"].Type)
		assert.Equal(t, schema.TypeBool, inner["all"].Type)
		assert.Equal(t, schema.TypeSet, inner["things"].Type)
	})

	t.Run("absent branches are not declared", func(t *testing.T) {
		inner := NoneAllRefsBlockInnerSchema(twoBranchTestConfig())

		require.Len(t, inner, 2)
		assert.Contains(t, inner, "none")
		assert.NotContains(t, inner, "all")
		assert.Equal(t, []string{"allowed_things.0.none", "allowed_things.0.things"}, inner["things"].ExactlyOneOf)
	})

	t.Run("refs key names the attribute", func(t *testing.T) {
		cfg := threeBranchTestConfig()
		cfg.RefsKey = "secrets"

		inner := NoneAllRefsBlockInnerSchema(cfg)

		assert.Contains(t, inner, "secrets")
		assert.NotContains(t, inner, "things")
		assert.NotContains(t, inner, "refs")
	})

	t.Run("descriptions are passed through", func(t *testing.T) {
		cfg := threeBranchTestConfig()
		cfg.NoneDescription, cfg.AllDescription, cfg.RefsDescription = "no things", "all things", "some things"

		inner := NoneAllRefsBlockInnerSchema(cfg)

		assert.Equal(t, "no things", inner["none"].Description)
		assert.Equal(t, "all things", inner["all"].Description)
		assert.Equal(t, "some things", inner["things"].Description)
	})

	t.Run("refs elements are strings validated by the given validator", func(t *testing.T) {
		cfg := threeBranchTestConfig()
		cfg.RefsElemValidator = func(any, cty.Path) diag.Diagnostics { return nil }

		inner := NoneAllRefsBlockInnerSchema(cfg)

		elem, ok := inner["things"].Elem.(*schema.Schema)
		require.True(t, ok, "expected the refs element to be a *schema.Schema, got %T", inner["things"].Elem)
		assert.Equal(t, schema.TypeString, elem.Type)
		assert.NotNil(t, elem.ValidateDiagFunc)
	})

	t.Run("validator and diff suppression are optional", func(t *testing.T) {
		inner := NoneAllRefsBlockInnerSchema(threeBranchTestConfig())

		elem, ok := inner["things"].Elem.(*schema.Schema)
		require.True(t, ok)
		assert.Nil(t, elem.ValidateDiagFunc)
		assert.Nil(t, inner["things"].DiffSuppressFunc)
	})

	t.Run("diff suppression is passed through to the refs attribute", func(t *testing.T) {
		cfg := threeBranchTestConfig()
		cfg.RefsDiffSuppress = func(string, string, string, *schema.ResourceData) bool { return true }

		inner := NoneAllRefsBlockInnerSchema(cfg)

		assert.NotNil(t, inner["things"].DiffSuppressFunc)
	})
}

func Test_NoneAllRefsBlockConfig_BranchFromConfig(t *testing.T) {
	testCases := []struct {
		Name           string
		Config         NoneAllRefsBlockConfig
		Input          any
		ExpectedBranch NoneAllRefsBranch
		ExpectedRefs   []string
		Error          string
	}{
		{
			Name:           "none",
			Config:         threeBranchTestConfig(),
			Input:          []any{map[string]any{"none": true, "all": false, "things": refsSet()}},
			ExpectedBranch: NoneAllRefsBranchNone,
		},
		{
			Name:           "all",
			Config:         threeBranchTestConfig(),
			Input:          []any{map[string]any{"none": false, "all": true, "things": refsSet()}},
			ExpectedBranch: NoneAllRefsBranchAll,
		},
		{
			Name:           "refs",
			Config:         threeBranchTestConfig(),
			Input:          []any{map[string]any{"none": false, "all": false, "things": refsSet("a", "b")}},
			ExpectedBranch: NoneAllRefsBranchRefs,
			ExpectedRefs:   []string{"a", "b"},
		},
		{
			Name:           "none wins over a populated refs set",
			Config:         threeBranchTestConfig(),
			Input:          []any{map[string]any{"none": true, "all": true, "things": refsSet("a")}},
			ExpectedBranch: NoneAllRefsBranchNone,
		},
		{
			Name:           "all wins over a populated refs set",
			Config:         threeBranchTestConfig(),
			Input:          []any{map[string]any{"none": false, "all": true, "things": refsSet("a")}},
			ExpectedBranch: NoneAllRefsBranchAll,
		},
		{
			Name:           "none is ignored when the branch is not declared",
			Config:         NoneAllRefsBlockConfig{AttrPath: "allowed_things", WithAll: true, RefsKey: "things"},
			Input:          []any{map[string]any{"none": true, "all": false, "things": refsSet("a")}},
			ExpectedBranch: NoneAllRefsBranchRefs,
			ExpectedRefs:   []string{"a"},
		},
		{
			Name:           "all is ignored when the branch is not declared",
			Config:         twoBranchTestConfig(),
			Input:          []any{map[string]any{"none": false, "all": true, "things": refsSet("a")}},
			ExpectedBranch: NoneAllRefsBranchRefs,
			ExpectedRefs:   []string{"a"},
		},
		{
			Name:   "an empty refs set selects no branch",
			Config: threeBranchTestConfig(),
			Input:  []any{map[string]any{"none": false, "all": false, "things": refsSet()}},
			Error:  "allowed_things block has no recognized field set",
		},
		{
			Name:   "no branch set at all",
			Config: threeBranchTestConfig(),
			Input:  []any{map[string]any{}},
			Error:  "allowed_things block has no recognized field set",
		},
		{
			Name:   "an empty block",
			Config: threeBranchTestConfig(),
			Input:  []any{},
			Error:  "allowed_things block is empty",
		},
		{
			Name:   "a value that is not a block",
			Config: threeBranchTestConfig(),
			Input:  "allowed_things",
			Error:  "allowed_things block is empty",
		},
		{
			Name:   "a block whose element is not a map",
			Config: threeBranchTestConfig(),
			Input:  []any{"none"},
			Error:  "allowed_things block has an unexpected shape",
		},
		{
			Name:   "refs that are not a set",
			Config: threeBranchTestConfig(),
			Input:  []any{map[string]any{"things": []any{"a"}}},
			Error:  "allowed_things block has no recognized field set",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			branch, refs, err := tc.Config.BranchFromConfig(tc.Input)

			if tc.Error != "" {
				assert.ErrorContains(t, err, tc.Error)
				assert.Empty(t, branch)
				assert.Nil(t, refs)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.ExpectedBranch, branch)
			assert.ElementsMatch(t, tc.ExpectedRefs, refs)
		})
	}
}

type testRefsRequest struct {
	none bool
	all  bool
	refs []string
}

func newTestRefsRequest() *testRefsRequest {
	return &testRefsRequest{}
}

func parseTestRef(s string) (string, error) {
	if s == "invalid" {
		return "", errors.New("invalid ref")
	}
	return strings.ToUpper(s), nil
}

func testRefsRequestSetters() (func(*testRefsRequest), func(*testRefsRequest), func(*testRefsRequest, []string)) {
	return func(r *testRefsRequest) { r.none = true },
		func(r *testRefsRequest) { r.all = true },
		func(r *testRefsRequest, refs []string) { r.refs = refs }
}

func Test_NoneAllRefsRequest(t *testing.T) {
	setNone, setAll, setRefs := testRefsRequestSetters()

	t.Run("none calls only the none setter", func(t *testing.T) {
		request, err := NoneAllRefsRequest(threeBranchTestConfig(), []any{map[string]any{"none": true}},
			newTestRefsRequest, parseTestRef, setNone, setAll, setRefs)

		require.NoError(t, err)
		assert.Equal(t, testRefsRequest{none: true}, request)
	})

	t.Run("all calls only the all setter", func(t *testing.T) {
		request, err := NoneAllRefsRequest(threeBranchTestConfig(), []any{map[string]any{"all": true}},
			newTestRefsRequest, parseTestRef, setNone, setAll, setRefs)

		require.NoError(t, err)
		assert.Equal(t, testRefsRequest{all: true}, request)
	})

	t.Run("refs are parsed before they are set", func(t *testing.T) {
		request, err := NoneAllRefsRequest(threeBranchTestConfig(), []any{map[string]any{"things": refsSet("a", "b")}},
			newTestRefsRequest, parseTestRef, setNone, setAll, setRefs)

		require.NoError(t, err)
		assert.False(t, request.none)
		assert.False(t, request.all)
		assert.ElementsMatch(t, []string{"A", "B"}, request.refs)
	})

	t.Run("a parse error is returned with a zero request", func(t *testing.T) {
		request, err := NoneAllRefsRequest(threeBranchTestConfig(), []any{map[string]any{"things": refsSet("invalid")}},
			newTestRefsRequest, parseTestRef, setNone, setAll, setRefs)

		assert.ErrorContains(t, err, "invalid ref")
		assert.Equal(t, testRefsRequest{}, request)
	})

	t.Run("a branch resolution error is propagated", func(t *testing.T) {
		request, err := NoneAllRefsRequest(threeBranchTestConfig(), []any{},
			newTestRefsRequest, parseTestRef, setNone, setAll, setRefs)

		assert.ErrorContains(t, err, "allowed_things block is empty")
		assert.Equal(t, testRefsRequest{}, request)
	})

	t.Run("a declared none branch without a setter is an error, not a silent no-op", func(t *testing.T) {
		request, err := NoneAllRefsRequest(threeBranchTestConfig(), []any{map[string]any{"none": true}},
			newTestRefsRequest, parseTestRef, nil, setAll, setRefs)

		assert.ErrorContains(t, err, "allowed_things block resolved to the none branch, but no setter was provided")
		assert.Equal(t, testRefsRequest{}, request)
	})

	t.Run("a declared all branch without a setter is an error, not a silent no-op", func(t *testing.T) {
		request, err := NoneAllRefsRequest(threeBranchTestConfig(), []any{map[string]any{"all": true}},
			newTestRefsRequest, parseTestRef, setNone, nil, setRefs)

		assert.ErrorContains(t, err, "allowed_things block resolved to the all branch, but no setter was provided")
		assert.Equal(t, testRefsRequest{}, request)
	})

	t.Run("a refs branch without a setter is an error, not a panic", func(t *testing.T) {
		request, err := NoneAllRefsRequest(threeBranchTestConfig(), []any{map[string]any{"things": refsSet("a")}},
			newTestRefsRequest, parseTestRef, setNone, setAll, nil)

		assert.ErrorContains(t, err, "allowed_things block resolved to the refs branch, but no setter was provided")
		assert.Equal(t, testRefsRequest{}, request)
	})
}

func Test_NoneAllRefsBlockConfig_StateFromDescribeOutput(t *testing.T) {
	t.Run("an empty or nil describe output is the none branch", func(t *testing.T) {
		for name, items := range map[string][]string{
			"an empty describe output": {},
			"a nil describe output":    nil,
		} {
			t.Run(name, func(t *testing.T) {
				cfg := threeBranchTestConfig()

				block := stateBlock(t, cfg.StateFromDescribeOutput(items))

				assert.Equal(t, map[string]any{"none": true}, block)
			})
		}
	})

	t.Run("a single ALL is the all branch, case-insensitively", func(t *testing.T) {
		for _, item := range []string{"ALL", "all", "All"} {
			t.Run(item, func(t *testing.T) {
				block := stateBlock(t, threeBranchTestConfig().StateFromDescribeOutput([]string{item}))

				assert.Equal(t, map[string]any{"all": true}, block)
			})
		}
	})

	t.Run("other items are the refs branch", func(t *testing.T) {
		block := stateBlock(t, threeBranchTestConfig().StateFromDescribeOutput([]string{"a", "b"}))

		assert.ElementsMatch(t, []string{"a", "b"}, stateRefs(t, block, "things"))
		assert.NotContains(t, block, "none")
		assert.NotContains(t, block, "all")
	})

	t.Run("ALL alongside other items is the refs branch", func(t *testing.T) {
		block := stateBlock(t, threeBranchTestConfig().StateFromDescribeOutput([]string{"ALL", "a"}))

		assert.ElementsMatch(t, []string{"ALL", "a"}, stateRefs(t, block, "things"))
	})

	t.Run("ALL is a ref when the all branch is not declared", func(t *testing.T) {
		block := stateBlock(t, twoBranchTestConfig().StateFromDescribeOutput([]string{"ALL"}))

		assert.ElementsMatch(t, []string{"ALL"}, stateRefs(t, block, "things"))
	})

	t.Run("an empty describe output is the refs branch when the none branch is not declared", func(t *testing.T) {
		cfg := NoneAllRefsBlockConfig{AttrPath: "allowed_things", WithAll: true, RefsKey: "things"}

		block := stateBlock(t, cfg.StateFromDescribeOutput([]string{}))

		assert.Empty(t, stateRefs(t, block, "things"))
	})

	t.Run("the refs key names the attribute", func(t *testing.T) {
		cfg := threeBranchTestConfig()
		cfg.RefsKey = "secrets"

		block := stateBlock(t, cfg.StateFromDescribeOutput([]string{"a"}))

		assert.Contains(t, block, "secrets")
		assert.NotContains(t, block, "things")
	})
}
