package main

import (
	"flag"
	"strings"

	"github.com/lordvidex/cpgen"
)

var (
	pq                  bool
	unionFind           bool
	count               uint
	only                string
	sieveOfEratosthenes bool
	cf                  bool
	fileIO              bool
)

func init() {
	flag.BoolVar(&pq, "pq", false, "Add code for priority queues")
	flag.BoolVar(&cf, "cf", false, "main should contain t testcases")
	flag.BoolVar(&unionFind, "uf", false, "Add code for union find")
	flag.BoolVar(&sieveOfEratosthenes, "sieve", false, "Add code for sieve of eratosthenes")
	flag.UintVar(&count, "c", 0, "number of solution files")
	flag.StringVar(
		&only,
		"only",
		"",
		"generate solution files for files (filenames should be comma separeted) e.g. 'a,b,c'",
	)
	flag.BoolVar(&fileIO, "file", false, "read input from 'input.txt' and read output to 'output.txt'")
}

func main() {
	flag.Parse()
	if only == "" && count < 1 {
		count = 5
	}
	config := cpgen.Config{
		Pq: pq,
		Uf: unionFind,
		Sv: sieveOfEratosthenes,
		Cf: cf,
	}
	if fileIO {
		config.FileIO = &cpgen.IO{
			Input:  "input.txt",
			Output: "output.txt",
		}
	}
	files := make([]string, 0)
	if count > 0 {
		for i := uint(0); i < count; i++ {
			questionNumber := string(rune(i + 'a'))
			files = append(files, questionNumber)
		}
	}
	if only != "" {
		files = append(files, strings.Split(only, ",")...)
	}
	c := cpgen.Generate(files, config, "")
	// wait for the task to finish
	for range c {
	}
}
