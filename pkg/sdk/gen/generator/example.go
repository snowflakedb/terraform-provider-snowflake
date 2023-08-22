package generator

var DatabaseRoleInterface = Interface{
	Name:         "DatabaseRoles",
	nameSingular: "DatabaseRole",
	Operations: []*Operation{
		{
			Name:            "Create",
			ObjectInterface: nil,
			Doc:             "https://docs.snowflake.com/en/sql-reference/sql/create-database-role",
			OptsStructFields: []*Field{
				{
					Name: "create",
					Kind: "bool",
					tags: map[string][]string{
						"ddl": {"static"},
						"sql": {"CREATE"},
					},
				},
				{
					Name: "OrReplace",
					Kind: "*bool",
					tags: map[string][]string{
						"ddl": {"keyword"},
						"sql": {"OR REPLACE"},
					},
				},
				{
					Name: "databaseRole",
					Kind: "bool",
					tags: map[string][]string{
						"ddl": {"static"},
						"sql": {"DATABASE ROLE"},
					},
				},
				{
					Name: "IfNotExists",
					Kind: "*bool",
					tags: map[string][]string{
						"ddl": {"keyword"},
						"sql": {"IF NOT EXISTS"},
					},
				},
				{
					Name: "name",
					Kind: "*bool",
					tags: map[string][]string{
						"ddl": {"identifier"},
					},
				},
				{
					Name: "Comment",
					Kind: "*string",
					tags: map[string][]string{
						"ddl": {"parameter", "single_quotes"},
						"sql": {"COMMENT"},
					},
				},
			},
			Validations: []*Validation{
				{
					Type:       ValidIdentifier,
					fieldNames: []string{"name"},
				},
				{
					Type:       ConflictingFields,
					fieldNames: []string{"OrReplace", "IfNotExists"},
				},
			},
		},
		{
			Name:            "Alter",
			ObjectInterface: nil,
			Doc:             "https://docs.snowflake.com/en/sql-reference/sql/alter-database-role",
			OptsStructFields: []*Field{
				{
					Name: "alter",
					Kind: "bool",
					tags: map[string][]string{
						"ddl": {"static"},
						"sql": {"ALTER"},
					},
				},
				{
					Name: "databaseRole",
					Kind: "bool",
					tags: map[string][]string{
						"ddl": {"static"},
						"sql": {"DATABASE ROLE"},
					},
				},
				{
					Name: "IfExists",
					Kind: "*bool",
					tags: map[string][]string{
						"ddl": {"keyword"},
						"sql": {"IF EXISTS"},
					},
				},
				{
					Name: "name",
					Kind: "DatabaseObjectIdentifier",
					tags: map[string][]string{
						"ddl": {"identifier"},
					},
				},
				{
					Name: "Rename",
					Kind: "*DatabaseRoleRename",
					tags: map[string][]string{
						"ddl": {"list,no_parentheses"},
						"sql": {"RENAME TO"},
					},
					Fields: []*Field{
						{
							Name: "Name",
							Kind: "DatabaseObjectIdentifier",
							tags: map[string][]string{
								"ddl": {"identifier"},
							},
						},
					},
					Validations: []*Validation{
						{
							Type:       ValidIdentifier,
							fieldNames: []string{"Name"},
						},
					},
				},
				{
					Name: "Set",
					Kind: "*DatabaseRoleSet",
					tags: map[string][]string{
						"ddl": {"list,no_parentheses"},
						"sql": {"SET"},
					},
					Fields: []*Field{
						{
							Name: "Comment",
							Kind: "string",
							tags: map[string][]string{
								"ddl": {"parameter", "single_quotes"},
								"sql": {"COMMENT"},
							},
						},
					},
				},
				{
					Name: "Unset",
					Kind: "*DatabaseRoleUnset",
					tags: map[string][]string{
						"ddl": {"list,no_parentheses"},
						"sql": {"UNSET"},
					},
					Fields: []*Field{
						{
							Name: "Comment",
							Kind: "bool",
							tags: map[string][]string{
								"ddl": {"keyword"},
								"sql": {"COMMENT"},
							},
						},
					},
					Validations: []*Validation{
						{
							Type:       AtLeastOneValueSet,
							fieldNames: []string{"Comment"},
						},
					},
				},
			},
			Validations: []*Validation{
				{
					Type:       ValidIdentifier,
					fieldNames: []string{"name"},
				},
				{
					Type:       ExactlyOneValueSet,
					fieldNames: []string{"Rename", "Set", "Unset"},
				},
			},
		},
	},
}
