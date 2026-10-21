package defs

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_ComputePoolAttributeSql(t *testing.T) {
	require.NotPanics(t, func() { ssot.ValidateObject(ComputePool) })

	t.Run("for application is create identifier", func(t *testing.T) {
		attr := ComputePoolAttrs.ForApplication
		require.Equal(t, ComputePool.IdType, attr.CreateSql.IdType)
		require.Equal(t, ssot.SqlKindIdentifier, attr.CreateSql.Kind)
		require.Equal(t, "FOR APPLICATION", attr.CreateSql.SQLText(attr))
		require.False(t, attr.SetSql.Present)
		require.False(t, attr.UnsetSql.Present)
	})

	t.Run("required numbers are optional on set", func(t *testing.T) {
		for _, attr := range []ssot.Attribute{ComputePoolAttrs.MinNodes, ComputePoolAttrs.MaxNodes} {
			require.True(t, attr.CreateSql.Present)
			require.Nil(t, attr.CreateSql.Required)
			require.Equal(t, ssot.SqlKindParameter, attr.CreateSql.Kind)
			require.Empty(t, attr.CreateSql.Quotes)
			require.Equal(t, ssot.SqlKindParameter, attr.SetSql.Kind)
			require.NotNil(t, attr.SetSql.Required)
			require.False(t, *attr.SetSql.Required)
			require.False(t, attr.UnsetSql.Present)
		}
	})

	t.Run("instance family and initially suspended are create only", func(t *testing.T) {
		family := ComputePoolAttrs.InstanceFamily
		require.Equal(t, ssot.SqlKindParameter, family.CreateSql.Kind)
		require.Equal(t, ssot.SqlQuotesNo, family.CreateSql.Quotes)
		require.Nil(t, family.CreateSql.Required)
		require.False(t, family.SetSql.Present)
		require.False(t, family.UnsetSql.Present)

		suspended := ComputePoolAttrs.InitiallySuspended
		require.Equal(t, ssot.SqlKindParameter, suspended.CreateSql.Kind)
		require.False(t, suspended.SetSql.Present)
		require.False(t, suspended.UnsetSql.Present)
	})

	t.Run("unset keywords match the value sql text", func(t *testing.T) {
		for _, attr := range []ssot.Attribute{
			ComputePoolAttrs.AutoResume,
			ComputePoolAttrs.AutoSuspendSecs,
			ComputePoolAttrs.BackupInstanceFamilies,
			ComputePoolAttrs.Comment,
		} {
			require.True(t, attr.CreateSql.Present)
			require.True(t, attr.SetSql.Present)
			require.Equal(t, ssot.SqlKindKeyword, attr.UnsetSql.Kind)
			require.Equal(t, attr.CreateSql.SQLText(attr), attr.UnsetSql.SQLText(attr))
		}
	})

	t.Run("comment quotes and backup list item", func(t *testing.T) {
		comment := ComputePoolAttrs.Comment
		require.Equal(t, ssot.SqlQuotesSingle, comment.CreateSql.Quotes)
		require.Equal(t, ssot.SqlQuotesSingle, comment.SetSql.Quotes)

		backup := ComputePoolAttrs.BackupInstanceFamilies
		require.Equal(t, "ComputePoolBackupInstanceFamilyListItem", backup.CreateSql.ItemStruct)
		require.Equal(t, ssot.SqlParensYes, backup.CreateSql.Parentheses)
		require.Equal(t, ssot.SqlParensYes, backup.SetSql.Parentheses)
	})

	t.Run("object attributes keep the clauses", func(t *testing.T) {
		byName := make(map[string]ssot.Attribute, len(ComputePool.Attributes))
		for _, attr := range ComputePool.Attributes {
			byName[attr.Name] = attr
		}
		require.Equal(t, ComputePoolAttrs.MinNodes.CreateSql, byName["min_nodes"].CreateSql)
		require.Equal(t, ComputePoolAttrs.MinNodes.SetSql, byName["min_nodes"].SetSql)
		require.Equal(t, ComputePoolAttrs.ForApplication.CreateSql, byName["for_application"].CreateSql)
		require.False(t, ComputePoolAttrs.Name.CreateSql.Present)
	})
}
