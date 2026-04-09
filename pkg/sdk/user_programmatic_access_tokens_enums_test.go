package sdk

import "testing"

func Test_ToProgrammaticAccessTokenStatus(t *testing.T) {
	testEnumConversion(t, AllProgrammaticAccessTokenStatuses, ToProgrammaticAccessTokenStatus)
}
