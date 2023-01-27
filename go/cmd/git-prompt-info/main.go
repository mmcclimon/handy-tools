package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("git", "--no-optional-locks", "status", "--branch", "--porcelain=v2")
	out, err := cmd.Output()

	if err != nil {
		fmt.Println("0")
		return
	}

	var sha, head string
	var isDirty int

	scanner := bufio.NewScanner(bytes.NewReader(out))

line:
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "# branch.oid"):
			sha = strings.Split(line, " ")[2]
		case strings.HasPrefix(line, "# branch.head"):
			head = strings.Split(line, " ")[2]
		case !strings.HasPrefix(line, "# "):
			isDirty = 1
			break line
		}
	}

	prep := "on"
	branch := head

	if head == "(detached)" {
		prep = "at"
		branch = sha[0:8]
	}

	fmt.Printf("1 %s %s %d\n", prep, branch, isDirty)
}
