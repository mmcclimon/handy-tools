package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: timediff START END")
		os.Exit(1)
	}

	start := parse(os.Args[1])
	end := parse(os.Args[2])

	fmt.Println(end.Sub(start).Truncate(time.Second))
}

func parse(what string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, what)
	if err != nil {
		fmt.Printf("could not parse %q as time: %s\n", what, err)
		os.Exit(1)
	}

	return t
}
