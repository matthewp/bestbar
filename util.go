package bestbar

import (
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

// stringWidth returns the number of horizontal cells needed to print the given
// text. It splits the text into its grapheme clusters, calculates each
// cluster's width, and adds them up to a total.
func stringWidth(text string) (width int) {
	g := uniseg.NewGraphemes(text)
	for g.Next() {
		var chWidth int
		for _, r := range g.Runes() {
			chWidth = runewidth.RuneWidth(r)
			if chWidth > 0 {
				break // Our best guess at this point is to use the width of the first non-zero-width rune.
			}
		}
		width += chWidth
	}
	return
}

func formatLabelWithShortcut(label string, shortcut rune) string {
	if shortcut != 0 {
		// Create label
		foundshortcut := false
		var sb strings.Builder
		for _, c := range label {
			//"[#be0000::b]%c[-:-:-]%s"
			if !foundshortcut && c == shortcut {
				sb.WriteString("[#be0000::b]")
				sb.WriteRune(c)
				sb.WriteString("[-:-:-]")
			} else {
				sb.WriteRune(c)
			}
		}
		l := sb.String()
		return l
	}
	return label
}
