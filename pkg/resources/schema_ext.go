package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Constructor
// =============================================================================

var (
	schemaCreateExt = TrackingCreateWrapper(resources.Schema, CreateSchema)
	schemaReadExt   = TrackingReadWrapper(resources.Schema, ReadSchemaFunc(true))
	schemaUpdateExt = TrackingUpdateWrapper(resources.Schema, UpdateSchema)
	schemaDeleteExt = TrackingDeleteWrapper(resources.Schema, ResourceDeleteContextFunc(
		schemaParseIdExt,
		schemaDropSafelyInSdkExt,
	))
	schemaDescriptionExt   = "Resource used to manage schema objects. For more information, check [schema documentation](https://docs.snowflake.com/en/sql-reference/sql/create-schema)."
	schemaCustomizeDiffExt = TrackingCustomDiffWrapper(resources.Schema, customdiff.All(
		TemporaryWorkaroundIdentifierForceNewIfHierarchyRenamesExperimentNotEnabled("database"),
		ComputedIfAnyAttributeChanged(schemaSchema, ShowOutputAttributeName, "name", "database", "comment", "with_managed_access", "is_transient"),
		ComputedIfAnyAttributeChanged(schemaSchema, DescribeOutputAttributeName, "name", "database"),
		ComputedIfAnyAttributeChanged(schemaSchema, FullyQualifiedNameAttributeName, "name", "database"),
		ComputedIfAnyAttributeChanged(schemaParametersAttributesSchema, ParametersAttributeName, schemaParametersFieldNames...),
		schemaParametersCustomDiff,
	))
	schemaImporterExt = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.Schema, ImportSchema),
	}
	schemaSchemaVersionExt  = 2
	schemaStateUpgradersExt = []schema.StateUpgrader{
		{
			Version: 0,
			// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
			Type:    cty.EmptyObject,
			Upgrade: v093SchemaStateUpgrader,
		},
		{
			Version: 1,
			// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
			Type:    cty.EmptyObject,
			Upgrade: migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName,
		},
	}
)

// =============================================================================
// Request / mapping / Read-Update-Import helpers
// =============================================================================

var schemaSetStateToValueFromConfigFieldsExt = []string{
	"is_transient",
	"with_managed_access",
}

type schemaAlterRequestsExt struct {
	set   *sdk.SchemaSetRequest
	unset *sdk.SchemaUnsetRequest
}

func schemaParseIdFromConfigExt(d *schema.ResourceData) (sdk.DatabaseObjectIdentifier, error) {
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	return sdk.NewDatabaseObjectIdentifier(database, name), nil
}

func schemaParseIdExt(id string) (sdk.DatabaseObjectIdentifier, error) {
	return sdk.ParseDatabaseObjectIdentifier(id)
}

func schemaBeforeCreateExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.DatabaseObjectIdentifier) (done bool, diags diag.Diagnostics) {
	name := d.Get("name").(string)
	if !strings.EqualFold(strings.TrimSpace(name), "PUBLIC") {
		return false, nil
	}

	client := meta.(*provider.Context).Client
	_, err := client.Schemas.ShowByID(ctx, id)
	if err != nil && !errors.Is(err, sdk.ErrObjectNotFound) {
		return true, diag.FromErr(err)
	} else if err == nil {
		// there is already a PUBLIC schema, so we need to alter it instead
		log.Printf("[DEBUG] found PUBLIC schema during creation, updating...")
		d.SetId(helpers.EncodeResourceIdentifier(id))
		return true, UpdateSchema(ctx, d, meta)
	}
	return false, nil
}

func schemaNewCreateRequestExt(d *schema.ResourceData, id sdk.DatabaseObjectIdentifier) (*sdk.CreateSchemaRequest, error) {
	return sdk.NewCreateSchemaRequest(id), nil
}

func schemaApplyCreateOptionalsExt(d *schema.ResourceData, request *sdk.CreateSchemaRequest) error {
	if v := GetConfigPropertyAsPointerAllowingZeroValue[string](d, "comment"); v != nil {
		request.WithComment(*v)
	}
	if v := d.Get("is_transient").(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return err
		}
		request.WithTransient(parsed)
	}
	if v := d.Get("with_managed_access").(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return err
		}
		request.WithWithManagedAccess(parsed)
	}
	return nil
}

func schemaHandleExternalChangesExt(d *schema.ResourceData, s *sdk.Schema, withExternalChangesMarking bool) error {
	if !withExternalChangesMarking {
		return nil
	}
	return handleExternalChangesToObjectInShow(d, schemaOutputMappingsExt(s)...)
}

func schemaOutputMappingsExt(s *sdk.Schema) []outputMapping {
	return []outputMapping{
		{"options", "is_transient", s.IsTransient(), booleanStringFromBool(s.IsTransient()), func(x any) any {
			return slices.Contains(sdk.ParseCommaSeparatedStringArray(x.(string), false), "TRANSIENT")
		}},
		{"options", "with_managed_access", s.IsManagedAccess(), booleanStringFromBool(s.IsManagedAccess()), func(x any) any {
			return slices.Contains(sdk.ParseCommaSeparatedStringArray(x.(string), false), "MANAGED ACCESS")
		}},
	}
}

func schemaSetConfigFieldsExt(d *schema.ResourceData, s *sdk.Schema) error {
	return d.Set("comment", s.Comment)
}

func schemaSetDescribeOutputExt(d *schema.ResourceData, id sdk.DatabaseObjectIdentifier, details []sdk.SchemaDetails, describeErr error) error {
	if describeErr != nil {
		log.Printf("[DEBUG] describing schema: %s, err: %s", id.FullyQualifiedName(), describeErr)
		return nil
	}
	return d.Set(DescribeOutputAttributeName, schemas.SchemaDetailsListToSchema(details))
}

