package resources

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type NoneAllRefsBlockConfig struct {
	AttrPath string
	WithNone bool
	WithAll  bool
	RefsKey  string

	NoneDescription string
	AllDescription  string
	RefsDescription string

	RefsElemValidator schema.SchemaValidateDiagFunc
	RefsDiffSuppress  schema.SchemaDiffSuppressFunc
}

type NoneAllRefsBranch string

const (
	NoneAllRefsBranchNone NoneAllRefsBranch = "none"
	NoneAllRefsBranchAll  NoneAllRefsBranch = "all"
	NoneAllRefsBranchRefs NoneAllRefsBranch = "refs"
)

func (cfg NoneAllRefsBlockConfig) branchNames() []string {
	branches := make([]string, 0, 3)
	if cfg.WithNone {
		branches = append(branches, "none")
	}
	if cfg.WithAll {
		branches = append(branches, "all")
	}
	return append(branches, cfg.RefsKey)
}

func (cfg NoneAllRefsBlockConfig) exactlyOneOf() []string {
	return collections.Map(cfg.branchNames(), func(branch string) string {
		return fmt.Sprintf("%s.0.%s", cfg.AttrPath, branch)
	})
}

func NoneAllRefsBlockInnerSchema(cfg NoneAllRefsBlockConfig) map[string]*schema.Schema {
	exactlyOneOf := cfg.exactlyOneOf()

	inner := map[string]*schema.Schema{
		cfg.RefsKey: {
			Type:         schema.TypeSet,
			Optional:     true,
			MinItems:     1,
			Description:  cfg.RefsDescription,
			ExactlyOneOf: exactlyOneOf,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: cfg.RefsElemValidator,
			},
			DiffSuppressFunc: cfg.RefsDiffSuppress,
		},
	}
	if cfg.WithNone {
		inner["none"] = &schema.Schema{
			Type:         schema.TypeBool,
			Optional:     true,
			Description:  cfg.NoneDescription,
			ExactlyOneOf: exactlyOneOf,
		}
	}
	if cfg.WithAll {
		inner["all"] = &schema.Schema{
			Type:         schema.TypeBool,
			Optional:     true,
			Description:  cfg.AllDescription,
			ExactlyOneOf: exactlyOneOf,
		}
	}
	return inner
}

func (cfg NoneAllRefsBlockConfig) BranchFromConfig(v any) (NoneAllRefsBranch, []string, error) {
	list, ok := v.([]any)
	if !ok || len(list) == 0 {
		return "", nil, fmt.Errorf("%s block is empty", cfg.AttrPath)
	}
	block, ok := list[0].(map[string]any)
	if !ok {
		return "", nil, fmt.Errorf("%s block has an unexpected shape", cfg.AttrPath)
	}

	if cfg.WithNone {
		if none, _ := block["none"].(bool); none {
			return NoneAllRefsBranchNone, nil, nil
		}
	}
	if cfg.WithAll {
		if all, _ := block["all"].(bool); all {
			return NoneAllRefsBranchAll, nil, nil
		}
	}
	if refs, ok := block[cfg.RefsKey].(*schema.Set); ok && refs.Len() > 0 {
		return NoneAllRefsBranchRefs, expandStringList(refs.List()), nil
	}

	return "", nil, fmt.Errorf("%s block has no recognized field set", cfg.AttrPath)
}

func NoneAllRefsRequest[R, I any](
	cfg NoneAllRefsBlockConfig,
	v any,
	newRequest func() *R,
	parse func(string) (I, error),
	setNone func(*R),
	setAll func(*R),
	setRefs func(*R, []I),
) (R, error) {
	var zero R

	branch, rawRefs, err := cfg.BranchFromConfig(v)
	if err != nil {
		return zero, err
	}

	request := newRequest()
	switch branch {
	case NoneAllRefsBranchNone:
		if setNone == nil {
			return zero, fmt.Errorf("%s block resolved to the none branch, but no setter was provided", cfg.AttrPath)
		}
		setNone(request)
	case NoneAllRefsBranchAll:
		if setAll == nil {
			return zero, fmt.Errorf("%s block resolved to the all branch, but no setter was provided", cfg.AttrPath)
		}
		setAll(request)
	default:
		if setRefs == nil {
			return zero, fmt.Errorf("%s block resolved to the refs branch, but no setter was provided", cfg.AttrPath)
		}
		ids, err := collections.MapErr(rawRefs, parse)
		if err != nil {
			return zero, err
		}
		setRefs(request, ids)
	}
	return *request, nil
}

func (cfg NoneAllRefsBlockConfig) StateFromDescribeOutput(items []string) []any {
	if cfg.WithNone && len(items) == 0 {
		return []any{map[string]any{"none": true}}
	}
	if cfg.WithAll && len(items) == 1 && strings.EqualFold(items[0], "ALL") {
		return []any{map[string]any{"all": true}}
	}

	elems := make([]any, len(items))
	for i, item := range items {
		elems[i] = item
	}
	return []any{map[string]any{
		cfg.RefsKey: schema.NewSet(schema.HashString, elems),
	}}
}
