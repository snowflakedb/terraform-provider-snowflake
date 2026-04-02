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

func (p *PipeModel) WithCopyStatementCopyFromStageIntoTable(tableReference string, stageReference string) *PipeModel {
	copyStatement := fmt.Sprintf(`COPY INTO ${%[1]s.fully_qualified_name}
FROM @${%[2]s.fully_qualified_name}
FILE_FORMAT = (TYPE = CSV)`, tableReference, stageReference)
	return p.WithCopyStatement(copyStatement)
}
