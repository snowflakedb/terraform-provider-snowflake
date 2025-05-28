package sdk

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
)

func TestServices_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	computePoolId := randomAccountObjectIdentifier()
	// Minimal valid CreateServiceOptions
	defaultOpts := func() *CreateServiceOptions {
		return &CreateServiceOptions{
			name:          id,
			InComputePool: computePoolId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateServiceOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: min ready instances must be greater than 0", func(t *testing.T) {
		opts := defaultOpts()
		opts.MinReadyInstances = Pointer(0)
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("CreateServiceOptions", "MinReadyInstances", IntErrGreater, 0))
	})

	t.Run("validation: min instances must be greater than 0", func(t *testing.T) {
		opts := defaultOpts()
		opts.MinInstances = Pointer(0)
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("CreateServiceOptions", "MinInstances", IntErrGreater, 0))
	})

	t.Run("validation: max instances must be greater than 0", func(t *testing.T) {
		opts := defaultOpts()
		opts.MaxInstances = Pointer(0)
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("CreateServiceOptions", "MaxInstances", IntErrGreater, 0))
	})

	t.Run("validation: min instances must be greater than or equal to min ready instances", func(t *testing.T) {
		opts := defaultOpts()
		opts.MinReadyInstances = Pointer(3)
		opts.MinInstances = Pointer(2)
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("CreateServiceOptions", "MinInstances", IntErrGreaterOrEqual, 3))
	})

	t.Run("validation: max instances must be greater than or equal to min ready instances", func(t *testing.T) {
		opts := defaultOpts()
		opts.MinReadyInstances = Pointer(3)
		opts.MaxInstances = Pointer(2)
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("CreateServiceOptions", "MaxInstances", IntErrGreaterOrEqual, 3))
	})

	t.Run("validation: max instances must be greater than or equal to min instances", func(t *testing.T) {
		opts := defaultOpts()
		opts.MinInstances = Pointer(3)
		opts.MaxInstances = Pointer(2)
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("CreateServiceOptions", "MaxInstances", IntErrGreaterOrEqual, 3))
	})

	t.Run("validation: invalid auto suspend secs", func(t *testing.T) {
		opts := defaultOpts()
		opts.AutoSuspendSecs = Pointer(-1)
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("CreateServiceOptions", "AutoSuspendSecs", IntErrGreaterOrEqual, 0))
	})

	t.Run("validation: conflicting fields for [opts.FromSpecification opts.FromSpecificationTemplate]", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecification = &ServiceFromSpecification{
			Specification: String("{}"),
		}
		opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
			SpecificationTemplate: String("{}"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateServiceOptions", "FromSpecification", "FromSpecificationTemplate"))
	})

	t.Run("validation: empty opts.FromSpecification", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecification = &ServiceFromSpecification{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateServiceOptions.FromSpecification", "SpecificationFile", "Specification"))
	})

	t.Run("validation: conflicting fields for [opts.FromSpecification.Stage opts.FromSpecification.Specification]", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecification = &ServiceFromSpecification{
			Specification:     String("{}"),
			SpecificationFile: String("spec.yaml"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateServiceOptions.FromSpecification", "SpecificationFile", "Specification"))
	})

	t.Run("validation: stage present without specification", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecification = &ServiceFromSpecification{
			Stage: Pointer("@stage"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateServiceOptions.FromSpecification", "SpecificationFile", "Specification"))
	})

	t.Run("validation: all specification fields present", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecification = &ServiceFromSpecification{
			Stage:             Pointer("@stage"),
			Specification:     String("{}"),
			SpecificationFile: String("spec.yaml"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateServiceOptions.FromSpecification", "SpecificationFile", "Specification"))
	})

	t.Run("validation: empty opts.FromSpecificationTemplate", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateServiceOptions.FromSpecificationTemplate", "SpecificationTemplateFile", "SpecificationTemplate"))
	})

	t.Run("validation: conflicting fields for [opts.FromSpecificationTemplate.Stage opts.FromSpecificationTemplate.SpecificationTemplate]", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
			SpecificationTemplate:     String("{}"),
			SpecificationTemplateFile: String("spec.yaml"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateServiceOptions.FromSpecificationTemplate", "SpecificationTemplateFile", "SpecificationTemplate"))
	})

	t.Run("validation: stage present without specification template", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
			Stage: Pointer("@stage"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateServiceOptions.FromSpecificationTemplate", "SpecificationTemplateFile", "SpecificationTemplate"))
	})

	t.Run("validation: all specification template fields present", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
			Stage:                     Pointer("@stage"),
			SpecificationTemplate:     String("{}"),
			SpecificationTemplateFile: String("spec.yaml"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateServiceOptions.FromSpecificationTemplate", "SpecificationTemplateFile", "SpecificationTemplate"))
	})

	t.Run("with if not exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.FromSpecification = &ServiceFromSpecification{
			SpecificationFile: String("spec.yaml"),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE SERVICE IF NOT EXISTS %s IN COMPUTE POOL %s FROM SPECIFICATION_FILE = 'spec.yaml'", id.FullyQualifiedName(), computePoolId.FullyQualifiedName())
	})

	t.Run("from specification file", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecification = &ServiceFromSpecification{
			SpecificationFile: String("spec.yaml"),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE SERVICE %s IN COMPUTE POOL %s FROM SPECIFICATION_FILE = 'spec.yaml'", id.FullyQualifiedName(), computePoolId.FullyQualifiedName())
	})

	t.Run("with specification file on stage", func(t *testing.T) {
		stageId := NewSchemaObjectIdentifier("db", "schema", "stage")
		opts := defaultOpts()
		opts.FromSpecification = &ServiceFromSpecification{
			Stage:             Pointer(fmt.Sprintf("@%s", stageId.FullyQualifiedName())),
			SpecificationFile: String("spec.yaml"),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE SERVICE %s IN COMPUTE POOL %s FROM @%s SPECIFICATION_FILE = 'spec.yaml'",
			id.FullyQualifiedName(), computePoolId.FullyQualifiedName(), stageId.FullyQualifiedName())
	})

	t.Run("from specification", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecification = &ServiceFromSpecification{
			Specification: String("{}"),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE SERVICE %s IN COMPUTE POOL %s FROM SPECIFICATION '{}'", id.FullyQualifiedName(), computePoolId.FullyQualifiedName())
	})

	t.Run("from specification template file", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
			SpecificationTemplateFile: String("spec.yaml"),
			Using: []ListItem{
				{
					Key:   "string",
					Value: `"bar"`,
				},
				{
					Key:   "int",
					Value: 42,
				},
				{
					Key:   "bool",
					Value: true,
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE SERVICE %s IN COMPUTE POOL %s FROM SPECIFICATION_TEMPLATE_FILE = 'spec.yaml' USING ("string" => "bar", "int" => 42, "bool" => true)`, id.FullyQualifiedName(), computePoolId.FullyQualifiedName())
	})

	t.Run("from specification template file on stage", func(t *testing.T) {
		stageId := NewSchemaObjectIdentifier("db", "schema", "stage")
		opts := defaultOpts()
		opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
			Stage:                     Pointer(fmt.Sprintf("@%s", stageId.FullyQualifiedName())),
			SpecificationTemplateFile: String("spec.yaml"),
			Using: []ListItem{
				{
					Key:   "string",
					Value: `"bar"`,
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE SERVICE %s IN COMPUTE POOL %s FROM @%s SPECIFICATION_TEMPLATE_FILE = 'spec.yaml' USING ("string" => "bar")`,
			id.FullyQualifiedName(), computePoolId.FullyQualifiedName(), stageId.FullyQualifiedName())
	})

	t.Run("from specification template", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
			SpecificationTemplate: String("{}"),
			Using: []ListItem{
				{
					Key:   "string",
					Value: `"bar"`,
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE SERVICE %s IN COMPUTE POOL %s FROM SPECIFICATION_TEMPLATE '{}' USING ("string" => "bar")`, id.FullyQualifiedName(), computePoolId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		warehouseId := NewAccountObjectIdentifier("my_warehouse")
		integration1Id := NewAccountObjectIdentifier("integration1")
		comment := random.Comment()

		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.FromSpecification = &ServiceFromSpecification{
			Specification: String("{}"),
		}
		opts.AutoSuspendSecs = Pointer(600)
		opts.ExternalAccessIntegrations = &ServiceExternalAccessIntegrations{
			ExternalAccessIntegrations: []AccountObjectIdentifier{
				integration1Id,
			},
		}
		opts.AutoResume = Bool(true)
		opts.MinInstances = Pointer(1)
		opts.MinReadyInstances = Pointer(1)
		opts.MaxInstances = Pointer(3)
		opts.QueryWarehouse = &warehouseId
		opts.Tag = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag1"),
				Value: "value1",
			},
		}
		opts.Comment = &comment

		assertOptsValidAndSQLEquals(t, opts, "CREATE SERVICE IF NOT EXISTS %s IN COMPUTE POOL %s FROM SPECIFICATION '{}' AUTO_SUSPEND_SECS = 600 "+
			"EXTERNAL_ACCESS_INTEGRATIONS = (%s) AUTO_RESUME = true MIN_INSTANCES = 1 MIN_READY_INSTANCES = 1 MAX_INSTANCES = 3 "+
			"QUERY_WAREHOUSE = %s TAG (\"tag1\" = 'value1') COMMENT = '%s'",
			id.FullyQualifiedName(), computePoolId.FullyQualifiedName(), integration1Id.FullyQualifiedName(),
			warehouseId.FullyQualifiedName(), comment)
	})
}

func TestServices_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	// Minimal valid AlterServiceOptions
	defaultOpts := func() *AlterServiceOptions {
		return &AlterServiceOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterServiceOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one property should be set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterServiceOptions", "Resume", "Suspend", "FromSpecification", "FromSpecificationTemplate", "Restore", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: min ready instances must be greater than 0", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ServiceSet{
			MinReadyInstances: Pointer(0),
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("AlterServiceOptions.Set", "MinReadyInstances", IntErrGreater, 0))
	})

	t.Run("validation: min instances must be greater than 0", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ServiceSet{
			MinInstances: Pointer(0),
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("AlterServiceOptions.Set", "MinInstances", IntErrGreater, 0))
	})

	t.Run("validation: max instances must be greater than 0", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ServiceSet{
			MaxInstances: Pointer(0),
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("AlterServiceOptions.Set", "MaxInstances", IntErrGreater, 0))
	})

	t.Run("validation: min instances must be greater than or equal to min ready instances", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ServiceSet{
			MinReadyInstances: Pointer(3),
			MinInstances:      Pointer(2),
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("AlterServiceOptions.Set", "MinInstances", IntErrGreaterOrEqual, 3))
	})

	t.Run("validation: max instances must be greater than or equal to min ready instances", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ServiceSet{
			MinReadyInstances: Pointer(3),
			MaxInstances:      Pointer(2),
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("AlterServiceOptions.Set", "MaxInstances", IntErrGreaterOrEqual, 3))
	})

	t.Run("validation: max instances must be greater than or equal to min instances", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ServiceSet{
			MinInstances: Pointer(3),
			MaxInstances: Pointer(2),
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("AlterServiceOptions.Set", "MaxInstances", IntErrGreaterOrEqual, 3))
	})

	t.Run("validation: invalid auto suspend secs", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ServiceSet{
			AutoSuspendSecs: Pointer(-1),
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("AlterServiceOptions.Set", "AutoSuspendSecs", IntErrGreaterOrEqual, 0))
	})

	t.Run("validation: at least one property should be set in Set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ServiceSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterServiceOptions.Set", "MinInstances", "MaxInstances", "AutoSuspendSecs", "MinReadyInstances", "QueryWarehouse", "AutoResume", "ExternalAccessIntegrations", "Comment"))
	})

	t.Run("validation: at least one property should be set in Unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ServiceUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterServiceOptions.Unset", "MinInstances", "AutoSuspendSecs", "MaxInstances", "MinReadyInstances", "QueryWarehouse", "AutoResume", "ExternalAccessIntegrations", "Comment"))
	})

	t.Run("validation: conflicting fields for [opts.FromSpecification opts.FromSpecificationTemplate]", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecification = &ServiceFromSpecification{
			Specification: String("{}"),
		}
		opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
			SpecificationTemplate: String("{}"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterServiceOptions", "Resume", "Suspend", "FromSpecification", "FromSpecificationTemplate", "Restore", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: empty opts.FromSpecification", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecification = &ServiceFromSpecification{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterServiceOptions.FromSpecification", "SpecificationFile", "Specification"))
	})

	t.Run("validation: conflicting fields for [opts.FromSpecification.Stage opts.FromSpecification.Specification]", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecification = &ServiceFromSpecification{
			Specification:     String("{}"),
			SpecificationFile: String("spec.yaml"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterServiceOptions.FromSpecification", "SpecificationFile", "Specification"))
	})

	t.Run("validation: stage present without specification", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecification = &ServiceFromSpecification{
			Stage: Pointer("@stage"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterServiceOptions.FromSpecification", "SpecificationFile", "Specification"))
	})

	t.Run("validation: all specification fields present", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecification = &ServiceFromSpecification{
			Stage:             Pointer("@stage"),
			Specification:     String("{}"),
			SpecificationFile: String("spec.yaml"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterServiceOptions.FromSpecification", "SpecificationFile", "Specification"))
	})

	t.Run("validation: empty opts.FromSpecificationTemplate", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterServiceOptions.FromSpecificationTemplate", "SpecificationTemplateFile", "SpecificationTemplate"))
	})

	t.Run("validation: conflicting fields for [opts.FromSpecificationTemplate.Stage opts.FromSpecificationTemplate.SpecificationTemplate]", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
			SpecificationTemplate:     String("{}"),
			SpecificationTemplateFile: String("spec.yaml"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterServiceOptions.FromSpecificationTemplate", "SpecificationTemplateFile", "SpecificationTemplate"))
	})

	t.Run("validation: stage present without specification template", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
			Stage: Pointer("@stage"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterServiceOptions.FromSpecificationTemplate", "SpecificationTemplateFile", "SpecificationTemplate"))
	})

	t.Run("validation: all specification template fields present", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
			Stage:                     Pointer("@stage"),
			SpecificationTemplate:     String("{}"),
			SpecificationTemplateFile: String("spec.yaml"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterServiceOptions.FromSpecificationTemplate", "SpecificationTemplateFile", "SpecificationTemplate"))
	})
	t.Run("suspend", func(t *testing.T) {
		opts := defaultOpts()
		opts.Suspend = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER SERVICE %s SUSPEND", id.FullyQualifiedName())
	})

	t.Run("resume", func(t *testing.T) {
		opts := defaultOpts()
		opts.Resume = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER SERVICE %s RESUME", id.FullyQualifiedName())
	})

	t.Run("with if exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Suspend = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER SERVICE IF EXISTS %s SUSPEND", id.FullyQualifiedName())
	})

	t.Run("from specification file", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecification = &ServiceFromSpecification{
			SpecificationFile: String("spec.yaml"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SERVICE %s FROM SPECIFICATION_FILE = 'spec.yaml'", id.FullyQualifiedName())
	})

	t.Run("from specification file on stage", func(t *testing.T) {
		stageId := NewSchemaObjectIdentifier("db", "schema", "stage")
		opts := defaultOpts()
		opts.FromSpecification = &ServiceFromSpecification{
			Stage:             Pointer(fmt.Sprintf("@%s", stageId.FullyQualifiedName())),
			SpecificationFile: String("spec.yaml"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SERVICE %s FROM @%s SPECIFICATION_FILE = 'spec.yaml'",
			id.FullyQualifiedName(), stageId.FullyQualifiedName())
	})

	t.Run("from specification", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecification = &ServiceFromSpecification{
			Specification: String("{}"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SERVICE %s FROM SPECIFICATION '{}'", id.FullyQualifiedName())
	})

	t.Run("from specification template file", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
			SpecificationTemplateFile: String("spec.yaml"),
			Using: []ListItem{
				{
					Key:   "string",
					Value: `"bar"`,
				},
				{
					Key:   "int",
					Value: 42,
				},
				{
					Key:   "bool",
					Value: true,
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SERVICE %s FROM SPECIFICATION_TEMPLATE_FILE = 'spec.yaml' USING ("string" => "bar", "int" => 42, "bool" => true)`, id.FullyQualifiedName())
	})

	t.Run("with specification template file on stage", func(t *testing.T) {
		stageId := NewSchemaObjectIdentifier("db", "schema", "stage")
		opts := defaultOpts()
		opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
			Stage:                     Pointer(fmt.Sprintf("@%s", stageId.FullyQualifiedName())),
			SpecificationTemplateFile: String("spec.yaml"),
			Using: []ListItem{
				{
					Key:   "string",
					Value: `"bar"`,
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SERVICE %s FROM @%s SPECIFICATION_TEMPLATE_FILE = 'spec.yaml' USING ("string" => "bar")`,
			id.FullyQualifiedName(), stageId.FullyQualifiedName())
	})

	t.Run("from specification template", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromSpecificationTemplate = &ServiceFromSpecificationTemplate{
			SpecificationTemplate: String("{}"),
			Using: []ListItem{
				{
					Key:   "string",
					Value: `"bar"`,
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SERVICE %s FROM SPECIFICATION_TEMPLATE '{}' USING ("string" => "bar")`, id.FullyQualifiedName())
	})

	t.Run("with restore", func(t *testing.T) {
		opts := defaultOpts()
		opts.Restore = &Restore{
			Volume:       "vol1",
			Instances:    []int{0, 1},
			FromSnapshot: "snap1",
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SERVICE %s RESTORE VOLUME "vol1" INSTANCES 0, 1 FROM SNAPSHOT "snap1"`, id.FullyQualifiedName())
	})

	t.Run("with set", func(t *testing.T) {
		warehouseId := NewAccountObjectIdentifier("my_warehouse")
		comment := random.Comment()
		integration1Id := NewAccountObjectIdentifier("integration1")
		integration2Id := NewAccountObjectIdentifier("integration2")
		opts := defaultOpts()
		opts.Set = &ServiceSet{
			MinInstances:      Pointer(2),
			MaxInstances:      Pointer(5),
			AutoSuspendSecs:   Pointer(600),
			MinReadyInstances: Pointer(1),
			QueryWarehouse:    Pointer(warehouseId),
			AutoResume:        Bool(true),
			ExternalAccessIntegrations: &ServiceExternalAccessIntegrations{
				ExternalAccessIntegrations: []AccountObjectIdentifier{
					integration1Id,
					integration2Id,
				},
			},
			Comment: &comment,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SERVICE %s SET MIN_INSTANCES = 2 MAX_INSTANCES = 5 AUTO_SUSPEND_SECS = 600 MIN_READY_INSTANCES = 1 QUERY_WAREHOUSE = %v AUTO_RESUME = true`+
			` EXTERNAL_ACCESS_INTEGRATIONS = (%s, %s) COMMENT = '%v'`, id.FullyQualifiedName(), warehouseId.FullyQualifiedName(), integration1Id.FullyQualifiedName(), integration2Id.FullyQualifiedName(), comment)
	})

	t.Run("with unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ServiceUnset{
			MinInstances:               Bool(true),
			AutoSuspendSecs:            Bool(true),
			MaxInstances:               Bool(true),
			MinReadyInstances:          Bool(true),
			QueryWarehouse:             Bool(true),
			AutoResume:                 Bool(true),
			ExternalAccessIntegrations: Bool(true),
			Comment:                    Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SERVICE %s UNSET MIN_INSTANCES, AUTO_SUSPEND_SECS, MAX_INSTANCES, MIN_READY_INSTANCES, QUERY_WAREHOUSE, AUTO_RESUME, EXTERNAL_ACCESS_INTEGRATIONS, COMMENT", id.FullyQualifiedName())
	})

	t.Run("with set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag1"),
				Value: "value1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SERVICE %s SET TAG \"tag1\" = 'value1'", id.FullyQualifiedName())
	})

	t.Run("with unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag1"),
			NewAccountObjectIdentifier("tag2"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SERVICE %s UNSET TAG \"tag1\", \"tag2\"", id.FullyQualifiedName())
	})
}

