"""Shared part of the generated file stats: rules, report, and the scan algorithm."""

from __future__ import annotations

import hashlib
import json
import re
from dataclasses import asdict, dataclass
from enum import StrEnum
from typing import Any, Protocol


class Classification(StrEnum):
    GENERATED = "generated"
    HANDWRITTEN = "handwritten"
    EXCLUDED = "excluded"
    IGNORED = "ignored"  # outside of the scan patterns; not counted at all


@dataclass(frozen=True)
class ScanConfig:
    scan_patterns: list[re.Pattern[str]]
    generated_patterns: list[re.Pattern[str]]
    exclude_patterns: list[re.Pattern[str]]

    @classmethod
    def from_strings(cls, scan: list[str], generated: list[str], exclude: list[str]) -> ScanConfig:
        return cls(
            scan_patterns=[re.compile(pattern) for pattern in scan],
            generated_patterns=[re.compile(pattern) for pattern in generated],
            exclude_patterns=[re.compile(pattern) for pattern in exclude],
        )

    @property
    def ruleset_version(self) -> str:
        # A short hash of the rules, so reports made with different rules can be told apart.
        rules = [
            [pattern.pattern for pattern in self.scan_patterns],
            [pattern.pattern for pattern in self.generated_patterns],
            [pattern.pattern for pattern in self.exclude_patterns],
        ]
        return hashlib.sha256(json.dumps(rules).encode()).hexdigest()[:8]


class ScanBackend(Protocol):
    @property
    def repository(self) -> str: ...

    @property
    def branch(self) -> str: ...

    def get_revision(self) -> str: ...

    def list_paths(self) -> list[str]: ...

    def line_counts(self, paths: list[str]) -> dict[str, int]: ...

    def save_report(self, report: Report) -> None: ...


@dataclass
class Totals:
    files: int = 0
    lines: int = 0


@dataclass
class FileResult:
    path: str
    classification: Classification
    lines: int


@dataclass
class Report:
    scope: str
    repository: str
    branch: str
    commit_sha: str
    ruleset_version: str
    generated: Totals
    handwritten: Totals
    excluded: Totals
    files: list[FileResult]

    @property
    def total_files(self) -> int:
        return self.generated.files + self.handwritten.files

    @property
    def total_lines(self) -> int:
        return self.generated.lines + self.handwritten.lines

    @property
    def generated_percent(self) -> float:
        return _percent(self.generated.files, self.total_files)

    @property
    def generated_line_percent(self) -> float:
        return _percent(self.generated.lines, self.total_lines)

    def to_row(self) -> dict[str, Any]:
        return {
            "scope": self.scope,
            "repository": self.repository,
            "branch": self.branch,
            "commit_sha": self.commit_sha,
            "ruleset_version": self.ruleset_version,
            "total_files": self.total_files,
            "generated_files": self.generated.files,
            "non_generated_files": self.handwritten.files,
            "excluded_files": self.excluded.files,
            "total_lines": self.total_lines,
            "generated_lines": self.generated.lines,
            "non_generated_lines": self.handwritten.lines,
            "excluded_lines": self.excluded.lines,
            "generated_percent": self.generated_percent,
            "generated_line_percent": self.generated_line_percent,
        }

    def to_dict(self) -> dict[str, Any]:
        return {**self.to_row(), "files": [asdict(file) for file in self.files]}


def run_scan(config: ScanConfig, backend: ScanBackend, scope: str, show_files: bool = False) -> Report:
    # 1. Classify every path, keeping only the ones in scope.
    classifications: dict[str, Classification] = {}
    for path in sorted(backend.list_paths()):
        classification = classify(path, config)
        if classification != Classification.IGNORED:
            classifications[path] = classification

    # 2. Count the lines of all kept files in one call (one query in Snowflake).
    lines = backend.line_counts(list(classifications))
    files = [FileResult(path, classification, lines.get(path, 0)) for path, classification in classifications.items()]

    # 3. Sum the files and lines per classification.
    report = Report(
        scope=scope,
        repository=backend.repository,
        branch=backend.branch,
        commit_sha=backend.get_revision(),
        ruleset_version=config.ruleset_version,
        generated=Totals(),
        handwritten=Totals(),
        excluded=Totals(),
        files=files if show_files else [],
    )
    totals = {
        Classification.GENERATED: report.generated,
        Classification.HANDWRITTEN: report.handwritten,
        Classification.EXCLUDED: report.excluded,
    }
    for file in files:
        totals[file.classification].files += 1
        totals[file.classification].lines += file.lines

    # 4. Print or store the report.
    backend.save_report(report)
    return report


def classify(path: str, config: ScanConfig) -> Classification:
    if not _matches(path, config.scan_patterns):
        return Classification.IGNORED
    if _matches(path, config.exclude_patterns):
        return Classification.EXCLUDED
    if _matches(path, config.generated_patterns):
        return Classification.GENERATED
    return Classification.HANDWRITTEN


def _matches(path: str, patterns: list[re.Pattern[str]]) -> bool:
    return any(pattern.search(path) for pattern in patterns)


def _percent(part: int, whole: int) -> float:
    return round(100 * part / whole, 2) if whole else 0.0
