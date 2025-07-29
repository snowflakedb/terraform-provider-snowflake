package testacc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/snowflakedb/gosnowflake"
)

type RestApiPocConfig struct {
	Account string
	Token   string
}

func RestApiPocConfigFromDriverConfig(driverConfig *gosnowflake.Config) (*RestApiPocConfig, error) {
	res := &RestApiPocConfig{
		Account: strings.ToLower(driverConfig.Account),
	}
	if driverConfig.Token == "" {
		return nil, fmt.Errorf("token is currently required for REST API PoC client initialization")
	} else {
		res.Token = driverConfig.Token
	}

	return res, nil
}

// TODO [mux-PR]: verify connection after creation
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
	req, err := http.NewRequestWithContext(ctx, method, c.url+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("X-Snowflake-Authorization-Token-Type", "PROGRAMMATIC_ACCESS_TOKEN")
	req.Header.Set("Accept", "application/json")

	values := req.URL.Query()
	for k, v := range queryParams {
		values.Add(k, v)
	}
	req.URL.RawQuery = values.Encode()
	accTestLog.Printf("[DEBUG] Sending request [%s] %s", method, req.URL)

	return c.httpClient.Do(req)
}

func post[T any](ctx context.Context, client *RestApiPocClient, path string, object T) (*Response, error) {
	return postOrPut(ctx, client, http.MethodPost, path, object)
}

func put[T any](ctx context.Context, client *RestApiPocClient, path string, object T) (*Response, error) {
	return postOrPut(ctx, client, http.MethodPut, path, object)
}

// TODO [mux-PR]: potentially merge postOrPut, get, and handleDelete
// TODO [mux-PR]: improve status codes handling
func postOrPut[T any](ctx context.Context, client *RestApiPocClient, method string, path string, object T) (*Response, error) {
	body, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	resp, err := client.doRequest(ctx, method, path, bytes.NewBuffer(body), map[string]string{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	accTestLog.Printf("[DEBUG] Response status for request [%s] %s: %s", method, resp.Request.URL, resp.Status)

	response := &Response{}
	if err = json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		d, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, fmt.Errorf("unexpected status code: %d, response: %v", resp.StatusCode, response)
		}
		return nil, fmt.Errorf("unexpected status code: %d, response: %v, dump: %q", resp.StatusCode, response, d)
	}

	return response, nil
}

type Response struct {
	State   string `json:"state"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// TODO [mux-PR]: improve status codes handling
func get[T any](ctx context.Context, client *RestApiPocClient, path string) (*T, error) {
	method := http.MethodGet
	resp, err := client.doRequest(ctx, method, path, nil, map[string]string{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	accTestLog.Printf("[DEBUG] Response status for request [%s] %s: %s", method, resp.Request.URL, resp.Status)

	if resp.StatusCode == http.StatusNotFound {
		// using the existing SDK error for now
		return nil, sdk.ErrObjectNotFound
	}

	if resp.StatusCode != http.StatusOK {
		d, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("unexpected status code: %d, dump: %q", resp.StatusCode, d)
	}

	var response T
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

// TODO [mux-PR]: improve status codes handling
func handleDelete(ctx context.Context, client *RestApiPocClient, path string, queryParams map[string]string) (*Response, error) {
	method := http.MethodDelete
	resp, err := client.doRequest(ctx, method, path, nil, queryParams)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	accTestLog.Printf("[DEBUG] Response status for request [%s] %s: %s", method, resp.Request.URL, resp.Status)

	response := &Response{}
	if err = json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		d, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, fmt.Errorf("unexpected status code: %d, response: %v", resp.StatusCode, response)
		}
		return nil, fmt.Errorf("unexpected status code: %d, response: %v, dump: %q", resp.StatusCode, response, d)
	}

	return response, nil
}
