package main

import (
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/unix"
)

// This program exists so that when I type `gi tpull` it does the right thing.
func main() {
	git, err := exec.LookPath("git")
	if err != nil {
		panic(err)
	}

	os.Args[0] = git
	os.Args[1] = strings.TrimPrefix(os.Args[1], "t")

	if err := unix.Exec(git, os.Args, os.Environ()); err != nil {
		panic(err)
	}
}
