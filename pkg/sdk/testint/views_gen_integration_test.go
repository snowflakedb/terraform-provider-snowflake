//go:build non_account_level_tests

package testint

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/tracking"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Views(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	table, tableCleanup := testClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)

	sql := fmt.Sprintf("SELECT id FROM %s", table.ID().FullyQualifiedName())

	assertViewWithOptions := func(t *testing.T, view *sdk.View, id sdk.SchemaObjectIdentifier, isSecure bool, comment string) {
		t.Helper()
		assertThatObject(t, objectassert.ViewFromObject(t, view).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasKind("").
			HasReserved("").
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasNonEmptyText().
			HasIsSecure(isSecure).
			HasIsMaterialized(false).
			HasOwnerRoleType("ROLE").
			HasChangeTracking("OFF"))
	}

	assertView := func(t *testing.T, view *sdk.View, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		assertViewWithOptions(t, view, id, false, "")
	}

	assertViewTerse := func(t *testing.T, view *sdk.View, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		assertThatObject(t, objectassert.ViewFromObject(t, view).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasKind("VIEW").
			HasDatabaseName(testClientHelper().Ids.DatabaseId().Name()).
			HasSchemaName(testClientHelper().Ids.SchemaId().Name()).
			// all below are not contained in the terse response, that's why all of them we expect to be empty
			HasReserved("").
			HasOwner("").
			HasComment("").
			HasText("").
			HasIsSecure(false).
			HasIsMaterialized(false).
			HasOwnerRoleType("").
			HasChangeTracking(""))
	}

	assertViewDetailsRow := func(t *testing.T, viewDetails *sdk.ViewDetails) {
		t.Helper()
		assert.Equal(t, sdk.ViewDetails{
			Name:       "ID",
			Type:       "NUMBER(38,0)",
			Kind:       "COLUMN",
			IsNullable: true,
			Default:    nil,
			IsPrimary:  false,
			IsUnique:   false,
			Check:      nil,
			Expression: nil,
			Comment:    nil,
			PolicyName: nil,
		}, *viewDetails)
	}

	assertPolicyReference := func(t *testing.T, policyRef sdk.PolicyReference,
		policyId sdk.SchemaObjectIdentifier,
		policyKind sdk.PolicyKind,
		viewId sdk.SchemaObjectIdentifier,
		refColumnName *string,
	) {
		t.Helper()
		assert.Equal(t, policyId.Name(), policyRef.PolicyName)
		assert.Equal(t, policyKind, policyRef.PolicyKind)
		assert.Equal(t, viewId.Name(), policyRef.RefEntityName)
		assert.Equal(t, "VIEW", policyRef.RefEntityDomain)
		assert.Equal(t, "ACTIVE", *policyRef.PolicyStatus)
		if refColumnName != nil {
			assert.NotNil(t, policyRef.RefColumnName)
			assert.Equal(t, *refColumnName, *policyRef.RefColumnName)
		} else {
			assert.Nil(t, policyRef.RefColumnName)
		}
	}

	assertDataMetricFunctionReference := func(t *testing.T, dataMetricFunctionReference sdk.DataMetricFunctionReference,
		viewId sdk.SchemaObjectIdentifier,
		schedule string,
	) {
		t.Helper()
		assert.Equal(t, viewId.Name(), dataMetricFunctionReference.RefEntityName)
		assert.Equal(t, "View", dataMetricFunctionReference.RefEntityDomain)
		assert.Equal(t, schedule, dataMetricFunctionReference.Schedule)
	}

	cleanupViewProvider := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.Views.Drop(ctx, sdk.NewDropViewRequest(id))
			require.NoError(t, err)
		}
	}

	createViewBasicRequest := func(t *testing.T) *sdk.CreateViewRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		return sdk.NewCreateViewRequest(id, sql)
	}

	createViewWithRequest := func(t *testing.T, request *sdk.CreateViewRequest) *sdk.View {
		t.Helper()
		id := request.GetName()

		err := client.Views.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupViewProvider(id))

		view, err := client.Views.ShowByID(ctx, id)
		require.NoError(t, err)

		return view
	}

	createView := func(t *testing.T) *sdk.View {
		t.Helper()
		return createViewWithRequest(t, createViewBasicRequest(t))
	}

	t.Run("create view: no optionals", func(t *testing.T) {
		request := createViewBasicRequest(t)

		view := createViewWithRequest(t, request)

		assertView(t, view, request.GetName())
	})

	// source https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2085
	t.Run("create view: no table reference", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateViewRequest(id, "SELECT NULL AS TYPE")

		view := createViewWithRequest(t, request)

		assertView(t, view, request.GetName())
	})

	t.Run("create view: with usage tracking comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		plainQuery := "SELECT NULL AS TYPE"
		query, err := tracking.AppendMetadata(plainQuery, tracking.NewVersionedResourceMetadata(resources.View, tracking.CreateOperation))
		require.NoError(t, err)
		request := sdk.NewCreateViewRequest(id, query)

		view := createViewWithRequest(t, request)

		assertView(t, view, request.GetName())
		assert.Equal(t, fmt.Sprintf("CREATE VIEW %s AS %s", id.FullyQualifiedName(), plainQuery), view.Text)
	})

	t.Run("create view: almost complete case - without masking and projection policies", func(t *testing.T) {
		rowAccessPolicy, rowAccessPolicyCleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
		t.Cleanup(rowAccessPolicyCleanup)

		aggregationPolicy, aggregationPolicyCleanup := testClientHelper().AggregationPolicy.CreateAggregationPolicy(t)
		t.Cleanup(aggregationPolicyCleanup)

		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		request := createViewBasicRequest(t).
			WithOrReplace(true).
			WithSecure(true).
			WithTemporary(true).
			WithColumns([]sdk.ViewColumnRequest{
				*sdk.NewViewColumnRequest("column_with_comment").WithComment("column comment"),
			}).
			WithCopyGrants(true).
			WithComment("comment").
			WithRowAccessPolicy(*sdk.NewViewRowAccessPolicyRequest(rowAccessPolicy.ID(), []sdk.Column{{Value: "column_with_comment"}})).
			WithAggregationPolicy(*sdk.NewViewAggregationPolicyRequest(aggregationPolicy).WithEntityKey([]sdk.Column{{Value: "column_with_comment"}})).
			WithTag([]sdk.TagAssociation{{
				Name:  tag.ID(),
				Value: "v2",
			}})

		id := request.GetName()

		view := createViewWithRequest(t, request)

		assertViewWithOptions(t, view, id, true, "comment")
		rowAccessPolicyReferences, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, view.ID(), sdk.PolicyEntityDomainView)
		require.NoError(t, err)
		assert.Len(t, rowAccessPolicyReferences, 2)
		slices.SortFunc(rowAccessPolicyReferences, func(x, y sdk.PolicyReference) int {
			return cmp.Compare(x.PolicyKind, y.PolicyKind)
		})

		assertPolicyReference(t, rowAccessPolicyReferences[0], aggregationPolicy, sdk.PolicyKindAggregationPolicy, view.ID(), nil)

		assertPolicyReference(t, rowAccessPolicyReferences[1], rowAccessPolicy.ID(), sdk.PolicyKindRowAccessPolicy, view.ID(), nil)
		require.NotNil(t, rowAccessPolicyReferences[1].RefArgColumnNames)
		refArgColumnNames := sdk.ParseCommaSeparatedStringArray(*rowAccessPolicyReferences[1].RefArgColumnNames, true)
		assert.Len(t, refArgColumnNames, 1)
		assert.Equal(t, "column_with_comment", refArgColumnNames[0])
	})

	t.Run("create view: masking and projection policies", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicyIdentity(t, testdatatypes.DataTypeNumber)
		t.Cleanup(maskingPolicyCleanup)

		projectionPolicy, projectionPolicyCleanup := testClientHelper().ProjectionPolicy.CreateProjectionPolicy(t)
		t.Cleanup(projectionPolicyCleanup)

		request := createViewBasicRequest(t).
			WithOrReplace(true).
			WithRecursive(true).
			WithColumns([]sdk.ViewColumnRequest{
				*sdk.NewViewColumnRequest("col1").WithMaskingPolicy(
					*sdk.NewViewColumnMaskingPolicyRequest(maskingPolicy.ID()).WithUsing([]sdk.Column{{Value: "col1"}}),
				).WithProjectionPolicy(
					*sdk.NewViewColumnProjectionPolicyRequest(projectionPolicy),
				),
			})

		id := request.GetName()

		view := createViewWithRequest(t, request)

		assertViewWithOptions(t, view, id, false, "")
		assert.Contains(t, view.Text, "RECURSIVE VIEW")
		rowAccessPolicyReferences, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, view.ID(), sdk.PolicyEntityDomainView)
		require.NoError(t, err)
		assert.Len(t, rowAccessPolicyReferences, 2)
		slices.SortFunc(rowAccessPolicyReferences, func(x, y sdk.PolicyReference) int {
			return cmp.Compare(x.PolicyKind, y.PolicyKind)
		})

		assertPolicyReference(t, rowAccessPolicyReferences[0], maskingPolicy.ID(), sdk.PolicyKindMaskingPolicy, view.ID(), sdk.Pointer("col1"))
		assertPolicyReference(t, rowAccessPolicyReferences[1], projectionPolicy, sdk.PolicyKindProjectionPolicy, view.ID(), sdk.Pointer("col1"))
	})

	t.Run("drop view: existing", func(t *testing.T) {
		request := createViewBasicRequest(t)
		id := request.GetName()

		err := client.Views.Create(ctx, request)
		require.NoError(t, err)

		err = client.Views.Drop(ctx, sdk.NewDropViewRequest(id))
		require.NoError(t, err)

		_, err = client.Views.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop view: non-existing", func(t *testing.T) {
		err := client.Views.Drop(ctx, sdk.NewDropViewRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter view: rename", func(t *testing.T) {
		createRequest := createViewBasicRequest(t)
		id := createRequest.GetName()

		err := client.Views.Create(ctx, createRequest)
		require.NoError(t, err)

		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		alterRequest := sdk.NewAlterViewRequest(id).WithRenameTo(newId)

		err = client.Views.Alter(ctx, alterRequest)
		if err != nil {
			t.Cleanup(cleanupViewProvider(id))
		} else {
			t.Cleanup(cleanupViewProvider(newId))
		}
		require.NoError(t, err)

		_, err = client.Views.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		view, err := client.Views.ShowByID(ctx, newId)
		require.NoError(t, err)

		assertView(t, view, newId)
	})

	t.Run("alter view: set and unset values", func(t *testing.T) {
		view := createView(t)
		id := view.ID()

		alterRequest := sdk.NewAlterViewRequest(id).WithSetComment("new comment")
		err := client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err := client.Views.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "new comment", alteredView.Comment)

		alterRequest = sdk.NewAlterViewRequest(id).WithSetSecure(true)
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err = client.Views.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.True(t, alteredView.IsSecure)

		alterRequest = sdk.NewAlterViewRequest(id).WithSetChangeTracking(true)
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err = client.Views.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "ON", alteredView.ChangeTracking)

		alterRequest = sdk.NewAlterViewRequest(id).WithUnsetComment(true)
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err = client.Views.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "", alteredView.Comment)

		alterRequest = sdk.NewAlterViewRequest(id).WithUnsetSecure(true)
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err = client.Views.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.False(t, alteredView.IsSecure)

		alterRequest = sdk.NewAlterViewRequest(id).WithSetChangeTracking(false)
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err = client.Views.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "OFF", alteredView.ChangeTracking)
	})

	t.Run("alter view: set and unset masking policy on column", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicyIdentity(t, testdatatypes.DataTypeNumber)
		t.Cleanup(maskingPolicyCleanup)

		view := createView(t)
		id := view.ID()

		alterRequest := sdk.NewAlterViewRequest(id).WithSetMaskingPolicyOnColumn(
			*sdk.NewViewSetColumnMaskingPolicyRequest("ID", maskingPolicy.ID()),
		)
		err := client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		policyReferences, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, view.ID(), sdk.PolicyEntityDomainView)
		require.NoError(t, err)
		require.Len(t, policyReferences, 1)

		assertPolicyReference(t, policyReferences[0], maskingPolicy.ID(), sdk.PolicyKindMaskingPolicy, view.ID(), sdk.Pointer("ID"))

		alterRequest = sdk.NewAlterViewRequest(id).WithUnsetMaskingPolicyOnColumn(
			*sdk.NewViewUnsetColumnMaskingPolicyRequest("ID"),
		)
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, view.ID(), sdk.PolicyEntityDomainView)
		require.NoError(t, err)
		require.Empty(t, references)
	})

	t.Run("alter view: set and unset projection policy on column", func(t *testing.T) {
		projectionPolicy, projectionPolicyCleanup := testClientHelper().ProjectionPolicy.CreateProjectionPolicy(t)
		t.Cleanup(projectionPolicyCleanup)

		view := createView(t)
		id := view.ID()

		alterRequest := sdk.NewAlterViewRequest(id).WithSetProjectionPolicyOnColumn(
			*sdk.NewViewSetProjectionPolicyRequest("ID", projectionPolicy),
		)
		err := client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		rowAccessPolicyReferences, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, view.ID(), sdk.PolicyEntityDomainView)
		require.NoError(t, err)
		require.Len(t, rowAccessPolicyReferences, 1)

		assertPolicyReference(t, rowAccessPolicyReferences[0], projectionPolicy, sdk.PolicyKindProjectionPolicy, view.ID(), sdk.Pointer("ID"))

		alterRequest = sdk.NewAlterViewRequest(id).WithUnsetProjectionPolicyOnColumn(
			*sdk.NewViewUnsetProjectionPolicyRequest("ID"),
		)
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, view.ID(), sdk.PolicyEntityDomainView)
		require.NoError(t, err)
		require.Empty(t, references)
	})

	t.Run("alter view: add and drop row access policies", func(t *testing.T) {
		rowAccessPolicy, rowAccessPolicyCleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
		t.Cleanup(rowAccessPolicyCleanup)
		rowAccessPolicy2, rowAccessPolicy2Cleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
		t.Cleanup(rowAccessPolicy2Cleanup)

		view := createView(t)
		id := view.ID()

		// add policy
		alterRequest := sdk.NewAlterViewRequest(id).WithAddRowAccessPolicy(*sdk.NewViewAddRowAccessPolicyRequest(rowAccessPolicy.ID(), []sdk.Column{{Value: "ID"}}))
		err := client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		rowAccessPolicyReference, err := testClientHelper().PolicyReferences.GetPolicyReference(t, view.ID(), sdk.PolicyEntityDomainView)
		require.NoError(t, err)

		assertPolicyReference(t, *rowAccessPolicyReference, rowAccessPolicy.ID(), sdk.PolicyKindRowAccessPolicy, view.ID(), nil)

		// remove policy
		alterRequest = sdk.NewAlterViewRequest(id).WithDropRowAccessPolicy(*sdk.NewViewDropRowAccessPolicyRequest(rowAccessPolicy.ID()))
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, view.ID(), sdk.PolicyEntityDomainView)
		require.NoError(t, err)
		require.Empty(t, references)

		// add policy again
		alterRequest = sdk.NewAlterViewRequest(id).WithAddRowAccessPolicy(*sdk.NewViewAddRowAccessPolicyRequest(rowAccessPolicy.ID(), []sdk.Column{{Value: "ID"}}))
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		rowAccessPolicyReference, err = testClientHelper().PolicyReferences.GetPolicyReference(t, view.ID(), sdk.PolicyEntityDomainView)
		require.NoError(t, err)
		assert.Equal(t, rowAccessPolicy.ID().Name(), rowAccessPolicyReference.PolicyName)

		// drop and add other policy simultaneously
		alterRequest = sdk.NewAlterViewRequest(id).WithDropAndAddRowAccessPolicy(*sdk.NewViewDropAndAddRowAccessPolicyRequest(
			*sdk.NewViewDropRowAccessPolicyRequest(rowAccessPolicy.ID()),
			*sdk.NewViewAddRowAccessPolicyRequest(rowAccessPolicy2.ID(), []sdk.Column{{Value: "ID"}}),
		))
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		rowAccessPolicyReference, err = testClientHelper().PolicyReferences.GetPolicyReference(t, view.ID(), sdk.PolicyEntityDomainView)
		require.NoError(t, err)
		assert.Equal(t, rowAccessPolicy2.ID().Name(), rowAccessPolicyReference.PolicyName)

		// drop all policies
		alterRequest = sdk.NewAlterViewRequest(id).WithDropAllRowAccessPolicies(true)
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		references, err = testClientHelper().PolicyReferences.GetPolicyReferences(t, view.ID(), sdk.PolicyEntityDomainView)
		require.NoError(t, err)
		require.Empty(t, references)
	})

	t.Run("alter view: add and drop data metrics", func(t *testing.T) {
		view := createView(t)
		id := view.ID()

		dataMetricFunction, dataMetricFunctionCleanup := testClientHelper().DataMetricFunctionClient.CreateDataMetricFunction(t, id)
		t.Cleanup(dataMetricFunctionCleanup)
		dataMetricFunction2, dataMetricFunction2Cleanup := testClientHelper().DataMetricFunctionClient.CreateDataMetricFunction(t, id)
		t.Cleanup(dataMetricFunction2Cleanup)

		// set cron schedule
		cron := "*/5 * * * * UTC"
		alterRequest := sdk.NewAlterViewRequest(id).WithSetDataMetricSchedule(*sdk.NewViewSetDataMetricScheduleRequest("USING CRON " + cron))
		err := client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		// add data metric function
		alterRequest = sdk.NewAlterViewRequest(id).WithAddDataMetricFunction(*sdk.NewViewAddDataMetricFunctionRequest([]sdk.ViewDataMetricFunction{
			{
				DataMetricFunction: dataMetricFunction,
				On:                 []sdk.Column{{Value: "ID"}},
			},
		}))
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		dataMetricFunctionReferences := testClientHelper().DataMetricFunctionReferences.GetDataMetricFunctionReferences(t, view.ID(), sdk.DataMetricFunctionRefEntityDomainView)
		require.Len(t, dataMetricFunctionReferences, 1)

		assertDataMetricFunctionReference(t, dataMetricFunctionReferences[0], view.ID(), cron)

		// remove function
		alterRequest = sdk.NewAlterViewRequest(id).WithDropDataMetricFunction(*sdk.NewViewDropDataMetricFunctionRequest([]sdk.ViewDataMetricFunction{
			{
				DataMetricFunction: dataMetricFunction,
				On:                 []sdk.Column{{Value: "ID"}},
			},
		}))
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		dataMetricFunctionReferences = testClientHelper().DataMetricFunctionReferences.GetDataMetricFunctionReferences(t, view.ID(), sdk.DataMetricFunctionRefEntityDomainView)
		require.NoError(t, err)
		require.Empty(t, dataMetricFunctionReferences)

		// add two functions
		alterRequest = sdk.NewAlterViewRequest(id).WithAddDataMetricFunction(*sdk.NewViewAddDataMetricFunctionRequest([]sdk.ViewDataMetricFunction{
			{
				DataMetricFunction: dataMetricFunction,
				On:                 []sdk.Column{{Value: "ID"}},
			},
			{
				DataMetricFunction: dataMetricFunction2,
				On:                 []sdk.Column{{Value: "ID"}},
			},
		}))
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		dataMetricFunctionReferences = testClientHelper().DataMetricFunctionReferences.GetDataMetricFunctionReferences(t, view.ID(), sdk.DataMetricFunctionRefEntityDomainView)
		require.Len(t, dataMetricFunctionReferences, 2)

		assertDataMetricFunctionReference(t, dataMetricFunctionReferences[0], view.ID(), cron)
		assertDataMetricFunctionReference(t, dataMetricFunctionReferences[1], view.ID(), cron)

		// drop all functions
		alterRequest = sdk.NewAlterViewRequest(id).WithDropAllRowAccessPolicies(true)
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, view.ID(), sdk.PolicyEntityDomainView)
		require.NoError(t, err)
		require.Empty(t, references)
	})

	t.Run("alter view: set and unset aggregation policies", func(t *testing.T) {
		aggregationPolicy, aggregationPolicyCleanup := testClientHelper().AggregationPolicy.CreateAggregationPolicy(t)
		t.Cleanup(aggregationPolicyCleanup)
		aggregationPolicy2, aggregationPolicy2Cleanup := testClientHelper().AggregationPolicy.CreateAggregationPolicy(t)
		t.Cleanup(aggregationPolicy2Cleanup)

		view := createView(t)
		id := view.ID()

		// set policy
		alterRequest := sdk.NewAlterViewRequest(id).WithSetAggregationPolicy(*sdk.NewViewSetAggregationPolicyRequest(aggregationPolicy).WithEntityKey([]sdk.Column{{Value: "ID"}}))
		err := client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		rowAccessPolicyReferences, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, view.ID(), sdk.PolicyEntityDomainView)
		require.NoError(t, err)
		require.Len(t, rowAccessPolicyReferences, 1)

		assertPolicyReference(t, rowAccessPolicyReferences[0], aggregationPolicy, sdk.PolicyKindAggregationPolicy, view.ID(), nil)

		// set policy with force
		alterRequest = sdk.NewAlterViewRequest(id).WithSetAggregationPolicy(*sdk.NewViewSetAggregationPolicyRequest(aggregationPolicy2).
			WithEntityKey([]sdk.Column{{Value: "ID"}}).
			WithForce(true))
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		rowAccessPolicyReferences, err = testClientHelper().PolicyReferences.GetPolicyReferences(t, view.ID(), sdk.PolicyEntityDomainView)
		require.NoError(t, err)
		require.Len(t, rowAccessPolicyReferences, 1)

		assertPolicyReference(t, rowAccessPolicyReferences[0], aggregationPolicy2, sdk.PolicyKindAggregationPolicy, view.ID(), nil)

		// remove policy
		alterRequest = sdk.NewAlterViewRequest(id).WithUnsetAggregationPolicy(*sdk.NewViewUnsetAggregationPolicyRequest())
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		references, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, view.ID(), sdk.PolicyEntityDomainView)
		require.NoError(t, err)
		require.Empty(t, references)
	})

	t.Run("show view: default", func(t *testing.T) {
		view1 := createView(t)
		view2 := createView(t)

		showRequest := sdk.NewShowViewRequest()
		returnedViews, err := client.Views.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.LessOrEqual(t, 2, len(returnedViews))
		assert.Contains(t, returnedViews, *view1)
		assert.Contains(t, returnedViews, *view2)
	})

	t.Run("show view: terse", func(t *testing.T) {
		view := createView(t)

		showRequest := sdk.NewShowViewRequest().WithTerse(true)
		returnedViews, err := client.Views.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.LessOrEqual(t, 1, len(returnedViews))
		assertViewTerse(t, &returnedViews[0], view.ID())
	})

	t.Run("show view: with options", func(t *testing.T) {
		view1 := createView(t)
		view2 := createView(t)

		showRequest := sdk.NewShowViewRequest().
			WithLike(sdk.Like{Pattern: &view1.Name}).
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: testClientHelper().Ids.SchemaId()}}).
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(5)})
		returnedViews, err := client.Views.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Len(t, returnedViews, 1)
		assert.Contains(t, returnedViews, *view1)
		assert.NotContains(t, returnedViews, *view2)
	})

	// proves issue https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2506
	t.Run("show view by id: same name in different schemas", func(t *testing.T) {
		// we assume that SF returns views alphabetically
		schemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifierWithPrefix("aaaa")
		schema, schemaCleanup := testClientHelper().Schema.CreateSchemaWithIdentifier(t, schemaId)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		request1 := sdk.NewCreateViewRequest(id1, sql)
		request2 := sdk.NewCreateViewRequest(id2, sql)

		createViewWithRequest(t, request1)
		createViewWithRequest(t, request2)

		returnedView1, err := client.Views.ShowByID(ctx, id1)
		require.NoError(t, err)
		returnedView2, err := client.Views.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id1, returnedView1.ID())
		require.Equal(t, id2, returnedView2.ID())
	})

	t.Run("describe view", func(t *testing.T) {
		view := createView(t)

		returnedViewDetails, err := client.Views.Describe(ctx, view.ID())
		require.NoError(t, err)

		assert.Len(t, returnedViewDetails, 1)
		assertViewDetailsRow(t, &returnedViewDetails[0])
	})

	t.Run("describe view: non-existing", func(t *testing.T) {
		_, err := client.Views.Describe(ctx, NonExistingSchemaObjectIdentifier)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_ViewsShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	table, tableCleanup := testClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)

	sql := fmt.Sprintf("SELECT id FROM %s", table.ID().FullyQualifiedName())

	cleanupViewHandle := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.Views.Drop(ctx, sdk.NewDropViewRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createViewHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		err := client.Views.Create(ctx, sdk.NewCreateViewRequest(id, sql))
		require.NoError(t, err)
		t.Cleanup(cleanupViewHandle(id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		createViewHandle(t, id1)
		createViewHandle(t, id2)

		e1, err := client.Views.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Views.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
