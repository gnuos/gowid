// Copyright 2019-2022 Graham Clark. All rights reserved.  Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// The second example from the gowid tutorial.
package main

import (
	"fmt"

	"github.com/gdamore/tcell/v3"
	"github.com/gnuos/gowid"
	"github.com/gnuos/gowid/examples"
	"github.com/gnuos/gowid/widgets/text"
)

//======================================================================

var txt *text.Widget

func unhandled(app gowid.IApp, ev any) bool {
	if evk, ok := ev.(*tcell.EventKey); ok {
		switch evk.Str() {
		case "q", "Q":
			app.Quit()
		default:
			txt.SetText(fmt.Sprintf("hello world - %s", evk.Str()), app)
		}
	}
	return true
}

func main() {
	txt = text.New("hello world")
	app, err := gowid.NewApp(gowid.AppArgs{View: txt})
	examples.ExitOnErr(err)
	app.MainLoop(gowid.UnhandledInputFunc(unhandled))
}

//======================================================================
// Local Variables:
// mode: Go
// fill-column: 110
// End:
