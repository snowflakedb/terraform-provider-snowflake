#!/usr/bin/env python3
"""Count generated vs hand-written files and lines in a local Git checkout."""

from __future__ import annotations

import argparse
import json
import re
import subprocess
import sys
from pathlib import Path
from typing import Sequence

from stats_common import Report, ScanConfig, run_scan


class LocalBackend:
    def __init__(self, root: Path, output_format: str, print_header: bool) -> None:
        self.root = root.resolve()
        self.output_format = output_format
        self.print_header = print_header

    @property
    def repository(self) -> str:
        # The directory name is useless for worktrees, so prefer the name from the origin URL.
        try:
            return self._git("remote", "get-url", "origin").removesuffix(".git").rsplit("/", 1)[-1]
        except subprocess.CalledProcessError:
            return self.root.name

    @property
    def branch(self) -> str:
        return self._git("branch", "--show-current") or "DETACHED"

    def get_revision(self) -> str:
        return self._git("rev-parse", "HEAD")

    def list_paths(self) -> list[str]:
        return [path for path in self._git("ls-files", "-z").split("\0") if path]

    def line_counts(self, paths: list[str]) -> dict[str, int]:
        counts: dict[str, int] = {}
        for path in paths:
            file_path = self.root / path
            if file_path.is_file():  # tracked, but possibly deleted in the working tree
                counts[path] = len(file_path.read_bytes().splitlines())
        return counts

    def save_report(self, report: Report) -> None:
        if self.output_format == "json":
            print(json.dumps(report.to_dict(), indent=2))
        else:
            print_plain_report(report, self.print_header)

    def _git(self, *args: str) -> str:
        output = subprocess.check_output(["git", "-C", str(self.root), *args])
        return output.decode("utf-8", errors="surrogateescape").strip("\n")


def parse_cli_args(argv: Sequence[str] | None = None) -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--root", type=Path, default=Path("."), help="Git checkout to scan (default: current directory)")
    parser.add_argument("--scope", default="all", help="Label shown in the report (default: all)")
    parser.add_argument("--scan", action="append", help=r"Regex for paths to count; repeatable (default: \.go$)")
    parser.add_argument("--generated", action="append", help=r"Regex for generated paths; repeatable (default: _gen(_test)?\.go$)")
    parser.add_argument("--exclude", action="append", help=r"Regex for paths to skip; repeatable (default: (^|/)testdata/)")
    parser.add_argument("--format", choices=("plain", "json"), default="plain", help="Output format (default: plain)")
    parser.add_argument("--no-header", dest="print_header", action="store_false", help="Plain format: print only the row")
    parser.add_argument("--show-files", action="store_true", help="List every counted file")
    args = parser.parse_args(argv)

    args.scan = args.scan or [r"\.go$"]
    args.generated = args.generated or [r"_gen(_test)?\.go$"]
    args.exclude = args.exclude or [r"(^|/)testdata/"]
    return args


def local_main(argv: Sequence[str] | None = None) -> int:
    args = parse_cli_args(argv)
    try:
        config = ScanConfig.from_strings(args.scan, args.generated, args.exclude)
        run_scan(config, LocalBackend(args.root, args.format, args.print_header), args.scope, args.show_files)
        return 0
    except (OSError, re.error, subprocess.CalledProcessError) as exc:
        print(f"error: {exc}", file=sys.stderr)
        return 2


def print_plain_report(report: Report, print_header: bool) -> None:
    scope_width = 24

    if print_header:
        print(f"Repository: {report.repository}")
        print(f"Branch:     {report.branch}")
        print(f"Commit:     {report.commit_sha}")
        print(f"Ruleset:    {report.ruleset_version}")
        print()
        print(f"{'':<{scope_width}} | {'FILES':^29} | {'LINES':^31} | {'EXCLUDED':^17}")
        print(
            f"{'SCOPE':<{scope_width}} | {'generated':>9} {'handwritten':>11} {'gen %':>7}"
            f" | {'generated':>11} {'handwritten':>11} {'gen %':>7} | {'files':>8} {'lines':>8}"
        )
        print(f"{'-' * scope_width}-+-{'-' * 29}-+-{'-' * 31}-+-{'-' * 17}")

    print(
        f"{report.scope:<{scope_width}} | {report.generated.files:>9} {report.handwritten.files:>11} {report.generated_percent:>7}"
        f" | {report.generated.lines:>11} {report.handwritten.lines:>11} {report.generated_line_percent:>7}"
        f" | {report.excluded.files:>8} {report.excluded.lines:>8}"
    )

    if report.files:
        print()
        print(f"{'CLASSIFICATION':<14} {'LINES':>8}  PATH")
        for file in report.files:
            print(f"{file.classification:<14} {file.lines:>8}  {file.path}")


if __name__ == "__main__":
    raise SystemExit(local_main())
