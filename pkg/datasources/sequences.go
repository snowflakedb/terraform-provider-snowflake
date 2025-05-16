package datasources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var sequencesSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the schemas from.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema from which to return the sequences from.",
	},
	"sequences": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The sequences in the schema",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"database": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"schema": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"comment": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

func Sequences() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.SequencesDatasource), TrackingReadWrapper(datasources.Sequences, ReadSequences)),
		Schema:      sequencesSchema,
	}
}

func ReadSequences(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)

	req := sdk.NewShowSequenceRequest().WithIn(sdk.In{
		Schema: sdk.NewDatabaseObjectIdentifier(databaseName, schemaName),
	})
	seqs, err := client.Sequences.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	sequences := []map[string]interface{}{}
	for _, seq := range seqs {
		sequenceMap := map[string]interface{}{}
		sequenceMap["name"] = seq.Name
		sequenceMap["database"] = seq.DatabaseName
		sequenceMap["schema"] = seq.SchemaName
		sequenceMap["comment"] = seq.Comment

		sequences = append(sequences, sequenceMap)
	}

	d.SetId(fmt.Sprintf(`%v|%v`, databaseName, schemaName))
	return diag.FromErr(d.Set("sequences", sequences))
}
