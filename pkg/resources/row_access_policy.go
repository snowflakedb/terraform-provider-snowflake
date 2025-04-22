package resources

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var rowAccessPolicySchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the row access policy; must be unique for the database and schema in which the row access policy is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the row access policy."),
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the row access policy."),
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"argument": {
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The argument name",
					ForceNew:    true,
				},
				// TODO(SNOW-1596962): Fully support VECTOR data type sdk.ParseFunctionArgumentsFromString could be a base for another function that takes argument names into consideration.
				"type": {
					Type:             schema.TypeString,
					Required:         true,
					DiffSuppressFunc: DiffSuppressDataTypes,
					ValidateDiagFunc: IsDataTypeValid,
					Description:      dataTypeFieldDescription("The argument type. VECTOR data types are not yet supported."),
					ForceNew:         true,
				},
			},
		},
		Required:    true,
		Description: "List of the arguments for the row access policy. A signature specifies a set of attributes that must be considered to determine whether the row is accessible. The attribute values come from the database object (e.g. table or view) to be protected by the row access policy. If any argument name or type is changed, the resource is recreated.",
		ForceNew:    true,
	},
	"body": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      diffSuppressStatementFieldDescription("Specifies the SQL expression. The expression can be any boolean-valued SQL expression."),
		DiffSuppressFunc: DiffSuppressStatement,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the row access policy.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW ROW ACCESS POLICIES` for the given row access policy.",
		Elem: &schema.Resource{
			Schema: schemas.ShowRowAccessPolicySchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE ROW ACCESS POLICY` for the given row access policy.",
		Elem: &schema.Resource{
			Schema: schemas.RowAccessPolicyDescribeSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// RowAccessPolicy returns a pointer to the resource representing a row access policy.
func RowAccessPolicy() *schema.Resource {
	// TODO(SNOW-1818849): unassign policies before dropping
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.RowAccessPolicies.DropSafely
		},
	)

	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: TrackingCreateWrapper(resources.RowAccessPolicy, CreateRowAccessPolicy),
		ReadContext:   TrackingReadWrapper(resources.RowAccessPolicy, ReadRowAccessPolicy),
		UpdateContext: TrackingUpdateWrapper(resources.RowAccessPolicy, UpdateRowAccessPolicy),
		DeleteContext: TrackingDeleteWrapper(resources.RowAccessPolicy, deleteFunc),
		Description:   "Resource used to manage row access policy objects. For more information, check [row access policy documentation](https://docs.snowflake.com/en/sql-reference/sql/create-row-access-policy).",

		Schema: rowAccessPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.RowAccessPolicy, ImportRowAccessPolicy),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.RowAccessPolicy, customdiff.All(
			ComputedIfAnyAttributeChanged(rowAccessPolicySchema, ShowOutputAttributeName, "comment", "name"),
			ComputedIfAnyAttributeChanged(rowAccessPolicySchema, DescribeOutputAttributeName, "body", "name", "signature"),
			ComputedIfAnyAttributeChanged(rowAccessPolicySchema, FullyQualifiedNameAttributeName, "name"),
		)),

		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v0_95_0_RowAccessPolicyStateUpgrader,
			},
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportRowAccessPolicy(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] Starting row access policy import")
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	policy, err := client.RowAccessPolicies.ShowByID(ctx, id)
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
	if err := d.Set("comment", policy.Comment); err != nil {
		return nil, err
	}
	policyDescription, err := client.RowAccessPolicies.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := d.Set("body", policyDescription.Body); err != nil {
		return nil, err
	}
	if err := d.Set("argument", schemas.RowAccessPolicyArgumentsToSchema(policyDescription.Signature)); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func CreateRowAccessPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	arguments := d.Get("argument").([]any)
	rowAccessExpression := d.Get("body").(string)

	args := make([]sdk.CreateRowAccessPolicyArgsRequest, 0)
	for _, arg := range arguments {
		v := arg.(map[string]any)
		dataType, err := datatypes.ParseDataType(v["type"].(string))
		if err != nil {
			return diag.FromErr(err)
		}
		args = append(args, *sdk.NewCreateRowAccessPolicyArgsRequest(v["name"].(string), dataType))
	}

	createRequest := sdk.NewCreateRowAccessPolicyRequest(id, args, rowAccessExpression)

	// Set optionals
	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(sdk.String(v.(string)))
	}

	err := client.RowAccessPolicies.Create(ctx, createRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating row access policy %v err = %w", name, err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadRowAccessPolicy(ctx, d, meta)
}

func ReadRowAccessPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	rowAccessPolicy, err := client.RowAccessPolicies.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query row access policy. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Row access policy id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", rowAccessPolicy.Comment); err != nil {
		return diag.FromErr(err)
	}

	rowAccessPolicyDescription, err := client.RowAccessPolicies.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("body", rowAccessPolicyDescription.Body); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("argument", schemas.RowAccessPolicyArgumentsToSchema(rowAccessPolicyDescription.Signature)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.RowAccessPolicyToSchema(rowAccessPolicy)}); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set(DescribeOutputAttributeName, []map[string]any{schemas.RowAccessPolicyDescriptionToSchema(*rowAccessPolicyDescription)}); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// UpdateRowAccessPolicy implements schema.UpdateFunc.
func UpdateRowAccessPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		err := client.RowAccessPolicies.Alter(ctx, sdk.NewAlterRowAccessPolicyRequest(id).WithRenameTo(&newId))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error renaming view %v err = %w", d.Id(), err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	if d.HasChange("comment") {
		comment := d.Get("comment")
		if c := comment.(string); c == "" {
			err := client.RowAccessPolicies.Alter(ctx, sdk.NewAlterRowAccessPolicyRequest(id).WithUnsetComment(sdk.Bool(true)))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error unsetting comment for row access policy on %v err = %w", d.Id(), err))
			}
		} else {
			err := client.RowAccessPolicies.Alter(ctx, sdk.NewAlterRowAccessPolicyRequest(id).WithSetComment(sdk.String(c)))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error updating comment for row access policy on %v err = %w", d.Id(), err))
			}
		}
	}

	if d.HasChange("body") {
		rowAccessExpression := d.Get("body").(string)
		err := client.RowAccessPolicies.Alter(ctx, sdk.NewAlterRowAccessPolicyRequest(id).WithSetBody(sdk.String(rowAccessExpression)))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating row access policy expression on %v err = %w", d.Id(), err))
		}
	}

	return ReadRowAccessPolicy(ctx, d, meta)
}
