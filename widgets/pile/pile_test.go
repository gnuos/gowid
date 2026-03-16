// Copyright 2019-2022 Graham Clark. All rights reserved.  Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package pile

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v3"
	"github.com/gnuos/gowid"
	"github.com/gnuos/gowid/gwtest"
	"github.com/gnuos/gowid/widgets/button"
	"github.com/gnuos/gowid/widgets/fill"
	"github.com/gnuos/gowid/widgets/framed"
	"github.com/gnuos/gowid/widgets/list"
	"github.com/gnuos/gowid/widgets/selectable"
	"github.com/gnuos/gowid/widgets/text"
	"github.com/stretchr/testify/assert"
)

// Test that a mouse wheel down inside the region of the list within a pile
// is correctly translated and passed to the list.
func TestPile5(t *testing.T) {
	bws := make([]gowid.IWidget, 50)
	for i := range len(bws) {
		bws[i] = button.New(text.New(fmt.Sprintf("%03d", i)))
	}

	walker := list.NewSimpleListWalker(bws)
	lb := list.New(walker)
	// Framed is needed because it validates the mouse y coordinate before passing it on
	// to its subwidget.
	flb := framed.New(lb)

	pws := make([]gowid.IContainerWidget, 3)
	pws[0] = &gowid.ContainerWidget{IWidget: button.New(text.New("top  ")), D: gowid.RenderWithUnits{U: 1}}
	pws[1] = &gowid.ContainerWidget{IWidget: flb, D: gowid.RenderWithWeight{W: 1}} // WEIGHT!!
	pws[2] = &gowid.ContainerWidget{IWidget: button.New(text.New("bot  ")), D: gowid.RenderWithUnits{U: 1}}

	pl := New(pws)

	sta := make([]string, 0)
	sta = append(sta, "<top  >")
	sta = append(sta, "-------")
	for i := range 6 {
		sta = append(sta, fmt.Sprintf("|<%03d>|", i))
	}
	sta = append(sta, "-------")
	sta = append(sta, "<bot  >")

	csize := gowid.RenderBox{C: 7, R: 10}
	c := pl.Render(csize, gowid.Focused, gwtest.D)
	assert.Equal(t, strings.Join(sta, "\n"), c.String())

	assert.Equal(t, 0, pl.Focus())

	evmdown := tcell.NewEventMouse(1, 4, tcell.WheelDown, 0)

	pl.UserInput(evmdown, csize, gowid.Focused, gwtest.D)
	assert.Equal(t, 1, pl.Focus()) // Now at list widget
	assert.Equal(t, 0, lb.Walker().Focus().(list.ListPos).ToInt())

	pl.UserInput(evmdown, csize, gowid.Focused, gwtest.D)
	assert.Equal(t, 1, pl.Focus()) // Now at list widget
	assert.Equal(t, 1, lb.Walker().Focus().(list.ListPos).ToInt())

	for range 40 {
		pl.UserInput(evmdown, csize, gowid.Focused, gwtest.D)
	}

	assert.Equal(t, 1, pl.Focus()) // Now at list widget
	assert.Equal(t, 41, lb.Walker().Focus().(list.ListPos).ToInt())
}

