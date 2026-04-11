package model

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

func (i *ImageRepositoryModel) WithEncryptionEnum(encryptionType sdk.ImageRepositoryEncryptionType) *ImageRepositoryModel {
	return i.WithEncryption(string(encryptionType))
}
