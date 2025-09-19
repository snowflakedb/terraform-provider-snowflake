//go:build !account_level_tests

package testint

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_GetTag(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	tagTest, tagCleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	t.Run("masking policy tag", func(t *testing.T) {
		maskingPolicyTest, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicy(t)
		t.Cleanup(maskingPolicyCleanup)

		tagValue := random.String()
		err := client.MaskingPolicies.Alter(ctx, maskingPolicyTest.ID(), &sdk.AlterMaskingPolicyOptions{
			SetTag: []sdk.TagAssociation{
				{
					Name:  tagTest.ID(),
					Value: tagValue,
				},
			},
		})
		require.NoError(t, err)
		s, err := client.SystemFunctions.GetTag(ctx, tagTest.ID(), maskingPolicyTest.ID(), sdk.ObjectTypeMaskingPolicy)
		require.NoError(t, err)
		assert.Equal(t, &tagValue, s)
	})

	t.Run("masking policy with no set tag", func(t *testing.T) {
		maskingPolicyTest, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicy(t)
		t.Cleanup(maskingPolicyCleanup)

		s, err := client.SystemFunctions.GetTag(ctx, tagTest.ID(), maskingPolicyTest.ID(), sdk.ObjectTypeMaskingPolicy)
		require.NoError(t, err)
		assert.Nil(t, s)
	})
	t.Run("unsupported object type", func(t *testing.T) {
		_, err := client.SystemFunctions.GetTag(ctx, tagTest.ID(), testClientHelper().Ids.RandomAccountObjectIdentifier(), sdk.ObjectTypeSequence)
		require.ErrorContains(t, err, "tagging for object type SEQUENCE is not supported")
	})
}

func TestInt_PipeStatus(t *testing.T) {
	client := testClient(t)

	schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	table, tableCleanup := testClientHelper().Table.CreateInSchema(t, schema.ID())
	t.Cleanup(tableCleanup)

	stage, stageCleanup := testClientHelper().Stage.CreateStageInSchema(t, schema.ID())
	t.Cleanup(stageCleanup)

	copyStatement := createPipeCopyStatement(t, table, stage)
	pipe, pipeCleanup := testClientHelper().Pipe.CreatePipe(t, copyStatement)
	t.Cleanup(pipeCleanup)

	pipeExecutionState, err := client.SystemFunctions.PipeStatus(pipe.ID())
	require.NoError(t, err)
	require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)

	// Pause the pipe
	ctx := context.Background()
	err = client.Pipes.Alter(ctx, pipe.ID(), &sdk.AlterPipeOptions{
		Set: &sdk.PipeSet{
			PipeExecutionPaused: sdk.Bool(true),
		},
	})
	require.NoError(t, err)

	pipeExecutionState, err = client.SystemFunctions.PipeStatus(pipe.ID())
	require.NoError(t, err)
	require.Equal(t, sdk.PausedPipeExecutionState, pipeExecutionState)

	// Unpause the pipe
	err = client.Pipes.Alter(ctx, pipe.ID(), &sdk.AlterPipeOptions{
		Set: &sdk.PipeSet{
			PipeExecutionPaused: sdk.Bool(false),
		},
	})
	require.NoError(t, err)

	pipeExecutionState, err = client.SystemFunctions.PipeStatus(pipe.ID())
	require.NoError(t, err)
	require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)
}

func TestInt_PipeForceResume(t *testing.T) {
	client := testClient(t)

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	table, tableCleanup := testClientHelper().Table.CreateInSchema(t, schema.ID())
	t.Cleanup(tableCleanup)

	stage, stageCleanup := testClientHelper().Stage.CreateStageInSchema(t, schema.ID())
	t.Cleanup(stageCleanup)

	copyStatement := createPipeCopyStatement(t, table, stage)
	pipe, pipeCleanup := testClientHelper().Pipe.CreatePipe(t, copyStatement)
	t.Cleanup(pipeCleanup)

	pipeExecutionState, err := client.SystemFunctions.PipeStatus(pipe.ID())
	require.NoError(t, err)
	require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)

	ctx := context.Background()
	err = client.Pipes.Alter(ctx, pipe.ID(), &sdk.AlterPipeOptions{
		Set: &sdk.PipeSet{
			PipeExecutionPaused: sdk.Bool(true),
		},
	})
	require.NoError(t, err)

	// Move the ownership to the role and back to the currently used role by the client
	err = client.Grants.GrantOwnership(
		ctx,
		sdk.OwnershipGrantOn{
			Object: &sdk.Object{
				ObjectType: sdk.ObjectTypePipe,
				Name:       pipe.ID(),
			},
		},
		sdk.OwnershipGrantTo{
			AccountRoleName: sdk.Pointer(role.ID()),
		},
		new(sdk.GrantOwnershipOptions),
	)
	require.NoError(t, err)

	currentRole := testClientHelper().Context.CurrentRole(t)

	err = client.Grants.GrantOwnership(
		ctx,
		sdk.OwnershipGrantOn{
			Object: &sdk.Object{
				ObjectType: sdk.ObjectTypePipe,
				Name:       pipe.ID(),
			},
		},
		sdk.OwnershipGrantTo{
			AccountRoleName: sdk.Pointer(currentRole),
		},
		new(sdk.GrantOwnershipOptions),
	)
	require.NoError(t, err)

	// Try to resume with ALTER (error)
	err = client.Pipes.Alter(ctx, pipe.ID(), &sdk.AlterPipeOptions{
		Set: &sdk.PipeSet{
			PipeExecutionPaused: sdk.Bool(false),
		},
	})
	require.ErrorContains(t, err, fmt.Sprintf("Pipe %s cannot be resumed as ownership had changed. Resuming pipe may load files inserted by previous owner into table. To forceresume pipe use SYSTEM$PIPE_FORCE_RESUME('%s')", pipe.Name, pipe.Name))

	// Resume with system func (success)
	err = client.SystemFunctions.PipeForceResume(pipe.ID(), nil)
	require.NoError(t, err)

	pipeExecutionState, err = client.SystemFunctions.PipeStatus(pipe.ID())
	require.NoError(t, err)
	require.Equal(t, sdk.RunningPipeExecutionState, pipeExecutionState)
}

func TestInt_BcrBundles(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("get bundle status", func(t *testing.T) {
		status, err := client.SystemFunctions.BehaviorChangeBundleStatus(ctx, "2025_01")
		require.NoError(t, err)
		assert.Equal(t, sdk.BehaviorChangeBundleStatusReleased, status)
	})
}
