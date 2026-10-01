package budget

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type fileFormat struct {
	Budgets []struct {
		Title  string   `yaml:"title"`
		Period string   `yaml:"period"`
		Days   []string `yaml:"days"`
		Amount string   `yaml:"amount"`
	} `yaml:"budgets"`
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
	return filepath.Join(dir, "sofar-sogood", "budgets.yaml"), nil
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
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, err
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
		amt, err := ParseAmount(raw.Amount)
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
