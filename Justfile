set windows-shell := ["bash", "-uc"]

# Fyne's OpenGL bindings need cgo (and so a C compiler)
export CGO_ENABLED := "1"

# The version comes from the nearest semantic-version git tag, or "dev" if none is available.
version := `git describe --tags --match 'v[0-9]*.[0-9]*.[0-9]*' --dirty 2>/dev/null || echo dev`

[private]
default:
    @just --list

# Build the widget binary
build:
    go build -ldflags "-X main.version={{ version }}" -o sofar-sogood .

# Build and run, optionally with a config file: just run examples/budgets.toml
run config="": build
    ./sofar-sogood {{ if config != "" { "--config " + config } else { "" } }}

# Run the unit tests
unittest:
    go test -cover ./...
    @echo

# Run go vet
vet:
    go vet ./...

# Install the tools needed by the other recipes, if they are missing
install-deps:
    #!/usr/bin/env bash
    set -euo pipefail
    export PATH="$(go env GOPATH)/bin:$PATH"
    # Install if missing, or if an older (v1) major version is installed.
    if ! command -v golangci-lint >/dev/null 2>&1 || ! golangci-lint --version | grep -Eq 'version v?2\.'; then
        echo "Installing golangci-lint..."
        go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
    fi

# Run the linter
lint: install-deps
    #!/usr/bin/env bash
    set -euo pipefail
    # The installer puts golangci-lint in GOPATH/bin, which may not be on the PATH.
    export PATH="$(go env GOPATH)/bin:$PATH"
    golangci-lint run ./...

# Check the formatting without changing any files; use "just fmt" to fix it
fmt-check:
    #!/usr/bin/env bash
    set -euo pipefail
    unformatted=$(gofmt -l .)
    if [ -n "$unformatted" ]; then
        echo "These files are not gofmt-formatted:" >&2
        echo "$unformatted" >&2
        exit 1
    fi

# Reformat all Go files in place
fmt:
    gofmt -w .

# Run all the checks: unit tests, vet, linter and formatting
test: unittest vet lint fmt-check

# Generate a coverage report for the unit tests
test-coverage:
    rm -rf _build
    mkdir -p _build/
    go test -cover -coverprofile=_build/unittest.out ./...
    go tool cover -html=_build/unittest.out -o _build/unittest.html

# Print the version that a build would be stamped with
get-version:
    @echo {{ version }}

# Install the widget into your bin directory.
install: build
    # Install to GOBIN if set, otherwise fall back to $(go env GOPATH)/bin.
    if [ -n "$(go env GOBIN)" ]; then \
      cp sofar-sogood "$(go env GOBIN)/sofar-sogood"; \
    else \
      cp sofar-sogood "$(go env GOPATH)/bin/sofar-sogood"; \
    fi

# Remove the binary and the test coverage reports
clean:
    rm -f sofar-sogood
    rm -rf _build

# Initialize decision records
init-decisions:
    python3 scripts/decisions.py --init

# Add a new decision record
add-decision TOPIC:
    python3 scripts/decisions.py --add "{{TOPIC}}"

# Check that the top CHANGELOG section is a release, not "Unreleased".
shippable:
    python3 scripts/check-changelog.py

# Sign and push a release tag, triggering the release workflow, which builds the binaries
# and publishes a GitHub release using the matching CHANGELOG section as its notes.
# Usage: just release v0.2.0
release VERSION:
    #!/usr/bin/env bash
    set -euo pipefail
    # Refuse to tag a dirty working tree: uncommitted changes would not be part of the release.
    if ! git diff --quiet || ! git diff --cached --quiet; then
        echo "error: working tree has uncommitted changes; commit or stash them before tagging." >&2
        exit 1
    fi
    python3 scripts/check-changelog.py "{{VERSION}}"
    echo "Signing and pushing tag {{VERSION}}..."
    git tag -s "{{VERSION}}" -m "Release {{VERSION}}"
    git push origin "{{VERSION}}"
    echo "Tag pushed. Monitor CI at: https://github.com/sfkleach/sofar-sogood/actions"