func schemaHierarchyRenameExt(ctx context.Context, d *schema.ResourceData, meta any, id *sdk.DatabaseObjectIdentifier) diag.Diagnostics {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client

	if providerCtx.Experiments.IsEnabled(experimentalfeatures.HierarchyRenames) && d.HasChange("database") {
		schemaRenameFn := func(currentId, targetId sdk.DatabaseObjectIdentifier) func() error {
			return func() error {
				return client.Schemas.Alter(ctx, sdk.NewAlterSchemaRequest(currentId).WithRenameTo(targetId))
			}
		}

		if diags := handleTwoLevelHierarchyRename(
			ctx, d, client, id,
			schemaRenameFn,
			client.Schemas.ShowByID,
			func(id sdk.DatabaseObjectIdentifier) string { return helpers.EncodeResourceIdentifier(id) },
			"schema",
		); diags != nil {
			return diags
		}
	}

	if d.HasChange("name") && !d.GetRawState().IsNull() {
		newId := sdk.NewDatabaseObjectIdentifier(d.Get("database").(string), d.Get("name").(string))
		err := client.Schemas.Alter(ctx, sdk.NewAlterSchemaRequest(*id).WithRenameTo(newId))
		if err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
		d.SetId(helpers.EncodeResourceIdentifier(newId))
		*id = newId
	}
	return nil
}

func schemaBeforeAlterExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.DatabaseObjectIdentifier) diag.Diagnostics {
	if !d.HasChange("with_managed_access") {
		return nil
	}

	client := meta.(*provider.Context).Client
	if v := d.Get("with_managed_access").(string); v != BooleanDefault {
		var err error
		parsed, err := booleanStringToBool(v)
		if err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
		if parsed {
			err = client.Schemas.Alter(ctx, sdk.NewAlterSchemaRequest(id).WithEnableManagedAccess(true))
		} else {
			err = client.Schemas.Alter(ctx, sdk.NewAlterSchemaRequest(id).WithDisableManagedAccess(true))
		}
		if err != nil {
			d.Partial(true)
			return diag.FromErr(fmt.Errorf("error handling with_managed_access on %v err = %w", d.Id(), err))
		}
	} else {
		// managed access can not be UNSET to a default value
		if err := client.Schemas.Alter(ctx, sdk.NewAlterSchemaRequest(id).WithDisableManagedAccess(true)); err != nil {
			d.Partial(true)
			return diag.FromErr(fmt.Errorf("error handling with_managed_access on %v err = %w", d.Id(), err))
		}
	}
	return nil
}

func schemaApplySetUnsetExt(d *schema.ResourceData) (*schemaAlterRequestsExt, error) {
	set, unset := sdk.NewSchemaSetRequest(), sdk.NewSchemaUnsetRequest()
	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		if len(comment) > 0 {
			set.Comment = &comment
		} else {
			unset.Comment = sdk.Bool(true)
		}
	}
	return &schemaAlterRequestsExt{set: set, unset: unset}, nil
}

func schemaSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.DatabaseObjectIdentifier, s *sdk.Schema) error {
	return errors.Join(
		d.Set("name", id.Name()),
		d.Set("database", id.DatabaseName()),
		d.Set("comment", s.Comment),
		d.Set("is_transient", booleanStringFromBool(s.IsTransient())),
		d.Set("with_managed_access", booleanStringFromBool(s.IsManagedAccess())),
	)
}

// =============================================================================
// SDK incision points
// =============================================================================

func schemaCreateInSdkExt(ctx context.Context, meta any, req *sdk.CreateSchemaRequest) error {
	client := meta.(*provider.Context).Client
	return client.Schemas.Create(ctx, req)
}

func schemaShowByIdSafelyInSdkExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.DatabaseObjectIdentifier) (*sdk.Schema, diag.Diagnostics) {
	client := meta.(*provider.Context).Client
	s, err := client.Schemas.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query schema. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Schema id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}
	return s, nil
}

func schemaShowByIdInSdkExt(ctx context.Context, meta any, id sdk.DatabaseObjectIdentifier) (*sdk.Schema, error) {
	client := meta.(*provider.Context).Client
	return client.Schemas.ShowByID(ctx, id)
}

func schemaDescribeInSdkExt(ctx context.Context, meta any, id sdk.DatabaseObjectIdentifier) ([]sdk.SchemaDetails, error) {
	client := meta.(*provider.Context).Client
	return client.Schemas.Describe(ctx, id)
}

func schemaAlterSetInSdkExt(ctx context.Context, meta any, id sdk.DatabaseObjectIdentifier, reqs *schemaAlterRequestsExt) error {
	if *reqs.set != *sdk.NewSchemaSetRequest() {
		client := meta.(*provider.Context).Client
		return client.Schemas.Alter(ctx, sdk.NewAlterSchemaRequest(id).WithSet(*reqs.set))
	}
	return nil
}

func schemaAlterUnsetInSdkExt(ctx context.Context, meta any, id sdk.DatabaseObjectIdentifier, reqs *schemaAlterRequestsExt) error {
	if *reqs.unset != *sdk.NewSchemaUnsetRequest() {
		client := meta.(*provider.Context).Client
		return client.Schemas.Alter(ctx, sdk.NewAlterSchemaRequest(id).WithUnset(*reqs.unset))
	}
	return nil
}

func schemaDropSafelyInSdkExt(client *sdk.Client) DropSafelyFunc[sdk.DatabaseObjectIdentifier] {
	return client.Schemas.DropSafely
}
