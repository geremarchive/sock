package funcs

import (
	"github.com/gdamore/tcell"
	"strings"
	"fmt"
)

func Addstr(s tcell.Screen, style tcell.Style, x int, y int, text string) {
	for i := x; i < len(text)+x; i++ {
		s.SetContent(i, y, rune(text[i-x]), []rune(""), style)
	}
}

func DrawScreen(s tcell.Screen, escape bool, center bool, bg string, color string, message string) {
	width, height := s.Size()

	if bg != "" {
		if escape {
			fmt.Print("\033]11;" + strings.ToUpper(HexHash(bg)) + "\007")
		} else {
			s.Fill(' ', tcell.StyleDefault.Background(tcell.GetColor(HexHash(bg))))
		}
	}

	if center {
		if escape {
			Addstr(s, tcell.StyleDefault.Foreground(tcell.GetColor(HexHash(color))),
				 (width / 2) - len(message) / 2, (height / 2) - 1,
				 message)
		} else {
			Addstr(s, tcell.StyleDefault.Foreground(tcell.GetColor(HexHash(color))).Background(tcell.GetColor(HexHash(bg))),
				 (width / 2) - len(message) / 2, (height / 2) - 1,
				 message)
		}
	} else {
		if escape {
			Addstr(s, tcell.StyleDefault.Foreground(tcell.GetColor(HexHash(color))), 0, 0, message)
		} else {
			Addstr(s, tcell.StyleDefault.Foreground(tcell.GetColor(HexHash(color))).Background(tcell.GetColor(HexHash(bg))), 0, 0, message)
		}
	}

	s.Show()
}

func HexHash(hex string) string {
	if len(hex) > 0 {
		if hex[0] == '#' {
			return hex
		} else {
			return "#" + hex
		}
	} else {
		return hex
	}
}
