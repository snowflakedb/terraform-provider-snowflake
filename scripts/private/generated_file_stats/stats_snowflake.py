"""Count generated vs hand-written files and lines in a Snowflake Git repository (stored procedure handler)."""

from __future__ import annotations

import posixpath
import re
from typing import Any

from snowflake.snowpark.functions import current_timestamp

from stats_common import Report, ScanConfig, run_scan


class SnowflakeBackend:
    def __init__(
        self,
        session: Any,
        git_repository: str,
        branch: str,
        target_table: str,
        line_file_format: str,
        repository_label: str | None,
    ) -> None:
        self.session = session
        self.git_repository = git_repository
        self.target_table = target_table
        self.line_file_format = line_file_format
        self._branch = branch
        self._repository = repository_label or git_repository

    @property
    def repository(self) -> str:
        return self._repository

    @property
    def branch(self) -> str:
        return self._branch

    def get_revision(self) -> str:
        for row in self.session.sql(f"SHOW GIT BRANCHES IN {self.git_repository}").collect():
            if _row_field(row, "name") == self.branch:
                return _row_field(row, "commit_hash")
        raise ValueError(f"Could not find commit hash for branch {self.branch!r}")

    def list_paths(self) -> list[str]:
        rows = self.session.sql(f"LIST {self._branch_location()}").collect()
        return [self._repository_path(_row_field(row, "name")) for row in rows]

    def line_counts(self, paths: list[str]) -> dict[str, int]:
        if not paths:
            return {}

        # 1. Read only the directory shared by all paths (e.g. pkg/sdk/), not the whole repository.
        directory = posixpath.commonpath([posixpath.dirname(path) for path in paths])
        location = self._branch_location() + posixpath.join(directory, "")

        # 2. Count the lines of every file there; the file format turns each line into one row.
        rows = self.session.sql(
            f"SELECT METADATA$FILENAME AS path, COUNT(*) AS lines"
            f" FROM {location} (FILE_FORMAT => '{self.line_file_format}')"
            f" GROUP BY 1"
        ).collect()

        # 3. Keep the requested paths; the directory can contain more files.
        requested = set(paths)
        counts: dict[str, int] = {}
        for row in rows:
            path = self._repository_path(_row_field(row, "path"))
            if path in requested:
                counts[path] = int(_row_field(row, "lines"))

        if rows and not counts:  # better fail than report zeros when the file names change their form
            raise ValueError(f"No requested path in the line counts, e.g. {_row_field(rows[0], 'path')!r}")
        return counts

    def save_report(self, report: Report) -> None:
        row = self.session.create_dataframe([report.to_row()]).with_column("captured_at", current_timestamp())
        row.write.save_as_table(self.target_table, mode="append", column_order="name")

    def _branch_location(self) -> str:
        # The branch is quoted, because branch names can contain "/".
        return f'@{self.git_repository}/branches/"{self.branch}"/'

    def _repository_path(self, name: str) -> str:
        # LIST returns <repository>/branches/"<branch>"/<path>, METADATA$FILENAME git_repository/<id>/blobs/<hash>/<path>.
        for marker in (f'branches/"{self.branch}"/', f"branches/{self.branch}/"):
            if marker in name:
                return name.split(marker, 1)[1]
        match = re.match(r"git_repository/[^/]+/blobs/[0-9a-f]+/", name)
        if match:
            return name[match.end():]
        return name


def snowflake_handler(session: Any, branch: str, config: dict[str, Any]) -> dict[str, Any]:
    # Expected CALL arguments (the regexes are the same as in Makefile.private.mk):
    #   branch: the branch to scan, e.g. 'dev'
    #   config: OBJECT_CONSTRUCT(
    #       'git_repository',     '<database>.<schema>.<git repository>',
    #       'target_table',       '<database>.<schema>.<table>',
    #       'line_file_format',   '<database>.<schema>.<file format>',
    #       'scope',              'sdk',                                   -- optional, default 'all'
    #       'repository_label',   'terraform-provider-snowflake',          -- optional, default git_repository
    #       'scan_patterns',      ARRAY_CONSTRUCT('^pkg/sdk/[^/]+\\.go$'),
    #       'generated_patterns', ARRAY_CONSTRUCT('_gen(_test)?\\.go$'),
    #       'exclude_patterns',   ARRAY_CONSTRUCT('(^|/)testdata/'),     -- optional
    #       'show_files',         FALSE                                    -- optional
    #   )
    scan_config = ScanConfig.from_strings(
        scan=config["scan_patterns"],
        generated=config["generated_patterns"],
        exclude=config.get("exclude_patterns", []),
    )
    backend = SnowflakeBackend(
        session=session,
        git_repository=config["git_repository"],
        branch=branch,
        target_table=config["target_table"],
        line_file_format=config["line_file_format"],
        repository_label=config.get("repository_label"),
    )
    report = run_scan(scan_config, backend, config.get("scope", "all"), config.get("show_files", False))
    return report.to_dict()


def _row_field(row: Any, name: str) -> str:
    # Snowpark Row keys are upper or lower case, depending on the command.
    values = {str(key).lower(): value for key, value in row.as_dict().items()}
    return str(values[name.lower()])
