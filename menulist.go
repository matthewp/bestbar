package bestbar

import (
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	initialSelectionMade int = 1 << 0
	inFakeSelection      int = 1 << 1
	inDrag               int = 1 << 2
)

type MenuList struct {
	*tview.Box
	btn                     *tview.Button
	group                   *tview.Box
	list                    *tview.List
	shadow                  *tview.Box
	backgroundColor         tcell.Color
	selectedBackgroundColor tcell.Color
	states                  int
	index                   int
	shortkeys               map[rune]func()
	items                   map[int]*MenuListItem
	drawFn                  func()
	beforeSelectedFn        func()
}

func NewMenuList(label string, index int) *MenuList {
	l := tview.NewList()
	l.ShowSecondaryText(false)
	l.SetHighlightFullLine(true)

	group := tview.NewBox()
	group.SetBorder(true)
	group.SetBorderPadding(0, 0, 1, 1)
	group.SetBackgroundColor(tcell.ColorDarkOliveGreen)

	b := tview.NewBox()
	b.SetBackgroundColor(tcell.ColorDarkSlateGrey)

	shadow := tview.NewBox()
	shadow.SetBackgroundColor(Colors.Charcoal)

	return &MenuList{
		Box:                     b,
		list:                    l,
		index:                   index,
		group:                   group,
		btn:                     tview.NewButton(label),
		shadow:                  shadow,
		backgroundColor:         tcell.ColorDefault,
		selectedBackgroundColor: tcell.ColorDefault,
		shortkeys:               make(map[rune]func()),
		items:                   make(map[int]*MenuListItem),
		states:                  0,
	}
}

func (m *MenuList) SetBackgroundColor(color tcell.Color) {
	m.backgroundColor = color
	if m.selectedBackgroundColor == tcell.ColorDefault {
		m.list.SetSelectedBackgroundColor(color)
	}
	m.list.SetBackgroundColor(color)
	m.group.SetBackgroundColor(color)
	m.Box.SetBackgroundColor(color)
}

func (m *MenuList) SetSelectedBackgroundColor(color tcell.Color) {
	m.selectedBackgroundColor = color
}

func (m *MenuList) SetTextColor(color tcell.Color) {
	m.list.SetMainTextColor(color)
	m.group.SetBorderColor(color)
}

func (m *MenuList) ItemIsHighlighted() bool {
	return false
}

func (m *MenuList) SetDrawFunc(drawFn func()) {
	m.drawFn = drawFn
}

func (m *MenuList) SetBeforeSelectedFunc(beforeSelectedFn func()) {
	m.beforeSelectedFn = beforeSelectedFn
}

func (m *MenuList) SetButtonLabel(label string, shortcut rune) {
	fmt := formatLabelWithShortcut(label, shortcut)
	m.Button().SetLabel(fmt)
}

func (m *MenuList) Button() *tview.Button {
	return m.btn
}

func (m *MenuList) Index() int {
	return m.index
}

func (m *MenuList) AddItem(label string, shortcut rune, selected func()) *MenuList {
	idx := m.list.GetItemCount()

	mli := NewMenuListItem(label, idx)
	mli.OnSelected = func() {
		if m.beforeSelectedFn != nil {
			m.beforeSelectedFn()
		}
		if selected != nil {
			selected()
		}
	}
	mli.SetSetTextFunc(func(i int, s string) {
		m.list.SetItemText(i, s, "")
	})
	if shortcut != 0 {
		lowershortcut := unicode.ToLower(shortcut)

		if selected != nil {
			m.shortkeys[lowershortcut] = mli.OnSelected
		}
	}
	formattedLabel := formatLabelWithShortcut(label, shortcut)

	m.items[idx] = mli
	m.list.AddItem(formattedLabel, "", 0, func() {
		if m.states&inFakeSelection == 0 && mli.OnSelected != nil {
			mli.OnSelected()
		}
	})

	return m
}

func (m *MenuList) GetItem(index int) *MenuListItem {
	return m.items[index]
}

