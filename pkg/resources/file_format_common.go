package resources

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var fileFormatCommonSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the file format; must be unique for the database and schema in which the file format is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the file format."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the file format."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the file format.",
	},
	"type": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Specifies the type of the file format. This field is used to detect when the file format type was changed outside of Terraform and to recreate the resource when that happens.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW FILE FORMATS` for this file format.",
		Elem: &schema.Resource{
			Schema: schemas.ShowFileFormatSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// isFileFormatAutoSentinel reports whether v is the special "AUTO" value Snowflake uses to
// mean "detect the format automatically", as opposed to a literal format string.
func isFileFormatAutoSentinel(v string) bool {
	return strings.ToUpper(v) == "AUTO"
}

func fileFormatStringOrAutoMapper(v string) (sdk.StageFileFormatStringOrAutoRequest, error) {
	if isFileFormatAutoSentinel(v) {
		return *sdk.NewStageFileFormatStringOrAutoRequest().WithAuto(true), nil
	}
	return *sdk.NewStageFileFormatStringOrAutoRequest().WithValue(v), nil
}
