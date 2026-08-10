package sdk

func init() {
	computePoolId := randomAccountObjectIdentifier()
	warehouseId := NewAccountObjectIdentifier("my_warehouse")
	integration1Id := NewAccountObjectIdentifier("integration1")
	integration2Id := NewAccountObjectIdentifier("integration2")
	tagId := NewAccountObjectIdentifier("tag1")
	tagId2 := NewAccountObjectIdentifier("tag2")
	comment := "comment"
	stageId := NewSchemaObjectIdentifier("db", "schema", "stage")
	stageLocation := NewStageLocation(stageId, "/path/to/spec")
	snapshotId := randomSchemaObjectIdentifier()

	servicesTests.Create.
		withDefaultOpts(func() *CreateServiceOptions {
			return &CreateServiceOptions{
				name:          servicesTestIdSchemaObjectIdentifier,
				InComputePool: computePoolId,
			}
		}).
		withAdditionalValidationCase(
			"validation_Create_MinReadyInstances_greaterThan0",
			func(opts *CreateServiceOptions) { opts.MinReadyInstances = new(0) },
			errIntValue("CreateServiceOptions", "MinReadyInstances", IntErrGreater, 0),
		).
		withAdditionalValidationCase(
			"validation_Create_MinInstances_greaterThan0",
			func(opts *CreateServiceOptions) { opts.MinInstances = new(0) },
			errIntValue("CreateServiceOptions", "MinInstances", IntErrGreater, 0),
		).
		withAdditionalValidationCase(
			"validation_Create_MaxInstances_greaterThan0",
			func(opts *CreateServiceOptions) { opts.MaxInstances = new(0) },
			errIntValue("CreateServiceOptions", "MaxInstances", IntErrGreater, 0),
		).
		withAdditionalValidationCase(
			"validation_Create_MinInstances_greaterOrEqualMinReadyInstances",
			func(opts *CreateServiceOptions) { opts.MinReadyInstances = new(3); opts.MinInstances = new(2) },
			errIntValue("CreateServiceOptions", "MinInstances", IntErrGreaterOrEqual, 3),
		).
		withAdditionalValidationCase(
			"validation_Create_MaxInstances_greaterOrEqualMinReadyInstances",
			func(opts *CreateServiceOptions) { opts.MinReadyInstances = new(3); opts.MaxInstances = new(2) },
			errIntValue("CreateServiceOptions", "MaxInstances", IntErrGreaterOrEqual, 3),
		).
		withAdditionalValidationCase(
			"validation_Create_MaxInstances_greaterOrEqualMinInstances",
			func(opts *CreateServiceOptions) { opts.MinInstances = new(3); opts.MaxInstances = new(2) },
			errIntValue("CreateServiceOptions", "MaxInstances", IntErrGreaterOrEqual, 3),
		).
		withAdditionalValidationCase(
			"validation_Create_AutoSuspendSecs_greaterOrEqual0",
			func(opts *CreateServiceOptions) { opts.AutoSuspendSecs = new(-1) },
			errIntValue("CreateServiceOptions", "AutoSuspendSecs", IntErrGreaterOrEqual, 0),
		).
		withModify(
			case_Services_validation_Create_opts_FromSpecification_ConflictingFields,
			func(opts *CreateServiceOptions) {
				opts.FromSpecification = &ServiceFromSpecification{Location: &stageLocation, Specification: new("{}")}
			},
		).
		withModify(
			case_Services_validation_Create_opts_FromSpecificationTemplate_ConflictingFields,
			func(opts *CreateServiceOptions) {
				opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{Location: &stageLocation, SpecificationTemplate: new("{}")}
			},
		).
		withModifyAndExpectedSqlf(
			case_Services_sql_Create_basic,
			func(opts *CreateServiceOptions) {
				opts.FromSpecification = &ServiceFromSpecification{SpecificationFile: new("spec.yaml")}
			},
			"CREATE SERVICE %s IN COMPUTE POOL %s FROM SPECIFICATION_FILE = 'spec.yaml'",
			servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(), computePoolId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Services_sql_Create_all,
			func(opts *CreateServiceOptions) {
				opts.IfNotExists = new(true)
				opts.FromSpecification = &ServiceFromSpecification{Specification: new("SPEC")}
				opts.AutoSuspendSecs = new(600)
				opts.ExternalAccessIntegrations = &ServiceExternalAccessIntegrations{ExternalAccessIntegrations: []AccountObjectIdentifier{integration1Id}}
				opts.AutoResume = new(true)
				opts.MinInstances = new(1)
				opts.MinReadyInstances = new(1)
				opts.MaxInstances = new(3)
				opts.QueryWarehouse = &warehouseId
				opts.Tag = []TagAssociation{{Name: tagId, Value: "value1"}}
				opts.Comment = &comment
			},
			"CREATE SERVICE IF NOT EXISTS %s IN COMPUTE POOL %s FROM SPECIFICATION $$SPEC$$ AUTO_SUSPEND_SECS = 600 "+
				"EXTERNAL_ACCESS_INTEGRATIONS = (%s) AUTO_RESUME = true MIN_INSTANCES = 1 MIN_READY_INSTANCES = 1 MAX_INSTANCES = 3 "+
				"QUERY_WAREHOUSE = %s TAG (%s = 'value1') COMMENT = '%s'",
			servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(), computePoolId.FullyQualifiedName(), integration1Id.FullyQualifiedName(),
			warehouseId.FullyQualifiedName(), tagId.FullyQualifiedName(), comment,
		).
		withAdditionalSqlCasef(
			"sql_Create_withIfNotExists",
			func(opts *CreateServiceOptions) {
				opts.IfNotExists = new(true)
				opts.FromSpecification = &ServiceFromSpecification{SpecificationFile: new("spec.yaml")}
			},
			"CREATE SERVICE IF NOT EXISTS %s IN COMPUTE POOL %s FROM SPECIFICATION_FILE = 'spec.yaml'",
			servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(), computePoolId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_withSpecificationFileOnStage",
			func(opts *CreateServiceOptions) {
				opts.FromSpecification = &ServiceFromSpecification{Location: &stageLocation, SpecificationFile: new("spec.yaml")}
			},
			`CREATE SERVICE %s IN COMPUTE POOL %s FROM '@\"%s\".\"%s\".\"%s\"//path/to/spec' SPECIFICATION_FILE = 'spec.yaml'`,
			servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(), computePoolId.FullyQualifiedName(), stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(),
		).
		withAdditionalSqlCasef(
			"sql_Create_fromSpecification",
			func(opts *CreateServiceOptions) {
				opts.FromSpecification = &ServiceFromSpecification{Specification: new("SPEC")}
			},
			"CREATE SERVICE %s IN COMPUTE POOL %s FROM SPECIFICATION $$SPEC$$",
			servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(), computePoolId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_fromSpecificationTemplateFile",
			func(opts *CreateServiceOptions) {
				opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
					SpecificationTemplateFile: new("spec.yaml"),
					Using: []ListItem{
						{Key: "string", Value: `"bar"`},
						{Key: "int", Value: 42},
						{Key: "bool", Value: true},
					},
				}
			},
			`CREATE SERVICE %s IN COMPUTE POOL %s FROM SPECIFICATION_TEMPLATE_FILE = 'spec.yaml' USING ("string" => "bar", "int" => 42, "bool" => true)`,
			servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(), computePoolId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_fromSpecificationTemplateFileOnStage",
			func(opts *CreateServiceOptions) {
				opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
					Location:                  &stageLocation,
					SpecificationTemplateFile: new("spec.yaml"),
					Using:                     []ListItem{{Key: "string", Value: `"bar"`}},
				}
			},
			`CREATE SERVICE %s IN COMPUTE POOL %s FROM '@\"%s\".\"%s\".\"%s\"//path/to/spec' SPECIFICATION_TEMPLATE_FILE = 'spec.yaml' USING ("string" => "bar")`,
			servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(), computePoolId.FullyQualifiedName(), stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(),
		).
		withAdditionalSqlCasef(
			"sql_Create_fromSpecificationTemplate",
			func(opts *CreateServiceOptions) {
				opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
					SpecificationTemplate: new("SPEC"),
					Using:                 []ListItem{{Key: "string", Value: `"bar"`}},
				}
			},
			`CREATE SERVICE %s IN COMPUTE POOL %s FROM SPECIFICATION_TEMPLATE $$SPEC$$ USING ("string" => "bar")`,
			servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(), computePoolId.FullyQualifiedName(),
		)

	servicesTests.Alter.
		withAdditionalValidationCase(
			"validation_Alter_Set_MinReadyInstances_greaterThan0",
			func(opts *AlterServiceOptions) { opts.Set = &ServiceSet{MinReadyInstances: new(0)} },
			errIntValue("AlterServiceOptions.Set", "MinReadyInstances", IntErrGreater, 0),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_MinInstances_greaterThan0",
			func(opts *AlterServiceOptions) { opts.Set = &ServiceSet{MinInstances: new(0)} },
			errIntValue("AlterServiceOptions.Set", "MinInstances", IntErrGreater, 0),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_MaxInstances_greaterThan0",
			func(opts *AlterServiceOptions) { opts.Set = &ServiceSet{MaxInstances: new(0)} },
			errIntValue("AlterServiceOptions.Set", "MaxInstances", IntErrGreater, 0),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_MinInstances_greaterOrEqualMinReadyInstances",
			func(opts *AlterServiceOptions) {
				opts.Set = &ServiceSet{MinReadyInstances: new(3), MinInstances: new(2)}
			},
			errIntValue("AlterServiceOptions.Set", "MinInstances", IntErrGreaterOrEqual, 3),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_MaxInstances_greaterOrEqualMinReadyInstances",
			func(opts *AlterServiceOptions) {
				opts.Set = &ServiceSet{MinReadyInstances: new(3), MaxInstances: new(2)}
			},
			errIntValue("AlterServiceOptions.Set", "MaxInstances", IntErrGreaterOrEqual, 3),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_MaxInstances_greaterOrEqualMinInstances",
			func(opts *AlterServiceOptions) {
				opts.Set = &ServiceSet{MinInstances: new(3), MaxInstances: new(2)}
			},
			errIntValue("AlterServiceOptions.Set", "MaxInstances", IntErrGreaterOrEqual, 3),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_AutoSuspendSecs_greaterOrEqual0",
			func(opts *AlterServiceOptions) { opts.Set = &ServiceSet{AutoSuspendSecs: new(-1)} },
			errIntValue("AlterServiceOptions.Set", "AutoSuspendSecs", IntErrGreaterOrEqual, 0),
		).
		withModify(
			case_Services_validation_Alter_opts_FromSpecification_ConflictingFields,
			func(opts *AlterServiceOptions) {
				opts.FromSpecification = &ServiceFromSpecification{Location: &stageLocation, Specification: new("{}")}
			},
		).
		withModify(
			case_Services_validation_Alter_opts_FromSpecificationTemplate_ConflictingFields,
			func(opts *AlterServiceOptions) {
				opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{Location: &stageLocation, SpecificationTemplate: new("{}")}
			},
		).
		withModifyAndExpectedSqlf(
			case_Services_sql_Alter_Resume,
			func(opts *AlterServiceOptions) { opts.Resume = new(true) },
			"ALTER SERVICE %s RESUME", servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Services_sql_Alter_Suspend,
			func(opts *AlterServiceOptions) { opts.Suspend = new(true) },
			"ALTER SERVICE %s SUSPEND", servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Services_sql_Alter_FromSpecification,
			func(opts *AlterServiceOptions) {
				opts.FromSpecification = &ServiceFromSpecification{SpecificationFile: new("spec.yaml")}
			},
			"ALTER SERVICE %s FROM SPECIFICATION_FILE = 'spec.yaml'", servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Services_sql_Alter_FromSpecificationTemplate,
			func(opts *AlterServiceOptions) {
				opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
					SpecificationTemplateFile: new("spec.yaml"),
					Using: []ListItem{
						{Key: "string", Value: `"bar"`},
						{Key: "int", Value: 42},
						{Key: "bool", Value: true},
					},
				}
			},
			`ALTER SERVICE %s FROM SPECIFICATION_TEMPLATE_FILE = 'spec.yaml' USING ("string" => "bar", "int" => 42, "bool" => true)`,
			servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Services_sql_Alter_Restore,
			func(opts *AlterServiceOptions) {
				opts.Restore = &Restore{
					Volume:       "vol1",
					Instances:    []int{0, 1},
					FromSnapshot: snapshotId,
				}
			},
			`ALTER SERVICE %s RESTORE VOLUME "vol1" INSTANCES 0, 1 FROM SNAPSHOT %s`,
			servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(), snapshotId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Services_sql_Alter_Set,
			func(opts *AlterServiceOptions) {
				opts.Set = &ServiceSet{
					MinInstances:      new(2),
					MaxInstances:      new(5),
					AutoSuspendSecs:   new(600),
					MinReadyInstances: new(1),
					QueryWarehouse:    &warehouseId,
					AutoResume:        new(true),
					ExternalAccessIntegrations: &ServiceExternalAccessIntegrations{
						ExternalAccessIntegrations: []AccountObjectIdentifier{integration1Id, integration2Id},
					},
					Comment: &comment,
				}
			},
			`ALTER SERVICE %s SET MIN_INSTANCES = 2 MAX_INSTANCES = 5 AUTO_SUSPEND_SECS = 600 MIN_READY_INSTANCES = 1 QUERY_WAREHOUSE = %s AUTO_RESUME = true`+
				` EXTERNAL_ACCESS_INTEGRATIONS = (%s, %s) COMMENT = '%s'`,
			servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(), warehouseId.FullyQualifiedName(), integration1Id.FullyQualifiedName(), integration2Id.FullyQualifiedName(), comment,
		).
		withModifyAndExpectedSqlf(
			case_Services_sql_Alter_Unset,
			func(opts *AlterServiceOptions) {
				opts.Unset = &ServiceUnset{
					MinInstances:               new(true),
					AutoSuspendSecs:            new(true),
					MaxInstances:               new(true),
					MinReadyInstances:          new(true),
					QueryWarehouse:             new(true),
					AutoResume:                 new(true),
					ExternalAccessIntegrations: new(true),
					Comment:                    new(true),
				}
			},
			"ALTER SERVICE %s UNSET MIN_INSTANCES, AUTO_SUSPEND_SECS, MAX_INSTANCES, MIN_READY_INSTANCES, QUERY_WAREHOUSE, AUTO_RESUME, EXTERNAL_ACCESS_INTEGRATIONS, COMMENT",
			servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Services_sql_Alter_SetTags,
			func(opts *AlterServiceOptions) { opts.SetTags = []TagAssociation{{Name: tagId, Value: "value1"}} },
			"ALTER SERVICE %s SET TAG %s = 'value1'", servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Services_sql_Alter_UnsetTags,
			func(opts *AlterServiceOptions) { opts.UnsetTags = []ObjectIdentifier{tagId, tagId2} },
			"ALTER SERVICE %s UNSET TAG %s, %s", servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_withIfExists",
			func(opts *AlterServiceOptions) { opts.IfExists = new(true); opts.Suspend = new(true) },
			"ALTER SERVICE IF EXISTS %s SUSPEND", servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_fromSpecificationFileOnStage",
			func(opts *AlterServiceOptions) {
				opts.FromSpecification = &ServiceFromSpecification{Location: &stageLocation, SpecificationFile: new("spec.yaml")}
			},
			`ALTER SERVICE %s FROM '@\"%s\".\"%s\".\"%s\"//path/to/spec' SPECIFICATION_FILE = 'spec.yaml'`,
			servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(), stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_fromSpecification",
			func(opts *AlterServiceOptions) {
				opts.FromSpecification = &ServiceFromSpecification{Specification: new("SPEC")}
			},
			"ALTER SERVICE %s FROM SPECIFICATION $$SPEC$$", servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_withSpecificationTemplateFileOnStage",
			func(opts *AlterServiceOptions) {
				opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
					Location:                  &stageLocation,
					SpecificationTemplateFile: new("spec.yaml"),
					Using:                     []ListItem{{Key: "string", Value: `"bar"`}},
				}
			},
			`ALTER SERVICE %s FROM '@\"%s\".\"%s\".\"%s\"//path/to/spec' SPECIFICATION_TEMPLATE_FILE = 'spec.yaml' USING ("string" => "bar")`,
			servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(), stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_fromSpecificationTemplate",
			func(opts *AlterServiceOptions) {
				opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
					SpecificationTemplate: new("SPEC"),
					Using:                 []ListItem{{Key: "string", Value: `"bar"`}},
				}
			},
			`ALTER SERVICE %s FROM SPECIFICATION_TEMPLATE $$SPEC$$ USING ("string" => "bar")`, servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	servicesTests.Drop.
		withExpectedSqlf(case_Services_sql_Drop_basic, "DROP SERVICE %s", servicesTestIdSchemaObjectIdentifier.FullyQualifiedName()).
		withModifyAndExpectedSqlf(
			case_Services_sql_Drop_all,
			func(opts *DropServiceOptions) { opts.IfExists = new(true); opts.Force = new(true) },
			"DROP SERVICE IF EXISTS %s FORCE", servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	servicesTests.Show.
		withExpectedSql(case_Services_sql_Show_basic, "SHOW SERVICES").
		withModifyAndExpectedSqlf(
			case_Services_sql_Show_all,
			func(opts *ShowServiceOptions) {
				opts.Like = &Like{Pattern: new("service_*")}
				opts.In = &ServiceIn{In: In{Database: NewAccountObjectIdentifier("database")}}
				opts.StartsWith = new("my_prefix")
				opts.Limit = &LimitFrom{Rows: new(10), From: new("service1")}
			},
			`SHOW SERVICES LIKE 'service_*' IN DATABASE "database" STARTS WITH 'my_prefix' LIMIT 10 FROM 'service1'`,
		).
		withModifyAndExpectedSqlf(
			case_Services_sql_Show_Like,
			func(opts *ShowServiceOptions) { opts.Like = &Like{Pattern: new("service_*")} },
			"SHOW SERVICES LIKE 'service_*'",
		).
		withModifyAndExpectedSqlf(
			case_Services_sql_Show_In,
			func(opts *ShowServiceOptions) {
				opts.In = &ServiceIn{In: In{Database: NewAccountObjectIdentifier("database")}}
			},
			`SHOW SERVICES IN DATABASE "database"`,
		).
		withModifyAndExpectedSqlf(
			case_Services_sql_Show_StartsWith,
			func(opts *ShowServiceOptions) { opts.StartsWith = new("my_prefix") },
			"SHOW SERVICES STARTS WITH 'my_prefix'",
		).
		withModifyAndExpectedSqlf(
			case_Services_sql_Show_Limit,
			func(opts *ShowServiceOptions) { opts.Limit = &LimitFrom{Rows: new(10)} },
			"SHOW SERVICES LIMIT 10",
		).
		withAdditionalSqlCasef(
			"sql_Show_withJobsOption",
			func(opts *ShowServiceOptions) { opts.Job = new(true) },
			"SHOW JOB SERVICES",
		).
		withAdditionalSqlCasef(
			"sql_Show_withExcludeJobs",
			func(opts *ShowServiceOptions) { opts.ExcludeJobs = new(true) },
			"SHOW SERVICES EXCLUDE JOBS",
		).
		withAdditionalSqlCasef(
			"sql_Show_inComputePool",
			func(opts *ShowServiceOptions) {
				opts.In = &ServiceIn{ComputePool: NewAccountObjectIdentifier("compute_pool")}
			},
			`SHOW SERVICES IN COMPUTE POOL "compute_pool"`,
		).
		withAdditionalSqlCasef(
			"sql_Show_withLimitAndFrom",
			func(opts *ShowServiceOptions) { opts.Limit = &LimitFrom{Rows: new(10), From: new("service1")} },
			"SHOW SERVICES LIMIT 10 FROM 'service1'",
		)

	servicesTests.Describe.
		withExpectedSqlf(case_Services_sql_Describe_basic, "DESCRIBE SERVICE %s", servicesTestIdSchemaObjectIdentifier.FullyQualifiedName())

	servicesTests.ExecuteJob.
		withDefaultOpts(func() *ExecuteJobServiceOptions {
			return &ExecuteJobServiceOptions{
				Name:          servicesTestIdSchemaObjectIdentifier,
				InComputePool: computePoolId,
			}
		}).
		withModify(
			case_Services_validation_ExecuteJob_opts_JobServiceFromSpecification_ExactlyOneValueSet_Location_Specification_MoreThanOneSet,
			func(opts *ExecuteJobServiceOptions) {
				opts.JobServiceFromSpecification = &JobServiceFromSpecification{Location: &stageLocation, Specification: new("{}")}
			},
		).
		withModify(
			case_Services_validation_ExecuteJob_opts_JobServiceFromSpecificationTemplate_ExactlyOneValueSet_Location_SpecificationTemplate_MoreThanOneSet,
			func(opts *ExecuteJobServiceOptions) {
				opts.JobServiceFromSpecificationTemplate = &JobServiceFromSpecificationTemplate{Location: &stageLocation, SpecificationTemplate: new("{}")}
			},
		).
		withModifyAndExpectedSqlf(
			case_Services_sql_ExecuteJob_basic,
			func(opts *ExecuteJobServiceOptions) {
				opts.JobServiceFromSpecification = &JobServiceFromSpecification{Specification: new("SPEC")}
			},
			"EXECUTE JOB SERVICE IN COMPUTE POOL %s NAME = %s FROM SPECIFICATION $$SPEC$$",
			computePoolId.FullyQualifiedName(), servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_ExecuteJob_withSpecificationFileOnStage",
			func(opts *ExecuteJobServiceOptions) {
				opts.JobServiceFromSpecification = &JobServiceFromSpecification{Location: &stageLocation, SpecificationFile: new("spec.yaml")}
			},
			`EXECUTE JOB SERVICE IN COMPUTE POOL %s NAME = %s FROM '@\"%s\".\"%s\".\"%s\"//path/to/spec' SPECIFICATION_FILE = 'spec.yaml'`,
			computePoolId.FullyQualifiedName(), servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(), stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(),
		).
		withAdditionalSqlCasef(
			"sql_ExecuteJob_fromSpecificationTemplateFileOnStage",
			func(opts *ExecuteJobServiceOptions) {
				opts.JobServiceFromSpecificationTemplate = &JobServiceFromSpecificationTemplate{
					Location:                  &stageLocation,
					SpecificationTemplateFile: new("spec.yaml"),
					Using:                     []ListItem{{Key: "string", Value: `"bar"`}},
				}
			},
			`EXECUTE JOB SERVICE IN COMPUTE POOL %s NAME = %s FROM '@\"%s\".\"%s\".\"%s\"//path/to/spec' SPECIFICATION_TEMPLATE_FILE = 'spec.yaml' USING ("string" => "bar")`,
			computePoolId.FullyQualifiedName(), servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(), stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(),
		).
		withAdditionalSqlCasef(
			"sql_ExecuteJob_fromSpecificationTemplate",
			func(opts *ExecuteJobServiceOptions) {
				opts.JobServiceFromSpecificationTemplate = &JobServiceFromSpecificationTemplate{
					SpecificationTemplate: new("SPEC"),
					Using:                 []ListItem{{Key: "string", Value: `"bar"`}},
				}
			},
			`EXECUTE JOB SERVICE IN COMPUTE POOL %s NAME = %s FROM SPECIFICATION_TEMPLATE $$SPEC$$ USING ("string" => "bar")`,
			computePoolId.FullyQualifiedName(), servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_ExecuteJob_allOptions",
			func(opts *ExecuteJobServiceOptions) {
				opts.Async = new(true)
				opts.JobServiceFromSpecification = &JobServiceFromSpecification{Specification: new("SPEC")}
				opts.ExternalAccessIntegrations = &ServiceExternalAccessIntegrations{ExternalAccessIntegrations: []AccountObjectIdentifier{integration1Id}}
				opts.QueryWarehouse = &warehouseId
				opts.Tag = []TagAssociation{{Name: tagId, Value: "value1"}}
				opts.Comment = &comment
			},
			"EXECUTE JOB SERVICE IN COMPUTE POOL %s NAME = %s ASYNC = true QUERY_WAREHOUSE = %s COMMENT = '%s' "+
				"EXTERNAL_ACCESS_INTEGRATIONS = (%s) FROM SPECIFICATION $$SPEC$$ TAG (%s = 'value1')",
			computePoolId.FullyQualifiedName(), servicesTestIdSchemaObjectIdentifier.FullyQualifiedName(), warehouseId.FullyQualifiedName(), comment, integration1Id.FullyQualifiedName(), tagId.FullyQualifiedName(),
		)
}
