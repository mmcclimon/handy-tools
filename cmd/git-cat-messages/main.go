package main

import (
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

var temp *os.File

func main() {
	// simple for now, more later
	gitArgs := []string{"log", "main..", "--format=%s%n%n%b%n%n", "--reverse"}
	diffArgs := []string{"diff", "main.."}

	var err error
	temp, err = os.CreateTemp("", "catmsg-")
	handleErr(err)

	git := exec.Command("git", gitArgs...)

	out, err := git.Output()
	handleErr(err)

	// deliberately ignoring errors here; it's fine.
	temp.Write(out)

	// and print the diff.
	git = exec.Command("git", diffArgs...)
	out, err = git.Output()
	handleErr(err)

	temp.WriteString("# ------------------------ >8 ------------------------\n")
	temp.Write(out)
	temp.WriteString("\n\n# vim: ft=gitcommit\n")

	// we're gonna print it, and then exec vim for it
	fmt.Println(temp.Name())

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	vim, err := exec.LookPath(editor)
	handleErr(err)

	err = unix.Exec(vim, []string{editor, temp.Name()}, os.Environ())
	handleErr(err)
}

func handleErr(err error) {
	if err == nil {
		return
	}

	// clean up tempfile on error
	if temp != nil {
		defer os.Remove(temp.Name())
	}

	fmt.Fprintln(os.Stdout, err)
	os.Exit(1)
}