func TestPile2(t *testing.T) {
	btns := make([]gowid.IContainerWidget, 0)
	//clicks := make([]*gwtest.ButtonTester, 0)

	for range 3 {
		btn := button.New(text.New("abc"))
		click := &gwtest.ButtonTester{Gotit: false}
		btn.OnClick(click)
		btns = append(btns, &gowid.ContainerWidget{IWidget: btn, D: gowid.RenderFixed{}})
		//clicks = append(clicks, click)
	}

	pl := New(btns)

	st1 := "<abc>\n<abc>\n<abc>"
	st2 := "<abc> \n<abc> \n<abc> "

	c := pl.Render(gowid.RenderFixed{}, gowid.Focused, gwtest.D)
	assert.Equal(t, c.String(), st1)

	c2 := pl.Render(gowid.RenderFlowWith{C: 6}, gowid.Focused, gwtest.D)
	assert.Equal(t, c2.String(), st2)

	assert.Equal(t, pl.Focus(), 0)

	// evright := tcell.NewEventKey(tcell.KeyRight, ' ', tcell.ModNone)
	// evleft := tcell.NewEventKey(tcell.KeyLeft, ' ', tcell.ModNone)
	// evdown := tcell.NewEventKey(tcell.KeyDown, ' ', tcell.ModNone)
	// evspace := tcell.NewEventKey(tcell.KeyRune, ' ', tcell.ModNone)
	evmdown := tcell.NewEventMouse(1, 1, tcell.WheelDown, 0)
	evmup := tcell.NewEventMouse(1, 1, tcell.WheelUp, 0)
	// evmright := tcell.NewEventMouse(1, 1, tcell.WheelRight, 0)
	// evmleft := tcell.NewEventMouse(1, 1, tcell.WheelLeft, 0)

	cbcalled := false

	pl.OnFocusChanged(gowid.WidgetCallback{Name: "cb",
		WidgetChangedFunction: func(app gowid.IApp, w gowid.IWidget) {
			assert.Equal(t, w, pl)
			cbcalled = true
		}})

	pl.UserInput(evmdown, gowid.RenderFixed{}, gowid.Focused, gwtest.D)
	assert.Equal(t, 1, pl.Focus())
	assert.Equal(t, true, cbcalled)
	cbcalled = false
	pl.UserInput(evmdown, gowid.RenderFixed{}, gowid.Focused, gwtest.D)
	assert.Equal(t, 2, pl.Focus())
	assert.Equal(t, true, cbcalled)
	cbcalled = false
	pl.UserInput(evmdown, gowid.RenderFixed{}, gowid.Focused, gwtest.D)
	assert.Equal(t, 2, pl.Focus())
	assert.Equal(t, false, cbcalled)

	pl.UserInput(evmup, gowid.RenderFixed{}, gowid.Focused, gwtest.D)
	assert.Equal(t, 1, pl.Focus())
	assert.Equal(t, true, cbcalled)
	cbcalled = false
	pl.UserInput(evmup, gowid.RenderFixed{}, gowid.Focused, gwtest.D)
	assert.Equal(t, 0, pl.Focus())
	assert.Equal(t, true, cbcalled)
	cbcalled = false
	pl.UserInput(evmup, gowid.RenderFixed{}, gowid.Focused, gwtest.D)
	assert.Equal(t, 0, pl.Focus())
	assert.Equal(t, false, cbcalled)
}

func TestPile1(t *testing.T) {
	w1 := New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: fill.New('x'), D: gowid.RenderWithUnits{U: 2}},
		&gowid.ContainerWidget{IWidget: fill.New('y'), D: gowid.RenderWithUnits{U: 2}},
	})
	c1 := w1.Render(gowid.RenderBox{C: 3, R: 4}, gowid.Focused, gwtest.D)
	assert.Equal(t, c1.String(), "xxx\nxxx\nyyy\nyyy")

	w2 := New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: fill.New('x'), D: gowid.RenderWithUnits{U: 1}},
		&gowid.ContainerWidget{IWidget: fill.New('y'), D: gowid.RenderWithUnits{U: 2}},
	})
	c2 := w2.Render(gowid.RenderFlowWith{C: 3}, gowid.Focused, gwtest.D)
	assert.Equal(t, c2.String(), "xxx\nyyy\nyyy")

	w3 := New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: fill.New('x'), D: gowid.RenderWithWeight{W: 1}},
		&gowid.ContainerWidget{IWidget: fill.New('y'), D: gowid.RenderWithWeight{W: 2}},
	})
	assert.Panics(t, func() {
		w3.Render(gowid.RenderFlowWith{C: 3}, gowid.Focused, gwtest.D)
	})

	w4 := New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: fill.New('x'), D: gowid.RenderWithRatio{R: 0.25}},
		&gowid.ContainerWidget{IWidget: fill.New('y'), D: gowid.RenderWithRatio{R: 0.5}},
	})

	c4 := w4.Render(gowid.RenderBox{C: 3, R: 3}, gowid.Focused, gwtest.D)
	assert.Equal(t, c4.String(), "xxx\nyyy\nyyy")

	c41 := w4.Render(gowid.RenderBox{C: 3, R: 4}, gowid.Focused, gwtest.D)
	assert.Equal(t, c41.String(), "xxx\nyyy\nyyy\n   ")

	for _, w := range []gowid.IWidget{w1, w2, w4} {
		gwtest.RenderBoxManyTimes(t, w, 0, 10, 0, 10)
	}
	gwtest.RenderFlowManyTimes(t, w2, 0, 10)
}

