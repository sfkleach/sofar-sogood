# Change Log for So far, so good

Following the style in https://keepachangelog.com/en/1.0.0/

Release headings take the form `## vX.Y.Z, Short title YYYY-MM-DD`. The top section must be a release heading, not
`Unreleased`, before a tag is pushed; `just shippable` checks this.

## v0.2.0, TOML config 2026-10-02

### Changed

- **Breaking:** the config file is now TOML, not YAML. It is called `budgets.toml` and each budget is a `[[budgets]]`
  table, with amounts that have a symbol or unit written as quoted strings. See `examples/budgets.toml`. Existing
  `budgets.yaml` files are no longer read and must be converted.
- Unknown keys in the config file are now reported as errors instead of being ignored, so a typo such as `perod`
  is caught. A bare number is accepted as an amount, for example `amount = 800`.
- **Breaking:** command-line options now use the GNU style, with double dashes: `--config` and `--version`, and
  `-v` for short. The old single-dash `-config` and `-version` are no longer accepted. The options are handled by
  Cobra, which also provides `--help`.

### Added

- `sofar_sogood.py`, a console report that needs only the Python standard library (Python 3.11 or later), for machines
  where the widget cannot run, such as those without OpenGL. Its `--info` option explains the report and shows where
  the config file is.

## v0.1.0, First release 2026-10-02

### Added

- Desktop widget (Fyne) showing, for each budget in a YAML config file, how much should have been spent by the close
  of play today, based on the working days of the renewal period so far.
- Weekly, monthly, quarterly and annual renewal periods, with configurable working days (default Monday to Friday).
- Amounts with a currency symbol or unit, such as `$200`, `£1,000` or `800 credits`; see `docs/specs/amount-format.md`.
- Refresh, Info and Dismiss buttons. The Info pane explains the widget and shows the config file in use (and the
  default config file, if that is in a different folder), the amount format and the version.
- The widget refreshes itself once an hour, on the hour, so it can be left running overnight.
- `-config` and `-version` command-line flags.
- GitHub workflows for build-and-test on Linux (amd64 and arm64), macOS and Windows, and for building releases from
  `vX.Y.Z` tags. Release notes come from this changelog and cover the Windows Defender and macOS Gatekeeper caveats.
- `install.sh` to install the latest release on Linux and macOS (and from a Windows bash shell).
- `just test` runs the unit tests, `go vet`, golangci-lint and a gofmt check (without reformatting), and CI uses it.
  `just install-deps` installs golangci-lint v2 if it is missing.
- Decision records under `docs/decisions/`, with `just add-decision`.
