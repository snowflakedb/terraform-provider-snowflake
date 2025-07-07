package common

import (
	"encoding/json"
	"io"
	"net/http"
)

func Get[T any](serverUrl string, path string, target *T) error {
	resp, err := http.Get(serverUrl + "/" + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, &target)
}
