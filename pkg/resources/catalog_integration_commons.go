package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func catalogIntegrationCommonSchema(describeSchema map[string]*schema.Schema) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Specifies the identifier (i.e. name) of the catalog integration; must be unique in your account.",
		},
		"enabled": {
			Type:     schema.TypeBool,
			Required: true,
			Description: joinWithSpace("Specifies whether the catalog integration is available for use for Iceberg tables.",
				"`true` allows users to create new Iceberg tables that reference this integration. Existing Iceberg tables that reference this integration function normally.",
				"`false` prevents users from creating new Iceberg tables that reference this integration. Existing Iceberg tables that reference this integration cannot access the catalog in the table definition."),
		},
		"refresh_interval_seconds": {
			Type:     schema.TypeInt,
			Optional: true,
			Description: joinWithSpace("Specifies the number of seconds to wait between attempts to poll the external Iceberg catalog for metadata updates for automated refresh.",
				"For Delta-based tables, specifies the number of seconds to wait between attempts to poll your external cloud storage for new metadata."),
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("refresh_interval_seconds"),
			ValidateFunc:     validation.IntAtLeast(1),
		},
		"comment": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "Specifies a comment for the catalog integration.",
		},
		ShowOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `SHOW CATALOG INTEGRATIONS` for the given catalog integration.",
			Elem: &schema.Resource{
				Schema: schemas.ShowCatalogIntegrationSchema,
			},
		},
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE CATALOG INTEGRATION` for the given catalog integration.",
			Elem: &schema.Resource{
				Schema: describeSchema,
			},
		},
		FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	}
}

func handleCatalogIntegrationUpdate(ctx context.Context, d *schema.ResourceData, meta any) error {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return err
	}

	set := sdk.NewCatalogIntegrationSetRequest()

	errs := errors.Join(
		booleanAttributeUpdateSetOnly(d, "enabled", &set.Enabled),
		// TODO [SNOW-3243983]: UNSET not implemented
		intAttributeUnsetFallbackUpdateWithZeroDefault(d, "refresh_interval_seconds", &set.RefreshIntervalSeconds, 30),
		stringAttributeUpdateSetOnly(d, "comment", &set.Comment),
	)
	if errs != nil {
		return errs
	}

	if !reflect.DeepEqual(*set, *sdk.NewCatalogIntegrationSetRequest()) {
		req := sdk.NewAlterCatalogIntegrationRequest(id).WithSet(*set)
		if err := client.CatalogIntegrations.Alter(ctx, req); err != nil {
			return fmt.Errorf("error updating catalog integration (%s), err = %w", d.Id(), err)
		}
	}
	return nil
}

func deleteCatalogIntegrationFunc() schema.DeleteContextFunc {
	return ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.CatalogIntegrations.DropSafely
		},
	)
}
