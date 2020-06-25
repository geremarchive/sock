package funcs

import (
	flag "github.com/spf13/pflag"
	"fmt"
	"os"
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
--bold, -B: Use bold text
--all, -a: Lock all terminals
--check, -k: Checks if all terminals are locked, if they are, it locks the terminal`

type Options struct {
	Message string
	Color string
	Bg string
	Og string

	Escape bool
	Center bool
	Bold bool
	All bool
	Check bool
}

func Args() (Options) {
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
	)

	flag.StringVarP(&message, "message", "m", "", "Set the message")
	flag.StringVarP(&color, "color", "c", "", "Set the color of the message")
	flag.StringVarP(&bg, "bg", "b", "", "Set the color of the background")
	flag.StringVarP(&og, "og", "o", "", "Set the original color (only needed if you are using the '-e' option)")

	escape := flag.BoolP("escape", "e", false, "Set the background color using escape seqences (check if your terminal supports this)")
	center := flag.BoolP("center", "C", false, "Center the text")
	bold := flag.BoolP("bold", "B", false, "Use bold text")

	all := flag.BoolP("all", "a", false, "Lock all terminals")
	check := flag.BoolP("check", "k", false, "Checks if all terminals are locked")

	flag.Parse()

	return Options {message, color, bg, og, *escape, *center, *bold, *all, *check}
}

func (o Options) CheckLock() {
	if o.Check {
		if _, err := os.Stat("/tmp/locked.sock"); err != nil {
			os.Exit(0)
		}
	}
}

func (o Options) Lock() (pass string) {
	if os.Geteuid() == 0 {
		var err error

		pass, err = GetCrypt("root")

		if err != nil || pass == "" {
			fmt.Println("Couldn't get root's encrypted password")
			os.Exit(0)
		}
	} else {
		fmt.Println("Must be run as root!")
		os.Exit(0)
	}

	if o.All {
		f, err := os.Create("/tmp/locked.sock")

		if err != nil {
			fmt.Println("Couldn't lock all terminals!")
			f.Close()
		}

		f.Close()
	}

	return
}
