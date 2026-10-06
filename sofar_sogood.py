#!/usr/bin/env python3
"""Print a text report of how much of each budget should have been spent by today.

This is a console version of the So far, so good widget, for machines where the
widget cannot run (for example, where OpenGL is unavailable). It uses only the
Python standard library and reads the same configuration file as the widget.

Usage: sofar_sogood.py [--config PATH] [--date YYYY-MM-DD] [--info] [--version]

The configuration file is TOML, read with the standard library's tomllib, so Python 3.11 or later is needed:

    [[budgets]]
    title = "Claude credits"
    period = "monthly"        # weekly | monthly | quarterly | annual
    amount = "800 credits"    # "$200", "£1,000", "800 credits", ...
    days = ["Mon", "Tue", "Wed", "Thu", "Fri"]   # optional, the default is Mon-Fri
"""
import argparse
import calendar
import datetime
import os
import re
import sys
import textwrap

try:
    import tomllib
except ImportError:  # Python before 3.11.
    sys.exit("sofar_sogood.py needs Python 3.11 or later, for the tomllib module.")

VERSION = "dev"

PERIODS = ("weekly", "monthly", "quarterly", "annual")

# Python numbers weekdays from Monday = 0, which also matches the order used in the config.
DAY_NAMES = {"mon": 0, "tue": 1, "wed": 2, "thu": 3, "fri": 4, "sat": 5, "sun": 6}
DEFAULT_DAYS = frozenset(range(5))


class ConfigError(Exception):
    """Raised when the configuration file cannot be understood."""


# --------------------------------------------------------------------------------------
# Locating the config file.
# --------------------------------------------------------------------------------------


def default_config_path():
    """Return the default config location, matching the widget on each platform."""
    if sys.platform == "win32":
        base = os.environ.get("APPDATA") or os.path.expanduser("~")
    elif sys.platform == "darwin":
        base = os.path.expanduser("~/Library/Application Support")
    else:
        base = os.environ.get("XDG_CONFIG_HOME") or os.path.expanduser("~/.config")
    return os.path.join(base, "sofar-sogood", "budgets.toml")


# --------------------------------------------------------------------------------------
# Budgets.
# --------------------------------------------------------------------------------------


class Amount:
    """A quantity with a display prefix (such as '$') or suffix (such as 'credits')."""

    def __init__(self, value, prefix="", suffix=""):
        self.value = value
        self.prefix = prefix
        self.suffix = suffix

    @classmethod
    def parse(cls, text):
        """Parse strings such as '$200', '£1,250.50' or '800 credits'."""
        text = str(text).strip()
        match = re.search(r"[0-9.,-]", text)
        if not match:
            raise ConfigError("amount %r has no number" % text)
        start = match.start()
        end = start
        while end < len(text) and text[end] in "0123456789.,-":
            end += 1
        prefix = text[:start].strip()
        suffix = text[end:].strip()
        if prefix and suffix:
            raise ConfigError("amount %r: use either a prefix or a suffix, not both" % text)
        try:
            value = float(text[start:end].replace(",", ""))
        except ValueError:
            raise ConfigError("amount %r: not a valid number" % text) from None
        if value < 0:
            raise ConfigError("amount %r must not be negative" % text)
        return cls(value, prefix, suffix)

    def format(self, value):
        """Render value with this amount's prefix or suffix, rounded to two decimal places."""
        rounded = int(value * 100 + 0.5) / 100
        number = ("%.2f" % rounded).rstrip("0").rstrip(".")
        if self.suffix:
            return "%s %s" % (number, self.suffix)
        return self.prefix + number


class Budget:
    """One entry from the config file."""

    def __init__(self, title, period, days, amount):
        self.title = title
        self.period = period
        self.days = days
        self.amount = amount

    @classmethod
    def from_mapping(cls, number, raw):
        """Validate a raw mapping from the config file."""
        for key in raw:
            if key not in ("title", "period", "amount", "days"):
                raise ConfigError("budget %d: unknown key %r" % (number, key))
        title = str(raw.get("title", "")).strip()
        where = "budget %d (%r)" % (number, title)
        if not title:
            raise ConfigError("budget %d: missing title" % number)
        period = str(raw.get("period", "")).strip().lower()
        if period not in PERIODS:
            raise ConfigError("%s: period must be weekly, monthly, quarterly or annual" % where)
        try:
            amount = Amount.parse(raw.get("amount", ""))  # A bare TOML number is accepted as well as a string.
        except ConfigError as err:
            raise ConfigError("%s: %s" % (where, err)) from None
        names = raw.get("days")
        if names is None or names == []:
            days = DEFAULT_DAYS
        else:
            if isinstance(names, str):
                names = [names]
            days = set()
            for name in names:
                day = DAY_NAMES.get(str(name).strip().lower()[:3]) if len(str(name).strip()) >= 3 else None
                if day is None:
                    raise ConfigError("%s: unknown day %r" % (where, name))
                days.add(day)
        return cls(title, period, frozenset(days), amount)

    def bounds(self, today):
        """Return the first and last day of the renewal period containing today."""
        if self.period == "weekly":
            first = today - datetime.timedelta(days=today.weekday())
            return first, first + datetime.timedelta(days=6)
        if self.period == "monthly":
            last_day = calendar.monthrange(today.year, today.month)[1]
            return today.replace(day=1), today.replace(day=last_day)
        if self.period == "quarterly":
            first_month = (today.month - 1) // 3 * 3 + 1
            last_month = first_month + 2
            last_day = calendar.monthrange(today.year, last_month)[1]
            return datetime.date(today.year, first_month, 1), datetime.date(today.year, last_month, last_day)
        return datetime.date(today.year, 1, 1), datetime.date(today.year, 12, 31)

    def working_days(self, first, last):
        """Count the working days in the inclusive range first to last."""
        count = 0
        day = first
        while day <= last:
            if day.weekday() in self.days:
                count += 1
            day += datetime.timedelta(days=1)
        return count

    def estimate(self, today):
        """Return (spent, days_elapsed, days_total) as at the close of play today."""
        first, last = self.bounds(today)
        elapsed = self.working_days(first, today)
        total = self.working_days(first, last)
        spent = self.amount.value * elapsed / total if total else 0.0
        return spent, elapsed, total


