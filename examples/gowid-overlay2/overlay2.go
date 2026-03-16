// Copyright 2019-2022 Graham Clark. All rights reserved.  Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// A demonstration of gowid's overlay, fill, asciigraph and radio widgets.
package main

import (
	"math/rand"

	"github.com/gdamore/tcell/v3"
	"github.com/gnuos/gowid"
	"github.com/gnuos/gowid/examples"
	"github.com/gnuos/gowid/gwutil"
	"github.com/gnuos/gowid/widgets/asciigraph"
	"github.com/gnuos/gowid/widgets/checkbox"
	"github.com/gnuos/gowid/widgets/columns"
	"github.com/gnuos/gowid/widgets/divider"
	"github.com/gnuos/gowid/widgets/framed"
	"github.com/gnuos/gowid/widgets/hpadding"
	"github.com/gnuos/gowid/widgets/overlay"
	"github.com/gnuos/gowid/widgets/pile"
	"github.com/gnuos/gowid/widgets/radio"
	"github.com/gnuos/gowid/widgets/styled"
	"github.com/gnuos/gowid/widgets/text"
	"github.com/gnuos/gowid/widgets/vpadding"
	asc "github.com/guptarohit/asciigraph"
	log "github.com/sirupsen/logrus"
)

//======================================================================

var ov *overlay.Widget
var ovh, ovw int = 50, 50

//======================================================================

type handler struct{}

func (h handler) UnhandledInput(app gowid.IApp, ev any) bool {
	handled := false
	if evk, ok := ev.(*tcell.EventKey); ok {
		handled = true
		if evk.Key() == tcell.KeyCtrlC || evk.Key() == tcell.KeyEsc || evk.Str() == "q" || evk.Str() == "Q" {
			app.Quit()
		} else if evk.Key() == tcell.KeyUp || evk.Str() == "u" {
			ovh = gwutil.Min(100, ovh+1)
			ov.SetHeight(gowid.RenderWithRatio{R: float64(ovh) / 100.0}, app)
		} else if evk.Key() == tcell.KeyDown || evk.Str() == "d" {
			ovh = gwutil.Max(0, ovh-1)
			ov.SetHeight(gowid.RenderWithRatio{R: float64(ovh) / 100.0}, app)
		} else if evk.Key() == tcell.KeyRight {
			ovw = gwutil.Min(100, ovw+1)
			ov.SetWidth(gowid.RenderWithRatio{R: float64(ovw) / 100.0}, app)
		} else if evk.Key() == tcell.KeyLeft {
			ovw = gwutil.Max(0, ovw-1)
			ov.SetWidth(gowid.RenderWithRatio{R: float64(ovw) / 100.0}, app)
		} else {
			handled = false
		}
	}
	return handled
}

//======================================================================

func main() {

	f := examples.RedirectLogger("overlay2.log")
	defer f.Close()

	palette := gowid.Palette{
		"red": gowid.MakePaletteEntry(gowid.ColorRed, gowid.ColorDefault),
	}

	fixed := gowid.RenderFixed{}

	rbgroup := make([]radio.IWidget, 0)
	rb1 := radio.New(&rbgroup)
	rbt1 := text.New(" option1 ")
	rb2 := radio.New(&rbgroup)
	rbt2 := text.New(" option2 ")
	rb3 := radio.New(&rbgroup)
	rbt3 := text.New(" option3 ")

	data := []float64{2, 1, 1, 2, -2, 5, 7, 11, 3, 7, 1, 4, 7, 2, 2, 9}
	data2 := []float64{9, 2, 2, 7, 4, 1, 7, 3, 11, 7, 5, -2, 2, 1, 1, 2}
	conf := []asc.Option{}
	graph := asciigraph.New(data, conf)

	callback := func(app gowid.IApp, target gowid.IWidget) {
		if rb1.IsChecked() {
			graph.SetData(data, app)
		}
		if rb2.IsChecked() {
			graph.SetData(data2, app)
		}
		if rb3.IsChecked() {
			data3 := make([]float64, 40)
			for i := range len(data3) {
				data3[i] = gwutil.Round(rand.Float64() * 14)
			}
			graph.SetData(data3, app)
		}
	}

	rb1.OnClick(gowid.WidgetCallback{Name: gowid.ClickCB{}, WidgetChangedFunction: callback})
	rb2.OnClick(gowid.WidgetCallback{Name: gowid.ClickCB{}, WidgetChangedFunction: callback})
	rb3.OnClick(gowid.WidgetCallback{Name: gowid.ClickCB{}, WidgetChangedFunction: callback})

	c2cols := []gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: rb1, D: fixed},
		&gowid.ContainerWidget{IWidget: rbt1, D: fixed},
		&gowid.ContainerWidget{IWidget: rb2, D: fixed},
		&gowid.ContainerWidget{IWidget: rbt2, D: fixed},
		&gowid.ContainerWidget{IWidget: rb3, D: fixed},
		&gowid.ContainerWidget{IWidget: rbt3, D: fixed},
	}
	cols := columns.New(c2cols)

	rows := pile.New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: cols, D: gowid.RenderWithUnits{U: 1}},
		&gowid.ContainerWidget{IWidget: divider.NewUnicode(), D: gowid.RenderFlow{}},
		&gowid.ContainerWidget{IWidget: graph, D: gowid.RenderWithWeight{W: 1}},
	})

	fcols := framed.NewUnicodeAlt(framed.NewUnicodeAlt(rows))
	top := styled.New(fcols, gowid.MakePaletteRef("red"))
	bottom := vpadding.New(hpadding.New(checkbox.New(false), gowid.HAlignLeft{}, gowid.RenderFixed{}), gowid.VAlignTop{}, gowid.RenderFlow{})

	ov = overlay.New(top, bottom,
		gowid.VAlignMiddle{}, gowid.RenderWithRatio{R: 0.5},
		gowid.HAlignMiddle{}, gowid.RenderWithRatio{R: 0.5})

	app, err := gowid.NewApp(gowid.AppArgs{
		View:    ov,
		Palette: &palette,
		Log:     log.StandardLogger(),
	})
	examples.ExitOnErr(err)

	app.MainLoop(handler{})
}

//======================================================================
// Local Variables:
// mode: Go
// fill-column: 110
// End:
