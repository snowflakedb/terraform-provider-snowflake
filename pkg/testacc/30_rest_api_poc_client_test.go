package testacc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/snowflakedb/gosnowflake"
)

type RestApiPocConfig struct {
	Account string
	Token   string
}

func RestApiPocConfigFromDriverConfig(driverConfig *gosnowflake.Config) (*RestApiPocConfig, error) {
	res := &RestApiPocConfig{
		Account: driverConfig.Account,
	}
	if driverConfig.Token == "" {
		return nil, fmt.Errorf("token is currently required for REST API PoC client initialization")
	} else {
		res.Token = driverConfig.Token
	}

	return res, nil
}

func NewRestApiPocClient(config *RestApiPocConfig) (*RestApiPocClient, error) {
	c := &RestApiPocClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		url:        fmt.Sprintf("https://%s.snowflakecomputing.com/api/v2/", config.Account),
		token:      config.Token,
	}

	c.Warehouses = warehousesPoc{client: c}

	return c, nil
}

type RestApiPocClient struct {
	httpClient *http.Client
	url        string
	token      string

	Warehouses WarehousesPoc
}

func (c *RestApiPocClient) doRequest(ctx context.Context, method string, path string, body io.Reader, queryParams map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.url+"/"+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Add("X-Snowflake-Authorization-Token-Type", "PROGRAMMATIC_ACCESS_TOKEN")

	values := req.URL.Query()
	for k, v := range queryParams {
		values.Add(k, v)
	}
	req.URL.RawQuery = values.Encode()

	return c.httpClient.Do(req)
}

func post[T any](ctx context.Context, client *RestApiPocClient, path string, object T) (*Response, error) {
	return postOrPut(ctx, client, http.MethodPost, path, object)
}

func put[T any](ctx context.Context, client *RestApiPocClient, path string, object T) (*Response, error) {
	return postOrPut(ctx, client, http.MethodPut, path, object)
}

// TODO [this PR]: add status codes handling
func postOrPut[T any](ctx context.Context, client *RestApiPocClient, method string, path string, object T) (*Response, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(object)
	if err != nil {
		return nil, err
	}

	resp, err := client.doRequest(ctx, method, path, body, map[string]string{})
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

// TODO [this PR]: add status codes handling
func get[T any](ctx context.Context, client *RestApiPocClient, path string) (*T, error) {
	resp, err := client.doRequest(ctx, http.MethodGet, path, nil, map[string]string{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response T
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func runDelete(ctx context.Context, client *RestApiPocClient, path string, queryParams map[string]string) (*Response, error) {
	resp, err := client.doRequest(ctx, http.MethodDelete, path, nil, queryParams)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response Response
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}
