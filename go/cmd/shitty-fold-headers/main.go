package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

var blank = regexp.MustCompile(`\A\s*\z`)
var headerStart = regexp.MustCompile(`\A[-A-Za-z0-9]+:\s`)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	inHeaders := true

	for scanner.Scan() {
		line := scanner.Text()

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
