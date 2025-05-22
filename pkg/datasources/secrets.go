package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var secretsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC SECRET for each secret returned by SHOW SECRETS. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"in": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "IN clause to filter the list of secrets",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"account": {
					Type:         schema.TypeBool,
					Optional:     true,
					Description:  "Returns records for the entire account.",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
				},
				"database": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "Returns records for the current database in use or for a specified database.",
					ExactlyOneOf:     []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
					ValidateDiagFunc: resources.IsValidIdentifier[sdk.AccountObjectIdentifier](),
				},
				"schema": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "Returns records for the current schema in use or a specified schema. Use fully qualified name.",
					ExactlyOneOf:     []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
					ValidateDiagFunc: resources.IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
				},
				"application": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "Returns records for the specified application.",
					ExactlyOneOf:     []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
					ValidateDiagFunc: resources.IsValidIdentifier[sdk.AccountObjectIdentifier](),
				},
				"application_package": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "Returns records for the specified application package.",
					ExactlyOneOf:     []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
					ValidateDiagFunc: resources.IsValidIdentifier[sdk.AccountObjectIdentifier](),
				},
			},
		},
	},
	"like": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).",
	},
	"secrets": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all secrets details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW SECRETS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowSecretSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE SECRET.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeSecretSchema,
					},
				},
			},
		},
	},
}

func Secrets() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.Secrets, ReadSecrets),
		Schema:      secretsSchema,
		Description: "Data source used to get details of filtered secrets. Filtering is aligned with the current possibilities for [SHOW SECRETS](https://docs.snowflake.com/en/sql-reference/sql/show-secrets) query. The results of SHOW and DESCRIBE are encapsulated in one output collection `secrets`.",
	}
}

func ReadSecrets(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.NewShowSecretRequest()

	handleLike(d, &req.Like)
	err := handleExtendedIn(d, &req.In)
	if err != nil {
		return diag.FromErr(err)
	}

	secrets, err := client.Secrets.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("secrets_read")

	flattenedSecrets := make([]map[string]any, len(secrets))
	for i, secret := range secrets {
		secret := secret
		var secretDescriptions []map[string]any
		if d.Get("with_describe").(bool) {
			describeOutput, err := client.Secrets.Describe(ctx, secret.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			secretDescriptions = []map[string]any{schemas.SecretDescriptionToSchema(*describeOutput)}
		}

		flattenedSecrets[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.SecretToSchema(&secret)},
			resources.DescribeOutputAttributeName: secretDescriptions,
		}
	}
	if err := d.Set("secrets", flattenedSecrets); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
