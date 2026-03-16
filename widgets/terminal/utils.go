// Copyright 2019-2022 Graham Clark. All rights reserved.  Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Based heavily on vterm.py from urwid

package terminal

import (
	"fmt"

	"github.com/gdamore/tcell/v3"
	"github.com/gnuos/gowid"
	log "github.com/sirupsen/logrus"
)

//======================================================================

type EventNotSupported struct {
	Event any
}

var _ error = EventNotSupported{}

func (e EventNotSupported) Error() string {
	return fmt.Sprintf("Terminal input event %v of type %T not supported yet", e.Event, e.Event)
}

// TCellEventToBytes converts TCell's representation of a terminal event to
// the string of bytes that would be the equivalent event according to the
// supplied Terminfo object. It returns a tuple of the byte slice
// representing the terminal event (if successful), and a bool (denoting
// success or failure). This function is used by the TerminalWidget. Its
// subprocess is connected to a tty controlled by gowid. Events from the
// user are parsed by gowid via TCell - they are then translated by this
// function before being written to the TerminalWidget subprocess's tty.
func TCellEventToBytes(ev any, mouse IMouseSupport, last gowid.MouseState, paster IPaste) ([]byte, bool) {
	res := make([]byte, 0)
	res2 := false

	switch ev := ev.(type) {
	case *tcell.EventPaste:
		res2 = true
		if paster.PasteState() {
			// Already saw start
			res = append(res, []byte("\x1b[200~")...)
			paster.PasteState(false)
		} else {
			res = append(res, []byte("\x1b[200~")...)
			paster.PasteState(true)
		}
	case *tcell.EventKey:
		if ev.Key() < ' ' {
			str := []rune{rune(ev.Key())}
			res = append(res, string(str)...)
			res2 = true
		} else {
			res2 = true
			switch ev.Key() {
			case tcell.KeyRune:
				str := []rune(ev.Str())
				res = append(res, string(str)...)
			case tcell.KeyCR:
				str := []rune{rune(tcell.KeyCR)}
				res = append(res, string(str)...)
			case tcell.KeyF1:
				str := []rune{rune(tcell.KeyF1)}
				res = append(res, string(str)...)
			case tcell.KeyF2:
				str := []rune{rune(tcell.KeyF2)}
				res = append(res, string(str)...)
			case tcell.KeyF3:
				str := []rune{rune(tcell.KeyF3)}
				res = append(res, string(str)...)
			case tcell.KeyF4:
				str := []rune{rune(tcell.KeyF4)}
				res = append(res, string(str)...)
			case tcell.KeyF5:
				str := []rune{rune(tcell.KeyF5)}
				res = append(res, string(str)...)
			case tcell.KeyF6:
				str := []rune{rune(tcell.KeyF6)}
				res = append(res, string(str)...)
			case tcell.KeyF7:
				str := []rune{rune(tcell.KeyF7)}
				res = append(res, string(str)...)
			case tcell.KeyF8:
				str := []rune{rune(tcell.KeyF8)}
				res = append(res, string(str)...)
			case tcell.KeyF9:
				str := []rune{rune(tcell.KeyF9)}
				res = append(res, string(str)...)
			case tcell.KeyF10:
				str := []rune{rune(tcell.KeyF10)}
				res = append(res, string(str)...)
			case tcell.KeyF11:
				str := []rune{rune(tcell.KeyF11)}
				res = append(res, string(str)...)
			case tcell.KeyF12:
				str := []rune{rune(tcell.KeyF12)}
				res = append(res, string(str)...)
			case tcell.KeyF13:
				str := []rune{rune(tcell.KeyF13)}
				res = append(res, string(str)...)
			case tcell.KeyF14:
				str := []rune{rune(tcell.KeyF14)}
				res = append(res, string(str)...)
			case tcell.KeyF15:
				str := []rune{rune(tcell.KeyF15)}
				res = append(res, string(str)...)
			case tcell.KeyF16:
				str := []rune{rune(tcell.KeyF16)}
				res = append(res, string(str)...)
			case tcell.KeyF17:
				str := []rune{rune(tcell.KeyF17)}
				res = append(res, string(str)...)
			case tcell.KeyF18:
				str := []rune{rune(tcell.KeyF18)}
				res = append(res, string(str)...)
			case tcell.KeyF19:
				str := []rune{rune(tcell.KeyF19)}
				res = append(res, string(str)...)
			case tcell.KeyF20:
				str := []rune{rune(tcell.KeyF20)}
				res = append(res, string(str)...)
			case tcell.KeyF21:
				str := []rune{rune(tcell.KeyF21)}
				res = append(res, string(str)...)
			case tcell.KeyF22:
				str := []rune{rune(tcell.KeyF22)}
				res = append(res, string(str)...)
			case tcell.KeyF23:
				str := []rune{rune(tcell.KeyF23)}
				res = append(res, string(str)...)
			case tcell.KeyF24:
				str := []rune{rune(tcell.KeyF24)}
				res = append(res, string(str)...)
			case tcell.KeyF25:
				str := []rune{rune(tcell.KeyF25)}
				res = append(res, string(str)...)
			case tcell.KeyF26:
				str := []rune{rune(tcell.KeyF26)}
				res = append(res, string(str)...)
			case tcell.KeyF27:
				str := []rune{rune(tcell.KeyF27)}
				res = append(res, string(str)...)
			case tcell.KeyF28:
				str := []rune{rune(tcell.KeyF28)}
				res = append(res, string(str)...)
			case tcell.KeyF29:
				str := []rune{rune(tcell.KeyF29)}
				res = append(res, string(str)...)
			case tcell.KeyF30:
				str := []rune{rune(tcell.KeyF30)}
				res = append(res, string(str)...)
			case tcell.KeyF31:
				str := []rune{rune(tcell.KeyF31)}
				res = append(res, string(str)...)
			case tcell.KeyF32:
				str := []rune{rune(tcell.KeyF32)}
				res = append(res, string(str)...)
			case tcell.KeyF33:
				str := []rune{rune(tcell.KeyF33)}
				res = append(res, string(str)...)
			case tcell.KeyF34:
				str := []rune{rune(tcell.KeyF34)}
				res = append(res, string(str)...)
			case tcell.KeyF35:
				str := []rune{rune(tcell.KeyF35)}
				res = append(res, string(str)...)
			case tcell.KeyF36:
				str := []rune{rune(tcell.KeyF36)}
				res = append(res, string(str)...)
			case tcell.KeyF37:
				str := []rune{rune(tcell.KeyF37)}
				res = append(res, string(str)...)
			case tcell.KeyF38:
				str := []rune{rune(tcell.KeyF38)}
				res = append(res, string(str)...)
			case tcell.KeyF39:
				str := []rune{rune(tcell.KeyF39)}
				res = append(res, string(str)...)
			case tcell.KeyF40:
				str := []rune{rune(tcell.KeyF40)}
				res = append(res, string(str)...)
			case tcell.KeyF41:
				str := []rune{rune(tcell.KeyF41)}
				res = append(res, string(str)...)
			case tcell.KeyF42:
				str := []rune{rune(tcell.KeyF42)}
				res = append(res, string(str)...)
			case tcell.KeyF43:
				str := []rune{rune(tcell.KeyF43)}
				res = append(res, string(str)...)
			case tcell.KeyF44:
				str := []rune{rune(tcell.KeyF44)}
				res = append(res, string(str)...)
			case tcell.KeyF45:
				str := []rune{rune(tcell.KeyF45)}
				res = append(res, string(str)...)
			case tcell.KeyF46:
				str := []rune{rune(tcell.KeyF46)}
				res = append(res, string(str)...)
			case tcell.KeyF47:
				str := []rune{rune(tcell.KeyF47)}
				res = append(res, string(str)...)
			case tcell.KeyF48:
				str := []rune{rune(tcell.KeyF48)}
				res = append(res, string(str)...)
			case tcell.KeyF49:
				str := []rune{rune(tcell.KeyF49)}
				res = append(res, string(str)...)
			case tcell.KeyF50:
				str := []rune{rune(tcell.KeyF50)}
				res = append(res, string(str)...)
			case tcell.KeyF51:
				str := []rune{rune(tcell.KeyF51)}
				res = append(res, string(str)...)
			case tcell.KeyF52:
				str := []rune{rune(tcell.KeyF52)}
				res = append(res, string(str)...)
			case tcell.KeyF53:
				str := []rune{rune(tcell.KeyF53)}
				res = append(res, string(str)...)
			case tcell.KeyF54:
				str := []rune{rune(tcell.KeyF54)}
				res = append(res, string(str)...)
			case tcell.KeyF55:
				str := []rune{rune(tcell.KeyF55)}
				res = append(res, string(str)...)
			case tcell.KeyF56:
				str := []rune{rune(tcell.KeyF56)}
				res = append(res, string(str)...)
			case tcell.KeyF57:
				str := []rune{rune(tcell.KeyF57)}
				res = append(res, string(str)...)
			case tcell.KeyF58:
				str := []rune{rune(tcell.KeyF58)}
				res = append(res, string(str)...)
			case tcell.KeyF59:
				str := []rune{rune(tcell.KeyF59)}
				res = append(res, string(str)...)
			case tcell.KeyF60:
				str := []rune{rune(tcell.KeyF60)}
				res = append(res, string(str)...)
			case tcell.KeyF61:
				str := []rune{rune(tcell.KeyF61)}
				res = append(res, string(str)...)
			case tcell.KeyF62:
				str := []rune{rune(tcell.KeyF62)}
				res = append(res, string(str)...)
			case tcell.KeyF63:
				str := []rune{rune(tcell.KeyF63)}
				res = append(res, string(str)...)
			case tcell.KeyF64:
				str := []rune{rune(tcell.KeyF64)}
				res = append(res, string(str)...)
			case tcell.KeyInsert:
				str := []rune{rune(tcell.KeyInsert)}
				res = append(res, string(str)...)
			case tcell.KeyDelete:
				str := []rune{rune(tcell.KeyDelete)}
				res = append(res, string(str)...)
			case tcell.KeyHome:
				str := []rune{rune(tcell.KeyHome)}
				res = append(res, string(str)...)
			case tcell.KeyEnd:
				str := []rune{rune(tcell.KeyEnd)}
				res = append(res, string(str)...)
			case tcell.KeyHelp:
				str := []rune{rune(tcell.KeyHelp)}
				res = append(res, string(str)...)
			case tcell.KeyPgUp:
				str := []rune{rune(tcell.KeyPgUp)}
				res = append(res, string(str)...)
			case tcell.KeyPgDn:
				str := []rune{rune(tcell.KeyPgDn)}
				res = append(res, string(str)...)
			case tcell.KeyUp:
				str := []rune{rune(tcell.KeyUp)}
				res = append(res, string(str)...)
			case tcell.KeyDown:
				str := []rune{rune(tcell.KeyDown)}
				res = append(res, string(str)...)
			case tcell.KeyLeft:
				str := []rune{rune(tcell.KeyLeft)}
				res = append(res, string(str)...)
			case tcell.KeyRight:
				str := []rune{rune(tcell.KeyRight)}
				res = append(res, string(str)...)
			case tcell.KeyBacktab:
				str := []rune{rune(tcell.KeyBacktab)}
				res = append(res, string(str)...)
			case tcell.KeyExit:
				str := []rune{rune(tcell.KeyExit)}
				res = append(res, string(str)...)
			case tcell.KeyClear:
				str := []rune{rune(tcell.KeyClear)}
				res = append(res, string(str)...)
			case tcell.KeyPrint:
				str := []rune{rune(tcell.KeyPrint)}
				res = append(res, string(str)...)
			case tcell.KeyCancel:
				str := []rune{rune(tcell.KeyCancel)}
				res = append(res, string(str)...)
			case tcell.KeyDEL:
				str := []rune{rune(tcell.KeyBackspace)}
				res = append(res, string(str)...)
			case tcell.KeyBackspace:
				str := []rune{rune(tcell.KeyBackspace)}
				res = append(res, string(str)...)
			default:
				res2 = false
				panic(EventNotSupported{Event: ev})
			}
		}
	case *tcell.EventMouse:
		if mouse.MouseEnabled() {
			var data string

			btnind := 0
			switch ev.Buttons() {
			case tcell.Button1:
				btnind = 0
			case tcell.Button2:
				btnind = 1
			case tcell.Button3:
				btnind = 2
			case tcell.WheelUp:
				btnind = 64
			case tcell.WheelDown:
				btnind = 65
			}

			lastind := 0
			if last.LeftIsClicked() {
				lastind = 0
			} else if last.MiddleIsClicked() {
				lastind = 1
			} else if last.RightIsClicked() {
				lastind = 2
			}

			switch ev.Buttons() {
			case tcell.Button1, tcell.Button2, tcell.Button3, tcell.WheelUp, tcell.WheelDown:
				mx, my := ev.Position()
				btn := btnind
				if (last.LeftIsClicked() && (ev.Buttons() == tcell.Button1)) ||
					(last.MiddleIsClicked() && (ev.Buttons() == tcell.Button2)) ||
					(last.RightIsClicked() && (ev.Buttons() == tcell.Button3)) {
					// assume the mouse pointer has been moved with button down, a "drag"
					btn += 32
				}
				if mouse.MouseIsSgr() {
					data = fmt.Sprintf("\033[<%d;%d;%dM", btn, mx+1, my+1)
				} else {
					data = fmt.Sprintf("\033[M%c%c%c", btn+32, mx+33, my+33)
				}
				res = append(res, data...)
				res2 = true
			case tcell.ButtonNone:
				// TODO - how to report no press?
				mx, my := ev.Position()

				if last.LeftIsClicked() || last.MiddleIsClicked() || last.RightIsClicked() {
					// 0 means left mouse button, m means released
					if mouse.MouseIsSgr() {
						data = fmt.Sprintf("\033[<%d;%d;%dm", lastind, mx+1, my+1)
					} else if mouse.MouseReportAny() {
						data = fmt.Sprintf("\033[M%c%c%c", 35, mx+33, my+33)
					}
				} else if mouse.MouseReportAny() {
					if mouse.MouseIsSgr() {
						// +32 for motion, +3 for no button
						data = fmt.Sprintf("\033[<35;%d;%dm", mx+1, my+1)
					} else {
						data = fmt.Sprintf("\033[M%c%c%c", 35+32, mx+33, my+33)
					}
				}
				res = append(res, data...)
				res2 = true
			}
		}
	default:
		log.WithField("event", ev).Info("Event not implemented")
	}
	return res, res2
}
