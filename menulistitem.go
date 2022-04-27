package bestbar

type MenuListItem struct {
	OnSelected  func()
	label       string
	index       int
	setItemText func(int, string)
}

func NewMenuListItem(label string, index int) *MenuListItem {
	return &MenuListItem{
		label: label,
		index: index,
	}
}

func (mli *MenuListItem) SetSetTextFunc(setText func(int, string)) {
	mli.setItemText = setText
}

func (mli *MenuListItem) SetLabel(label string, shortcut rune) {
	fmt := formatLabelWithShortcut(label, shortcut)
	mli.label = fmt
	if mli.setItemText != nil {
		mli.setItemText(mli.index, fmt)
	}
}

func (mli *MenuListItem) Label() string {
	return mli.label
}
