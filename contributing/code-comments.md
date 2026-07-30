## Code comments

### Guidelines for In-Code Comments

The guiding principle is simple: **write a comment only when the code itself cannot fully express what is going on.** A well-named function, a clear variable, and a readable test already communicate the *what*.
Comments exist to convey the *why* — hidden constraints, non-obvious behavior, deliberate trade-offs, and workarounds for specific bugs.

---

### Do not comment the obvious

If a reader can understand the intent from the identifier name and the surrounding code, the comment adds noise rather than signal.

❌ **Restating the function name**
```go
// WithComment sets the comment field.
func (m *WarehouseModel) WithComment(comment string) *WarehouseModel { ... }
```

❌ **Restating the type name**
```go
// DefaultValueConfig represents the configuration for default value behavior.
type DefaultValueConfig struct { ... }
```

❌ **Describing parameters when the names speak for themselves**
```go
// sourceId is the source object, targetId is the target object, roles are the roles to grant.
func GrantOnObject(sourceId AccountObjectIdentifier, targetId AccountObjectIdentifier, roles []string) { ... }
```

---

### Do comment the non-obvious

Add a comment when a future reader would otherwise be surprised, confused, or tempted to "simplify" the code in a way that breaks it.

✅ **Labeling acceptance test steps**
```go
// Create - without optionals
{
    Config: basicModel(...).HCL(),
    ...
},
// Update - set optionals
{
    Config: completeModel(...).HCL(),
    ...
},
```
Acceptance tests contain multiple sequential steps that are not individually named. A short label on each step makes the intent scannable without opening the config or counting indices.

✅ **A Snowflake API constraint that is not reflected in the schema**
```go
// These fields cannot be set at CREATE time; they must be applied via ALTER after the object reaches READY state.
```

✅ **A known fragility or workaround for a library limitation**
```go
// Uses ObjectVariable rather than MapVariable because MapVariable.MarshalJSON
// requires all values to share the same underlying type; this block mixes string, bool, and list.
```

✅ **Surprising Snowflake behavior that affects test assertions**
```go
// Not asserted: DESCRIBE has a propagation lag for this field and consistently
// returns an empty value immediately after CREATE/ALTER; the assertion would be flaky.
```

✅ **A fragile match that could misfire**
```go
// Note: this is a substring match against a Snowflake error message. The message
// may change in future Snowflake versions or match an unrelated error accidentally.
// Treat any retry triggered here as best-effort.
```

✅ **A mutually exclusive constraint that the schema cannot enforce**
```go
// The framework does not support ConflictsWith for nested attributes.
// Mutual exclusivity between these fields is validated manually here.
```

✅ **A TODO with a tracking reference**
```go
// TODO(SNOW-XXXXXXX): this field cannot be read back via DESCRIBE; tracked for a follow-up.
```

---

### Comments in generated files belong in the def file

Files marked `// Code generated … DO NOT EDIT.` are overwritten on every `make generate-sdk` run. Any comment you write there will be lost. Put explanatory notes in the object's `*_def.go` file instead.

❌ **Comment added directly to a generated file**
```go
// pkg/sdk/application_roles_gen.go  (DO NOT EDIT)
// Only query operations are allowed as other operations are not accessible from program context.
type ApplicationRoles interface { ... }
```

✅ **Comment placed in the def file** (real example from [`application_roles_def.go`](../pkg/sdk/generator/defs/application_roles_def.go))
```go
// applicationRolesDef creates an interface that allows for querying application roles.
// It does not allow for other DDL queries (CREATE, ALTER, DROP, ...) to be called, because
// they are not possible to be called from the program level. Application roles are a special
// case where they're only usable inside application context (e.g. setup.sql). Right now, they
// can only be manipulated from the program context by applying debug_mode to the application,
// but it's a hacky solution and even then GRANT and REVOKE options are limited.
// That's why we're only exposing SHOW operations — the only ones allowed from the program context.
var applicationRolesDef = g.NewInterface(...)
```
