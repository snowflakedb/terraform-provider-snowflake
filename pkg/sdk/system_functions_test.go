package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToBehaviorChangeBundleStatus(t *testing.T) {
	type test struct {
		input string
		want  BehaviorChangeBundleStatus
	}

	valid := []test{
		{input: "enabled", want: BehaviorChangeBundleStatusEnabled},
		{input: "ENABLED", want: BehaviorChangeBundleStatusEnabled},
		{input: "DISABLED", want: BehaviorChangeBundleStatusDisabled},
		{input: "RELEASED", want: BehaviorChangeBundleStatusReleased},
	}

	invalid := []string{
		"",
		"foo",
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToBehaviorChangeBundleStatus(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, in := range invalid {
		t.Run(in, func(t *testing.T) {
			_, err := ToBehaviorChangeBundleStatus(in)
			require.Error(t, err)
		})
	}
}

func Test_parseClusteringInformation(t *testing.T) {
	t.Run("valid output", func(t *testing.T) {
		// Output captured from SYSTEM$CLUSTERING_INFORMATION on a clustered table.
		raw := `{
  "cluster_by_keys" : "LINEAR(REGION)",
  "total_partition_count" : 1,
  "total_constant_partition_count" : 0,
  "average_overlaps" : 0.0,
  "average_depth" : 1.0,
  "partition_depth_histogram" : {
    "00000" : 0,
    "00001" : 1,
    "00002" : 0,
    "00016" : 0
  },
  "clustering_errors" : [ { "timestamp" : "2024-01-01 00:00:00", "error" : "some error" } ]
}`
		info, err := parseClusteringInformation(raw)
		require.NoError(t, err)
		require.NotNil(t, info)
		require.Equal(t, "LINEAR(REGION)", info.ClusterByKeys)
		require.Equal(t, 1, info.TotalPartitionCount)
		require.Equal(t, 0, info.TotalConstantPartitionCount)
		require.InDelta(t, 0.0, info.AverageOverlaps, 0.0001)
		require.InDelta(t, 1.0, info.AverageDepth, 0.0001)
		require.Equal(t, 1, info.PartitionDepthHistogram["00001"])
		require.Len(t, info.ClusteringErrors, 1)
		require.Equal(t, "2024-01-01 00:00:00", info.ClusteringErrors[0].Timestamp)
		require.Equal(t, "some error", info.ClusteringErrors[0].Error)
	})

	t.Run("invalid json", func(t *testing.T) {
		_, err := parseClusteringInformation("not json")
		require.Error(t, err)
	})
}
