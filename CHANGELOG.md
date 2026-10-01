# Change Log for So far, so good

Following the style in https://keepachangelog.com/en/1.0.0/

Release headings take the form `## vX.Y.Z, Short title YYYY-MM-DD`. The top section must be a release heading, not
`Unreleased`, before a tag is pushed; `just shippable` checks this.

## Unreleased

### Added

- Desktop widget (Fyne) showing, for each budget in a YAML config file, how much should have been spent by the close
  of play today, based on the working days of the renewal period so far.
- Weekly, monthly, quarterly and annual renewal periods, with configurable working days (default Monday to Friday).
- Amounts with a currency symbol or unit, such as `$200`, `£1,000` or `800 credits`.
- Refresh, Info and Dismiss buttons. The Info pane explains the widget and shows the config file location and version.
- `-config` and `-version` command-line flags.
- GitHub workflows for build-and-test on Linux, macOS and Windows, and for building releases from `vX.Y.Z` tags.
