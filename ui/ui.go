package ui

import (
	"fmt"
	"image/png"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

// UI Wraps the entire UI
type UI struct {
	progressBar   *widget.ProgressBar
	progressValue float64
	image         *canvas.Image
	box           packr.Box
}

// New creates a new UI
func New() (*UI, error) {
	ui := new(UI)
	ui.box = packr.NewBox("./assets")
	return ui, nil
}

// Start starts a new UI instance
func (ui *UI) Start() error {
	app := app.New()

	w := app.NewWindow("EQGamePatch")
	ui.progressBar = widget.NewProgressBar()
	/*w.SetContent(widget.NewVBox(
		//widget.NewLabel("Hello Fyne!"),
		ui.progressBar,
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))*/
	buttons := widget.NewHBox(
		widget.NewButton("Play", func() {
			fmt.Println("play")
		}),
		layout.NewSpacer(),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	)

	ui.image = &canvas.Image{FillMode: canvas.ImageFillOriginal}
	r, err := ui.box.Open("rof.png")
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", "rof.png")
	}
	ui.image.Image, err = png.Decode(r)
	if err != nil {
		return errors.Wrap(err, "read png")
	}
	r.Close()

	ui.image.Resize(fyne.Size{Width: 400, Height: 450})

	/*w.SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		ui.image,

		fyne.NewContainerWithLayout(layout.NewGridLayout(1)), buttons),
	)*/
	w.SetContent(widget.NewVBox(
		//widget.NewLabel("Hello Fyne!"),
		ui.image,
		ui.progressBar,
		buttons,
	))
	w.Resize(fyne.Size{Width: 410, Height: 505})
	w.SetFixedSize(true)

	go func() {
		for {
			time.Sleep(250 * time.Millisecond)
			if ui.progressValue >= 1 {
				ui.progressValue = 0
			}
			ui.progressValue += 0.01
			ui.progressBar.SetValue(ui.progressValue)
		}
	}()
	w.ShowAndRun()
	return nil
}
