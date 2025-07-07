package testfunctional_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

var server *httptest.Server
var serverCleanup func()
var allTestHandlers = make(map[string]http.Handler)

type testHandler struct{}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	d, err := w.Write([]byte(`{"message": "example"}`))
	functionalTestLog.Printf("[DEBUG] Bytes written: %d, err: %v", d, err)
}

func setUpTestHttpServer() error {
	allTestHandlers["example_test"] = &testHandler{}

	mux := http.NewServeMux()
	for path, handler := range allTestHandlers {
		mux.Handle("/"+path, handler)
	}

	server = httptest.NewServer(mux)
	serverCleanup = func() {
		server.Close()
	}

	functionalTestLog.Printf("[INFO] Started a server at %s", server.URL)

	msg, err := fetchTest1Message(server.URL)
	if err != nil {
		functionalTestLog.Printf("[DEBUG] Connection error: %v", err)
		return fmt.Errorf("error fetching test message from test http server: %w", err)
	} else {
		functionalTestLog.Printf("[DEBUG] Test message received `%s`", msg)
	}
	return nil
}

func fetchTest1Message(baseUrl string) (string, error) {
	resp, err := http.Get(baseUrl + "/example_test")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(buf[:]), nil
}