func TestPile3(t *testing.T) {
	w1 := New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: fill.New('x'), D: gowid.RenderWithUnits{U: 2}},
		&gowid.ContainerWidget{IWidget: text.New("y"), D: gowid.RenderFlow{}},
	})
	// Test that a pile can render in flow mode with a single embedded flow widget
	c1 := w1.Render(gowid.RenderFlowWith{C: 3}, gowid.Focused, gwtest.D)
	assert.Equal(t, c1.String(), "xxx\nxxx\ny  ")

	w1 = New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: fill.New('x'), D: gowid.RenderWithUnits{U: 2}},
		&gowid.ContainerWidget{IWidget: text.New("y"), D: gowid.RenderWithWeight{W: 1}},
	})
	// Test that a pile can render in flow mode with a single embedded flow widget
	c1 = w1.Render(gowid.RenderFlowWith{C: 3}, gowid.Focused, gwtest.D)
	assert.Equal(t, c1.String(), "xxx\nxxx\ny  ")

	w1 = New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: fill.New('x'), D: gowid.RenderWithUnits{U: 2}},
		&gowid.ContainerWidget{IWidget: text.New("y"), D: gowid.RenderWithWeight{W: 1}},
		&gowid.ContainerWidget{IWidget: text.New("z"), D: gowid.RenderWithWeight{W: 1}},
	})
	// Two weight widgets don't work in flow mode, how do you restrict their vertical ratio?
	assert.Panics(t, func() {
		w1.Render(gowid.RenderFlowWith{C: 3}, gowid.Focused, gwtest.D)
	})

}

func makep(c rune) gowid.IWidget {
	return selectable.New(fill.New(c))
}

func makepfixed(c rune) gowid.IContainerWidget {
	return &gowid.ContainerWidget{
		IWidget: makep(c),
		D:       gowid.RenderFixed{},
	}
}

type renderWeightUpTo struct {
	gowid.RenderWithWeight
	max int
}

func (s renderWeightUpTo) MaxUnits() int {
	return s.max
}

func weightupto(w int, max int) renderWeightUpTo {
	return renderWeightUpTo{gowid.RenderWithWeight{W: w}, max}
}

func TestPile4(t *testing.T) {
	subs := []gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: makep('x'), D: gowid.RenderWithWeight{W: 1}},
		&gowid.ContainerWidget{IWidget: makep('y'), D: gowid.RenderWithWeight{W: 1}},
		&gowid.ContainerWidget{IWidget: makep('z'), D: gowid.RenderWithWeight{W: 1}},
	}
	w := New(subs)
	c := w.Render(gowid.RenderBox{C: 1, R: 12}, gowid.Focused, gwtest.D)
	assert.Equal(t, `
x
x
x
x
y
y
y
y
z
z
z
z`[1:], c.String())
	subs[2] = &gowid.ContainerWidget{IWidget: makep('z'), D: renderWeightUpTo{gowid.RenderWithWeight{W: 1}, 2}}
	w = New(subs)
	c = w.Render(gowid.RenderBox{C: 1, R: 12}, gowid.Focused, gwtest.D)
	assert.Equal(t, `
x
x
x
x
x
y
y
y
y
y
z
z`[1:], c.String())

}

func TestPile6(t *testing.T) {
	subs := []gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: text.New("foo"), D: gowid.RenderFixed{}},
		&gowid.ContainerWidget{IWidget: text.New("bar"), D: gowid.RenderWithWeight{W: 1}},
		&gowid.ContainerWidget{IWidget: text.New("baz"), D: gowid.RenderFixed{}},
	}
	w := New(subs)
	c := w.Render(gowid.RenderBox{C: 3, R: 5}, gowid.Focused, gwtest.D)
	assert.Equal(t, `
foo
bar
   
   
baz`[1:], c.String())
}

//======================================================================
// Local Variables:
// mode: Go
// fill-column: 110
// End:
