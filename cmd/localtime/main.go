package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	times := make([]time.Time, 0, len(os.Args))

	for _, timestr := range os.Args[1:] {
		epoch, err := strconv.ParseInt(timestr, 10, 64)
		if err != nil {
			fmt.Printf("bad time: %s\n", err)
			os.Exit(1)
		}

		times = append(times, time.Unix(epoch, 0))
	}

	if len(times) == 0 {
		times = append(times, time.Now())
	}

	for _, t := range times {
		fmt.Println(t.Format("Jan 2, 2006 15:04:05 MST"))
	}
}
