# 0004 - Cobra for command-line options, 2026-10-02

## Issue
Go's standard `flag` package is used with single-dash long options such as `-config` and `-version`. Almost every
other command-line tool, and the Python script in this repository, uses double dashes for long options.

## Factors
- Consistency between the widget and `sofar_sogood.py`, and with common practice (`--config`, `--version`).
- The project owner's usual practice is to use Cobra.
- The widget has very few options, so the extra dependency should be small in effect.

## Options and Outcome
Use Cobra (`github.com/spf13/cobra`), which uses POSIX-style options through pflag, and gives `--help` and `--version`
for free. The alternative of keeping `flag` was rejected because the standard library's own help output also shows
single-dash long options.

## Consequences
- The options are now `--config`, `--version` (and `-v`) and `--help`. The old single-dash `-config` and `-version`
  are no longer accepted.
- Cobra's version output is configured to print just the version, as before.
- Cobra and pflag are added as dependencies, with Windows-only `mousetrap` pulled in by Cobra.