func (m *MenuList) SetInDrag(dragging bool) {
	if dragging {
		m.states |= inDrag
	} else {
		m.states = m.states &^ inDrag
	}
}

func (m *MenuList) HasFocus() bool {
	return m.Box.HasFocus() || m.group.HasFocus() || m.list.HasFocus()
}

func (m *MenuList) InputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch key := event.Key(); key {
	case tcell.KeyDown:
		cnt := m.list.GetItemCount()
		idx := m.list.GetCurrentItem()
		if m.states&initialSelectionMade == 0 {
			idx = -1
			m.states |= initialSelectionMade
		}
		var newIdx int
		if idx < (cnt - 1) {
			newIdx = idx + 1
		} else {
			newIdx = 0
		}

		if newIdx != idx {
			m.list.SetCurrentItem(newIdx)
			m.list.SetSelectedBackgroundColor(m.selectedBackgroundColor)
		}

		return nil
	case tcell.KeyUp:
		idx := m.list.GetCurrentItem()
		if m.states&initialSelectionMade == 0 {
			idx = 1
			m.states |= initialSelectionMade
		}
		var newIdx int
		if idx > 0 {
			newIdx = idx - 1
		} else {
			cnt := m.list.GetItemCount()
			newIdx = cnt - 1
		}

		if newIdx != idx {
			m.list.SetCurrentItem(newIdx)
			m.list.SetSelectedBackgroundColor(m.selectedBackgroundColor)
		}

		return nil
	case tcell.KeyEnter:
		idx := m.list.GetCurrentItem()
		if mli, ok := m.items[idx]; ok {
			mli.OnSelected()
			return nil
		}
		return nil
	}

	// Run shortcut
	if val, ok := m.shortkeys[event.Rune()]; ok {
		val()
		return nil
	}

	return event
}

func (m *MenuList) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return m.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		switch action {
		case tview.MouseMove:
			if m.list.InRect(event.Position()) {
				m.list.SetSelectedBackgroundColor(m.selectedBackgroundColor)

				// Send mousemoves as clicks so it will select the item
				m.states |= inFakeSelection
				consumed, capture = m.list.MouseHandler()(tview.MouseLeftClick, event, setFocus)
				m.states = m.states &^ inFakeSelection

				if consumed {
					m.states |= initialSelectionMade
				}

				return consumed, capture

			} else {
				m.list.SetSelectedBackgroundColor(m.backgroundColor)
				if m.drawFn != nil {
					go m.drawFn()
				}
			}
			break
		case tview.MouseLeftUp:
			if m.list.InRect(event.Position()) && m.states&inDrag != 0 {
				// Click it..
				return m.list.MouseHandler()(tview.MouseLeftClick, event, setFocus)
			}
			break
		}

		return m.list.MouseHandler()(action, event, setFocus)
	})
}

func (m *MenuList) Draw(screen tcell.Screen) {
	boxPadding := 1
	groupPadding := 1

	y := 1

	// The x is to the start of the button
	x, _, _, _ := m.btn.GetRect()

	// The height should be the number of items, basically
	itemCount := m.list.GetItemCount()
	itemHeight := 1
	listHeight := itemCount * itemHeight

	// The width is the width of the list
	_, _, listWidth, _ := m.list.GetInnerRect()
	for _, item := range m.items {
		w := stringWidth(item.Label())
		if w > listWidth {
			listWidth = w
		}
	}

	width := listWidth + (boxPadding * 2) + (groupPadding * 2)
	height := listHeight + (groupPadding * 2)

	m.shadow.SetRect(x+2, y+1, width, height)
	m.shadow.Draw(screen)

	m.SetRect(x, y, width, height)
	m.Box.DrawForSubclass(screen, m)

	m.group.SetRect(x+boxPadding, y, width-(boxPadding*2), height)
	m.group.Draw(screen)

	m.list.SetRect(x+boxPadding+groupPadding, y+groupPadding, width-((boxPadding*2)+(groupPadding*2)), height-(groupPadding*2))
	m.list.Draw(screen)
}
