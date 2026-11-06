package testenvs_test

import (
	"sync"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/stretchr/testify/require"
)

func Test_GetOrSkipTest(t *testing.T) {
	// runGetOrSkipInGoroutineAndWaitForCompletion is needed because underneath we test t.Skipf, that leads to t.SkipNow() that in turn call runtime.Goexit()
	// so we need to be wrapped in a Goroutine.
	runGetOrSkipInGoroutineAndWaitForCompletion := func(t *testing.T) string {
		t.Helper()
		var env string
		var wg sync.WaitGroup
		wg.Go(func() {
			env = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)
		})
		wg.Wait()
		return env
	}

	t.Run("skip test if missing", func(t *testing.T) {
		t.Setenv(string(testenvs.TestAccountCreate), "")

		tut := &testing.T{}
		env := runGetOrSkipInGoroutineAndWaitForCompletion(tut)

		require.True(t, tut.Skipped())
		require.Empty(t, env)
	})

	t.Run("get env if exists", func(t *testing.T) {
		t.Setenv(string(testenvs.TestAccountCreate), "user")

		tut := &testing.T{}
		env := runGetOrSkipInGoroutineAndWaitForCompletion(tut)

		require.False(t, tut.Skipped())
		require.Equal(t, "user", env)
	})
}

func Test_SkipTestIfSet(t *testing.T) {
	// runSkipTestIfSetInGoroutineAndWaitForCompletion is needed because underneath we test t.Skipf, that leads to t.SkipNow() that in turn call runtime.Goexit()
	// so we need to be wrapped in a Goroutine.
	runSkipTestIfSetInGoroutineAndWaitForCompletion := func(t *testing.T, env testenvs.Env) {
		t.Helper()
		var wg sync.WaitGroup
		wg.Go(func() {
			testenvs.SkipTestIfSet(t, env, "some good reason")
		})
		wg.Wait()
	}

	t.Run("skip test if env is set", func(t *testing.T) {
		t.Setenv(string(testenvs.TestAccountCreate), "1")

		tut := &testing.T{}
		runSkipTestIfSetInGoroutineAndWaitForCompletion(tut, testenvs.TestAccountCreate)

		require.True(t, tut.Skipped())
	})

	t.Run("do not skip if env not set", func(t *testing.T) {
		t.Setenv(string(testenvs.TestAccountCreate), "")

		tut := &testing.T{}
		runSkipTestIfSetInGoroutineAndWaitForCompletion(tut, testenvs.TestAccountCreate)

		require.False(t, tut.Skipped())
	})
}

func Test_SkipTestIfSetTo(t *testing.T) {
	// runSkipTestIfSetToInGoroutineAndWaitForCompletion is needed because underneath we test t.Skipf, that leads to t.SkipNow() that in turn call runtime.Goexit()
	// so we need to be wrapped in a Goroutine.
	runSkipTestIfSetToInGoroutineAndWaitForCompletion := func(t *testing.T, env testenvs.Env, value string) {
		t.Helper()
		var wg sync.WaitGroup
		wg.Go(func() {
			testenvs.SkipTestIfSetTo(t, env, value, "some good reason")
		})
		wg.Wait()
	}

	t.Run("skip test if env is set to specific value", func(t *testing.T) {
		t.Setenv(string(testenvs.TestAccountCreate), "specific_value")

		tut := &testing.T{}
		runSkipTestIfSetToInGoroutineAndWaitForCompletion(tut, testenvs.TestAccountCreate, "specific_value")

		require.True(t, tut.Skipped())
	})

	t.Run("do not skip if env is set to different value", func(t *testing.T) {
		t.Setenv(string(testenvs.TestAccountCreate), "different_value")

		tut := &testing.T{}
		runSkipTestIfSetToInGoroutineAndWaitForCompletion(tut, testenvs.TestAccountCreate, "specific_value")

		require.False(t, tut.Skipped())
	})

	t.Run("do not skip if env not set", func(t *testing.T) {
		t.Setenv(string(testenvs.TestAccountCreate), "")

		tut := &testing.T{}
		runSkipTestIfSetToInGoroutineAndWaitForCompletion(tut, testenvs.TestAccountCreate, "specific_value")

		require.False(t, tut.Skipped())
	})
}

func Test_SkipTestIfValueIn(t *testing.T) {
	// runSkipTestIfValueInInGoroutineAndWaitForCompletion is needed because underneath we test t.Skipf, that leads to t.SkipNow() that in turn call runtime.Goexit()
	// so we need to be wrapped in a Goroutine.
	runSkipTestIfValueInInGoroutineAndWaitForCompletion := func(t *testing.T, env testenvs.Env, values []string) {
		t.Helper()
		var wg sync.WaitGroup
		wg.Go(func() {
			testenvs.SkipTestIfValueIn(t, env, values, "some good reason")
		})
		wg.Wait()
	}

	t.Run("skip test if env is set to one of the values", func(t *testing.T) {
		t.Setenv(string(testenvs.TestAccountCreate), "value2")

		tut := &testing.T{}
		runSkipTestIfValueInInGoroutineAndWaitForCompletion(tut, testenvs.TestAccountCreate, []string{"value1", "value2", "value3"})

		require.True(t, tut.Skipped())
	})

	t.Run("do not skip if env is set to value not in list", func(t *testing.T) {
		t.Setenv(string(testenvs.TestAccountCreate), "different_value")

		tut := &testing.T{}
		runSkipTestIfValueInInGoroutineAndWaitForCompletion(tut, testenvs.TestAccountCreate, []string{"value1", "value2", "value3"})

		require.False(t, tut.Skipped())
	})

	t.Run("do not skip if env not set", func(t *testing.T) {
		t.Setenv(string(testenvs.TestAccountCreate), "")

		tut := &testing.T{}
		runSkipTestIfValueInInGoroutineAndWaitForCompletion(tut, testenvs.TestAccountCreate, []string{"value1", "value2", "value3"})

		require.False(t, tut.Skipped())
	})
}

func Test_Assertions(t *testing.T) {
	// runAssertionInGoroutineAndWaitForCompletion is needed because underneath we test require, that leads to t.FailNow() that in turn call runtime.Goexit()
	// so we need to be wrapped in a Goroutine.
	runAssertionInGoroutineAndWaitForCompletion := func(assertion func()) {
		var wg sync.WaitGroup
		wg.Go(func() {
			assertion()
		})
		wg.Wait()
	}

	t.Run("test if env does not exist", func(t *testing.T) {
		t.Setenv(string(testenvs.TestAccountCreate), "")

		tut1 := &testing.T{}
		runAssertionInGoroutineAndWaitForCompletion(func() { testenvs.AssertEnvNotSet(tut1, string(testenvs.TestAccountCreate)) })

		tut2 := &testing.T{}
		runAssertionInGoroutineAndWaitForCompletion(func() { testenvs.AssertEnvSet(tut2, string(testenvs.TestAccountCreate)) })

		require.False(t, tut1.Failed())
		require.True(t, tut2.Failed())
	})

	t.Run("test if env exists", func(t *testing.T) {
		t.Setenv(string(testenvs.TestAccountCreate), "user")

		tut1 := &testing.T{}
		runAssertionInGoroutineAndWaitForCompletion(func() { testenvs.AssertEnvNotSet(tut1, string(testenvs.TestAccountCreate)) })

		tut2 := &testing.T{}
		runAssertionInGoroutineAndWaitForCompletion(func() { testenvs.AssertEnvSet(tut2, string(testenvs.TestAccountCreate)) })

		require.True(t, tut1.Failed())
		require.False(t, tut2.Failed())
	})
}
