package sdk

func init() {
	externalVolumeId := NewAccountObjectIdentifier("vol1")
	catalogId := NewAccountObjectIdentifier("cat1")
	tagId1 := randomSchemaObjectIdentifier()
	tagId2 := randomSchemaObjectIdentifier()
	contactId := randomSchemaObjectIdentifier()
	aggregationPolicyId := randomSchemaObjectIdentifier()
	rowAccessPolicyId := randomSchemaObjectIdentifier()
	maskingPolicyId := randomSchemaObjectIdentifier()
	projectionPolicyId := randomSchemaObjectIdentifier()
	fkRefId := randomSchemaObjectIdentifier()
	joinPolicyId := randomSchemaObjectIdentifier()
	rowAccessPolicy1Id := randomSchemaObjectIdentifier()
	rowAccessPolicy2Id := randomSchemaObjectIdentifier()

	icebergTablesTests.Create.
		withModify(case_IcebergTables_validation_Create_opts_ConflictingFields_PartitionBy_ClusterBy, func(opts *CreateIcebergTableOptions) {
			opts.PartitionBy = []IcebergTablePartitionExpression{{Identity: new("ID")}}
			opts.ClusterBy = []string{`"ID"`}
		}).
		withModify(case_IcebergTables_validation_Create_opts_PartitionBy_ExactlyOneValueSet_MoreThanOneSet, func(opts *CreateIcebergTableOptions) {
			opts.PartitionBy = []IcebergTablePartitionExpression{{
				Identity: new("ID"),
				Bucket:   &IcebergTablePartitionBucket{Args: IcebergTablePartitionBucketArgs{NumBuckets: 4, Column: "NAME"}},
			}}
		}).
		withModify(case_IcebergTables_validation_Create_opts_PartitionBy_ExactlyOneValueSet_OneValidOneInvalid, func(opts *CreateIcebergTableOptions) {
			opts.PartitionBy = []IcebergTablePartitionExpression{{Identity: new("ID")}, {}}
		}).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_ExactlyOneOf_NoneSet",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{}}},
				}
			},
			errExactlyOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint", "UniquePK", "FK", "CH"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_MoreThanOneOf_UniquePK_FK",
			func(opts *CreateIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						UniquePK: &TableColumnInlineUniquePK{Unique: new(true)},
						FK:       &TableColumnInlineFK{References: refId},
					}}},
				}
			},
			errExactlyOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint", "UniquePK", "FK", "CH"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_MoreThanOneOf_UniquePK_CH",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						UniquePK: &TableColumnInlineUniquePK{Unique: new(true)},
						CH:       &TableColumnInlineCH{Expression: "ID > 0"},
					}}},
				}
			},
			errExactlyOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint", "UniquePK", "FK", "CH"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_MoreThanOneOf_FK_CH",
			func(opts *CreateIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						FK: &TableColumnInlineFK{References: refId},
						CH: &TableColumnInlineCH{Expression: "ID > 0"},
					}}},
				}
			},
			errExactlyOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint", "UniquePK", "FK", "CH"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_MoreThanOneOf_UniquePK_FK_CH",
			func(opts *CreateIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						UniquePK: &TableColumnInlineUniquePK{Unique: new(true)},
						FK:       &TableColumnInlineFK{References: refId},
						CH:       &TableColumnInlineCH{Expression: "ID > 0"},
					}}},
				}
			},
			errExactlyOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint", "UniquePK", "FK", "CH"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_UniquePK_ExactlyOneOf_Unique_PrimaryKey",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						UniquePK: &TableColumnInlineUniquePK{},
					}}},
				}
			},
			errExactlyOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint.UniquePK", "Unique", "PrimaryKey"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_UniquePK_Enforced_NotEnforced",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						UniquePK: &TableColumnInlineUniquePK{Unique: new(true), Enforced: new(true), NotEnforced: new(true)},
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint.UniquePK", "Enforced", "NotEnforced"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_UniquePK_Deferrable_NotDeferrable",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						UniquePK: &TableColumnInlineUniquePK{Unique: new(true), Deferrable: new(true), NotDeferrable: new(true)},
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint.UniquePK", "Deferrable", "NotDeferrable"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_UniquePK_InitiallyDeferred_InitiallyImmediate",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						UniquePK: &TableColumnInlineUniquePK{Unique: new(true), InitiallyDeferred: new(true), InitiallyImmediate: new(true)},
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint.UniquePK", "InitiallyDeferred", "InitiallyImmediate"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_UniquePK_Enable_Disable",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						UniquePK: &TableColumnInlineUniquePK{Unique: new(true), Enable: new(true), Disable: new(true)},
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint.UniquePK", "Enable", "Disable"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_UniquePK_Validate_Novalidate",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						UniquePK: &TableColumnInlineUniquePK{Unique: new(true), Validate: new(true), Novalidate: new(true)},
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint.UniquePK", "Validate", "Novalidate"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_UniquePK_Rely_Norely",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						UniquePK: &TableColumnInlineUniquePK{Unique: new(true), Rely: new(true), Norely: new(true)},
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint.UniquePK", "Rely", "Norely"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_FK_ValidReferences",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						FK: &TableColumnInlineFK{References: emptySchemaObjectIdentifier},
					}}},
				}
			},
			ErrInvalidObjectIdentifier,
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_FK_Enforced_NotEnforced",
			func(opts *CreateIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						FK: &TableColumnInlineFK{References: refId, Enforced: new(true), NotEnforced: new(true)},
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint.FK", "Enforced", "NotEnforced"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_FK_Deferrable_NotDeferrable",
			func(opts *CreateIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						FK: &TableColumnInlineFK{References: refId, Deferrable: new(true), NotDeferrable: new(true)},
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint.FK", "Deferrable", "NotDeferrable"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_FK_InitiallyDeferred_InitiallyImmediate",
			func(opts *CreateIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						FK: &TableColumnInlineFK{References: refId, InitiallyDeferred: new(true), InitiallyImmediate: new(true)},
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint.FK", "InitiallyDeferred", "InitiallyImmediate"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_FK_Enable_Disable",
			func(opts *CreateIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						FK: &TableColumnInlineFK{References: refId, Enable: new(true), Disable: new(true)},
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint.FK", "Enable", "Disable"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_FK_Validate_Novalidate",
			func(opts *CreateIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						FK: &TableColumnInlineFK{References: refId, Validate: new(true), Novalidate: new(true)},
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint.FK", "Validate", "Novalidate"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_FK_Rely_Norely",
			func(opts *CreateIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						FK: &TableColumnInlineFK{References: refId, Rely: new(true), Norely: new(true)},
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint.FK", "Rely", "Norely"),
		).
		withAdditionalValidationCase(
			"validation_Create_InlineConstraint_CH_EnableValidate_EnableNovalidate",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber, InlineConstraint: &TableColumnInlineConstraint{
						CH: &TableColumnInlineCH{Expression: "ID > 0", EnableValidate: new(true), EnableNovalidate: new(true)},
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.Columns[0].InlineConstraint.CH", "EnableValidate", "EnableNovalidate"),
		).
		withAdditionalValidationCase(
			"validation_Create_OutOfLineConstraint_ExactlyOneOf_NoneSet",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns:             []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{}},
				}
			},
			errExactlyOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.OutOfLineConstraint[0]", "UniquePK", "FK", "CH"),
		).
		withAdditionalValidationCase(
			"validation_Create_OutOfLineConstraint_MoreThanOneOf_UniquePK_FK",
			func(opts *CreateIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{
						UniquePK: &TableOutOfLineUniquePK{Unique: new(true), Columns: []Column{{Value: "ID"}}},
						FK:       &TableOutOfLineFK{References: refId, Columns: []Column{{Value: "ID"}}},
					}},
				}
			},
			errExactlyOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.OutOfLineConstraint[0]", "UniquePK", "FK", "CH"),
		).
		withAdditionalValidationCase(
			"validation_Create_OutOfLineConstraint_UniquePK_ExactlyOneOf_Unique_PrimaryKey",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns:             []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{UniquePK: &TableOutOfLineUniquePK{Columns: []Column{{Value: "ID"}}}}},
				}
			},
			errExactlyOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.OutOfLineConstraint[0].UniquePK", "Unique", "PrimaryKey"),
		).
		withAdditionalValidationCase(
			"validation_Create_OutOfLineConstraint_UniquePK_Enforced_NotEnforced",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{UniquePK: &TableOutOfLineUniquePK{
						Unique: new(true), Columns: []Column{{Value: "ID"}}, Enforced: new(true), NotEnforced: new(true),
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.OutOfLineConstraint[0].UniquePK", "Enforced", "NotEnforced"),
		).
		withAdditionalValidationCase(
			"validation_Create_OutOfLineConstraint_UniquePK_Deferrable_NotDeferrable",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{UniquePK: &TableOutOfLineUniquePK{
						Unique: new(true), Columns: []Column{{Value: "ID"}}, Deferrable: new(true), NotDeferrable: new(true),
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.OutOfLineConstraint[0].UniquePK", "Deferrable", "NotDeferrable"),
		).
		withAdditionalValidationCase(
			"validation_Create_OutOfLineConstraint_UniquePK_InitiallyDeferred_InitiallyImmediate",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{UniquePK: &TableOutOfLineUniquePK{
						Unique: new(true), Columns: []Column{{Value: "ID"}}, InitiallyDeferred: new(true), InitiallyImmediate: new(true),
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.OutOfLineConstraint[0].UniquePK", "InitiallyDeferred", "InitiallyImmediate"),
		).
		withAdditionalValidationCase(
			"validation_Create_OutOfLineConstraint_UniquePK_Enable_Disable",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{UniquePK: &TableOutOfLineUniquePK{
						Unique: new(true), Columns: []Column{{Value: "ID"}}, Enable: new(true), Disable: new(true),
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.OutOfLineConstraint[0].UniquePK", "Enable", "Disable"),
		).
		withAdditionalValidationCase(
			"validation_Create_OutOfLineConstraint_UniquePK_Validate_Novalidate",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{UniquePK: &TableOutOfLineUniquePK{
						Unique: new(true), Columns: []Column{{Value: "ID"}}, Validate: new(true), Novalidate: new(true),
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.OutOfLineConstraint[0].UniquePK", "Validate", "Novalidate"),
		).
		withAdditionalValidationCase(
			"validation_Create_OutOfLineConstraint_UniquePK_Rely_Norely",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{UniquePK: &TableOutOfLineUniquePK{
						Unique: new(true), Columns: []Column{{Value: "ID"}}, Rely: new(true), Norely: new(true),
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.OutOfLineConstraint[0].UniquePK", "Rely", "Norely"),
		).
		withAdditionalValidationCase(
			"validation_Create_OutOfLineConstraint_FK_ValidReferences",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{FK: &TableOutOfLineFK{
						Columns: []Column{{Value: "ID"}}, References: emptySchemaObjectIdentifier,
					}}},
				}
			},
			ErrInvalidObjectIdentifier,
		).
		withAdditionalValidationCase(
			"validation_Create_OutOfLineConstraint_FK_Enforced_NotEnforced",
			func(opts *CreateIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{FK: &TableOutOfLineFK{
						Columns: []Column{{Value: "ID"}}, References: refId, Enforced: new(true), NotEnforced: new(true),
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.OutOfLineConstraint[0].FK", "Enforced", "NotEnforced"),
		).
		withAdditionalValidationCase(
			"validation_Create_OutOfLineConstraint_FK_Deferrable_NotDeferrable",
			func(opts *CreateIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{FK: &TableOutOfLineFK{
						Columns: []Column{{Value: "ID"}}, References: refId, Deferrable: new(true), NotDeferrable: new(true),
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.OutOfLineConstraint[0].FK", "Deferrable", "NotDeferrable"),
		).
		withAdditionalValidationCase(
			"validation_Create_OutOfLineConstraint_FK_InitiallyDeferred_InitiallyImmediate",
			func(opts *CreateIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{FK: &TableOutOfLineFK{
						Columns: []Column{{Value: "ID"}}, References: refId, InitiallyDeferred: new(true), InitiallyImmediate: new(true),
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.OutOfLineConstraint[0].FK", "InitiallyDeferred", "InitiallyImmediate"),
		).
		withAdditionalValidationCase(
			"validation_Create_OutOfLineConstraint_FK_Enable_Disable",
			func(opts *CreateIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{FK: &TableOutOfLineFK{
						Columns: []Column{{Value: "ID"}}, References: refId, Enable: new(true), Disable: new(true),
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.OutOfLineConstraint[0].FK", "Enable", "Disable"),
		).
		withAdditionalValidationCase(
			"validation_Create_OutOfLineConstraint_FK_Validate_Novalidate",
			func(opts *CreateIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{FK: &TableOutOfLineFK{
						Columns: []Column{{Value: "ID"}}, References: refId, Validate: new(true), Novalidate: new(true),
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.OutOfLineConstraint[0].FK", "Validate", "Novalidate"),
		).
		withAdditionalValidationCase(
			"validation_Create_OutOfLineConstraint_FK_Rely_Norely",
			func(opts *CreateIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{FK: &TableOutOfLineFK{
						Columns: []Column{{Value: "ID"}}, References: refId, Rely: new(true), Norely: new(true),
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.OutOfLineConstraint[0].FK", "Rely", "Norely"),
		).
		withAdditionalValidationCase(
			"validation_Create_OutOfLineConstraint_CH_EnableValidate_EnableNovalidate",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{CH: &TableOutOfLineCH{
						Expression: "ID > 0", EnableValidate: new(true), EnableNovalidate: new(true),
					}}},
				}
			},
			errOneOf("CreateIcebergTableOptions.ColumnsAndConstraints.OutOfLineConstraint[0].CH", "EnableValidate", "EnableNovalidate"),
		).
		withExpectedSqlf(
			case_IcebergTables_sql_Create_basic,
			`CREATE ICEBERG TABLE %s`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Create_all,
			func(opts *CreateIcebergTableOptions) {
				opts.IfNotExists = new(true)
				opts.Transient = new(true)
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{
						{
							Name:       "ID",
							ColumnType: dataTypeNumber,
							DefaultValue: &ColumnDefaultValue{
								Expression: new("1"),
							},
							NotNull: new(true),
							MaskingPolicy: &TableColumnMaskingPolicy{
								MaskingPolicy: maskingPolicyId,
								Using:         []Column{{"ID"}, {"NAME"}},
							},
							ProjectionPolicy: &TableColumnProjectionPolicy{
								ProjectionPolicy: projectionPolicyId,
							},
							Tag: []TagAssociation{
								{Name: tagId1, Value: "v1"},
								{Name: tagId2, Value: "v2"},
							},
							Comment: new("id column"),
						},
						{
							Name:       "NAME",
							ColumnType: dataTypeVarchar,
							Comment:    new("name column"),
						},
					},
				}
				opts.PartitionBy = []IcebergTablePartitionExpression{
					{Identity: new("ID")},
					{Bucket: &IcebergTablePartitionBucket{Args: IcebergTablePartitionBucketArgs{NumBuckets: 4, Column: "NAME"}}},
					{Truncate: &IcebergTablePartitionTruncate{Args: IcebergTablePartitionTruncateArgs{Width: 10, Column: "C1"}}},
					{Year: &IcebergTablePartitionYear{Args: IcebergTablePartitionTimeArgs{Column: "C2"}}},
					{Month: &IcebergTablePartitionMonth{Args: IcebergTablePartitionTimeArgs{Column: "C3"}}},
					{Day: &IcebergTablePartitionDay{Args: IcebergTablePartitionTimeArgs{Column: "C4"}}},
					{Hour: &IcebergTablePartitionHour{Args: IcebergTablePartitionTimeArgs{Column: "C5"}}},
				}
				opts.PathLayout = new(IcebergTablePathLayoutHierarchical)
				opts.ExternalVolume = new(NewAccountObjectIdentifier("vol1"))
				opts.Catalog = new(IcebergTableCatalogSnowflake)
				opts.BaseLocation = new("base/loc")
				opts.TargetFileSize = new(IcebergTableTargetFileSize64mb)
				opts.CatalogSync = new("integration1")
				opts.StorageSerializationPolicy = new(StorageSerializationPolicyOptimized)
				opts.DataRetentionTimeInDays = new(1)
				opts.MaxDataExtensionTimeInDays = new(2)
				opts.ChangeTracking = new(true)
				opts.ErrorLogging = new(true)
				opts.Comment = new("test comment")
				opts.IcebergVersion = new(2)
				opts.EnableIcebergMergeOnRead = new(true)
				opts.RowAccessPolicy = &IcebergTableRowAccessPolicy{
					Name: rowAccessPolicyId,
					On:   []Column{{"ID"}, {"NAME"}},
				}
				opts.AggregationPolicy = &IcebergTableAggregationPolicy{
					AggregationPolicy: aggregationPolicyId,
					EntityKey:         []Column{{"ID"}},
				}
				opts.Tag = []TagAssociation{
					{Name: tagId1, Value: "v1"},
					{Name: tagId2, Value: "v2"},
				}
				opts.EnableDataCompaction = new(true)
				opts.Contact = []TableContact{
					{Purpose: "SUPPORT", Contact: contactId},
					{Purpose: "ACCESS_APPROVAL", Contact: contactId},
				}
			},
			`CREATE TRANSIENT ICEBERG TABLE IF NOT EXISTS %s `+
				`("ID" NUMBER(38, 0) DEFAULT 1 NOT NULL MASKING POLICY %s USING ("ID", "NAME") PROJECTION POLICY %s TAG (%s = 'v1', %s = 'v2') COMMENT 'id column', `+
				`"NAME" VARCHAR(16777216) COMMENT 'name column') `+
				`PARTITION BY ("ID", BUCKET (4, "NAME"), TRUNCATE (10, "C1"), YEAR ("C2"), MONTH ("C3"), DAY ("C4"), HOUR ("C5")) `+
				`PATH_LAYOUT = HIERARCHICAL `+
				`EXTERNAL_VOLUME = '\"vol1\"' `+
				`CATALOG = 'SNOWFLAKE' `+
				`BASE_LOCATION = 'base/loc' `+
				`TARGET_FILE_SIZE = '64MB' `+
				`CATALOG_SYNC = 'integration1' `+
				`STORAGE_SERIALIZATION_POLICY = OPTIMIZED `+
				`DATA_RETENTION_TIME_IN_DAYS = 1 `+
				`MAX_DATA_EXTENSION_TIME_IN_DAYS = 2 `+
				`CHANGE_TRACKING = true `+
				`ERROR_LOGGING = true `+
				`COMMENT = 'test comment' `+
				`ICEBERG_VERSION = 2 `+
				`ENABLE_ICEBERG_MERGE_ON_READ = true `+
				`ROW ACCESS POLICY %s ON ("ID", "NAME") `+
				`AGGREGATION POLICY %s ENTITY KEY ("ID") `+
				`TAG (%s = 'v1', %s = 'v2') `+
				`ENABLE_DATA_COMPACTION = true `+
				`WITH CONTACT (SUPPORT = %s, ACCESS_APPROVAL = %s)`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			maskingPolicyId.FullyQualifiedName(),
			projectionPolicyId.FullyQualifiedName(),
			tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
			rowAccessPolicyId.FullyQualifiedName(),
			aggregationPolicyId.FullyQualifiedName(),
			tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
			contactId.FullyQualifiedName(), contactId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateIcebergTableOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
			},
			`CREATE OR REPLACE ICEBERG TABLE %s COPY GRANTS`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_clusterBy",
			func(opts *CreateIcebergTableOptions) { opts.ClusterBy = []string{`"col1"`, `"col2"`} },
			`CREATE ICEBERG TABLE %s CLUSTER BY ("col1", "col2")`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_inlineUnique_allFields",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{
						Name:       "ID",
						ColumnType: dataTypeNumber,
						InlineConstraint: &TableColumnInlineConstraint{
							UniquePK: &TableColumnInlineUniquePK{
								Name:              new("uq_id"),
								Unique:            new(true),
								Enforced:          new(true),
								Deferrable:        new(true),
								InitiallyDeferred: new(true),
								Enable:            new(true),
								Validate:          new(true),
								Rely:              new(true),
							},
						},
					}},
				}
			},
			`CREATE ICEBERG TABLE %s ("ID" NUMBER(38, 0) CONSTRAINT "uq_id" UNIQUE ENFORCED DEFERRABLE INITIALLY DEFERRED ENABLE VALIDATE RELY)`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_inlinePrimaryKey_negatedFlags",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{
						Name:       "ID",
						ColumnType: dataTypeNumber,
						InlineConstraint: &TableColumnInlineConstraint{
							UniquePK: &TableColumnInlineUniquePK{
								Name:               new("pk_id"),
								PrimaryKey:         new(true),
								NotEnforced:        new(true),
								NotDeferrable:      new(true),
								InitiallyImmediate: new(true),
								Disable:            new(true),
								Novalidate:         new(true),
								Norely:             new(true),
							},
						},
					}},
				}
			},
			`CREATE ICEBERG TABLE %s ("ID" NUMBER(38, 0) CONSTRAINT "pk_id" PRIMARY KEY NOT ENFORCED NOT DEFERRABLE INITIALLY IMMEDIATE DISABLE NOVALIDATE NORELY)`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_inlineForeignKey_allFields",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{
						Name:       "ID",
						ColumnType: dataTypeNumber,
						InlineConstraint: &TableColumnInlineConstraint{
							FK: &TableColumnInlineFK{
								Name:       new("fk_id"),
								ForeignKey: new(true),
								References: fkRefId,
								RefColumn:  []Column{{Value: "REF_COL"}},
								Match:      new(FullMatchType),
								On: &ForeignKeyOnAction{
									OnUpdate: new(ForeignKeySetNullAction),
									OnDelete: new(ForeignKeyRestrictAction),
								},
								Enforced:          new(true),
								Deferrable:        new(true),
								InitiallyDeferred: new(true),
								Enable:            new(true),
								Validate:          new(true),
								Rely:              new(true),
							},
						},
					}},
				}
			},
			`CREATE ICEBERG TABLE %s ("ID" NUMBER(38, 0) CONSTRAINT "fk_id" FOREIGN KEY REFERENCES %s ("REF_COL") MATCH FULL ON UPDATE SET NULL ON DELETE RESTRICT ENFORCED DEFERRABLE INITIALLY DEFERRED ENABLE VALIDATE RELY)`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			fkRefId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_inlineCheck_enableValidate",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{
						Name:       "ID",
						ColumnType: dataTypeNumber,
						InlineConstraint: &TableColumnInlineConstraint{
							CH: &TableColumnInlineCH{
								Name:           new("ck_id"),
								Expression:     "ID > 0",
								EnableValidate: new(true),
							},
						},
					}},
				}
			},
			`CREATE ICEBERG TABLE %s ("ID" NUMBER(38, 0) CONSTRAINT "ck_id" CHECK ( ID > 0 ) ENABLE VALIDATE)`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_inlineCheck_enableNovalidate",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{
						Name:       "ID",
						ColumnType: dataTypeNumber,
						InlineConstraint: &TableColumnInlineConstraint{
							CH: &TableColumnInlineCH{
								Name:             new("ck_id"),
								Expression:       "ID > 0",
								EnableNovalidate: new(true),
							},
						},
					}},
				}
			},
			`CREATE ICEBERG TABLE %s ("ID" NUMBER(38, 0) CONSTRAINT "ck_id" CHECK ( ID > 0 ) ENABLE NOVALIDATE)`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_outOfLine_allVariants",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{
						{UniquePK: &TableOutOfLineUniquePK{Name: new("uq_c"), Unique: new(true), Columns: []Column{{Value: "ID"}}}},
						{UniquePK: &TableOutOfLineUniquePK{Name: new("pk_c"), PrimaryKey: new(true), Columns: []Column{{Value: "ID"}}}},
						{FK: &TableOutOfLineFK{
							Name:       new("fk_c"),
							Columns:    []Column{{Value: "ID"}},
							References: fkRefId,
							RefColumns: []Column{{Value: "COL_A"}, {Value: "COL_B"}},
							Match:      new(SimpleMatchType),
							On: &ForeignKeyOnAction{
								OnUpdate: new(ForeignKeyCascadeAction),
								OnDelete: new(ForeignKeyNoAction),
							},
						}},
					},
				}
			},
			`CREATE ICEBERG TABLE %s ("ID" NUMBER(38, 0), CONSTRAINT "uq_c" UNIQUE ("ID"), CONSTRAINT "pk_c" PRIMARY KEY ("ID"), CONSTRAINT "fk_c" FOREIGN KEY ("ID") REFERENCES %s ("COL_A", "COL_B") MATCH SIMPLE ON UPDATE CASCADE ON DELETE NO ACTION)`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			fkRefId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_outOfLineUnique_allFieldsWithComment",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{
						UniquePK: &TableOutOfLineUniquePK{
							Name:              new("uq_full"),
							Unique:            new(true),
							Columns:           []Column{{Value: "ID"}},
							Enforced:          new(true),
							Deferrable:        new(true),
							InitiallyDeferred: new(true),
							Enable:            new(true),
							Validate:          new(true),
							Rely:              new(true),
							Comment:           new("uq note"),
						},
					}},
				}
			},
			`CREATE ICEBERG TABLE %s ("ID" NUMBER(38, 0), CONSTRAINT "uq_full" UNIQUE ("ID") ENFORCED DEFERRABLE INITIALLY DEFERRED ENABLE VALIDATE RELY COMMENT 'uq note')`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_outOfLinePrimaryKey_negatedFlagsWithComment",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{
						UniquePK: &TableOutOfLineUniquePK{
							Name:               new("pk_full"),
							PrimaryKey:         new(true),
							Columns:            []Column{{Value: "ID"}},
							NotEnforced:        new(true),
							NotDeferrable:      new(true),
							InitiallyImmediate: new(true),
							Disable:            new(true),
							Novalidate:         new(true),
							Norely:             new(true),
							Comment:            new("pk note"),
						},
					}},
				}
			},
			`CREATE ICEBERG TABLE %s ("ID" NUMBER(38, 0), CONSTRAINT "pk_full" PRIMARY KEY ("ID") NOT ENFORCED NOT DEFERRABLE INITIALLY IMMEDIATE DISABLE NOVALIDATE NORELY COMMENT 'pk note')`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_outOfLineForeignKey_allFieldsWithComment",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{
						FK: &TableOutOfLineFK{
							Name:       new("fk_full"),
							Columns:    []Column{{Value: "ID"}},
							References: fkRefId,
							RefColumns: []Column{{Value: "COL_A"}},
							Match:      new(PartialMatchType),
							On: &ForeignKeyOnAction{
								OnUpdate: new(ForeignKeySetNullAction),
								OnDelete: new(ForeignKeySetDefaultAction),
							},
							Enforced:          new(true),
							Deferrable:        new(true),
							InitiallyDeferred: new(true),
							Enable:            new(true),
							Validate:          new(true),
							Rely:              new(true),
							Comment:           new("fk note"),
						},
					}},
				}
			},
			`CREATE ICEBERG TABLE %s ("ID" NUMBER(38, 0), CONSTRAINT "fk_full" FOREIGN KEY ("ID") REFERENCES %s ("COL_A") MATCH PARTIAL ON UPDATE SET NULL ON DELETE SET DEFAULT ENFORCED DEFERRABLE INITIALLY DEFERRED ENABLE VALIDATE RELY COMMENT 'fk note')`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			fkRefId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_outOfLineCheck_enableValidate",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{CH: &TableOutOfLineCH{
						Name:           new("ck_out"),
						Expression:     "ID > 0",
						EnableValidate: new(true),
					}}},
				}
			},
			`CREATE ICEBERG TABLE %s ("ID" NUMBER(38, 0), CONSTRAINT "ck_out" CHECK ( ID > 0 ) ENABLE VALIDATE)`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_outOfLineCheck_enableNovalidate",
			func(opts *CreateIcebergTableOptions) {
				opts.ColumnsAndConstraints = IcebergTableColumnsAndConstraints{
					Columns: []IcebergTableColumn{{Name: "ID", ColumnType: dataTypeNumber}},
					OutOfLineConstraint: []TableOutOfLineConstraint{{CH: &TableOutOfLineCH{
						Name:             new("ck_out"),
						Expression:       "ID > 0",
						EnableNovalidate: new(true),
					}}},
				}
			},
			`CREATE ICEBERG TABLE %s ("ID" NUMBER(38, 0), CONSTRAINT "ck_out" CHECK ( ID > 0 ) ENABLE NOVALIDATE)`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	icebergTablesTests.CreateFromIcebergFiles.
		withDefaultOpts(func() *CreateFromIcebergFilesIcebergTableOptions {
			return &CreateFromIcebergFilesIcebergTableOptions{
				name:             icebergTablesTestIdSchemaObjectIdentifier,
				MetadataFilePath: "metadata/v1.metadata.json",
			}
		}).
		withExpectedSqlf(
			case_IcebergTables_sql_CreateFromIcebergFiles_basic,
			`CREATE ICEBERG TABLE %s METADATA_FILE_PATH = 'metadata/v1.metadata.json'`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_CreateFromIcebergFiles_all,
			func(opts *CreateFromIcebergFilesIcebergTableOptions) {
				opts.IfNotExists = new(true)
				opts.ExternalVolume = &externalVolumeId
				opts.Catalog = &catalogId
				opts.ReplaceInvalidCharacters = new(true)
				opts.Comment = new("some comment")
				opts.Tag = []TagAssociation{
					{Name: tagId1, Value: "v1"},
					{Name: tagId2, Value: "v2"},
				}
				opts.Contact = []TableContact{{Purpose: "SUPPORT", Contact: contactId}}
			},
			`CREATE ICEBERG TABLE IF NOT EXISTS %s `+
				`EXTERNAL_VOLUME = '\"%s\"' `+
				`CATALOG = '\"%s\"' `+
				`METADATA_FILE_PATH = 'metadata/v1.metadata.json' `+
				`REPLACE_INVALID_CHARACTERS = true `+
				`COMMENT = 'some comment' `+
				`TAG (%s = 'v1', %s = 'v2') `+
				`WITH CONTACT (SUPPORT = %s)`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			externalVolumeId.Name(),
			catalogId.Name(),
			tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
			contactId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateFromIcebergFiles_orReplace",
			func(opts *CreateFromIcebergFilesIcebergTableOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE ICEBERG TABLE %s METADATA_FILE_PATH = 'metadata/v1.metadata.json'`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	icebergTablesTests.CreateFromDeltaLake.
		withDefaultOpts(func() *CreateFromDeltaLakeIcebergTableOptions {
			return &CreateFromDeltaLakeIcebergTableOptions{
				name:         icebergTablesTestIdSchemaObjectIdentifier,
				BaseLocation: "my/base/location",
			}
		}).
		withExpectedSqlf(
			case_IcebergTables_sql_CreateFromDeltaLake_basic,
			`CREATE ICEBERG TABLE %s BASE_LOCATION = 'my/base/location'`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_CreateFromDeltaLake_all,
			func(opts *CreateFromDeltaLakeIcebergTableOptions) {
				opts.IfNotExists = new(true)
				opts.ExternalVolume = &externalVolumeId
				opts.Catalog = &catalogId
				opts.ReplaceInvalidCharacters = new(true)
				opts.AutoRefresh = new(true)
				opts.Comment = new("some comment")
				opts.Tag = []TagAssociation{
					{Name: tagId1, Value: "v1"},
					{Name: tagId2, Value: "v2"},
				}
				opts.Contact = []TableContact{{Purpose: "SUPPORT", Contact: contactId}}
			},
			`CREATE ICEBERG TABLE IF NOT EXISTS %s `+
				`EXTERNAL_VOLUME = '\"%s\"' `+
				`CATALOG = '\"%s\"' `+
				`BASE_LOCATION = 'my/base/location' `+
				`REPLACE_INVALID_CHARACTERS = true `+
				`AUTO_REFRESH = true `+
				`COMMENT = 'some comment' `+
				`TAG (%s = 'v1', %s = 'v2') `+
				`WITH CONTACT (SUPPORT = %s)`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			externalVolumeId.Name(),
			catalogId.Name(),
			tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
			contactId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateFromDeltaLake_orReplace",
			func(opts *CreateFromDeltaLakeIcebergTableOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE ICEBERG TABLE %s BASE_LOCATION = 'my/base/location'`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	icebergTablesTests.CreateFromIcebergRest.
		withDefaultOpts(func() *CreateFromIcebergRestIcebergTableOptions {
			return &CreateFromIcebergRestIcebergTableOptions{
				name:             icebergTablesTestIdSchemaObjectIdentifier,
				CatalogTableName: "my_remote_table",
			}
		}).
		withExpectedSqlf(
			case_IcebergTables_sql_CreateFromIcebergRest_basic,
			`CREATE ICEBERG TABLE %s CATALOG_TABLE_NAME = 'my_remote_table'`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_CreateFromIcebergRest_all,
			func(opts *CreateFromIcebergRestIcebergTableOptions) {
				opts.IfNotExists = new(true)
				opts.ExternalVolume = &externalVolumeId
				opts.Catalog = &catalogId
				opts.CatalogNamespace = new("my_namespace")
				opts.PathLayout = new(IcebergTablePathLayoutHierarchical)
				opts.TargetFileSize = new(IcebergTableTargetFileSize64mb)
				opts.ReplaceInvalidCharacters = new(true)
				opts.AutoRefresh = new(true)
				opts.Comment = new("some comment")
				opts.StorageSerializationPolicy = new(StorageSerializationPolicyOptimized)
				opts.IcebergMergeOnReadBehavior = new(IcebergTableIcebergMergeOnReadBehaviorEnabled)
				opts.EnableIcebergMergeOnRead = new(true)
				opts.Tag = []TagAssociation{
					{Name: tagId1, Value: "v1"},
					{Name: tagId2, Value: "v2"},
				}
				opts.Contact = []TableContact{{Purpose: "SUPPORT", Contact: contactId}}
			},
			`CREATE ICEBERG TABLE IF NOT EXISTS %s `+
				`EXTERNAL_VOLUME = '\"%s\"' `+
				`CATALOG = '\"%s\"' `+
				`CATALOG_TABLE_NAME = 'my_remote_table' `+
				`CATALOG_NAMESPACE = 'my_namespace' `+
				`PATH_LAYOUT = HIERARCHICAL `+
				`TARGET_FILE_SIZE = '64MB' `+
				`REPLACE_INVALID_CHARACTERS = true `+
				`AUTO_REFRESH = true `+
				`COMMENT = 'some comment' `+
				`STORAGE_SERIALIZATION_POLICY = OPTIMIZED `+
				`ICEBERG_MERGE_ON_READ_BEHAVIOR = 'ENABLED' `+
				`ENABLE_ICEBERG_MERGE_ON_READ = true `+
				`TAG (%s = 'v1', %s = 'v2') `+
				`WITH CONTACT (SUPPORT = %s)`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			externalVolumeId.Name(),
			catalogId.Name(),
			tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
			contactId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateFromIcebergRest_orReplace",
			func(opts *CreateFromIcebergRestIcebergTableOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE ICEBERG TABLE %s CATALOG_TABLE_NAME = 'my_remote_table'`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	icebergTablesTests.CreateFromAwsGlue.
		withDefaultOpts(func() *CreateFromAwsGlueIcebergTableOptions {
			return &CreateFromAwsGlueIcebergTableOptions{
				name:             icebergTablesTestIdSchemaObjectIdentifier,
				CatalogTableName: "my_remote_table",
			}
		}).
		withExpectedSqlf(
			case_IcebergTables_sql_CreateFromAwsGlue_basic,
			`CREATE ICEBERG TABLE %s CATALOG_TABLE_NAME = 'my_remote_table'`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_CreateFromAwsGlue_all,
			func(opts *CreateFromAwsGlueIcebergTableOptions) {
				opts.IfNotExists = new(true)
				opts.ExternalVolume = &externalVolumeId
				opts.Catalog = &catalogId
				opts.CatalogNamespace = new("my_namespace")
				opts.ReplaceInvalidCharacters = new(true)
				opts.AutoRefresh = new(true)
				opts.Comment = new("some comment")
				opts.Tag = []TagAssociation{
					{Name: tagId1, Value: "v1"},
					{Name: tagId2, Value: "v2"},
				}
				opts.Contact = []TableContact{{Purpose: "SUPPORT", Contact: contactId}}
			},
			`CREATE ICEBERG TABLE IF NOT EXISTS %s `+
				`EXTERNAL_VOLUME = '\"%s\"' `+
				`CATALOG = '\"%s\"' `+
				`CATALOG_TABLE_NAME = 'my_remote_table' `+
				`CATALOG_NAMESPACE = 'my_namespace' `+
				`REPLACE_INVALID_CHARACTERS = true `+
				`AUTO_REFRESH = true `+
				`COMMENT = 'some comment' `+
				`TAG (%s = 'v1', %s = 'v2') `+
				`WITH CONTACT (SUPPORT = %s)`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			externalVolumeId.Name(),
			catalogId.Name(),
			tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
			contactId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateFromAwsGlue_orReplace",
			func(opts *CreateFromAwsGlueIcebergTableOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE ICEBERG TABLE %s CATALOG_TABLE_NAME = 'my_remote_table'`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	icebergTablesTests.Alter.
		withModify(case_IcebergTables_validation_Alter_opts_AlterColumnAction_ExactlyOneValueSet_MoreThanOneSet, func(opts *AlterIcebergTableOptions) {
			opts.AlterColumnAction = []IcebergTableAlterColumnAction{{ColumnName: "col1", SetNotNull: new(true), DropNotNull: new(true)}}
		}).
		withModify(case_IcebergTables_validation_Alter_opts_AlterColumnAction_ExactlyOneValueSet_OneValidOneInvalid, func(opts *AlterIcebergTableOptions) {
			opts.AlterColumnAction = []IcebergTableAlterColumnAction{{ColumnName: "col1", SetNotNull: new(true)}, {ColumnName: "col2"}}
		}).
		withModify(case_IcebergTables_validation_Alter_opts_ClusteringAction_ExactlyOneValueSet_MoreThanOneSet, func(opts *AlterIcebergTableOptions) {
			opts.ClusteringAction = &IcebergTableClusteringAction{ClusterBy: []string{`"col1"`}, DropClusteringKey: new(true)}
		}).
		withModify(case_IcebergTables_validation_Alter_opts_SearchOptimizationAction_Drop_On_ExactlyOneValueSet_MoreThanOneSet, func(opts *AlterIcebergTableOptions) {
			opts.SearchOptimizationAction = &TableSearchOptimizationAction{Drop: &TableDropSearchOptimization{
				On: []TableDropSearchOptimizationOn{{ColumnName: new("col1"), ExpressionId: new("expr_1")}},
			}}
		}).
		withModify(case_IcebergTables_validation_Alter_opts_SearchOptimizationAction_Drop_On_ExactlyOneValueSet_OneValidOneInvalid, func(opts *AlterIcebergTableOptions) {
			opts.SearchOptimizationAction = &TableSearchOptimizationAction{Drop: &TableDropSearchOptimization{
				On: []TableDropSearchOptimizationOn{{ColumnName: new("col1")}, {}},
			}}
		}).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_ExactlyOneOf_NoneSet",
			func(opts *AlterIcebergTableOptions) {
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{}}
			},
			errExactlyOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint", "UniquePK", "FK", "CH"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_MoreThanOneOf_UniquePK_FK",
			func(opts *AlterIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					UniquePK: &TableColumnInlineUniquePK{Unique: new(true)},
					FK:       &TableColumnInlineFK{References: refId},
				}}
			},
			errExactlyOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint", "UniquePK", "FK", "CH"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_MoreThanOneOf_UniquePK_CH",
			func(opts *AlterIcebergTableOptions) {
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					UniquePK: &TableColumnInlineUniquePK{Unique: new(true)},
					CH:       &TableColumnInlineCH{Expression: "NEW_COL IS NOT NULL"},
				}}
			},
			errExactlyOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint", "UniquePK", "FK", "CH"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_MoreThanOneOf_FK_CH",
			func(opts *AlterIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					FK: &TableColumnInlineFK{References: refId},
					CH: &TableColumnInlineCH{Expression: "NEW_COL IS NOT NULL"},
				}}
			},
			errExactlyOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint", "UniquePK", "FK", "CH"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_MoreThanOneOf_UniquePK_FK_CH",
			func(opts *AlterIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					UniquePK: &TableColumnInlineUniquePK{Unique: new(true)},
					FK:       &TableColumnInlineFK{References: refId},
					CH:       &TableColumnInlineCH{Expression: "NEW_COL IS NOT NULL"},
				}}
			},
			errExactlyOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint", "UniquePK", "FK", "CH"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_UniquePK_ExactlyOneOf_Unique_PrimaryKey",
			func(opts *AlterIcebergTableOptions) {
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					UniquePK: &TableColumnInlineUniquePK{},
				}}
			},
			errExactlyOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.UniquePK", "Unique", "PrimaryKey"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_UniquePK_Enforced_NotEnforced",
			func(opts *AlterIcebergTableOptions) {
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					UniquePK: &TableColumnInlineUniquePK{Unique: new(true), Enforced: new(true), NotEnforced: new(true)},
				}}
			},
			errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.UniquePK", "Enforced", "NotEnforced"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_UniquePK_Deferrable_NotDeferrable",
			func(opts *AlterIcebergTableOptions) {
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					UniquePK: &TableColumnInlineUniquePK{Unique: new(true), Deferrable: new(true), NotDeferrable: new(true)},
				}}
			},
			errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.UniquePK", "Deferrable", "NotDeferrable"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_UniquePK_InitiallyDeferred_InitiallyImmediate",
			func(opts *AlterIcebergTableOptions) {
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					UniquePK: &TableColumnInlineUniquePK{Unique: new(true), InitiallyDeferred: new(true), InitiallyImmediate: new(true)},
				}}
			},
			errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.UniquePK", "InitiallyDeferred", "InitiallyImmediate"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_UniquePK_Enable_Disable",
			func(opts *AlterIcebergTableOptions) {
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					UniquePK: &TableColumnInlineUniquePK{Unique: new(true), Enable: new(true), Disable: new(true)},
				}}
			},
			errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.UniquePK", "Enable", "Disable"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_UniquePK_Validate_Novalidate",
			func(opts *AlterIcebergTableOptions) {
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					UniquePK: &TableColumnInlineUniquePK{Unique: new(true), Validate: new(true), Novalidate: new(true)},
				}}
			},
			errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.UniquePK", "Validate", "Novalidate"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_UniquePK_Rely_Norely",
			func(opts *AlterIcebergTableOptions) {
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					UniquePK: &TableColumnInlineUniquePK{Unique: new(true), Rely: new(true), Norely: new(true)},
				}}
			},
			errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.UniquePK", "Rely", "Norely"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_FK_ValidReferences",
			func(opts *AlterIcebergTableOptions) {
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					FK: &TableColumnInlineFK{References: emptySchemaObjectIdentifier},
				}}
			},
			ErrInvalidObjectIdentifier,
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_FK_Enforced_NotEnforced",
			func(opts *AlterIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					FK: &TableColumnInlineFK{References: refId, Enforced: new(true), NotEnforced: new(true)},
				}}
			},
			errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.FK", "Enforced", "NotEnforced"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_FK_Deferrable_NotDeferrable",
			func(opts *AlterIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					FK: &TableColumnInlineFK{References: refId, Deferrable: new(true), NotDeferrable: new(true)},
				}}
			},
			errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.FK", "Deferrable", "NotDeferrable"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_FK_InitiallyDeferred_InitiallyImmediate",
			func(opts *AlterIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					FK: &TableColumnInlineFK{References: refId, InitiallyDeferred: new(true), InitiallyImmediate: new(true)},
				}}
			},
			errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.FK", "InitiallyDeferred", "InitiallyImmediate"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_FK_Enable_Disable",
			func(opts *AlterIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					FK: &TableColumnInlineFK{References: refId, Enable: new(true), Disable: new(true)},
				}}
			},
			errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.FK", "Enable", "Disable"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_FK_Validate_Novalidate",
			func(opts *AlterIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					FK: &TableColumnInlineFK{References: refId, Validate: new(true), Novalidate: new(true)},
				}}
			},
			errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.FK", "Validate", "Novalidate"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_FK_Rely_Norely",
			func(opts *AlterIcebergTableOptions) {
				refId := randomSchemaObjectIdentifier()
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					FK: &TableColumnInlineFK{References: refId, Rely: new(true), Norely: new(true)},
				}}
			},
			errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.FK", "Rely", "Norely"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddColumnAction_InlineConstraint_CH_EnableValidate_EnableNovalidate",
			func(opts *AlterIcebergTableOptions) {
				opts.AddColumnAction = &IcebergTableAddColumnAction{Name: "NEW_COL", ColumnType: dataTypeVarchar, InlineConstraint: &TableColumnInlineConstraint{
					CH: &TableColumnInlineCH{Expression: "NEW_COL > 0", EnableValidate: new(true), EnableNovalidate: new(true)},
				}}
			},
			errOneOf("AlterIcebergTableOptions.AddColumnAction.InlineConstraint.CH", "EnableValidate", "EnableNovalidate"),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_AddColumnAction,
			func(opts *AlterIcebergTableOptions) {
				opts.IfExists = new(true)
				opts.AddColumnAction = &IcebergTableAddColumnAction{
					IfNotExists: new(true),
					Name:        "NEW_COL",
					ColumnType:  dataTypeVarchar,
					DefaultValue: &ColumnDefaultValue{
						Expression: new("'a'"),
					},
					MaskingPolicy: &TableColumnMaskingPolicy{
						MaskingPolicy: maskingPolicyId,
						Using:         []Column{{"NEW_COL"}, {"OTHER"}},
					},
					ProjectionPolicy: &TableColumnProjectionPolicy{
						ProjectionPolicy: projectionPolicyId,
					},
					Tag: []TagAssociation{
						{Name: tagId1, Value: "v1"},
						{Name: tagId2, Value: "v2"},
					},
				}
			},
			`ALTER ICEBERG TABLE IF EXISTS %s ADD COLUMN IF NOT EXISTS "NEW_COL" VARCHAR(16777216) DEFAULT 'a' MASKING POLICY %s USING ("NEW_COL", "OTHER") PROJECTION POLICY %s TAG (%s = 'v1', %s = 'v2')`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			maskingPolicyId.FullyQualifiedName(),
			projectionPolicyId.FullyQualifiedName(),
			tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_addColumn_inlineUnique",
			func(opts *AlterIcebergTableOptions) {
				opts.AddColumnAction = &IcebergTableAddColumnAction{
					Name:       "NEW_COL",
					ColumnType: dataTypeVarchar,
					InlineConstraint: &TableColumnInlineConstraint{
						UniquePK: &TableColumnInlineUniquePK{
							Name:              new("uq_new"),
							Unique:            new(true),
							Enforced:          new(true),
							Deferrable:        new(true),
							InitiallyDeferred: new(true),
							Enable:            new(true),
							Validate:          new(true),
							Rely:              new(true),
						},
					},
				}
			},
			`ALTER ICEBERG TABLE %s ADD COLUMN "NEW_COL" VARCHAR(16777216) CONSTRAINT "uq_new" UNIQUE ENFORCED DEFERRABLE INITIALLY DEFERRED ENABLE VALIDATE RELY`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_addColumn_inlineForeignKey",
			func(opts *AlterIcebergTableOptions) {
				opts.AddColumnAction = &IcebergTableAddColumnAction{
					Name:       "NEW_COL",
					ColumnType: dataTypeVarchar,
					InlineConstraint: &TableColumnInlineConstraint{
						FK: &TableColumnInlineFK{
							Name:       new("fk_new"),
							ForeignKey: new(true),
							References: fkRefId,
							RefColumn:  []Column{{Value: "REF_COL"}},
							Match:      new(PartialMatchType),
							On: &ForeignKeyOnAction{
								OnUpdate: new(ForeignKeySetDefaultAction),
								OnDelete: new(ForeignKeyCascadeAction),
							},
							NotEnforced:        new(true),
							NotDeferrable:      new(true),
							InitiallyImmediate: new(true),
							Disable:            new(true),
							Novalidate:         new(true),
							Norely:             new(true),
						},
					},
				}
			},
			`ALTER ICEBERG TABLE %s ADD COLUMN "NEW_COL" VARCHAR(16777216) CONSTRAINT "fk_new" FOREIGN KEY REFERENCES %s ("REF_COL") MATCH PARTIAL ON UPDATE SET DEFAULT ON DELETE CASCADE NOT ENFORCED NOT DEFERRABLE INITIALLY IMMEDIATE DISABLE NOVALIDATE NORELY`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			fkRefId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_addColumn_inlineCheck",
			func(opts *AlterIcebergTableOptions) {
				opts.AddColumnAction = &IcebergTableAddColumnAction{
					Name:       "NEW_COL",
					ColumnType: dataTypeVarchar,
					InlineConstraint: &TableColumnInlineConstraint{
						CH: &TableColumnInlineCH{
							Name:           new("ck_new"),
							Expression:     "LENGTH(NEW_COL) > 0",
							EnableValidate: new(true),
						},
					},
				}
			},
			`ALTER ICEBERG TABLE %s ADD COLUMN "NEW_COL" VARCHAR(16777216) CONSTRAINT "ck_new" CHECK ( LENGTH(NEW_COL) > 0 ) ENABLE VALIDATE`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_DropColumnAction,
			func(opts *AlterIcebergTableOptions) {
				opts.DropColumnAction = &TableDropColumnAction{
					IfExists: new(true),
					Columns:  []Column{{"col1"}, {"col2"}},
				}
			},
			`ALTER ICEBERG TABLE %s DROP COLUMN IF EXISTS "col1", "col2"`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_RenameColumnAction,
			func(opts *AlterIcebergTableOptions) {
				opts.RenameColumnAction = &TableRenameColumnAction{OldName: "old_col", NewName: "new_col"}
			},
			`ALTER ICEBERG TABLE %s RENAME COLUMN "old_col" TO "new_col"`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_AlterColumnAction,
			func(opts *AlterIcebergTableOptions) {
				opts.AlterColumnAction = []IcebergTableAlterColumnAction{
					{ColumnName: "col1", SetNotNull: new(true)},
					{ColumnName: "col2", SetNotNull: new(true)},
				}
			},
			`ALTER ICEBERG TABLE %s ALTER COLUMN "col1" SET NOT NULL, COLUMN "col2" SET NOT NULL`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_alterColumn_dropNotNull",
			func(opts *AlterIcebergTableOptions) {
				opts.AlterColumnAction = []IcebergTableAlterColumnAction{
					{ColumnName: "col1", DropNotNull: new(true)},
					{ColumnName: "col2", DropNotNull: new(true)},
				}
			},
			`ALTER ICEBERG TABLE %s ALTER COLUMN "col1" DROP NOT NULL, COLUMN "col2" DROP NOT NULL`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_alterColumn_setDataType",
			func(opts *AlterIcebergTableOptions) {
				opts.AlterColumnAction = []IcebergTableAlterColumnAction{
					{ColumnName: "col1", DataType: &dataTypeVarchar_100},
					{ColumnName: "col2", DataType: &dataTypeNumber_36_2},
				}
			},
			`ALTER ICEBERG TABLE %s ALTER COLUMN "col1" SET DATA TYPE VARCHAR(100), COLUMN "col2" SET DATA TYPE NUMBER(36, 2)`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_alterColumn_commentUnsetComment",
			func(opts *AlterIcebergTableOptions) {
				opts.AlterColumnAction = []IcebergTableAlterColumnAction{
					{ColumnName: "col1", Comment: new("comment1")},
					{ColumnName: "col2", UnsetComment: new(true)},
				}
			},
			`ALTER ICEBERG TABLE %s ALTER COLUMN "col1" COMMENT 'comment1', COLUMN "col2" UNSET COMMENT`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_alterColumn_setDropWriteDefault",
			func(opts *AlterIcebergTableOptions) {
				opts.AlterColumnAction = []IcebergTableAlterColumnAction{
					{ColumnName: "col1", SetWriteDefault: new("1")},
					{ColumnName: "col2", DropWriteDefault: new(true)},
				}
			},
			`ALTER ICEBERG TABLE %s ALTER COLUMN "col1" SET WRITE DEFAULT 1, COLUMN "col2" DROP WRITE DEFAULT`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_SetMaskingPolicyOnColumn,
			func(opts *AlterIcebergTableOptions) {
				opts.SetMaskingPolicyOnColumn = &TableSetColumnMaskingPolicy{
					Name:          "col1",
					MaskingPolicy: maskingPolicyId,
					Using:         []Column{{"col1"}, {"col2"}},
					Force:         new(true),
				}
			},
			`ALTER ICEBERG TABLE %s ALTER COLUMN "col1" SET MASKING POLICY %s USING ("col1", "col2") FORCE`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(), maskingPolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_UnsetMaskingPolicyOnColumn,
			func(opts *AlterIcebergTableOptions) {
				opts.UnsetMaskingPolicyOnColumn = &TableUnsetColumnMaskingPolicy{Name: "col1"}
			},
			`ALTER ICEBERG TABLE %s ALTER COLUMN "col1" UNSET MASKING POLICY`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_SetProjectionPolicyOnColumn,
			func(opts *AlterIcebergTableOptions) {
				opts.SetProjectionPolicyOnColumn = &TableSetColumnProjectionPolicy{
					Name:             "col1",
					ProjectionPolicy: projectionPolicyId,
					Force:            new(true),
				}
			},
			`ALTER ICEBERG TABLE %s ALTER COLUMN "col1" SET PROJECTION POLICY %s FORCE`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(), projectionPolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_UnsetProjectionPolicyOnColumn,
			func(opts *AlterIcebergTableOptions) {
				opts.UnsetProjectionPolicyOnColumn = &TableUnsetColumnProjectionPolicy{Name: "col1"}
			},
			`ALTER ICEBERG TABLE %s ALTER COLUMN "col1" UNSET PROJECTION POLICY`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_SetTagsOnColumn,
			func(opts *AlterIcebergTableOptions) {
				opts.SetTagsOnColumn = &TableSetColumnTags{
					Name: "col1",
					SetTags: []TagAssociation{
						{Name: tagId1, Value: "v1"},
						{Name: tagId2, Value: "v2"},
					},
				}
			},
			`ALTER ICEBERG TABLE %s ALTER COLUMN "col1" SET TAG %s = 'v1', %s = 'v2'`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_UnsetTagsOnColumn,
			func(opts *AlterIcebergTableOptions) {
				opts.UnsetTagsOnColumn = &TableUnsetColumnTags{
					Name:      "col1",
					UnsetTags: []ObjectIdentifier{tagId1, tagId2},
				}
			},
			`ALTER ICEBERG TABLE %s ALTER COLUMN "col1" UNSET TAG %s, %s`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_ClusteringAction,
			func(opts *AlterIcebergTableOptions) {
				opts.ClusteringAction = &IcebergTableClusteringAction{ClusterBy: []string{`"col1"`, `"col2"`}}
			},
			`ALTER ICEBERG TABLE %s CLUSTER BY ("col1", "col2")`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_clustering_suspendRecluster",
			func(opts *AlterIcebergTableOptions) {
				opts.ClusteringAction = &IcebergTableClusteringAction{ChangeReclusterState: &IcebergTableReclusterChangeState{State: new(ReclusterStateSuspend)}}
			},
			`ALTER ICEBERG TABLE %s SUSPEND RECLUSTER`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_clustering_resumeRecluster",
			func(opts *AlterIcebergTableOptions) {
				opts.ClusteringAction = &IcebergTableClusteringAction{ChangeReclusterState: &IcebergTableReclusterChangeState{State: new(ReclusterStateResume)}}
			},
			`ALTER ICEBERG TABLE %s RESUME RECLUSTER`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_clustering_dropClusteringKey",
			func(opts *AlterIcebergTableOptions) {
				opts.ClusteringAction = &IcebergTableClusteringAction{DropClusteringKey: new(true)}
			},
			`ALTER ICEBERG TABLE %s DROP CLUSTERING KEY`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_Set,
			func(opts *AlterIcebergTableOptions) {
				opts.IfExists = new(true)
				opts.Set = &IcebergTableSetProperties{
					ReplaceInvalidCharacters:   new(true),
					CatalogSync:                new("integration1"),
					DataRetentionTimeInDays:    new(7),
					MaxDataExtensionTimeInDays: new(14),
					AutoRefresh:                new(true),
					TargetFileSize:             new(IcebergTableTargetFileSize128mb),
					Contact: []TableContact{
						{Purpose: "SUPPORT", Contact: contactId},
						{Purpose: "ACCESS_APPROVAL", Contact: contactId},
					},
					LogEventLevel:            new(IcebergTableLogEventLevelError),
					ErrorLogging:             new(true),
					EnableDataCompaction:     new(true),
					EnableIcebergMergeOnRead: new(true),
					Comment:                  new("updated comment"),
				}
			},
			`ALTER ICEBERG TABLE IF EXISTS %s SET `+
				`REPLACE_INVALID_CHARACTERS = true `+
				`CATALOG_SYNC = 'integration1' `+
				`DATA_RETENTION_TIME_IN_DAYS = 7 `+
				`MAX_DATA_EXTENSION_TIME_IN_DAYS = 14 `+
				`AUTO_REFRESH = true `+
				`TARGET_FILE_SIZE = '128MB' `+
				`CONTACT (SUPPORT = %s, ACCESS_APPROVAL = %s) `+
				`LOG_EVENT_LEVEL = ERROR `+
				`ERROR_LOGGING = true `+
				`ENABLE_DATA_COMPACTION = true `+
				`ENABLE_ICEBERG_MERGE_ON_READ = true `+
				`COMMENT = 'updated comment'`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			contactId.FullyQualifiedName(), contactId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_Unset,
			func(opts *AlterIcebergTableOptions) {
				opts.IfExists = new(true)
				opts.Unset = &IcebergTableUnsetProperties{
					ReplaceInvalidCharacters:   new(true),
					CatalogSync:                new(true),
					DataRetentionTimeInDays:    new(true),
					MaxDataExtensionTimeInDays: new(true),
					TargetFileSize:             new(true),
					LogEventLevel:              new(true),
					ErrorLogging:               new(true),
					EnableDataCompaction:       new(true),
					EnableIcebergMergeOnRead:   new(true),
					Comment:                    new(true),
				}
			},
			`ALTER ICEBERG TABLE IF EXISTS %s UNSET `+
				`REPLACE_INVALID_CHARACTERS, `+
				`CATALOG_SYNC, `+
				`DATA_RETENTION_TIME_IN_DAYS, `+
				`MAX_DATA_EXTENSION_TIME_IN_DAYS, `+
				`TARGET_FILE_SIZE, `+
				`LOG_EVENT_LEVEL, `+
				`ERROR_LOGGING, `+
				`ENABLE_DATA_COMPACTION, `+
				`ENABLE_ICEBERG_MERGE_ON_READ, `+
				`COMMENT`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_SetTags,
			func(opts *AlterIcebergTableOptions) {
				opts.SetTags = []TagAssociation{
					{Name: tagId1, Value: "v1"},
					{Name: tagId2, Value: "v2"},
				}
			},
			`ALTER ICEBERG TABLE %s SET TAG %s = 'v1', %s = 'v2'`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_UnsetTags,
			func(opts *AlterIcebergTableOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId1, tagId2}
			},
			`ALTER ICEBERG TABLE %s UNSET TAG %s, %s`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_AddRowAccessPolicy,
			func(opts *AlterIcebergTableOptions) {
				opts.AddRowAccessPolicy = &ViewAddRowAccessPolicy{
					RowAccessPolicy: rowAccessPolicy1Id,
					On:              []Column{{"col1"}, {"col2"}},
				}
			},
			`ALTER ICEBERG TABLE %s ADD ROW ACCESS POLICY %s ON ("col1", "col2")`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(), rowAccessPolicy1Id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_DropRowAccessPolicy,
			func(opts *AlterIcebergTableOptions) {
				opts.DropRowAccessPolicy = &ViewDropRowAccessPolicy{RowAccessPolicy: rowAccessPolicy1Id}
			},
			`ALTER ICEBERG TABLE %s DROP ROW ACCESS POLICY %s`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(), rowAccessPolicy1Id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_DropAndAddRowAccessPolicy,
			func(opts *AlterIcebergTableOptions) {
				opts.DropAndAddRowAccessPolicy = &IcebergTableDropAndAddRowAccessPolicy{
					Drop: IcebergTableDropRowAccessPolicy{RowAccessPolicy: rowAccessPolicy1Id},
					Add: IcebergTableAddRowAccessPolicy{
						RowAccessPolicy: rowAccessPolicy2Id,
						On:              []Column{{"col1"}, {"col2"}},
					},
				}
			},
			`ALTER ICEBERG TABLE %s DROP ROW ACCESS POLICY %s, ADD ROW ACCESS POLICY %s ON ("col1", "col2")`,
			icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(), rowAccessPolicy1Id.FullyQualifiedName(), rowAccessPolicy2Id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_DropAllRowAccessPolicies,
			func(opts *AlterIcebergTableOptions) { opts.DropAllRowAccessPolicies = new(true) },
			`ALTER ICEBERG TABLE %s DROP ALL ROW ACCESS POLICIES`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_SetAggregationPolicy,
			func(opts *AlterIcebergTableOptions) {
				opts.SetAggregationPolicy = &ViewSetAggregationPolicy{
					AggregationPolicy: aggregationPolicyId,
					EntityKey:         []Column{{"col1"}, {"col2"}},
					Force:             new(true),
				}
			},
			`ALTER ICEBERG TABLE %s SET AGGREGATION POLICY %s ENTITY KEY ("col1", "col2") FORCE`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(), aggregationPolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_UnsetAggregationPolicy,
			func(opts *AlterIcebergTableOptions) { opts.UnsetAggregationPolicy = &ViewUnsetAggregationPolicy{} },
			`ALTER ICEBERG TABLE %s UNSET AGGREGATION POLICY`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_SetJoinPolicy,
			func(opts *AlterIcebergTableOptions) {
				opts.SetJoinPolicy = &TableSetJoinPolicy{JoinPolicy: joinPolicyId, Force: new(true)}
			},
			`ALTER ICEBERG TABLE %s SET JOIN POLICY %s FORCE`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(), joinPolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_UnsetJoinPolicy,
			func(opts *AlterIcebergTableOptions) { opts.UnsetJoinPolicy = &TableUnsetJoinPolicy{} },
			`ALTER ICEBERG TABLE %s UNSET JOIN POLICY`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Alter_SearchOptimizationAction,
			func(opts *AlterIcebergTableOptions) {
				opts.SearchOptimizationAction = &TableSearchOptimizationAction{
					Add: &TableAddSearchOptimization{
						On: []TableSearchMethodWithTarget{
							{
								Method: TableSearchMethodEquality,
								Args:   TableSearchMethodArgs{Targets: []string{"col1", "col2"}, Analyzer: new("DEFAULT_ANALYZER")},
							},
							{
								Method: TableSearchMethodSubstring,
								Args:   TableSearchMethodArgs{Targets: []string{"*"}},
							},
						},
					},
				}
			},
			`ALTER ICEBERG TABLE %s ADD SEARCH OPTIMIZATION ON EQUALITY (col1, col2, ANALYZER => 'DEFAULT_ANALYZER'), SUBSTRING (*)`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_searchOptimization_dropWithSearchMethods",
			func(opts *AlterIcebergTableOptions) {
				opts.SearchOptimizationAction = &TableSearchOptimizationAction{
					Drop: &TableDropSearchOptimization{
						On: []TableDropSearchOptimizationOn{
							{SearchMethodWithTarget: &TableSearchMethodWithTarget{Method: TableSearchMethodFullText, Args: TableSearchMethodArgs{Targets: []string{"col1", "col2"}}}},
							{SearchMethodWithTarget: &TableSearchMethodWithTarget{Method: TableSearchMethodSubstring, Args: TableSearchMethodArgs{Targets: []string{"*"}}}},
						},
					},
				}
			},
			`ALTER ICEBERG TABLE %s DROP SEARCH OPTIMIZATION ON FULL_TEXT (col1, col2), SUBSTRING (*)`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_searchOptimization_dropMixedForms",
			func(opts *AlterIcebergTableOptions) {
				opts.SearchOptimizationAction = &TableSearchOptimizationAction{
					Drop: &TableDropSearchOptimization{
						On: []TableDropSearchOptimizationOn{
							{SearchMethodWithTarget: &TableSearchMethodWithTarget{Method: TableSearchMethodEquality, Args: TableSearchMethodArgs{Targets: []string{"col1"}}}},
							{ColumnName: new("col2")},
							{ExpressionId: new("expr_123")},
						},
					},
				}
			},
			`ALTER ICEBERG TABLE %s DROP SEARCH OPTIMIZATION ON EQUALITY (col1), col2, expr_123`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	icebergTablesTests.Drop.
		withExpectedSqlf(
			case_IcebergTables_sql_Drop_basic,
			`DROP ICEBERG TABLE %s`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Drop_all,
			func(opts *DropIcebergTableOptions) {
				opts.IfExists = new(true)
				opts.Cascade = new(true)
			},
			`DROP ICEBERG TABLE IF EXISTS %s CASCADE`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Drop_restrict",
			func(opts *DropIcebergTableOptions) {
				opts.IfExists = new(true)
				opts.Restrict = new(true)
			},
			`DROP ICEBERG TABLE IF EXISTS %s RESTRICT`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	icebergTablesTests.Describe.
		withExpectedSqlf(
			case_IcebergTables_sql_Describe_basic,
			`DESCRIBE ICEBERG TABLE %s`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Describe_columns",
			func(opts *DescribeIcebergTableOptions) { opts.DescribeType = new(IcebergTableDescribeTypeColumns) },
			`DESCRIBE ICEBERG TABLE %s TYPE = COLUMNS`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Describe_stage",
			func(opts *DescribeIcebergTableOptions) { opts.DescribeType = new(IcebergTableDescribeTypeStage) },
			`DESCRIBE ICEBERG TABLE %s TYPE = STAGE`, icebergTablesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	icebergTablesTests.Show.
		withExpectedSqlf(
			case_IcebergTables_sql_Show_basic,
			`SHOW ICEBERG TABLES`,
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Show_all,
			func(opts *ShowIcebergTableOptions) {
				opts.Like = &Like{Pattern: new("some_pattern")}
				opts.In = &In{Schema: NewDatabaseObjectIdentifier("db", "schema")}
				opts.StartsWith = new("prefix")
				opts.Limit = &LimitFrom{Rows: new(10), From: new("table_name")}
			},
			`SHOW ICEBERG TABLES LIKE 'some_pattern' IN SCHEMA "db"."schema" STARTS WITH 'prefix' LIMIT 10 FROM 'table_name'`,
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Show_Like,
			func(opts *ShowIcebergTableOptions) { opts.Like = &Like{Pattern: new("pattern_test")} },
			`SHOW ICEBERG TABLES LIKE 'pattern_test'`,
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Show_In,
			func(opts *ShowIcebergTableOptions) { opts.In = &In{Database: NewAccountObjectIdentifier("test_db")} },
			`SHOW ICEBERG TABLES IN DATABASE "test_db"`,
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Show_StartsWith,
			func(opts *ShowIcebergTableOptions) { opts.StartsWith = new("prefix") },
			`SHOW ICEBERG TABLES STARTS WITH 'prefix'`,
		).
		withModifyAndExpectedSqlf(
			case_IcebergTables_sql_Show_Limit,
			func(opts *ShowIcebergTableOptions) { opts.Limit = &LimitFrom{Rows: new(5)} },
			`SHOW ICEBERG TABLES LIMIT 5`,
		).
		withAdditionalSqlCasef(
			"sql_Show_InAccount",
			func(opts *ShowIcebergTableOptions) { opts.In = &In{Account: new(true)} },
			`SHOW ICEBERG TABLES IN ACCOUNT`,
		)
}
