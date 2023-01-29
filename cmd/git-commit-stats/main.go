package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/maps"
)

type Stats struct {
	files      map[string]PlusMinus
	totals     PlusMinus
	first      time.Time
	last       time.Time
	numCommits int
}

type PlusMinus struct {
	plus  int
	minus int
}

func main() {
	var suppressFiles bool

	args := []string{"log", "--numstat", "--no-merges", "--format=COMMIT %h %ct"}

	// silly opt parsing
	for _, arg := range os.Args[1:] {
		if arg == "--no-files" {
			suppressFiles = true
		} else {
			args = append(args, arg)
		}
	}

	stats := computeStats(args)
	outputStats(stats, suppressFiles)
}

func computeStats(gitArgs []string) Stats {
	handleErr := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	git := exec.Command("git", gitArgs...)

	stdout, err := git.StdoutPipe()
	handleErr(err)
	scanner := bufio.NewScanner(stdout)

	err = git.Start()
	handleErr(err)

	var currentCommit string
	var first, last int
	totals := PlusMinus{}
	files := make(map[string]PlusMinus)
	commits := make(map[string]int)

	headerRe := regexp.MustCompile(`^COMMIT (.*?) (.*)`)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// process header lines
		if m := headerRe.FindStringSubmatch(line); m != nil {
			currentCommit = m[1]
			epoch, _ := strconv.Atoi(m[2])

			if first == 0 || epoch < first {
				first = epoch
			}

			if last == 0 || epoch > last {
				last = epoch
			}

			continue
		}

		// process list of files
		bits := strings.SplitN(line, "\t", 3)
		plus := bits[0]
		minus := bits[1]
		path := bits[2]

		if plus == "-" || minus == "-" || (plus == "0" && minus == "0") {
			continue // not interesting
		}

		plusN, _ := strconv.Atoi(plus)
		minusN, _ := strconv.Atoi(minus)

		totals.plus += plusN
		totals.minus += minusN

		stats := files[path]
		stats.plus += plusN
		stats.minus += minusN
		files[path] = stats

		commits[currentCommit]++
	}

	if err := git.Wait(); err != nil {
		log.Fatal(err)
	}

	return Stats{
		files:      files,
		totals:     totals,
		numCommits: len(commits),
		first:      time.Unix(int64(first), 0),
		last:       time.Unix(int64(last), 0),
	}
}

func outputStats(stats Stats, suppressFiles bool) {
	if !suppressFiles {
		keys := maps.Keys(stats.files)
		sort.Strings(keys)

		for _, path := range keys {
			pm := stats.files[path]
			plus := fmt.Sprintf("+%d", pm.plus)
			minus := fmt.Sprintf("-%d", pm.minus)
			fmt.Printf("%5s / %5s  %s\n", plus, minus, path)
		}

		fmt.Println("")
	}

	totals := stats.totals
	average := float64(totals.plus+totals.minus) / float64(stats.numCommits)
	fmt.Printf("%d commits total(+%d/-%d)\n", stats.numCommits, totals.plus, totals.minus)
	fmt.Printf("avg of %.2f lines diff per commit\n", average)

	nDays := int(stats.last.Sub(stats.first).Seconds() / 86400)

	fmt.Printf(
		"avg of %.2f commits per day, from %s to %s (%s)\n",
		float64(stats.numCommits)/float64(nDays),
		stats.first.Format("2006-01-02"),
		stats.last.Format("2006-01-02"),
		formatDuration(stats.last.Sub(stats.first)),
	)
}

func formatDuration(dur time.Duration) string {
	days := dur.Seconds() / 86400
	dayNoun := "days"

	if int(days)%365 == 1 {
		dayNoun = "day"
	}

	if days < 365 {
		return fmt.Sprintf("%d %s", int(days), dayNoun)
	}

	years := int(days) / 365
	yearNoun := "year"

	if years > 1 {
		yearNoun = "years"
	}

	return fmt.Sprintf("%d %s and %d %s", years, yearNoun, int(days)%365, dayNoun)
}
