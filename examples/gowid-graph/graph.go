// Copyright 2019-2022 Graham Clark. All rights reserved.  Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// A port of urwid's graph.py example using gowid widgets.
package main

import (
	"math"
	"time"

	"github.com/gnuos/gowid"
	"github.com/gnuos/gowid/examples"
	"github.com/gnuos/gowid/gwutil"
	"github.com/gnuos/gowid/widgets/bargraph"
	"github.com/gnuos/gowid/widgets/button"
	"github.com/gnuos/gowid/widgets/clicktracker"
	"github.com/gnuos/gowid/widgets/columns"
	"github.com/gnuos/gowid/widgets/divider"
	"github.com/gnuos/gowid/widgets/fill"
	"github.com/gnuos/gowid/widgets/framed"
	"github.com/gnuos/gowid/widgets/hpadding"
	"github.com/gnuos/gowid/widgets/pile"
	"github.com/gnuos/gowid/widgets/progress"
	"github.com/gnuos/gowid/widgets/radio"
	"github.com/gnuos/gowid/widgets/shadow"
	"github.com/gnuos/gowid/widgets/styled"
	"github.com/gnuos/gowid/widgets/text"
	log "github.com/sirupsen/logrus"
)

//======================================================================

var app *gowid.App
var controller *GraphController

//======================================================================

func sin100(x int) int {
	return int(50 + 50*math.Sin(float64(x)*math.Pi/50.0))
}

func round(f float64) float64 {
	return math.Floor(f + 0.5)
}

func sum(x ...int) int {
	res := 0
	for _, i := range x {
		res += i
	}
	return res
}

//======================================================================

type GraphModel struct {
	Data        map[string][]int
	Modes       []string
	CurrentMode string
}

func NewGraphModel() *GraphModel {
	modes := make([]string, 0)
	data := make(map[string][]int)
	var a1 []int
	for i := range 50 {
		a1 = append(a1, i*2)
	}
	a2 := append(a1, a1...)
	data["Saw"] = a2
	modes = append(modes, "Saw")

	var a3 []int
	var a4 []int
	for range 30 {
		a3 = append(a3, 0)
		a4 = append(a4, 100)
	}
	data["Square"] = append(a3, a4...)
	modes = append(modes, "Square")

	var a5 []int
	for i := range 100 {
		a5 = append(a5, sin100(i))
	}
	data["Sine 1"] = a5
	modes = append(modes, "Sine 1")

	var a6 []int
	for i := range 100 {
		a6 = append(a6, (sin100(i)+sin100(i*2))/2)
	}
	data["Sine 2"] = a6
	modes = append(modes, "Sine 2")

	var a7 []int
	for i := range 100 {
		a7 = append(a7, (sin100(i)+sin100(i*3))/2)
	}
	data["Sine 3"] = a7
	modes = append(modes, "Sine 3")

	return &GraphModel{data, modes, modes[0]}
}

func (g *GraphModel) SetMode(mode string) {
	g.CurrentMode = mode
}

func (g *GraphModel) GetData(offset, r int) ([]int, int, int) {
	l := make([]int, 0)
	d := g.Data[g.CurrentMode]
	for r > 0 {
		offset = offset % len(d)
		segment := d[offset:gwutil.Min(offset+r, len(d))]
		if len(segment) == 0 {
			break
		}
		r = r - len(segment)
		offset += len(segment)
		l = append(l, segment...)
	}
	return l, 100, len(d)
}

//======================================================================

type GraphView struct {
	*styled.Widget
	controller *GraphController
	started    bool
	startTime  *time.Time
	offset     int
	lastOffset *int
	graph      *bargraph.Widget
	pb         *progress.Widget
}

func NewGraphView(controller *GraphController) *GraphView {
	t := time.Now()

	pb := progress.New(progress.Options{
		Normal:   gowid.MakePaletteRef("pg normal"),
		Complete: gowid.MakePaletteRef("pg complete"),
	})
	graph := MakeBarGraph()
	controls := MakeBarGraphControls(controller, pb)

	weight1 := gowid.RenderWithWeight{W: 1}
	weight2 := gowid.RenderWithWeight{W: 2}
	unit1 := gowid.RenderWithUnits{U: 1}

	vline := styled.New(fill.New('│'), gowid.MakePaletteRef("line"))

	view := styled.New(
		framed.NewSpace(
			shadow.New(
				styled.New(
					framed.NewUnicode(
						columns.New([]gowid.IContainerWidget{
							&gowid.ContainerWidget{IWidget: graph, D: weight2},
							&gowid.ContainerWidget{IWidget: vline, D: unit1},
							&gowid.ContainerWidget{IWidget: controls, D: weight1},
						}),
					),
					gowid.MakePaletteRef("body"),
				),
				1,
			),
		),
		gowid.MakePaletteRef("screen edge"),
	)

	res := &GraphView{
		Widget:     view,
		controller: controller,
		startTime:  &t,
		graph:      graph,
		pb:         pb,
	}
	return res
}

