package sdk

import (
	"context"
	"fmt"
)

// TODO [next PRs]: define these validations too (adjust generator if needed)
func (opts *CreateIcebergTableOptions) additionalValidations() error {
	var errs []error
	for i, col := range opts.ColumnsAndConstraints.Columns {
		if valueSet(col.InlineConstraint) {
			path := fmt.Sprintf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[%d].InlineConstraint", i)
			if !exactlyOneValueSet(col.InlineConstraint.UniquePK, col.InlineConstraint.FK, col.InlineConstraint.CH) {
				errs = append(errs, errExactlyOneOf(path, "UniquePK", "FK", "CH"))
			}
			if valueSet(col.InlineConstraint.UniquePK) {
				upPath := path + ".UniquePK"
				if everyValueSet(col.InlineConstraint.UniquePK.Enforced, col.InlineConstraint.UniquePK.NotEnforced) {
					errs = append(errs, errOneOf(upPath, "Enforced", "NotEnforced"))
				}
				if everyValueSet(col.InlineConstraint.UniquePK.Deferrable, col.InlineConstraint.UniquePK.NotDeferrable) {
					errs = append(errs, errOneOf(upPath, "Deferrable", "NotDeferrable"))
				}
				if everyValueSet(col.InlineConstraint.UniquePK.InitiallyDeferred, col.InlineConstraint.UniquePK.InitiallyImmediate) {
					errs = append(errs, errOneOf(upPath, "InitiallyDeferred", "InitiallyImmediate"))
				}
				if everyValueSet(col.InlineConstraint.UniquePK.Enable, col.InlineConstraint.UniquePK.Disable) {
					errs = append(errs, errOneOf(upPath, "Enable", "Disable"))
				}
				if everyValueSet(col.InlineConstraint.UniquePK.Validate, col.InlineConstraint.UniquePK.Novalidate) {
					errs = append(errs, errOneOf(upPath, "Validate", "Novalidate"))
				}
				if everyValueSet(col.InlineConstraint.UniquePK.Rely, col.InlineConstraint.UniquePK.Norely) {
					errs = append(errs, errOneOf(upPath, "Rely", "Norely"))
				}
				if !exactlyOneValueSet(col.InlineConstraint.UniquePK.Unique, col.InlineConstraint.UniquePK.PrimaryKey) {
					errs = append(errs, errExactlyOneOf(upPath, "Unique", "PrimaryKey"))
				}
			}
			if valueSet(col.InlineConstraint.FK) {
				fkPath := path + ".FK"
				if everyValueSet(col.InlineConstraint.FK.Enforced, col.InlineConstraint.FK.NotEnforced) {
					errs = append(errs, errOneOf(fkPath, "Enforced", "NotEnforced"))
				}
				if everyValueSet(col.InlineConstraint.FK.Deferrable, col.InlineConstraint.FK.NotDeferrable) {
					errs = append(errs, errOneOf(fkPath, "Deferrable", "NotDeferrable"))
				}
				if everyValueSet(col.InlineConstraint.FK.InitiallyDeferred, col.InlineConstraint.FK.InitiallyImmediate) {
					errs = append(errs, errOneOf(fkPath, "InitiallyDeferred", "InitiallyImmediate"))
				}
				if everyValueSet(col.InlineConstraint.FK.Enable, col.InlineConstraint.FK.Disable) {
					errs = append(errs, errOneOf(fkPath, "Enable", "Disable"))
				}
				if everyValueSet(col.InlineConstraint.FK.Validate, col.InlineConstraint.FK.Novalidate) {
					errs = append(errs, errOneOf(fkPath, "Validate", "Novalidate"))
				}
				if everyValueSet(col.InlineConstraint.FK.Rely, col.InlineConstraint.FK.Norely) {
					errs = append(errs, errOneOf(fkPath, "Rely", "Norely"))
				}
				if !ValidObjectIdentifier(col.InlineConstraint.FK.References) {
					errs = append(errs, ErrInvalidObjectIdentifier)
				}
			}
			if valueSet(col.InlineConstraint.CH) {
				chPath := path + ".CH"
				if everyValueSet(col.InlineConstraint.CH.EnableValidate, col.InlineConstraint.CH.EnableNovalidate) {
					errs = append(errs, errOneOf(chPath, "EnableValidate", "EnableNovalidate"))
				}
			}
		}
	}
	// Adjusted manually: OutOfLineConstraint is a slice, validate each element
	for i, oc := range opts.ColumnsAndConstraints.OutOfLineConstraint {
		path := fmt.Sprintf("CreateIcebergTableOptions.ColumnsAndConstraints.OutOfLineConstraint[%d]", i)
		if !exactlyOneValueSet(oc.UniquePK, oc.FK, oc.CH) {
			errs = append(errs, errExactlyOneOf(path, "UniquePK", "FK", "CH"))
		}
		if valueSet(oc.UniquePK) {
			upPath := path + ".UniquePK"
			if everyValueSet(oc.UniquePK.Enforced, oc.UniquePK.NotEnforced) {
				errs = append(errs, errOneOf(upPath, "Enforced", "NotEnforced"))
			}
			if everyValueSet(oc.UniquePK.Deferrable, oc.UniquePK.NotDeferrable) {
				errs = append(errs, errOneOf(upPath, "Deferrable", "NotDeferrable"))
			}
			if everyValueSet(oc.UniquePK.InitiallyDeferred, oc.UniquePK.InitiallyImmediate) {
				errs = append(errs, errOneOf(upPath, "InitiallyDeferred", "InitiallyImmediate"))
			}
			if everyValueSet(oc.UniquePK.Enable, oc.UniquePK.Disable) {
				errs = append(errs, errOneOf(upPath, "Enable", "Disable"))
			}
			if everyValueSet(oc.UniquePK.Validate, oc.UniquePK.Novalidate) {
				errs = append(errs, errOneOf(upPath, "Validate", "Novalidate"))
			}
			if everyValueSet(oc.UniquePK.Rely, oc.UniquePK.Norely) {
				errs = append(errs, errOneOf(upPath, "Rely", "Norely"))
			}
			if !exactlyOneValueSet(oc.UniquePK.Unique, oc.UniquePK.PrimaryKey) {
				errs = append(errs, errExactlyOneOf(upPath, "Unique", "PrimaryKey"))
			}
		}
		if valueSet(oc.FK) {
			fkPath := path + ".FK"
			if everyValueSet(oc.FK.Enforced, oc.FK.NotEnforced) {
				errs = append(errs, errOneOf(fkPath, "Enforced", "NotEnforced"))
			}
			if everyValueSet(oc.FK.Deferrable, oc.FK.NotDeferrable) {
				errs = append(errs, errOneOf(fkPath, "Deferrable", "NotDeferrable"))
			}
			if everyValueSet(oc.FK.InitiallyDeferred, oc.FK.InitiallyImmediate) {
				errs = append(errs, errOneOf(fkPath, "InitiallyDeferred", "InitiallyImmediate"))
			}
			if everyValueSet(oc.FK.Enable, oc.FK.Disable) {
				errs = append(errs, errOneOf(fkPath, "Enable", "Disable"))
			}
			if everyValueSet(oc.FK.Validate, oc.FK.Novalidate) {
				errs = append(errs, errOneOf(fkPath, "Validate", "Novalidate"))
			}
			if everyValueSet(oc.FK.Rely, oc.FK.Norely) {
				errs = append(errs, errOneOf(fkPath, "Rely", "Norely"))
			}
			if !ValidObjectIdentifier(oc.FK.References) {
				errs = append(errs, ErrInvalidObjectIdentifier)
			}
		}
		if valueSet(oc.CH) {
			chPath := path + ".CH"
			if everyValueSet(oc.CH.EnableValidate, oc.CH.EnableNovalidate) {
				errs = append(errs, errOneOf(chPath, "EnableValidate", "EnableNovalidate"))
			}
		}
	}
	return JoinErrors(errs...)
}

func (opts *AlterIcebergTableOptions) additionalValidations() error {
	var errs []error
	if valueSet(opts.SearchOptimizationAction) {
		// Adjusted manually: each Drop.On entry must have exactly one of SearchMethodWithTarget, ColumnName, or ExpressionId set.
		if valueSet(opts.SearchOptimizationAction.Drop) {
			for _, on := range opts.SearchOptimizationAction.Drop.On {
				if !exactlyOneValueSet(on.SearchMethodWithTarget, on.ColumnName, on.ExpressionId) {
					errs = append(errs, errExactlyOneOf("AlterIcebergTableOptions.SearchOptimizationAction.Drop.On", "SearchMethodWithTarget", "ColumnName", "ExpressionId"))
				}
			}
		}
	}
	if valueSet(opts.AddColumnAction) {
		if valueSet(opts.AddColumnAction.InlineConstraint) {
			if !exactlyOneValueSet(opts.AddColumnAction.InlineConstraint.UniquePK, opts.AddColumnAction.InlineConstraint.FK, opts.AddColumnAction.InlineConstraint.CH) {
				errs = append(errs, errExactlyOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint", "UniquePK", "FK", "CH"))
			}
			if valueSet(opts.AddColumnAction.InlineConstraint.UniquePK) {
				if everyValueSet(opts.AddColumnAction.InlineConstraint.UniquePK.Enforced, opts.AddColumnAction.InlineConstraint.UniquePK.NotEnforced) {
					errs = append(errs, errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.UniquePK", "Enforced", "NotEnforced"))
				}
				if everyValueSet(opts.AddColumnAction.InlineConstraint.UniquePK.Deferrable, opts.AddColumnAction.InlineConstraint.UniquePK.NotDeferrable) {
					errs = append(errs, errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.UniquePK", "Deferrable", "NotDeferrable"))
				}
				if everyValueSet(opts.AddColumnAction.InlineConstraint.UniquePK.InitiallyDeferred, opts.AddColumnAction.InlineConstraint.UniquePK.InitiallyImmediate) {
					errs = append(errs, errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.UniquePK", "InitiallyDeferred", "InitiallyImmediate"))
				}
				if everyValueSet(opts.AddColumnAction.InlineConstraint.UniquePK.Enable, opts.AddColumnAction.InlineConstraint.UniquePK.Disable) {
					errs = append(errs, errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.UniquePK", "Enable", "Disable"))
				}
				if everyValueSet(opts.AddColumnAction.InlineConstraint.UniquePK.Validate, opts.AddColumnAction.InlineConstraint.UniquePK.Novalidate) {
					errs = append(errs, errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.UniquePK", "Validate", "Novalidate"))
				}
				if everyValueSet(opts.AddColumnAction.InlineConstraint.UniquePK.Rely, opts.AddColumnAction.InlineConstraint.UniquePK.Norely) {
					errs = append(errs, errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.UniquePK", "Rely", "Norely"))
				}
				if !exactlyOneValueSet(opts.AddColumnAction.InlineConstraint.UniquePK.Unique, opts.AddColumnAction.InlineConstraint.UniquePK.PrimaryKey) {
					errs = append(errs, errExactlyOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.UniquePK", "Unique", "PrimaryKey"))
				}
			}
			if valueSet(opts.AddColumnAction.InlineConstraint.FK) {
				if everyValueSet(opts.AddColumnAction.InlineConstraint.FK.Enforced, opts.AddColumnAction.InlineConstraint.FK.NotEnforced) {
					errs = append(errs, errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.FK", "Enforced", "NotEnforced"))
				}
				if everyValueSet(opts.AddColumnAction.InlineConstraint.FK.Deferrable, opts.AddColumnAction.InlineConstraint.FK.NotDeferrable) {
					errs = append(errs, errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.FK", "Deferrable", "NotDeferrable"))
				}
				if everyValueSet(opts.AddColumnAction.InlineConstraint.FK.InitiallyDeferred, opts.AddColumnAction.InlineConstraint.FK.InitiallyImmediate) {
					errs = append(errs, errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.FK", "InitiallyDeferred", "InitiallyImmediate"))
				}
				if everyValueSet(opts.AddColumnAction.InlineConstraint.FK.Enable, opts.AddColumnAction.InlineConstraint.FK.Disable) {
					errs = append(errs, errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.FK", "Enable", "Disable"))
				}
				if everyValueSet(opts.AddColumnAction.InlineConstraint.FK.Validate, opts.AddColumnAction.InlineConstraint.FK.Novalidate) {
					errs = append(errs, errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.FK", "Validate", "Novalidate"))
				}
				if everyValueSet(opts.AddColumnAction.InlineConstraint.FK.Rely, opts.AddColumnAction.InlineConstraint.FK.Norely) {
					errs = append(errs, errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.FK", "Rely", "Norely"))
				}
				if !ValidObjectIdentifier(opts.AddColumnAction.InlineConstraint.FK.References) {
					errs = append(errs, ErrInvalidObjectIdentifier)
				}
			}
			if valueSet(opts.AddColumnAction.InlineConstraint.CH) {
				if everyValueSet(opts.AddColumnAction.InlineConstraint.CH.EnableValidate, opts.AddColumnAction.InlineConstraint.CH.EnableNovalidate) {
					errs = append(errs, errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.CH", "EnableValidate", "EnableNovalidate"))
				}
			}
		}
	}
	return JoinErrors(errs...)
}

func (r *CreateIcebergTableRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

// icebergTableExternalVolumeQuoted formats an AccountObjectIdentifier for the
// EXTERNAL_VOLUME clause of CREATE ICEBERG TABLE, which expects a single-quoted
// string literal whose content is the double-quoted volume name (e.g. '"vol1"').
//
// TODO(SNOW-2236323): Use a proper generation option instead.
// We need to use a custom parsing here, see SNOW-1833593 for more details.
func icebergTableExternalVolumeQuoted(id *AccountObjectIdentifier) *string {
	if id == nil {
		return nil
	}
	return new(fmt.Sprintf("'%s'", id.FullyQualifiedName()))
}

// TableColumnInlineConstraintFromRequest converts an inline constraint
// Request into its Options counterpart for emission as DDL.
func TableColumnInlineConstraintFromRequest(r *TableColumnInlineConstraintRequest) *TableColumnInlineConstraint {
	if r == nil {
		return nil
	}
	out := &TableColumnInlineConstraint{}
	if r.UniquePK != nil {
		out.UniquePK = &TableColumnInlineUniquePK{
			Name:               r.UniquePK.Name,
			Unique:             r.UniquePK.Unique,
			PrimaryKey:         r.UniquePK.PrimaryKey,
			Enforced:           r.UniquePK.Enforced,
			NotEnforced:        r.UniquePK.NotEnforced,
			Deferrable:         r.UniquePK.Deferrable,
			NotDeferrable:      r.UniquePK.NotDeferrable,
			InitiallyDeferred:  r.UniquePK.InitiallyDeferred,
			InitiallyImmediate: r.UniquePK.InitiallyImmediate,
			Enable:             r.UniquePK.Enable,
			Disable:            r.UniquePK.Disable,
			Validate:           r.UniquePK.Validate,
			Novalidate:         r.UniquePK.Novalidate,
			Rely:               r.UniquePK.Rely,
			Norely:             r.UniquePK.Norely,
		}
	}
	if r.FK != nil {
		out.FK = &TableColumnInlineFK{
			Name:               r.FK.Name,
			ForeignKey:         r.FK.ForeignKey,
			References:         r.FK.References,
			RefColumn:          r.FK.RefColumn,
			Match:              r.FK.Match,
			On:                 r.FK.On,
			Enforced:           r.FK.Enforced,
			NotEnforced:        r.FK.NotEnforced,
			Deferrable:         r.FK.Deferrable,
			NotDeferrable:      r.FK.NotDeferrable,
			InitiallyDeferred:  r.FK.InitiallyDeferred,
			InitiallyImmediate: r.FK.InitiallyImmediate,
			Enable:             r.FK.Enable,
			Disable:            r.FK.Disable,
			Validate:           r.FK.Validate,
			Novalidate:         r.FK.Novalidate,
			Rely:               r.FK.Rely,
			Norely:             r.FK.Norely,
		}
	}
	if r.CH != nil {
		out.CH = &TableColumnInlineCH{
			Name:             r.CH.Name,
			Expression:       r.CH.Expression,
			EnableValidate:   r.CH.EnableValidate,
			EnableNovalidate: r.CH.EnableNovalidate,
		}
	}
	return out
}

// TableOutOfLineConstraintFromRequest converts an out-of-line constraint
// Request into its Options counterpart for emission as DDL.
func TableOutOfLineConstraintFromRequest(r *TableOutOfLineConstraintRequest) *TableOutOfLineConstraint {
	if r == nil {
		return nil
	}
	out := &TableOutOfLineConstraint{}
	if r.UniquePK != nil {
		out.UniquePK = &TableOutOfLineUniquePK{
			Name:               r.UniquePK.Name,
			Unique:             r.UniquePK.Unique,
			PrimaryKey:         r.UniquePK.PrimaryKey,
			Columns:            r.UniquePK.Columns,
			Enforced:           r.UniquePK.Enforced,
			NotEnforced:        r.UniquePK.NotEnforced,
			Deferrable:         r.UniquePK.Deferrable,
			NotDeferrable:      r.UniquePK.NotDeferrable,
			InitiallyDeferred:  r.UniquePK.InitiallyDeferred,
			InitiallyImmediate: r.UniquePK.InitiallyImmediate,
			Enable:             r.UniquePK.Enable,
			Disable:            r.UniquePK.Disable,
			Validate:           r.UniquePK.Validate,
			Novalidate:         r.UniquePK.Novalidate,
			Rely:               r.UniquePK.Rely,
			Norely:             r.UniquePK.Norely,
			Comment:            r.UniquePK.Comment,
		}
	}
	if r.FK != nil {
		out.FK = &TableOutOfLineFK{
			Name:               r.FK.Name,
			Columns:            r.FK.Columns,
			References:         r.FK.References,
			RefColumns:         r.FK.RefColumns,
			Match:              r.FK.Match,
			On:                 r.FK.On,
			Enforced:           r.FK.Enforced,
			NotEnforced:        r.FK.NotEnforced,
			Deferrable:         r.FK.Deferrable,
			NotDeferrable:      r.FK.NotDeferrable,
			InitiallyDeferred:  r.FK.InitiallyDeferred,
			InitiallyImmediate: r.FK.InitiallyImmediate,
			Enable:             r.FK.Enable,
			Disable:            r.FK.Disable,
			Validate:           r.FK.Validate,
			Novalidate:         r.FK.Novalidate,
			Rely:               r.FK.Rely,
			Norely:             r.FK.Norely,
			Comment:            r.FK.Comment,
		}
	}
	if r.CH != nil {
		out.CH = &TableOutOfLineCH{
			Name:             r.CH.Name,
			Expression:       r.CH.Expression,
			EnableValidate:   r.CH.EnableValidate,
			EnableNovalidate: r.CH.EnableNovalidate,
		}
	}
	return out
}

func (v *icebergTables) ShowParameters(ctx context.Context, id SchemaObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Table: id,
		},
	})
}

type IcebergTablePartitionSpec struct {
	SpecId int                              `json:"spec-id"`
	Fields []IcebergTablePartitionSpecField `json:"fields"`
}

type IcebergTablePartitionSpecField struct {
	Name      string `json:"name"`
	Transform string `json:"transform"`
	SourceId  int    `json:"source-id"`
	FieldId   int    `json:"field-id"`
}
