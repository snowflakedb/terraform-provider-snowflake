# Private scripts

Excluded from the public mirror (`scripts/private/**` in `scripts/mirroring/copybara/copy.bara.sky`).

## generated_file_stats

Counts generated vs hand-written Go files and lines among the Git-tracked files. One run reports one scope; the Makefile
target runs it once per package our generators write into (`pkg/sdk`, `pkg/schemas`, `pkg/resources`, and the generated
assertion and config model packages in `pkg/acceptance/bettertestspoc`), plus a `total` row for all of them together.
Standard library only (Python 3.11+), in three files:

| File                 | Contents                                                                  |
|----------------------|---------------------------------------------------------------------------|
| `stats_common.py`    | Rules (`ScanConfig`), the report, and the scan algorithm (`run_scan`).    |
| `stats_local.py`     | Local entry point: `LocalBackend` and the CLI.                            |
| `stats_snowflake.py` | Stored procedure entry point: `SnowflakeBackend` and `snowflake_handler`. |

The modules import each other by plain name (`from stats_common import ...`); the `stats_` prefix keeps them from
shadowing other modules (a `snowflake.py` would hide the `snowflake.snowpark` package in Snowflake).

```bash
make -f Makefile.private.mk generated-file-stats                     # the usual report, current checkout
make -f Makefile.private.mk generated-file-stats ARGS="--format json"
python3 scripts/private/generated_file_stats/stats_local.py --help   # a single run
```

### Options

Regex values are matched with `re.search` against repo-relative paths with `/` separators. Flags that are not passed fall
back to the defaults below (defined in `parse_cli_args`); the Makefile target passes its own values.

| Option         | Default             | Accepted values                      | Description                                                                              |
|----------------|---------------------|--------------------------------------|------------------------------------------------------------------------------------------|
| `--root`       | `.`                 | path to a Git checkout               | Repository to scan; only Git-tracked files are counted.                                  |
| `--scope`      | `all`               | any label                            | Name shown in the `SCOPE` column (and stored in the table); it does not filter anything. |
| `--scan`       | `\.go$`             | regex, repeatable (any match counts) | Files that are counted at all. Limit it to scope a run to part of the tree.              |
| `--generated`  | `_gen(_test)?\.go$` | regex, repeatable                    | Counted files that are generated; the rest (including `*_ext.go`) are hand-written.      |
| `--exclude`    | `(^\|/)testdata/`   | regex, repeatable                    | Files in scan scope that are skipped; only reported in the excluded columns.             |
| `--format`     | `plain`             | `plain`, `json`                      | Table for people, or JSON (with the file list when `--show-files` is set) for tools.     |
| `--no-header`  | off                 | flag                                 | Plain format only: print just the row, to stack several runs under one header.           |
| `--show-files` | off                 | flag                                 | List every counted file with its classification and lines (to check the regexes).        |
| `-h`, `--help` |                     | flag                                 | Print the built-in help.                                                                 |

Scopes are flat directories (`^pkg/sdk/[^/]+\.go$`), so a generator's own code in its `gen/` subdirectory is not counted
as hand-written code of the package. When a generator starts writing into a new package, add a line for it to the
`generated-file-stats` target and its `--scan` to the `total` line.
`ruleset_version` in the report is a hash of the three regex lists, so reports made with different rules can be told apart.

### Running in Snowflake

For more information head over to the team's Snowsight workspace in Snowhouse.

### Extending

Both entry points share `run_scan`; a new platform only needs a module with a class implementing the `ScanBackend`
methods (see `stats_common.py`).
A new category (e.g. split tests from production code) is one more `Classification` value, a branch in `classify`, and
one more `Totals` in `Report` and `run_scan`.
