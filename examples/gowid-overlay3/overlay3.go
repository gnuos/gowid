// Copyright 2019-2022 Graham Clark. All rights reserved.  Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// A demonstration of gowid's overlay and fill widgets.
package main

import (
	"github.com/gdamore/tcell/v3"
	"github.com/gnuos/gowid"
	"github.com/gnuos/gowid/examples"
	"github.com/gnuos/gowid/widgets/fill"
	"github.com/gnuos/gowid/widgets/framed"
	"github.com/gnuos/gowid/widgets/overlay"
	"github.com/gnuos/gowid/widgets/styled"
	"github.com/gnuos/gowid/widgets/text"
	log "github.com/sirupsen/logrus"
)

//======================================================================

var ov *overlay.Widget

//======================================================================

type handler struct{}

func (h handler) UnhandledInput(app gowid.IApp, ev any) bool {
	handled := false
	if evk, ok := ev.(*tcell.EventKey); ok {
		handled = true
		if evk.Key() == tcell.KeyCtrlC || evk.Key() == tcell.KeyEsc || evk.Str() == "q" || evk.Str() == "Q" {
			app.Quit()
		} else {
			handled = false
		}
	}
	return handled
}

//======================================================================

func main() {

	f := examples.RedirectLogger("overlay1.log")
	defer f.Close()

	palette := gowid.Palette{
		"red": gowid.MakePaletteEntry(gowid.ColorDefault, gowid.ColorRed),
	}

	top := styled.New(
		framed.NewUnicode(
			text.New("hello"),
		),
		gowid.MakePaletteRef("red"),
	)
	bottom := fill.New(' ')

	ov = overlay.New(top, bottom,
		gowid.VAlignMiddle{}, gowid.RenderFixed{},
		gowid.HAlignMiddle{}, gowid.RenderFixed{},
	)

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
