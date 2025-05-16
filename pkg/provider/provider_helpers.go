package provider

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/snowflakedb/gosnowflake"
)

type protocol string

const (
	// these values are lower case on purpose to match gosnowflake case
	protocolHttp  protocol = "http"
	protocolHttps protocol = "https"
)

var allProtocols = []protocol{
	protocolHttp,
	protocolHttps,
}

func toProtocol(s string) (protocol, error) {
	lowerCase := strings.ToLower(s)
	switch lowerCase {
	case string(protocolHttp),
		string(protocolHttps):
		return protocol(lowerCase), nil
	default:
		return "", fmt.Errorf("invalid protocol: %s", s)
	}
}

func getPrivateKey(privateKeyString, privateKeyPassphrase string) (*rsa.PrivateKey, error) {
	if privateKeyString == "" {
		return nil, nil
	}
	privateKeyBytes := []byte(privateKeyString)
	return sdk.ParsePrivateKey(privateKeyBytes, []byte(privateKeyPassphrase))
}

type GetRefreshTokenResponseBody struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func GetAccessTokenWithRefreshToken(
	tokenEndPoint string,
	clientID string,
	clientSecret string,
	refreshToken string,
	redirectURI string,
) (string, error) {
	client := &http.Client{}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("redirect_uri", redirectURI)
	body := strings.NewReader(data.Encode())

	request, err := http.NewRequest("POST", tokenEndPoint, body)
	if err != nil {
		return "", fmt.Errorf("preparing a HTTP request: %w", err)
	}
	request.SetBasicAuth(clientID, clientSecret)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("response status returned an err = %w", err)
	}
	if response.StatusCode != 200 {
		return "", fmt.Errorf("response status code: %s: %s", strconv.Itoa(response.StatusCode), http.StatusText(response.StatusCode))
	}
	defer response.Body.Close()
	dat, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("response body was not able to be parsed err = %w", err)
	}
	var result GetRefreshTokenResponseBody
	err = json.Unmarshal(dat, &result)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON from Snowflake err = %w", err)
	}
	return result.AccessToken, nil
}

func envNameFieldDescription(description, envName string) string {
	return fmt.Sprintf("%s Can also be sourced from the `%s` environment variable.", description, envName)
}

// TODO(SNOW-1787926): reuse these handlers with the ones in resources
func handleStringField(d *schema.ResourceData, key string, field *string) error {
	if v, ok := d.GetOk(key); ok {
		*field = v.(string)
	}
	return nil
}

func handleBoolField(d *schema.ResourceData, key string, field *bool) error {
	if v, ok := d.GetOk(key); ok {
		*field = v.(bool)
	}
	return nil
}

func handleDurationInSecondsAttribute(d *schema.ResourceData, key string, field *time.Duration) error {
	if v, ok := d.GetOk(key); ok {
		*field = time.Second * time.Duration(int64(v.(int)))
	}
	return nil
}

func handleIntAttribute(d *schema.ResourceData, key string, field *int) error {
	if v, ok := d.GetOk(key); ok {
		*field = v.(int)
	}
	return nil
}

func handleBooleanStringAttribute(d *schema.ResourceData, key string, field *gosnowflake.ConfigBool) error {
	if v := d.Get(key).(string); v != provider.BooleanDefault {
		parsed, err := provider.BooleanStringToBool(v)
		if err != nil {
			return err
		}
		if parsed {
			*field = gosnowflake.ConfigBoolTrue
		} else {
			*field = gosnowflake.ConfigBoolFalse
		}
	}
	return nil
}
