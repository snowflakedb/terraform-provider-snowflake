package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// icebergTableCommonSchema returns the schema fields shared by all Iceberg table resources
func icebergTableCommonSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"database": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			Description:      blocklistedCharactersFieldDescription("The database in which to create the Iceberg table."),
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"schema": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			Description:      blocklistedCharactersFieldDescription("The schema in which to create the Iceberg table."),
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"name": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the Iceberg table; must be unique for the schema in which the Iceberg table is created."),
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"comment": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies a comment for the Iceberg table.",
		},
		FullyQualifiedNameAttributeName: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fully qualified name of the resource. For more information, see [object name resolution](https://docs.snowflake.com/en/sql-reference/name-resolution).",
		},
		ShowOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `SHOW ICEBERG TABLES` for the given Iceberg table. Note that this value will be only recomputed whenever values of fields affecting the output change.",
			Elem:        &schema.Resource{Schema: schemas.ShowIcebergTableSchema},
		},
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE ICEBERG TABLE` for the given Iceberg table.",
			Elem:        &schema.Resource{Schema: schemas.DescribeIcebergTableSchema},
		},
	}
}

// icebergTablePartitionKinds lists the top-level fields of a single partition_by block; exactly one must be set.
var icebergTablePartitionKinds = []string{"identity", "bucket", "truncate", "year", "month", "day", "hour"}

func icebergTablePartitionBySchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		ForceNew:      true,
		ConflictsWith: []string{"cluster_by"},
		Description:   "Defines the partitioning for the Iceberg table. Cannot be changed after creation. Exactly one of identity, bucket, truncate, year, month, day, or hour must be set for each entry. Cannot be used together with `cluster_by`.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"identity": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Name of the column to use as-is for partitioning.",
				},
				"bucket": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					MaxItems:    1,
					Description: "Partitions the table by hashing the column into a fixed number of buckets.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"num_buckets": {Type: schema.TypeInt, Required: true, ForceNew: true, Description: "Number of buckets to hash the column values into."},
							"column":      {Type: schema.TypeString, Required: true, ForceNew: true, Description: "Name of the column to bucket."},
						},
					},
				},
				"truncate": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					MaxItems:    1,
					Description: "Partitions the table by truncating the column value to a fixed width.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"width":  {Type: schema.TypeInt, Required: true, ForceNew: true, Description: "Width to truncate the column value to."},
							"column": {Type: schema.TypeString, Required: true, ForceNew: true, Description: "Name of the column to truncate."},
						},
					},
				},
				"year":  icebergTablePartitionTimeSchema("Partitions the table by the year component of the column."),
				"month": icebergTablePartitionTimeSchema("Partitions the table by the month component of the column."),
				"day":   icebergTablePartitionTimeSchema("Partitions the table by the day component of the column."),
				"hour":  icebergTablePartitionTimeSchema("Partitions the table by the hour component of the column."),
			},
		},
	}
}

func icebergTablePartitionTimeSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		MaxItems:    1,
		Description: description,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"column": {Type: schema.TypeString, Required: true, ForceNew: true, Description: "Name of the date/timestamp column to partition by."},
			},
		},
	}
}

// validateIcebergTablePartitionByDiff ensures that every partition_by entry sets exactly one of
// identity, bucket, truncate, year, month, day, or hour.
func validateIcebergTablePartitionByDiff(d *schema.ResourceDiff) error {
	entries := d.Get("partition_by").([]any)
	for i, e := range entries {
		entry := e.(map[string]any)
		set := 0
		for _, kind := range icebergTablePartitionKinds {
			v := entry[kind]
			switch value := v.(type) {
			case string:
				if value != "" {
					set++
				}
			case []any:
				if len(value) > 0 {
					set++
				}
			}
		}
		if set != 1 {
			return fmt.Errorf("partition_by.%d: exactly one of %v must be set, got %d set", i, icebergTablePartitionKinds, set)
		}
	}
	return nil
}

// parseIcebergTablePartitionBy parses the partition_by blocks from the resource data to SDK objects.
func parseIcebergTablePartitionBy(d *schema.ResourceData) ([]sdk.IcebergTablePartitionExpressionRequest, error) {
	entries := d.Get("partition_by").([]any)
	requests := make([]sdk.IcebergTablePartitionExpressionRequest, len(entries))

	for i := range entries {
		prefix := fmt.Sprintf("partition_by.%d.", i)
		req := sdk.NewIcebergTablePartitionExpressionRequest()

		err := errors.Join(
			attributeMappedValueCreateBuilderNested(d, prefix+"identity", req.WithIdentity, func(d *schema.ResourceData) (string, error) {
				return d.Get(prefix + "identity").(string), nil
			}),
			attributeMappedValueCreateBuilderNested(d, prefix+"bucket", req.WithBucket, parseIcebergTablePartitionBucket(prefix+"bucket.0.")),
			attributeMappedValueCreateBuilderNested(d, prefix+"truncate", req.WithTruncate, parseIcebergTablePartitionTruncate(prefix+"truncate.0.")),
			attributeMappedValueCreateBuilderNested(d, prefix+"year", req.WithYear, parseIcebergTablePartitionTime(prefix+"year.0.", func(args sdk.IcebergTablePartitionTimeArgsRequest) sdk.IcebergTablePartitionYearRequest {
				return *sdk.NewIcebergTablePartitionYearRequest().WithArgs(args)
			})),
			attributeMappedValueCreateBuilderNested(d, prefix+"month", req.WithMonth, parseIcebergTablePartitionTime(prefix+"month.0.", func(args sdk.IcebergTablePartitionTimeArgsRequest) sdk.IcebergTablePartitionMonthRequest {
				return *sdk.NewIcebergTablePartitionMonthRequest().WithArgs(args)
			})),
			attributeMappedValueCreateBuilderNested(d, prefix+"day", req.WithDay, parseIcebergTablePartitionTime(prefix+"day.0.", func(args sdk.IcebergTablePartitionTimeArgsRequest) sdk.IcebergTablePartitionDayRequest {
				return *sdk.NewIcebergTablePartitionDayRequest().WithArgs(args)
			})),
			attributeMappedValueCreateBuilderNested(d, prefix+"hour", req.WithHour, parseIcebergTablePartitionTime(prefix+"hour.0.", func(args sdk.IcebergTablePartitionTimeArgsRequest) sdk.IcebergTablePartitionHourRequest {
				return *sdk.NewIcebergTablePartitionHourRequest().WithArgs(args)
			})),
		)
		if err != nil {
			return nil, err
		}

		requests[i] = *req
	}
	return requests, nil
}

// parseIcebergTablePartitionBucket returns a mapper parsing a bucket partition block at the given prefix.
func parseIcebergTablePartitionBucket(prefix string) func(d *schema.ResourceData) (sdk.IcebergTablePartitionBucketRequest, error) {
	return func(d *schema.ResourceData) (sdk.IcebergTablePartitionBucketRequest, error) {
		return *sdk.NewIcebergTablePartitionBucketRequest().WithArgs(
			*sdk.NewIcebergTablePartitionBucketArgsRequest(d.Get(prefix+"num_buckets").(int), d.Get(prefix+"column").(string)),
		), nil
	}
}

// parseIcebergTablePartitionTruncate returns a mapper parsing a truncate partition block at the given prefix.
func parseIcebergTablePartitionTruncate(prefix string) func(d *schema.ResourceData) (sdk.IcebergTablePartitionTruncateRequest, error) {
	return func(d *schema.ResourceData) (sdk.IcebergTablePartitionTruncateRequest, error) {
		return *sdk.NewIcebergTablePartitionTruncateRequest().WithArgs(
			*sdk.NewIcebergTablePartitionTruncateArgsRequest(d.Get(prefix+"width").(int), d.Get(prefix+"column").(string)),
		), nil
	}
}

// parseIcebergTablePartitionTime returns a mapper parsing a time-based (year/month/day/hour) partition block at the given prefix.
func parseIcebergTablePartitionTime[T any](prefix string, build func(sdk.IcebergTablePartitionTimeArgsRequest) T) func(d *schema.ResourceData) (T, error) {
	return func(d *schema.ResourceData) (T, error) {
		return build(*sdk.NewIcebergTablePartitionTimeArgsRequest(d.Get(prefix + "column").(string))), nil
	}
}

func icebergTablePathLayoutSchema() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		ValidateDiagFunc: StringInSlice(sdk.AsStringList(sdk.AllIcebergTablePathLayouts), true),
		Description:      externalChangesNotDetectedFieldDescription(fmt.Sprintf("Specifies the storage layout for the Iceberg table's Parquet files. Valid values are: %v. Cannot be changed after creation.", sdk.AllIcebergTablePathLayouts)),
	}
}

func icebergTableDeleteFunc() schema.DeleteContextFunc {
	return ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.IcebergTables.DropSafely
		},
	)
}

func importIcebergTable(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	_, err = client.IcebergTables.ShowByIDSafely(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, err := ImportName[sdk.SchemaObjectIdentifier](ctx, d, nil); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func readIcebergTable(ctx context.Context, d *schema.ResourceData, meta any, setExtra func(d *schema.ResourceData, table *sdk.IcebergTable, details []sdk.IcebergTableDetails) error) diag.Diagnostics {
	return readIcebergTableWithParameterHandler(ctx, d, meta, handleIcebergTableParameterRead, schemas.IcebergTableExternallyManagedParametersToSchema, setExtra)
}

func readIcebergTableWithParameterHandler(ctx context.Context, d *schema.ResourceData, meta any, handleParameterRead func(d *schema.ResourceData, parameters []*sdk.Parameter) error, parametersToSchema func([]*sdk.Parameter, *provider.Context) map[string]any, setExtra func(d *schema.ResourceData, table *sdk.IcebergTable, details []sdk.IcebergTableDetails) error) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	table, err := client.IcebergTables.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query Iceberg table. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Iceberg table id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	details, err := client.IcebergTables.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not describe Iceberg table (%s), err = %w", id.FullyQualifiedName(), err))
	}

	parameters, err := client.IcebergTables.ShowParameters(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not show parameters for Iceberg table (%s), err = %w", id.FullyQualifiedName(), err))
	}

	var comment string
	if table.Comment != nil {
		comment = *table.Comment
	}

	providerCtx := meta.(*provider.Context)
	if setExtra != nil {
		if err := setExtra(d, table, details); err != nil {
			return diag.FromErr(err)
		}
	}
	errs := errors.Join(
		d.Set("database", table.DatabaseName),
		d.Set("schema", table.SchemaName),
		d.Set("name", table.Name),
		d.Set("comment", comment),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.IcebergTableToSchema(table)}),
		d.Set(DescribeOutputAttributeName, schemas.IcebergTableDetailsToSchema(details)),
		d.Set(ParametersAttributeName, []map[string]any{parametersToSchema(parameters, providerCtx)}),
		handleParameterRead(d, parameters),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	return nil
}

// handleIcebergTableParametersCreate populates the parameter-backed fields shared by all Iceberg
// table create requests (external_volume, catalog, replace_invalid_characters).
func handleIcebergTableParametersCreate(d *schema.ResourceData, externalVolume, catalog **sdk.AccountObjectIdentifier, replaceInvalidCharacters **bool) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreateWithMapping(d, sdk.IcebergTableParameterExternalVolume, externalVolume, sdk.ParseAccountObjectIdentifier),
		handleParameterCreateWithMapping(d, sdk.IcebergTableParameterCatalog, catalog, sdk.ParseAccountObjectIdentifier),
		handleParameterCreate(d, sdk.IcebergTableParameterReplaceInvalidCharacters, replaceInvalidCharacters),
	)
}

func handleIcebergTableCommonUpdate(d *schema.ResourceData, set *sdk.IcebergTableSetPropertiesRequest, unset *sdk.IcebergTableUnsetPropertiesRequest) diag.Diagnostics {
	if errs := errors.Join(
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
	); errs != nil {
		return diag.FromErr(errs)
	}
	if diags := handleParameterUpdate(d, sdk.IcebergTableParameterReplaceInvalidCharacters, &set.ReplaceInvalidCharacters, &unset.ReplaceInvalidCharacters); len(diags) > 0 {
		return diags
	}
	return nil
}

func applyIcebergTableAlter(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier, set *sdk.IcebergTableSetPropertiesRequest, unset *sdk.IcebergTableUnsetPropertiesRequest) error {
	if !reflect.DeepEqual(*set, *sdk.NewIcebergTableSetPropertiesRequest()) {
		if err := client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).WithSet(*set)); err != nil {
			return err
		}
	}
	if !reflect.DeepEqual(*unset, *sdk.NewIcebergTableUnsetPropertiesRequest()) {
		if err := client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).WithUnset(*unset)); err != nil {
			return err
		}
	}
	return nil
}
