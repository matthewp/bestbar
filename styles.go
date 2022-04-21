package bestbar

import (
	"github.com/gdamore/tcell/v2"
)

var Styles = struct {
	BackgroundColor            tcell.Color
	TextColor                  tcell.Color
	TextHotKeyColor            tcell.Color
	MenuBackgroundColorActived tcell.Color
}{
	tcell.ColorLightGray,
	tcell.ColorBlack.TrueColor(),
	tcell.NewRGBColor(190, 0, 0),
	tcell.NewRGBColor(87, 192, 56),
}

var Colors = struct {
	Charcoal tcell.Color
}{
	tcell.NewRGBColor(40, 35, 29),
}
