package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToColumnConstraintType(t *testing.T) {
	type test struct {
		input string
		want  ColumnConstraintType
	}

	positiveTests := []test{
		{input: "UNIQUE", want: ColumnConstraintTypeUnique},
		{input: "unique", want: ColumnConstraintTypeUnique},
		{input: "uNiQuE", want: ColumnConstraintTypeUnique},
		{input: "PRIMARY KEY", want: ColumnConstraintTypePrimaryKey},
		{input: "primary key", want: ColumnConstraintTypePrimaryKey},
		{input: "PRIMARY key", want: ColumnConstraintTypePrimaryKey},
		{input: "FOREIGN KEY", want: ColumnConstraintTypeForeignKey},
		{input: "foreign key", want: ColumnConstraintTypeForeignKey},
		{input: "foreign KEY", want: ColumnConstraintTypeForeignKey},
	}

	negativeTests := []test{
		{input: "foreign key "},
		{input: "not null"},
		{input: "NOT NULL"},
		{input: "abc"},
		{input: ""},
	}

	for _, tc := range positiveTests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToColumnConstraintType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range negativeTests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToColumnConstraintType(tc.input)
			require.Error(t, err)
			require.Empty(t, got)
		})
	}
}

func Test_ToMatchType(t *testing.T) {
	type test struct {
		input string
		want  MatchType
	}

	positiveTests := []test{
		{input: string(FullMatchType), want: FullMatchType},
		{input: "FULL", want: FullMatchType},
		{input: "full", want: FullMatchType},
		{input: string(SimpleMatchType), want: SimpleMatchType},
		{input: "SIMPLE", want: SimpleMatchType},
		{input: "simple", want: SimpleMatchType},
		{input: string(PartialMatchType), want: PartialMatchType},
		{input: "PARTIAL", want: PartialMatchType},
		{input: "partial", want: PartialMatchType},
	}

	negativeTests := []test{
		{input: "full "},
		{input: " PARTIAL"},
		{input: "NOT NULL"},
		{input: "abc"},
		{input: ""},
	}

	for _, tc := range positiveTests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToMatchType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range negativeTests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToMatchType(tc.input)
			require.Error(t, err)
			require.Empty(t, got)
		})
	}
}

func Test_ToForeignKeyAction(t *testing.T) {
	type test struct {
		input string
		want  ForeignKeyAction
	}

	positiveTests := []test{
		{input: string(ForeignKeyCascadeAction), want: ForeignKeyCascadeAction},
		{input: "CASCADE", want: ForeignKeyCascadeAction},
		{input: "cascade", want: ForeignKeyCascadeAction},
		{input: string(ForeignKeySetNullAction), want: ForeignKeySetNullAction},
		{input: "SET NULL", want: ForeignKeySetNullAction},
		{input: "set null", want: ForeignKeySetNullAction},
		{input: string(ForeignKeySetDefaultAction), want: ForeignKeySetDefaultAction},
		{input: "SET DEFAULT", want: ForeignKeySetDefaultAction},
		{input: "set default", want: ForeignKeySetDefaultAction},
		{input: string(ForeignKeyRestrictAction), want: ForeignKeyRestrictAction},
		{input: "RESTRICT", want: ForeignKeyRestrictAction},
		{input: "restrict", want: ForeignKeyRestrictAction},
		{input: string(ForeignKeyNoAction), want: ForeignKeyNoAction},
		{input: "NO ACTION", want: ForeignKeyNoAction},
		{input: "no action", want: ForeignKeyNoAction},
	}

	negativeTests := []test{
		{input: "no action "},
		{input: " RESTRICT"},
		{input: "not null"},
		{input: "abc"},
		{input: ""},
	}

	for _, tc := range positiveTests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToForeignKeyAction(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range negativeTests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToForeignKeyAction(tc.input)
			require.Error(t, err)
			require.Empty(t, got)
		})
	}
}
