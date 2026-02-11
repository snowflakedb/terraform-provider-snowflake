package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
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

var storageIntegrationAzureSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: blocklistedCharactersFieldDescription("String that specifies the identifier (i.e. name) for the integration; must be unique in your account."),
	},
	"enabled": {
		Type:     schema.TypeBool,
		Required: true,
		Description: joinWithSpace(
			"Specifies whether this storage integration is available for usage in stages.",
			"`TRUE` allows users to create new stages that reference this integration. Existing stages that reference this integration function normally.",
			"`FALSE` prevents users from creating new stages that reference this integration. Existing stages that reference this integration cannot access the storage location in the stage definition.",
		),
	},
	"storage_allowed_locations": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: "Explicitly limits external stages that use the integration to reference one or more storage locations.",
		MinItems:    1,
	},
	"storage_blocked_locations": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Explicitly prohibits external stages that use the integration from referencing one or more storage locations.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the storage integration.",
	},
	"use_privatelink_endpoint": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("use_privatelink_endpoint"),
		Description:      booleanStringFieldDescription("Specifies whether to use outbound private connectivity to harden the security posture."),
	},
	"azure_tenant_id": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotEmpty,
		Description:  "Specifies the ID for your Office 365 tenant that the allowed and blocked storage accounts belong to.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW STORAGE INTEGRATIONS` for the given storage integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowStorageIntegrationSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE STORAGE INTEGRATION` for the given storage integration.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeStorageIntegrationAzureDetailsSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// TODO [next PR]: react to external provider type change (recreate)
func StorageIntegrationAzure() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.StorageIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.StorageIntegrationAzureResource), TrackingCreateWrapper(resources.StorageIntegrationAzure, CreateStorageIntegrationAzure)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.StorageIntegrationAzureResource), TrackingReadWrapper(resources.StorageIntegrationAzure, GetReadStorageIntegrationAzureFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.StorageIntegrationAzureResource), TrackingUpdateWrapper(resources.StorageIntegrationAzure, UpdateStorageIntegrationAzure)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.StorageIntegrationAzureResource), TrackingDeleteWrapper(resources.StorageIntegrationAzure, deleteFunc)),
		Description:   "Resource used to manage Azure storage integration objects. For more information, check [storage integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-storage-integration).",

		Schema: storageIntegrationAzureSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.StorageIntegrationAzure, ImportStorageIntegrationAzure),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(storageIntegrationAzureSchema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(storageIntegrationAzureSchema, DescribeOutputAttributeName, "enabled", "storage_allowed_locations", "storage_blocked_locations", "comment", "use_privatelink_endpoint", "azure_tenant_id"),
		),
	}
}

func ImportStorageIntegrationAzure(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	azureDetails, err := client.StorageIntegrations.DescribeAzureDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	if azureDetails.Provider != "AZURE" {
		return nil, fmt.Errorf("expected AZURE storage provider got %s", azureDetails.Provider)
	}

	errs := errors.Join(
		d.Set("name", azureDetails.ID().Name()),
		d.Set("use_privatelink_endpoint", booleanStringFromBool(azureDetails.UsePrivatelinkEndpoint)),
	)
	if errs != nil {
		return nil, errs
	}
	return []*schema.ResourceData{d}, nil
}

func GetReadStorageIntegrationAzureFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		s, err := client.StorageIntegrations.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query azure storage integration. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Azure storage integration id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		azureDetails, err := client.StorageIntegrations.DescribeAzureDetails(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not describe azure storage integration (%s), err = %w", d.Id(), err))
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInFlatDescribe(d,
				outputMapping{"use_privatelink_endpoint", "use_privatelink_endpoint", azureDetails.UsePrivatelinkEndpoint, booleanStringFromBool(azureDetails.UsePrivatelinkEndpoint), nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		errs := errors.Join(
			// not reading name on purpose (we never update the name externally)
			d.Set("enabled", s.Enabled),
			d.Set("storage_allowed_locations", azureDetails.AllowedLocations),
			d.Set("storage_blocked_locations", azureDetails.BlockedLocations),
			d.Set("comment", s.Comment),
			// not reading use_privatelink_endpoint on purpose (handled as external change to describe output)
			d.Set("azure_tenant_id", azureDetails.TenantId),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		)

		errs = errors.Join(errs,
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.StorageIntegrationToSchema(s)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.StorageIntegrationAzureDetailsToSchema(azureDetails)}),
		)

		return diag.FromErr(errs)
	}
}

