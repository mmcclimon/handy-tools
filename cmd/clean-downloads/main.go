package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mmcclimon/handy-tools/internal/cli"
	"golang.org/x/exp/slices"
)

var opts struct {
	desktop      bool
	really       bool
	leaveFolders bool
}

var downloadDir = "/Users/michael/Downloads"
var desktopDir = "/Users/michael/Desktop"
var bydateDir = filepath.Join(downloadDir, "bydate")

func main() {
	flags := cli.NewFlagSet("clean-downloads")
	flags.BoolVar(&opts.desktop, "desktop", false, "run on desktop, not downloads")
	flags.BoolVarP(&opts.leaveFolders, "leave-folders", "f", false, "do not move folders (default with --desktop)")
	flags.BoolVarP(&opts.really, "really", "r", false, "actually move stuff")

	cli.ParseFlags(flags)

	dir := downloadDir

	if opts.desktop {
		dir = desktopDir
	}

	cleanDir(dir)
	cleanBydate(dir)
}

func cleanDir(dir string) {
	leaveFolders := opts.leaveFolders || opts.desktop

	for _, path := range readDir(dir) {
		basename := path.Name()

		if basename == ".DS_Store" ||
			basename == "bydate" && path.IsDir() ||
			leaveFolders && path.IsDir() {
			continue
		}

		fullPath := filepath.Join(dir, basename)

		info, err := path.Info()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		daysOld := time.Since(info.ModTime()).Hours() / 24
		if daysOld < 2 {
			continue
		}

		ymd := info.ModTime().Local().Format(time.DateOnly)
		datedir := filepath.Join(bydateDir, ymd)

		if !opts.really {
			fmt.Printf("would move %s to %s\n", basename, prettyPath(datedir))
			continue
		}

		if err := os.MkdirAll(datedir, 0755); err != nil {
			log.Fatalf("error making dir: %s", err)
		}

		if err := os.Rename(fullPath, filepath.Join(datedir, basename)); err != nil {
			log.Fatalf("error moving file: %s", err)
		}
	}
}

func cleanBydate(dir string) {
	then := time.Now().Add(time.Duration(-1 * time.Hour * 24 * 30)).Format(time.DateOnly)

	for _, child := range readDir(bydateDir) {
		if isHiddenFile(child) {
			continue
		}

		abspath := filepath.Join(bydateDir, child.Name())
		maybeRemoveEmptyDir(child, abspath)

		bits := strings.Split(child.Name(), "-")

		if len(bits) <= 2 {
			continue // already a filed month
		}

		monthPath := filepath.Join(bydateDir, strings.Join(bits[0:2], "-"))

		if child.Name() < then {
			moveFiles(abspath, monthPath)
		}
	}
}

func maybeRemoveEmptyDir(dir fs.DirEntry, abspath string) {
	if !dir.IsDir() {
		return
	}

	files := readDir(abspath)

	// we remove if it's totally empty, or only has hidden files
	hasVisibleFiles := slices.ContainsFunc(files, func(de fs.DirEntry) bool {
		return !isHiddenFile(de)
	})

	if hasVisibleFiles {
		return
	}

	if !opts.really {
		fmt.Printf("would remove empty dir %s\n", prettyPath(abspath))
		return
	}

	if err := os.RemoveAll(abspath); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func moveFiles(dayPath, monthPath string) {
	for _, child := range readDir(dayPath) {
		abs := filepath.Join(dayPath, child.Name())

		if !opts.really {
			fmt.Printf("would move %s to %s\n", prettyPath(abs), prettyPath(monthPath))
			continue
		}

		if err := os.MkdirAll(monthPath, 0755); err != nil {
			log.Fatalf("error making dir: %s", err)
		}

		if err := os.Rename(abs, filepath.Join(monthPath, child.Name())); err != nil {
			log.Fatalf("error moving file: %s", err)
		}
	}

	if !opts.really {
		return
	}

	if err := os.Remove(dayPath); err != nil {
		log.Fatalf("error removing day dir: %s", err)
	}
}

// Utilities

// readDir is just os.ReadDir with fatal errors
func readDir(dir string) []fs.DirEntry {
	files, err := os.ReadDir(dir)

	if err != nil {
		log.Fatalf("error reading dir: %s", err)
	}

	return files
}

func isHiddenFile(dirent fs.DirEntry) bool {
	return dirent.Name()[0] == '.'
}

func prettyPath(path string) string {
	s := strings.TrimPrefix(path, "/Users/michael/Downloads/")
	return strings.TrimPrefix(s, "/Users/michael/Desktop/")
}
