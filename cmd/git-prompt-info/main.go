package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	var isDirty, isWeird int

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

	// Check for weird cases (in the middle of rebase or whatever).
	weirdIndicators := []string{
		"rebase-apply",
		"rebase-merge",
		"MERGE_HEAD",
		"CHERRY_PICK_HEAD",
		"REVERT_HEAD",
	}

	cmd = exec.Command("git", "rev-parse", "--git-dir")
	out, err = cmd.Output()

	if err != nil {
		fmt.Println("0")
		return
	}

	gitDir := strings.TrimSpace(string(out))

	for _, f := range weirdIndicators {
		_, err := os.Stat(filepath.Join(gitDir, f))
		if err == nil {
			isWeird = 1
			break
		}
	}

	prep := "on"
	branch := head

	if head == "(detached)" {
		prep = "at"
		branch = sha[0:8]
	}

	fmt.Printf("1 %s %s %d %d\n", prep, branch, isDirty, isWeird)
}
