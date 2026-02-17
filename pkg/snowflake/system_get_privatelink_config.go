package snowflake

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

func SystemGetPrivateLinkConfigQuery() string {
	return `SELECT SYSTEM$GET_PRIVATELINK_CONFIG() AS "CONFIG"`
}

type RawPrivateLinkConfig struct {
	Config string `db:"CONFIG"`
}

type privateLinkConfigInternal struct {
	AccountName               string `json:"privatelink-account-name"`
	AwsVpceID                 string `json:"privatelink-vpce-id,omitempty"`
	AzurePrivateLinkServiceID string `json:"privatelink-pls-id,omitempty"`
	AccountURL                string `json:"privatelink-account-url"`
	OCSPURL                   string `json:"privatelink-ocsp-url,omitempty"`
	TypodOCSPURL              string `json:"privatelink_ocsp-url,omitempty"` // because snowflake returns this for AWS, but don't have an Azure account to verify against
	InternalStage             string `json:"privatelink-internal-stage,omitempty"`
	SnowsightURL              string `json:"snowsight-privatelink-url,omitempty"`
	RegionlessSnowsightURL    string `json:"regionless-snowsight-privatelink-url,omitempty"`
	RegionlessAccountURL      string `json:"regionless-privatelink-account-url,omitempty"`
	ConnectionURLs            string `json:"privatelink-connection-urls,omitempty"`
}

type PrivateLinkConfig struct {
	AccountName               string
	AwsVpceID                 string
	AzurePrivateLinkServiceID string
	AccountURL                string
	OCSPURL                   string
	InternalStage             string
	SnowsightURL              string
	RegionlessSnowsightURL    string
	RegionlessAccountURL      string
	ConnectionURLs            string
}

func ScanPrivateLinkConfig(row *sqlx.Row) (*RawPrivateLinkConfig, error) {
	config := &RawPrivateLinkConfig{}
	err := row.StructScan(config)
	return config, err
}

func (r *RawPrivateLinkConfig) GetStructuredConfig() (*PrivateLinkConfig, error) {
	config := &privateLinkConfigInternal{}
	err := json.Unmarshal([]byte(r.Config), config)
	if err != nil {
		return nil, err
	}

	return config.getPrivateLinkConfig()
}

func (i *privateLinkConfigInternal) getPrivateLinkConfig() (*PrivateLinkConfig, error) {
	config := &PrivateLinkConfig{
		i.AccountName,
		i.AwsVpceID,
		i.AzurePrivateLinkServiceID,
		i.AccountURL,
		i.OCSPURL,
		i.InternalStage,
		i.SnowsightURL,
		i.RegionlessSnowsightURL,
		i.RegionlessAccountURL,
		i.ConnectionURLs,
	}

	if i.TypodOCSPURL != "" {
		config.OCSPURL = i.TypodOCSPURL
	}

	return config, nil
}
