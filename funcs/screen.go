package funcs

import (
	"github.com/gdamore/tcell"
	"strings"
	"fmt"
	"os"
)

func Addstr(s tcell.Screen, style tcell.Style, x int, y int, text string) {
	for i := x; i < len(text)+x; i++ {
		s.SetContent(i, y, rune(text[i-x]), []rune(""), style)
	}
}

func (o Options) DrawScreen(s tcell.Screen) {
	width, height := s.Size()

	if o.Bg != "" {
		if o.Escape {
			fmt.Print("\033]11;" + strings.ToUpper(HexHash(o.Bg)) + "\007")
		} else {
			s.Fill(' ', tcell.StyleDefault.Background(tcell.GetColor(HexHash(o.Bg))))
		}
	}

	if o.Center {
		if o.Escape {
			Addstr(s, tcell.StyleDefault.Foreground(tcell.GetColor(HexHash(o.Color))).Bold(o.Bold),
				 (width / 2) - len(o.Message) / 2, (height / 2) - 1,
				 o.Message)
		} else {
			Addstr(s, tcell.StyleDefault.Foreground(tcell.GetColor(HexHash(o.Color))).Background(tcell.GetColor(HexHash(o.Bg))).Bold(o.Bold),
				 (width / 2) - len(o.Message) / 2, (height / 2) - 1,
				 o.Message)
		}
	} else {
		if o.Escape {
			Addstr(s, tcell.StyleDefault.Foreground(tcell.GetColor(HexHash(o.Color))).Bold(o.Bold), 0, 0, o.Message)
		} else {
			Addstr(s, tcell.StyleDefault.Foreground(tcell.GetColor(HexHash(o.Color))).Background(tcell.GetColor(HexHash(o.Bg))).Bold(o.Bold), 0, 0, o.Message)
		}
	}

	s.Show()
}

func (o Options) Start(crypt string) error {
	s, err := tcell.NewScreen()

	if err != nil {
		return err
	}

	s.Init()
	width, height := s.Size()

	defer s.Fini()

	o.DrawScreen(s)

	if o.All || o.Check {
		go func() {
			for {
				if _, err := os.Stat("/tmp/locked.sock"); err != nil {
					break
				}
			}

			s.Fini()

			if o.Escape {
				fmt.Print("\033]11;" + strings.ToUpper(HexHash(o.Og)) + "\007")
			}

			os.Exit(0)
		}()
	}

	var chars string

	for {
		nw, nh := s.Size()

		if nw != width || nh != height {
			width, height = nw, nh

			s.Clear()
			o.DrawScreen(s)
		}

		if !(o.Check) {
			input := s.PollEvent()
			switch input := input.(type) {
				case *tcell.EventKey:
					if input.Key() == tcell.KeyEnter {
						matches, _ := MatchCrypt(chars, crypt)

						if matches {
							s.Fini()

							if o.Escape {
								fmt.Print("\033]11;" + strings.ToUpper(HexHash(o.Og)) + "\007")
							}

							if _, err := os.Stat("/tmp/locked.sock"); err == nil {
								if err := os.Remove("/tmp/locked.sock"); err != nil {
									fmt.Println("sock: couldn't delete /tmp/locked.sock")
								}
							}

							os.Exit(0)
						} else {
							chars = ""
						}
					} else {
						chars += string(input.Rune())
					}
			}
		}
	}

	return nil
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
