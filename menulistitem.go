package bestbar

type MenuListItem struct {
	OnSelected func()
	label      string
	index      int
}

func NewMenuListItem(label string, index int) *MenuListItem {
	return &MenuListItem{
		label: label,
		index: index,
	}
}
func (mli *MenuListItem) SetLabel(label string) {
	mli.label = label
}

func (mli *MenuListItem) Label() string {
	return mli.label
}
