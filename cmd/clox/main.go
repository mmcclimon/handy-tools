package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	us, _ := time.LoadLocation("America/New_York")
	au, _ := time.LoadLocation("Australia/Melbourne")

	now := time.Now()

	if len(os.Args) > 1 {
		today := now.Format("2006-01-02")
		when, err := time.ParseInLocation("2006-01-02 15:04", today+" "+os.Args[1], us)

		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing time, using now (%s)\n", err)
		} else {
			now = when
		}
	}

	format := "Monday, 15:04"

	fmt.Printf("%s : %s <-- you are here\n", "🇺🇸", now.In(us).Format(format))
	fmt.Printf("%s : %s\n", "🇺🇳", now.In(time.UTC).Format(format))
	fmt.Printf("%s : %s\n", "🇦🇺", now.In(au).Format(format))
}
