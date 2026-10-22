package sdk

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testTasks struct {
	tasks

	stubbedTasks map[string]*Task
}

func (v *testTasks) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Task, error) {
	t, ok := v.stubbedTasks[id.Name()]
	if !ok {
		return nil, errors.New("no task configured, check test config")
	}
	return t, nil
}

func TestTasks_GetRootTasks(t *testing.T) {
	db := "database"
	sc := "schema"

	setUpTasks := func(p map[string][]string) map[string]*Task {
		r := make(map[string]*Task)
		for k, v := range p {
			if k == "initial" || k == "expected" {
				continue
			}
			t := Task{DatabaseName: db, SchemaName: sc, Name: k}
			predecessors := make([]SchemaObjectIdentifier, len(v))
			for i, predecessor := range v {
				predecessors[i] = NewSchemaObjectIdentifier(db, sc, predecessor)
			}
			t.Predecessors = predecessors
			r[k] = &t
		}
		return r
	}

	// To increase readability tests are defined as maps (anonymous structs looked much worse in this case).
	// What are the contents of the map:
	// - key "initial" -> one element list with the task for which we will be getting root tasks
	// - key "expected" -> list of the expected root tasks
	// - any other key is considered as a task definition, that contains all direct predecessors (empty list for root task).
	tests := []map[string][]string{
		{"t1": {}, "initial": {"t1"}, "expected": {"t1"}},
		{"t1": {"t2"}, "t2": {"t3"}, "t3": {}, "initial": {"t1"}, "expected": {"t3"}},
		{"t1": {"t2", "t3"}, "t2": {"t3"}, "t3": {}, "initial": {"t1"}, "expected": {"t3"}},
		{"t1": {"t2", "t3"}, "t2": {}, "t3": {}, "initial": {"t1"}, "expected": {"t2", "t3"}},
		{"t1": {}, "t2": {}, "initial": {"t1"}, "expected": {"t1"}},
		{"t1": {"t2", "t3", "t4"}, "t2": {}, "t3": {}, "t4": {}, "initial": {"t1"}, "expected": {"t2", "t3", "t4"}},
		{"t1": {"t2", "t3", "t4"}, "t2": {}, "t3": {"t2"}, "t4": {"t3"}, "initial": {"t1"}, "expected": {"t2"}},
		{"r": {}, "t1": {"t2", "r"}, "t2": {"t3"}, "t3": {"t1"}, "initial": {"t1"}, "expected": {"r"}}, // cycle -> failing for the old implementation
		{"r": {}, "t1": {"t2", "r"}, "t2": {"t3"}, "t3": {"t1"}, "initial": {"t3"}, "expected": {"r"}}, // cycle -> failing for the old implementation
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test case [%v]", i), func(t *testing.T) {
			ctx := context.Background()
			initial, ok := tt["initial"]
			if !ok {
				t.FailNow()
			}
			expected, ok := tt["expected"]
			if !ok {
				t.FailNow()
			}
			client := new(testTasks)
			client.stubbedTasks = setUpTasks(tt)

			rootTasks, err := GetRootTasks(client, ctx, NewSchemaObjectIdentifier(db, sc, initial[0]))
			require.NoError(t, err)
			for _, v := range rootTasks {
				assert.Contains(t, expected, v.Name)
			}
			require.Len(t, rootTasks, len(expected))
		})
	}
}

func Test_getPredecessors(t *testing.T) {
	special := "!@#$%&*+-_=?:;,.|(){}<>"

	tests := []struct {
		predecessorsRaw      string
		expectedPredecessors []string
	}{
		{predecessorsRaw: "[]", expectedPredecessors: []string{}},
		{predecessorsRaw: "[\n  \"\\\"qgb)Z1KcNWJ(\\\".\\\"glN@JtR=7dzP$7\\\".\\\"Ls.T7-(bt{.lWd@DRWkyA6<6hNdh\\\"\"\n]", expectedPredecessors: []string{"Ls.T7-(bt{.lWd@DRWkyA6<6hNdh"}},
		{predecessorsRaw: "[\n  \"\\\"qgb)Z1KcNWJ(\\\".\\\"glN@JtR=7dzP$7\\\".Ls.T7-(bt{.lWd@DRWkyA6<6hNdh\"\n]", expectedPredecessors: []string{"Ls.T7-(bt{.lWd@DRWkyA6<6hNdh"}},
		{predecessorsRaw: fmt.Sprintf("[\n  \"\\\"a\\\".\\\"b\\\".\\\"%s\\\"\"\n]", special), expectedPredecessors: []string{special}},
		{predecessorsRaw: "[\n  \"\\\"a\\\".\\\"b\\\".\\\"c\\\"\",\"\\\"a\\\".\\\"b\\\".\\\"d\\\"\",\"\\\"a\\\".\\\"b\\\".\\\"e\\\"\"\n]", expectedPredecessors: []string{"c", "d", "e"}},
		{predecessorsRaw: `["\"a\".\"b\".\".PHo,k:%Sz8tdx,9?23xTsgHLYxe\""]`, expectedPredecessors: []string{".PHo,k:%Sz8tdx,9?23xTsgHLYxe"}},
		{predecessorsRaw: `["MY_DB.MY_SCH.MY_PARENT_TASK"]`, expectedPredecessors: []string{"MY_PARENT_TASK"}},
		{predecessorsRaw: `["CTG_DEV_PLK.INGESTION.COPY_FROM_KINESIS_S3_DELIVERY"]`, expectedPredecessors: []string{"COPY_FROM_KINESIS_S3_DELIVERY"}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test number %d for input: [%s]", i, tt.predecessorsRaw), func(t *testing.T) {
			got, err := getPredecessors(tt.predecessorsRaw)
			require.NoError(t, err)
			require.Equal(t, tt.expectedPredecessors, got)
		})
	}

	t.Run("incorrect json", func(t *testing.T) {
		_, err := getPredecessors("[{]")
		require.ErrorContains(t, err, "invalid character ']'")
	})
}
