package sdk

import (
	"context"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ NetworkRules = (*networkRules)(nil)
var _ convertibleRow[NetworkRule] = new(ShowNetworkRulesRow)
var _ convertibleRow[NetworkRuleDetails] = new(DescNetworkRulesRow)

type networkRules struct {
	client *Client
}

func (v *networkRules) Create(ctx context.Context, request *CreateNetworkRuleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *networkRules) Alter(ctx context.Context, request *AlterNetworkRuleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *networkRules) Drop(ctx context.Context, request *DropNetworkRuleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *networkRules) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropNetworkRuleRequest(id).WithIfExists(Bool(true))) }, ctx, id)
}

func (v *networkRules) Show(ctx context.Context, request *ShowNetworkRuleRequest) ([]NetworkRule, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[ShowNetworkRulesRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[ShowNetworkRulesRow, NetworkRule](dbRows)
}

func (v *networkRules) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*NetworkRule, error) {
	request := NewShowNetworkRuleRequest().
		WithIn(In{Schema: id.SchemaId()}).
		WithLike(Like{Pattern: String(id.Name())})
	networkRules, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(networkRules, func(r NetworkRule) bool { return r.Name == id.Name() })
}

func (v *networkRules) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*NetworkRule, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *networkRules) Describe(ctx context.Context, id SchemaObjectIdentifier) (*NetworkRuleDetails, error) {
	opts := &DescribeNetworkRuleOptions{
		name: id,
	}
	result, err := validateAndQueryOne[DescNetworkRulesRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return conversionErrorWrapped(result.convert())
}

func (r *CreateNetworkRuleRequest) toOpts() *CreateNetworkRuleOptions {
	opts := &CreateNetworkRuleOptions{
		OrReplace: r.OrReplace,
		name:      r.name,
		Type:      r.Type,
		ValueList: r.ValueList,
		Mode:      r.Mode,
		Comment:   r.Comment,
	}
	return opts
}

func (r *AlterNetworkRuleRequest) toOpts() *AlterNetworkRuleOptions {
	opts := &AlterNetworkRuleOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	if r.Set != nil {
		opts.Set = &NetworkRuleSet{
			ValueList: r.Set.ValueList,
			Comment:   r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &NetworkRuleUnset{
			ValueList: r.Unset.ValueList,
			Comment:   r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropNetworkRuleRequest) toOpts() *DropNetworkRuleOptions {
	opts := &DropNetworkRuleOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowNetworkRuleRequest) toOpts() *ShowNetworkRuleOptions {
	opts := &ShowNetworkRuleOptions{
		Like:       r.Like,
		In:         r.In,
		StartsWith: r.StartsWith,
		Limit:      r.Limit,
	}
	return opts
}

func (row ShowNetworkRulesRow) convert() (*NetworkRule, error) {
	return &NetworkRule{
		CreatedOn:          row.CreatedOn,
		Name:               row.Name,
		DatabaseName:       row.DatabaseName,
		SchemaName:         row.SchemaName,
		Owner:              row.Owner,
		Comment:            row.Comment,
		Type:               NetworkRuleType(row.Type),
		Mode:               NetworkRuleMode(row.Mode),
		EntriesInValueList: row.EntriesInValueList,
		OwnerRoleType:      row.OwnerRoleType,
	}, nil
}

func (r *DescribeNetworkRuleRequest) toOpts() *DescribeNetworkRuleOptions {
	opts := &DescribeNetworkRuleOptions{
		name: r.name,
	}
	return opts
}

func (row DescNetworkRulesRow) convert() (*NetworkRuleDetails, error) {
	valueList := strings.Split(row.ValueList, ",")
	if len(valueList) == 1 && valueList[0] == "" {
		valueList = []string{}
	}
	return &NetworkRuleDetails{
		CreatedOn:    row.CreatedOn,
		Name:         row.Name,
		DatabaseName: row.DatabaseName,
		SchemaName:   row.SchemaName,
		Owner:        row.Owner,
		Comment:      row.Comment,
		Type:         NetworkRuleType(row.Type),
		Mode:         NetworkRuleMode(row.Mode),
		ValueList:    valueList,
	}, nil
}
