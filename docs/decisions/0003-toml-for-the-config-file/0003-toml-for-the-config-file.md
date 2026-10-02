# 0003 - TOML for the config file, 2026-10-02

## Issue
The budgets are defined in a config file. The first version used YAML. We need to decide whether to keep it, given
that a console version of the report, `sofar_sogood.py`, was added for machines where the widget cannot run.

## Factors
- The Python script should need only the standard library, which has a TOML parser (`tomllib`, from Python 3.11) but
  no YAML parser.
- YAML coerces unquoted values, so amounts such as `1,250` or `5e3` needed quoting advice, and indentation matters.
- Go's `github.com/BurntSushi/toml` was already in the dependency tree through Fyne.
- No stable release had been made, so there were no existing config files to keep working.

## Options
1. Keep YAML, and write a small parser for the subset we use in the Python script.
2. Switch to TOML.
3. Switch to JSON, which both languages read in the standard library.

## Pros and Cons of Options

### Option 1: YAML
- Pros: no change to the Go code or the example file.
- Cons: a hand-written YAML subset parser in Python, which will not understand anything beyond what it was written
  for; YAML's implicit typing and indentation rules.

### Option 2: TOML
- Pros: parsed by the Python standard library and by a library we already depend on; explicit types; supports
  comments; a repeated `[[budgets]]` table reads naturally.
- Cons: amounts with a symbol or unit must be quoted; requires Python 3.11 or later for the script.

### Option 3: JSON
- Pros: parsed by the standard library everywhere.
- Cons: no comments, which are useful in a hand-edited file; more punctuation to get wrong.

## Outcome and Consequences
Option 2, TOML. The file is `budgets.toml` and each budget is a `[[budgets]]` table. Unknown keys are reported as
errors, so a typo such as `perod` is not silently ignored. A bare number is accepted for `amount` as well as a
quoted string. YAML config files are no longer read.

## Additional Notes
`sofar_sogood.py` exits with a clear message on a Python older than 3.11.
