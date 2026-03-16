// Copyright 2019-2022 Graham Clark. All rights reserved.  Use of this source
// code is governed by the MIT license that can be found in the LICENSE
// file.

package overlay

import (
	"testing"

	"github.com/gnuos/gowid"
	"github.com/gnuos/gowid/gwtest"
	"github.com/gnuos/gowid/widgets/styled"
	"github.com/gnuos/gowid/widgets/text"
	"github.com/gdamore/tcell/v3"
	"github.com/stretchr/testify/assert"
)

func TestOverlay1(t *testing.T) {
	tw := text.New("top")
	bw := text.New("bottom")
	ov := New(tw, bw, gowid.VAlignTop{}, gowid.RenderFixed{}, gowid.HAlignLeft{}, gowid.RenderFixed{})
	c := ov.Render(gowid.RenderFlowWith{C: 6}, gowid.Focused, gwtest.D)
	assert.Equal(t, "toptom", c.String())

	bwStyled := styled.New(bw, gowid.MakeStyledAs(gowid.StyleBold))

	// When the widget is created this way, the style from the lower widget bleeds through
	ov = New(tw, bwStyled, gowid.VAlignTop{}, gowid.RenderFixed{}, gowid.HAlignLeft{}, gowid.RenderFixed{})
	c = ov.Render(gowid.RenderFlowWith{C: 6}, gowid.Focused, gwtest.D)
	assert.Equal(t, "toptom", c.String())
	assert.Equal(t, tcell.AttrBold, c.CellAt(0, 0).Style().OnOff&tcell.AttrBold)

	// When the widget is created this way, the style from the upper widget is set unilaterally
	ov = New(tw, bwStyled, gowid.VAlignTop{}, gowid.RenderFixed{}, gowid.HAlignLeft{}, gowid.RenderFixed{},
		Options{
			IgnoreLowerStyle: true,
		})
	c = ov.Render(gowid.RenderFlowWith{C: 6}, gowid.Focused, gwtest.D)
	assert.Equal(t, "toptom", c.String())
	assert.Equal(t, tcell.AttrMask(0), c.CellAt(0, 0).Style().OnOff&tcell.AttrBold)
}
