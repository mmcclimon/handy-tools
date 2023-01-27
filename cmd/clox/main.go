package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	format := "Monday, 15:04"

	us, _ := time.LoadLocation("America/New_York")
	au, _ := time.LoadLocation("Australia/Melbourne")

	fmt.Printf("🇺🇸 : %s <-- you are here\n", now.In(us).Format(format))
	fmt.Printf("🌐 : %s\n", now.In(time.UTC).Format(format))
	fmt.Printf("🇦🇺 : %s\n", now.In(au).Format(format))
}
