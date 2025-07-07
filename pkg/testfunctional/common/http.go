package common

import (
	"bytes"
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

func Post[T any](serverUrl string, path string, target T) error {
	body, err := json.Marshal(target)
	if err != nil {
		return err
	}

	resp, err := http.Post(serverUrl+"/"+path, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
