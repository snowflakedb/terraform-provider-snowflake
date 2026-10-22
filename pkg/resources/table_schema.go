package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO [SNOW-1348114]: old implementation was quoting every column, SDK is not quoting them, therefore they are quoted here: decide if we quote columns or not
// TODO [SNOW-1031688]: move data manipulation logic to the SDK - SQL generation or builders part (e.g. different default types/identity)
var tableSchema = collections.MergeMaps(
	tableIdentitySchemaExt,
	tableAttributesSchemaExt,
	map[string]*schema.Schema{FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema},
)