func (g *GraphView) Selectable() bool {
	return true
}

func (g *GraphView) UserInput(ev any, size gowid.IRenderSize, focus gowid.Selector, app gowid.IApp) bool {
	return g.Widget.UserInput(ev, size, focus, app)
}

func MakeBarGraph() *bargraph.Widget {
	w := bargraph.New([]gowid.IColor{
		gowid.NewUrwidColor("dark gray"),
		gowid.NewUrwidColor("dark blue"),
		gowid.NewUrwidColor("dark cyan"),
	})

	return w
}

func MakeBarGraphControls(controller *GraphController, pb progress.IWidget) *pile.Widget {

	modeButtons := make([]gowid.IContainerWidget, 0)
	rbgroup := make([]radio.IWidget, 0)
	p := gowid.RenderFixed{}
	f := gowid.RenderFlow{}
	var firstrb *radio.Widget

	for _, mode := range controller.model.Modes {
		capturedMode := mode
		rb1 := radio.New(&rbgroup)
		if firstrb == nil {
			firstrb = rb1
		}
		rbt1 := text.New(" " + mode)
		rb1.OnClick(gowid.WidgetCallback{Name: "cb",
			WidgetChangedFunction: func(app gowid.IApp, w gowid.IWidget) {
				controller.model.SetMode(capturedMode)
				controller.view.UpdateGraph(true, app)
				controller.view.lastOffset = nil
			}})
		modeButton := make([]gowid.IContainerWidget, 0)
		modeButton = append(modeButton, &gowid.ContainerWidget{IWidget: rb1, D: p})
		modeButton = append(modeButton, &gowid.ContainerWidget{IWidget: rbt1, D: p})
		modeButtonCols := styled.NewExt(columns.New(modeButton),
			gowid.MakePaletteRef("button normal"),
			gowid.MakePaletteRef("button select"))
		modeButtons = append(modeButtons, &gowid.ContainerWidget{IWidget: modeButtonCols, D: f})
	}

	animateText := text.New("Start")
	animateButton := button.New(animateText)
	resetText := text.New("Reset")
	resetButton := button.New(resetText)
	quitText := text.New("Quit")
	quitButton := button.New(quitText)

	animateButton.OnClick(gowid.WidgetCallback{Name: "cb",
		WidgetChangedFunction: func(app gowid.IApp, w gowid.IWidget) {
			if animateText.Content().Length() == 5 {
				controller.AnimateGraph(app)
				animateText.SetText("Stop", app)
			} else {
				controller.StopAnimation()
				animateText.SetText("Start", app)
			}
		}})

	resetButton.OnClick(gowid.WidgetCallback{Name: "cb",
		WidgetChangedFunction: func(app gowid.IApp, w gowid.IWidget) {
			controller.ResetGraph(app)
		}})

	quitButton.OnClick(gowid.WidgetCallback{Name: "cb",
		WidgetChangedFunction: func(app gowid.IApp, w gowid.IWidget) {
			app.Quit()
		}})

	animateButtonStyled := styled.NewExt(animateButton,
		gowid.MakePaletteRef("button normal"),
		gowid.MakePaletteRef("button select"))
	resetButtonStyled := styled.NewExt(resetButton,
		gowid.MakePaletteRef("button normal"),
		gowid.MakePaletteRef("button select"))
	quitButtonStyled := styled.NewExt(quitButton,
		gowid.MakePaletteRef("button normal"),
		gowid.MakePaletteRef("button select"))

	animateButtonTracker := clicktracker.New(animateButtonStyled)
	resetButtonTracker := clicktracker.New(resetButtonStyled)
	quitButtonTracker := clicktracker.New(quitButtonStyled)

	buttonGrid := columns.New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: hpadding.New(animateButtonTracker, gowid.HAlignMiddle{}, p), D: gowid.RenderWithWeight{W: 1}},
		&gowid.ContainerWidget{IWidget: hpadding.New(resetButtonTracker, gowid.HAlignMiddle{}, p), D: gowid.RenderWithWeight{W: 1}},
	})

	controls := make([]gowid.IContainerWidget, 0, 7+len(modeButtons))
	controls = append(controls, modeButtons...)
	controls = append(controls, &gowid.ContainerWidget{IWidget: divider.NewBlank(), D: f})
	controls = append(controls, &gowid.ContainerWidget{IWidget: hpadding.New(text.New("Animation"), gowid.HAlignMiddle{}, p), D: f})
	controls = append(controls, &gowid.ContainerWidget{IWidget: buttonGrid, D: f})
	controls = append(controls, &gowid.ContainerWidget{IWidget: divider.NewBlank(), D: f})
	controls = append(controls, &gowid.ContainerWidget{IWidget: pb, D: f})
	controls = append(controls, &gowid.ContainerWidget{IWidget: divider.NewBlank(), D: f})
	controls = append(controls, &gowid.ContainerWidget{IWidget: hpadding.New(quitButtonTracker, gowid.HAlignMiddle{}, p), D: f})

	controlsPile := pile.New(controls)

	return controlsPile
}

