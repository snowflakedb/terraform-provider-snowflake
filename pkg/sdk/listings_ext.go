package sdk

import "context"

func (v *listings) dropSafelyHook(ctx context.Context, id AccountObjectIdentifier) error {
	if l, err := v.ShowByIDSafely(ctx, id); err == nil {
		if l.State == ListingStatePublished {
			return v.Alter(ctx, NewAlterListingRequest(id).WithIfExists(true).WithUnpublish(true))
		}
	}
	return nil
}
