package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/mmcclimon/handy-tools/internal/cli"
)

var now = time.Now()
var opts = options{}
var root = "/Users/michael/.Trash"

type options struct {
	quiet  bool
	really bool
}

func main() {
	flags := cli.NewFlagSet("clean-trash")
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "be quiet")
	flags.BoolVar(&opts.really, "really", false, "actully delete stuff")
	cli.ParseFlags(flags)

	err := filepath.WalkDir(root, processPath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func processPath(path string, dirent fs.DirEntry, err error) error {
	if path == root {
		return nil
	}

	info, _ := dirent.Info()
	daysOld := now.Sub(info.ModTime()).Seconds() / 86400

	var toReturn error
	kind := "file"

	if dirent.IsDir() {
		toReturn = fs.SkipDir // do not descend
		kind = "dir"
	}

	if daysOld < 30 {
		return toReturn
	}

	verb := "would delete"
	if opts.really {
		verb = "deleted"
	}

	toLog := fmt.Sprintf(
		"%s %s %s (%d days old)",
		verb,
		kind,
		filepath.Base(path),
		int(daysOld),
	)

	if opts.really {
		err := os.RemoveAll(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error removing %s: %s\n", path, err)
			return toReturn
		}
	}

	if !opts.quiet {
		fmt.Println(toLog)
	}

	return toReturn
}
