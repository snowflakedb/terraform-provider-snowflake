package common

type TestProviderContext struct {
	testName  string
	serverUrl string
}

func NewTestProviderContext(testName string, serverUrl string) *TestProviderContext {
	return &TestProviderContext{
		testName:  testName,
		serverUrl: serverUrl,
	}
}

func (c TestProviderContext) ServerUrl() string {
	return c.serverUrl
}
