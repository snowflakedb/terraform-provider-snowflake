package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var schemaSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the schema; must be unique for the database in which the schema is created. When the name is `PUBLIC`, during creation the provider checks if this schema has already been created and, in such case, `ALTER` is used to match the desired state."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the schema."),
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"with_managed_access": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      booleanStringFieldDescription("Specifies a managed schema. Managed access schemas centralize privilege management with the schema owner."),
		ValidateDiagFunc: validateBooleanString,
		Default:          BooleanDefault,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShowWithMapping("options", func(x any) any {
			return slices.Contains(sdk.ParseCommaSeparatedStringArray(x.(string), false), "MANAGED ACCESS")
		}),
	},
	"is_transient": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      booleanStringFieldDescription("Specifies the schema as transient. Transient schemas do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss."),
		ValidateDiagFunc: validateBooleanString,
		Default:          BooleanDefault,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShowWithMapping("options", func(x any) any {
			return slices.Contains(sdk.ParseCommaSeparatedStringArray(x.(string), false), "TRANSIENT")
		}),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the schema.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW SCHEMA` for the given object.",
		Elem: &schema.Resource{
			Schema: schemas.ShowSchemaSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE SCHEMA` for the given object. In order to handle this output, one must grant sufficient privileges, e.g. [grant_ownership](./grant_ownership) on all objects in the schema.",
		Elem: &schema.Resource{
			Schema: schemas.SchemaDescribeSchema,
		},
	},
	ParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PARAMETERS IN SCHEMA` for the given object.",
		Elem: &schema.Resource{
			Schema: schemas.ShowSchemaParametersSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// Schema returns a pointer to the resource representing a schema.
func Schema() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseDatabaseObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.DatabaseObjectIdentifier] {
			return client.Schemas.DropSafely
		},
	)

	return &schema.Resource{
		SchemaVersion: 2,

		CreateContext: TrackingCreateWrapper(resources.Schema, CreateContextSchema),
		ReadContext:   TrackingReadWrapper(resources.Schema, ReadContextSchema(true)),
		UpdateContext: TrackingUpdateWrapper(resources.Schema, UpdateContextSchema),
		DeleteContext: TrackingDeleteWrapper(resources.Schema, deleteFunc),
		Description:   "Resource used to manage schema objects. For more information, check [schema documentation](https://docs.snowflake.com/en/sql-reference/sql/create-schema).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.Schema, customdiff.All(
			ComputedIfAnyAttributeChanged(schemaSchema, ShowOutputAttributeName, "name", "comment", "with_managed_access", "is_transient"),
			ComputedIfAnyAttributeChanged(schemaSchema, DescribeOutputAttributeName, "name"),
			ComputedIfAnyAttributeChanged(schemaSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(schemaParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllSchemaParameters), strings.ToLower)...),
			schemaParametersCustomDiff,
		)),

		Schema: collections.MergeMaps(schemaSchema, schemaParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.Schema, ImportSchema),
		},

		StateUpgraders: []schema.StateUpgrader{
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
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportSchema(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] Starting schema import")
	client := meta.(*provider.Context).Client
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

	s, err := client.Schemas.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := d.Set("comment", s.Comment); err != nil {
		return nil, err
	}

	if err := d.Set("is_transient", booleanStringFromBool(s.IsTransient())); err != nil {
		return nil, err
	}

	if err := d.Set("with_managed_access", booleanStringFromBool(s.IsManagedAccess())); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func CreateContextSchema(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	id := sdk.NewDatabaseObjectIdentifier(database, name)

	if strings.EqualFold(strings.TrimSpace(name), "PUBLIC") {
		_, err := client.Schemas.ShowByID(ctx, id)
		if err != nil && !errors.Is(err, sdk.ErrObjectNotFound) {
			return diag.FromErr(err)
		} else if err == nil {
			// there is already a PUBLIC schema, so we need to alter it instead
			log.Printf("[DEBUG] found PUBLIC schema during creation, updating...")
			d.SetId(helpers.EncodeResourceIdentifier(id))
			return UpdateContextSchema(ctx, d, meta)
		}
	}

	opts := &sdk.CreateSchemaOptions{
		Comment: GetConfigPropertyAsPointerAllowingZeroValue[string](d, "comment"),
	}
	if parametersCreateDiags := handleSchemaParametersCreate(d, opts); len(parametersCreateDiags) > 0 {
		return parametersCreateDiags
	}

	if v := d.Get("is_transient").(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		opts.Transient = sdk.Bool(parsed)
	}
	if v := d.Get("with_managed_access").(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		opts.WithManagedAccess = sdk.Bool(parsed)
	}
	if err := client.Schemas.Create(ctx, id, opts); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to create schema.",
				Detail:   fmt.Sprintf("schema name: %s, err: %s", id.FullyQualifiedName(), err),
			},
		}
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextSchema(false)(ctx, d, meta)
}

func ReadContextSchema(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseDatabaseObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		schema, err := client.Schemas.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query schema. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Schema id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("comment", schema.Comment); err != nil {
			return diag.FromErr(err)
		}

		schemaParameters, err := client.Schemas.ShowParameters(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		if diags := handleSchemaParameterRead(d, schemaParameters); diags != nil {
			return diags
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"options", "is_transient", schema.IsTransient(), booleanStringFromBool(schema.IsTransient()), func(x any) any {
					return slices.Contains(sdk.ParseCommaSeparatedStringArray(x.(string), false), "TRANSIENT")
				}},
				outputMapping{"options", "with_managed_access", schema.IsManagedAccess(), booleanStringFromBool(schema.IsManagedAccess()), func(x any) any {
					return slices.Contains(sdk.ParseCommaSeparatedStringArray(x.(string), false), "MANAGED ACCESS")
				}},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, schemaSchema, []string{
			"is_transient",
			"with_managed_access",
		}); err != nil {
			return diag.FromErr(err)
		}

		describeResult, err := client.Schemas.Describe(ctx, schema.ID())
		if err != nil {
			log.Printf("[DEBUG] describing schema: %s, err: %s", id.FullyQualifiedName(), err)
		} else {
			if err = d.Set(DescribeOutputAttributeName, schemas.SchemaDescriptionToSchema(describeResult)); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.SchemaToSchema(schema)}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(ParametersAttributeName, []map[string]any{schemas.SchemaParametersToSchema(schemaParameters)}); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
}

func UpdateContextSchema(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseDatabaseObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") && !d.GetRawState().IsNull() {
		newId := sdk.NewDatabaseObjectIdentifier(d.Get("database").(string), d.Get("name").(string))
		err := client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
			NewName: sdk.Pointer(newId),
		})
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	if d.HasChange("with_managed_access") {
		if v := d.Get("with_managed_access").(string); v != BooleanDefault {
			var err error
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			if parsed {
				err = client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
					EnableManagedAccess: sdk.Pointer(true),
				})
			} else {
				err = client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
					DisableManagedAccess: sdk.Pointer(true),
				})
			}
			if err != nil {
				return diag.FromErr(fmt.Errorf("error handling with_managed_access on %v err = %w", d.Id(), err))
			}
		} else {
			// managed access can not be UNSET to a default value
			if err := client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
				DisableManagedAccess: sdk.Pointer(true),
			}); err != nil {
				return diag.FromErr(fmt.Errorf("error handling with_managed_access on %v err = %w", d.Id(), err))
			}
		}
	}

	set := new(sdk.SchemaSet)
	unset := new(sdk.SchemaUnset)

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		if len(comment) > 0 {
			set.Comment = &comment
		} else {
			unset.Comment = sdk.Bool(true)
		}
	}

	if updateParamDiags := handleSchemaParametersChanges(d, set, unset); len(updateParamDiags) > 0 {
		return updateParamDiags
	}
	if (*set != sdk.SchemaSet{}) {
		err := client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
			Set: set,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if (*unset != sdk.SchemaUnset{}) {
		err := client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
			Unset: unset,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextSchema(false)(ctx, d, meta)
}
