package common

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Get[T any](serverUrl string, path string, target *T) error {
	resp, err := http.Get(serverUrl + "/" + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func Post[T any](serverUrl string, path string, target T) error {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(target)
	if err != nil {
		return err
	}

	resp, err := http.Post(serverUrl+"/"+path, "application/json", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
