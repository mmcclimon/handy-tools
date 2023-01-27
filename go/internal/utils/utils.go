package utils

import (
	"bufio"
	"fmt"
	"os"
)

func DiamondOperator() <-chan string {
	files := os.Args[1:]

	if len(files) == 0 {
		files = append(files, "-")
	}

	out := make(chan string)

	go readFiles(files, out)

	return out
}

func readFiles(files []string, ch chan<- string) {
	for _, filename := range files {
		f := fileForName(filename)
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}

	close(ch)
}

func fileForName(filename string) *os.File {
	if filename == "-" {
		return os.Stdin
	}

	f, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return f
}
