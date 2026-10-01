// Command sofar-sogood is a small desktop widget showing how much of each
// configured budget should have been spent by today.
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/sfkleach/sofar-sogood/budget"
)

func main() {
	path := flag.String("config", "", "path to budgets.yaml (default: user config dir)")
	flag.Parse()
	if *path == "" {
		p, err := budget.DefaultPath()
		if err != nil {
			log.Fatal(err)
		}
		*path = p
	}

	a := app.New()
	w := a.NewWindow("So far, so good")

	lines := container.NewVBox()
	render := func() {
		lines.Objects = nil
		budgets, err := budget.Load(*path)
		switch {
		case err != nil:
			lines.Add(widget.NewLabel(fmt.Sprintf("config error: %v", err)))
		case len(budgets) == 0:
			lines.Add(widget.NewLabel("no budgets"))
		default:
			now := time.Now()
			for _, b := range budgets {
				lines.Add(widget.NewLabel(b.Estimate(now).Line()))
			}
		}
		lines.Refresh()
		w.Resize(fyne.NewSize(w.Content().MinSize().Width, w.Content().MinSize().Height))
	}

	info := widget.NewLabel(fmt.Sprintf(
		"So far, so good helps you track AI credit usage against a budget.\n\n"+
			"Each budget shows how much you should have used by the close of play "+
			"today, spreading the renewal amount evenly over the working days "+
			"of the period.\n\n"+
			"Budgets are defined in this config file:\n%s\n\n"+
			"Each amount is a number with an optional symbol or unit, "+
			"such as \"$200\", \"£1,000\" or \"800 credits\".\n\n"+
			"(Pass -config <path> to use a different file.)", *path))
	info.Wrapping = fyne.TextWrapWord

	showingInfo := false
	var infoBtn *widget.Button
	body := container.NewStack(lines)
	infoBtn = widget.NewButton("Info", func() {
		showingInfo = !showingInfo
		if showingInfo {
			body.Objects = []fyne.CanvasObject{info}
			infoBtn.SetText("Back")
		} else {
			body.Objects = []fyne.CanvasObject{lines}
			infoBtn.SetText("Info")
		}
		body.Refresh()
		w.Resize(fyne.NewSize(max(w.Content().MinSize().Width, 360), w.Content().MinSize().Height))
	})
	refresh := widget.NewButton("Refresh", render)
	dismiss := widget.NewButton("Dismiss", a.Quit)
	w.SetContent(container.NewVBox(body, container.NewGridWithColumns(3, refresh, infoBtn, dismiss)))
	render()
	w.SetFixedSize(false)
	w.ShowAndRun()
}
