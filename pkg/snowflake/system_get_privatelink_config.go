package snowflake

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

func SystemGetPrivateLinkConfigQuery() string {
	return `SELECT SYSTEM$GET_PRIVATELINK_CONFIG() AS "config"`
}

type RawPrivateLinkConfig struct {
	Config string `db:"config"`
}

type privateLinkConfigInternal struct {
	AccountName               string `json:"privatelink-account-name"`
	AccountURL                string `json:"privatelink-account-url"`
	AppServiceURL             string `json:"app-service-privatelink-url,omitempty"`
	AwsVpceID                 string `json:"privatelink-vpce-id,omitempty"`
	AzurePrivateLinkServiceID string `json:"privatelink-pls-id,omitempty"`
	InternalStage             string `json:"privatelink-internal-stage,omitempty"`
	OCSPURL                   string `json:"privatelink-ocsp-url,omitempty"`
	OpenflowURL               string `json:"openflow-privatelink-url,omitempty"`
	OpenflowTelemetryURL      string `json:"external-telemetry-privatelink-url,omitempty"`
	RegionlessAccountURL      string `json:"regionless-privatelink-account-url,omitempty"`
	RegionlessSnowsightURL    string `json:"regionless-snowsight-privatelink-url,omitempty"`
	SnowparkCSAuthURL         string `json:"spcs-auth-privatelink-url,omitempty"`
	SnowparkCSRegistryURL     string `json:"spcs-registry-privatelink-url,omitempty"`
	SnowsightURL              string `json:"snowsight-privatelink-url,omitempty"`
	TypodOCSPURL              string `json:"privatelink_ocsp-url,omitempty"` // because snowflake returns this for AWS, but don't have an Azure account to verify against
}

type PrivateLinkConfig struct {
	AccountName               string
	AccountURL                string
	AppServiceURL             string
	AwsVpceID                 string
	AzurePrivateLinkServiceID string
	InternalStage             string
	OCSPURL                   string
	OpenflowURL               string
	OpenflowTelemetryURL      string
	RegionlessAccountURL      string
	RegionlessSnowsightURL    string
	SnowparkCSAuthURL         string
	SnowparkCSRegistryURL     string
	SnowsightURL              string
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
		i.AccountURL,
		i.AppServiceURL,
		i.AwsVpceID,
		i.AzurePrivateLinkServiceID,
		i.InternalStage,
		i.OCSPURL,
		i.OpenflowURL,
		i.OpenflowTelemetryURL,
		i.RegionlessAccountURL,
		i.RegionlessSnowsightURL,
		i.SnowparkCSAuthURL,
		i.SnowparkCSRegistryURL,
		i.SnowsightURL,
	}

	if i.TypodOCSPURL != "" {
		config.OCSPURL = i.TypodOCSPURL
	}

	return config, nil
}
