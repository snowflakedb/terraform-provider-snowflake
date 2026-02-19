package sdk

import (
	"errors"
	"fmt"
	"strings"
)

type TableRowAccessPolicy struct {
	rowAccessPolicy bool                   `ddl:"static" sql:"ROW ACCESS POLICY"`
	Name            SchemaObjectIdentifier `ddl:"identifier"`
	On              []string               `ddl:"keyword,parentheses" sql:"ON"`
}

// ColumnInlineConstraint is based on https://docs.snowflake.com/en/sql-reference/sql/create-table-constraint#inline-unique-primary-foreign-key.
type ColumnInlineConstraint struct {
	Name       *string              `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	Type       ColumnConstraintType `ddl:"keyword"`
	ForeignKey *InlineForeignKey    `ddl:"keyword" sql:"FOREIGN KEY"`

	// optional
	Enforced           *bool `ddl:"keyword" sql:"ENFORCED"`
	NotEnforced        *bool `ddl:"keyword" sql:"NOT ENFORCED"`
	Deferrable         *bool `ddl:"keyword" sql:"DEFERRABLE"`
	NotDeferrable      *bool `ddl:"keyword" sql:"NOT DEFERRABLE"`
	InitiallyDeferred  *bool `ddl:"keyword" sql:"INITIALLY DEFERRED"`
	InitiallyImmediate *bool `ddl:"keyword" sql:"INITIALLY IMMEDIATE"`
	Enable             *bool `ddl:"keyword" sql:"ENABLE"`
	Disable            *bool `ddl:"keyword" sql:"DISABLE"`
	Validate           *bool `ddl:"keyword" sql:"VALIDATE"`
	NoValidate         *bool `ddl:"keyword" sql:"NOVALIDATE"`
	Rely               *bool `ddl:"keyword" sql:"RELY"`
	NoRely             *bool `ddl:"keyword" sql:"NORELY"`
}

func (v *ColumnInlineConstraint) validate() error {
	var errs []error
	switch v.Type {
	case ColumnConstraintTypeForeignKey:
		if !valueSet(v.ForeignKey) {
			errs = append(errs, errNotSet("ColumnInlineConstraint", "ForeignKey"))
		} else {
			if err := v.ForeignKey.validate(); err != nil {
				errs = append(errs, err)
			}
		}
	case ColumnConstraintTypeUnique, ColumnConstraintTypePrimaryKey:
		if valueSet(v.ForeignKey) {
			errs = append(errs, errSet("ColumnInlineConstraint", "ForeignKey"))
		}
	default:
		errs = append(errs, errInvalidValue("ColumnInlineConstraint", "Type", string(v.Type)))
	}
	if moreThanOneValueSet(v.Enforced, v.NotEnforced) {
		errs = append(errs, errMoreThanOneOf("ColumnInlineConstraint", "Enforced", "NotEnforced"))
	}
	if moreThanOneValueSet(v.Deferrable, v.NotDeferrable) {
		errs = append(errs, errMoreThanOneOf("ColumnInlineConstraint", "Deferrable", "NotDeferrable"))
	}
	if moreThanOneValueSet(v.InitiallyDeferred, v.InitiallyImmediate) {
		errs = append(errs, errMoreThanOneOf("ColumnInlineConstraint", "InitiallyDeferred", "InitiallyImmediate"))
	}
	if moreThanOneValueSet(v.Enable, v.Disable) {
		errs = append(errs, errMoreThanOneOf("ColumnInlineConstraint", "Enable", "Disable"))
	}
	if moreThanOneValueSet(v.Validate, v.NoValidate) {
		errs = append(errs, errMoreThanOneOf("ColumnInlineConstraint", "Validate", "Novalidate"))
	}
	if moreThanOneValueSet(v.Rely, v.NoRely) {
		errs = append(errs, errMoreThanOneOf("ColumnInlineConstraint", "Rely", "Norely"))
	}
	return errors.Join(errs...)
}

type ColumnConstraintType string

const (
	ColumnConstraintTypeUnique     ColumnConstraintType = "UNIQUE"
	ColumnConstraintTypePrimaryKey ColumnConstraintType = "PRIMARY KEY"
	ColumnConstraintTypeForeignKey ColumnConstraintType = "FOREIGN KEY"
)

var AllColumnConstraintTypes = []ColumnConstraintType{ColumnConstraintTypeUnique, ColumnConstraintTypePrimaryKey, ColumnConstraintTypeForeignKey}

func ToColumnConstraintType(s string) (ColumnConstraintType, error) {
	cType := strings.ToUpper(s)

	switch cType {
	case string(ColumnConstraintTypeUnique):
		return ColumnConstraintTypeUnique, nil
	case string(ColumnConstraintTypePrimaryKey):
		return ColumnConstraintTypePrimaryKey, nil
	case string(ColumnConstraintTypeForeignKey):
		return ColumnConstraintTypeForeignKey, nil
	}

	return "", fmt.Errorf("invalid column constraint type: %s", s)
}

type InlineForeignKey struct {
	TableName  string              `ddl:"keyword" sql:"REFERENCES"`
	ColumnName []string            `ddl:"keyword,parentheses"`
	Match      *MatchType          `ddl:"keyword" sql:"MATCH"`
	On         *ForeignKeyOnAction `ddl:"keyword" sql:"ON"`
}

func (v *InlineForeignKey) validate() error {
	var errs []error
	if !valueSet(v.TableName) {
		errs = append(errs, errNotSet("InlineForeignKey", "TableName"))
	}
	return errors.Join(errs...)
}

type MatchType string

var (
	FullMatchType    MatchType = "FULL"
	SimpleMatchType  MatchType = "SIMPLE"
	PartialMatchType MatchType = "PARTIAL"
)

var AllMatchTypes = []MatchType{FullMatchType, SimpleMatchType, PartialMatchType}

func ToMatchType(s string) (MatchType, error) {
	cType := strings.ToUpper(s)

	switch cType {
	case string(FullMatchType):
		return FullMatchType, nil
	case string(SimpleMatchType):
		return SimpleMatchType, nil
	case string(PartialMatchType):
		return PartialMatchType, nil
	}

	return "", fmt.Errorf("invalid match type: %s", s)
}

type ForeignKeyOnAction struct {
	OnUpdate *ForeignKeyAction `ddl:"parameter,no_equals" sql:"ON UPDATE"`
	OnDelete *ForeignKeyAction `ddl:"parameter,no_equals" sql:"ON DELETE"`
}

type ForeignKeyAction string

const (
	ForeignKeyCascadeAction    ForeignKeyAction = "CASCADE"
	ForeignKeySetNullAction    ForeignKeyAction = "SET NULL"
	ForeignKeySetDefaultAction ForeignKeyAction = "SET DEFAULT"
	ForeignKeyRestrictAction   ForeignKeyAction = "RESTRICT"
	ForeignKeyNoAction         ForeignKeyAction = "NO ACTION"
)

var AllForeignKeyActions = []ForeignKeyAction{ForeignKeyCascadeAction, ForeignKeySetNullAction, ForeignKeySetDefaultAction, ForeignKeyRestrictAction, ForeignKeyNoAction}

func ToForeignKeyAction(s string) (ForeignKeyAction, error) {
	cType := strings.ToUpper(s)

	switch cType {
	case string(ForeignKeyCascadeAction):
		return ForeignKeyCascadeAction, nil
	case string(ForeignKeySetNullAction):
		return ForeignKeySetNullAction, nil
	case string(ForeignKeySetDefaultAction):
		return ForeignKeySetDefaultAction, nil
	case string(ForeignKeyRestrictAction):
		return ForeignKeyRestrictAction, nil
	case string(ForeignKeyNoAction):
		return ForeignKeyNoAction, nil
	}

	return "", fmt.Errorf("invalid column constraint type: %s", s)
}

func AsStringList[T ~string](input []T) []string {
	output := make([]string, len(input))
	for i, element := range input {
		output[i] = string(element)
	}
	return output
}
