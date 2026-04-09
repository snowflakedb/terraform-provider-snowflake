//go:build non_account_level_tests

package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_TagReferences(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("table domain: tag set manually on table", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		tagValue := "production"
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("id", sdk.DataTypeNumber),
		}
		table, tableCleanup := testClientHelper().Table.CreateWithRequest(t,
			sdk.NewCreateTableRequest(testClientHelper().Ids.RandomSchemaObjectIdentifier(), columns).
				WithTags([]sdk.TagAssociationRequest{
					*sdk.NewTagAssociationRequest(tag.ID(), tagValue),
				}),
		)
		t.Cleanup(tableCleanup)

		refs, err := client.TagReferences.GetForEntity(ctx, sdk.NewGetForEntityTagReferenceRequest(table.ID(), sdk.TagReferenceObjectDomainTable))
		require.NoError(t, err)

		ref, err := collections.FindFirst(refs, func(r sdk.TagReference) bool {
			return r.TagName == tag.ID().Name()
		})
		require.NoError(t, err)

		assertThatObject(t, objectassert.TagReferenceFromObject(t, ref).
			HasTagDatabase(tag.ID().DatabaseName()).
			HasTagSchema(tag.ID().SchemaName()).
			HasTagName(tag.ID().Name()).
			HasTagValue(tagValue).
			HasLevel(sdk.TagReferenceObjectDomainTable).
			HasDomain(sdk.TagReferenceObjectDomainTable).
			HasObjectName(table.ID().Name()).
			HasObjectDatabase(table.ID().DatabaseName()).
			HasObjectSchema(table.ID().SchemaName()).
			HasApplyMethod(sdk.TagReferenceApplyMethodManual).
			HasColumnNameNil(),
		)
	})

	t.Run("table domain: multiple tags on same object", func(t *testing.T) {
		tag1, tag1Cleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tag1Cleanup)

		tag2, tag2Cleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tag2Cleanup)

		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("id", sdk.DataTypeNumber),
		}
		table, tableCleanup := testClientHelper().Table.CreateWithRequest(t,
			sdk.NewCreateTableRequest(testClientHelper().Ids.RandomSchemaObjectIdentifier(), columns).
				WithTags([]sdk.TagAssociationRequest{
					*sdk.NewTagAssociationRequest(tag1.ID(), "v1"),
					*sdk.NewTagAssociationRequest(tag2.ID(), "v2"),
				}),
		)
		t.Cleanup(tableCleanup)

		refs, err := client.TagReferences.GetForEntity(ctx, sdk.NewGetForEntityTagReferenceRequest(table.ID(), sdk.TagReferenceObjectDomainTable))
		require.NoError(t, err)
		require.Len(t, refs, 2)
	})

	t.Run("table domain: inherited tag from schema", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		tagValue := "inherited_value"
		err := client.Tags.Set(ctx, sdk.NewSetTagRequest(sdk.ObjectTypeSchema, testClientHelper().Ids.SchemaId()).WithSetTags([]sdk.TagAssociation{
			{Name: tag.ID(), Value: tagValue},
		}))
		require.NoError(t, err)

		table, tableCleanup := testClientHelper().Table.Create(t)
		t.Cleanup(tableCleanup)

		refs, err := client.TagReferences.GetForEntity(ctx, sdk.NewGetForEntityTagReferenceRequest(table.ID(), sdk.TagReferenceObjectDomainTable))
		require.NoError(t, err)

		ref, err := collections.FindFirst(refs, func(r sdk.TagReference) bool {
			return r.TagName == tag.ID().Name()
		})
		require.NoError(t, err)

		assertThatObject(t, objectassert.TagReferenceFromObject(t, ref).
			HasTagName(tag.ID().Name()).
			HasTagValue(tagValue).
			HasLevel(sdk.TagReferenceObjectDomainSchema).
			HasDomain(sdk.TagReferenceObjectDomainTable).
			HasApplyMethod(sdk.TagReferenceApplyMethodInherited),
		)
	})

	t.Run("view domain: propagated tag from source table", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateWithRequest(t,
			sdk.NewCreateTagRequest(testClientHelper().Ids.RandomSchemaObjectIdentifier()).
				WithAllowedValues([]string{"propagated_value"}).
				WithPropagate(*sdk.NewTagPropagateRequest(sdk.TagPropagationOnDependency)),
		)
		t.Cleanup(tagCleanup)

		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("id", sdk.DataTypeNumber),
		}
		table, tableCleanup := testClientHelper().Table.CreateWithRequest(t,
			sdk.NewCreateTableRequest(testClientHelper().Ids.RandomSchemaObjectIdentifier(), columns).
				WithTags([]sdk.TagAssociationRequest{
					*sdk.NewTagAssociationRequest(tag.ID(), "propagated_value"),
				}),
		)
		t.Cleanup(tableCleanup)

		view, viewCleanup := testClientHelper().View.CreateView(t, fmt.Sprintf("SELECT * FROM %s", table.ID().FullyQualifiedName()))
		t.Cleanup(viewCleanup)

		refs, err := client.TagReferences.GetForEntity(ctx, sdk.NewGetForEntityTagReferenceRequest(view.ID(), sdk.TagReferenceObjectDomainTable))
		require.NoError(t, err)

		ref, err := collections.FindFirst(refs, func(r sdk.TagReference) bool {
			return r.TagName == tag.ID().Name()
		})
		require.NoError(t, err)

		assertThatObject(t, objectassert.TagReferenceFromObject(t, ref).
			HasTagName(tag.ID().Name()).
			HasTagValue("propagated_value").
			HasLevel(sdk.TagReferenceObjectDomainTable).
			HasDomain(sdk.TagReferenceObjectDomainTable).
			HasApplyMethod(sdk.TagReferenceApplyMethodPropagated),
		)
	})

	t.Run("table domain: no tags set", func(t *testing.T) {
		table, tableCleanup := testClientHelper().Table.Create(t)
		t.Cleanup(tableCleanup)

		refs, err := client.TagReferences.GetForEntity(ctx, sdk.NewGetForEntityTagReferenceRequest(table.ID(), sdk.TagReferenceObjectDomainTable))
		require.NoError(t, err)
		assert.Empty(t, refs)
	})

	t.Run("error: non-existent object", func(t *testing.T) {
		_, err := client.TagReferences.GetForEntity(ctx, sdk.NewGetForEntityTagReferenceRequest(NonExistingSchemaObjectIdentifier, sdk.TagReferenceObjectDomainTable))
		require.Error(t, err)
	})
}
