package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

var blank = regexp.MustCompile(`\A\s*\z`)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	prevWasEmpty := true

	for scanner.Scan() {
		line := scanner.Text()

		isEmpty := blank.MatchString(line)

		if isEmpty && prevWasEmpty {
			continue // no double-blanks
		}

		if !isEmpty && !prevWasEmpty {
			fmt.Println("") // ensure blank
		}

		fmt.Println(line)
		prevWasEmpty = isEmpty
	}
}
