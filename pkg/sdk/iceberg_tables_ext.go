package sdk

import "fmt"

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
	return Pointer(fmt.Sprintf("'%s'", id.FullyQualifiedName()))
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
