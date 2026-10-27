package resources

import (
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	codeBundleCreateExt      = TrackingCreateWrapper(resources.CodeBundle, CreateCodeBundle)
	codeBundleReadExt        = TrackingReadWrapper(resources.CodeBundle, ReadCodeBundleFunc(true))
	codeBundleUpdateExt      = TrackingUpdateWrapper(resources.CodeBundle, UpdateCodeBundle)
	codeBundleDescriptionExt = "Resource used to manage code bundles. For more information, check [code bundle documentation](https://docs.snowflake.com/en/developer-guide/code-bundles/code-bundles)."
	codeBundleImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.CodeBundle, ImportCodeBundle),
	}
)

type codeBundleAlterRequestsExt struct {
	set   *sdk.CodeBundleSetRequest
	unset *sdk.CodeBundleUnsetRequest
}

// from is not returned by SHOW. Read keeps the configured value.
var codeBundleSetStateToValueFromConfigFieldsExt = []string{
	"from",
}

func codeBundleSetConfigFieldsExt(d *schema.ResourceData, codeBundle *sdk.CodeBundle) error {
	comment := ""
	if codeBundle.Comment != nil {
		comment = *codeBundle.Comment
	}
	return errors.Join(
		d.Set("name", codeBundle.Name),
		d.Set("database", codeBundle.DatabaseName),
		d.Set("schema", codeBundle.SchemaName),
		d.Set("comment", comment),
	)
}

func codeBundleSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier, codeBundle *sdk.CodeBundle, _ bool) error {
	comment := ""
	if codeBundle.Comment != nil {
		comment = *codeBundle.Comment
	}
	return errors.Join(
		d.Set("name", id.Name()),
		d.Set("database", id.DatabaseName()),
		d.Set("schema", id.SchemaName()),
		d.Set("comment", comment),
	)
}
