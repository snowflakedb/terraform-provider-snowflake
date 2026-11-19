package sdk

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

func (v *tables) ShowParameters(ctx context.Context, id SchemaObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Table: id,
		},
	})
}

func (r *CreateTableRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

// SequenceName is a SET DEFAULT value that renders as <name>.NEXTVAL.
// OptionalEnumLegacy does not emit the type, so it lives here next to String().
type SequenceName string

func (sn SequenceName) String() string {
	return fmt.Sprintf("%s.NEXTVAL", string(sn))
}

// GetClusterByKeys converts the SHOW TABLES result for ClusterBy and converts it to list of keys.
func (v *Table) GetClusterByKeys() []string {
	if v.ClusterBy == "" {
		return nil
	}

	statementWithoutLinear := strings.TrimSuffix(strings.Replace(v.ClusterBy, "LINEAR(", "", 1), ")")
	return splitOuter(statementWithoutLinear, '(', ')')
}

func (opts *CreateTableOptions) additionalValidations() error {
	var errs []error
	if len(opts.ColumnsAndConstraints.Columns) == 0 {
		errs = append(errs, errNotSet("CreateTableOptions", "Columns"))
	}
	if opts.RowAccessPolicy != nil {
		if !ValidObjectIdentifier(opts.RowAccessPolicy.Name) {
			errs = append(errs, errInvalidIdentifier("CreateTableOptions.RowAccessPolicy", "Name"))
		}
	}
	return JoinErrors(errs...)
}

// The ADD branch of ALTER TABLE ... ADD COLUMN for external tables is the only ExternalTableAction branch, so it is required.
// ChildrenRule needs two alternatives, so the SSOT definition cannot express this.
func (opts *AlterTableOptions) additionalValidations() error {
	if valueSet(opts.ExternalTableAction) && !valueSet(opts.ExternalTableAction.Add) {
		return errNotSet("AlterTableOptions.ExternalTableAction", "Add")
	}
	return nil
}

func (opts *ShowTableOptions) additionalValidations() error {
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		return ErrPatternRequiredForLikeKeyword
	}
	return nil
}

// additionalConvert splits COLLATE off the DESCRIBE TABLE type column.
func (r tableColumnDetailsRow) additionalConvert(result *TableColumnDetails) error {
	type_, collation := r.splitTypeAndCollation()
	result.Type = type_
	result.Collation = collation
	return nil
}

func (r tableColumnDetailsRow) splitTypeAndCollation() (DataType, *string) {
	collateRegexp := regexp.MustCompile(`COLLATE +'([a-zA-Z0-9_-]*)'`)
	matches := collateRegexp.FindStringSubmatch(r.Type)

	if len(matches) == 2 {
		collation := matches[1]
		type_ := DataType(strings.TrimSpace(collateRegexp.ReplaceAllString(r.Type, "")))
		return type_, &collation
	}
	return DataType(r.Type), nil
}
