package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO [SNOW-1348325,SNOW-1501905]: Add this parameter to SDK
const ReplicableWithFailoverGroups = "REPLICABLE_WITH_FAILOVER_GROUPS"

var objectParameterSchema = map[string]*schema.Schema{
	"key": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Name of object parameter. Valid values are those in [object parameters](https://docs.snowflake.com/en/sql-reference/parameters.html#object-parameters).",
	},
	"value": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Value of object parameter, as a string. Constraints are the same as those for the parameters in Snowflake documentation.",
	},
	"on_account": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "If true, the object parameter will be set on the account level.",
	},
	"object_type": {
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		Description:  "Type of object to which the parameter applies. Valid values are those in [object types](https://docs.snowflake.com/en/sql-reference/parameters.html#object-types). If no value is provided, then the resource will default to setting the object parameter at account level.",
		RequiredWith: []string{"object_identifier"},
	},
	"object_identifier": {
		Type:         schema.TypeList,
		Optional:     true,
		MinItems:     1,
		Description:  "Specifies the object identifier for the object parameter. If no value is provided, then the resource will default to setting the object parameter at account level.",
		RequiredWith: []string{"object_type"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Name of the object to set the parameter for.",
				},
				"database": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Name of the database that the object was created in.",
				},
				"schema": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Name of the schema that the object was created in.",
				},
			},
		},
	},
}

func ObjectParameter() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ObjectParameterResource), TrackingCreateWrapper(resources.ObjectParameter, CreateObjectParameter)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ObjectParameterResource), TrackingReadWrapper(resources.ObjectParameter, ReadObjectParameter)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ObjectParameterResource), TrackingUpdateWrapper(resources.ObjectParameter, UpdateObjectParameter)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ObjectParameterResource), TrackingDeleteWrapper(resources.ObjectParameter, DeleteObjectParameter)),

		Schema: objectParameterSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

// CreateObjectParameter implements schema.CreateFunc.
func CreateObjectParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	key := d.Get("key").(string)
	value := d.Get("value").(string)

	parameter := sdk.ObjectParameter(key)

	o := sdk.Object{}
	if v, ok := d.GetOk("object_identifier"); ok {
		objectDatabase, objectSchema, objectName := expandObjectIdentifier(v.([]interface{}))
		fullyQualifierObjectIdentifier := FormatFullyQualifiedObjectID(objectDatabase, objectSchema, objectName)
		fullyQualifierObjectIdentifier = strings.Trim(fullyQualifierObjectIdentifier, "\"")
		o.Name = sdk.NewObjectIdentifierFromFullyQualifiedName(fullyQualifierObjectIdentifier)
		o.ObjectType = sdk.ObjectType(d.Get("object_type").(string))
	}

	onAccount := d.Get("on_account").(bool)
	if onAccount {
		err := client.Parameters.SetObjectParameterOnAccount(ctx, parameter, value)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error creating object parameter err = %w", err))
		}
	} else {
		err := client.Parameters.SetObjectParameterOnObject(ctx, o, parameter, value)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error setting object parameter err = %w", err))
		}
	}

	var id string
	if onAccount {
		id = fmt.Sprintf("%v||", key)
	} else {
		id = fmt.Sprintf("%v|%v|%v", key, o.ObjectType, o.Name.FullyQualifiedName())
	}

	d.SetId(id)
	var err error
	var p *sdk.Parameter
	if onAccount {
		p, err = client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameter(key))
	} else {
		p, err = client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameter(key), o)
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading object parameter err = %w", err))
	}
	err = d.Set("value", p.Value)
	if err != nil {
		return diag.FromErr(err)
	}
	return ReadObjectParameter(ctx, d, meta)
}

// ReadObjectParameter implements schema.ReadFunc.
func ReadObjectParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := d.Id()
	parts := strings.Split(id, "|")
	if len(parts) != 3 {
		parts = strings.Split(id, "❄️") // for backwards compatibility
	}
	if len(parts) != 3 {
		return diag.FromErr(fmt.Errorf("unexpected format of ID (%v), expected key|object_type|object_identifier", id))
	}
	key := parts[0]
	var p *sdk.Parameter
	var err error
	if parts[1] == "" {
		p, err = client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameter(key))
	} else {
		objectType := sdk.ObjectType(parts[1])
		objectIdentifier := parts[2]
		objectIdentifier = strings.Trim(objectIdentifier, "\"")
		p, err = client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameter(key), sdk.Object{
			ObjectType: objectType,
			Name:       sdk.NewObjectIdentifierFromFullyQualifiedName(objectIdentifier),
		})
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading object parameter err = %w", err))
	}
	if err := d.Set("value", p.Value); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// UpdateObjectParameter implements schema.UpdateFunc.
func UpdateObjectParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return CreateObjectParameter(ctx, d, meta)
}

// DeleteObjectParameter implements schema.DeleteFunc.
func DeleteObjectParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	key := d.Get("key").(string)

	onAccount := d.Get("on_account").(bool)
	if onAccount {
		defaultParameter, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameter(key))
		if err != nil {
			return diag.FromErr(err)
		}
		defaultValue := defaultParameter.Default
		err = client.Parameters.SetAccountParameter(ctx, sdk.AccountParameter(key), defaultValue)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error resetting account parameter err = %w", err))
		}
	} else {
		v := d.Get("object_identifier")
		objectDatabase, objectSchema, objectName := expandObjectIdentifier(v.([]interface{}))
		fullyQualifierObjectIdentifier := FormatFullyQualifiedObjectID(objectDatabase, objectSchema, objectName)
		fullyQualifierObjectIdentifier = strings.Trim(fullyQualifierObjectIdentifier, "\"")
		o := sdk.Object{
			ObjectType: sdk.ObjectType(d.Get("object_type").(string)),
			Name:       sdk.NewObjectIdentifierFromFullyQualifiedName(fullyQualifierObjectIdentifier),
		}
		objectParameter := sdk.ObjectParameter(key)
		defaultParameter, err := client.Parameters.ShowObjectParameter(ctx, objectParameter, o)
		if err != nil {
			return diag.FromErr(err)
		}
		defaultValue := defaultParameter.Default
		// TODO [SNOW-2395064,SNOW-1348325]: Remove the workaround when default value on REPLICABLE_WITH_FAILOVER_GROUPS parameter in Snowflake is settable or during snowflake_object_parameter rework (using UNSET)
		if key == ReplicableWithFailoverGroups {
			defaultValue = "YES"
		}
		err = client.Parameters.SetObjectParameterOnObject(ctx, o, objectParameter, defaultValue)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error resetting object parameter err = %w", err))
		}
	}
	d.SetId("")
	return nil
}
