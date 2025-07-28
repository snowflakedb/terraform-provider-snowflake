package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeCortexSearchServiceSchema represents output of DESCRIBE query for the single CortexSearchService.
var DescribeCortexSearchServiceSchema = map[string]*schema.Schema{
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"database_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"schema_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"target_lag": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"warehouse": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"search_column": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"attribute_columns": {
		Type:     schema.TypeList,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"columns": {
		Type:     schema.TypeList,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"definition": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"service_query_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"data_timestamp": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"source_data_num_rows": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"indexing_state": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"indexing_error": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"embedding_model": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = DescribeCortexSearchServiceSchema

func CortexSearchServiceDetailsToSchema(details *sdk.CortexSearchServiceDetails) map[string]any {
	detailsSchema := make(map[string]any)
	detailsSchema["created_on"] = details.CreatedOn
	detailsSchema["name"] = details.Name
	detailsSchema["database_name"] = details.DatabaseName
	detailsSchema["schema_name"] = details.SchemaName
	detailsSchema["target_lag"] = details.TargetLag
	detailsSchema["warehouse"] = details.Warehouse
	if details.SearchColumn != nil {
		detailsSchema["search_column"] = *details.SearchColumn
	}
	detailsSchema["attribute_columns"] = details.AttributeColumns
	detailsSchema["columns"] = details.Columns
	if details.Definition != nil {
		detailsSchema["definition"] = *details.Definition
	}
	if details.Comment != nil {
		detailsSchema["comment"] = *details.Comment
	}
	detailsSchema["service_query_url"] = details.ServiceQueryUrl
	detailsSchema["data_timestamp"] = details.DataTimestamp
	detailsSchema["source_data_num_rows"] = details.SourceDataNumRows
	detailsSchema["indexing_state"] = details.IndexingState
	if details.IndexingError != nil {
		detailsSchema["indexing_error"] = *details.IndexingError
	}
	if details.EmbeddingModel != nil {
		detailsSchema["embedding_model"] = *details.EmbeddingModel
	}
	return detailsSchema
}

var _ = CortexSearchServiceDetailsToSchema
