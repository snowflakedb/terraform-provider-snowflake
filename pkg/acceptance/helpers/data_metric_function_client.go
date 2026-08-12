package helpers

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// TODO(SNOW-1564959): change raw sqls to proper client
type DataMetricFunctionClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewDataMetricFunctionClient(context *TestClientContext, idsGenerator *IdsGenerator) *DataMetricFunctionClient {
	return &DataMetricFunctionClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *DataMetricFunctionClient) client() *sdk.Client {
	return c.context.client
}

func (c *DataMetricFunctionClient) CreateDataMetricFunction(t *testing.T, viewID sdk.SchemaObjectIdentifier) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`CREATE DATA METRIC FUNCTION %s(arg_t TABLE (arg_c INT))
RETURNS NUMBER AS
'SELECT COUNT(*) FROM arg_t
   WHERE arg_c IN (SELECT id FROM %s)'`, id.FullyQualifiedName(), viewID.FullyQualifiedName()))
	require.NoError(t, err)

	idWithArgs := sdk.NewSchemaObjectIdentifierWithArgumentsInSchema(id.SchemaId(), id.Name(), sdk.DataType("TABLE(INT)"))
	return id, c.DropFunc(t, idWithArgs)
}

// CreateWithArguments creates a data metric function with passed arguments via raw SQL (no SDK client yet).
// createArguments are argument definitions as written in CREATE DDL (e.g. []string{"ARG_T TABLE(ARG_C DATE)"}).
// signatureArguments are Snowflake's abbreviated signature form used in object identifiers returned by SHOW, DESC, SHOW GRANTS, etc. (e.g. []string{"TABLE(DATE)"}).
func (c *DataMetricFunctionClient) CreateWithArguments(t *testing.T, createArguments, signatureArguments []string) (sdk.SchemaObjectIdentifierWithArguments, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	_, err := c.client().ExecForTests(ctx, fmt.Sprintf(
		`CREATE OR REPLACE DATA METRIC FUNCTION %s(%s) RETURNS NUMBER AS $$ SELECT COUNT(*) FROM ARG_T $$`,
		id.FullyQualifiedName(),
		strings.Join(createArguments, ", "),
	))
	require.NoError(t, err)

	argumentDataTypes := make([]sdk.DataType, len(signatureArguments))
	for i, arg := range signatureArguments {
		argumentDataTypes[i] = sdk.DataType(arg)
	}
	idWithArgs := sdk.NewSchemaObjectIdentifierWithArgumentsInSchema(id.SchemaId(), id.Name(), argumentDataTypes...)

	return idWithArgs, c.DropFunc(t, idWithArgs)
}

func (c *DataMetricFunctionClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`DROP FUNCTION IF EXISTS %s`, id.FullyQualifiedName()))
		require.NoError(t, err)
	}
}
