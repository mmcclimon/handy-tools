package main

import (
	"fmt"
	"regexp"

	"github.com/mmcclimon/handy-tools/internal/utils"
)

var blank = regexp.MustCompile(`\A\s*\z`)
var headerStart = regexp.MustCompile(`\A[-A-Za-z0-9]+:\s`)

func main() {
	ch := utils.DiamondOperator()

	inHeaders := true

	for line := range ch {
		if blank.MatchString(line) {
			// blank line, all done
			inHeaders = false
		}

		if !inHeaders || headerStart.MatchString(line) {
			fmt.Println(line)
		} else {
			fmt.Printf("    %s\n", line)
		}
	}
}
