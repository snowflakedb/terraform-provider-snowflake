package model

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func PipeWithId(
	resourceName string,
	pipeId sdk.SchemaObjectIdentifier,
	copyStatement string,
) *PipeModel {
	return Pipe(resourceName, pipeId.DatabaseName(), pipeId.SchemaName(), pipeId.Name(), copyStatement)
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
