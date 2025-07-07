package testfunctional_test

import (
	"io"
	"net/http"
	"net/http/httptest"
)

var server *httptest.Server
var serverCleanup func()
var allTestHandlers = make(map[string]http.Handler)

type test1Handler struct{}

func (h *test1Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	d, err := w.Write([]byte(`{"message": "test1"}`))
	functionalTestLog.Printf("[DEBUG] Bytes written: %d, err: %v", d, err)
}

// TODO [mux-PR]: move from init
func init() {
	mux := http.NewServeMux()

	allTestHandlers["test1"] = &test1Handler{}

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
	} else {
		functionalTestLog.Printf("[DEBUG] Test message received `%s`", msg)
	}
}

func fetchTest1Message(baseUrl string) (string, error) {
	resp, err := http.Get(baseUrl + "/test1")
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
