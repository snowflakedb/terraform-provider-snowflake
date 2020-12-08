package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pkg/errors"
)

const (
	taskIDDelimiter = '|'
)

var taskSchema = map[string]*schema.Schema{
	"enabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies if the task should be started (enabled) after creation or should remain suspended (default).",
	},
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the task; must be unique for the database and schema in which the task is created.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the task.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the task.",
		ForceNew:    true,
	},
	"warehouse": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The warehouse the task will use.",
		ForceNew:    false,
	},
	"schedule": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The schedule for periodically running the task. This can be a cron or interval in minutes.",
	},
	"session_parameters": {
		Type:        schema.TypeMap,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies session parameters to set for the session when the task runs. A task supports all session parameters.",
	},
	"user_task_timeout_ms": {
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 86400000),
		Description:  "Specifies the time limit on a single run of the task before it times out (in milliseconds).",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the task.",
	},
	"after": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the predecessor task in the same database and schema of the current task. When a run of the predecessor task finishes successfully, it triggers this task (after a brief lag).",
	},
	"when": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a Boolean SQL expression; multiple conditions joined with AND/OR are supported.",
	},
	"sql_statement": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Any single SQL statement, or a call to a stored procedure, executed when the task runs.",
		ForceNew:    false,
	},
}

type taskID struct {
	DatabaseName string
	SchemaName   string
	TaskName     string
}

//String() takes in a taskID object and returns a pipe-delimited string:
//DatabaseName|SchemaName|TaskName
func (t *taskID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = taskIDDelimiter
	dataIdentifiers := [][]string{{t.DatabaseName, t.SchemaName, t.TaskName}}
	err := csvWriter.WriteAll(dataIdentifiers)
	if err != nil {
		return "", err
	}
	strTaskID := strings.TrimSpace(buf.String())
	return strTaskID, nil
}

// difference find keys in a but not in b
func difference(a, b map[string]interface{}) map[string]interface{} {
	diff := make(map[string]interface{})
	for k := range a {
		if _, ok := b[k]; !ok {
			diff[k] = a[k]
		}
	}
	return diff
}

// getActiveRootTask tries to retrieve the root of current task or returns the current (standalone) task
func getActiveRootTask(data *schema.ResourceData, meta interface{}) (*snowflake.TaskBuilder, error) {
	log.Println("[DEBUG] retrieving root task")

	db := meta.(*sql.DB)
	database := data.Get("database").(string)
	dbSchema := data.Get("schema").(string)
	name := data.Get("name").(string)
	after := data.Get("after").(string)

	if name == "" {
		return nil, nil
	}

	// always start from first predecessor
	// or the current task when standalone
	if after != "" {
		name = after
	}

	for {
		builder := snowflake.Task(name, database, dbSchema)
		q := builder.Show()
		row := snowflake.QueryRow(db, q)
		task, err := snowflake.ScanTask(row)

		if err != nil && name != data.Get("name").(string) {
			return nil, errors.Wrapf(err, "failed to locate the root node of: %v", name)
		}

		if task.Predecessors == nil {
			log.Println(fmt.Sprintf("[DEBUG] found root task: %v", name))
			// we only want to deal with suspending the root task when its enabled (started)
			if task.IsEnabled() {
				return snowflake.Task(name, database, dbSchema), nil
			}
			return nil, nil
		}

		name = task.GetPredecessorName()
	}
}

// getActiveRootTaskAndSuspend retrieves the root task and suspends it
func getActiveRootTaskAndSuspend(data *schema.ResourceData, meta interface{}) (*snowflake.TaskBuilder, error) {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)

	root, err := getActiveRootTask(data, meta)
	if err != nil {
		return nil, errors.Wrapf(err, "error retrieving root task %v", name)
	}

	if root != nil {
		qr := root.Suspend()
		err = snowflake.Exec(db, qr)
		if err != nil {
			return nil, errors.Wrapf(err, "error suspending root task %v", name)
		}
	}

	return root, nil
}

func resumeTask(root *snowflake.TaskBuilder, meta interface{}) {
	if root == nil {
		return
	}

	if root.IsDisabled() {
		return
	}

	db := meta.(*sql.DB)
	qr := root.Resume()
	err := snowflake.Exec(db, qr)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "error resuming root task %v", root.QualifiedName()))
	}
}

// taskIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|TaskName
// and returns a taskID object
func taskIDFromString(stringID string) (*taskID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = pipeIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per task")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("3 fields allowed")
	}

	taskResult := &taskID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		TaskName:     lines[0][2],
	}
	return taskResult, nil
}

// Task returns a pointer to the resource representing a task
func Task() *schema.Resource {
	return &schema.Resource{
		Create: CreateTask,
		Read:   ReadTask,
		Update: UpdateTask,
		Delete: DeleteTask,

		Schema: taskSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// ReadTask implements schema.ReadFunc
func ReadTask(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	taskID, err := taskIDFromString(d.Id())
	if err != nil {
		return err
	}

	database := taskID.DatabaseName
	schema := taskID.SchemaName
	name := taskID.TaskName

	builder := snowflake.Task(name, database, schema)
	q := builder.Show()
	row := snowflake.QueryRow(db, q)
	t, err := snowflake.ScanTask(row)
	if err == sql.ErrNoRows {
		log.Printf("[WARN] task (%s) not found, removing from state file", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	err = d.Set("enabled", t.IsEnabled())
	if err != nil {
		return err
	}

	err = d.Set("name", t.Name)
	if err != nil {
		return err
	}

	err = d.Set("database", t.DatabaseName)
	if err != nil {
		return err
	}

	err = d.Set("schema", t.SchemaName)
	if err != nil {
		return err
	}

	err = d.Set("warehouse", t.Warehouse)
	if err != nil {
		return err
	}

	err = d.Set("schedule", t.Schedule)
	if err != nil {
		return err
	}

	err = d.Set("comment", t.Comment)
	if err != nil {
		return err
	}

	if t.Predecessors != nil {
		err = d.Set("after", t.GetPredecessorName())
		if err != nil {
			return err
		}
	}

	err = d.Set("when", t.Condition)
	if err != nil {
		return err
	}

	err = d.Set("sql_statement", t.Definition)
	if err != nil {
		return err
	}

	q = builder.ShowParameters()
	paramRows, err := snowflake.Query(db, q)
	if err != nil {
		return err
	}
	params, err := snowflake.ScanTaskParameters(paramRows)
	if err != nil {
		return err
	}

	if len(params) > 0 {
		paramMap := map[string]interface{}{}
		for _, param := range params {
			log.Printf("[TRACE] %+v\n", param)
			if param.Value == param.DefaultValue {
				continue
			}

			paramMap[param.Key] = param.Value
		}

		err := d.Set("session_parameters", paramMap)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateTask implements schema.CreateFunc
func CreateTask(d *schema.ResourceData, meta interface{}) error {

	var err error
	db := meta.(*sql.DB)
	database := d.Get("database").(string)
	dbSchema := d.Get("schema").(string)
	name := d.Get("name").(string)
	warehouse := d.Get("warehouse").(string)
	sql := d.Get("sql_statement").(string)
	enabled := d.Get("enabled").(bool)

	builder := snowflake.Task(name, database, dbSchema)
	builder.WithWarehouse(warehouse)
	builder.WithStatement(sql)

	// Set optionals
	if v, ok := d.GetOk("schedule"); ok {
		builder.WithSchedule(v.(string))
	}

	if v, ok := d.GetOk("session_parameters"); ok {
		builder.WithSessionParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("user_task_timeout_ms"); ok {
		builder.WithTimeout(v.(int))
	}

	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := d.GetOk("after"); ok {
		root, err := getActiveRootTaskAndSuspend(d, meta)
		if err != nil {
			return err
		}
		defer resumeTask(root, meta)

		builder.WithDependency(v.(string))
	}

	if v, ok := d.GetOk("when"); ok {
		builder.WithCondition(v.(string))
	}

	q := builder.Create()
	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error creating task %v", name)
	}

	if enabled {
		q = builder.Resume()
		err = snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error starting task %v", name)
		}
	}

	taskID := &taskID{
		DatabaseName: database,
		SchemaName:   dbSchema,
		TaskName:     name,
	}
	dataIDInput, err := taskID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadTask(d, meta)
}

// UpdateTask implements schema.UpdateFunc
func UpdateTask(d *schema.ResourceData, meta interface{}) error {
	taskID, err := taskIDFromString(d.Id())
	if err != nil {
		return err
	}

	db := meta.(*sql.DB)
	database := taskID.DatabaseName
	dbSchema := taskID.SchemaName
	name := taskID.TaskName
	builder := snowflake.Task(name, database, dbSchema)

	root, err := getActiveRootTaskAndSuspend(d, meta)
	if err != nil {
		return err
	}
	defer resumeTask(root, meta)

	if d.HasChange("warehouse") {
		new := d.Get("warehouse")
		q := builder.ChangeWarehouse(new.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating warehouse on task %v", d.Id())
		}
	}

	if d.HasChange("schedule") {
		var q string
		old, new := d.GetChange("schedule")
		if old != "" && new == "" {
			q = builder.RemoveSchedule()
		} else {
			q = builder.ChangeSchedule(new.(string))
		}
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating schedule on task %v", d.Id())
		}
	}

	if d.HasChange("user_task_timeout_ms") {
		var q string
		old, new := d.GetChange("user_task_timeout_ms")
		if old.(int) > 0 && new.(int) == 0 {
			q = builder.RemoveTimeout()
		} else {
			q = builder.ChangeTimeout(new.(int))
		}
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating user task timeout on task %v", d.Id())
		}
	}

	if d.HasChange("comment") {
		var q string
		old, new := d.GetChange("comment")
		if old != "" && new == "" {
			q = builder.RemoveComment()
		} else {
			q = builder.ChangeComment(new.(string))
		}
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating comment on task %v", d.Id())
		}
	}

	if d.HasChange("after") {
		var (
			q   string
			err error
		)

		old, new := d.GetChange("after")
		enabled := d.Get("enabled").(bool)

		if enabled {
			q = builder.Suspend()
			err = snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error suspending task %v", d.Id())
			}
			defer resumeTask(builder, meta)
		}

		if old != "" {
			q = builder.RemoveDependency(old.(string))
			err = snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error removing old after dependency from task %v", d.Id())
			}
		}

		if new != "" {
			q = builder.AddDependency(new.(string))
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error adding after dependency on task %v", d.Id())
			}
		}
	}

	if d.HasChange("session_parameters") {
		var q string
		o, n := d.GetChange("session_parameters")

		if o == nil {
			o = make(map[string]interface{})
		}
		if n == nil {
			n = make(map[string]interface{})
		}
		os := o.(map[string]interface{})
		ns := n.(map[string]interface{})

		remove := difference(os, ns)
		add := difference(ns, os)

		if len(remove) > 0 {
			q = builder.RemoveSessionParameters(remove)
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error removing session_parameters on task %v", d.Id())
			}
		}

		if len(add) > 0 {
			q = builder.AddSessionParameters(add)
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error adding session_parameters to task %v", d.Id())
			}
		}
	}

	if d.HasChange("when") {
		new := d.Get("when")
		q := builder.ChangeCondition(new.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating when condition on task %v", d.Id())
		}
	}

	if d.HasChange("sql_statement") {
		new := d.Get("sql_statement")
		q := builder.ChangeSqlStatement(new.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating sql statement on task %v", d.Id())
		}
	}

	if d.HasChange("enabled") {
		var q string
		n := d.Get("enabled")
		enable := n.(bool)

		if enable {
			q = builder.Resume()
		} else {
			q = builder.Suspend()
			// make sure defer doesn't enable task again
			// when standalone or root task and status is supsended
			if root != nil && builder.QualifiedName() == root.QualifiedName() {
				root = root.SetDisabled()
			}
		}

		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating task state %v", d.Id())
		}
	}

	return ReadTask(d, meta)
}

// DeleteTask implements schema.DeleteFunc
func DeleteTask(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	taskID, err := taskIDFromString(d.Id())
	if err != nil {
		return err
	}

	database := taskID.DatabaseName
	schema := taskID.SchemaName
	name := taskID.TaskName

	root, err := getActiveRootTaskAndSuspend(d, meta)
	if err != nil {
		return err
	}

	// only resume the root when not a standalone task
	if root != nil && name != root.Name() {
		defer resumeTask(root, meta)
	}

	q := snowflake.Task(name, database, schema).Drop()
	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting task %v", d.Id())
	}

	d.SetId("")

	return nil
}

// TaskExists implements schema.ExistsFunc
func TaskExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	taskID, err := taskIDFromString(data.Id())
	if err != nil {
		return false, err
	}

	database := taskID.DatabaseName
	schema := taskID.SchemaName
	name := taskID.TaskName

	q := snowflake.Task(name, database, schema).Show()
	rows, err := db.Query(q)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, nil
}
