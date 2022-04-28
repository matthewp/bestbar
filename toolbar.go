package bestbar

import (
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Toolbar struct {
	*tview.Flex
	titlebar  *tview.TextView
	buttons   *tview.Flex
	ml        *MenuList
	lists     []*MenuList
	activebtn *tview.Button
	shortkeys map[rune]func()
	drawFn    func()
}

func NewToolbar() *Toolbar {
	flex := tview.NewFlex()

	buttons := tview.NewFlex()
	buttons.SetBackgroundColor(Styles.BackgroundColor)
	buttons.AddItem(tview.NewBox().SetBackgroundColor(Styles.BackgroundColor), 0, 1, false)

	titlebar := tview.NewTextView()
	titlebar.SetBackgroundColor(Styles.BackgroundColor)
	titlebar.SetTextColor(Styles.TextColor)
	titlebar.SetTextAlign(tview.AlignRight)
	titlebar.SetBorderPadding(0, 0, 0, 2)

	flex.AddItem(buttons, 0, 2, false)
	flex.AddItem(titlebar, 0, 1, false)

	t := &Toolbar{
		Flex:      flex,
		buttons:   buttons,
		ml:        nil,
		shortkeys: make(map[rune]func()),
		titlebar:  titlebar,
	}

	return t
}

func (t *Toolbar) SetDrawFunc(drawFn func()) {
	t.drawFn = drawFn
}

func (t *Toolbar) Redraw() {
	if t.drawFn != nil {
		t.drawFn()
	}
}

func (t *Toolbar) AddMenuList(label string, shortcut rune) *MenuList {
	activeColor := Styles.MenuBackgroundColorActived
	inactiveColor := Styles.BackgroundColor

	lf := formatLabelWithShortcut(label, shortcut)
	ml := NewMenuList(lf, len(t.lists))
	ml.SetBackgroundColor(inactiveColor)
	ml.SetTextColor(Styles.TextColor)
	ml.SetSelectedBackgroundColor(activeColor)
	ml.SetDrawFunc(func() {
		t.Redraw()
	})
	ml.SetBeforeSelectedFunc(func() {
		ml.SetInDrag(false)
		t.CloseMenuList()
	})
	t.lists = append(t.lists, ml)

	btn := ml.Button()
	btn.SetBackgroundColor(inactiveColor)
	btn.SetBackgroundColorActivated(activeColor)
	btn.SetLabelColorActivated(Styles.TextColor)
	btn.SetLabelColor(Styles.TextColor)
	btn.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		switch action {
		case tview.MouseLeftDown:
			if btn.InRect(event.Position()) {
				btn.SetBackgroundColorActivated(activeColor)
				btn.Focus(func(p tview.Primitive) {})
				t.ml = ml
				t.ml.SetInDrag(true)
				go t.Redraw()
			}
			break
		case tview.MouseLeftClick:
			if btn.InRect(event.Position()) {
				if t.activebtn == btn {
					btn.SetBackgroundColor(inactiveColor)
					btn.SetBackgroundColorActivated(inactiveColor)
					btn.Blur()
					t.ml.SetInDrag(false)
					t.ml = nil
					t.activebtn = nil
				} else {
					btn.SetBackgroundColor(activeColor)
					if t.activebtn != nil {
						t.activebtn.SetBackgroundColor(inactiveColor)
						t.activebtn.Blur()
					}
					t.ml = ml
					t.ml.SetInDrag(false)
					t.activebtn = btn
				}
			} else if t.activebtn == btn {
				btn.SetBackgroundColor(inactiveColor)
				btn.Blur()
				t.ml.SetInDrag(false)
				t.ml = nil
				t.activebtn = nil
			}

			break
		case tview.MouseMove:
			if btn.HasFocus() {
				if event.Buttons()&tcell.Button1 == 1 {
					if btn.InRect(event.Position()) {
						t.ml = ml
						t.ml.SetInDrag(true)
					} else {
						btn.Blur()
						t.ml.SetInDrag(false)
						t.ml = nil
						go t.Redraw()
					}
				}
			} else if tcell.Button1&event.Buttons() == 1 {
				if btn.InRect(event.Position()) {
					btn.Focus(func(p tview.Primitive) {})
					t.ml = ml
					t.ml.SetInDrag(true)
					go t.Redraw()
				}
			}
			break
		}

		return action, event
	})

	t.shortkeys[unicode.ToLower(rune(shortcut))] = func() {
		t.ToggleMenuList(ml)
	}

	idx := t.buttons.GetItemCount()
	filler := t.buttons.GetItem(idx - 1)
	t.buttons.RemoveItem(filler)
	w := tview.TaggedStringWidth(label) + 4
	t.buttons.AddItem(btn, w, 0, false)
	t.buttons.AddItem(filler, 0, 1, false)

	return ml
}

