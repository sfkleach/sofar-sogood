# 0001 - Version from git tag, 2026-10-01

## Issue
The widget needs a version number that is shown in the Info pane and by `--version`, and that matches the release
tag. We need to decide where the version is stored.

## Factors
- The release workflow is triggered by pushing a semantic-version tag such as `v1.2.3`.
- A version stored in a file must be kept in step with the tag, usually with a commit just to bump it.
- A local build with no tag available should still work.

## Options
1. Store the version in a file (as `monogram-go` does with `version.txt` and `mg/version.go`), bumped by a script.
2. Take the version from the git tag at build time, stamping it into the binary with
   `-ldflags "-X main.version=..."`.

## Pros and Cons of Options

### Option 1: version file
- Pros: the version is visible in the source tree; `go run` shows it.
- Cons: a bump commit for every release; the file and the tag can disagree.

### Option 2: git tag
- Pros: a single source of truth; no bump commits; the release workflow stamps the tag directly.
- Cons: a build outside a git checkout or without tags shows `dev`.

## Outcome and Consequences
Option 2. `main.version` defaults to `dev`. `just build` derives the version from `git describe`, and the release
workflow passes the tag. `scripts/check-changelog.py` checks that the top CHANGELOG section matches the tag, which
replaces the consistency check that a version file would otherwise need.

## Additional Notes
Because Fyne needs cgo, release binaries are built on a native runner per platform rather than cross-compiled.
