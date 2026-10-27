package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	codeBundleCreateExt      = PreviewFeatureCreateContextWrapper(string(previewfeatures.CodeBundleResource), TrackingCreateWrapper(resources.CodeBundle, CreateCodeBundle))
	codeBundleReadExt        = PreviewFeatureReadContextWrapper(string(previewfeatures.CodeBundleResource), TrackingReadWrapper(resources.CodeBundle, ReadCodeBundleFunc(true)))
	codeBundleUpdateExt      = PreviewFeatureUpdateContextWrapper(string(previewfeatures.CodeBundleResource), TrackingUpdateWrapper(resources.CodeBundle, UpdateCodeBundle))
	codeBundleDescriptionExt = "Resource used to manage code bundles. For more information, check [code bundle documentation](https://docs.snowflake.com/en/developer-guide/code-bundles/code-bundles)."
	codeBundleImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.CodeBundle, ImportCodeBundle),
	}
)

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

func codeBundleAfterAlterExt(context.Context, *schema.ResourceData, any, sdk.SchemaObjectIdentifier) diag.Diagnostics {
	return nil
}