func TestServices_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	// Minimal valid DropServiceOptions
	defaultOpts := func() *DropServiceOptions {
		return &DropServiceOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropServiceOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DROP SERVICE %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Force = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP SERVICE IF EXISTS %s FORCE", id.FullyQualifiedName())
	})
}

func TestServices_Show(t *testing.T) {
	// Minimal valid ShowServiceOptions
	defaultOpts := func() *ShowServiceOptions {
		return &ShowServiceOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowServiceOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW SERVICES")
	})

	t.Run("validation: conflicting fields for [opts.Job opts.ExcludeJobs]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Job = Bool(true)
		opts.ExcludeJobs = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("ShowServiceOptions", "Job", "ExcludeJobs"))
	})

	t.Run("with jobs option", func(t *testing.T) {
		opts := defaultOpts()
		opts.Job = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "SHOW JOB SERVICES")
	})

	t.Run("with exclude jobs", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExcludeJobs = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "SHOW SERVICES EXCLUDE JOBS")
	})

	t.Run("with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{Pattern: String("service_*")}
		assertOptsValidAndSQLEquals(t, opts, "SHOW SERVICES LIKE 'service_*'")
	})

	t.Run("in schema", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &ServiceIn{
			In: In{
				Database: NewAccountObjectIdentifier("database"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW SERVICES IN DATABASE "database"`)
	})

	t.Run("in compute pool", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &ServiceIn{ComputePool: NewAccountObjectIdentifier("compute_pool")}
		assertOptsValidAndSQLEquals(t, opts, `SHOW SERVICES IN COMPUTE POOL "compute_pool"`)
	})

	t.Run("with starts with", func(t *testing.T) {
		opts := defaultOpts()
		opts.StartsWith = String("my_prefix")
		assertOptsValidAndSQLEquals(t, opts, "SHOW SERVICES STARTS WITH 'my_prefix'")
	})

	t.Run("with limit", func(t *testing.T) {
		opts := defaultOpts()
		opts.Limit = &LimitFrom{Rows: Pointer(10)}
		assertOptsValidAndSQLEquals(t, opts, "SHOW SERVICES LIMIT 10")
	})

	t.Run("with limit and from", func(t *testing.T) {
		opts := defaultOpts()
		opts.Limit = &LimitFrom{Rows: Pointer(10), From: String("service1")}
		assertOptsValidAndSQLEquals(t, opts, "SHOW SERVICES LIMIT 10 FROM 'service1'")
	})
}

func TestServices_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	// Minimal valid DescribeServiceOptions
	defaultOpts := func() *DescribeServiceOptions {
		return &DescribeServiceOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeServiceOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE SERVICE %s", id.FullyQualifiedName())
	})
}
