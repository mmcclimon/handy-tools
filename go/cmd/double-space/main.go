package main

import (
	"fmt"
	"regexp"

	"github.com/mmcclimon/handy-tools/internal/utils"
)

var blank = regexp.MustCompile(`\A\s*\z`)

func main() {
	lines := utils.DiamondOperator()

	prevWasEmpty := true

	for line := range lines {
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
