package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const day = 24 * time.Hour

var formats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	time.DateTime,
	time.DateOnly,
	time.TimeOnly,
}

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
	if strings.ToLower(what) == "now" {
		return time.Now()
	}

	for _, format := range formats {
		t, err := time.Parse(format, what)
		if err == nil {
			return t
		}
	}

	fmt.Printf("could not parse %q as time\n", what)
	os.Exit(1)
	panic("unreachable")
}
