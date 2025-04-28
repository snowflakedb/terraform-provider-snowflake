package resources

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/tracking"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	FullyQualifiedNameAttributeName = "fully_qualified_name"
	AtAttributeName                 = "at"
	BeforeAttributeName             = "before"
)

// DiffSuppressStatement will suppress diffs between statements if they differ in only case or in
// runs of whitespace (\s+ = \s). This is needed because the snowflake api does not faithfully
// round-trip queries, so we cannot do a simple character-wise comparison to detect changes.
//
// Warnings: We will have false positives in cases where a change in case or run of whitespace is
// semantically significant.
//
// If we can find a sql parser that can handle the snowflake dialect then we should switch to parsing
// queries and either comparing ASTs or emitting a canonical serialization for comparison. I couldn't
// find such a library.
func DiffSuppressStatement(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(normalizeQuery(old), normalizeQuery(new))
}

var space = regexp.MustCompile(`\s+`)

func normalizeQuery(str string) string {
	return strings.TrimSpace(space.ReplaceAllString(str, " "))
}

// TODO [SNOW-1325214]: address during stage resource rework
func suppressQuoting(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	if oldValue == "" || newValue == "" {
		return false
	} else {
		oldWithoutQuotes := strings.ReplaceAll(oldValue, "'", "")
		newWithoutQuotes := strings.ReplaceAll(newValue, "'", "")
		return oldWithoutQuotes == newWithoutQuotes
	}
}

func ctyValToSliceString(valueElems []cty.Value) []string {
	elems := make([]string, len(valueElems))
	for i, v := range valueElems {
		elems[i] = v.AsString()
	}
	return elems
}

func ImportName[T sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier | sdk.AccountIdentifier](ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	switch any(new(T)).(type) {
	case *sdk.AccountObjectIdentifier:
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return nil, err
		}

		if err := d.Set("name", id.Name()); err != nil {
			return nil, err
		}
	case *sdk.DatabaseObjectIdentifier:
		id, err := sdk.ParseDatabaseObjectIdentifier(d.Id())
		if err != nil {
			return nil, err
		}

		if err := d.Set("name", id.Name()); err != nil {
			return nil, err
		}

		if err := d.Set("database", id.DatabaseName()); err != nil {
			return nil, err
		}
	case *sdk.SchemaObjectIdentifier:
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return nil, err
		}

		if err := d.Set("name", id.Name()); err != nil {
			return nil, err
		}

		if err := d.Set("database", id.DatabaseName()); err != nil {
			return nil, err
		}

		if err := d.Set("schema", id.SchemaName()); err != nil {
			return nil, err
		}
	case *sdk.AccountIdentifier:
		id, err := sdk.ParseAccountIdentifier(d.Id())
		if err != nil {
			return nil, err
		}

		if err := d.Set("name", id.AccountName()); err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
}

func TrackingImportWrapper(resourceName resources.Resource, importImplementation schema.StateContextFunc) schema.StateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
		ctx = tracking.NewContext(ctx, tracking.NewVersionedResourceMetadata(resourceName, tracking.ImportOperation))
		return importImplementation(ctx, d, meta)
	}
}

func TrackingCreateWrapper(resourceName resources.Resource, createImplementation schema.CreateContextFunc) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		ctx = tracking.NewContext(ctx, tracking.NewVersionedResourceMetadata(resourceName, tracking.CreateOperation))
		return createImplementation(ctx, d, meta)
	}
}

func TrackingReadWrapper(resourceName resources.Resource, readImplementation schema.ReadContextFunc) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		ctx = tracking.NewContext(ctx, tracking.NewVersionedResourceMetadata(resourceName, tracking.ReadOperation))
		return readImplementation(ctx, d, meta)
	}
}

func TrackingUpdateWrapper(resourceName resources.Resource, updateImplementation schema.UpdateContextFunc) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		ctx = tracking.NewContext(ctx, tracking.NewVersionedResourceMetadata(resourceName, tracking.UpdateOperation))
		return updateImplementation(ctx, d, meta)
	}
}

func TrackingDeleteWrapper(resourceName resources.Resource, deleteImplementation schema.DeleteContextFunc) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		ctx = tracking.NewContext(ctx, tracking.NewVersionedResourceMetadata(resourceName, tracking.DeleteOperation))
		return deleteImplementation(ctx, d, meta)
	}
}

func TrackingCustomDiffWrapper(resourceName resources.Resource, customdiffImplementation schema.CustomizeDiffFunc) schema.CustomizeDiffFunc {
	return func(ctx context.Context, diff *schema.ResourceDiff, meta any) error {
		ctx = tracking.NewContext(ctx, tracking.NewVersionedResourceMetadata(resourceName, tracking.CustomDiffOperation))
		return customdiffImplementation(ctx, diff, meta)
	}
}

func ensureResourceIsEnabled(featureRaw string, meta any) error {
	enabled := meta.(*provider.Context).EnabledFeatures
	feature, err := previewfeatures.StringToFeature(featureRaw)
	if err != nil {
		return err
	}
	if err := previewfeatures.EnsurePreviewFeatureEnabled(feature, enabled); err != nil {
		return err
	}
	return nil
}

func PreviewFeatureCreateContextWrapper(featureRaw string, createFunc schema.CreateContextFunc) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		if err := ensureResourceIsEnabled(featureRaw, meta); err != nil {
			return diag.FromErr(err)
		}
		return createFunc(ctx, d, meta)
	}
}

func PreviewFeatureReadContextWrapper(featureRaw string, readFunc schema.ReadContextFunc) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		if err := ensureResourceIsEnabled(featureRaw, meta); err != nil {
			return diag.FromErr(err)
		}
		return readFunc(ctx, d, meta)
	}
}

func PreviewFeatureUpdateContextWrapper(featureRaw string, updateFunc schema.UpdateContextFunc) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		if err := ensureResourceIsEnabled(featureRaw, meta); err != nil {
			return diag.FromErr(err)
		}
		return updateFunc(ctx, d, meta)
	}
}

func PreviewFeatureDeleteContextWrapper(featureRaw string, deleteFunc schema.DeleteContextFunc) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		if err := ensureResourceIsEnabled(featureRaw, meta); err != nil {
			return diag.FromErr(err)
		}
		return deleteFunc(ctx, d, meta)
	}
}

const (
	defaultCreateTimeout = 20 * time.Minute
	defaultReadTimeout   = 20 * time.Minute
	defaultUpdateTimeout = 20 * time.Minute
	defaultDeleteTimeout = 20 * time.Minute
)

var defaultTimeouts = &schema.ResourceTimeout{
	Create: schema.DefaultTimeout(defaultCreateTimeout),
	Read:   schema.DefaultTimeout(defaultReadTimeout),
	Update: schema.DefaultTimeout(defaultUpdateTimeout),
	Delete: schema.DefaultTimeout(defaultDeleteTimeout),
}
