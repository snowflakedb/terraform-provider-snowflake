package model

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func PipeWithId(
	resourceName string,
	pipeId sdk.SchemaObjectIdentifier,
	copyStatement string,
) *PipeModel {
	return Pipe(resourceName, pipeId.DatabaseName(), pipeId.SchemaName(), pipeId.Name(), copyStatement).WithCopyStatementMultiline(copyStatement)
}

// PipeWithIdCopyFromStageIntoTable is like PipeWithId but sets copy_statement to a COPY INTO table FROM @stage (CSV) statement.
// tableReference and stageReference are Terraform resource addresses (e.g. snowflake_table.test).
func PipeWithIdCopyFromStageIntoTable(
	resourceName string,
	pipeId sdk.SchemaObjectIdentifier,
	tableReference string,
	stageReference string,
) *PipeModel {
	return Pipe(
		resourceName,
		pipeId.DatabaseName(),
		pipeId.SchemaName(),
		pipeId.Name(),
		pipeCopyStatementFromStageIntoTable(tableReference, stageReference),
	)
}

func pipeCopyStatementFromStageIntoTable(tableReference string, stageReference string) string {
	return fmt.Sprintf(`COPY INTO ${%[1]s.fully_qualified_name}
FROM @${%[2]s.fully_qualified_name}
FILE_FORMAT = (TYPE = CSV)`, tableReference, stageReference)
}

// WithCopyStatementMultiline preserves quoted identifiers when converting the model to HCL.
func (p *PipeModel) WithCopyStatementMultiline(copyStatement string) *PipeModel {
	statement := strings.ReplaceAll(copyStatement, "\"", string(config.SnowflakeProviderConfigQuoteMarker))
	return p.WithCopyStatementValue(config.MultilineWrapperVariable(statement))
}
