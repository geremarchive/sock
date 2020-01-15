package main

import (
	"github.com/gdamore/tcell"
	fu "sock/funcs"
	flag "github.com/spf13/pflag"
	"strings"
	"os"
	"fmt"
)

const help = `Usage: sock [OPTION]
A simple terminal locker

--help, -h: Display this information
--message=[STRING], -m: Set the message
--color=[HEX], -c: Set the color of the message
--bg=[HEX], -b: Set the color of the background
--og=[HEX], -o: Set the original color (only needed if you are using the '-e' option)
--escape, -e: Set the background color using escape seqences (check if your terminal supports this)
--center, -C: Center the text
--all, -a: Lock all terminals
--check, -k: Checks if all terminals are locked, if they are, it locks the terminal`

func main() {
	if len(os.Args) == 2 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			fmt.Println(help)
			os.Exit(0)
		}
	}

	var (
		message string
		color string

		bg string
		og string

		chars string
		crypt string
	)

	flag.StringVarP(&message, "message", "m", "", "Set the message")
	flag.StringVarP(&color, "color", "c", "", "Set the color of the message")
	flag.StringVarP(&bg, "bg", "b", "", "Set the color of the background")
	flag.StringVarP(&og, "og", "o", "", "Set the original color (only needed if you are using the '-e' option)")

	escape := flag.BoolP("escape", "e", false, "Set the background color using escape seqences (check if your terminal supports this)")
	center := flag.BoolP("center", "C", false, "Center the text")

	all := flag.BoolP("all", "a", false, "Lock all terminals")
	check := flag.BoolP("check", "k", false, "Checks if all terminals are locked")

	flag.Parse()
	if *check {
		if _, err := os.Stat("/tmp/locked.sock"); err != nil {
			os.Exit(0)
		}
	} else {
		if os.Geteuid() == 0 {
			var err error
			crypt, err = fu.GetCrypt("root")

			if err != nil || crypt == "" {
				fmt.Println("Couldn't get root's encrypted password")
				os.Exit(0)
			}
		} else {
			fmt.Println("Must be run as root!")
			os.Exit(0)
		}
	}

	if *all {
		f, err := os.Create("/tmp/locked.sock")

		if err != nil {
			fmt.Println("Couldn't lock all terminals!")
			f.Close()
		}

		f.Close()
	}

	s, err := tcell.NewScreen()

	if err != nil {
		panic(err)
	}

	s.Init()
	width, height := s.Size()

	defer s.Fini()

	fu.DrawScreen(s, *escape, *center, bg, color, message)

	if *all || *check {
		go func() {
			for {
				if _, err := os.Stat("/tmp/locked.sock"); err != nil {
					break
				}
			}

			s.Fini()

			if *escape {
				fmt.Print("\033]11;" + strings.ToUpper(fu.HexHash(og)) + "\007")
			}

			os.Exit(0)
		}()
	}


	for {
		nw, nh := s.Size()

		if nw != width || nh != height {
			width, height = nw, nh

			s.Clear()
			fu.DrawScreen(s, *escape, *center, bg, color, message)
		}

		if !(*check) {
			input := s.PollEvent()
			switch input := input.(type) {
				case *tcell.EventKey:
					if input.Key() == tcell.KeyEnter {
						matches, err := fu.MatchCrypt(chars, crypt)

						if err != nil {
							panic(err)
						}

						if matches {
							s.Fini()

							if *escape {
								fmt.Print("\033]11;" + strings.ToUpper(fu.HexHash(og)) + "\007")
							}

							if _, err := os.Stat("/tmp/locked.sock"); err == nil {
								if err := os.Remove("/tmp/locked.sock"); err != nil {
									fmt.Println("Couldn't delete /tmp/locked.sock")
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
}

func closeAll(s tcell.Screen, esc bool, og string) {
	for {
		if _, err := os.Stat("/tmp/locked.sock"); err != nil {
			break
		}
	}

	s.Fini()

	if esc {
		fmt.Print("\033]11;" + strings.ToUpper(og) + "\007")
	}

	os.Exit(0)
}
