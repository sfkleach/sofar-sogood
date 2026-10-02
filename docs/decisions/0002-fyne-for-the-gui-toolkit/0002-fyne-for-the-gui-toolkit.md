# 0002 - Fyne for the GUI toolkit, 2026-10-02

## Issue
The widget must run on Windows, macOS and Linux and must be implemented in Go. We need to choose a GUI toolkit.

## Factors
- The brief requires Go, which rules out toolkits whose application code is in another language.
- The UI is small: a few labels and three buttons.
- A frameless, always-on-top window would suit a widget, but it was not an explicit requirement.
- Release builds and CI should stay simple.

## Options
1. Fyne, a pure-Go toolkit that draws its own widgets with OpenGL.
2. Flutter, a Dart toolkit with prebuilt engines for all three platforms.
3. Wails, a Go backend with a web-view front end.
4. Gio, a lower-level pure-Go immediate-mode toolkit.
5. A system-tray tool rather than a window.

## Pros and Cons of Options

### Option 1: Fyne
- Pros: established and popular in Go; a simple API that is enough for this UI; one codebase and one language.
- Cons: no supported borderless or always-on-top window; needs cgo and a C compiler, so release builds must run on a
  native runner per platform and Linux needs extra development packages; the look is consistent but not native.

### Option 2: Flutter
- Pros: the best polish and tooling; window-management plugins allow widget-style windows; no cgo.
- Cons: it is Dart, so it would break the requirement to use Go. It could only be a front end for a separate Go
  backend, which is far more than this widget needs.

### Option 3: Wails
- Pros: supports frameless windows; the front end can be styled freely.
- Cons: needs WebView2 on Windows and webkit2gtk on Linux; adds a web front-end toolchain.

### Option 4: Gio
- Pros: pure Go with a lighter build.
- Cons: a lower-level API, so more code for the same UI.

### Option 5: System-tray tool
- Pros: very small.
- Cons: not a persistent pane, so it does not match the idea of a widget.

## Outcome and Consequences
Option 1, Fyne. It meets the Go requirement with the least code. The consequences are those listed above: the widget
is an ordinary window rather than a frameless one, and CI and release builds use native runners.

The comparison with the other options was not made at the start of the project; this record was written afterwards.

## Additional Notes
All of the budget logic is in the `budget` package, independent of the UI, so replacing the UI layer, for example
with Wails if a frameless widget becomes important, would mostly mean rewriting `main.go`.
