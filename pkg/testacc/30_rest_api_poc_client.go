package testacc

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type RestApiPocConfig struct {
	// TODO [this PR]: fill config out (maybe from golang config?)
}

// TODO [this PR]: build url (account + /api/v2/)
func NewRestApiPocClient(config *RestApiPocConfig) (*RestApiPocClient, error) {
	c := &RestApiPocClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	c.Warehouses = warehousesPoc{client: c}

	return c, nil
}

type RestApiPocClient struct {
	httpClient *http.Client
	url        string

	Warehouses WarehousesPoc
}

// TODO [this PR]: token
// TODO [this PR]: move parts to client
// TODO [this PR]: add status codes handling
func post[T any](ctx context.Context, httpClient *http.Client, url string, path string, object T) (*Response, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(object)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/"+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer token")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := &Response{}
	if err = json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}
	return response, nil
}

type Response struct {
	State string `json:"state"`
}
