package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var catalogLinkedDatabaseSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the database; must be unique for your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"linked_catalog": {
		Type:        schema.TypeList,
		Required:    true,
		MaxItems:    1,
		Description: "Specifies how the database is linked to a remote catalog.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"catalog": {
					Type:             schema.TypeString,
					Required:         true,
					ForceNew:         true,
					ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
					DiffSuppressFunc: suppressIdentifierQuoting,
					Description:      "Specifies the catalog integration used to connect to the remote catalog. Changing this value recreates the database.",
				},
				"allowed_namespaces": {
					Type:        schema.TypeSet,
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "Specifies the namespaces from the remote catalog allowed to be synced into the database. If both `allowed_namespaces` and `blocked_namespaces` are specified, `blocked_namespaces` takes precedence.",
				},
				"blocked_namespaces": {
					Type:        schema.TypeSet,
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "Specifies the namespaces from the remote catalog blocked from being synced into the database. If both `allowed_namespaces` and `blocked_namespaces` are specified, `blocked_namespaces` takes precedence.",
				},
				"allowed_write_operations": {
					Type:             schema.TypeString,
					Optional:         true,
					Computed:         true,
					ValidateDiagFunc: sdkValidation(sdk.ToCatalogLinkedDatabaseAllowedWriteOperations),
					DiffSuppressFunc: NormalizeAndCompare(sdk.ToCatalogLinkedDatabaseAllowedWriteOperations),
					Description:      fmt.Sprintf("Specifies the write operations allowed on the tables in the database. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllCatalogLinkedDatabaseAllowedWriteOperations)),
				},
				"namespace_mode": {
					Type:             schema.TypeString,
					Optional:         true,
					Computed:         true,
					ForceNew:         true,
					ValidateDiagFunc: sdkValidation(sdk.ToCatalogLinkedDatabaseNamespaceMode),
					DiffSuppressFunc: NormalizeAndCompare(sdk.ToCatalogLinkedDatabaseNamespaceMode),
					Description:      fmt.Sprintf("Specifies how nested namespaces in the remote catalog are handled. Valid values are (case-insensitive): %s. Changing this value recreates the database.", possibleValuesListed(sdk.AllCatalogLinkedDatabaseNamespaceModes)),
				},
				"namespace_flatten_delimiter": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: "Specifies the delimiter used to flatten nested namespaces. Only applicable when `namespace_mode` is `FLATTEN_NESTED_NAMESPACE`. Changing this value recreates the database.",
				},
				"sync_interval_seconds": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
					// Server-side limits can change, so the documented 30-86400 range stays in the description only.
					ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
					Description:      "Specifies how often, in seconds, the database is synced with the remote catalog. Valid values are between `30` and `86400`.",
				},
			},
		},
	},
	"external_volume": {
		Type:             schema.TypeString,
		Optional:         true,
		Computed:         true,
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      relatedResourceDescription("Specifies the external volume where the data of the linked tables is stored. Required unless the catalog vends credentials. Changing this value recreates the database; removing it from the configuration keeps the current volume.", resources.ExternalVolume),
	},
	"catalog_case_sensitivity": {
		Type:             schema.TypeString,
		Optional:         true,
		Computed:         true,
		ForceNew:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToDatabaseCatalogCaseSensitivity),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToDatabaseCatalogCaseSensitivity),
		Description:      fmt.Sprintf("Specifies how the case of the identifiers from the remote catalog is handled. Valid values are (case-insensitive): %s. Changing this value recreates the database.", possibleValuesListed(sdk.AllDatabaseCatalogCaseSensitivities)),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the database.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW DATABASES` for this database.",
		Elem: &schema.Resource{
			Schema: schemas.ShowDatabaseSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func CatalogLinkedDatabase() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.Databases.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.CatalogLinkedDatabaseResource), TrackingCreateWrapper(resources.CatalogLinkedDatabase, CreateCatalogLinkedDatabase)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.CatalogLinkedDatabaseResource), TrackingReadWrapper(resources.CatalogLinkedDatabase, ReadCatalogLinkedDatabase)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.CatalogLinkedDatabaseResource), TrackingUpdateWrapper(resources.CatalogLinkedDatabase, UpdateCatalogLinkedDatabase)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.CatalogLinkedDatabaseResource), TrackingDeleteWrapper(resources.CatalogLinkedDatabase, deleteFunc)),
		Description:   "Resource used to manage catalog-linked databases, i.e. databases synced with an external Iceberg REST catalog. For more information, check [CREATE DATABASE (catalog-linked) documentation](https://docs.snowflake.com/en/sql-reference/sql/create-database-catalog-linked). Note that a reachable external Iceberg REST catalog and a corresponding catalog integration are required.",

		Schema: catalogLinkedDatabaseSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.CatalogLinkedDatabase, ImportName[sdk.AccountObjectIdentifier]),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: TrackingCustomDiffWrapper(resources.CatalogLinkedDatabase, customdiff.All(
			validateNamespaceFlattenDelimiterUsage,
			ComputedIfAnyAttributeChanged(catalogLinkedDatabaseSchema, ShowOutputAttributeName, "name", "comment"),
			ComputedIfAnyAttributeChanged(catalogLinkedDatabaseSchema, FullyQualifiedNameAttributeName, "name"),
		)),
	}
}

// Outside FLATTEN_NESTED_NAMESPACE mode Snowflake reports the delimiter as null, which would
// permanently diff this ForceNew field. HasChange is needed because diff.Get falls back to the
// state value for Computed attributes absent from the configuration.
func validateNamespaceFlattenDelimiterUsage(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	if !diff.HasChange("linked_catalog.0.namespace_flatten_delimiter") {
		return nil
	}
	if diff.Get("linked_catalog.0.namespace_flatten_delimiter").(string) == "" {
		return nil
	}
	namespaceMode, err := sdk.ToCatalogLinkedDatabaseNamespaceMode(diff.Get("linked_catalog.0.namespace_mode").(string))
	if err != nil || namespaceMode != sdk.CatalogLinkedDatabaseNamespaceModeFlattenNestedNamespace {
		return fmt.Errorf("namespace_flatten_delimiter can only be set when namespace_mode is %s", sdk.CatalogLinkedDatabaseNamespaceModeFlattenNestedNamespace)
	}
	return nil
}

func CreateCatalogLinkedDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	linkedCatalogRequest, err := buildLinkedCatalogRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}
	request := sdk.NewCreateCatalogLinkedDatabaseRequest(id).WithLinkedCatalog(*linkedCatalogRequest)

	errs := errors.Join(
		attributeMappedValueCreateBuilder(d, "external_volume", request.WithExternalVolume, sdk.ParseAccountObjectIdentifier),
		attributeMappedValueCreateBuilder(d, "catalog_case_sensitivity", request.WithCatalogCaseSensitivity, sdk.ToDatabaseCatalogCaseSensitivity),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.Databases.CreateCatalogLinked(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating catalog-linked database %v, err = %w", id.FullyQualifiedName(), err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadCatalogLinkedDatabase(ctx, d, meta)
}

func buildLinkedCatalogRequest(d *schema.ResourceData) (*sdk.LinkedCatalogRequest, error) {
	linkedCatalog := d.Get("linked_catalog").([]any)[0].(map[string]any)

	catalogId, err := sdk.ParseAccountObjectIdentifier(linkedCatalog["catalog"].(string))
	if err != nil {
		return nil, err
	}
	request := sdk.NewLinkedCatalogRequest().WithCatalog(catalogId)

	if v := linkedCatalog["allowed_namespaces"].(*schema.Set); v.Len() > 0 {
		request.WithAllowedNamespaces(namespacesToStringListItemWrappers(expandStringList(v.List())))
	}
	if v := linkedCatalog["blocked_namespaces"].(*schema.Set); v.Len() > 0 {
		request.WithBlockedNamespaces(namespacesToStringListItemWrappers(expandStringList(v.List())))
	}
	if v := linkedCatalog["allowed_write_operations"].(string); v != "" {
		allowedWriteOperations, err := sdk.ToCatalogLinkedDatabaseAllowedWriteOperations(v)
		if err != nil {
			return nil, err
		}
		request.WithAllowedWriteOperations(allowedWriteOperations)
	}
	if v := linkedCatalog["namespace_mode"].(string); v != "" {
		namespaceMode, err := sdk.ToCatalogLinkedDatabaseNamespaceMode(v)
		if err != nil {
			return nil, err
		}
		request.WithNamespaceMode(namespaceMode)
	}
	if v := linkedCatalog["namespace_flatten_delimiter"].(string); v != "" {
		request.WithNamespaceFlattenDelimiter(v)
	}
	if v := linkedCatalog["sync_interval_seconds"].(int); v != 0 {
		request.WithSyncIntervalSeconds(v)
	}
	return request, nil
}

func namespacesToStringListItemWrappers(namespaces []string) []sdk.StringListItemWrapper {
	return collections.Map(namespaces, func(namespace string) sdk.StringListItemWrapper {
		return sdk.StringListItemWrapper{Value: namespace}
	})
}

func ReadCatalogLinkedDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	database, err := client.Databases.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query catalog-linked database. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Catalog-linked database id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	if database.Kind == nil || *database.Kind != sdk.DatabaseKindCatalogLinkedDatabase {
		return diag.FromErr(fmt.Errorf("database %s is not a catalog-linked database; use the snowflake_database resource (or the matching variant) to manage it", id.FullyQualifiedName()))
	}

	config, err := client.SystemFunctions.GetCatalogLinkedDatabaseConfig(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	linkedCatalog := map[string]any{
		"catalog":            config.CatalogIntegration.Name(),
		"allowed_namespaces": config.AllowedNamespaces,
		"blocked_namespaces": config.BlockedNamespaces,
	}
	if config.AllowedWriteOperations != nil {
		linkedCatalog["allowed_write_operations"] = string(*config.AllowedWriteOperations)
	}
	if config.NamespaceMode != nil {
		linkedCatalog["namespace_mode"] = string(*config.NamespaceMode)
	}
	if config.NamespaceFlattenDelimiter != nil {
		linkedCatalog["namespace_flatten_delimiter"] = *config.NamespaceFlattenDelimiter
	}
	if config.SyncIntervalSeconds != nil {
		linkedCatalog["sync_interval_seconds"] = *config.SyncIntervalSeconds
	}

	externalVolume := ""
	if config.ExternalVolume != nil {
		externalVolume = config.ExternalVolume.Name()
	}
	catalogCaseSensitivity := ""
	if config.CatalogCaseSensitivity != nil {
		catalogCaseSensitivity = string(*config.CatalogCaseSensitivity)
	}

	errs := errors.Join(
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set("name", database.Name),
		d.Set("linked_catalog", []any{linkedCatalog}),
		d.Set("external_volume", externalVolume),
		d.Set("catalog_case_sensitivity", catalogCaseSensitivity),
		d.Set("comment", database.Comment),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.DatabaseToSchema(database)}),
	)

	return diag.FromErr(errs)
}

func UpdateCatalogLinkedDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		if err := client.Databases.Alter(ctx, sdk.NewAlterDatabaseRequest(id).WithRenameTo(newId)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	if diags := updateCatalogLinkedDatabaseNamespaces(ctx, client, d, id, "linked_catalog.0.allowed_namespaces",
		func(request *sdk.AlterCatalogLinkedDatabaseRequest, namespaces []sdk.StringListItemWrapper) {
			request.WithAddToAllowedNamespaces(*sdk.NewAddToAllowedNamespacesRequest(namespaces))
		},
		func(request *sdk.AlterCatalogLinkedDatabaseRequest, namespaces []sdk.StringListItemWrapper) {
			request.WithRemoveFromAllowedNamespaces(*sdk.NewRemoveFromAllowedNamespacesRequest(namespaces))
		},
		func(request *sdk.AlterCatalogLinkedDatabaseRequest) {
			request.WithUnsetAllowedNamespaces(true)
		},
	); diags != nil {
		return diags
	}

	if diags := updateCatalogLinkedDatabaseNamespaces(ctx, client, d, id, "linked_catalog.0.blocked_namespaces",
		func(request *sdk.AlterCatalogLinkedDatabaseRequest, namespaces []sdk.StringListItemWrapper) {
			request.WithAddToBlockedNamespaces(*sdk.NewAddToBlockedNamespacesRequest(namespaces))
		},
		func(request *sdk.AlterCatalogLinkedDatabaseRequest, namespaces []sdk.StringListItemWrapper) {
			request.WithRemoveFromBlockedNamespaces(*sdk.NewRemoveFromBlockedNamespacesRequest(namespaces))
		},
		func(request *sdk.AlterCatalogLinkedDatabaseRequest) {
			request.WithUnsetBlockedNamespaces(true)
		},
	); diags != nil {
		return diags
	}

	setRequest := sdk.NewCatalogLinkedDatabaseSetRequest()
	if errs := errors.Join(
		intAttributeUpdateSetOnly(d, "linked_catalog.0.sync_interval_seconds", &setRequest.SyncIntervalSeconds),
		attributeMappedValueUpdateSetOnly(d, "linked_catalog.0.allowed_write_operations", &setRequest.AllowedWriteOperations, sdk.ToCatalogLinkedDatabaseAllowedWriteOperations),
	); errs != nil {
		return diag.FromErr(errs)
	}
	if (*setRequest != sdk.CatalogLinkedDatabaseSetRequest{}) {
		if err := client.Databases.AlterCatalogLinked(ctx, sdk.NewAlterCatalogLinkedDatabaseRequest(id).WithSet(*setRequest)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting linked catalog properties for catalog-linked database %v, err = %w", id.FullyQualifiedName(), err))
		}
	}

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		if len(comment) > 0 {
			err = client.Databases.Alter(ctx, sdk.NewAlterDatabaseRequest(id).WithSet(*sdk.NewDatabaseSetRequest().WithComment(comment)))
		} else {
			err = client.Databases.Alter(ctx, sdk.NewAlterDatabaseRequest(id).WithUnset(*sdk.NewDatabaseUnsetRequest().WithComment(true)))
		}
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadCatalogLinkedDatabase(ctx, d, meta)
}

func updateCatalogLinkedDatabaseNamespaces(
	ctx context.Context,
	client *sdk.Client,
	d *schema.ResourceData,
	id sdk.AccountObjectIdentifier,
	attributePath string,
	withAdd func(*sdk.AlterCatalogLinkedDatabaseRequest, []sdk.StringListItemWrapper),
	withRemove func(*sdk.AlterCatalogLinkedDatabaseRequest, []sdk.StringListItemWrapper),
	withUnset func(*sdk.AlterCatalogLinkedDatabaseRequest),
) diag.Diagnostics {
	if !d.HasChange(attributePath) {
		return nil
	}

	oldValue, newValue := d.GetChange(attributePath)
	oldNamespaces, newNamespaces := oldValue.(*schema.Set), newValue.(*schema.Set)

	if newNamespaces.Len() == 0 {
		request := sdk.NewAlterCatalogLinkedDatabaseRequest(id)
		withUnset(request)
		if err := client.Databases.AlterCatalogLinked(ctx, request); err != nil {
			return diag.FromErr(fmt.Errorf("error unsetting namespaces for catalog-linked database %v, err = %w", id.FullyQualifiedName(), err))
		}
		return nil
	}

	// Each ALTER accepts exactly one action.
	if added := newNamespaces.Difference(oldNamespaces); added.Len() > 0 {
		request := sdk.NewAlterCatalogLinkedDatabaseRequest(id)
		withAdd(request, namespacesToStringListItemWrappers(expandStringList(added.List())))
		if err := client.Databases.AlterCatalogLinked(ctx, request); err != nil {
			return diag.FromErr(fmt.Errorf("error adding namespaces for catalog-linked database %v, err = %w", id.FullyQualifiedName(), err))
		}
	}
	if removed := oldNamespaces.Difference(newNamespaces); removed.Len() > 0 {
		request := sdk.NewAlterCatalogLinkedDatabaseRequest(id)
		withRemove(request, namespacesToStringListItemWrappers(expandStringList(removed.List())))
		if err := client.Databases.AlterCatalogLinked(ctx, request); err != nil {
			return diag.FromErr(fmt.Errorf("error removing namespaces for catalog-linked database %v, err = %w", id.FullyQualifiedName(), err))
		}
	}
	return nil
}
