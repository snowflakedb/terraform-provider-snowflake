package testacc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
