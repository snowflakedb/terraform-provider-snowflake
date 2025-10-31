package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var secretGenericStringSchema = func() map[string]*schema.Schema {
	secretGenericString := map[string]*schema.Schema{
		"secret_string": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: externalChangesNotDetectedFieldDescription("Specifies the string to store in the secret. The string can be an API token or a string of sensitive value that can be used in the handler code of a UDF or stored procedure. For details, see [Creating and using an external access integration](https://docs.snowflake.com/en/developer-guide/external-network-access/creating-using-external-network-access). You should not use this property to store any kind of OAuth token; use one of the other secret types for your OAuth use cases."),
		},
	}
	return collections.MergeMaps(secretCommonSchema, secretGenericString)
}()

func SecretWithGenericString() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.SecretWithGenericString, CreateContextSecretWithGenericString),
		ReadContext:   TrackingReadWrapper(resources.SecretWithGenericString, ReadContextSecretWithGenericString),
		UpdateContext: TrackingUpdateWrapper(resources.SecretWithGenericString, UpdateContextSecretWithGenericString),
		DeleteContext: TrackingDeleteWrapper(resources.SecretWithGenericString, DeleteContextSecret),
		Description:   "Resource used to manage secret objects with Generic String. For more information, check [secret documentation](https://docs.snowflake.com/en/sql-reference/sql/create-secret).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.SecretWithGenericString, customdiff.All(
			ComputedIfAnyAttributeChanged(secretGenericStringSchema, ShowOutputAttributeName, "comment"),
			ComputedIfAnyAttributeChanged(secretGenericStringSchema, DescribeOutputAttributeName),
			RecreateWhenSecretTypeChangedExternally(sdk.SecretTypeGenericString),
		)),

		Schema: secretGenericStringSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.SecretWithGenericString, ImportSecretWithGenericString),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportSecretWithGenericString(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	_, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	if err := handleSecretImport(d); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateContextSecretWithGenericString(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName, schemaName, name := d.Get("database").(string), d.Get("schema").(string), d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	secretSting := d.Get("secret_string").(string)

	request := sdk.NewCreateWithGenericStringSecretRequest(id, secretSting)
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	err := client.Secrets.CreateWithGenericString(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextSecretWithGenericString(ctx, d, meta)
}

func ReadContextSecretWithGenericString(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	secret, err := client.Secrets.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query secret with generic string. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Secret with generic string id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to query secret generic string.",
				Detail:   fmt.Sprintf("Secret with generic string id: %s, Err: %s", id.FullyQualifiedName(), err),
			},
		}
	}

	secretDescription, err := client.Secrets.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := handleSecretRead(d, id, secret, secretDescription); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func UpdateContextSecretWithGenericString(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := &sdk.SecretSetRequest{}
	unset := &sdk.SecretUnsetRequest{}
	handleSecretUpdate(d, set, unset)
	setForGenericString := &sdk.SetForGenericStringRequest{}

	if d.HasChange("secret_string") {
		secretString := d.Get("secret_string").(string)
		setForGenericString.WithSecretString(secretString)
	}

	if !reflect.DeepEqual(*setForGenericString, sdk.SetForGenericStringRequest{}) {
		set.WithSetForFlow(sdk.SetForFlowRequest{SetForGenericString: setForGenericString})
	}

	if !reflect.DeepEqual(*set, sdk.SecretSetRequest{}) {
		if err := client.Secrets.Alter(ctx, sdk.NewAlterSecretRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	if !reflect.DeepEqual(*unset, sdk.SecretUnsetRequest{}) {
		if err := client.Secrets.Alter(ctx, sdk.NewAlterSecretRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextSecretWithGenericString(ctx, d, meta)
}
