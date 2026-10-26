package defs

import (
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var externalVolumesDef = g.NewInterface(
	"ExternalVolumes",
	"ExternalVolume",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	CreateOperationFromObject(ssotdefs.ExternalVolumeObject).
	AlterOperationFromObject(ssotdefs.ExternalVolumeObject).
	DropOperationFromObject(ssotdefs.ExternalVolumeObject).
	DescribeOperationFromObject(ssotdefs.ExternalVolumeObject).
	ShowOperationFromObject(ssotdefs.ExternalVolumeObject).
	WithEnums(
		ssotdefs.StorageProvider,
		ssotdefs.S3StorageProvider,
		ssotdefs.S3EncryptionType,
		ssotdefs.GCSEncryptionType,
		ssotdefs.AzureEncryptionType,
		ssotdefs.S3CompatEncryptionType,
	)
