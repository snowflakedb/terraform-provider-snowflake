package model

import (
	"encoding/json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

type ListingModel struct {
	Name               tfconfig.Variable `json:"name,omitempty"`
	ApplicationPackage tfconfig.Variable `json:"application_package,omitempty"`
	Comment            tfconfig.Variable `json:"comment,omitempty"`
	FullyQualifiedName tfconfig.Variable `json:"fully_qualified_name,omitempty"`
	Manifest           tfconfig.Variable `json:"manifest,omitempty"`
	Publish            tfconfig.Variable `json:"publish,omitempty"`
	Review             tfconfig.Variable `json:"review,omitempty"`
	Share              tfconfig.Variable `json:"share,omitempty"`

	DynamicBlock *config.DynamicBlock `json:"dynamic,omitempty"`

	*config.ResourceModelMeta
}

// TODO(SNOW-1501905): Adjusted manually (because the model generator does not support required list attributes yet).
// Current overrides are not sufficient to generate the model with the required list attribute that would use tfconfig.Variable type.
// The WithManifestValue method should be used in this case as WithManifest is not generated at all.

func ListingWithInlineManifest(
	resourceName string,
	name string,
	manifest string,
) *ListingModel {
	l := &ListingModel{ResourceModelMeta: config.Meta(resourceName, resources.Listing)}
	l.WithName(name)
	l.WithManifestValue(tfconfig.ListVariable(
		tfconfig.MapVariable(map[string]tfconfig.Variable{
			"from_string": tfconfig.StringVariable(manifest),
		}),
	))
	return l
}

func ListingWithStagedManifest(
	resourceName string,
	name string,
	stageId sdk.SchemaObjectIdentifier,
) *ListingModel {
	l := &ListingModel{ResourceModelMeta: config.Meta(resourceName, resources.Listing)}
	l.WithName(name)
	l.WithManifestValue(tfconfig.ListVariable(
		tfconfig.MapVariable(map[string]tfconfig.Variable{
			"from_stage": tfconfig.ListVariable(
				tfconfig.MapVariable(map[string]tfconfig.Variable{
					"stage": tfconfig.StringVariable(stageId.FullyQualifiedName()),
				}),
			),
		}),
	))
	return l
}

///////////////////////////////////////////////////////////////////////
// set proper json marshaling, handle depends on and dynamic blocks //
///////////////////////////////////////////////////////////////////////

func (l *ListingModel) MarshalJSON() ([]byte, error) {
	type Alias ListingModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(l),
		DependsOn: l.DependsOn(),
	})
}

func (l *ListingModel) WithDependsOn(values ...string) *ListingModel {
	l.SetDependsOn(values...)
	return l
}

func (l *ListingModel) WithDynamicBlock(dynamicBlock *config.DynamicBlock) *ListingModel {
	l.DynamicBlock = dynamicBlock
	return l
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

func (l *ListingModel) WithName(name string) *ListingModel {
	l.Name = tfconfig.StringVariable(name)
	return l
}

func (l *ListingModel) WithApplicationPackage(applicationPackage string) *ListingModel {
	l.ApplicationPackage = tfconfig.StringVariable(applicationPackage)
	return l
}

func (l *ListingModel) WithComment(comment string) *ListingModel {
	l.Comment = tfconfig.StringVariable(comment)
	return l
}

func (l *ListingModel) WithFullyQualifiedName(fullyQualifiedName string) *ListingModel {
	l.FullyQualifiedName = tfconfig.StringVariable(fullyQualifiedName)
	return l
}

// manifest attribute type is not yet supported, so WithManifest can't be generated

func (l *ListingModel) WithPublish(publish string) *ListingModel {
	l.Publish = tfconfig.StringVariable(publish)
	return l
}

func (l *ListingModel) WithReview(review string) *ListingModel {
	l.Review = tfconfig.StringVariable(review)
	return l
}

func (l *ListingModel) WithShare(share string) *ListingModel {
	l.Share = tfconfig.StringVariable(share)
	return l
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (l *ListingModel) WithNameValue(value tfconfig.Variable) *ListingModel {
	l.Name = value
	return l
}

func (l *ListingModel) WithApplicationPackageValue(value tfconfig.Variable) *ListingModel {
	l.ApplicationPackage = value
	return l
}

func (l *ListingModel) WithCommentValue(value tfconfig.Variable) *ListingModel {
	l.Comment = value
	return l
}

func (l *ListingModel) WithFullyQualifiedNameValue(value tfconfig.Variable) *ListingModel {
	l.FullyQualifiedName = value
	return l
}

func (l *ListingModel) WithManifestValue(value tfconfig.Variable) *ListingModel {
	l.Manifest = value
	return l
}

func (l *ListingModel) WithPublishValue(value tfconfig.Variable) *ListingModel {
	l.Publish = value
	return l
}

func (l *ListingModel) WithReviewValue(value tfconfig.Variable) *ListingModel {
	l.Review = value
	return l
}

func (l *ListingModel) WithShareValue(value tfconfig.Variable) *ListingModel {
	l.Share = value
	return l
}
