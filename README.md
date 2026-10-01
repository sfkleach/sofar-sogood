# So far, so good

[![Build and Test](https://github.com/sfkleach/sofar-sogood/actions/workflows/build-and-test.yaml/badge.svg?branch=main)](https://github.com/sfkleach/sofar-sogood/actions/workflows/build-and-test.yaml)

So far, so good is a small desktop widget for Windows, macOS and Linux that helps you track AI credit usage against a
budget. It reads a YAML file listing your budgets, each with a title, a renewal period (weekly, monthly, quarterly or
annual), an amount such as `$200` or `800 credits`, and optionally the days of the week you work (default Monday to
Friday). For each budget it shows one line: how much you should have used by the close of play today, spreading the
renewal amount evenly over the working days of the period.

The config file lives in your user config directory, for example `~/.config/sofar-sogood/budgets.yaml` on Linux, or
anywhere you like via `-config`. See [examples/budgets.yaml](examples/budgets.yaml) and
[docs/specs/amount-format.md](docs/specs/amount-format.md). To build it you need Go and a C compiler (it uses
[Fyne](https://fyne.io)); run `just build`, or `just run examples/budgets.yaml` to try it out.
