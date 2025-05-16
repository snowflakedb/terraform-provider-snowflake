package testenvs_test

import (
	"sync"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/stretchr/testify/require"
)

func Test_GetOrSkipTest(t *testing.T) {
	// runGetOrSkipInGoroutineAndWaitForCompletion is needed because underneath we test t.Skipf, that leads to t.SkipNow() that in turn call runtime.Goexit()
	// so we need to be wrapped in a Goroutine.
	runGetOrSkipInGoroutineAndWaitForCompletion := func(t *testing.T) string {
		t.Helper()
		var env string
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			env = testenvs.GetOrSkipTest(t, testenvs.BusinessCriticalAccount)
		}()
		wg.Wait()
		return env
	}

	t.Run("skip test if missing", func(t *testing.T) {
		t.Setenv(string(testenvs.BusinessCriticalAccount), "")

		tut := &testing.T{}
		env := runGetOrSkipInGoroutineAndWaitForCompletion(tut)

		require.True(t, tut.Skipped())
		require.Empty(t, env)
	})

	t.Run("get env if exists", func(t *testing.T) {
		t.Setenv(string(testenvs.BusinessCriticalAccount), "user")

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
		wg.Add(1)
		go func() {
			defer wg.Done()
			testenvs.SkipTestIfSet(t, env, "some good reason")
		}()
		wg.Wait()
	}

	t.Run("skip test if env is set", func(t *testing.T) {
		t.Setenv(string(testenvs.BusinessCriticalAccount), "1")

		tut := &testing.T{}
		runSkipTestIfSetInGoroutineAndWaitForCompletion(tut, testenvs.BusinessCriticalAccount)

		require.True(t, tut.Skipped())
	})

	t.Run("do not skip if env not set", func(t *testing.T) {
		t.Setenv(string(testenvs.BusinessCriticalAccount), "")

		tut := &testing.T{}
		runSkipTestIfSetInGoroutineAndWaitForCompletion(tut, testenvs.BusinessCriticalAccount)

		require.False(t, tut.Skipped())
	})
}

func Test_Assertions(t *testing.T) {
	// runAssertionInGoroutineAndWaitForCompletion is needed because underneath we test require, that leads to t.FailNow() that in turn call runtime.Goexit()
	// so we need to be wrapped in a Goroutine.
	runAssertionInGoroutineAndWaitForCompletion := func(assertion func()) {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			assertion()
		}()
		wg.Wait()
	}

	t.Run("test if env does not exist", func(t *testing.T) {
		t.Setenv(string(testenvs.BusinessCriticalAccount), "")

		tut1 := &testing.T{}
		runAssertionInGoroutineAndWaitForCompletion(func() { testenvs.AssertEnvNotSet(tut1, string(testenvs.BusinessCriticalAccount)) })

		tut2 := &testing.T{}
		runAssertionInGoroutineAndWaitForCompletion(func() { testenvs.AssertEnvSet(tut2, string(testenvs.BusinessCriticalAccount)) })

		require.False(t, tut1.Failed())
		require.True(t, tut2.Failed())
	})

	t.Run("test if env exists", func(t *testing.T) {
		t.Setenv(string(testenvs.BusinessCriticalAccount), "user")

		tut1 := &testing.T{}
		runAssertionInGoroutineAndWaitForCompletion(func() { testenvs.AssertEnvNotSet(tut1, string(testenvs.BusinessCriticalAccount)) })

		tut2 := &testing.T{}
		runAssertionInGoroutineAndWaitForCompletion(func() { testenvs.AssertEnvSet(tut2, string(testenvs.BusinessCriticalAccount)) })

		require.True(t, tut1.Failed())
		require.False(t, tut2.Failed())
	})
}
