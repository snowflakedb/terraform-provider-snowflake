package resourceshowoutputassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func CatalogIntegrationsDatasourceAwsGlueDescribeOutput(t *testing.T, datasourceReference string) *CatalogIntegrationAwsGlueDescribeOutputAssert {
	t.Helper()
	c := CatalogIntegrationAwsGlueDescribeOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "describe_output", "catalog_integrations.0."),
	}
	c.AddAssertion(assert.ValueSet("describe_output.#", "1"))
	return &c
}

func CatalogIntegrationsDatasourceObjectStorageDescribeOutput(t *testing.T, datasourceReference string) *CatalogIntegrationObjectStorageDescribeOutputAssert {
	t.Helper()
	c := CatalogIntegrationObjectStorageDescribeOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "describe_output", "catalog_integrations.0."),
	}
	c.AddAssertion(assert.ValueSet("describe_output.#", "1"))
	return &c
}

func CatalogIntegrationsDatasourceOpenCatalogDescribeOutput(t *testing.T, datasourceReference string) *CatalogIntegrationOpenCatalogDescribeOutputAssert {
	t.Helper()
	c := CatalogIntegrationOpenCatalogDescribeOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "describe_output", "catalog_integrations.0."),
	}
	c.AddAssertion(assert.ValueSet("describe_output.#", "1"))
	return &c
}

func CatalogIntegrationsDatasourceIcebergRestDescribeOutput(t *testing.T, datasourceReference string) *CatalogIntegrationIcebergRestDescribeOutputAssert {
	t.Helper()
	c := CatalogIntegrationIcebergRestDescribeOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "describe_output", "catalog_integrations.0."),
	}
	c.AddAssertion(assert.ValueSet("describe_output.#", "1"))
	return &c
}

func OAuthRestAuthenticationDatasourceDescribeOutput(t *testing.T, datasourceReference string, containingField string) *OAuthRestAuthenticationDescribeOutputAssert {
	t.Helper()

	o := OAuthRestAuthenticationDescribeOutputAssert{
		ResourceAssert:  assert.NewDatasourceAssert(datasourceReference, "describe_output.0."+containingField, "catalog_integrations.0."),
		containingField: containingField,
	}
	o.AddAssertion(assert.ValueSet(fmt.Sprintf("describe_output.0.%s.#", containingField), "1"))
	return &o
}

func SigV4RestAuthenticationDatasourceDescribeOutput(t *testing.T, datasourceReference string) *SigV4RestAuthenticationDescribeOutputAssert {
	t.Helper()

	s := SigV4RestAuthenticationDescribeOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "describe_output.0.sigv4_rest_authentication", "catalog_integrations.0."),
	}
	s.AddAssertion(assert.ValueSet("describe_output.0.sigv4_rest_authentication.#", "1"))
	return &s
}

func OpenCatalogRestConfigDatasourceDescribeOutput(t *testing.T, datasourceReference string) *OpenCatalogRestConfigDescribeOutputAssert {
	t.Helper()

	o := OpenCatalogRestConfigDescribeOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "describe_output.0.rest_config", "catalog_integrations.0."),
	}
	o.AddAssertion(assert.ValueSet("describe_output.0.rest_config.#", "1"))
	return &o
}

func IcebergRestRestConfigDatasourceDescribeOutput(t *testing.T, datasourceReference string) *IcebergRestRestConfigDescribeOutputAssert {
	t.Helper()

	i := IcebergRestRestConfigDescribeOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "describe_output.0.rest_config", "catalog_integrations.0."),
	}
	i.AddAssertion(assert.ValueSet("describe_output.0.rest_config.#", "1"))
	return &i
}