func CreateStorageIntegrationAzure(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(name)
	enabled := d.Get("enabled").(bool)
	storageAllowedLocations, _ := parseLocations(d.Get("storage_allowed_locations").(*schema.Set).List())
	azureTenantId := d.Get("azure_tenant_id").(string)

	request := sdk.NewCreateStorageIntegrationRequest(id, enabled, storageAllowedLocations)
	azureRequest := sdk.NewAzureStorageParamsRequest(azureTenantId)
	errs := errors.Join(
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		func() error {
			if _, ok := d.GetOk("storage_blocked_locations"); ok {
				storageBlockedLocations, _ := parseLocations(d.Get("storage_blocked_locations").(*schema.Set).List())
				request.WithStorageBlockedLocations(storageBlockedLocations)
			}
			return nil
		}(),
		booleanStringAttributeCreateBuilder(d, "use_privatelink_endpoint", azureRequest.WithUsePrivatelinkEndpoint),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.StorageIntegrations.Create(ctx, request.WithAzureStorageProviderParams(*azureRequest)); err != nil {
		return diag.FromErr(fmt.Errorf("error creating storage integration azure: %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return GetReadStorageIntegrationAzureFunc(false)(ctx, d, meta)
}

func UpdateStorageIntegrationAzure(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set, unset := sdk.NewStorageIntegrationSetRequest(), sdk.NewStorageIntegrationUnsetRequest()
	azureSetParams := sdk.NewSetAzureStorageParamsRequest()

	errs := errors.Join(
		booleanAttributeUpdateSetOnly(d, "enabled", &set.Enabled),
		// TODO [next PRs]: extract helpers for lists with builders
		func() error {
			if d.HasChange("storage_allowed_locations") {
				if v, ok := d.GetOk("storage_allowed_locations"); ok {
					locations, err := parseLocations(v.(*schema.Set).List())
					if err != nil {
						return err
					}
					set.WithStorageAllowedLocations(locations)
				}
			}
			return nil
		}(),
		func() error {
			if d.HasChange("storage_blocked_locations") {
				v := d.Get("storage_blocked_locations")
				if len(v.(*schema.Set).List()) > 0 {
					locations, err := parseLocations(v.(*schema.Set).List())
					if err != nil {
						return err
					}
					set.WithStorageBlockedLocations(locations)
				} else {
					unset.WithStorageBlockedLocations(true)
				}
			}
			return nil
		}(),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
		// TODO [SNOW-2356049]: implement & use booleanStringAttributeUnsetBuilder when UNSET starts working correctly
		booleanStringAttributeUnsetFallbackUpdateBuilder(d, "use_privatelink_endpoint", azureSetParams.WithUsePrivatelinkEndpoint, false),
		stringAttributeUpdateSetOnlyNotEmpty(d, "azure_tenant_id", &azureSetParams.AzureTenantId),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(*azureSetParams, *sdk.NewSetAzureStorageParamsRequest()) {
		set.WithAzureParams(*azureSetParams)
	}
	if !reflect.DeepEqual(*set, *sdk.NewStorageIntegrationSetRequest()) {
		req := sdk.NewAlterStorageIntegrationRequest(id).WithSet(*set)
		if err := client.StorageIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating azure storage integration, err = %w", err))
		}
	}

	if !reflect.DeepEqual(*unset, *sdk.NewStorageIntegrationUnsetRequest()) {
		req := sdk.NewAlterStorageIntegrationRequest(id).WithUnset(*unset)
		if err := client.StorageIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating azure storage integration, err = %w", err))
		}
	}

	return GetReadStorageIntegrationAzureFunc(false)(ctx, d, meta)
}