def load_budgets(path):
    """Load and validate the config file. A missing file counts as no budgets."""
    try:
        with open(path, "rb") as f:
            data = tomllib.load(f)
    except FileNotFoundError:
        return []
    except tomllib.TOMLDecodeError as err:
        raise ConfigError(str(err)) from None
    # Reject unknown keys so that a typo such as "perod" is reported instead of silently ignored.
    for key in data:
        if key != "budgets":
            raise ConfigError("unknown key %r" % key)
    raw_budgets = data.get("budgets", [])
    if not isinstance(raw_budgets, list):
        raise ConfigError("'budgets' must be an array of tables, written [[budgets]]")
    return [Budget.from_mapping(n, raw) for n, raw in enumerate(raw_budgets, start=1)]


# --------------------------------------------------------------------------------------
# The report.
# --------------------------------------------------------------------------------------


def report(budgets, today, path):
    """Return the text of the report."""
    out = ["So far, so good: %s" % today.strftime("%A %Y-%m-%d"), "Config: %s" % path, ""]
    if not budgets:
        out.append("no budgets")
        return "\n".join(out)
    rows = []
    for b in budgets:
        spent, elapsed, total = b.estimate(today)
        percent = spent / b.amount.value * 100 if b.amount.value else 0.0
        rows.append((b.title, "%s of %s" % (b.amount.format(spent), b.amount.format(b.amount.value)),
                     "(%d/%d days, %.0f%%)" % (elapsed, total, percent)))
    widths = [max(len(row[col]) for row in rows) for col in range(3)]
    for title, spend, days in rows:
        out.append("%s  %s  %s" % (title.ljust(widths[0]), spend.ljust(widths[1]), days))
    return "\n".join(out)


def same_folder(a, b):
    """Return True if two file paths are in the same folder."""
    return os.path.dirname(os.path.abspath(a)) == os.path.dirname(os.path.abspath(b))


def info_text(path):
    """Return the same explanation that the widget shows in its Info pane."""
    default = default_config_path()
    def wrap(text):
        """Wrap prose to fit a console; the paths below are deliberately left unwrapped."""
        return textwrap.fill(text, 78)

    paragraphs = [
        wrap("So far, so good helps you track AI credit usage against a budget."),
        wrap("Each budget shows how much you should have used by the close of play today, spreading the renewal "
             "amount evenly over the working days of the period."),
        "Budgets are defined in this config file:\n%s" % path,
    ]
    # The default location is only worth mentioning if the file in use is somewhere else.
    if not same_folder(path, default):
        paragraphs.append("The default config file is:\n%s" % default)
    paragraphs += [
        wrap('Each amount is a number with an optional symbol or unit, such as "$200", "£1,000" or "800 credits".'),
        "(Pass --config PATH to use a different file.)",
        "Version: %s" % VERSION,
    ]
    return "\n\n".join(paragraphs)


def main(argv=None):
    parser = argparse.ArgumentParser(description="Print a text report of budget spend to date.")
    parser.add_argument("--config", help="path to budgets.toml (default: the user config directory)")
    parser.add_argument("--date", help="report as at this date, YYYY-MM-DD (default: today)")
    parser.add_argument("--info", action="store_true", help="explain what this does, and where the config file is")
    parser.add_argument("--version", action="version", version=VERSION)
    args = parser.parse_args(argv)

    path = args.config or default_config_path()
    if args.info:
        print(info_text(path))
        return 0
    try:
        today = datetime.date.fromisoformat(args.date) if args.date else datetime.date.today()
        budgets = load_budgets(path)
    except ValueError:
        print("error: --date must be in the form YYYY-MM-DD", file=sys.stderr)
        return 2
    except ConfigError as err:
        print("config error: %s" % err, file=sys.stderr)
        return 1
    print(report(budgets, today, path))
    return 0


if __name__ == "__main__":
    sys.exit(main())
