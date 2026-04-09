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

		table, tableCleanup := testClientHelper().Table.Create(t)
		t.Cleanup(tableCleanup)

		tagValue := "production"
		err := client.Tags.Set(ctx, sdk.NewSetTagRequest(sdk.ObjectTypeTable, table.ID()).WithSetTags([]sdk.TagAssociation{
			{Name: tag.ID(), Value: tagValue},
		}))
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.Tags.Unset(ctx, sdk.NewUnsetTagRequest(sdk.ObjectTypeTable, table.ID()).WithUnsetTags([]sdk.ObjectIdentifier{tag.ID()}))
			require.NoError(t, err)
		})

		refs, err := client.TagReferences.GetForEntity(ctx, sdk.NewGetForEntityTagReferenceRequest(table.ID(), sdk.TagReferenceObjectDomainTable))
		require.NoError(t, err)
		require.Len(t, refs, 1)

		assertThatObject(t, objectassert.TagReferenceFromObject(t, &refs[0]).
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

	t.Run("table domain: no tags set", func(t *testing.T) {
		table, tableCleanup := testClientHelper().Table.Create(t)
		t.Cleanup(tableCleanup)

		refs, err := client.TagReferences.GetForEntity(ctx, sdk.NewGetForEntityTagReferenceRequest(table.ID(), sdk.TagReferenceObjectDomainTable))
		require.NoError(t, err)
		assert.Empty(t, refs)
	})

	t.Run("table domain: multiple tags on same object", func(t *testing.T) {
		tag1, tag1Cleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tag1Cleanup)

		tag2, tag2Cleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tag2Cleanup)

		table, tableCleanup := testClientHelper().Table.Create(t)
		t.Cleanup(tableCleanup)

		err := client.Tags.Set(ctx, sdk.NewSetTagRequest(sdk.ObjectTypeTable, table.ID()).WithSetTags([]sdk.TagAssociation{
			{Name: tag1.ID(), Value: "v1"},
			{Name: tag2.ID(), Value: "v2"},
		}))
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.Tags.Unset(ctx, sdk.NewUnsetTagRequest(sdk.ObjectTypeTable, table.ID()).WithUnsetTags([]sdk.ObjectIdentifier{tag1.ID(), tag2.ID()}))
			require.NoError(t, err)
		})

		refs, err := client.TagReferences.GetForEntity(ctx, sdk.NewGetForEntityTagReferenceRequest(table.ID(), sdk.TagReferenceObjectDomainTable))
		require.NoError(t, err)
		require.Len(t, refs, 2)
	})

	t.Run("table domain: inherited tag from schema", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		schemaId := testClientHelper().Ids.SchemaId()
		tagValue := "inherited_value"
		err := client.Tags.Set(ctx, sdk.NewSetTagRequest(sdk.ObjectTypeSchema, schemaId).WithSetTags([]sdk.TagAssociation{
			{Name: tag.ID(), Value: tagValue},
		}))
		require.NoError(t, err)

		t.Cleanup(func() {
			err := client.Tags.Unset(ctx, sdk.NewUnsetTagRequest(sdk.ObjectTypeSchema, schemaId).WithUnsetTags([]sdk.ObjectIdentifier{tag.ID()}))
			require.NoError(t, err)
		})

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

		table, tableCleanup := testClientHelper().Table.Create(t)
		t.Cleanup(tableCleanup)

		err := client.Tags.Set(ctx, sdk.NewSetTagRequest(sdk.ObjectTypeTable, table.ID()).WithSetTags([]sdk.TagAssociation{
			{Name: tag.ID(), Value: "propagated_value"},
		}))
		require.NoError(t, err)

		view, viewCleanup := testClientHelper().View.CreateView(t, fmt.Sprintf("SELECT * FROM %s", table.ID().FullyQualifiedName()))
		t.Cleanup(viewCleanup)

		refs, err := client.TagReferences.GetForEntity(ctx, sdk.NewGetForEntityTagReferenceRequest(view.ID(), sdk.TagReferenceObjectDomainTable))
		require.NoError(t, err)
		require.Len(t, refs, 1)

		assertThatObject(t, objectassert.TagReferenceFromObject(t, &refs[0]).
			HasTagName(tag.ID().Name()).
			HasTagValue("propagated_value").
			HasLevel(sdk.TagReferenceObjectDomainTable).
			HasDomain(sdk.TagReferenceObjectDomainTable).
			HasApplyMethod(sdk.TagReferenceApplyMethodPropagated),
		)
	})

	t.Run("error: non-existent object", func(t *testing.T) {
		nonExistentId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		_, err := client.TagReferences.GetForEntity(ctx, sdk.NewGetForEntityTagReferenceRequest(nonExistentId, sdk.TagReferenceObjectDomainTable))
		require.Error(t, err)
	})
}
