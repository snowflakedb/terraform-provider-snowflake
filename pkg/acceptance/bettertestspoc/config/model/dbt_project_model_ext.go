package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

// BasicDbtProjectModel creates a basic DBT project model with required fields
func BasicDbtProjectModel(
	resourceName string,
	database string,
	schema string,
	name string,
) *DbtProjectModel {
	return DbtProject(resourceName, database, schema, name).
		WithDefaultVersion("LAST").
		WithComment("Test DBT project")
}

// BasicDbtProjectModelWithDefaultMeta creates a basic DBT project model with default meta
func BasicDbtProjectModelWithDefaultMeta(
	database string,
	schema string,
	name string,
) *DbtProjectModel {
	return DbtProjectWithDefaultMeta(database, schema, name).
		WithDefaultVersion("LAST").
		WithComment("Test DBT project")
}

// WithGitSource adds a git_source block to the DBT project model
func (d *DbtProjectModel) WithGitSource(
	repositoryUrl string,
	stage string,
	branch string,
	path string,
	stagePath string,
) *DbtProjectModel {
	gitSourceConfig := map[string]tfconfig.Variable{
		"repository_url": tfconfig.StringVariable(repositoryUrl),
		"stage":          tfconfig.StringVariable(stage),
	}
	
	if branch != "" {
		gitSourceConfig["branch"] = tfconfig.StringVariable(branch)
	}
	
	if path != "" {
		gitSourceConfig["path"] = tfconfig.StringVariable(path)
	}
	
	if stagePath != "" {
		gitSourceConfig["stage_path"] = tfconfig.StringVariable(stagePath)
	}
	
	d.GitSource = tfconfig.ListVariable(
		tfconfig.ObjectVariable(gitSourceConfig),
	)
	return d
}

// WithGitSourceTag adds a git_source block with tag to the DBT project model
func (d *DbtProjectModel) WithGitSourceTag(
	repositoryUrl string,
	stage string,
	tag string,
	path string,
	stagePath string,
) *DbtProjectModel {
	gitSourceConfig := map[string]tfconfig.Variable{
		"repository_url": tfconfig.StringVariable(repositoryUrl),
		"stage":          tfconfig.StringVariable(stage),
		"tag":            tfconfig.StringVariable(tag),
	}
	
	if path != "" {
		gitSourceConfig["path"] = tfconfig.StringVariable(path)
	}
	
	if stagePath != "" {
		gitSourceConfig["stage_path"] = tfconfig.StringVariable(stagePath)
	}
	
	d.GitSource = tfconfig.ListVariable(
		tfconfig.ObjectVariable(gitSourceConfig),
	)
	return d
}

// WithFromStage adds a from block to the DBT project model
func (d *DbtProjectModel) WithFromStage(stage string) *DbtProjectModel {
	fromConfig := map[string]tfconfig.Variable{
		"stage": tfconfig.StringVariable(stage),
	}
	
	d.From = tfconfig.ListVariable(
		tfconfig.ObjectVariable(fromConfig),
	)
	return d
}