func (t *Toolbar) Draw(screen tcell.Screen) {
	screenWidth, screenHeight := screen.Size()

	width, _ := screenWidth, screenHeight

	t.SetRect(0, 0, width, 1)
	t.Flex.SetRect(0, 0, width, 1)

	if t.ml != nil {
		t.ml.Draw(screen)
	}

	t.Flex.Draw(screen)
}

func (t *Toolbar) ToggleMenuList(ml *MenuList) {
	if t.ml == ml {
		t.CloseMenuList()
	} else {
		t.MakeMenuListActive(ml)
	}
}

func (t *Toolbar) MakeMenuListActive(ml *MenuList) {
	btn := ml.Button()
	activeColor := Styles.MenuBackgroundColorActived
	inactiveColor := Styles.BackgroundColor

	btn.SetBackgroundColor(activeColor)
	if t.activebtn != nil {
		t.activebtn.SetBackgroundColor(inactiveColor)
		t.activebtn.Blur()
	}
	t.ml = ml
	t.activebtn = btn
}

func (t *Toolbar) CloseMenuList() {
	if t.ml != nil {
		btn := t.ml.Button()
		btn.Blur()
		btn.SetBackgroundColor(Styles.BackgroundColor)
		t.ml = nil
		go t.Redraw()
	}
}

func (t *Toolbar) goToMenuList(forward bool) {
	l := len(t.lists)
	if t.ml == nil || l == 0 {
		return
	}

	idx := t.ml.Index()

	if !forward {
		if idx > 0 {
			t.MakeMenuListActive(t.lists[idx-1])
		} else {
			t.MakeMenuListActive(t.lists[l-1])
		}
	} else {
		if idx < (l - 1) {
			t.MakeMenuListActive(t.lists[idx+1])
		} else {
			t.MakeMenuListActive(t.lists[0])
		}
	}
}

func (t *Toolbar) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return t.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if !consumed {
			switch action {
			case tview.MouseLeftClick:
				if t.HasOpenMenu() {
					if t.ml.InRect(event.Position()) {
						return t.ml.MouseHandler()(action, event, setFocus)
					} else {
						consumed, capture = t.Flex.MouseHandler()(action, event, setFocus)
						if !consumed {
							t.ml.Blur()
							t.ml = nil
						}
						return consumed, capture
					}
				} else {
					consumed, capture = t.Flex.MouseHandler()(action, event, setFocus)
					if consumed && t.ml != nil {
						setFocus(t.ml)
					}
					return consumed, capture
				}
			}

			if t.ml != nil {
				consumed, capture = t.ml.MouseHandler()(action, event, setFocus)
				if consumed {
					return consumed, capture
				}
			}
		}

		return t.Flex.MouseHandler()(action, event, setFocus)
	})
}

func (t *Toolbar) HasOpenMenu() bool {
	return t.ml != nil
}

func (t *Toolbar) InputCapture(event *tcell.EventKey) *tcell.EventKey {
	if t.HasOpenMenu() {
		switch event.Key() {
		case tcell.KeyEsc:
			t.CloseMenuList()
			return nil
		case tcell.KeyLeft:
			t.goToMenuList(false)
			return nil
		case tcell.KeyRight:
			t.goToMenuList(true)
			return nil
		}

		event = t.ml.InputCapture(event)

		if event != nil {
			// Run shortcut
			if val, ok := t.shortkeys[event.Rune()]; ok {
				val()
				return nil
			}
		}
		return event
	} else {
		// Run shortcut
		if val, ok := t.shortkeys[event.Rune()]; ok {
			val()
			return nil
		}
	}
	return event
}

func (t *Toolbar) SetTitle(title string) {
	t.titlebar.SetText(title)
}

func (t *Toolbar) SetTitleTextColor(color tcell.Color) {
	t.titlebar.SetTextColor(color)
}
