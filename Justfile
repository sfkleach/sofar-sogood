# Fyne's OpenGL bindings need cgo (and so a C compiler)
export CGO_ENABLED := "1"

# The version comes from the nearest semantic-version git tag, or "dev" if none is available.
version := `git describe --tags --match 'v[0-9]*.[0-9]*.[0-9]*' --dirty 2>/dev/null || echo dev`

default: build

# Build the widget binary
build:
    go build -ldflags "-X main.version={{ version }}" -o sofar-sogood .

# Build and run, optionally with a config file: just run examples/budgets.yaml
run config="": build
    ./sofar-sogood {{ if config != "" { "-config " + config } else { "" } }}

# Run the unit tests
test:
    go test ./...
