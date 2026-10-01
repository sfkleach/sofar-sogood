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

	refresh := widget.NewButton("Refresh", render)
	dismiss := widget.NewButton("Dismiss", a.Quit)
	w.SetContent(container.NewVBox(lines, container.NewGridWithColumns(2, refresh, dismiss)))
	render()
	w.SetFixedSize(false)
	w.ShowAndRun()
}
