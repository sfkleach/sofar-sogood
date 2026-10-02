package budget

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type fileFormat struct {
	Budgets []struct {
		Title  string   `toml:"title"`
		Period string   `toml:"period"`
		Days   []string `toml:"days"`
		// Amount is usually a string such as "$200", but a bare TOML number is accepted too.
		Amount any `toml:"amount"`
	} `toml:"budgets"`
}

var dayNames = map[string]time.Weekday{
	"mon": time.Monday, "tue": time.Tuesday, "wed": time.Wednesday,
	"thu": time.Thursday, "fri": time.Friday, "sat": time.Saturday, "sun": time.Sunday,
}

// DefaultPath is where the configuration file lives unless overridden.
func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "sofar-sogood", "budgets.toml"), nil
}

// SameFolder reports whether two file paths are in the same folder.
func SameFolder(a, b string) bool {
	dir := func(p string) string {
		// Fall back to the cleaned path if it cannot be made absolute, which only
		// happens when the working directory is unavailable.
		if abs, err := filepath.Abs(p); err == nil {
			p = abs
		}
		return filepath.Dir(p)
	}
	return dir(a) == dir(b)
}

// Load reads and validates the configuration file. A missing file is
// treated as an empty list of budgets.
func Load(path string) ([]Budget, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

// Parse validates configuration file contents.
func Parse(data []byte) ([]Budget, error) {
	var f fileFormat
	md, err := toml.Decode(string(data), &f)
	if err != nil {
		return nil, err
	}
	// Reject unknown keys so that a typo such as "perod" is reported instead of silently ignored.
	if undecoded := md.Undecoded(); len(undecoded) > 0 {
		return nil, fmt.Errorf("unknown key %q", undecoded[0].String())
	}
	var out []Budget
	for i, raw := range f.Budgets {
		where := fmt.Sprintf("budget %d (%q)", i+1, raw.Title)
		if strings.TrimSpace(raw.Title) == "" {
			return nil, fmt.Errorf("budget %d: missing title", i+1)
		}
		period := Period(strings.ToLower(strings.TrimSpace(raw.Period)))
		switch period {
		case Weekly, Monthly, Quarterly, Annual:
		default:
			return nil, fmt.Errorf("%s: period must be weekly, monthly, quarterly or annual", where)
		}
		amt, err := amountFromConfig(raw.Amount)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", where, err)
		}
		days := map[time.Weekday]bool{}
		if len(raw.Days) == 0 {
			for _, d := range []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday} {
				days[d] = true
			}
		}
		for _, name := range raw.Days {
			if len(name) < 3 {
				return nil, fmt.Errorf("%s: unknown day %q", where, name)
			}
			d, ok := dayNames[strings.ToLower(name[:3])]
			if !ok {
				return nil, fmt.Errorf("%s: unknown day %q", where, name)
			}
			days[d] = true
		}
		out = append(out, Budget{Title: raw.Title, Period: period, Days: days, Amount: amt})
	}
	return out, nil
}

// amountFromConfig converts the decoded TOML value of the amount key.
func amountFromConfig(v any) (Amount, error) {
	switch v := v.(type) {
	case string:
		return ParseAmount(v)
	case int64:
		return ParseAmount(strconv.FormatInt(v, 10))
	case float64:
		return ParseAmount(strconv.FormatFloat(v, 'f', -1, 64))
	case nil:
		return Amount{}, fmt.Errorf("missing amount")
	default:
		return Amount{}, fmt.Errorf("amount must be a string or a number")
	}
}