func (v *GraphView) GetOffsetNow() int {
	if v.startTime == nil {
		return 0
	}
	if !v.started {
		return v.offset
	}
	tdelta := time.Since(*v.startTime)
	tdelta = tdelta * 5
	x := v.offset + (int(round(tdelta.Seconds())))
	return x
}

func (v *GraphView) UpdateGraph(forceUpdate bool, app gowid.IApp) bool {
	o := v.GetOffsetNow()
	if v.lastOffset != nil && o == *v.lastOffset && !forceUpdate {
		return false
	}
	v.lastOffset = &o
	gspb := 10
	r := gspb * 5
	d, maxValue, repeat := v.controller.GetData(o, r)
	l := make([][]int, 0)
	for n := range 5 {
		value := sum(d[n*gspb:(n+1)*gspb]...) / gspb
		// toggle between two bar types
		if n&1 == 1 {
			l = append(l, []int{0, value})
		} else {
			l = append(l, []int{value, 0})
		}
	}
	v.graph.SetData(l, maxValue, app)

	var prog int
	// also update progress
	if (o/repeat)&1 == 1 {
		// show 100% for first half, 0 for second half
		if o%repeat > repeat {
			prog = 0
		} else {
			prog = 100
		}
	} else {
		prog = ((o % repeat) * 100) / repeat
	}
	v.pb.SetProgress(app, prog)

	return true
}

//======================================================================

type GraphController struct {
	model  *GraphModel
	view   *GraphView
	mode   string
	ticker *time.Ticker
}

func NewGraphController() *GraphController {
	res := &GraphController{NewGraphModel(), nil, "", nil}
	view := NewGraphView(res)
	res.view = view
	res.mode = res.model.Modes[0]
	return res
}

func (g *GraphController) GetData(offset, r int) ([]int, int, int) {
	return g.model.GetData(offset, r)
}

func (g *GraphController) ResetGraph(app gowid.IApp) {
	t := time.Now()
	g.view.startTime = &t
	g.view.offset = 0
	g.view.UpdateGraph(true, app)
}

func (g *GraphController) AnimateGraph(app gowid.IApp) {
	t := time.Now()
	g.view.startTime = &t
	g.ticker = time.NewTicker(time.Millisecond * 200)
	g.view.started = true
	go func() {
		for range g.ticker.C {
			app.Run(gowid.RunFunction(func(app gowid.IApp) {
				g.view.UpdateGraph(true, app)
				app.Redraw()
			}))
		}
	}()
}

func (g *GraphController) StopAnimation() {
	g.ticker.Stop()
	g.view.offset = g.view.GetOffsetNow()
	g.view.started = false
}

//======================================================================

func main() {
	var err error

	f := examples.RedirectLogger("graph.log")
	defer f.Close()

	palette := gowid.Palette{
		"body":          gowid.MakeStyledPaletteEntry(gowid.NewUrwidColor("black"), gowid.NewUrwidColor("light gray"), gowid.StyleBold),
		"line":          gowid.MakePaletteRef("body"),
		"button normal": gowid.MakeStyledPaletteEntry(gowid.NewUrwidColor("light gray"), gowid.NewUrwidColor("dark blue"), gowid.StyleBold),
		"button select": gowid.MakePaletteEntry(gowid.NewUrwidColor("white"), gowid.NewUrwidColor("dark green")),
		"pg normal":     gowid.MakeStyledPaletteEntry(gowid.NewUrwidColor("white"), gowid.NewUrwidColor("black"), gowid.StyleBold),
		"pg complete":   gowid.MakeStyleMod(gowid.MakePaletteRef("pg normal"), gowid.MakeBackground(gowid.NewUrwidColor("dark magenta"))),
		"screen edge":   gowid.MakePaletteEntry(gowid.NewUrwidColor("light blue"), gowid.NewUrwidColor("dark cyan")),
	}

	controller = NewGraphController()

	app, err = gowid.NewApp(gowid.AppArgs{
		View:    controller.view,
		Palette: &palette,
		Log:     log.StandardLogger(),
	})
	examples.ExitOnErr(err)

	controller.ResetGraph(app)

	app.SimpleMainLoop()
}

//======================================================================
// Local Variables:
// mode: Go
// fill-column: 110
// End:
