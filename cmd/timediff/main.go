package main

import (
	"fmt"
	"os"
	"time"
)

const day = 24 * time.Hour

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: timediff START END")
		os.Exit(1)
	}

	start := parse(os.Args[1])
	end := parse(os.Args[2])

	delta := end.Sub(start).Truncate(time.Second)
	var days string

	if delta > day {
		d := time.Duration(delta / day)
		capped := delta - d*day
		days = fmt.Sprintf(" (%dd + %s)", d, capped)
	}

	fmt.Printf("%s%s\n", delta, days)
}

func parse(what string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, what)
	if err != nil {
		fmt.Printf("could not parse %q as time: %s\n", what, err)
		os.Exit(1)
	}

	return t
}
