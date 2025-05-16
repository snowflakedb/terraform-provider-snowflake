package genhelpers

import "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"

type SdkObjectDetails struct {
	IdType     string
	ObjectType sdk.ObjectType
	StructDetails
}
