// Copyright 2019-2022 Graham Clark. All rights reserved.  Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// A demonstration of gowid's menu, list and overlay widgets.
package main

import (
	"fmt"

	"github.com/gdamore/tcell/v3"
	"github.com/gnuos/gowid"
	"github.com/gnuos/gowid/examples"
	"github.com/gnuos/gowid/widgets/button"
	"github.com/gnuos/gowid/widgets/checkbox"
	"github.com/gnuos/gowid/widgets/columns"
	"github.com/gnuos/gowid/widgets/hpadding"
	"github.com/gnuos/gowid/widgets/list"
	"github.com/gnuos/gowid/widgets/menu"
	"github.com/gnuos/gowid/widgets/pile"
	"github.com/gnuos/gowid/widgets/styled"
	"github.com/gnuos/gowid/widgets/text"
	"github.com/gnuos/gowid/widgets/vpadding"
	log "github.com/sirupsen/logrus"
)

//======================================================================

var menu1 *menu.Widget
var menu2 *menu.Widget

//======================================================================

type handler struct{}

func (h handler) UnhandledInput(app gowid.IApp, ev any) bool {
	handled := false
	if evk, ok := ev.(*tcell.EventKey); ok {
		handled = true
		if evk.Key() == tcell.KeyCtrlC || evk.Str() == "q" || evk.Str() == "Q" {
			app.Quit()
		} else if evk.Key() == tcell.KeyEsc {
			if !menu1.IsOpen() {
				app.Quit()
			} else {
				menu1.Close(app)
			}
		} else {
			handled = false
		}
	}
	return handled
}

//======================================================================

func main() {
	f := examples.RedirectLogger("menu.log")
	defer f.Close()

	palette := gowid.Palette{
		"red":   gowid.MakePaletteEntry(gowid.ColorRed, gowid.ColorDarkBlue),
		"green": gowid.MakePaletteEntry(gowid.ColorGreen, gowid.ColorDarkBlue),
		"white": gowid.MakePaletteEntry(gowid.ColorWhite, gowid.ColorCyan),
	}

	fixed := gowid.RenderFixed{}

	menu2Widgets := make([]gowid.IWidget, 0)

	for i := range 10 {
		clickme := button.New(text.New(fmt.Sprintf("subwidget %d", i)))
		clickmeStyled := styled.NewInvertedFocus(clickme, gowid.MakePaletteRef("green"))
		clickme.OnClick(gowid.WidgetCallback{Name: gowid.ClickCB{}, WidgetChangedFunction: func(app gowid.IApp, target gowid.IWidget) {
			log.Infof("SUBMENU button CLICKED")
		}})
		cols := columns.New([]gowid.IContainerWidget{
			&gowid.ContainerWidget{IWidget: clickmeStyled, D: fixed},
		})

		menu2Widgets = append(menu2Widgets, cols)
	}

	walker2 := list.NewSimpleListWalker(menu2Widgets)
	menuListBox2 := styled.New(list.New(walker2), gowid.MakePaletteRef("green"))

	menu1Widgets := make([]gowid.IWidget, 0)
	for i := range 40 {
		content := text.NewContent([]text.ContentSegment{
			text.StringContent(fmt.Sprintf("widget %d", i)),
		})
		txt := styled.NewInvertedFocus(text.NewFromContent(content), gowid.MakePaletteRef("red"))
		btn := button.NewBare(txt)
		btnSite := menu.NewSite()
		checkme := checkbox.New(false)
		checkmeStyled := styled.NewInvertedFocus(checkme, gowid.MakePaletteRef("red"))
		checkme.OnClick(gowid.WidgetCallback{Name: gowid.ClickCB{}, WidgetChangedFunction: func(app gowid.IApp, target gowid.IWidget) {
			log.Infof("MENU checkbox CLICKED")
		}})
		btn.OnClick(gowid.WidgetCallback{Name: gowid.ClickCB{}, WidgetChangedFunction: func(app gowid.IApp, target gowid.IWidget) {
			if menu2.IsOpen() {
				menu2.Close(app)
			} else {
				menu2.Open(btnSite, app)
			}
		}})
		cols := columns.New([]gowid.IContainerWidget{
			&gowid.ContainerWidget{IWidget: checkmeStyled, D: fixed},
			&gowid.ContainerWidget{IWidget: btn, D: fixed},
			&gowid.ContainerWidget{IWidget: btnSite, D: fixed},
		})

		menu1Widgets = append(menu1Widgets, cols)
	}

	walker1 := list.NewSimpleListWalker(menu1Widgets)
	menuListBox1 := styled.New(list.New(walker1), gowid.MakePaletteRef("red"))

	menu1 = menu.New("main", menuListBox1, gowid.RenderWithUnits{U: 16})
	menu2 = menu.New("main2", menuListBox2, gowid.RenderWithUnits{U: 16})

	clickToOpenWidgets := make([]gowid.IContainerWidget, 0)
	// Make the on screen buttons to click to open the menu
	for i := range 20 {
		btn := button.New(text.New(fmt.Sprintf("clickety%d", i)))
		btnStyled := styled.NewExt(btn, gowid.MakePaletteRef("red"), gowid.MakePaletteRef("white"))
		btnSite := menu.NewSite(menu.SiteOptions{YOffset: 1})
		btn.OnClick(gowid.WidgetCallback{Name: gowid.ClickCB{}, WidgetChangedFunction: func(app gowid.IApp, target gowid.IWidget) {
			menu1.Open(btnSite, app)
		}})
		clickToOpenWidgets = append(clickToOpenWidgets, &gowid.ContainerWidget{IWidget: btnSite, D: fixed})
		clickToOpenWidgets = append(clickToOpenWidgets, &gowid.ContainerWidget{IWidget: btnStyled, D: fixed})
	}
	clickToOpenCols := columns.New(clickToOpenWidgets)

	check := checkbox.New(false)

	view1 := pile.New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: clickToOpenCols, D: fixed},
		&gowid.ContainerWidget{IWidget: check, D: fixed},
	})

	view := vpadding.New(
		hpadding.New(view1, gowid.HAlignLeft{}, fixed),
		gowid.VAlignTop{Margin: 2},
		gowid.RenderFlow{},
	)

	app, err := gowid.NewApp(gowid.AppArgs{
		View:    view,
		Palette: &palette,
		Log:     log.StandardLogger(),
	})
	examples.ExitOnErr(err)

	// Required for menus to appear overlaid on top of the main view.
	app.RegisterMenu(menu1)
	app.RegisterMenu(menu2)

	app.MainLoop(handler{})
}

//======================================================================
// Local Variables:
// mode: Go
// fill-column: 110
// End:
