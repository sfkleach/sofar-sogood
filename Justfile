# Fyne's OpenGL bindings need cgo (and so a C compiler)
export CGO_ENABLED := "1"

default: build

# Build the widget binary
build:
    go build -o sofar-sogood .

# Build and run, optionally with a config file: just run examples/budgets.yaml
run config="": build
    ./sofar-sogood {{ if config != "" { "-config " + config } else { "" } }}

# Run the unit tests
test:
    go test ./...
