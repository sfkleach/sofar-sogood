package budget

import (
	"testing"
	"time"
)

func day(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 15, 30, 0, 0, time.UTC)
}

func TestParseAmount(t *testing.T) {
	cases := []struct {
		in     string
		val    float64
		prefix string
		suffix string
	}{
		{"$200", 200, "$", ""},
		{"£1,250.50", 1250.5, "£", ""},
		{"800 credits", 800, "", "credits"},
		{"  42  ", 42, "", ""},
	}
	for _, c := range cases {
		a, err := ParseAmount(c.in)
		if err != nil {
			t.Fatalf("%q: %v", c.in, err)
		}
		if a.Value != c.val || a.Prefix != c.prefix || a.Suffix != c.suffix {
			t.Errorf("%q: got %+v", c.in, a)
		}
	}
	for _, bad := range []string{"", "lots", "$5 credits", "-3"} {
		if _, err := ParseAmount(bad); err == nil {
			t.Errorf("%q: expected error", bad)
		}
	}
}

func TestFormat(t *testing.T) {
	if got := (Amount{Prefix: "$"}).Format(66.666); got != "$66.67" {
		t.Errorf("got %s", got)
	}
	if got := (Amount{Suffix: "credits"}).Format(400); got != "400 credits" {
		t.Errorf("got %s", got)
	}
}

func TestEstimateWeekly(t *testing.T) {
	bs, err := Parse([]byte(`[[budgets]]
title = "Lunch"
period = "weekly"
amount = "$100"
`))
	if err != nil {
		t.Fatal(err)
	}
	// Wed 2026-09-30: Mon, Tue, Wed done = 3 of 5.
	e := bs[0].Estimate(day(2026, 9, 30))
	if e.DaysElapsed != 3 || e.DaysTotal != 5 || e.Spent != 60 {
		t.Errorf("got %+v", e)
	}
	// Sunday is not a working day: the week is complete.
	e = bs[0].Estimate(day(2026, 10, 4))
	if e.DaysElapsed != 5 || e.Spent != 100 {
		t.Errorf("got %+v", e)
	}
}

func TestEstimateMonthlyWithSaturdays(t *testing.T) {
	bs, err := Parse([]byte(`[[budgets]]
title = "Cloud"
period = "monthly"
days = ["Mon", "Sat"]
amount = "800 credits"
`))
	if err != nil {
		t.Fatal(err)
	}
	// October 2026: Mondays 5,12,19,26; Saturdays 3,10,17,24,31 = 9 days.
	// By Sat 10th: Sat 3, Mon 5, Sat 10 = 3.
	e := bs[0].Estimate(day(2026, 10, 10))
	if e.DaysTotal != 9 || e.DaysElapsed != 3 {
		t.Fatalf("got %+v", e)
	}
	if got, want := e.Line(), "Cloud: 266.67 credits of 800 credits (3/9 days, 33%)"; got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestBounds(t *testing.T) {
	f, l := Quarterly.Bounds(day(2026, 8, 15))
	if f.Month() != time.July || f.Day() != 1 || l.Month() != time.September || l.Day() != 30 {
		t.Errorf("quarter: %v %v", f, l)
	}
	f, l = Annual.Bounds(day(2026, 8, 15))
	if f.YearDay() != 1 || l.YearDay() != 365 {
		t.Errorf("year: %v %v", f, l)
	}
	f, _ = Weekly.Bounds(day(2026, 10, 4)) // Sunday
	if f.Weekday() != time.Monday || f.Day() != 28 {
		t.Errorf("week: %v", f)
	}
}

func TestParseErrors(t *testing.T) {
	for _, bad := range []string{
		"[[budgets]]\nperiod = \"weekly\"\namount = \"$5\"",
		"[[budgets]]\ntitle = \"x\"\nperiod = \"daily\"\namount = \"$5\"",
		"[[budgets]]\ntitle = \"x\"\nperiod = \"weekly\"\namount = \"$5\"\ndays = [\"Funday\"]",
		"[[budgets]]\ntitle = \"x\"\nperiod = \"weekly\"\namount = \"lots\"",
		"[[budgets]]\ntitle = \"x\"\nperiod = \"weekly\"",
		"[[budgets]]\ntitle = \"x\"\nperiod = \"weekly\"\namount = \"$5\"\nperod = \"typo\"",
		"[[budgets]\n",
	} {
		if _, err := Parse([]byte(bad)); err == nil {
			t.Errorf("expected error for %q", bad)
		}
	}
}

func TestParseNumericAmount(t *testing.T) {
	bs, err := Parse([]byte("[[budgets]]\ntitle = \"x\"\nperiod = \"weekly\"\namount = 250\n"))
	if err != nil {
		t.Fatal(err)
	}
	if bs[0].Amount.Value != 250 || bs[0].Amount.Prefix != "" || bs[0].Amount.Suffix != "" {
		t.Errorf("got %+v", bs[0].Amount)
	}
}

func TestSameFolder(t *testing.T) {
	if !SameFolder("/a/b/budgets.toml", "/a/b/../b/other.toml") {
		t.Error("expected the same folder")
	}
	if SameFolder("/a/b/budgets.toml", "/a/c/budgets.toml") {
		t.Error("expected different folders")
	}
}
