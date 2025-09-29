package model

import tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

func S3StorageIntegration(resourceName string, name string, awsRoleArn string) *StorageIntegrationModel {
	return StorageIntegration(resourceName, name, []string{"s3://foo/"}, "S3").
		WithStorageAwsRoleArn(awsRoleArn)
}

func (s *StorageIntegrationModel) WithStorageAllowedLocations(location []string) *StorageIntegrationModel {
	variables := make([]tfconfig.Variable, len(location))
	for i, v := range location {
		variables[i] = tfconfig.StringVariable(v)
	}
	s.StorageAllowedLocations = tfconfig.SetVariable(variables...)
	return s
}
