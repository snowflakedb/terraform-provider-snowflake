package datasources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Services_ToShowServicesType(t *testing.T) {
	type test struct {
		input string
		want  ShowServicesType
	}

	valid := []test{
		// case insensitive.
		{input: "all", want: ShowServicesTypeAll},

		// Supported Values
		{input: "ALL", want: ShowServicesTypeAll},
		{input: "JOBS_ONLY", want: ShowServicesTypeJobsOnly},
		{input: "SERVICES_ONLY", want: ShowServicesTypeServicesOnly},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := toShowServicesType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := toShowServicesType(tc.input)
			require.Error(t, err)
		})
	}
}
