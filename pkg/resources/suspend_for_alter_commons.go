package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	tfschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SuspensionForAlter keeps an object suspended for the alter statements of a generated Update (see SuspendForAlter in
// pkg/internal/ssot). The update restores the object before its final Read, and the deferred CleanUp resumes it when the
// update failed before that.
type SuspensionForAlter struct {
	id sdk.ObjectIdentifier
	// resume brings the object back to the state it had before the update. It is nil when nothing was suspended.
	resume   func(context.Context) error
	restored bool
}

// SuspendForAlter calls the SDK SuspendIfStarted of the object when one of the attributes changed. The returned suspension
// is never nil, also with an error, because the object may already be suspended then.
func SuspendForAlter[ID sdk.ObjectIdentifier](
	ctx context.Context,
	d *tfschema.ResourceData,
	suspendIfStarted func(context.Context, ID) (func(context.Context) error, error),
	id ID,
	attributes ...string,
) (*SuspensionForAlter, error) {
	if !d.HasChanges(attributes...) {
		return &SuspensionForAlter{}, nil
	}
	resume, err := suspendIfStarted(ctx, id)
	if resume == nil {
		resume = func(context.Context) error { return nil }
	}
	return &SuspensionForAlter{id: id, resume: resume}, err
}

// SuspendForAlterIf is SuspendForAlter for an object whose resource also decides whether the change needs the
// suspension (SuspendForAlter.NeededExt). Without that, nothing is suspended.
func SuspendForAlterIf[ID sdk.ObjectIdentifier](
	ctx context.Context,
	d *tfschema.ResourceData,
	needed bool,
	suspendIfStarted func(context.Context, ID) (func(context.Context) error, error),
	id ID,
	attributes ...string,
) (*SuspensionForAlter, error) {
	if !needed {
		return &SuspensionForAlter{}, nil
	}
	return SuspendForAlter(ctx, d, suspendIfStarted, id, attributes...)
}

// CleanUp is deferred right after SuspendForAlter. When the update failed, it marks the state partial, so Terraform keeps
// the old values, and resumes the object unless it was already restored. A failed resume is added next to the original
// error.
func (s *SuspensionForAlter) CleanUp(ctx context.Context, d *tfschema.ResourceData, diags *diag.Diagnostics) {
	if s.resume == nil {
		return
	}
	if diags.HasError() {
		d.Partial(true)
	}
	if s.restored {
		return
	}
	if err := s.resume(ctx); err != nil {
		*diags = append(*diags, diag.Diagnostic{Severity: diag.Error, Summary: fmt.Sprintf("failed to resume %s after a failed update: %s", s.id.FullyQualifiedName(), err.Error())})
	}
}

// Restore brings the object back to the state it had before the update. A failure is a warning, because the update
// succeeded.
func (s *SuspensionForAlter) Restore(ctx context.Context) diag.Diagnostics {
	return s.restoreWith(ctx, s.resume)
}

// RestoreIfEnabled resumes the object when enabled is true and leaves it suspended otherwise, for an object whose enabled
// attribute owns its state. A failure is a warning, because the update succeeded.
func (s *SuspensionForAlter) RestoreIfEnabled(ctx context.Context, enabled bool, resume func(context.Context) error) diag.Diagnostics {
	if !enabled {
		resume = func(context.Context) error { return nil }
	}
	return s.restoreWith(ctx, resume)
}

func (s *SuspensionForAlter) restoreWith(ctx context.Context, resume func(context.Context) error) diag.Diagnostics {
	if s.resume == nil {
		return nil
	}
	s.restored = true
	if err := resume(ctx); err != nil {
		return diag.Diagnostics{{Severity: diag.Warning, Summary: fmt.Sprintf("failed to resume %s: %s", s.id.FullyQualifiedName(), err.Error())}}
	}
	return nil
}
