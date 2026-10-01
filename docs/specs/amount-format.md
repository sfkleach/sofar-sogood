# Budget amount format

The `amount` field of a budget is a single number with an optional currency symbol before it or
a unit after it. It is parsed by `ParseAmount` in `budget/budget.go`.

## Accepted forms

| Example      | Meaning                         |
| ------------ | ------------------------------- |
| `$200`       | Number with a symbol prefix     |
| `£1,250.50`  | Prefix, thousands separator and decimals |
| `800 credits`| Number with a unit suffix       |
| `800credits` | Same, the space is optional     |
| `42`         | Bare number, no symbol or unit  |

- Commas are treated as thousands separators and removed before parsing.
- The decimal separator is `.`.
- Surrounding whitespace is ignored.

## Rejected forms

- No number at all, such as `lots`.
- Both a prefix and a suffix, such as `$5 credits`.
- A negative number.
- A malformed number, such as `1.2.3`.
- Any text after the number that is not a plain unit, such as `$200 USD`.

## Display

The prefix or suffix is reused when showing estimates, so `$200` is displayed as `$66.67 of $200` and
`800 credits` as `266.67 credits of 800 credits`. Values are rounded to two decimal places for display.

## Recommendation

Quote the value in the YAML file, for example `amount: "$200"`, so that values such as `1,250` or
`5e3` are not interpreted by YAML before they reach the parser.
