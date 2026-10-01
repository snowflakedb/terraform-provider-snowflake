package sdk

import (
	"fmt"
	"regexp"
	"strings"
)

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
	for columnIdx, column := range opts.ColumnsAndConstraints.Columns {
		if column.InlineConstraint != nil {
			if err := column.InlineConstraint.validate(); err != nil {
				errs = append(errs, err)
			}
		}
		for tagIdx, tag := range column.Tag {
			if !ValidObjectIdentifier(tag.Name) {
				errs = append(errs, errInvalidIdentifier(fmt.Sprintf("CreateTableOptions.ColumnsAndConstraints.Columns[%d].Tag[%d]", columnIdx, tagIdx), "Name"))
			}
		}
	}
	for _, outOfLineConstraint := range opts.ColumnsAndConstraints.OutOfLineConstraint {
		if err := outOfLineConstraint.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if opts.RowAccessPolicy != nil {
		if !ValidObjectIdentifier(opts.RowAccessPolicy.Name) {
			errs = append(errs, errInvalidIdentifier("CreateTableOptions.RowAccessPolicy", "Name"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterTableOptions) additionalValidations() error {
	var errs []error
	if constraintAction := opts.ConstraintAction; valueSet(constraintAction) {
		if addAction := constraintAction.Add; valueSet(addAction) {
			if err := addAction.validate(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return JoinErrors(errs...)
}

func (opts *ShowTableOptions) additionalValidations() error {
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		return ErrPatternRequiredForLikeKeyword
	}
	return nil
}

func (v *OutOfLineConstraint) validate() error {
	var errs []error
	switch v.ConstraintType {
	case ColumnConstraintTypeForeignKey:
		if !valueSet(v.ForeignKey) {
			errs = append(errs, errNotSet("OutOfLineConstraint", "ForeignKey"))
		} else {
			if err := v.ForeignKey.validate(); err != nil {
				errs = append(errs, err)
			}
		}
	case ColumnConstraintTypeUnique, ColumnConstraintTypePrimaryKey:
		if valueSet(v.ForeignKey) {
			errs = append(errs, errSet("OutOfLineConstraint", "ForeignKey"))
		}
	default:
		errs = append(errs, errInvalidValue("OutOfLineConstraint", "ConstraintType", string(v.ConstraintType)))
	}
	return JoinErrors(errs...)
}

func (v *OutOfLineForeignKey) validate() error {
	var errs []error
	if !valueSet(v.TableName) {
		errs = append(errs, errNotSet("OutOfLineForeignKey", "TableName"))
	}
	return JoinErrors(errs...)
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
	matches := collateRegexp.FindStringSubmatch(string(r.Type))

	if len(matches) == 2 {
		collation := matches[1]
		type_ := DataType(strings.TrimSpace(collateRegexp.ReplaceAllString(string(r.Type), "")))
		return type_, &collation
	}
	return r.Type, nil
}
