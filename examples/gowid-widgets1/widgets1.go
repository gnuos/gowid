// Copyright 2019-2022 Graham Clark. All rights reserved.  Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// A gowid test app which exercises the pile, button, edit and progress widgets.
package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gnuos/gowid"
	"github.com/gnuos/gowid/examples"
	"github.com/gnuos/gowid/gwutil"
	"github.com/gnuos/gowid/widgets/button"
	"github.com/gnuos/gowid/widgets/divider"
	"github.com/gnuos/gowid/widgets/edit"
	"github.com/gnuos/gowid/widgets/framed"
	"github.com/gnuos/gowid/widgets/holder"
	"github.com/gnuos/gowid/widgets/palettemap"
	"github.com/gnuos/gowid/widgets/pile"
	"github.com/gnuos/gowid/widgets/progress"
	"github.com/gnuos/gowid/widgets/styled"
	"github.com/gnuos/gowid/widgets/text"
	"github.com/gnuos/gowid/widgets/vpadding"
	"github.com/sirupsen/logrus"
)

//======================================================================
// An example of how to override

type PBWidget struct {
	*progress.Widget
}

func NewPB() *PBWidget {
	return &PBWidget{progress.New(progress.Options{
		Normal:   gowid.MakeEmptyPalette(),
		Complete: gowid.MakePaletteRef("invred"),
	})}
}

func (w *PBWidget) Text() string {
	cur, done := w.Progress(), w.Target()
	percent := gwutil.Min(100, gwutil.Max(0, cur*100/done))
	return fmt.Sprintf("At %d %% (%d/%d)", percent, cur, done)
}

func (w *PBWidget) Render(size gowid.IRenderSize, focus gowid.Selector, app gowid.IApp) gowid.ICanvas {
	return progress.Render(w, size, focus, app)
}

//======================================================================

type handler struct{}

func (h handler) UnhandledInput(app gowid.IApp, ev any) bool {
	if evk, ok := ev.(*tcell.EventKey); ok {
		if evk.Key() == tcell.KeyCtrlC || evk.Str() == "q" || evk.Str() == "Q" {
			app.Quit()
			return true
		}
	}
	return false
}

//======================================================================

func main() {

	f := examples.RedirectLogger("widgets1.log")
	defer f.Close()

	styles := gowid.Palette{
		"banner":        gowid.MakePaletteEntry(gowid.ColorBlack, gowid.ColorWhite),
		"streak":        gowid.MakePaletteEntry(gowid.ColorBlack, gowid.ColorRed),
		"bg":            gowid.MakePaletteEntry(gowid.ColorBlack, gowid.ColorBlue),
		"test1focus":    gowid.MakePaletteEntry(gowid.ColorBlue, gowid.ColorBlack),
		"test1notfocus": gowid.MakePaletteEntry(gowid.ColorGreen, gowid.ColorBlack),
		"red":           gowid.MakePaletteEntry(gowid.ColorRed, gowid.ColorBlack),
		"invred":        gowid.MakePaletteEntry(gowid.ColorBlack, gowid.ColorRed),
		"magenta":       gowid.MakePaletteEntry(gowid.ColorMagenta, gowid.ColorBlack),
		"cyan":          gowid.MakePaletteEntry(gowid.ColorCyan, gowid.ColorBlack),
	}

	flowme := gowid.RenderFlow{}
	pb1 := NewPB()

	nl := gowid.MakePaletteRef

	mh := text.NewContent([]text.ContentSegment{
		text.StyledContent("abc", nl("invred")),
		text.StringContent("def"),
		text.StyledContent("ghijk", nl("cyan")),
		text.StyledContent("lmnopq", nl("magenta")),
	})

	mti := text.NewFromContent(mh)
	mtj := palettemap.New(mti, palettemap.Map{}, palettemap.Map{"invred": "red"})
	mt := holder.New(mti)

	xt := text.New("something else")
	xt2 := styled.New(xt, gowid.MakePaletteEntry(gowid.NewUrwidColor("dark red"), gowid.NewUrwidColor("light red")))

	tw1 := text.New("click me or double-click me█ █xx")
	tw := styled.NewWithRanges(tw1,
		[]styled.AttributeRange{{Start: 0, End: 2, Styler: nl("test1notfocus")}},
		[]styled.AttributeRange{{Start: 0, End: -1, Styler: nl("test1focus")}},
	)

	bw1i := button.New(tw, button.Options{
		Decoration:       button.NormalDecoration,
		DoubleClickDelay: 200 * time.Millisecond,
	})
	bw1 := holder.New(bw1i)

	dv1 := divider.NewAscii()
	e1e := edit.New(edit.Options{Caption: "Name:", Text: "(1)abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzab(2)CDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCD(3)efghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdef(4)&^&^&^&^&GHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZ"})
	e1 := vpadding.New(e1e, gowid.VAlignTop{}, gowid.RenderWithUnits{U: 2})

	e2 := edit.New(edit.Options{Caption: "Password:", Text: "foobar", Mask: edit.MakeMask('*')})

	e2e := edit.New(edit.Options{Caption: "Domain:", Text: "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"})

	bw1i.OnClick(gowid.WidgetCallback{Name: "cb", WidgetChangedFunction: func(app gowid.IApp, w gowid.IWidget) {
		pb1.SetProgress(app, pb1.Progress()+1)
		if mt.SubWidget() == mti {
			mt.SetSubWidget(mtj, app)
		} else {
			mt.SetSubWidget(mti, app)
		}
	}})

	bw1i.OnDoubleClick(gowid.WidgetCallback{Name: "cb",
		WidgetChangedFunction: func(app gowid.IApp, w gowid.IWidget) {
			logrus.Infof("GCLA: got double click")
			pb1.SetProgress(app, 0)
		}})

	pw := pile.New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: pb1, D: flowme},
		&gowid.ContainerWidget{IWidget: dv1, D: flowme},
		&gowid.ContainerWidget{IWidget: bw1, D: flowme},
		&gowid.ContainerWidget{IWidget: dv1, D: flowme},
		&gowid.ContainerWidget{IWidget: e1, D: flowme},
		&gowid.ContainerWidget{IWidget: dv1, D: flowme},
		&gowid.ContainerWidget{IWidget: e2, D: flowme},
		&gowid.ContainerWidget{IWidget: dv1, D: flowme},
		&gowid.ContainerWidget{IWidget: e2e, D: flowme},
		&gowid.ContainerWidget{IWidget: dv1, D: flowme},
		&gowid.ContainerWidget{IWidget: mt, D: flowme},
		&gowid.ContainerWidget{IWidget: dv1, D: flowme},
		&gowid.ContainerWidget{IWidget: xt2, D: flowme},
	})
	twi := styled.New(text.New(" widgets1 "), gowid.MakePaletteRef("magenta"))
	params := framed.Options{
		TitleWidget: twi,
	}
	fw1 := framed.New(pw, params)
	fw := framed.NewUnicode(fw1)
	pw2 := vpadding.New(fw, gowid.VAlignMiddle{}, gowid.RenderFlow{})

	app, err := gowid.NewApp(gowid.AppArgs{
		View:    pw2,
		Palette: &styles,
		Log:     logrus.StandardLogger(),
	})
	examples.ExitOnErr(err)

	app.MainLoop(handler{})
}

//======================================================================
// Local Variables:
// mode: Go
// fill-column: 110
// End:
