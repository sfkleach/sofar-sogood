// Package budget loads budget definitions and computes how much of each
// should have been spent by a given day.
package budget

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Period is how often a budget renews. Periods are calendar-aligned:
// weeks start on Monday, quarters on Jan/Apr/Jul/Oct 1st.
type Period string

const (
	Weekly    Period = "weekly"
	Monthly   Period = "monthly"
	Quarterly Period = "quarterly"
	Annual    Period = "annual"
)

// Amount is a quantity with a display prefix ("$") or suffix ("credits").
type Amount struct {
	Value  float64
	Prefix string
	Suffix string
}

// ParseAmount parses strings such as "$200", "£1,250.50" or "800 credits".
func ParseAmount(s string) (Amount, error) {
	orig := s
	s = strings.TrimSpace(s)
	isNum := func(r rune) bool { return unicode.IsDigit(r) || r == '.' || r == ',' || r == '-' }

	start := strings.IndexFunc(s, isNum)
	if start < 0 {
		return Amount{}, fmt.Errorf("amount %q has no number", orig)
	}
	end := start
	for end < len(s) && isNum(rune(s[end])) {
		end++
	}
	prefix := strings.TrimSpace(s[:start])
	suffix := strings.TrimSpace(s[end:])
	if prefix != "" && suffix != "" {
		return Amount{}, fmt.Errorf("amount %q: use either a prefix or a suffix, not both", orig)
	}
	v, err := strconv.ParseFloat(strings.ReplaceAll(s[start:end], ",", ""), 64)
	if err != nil {
		return Amount{}, fmt.Errorf("amount %q: %w", orig, err)
	}
	if v < 0 {
		return Amount{}, fmt.Errorf("amount %q must not be negative", orig)
	}
	return Amount{Value: v, Prefix: prefix, Suffix: suffix}, nil
}

// Format renders v using this amount's prefix/suffix.
func (a Amount) Format(v float64) string {
	num := strconv.FormatFloat(math.Round(v*100)/100, 'f', -1, 64)
	if a.Suffix != "" {
		return num + " " + a.Suffix
	}
	return a.Prefix + num
}

// Budget is one entry from the configuration file.
type Budget struct {
	Title  string
	Period Period
	Days   map[time.Weekday]bool
	Amount Amount
}

// Bounds returns the first and last day (midnight, inclusive) of the
// renewal period containing now.
func (p Period) Bounds(now time.Time) (first, last time.Time) {
	y, m, d := now.Date()
	loc := now.Location()
	switch p {
	case Weekly:
		offset := (int(now.Weekday()) + 6) % 7 // Monday = 0
		first = time.Date(y, m, d-offset, 0, 0, 0, 0, loc)
		last = first.AddDate(0, 0, 6)
	case Monthly:
		first = time.Date(y, m, 1, 0, 0, 0, 0, loc)
		last = first.AddDate(0, 1, -1)
	case Quarterly:
		qm := time.Month((int(m)-1)/3*3 + 1)
		first = time.Date(y, qm, 1, 0, 0, 0, 0, loc)
		last = first.AddDate(0, 3, -1)
	default: // Annual
		first = time.Date(y, 1, 1, 0, 0, 0, 0, loc)
		last = time.Date(y, 12, 31, 0, 0, 0, 0, loc)
	}
	return first, last
}

// WorkingDays counts the working days in [from, to], inclusive.
func (b Budget) WorkingDays(from, to time.Time) int {
	n := 0
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		if b.Days[d.Weekday()] {
			n++
		}
	}
	return n
}

// Estimate is the expected position of a budget on a particular day.
type Estimate struct {
	Budget      Budget
	DaysElapsed int // working days up to and including today
	DaysTotal   int // working days in the whole period
	Spent       float64
}

// Estimate computes the expected spend by the close of play on now's date.
// Today counts in full if it is a working day, as does every earlier working
// day of the period.
func (b Budget) Estimate(now time.Time) Estimate {
	first, last := b.Period.Bounds(now)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	e := Estimate{
		Budget:      b,
		DaysElapsed: b.WorkingDays(first, today),
		DaysTotal:   b.WorkingDays(first, last),
	}
	if e.DaysTotal > 0 {
		e.Spent = b.Amount.Value * float64(e.DaysElapsed) / float64(e.DaysTotal)
	}
	return e
}

// Line renders the estimate as a single line of text.
func (e Estimate) Line() string {
	percent := 0.0
	if e.Budget.Amount.Value > 0 {
		percent = e.Spent / e.Budget.Amount.Value * 100
	}
	return fmt.Sprintf("%s: %s of %s (%d/%d days, %.0f%%)",
		e.Budget.Title,
		e.Budget.Amount.Format(e.Spent),
		e.Budget.Amount.Format(e.Budget.Amount.Value),
		e.DaysElapsed, e.DaysTotal, percent)
}
