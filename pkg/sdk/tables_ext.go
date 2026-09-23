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
	for _, column := range opts.ColumnsAndConstraints.Columns {
		if column.InlineConstraint != nil {
			if err := column.InlineConstraint.validate(); err != nil {
				errs = append(errs, err)
			}
		}
		if column.DefaultValue != nil {
			if ok := exactlyOneValueSet(
				column.DefaultValue.Expression,
				column.DefaultValue.Identity,
			); !ok {
				errs = append(errs, errExactlyOneOf("DefaultValue", "Expression", "Identity"))
			}
			if identity := column.DefaultValue.Identity; valueSet(identity) {
				if moreThanOneValueSet(identity.Order, identity.Noorder) {
					errs = append(errs, errMoreThanOneOf("Identity", "Order", "Noorder"))
				}
			}
		}
		if column.MaskingPolicy != nil {
			if !ValidObjectIdentifier(column.MaskingPolicy.Name) {
				errs = append(errs, errInvalidIdentifier("ColumnMaskingPolicy", "Name"))
			}
		}
		for _, tag := range column.Tag {
			if !ValidObjectIdentifier(tag.Name) {
				errs = append(errs, errInvalidIdentifier("TagAssociation", "Name"))
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
			errs = append(errs, errInvalidIdentifier("TableRowAccessPolicy", "Name"))
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
	if len(v.Columns) == 0 {
		errs = append(errs, errNotSet("OutOfLineConstraint", "Columns"))
	}
	if moreThanOneValueSet(v.Enforced, v.NotEnforced) {
		errs = append(errs, errMoreThanOneOf("OutOfLineConstraint", "Enforced", "NotEnforced"))
	}
	if moreThanOneValueSet(v.Deferrable, v.NotDeferrable) {
		errs = append(errs, errMoreThanOneOf("OutOfLineConstraint", "Deferrable", "NotDeferrable"))
	}
	if moreThanOneValueSet(v.InitiallyDeferred, v.InitiallyImmediate) {
		errs = append(errs, errMoreThanOneOf("OutOfLineConstraint", "InitiallyDeferred", "InitiallyImmediate"))
	}
	if moreThanOneValueSet(v.Enable, v.Disable) {
		errs = append(errs, errMoreThanOneOf("OutOfLineConstraint", "Enable", "Disable"))
	}
	if moreThanOneValueSet(v.Validate, v.Novalidate) {
		errs = append(errs, errMoreThanOneOf("OutOfLineConstraint", "Validate", "Novalidate"))
	}
	if moreThanOneValueSet(v.Rely, v.Norely) {
		errs = append(errs, errMoreThanOneOf("OutOfLineConstraint", "Rely", "Norely"))
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
